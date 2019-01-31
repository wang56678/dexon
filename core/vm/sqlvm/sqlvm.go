package sqlvm

import (
	"math/big"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/params"
)

type SQLVM struct {
	// Context provides auxiliary blockchain related information
	vm.Context
	// StateDB gives access to the underlying state
	StateDB vm.StateDB
	// Depth is the current call stack
	depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// abort is used to abort the SQLVM calling operations
	// NOTE: must be set atomically
	abort int32
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

func init() {
	vm.Register(vm.SQLVM, &SQLVM{})
}

func (sqlvm *SQLVM) Create(caller vm.ContractRef, code []byte, gas uint64, value *big.Int,
	in vm.Interpreter) ([]byte, common.Address, uint64, error) {
	// todo (jm) need to implemnt
	return nil, common.Address{}, gas, nil
}

func (sqlvm *SQLVM) Create2(caller vm.ContractRef, code []byte, gas uint64, endowment *big.Int, salt *big.Int,
	in vm.Interpreter) ([]byte, common.Address, uint64, error) {
	// todo (jm) need to implemnt
	return nil, common.Address{}, gas, nil
}
func (sqlvm *SQLVM) Call(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int,
	in vm.Interpreter) ([]byte, uint64, error) {
	// todo (jm) need to implemnt
	return nil, gas, nil
}
func (sqlvm *SQLVM) CallCode(caller vm.ContractRef, addr common.Address, input []byte, gas uint64,
	value *big.Int, in vm.Interpreter) ([]byte, uint64, error) {
	// todo (jm) need to implemnt
	return nil, gas, nil
}

func (sqlvm *SQLVM) DelegateCall(caller vm.ContractRef, addr common.Address, input []byte, gas uint64,
	in vm.Interpreter) ([]byte, uint64, error) {
	// todo (jm) need to implemnt
	return nil, gas, nil
}
func (sqlvm *SQLVM) StaticCall(caller vm.ContractRef, addr common.Address, input []byte, gas uint64,
	in vm.Interpreter) ([]byte, uint64, error) {
	// todo (jm) need to implemnt
	return nil, gas, nil
}
