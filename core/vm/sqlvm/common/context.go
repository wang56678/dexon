package common

import "github.com/dexon-foundation/dexon/core/vm"

// Option is collection of SQL options.
type Option struct {
	SafeMath bool
}

// Context holds SQL VM required params.
type Context struct {
	vm.Context

	Storage  *Storage
	Contract *vm.Contract

	Opt Option
}
