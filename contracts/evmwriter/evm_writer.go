package evmwriter

import (
	"bytes"
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

var (
	setBalanceMethodID []byte
	copyCodeMethodID   []byte
	swapCodeMethodID   []byte
	setStorageMethodID []byte
	incNonceMethodID   []byte
)

func init() {
	for name, constID := range map[string]*[]byte{
		"setBalance": &setBalanceMethodID,
		"copyCode":   &copyCodeMethodID,
		"swapCode":   &swapCodeMethodID,
		"setStorage": &setStorageMethodID,
		"incNonce":   &incNonceMethodID,
	} {
		method, exist := abis.EVMWriter.Methods[name]
		if !exist {
			panic("unknown EvmWriter method")
		}

		*constID = make([]byte, len(method.ID))
		copy(*constID, method.ID)
	}
}

type PreCompiledContract struct{}

func (_ PreCompiledContract) Run(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if caller != contracts.NodeDriverSmartContractAddress {
		return nil, 0, vm.ErrExecutionReverted
	}
	if len(input) < 4 {
		return nil, 0, vm.ErrExecutionReverted
	}
	if bytes.Equal(input[:4], setBalanceMethodID) {
		input = input[4:]
		// setBalance
		if suppliedGas < params.CallValueTransferGas {
			return nil, 0, vm.ErrOutOfGas
		}
		suppliedGas -= params.CallValueTransferGas
		if len(input) != 64 {
			return nil, 0, vm.ErrExecutionReverted
		}

		acc := common.BytesToAddress(input[12:32])
		input = input[32:]
		value := new(big.Int).SetBytes(input[:32])

		if acc == evm.TxContext.Origin {
			// Origin balance shouldn't decrease during his transaction
			return nil, 0, vm.ErrExecutionReverted
		}

		balance := evm.StateDB.GetBalance(acc)
		if balance.Cmp(value) >= 0 {
			diff := new(big.Int).Sub(balance, value)
			evm.StateDB.SubBalance(acc, diff)
		} else {
			diff := new(big.Int).Sub(value, balance)
			evm.StateDB.AddBalance(acc, diff)
		}
	} else if bytes.Equal(input[:4], copyCodeMethodID) {
		input = input[4:]
		// copyCode
		if suppliedGas < params.CreateGas {
			return nil, 0, vm.ErrOutOfGas
		}
		suppliedGas -= params.CreateGas
		if len(input) != 64 {
			return nil, 0, vm.ErrExecutionReverted
		}

		accTo := common.BytesToAddress(input[12:32])
		input = input[32:]
		accFrom := common.BytesToAddress(input[12:32])

		code := evm.StateDB.GetCode(accFrom)
		if code == nil {
			code = []byte{}
		}
		cost := uint64(len(code)) * (params.CreateDataGas + params.MemoryGas)
		if suppliedGas < cost {
			return nil, 0, vm.ErrOutOfGas
		}
		suppliedGas -= cost
		if accTo != accFrom {
			evm.StateDB.SetCode(accTo, code)
		}
	} else if bytes.Equal(input[:4], swapCodeMethodID) {
		input = input[4:]
		// swapCode
		cost := 2 * params.CreateGas
		if suppliedGas < cost {
			return nil, 0, vm.ErrOutOfGas
		}
		suppliedGas -= cost
		if len(input) != 64 {
			return nil, 0, vm.ErrExecutionReverted
		}

		acc0 := common.BytesToAddress(input[12:32])
		input = input[32:]
		acc1 := common.BytesToAddress(input[12:32])
		code0 := evm.StateDB.GetCode(acc0)
		if code0 == nil {
			code0 = []byte{}
		}
		code1 := evm.StateDB.GetCode(acc1)
		if code1 == nil {
			code1 = []byte{}
		}
		cost0 := uint64(len(code0)) * (params.CreateDataGas + params.MemoryGas)
		cost1 := uint64(len(code1)) * (params.CreateDataGas + params.MemoryGas)
		cost = (cost0 + cost1) / 2 // 50% discount because trie size won't increase after pruning
		if suppliedGas < cost {
			return nil, 0, vm.ErrOutOfGas
		}
		suppliedGas -= cost
		if acc0 != acc1 {
			evm.StateDB.SetCode(acc0, code1)
			evm.StateDB.SetCode(acc1, code0)
		}
	} else if bytes.Equal(input[:4], setStorageMethodID) {
		input = input[4:]
		// setStorage
		if suppliedGas < params.SstoreSetGasEIP2200 {
			return nil, 0, vm.ErrOutOfGas
		}
		suppliedGas -= params.SstoreSetGasEIP2200
		if len(input) != 96 {
			return nil, 0, vm.ErrExecutionReverted
		}
		acc := common.BytesToAddress(input[12:32])
		input = input[32:]
		key := common.BytesToHash(input[:32])
		input = input[32:]
		value := common.BytesToHash(input[:32])

		evm.StateDB.SetState(acc, key, value)
	} else if bytes.Equal(input[:4], incNonceMethodID) {
		input = input[4:]
		// incNonce
		if suppliedGas < params.CallValueTransferGas {
			return nil, 0, vm.ErrOutOfGas
		}
		suppliedGas -= params.CallValueTransferGas
		if len(input) != 64 {
			return nil, 0, vm.ErrExecutionReverted
		}

		acc := common.BytesToAddress(input[12:32])
		input = input[32:]
		value := new(big.Int).SetBytes(input[:32])

		if acc == evm.TxContext.Origin {
			// Origin nonce shouldn't change during his transaction
			return nil, 0, vm.ErrExecutionReverted
		}

		if value.Cmp(common.Big256) >= 0 {
			// Don't allow large nonce increasing to prevent a nonce overflow
			return nil, 0, vm.ErrExecutionReverted
		}
		if value.Sign() <= 0 {
			return nil, 0, vm.ErrExecutionReverted
		}

		evm.StateDB.SetNonce(acc, evm.StateDB.GetNonce(acc)+value.Uint64())
	} else {
		return nil, 0, vm.ErrExecutionReverted
	}
	return nil, suppliedGas, nil
}
