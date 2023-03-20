package pyag

import (
	"github.com/artheranet/arthera-node/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func getStorageKey(addr common.Address) common.Hash {
	return crypto.Keccak256Hash(
		common.LeftPadBytes(addr.Bytes(), 32),
		common.LeftPadBytes(big.NewInt(0).Bytes(), 32),
	)
}

func GetDeployer(targetContract common.Address, statedb vm.StateDB) common.Address {
	storageKey := getStorageKey(targetContract)
	hashBytes := statedb.GetState(contracts.PayAsYouGoRebatesSmartContractAddress, storageKey)
	return common.BytesToAddress(hashBytes[(common.HashLength - common.AddressLength):])
}

func AddDeployer(targetContract common.Address, deployer common.Address, statedb *state.StateDB) {
	hash := getStorageKey(targetContract)
	statedb.SetState(
		contracts.PayAsYouGoRebatesSmartContractAddress,
		hash,
		common.BytesToHash(common.LeftPadBytes(deployer.Bytes()[:], 32)),
	)
}
