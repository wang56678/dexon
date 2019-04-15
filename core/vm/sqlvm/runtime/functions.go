package runtime

import (
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
)

type fn func(*common.Context, []*Operand, uint64) (*Operand, error)

var (
	fnTable = map[string]fn{}
)
