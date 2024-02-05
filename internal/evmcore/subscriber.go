package evmcore

import (
	"github.com/artheranet/arthera-node/contracts/subscriber"
	"github.com/artheranet/arthera-node/internal/evmcore/vmcontext"
	"github.com/artheranet/arthera-node/internal/inter"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

// InfiniteCap = 2 ** 256 - 1
var InfiniteCap = new(big.Int).Sub(new(big.Int).Exp(new(big.Int).SetUint64(2), new(big.Int).SetUint64(256), nil), new(big.Int).SetUint64(1))

func ValidateSubscriberBalance(from common.Address, tx *types.Transaction, state *state.StateDB, vmRunner vmcontext.EVMRunner) bool {
	senderAABalance := state.GetBalance(from)
	return ValidateSubscriberBalanceWithParam(senderAABalance, from, tx, state, vmRunner)
}

func ValidateSubscriberBalanceWithParam(senderAABalance *big.Int, from common.Address, tx *types.Transaction, state *state.StateDB, vmRunner vmcontext.EVMRunner) bool {
	// if the tx has a value, the sender must have enough balance to pay for it
	if senderAABalance.Cmp(tx.Value()) < 0 {
		return false
	}

	senderSub := GetSubscriptionData(from, false, vmRunner)
	var receiverSub *subscriber.Subscription = nil
	contractCreation := tx.To() == nil
	if !contractCreation && state.GetCodeSize(*tx.To()) > 0 {
		receiverSub = GetSubscriptionData(*tx.To(), true, vmRunner)
	}

	// check to see if the target is a contract account, and it can pay for all gas units from its subscription
	if SubscriptionDataValid(receiverSub) {
		receiverActiveSub := HasActiveSubscription(*tx.To(), true, vmRunner)
		if receiverActiveSub {
			// the dapp has an active subscription
			senderWhitelisted := IsWhitelistedForContract(*tx.To(), from, vmRunner)
			if senderWhitelisted {
				// the sender is whitelisted, the receiver needs to cover the gas
				subBalance := GetCappedBalance(receiverSub, *tx.To(), true, vmRunner)
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
		senderActiveSub := HasActiveSubscription(from, false, vmRunner)
		if senderActiveSub {
			// the sender has an active subscription that needs to cover the gas
			subBalance := GetCappedBalance(senderSub, from, false, vmRunner)
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

func GetCappedBalance(subscr *subscriber.Subscription, address common.Address, contractSub bool, vmRunner vmcontext.EVMRunner) *big.Int {
	capRemaining := GetCapRemaining(address, contractSub, vmRunner)
	subBalance := subscr.Balance
	if capRemaining.Cmp(InfiniteCap) != 0 {
		if subBalance.Cmp(capRemaining) > 0 {
			// if the sub balance > cap remaining, use the cap remaining
			subBalance = capRemaining
		}
	}
	return subBalance
}

func HasActiveSubscription(address common.Address, contractSub bool, vmRunner vmcontext.EVMRunner) bool {
	activeSub, err := subscriber.HasActiveSubscription(vmRunner, address, contractSub)
	if err != nil {
		log.Debug("Subscribers::HasActiveSubscription() failed", "error", err.Error())
		activeSub = false
	}
	return activeSub
}

func GetSubscriptionData(address common.Address, contractSub bool, vmRunner vmcontext.EVMRunner) *subscriber.Subscription {
	sub, err := subscriber.GetSubscriptionData(vmRunner, address, contractSub)
	if err != nil {
		log.Error("Subscribers::getSubscription() failed", "error", err.Error())
		sub = nil
	}
	return sub
}

func IsWhitelistedForContract(_contract common.Address, _account common.Address, vmRunner vmcontext.EVMRunner) bool {
	whitelisted, err := subscriber.IsWhitelistedForContract(vmRunner, _contract, _account)
	if err != nil {
		log.Error("Subscribers::isWhitelisted() failed", "error", err.Error())
		whitelisted = false
	}
	return whitelisted
}

func GetCapWindow(address common.Address, contractSub bool, vmRunner vmcontext.EVMRunner) inter.Timestamp {
	ts, err := subscriber.GetCapWindow(vmRunner, address, contractSub)
	if err != nil {
		log.Error("Subscribers::getCapWindow() failed", "error", err.Error())
		ts = inter.FromUnix(0)
	}
	return ts
}

func GetCapRemaining(address common.Address, contractSub bool, vmRunner vmcontext.EVMRunner) *big.Int {
	capRemaining, err := subscriber.GetCapRemaining(vmRunner, address, contractSub)
	if err != nil {
		log.Error("Subscribers::getCapRemaining() failed", "error", err.Error())
		capRemaining = big.NewInt(0)
	}
	return capRemaining
}

func DebitSubscription(target common.Address, units *big.Int, contractSub bool, vmRunner vmcontext.EVMRunner) *big.Int {
	if units.BitLen() == 0 {
		return big.NewInt(0)
	}
	result, err := subscriber.DebitSubscription(vmRunner, target, units, contractSub)
	if err != nil {
		log.Error("Subscribers::debitSubscription() failed", "error", err.Error())
		return units
	}
	return result
}

func CreditSubscription(target common.Address, units *big.Int, contractSub bool, vmRunner vmcontext.EVMRunner) {
	if units.BitLen() == 0 {
		return
	}
	err := subscriber.CreditSubscription(vmRunner, target, units, contractSub)
	if err != nil {
		log.Error("Subscribers::creditSubscription() failed", "error", err.Error())
	}
}

func SubscriptionDataValid(sub *subscriber.Subscription) bool {
	return sub != nil &&
		sub.Id != nil && sub.Id.BitLen() > 0 &&
		sub.PlanId != nil && sub.PlanId.BitLen() > 0 &&
		sub.StartTime != nil && sub.StartTime.BitLen() > 0 &&
		sub.EndTime != nil && sub.EndTime.BitLen() > 0
}

func GetSubscriberById(subId *big.Int, vmRunner vmcontext.EVMRunner) common.Address {
	sub, err := subscriber.GetSubscriberById(vmRunner, subId)
	if err != nil {
		log.Error("Subscribers::getSubscriberById() failed", "error", err.Error())
		sub = common.Address{}
	}
	return sub
}
