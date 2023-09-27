package abis

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var (
	ZeroAddress   = common.Address{}
	AbiUint8, _   = abi.NewType("uint8", "", nil)
	AbiUint256, _ = abi.NewType("uint256", "", nil)
	AbiString, _  = abi.NewType("string", "", nil)
	AbiAddress, _ = abi.NewType("address", "", nil)
	AbiBool, _    = abi.NewType("bool", "", nil)
	MaxUint256    = new(big.Int).Sub(new(big.Int).Exp(new(big.Int).SetUint64(2), new(big.Int).SetUint64(256), nil), new(big.Int).SetUint64(1))
)

func PackRevert(err string) []byte {
	errText, _ := abi.Arguments{{Type: AbiString}}.Pack(err)
	ret := []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	return append(ret, errText...)
}

func PackAbiBool(b bool) []byte {
	ret, _ := abi.Arguments{{Type: AbiBool}}.Pack(b)
	return ret
}

func PackAbiUint256(value *big.Int) []byte {
	ret, _ := abi.Arguments{{Type: AbiUint256}}.Pack(value)
	return ret
}

func PackAbiUint8(value uint8) []byte {
	ret, _ := abi.Arguments{{Type: AbiUint8}}.Pack(value)
	return ret
}

func PackAbiString(value string) []byte {
	ret, _ := abi.Arguments{{Type: AbiString}}.Pack(value)
	return ret
}

func PackAbiAddress(value common.Address) []byte {
	ret, _ := abi.Arguments{{Type: AbiAddress}}.Pack(value)
	return ret
}

func GetArrayData(state vm.StateDB, address common.Address, slot uint64, index uint64) common.Hash {
	arrayLen := state.GetState(address, common.BigToHash(new(big.Int).SetUint64(slot))).Big()
	if arrayLen.Cmp(new(big.Int).SetUint64(index)) <= 0 {
		return common.Hash{}
	}
	arrayStart := crypto.Keccak256Hash(common.LeftPadBytes(new(big.Int).SetUint64(slot).Bytes(), 32)).Big()
	arrayPos := new(big.Int).Add(arrayStart, new(big.Int).SetUint64(index))
	return state.GetState(address, crypto.Keccak256Hash(common.LeftPadBytes(arrayPos.Bytes(), 32)))
}

func SetArrayData(state vm.StateDB, address common.Address, slot uint64, index uint64, value common.Hash) {
	arrayLen := state.GetState(address, common.BigToHash(new(big.Int).SetUint64(slot))).Big()
	if arrayLen.Cmp(new(big.Int).SetUint64(index+1)) < 0 {
		// update array length
		arrayLen.SetUint64(index + 1)
		state.SetState(address, common.BigToHash(new(big.Int).SetUint64(slot)), common.BigToHash(arrayLen))
	}
	arrayStart := crypto.Keccak256Hash(common.LeftPadBytes(new(big.Int).SetUint64(slot).Bytes(), 32)).Big()
	arrayPos := new(big.Int).Add(arrayStart, new(big.Int).SetUint64(index))
	state.SetState(address, crypto.Keccak256Hash(common.LeftPadBytes(arrayPos.Bytes(), 32)), value)
}

func GetMappingData(state vm.StateDB, address common.Address, slot uint64, key []byte) common.Hash {
	storageKey := crypto.Keccak256Hash(
		common.LeftPadBytes(key, 32),
		common.LeftPadBytes(new(big.Int).SetUint64(slot).Bytes(), 32),
	)
	return state.GetState(address, storageKey)
}

func SetMappingData(state vm.StateDB, address common.Address, slot uint64, key []byte, value common.Hash) {
	storageKey := crypto.Keccak256Hash(
		common.LeftPadBytes(key, 32),
		common.LeftPadBytes(new(big.Int).SetUint64(slot).Bytes(), 32),
	)
	state.SetState(address, storageKey, value)
}

func GetDoubleMappingData(state vm.StateDB, address common.Address, slot uint64, key1 []byte, key2 []byte) common.Hash {
	storageKey := crypto.Keccak256Hash(
		common.LeftPadBytes(key2, 32),
		common.LeftPadBytes(key1, 32),
		common.LeftPadBytes(new(big.Int).SetUint64(slot).Bytes(), 32),
	)
	return state.GetState(address, storageKey)
}

func SetDoubleMappingData(state vm.StateDB, address common.Address, slot uint64, key1 []byte, key2 []byte, value common.Hash) {
	storageKey := crypto.Keccak256Hash(
		common.LeftPadBytes(key2, 32),
		common.LeftPadBytes(key1, 32),
		common.LeftPadBytes(new(big.Int).SetUint64(slot).Bytes(), 32),
	)
	state.SetState(address, storageKey, value)
}

func HasNumArgs(input []byte, num int) bool {
	return len(input) == 32*num
}

// GetUint256Arg pos starts from 0
func GetUint256Arg(input []byte, pos int) *big.Int {
	idx := pos * 32
	return new(big.Int).SetBytes(input[idx : idx+32])
}

// GetAddressArg pos starts from 0
func GetAddressArg(input []byte, pos int) common.Address {
	idx := pos * 32
	return common.BytesToAddress(input[idx+12 : idx+32])
}
