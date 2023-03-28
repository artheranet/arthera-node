package subscriber

import (
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/artheranet/arthera-node/contracts/runner"
	"github.com/artheranet/arthera-node/inter"
	"github.com/artheranet/arthera-node/params"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	hasActiveSubscription = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "hasActiveSubscription", params.MaxGasForHasActiveSubscription)
	debitSubscription     = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "debit", params.MaxGasForDebitSubscription)
	creditSubscription    = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "credit", params.MaxGasForCreditSubscription)
	getSubscriptionData   = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "getSubscriptionData", params.MaxGasForGetSub)
	getCapRemaining       = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "getCapRemaining", params.MaxGasForGetSub)
	getCapWindow          = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "getCapWindow", params.MaxGasForGetSub)
	isWhitelisted         = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "isWhitelisted", params.MaxGasForIsWhitelisted)
)

type Subscription struct {
	Id           *big.Int
	PlanId       *big.Int
	Balance      *big.Int
	StartTime    *big.Int
	EndTime      *big.Int
	LastCapReset *big.Int
	PeriodUsage  *big.Int
}

func HasActiveSubscription(evmRunner runner.EVMRunner, subscriber common.Address) (bool, error) {
	var result bool
	if subscriber == contracts.ZeroAddress {
		return false, nil
	}
	err := hasActiveSubscription.Query(evmRunner, &result, subscriber)
	if err != nil {
		return false, err
	}
	return result, nil
}

func DebitSubscription(evmRunner runner.EVMRunner, subscriber common.Address, units *big.Int) (*big.Int, error) {
	var result *big.Int
	if subscriber == contracts.ZeroAddress {
		return units, nil
	}
	err := debitSubscription.Execute(evmRunner, &result, big.NewInt(0), subscriber, units)
	if err != nil {
		return big.NewInt(0), err
	}

	return result, nil
}

func CreditSubscription(evmRunner runner.EVMRunner, subscriber common.Address, units *big.Int) error {
	if subscriber == contracts.ZeroAddress {
		return nil
	}
	err := creditSubscription.Execute(evmRunner, nil, big.NewInt(0), subscriber, units)
	if err != nil {
		return err
	}

	return nil
}

func GetSubscriptionData(evmRunner runner.EVMRunner, subscriber common.Address) (*Subscription, error) {
	var result Subscription
	if subscriber == contracts.ZeroAddress {
		return nil, nil
	}

	err := getSubscriptionData.Query(evmRunner, &result, subscriber)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetCapWindow(evmRunner runner.EVMRunner, subscriber common.Address) (inter.Timestamp, error) {
	var result *big.Int
	if subscriber == contracts.ZeroAddress {
		return inter.FromUnix(0), nil
	}
	err := getCapWindow.Query(evmRunner, &result, subscriber)
	if err != nil {
		return inter.FromUnix(0), err
	}

	return inter.FromUnix(result.Int64()), nil
}

func GetCapRemaining(evmRunner runner.EVMRunner, subscriber common.Address) (*big.Int, error) {
	var result *big.Int
	if subscriber == contracts.ZeroAddress {
		return big.NewInt(0), nil
	}
	err := getCapRemaining.Query(evmRunner, &result, subscriber)
	if err != nil {
		return big.NewInt(0), err
	}

	return result, nil
}

func IsWhitelisted(evmRunner runner.EVMRunner, subscriber common.Address, account common.Address) (bool, error) {
	var result bool
	if subscriber == contracts.ZeroAddress || account == contracts.ZeroAddress {
		return false, nil
	}
	err := isWhitelisted.Query(evmRunner, &result, subscriber, account)
	if err != nil {
		return false, err
	}
	return result, nil
}
