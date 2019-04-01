package runtime

import (
	"fmt"
	"strings"

	"github.com/dexon-foundation/decimal"

	dexCommon "github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

var tupleJoin = "|"

// OpFunction type
// data could be fields Fields, pattern []byte, order Orders
type OpFunction func(ctx *common.Context, ops []*Operand, registers []*Operand, output int) error

// Instruction represents single instruction with essential information
// collection.
type Instruction struct {
	Op       OpCode
	Input    []*Operand
	Output   int
	Position uint32 // ast tree position
}

// Raw with embedded big.Int value or byte slice which represents the real value
// of basic operand unit.
type Raw struct {
	Value decimal.Decimal
	Bytes []byte
}

func (r *Raw) String() string {
	return fmt.Sprintf("Value: %v, Bytes: %v", r.Value, r.Bytes)
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
	RegisterIndex uint
}

func (o *Operand) toUint64() (result []uint64, err error) {
	result = make([]uint64, len(o.Data))
	for i, tuple := range o.Data {
		result[i], err = ast.DecimalToUint64(tuple[0].Value)
		if err != nil {
			return
		}
	}
	return
}

func (o *Operand) toUint8() ([]uint8, error) {
	result := make([]uint8, len(o.Data))
	for i, tuple := range o.Data {
		u, err := ast.DecimalToUint64(tuple[0].Value)
		if err != nil {
			return nil, err
		}
		result[i] = uint8(u)
	}
	return result, nil
}

func opLoad(ctx *common.Context, input []*Operand, registers []*Operand, output int) error {
	tableIdx := input[0].Data[0][0].Value.IntPart()
	if tableIdx >= int64(len(ctx.Storage.Schema)) {
		return errors.ErrorCodeIndexOutOfRange
	}
	table := ctx.Storage.Schema[tableIdx]

	ids, err := input[1].toUint64()
	if err != nil {
		return err
	}
	fields, err := input[2].toUint8()
	if err != nil {
		return err
	}
	op := Operand{
		IsImmediate:   false,
		Data:          make([]Tuple, len(ids)),
		RegisterIndex: 0,
	}
	for i := range op.Data {
		op.Data[i] = make([]*Raw, len(fields))
	}
	op.Meta, err = table.GetFieldType(fields)
	if err != nil {
		return err
	}
	for i, id := range ids {
		slotDataCache := make(map[dexCommon.Hash]dexCommon.Hash)
		head := ctx.Storage.GetPrimaryKeyHash(table.Name, id)
		for j := range fields {
			col := table.Columns[int(fields[j])]
			byteOffset := col.ByteOffset
			slotOffset := col.SlotOffset
			dt := op.Meta[j]
			size := dt.Size()
			slot := ctx.Storage.ShiftHashUint64(head, uint64(slotOffset))
			slotData := getSlotData(ctx, slot, slotDataCache)
			bytes := slotData.Bytes()[byteOffset : byteOffset+size]
			op.Data[i][j], err = decode(ctx, dt, slot, bytes)
			if err != nil {
				return err
			}
		}
	}
	registers[output] = &op
	return nil
}

func getSlotData(ctx *common.Context, slot dexCommon.Hash,
	cache map[dexCommon.Hash]dexCommon.Hash) dexCommon.Hash {
	if d, exist := cache[slot]; exist {
		return d
	}
	cache[slot] = ctx.Storage.GetState(ctx.Contract.Address(), slot)
	return cache[slot]
}

// decode byte data to Raw format
func decode(ctx *common.Context, dt ast.DataType, slot dexCommon.Hash, bytes []byte) (*Raw, error) {
	rVal := &Raw{}
	major, _ := ast.DecomposeDataType(dt)
	switch major {
	case ast.DataTypeMajorDynamicBytes:
		rVal.Bytes = ctx.Storage.DecodeDByteBySlot(ctx.Contract.Address(), slot)
	case ast.DataTypeMajorFixedBytes, ast.DataTypeMajorAddress:
		rVal.Bytes = bytes
	case ast.DataTypeMajorBool, ast.DataTypeMajorInt, ast.DataTypeMajorUint:
		d, err := ast.DecimalDecode(dt, bytes)
		if err != nil {
			return nil, err
		}
		rVal.Value = d
	}
	if major.IsFixedRange() || major.IsUfixedRange() {
		d, err := ast.DecimalDecode(dt, bytes)
		if err != nil {
			return nil, err
		}
		rVal.Value = d
	}
	return rVal, nil
}
