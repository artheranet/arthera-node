package subscriber

import (
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/artheranet/arthera-node/contracts/runner"
	"github.com/artheranet/arthera-node/params"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	hasActiveSubscription = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "hasActiveSubscription", params.MaxGasForHasActiveSubscription)
	reduceBalance         = runner.NewBoundMethod(contracts.SubscribersSmartContractAddress, abis.Subscribers, "reduceBalance", params.MaxGasForReduceBalance)
)

func HasActiveSubscription(evmRunner runner.EVMRunner, subscriber common.Address) (bool, error) {
	var result bool
	if subscriber == contracts.ZeroAddress {
		return false, nil
	}
	err := hasActiveSubscription.Query(evmRunner, &result, big.NewInt(0), subscriber)
	if err != nil {
		return false, err
	}
	return result, nil
}

func ReduceBalance(evmRunner runner.EVMRunner, subscriber common.Address, units *big.Int) (*big.Int, error) {
	var result big.Int
	if subscriber == contracts.ZeroAddress {
		return units, nil
	}
	err := reduceBalance.Execute(evmRunner, &result, big.NewInt(0), subscriber, units)
	if err != nil {
		return big.NewInt(0), err
	}

	return &result, nil
}