package vm

import (
	"math/big"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/params"
)

const (
	// EVM enum
	EVM = uint8(iota)

	// SQLVM enum
	SQLVM
)

var (
	// MULTIVM flag
	MULTIVM = true
)

// NUMS represents the number of supported VM type.
const NUMS = 2

// VM Create and Call interface.
type VM interface {
	Create(ContractRef, []byte, uint64, *big.Int,
		*ExecPack) ([]byte, common.Address, uint64, error)
	Create2(ContractRef, []byte, uint64, *big.Int, *big.Int,
		*ExecPack) ([]byte, common.Address, uint64, error)
	Call(ContractRef, common.Address, []byte, uint64, *big.Int,
		*ExecPack) ([]byte, uint64, error)
	CallCode(ContractRef, common.Address, []byte, uint64,
		*big.Int, *ExecPack) ([]byte, uint64, error)
	DelegateCall(ContractRef, common.Address, []byte, uint64,
		*ExecPack) ([]byte, uint64, error)
	StaticCall(ContractRef, common.Address, []byte, uint64,
		*ExecPack) ([]byte, uint64, error)
}

type createFunc func(*Context, StateDB, *params.ChainConfig, interface{}) VM

// ExecPack contains runtime context, stateDB, chain config, VM list and VM configs.
type ExecPack struct {
	Context     *Context
	StateDB     StateDB
	ChainConfig *params.ChainConfig
	VMList      [NUMS]VM
	VMConfig    [NUMS]interface{}
}

var createFuncs [NUMS]createFunc

// Register registers VM create function.
func Register(idx uint8, c createFunc) {
	createFuncs[idx] = c
}

// NewExecPack creates a ExecPack instance, and create all VM instance.
func NewExecPack(context *Context, stateDB StateDB, chainConfig *params.ChainConfig, vmConfigs [NUMS]interface{}) ExecPack {
	p := ExecPack{
		Context:     context,
		StateDB:     stateDB,
		ChainConfig: chainConfig,
		VMConfig:    vmConfigs,
	}
	context.ExecPack = &p
	for i := 0; i < NUMS; i++ {
		if createFuncs[i] != nil {
			p.VMList[i] = createFuncs[i](context, stateDB, chainConfig, vmConfigs[i])
		}
	}
	return p
}

// Create is the entry for multiple VMs' Create.
func Create(caller ContractRef, code []byte, gas uint64, value *big.Int,
	p *ExecPack) (ret []byte, contractAddr common.Address,
	leftOverGas uint64, err error) {

	v, code := getVMAndCode(code)
	return p.VMList[v].Create(caller, code, gas, value, p)
}

// Create2 is the entry for multiple VMs' Create2.
func Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int,
	salt *big.Int, p *ExecPack) (ret []byte,
	contractAddr common.Address, leftOverGas uint64, err error) {

	v, code := getVMAndCode(code)
	return p.VMList[v].Create2(caller, code, gas, endowment, salt, p)
}

// Call is the entry for multiple VMs' Call.
func Call(caller ContractRef, addr common.Address, input []byte, gas uint64,
	value *big.Int, p *ExecPack) (ret []byte, leftOverGas uint64, err error) {

	code := p.StateDB.GetCode(addr)
	v, _ := getVMAndCode(code)
	return p.VMList[v].Call(caller, addr, input, gas, value, p)
}

// CallCode is the entry for multiple VMs' CallCode.
func CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64,
	value *big.Int, p *ExecPack) (ret []byte, leftOverGas uint64, err error) {

	code := p.StateDB.GetCode(addr)
	v, _ := getVMAndCode(code)
	return p.VMList[v].CallCode(caller, addr, input, gas, value, p)
}

// DelegateCall is the entry for multiple VMs' DelegateCall.
func DelegateCall(caller ContractRef, addr common.Address, input []byte,
	gas uint64, p *ExecPack) (ret []byte, leftOverGas uint64, err error) {

	code := p.StateDB.GetCode(addr)
	v, _ := getVMAndCode(code)
	return p.VMList[v].DelegateCall(caller, addr, input, gas, p)
}

// StaticCall is the entry for multiple VMs' StaticCall.
func StaticCall(caller ContractRef, addr common.Address, input []byte,
	gas uint64, p *ExecPack) (ret []byte, leftOverGas uint64, err error) {

	code := p.StateDB.GetCode(addr)
	v, _ := getVMAndCode(code)
	return p.VMList[v].StaticCall(caller, addr, input, gas, p)
}

func getVMAndCode(code []byte) (uint8, []byte) {
	if MULTIVM && len(code) > 0 {
		switch code[0] {
		case EVM, SQLVM:
			return code[0], code[1:]
		default:
			return EVM, code
		}
	}
	return EVM, code
}
