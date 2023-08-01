package evmcore

import (
	"github.com/artheranet/arthera-node/internal/evmcore/vmcontext"
	"github.com/artheranet/arthera-node/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	params2 "github.com/ethereum/go-ethereum/params"
	"math/big"
)

// VMAddress is the address the VM uses to make internal calls to contracts
var VMAddress = common.Address{}

type evmRunner struct {
	newEVM       func(from common.Address) *vm.EVM
	state        vm.StateDB
	dontMeterGas bool
}

func (ev *evmRunner) Execute(recipient common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, err error) {
	evm := ev.newEVM(VMAddress)
	if ev.dontMeterGas {
		evm.StopGasMetering()
	}
	ret, _, err = evm.Call(vm.AccountRef(evm.Origin), recipient, input, gas, value)
	return ret, err
}

func (ev *evmRunner) ExecuteFrom(sender, recipient common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, err error) {
	evm := ev.newEVM(sender)
	if ev.dontMeterGas {
		evm.StopGasMetering()
	}
	ret, _, err = evm.Call(vm.AccountRef(sender), recipient, input, gas, value)
	return ret, err
}

func (ev *evmRunner) Query(recipient common.Address, input []byte, gas uint64) (ret []byte, err error) {
	evm := ev.newEVM(VMAddress)
	if ev.dontMeterGas {
		evm.StopGasMetering()
	}
	ret, _, err = evm.StaticCall(vm.AccountRef(evm.Origin), recipient, input, gas)
	return ret, err
}

func (ev *evmRunner) StopGasMetering() {
	ev.dontMeterGas = true
}

func (ev *evmRunner) StartGasMetering() {
	ev.dontMeterGas = false
}

func NewEVMRunner(chain DummyChain, chainconfig *params2.ChainConfig, header *EvmHeader, state vm.StateDB) vmcontext.EVMRunner {
	return &evmRunner{
		state: state,
		newEVM: func(from common.Address) *vm.EVM {
			// The EVM Context requires a msg, but the actual field values don't really matter for this case.
			// Putting in zero values for gas price and tx fee recipient
			blockContext := NewEVMBlockContext(header, chain, nil)
			txContext := vm.TxContext{
				Origin:   from,
				GasPrice: common.Big0,
			}
			return vm.NewEVM(blockContext, txContext, state, chainconfig, params.DefaultVMConfig)
		},
	}
}
