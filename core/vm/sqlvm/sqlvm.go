package sqlvm

import (
	"math/big"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/crypto"
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

func (sqlvm *SQLVM) Create(caller vm.ContractRef, code []byte, gas uint64,
	value *big.Int) (ret []byte, contractAddr common.Address,
	leftOverGas uint64, err error) {

	contractAddr = crypto.CreateAddress(caller.Address(), sqlvm.StateDB.GetNonce(caller.Address()))
	return sqlvm.create(caller, &vm.CodeAndHash{Code: code}, gas, value, contractAddr)
}

// create creates a new contract using code as deployment code.
func (sqlvm *SQLVM) create(caller vm.ContractRef, codeAndHash *vm.CodeAndHash, gas uint64,
	value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the
	if sqlvm.depth > int(params.CallCreateDepth) {
		return nil, common.Address{}, gas, vm.ErrDepth
	}
	// TODO (JM) implement create database contract function
	return nil, common.Address{}, gas, nil
}
