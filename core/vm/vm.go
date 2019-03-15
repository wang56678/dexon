package vm

import (
	"math/big"

	"github.com/dexon-foundation/dexon/common"
)

const (
	EVM = byte(iota)
	SQLVM
)

var (
	MULTIVM = true
)

type VM interface {
	Create(ContractRef, []byte, uint64, *big.Int,
		Interpreter) ([]byte, common.Address, uint64, error)
	Create2(ContractRef, []byte, uint64, *big.Int, *big.Int,
		Interpreter) ([]byte, common.Address, uint64, error)
	Call(ContractRef, common.Address, []byte, uint64, *big.Int,
		Interpreter) ([]byte, uint64, error)
	CallCode(ContractRef, common.Address, []byte, uint64,
		*big.Int, Interpreter) ([]byte, uint64, error)
	DelegateCall(ContractRef, common.Address, []byte, uint64,
		Interpreter) ([]byte, uint64, error)
	StaticCall(ContractRef, common.Address, []byte, uint64,
		Interpreter) ([]byte, uint64, error)
}

type Interpreter interface {
	StateDB() StateDB
}

var vmList map[byte]VM

func init() {
	vmList = make(map[byte]VM)
}
func Register(vmType byte, vm VM) {
	vmList[vmType] = vm
}
func Create(caller ContractRef, code []byte, gas uint64, value *big.Int,
	interpreter Interpreter) (ret []byte, contractAddr common.Address,
	leftOverGas uint64, err error) {

	name, code := getVMAndCode(code)
	return vmList[name].Create(caller, code, gas, value, interpreter)
}

func Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int,
	salt *big.Int, interpreter Interpreter) (ret []byte,
	contractAddr common.Address, leftOverGas uint64, err error) {

	name, code := getVMAndCode(code)
	return vmList[name].Create2(caller, code, gas, endowment, salt, interpreter)
}

func Call(caller ContractRef, addr common.Address, input []byte, gas uint64,
	value *big.Int, interpreter Interpreter) (ret []byte, leftOverGas uint64, err error) {

	code := interpreter.StateDB().GetCode(addr)
	name, _ := getVMAndCode(code)
	return vmList[name].Call(caller, addr, input, gas, value, interpreter)
}

func CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64,
	value *big.Int, interpreter Interpreter) (ret []byte, leftOverGas uint64, err error) {

	code := interpreter.StateDB().GetCode(addr)
	name, _ := getVMAndCode(code)
	return vmList[name].CallCode(caller, addr, input, gas, value, interpreter)
}

func DelegateCall(caller ContractRef, addr common.Address, input []byte,
	gas uint64, interpreter Interpreter) (ret []byte, leftOverGas uint64, err error) {

	code := interpreter.StateDB().GetCode(addr)
	name, _ := getVMAndCode(code)
	return vmList[name].DelegateCall(caller, addr, input, gas, interpreter)
}

func StaticCall(caller ContractRef, addr common.Address, input []byte,
	gas uint64, interpreter Interpreter) (ret []byte, leftOverGas uint64, err error) {

	code := interpreter.StateDB().GetCode(addr)
	name, _ := getVMAndCode(code)
	return vmList[name].StaticCall(caller, addr, input, gas, interpreter)
}

func getVMAndCode(code []byte) (byte, []byte) {
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
