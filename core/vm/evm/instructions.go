// Copyright 2015 The go-ethereum Authors
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
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/common/math"
	"github.com/dexon-foundation/dexon/core/types"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/crypto"
	"github.com/dexon-foundation/dexon/params"
	"golang.org/x/crypto/sha3"
)

var (
	bigZero                  = new(big.Int)
	big2                     = big.NewInt(2)
	big256                   = big.NewInt(256)
	tt255                    = math.BigPow(2, 255)
	errWriteProtection       = errors.New("evm: write protection")
	errReturnDataOutOfBounds = errors.New("evm: return data out of bounds")
	errExecutionReverted     = errors.New("evm: execution reverted")
	errMaxCodeSizeExceeded   = errors.New("evm: max code size exceeded")
	power2                   = make([]*big.Int, 256)
)

func init() {
	cur := big.NewInt(1)
	for i := 0; i < 256; i++ {
		power2[i] = new(big.Int).Set(cur)
		cur = new(big.Int).Mul(cur, big2)
	}
}

func opAdd(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()
	math.U256(y.Add(x, y))

	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opSub(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()
	math.U256(y.Sub(x, y))

	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opMul(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Pop()
	stack.Push(math.U256(x.Mul(x, y)))

	interpreter.evm.IntPool.Put(y)

	return nil, nil
}

func opDiv(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()
	if y.Sign() != 0 {
		math.U256(y.Div(x, y))
	} else {
		y.SetUint64(0)
	}
	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opSdiv(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := math.S256(stack.Pop()), math.S256(stack.Pop())
	res := interpreter.evm.IntPool.GetZero()

	if y.Sign() == 0 || x.Sign() == 0 {
		stack.Push(res)
	} else {
		if x.Sign() != y.Sign() {
			res.Div(x.Abs(x), y.Abs(y))
			res.Neg(res)
		} else {
			res.Div(x.Abs(x), y.Abs(y))
		}
		stack.Push(math.U256(res))
	}
	interpreter.evm.IntPool.Put(x, y)
	return nil, nil
}

func opMod(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Pop()
	if y.Sign() == 0 {
		stack.Push(x.SetUint64(0))
	} else {
		stack.Push(math.U256(x.Mod(x, y)))
	}
	interpreter.evm.IntPool.Put(y)
	return nil, nil
}

func opSmod(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := math.S256(stack.Pop()), math.S256(stack.Pop())
	res := interpreter.evm.IntPool.GetZero()

	if y.Sign() == 0 {
		stack.Push(res)
	} else {
		if x.Sign() < 0 {
			res.Mod(x.Abs(x), y.Abs(y))
			res.Neg(res)
		} else {
			res.Mod(x.Abs(x), y.Abs(y))
		}
		stack.Push(math.U256(res))
	}
	interpreter.evm.IntPool.Put(x, y)
	return nil, nil
}

func opExp(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	base, exponent := stack.Pop(), stack.Pop()
	if base.Cmp(big2) == 0 && exponent.Cmp(big256) == -1 {
		exp := exponent.Int64()
		stack.Push(interpreter.evm.IntPool.Get().Set(power2[exp]))
	} else {
		stack.Push(math.Exp(base, exponent))
	}

	interpreter.evm.IntPool.Put(base, exponent)

	return nil, nil
}

func opSignExtend(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	back := stack.Pop()
	if back.Cmp(big.NewInt(31)) < 0 {
		bit := uint(back.Uint64()*8 + 7)
		num := stack.Pop()
		mask := back.Lsh(common.Big1, bit)
		mask.Sub(mask, common.Big1)
		if num.Bit(int(bit)) > 0 {
			num.Or(num, mask.Not(mask))
		} else {
			num.And(num, mask)
		}

		stack.Push(math.U256(num))
	}

	interpreter.evm.IntPool.Put(back)
	return nil, nil
}

func opNot(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x := stack.Peek()
	math.U256(x.Not(x))
	return nil, nil
}

func opLt(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()
	if x.Cmp(y) < 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opGt(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()
	if x.Cmp(y) > 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opSlt(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()

	xSign := x.Cmp(tt255)
	ySign := y.Cmp(tt255)

	switch {
	case xSign >= 0 && ySign < 0:
		y.SetUint64(1)

	case xSign < 0 && ySign >= 0:
		y.SetUint64(0)

	default:
		if x.Cmp(y) < 0 {
			y.SetUint64(1)
		} else {
			y.SetUint64(0)
		}
	}
	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opSgt(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()

	xSign := x.Cmp(tt255)
	ySign := y.Cmp(tt255)

	switch {
	case xSign >= 0 && ySign < 0:
		y.SetUint64(0)

	case xSign < 0 && ySign >= 0:
		y.SetUint64(1)

	default:
		if x.Cmp(y) > 0 {
			y.SetUint64(1)
		} else {
			y.SetUint64(0)
		}
	}
	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opEq(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()
	if x.Cmp(y) == 0 {
		y.SetUint64(1)
	} else {
		y.SetUint64(0)
	}
	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opIszero(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x := stack.Peek()
	if x.Sign() > 0 {
		x.SetUint64(0)
	} else {
		x.SetUint64(1)
	}
	return nil, nil
}

func opAnd(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Pop()
	stack.Push(x.And(x, y))

	interpreter.evm.IntPool.Put(y)
	return nil, nil
}

func opOr(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()
	y.Or(x, y)

	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opXor(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y := stack.Pop(), stack.Peek()
	y.Xor(x, y)

	interpreter.evm.IntPool.Put(x)
	return nil, nil
}

func opByte(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	th, val := stack.Pop(), stack.Peek()
	if th.Cmp(common.Big32) < 0 {
		b := math.Byte(val, 32, int(th.Int64()))
		val.SetUint64(uint64(b))
	} else {
		val.SetUint64(0)
	}
	interpreter.evm.IntPool.Put(th)
	return nil, nil
}

func opAddmod(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y, z := stack.Pop(), stack.Pop(), stack.Pop()
	if z.Cmp(bigZero) > 0 {
		x.Add(x, y)
		x.Mod(x, z)
		stack.Push(math.U256(x))
	} else {
		stack.Push(x.SetUint64(0))
	}
	interpreter.evm.IntPool.Put(y, z)
	return nil, nil
}

func opMulmod(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	x, y, z := stack.Pop(), stack.Pop(), stack.Pop()
	if z.Cmp(bigZero) > 0 {
		x.Mul(x, y)
		x.Mod(x, z)
		stack.Push(math.U256(x))
	} else {
		stack.Push(x.SetUint64(0))
	}
	interpreter.evm.IntPool.Put(y, z)
	return nil, nil
}

// opSHL implements Shift Left
// The SHL instruction (shift left) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the left by arg1 number of bits.
func opSHL(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := math.U256(stack.Pop()), math.U256(stack.Peek())
	defer interpreter.evm.IntPool.Put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		value.SetUint64(0)
		return nil, nil
	}
	n := uint(shift.Uint64())
	math.U256(value.Lsh(value, n))

	return nil, nil
}

// opSHR implements Logical Shift Right
// The SHR instruction (logical shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with zero fill.
func opSHR(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	// Note, second operand is left in the stack; accumulate result into it, and no need to push it afterwards
	shift, value := math.U256(stack.Pop()), math.U256(stack.Peek())
	defer interpreter.evm.IntPool.Put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		value.SetUint64(0)
		return nil, nil
	}
	n := uint(shift.Uint64())
	math.U256(value.Rsh(value, n))

	return nil, nil
}

// opSAR implements Arithmetic Shift Right
// The SAR instruction (arithmetic shift right) pops 2 values from the stack, first arg1 and then arg2,
// and pushes on the stack arg2 shifted to the right by arg1 number of bits with sign extension.
func opSAR(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	// Note, S256 returns (potentially) a new bigint, so we're popping, not peeking this one
	shift, value := math.U256(stack.Pop()), math.S256(stack.Pop())
	defer interpreter.evm.IntPool.Put(shift) // First operand back into the pool

	if shift.Cmp(common.Big256) >= 0 {
		if value.Sign() >= 0 {
			value.SetUint64(0)
		} else {
			value.SetInt64(-1)
		}
		stack.Push(math.U256(value))
		return nil, nil
	}
	n := uint(shift.Uint64())
	value.Rsh(value, n)
	stack.Push(math.U256(value))

	return nil, nil
}

func opSha3(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	offset, size := stack.Pop(), stack.Pop()
	data := memory.Get(offset.Int64(), size.Int64())

	if interpreter.hasher == nil {
		interpreter.hasher = sha3.NewLegacyKeccak256().(keccakState)
	} else {
		interpreter.hasher.Reset()
	}
	interpreter.hasher.Write(data)
	interpreter.hasher.Read(interpreter.hasherBuf[:])

	evm := interpreter.evm
	if evm.vmConfig.EnablePreimageRecording {
		evm.StateDB.AddPreimage(interpreter.hasherBuf, data)
	}
	stack.Push(interpreter.evm.IntPool.Get().SetBytes(interpreter.hasherBuf[:]))

	interpreter.evm.IntPool.Put(offset, size)
	return nil, nil
}

func opRand(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	evm := interpreter.evm

	nonce := evm.StateDB.GetNonce(evm.Origin)
	binaryOriginNonce := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(binaryOriginNonce, nonce)

	binaryUsedIndex := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(binaryUsedIndex, evm.RandCallIndex)

	evm.RandCallIndex++

	hash := crypto.Keccak256(
		evm.Randomness,
		evm.Origin.Bytes(),
		binaryOriginNonce,
		binaryUsedIndex)

	stack.Push(interpreter.evm.IntPool.Get().SetBytes(hash))
	return nil, nil
}

func opAddress(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(contract.Address().Big())
	return nil, nil
}

func opBalance(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	slot := stack.Peek()
	slot.Set(interpreter.evm.StateDB.GetBalance(common.BigToAddress(slot)))
	return nil, nil
}

func opOrigin(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.Origin.Big())
	return nil, nil
}

func opCaller(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(contract.Caller().Big())
	return nil, nil
}

func opCallValue(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.IntPool.Get().Set(contract.Value))
	return nil, nil
}

func opCallDataLoad(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.IntPool.Get().SetBytes(vm.GetDataBig(contract.Input, stack.Pop(), big32)))
	return nil, nil
}

func opCallDataSize(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.IntPool.Get().SetInt64(int64(len(contract.Input))))
	return nil, nil
}

func opCallDataCopy(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	var (
		memOffset  = stack.Pop()
		dataOffset = stack.Pop()
		length     = stack.Pop()
	)
	memory.Set(memOffset.Uint64(), length.Uint64(), vm.GetDataBig(contract.Input, dataOffset, length))

	interpreter.evm.IntPool.Put(memOffset, dataOffset, length)
	return nil, nil
}

func opReturnDataSize(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.IntPool.Get().SetUint64(uint64(len(interpreter.returnData))))
	return nil, nil
}

func opReturnDataCopy(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	var (
		memOffset  = stack.Pop()
		dataOffset = stack.Pop()
		length     = stack.Pop()
		end        = interpreter.evm.IntPool.Get().Add(dataOffset, length)
	)
	defer interpreter.evm.IntPool.Put(memOffset, dataOffset, length, end)

	if end.BitLen() > 64 || uint64(len(interpreter.returnData)) < end.Uint64() {
		return nil, errReturnDataOutOfBounds
	}
	memory.Set(memOffset.Uint64(), length.Uint64(), interpreter.returnData[dataOffset.Uint64():end.Uint64()])

	return nil, nil
}

func opExtCodeSize(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	slot := stack.Peek()
	slot.SetUint64(uint64(interpreter.evm.StateDB.GetCodeSize(common.BigToAddress(slot))))

	return nil, nil
}

func opCodeSize(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	l := interpreter.evm.IntPool.Get().SetInt64(int64(len(contract.Code)))
	stack.Push(l)

	return nil, nil
}

func opCodeCopy(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	var (
		memOffset  = stack.Pop()
		codeOffset = stack.Pop()
		length     = stack.Pop()
	)
	codeCopy := vm.GetDataBig(contract.Code, codeOffset, length)
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	interpreter.evm.IntPool.Put(memOffset, codeOffset, length)
	return nil, nil
}

func opExtCodeCopy(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	var (
		addr       = common.BigToAddress(stack.Pop())
		memOffset  = stack.Pop()
		codeOffset = stack.Pop()
		length     = stack.Pop()
	)
	codeCopy := vm.GetDataBig(interpreter.evm.StateDB.GetCode(addr), codeOffset, length)
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	interpreter.evm.IntPool.Put(memOffset, codeOffset, length)
	return nil, nil
}

// opExtCodeHash returns the code hash of a specified account.
// There are several cases when the function is called, while we can relay everything
// to `state.GetCodeHash` function to ensure the correctness.
//   (1) Caller tries to get the code hash of a normal contract account, state
// should return the relative code hash and set it as the result.
//
//   (2) Caller tries to get the code hash of a non-existent account, state should
// return common.Hash{} and zero will be set as the result.
//
//   (3) Caller tries to get the code hash for an account without contract code,
// state should return emptyCodeHash(0xc5d246...) as the result.
//
//   (4) Caller tries to get the code hash of a precompiled account, the result
// should be zero or emptyCodeHash.
//
// It is worth noting that in order to avoid unnecessary create and clean,
// all precompile accounts on mainnet have been transferred 1 wei, so the return
// here should be emptyCodeHash.
// If the precompile account is not transferred any amount on a private or
// customized chain, the return value will be zero.
//
//   (5) Caller tries to get the code hash for an account which is marked as suicided
// in the current transaction, the code hash of this account should be returned.
//
//   (6) Caller tries to get the code hash for an account which is marked as deleted,
// this account should be regarded as a non-existent account and zero should be returned.
func opExtCodeHash(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	slot := stack.Peek()
	address := common.BigToAddress(slot)
	if interpreter.evm.StateDB.Empty(address) {
		slot.SetUint64(0)
	} else {
		slot.SetBytes(interpreter.evm.StateDB.GetCodeHash(address).Bytes())
	}
	return nil, nil
}

func opGasprice(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.IntPool.Get().Set(interpreter.evm.GasPrice))
	return nil, nil
}

func opBlockhash(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	num := stack.Pop()

	n := interpreter.evm.IntPool.Get().Sub(interpreter.evm.BlockNumber, common.Big257)
	if num.Cmp(n) > 0 && num.Cmp(interpreter.evm.BlockNumber) < 0 {
		stack.Push(interpreter.evm.GetHash(num.Uint64()).Big())
	} else {
		stack.Push(interpreter.evm.IntPool.GetZero())
	}
	interpreter.evm.IntPool.Put(num, n)
	return nil, nil
}

func opCoinbase(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.Coinbase.Big())
	return nil, nil
}

func opTimestamp(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(math.U256(interpreter.evm.IntPool.Get().Set(interpreter.evm.Time)))
	return nil, nil
}

func opNumber(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(math.U256(interpreter.evm.IntPool.Get().Set(interpreter.evm.BlockNumber)))
	return nil, nil
}

func opDifficulty(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(math.U256(interpreter.evm.IntPool.Get().Set(interpreter.evm.Difficulty)))
	return nil, nil
}

func opGasLimit(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(math.U256(interpreter.evm.IntPool.Get().SetUint64(interpreter.evm.GasLimit)))
	return nil, nil
}

func opPop(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	interpreter.evm.IntPool.Put(stack.Pop())
	return nil, nil
}

func opMload(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	offset := stack.Pop()
	val := interpreter.evm.IntPool.Get().SetBytes(memory.Get(offset.Int64(), 32))
	stack.Push(val)

	interpreter.evm.IntPool.Put(offset)
	return nil, nil
}

func opMstore(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	// pop value of the stack
	mStart, val := stack.Pop(), stack.Pop()
	memory.Set32(mStart.Uint64(), val)

	interpreter.evm.IntPool.Put(mStart, val)
	return nil, nil
}

func opMstore8(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	off, val := stack.Pop().Int64(), stack.Pop().Int64()
	memory.Store[off] = byte(val & 0xff)

	return nil, nil
}

func opSload(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	loc := stack.Peek()
	val := interpreter.evm.StateDB.GetState(contract.Address(), common.BigToHash(loc))
	loc.SetBytes(val.Bytes())
	return nil, nil
}

func opSstore(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	loc := common.BigToHash(stack.Pop())
	val := stack.Pop()
	interpreter.evm.StateDB.SetState(contract.Address(), loc, common.BigToHash(val))

	interpreter.evm.IntPool.Put(val)
	return nil, nil
}

func opJump(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	pos := stack.Pop()
	if !validJumpdest(pos, contract) {
		nop := OpCode(contract.GetByte(pos.Uint64()))
		return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
	}
	*pc = pos.Uint64()

	interpreter.evm.IntPool.Put(pos)
	return nil, nil
}

func opJumpi(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	pos, cond := stack.Pop(), stack.Pop()
	if cond.Sign() != 0 {
		if !validJumpdest(pos, contract) {
			nop := OpCode(contract.GetByte(pos.Uint64()))
			return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
		}
		*pc = pos.Uint64()
	} else {
		*pc++
	}

	interpreter.evm.IntPool.Put(pos, cond)
	return nil, nil
}

func opJumpdest(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	return nil, nil
}

func opPc(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.IntPool.Get().SetUint64(*pc))
	return nil, nil
}

func opMsize(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.IntPool.Get().SetInt64(int64(memory.Len())))
	return nil, nil
}

func opGas(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	stack.Push(interpreter.evm.IntPool.Get().SetUint64(contract.Gas))
	return nil, nil
}

func opCreate(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	var (
		value        = stack.Pop()
		offset, size = stack.Pop(), stack.Pop()
		input        = memory.Get(offset.Int64(), size.Int64())
		gas          = contract.Gas
	)
	// size.Add(size, big1)
	if interpreter.evm.ChainConfig().IsEIP150(interpreter.evm.BlockNumber) {
		gas -= gas / 64
	}
	contract.UseGas(gas)
	res, addr, returnGas, suberr := vm.Create(contract, input, gas, value, interpreter.evm.ExecPack)
	// Push item on the stack based on the returned error. If the ruleset is
	// homestead we must check for CodeStoreOutOfGasError (homestead only
	// rule) and treat as an error, if the ruleset is frontier we must
	// ignore this error and pretend the operation was successful.
	if interpreter.evm.ChainConfig().IsHomestead(interpreter.evm.BlockNumber) && suberr == vm.ErrCodeStoreOutOfGas {
		stack.Push(interpreter.evm.IntPool.GetZero())
	} else if suberr != nil && suberr != vm.ErrCodeStoreOutOfGas {
		stack.Push(interpreter.evm.IntPool.GetZero())
	} else {
		stack.Push(addr.Big())
	}
	contract.Gas += returnGas
	interpreter.evm.IntPool.Put(value, offset, size)

	if suberr == errExecutionReverted {
		return res, nil
	}
	return nil, nil
}

func opCreate2(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	var (
		endowment    = stack.Pop()
		offset, size = stack.Pop(), stack.Pop()
		salt         = stack.Pop()
		input        = memory.Get(offset.Int64(), size.Int64())
		gas          = contract.Gas
	)

	// Apply EIP150
	gas -= gas / 64
	contract.UseGas(gas)
	res, addr, returnGas, suberr := vm.Create2(contract, input, gas, endowment, salt, interpreter.evm.ExecPack)
	// Push item on the stack based on the returned error.
	if suberr != nil {
		stack.Push(interpreter.evm.IntPool.GetZero())
	} else {
		stack.Push(addr.Big())
	}
	contract.Gas += returnGas
	interpreter.evm.IntPool.Put(endowment, offset, size, salt)

	if suberr == errExecutionReverted {
		return res, nil
	}
	return nil, nil
}

func opCall(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	// Pop gas. The actual gas in interpreter.evm.CallGasTemp.
	interpreter.evm.IntPool.Put(stack.Pop())
	gas := interpreter.evm.CallGasTemp
	// Pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop()
	toAddr := common.BigToAddress(addr)
	value = math.U256(value)
	// Get the arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	ret, returnGas, err := vm.Call(contract, toAddr, args, gas, value, interpreter.evm.ExecPack)
	if err != nil {
		stack.Push(interpreter.evm.IntPool.GetZero())
	} else {
		stack.Push(interpreter.evm.IntPool.Get().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.Gas += returnGas

	interpreter.evm.IntPool.Put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opCallCode(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.CallGasTemp.
	interpreter.evm.IntPool.Put(stack.Pop())
	gas := interpreter.evm.CallGasTemp
	// Pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop()
	toAddr := common.BigToAddress(addr)
	value = math.U256(value)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	ret, returnGas, err := vm.CallCode(contract, toAddr, args, gas, value, interpreter.evm.ExecPack)
	if err != nil {
		stack.Push(interpreter.evm.IntPool.GetZero())
	} else {
		stack.Push(interpreter.evm.IntPool.Get().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.Gas += returnGas

	interpreter.evm.IntPool.Put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opDelegateCall(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.CallGasTemp.
	interpreter.evm.IntPool.Put(stack.Pop())
	gas := interpreter.evm.CallGasTemp
	// Pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop()
	toAddr := common.BigToAddress(addr)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := vm.DelegateCall(contract, toAddr, args, gas, interpreter.evm.ExecPack)
	if err != nil {
		stack.Push(interpreter.evm.IntPool.GetZero())
	} else {
		stack.Push(interpreter.evm.IntPool.Get().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.Gas += returnGas

	interpreter.evm.IntPool.Put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opStaticCall(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	// Pop gas. The actual gas is in interpreter.evm.CallGasTemp.
	interpreter.evm.IntPool.Put(stack.Pop())
	gas := interpreter.evm.CallGasTemp
	// Pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop(), stack.Pop()
	toAddr := common.BigToAddress(addr)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := vm.StaticCall(contract, toAddr, args, gas, interpreter.evm.ExecPack)
	if err != nil {
		stack.Push(interpreter.evm.IntPool.GetZero())
	} else {
		stack.Push(interpreter.evm.IntPool.Get().SetUint64(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.Gas += returnGas

	interpreter.evm.IntPool.Put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opReturn(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	offset, size := stack.Pop(), stack.Pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())

	interpreter.evm.IntPool.Put(offset, size)
	return ret, nil
}

func opRevert(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	offset, size := stack.Pop(), stack.Pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())

	interpreter.evm.IntPool.Put(offset, size)
	return ret, nil
}

func opStop(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	return nil, nil
}

func opSuicide(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
	balance := interpreter.evm.StateDB.GetBalance(contract.Address())
	interpreter.evm.StateDB.AddBalance(common.BigToAddress(stack.Pop()), balance)

	interpreter.evm.StateDB.Suicide(contract.Address())
	return nil, nil
}

// following functions are used by the instruction jump  table

// make log instruction function
func makeLog(size int) executionFunc {
	return func(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
		topics := make([]common.Hash, size)
		mStart, mSize := stack.Pop(), stack.Pop()
		for i := 0; i < size; i++ {
			topics[i] = common.BigToHash(stack.Pop())
		}

		d := memory.Get(mStart.Int64(), mSize.Int64())
		interpreter.evm.StateDB.AddLog(&types.Log{
			Address: contract.Address(),
			Topics:  topics,
			Data:    d,
			// This is a non-consensus field, but assigned here because
			// core/state doesn't know the current block number.
			BlockNumber: interpreter.evm.BlockNumber.Uint64(),
		})

		interpreter.evm.IntPool.Put(mStart, mSize)
		return nil, nil
	}
}

// make push instruction function
func makePush(size uint64, pushByteSize int) executionFunc {
	return func(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
		codeLen := len(contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			startMin = int(*pc + 1)
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		integer := interpreter.evm.IntPool.Get()
		stack.Push(integer.SetBytes(common.RightPadBytes(contract.Code[startMin:endMin], pushByteSize)))

		*pc += size
		return nil, nil
	}
}

// make dup instruction function
func makeDup(size int64) executionFunc {
	return func(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
		stack.Dup(interpreter.evm.IntPool, int(size))
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	// switch n + 1 otherwise n would be swapped with n
	size++
	return func(pc *uint64, interpreter *EVMInterpreter, contract *vm.Contract, memory *vm.Memory, stack *vm.Stack) ([]byte, error) {
		stack.Swap(int(size))
		return nil, nil
	}
}
func validJumpdest(dest *big.Int, c *vm.Contract) bool {
	udest := dest.Uint64()
	// PC cannot go beyond len(code) and certainly can't be bigger than 63bits.
	// Don't bother checking for JUMPDEST in that case.
	if dest.BitLen() >= 63 || udest >= uint64(len(c.Code)) {
		return false
	}
	// Only JUMPDESTs allowed for destinations
	if OpCode(c.Code[udest]) != JUMPDEST {
		return false
	}
	// Do we have a contract hash already?
	if c.CodeHash != (common.Hash{}) {
		// Does parent context have the analysis?
		analysis, exist := c.Jumpdests[c.CodeHash]
		if !exist {
			// Do the analysis and save in parent context
			// We do not need to store it in c.analysis
			analysis = codeBitmap(c.Code)
			c.Jumpdests[c.CodeHash] = analysis
		}
		return analysis.CodeSegment(udest)
	}
	// We don't have the code hash, most likely a piece of initcode not already
	// in state trie. In that case, we do an analysis, and save it locally, so
	// we don't have to recalculate it for every JUMP instruction in the execution
	// However, we don't save it within the parent context
	if c.Analysis == nil {
		c.Analysis = codeBitmap(c.Code)
	}
	return c.Analysis.CodeSegment(udest)
}
