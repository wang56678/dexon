package ast

import (
	"database/sql"
	"fmt"
	"math"
	"math/big"

	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	se "github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

var (
	bigIntOne = big.NewInt(1)
)

// BoolValue represents a boolean value used by SQL three-valued logic.
type BoolValue uint8

// Define valid values for SQL boolean type. The zero value is invalid.
const (
	BoolValueTrue    BoolValue = 1
	BoolValueFalse   BoolValue = 2
	BoolValueUnknown BoolValue = 3
)

// DataTypeMajor defines type for high byte of DataType.
type DataTypeMajor uint8

// DataTypeMinor defines type for low byte of DataType.
type DataTypeMinor uint8

// DataType defines type for data type encoded.
type DataType uint16

// DataTypeMajor enums.
const (
	DataTypeMajorPending DataTypeMajor = iota
	DataTypeMajorSpecial
	DataTypeMajorBool
	DataTypeMajorAddress
	DataTypeMajorInt
	DataTypeMajorUint
	DataTypeMajorFixedBytes
	DataTypeMajorDynamicBytes
	DataTypeMajorFixed  DataTypeMajor = 0x10
	DataTypeMajorUfixed DataTypeMajor = 0x30
)

// DataTypeMinor enums.
const (
	DataTypeMinorDontCare    DataTypeMinor = 0x00
	DataTypeMinorSpecialNull DataTypeMinor = 0x00
)

// Special data types which are commonly used.
const (
	DataTypePending DataType = (DataType(DataTypeMajorPending) << 8) | DataType(DataTypeMinorDontCare)
	DataTypeNull    DataType = (DataType(DataTypeMajorSpecial) << 8) | DataType(DataTypeMinorSpecialNull)
	DataTypeBad     DataType = math.MaxUint16
)

// Valid returns whether a BoolValue is valid.
func (v BoolValue) Valid() bool {
	return v-1 < 3
}

var boolValueStringMap = [3]string{
	BoolValueTrue - 1:    "TRUE",
	BoolValueFalse - 1:   "FALSE",
	BoolValueUnknown - 1: "UNKNOWN",
}

// String returns a string for printing a BoolValue.
func (v BoolValue) String() string {
	return boolValueStringMap[v-1]
}

var boolValueNullBoolMap = [3]sql.NullBool{
	BoolValueTrue - 1:    {Valid: true, Bool: true},
	BoolValueFalse - 1:   {Valid: true, Bool: false},
	BoolValueUnknown - 1: {Valid: false, Bool: false},
}

// NullBool converts a BoolValue to a sql.NullBool.
func (v BoolValue) NullBool() sql.NullBool {
	return boolValueNullBoolMap[v-1]
}

var boolValueAndTruthTable = [3][3]BoolValue{
	BoolValueTrue - 1: {
		BoolValueTrue - 1:    BoolValueTrue,
		BoolValueFalse - 1:   BoolValueFalse,
		BoolValueUnknown - 1: BoolValueUnknown,
	},
	BoolValueFalse - 1: {
		BoolValueTrue - 1:    BoolValueFalse,
		BoolValueFalse - 1:   BoolValueFalse,
		BoolValueUnknown - 1: BoolValueFalse,
	},
	BoolValueUnknown - 1: {
		BoolValueTrue - 1:    BoolValueUnknown,
		BoolValueFalse - 1:   BoolValueFalse,
		BoolValueUnknown - 1: BoolValueUnknown,
	},
}

// And returns v AND v2.
func (v BoolValue) And(v2 BoolValue) BoolValue {
	return boolValueAndTruthTable[v-1][v2-1]
}

var boolValueOrTruthTable = [3][3]BoolValue{
	BoolValueTrue - 1: {
		BoolValueTrue - 1:    BoolValueTrue,
		BoolValueFalse - 1:   BoolValueTrue,
		BoolValueUnknown - 1: BoolValueTrue,
	},
	BoolValueFalse - 1: {
		BoolValueTrue - 1:    BoolValueTrue,
		BoolValueFalse - 1:   BoolValueFalse,
		BoolValueUnknown - 1: BoolValueUnknown,
	},
	BoolValueUnknown - 1: {
		BoolValueTrue - 1:    BoolValueTrue,
		BoolValueFalse - 1:   BoolValueUnknown,
		BoolValueUnknown - 1: BoolValueUnknown,
	},
}

// Or returns v OR v2.
func (v BoolValue) Or(v2 BoolValue) BoolValue {
	return boolValueOrTruthTable[v-1][v2-1]
}

var boolValueNotTruthTable = [3]BoolValue{
	BoolValueTrue - 1:    BoolValueFalse,
	BoolValueFalse - 1:   BoolValueTrue,
	BoolValueUnknown - 1: BoolValueUnknown,
}

// Not returns NOT v.
func (v BoolValue) Not() BoolValue {
	return boolValueNotTruthTable[v-1]
}

// DecomposeDataType to major and minor part with given data type.
func DecomposeDataType(t DataType) (DataTypeMajor, DataTypeMinor) {
	return DataTypeMajor(t >> 8), DataTypeMinor(t & 0xff)
}

// ComposeDataType to concrete type with major and minor part.
func ComposeDataType(major DataTypeMajor, minor DataTypeMinor) DataType {
	return (DataType(major) << 8) | DataType(minor)
}

// IsFixedRange checks if major is in range of DataTypeMajorFixed.
func (d DataTypeMajor) IsFixedRange() bool {
	return d >= DataTypeMajorFixed && d-DataTypeMajorFixed <= 0x1f
}

// IsUfixedRange checks if major is in range of DataTypeMajorUfixed.
func (d DataTypeMajor) IsUfixedRange() bool {
	return d >= DataTypeMajorUfixed && d-DataTypeMajorUfixed <= 0x1f
}

// Size return the bytes of the data type occupied.
func (dt DataType) Size() uint8 {
	major, minor := DecomposeDataType(dt)
	if major.IsFixedRange() {
		return uint8(major - DataTypeMajorFixed + 1)
	}
	if major.IsUfixedRange() {
		return uint8(major - DataTypeMajorUfixed + 1)
	}
	switch major {
	case DataTypeMajorBool:
		return 1
	case DataTypeMajorDynamicBytes:
		return common.HashLength
	case DataTypeMajorAddress:
		return common.AddressLength
	case DataTypeMajorInt, DataTypeMajorUint, DataTypeMajorFixedBytes:
		return uint8(minor + 1)
	default:
		panic(fmt.Sprintf("unknown data type %v", dt))
	}
}

// Valid checks whether a data type is set to a defined value.
func (dt DataType) Valid() bool {
	major, minor := DecomposeDataType(dt)
	switch major {
	case DataTypeMajorPending,
		DataTypeMajorBool,
		DataTypeMajorAddress,
		DataTypeMajorDynamicBytes:
		return true
	case DataTypeMajorSpecial:
		switch minor {
		case DataTypeMinorSpecialNull:
			return true
		}
	case DataTypeMajorInt,
		DataTypeMajorUint,
		DataTypeMajorFixedBytes:
		if minor <= 0x1f {
			return true
		}
	}
	if (major.IsFixedRange() || major.IsUfixedRange()) && minor <= 80 {
		return true
	}
	return false
}

// ValidColumn checks whether a data type is a valid column type.
func (dt DataType) ValidColumn() bool {
	major, minor := DecomposeDataType(dt)
	switch major {
	case DataTypeMajorBool,
		DataTypeMajorAddress,
		DataTypeMajorDynamicBytes:
		return true
	case DataTypeMajorInt,
		DataTypeMajorUint,
		DataTypeMajorFixedBytes:
		if minor <= 0x1f {
			return true
		}
	}
	if (major.IsFixedRange() || major.IsUfixedRange()) && minor <= 80 {
		return true
	}
	return false
}

// ValidExpr checks whether a data type is a valid type of an expression.
func (dt DataType) ValidExpr() bool {
	return dt.ValidColumn() || dt == DataTypeNull
}

// Equal checks whether two data types are equal and valid. If any of them is
// invalid, false is returned.
func (dt DataType) Equal(dt2 DataType) bool {
	// Rename variables.
	a := dt
	b := dt2
	// Process the common case.
	if a == b {
		return a.Valid()
	}
	// a ≠ b
	aMajor, _ := DecomposeDataType(a)
	bMajor, _ := DecomposeDataType(b)
	if aMajor != bMajor {
		return false
	}
	// a ≠ b, aMajor = bMajor ⇒ aMinor ≠ bMinor
	switch aMajor {
	case DataTypeMajorPending,
		DataTypeMajorBool,
		DataTypeMajorAddress,
		DataTypeMajorDynamicBytes:
		return true
	default:
		return false
	}
}

func (dt DataType) String() string {
	major, minor := DecomposeDataType(dt)
	switch major {
	case DataTypeMajorPending:
		return "<PENDING>"
	case DataTypeMajorSpecial:
		switch minor {
		case DataTypeMinorSpecialNull:
			return "NULL"
		}
	case DataTypeMajorBool:
		return "BOOL"
	case DataTypeMajorAddress:
		return "ADDRESS"
	case DataTypeMajorInt:
		if minor <= 0x1f {
			size := (uint32(minor) + 1) * 8
			return fmt.Sprintf("INT%d", size)
		}
	case DataTypeMajorUint:
		if minor <= 0x1f {
			size := (uint32(minor) + 1) * 8
			return fmt.Sprintf("UINT%d", size)
		}
	case DataTypeMajorFixedBytes:
		if minor <= 0x1f {
			size := uint32(minor) + 1
			return fmt.Sprintf("BYTES%d", size)
		}
	case DataTypeMajorDynamicBytes:
		return "BYTES"
	}
	switch {
	case major.IsFixedRange():
		if minor <= 80 {
			size := (uint32(major-DataTypeMajorFixed) + 1) * 8
			fractionalDigits := uint32(minor)
			return fmt.Sprintf("FIXED%dX%d", size, fractionalDigits)
		}
	case major.IsUfixedRange():
		if minor <= 80 {
			size := (uint32(major-DataTypeMajorUfixed) + 1) * 8
			fractionalDigits := uint32(minor)
			return fmt.Sprintf("UFIXED%dX%d", size, fractionalDigits)
		}
	}
	return "<INVALID>"
}

// GetNode constructs an AST node from a data type.
func (dt DataType) GetNode() TypeNode {
	major, minor := DecomposeDataType(dt)
	switch major {
	case DataTypeMajorBool:
		return &BoolTypeNode{}
	case DataTypeMajorAddress:
		return &AddressTypeNode{}
	case DataTypeMajorInt:
		if minor <= 0x1f {
			size := (uint32(minor) + 1) * 8
			return &IntTypeNode{Unsigned: false, Size: size}
		}
	case DataTypeMajorUint:
		if minor <= 0x1f {
			size := (uint32(minor) + 1) * 8
			return &IntTypeNode{Unsigned: true, Size: size}
		}
	case DataTypeMajorFixedBytes:
		if minor <= 0x1f {
			size := uint32(minor) + 1
			return &FixedBytesTypeNode{Size: size}
		}
	case DataTypeMajorDynamicBytes:
		return &DynamicBytesTypeNode{}
	}
	switch {
	case major.IsFixedRange():
		if minor <= 80 {
			size := (uint32(major-DataTypeMajorFixed) + 1) * 8
			fractionalDigits := uint32(minor)
			return &FixedTypeNode{
				Unsigned:         false,
				Size:             size,
				FractionalDigits: fractionalDigits,
			}
		}
	case major.IsUfixedRange():
		if minor <= 80 {
			size := (uint32(major-DataTypeMajorUfixed) + 1) * 8
			fractionalDigits := uint32(minor)
			return &FixedTypeNode{
				Unsigned:         true,
				Size:             size,
				FractionalDigits: fractionalDigits,
			}
		}
	}
	return nil
}

type decimalMinMaxPair struct {
	Min, Max decimal.Decimal
}

var decimalMinMaxMap = func() map[DataType]decimalMinMaxPair {
	m := make(map[DataType]decimalMinMaxPair)

	// bool
	{
		dt := ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
		m[dt] = decimalMinMaxPair{Min: dec.False, Max: dec.True}
	}

	// address
	{
		dt := ComposeDataType(DataTypeMajorAddress, DataTypeMinorDontCare)
		bigMax := new(big.Int)
		bigMax.Lsh(bigIntOne, common.AddressLength*8)
		bigMax.Sub(bigMax, bigIntOne)
		max := decimal.NewFromBigInt(bigMax, 0)
		m[dt] = decimalMinMaxPair{Min: decimal.Zero, Max: max}
	}

	// uint, bytes{X}
	for i := uint(0); i <= 0x1f; i++ {
		dtUint := ComposeDataType(DataTypeMajorUint, DataTypeMinor(i))
		dtBytes := ComposeDataType(DataTypeMajorFixedBytes, DataTypeMinor(i))
		bigMax := new(big.Int)
		bigMax.Lsh(bigIntOne, (i+1)*8)
		bigMax.Sub(bigMax, bigIntOne)
		max := decimal.NewFromBigInt(bigMax, 0)
		m[dtUint] = decimalMinMaxPair{Min: decimal.Zero, Max: max}
		m[dtBytes] = decimalMinMaxPair{Min: decimal.Zero, Max: max}
	}

	// int
	for i := uint(0); i <= 0x1f; i++ {
		dt := ComposeDataType(DataTypeMajorInt, DataTypeMinor(i))
		bigMax := new(big.Int)
		bigMax.Lsh(bigIntOne, (i+1)*8-1)
		bigMin := new(big.Int)
		bigMin.Neg(bigMax)
		bigMax.Sub(bigMax, bigIntOne)
		min := decimal.NewFromBigInt(bigMin, 0)
		max := decimal.NewFromBigInt(bigMax, 0)
		m[dt] = decimalMinMaxPair{Min: min, Max: max}
	}

	return m
}()

// GetMinMax returns min, max pair according to given data type.
func (dt DataType) GetMinMax() (decimal.Decimal, decimal.Decimal, bool) {
	pair, ok := decimalMinMaxMap[dt]
	if ok {
		return pair.Min, pair.Max, true
	}

	// Compute the range of fixed and ufixed types on demand.
	major, minor := DecomposeDataType(dt)
	switch {
	case major.IsFixedRange():
		mapInt := ComposeDataType(
			DataTypeMajorInt, DataTypeMinor(major-DataTypeMajorFixed))
		pair, ok = decimalMinMaxMap[mapInt]
		if !ok || minor > 80 {
			return decimal.Zero, decimal.Zero, false
		}
		min := pair.Min.Shift(-int32(minor))
		max := pair.Max.Shift(-int32(minor))
		return min, max, true

	case major.IsUfixedRange():
		mapUint := ComposeDataType(
			DataTypeMajorUint, DataTypeMinor(major-DataTypeMajorUfixed))
		pair, ok = decimalMinMaxMap[mapUint]
		if !ok || minor > 80 {
			return decimal.Zero, decimal.Zero, false
		}
		min := pair.Min.Shift(-int32(minor))
		max := pair.Max.Shift(-int32(minor))
		return min, max, true

	default:
		return decimal.Zero, decimal.Zero, false
	}
}

func decimalToBig(d decimal.Decimal) *big.Int {
	d = d.Rescale(0)
	return d.Coefficient()
}

// Don't handle overflow here.
func decimalEncode(size int, d decimal.Decimal) []byte {
	ret := make([]byte, size)
	s := d.Sign()
	if s == 0 {
		return ret
	}

	b := decimalToBig(d)

	if s > 0 {
		bs := b.Bytes()
		if size >= len(bs) {
			copy(ret[size-len(bs):], bs)
		} else {
			copy(ret, bs[len(bs)-size:])
		}
		return ret
	}

	b.Add(b, bigIntOne)
	bs := b.Bytes()
	if size >= len(bs) {
		copy(ret[size-len(bs):], bs)
	} else {
		copy(ret, bs[len(bs)-size:])
	}
	for idx := range ret {
		ret[idx] = ^ret[idx]
	}
	return ret
}

// Don't handle overflow here.
func decimalDecode(signed bool, bs []byte) decimal.Decimal {
	neg := false
	if signed && (bs[0]&0x80 != 0) {
		newbs := make([]byte, 0, len(bs))
		bs = append(newbs, bs...)
		neg = true
		for idx := range bs {
			bs[idx] = ^bs[idx]
		}
	}

	b := new(big.Int).SetBytes(bs)

	if neg {
		b.Add(b, bigIntOne)
		b.Neg(b)
	}

	return decimal.NewFromBigInt(b, 0)
}

// DecimalEncode encodes decimal to bytes depend on data type.
func DecimalEncode(dt DataType, d decimal.Decimal) ([]byte, bool) {
	major, minor := DecomposeDataType(dt)
	switch major {
	case DataTypeMajorInt,
		DataTypeMajorUint:
		return decimalEncode(int(minor)+1, d), true
	}
	switch {
	case major.IsFixedRange():
		return decimalEncode(
			int(major-DataTypeMajorFixed)+1,
			d.Shift(int32(minor))), true
	case major.IsUfixedRange():
		return decimalEncode(
			int(major-DataTypeMajorUfixed)+1,
			d.Shift(int32(minor))), true
	}

	return nil, false
}

// DecimalDecode decodes decimal from bytes.
func DecimalDecode(dt DataType, b []byte) (decimal.Decimal, bool) {
	major, minor := DecomposeDataType(dt)
	switch major {
	case DataTypeMajorInt:
		return decimalDecode(true, b), true
	case DataTypeMajorUint:
		return decimalDecode(false, b), true
	case DataTypeMajorBool:
		if b[0] == 0 {
			return dec.False, true
		}
		return dec.True, true
	}
	switch {
	case major.IsFixedRange():
		return decimalDecode(true, b).Shift(-int32(minor)), true
	case major.IsUfixedRange():
		return decimalDecode(false, b).Shift(-int32(minor)), true
	}

	return decimal.Zero, false
}

// DecimalToUint64 convert decimal to uint64.
// Negative case will return error, and decimal part will be trancated.
func DecimalToUint64(d decimal.Decimal) (uint64, error) {
	s := d.Sign()
	if s == 0 {
		return 0, nil
	}
	if s < 0 {
		return 0, se.ErrorCodeNegDecimalToUint64
	}
	return decimalToBig(d).Uint64(), nil
}
