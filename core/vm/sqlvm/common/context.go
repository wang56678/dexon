package common

import "github.com/dexon-foundation/dexon/core/vm"

// Context holds SQL VM required params.
type Context struct {
	vm.Context

	Storage  Storage
	Contract *vm.Contract
}
