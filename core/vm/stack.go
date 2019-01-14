// Copyright 2014 The go-ethereum Authors
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
	"fmt"
	"math/big"
)

// Stack is an object for basic stack operations. Items popped to the stack are
// expected to be changed and modified. stack does not take care of adding newly
// initialised objects.
type Stack struct {
	Data []*big.Int
}

func (st *Stack) Push(d *big.Int) {
	// NOTE push limit (1024) is checked in baseCheck
	//stackItem := new(big.Int).Set(d)
	//st.Data = append(st.Data, stackItem)
	st.Data = append(st.Data, d)
}
func (st *Stack) PushN(ds ...*big.Int) {
	st.Data = append(st.Data, ds...)
}

func (st *Stack) Pop() (ret *big.Int) {
	ret = st.Data[len(st.Data)-1]
	st.Data = st.Data[:len(st.Data)-1]
	return
}

func (st *Stack) Len() int {
	return len(st.Data)
}

func (st *Stack) Swap(n int) {
	st.Data[st.Len()-n], st.Data[st.Len()-1] = st.Data[st.Len()-1], st.Data[st.Len()-n]
}

func (st *Stack) Dup(pool *IntPool, n int) {
	st.Push(pool.Get().Set(st.Data[st.Len()-n]))
}

func (st *Stack) Peek() *big.Int {
	return st.Data[st.Len()-1]
}

// Back returns the n'th item in stack
func (st *Stack) Back(n int) *big.Int {
	return st.Data[st.Len()-n-1]
}

func (st *Stack) Require(n int) error {
	if st.Len() < n {
		return fmt.Errorf("stack underflow (%d <=> %d)", len(st.Data), n)
	}
	return nil
}

// Print dumps the content of the stack
func (st *Stack) Print() {
	fmt.Println("### stack ###")
	if len(st.Data) > 0 {
		for i, val := range st.Data {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
