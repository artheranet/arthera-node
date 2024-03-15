// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evmcore

import (
	"fmt"
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/pyag"
	"github.com/artheranet/arthera-node/contracts/subscriber"
	"github.com/artheranet/arthera-node/internal/evmcore/vmcontext"
	params2 "github.com/artheranet/arthera-node/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math"
	"math/big"
)

var emptyCodeHash = crypto.Keccak256Hash(nil)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==

	4a) Attempt to run transaction data
	4b) If valid, use result as code for the new state object

== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	gp               *GasPool
	msg              Message
	gas              uint64
	senderSpentGas   uint64
	receiverSpentGas uint64
	pyagSpentGas     uint64
	gasPrice         *big.Int
	initialGas       uint64
	value            *big.Int
	data             []byte
	state            vm.StateDB
	evm              *vm.EVM
	evmRunner        vmcontext.SharedEVMRunner
}

// Message represents a message sent to a contract.
type Message interface {
	From() common.Address
	To() *common.Address

	GasPrice() *big.Int
	GasFeeCap() *big.Int
	GasTipCap() *big.Int
	Gas() uint64
	Value() *big.Int

	Nonce() uint64
	IsFake() bool
	Data() []byte
	AccessList() types.AccessList
}

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	UsedGas    uint64 // Total used gas but include the refunded gas
	Err        error  // Any error encountered during the execution(listed in core/vm/errors.go)
	ReturnData []byte // Returned data from evm(function result or data supplied with revert opcode)
}

// Unwrap returns the internal evm error which allows us for further
// analysis outside.
func (result *ExecutionResult) Unwrap() error {
	return result.Err
}

// Failed returns the indicator whether the execution is successful or not
func (result *ExecutionResult) Failed() bool { return result.Err != nil }

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (result *ExecutionResult) Return() []byte {
	if result.Err != nil {
		return nil
	}
	return common.CopyBytes(result.ReturnData)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (result *ExecutionResult) Revert() []byte {
	if result.Err != vm.ErrExecutionReverted {
		return nil
	}
	return common.CopyBytes(result.ReturnData)
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
// it does not compute the actual gas consumed by executing the actual transaction
/*
intrinsic gas = minimum gas per tx type (TxGasContractCreation or TxGas)
					+ non-zero bytes txdata * TxDataNonZeroGasEIP2028
					+ zero-bytes txdata  * TxDataZeroGas
					+ len(access list) * TxAccessListAddressGas
					+ len(access list storage key) * TxAccessListStorageKeyGas
*/
func IntrinsicGas(data []byte, accessList types.AccessList, isContractCreation bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if isContractCreation {
		gas = params.TxGasContractCreation
	} else {
		gas = params.TxGas
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/params.TxDataNonZeroGasEIP2028 < nz {
			return 0, vm.ErrOutOfGas
		}
		gas += nz * params.TxDataNonZeroGasEIP2028

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, ErrGasUintOverflow
		}
		gas += z * params.TxDataZeroGas
	}
	if accessList != nil {
		gas += uint64(len(accessList)) * params.TxAccessListAddressGas
		gas += uint64(accessList.StorageKeys()) * params.TxAccessListStorageKeyGas
	}
	return gas, nil
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:        gp,
		evm:       evm,
		msg:       msg,
		gasPrice:  msg.GasPrice(),
		value:     msg.Value(),
		data:      msg.Data(),
		state:     evm.StateDB,
		evmRunner: vmcontext.SharedEVMRunner{EVM: evm},
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg Message, gp *GasPool) (*ExecutionResult, error) {
	res, err := NewStateTransition(evm, msg, gp).TransitionDb()
	if err != nil {
		log.Debug("Tx skipped", "err", err)
	}
	return res, err
}

// ApplyMessageWithoutGasPrice applies the given message with the gas price
// set to zero. It's only for use in eth_call and eth_estimateGas, so that they can be used
// with gas price set to zero if the sender doesn't have funds to pay for gas.
// Returns the gas used (which does not include gas refunds) and an error if it failed.
func ApplyMessageWithoutGasPrice(evm *vm.EVM, msg Message, gp *GasPool) (*ExecutionResult, error) {
	st := NewStateTransition(evm, msg, gp)
	st.gasPrice = big.NewInt(0)
	res, err := st.TransitionDb()
	if err != nil {
		log.Debug("Tx skipped", "err", err)
	}
	return res, err
}

// to returns the recipient of the message.
func (st *StateTransition) to() common.Address {
	if st.msg == nil || st.msg.To() == nil /* contract creation */ {
		return params2.ZeroAddress
	}
	return *st.msg.To()
}

func (st *StateTransition) buyGas(senderSub *subscriber.Subscription, receiverSub *subscriber.Subscription) error {
	pyagGasUnits := new(big.Int).SetUint64(st.msg.Gas())

	if st.gasPrice.BitLen() > 0 {
		log.Trace("Buying gas", "units", pyagGasUnits)
	}

	if st.gasPrice.BitLen() > 0 {
		// first check to see if the target is a contract account and it can pay for all gas units from its subscription
		if SubscriptionDataValid(receiverSub) {
			if st.hasActiveSubscription(receiverSub) {
				log.Trace("Receiver has an active subscription")
				// if the sender is whitelisted with the contract account
				if IsWhitelistedForContract(*st.msg.To(), st.msg.From(), &st.evmRunner) {
					// if the contract account has an active subscription, pay gas from its balance
					subBalance := GetCappedBalance(receiverSub, *st.msg.To(), true, &st.evmRunner)
					if subBalance.Cmp(pyagGasUnits) < 0 {
						// receiver's subscription balance is not enough
						// the overflowed value needs to be covered from the sender
						pyagGasUnits = pyagGasUnits.Sub(pyagGasUnits, subBalance)
						log.Trace("Receiver's subscription balance overflows", "balance", subBalance, "overflow", pyagGasUnits)
					} else {
						// the subscription has enough balance to cover the gas, nothing left to pay
						pyagGasUnits = big.NewInt(0)
						log.Trace("Receiver's subscription balance is enough", "balance", subBalance)
					}
				}
			}
		} else {
			// the contract account does not have a subscription, take fees from the user's subscription
			if st.hasActiveSubscription(senderSub) {
				log.Trace("Sender has an active subscription")
				subBalance := GetCappedBalance(senderSub, st.msg.From(), false, &st.evmRunner)
				if subBalance.Cmp(pyagGasUnits) < 0 {
					// the subscription balance is not enough
					// the overflowed value needs to be covered from Pay-as-You-Go
					pyagGasUnits = pyagGasUnits.Sub(pyagGasUnits, subBalance)
					log.Trace("Sender's subscription balance overflows", "balance", subBalance, "overflow", pyagGasUnits)
				} else {
					// the subscription has enough balance to cover the gas, nothing left to pay
					pyagGasUnits = big.NewInt(0)
					log.Trace("Sender's subscription balance is enough", "balance", subBalance)
				}
			}
		}
	}

	if st.gasPrice.BitLen() > 0 {
		log.Trace("Pay-as-You-Go required balance", "units", pyagGasUnits)
	}

	// at this point, gas units were deducted from existing subscriptions
	// check if there's anything else to pay under Pay-as-You-Go
	pyagGasValue := pyagGasUnits.Mul(pyagGasUnits, st.gasPrice)

	// and check if the sender has enough balance
	// Note: we don't need to check against gasFeeCap instead of gasPrice, as it's too aggressive in the asynchronous environment
	if have, want := st.state.GetBalance(st.msg.From()), pyagGasValue; have.Cmp(want) < 0 {
		return fmt.Errorf("%w: address %v have %v want %v", ErrInsufficientFunds, st.msg.From().Hex(), have, want)
	}

	// deduct the gas from the block gas counter
	if err := st.gp.SubGas(st.msg.Gas()); err != nil {
		return err
	}

	// set the initial gas specified in the message
	// everyone down the road will consume gas from st.gas
	st.gas += st.msg.Gas()

	// copy it to st.initialGas to keep the initial gas value specified in the message
	st.initialGas = st.msg.Gas()

	// debit the gas, redo the same logic as above
	pyagGasUnits = new(big.Int).SetUint64(st.msg.Gas())

	if st.gasPrice.BitLen() > 0 {
		if SubscriptionDataValid(receiverSub) && st.hasActiveSubscription(receiverSub) {
			if IsWhitelistedForContract(*st.msg.To(), st.msg.From(), &st.evmRunner) {
				// receiver pays from his subscription the entire cost is the caller is whitelisted
				// the caps are handled by the subscribers contract
				log.Trace("Debit from receiver's subscription", "units", pyagGasUnits)
				pyagGasUnits = DebitSubscription(*st.msg.To(), pyagGasUnits, true, &st.evmRunner)
				st.receiverSpentGas = st.msg.Gas() - pyagGasUnits.Uint64()
			}
		} else if st.hasActiveSubscription(senderSub) {
			// sender pays the rest from his subscription
			// the caps are handled by the subscribers contract
			log.Trace("Debit from sender's subscription", "units", pyagGasUnits)
			pyagGasUnits = DebitSubscription(st.msg.From(), pyagGasUnits, false, &st.evmRunner)
			st.senderSpentGas = st.msg.Gas() - pyagGasUnits.Uint64()
		}
	}

	// if there's anything else to pay not covered by subscriptions, do a standard (Pay-as-You-Go) payment
	if st.gasPrice.BitLen() > 0 {
		log.Trace("Debit from Pay-as-You-Go", "units", pyagGasUnits)
	}
	st.pyagSpentGas = pyagGasUnits.Uint64()

	pyagGasValue = pyagGasUnits.Mul(pyagGasUnits, st.gasPrice)

	st.state.SubBalance(st.msg.From(), pyagGasValue)

	return nil
}

func (st *StateTransition) preCheck(senderSub *subscriber.Subscription, receiverSub *subscriber.Subscription) error {
	// Only check transactions that are not fake
	if !st.msg.IsFake() {
		// Make sure this transaction's nonce is correct.
		stNonce := st.state.GetNonce(st.msg.From())
		if msgNonce := st.msg.Nonce(); stNonce < msgNonce {
			return fmt.Errorf("%w: address %v, tx: %d state: %d", ErrNonceTooHigh,
				st.msg.From().Hex(), msgNonce, stNonce)
		} else if stNonce > msgNonce {
			return fmt.Errorf("%w: address %v, tx: %d state: %d", ErrNonceTooLow,
				st.msg.From().Hex(), msgNonce, stNonce)
		} else if stNonce+1 < stNonce {
			return fmt.Errorf("%w: address %v, nonce: %d", ErrNonceMax,
				st.msg.From().Hex(), stNonce)
		}
		// Make sure the sender is an EOA
		if codeHash := st.state.GetCodeHash(st.msg.From()); codeHash != emptyCodeHash && codeHash != (common.Hash{}) {
			return fmt.Errorf("%w: address %v, codehash: %s", ErrSenderNoEOA,
				st.msg.From().Hex(), codeHash)
		}
	}
	// Note: we don't need to check gasFeeCap >= BaseFee, because it's already checked by epochcheck
	return st.buyGas(senderSub, receiverSub)
}

func (st *StateTransition) internal() bool {
	zeroAddr := common.Address{}
	return st.msg.From() == zeroAddr
}

// TransitionDb will transition the state by applying the current message and
// returning the evm execution result with following fields.
//
//   - used gas:
//     total gas used (including gas being refunded)
//   - returndata:
//     the returned data from evm
//   - concrete execution error:
//     various **EVM** error which aborts the execution,
//     e.g. ErrOutOfGas, ErrExecutionReverted
//
// However if any consensus issue encountered, return the error directly with
// nil evm execution result.
func (st *StateTransition) TransitionDb() (*ExecutionResult, error) {
	// First check this message satisfies all consensus rules before
	// applying the message. The rules include these clauses
	//
	// 1. the nonce of the message caller is correct
	// 2. caller has enough balance to cover transaction fee(gaslimit * gasprice)
	// 3. the amount of gas required is available in the block
	// 4. the purchased gas is enough to cover intrinsic usage
	// 5. there is no overflow when calculating intrinsic gas

	// Note: insufficient balance for **topmost** call isn't a consensus error in Arthera, unlike Ethereum
	// Such transaction will revert and consume sender's gas

	msg := st.msg
	sender := vm.AccountRef(msg.From())
	contractCreation := msg.To() == nil

	// check if the user has an active subscription
	senderSubscription1 := GetSubscriptionData(st.msg.From(), false, &st.evmRunner)
	var receiverSubscription *subscriber.Subscription = nil
	if !contractCreation && st.state.GetCodeSize(*st.msg.To()) > 0 {
		// receiver is a smart contract, retrieve its subscription
		receiverSubscription = GetSubscriptionData(*st.msg.To(), true, &st.evmRunner)
	}

	// Check clauses 1-3, buy gas if everything is correct
	if err := st.preCheck(senderSubscription1, receiverSubscription); err != nil {
		return nil, err
	}

	london := st.evm.ChainConfig().IsLondon(st.evm.Context.BlockNumber)

	// Check clauses 4-5, subtract intrinsic gas if everything is correct
	gas, err := IntrinsicGas(st.data, st.msg.AccessList(), contractCreation)
	if err != nil {
		return nil, err
	}
	if st.gas < gas {
		return nil, fmt.Errorf("%w: have %d, want %d", ErrIntrinsicGas, st.gas, gas)
	}

	// deduct the intrinsic gas
	st.gas -= gas

	// Set up the initial access list.
	if rules := st.evm.ChainConfig().Rules(st.evm.Context.BlockNumber); rules.IsBerlin {
		st.state.PrepareAccessList(msg.From(), msg.To(), vm.ActivePrecompiles(rules), msg.AccessList())
	}

	var senderHadActiveSubscription = st.hasActiveSubscription(senderSubscription1)

	var (
		ret   []byte
		vmerr error // vm errors do not effect consensus and are therefore not assigned to err
	)
	if contractCreation {
		// deduct Create gas
		ret, _, st.gas, vmerr = st.evm.Create(sender, st.data, st.gas, st.value)
	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(msg.From(), st.state.GetNonce(sender.Address())+1)
		// deduct Call gas
		ret, st.gas, vmerr = st.evm.Call(sender, st.to(), st.data, st.gas, st.value)
	}

	// 10% of unspent gas gets spent as a disincentive to militate against excessive transaction gas limits
	// The disincentive is required because Fantom is leaderless decentralized aBFT and blocks are not known in advance
	// to a validator (unlike Ethereum miner) until blocks are created from confirmed events.
	// There is no single proposer who originated transactions for a block and so validator doesn't know in advance
	// how much gas will be spent by a transaction, this is different from Ethereum and so adjustments were made
	// to address this issue. The penalty (10% charge of unspent gas) is introduced to avoid many of such cases.
	if !st.internal() {
		st.gas -= st.gas / 10
	}

	// check if the user has an active subscription again, because the transaction might be a terminate subscription
	senderSubscription2 := GetSubscriptionData(st.msg.From(), false, &st.evmRunner)
	receiverSubscription = nil
	if !contractCreation && st.state.GetCodeSize(*st.msg.To()) > 0 {
		// receiver is a smart contract, retrieve its subscription
		receiverSubscription = GetSubscriptionData(*st.msg.To(), true, &st.evmRunner)
	}

	if !london {
		// Before EIP-3529: refunds were capped to gasUsed / 2
		st.refundGas(params.RefundQuotient, senderHadActiveSubscription, senderSubscription1, senderSubscription2, receiverSubscription)
	} else {
		// After EIP-3529: refunds are capped to gasUsed / 5
		st.refundGas(params.RefundQuotientEIP3529, senderHadActiveSubscription, senderSubscription1, senderSubscription2, receiverSubscription)
	}

	// Pay-as-You-Go rebates
	if !contractCreation && !st.hasActiveSubscription(senderSubscription2) && !st.hasActiveSubscription(receiverSubscription) {
		if !contracts.IsSystemContract(st.to()) {
			// check to see if the destination address is eligible for Pay-as-You-Go rebates
			owner, _ := pyag.GetOwnerOfContract(&st.evmRunner, st.to())
			if owner != params2.ZeroAddress {
				deployerGas := st.gasUsed() / 10
				refund := new(big.Int).Mul(new(big.Int).SetUint64(deployerGas), st.gasPrice)
				err := pyag.AddReward(&st.evmRunner, owner, refund)
				if err == nil {
					st.state.AddBalance(contracts.PayAsYouGoGasRewardsContractAddress, refund)
				}
			}
		}
	}

	return &ExecutionResult{
		UsedGas:    st.gasUsed(),
		Err:        vmerr,
		ReturnData: ret,
	}, nil
}

// Returns the remaining gas, plus a refund to the sender, because the initial gas that was provided
// to the transaction might be bigger than was actually consumed
func (st *StateTransition) refundGas(
	refundQuotient uint64,
	senderHadActiveSubscription bool,
	prevSenderSubscription *subscriber.Subscription,
	senderSubscription *subscriber.Subscription,
	receiverSubscription *subscriber.Subscription,
) {
	// Apply refund counter, capped to a refund quotient
	refund := st.gasUsed() / refundQuotient
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gas += refund

	if st.gasPrice.BitLen() > 0 {
		// we have st.gas units to send back proportionally, exchanged at the original rate.

		if st.to() != params2.ZeroAddress {
			if st.hasActiveSubscription(receiverSubscription) {
				receiverGasRefund := st.gas * st.receiverSpentGas / st.initialGas
				if receiverGasRefund > 0 {
					log.Info("Credit receiver subscription", "refund (units)", receiverGasRefund)
					CreditSubscription(st.to(), new(big.Int).SetUint64(receiverGasRefund), true, &st.evmRunner)
				}
			} else {
				// if the receiver does hot have an active subscription, give the refund to PYAG
				st.pyagSpentGas += st.receiverSpentGas
			}
		}

		senderGasRefund := st.gas * st.senderSpentGas / st.initialGas
		if senderGasRefund > 0 {
			log.Info("Refunding gas", "sender", st.msg.From().String())

			if st.hasActiveSubscription(senderSubscription) {
				log.Info("Credit sender subscription", "refund (units)", senderGasRefund)
				CreditSubscription(st.msg.From(), new(big.Int).SetUint64(senderGasRefund), false, &st.evmRunner)
			} else {
				if st.isSubscribersCall() && senderHadActiveSubscription && !st.hasActiveSubscription(senderSubscription) {
					// the sender had an active subscription and no longer has it
					// we check to see if somebody else has it
					newSubscriber := GetSubscriberById(prevSenderSubscription.Id, &st.evmRunner)
					if newSubscriber != params2.ZeroAddress {
						// if yes, give the refund to the subscription
						CreditSubscription(newSubscriber, new(big.Int).SetUint64(senderGasRefund), false, &st.evmRunner)
					} else {
						// subscription was terminated, don't give a refund
					}
				} else {
					// if the sender does hot have an active subscription, give the refund to PYAG
					log.Info("Refunding Sender", "senderSpentGas", st.senderSpentGas)
					st.pyagSpentGas += st.senderSpentGas
				}
			}
		}

		pyagGasRefund := st.gas * st.pyagSpentGas / st.initialGas
		if pyagGasRefund > 0 {
			pyagRefund := new(big.Int).Mul(new(big.Int).SetUint64(pyagGasRefund), st.gasPrice)
			log.Info("Credit Pay-as-You-Go", "refund (units)", pyagGasRefund, "refund (wei)", pyagRefund.String())
			st.state.AddBalance(st.msg.From(), pyagRefund)
		}
	} else {
		remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
		st.state.AddBalance(st.msg.From(), remaining)
	}

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(st.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}

func (st *StateTransition) hasActiveSubscription(sub *subscriber.Subscription) bool {
	if sub == nil {
		return false
	}

	// if the gas price is zero, this is from eth_call or eth_estimateGas,
	// so subscriptions are not applied
	return st.msg.GasPrice().Cmp(big.NewInt(0)) > 0 &&
		sub != nil &&
		SubscriptionDataActive(sub, st.evm.Context.Time)
}

func (st *StateTransition) isSubscribersCall() bool {
	return st.to() == contracts.SubscribersSmartContractAddress
}
