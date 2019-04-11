// Copyright 2018 The dexon-consensus Authors
// This file is part of the dexon-consensus library.
//
// The dexon-consensus library is free software: you can redistribute it
// and/or modify it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// The dexon-consensus library is distributed in the hope that it will be
// useful, but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser
// General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the dexon-consensus library. If not, see
// <http://www.gnu.org/licenses/>.

package dvm

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/state"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/ethdb"
)

func TestWasmFinish(t *testing.T) {
	// (module
	//   (import "ethereum" "finish" (func (param i32 i32)))
	//   (func (export "main")
	//     i32.const 0
	//     i32.const 4
	//     call 0)
	//   (memory (export "memory") 1)
	//   (data (i32.const 0) "TEST"))
	code := "020061736d0100000001090260027f7f0060000002130108657468657265756d0666696e6973680000030201010503010001071102046d61696e0001066d656d6f727902000a0a0108004100410410000b0b0a010041000b0454455354"
	binary := common.Hex2Bytes(code)

	caller := getRandomAddress()
	statedb, err := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	config := Config{Metering: false}
	gas := uint64(1000000)
	value := big.NewInt(0)
	myDVM := NewDVM(statedb, config)
	ret, _, _, err := vm.Create(vm.AccountRef(caller), binary, gas, value, myDVM)

	assert.Equal(t, []byte("TEST"), ret)
	assert.NoError(t, err)
}

func TestWasmStorageStore(t *testing.T) {
	// (module
	//   (import "ethereum" "storageStore" (func (param i32 i32)))
	//   (func (export "main")
	//     i32.const 0
	//     i32.const 32
	//     call 0)
	//   (memory (export "memory") 1)
	//   (data (i32.const 0) "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"))
	code := "020061736d0100000001090260027f7f0060000002190108657468657265756d0c73746f7261676553746f72650000030201010503010001071102046d61696e0001066d656d6f727902000a0a0108004100412010000b0b46010041000b4041414141414141414141414141414141414141414141414141414141414141414646464646464646464646464646464646464646464646464646464646464646"
	binary := common.Hex2Bytes(code)

	caller := getRandomAddress()
	statedb, err := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	config := Config{Metering: false}
	gas := uint64(1000000)
	value := big.NewInt(0)
	myDVM := NewDVM(statedb, config)
	ret, addr, _, err := vm.Create(vm.AccountRef(caller), binary, gas, value, myDVM)

	expectedStorageResult := common.BytesToHash([]byte("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"))
	storageResult := myDVM.StateDB().GetState(addr, common.BytesToHash([]byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")))

	assert.Empty(t, ret)
	assert.Equal(t, expectedStorageResult, storageResult)
	assert.NoError(t, err)
}

func TestWasmERC20CreateToken(t *testing.T) {
	// create contract
	binary := common.Hex2Bytes(TestWasmERC20Code)
	caller := getRandomAddress()
	statedb, err := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	config := Config{Metering: false}
	gas := uint64(900000000)
	value := big.NewInt(0)
	myDVM := NewDVM(statedb, config)
	ret, addr, _, err := vm.Create(vm.AccountRef(caller), binary, gas, value, myDVM)
	assert.NoError(t, err)

	// constructor
	input := common.Hex2Bytes("90fa17bb")
	ret, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.NoError(t, err)

	// createToken
	accountAddr := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
	createTokenValue := "00000000000000000000000000000000000000000000000000000000000000FF"
	input = common.Hex2Bytes("a69e894e" + accountAddr + createTokenValue)
	ret, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.NoError(t, err)

	// balanceOf
	input = common.Hex2Bytes("70a08231" + accountAddr)
	ret, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.Equal(t, common.Hex2Bytes(createTokenValue), ret)
	assert.NoError(t, err)
}

func TestWasmERC20Transfer(t *testing.T) {
	// create contract
	binary := common.Hex2Bytes(TestWasmERC20Code)
	caller := getRandomAddress()
	statedb, err := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	config := Config{Metering: false}
	gas := uint64(900000000)
	value := big.NewInt(0)
	myDVM := NewDVM(statedb, config)
	ret, addr, _, err := vm.Create(vm.AccountRef(caller), binary, gas, value, myDVM)
	assert.NoError(t, err)

	// constructor
	input := common.Hex2Bytes("90fa17bb")
	ret, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.NoError(t, err)

	// create token:
	//   balance[0xABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB] = 0xFF
	fromAccountAddr := "000000000000000000000000ABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
	createTokenValue := "00000000000000000000000000000000000000000000000000000000000000FF"
	input = common.Hex2Bytes("a69e894e" + fromAccountAddr + createTokenValue)
	ret, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.NoError(t, err)

	// transfer token:
	//   balance[0xABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB] -= 0x0F
	//   balance[0xACCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC] += 0x0F
	toAccountAddr := "000000000000000000000000ACCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"
	transferTokenValue := "000000000000000000000000000000000000000000000000000000000000000F"
	input = common.Hex2Bytes("a9059cbb" + toAccountAddr + transferTokenValue)
	_, _, err = vm.Call(vm.AccountRef(common.HexToAddress(fromAccountAddr)), addr, input, gas, value, myDVM)
	assert.NoError(t, err)

	// check balance of two acounts
	//   balance[0xABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB] should be 0xF0
	//   balance[0xACCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC] should be 0x0F
	input = common.Hex2Bytes("70a08231" + fromAccountAddr)
	ret, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.Equal(t, common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000F0"), ret)
	assert.NoError(t, err)
	input = common.Hex2Bytes("70a08231" + toAccountAddr)
	ret, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.Equal(t, common.Hex2Bytes("000000000000000000000000000000000000000000000000000000000000000F"), ret)
	assert.NoError(t, err)
}

func BenchmarkWasmERC20Transfer(b *testing.B) {
	// create contract
	binary := common.Hex2Bytes(TestWasmERC20Code)
	caller := getRandomAddress()
	statedb, err := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	config := Config{Metering: false}
	gas := uint64(900000000)
	value := big.NewInt(0)
	myDVM := NewDVM(statedb, config)
	_, addr, _, err := vm.Create(vm.AccountRef(caller), binary, gas, value, myDVM)
	assert.NoError(b, err)

	// constructor
	input := common.Hex2Bytes("90fa17bb")
	_, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.NoError(b, err)

	// create token:
	fromAccountAddr := "000000000000000000000000ABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
	createTokenValue := "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
	input = common.Hex2Bytes("a69e894e" + fromAccountAddr + createTokenValue)
	_, _, err = vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	assert.NoError(b, err)

	// transfer token:
	toAccountAddr := "000000000000000000000000ACCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"
	transferTokenValue := "0000000000000000000000000000000000000000000000000000000000000001"
	input = common.Hex2Bytes("a9059cbb" + toAccountAddr + transferTokenValue)
	fromAccountRef := vm.AccountRef(common.HexToAddress(fromAccountAddr))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.Call(fromAccountRef, addr, input, gas, value, myDVM)
	}
	b.StopTimer()
}

func BenchmarkWasmERC20TransferSimulation(b *testing.B) {
	// void _main() {
	//     i32 mem[1024];
	//     i32 size = getCallDataSize();
	//     callDataCopy(mem, 0, 4);
	//     callDataCopy(mem, 4, 64);
	//     getCaller(mem);
	//     storageLoad(0, mem);
	//     storageLoad(0, mem);
	//     storageStore(mem, 0);
	//     storageStore(mem, 0);
	//     finish(mem, 32);
	// }

	// create contract
	code := "020061736d0100000001090260027f7f0060000002130108657468657265756d0666696e6973680000030201010503010001071102046d61696e0001066d656d6f727902000a0b010900410041a10210000b0ba902010041000ba202020061736d010000000117056000017f60017f0060037f7f7f0060027f7f00600000028a010608657468657265756d0b73746f726167654c6f6164000308657468657265756d0c63616c6c44617461436f7079000208657468657265756d0c73746f7261676553746f7265000308657468657265756d0967657443616c6c6572000108657468657265756d0f67657443616c6c4461746153697a65000008657468657265756d0666696e69736800030302010405030100010607017f014180160b071102066d656d6f72790200046d61696e00060a4c014a01017f2300210023004180206a240010041a20004100410410012000410441c000100120001003410020001000410020001000200041001002200041001002200041201005200024000b"
	binary := common.Hex2Bytes(code)
	caller := getRandomAddress()
	statedb, err := state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	config := Config{Metering: false}
	gas := uint64(900000000)
	value := big.NewInt(0)
	myDVM := NewDVM(statedb, config)
	_, addr, _, err := vm.Create(vm.AccountRef(caller), binary, gas, value, myDVM)
	assert.NoError(b, err)

	// transfer token:
	toAccountAddr := "000000000000000000000000ACCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"
	transferTokenValue := "0000000000000000000000000000000000000000000000000000000000000001"
	input := common.Hex2Bytes("a9059cbb" + toAccountAddr + transferTokenValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.Call(vm.AccountRef(caller), addr, input, gas, value, myDVM)
	}
	b.StopTimer()
}

func getRandomAddress() common.Address {
	randomBytes := make([]byte, common.HashLength)
	rand.Read(randomBytes)
	return common.BytesToAddress(randomBytes)
}
