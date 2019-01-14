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

package vm

import (
	"math/big"
	"sync"
)

var checkVal = big.NewInt(-42)

const PoolLimit = 256

// IntPool is a Pool of big integers that
// can be reused for all big.Int operations.
type IntPool struct {
	Pool *Stack
}

func newIntPool() *IntPool {
	return &IntPool{Pool: &Stack{Data: make([]*big.Int, 0, 1024)}}
}

// get retrieves a big int from the Pool, allocating one if the Pool is empty.
// Note, the returned int's value is arbitrary and will not be zeroed!
func (p *IntPool) Get() *big.Int {
	if p.Pool.Len() > 0 {
		return p.Pool.Pop()
	}
	return new(big.Int)
}

// getZero retrieves a big int from the Pool, setting it to zero or allocating
// a new one if the Pool is empty.
func (p *IntPool) GetZero() *big.Int {
	if p.Pool.Len() > 0 {
		return p.Pool.Pop().SetUint64(0)
	}
	return new(big.Int)
}

// put returns an allocated big int to the Pool to be later reused by get calls.
// Note, the values as saved as is; neither put nor get zeroes the ints out!
func (p *IntPool) Put(is ...*big.Int) {
	if len(p.Pool.Data) > PoolLimit {
		return
	}
	for _, i := range is {
		// verifyPool is a build flag. Pool verification makes sure the integrity
		// of the integer Pool by comparing values to a default value.
		if VerifyPool {
			i.Set(checkVal)
		}
		p.Pool.Push(i)
	}
}

// The IntPool Pool's default capacity
const PoolDefaultCap = 25

// IntPoolPool manages a Pool of IntPools.
type IntPoolPool struct {
	Pools []*IntPool
	lock  sync.Mutex
}

var PoolOfIntPools = &IntPoolPool{
	Pools: make([]*IntPool, 0, PoolDefaultCap),
}

// get is looking for an available Pool to return.
func (ipp *IntPoolPool) Get() *IntPool {
	ipp.lock.Lock()
	defer ipp.lock.Unlock()

	if len(PoolOfIntPools.Pools) > 0 {
		ip := ipp.Pools[len(ipp.Pools)-1]
		ipp.Pools = ipp.Pools[:len(ipp.Pools)-1]
		return ip
	}
	return newIntPool()
}

// put a Pool that has been allocated with get.
func (ipp *IntPoolPool) Put(ip *IntPool) {
	ipp.lock.Lock()
	defer ipp.lock.Unlock()

	if len(ipp.Pools) < cap(ipp.Pools) {
		ipp.Pools = append(ipp.Pools, ip)
	}
}
