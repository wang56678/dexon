// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"math/big"

	"github.com/dexon-foundation/dexon/common/math"
	"github.com/dexon-foundation/dexon/core/vm"
)

var (
	big1  = big.NewInt(1)
	big32 = big.NewInt(32)
)

func memorySha3(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), stack.Back(1))
}

func memoryCallDataCopy(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), stack.Back(2))
}

func memoryReturnDataCopy(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), stack.Back(2))
}

func memoryCodeCopy(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), stack.Back(2))
}

func memoryExtCodeCopy(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(1), stack.Back(3))
}

func memoryMLoad(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), big32)
}

func memoryMStore8(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), big1)
}

func memoryMStore(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), big32)
}

func memoryCreate(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(1), stack.Back(2))
}

func memoryCreate2(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(1), stack.Back(2))
}

func memoryCall(stack *vm.Stack) *big.Int {
	x := vm.CalcMemSize(stack.Back(5), stack.Back(6))
	y := vm.CalcMemSize(stack.Back(3), stack.Back(4))

	return math.BigMax(x, y)
}

func memoryDelegateCall(stack *vm.Stack) *big.Int {
	x := vm.CalcMemSize(stack.Back(4), stack.Back(5))
	y := vm.CalcMemSize(stack.Back(2), stack.Back(3))

	return math.BigMax(x, y)
}

func memoryStaticCall(stack *vm.Stack) *big.Int {
	x := vm.CalcMemSize(stack.Back(4), stack.Back(5))
	y := vm.CalcMemSize(stack.Back(2), stack.Back(3))

	return math.BigMax(x, y)
}

func memoryReturn(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), stack.Back(1))
}

func memoryRevert(stack *vm.Stack) *big.Int {
	return vm.CalcMemSize(stack.Back(0), stack.Back(1))
}

func memoryLog(stack *vm.Stack) *big.Int {
	mSize, mStart := stack.Back(1), stack.Back(0)
	return vm.CalcMemSize(mStart, mSize)
}
