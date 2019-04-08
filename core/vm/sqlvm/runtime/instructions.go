package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/dexon-foundation/decimal"

	dexCommon "github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	se "github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

var (
	byteLikeP    = byte('%')
	byteLikeU    = byte('_')
	byteDot      = byte('.')
	bytesLikeReg = []byte{'.', '*', '?'}
	bytesStart   = []byte{'^'}
	bytesEnd     = []byte{'$'}

	tupleJoin = "|"
)

// OpFunction type
// data could be fields Fields, pattern []byte, order Orders
type OpFunction func(ctx *common.Context, ops, registers []*Operand, output uint) error

// Instruction represents single instruction with essential information
// collection.
type Instruction struct {
	Op       OpCode
	Input    []*Operand
	Output   uint
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

func (r *Raw) isTrue() bool  { return dec.IsTrue(r.Value) }
func (r *Raw) isFalse() bool { return dec.IsFalse(r.Value) }

// Tuple is collection of Raw.
type Tuple []*Raw

func (t Tuple) String() string {
	rawStr := []string{}
	for i := 0; i < len(t); i++ {
		rawStr = append(rawStr, t[i].String())
	}
	return fmt.Sprintf("\n%v\n", strings.Join(rawStr, tupleJoin))
}

// Operand would be array-based value associated with meta to describe type of
// array element.
type Operand struct {
	IsImmediate   bool
	Meta          []ast.DataType
	Data          []Tuple
	RegisterIndex uint
}

func (op *Operand) toUint64() (result []uint64, err error) {
	result = make([]uint64, len(op.Data))
	for i, tuple := range op.Data {
		result[i], err = ast.DecimalToUint64(tuple[0].Value)
		if err != nil {
			return
		}
	}
	return
}

func (op *Operand) toUint8() ([]uint8, error) {
	result := make([]uint8, len(op.Data))
	for i, tuple := range op.Data {
		u, err := ast.DecimalToUint64(tuple[0].Value)
		if err != nil {
			return nil, err
		}
		result[i] = uint8(u)
	}
	return result, nil
}

func opLoad(ctx *common.Context, input []*Operand, registers []*Operand, output uint) error {
	tableIdx := input[0].Data[0][0].Value.IntPart()
	if tableIdx >= int64(len(ctx.Storage.Schema)) {
		return se.ErrorCodeIndexOutOfRange
	}
	table := ctx.Storage.Schema[tableIdx]
	tableRef := schema.TableRef(tableIdx)

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
		head := ctx.Storage.GetRowPathHash(tableRef, id)
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

func (op *Operand) clone(metaOnly bool) (op2 *Operand) {
	op2 = &Operand{
		Meta: op.cloneMeta(),
		// skip RegisterIndex since is only used when loading
		// skip IsImmediate flag since is only set by codegen
	}

	if metaOnly {
		return
	}

	op2.Data = make([]Tuple, len(op.Data))
	for i := 0; i < len(op.Data); i++ {
		op2.Data[i] = append(Tuple{}, op.Data[i]...)
	}
	return
}

func (op *Operand) cloneMeta() (meta []ast.DataType) {
	meta = make([]ast.DataType, len(op.Meta))
	copy(meta, op.Meta)
	return
}

// Equal compares underlying data level-by-level.
func (op *Operand) Equal(op2 *Operand) (equal bool) {
	if op2 == nil {
		return
	}

	equal = op.IsImmediate == op2.IsImmediate
	if !equal {
		return
	}

	equal = metaAllEq(op, op2)
	if !equal {
		return
	}

	equal = len(op.Data) == len(op2.Data)
	if !equal {
		return
	}

	for i := 0; i < len(op.Data); i++ {
		equal = op.Data[i].Equal(op2.Data[i], op.Meta)
		if !equal {
			return
		}
	}
	return
}

// Equal compares tuple values one by one.
func (t Tuple) Equal(t2 Tuple, meta []ast.DataType) (equal bool) {
	equal = len(t) == len(t2)
	if !equal {
		return
	}

	for i := 0; i < len(t); i++ {
		equal = t[i].Equal(t2[i], meta[i])
		if !equal {
			return
		}
	}
	return
}

// Equal compares raw by type.
func (r *Raw) Equal(r2 *Raw, dType ast.DataType) (equal bool) {
	major, _ := ast.DecomposeDataType(dType)
	switch major {
	case ast.DataTypeMajorDynamicBytes,
		ast.DataTypeMajorFixedBytes,
		ast.DataTypeMajorAddress:
		equal = bytes.Equal(r.Bytes, r2.Bytes)
	default:
		equal = r.Value.Cmp(r2.Value) == 0
	}
	return
}

var (
	rawFalse = &Raw{Value: dec.False}
	rawTrue  = &Raw{Value: dec.True}
)

func metaAllEq(op1, op2 *Operand) bool {
	if len(op1.Meta) != len(op2.Meta) {
		return false
	}

	for i := 0; i < len(op1.Meta); i++ {
		if op1.Meta[i] != op2.Meta[i] {
			return false
		}
	}
	return true
}

func metaAll(op *Operand, fn func(ast.DataType) bool) bool {
	for i := 0; i < len(op.Meta); i++ {
		if !fn(op.Meta[i]) {
			return false
		}
	}
	return true
}

func metaBool(dType ast.DataType) bool {
	dMajor, _ := ast.DecomposeDataType(dType)
	return dMajor == ast.DataTypeMajorBool
}

func metaAllBool(op *Operand) bool { return metaAll(op, metaBool) }

func metaArith(dType ast.DataType) bool {
	major, _ := ast.DecomposeDataType(dType)
	if major == ast.DataTypeMajorInt ||
		major == ast.DataTypeMajorUint {
		return true
	}
	return false
}

func metaAllArith(op *Operand) bool { return metaAll(op, metaArith) }

func findMaxDataLength(ops []*Operand) (l int, err error) {
	l = -1

	for i := 0; i < len(ops); i++ {
		if ops[i].IsImmediate {
			continue
		}

		if l == -1 {
			l = len(ops[i].Data)
			continue
		}

		if len(ops[i].Data) != l {
			err = se.ErrorCodeDataLengthNotMatch
			l = -1
			return
		}
	}

	if l == -1 {
		l = 1
	}
	return
}

func bool2Raw(b bool) (r *Raw) {
	if b {
		r = rawTrue
	} else {
		r = rawFalse
	}
	return
}

func value2ColIdx(v decimal.Decimal) (idx uint16) {
	if v.GreaterThan(dec.MaxUint16) {
		panic(errors.New("field index greater than uint16 max"))
	} else if v.LessThan(decimal.Zero) {
		panic(errors.New("field index less than 0"))
	}

	idx = uint16(v.IntPart())
	return
}

func metaDynBytes(dType ast.DataType) bool {
	dMajor, _ := ast.DecomposeDataType(dType)
	return dMajor == ast.DataTypeMajorDynamicBytes
}

func metaAllDynBytes(op *Operand) bool { return metaAll(op, metaDynBytes) }

func metaSignedNumeric(dType ast.DataType) bool {
	major, _ := ast.DecomposeDataType(dType)
	if major == ast.DataTypeMajorInt ||
		major == ast.DataTypeMajorFixed {
		return true
	}
	return false
}

func metaAllSignedNumeric(op *Operand) bool { return metaAll(op, metaSignedNumeric) }

func flowCheck(ctx *common.Context, v decimal.Decimal, dType ast.DataType) (err error) {
	if !ctx.Opt.SafeMath {
		return
	}

	min, max, ok := dType.GetMinMax()
	if !ok {
		err = se.ErrorCodeInvalidDataType
		return
	}

	if v.Cmp(max) > 0 {
		err = se.ErrorCodeOverflow
	} else if v.Cmp(min) < 0 {
		err = se.ErrorCodeUnderflow
	}
	return
}

func opAdd(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) || !metaAllArith(op1) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}

		raw, iErr := op1.Data[j].add(ctx, op2.Data[k], op1.Meta)
		if iErr != nil {
			err = iErr
			return
		}
		data[i] = raw
	}

	registers[output] = &Operand{Meta: op1.cloneMeta(), Data: data}
	return
}

func (t Tuple) add(ctx *common.Context, t2 Tuple, meta []ast.DataType) (t3 Tuple, err error) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		raw := t[i].add(t2[i])
		err = flowCheck(ctx, raw.Value, meta[i])
		if err != nil {
			return
		}
		t3[i] = raw
	}
	return
}

func (r *Raw) add(r2 *Raw) (r3 *Raw) {
	r3 = &Raw{Value: r.Value.Add(r2.Value)}
	return
}

func opMul(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) || !metaAllArith(op1) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}

		raw, iErr := op1.Data[j].mul(ctx, op2.Data[k], op1.Meta)
		if iErr != nil {
			err = iErr
			return
		}
		data[i] = raw
	}

	registers[output] = &Operand{Meta: op1.cloneMeta(), Data: data}
	return
}

func (t Tuple) mul(ctx *common.Context, t2 Tuple, meta []ast.DataType) (t3 Tuple, err error) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		raw := t[i].mul(t2[i])
		err = flowCheck(ctx, raw.Value, meta[i])
		if err != nil {
			return
		}

		t3[i] = raw
	}
	return
}

func (r *Raw) mul(r2 *Raw) (r3 *Raw) {
	r3 = &Raw{Value: r.Value.Mul(r2.Value)}
	return
}

func opSub(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) || !metaAllArith(op1) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}

		raw, iErr := op1.Data[j].sub(ctx, op2.Data[k], op1.Meta)
		if iErr != nil {
			err = iErr
			return
		}
		data[i] = raw
	}

	registers[output] = &Operand{Meta: op1.cloneMeta(), Data: data}
	return
}

func (t Tuple) sub(ctx *common.Context, t2 Tuple, meta []ast.DataType) (t3 Tuple, err error) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		raw := t[i].sub(t2[i])
		err = flowCheck(ctx, raw.Value, meta[i])
		if err != nil {
			return
		}

		t3[i] = raw
	}
	return
}

func (r *Raw) sub(r2 *Raw) (r3 *Raw) {
	r3 = &Raw{Value: r.Value.Sub(r2.Value)}
	return
}

func opDiv(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) || !metaAllArith(op1) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}
		raw, iErr := op1.Data[j].div(ctx, op2.Data[k], op1.Meta)
		if iErr != nil {
			err = iErr
			return
		}
		data[i] = raw
	}

	registers[output] = &Operand{Meta: op1.cloneMeta(), Data: data}
	return
}

func (t Tuple) div(ctx *common.Context, t2 Tuple, meta []ast.DataType) (t3 Tuple, err error) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		raw, iErr := t[i].div(t2[i])
		if iErr != nil {
			err = iErr
			return
		}

		iErr = flowCheck(ctx, raw.Value, meta[i])
		if iErr != nil {
			err = iErr
			return
		}

		t3[i] = raw
	}
	return
}

func (r *Raw) div(r2 *Raw) (r3 *Raw, err error) {
	if r2.Value.IsZero() {
		err = se.ErrorCodeDividedByZero
		return
	}

	q, _ := r.Value.QuoRem(r2.Value, 0)

	r3 = &Raw{Value: q}
	return
}

func opMod(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) || !metaAllArith(op1) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}
		raw, iErr := op1.Data[j].mod(ctx, op2.Data[k], op1.Meta)
		if iErr != nil {
			err = iErr
			return
		}
		data[i] = raw
	}

	registers[output] = &Operand{Meta: op1.cloneMeta(), Data: data}
	return
}

func (t Tuple) mod(ctx *common.Context, t2 Tuple, meta []ast.DataType) (t3 Tuple, err error) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		raw, iErr := t[i].mod(t2[i])
		if iErr != nil {
			err = iErr
			return

		}
		t3[i] = raw
	}
	return
}

func (r *Raw) mod(r2 *Raw) (r3 *Raw, err error) {
	if r2.Value.IsZero() {
		err = se.ErrorCodeDividedByZero
		return
	}

	_, qr := r.Value.QuoRem(r2.Value, 0)

	r3 = &Raw{Value: qr}
	return
}

func opLt(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}
		data[i] = op1.Data[j].lt(op2.Data[k])
	}

	boolType := ast.ComposeDataType(ast.DataTypeMajorBool, 0)
	meta := make([]ast.DataType, len(op1.Meta))
	for i := 0; i < len(op1.Meta); i++ {
		meta[i] = boolType
	}

	registers[output] = &Operand{Meta: meta, Data: data}
	return
}

func (t Tuple) lt(t2 Tuple) (t3 Tuple) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t3[i] = bool2Raw(t[i].lt(t2[i]))
	}
	return
}

func (r *Raw) lt(r2 *Raw) (lt bool) {
	if r.Bytes == nil {
		lt = r.Value.Cmp(r2.Value) < 0
	} else {
		lt = bytes.Compare(r.Bytes, r2.Bytes) < 0
	}
	return
}

func opGt(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}
		data[i] = op1.Data[j].gt(op2.Data[k])
	}

	boolType := ast.ComposeDataType(ast.DataTypeMajorBool, 0)
	meta := make([]ast.DataType, len(op1.Meta))
	for i := 0; i < len(op1.Meta); i++ {
		meta[i] = boolType
	}

	registers[output] = &Operand{Meta: meta, Data: data}
	return
}

func (t Tuple) gt(t2 Tuple) (t3 Tuple) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t3[i] = bool2Raw(t[i].gt(t2[i]))
	}
	return
}

func (r *Raw) gt(r2 *Raw) (gt bool) {
	if r.Bytes == nil {
		gt = r.Value.Cmp(r2.Value) > 0
	} else {
		gt = bytes.Compare(r.Bytes, r2.Bytes) > 0
	}
	return
}

func opEq(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}
		data[i] = op1.Data[j].eq(op2.Data[k], op1.Meta)
	}

	boolType := ast.ComposeDataType(ast.DataTypeMajorBool, 0)
	meta := make([]ast.DataType, len(op1.Meta))
	for i := 0; i < len(op1.Meta); i++ {
		meta[i] = boolType
	}

	registers[output] = &Operand{Meta: meta, Data: data}
	return
}

func (t Tuple) eq(t2 Tuple, meta []ast.DataType) (t3 Tuple) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t3[i] = bool2Raw(t[i].Equal(t2[i], meta[i]))
	}
	return
}

func opAnd(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	if !metaAllBool(op1) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}
		data[i] = op1.Data[j].and(op2.Data[k])
	}

	boolType := ast.ComposeDataType(ast.DataTypeMajorBool, 0)
	meta := make([]ast.DataType, l)
	for i := 0; i < l; i++ {
		meta[i] = boolType
	}

	registers[output] = &Operand{Meta: meta, Data: data}
	return
}

func (t Tuple) and(t2 Tuple) (t3 Tuple) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t3[i] = t[i].and(t2[i])
	}
	return
}

func (r *Raw) and(r2 *Raw) (r3 *Raw) {
	r3 = bool2Raw(r.isTrue() && r2.isTrue())
	return
}

func opOr(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	if !metaAllBool(op1) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	data := make([]Tuple, l)
	for i, j, k := 0, 0, 0; i < l; i, j, k = i+1, j+1, k+1 {
		if op1.IsImmediate {
			j = 0
		}
		if op2.IsImmediate {
			k = 0
		}
		data[i] = op1.Data[j].or(op2.Data[k])
	}

	boolType := ast.ComposeDataType(ast.DataTypeMajorBool, 0)
	meta := make([]ast.DataType, l)
	for i := 0; i < l; i++ {
		meta[i] = boolType
	}

	registers[output] = &Operand{Meta: meta, Data: data}
	return
}

func (t Tuple) or(t2 Tuple) (t3 Tuple) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t3[i] = t[i].or(t2[i])
	}
	return
}

func (r *Raw) or(r2 *Raw) (r3 *Raw) {
	r3 = bool2Raw(r.Value.Equal(dec.True) || r2.Value.Equal(dec.True))
	return
}

func opNot(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 1 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op := ops[0]

	if !metaAllBool(op) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	data := make([]Tuple, len(op.Data))
	for i := 0; i < len(op.Data); i++ {
		data[i] = op.Data[i].not()
	}

	registers[output] = &Operand{Meta: op.cloneMeta(), Data: data}
	return
}

func (t Tuple) not() (t2 Tuple) {
	t2 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t2[i] = t[i].not()
	}
	return
}

func (r *Raw) not() (r2 *Raw) {
	r2 = bool2Raw(r.Value.IsZero())
	return
}

func opUnion(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	if len(op1.Data) > len(op2.Data) {
		op1, op2 = op2, op1
	}

	tmpMap := make(map[string]struct{})
	for i := 0; i < len(op1.Data); i++ {
		tmpMap[op1.Data[i].String()] = struct{}{}
	}

	op3 := op1.clone(false)

	for i := 0; i < len(op2.Data); i++ {
		if _, ok := tmpMap[op2.Data[i].String()]; !ok {
			op3.Data = append(op3.Data, append(Tuple{}, op2.Data[i]...))
		}
	}

	orders := make([]sortOption, len(op3.Meta))
	for i := 0; i < len(orders); i++ {
		orders[i] = sortOption{Asc: true, Field: uint(i)}
	}

	sort.SliceStable(
		op3.Data,
		func(i, j int) bool { return op3.Data[i].less(op3.Data[j], orders) },
	)

	registers[output] = op3
	return
}

func opIntxn(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op1, op2 := ops[0], ops[1]

	if !metaAllEq(op1, op2) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	if len(op1.Data) > len(op2.Data) {
		op1, op2 = op2, op1
	}

	tmpMap := make(map[string]struct{})
	for i := 0; i < len(op1.Data); i++ {
		tmpMap[op1.Data[i].String()] = struct{}{}
	}

	op3 := &Operand{Meta: op1.cloneMeta(), Data: []Tuple{}}

	for i := 0; i < len(op2.Data); i++ {
		if _, ok := tmpMap[op2.Data[i].String()]; ok {
			op3.Data = append(op3.Data, append(Tuple{}, op2.Data[i]...))
		}
	}

	orders := make([]sortOption, len(op3.Meta))
	for i := 0; i < len(orders); i++ {
		orders[i] = sortOption{Asc: true, Field: uint(i)}
	}

	sort.SliceStable(
		op3.Data,
		func(i, j int) bool { return op3.Data[i].less(op3.Data[j], orders) },
	)

	registers[output] = op3
	return
}

func opLike(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 && len(ops) != 3 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op, pattern := ops[0], ops[1]

	var escape *Operand
	if len(ops) > 2 {
		escape = ops[2]
	}

	var cReg *regexp.Regexp

	matchWithI := pattern.IsImmediate && (escape == nil || escape.IsImmediate)
	if matchWithI {
		var escapeBytes []byte
		if escape != nil && len(escape.Data) > 0 {
			escapeBytes = escape.Data[0][0].Bytes
		}

		if len(escapeBytes) > 1 {
			err = se.ErrorCodeMultipleEscapeByte
			return
		}

		cReg, err = like2regexp(pattern.Data[0][0].Bytes, escapeBytes)
		if err != nil {
			return
		}

	}

	data := make([]Tuple, len(op.Data))
	if matchWithI {
		for i := 0; i < len(op.Data); i++ {
			raw, iErr := op.Data[i].like(cReg)
			if iErr != nil {
				err = iErr
				return
			}
			data[i] = raw
		}
	} else {
		var (
			pat []byte
			esc []byte
		)

		for i := 0; i < len(op.Data); i++ {
			if pattern.IsImmediate {
				pat = pattern.Data[0][0].Bytes
			} else {
				pat = pattern.Data[i][0].Bytes
			}

			if escape != nil {
				if escape.IsImmediate {
					esc = escape.Data[0][0].Bytes
				} else {
					esc = escape.Data[i][0].Bytes
				}
			} else {
				esc = []byte{}
			}

			if len(esc) > 1 {
				err = se.ErrorCodeMultipleEscapeByte
				return
			}

			reg, iErr := like2regexp(pat, esc)
			if iErr != nil {
				err = iErr
				return
			}

			raw, iErr := op.Data[i].like(reg)
			if iErr != nil {
				err = iErr
				return
			}

			data[i] = raw
		}
	}

	boolType := ast.ComposeDataType(ast.DataTypeMajorBool, 0)
	meta := make([]ast.DataType, len(op.Meta))
	for i := 0; i < len(meta); i++ {
		meta[i] = boolType
	}

	registers[output] = &Operand{Meta: meta, Data: data}
	return
}

// check parser/parser.go comment for string encoding
func encB(b []byte) []byte {
	encBuf := bytes.Buffer{}
	for _, c := range b {
		encBuf.WriteRune(rune(c))
	}
	return encBuf.Bytes()
}

func writeC2Buf(buf *bytes.Buffer, c byte) {
	if c < 0x80 {
		// quote for valid ascii
		buf.WriteString(regexp.QuoteMeta(string(c)))
	} else {
		buf.WriteRune(rune(c))
	}
}

func like2regexp(pattern []byte, escape []byte) (reg *regexp.Regexp, err error) {
	var (
		buf     = &bytes.Buffer{}
		escMode = len(escape) > 0
		isEsc   = false
		c       byte
	)

	for i := 0; i < len(pattern); i++ {
		c = pattern[i]

		if escMode && !isEsc && c == escape[0] {
			isEsc = true
			continue
		}

		if escMode && isEsc {
			isEsc = false
			writeC2Buf(buf, c)
			continue
		}

		switch c {
		case byteLikeP:
			buf.Write(bytesLikeReg)
		case byteLikeU:
			buf.WriteByte(byteDot)
		default:
			writeC2Buf(buf, c)
		}
	}

	if isEsc {
		err = se.ErrorCodePendingEscapeByte
		return
	}

	rPattern := buf.Bytes()

	if !bytes.HasPrefix(rPattern, bytesLikeReg) {
		rPattern = append(bytesStart, rPattern...)
	}

	if !bytes.HasSuffix(rPattern, bytesLikeReg) {
		rPattern = append(rPattern, bytesEnd...)
	}

	reg, err = regexp.Compile(string(rPattern))
	return
}

func (t Tuple) like(reg *regexp.Regexp) (t2 Tuple, err error) {
	t2 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t2[i] = t[i].like(reg)
	}
	return
}

func (r *Raw) like(reg *regexp.Regexp) (r2 *Raw) {
	r2 = bool2Raw(reg.Match(encB(r.Bytes)))
	return
}

func opZip(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) == 0 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}

	l, err := findMaxDataLength(ops)
	if err != nil {
		return
	}

	op3 := &Operand{Meta: make([]ast.DataType, 0), Data: make([]Tuple, l)}

	for i := 0; i < len(ops); i++ {
		op3.Meta = append(op3.Meta, ops[i].Meta...)
	}

	for i := 0; i < l; i++ {
		if ops[0].IsImmediate {
			op3.Data[i] = append(Tuple{}, ops[0].Data[0]...)
		} else {
			op3.Data[i] = append(Tuple{}, ops[0].Data[i]...)
		}

		for j := 1; j < len(ops); j++ {
			if ops[j].IsImmediate {
				op3.Data[i] = append(op3.Data[i], ops[j].Data[0]...)
			} else {
				op3.Data[i] = append(op3.Data[i], ops[j].Data[i]...)
			}
		}
	}

	registers[output] = op3
	return
}

func opField(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op, fields := ops[0], ops[1].Data[0]
	fLen := len(fields)

	var fieldIdx uint16
	meta, fieldIdxs := make([]ast.DataType, fLen), make([]uint16, fLen)
	for i := 0; i < fLen; i++ {
		fieldIdx = value2ColIdx(fields[i].Value)
		if len(op.Meta) <= int(fieldIdx) {
			err = se.ErrorCodeIndexOutOfRange
			return
		}
		meta[i], fieldIdxs[i] = op.Meta[fieldIdx], fieldIdx
	}

	data := make([]Tuple, len(op.Data))
	for i := 0; i < len(op.Data); i++ {
		tuple := make(Tuple, fLen)
		for j := 0; j < fLen; j++ {
			tuple[j] = op.Data[i][fieldIdxs[j]]
		}
		data[i] = tuple
	}

	registers[output] = &Operand{Meta: meta, Data: data}
	return
}

// in-place Op
func opPrune(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op, fields := ops[0], ops[1].Data[0]
	fLen := len(fields)

	var fieldIdx uint16
	fieldIdxs := make([]int, fLen)
	for i := 0; i < fLen; i++ {
		fieldIdx = value2ColIdx(fields[i].Value)
		if len(op.Meta) <= int(fieldIdx) {
			err = se.ErrorCodeIndexOutOfRange
			return
		}
		fieldIdxs[i] = int(fieldIdx)
	}

	op.Meta = pruneMeta(op.Meta, fieldIdxs)

	for i := 0; i < len(op.Data); i++ {
		op.Data[i] = pruneTuple(op.Data[i], fieldIdxs)
	}
	return
}

func pruneMeta(meta []ast.DataType, prune []int) []ast.DataType {
	for src, dst, pruneIdx := 0, 0, 0; src < len(meta); src++ {
		if pruneIdx < len(prune) && src == prune[pruneIdx] {
			pruneIdx++
			continue
		}
		meta[dst] = meta[src]
		dst++
	}
	return meta[:len(meta)-len(prune)]
}

func pruneTuple(t Tuple, prune []int) Tuple {
	for src, dst, pruneIdx := 0, 0, 0; src < len(t); src++ {
		if pruneIdx < len(prune) && src == prune[pruneIdx] {
			pruneIdx++
			continue
		}
		t[dst] = t[src]
		dst++
	}
	return t[:len(t)-len(prune)]
}

// in-place Op
func opCut(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op, slice := ops[0], ops[1].Data[0]

	maxL := uint16(len(op.Meta))
	start, end := value2ColIdx(slice[0].Value), maxL
	if len(slice) > 1 {
		end = value2ColIdx(slice[1].Value) + 1
	}

	if start > maxL || end > maxL || start > end {
		err = se.ErrorCodeIndexOutOfRange
		return
	}

	op.Meta = append(op.Meta[:start], op.Meta[end:]...)

	for i := 0; i < len(op.Data); i++ {
		op.Data[i] = append(op.Data[i][:start], op.Data[i][end:]...)
	}

	registers[output] = op
	return
}

// in-place Op
func opRange(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op, slice := ops[0], ops[1].Data[0]

	offset, err := ast.DecimalToUint64(slice[0].Value)
	if err != nil {
		return
	}

	if offset < uint64(len(op.Data)) {
		op.Data = op.Data[offset:]
	} else {
		op.Data = op.Data[:0]
	}

	if len(slice) > 1 && len(op.Data) > 0 {
		var limit uint64

		limit, err = ast.DecimalToUint64(slice[1].Value)
		if err != nil {
			return
		}

		if limit < uint64(len(op.Data)) {
			op.Data = op.Data[:limit]
		}
	}

	registers[output] = op
	return
}

// in-place Op
func opSort(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op := ops[0]

	// orders (ascending bool, field)
	orders := make([]sortOption, len(ops[1].Data))
	for i, o := range ops[1].Data {
		orders[i] = sortOption{
			Asc:   o[0].isTrue(),
			Field: uint(value2ColIdx(o[1].Value)),
		}
	}

	sort.SliceStable(
		op.Data,
		func(i, j int) bool { return op.Data[i].less(op.Data[j], orders) },
	)
	return
}

type sortOption struct {
	Asc   bool
	Field uint
}

func (t Tuple) less(t2 Tuple, orders []sortOption) bool {
	var r int

	for _, o := range orders {
		r = 0

		if o.Asc {
			r = t[o.Field].cmp(t2[o.Field])
		} else {
			r = t2[o.Field].cmp(t[o.Field])
		}

		switch r {
		case -1:
			return true
		case 1:
			return false
		}
	}
	return false
}

func (r *Raw) cmp(r2 *Raw) (v int) {
	if r.Bytes != nil {
		v = bytes.Compare(r.Bytes, r2.Bytes)
	} else {
		v = r.Value.Cmp(r2.Value)
	}
	return
}

func opFilter(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op, filters := ops[0], ops[1]

	op2 := &Operand{Meta: op.cloneMeta(), Data: make([]Tuple, 0)}

	for i := 0; i < len(filters.Data); i++ {
		if filters.Data[i][0].isTrue() {
			op2.Data = append(op2.Data, append(Tuple{}, op.Data[i]...))
		}
	}

	registers[output] = op2
	return
}

// Type check will ensure all cast is valid
func opCast(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}
	op := ops[0]
	dTypes := ops[1].Meta

	if len(dTypes) != len(op.Meta) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	op2 := &Operand{Meta: ops[1].cloneMeta(), Data: make([]Tuple, len(op.Data))}
	for i := 0; i < len(op.Data); i++ {
		op2.Data[i] = append(Tuple{}, op.Data[i]...)

		for j, dType := range dTypes {
			if op.Meta[j] == dType {
				continue
			}

			err = op2.Data[i][j].cast(ctx, op.Meta[j], dType)
			if err != nil {
				return
			}
		}
	}

	registers[output] = op2
	return
}

func (r *Raw) cast(ctx *common.Context, origin, target ast.DataType) (err error) {
	oMajor, _ := ast.DecomposeDataType(origin)

	// conversion table
	switch oMajor {
	case ast.DataTypeMajorInt, ast.DataTypeMajorUint:
		err = r.castInt(ctx, origin, target)
	case ast.DataTypeMajorFixedBytes:
		err = r.castFixedBytes(ctx, origin, target)
	case ast.DataTypeMajorAddress:
		err = r.castAddress(ctx, origin, target)
	case ast.DataTypeMajorBool:
		err = r.castBool(origin, target)
	case ast.DataTypeMajorDynamicBytes:
		err = r.castDynBytes(origin, target)
	default:
		err = se.ErrorCodeInvalidCastType
	}
	return
}

func (r *Raw) castValue(
	ctx *common.Context,
	origin, target ast.DataType,
	l int, signed, rPadding bool) (err error) {
	oBytes, err := ast.DecimalEncode(origin, r.Value)
	if err != nil {
		return
	}

	bytes2 := r.shiftBytes(oBytes, l, signed, rPadding)

	r.Value, err = ast.DecimalDecode(target, bytes2)
	if err != nil {
		return
	}

	err = flowCheck(ctx, r.Value, target)
	return
}

func (r *Raw) castInt(ctx *common.Context, origin, target ast.DataType) (err error) {
	oMajor, oMinor := ast.DecomposeDataType(origin)
	tMajor, tMinor := ast.DecomposeDataType(target)
	signed := oMajor == ast.DataTypeMajorInt

	switch tMajor {
	case ast.DataTypeMajorInt, ast.DataTypeMajorUint:
		err = r.castValue(ctx, origin, target, int(tMinor)+1, signed, false)
	case ast.DataTypeMajorAddress:
		r.Bytes, err = ast.DecimalEncode(origin, r.Value)
		if err != nil {
			return
		}

		if len(r.Bytes) > dexCommon.AddressLength {
			if r.Bytes[0]&0x80 != 0 && signed {
				err = se.ErrorCodeUnderflow
			} else {
				err = se.ErrorCodeOverflow
			}
			return
		}

		r.Bytes = r.shiftBytes(r.Bytes, dexCommon.AddressLength, signed, false)
		r.Value = decimal.Zero
	case ast.DataTypeMajorFixedBytes:
		if tMinor != oMinor {
			err = se.ErrorCodeInvalidCastType
			return
		}
		r.Bytes, err = ast.DecimalEncode(origin, r.Value)
		if err != nil {
			return
		}
		r.Value = decimal.Zero
	case ast.DataTypeMajorBool:
		r.Value = dec.Val2Bool(r.Value)
	default:
		err = se.ErrorCodeInvalidCastType
	}
	return
}

func (r *Raw) castFixedBytes(ctx *common.Context, origin, target ast.DataType) (err error) {
	_, oMinor := ast.DecomposeDataType(origin)
	tMajor, tMinor := ast.DecomposeDataType(target)
	switch tMajor {
	case ast.DataTypeMajorDynamicBytes:
	case ast.DataTypeMajorInt, ast.DataTypeMajorUint:
		if tMinor != oMinor {
			err = se.ErrorCodeInvalidCastType
			return
		}
		r.Value, err = ast.DecimalDecode(target, r.Bytes)
		if err != nil {
			return
		}
		r.Bytes = nil
	case ast.DataTypeMajorFixedBytes:
		r.Bytes = r.shiftBytes(r.Bytes, int(tMinor)+1, false, true)
	case ast.DataTypeMajorAddress:
		if oMinor != (dexCommon.AddressLength - 1) {
			err = se.ErrorCodeInvalidCastType
			return
		}
	default:
		err = se.ErrorCodeInvalidCastType
	}
	return
}

func (r *Raw) castAddress(ctx *common.Context, origin, target ast.DataType) (err error) {
	tMajor, tMinor := ast.DecomposeDataType(target)

	switch tMajor {
	case ast.DataTypeMajorAddress:
	case ast.DataTypeMajorInt, ast.DataTypeMajorUint:
		r.Value, err = ast.DecimalDecode(
			target,
			r.shiftBytes(r.Bytes, int(tMinor)+1, false, false),
		)
		if err != nil {
			return
		}
		err = flowCheck(ctx, r.Value, target)
		if err != nil {
			return
		}
		r.Bytes = nil
	case ast.DataTypeMajorFixedBytes:
		if tMinor != (dexCommon.AddressLength - 1) {
			err = se.ErrorCodeInvalidCastType
			return
		}
	default:
		err = se.ErrorCodeInvalidCastType
	}
	return
}

func (r *Raw) castBool(origin, target ast.DataType) (err error) {
	tMajor, _ := ast.DecomposeDataType(target)
	switch tMajor {
	case ast.DataTypeMajorBool, ast.DataTypeMajorInt, ast.DataTypeMajorUint:
	default:
		err = se.ErrorCodeInvalidCastType
	}
	return
}

func (r *Raw) castDynBytes(origin, target ast.DataType) (err error) {
	tMajor, tMinor := ast.DecomposeDataType(target)
	switch tMajor {
	case ast.DataTypeMajorDynamicBytes:
	case ast.DataTypeMajorFixedBytes:
		r.Bytes = r.shiftBytes(r.Bytes, int(tMinor)+1, false, true)
	default:
		err = se.ErrorCodeInvalidCastType
	}
	return
}

func (r *Raw) shiftBytes(src []byte, l int, signed, rPadding bool) (tgr []byte) {
	if len(src) >= l {
		if rPadding {
			tgr = src[:l]
		} else {
			tgr = src[len(src)-l:]
		}
		return
	}

	tgr = make([]byte, l)

	if rPadding {
		copy(tgr, src)
		return
	}

	copy(tgr[l-len(src):], src)

	if signed && src[0]&0x80 != 0 {
		for i := 0; i < l-len(src); i++ {
			tgr[i] = 0xff
		}
	}
	return
}

func opConcat(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 2 {
		err = se.ErrorCodeDataLengthNotMatch
		return
	}
	op, op2 := ops[0], ops[1]

	if !metaAllEq(op, op2) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	if !metaAllDynBytes(op) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	op3 := op.clone(true)
	op3.Data = make([]Tuple, len(op.Data))

	for i := 0; i < len(op.Data); i++ {
		op3.Data[i] = op.Data[i].concat(op2.Data[i])
	}

	registers[output] = op3
	return
}

func (t Tuple) concat(t2 Tuple) (t3 Tuple) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t3[i] = &Raw{Bytes: make([]byte, len(t[i].Bytes)+len(t2[i].Bytes))}
		copy(t3[i].Bytes[:len(t[i].Bytes)], t[i].Bytes)
		copy(t3[i].Bytes[len(t[i].Bytes):], t2[i].Bytes)
	}
	return
}

func opNeg(ctx *common.Context, ops, registers []*Operand, output uint) (err error) {
	if len(ops) != 1 {
		err = se.ErrorCodeDataLengthNotMatch
		return
	}
	op := ops[0]

	if !metaAllSignedNumeric(op) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	op2 := op.clone(true)
	op2.Data = make([]Tuple, len(op.Data))

	for i := 0; i < len(op.Data); i++ {
		op2.Data[i], err = op.Data[i].neg(ctx, op2.Meta)
		if err != nil {
			return
		}
	}

	registers[output] = op2
	return
}

func (t Tuple) neg(ctx *common.Context, meta []ast.DataType) (t2 Tuple, err error) {
	t2 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t2[i] = &Raw{Value: t[i].Value.Neg()}

		err = flowCheck(ctx, t2[i].Value, meta[i])
		if err != nil {
			return
		}
	}
	return
}
