// Copyright 2018 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// This file lists the EEI functions, so that they can be bound to any
// ewasm-compatible module, as well as the types of these functions

package dvm

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/crypto"
	"github.com/dexon-foundation/dexon/params"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

type terminationType int

// List of termination reasons
const (
	TerminateFinish = iota
	TerminateRevert
	TerminateSuicide
	TerminateInvalid
)

// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
var (
	emptyCodeHash           = crypto.Keccak256Hash(nil)
	errExecutionReverted    = errors.New("dvm: execution reverted")
	errMaxCodeSizeExceeded  = errors.New("dvm: max code size exceeded")
	u256Len                 = 32
	u128Len                 = 16
	maxCallDepth            = 1024
	sentinelContractAddress = "0x000000000000000000000000000000000000000a"
)

type DVM struct {
	// Context provides auxiliary blockchain related information
	vm.Context
	// StateDB gives access to the underlying state
	statedb vm.StateDB
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

	// vm is wasm VM to execute bytecode
	vm *exec.VM

	gasTable        params.GasTable
	contract        *vm.Contract
	returnData      []byte
	terminationType terminationType
	staticMode      bool

	metering           bool
	meteringContract   *vm.Contract
	meteringModule     *wasm.Module
	meteringStartIndex int64
}

type Config struct {
	Metering bool
}

func init() {
	vm.Register(vm.DVM, &DVM{})
}

// NewDVM creates a new wagon-based ewasm vm.
func NewDVM(statedb vm.StateDB, config Config) *DVM {
	ctx := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
	}
	dvm := &DVM{
		Context:  ctx,
		statedb:  statedb,
		metering: config.Metering,
	}

	if dvm.metering {
		meteringContractAddress := common.HexToAddress(sentinelContractAddress)
		meteringCode := dvm.StateDB().GetCode(meteringContractAddress)

		var err error
		dvm.meteringModule, err = wasm.ReadModule(bytes.NewReader(meteringCode), WrappedModuleResolver(dvm))
		if err != nil {
			panic(fmt.Sprintf("Error loading the metering contract: %v", err))
		}
		// TODO when the metering contract abides by that rule, check that it
		// only exports "main" and "memory".
		dvm.meteringStartIndex = int64(dvm.meteringModule.Export.Entries["main"].Index)
		mainSig := dvm.meteringModule.FunctionIndexSpace[dvm.meteringStartIndex].Sig
		if len(mainSig.ParamTypes) != 0 || len(mainSig.ReturnTypes) != 0 {
			panic(fmt.Sprintf("Invalid main function for the metering contract: index=%d sig=%v", dvm.meteringStartIndex, mainSig))
		}
	}

	return dvm
}

func (dvm *DVM) StateDB() vm.StateDB {
	return dvm.statedb
}

func (dvm *DVM) Create(caller vm.ContractRef, code []byte, gas uint64, value *big.Int, in vm.Interpreter) ([]byte, common.Address, uint64, error) {
	contractAddr := crypto.CreateAddress(caller.Address(), in.(*DVM).statedb.GetNonce(caller.Address()))
	return in.(*DVM).create(caller, &vm.CodeAndHash{Code: code}, gas, value, contractAddr)
}

func (dvm *DVM) Create2(caller vm.ContractRef, code []byte, gas uint64, endowment *big.Int, salt *big.Int, in vm.Interpreter) ([]byte, common.Address, uint64, error) {
	codeAndHash := &vm.CodeAndHash{Code: code}
	contractAddr := crypto.CreateAddress2(caller.Address(), common.BigToHash(salt), codeAndHash.Hash().Bytes())
	return dvm.create(caller, codeAndHash, gas, endowment, contractAddr)
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (dvm *DVM) Call(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int, in vm.Interpreter) ([]byte, uint64, error) {
	// TODO: do we need these checks?
	// 	if evm.vmConfig.NoRecursion && evm.depth > 0 {
	// 		return nil, gas, nil
	// 	}

	// Fail if we're trying to execute above the call depth limit
	if in.(*DVM).depth > int(params.CallCreateDepth) {
		return nil, gas, vm.ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	if !in.(*DVM).Context.CanTransfer(in.(*DVM).StateDB(), caller.Address(), value) {
		return nil, gas, vm.ErrInsufficientBalance
	}

	var (
		to       = vm.AccountRef(addr)
		snapshot = in.(*DVM).StateDB().Snapshot()
	)
	if !in.(*DVM).StateDB().Exist(addr) {
		precompiles := vm.PrecompiledContractsByzantium
		if precompiles[addr] == nil && value.Sign() == 0 {
			return nil, gas, nil
		}
		in.(*DVM).StateDB().CreateAccount(addr)
	}
	in.(*DVM).Transfer(in.(*DVM).StateDB(), caller.Address(), to.Address(), value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := vm.NewContract(caller, to, value, gas)
	code := in.(*DVM).StateDB().GetCode(addr)
	if len(code) > 0 && code[0] == vm.DVM && vm.MULTIVM {
		code = code[1:]
	}
	codeAndHash := vm.CodeAndHash{Code: code}
	contract.SetCodeOptionalHash(&addr, &codeAndHash)

	// Even if the account has no code, we need to continue because it might be a precompile
	ret, err := in.(*DVM).run(contract, input, false)

	// When an error was returned by the DVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		in.(*DVM).StateDB().RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// CallCode executes the contract associated with the addr with the given input
// as parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
//
// CallCode differs from Call in the sense that it executes the given address'
// code with the caller as context.
func (dvm *DVM) CallCode(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int, in vm.Interpreter) ([]byte, uint64, error) {
	// TODO: do we need these checks?
	// 	if evm.vmConfig.NoRecursion && evm.depth > 0 {
	// 		return nil, gas, nil
	// 	}

	// Fail if we're trying to execute above the call depth limit
	if dvm.depth > int(params.CallCreateDepth) {
		return nil, gas, vm.ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	if !dvm.Context.CanTransfer(dvm.StateDB(), caller.Address(), value) {
		return nil, gas, vm.ErrInsufficientBalance
	}

	var (
		snapshot = dvm.StateDB().Snapshot()
		to       = vm.AccountRef(caller.Address())
	)

	// initialise a new contract and set the code that is to be used by the
	// DVM. The contract is a scoped environment for this execution context
	// only.
	contract := vm.NewContract(caller, to, value, gas)
	code := dvm.StateDB().GetCode(addr)
	if len(code) > 0 && vm.MULTIVM {
		code = code[1:]
	}
	codeAndHash := vm.CodeAndHash{Code: code}
	contract.SetCodeOptionalHash(&addr, &codeAndHash)

	ret, err := dvm.run(contract, input, false)
	if err != nil {
		dvm.StateDB().RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (dvm *DVM) DelegateCall(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, in vm.Interpreter) ([]byte, uint64, error) {
	// TODO: do we need these checks?
	// 	if evm.vmConfig.NoRecursion && evm.depth > 0 {
	// 		return nil, gas, nil
	// 	}

	// Fail if we're trying to execute above the call depth limit
	if dvm.depth > int(params.CallCreateDepth) {
		return nil, gas, vm.ErrDepth
	}

	var (
		snapshot = dvm.StateDB().Snapshot()
		to       = vm.AccountRef(caller.Address())
	)

	// Initialise a new contract and make initialise the delegate values
	contract := vm.NewContract(caller, to, nil, gas).AsDelegate()
	code := dvm.StateDB().GetCode(addr)
	if len(code) > 0 && vm.MULTIVM {
		code = code[1:]
	}
	codeAndHash := vm.CodeAndHash{Code: code}
	contract.SetCodeOptionalHash(&addr, &codeAndHash)

	ret, err := dvm.run(contract, input, false)
	if err != nil {
		dvm.StateDB().RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (dvm *DVM) StaticCall(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, in vm.Interpreter) ([]byte, uint64, error) {
	// TODO: do we need these checks?
	// 	if evm.vmConfig.NoRecursion && evm.depth > 0 {
	// 		return nil, gas, nil
	// 	}

	// Fail if we're trying to execute above the call depth limit
	if dvm.depth > int(params.CallCreateDepth) {
		return nil, gas, vm.ErrDepth
	}

	var (
		to       = vm.AccountRef(addr)
		snapshot = dvm.StateDB().Snapshot()
	)

	// Initialise a new contract and set the code that is to be used by the
	// DVM. The contract is a scoped environment for this execution context
	// only.
	contract := vm.NewContract(caller, to, new(big.Int), gas)
	code := dvm.StateDB().GetCode(addr)
	if len(code) > 0 && vm.MULTIVM {
		code = code[1:]
	}
	codeAndHash := vm.CodeAndHash{Code: code}
	contract.SetCodeOptionalHash(&addr, &codeAndHash)

	// We do an AddBalance of zero here, just in order to trigger a touch.
	// This doesn't matter on Mainnet, where all empties are gone at the time of Byzantium,
	// but is the correct thing to do and matters on other networks, in tests, and potential
	// future scenarios
	dvm.StateDB().AddBalance(addr, new(big.Int))

	// When an error was returned by the DVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in Homestead this also counts for code storage gas errors.
	ret, err := dvm.run(contract, input, true)
	if err != nil {
		dvm.StateDB().RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

func (dvm *DVM) create(caller vm.ContractRef, codeAndHash *vm.CodeAndHash, gas uint64, value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if dvm.depth > int(params.CallCreateDepth) {
		return nil, common.Address{}, gas, vm.ErrDepth
	}
	if !dvm.Context.CanTransfer(dvm.statedb, caller.Address(), value) {
		return nil, common.Address{}, gas, vm.ErrInsufficientBalance
	}
	nonce := dvm.statedb.GetNonce(caller.Address())
	dvm.statedb.SetNonce(caller.Address(), nonce+1)

	// Ensure there's no existing contract already at the designated address
	contractHash := dvm.statedb.GetCodeHash(address)
	if dvm.statedb.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, vm.ErrContractAddressCollision
	}

	// Create a new account on the state
	snapshot := dvm.statedb.Snapshot()
	dvm.statedb.CreateAccount(address)
	dvm.statedb.SetNonce(address, 1)
	dvm.Context.Transfer(dvm.statedb, caller.Address(), address, value)

	// initialise a new contract and set the code that is to be used by the
	// DVM. The contract is a scoped environment for this execution context
	// only.
	contract := vm.NewContract(caller, vm.AccountRef(address), value, gas)
	meteredCode, err := dvm.PreContractCreation(codeAndHash.Code, contract)
	if err != nil {
		return nil, address, gas, nil
	}
	codeAndHash.Code = meteredCode
	contract.SetCodeOptionalHash(&address, codeAndHash)

	// TODO: do we need these checks?
	// 	if evm.vmConfig.NoRecursion && evm.depth > 0 {
	// 		return nil, address, gas, nil
	// 	}

	ret, err := dvm.run(contract, nil, false)

	// The new contract needs to be metered after it has executed the constructor
	if err != nil {
		if dvm.CanRun(contract.Code) {
			ret, err = dvm.PostContractCreation(ret)
		}
	}

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			dvm.statedb.SetCode(address, ret)
		} else {
			err = vm.ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the DVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || err != nil {
		dvm.statedb.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}

	return ret, address, contract.Gas, err
}

// Run loops and evaluates the contract's code with the given input data and returns
// the return byte-slice and an error if one occurred.
func (dvm *DVM) run(contract *vm.Contract, input []byte, ro bool) ([]byte, error) {
	// Take care of running precompile contracts
	if contract.CodeAddr != nil {
		precompiles := vm.PrecompiledContractsByzantium
		if p := precompiles[*contract.CodeAddr]; p != nil {
			return vm.RunPrecompiledContract(p, input, contract)
		}
	}

	// Increment the call depth which is restricted to 1024
	dvm.depth++
	defer func() { dvm.depth-- }()

	dvm.contract = contract
	dvm.contract.Input = input

	module, err := wasm.ReadModule(bytes.NewReader(contract.Code), WrappedModuleResolver(dvm))
	if err != nil {
		dvm.terminationType = TerminateInvalid
		return nil, fmt.Errorf("Error decoding module at address %s: %v", contract.Address().Hex(), err)
	}

	wavm, err := exec.NewVM(module)
	if err != nil {
		dvm.terminationType = TerminateInvalid
		return nil, fmt.Errorf("could not create the vm: %v", err)
	}
	wavm.RecoverPanic = true
	dvm.vm = wavm

	mainIndex, err := validateModule(module)
	if err != nil {
		dvm.terminationType = TerminateInvalid
		return nil, err
	}

	// Check input and output types
	sig := module.FunctionIndexSpace[mainIndex].Sig
	if len(sig.ParamTypes) == 0 && len(sig.ReturnTypes) == 0 {
		_, err = wavm.ExecCode(int64(mainIndex))

		if err != nil && err != errExecutionReverted {
			dvm.terminationType = TerminateInvalid
		}

		if dvm.StateDB().HasSuicided(contract.Address()) {
			dvm.StateDB().AddRefund(params.SuicideRefundGas)
			err = nil
		}

		return dvm.returnData, err
	}

	dvm.terminationType = TerminateInvalid
	return nil, errors.New("Could not find a suitable 'main' function in that contract")
}

func validateModule(m *wasm.Module) (int, error) {
	// A module should not have a start section
	if m.Start != nil {
		return -1, fmt.Errorf("Module has a start section")
	}

	// Only two exports are authorized: "main" and "memory"
	if m.Export == nil {
		return -1, fmt.Errorf("Module has no exports instead of 2")
	}
	if len(m.Export.Entries) != 2 {
		return -1, fmt.Errorf("Module has %d exports instead of 2", len(m.Export.Entries))
	}

	mainIndex := -1
	for name, entry := range m.Export.Entries {
		switch name {
		case "main":
			if entry.Kind != wasm.ExternalFunction {
				return -1, fmt.Errorf("Main is not a function in module")
			}
			mainIndex = int(entry.Index)
		case "memory":
			if entry.Kind != wasm.ExternalMemory {
				return -1, fmt.Errorf("'memory' is not a memory in module")
			}
		default:
			return -1, fmt.Errorf("A symbol named %s has been exported. Only main and memory should exist", name)
		}
	}

	if m.Import != nil {
	OUTER:
		for _, entry := range m.Import.Entries {
			if entry.ModuleName == "ethereum" {
				if entry.Type.Kind() == wasm.ExternalFunction {
					for _, name := range eeiFunctionList {
						if name == entry.FieldName {
							continue OUTER
						}
					}
					return -1, fmt.Errorf("%s could not be found in the list of ethereum-provided functions", entry.FieldName)
				}
			}
		}
	}

	return mainIndex, nil
}

// CanRun checks the binary for a WASM header and accepts the binary blob
// if it matches.
func (dvm *DVM) CanRun(file []byte) bool {
	// Check the header
	if len(file) < 4 || string(file[:4]) != "\000asm" {
		return false
	}

	return true
}

// PreContractCreation meters the contract's its init code before it
// is run.
func (dvm *DVM) PreContractCreation(code []byte, contract *vm.Contract) ([]byte, error) {
	savedContract := dvm.contract
	dvm.contract = contract

	defer func() {
		dvm.contract = savedContract
	}()

	if dvm.metering {
		metered, _, err := sentinel(dvm, code)
		if len(metered) < 5 || err != nil {
			return nil, fmt.Errorf("Error metering the init contract code, err=%v", err)
		}
		return metered, nil
	}
	return code, nil
}

// PostContractCreation meters the contract once its init code has
// been run. It also validates the module's format before it is to
// be committed to disk.
func (dvm *DVM) PostContractCreation(code []byte) ([]byte, error) {
	// If a REVERT has been encountered, then return the code and
	if dvm.terminationType == TerminateRevert {
		return nil, errExecutionReverted
	}

	if dvm.CanRun(code) {
		if dvm.metering {
			meteredCode, _, err := sentinel(dvm, code)
			code = meteredCode
			if len(code) < 5 || err != nil {
				return nil, fmt.Errorf("Error metering the generated contract code, err=%v", err)
			}

			if len(code) < 8 {
				return nil, fmt.Errorf("Invalid contract code")
			}
		}

		if len(code) > 8 {
			// Check the validity of the module
			m, err := wasm.DecodeModule(bytes.NewReader(code))
			if err != nil {
				return nil, fmt.Errorf("Error decoding the module produced by init code: %v", err)
			}

			_, err = validateModule(m)
			if err != nil {
				dvm.terminationType = TerminateInvalid
				return nil, err
			}
		}
	}

	return code, nil
}
