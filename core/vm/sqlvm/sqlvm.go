package sqlvm

import (
	"math/big"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/params"
)

// SQLVM implements the required VM interface.
type SQLVM struct {
	// Context provides auxiliary blockchain related information
	*vm.Context
	// StateDB gives access to the underlying state
	StateDB vm.StateDB

	// abort is used to abort the SQLVM calling operations
	// NOTE: must be set atomically
	abort int32
}

func init() {
	vm.Register(vm.SQLVM, NewSQLVM)
}

// NewSQLVM is the SQLVM constructor.
func NewSQLVM(context *vm.Context, stateDB vm.StateDB, chainConfig *params.ChainConfig, vmConfig interface{}) vm.VM {
	return &SQLVM{Context: context, StateDB: stateDB}
}

// Create creates SQL contract.
func (sqlvm *SQLVM) Create(caller vm.ContractRef, code []byte, gas uint64, value *big.Int,
	pack *vm.ExecPack) ([]byte, common.Address, uint64, error) {
	// todo (jm) need to implemnt
	return nil, common.Address{}, gas, nil
}

// Create2 mock interface.
func (sqlvm *SQLVM) Create2(caller vm.ContractRef, code []byte, gas uint64, endowment *big.Int, salt *big.Int,
	pack *vm.ExecPack) ([]byte, common.Address, uint64, error) {
	// todo (jm) need to implemnt
	return nil, common.Address{}, gas, nil
}

// Call is the entry to call SQLVM contract.
func (sqlvm *SQLVM) Call(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int,
	pack *vm.ExecPack) ([]byte, uint64, error) {
	// todo (jm) need to implemnt
	return nil, gas, nil
}

// CallCode mock interface.
func (sqlvm *SQLVM) CallCode(caller vm.ContractRef, addr common.Address, input []byte, gas uint64,
	value *big.Int, pack *vm.ExecPack) ([]byte, uint64, error) {
	// todo (jm) need to implemnt
	return nil, gas, nil
}

// DelegateCall mock interface.
func (sqlvm *SQLVM) DelegateCall(caller vm.ContractRef, addr common.Address, input []byte, gas uint64,
	pack *vm.ExecPack) ([]byte, uint64, error) {
	// todo (jm) need to implemnt
	return nil, gas, nil
}

// StaticCall is the entry for read-only call on SQL contract.
func (sqlvm *SQLVM) StaticCall(caller vm.ContractRef, addr common.Address, input []byte, gas uint64,
	pack *vm.ExecPack) ([]byte, uint64, error) {
	// todo (jm) need to implemnt
	return nil, gas, nil
}
