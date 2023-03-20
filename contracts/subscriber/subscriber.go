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
)

func HasActiveSubscription(evmRunner runner.EVMRunner, subscriber common.Address) (bool, error) {
	var result bool
	err := hasActiveSubscription.Execute(evmRunner, &result, big.NewInt(0), subscriber)

	if err != nil {
		return false, err
	}

	return result, nil
}
