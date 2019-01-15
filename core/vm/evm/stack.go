package evm

import (
	"math/big"
	"sync"

	"github.com/dexon-foundation/dexon/core/vm"
)

var stackPool = sync.Pool{
	New: func() interface{} {
		return &vm.Stack{Data: make([]*big.Int, 0, 1024)}
	},
}

func NewStack() *vm.Stack {
	stack := stackPool.Get().(*vm.Stack)
	stack.Data = stack.Data[:0]
	return stack
}

func Recyclestack(stack *vm.Stack) {
	stackPool.Put(stack)
}
