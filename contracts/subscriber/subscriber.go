package subscriber

import (
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/artheranet/arthera-node/contracts/runner"
	"github.com/artheranet/arthera-node/internal/evmcore/vmcontext"
	"github.com/artheranet/arthera-node/internal/inter"
	"github.com/artheranet/arthera-node/params"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	hasActiveSubscription    = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "hasActiveSubscription", params.MaxGasForHasActiveSubscription)
	debitSubscription        = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "debit", params.MaxGasForDebitSubscription)
	creditSubscription       = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "credit", params.MaxGasForCreditSubscription)
	getSubscriptionData      = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "getSubscriptionData", params.MaxGasForGetSub)
	getCapRemaining          = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "getCapRemaining", params.MaxGasForGetSub)
	getCapWindow             = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "getCapWindow", params.MaxGasForGetSub)
	isWhitelistedForContract = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "isWhitelistedForContract", params.MaxGasForIsWhitelisted)
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

func HasActiveSubscription(evmRunner vmcontext.EVMRunner, subscriber common.Address, contractSub bool) (bool, error) {
	var result bool
	if subscriber == params.ZeroAddress {
		return false, nil
	}
	evmRunner.StopGasMetering()
	defer evmRunner.StartGasMetering()
	err := hasActiveSubscription.Query(evmRunner, &result, subscriber, contractSub)
	if err != nil {
		return false, err
	}
	return result, nil
}

func DebitSubscription(evmRunner vmcontext.EVMRunner, subscriber common.Address, units *big.Int, contractSub bool) (*big.Int, error) {
	var result *big.Int
	if subscriber == params.ZeroAddress {
		return units, nil
	}
	evmRunner.StopGasMetering()
	defer evmRunner.StartGasMetering()
	err := debitSubscription.Execute(evmRunner, &result, big.NewInt(0), subscriber, units, contractSub)
	if err != nil {
		return big.NewInt(0), err
	}

	return result, nil
}

func CreditSubscription(evmRunner vmcontext.EVMRunner, subscriber common.Address, units *big.Int, contractSub bool) error {
	if subscriber == params.ZeroAddress {
		return nil
	}
	evmRunner.StopGasMetering()
	defer evmRunner.StartGasMetering()
	err := creditSubscription.Execute(evmRunner, nil, big.NewInt(0), subscriber, units, contractSub)
	if err != nil {
		return err
	}

	return nil
}

func GetSubscriptionData(evmRunner vmcontext.EVMRunner, subscriber common.Address, contractSub bool) (*Subscription, error) {
	var result Subscription
	if subscriber == params.ZeroAddress {
		return nil, nil
	}
	evmRunner.StopGasMetering()
	defer evmRunner.StartGasMetering()
	err := getSubscriptionData.Query(evmRunner, &result, subscriber, contractSub)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetCapWindow(evmRunner vmcontext.EVMRunner, subscriber common.Address, contractSub bool) (inter.Timestamp, error) {
	var result *big.Int
	if subscriber == params.ZeroAddress {
		return inter.FromUnix(0), nil
	}
	evmRunner.StopGasMetering()
	defer evmRunner.StartGasMetering()
	err := getCapWindow.Query(evmRunner, &result, subscriber, contractSub)
	if err != nil {
		return inter.FromUnix(0), err
	}

	return inter.FromUnix(result.Int64()), nil
}

func GetCapRemaining(evmRunner vmcontext.EVMRunner, subscriber common.Address, contractSub bool) (*big.Int, error) {
	var result *big.Int
	if subscriber == params.ZeroAddress {
		return big.NewInt(0), nil
	}
	evmRunner.StopGasMetering()
	defer evmRunner.StartGasMetering()
	err := getCapRemaining.Query(evmRunner, &result, subscriber, contractSub)
	if err != nil {
		return big.NewInt(0), err
	}

	return result, nil
}

func IsWhitelistedForContract(evmRunner vmcontext.EVMRunner, contract common.Address, account common.Address) (bool, error) {
	var result bool
	if contract == params.ZeroAddress || account == params.ZeroAddress {
		return false, nil
	}
	evmRunner.StopGasMetering()
	defer evmRunner.StartGasMetering()
	err := isWhitelistedForContract.Query(evmRunner, &result, contract, account)
	if err != nil {
		return false, err
	}
	return result, nil
}
