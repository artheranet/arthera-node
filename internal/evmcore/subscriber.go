package evmcore

import (
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/artheranet/arthera-node/contracts/native_sub"
	"github.com/artheranet/arthera-node/internal/evmcore/vmcontext"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
)

func ValidateSubscriberBalance(from common.Address, tx *types.Transaction, state *state.StateDB, vmRunner vmcontext.EVMRunner) bool {
	senderAABalance := state.GetBalance(from)
	return ValidateSubscriberBalanceWithParam(senderAABalance, from, tx, state, vmRunner)
}

func ValidateSubscriberBalanceWithParam(senderAABalance *big.Int, from common.Address, tx *types.Transaction, state *state.StateDB, vmRunner vmcontext.EVMRunner) bool {
	// if the tx has a value, the sender must have enough balance to pay for it
	if senderAABalance.Cmp(tx.Value()) < 0 {
		return false
	}

	senderSub := native_sub.GetSubscriptionData(state, from)
	var receiverSub *native_sub.Subscription = nil
	contractCreation := tx.To() == nil
	if !contractCreation && state.GetCodeSize(*tx.To()) > 0 {
		receiverSub = native_sub.GetSubscriptionData(state, *tx.To())
	}

	// check to see if the target is a contract account, and it can pay for all gas units from its subscription
	if native_sub.SubscriptionDataValid(receiverSub) {
		receiverActiveSub := native_sub.HasActiveSubscription(state, *tx.To())
		if receiverActiveSub {
			// the dapp has an active subscription
			senderWhitelisted := native_sub.IsWhitelisted(state, *tx.To(), from)
			if senderWhitelisted {
				// the sender is whitelisted, the receiver needs to cover the gas
				subBalance := GetCappedBalance(receiverSub, *tx.To(), state)
				if subBalance.Cmp(new(big.Int).SetUint64(tx.Gas())) < 0 {
					// the subscription balance is not enough to cover the gas
					// the user needs to have funds to cover the difference
					pyagUnits := new(big.Int).Sub(new(big.Int).SetUint64(tx.Gas()), subBalance)
					remainingCost := new(big.Int).Mul(tx.GasPrice(), pyagUnits)
					return senderAABalance.Cmp(remainingCost) >= 0
				} else {
					// the subscription balance is enough to cover the gas
					return true
				}
			} else {
				// the sender is not whitelisted, he pays the entire cost
				return senderAABalance.Cmp(tx.Cost()) >= 0
			}
		} else {
			// the receiver does not have an active subscription, the sender pays the entire cost
			return senderAABalance.Cmp(tx.Cost()) >= 0
		}
	} else {
		// the contract account does not have a subscription, check the sender's subscription
		senderActiveSub := native_sub.HasActiveSubscription(state, from)
		if senderActiveSub {
			// the sender has an active subscription that needs to cover the gas
			subBalance := GetCappedBalance(senderSub, from, state)
			if subBalance.Cmp(new(big.Int).SetUint64(tx.Gas())) < 0 {
				// the subscription balance is not enough to cover the gas
				// the user needs to have funds to cover the difference
				pyagUnits := new(big.Int).Sub(new(big.Int).SetUint64(tx.Gas()), subBalance)
				remainingCost := new(big.Int).Mul(tx.GasPrice(), pyagUnits)
				return senderAABalance.Cmp(remainingCost) >= 0
			} else {
				// the subscription balance is enough to cover the gas
				return true
			}
		} else {
			// the sender pays the entire cost
			return senderAABalance.Cmp(tx.Cost()) >= 0
		}
	}
}

func GetCappedBalance(subscr *native_sub.Subscription, address common.Address, state vm.StateDB) *big.Int {
	capRemaining := native_sub.GetCapRemaining(state, address)
	subBalance := subscr.Balance
	if capRemaining.Cmp(abis.MaxUint256) != 0 {
		if subBalance.Cmp(capRemaining) > 0 {
			// if the sub balance > cap remaining, use the cap remaining
			subBalance = capRemaining
		}
	}
	return subBalance
}
