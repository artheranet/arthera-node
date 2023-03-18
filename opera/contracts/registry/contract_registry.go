package registry

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// GetContractBin is the Staking contract genesis implementation bin code
// Built from arthera-contracts main, solc 0.8.17, optimize-runs 10000
func GetContractBin() []byte {
	return hexutil.MustDecode("0x6080604052348015600f57600080fd5b506004361060285760003560e01c8063669d8d4514602d575b600080fd5b60636038366004608c565b73ffffffffffffffffffffffffffffffffffffffff9081166000908152602081905260409020541690565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b600060208284031215609d57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811460c057600080fd5b939250505056fea2646970667358221220cfb4d1cdd167178f1bfa399df7fe5f35d21cb6c58c165888db928191daeebb7264736f6c63430008110033")
}

// ContractAddress is the Registry contract address
var ContractAddress = common.HexToAddress("0xfc00face00000000000000000000000000000002")

func getStorageKey(addr common.Address) common.Hash {
	return crypto.Keccak256Hash(
		common.LeftPadBytes(addr.Bytes(), 32),
		common.LeftPadBytes(big.NewInt(0).Bytes(), 32),
	)
}

func GetDeployer(targetContract common.Address, statedb vm.StateDB) common.Address {
	storageKey := getStorageKey(targetContract)
	hashBytes := statedb.GetState(ContractAddress, storageKey)
	return common.BytesToAddress(hashBytes[(common.HashLength - common.AddressLength):])
}

func AddDeployer(targetContract common.Address, deployer common.Address, statedb *state.StateDB) {
	hash := getStorageKey(targetContract)
	statedb.SetState(
		ContractAddress,
		hash,
		common.BytesToHash(common.LeftPadBytes(deployer.Bytes()[:], 32)),
	)
}
