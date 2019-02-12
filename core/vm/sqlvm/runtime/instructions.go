package runtime

import (
	"fmt"
	"strings"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
)

var tupleJoin = "|"

// OpFunction type
// data could be fields Fields, pattern []byte, order Orders
type OpFunction func(ctx *common.Context, ops []*Operand, registers []*Operand, output int) error

// Instruction represents single instruction with essential information
// collection.
type Instruction struct {
	op     OpCode
	input  []*Operand
	output int
}

// Raw with embedded big.Int value or byte slice which represents the real value
// of basic operand unit.
type Raw struct {
	MajorType ast.DataTypeMajor
	MinorType ast.DataTypeMinor

	// for not bytes
	Value interface{}
}

func (r *Raw) String() string {
	return fmt.Sprintf(
		"MajorType: %v, MinorType: %v, Value: %v",
		r.MajorType, r.MinorType, r.Value)
}

// Tuple is collection of Raw.
type Tuple []*Raw

func (t Tuple) String() string {
	rawStr := []string{}
	for i := 0; i < len(t); i++ {
		rawStr = append(rawStr, t[i].String())
	}
	return strings.Join(rawStr, tupleJoin)
}

// Operand would be array-based value associated with meta to describe type of
// array element.
type Operand struct {
	IsImmediate   bool
	Meta          []ast.DataType
	Data          []Tuple
	RegisterIndex *int
}
