package ast

import (
	"errors"
	"math/big"
	"reflect"

	"github.com/shopspring/decimal"

	"github.com/dexon-foundation/dexon/common"
)

var (
	bigIntOne = big.NewInt(1)
	bigIntTen = big.NewInt(10)
)

// Error defines.
var (
	ErrDataTypeEncode = errors.New("data type encode failed")
	ErrDataTypeDecode = errors.New("data type decode failed")
	ErrDecimalEncode  = errors.New("decimal encode failed")
	ErrDecimalDecode  = errors.New("decimal decode failed")
)

// DataTypeMajor defines type for high byte of DataType.
type DataTypeMajor uint8

// DataTypeMinor defines type for low byte of DataType.
type DataTypeMinor uint8

// DataType defines type for data type encoded.
type DataType uint16

// DataTypeMajor enums.
const (
	DataTypeMajorUnknown DataTypeMajor = iota
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

// DataTypeUnknown for unknown data type.
const DataTypeUnknown DataType = 0

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

// DataTypeEncode encodes data type node into DataType.
func DataTypeEncode(n interface{}) (DataType, error) {
	if n == nil {
		return DataTypeUnknown, ErrDataTypeEncode
	}
	if reflect.TypeOf(n).Kind() == reflect.Ptr {
		return DataTypeEncode(reflect.ValueOf(n).Elem())
	}

	switch t := n.(type) {
	case BoolTypeNode:
		return ComposeDataType(DataTypeMajorBool, 0), nil

	case AddressTypeNode:
		return ComposeDataType(DataTypeMajorAddress, 0), nil

	case IntTypeNode:
		if t.Size%8 != 0 || t.Size > 256 {
			return DataTypeUnknown, ErrDataTypeEncode
		}

		minor := DataTypeMinor((t.Size / 8) - 1)
		if t.Unsigned {
			return ComposeDataType(DataTypeMajorUint, minor), nil
		}
		return ComposeDataType(DataTypeMajorInt, minor), nil

	case FixedBytesTypeNode:
		if t.Size%8 != 0 || t.Size > 256 {
			return DataTypeUnknown, ErrDataTypeEncode
		}

		minor := DataTypeMinor((t.Size / 8) - 1)
		return ComposeDataType(DataTypeMajorFixedBytes, minor), nil

	case DynamicBytesTypeNode:
		return ComposeDataType(DataTypeMajorDynamicBytes, 0), nil

	case FixedTypeNode:
		if t.Size%8 != 0 || t.Size > 256 {
			return DataTypeUnknown, ErrDataTypeEncode
		}

		if t.FractionalDigits > 80 {
			return DataTypeUnknown, ErrDataTypeEncode
		}

		major := DataTypeMajor((t.Size / 8) - 1)
		minor := DataTypeMinor(t.FractionalDigits)
		if t.Unsigned {
			return ComposeDataType(DataTypeMajorUfixed+major, minor), nil
		}
		return ComposeDataType(DataTypeMajorFixed+major, minor), nil
	}

	return DataTypeUnknown, ErrDataTypeEncode
}

// DataTypeDecode decodes DataType into data type node.
func DataTypeDecode(t DataType) (interface{}, error) {
	major, minor := DecomposeDataType(t)
	switch major {
	// TODO(wmin0): define unsupported error for special type.
	case DataTypeMajorBool:
		if minor == 0 {
			return BoolTypeNode{}, nil
		}
	case DataTypeMajorAddress:
		if minor == 0 {
			return AddressTypeNode{}, nil
		}
	case DataTypeMajorInt:
		if minor <= 0x1f {
			size := (uint32(minor) + 1) * 8
			return IntTypeNode{Unsigned: false, Size: size}, nil
		}
	case DataTypeMajorUint:
		if minor <= 0x1f {
			size := (uint32(minor) + 1) * 8
			return IntTypeNode{Unsigned: true, Size: size}, nil
		}
	case DataTypeMajorFixedBytes:
		if minor <= 0x1f {
			size := (uint32(minor) + 1) * 8
			return FixedBytesTypeNode{Size: size}, nil
		}
	case DataTypeMajorDynamicBytes:
		if minor == 0 {
			return DynamicBytesTypeNode{}, nil
		}
	}
	switch {
	case major.IsFixedRange():
		if minor <= 80 {
			size := (uint32(major-DataTypeMajorFixed) + 1) * 8
			return FixedTypeNode{
				Unsigned:         false,
				Size:             size,
				FractionalDigits: uint32(minor),
			}, nil
		}
	case major.IsUfixedRange():
		if minor <= 80 {
			size := (uint32(major-DataTypeMajorUfixed) + 1) * 8
			return FixedTypeNode{
				Unsigned:         true,
				Size:             size,
				FractionalDigits: uint32(minor),
			}, nil
		}
	}
	return nil, ErrDataTypeDecode
}

// Don't handle overflow here.
func decimalEncode(size int, d decimal.Decimal) []byte {
	ret := make([]byte, size)
	s := d.Sign()
	if s == 0 {
		return ret
	}

	exp := new(big.Int).Exp(bigIntTen, big.NewInt(int64(d.Exponent())), nil)
	b := new(big.Int).Mul(d.Coefficient(), exp)

	if s > 0 {
		bs := b.Bytes()
		copy(ret[size-len(bs):], bs)
		return ret
	}

	b.Add(b, bigIntOne)
	bs := b.Bytes()
	copy(ret[size-len(bs):], bs)
	for idx := range ret {
		ret[idx] = ^ret[idx]
	}
	return ret
}

// Don't handle overflow here.
func decimalDecode(signed bool, bs []byte) decimal.Decimal {
	neg := false
	if signed && (bs[0]&0x80 != 0) {
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
func DecimalEncode(dt DataType, d decimal.Decimal) ([]byte, error) {
	major, minor := DecomposeDataType(dt)
	switch major {
	case DataTypeMajorInt,
		DataTypeMajorUint:
		return decimalEncode(int(minor)+1, d), nil
	case DataTypeMajorAddress:
		return decimalEncode(common.AddressLength, d), nil
	}
	switch {
	case major.IsFixedRange():
		return decimalEncode(
			int(major-DataTypeMajorFixed)+1,
			d.Shift(int32(minor))), nil
	case major.IsUfixedRange():
		return decimalEncode(
			int(major-DataTypeMajorUfixed)+1,
			d.Shift(int32(minor))), nil
	}

	return nil, ErrDecimalEncode
}

// DecimalDecode decodes decimal from bytes.
func DecimalDecode(dt DataType, b []byte) (decimal.Decimal, error) {
	major, minor := DecomposeDataType(dt)
	switch major {
	case DataTypeMajorInt:
		return decimalDecode(true, b), nil
	case DataTypeMajorUint,
		DataTypeMajorAddress:
		return decimalDecode(false, b), nil
	}
	switch {
	case major.IsFixedRange():
		return decimalDecode(true, b).Shift(-int32(minor)), nil
	case major.IsUfixedRange():
		return decimalDecode(false, b).Shift(-int32(minor)), nil
	}

	return decimal.Zero, ErrDecimalDecode
}
