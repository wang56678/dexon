package ast

import (
	"errors"
	"reflect"
)

// Error defines.
var (
	ErrDataTypeEncode = errors.New("data type encode failed")
	ErrDataTypeDecode = errors.New("data type decode failed")
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

func decomposeDataType(t DataType) (DataTypeMajor, DataTypeMinor) {
	return DataTypeMajor(t >> 8), DataTypeMinor(t & 0xff)
}

func composeDataType(major DataTypeMajor, minor DataTypeMinor) DataType {
	return (DataType(major) << 8) | DataType(minor)
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
		return composeDataType(DataTypeMajorBool, 0), nil

	case AddressTypeNode:
		return composeDataType(DataTypeMajorAddress, 0), nil

	case IntTypeNode:
		if t.Size%8 != 0 || t.Size > 256 {
			return DataTypeUnknown, ErrDataTypeEncode
		}

		minor := DataTypeMinor((t.Size / 8) - 1)
		if t.Unsigned {
			return composeDataType(DataTypeMajorUint, minor), nil
		}
		return composeDataType(DataTypeMajorInt, minor), nil

	case FixedBytesTypeNode:
		if t.Size%8 != 0 || t.Size > 256 {
			return DataTypeUnknown, ErrDataTypeEncode
		}

		minor := DataTypeMinor((t.Size / 8) - 1)
		return composeDataType(DataTypeMajorFixedBytes, minor), nil

	case DynamicBytesTypeNode:
		return composeDataType(DataTypeMajorDynamicBytes, 0), nil

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
			return composeDataType(DataTypeMajorUfixed+major, minor), nil
		}
		return composeDataType(DataTypeMajorFixed+major, minor), nil
	}

	return DataTypeUnknown, ErrDataTypeEncode
}

// DataTypeDecode decodes DataType into data type node.
func DataTypeDecode(t DataType) (interface{}, error) {
	major, minor := decomposeDataType(t)
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
	case major >= DataTypeMajorFixed && major-DataTypeMajorFixed <= 0x1f:
		if minor <= 80 {
			size := (uint32(major-DataTypeMajorFixed) + 1) * 8
			return FixedTypeNode{
				Unsigned:         false,
				Size:             size,
				FractionalDigits: uint32(minor),
			}, nil
		}
	case major >= DataTypeMajorUfixed && major-DataTypeMajorUfixed <= 0x1f:
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
