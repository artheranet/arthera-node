package pyag

import (
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/artheranet/arthera-node/contracts/runner"
	"github.com/artheranet/arthera-node/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var (
	getOwnerOfContract = runner.NewBoundMethod(contracts.PayAsYouGoGasRewardsContractAddress, abis.PayAsYouGoGasRewards, "getOwnerOfContract", params.MaxGasForSetOwnerOfContract)
	setOwnerOfContract = runner.NewBoundMethod(contracts.PayAsYouGoGasRewardsContractAddress, abis.PayAsYouGoGasRewards, "setOwnerOfContract", params.MaxGasForSetOwnerOfContract)
	addReward          = runner.NewBoundMethod(contracts.PayAsYouGoGasRewardsContractAddress, abis.PayAsYouGoGasRewards, "addReward", params.MaxGasForAddReward)
)

func getStorageKey(addr common.Address) common.Hash {
	return crypto.Keccak256Hash(
		common.LeftPadBytes(addr.Bytes(), 32),
		common.LeftPadBytes(big.NewInt(0).Bytes(), 32),
	)
}

func GetOwnerOfContractFast(targetContract common.Address, statedb vm.StateDB) common.Address {
	storageKey := getStorageKey(targetContract)
	hashBytes := statedb.GetState(contracts.PayAsYouGoGasRewardsContractAddress, storageKey)
	return common.BytesToAddress(hashBytes[(common.HashLength - common.AddressLength):])
}

func SetOwnerOfContractFast(targetContract common.Address, deployer common.Address, statedb *state.StateDB) {
	hash := getStorageKey(targetContract)
	statedb.SetState(
		contracts.PayAsYouGoGasRewardsContractAddress,
		hash,
		common.BytesToHash(common.LeftPadBytes(deployer.Bytes()[:], 32)),
	)
}

func AddReward(evmRunner runner.EVMRunner, contract common.Address, reward *big.Int) error {
	if contract == contracts.ZeroAddress {
		return nil
	}
	err := addReward.Execute(evmRunner, nil, big.NewInt(0), contract, reward)
	if err != nil {
		return err
	}

	return nil
}

func GetOwnerOfContract(evmRunner runner.EVMRunner, contract common.Address) (common.Address, error) {
	var result common.Address
	if contract == contracts.ZeroAddress {
		return contracts.ZeroAddress, nil
	}
	err := getOwnerOfContract.Query(evmRunner, &result, contract)
	if err != nil {
		return contracts.ZeroAddress, err
	}
	return result, nil
}

func SetOwnerOfContract(evmRunner runner.EVMRunner, contract common.Address, owner common.Address) error {
	if contract == contracts.ZeroAddress {
		return nil
	}
	err := setOwnerOfContract.Execute(evmRunner, nil, big.NewInt(0), contract, owner)
	if err != nil {
		return err
	}

	return nil
}