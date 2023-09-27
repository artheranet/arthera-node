package native_token

import (
	"bytes"
	"github.com/artheranet/arthera-node/contracts"
	"github.com/artheranet/arthera-node/contracts/abis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

const (
	AllowanceSlot = iota
)

var (
	transferEvent        = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	approvalEvent        = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
	totalSupplyMethodID  []byte
	balanceOfMethodID    []byte
	transferMethodID     []byte
	allowanceMethodID    []byte
	approveMethodID      []byte
	transferFromMethodID []byte
	nameMethodID         []byte
	symbolMethodID       []byte
	decimalsMethodID     []byte
)

func init() {
	for name, constID := range map[string]*[]byte{
		"totalSupply":  &totalSupplyMethodID,
		"balanceOf":    &balanceOfMethodID,
		"transfer":     &transferMethodID,
		"allowance":    &allowanceMethodID,
		"approve":      &approveMethodID,
		"transferFrom": &transferFromMethodID,
		"name":         &nameMethodID,
		"symbol":       &symbolMethodID,
		"decimals":     &decimalsMethodID,
	} {
		method, exist := abis.NativeToken.Methods[name]
		if !exist {
			panic("unknown IERC20Metadata method")
		}

		*constID = make([]byte, len(method.ID))
		copy(*constID, method.ID)
	}
}

type PreCompiledContract struct{}

func (_ PreCompiledContract) Run(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if len(input) < 4 {
		return nil, 0, vm.ErrExecutionReverted
	}
	if evm.StateDB.GetCodeSize(contracts.NativeTokenSmartContractAddress) == 0 {
		evm.StateDB.SetCode(contracts.NativeTokenSmartContractAddress, []byte{0})
	}
	methodId := input[:4]
	args := input[4:]
	if bytes.Equal(methodId, totalSupplyMethodID) {
		return totalSupply(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, balanceOfMethodID) {
		return balanceOf(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, transferMethodID) {
		return transfer(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, allowanceMethodID) {
		return allowance(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, approveMethodID) {
		return approve(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, transferFromMethodID) {
		return transferFrom(evm, caller, args, suppliedGas)
	} else if bytes.Equal(methodId, nameMethodID) {
		return abis.PackAbiString("Arthera"), suppliedGas, nil
	} else if bytes.Equal(methodId, symbolMethodID) {
		return abis.PackAbiString("AA"), suppliedGas, nil
	} else if bytes.Equal(methodId, decimalsMethodID) {
		return abis.PackAbiUint8(18), suppliedGas, nil
	} else {
		return nil, 0, vm.ErrExecutionReverted
	}
}

// function totalSupply() returns (uint256)
func totalSupply(evm *vm.EVM, caller common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	input, err := abis.Staking.Pack("totalSupply")
	if err != nil {
		return nil, 0, err
	}
	ret, gas, err := evm.StaticCall(vm.AccountRef(caller), contracts.StakingSmartContractAddress, input, suppliedGas)
	return ret, gas, err
}

// function balanceOf(address account) returns (uint256)
func balanceOf(evm *vm.EVM, _ common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 1) {
		return nil, 0, vm.ErrExecutionReverted
	}
	acc := abis.GetAddressArg(input, 0)
	balance := evm.StateDB.GetBalance(acc)
	return abis.PackAbiUint256(balance), suppliedGas, nil
}

// function transfer(address to, uint256 amount) returns (bool)
func transfer(evm *vm.EVM, owner common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if suppliedGas < params.CallValueTransferGas {
		return nil, 0, vm.ErrOutOfGas
	}
	suppliedGas -= params.CallValueTransferGas
	if !abis.HasNumArgs(input, 2) {
		return nil, 0, vm.ErrExecutionReverted
	}
	to := abis.GetAddressArg(input, 0)
	amount := abis.GetUint256Arg(input, 1)
	ret := _transfer(evm, owner, to, amount)
	if ret != nil {
		return ret, suppliedGas, vm.ErrExecutionReverted
	}
	return abis.PackAbiBool(true), suppliedGas, nil
}

// function approve(address spender, uint256 amount) returns (bool)
func approve(evm *vm.EVM, _ common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 2) {
		return nil, 0, vm.ErrExecutionReverted
	}
	if suppliedGas < params.SstoreSetGasEIP2200 {
		return nil, 0, vm.ErrOutOfGas
	}
	suppliedGas -= params.SstoreSetGasEIP2200
	spender := abis.GetAddressArg(input, 0)
	amount := abis.GetUint256Arg(input, 1)
	ret := _approve(evm, evm.TxContext.Origin, spender, amount)
	if ret != nil {
		return ret, suppliedGas, vm.ErrExecutionReverted
	}
	return abis.PackAbiBool(true), suppliedGas, nil
}

// function allowance(address owner, address spender) returns (uint256)
func allowance(evm *vm.EVM, _ common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 2) {
		return nil, 0, vm.ErrExecutionReverted
	}
	owner := abis.GetAddressArg(input, 0)
	spender := abis.GetAddressArg(input, 1)
	return abis.PackAbiUint256(_allowance(evm, owner, spender)), suppliedGas, nil
}

// function transferFrom(address from, address to, uint256 amount) returns (bool)
func transferFrom(evm *vm.EVM, spender common.Address, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
	if !abis.HasNumArgs(input, 3) {
		return nil, 0, vm.ErrExecutionReverted
	}
	from := abis.GetAddressArg(input, 0)
	to := abis.GetAddressArg(input, 1)
	amount := abis.GetUint256Arg(input, 2)
	ret := _spendAllowance(evm, from, spender, amount)
	if ret != nil {
		return ret, suppliedGas, vm.ErrExecutionReverted
	}
	_transfer(evm, from, to, amount)
	return abis.PackAbiBool(true), suppliedGas, nil
}

func _approve(evm *vm.EVM, owner common.Address, spender common.Address, amount *big.Int) []byte {
	if owner == abis.ZeroAddress {
		return abis.PackRevert("ERC20: approve from the zero address")
	}
	if spender == abis.ZeroAddress {
		return abis.PackRevert("ERC20: approve to the zero address")
	}
	abis.SetDoubleMappingData(
		evm.StateDB, contracts.NativeTokenSmartContractAddress, AllowanceSlot,
		owner.Bytes(), spender.Bytes(), common.BigToHash(amount),
	)
	event := createApprovalEvent(evm.TxContext.Origin, spender, amount, evm.Context.BlockNumber.Uint64())
	evm.StateDB.AddLog(&event)
	return nil
}

func _allowance(evm *vm.EVM, owner common.Address, spender common.Address) *big.Int {
	return abis.GetDoubleMappingData(
		evm.StateDB, contracts.NativeTokenSmartContractAddress, AllowanceSlot,
		owner.Bytes(), spender.Bytes(),
	).Big()
}

func _spendAllowance(evm *vm.EVM, owner common.Address, spender common.Address, amount *big.Int) []byte {
	currentAllowance := _allowance(evm, owner, spender)
	if currentAllowance.Cmp(abis.MaxUint256) != 0 {
		if currentAllowance.Cmp(amount) < 0 {
			return abis.PackRevert("ERC20: insufficient allowance")
		}
		_approve(evm, owner, spender, new(big.Int).Sub(currentAllowance, amount))
	}
	return nil
}

func _transfer(evm *vm.EVM, from common.Address, to common.Address, amount *big.Int) []byte {
	if from == abis.ZeroAddress {
		return abis.PackRevert("ERC20: transfer from the zero address")
	}
	if to == abis.ZeroAddress {
		return abis.PackRevert("ERC20: transfer to the zero address")
	}
	if amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}
	if evm.Context.CanTransfer(evm.StateDB, from, amount) {
		evm.Context.Transfer(evm.StateDB, from, to, amount)
		event := createTransferEvent(from, to, amount, evm.Context.BlockNumber.Uint64())
		evm.StateDB.AddLog(&event)
	} else {
		return abis.PackRevert("transfer amount exceeds balance")
	}
	return nil
}

func createTransferEvent(from common.Address, to common.Address, amount *big.Int, blockNumber uint64) types.Log {
	var topics = []common.Hash{
		transferEvent,
		common.BytesToHash(common.LeftPadBytes(from.Bytes(), 32)),
		common.BytesToHash(common.LeftPadBytes(to.Bytes(), 32)),
	}

	data := abis.PackAbiUint256(amount)

	return types.Log{
		Address:     contracts.NativeTokenSmartContractAddress,
		Topics:      topics,
		Data:        data,
		BlockNumber: blockNumber,
	}
}

func createApprovalEvent(owner common.Address, spender common.Address, amount *big.Int, blockNumber uint64) types.Log {
	var topics = []common.Hash{
		approvalEvent,
		common.BytesToHash(common.LeftPadBytes(owner.Bytes(), 32)),
		common.BytesToHash(common.LeftPadBytes(spender.Bytes(), 32)),
	}

	data := abis.PackAbiUint256(amount)

	return types.Log{
		Address:     contracts.NativeTokenSmartContractAddress,
		Topics:      topics,
		Data:        data,
		BlockNumber: blockNumber,
	}
}
