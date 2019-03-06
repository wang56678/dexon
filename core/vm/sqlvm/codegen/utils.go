package codegen

import (
	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/runtime"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

func newImmDummy() *runtime.Operand {
	return &runtime.Operand{
		IsImmediate: true,
		Meta:        []ast.DataType{},
		Data:        []runtime.Tuple{runtime.Tuple{}},
	}
}

func newImmType(t ...ast.DataType) *runtime.Operand {
	return &runtime.Operand{
		IsImmediate: true,
		Meta:        t,
	}
}

func newImmColDecimal(t ast.DataType, d ...decimal.Decimal) *runtime.Operand {
	data := make([]runtime.Tuple, len(d))
	for idx, v := range d {
		data[idx] = runtime.Tuple{
			&runtime.Raw{Value: v},
		}
	}
	return &runtime.Operand{
		IsImmediate: true,
		Meta:        []ast.DataType{t},
		Data:        data,
	}
}

func newImmRowDecimal(
	ts []ast.DataType, d ...decimal.Decimal) *runtime.Operand {
	data := make(runtime.Tuple, len(d))
	for idx, v := range d {
		data[idx] = &runtime.Raw{Value: v}
	}
	return &runtime.Operand{
		IsImmediate: true,
		Meta:        ts,
		Data:        []runtime.Tuple{data},
	}
}

func newImmColBytes(t ast.DataType, b ...[]byte) *runtime.Operand {
	data := make([]runtime.Tuple, len(b))
	for idx, v := range b {
		data[idx] = runtime.Tuple{
			&runtime.Raw{Bytes: v},
		}
	}
	return &runtime.Operand{
		IsImmediate: true,
		Meta:        []ast.DataType{t},
		Data:        data,
	}
}

func newImmRowBytes(ts []ast.DataType, b ...[]byte) *runtime.Operand {
	data := make(runtime.Tuple, len(b))
	for idx, v := range b {
		data[idx] = &runtime.Raw{Bytes: v}
	}
	return &runtime.Operand{
		IsImmediate: true,
		Meta:        ts,
		Data:        []runtime.Tuple{data},
	}
}

func newReg(i uint) *runtime.Operand {
	return &runtime.Operand{
		RegisterIndex: i,
	}
}

func expandImmCol(os []*runtime.Operand) *runtime.Operand {
	col := 0
	row := -1
	for _, o := range os {
		if !o.IsImmediate {
			panic("invalid operand in expandImmCol")
		}
		col += len(o.Meta)
		if row == -1 {
			row = len(o.Data)
		} else if row != len(o.Data) {
			panic("row not align in expandImmCol")
		}
	}
	result := &runtime.Operand{
		IsImmediate: true,
		Meta:        make([]ast.DataType, 0, col),
		Data:        make([]runtime.Tuple, row),
	}
	for _, o := range os {
		result.Meta = append(result.Meta, o.Meta...)
		for idx := range result.Data {
			if result.Data[idx] == nil {
				result.Data[idx] = make(runtime.Tuple, 0, col)
			}
			result.Data[idx] = append(result.Data[idx], o.Data[idx]...)
		}
	}
	return result
}

func expandImmRow(os []*runtime.Operand) *runtime.Operand {
	var meta []ast.DataType
	row := 0
	col := 0
	for _, o := range os {
		if !o.IsImmediate {
			panic("invalid operand in expandImmRow")
		}
		if meta == nil {
			meta = o.Meta
			col = len(meta)
		} else {
			// TODO(wmin0): utility?
			if len(meta) != len(o.Meta) {
				panic("invalid meta in expandImmRow")
			}
			for idx := range meta {
				if meta[idx] != o.Meta[idx] {
					panic("invalid meta in expandImmRow")
				}
			}
		}
		row += len(o.Data)
	}

	result := &runtime.Operand{
		IsImmediate: true,
		Meta:        make([]ast.DataType, 0, col),
		Data:        make([]runtime.Tuple, row),
	}
	result.Meta = append(result.Meta, meta...)
	for _, o := range os {
		result.Data = append(result.Data, o.Data...)
	}
	return result
}

func newRefDataType() ast.DataType {
	return ast.ComposeDataType(ast.DataTypeMajorUint, 0)
}

func newTableRefImm(ref schema.TableRef) *runtime.Operand {
	return newImmColDecimal(
		newRefDataType(),
		decimal.New(int64(ref), 0),
	)
}

func newColumnRefImm(ref schema.ColumnRef) *runtime.Operand {
	return newImmColDecimal(
		newRefDataType(),
		decimal.New(int64(ref), 0),
	)
}

func newColumnRefsImm(refs []schema.ColumnRef) *runtime.Operand {
	if len(refs) == 0 {
		return newImmDummy()
	}
	meta := make([]ast.DataType, len(refs))
	cols := make([]decimal.Decimal, len(refs))
	for idx, r := range refs {
		meta[idx] = newRefDataType()
		cols[idx] = decimal.New(int64(r), 0)
	}
	return newImmRowDecimal(meta, cols...)
}

func newIndexRefImm(
	table schema.TableRef, ref schema.IndexRef) *runtime.Operand {
	return newImmRowDecimal(
		[]ast.DataType{
			newRefDataType(),
			newRefDataType(),
		},
		decimal.New(int64(table), 0),
		decimal.New(int64(ref), 0),
	)
}

func newOrderImm(field int, desc bool) *runtime.Operand {
	d := dec.False
	if desc {
		d = dec.True
	}
	return newImmRowDecimal(
		[]ast.DataType{
			ast.ComposeDataType(ast.DataTypeMajorUint, 1),
			ast.ComposeDataType(ast.DataTypeMajorBool, ast.DataTypeMinorDontCare),
		},
		decimal.New(int64(field), 0),
		d,
	)
}
