package runtime

import (
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

// Run is runtime entrypoint.
func Run(stateDB vm.StateDB, ins []Instruction, registers []*Operand) (ret []byte, err error) {
	for _, in := range ins {
		opFunc := jumpTable[in.Op]
		errCode := opFunc(&common.Context{}, in.Input, registers, in.Output)
		if errCode != nil {
			err = errors.Error{
				Position: in.Position,
				Code:     errCode.(errors.ErrorCode),
				Category: errors.ErrorCategoryRuntime,
			}
			return nil, err
		}
	}
	// TODO: ret = ABIEncode(ins[len(ins)-1].Output)
	return
}
