package ast

import (
	"database/sql"
	"testing"

	"github.com/dexon-foundation/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/dexon-foundation/dexon/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

type TypesTestSuite struct{ suite.Suite }

func (s *TypesTestSuite) requireEncodeAndDecodeDecimalNoError(
	d DataType, t decimal.Decimal, bs int) {
	encode, ok := DecimalEncode(d, t)
	s.Require().True(ok)
	s.Require().Len(encode, bs)
	decode, ok := DecimalDecode(d, encode)
	s.Require().True(ok)
	s.Require().Equal(t.String(), decode.String())
}

func (s *TypesTestSuite) requireEncodeAndDecodeDataTypeNoError(
	d DataType, t TypeNode) {
	encode, code, message := t.GetType()
	s.Require().Zero(code)
	s.Require().Empty(message)
	s.Require().Equal(d, encode)
	decode := d.GetNode()
	s.Require().Equal(t, decode)
}

func (s *TypesTestSuite) requireEncodeDataTypeError(input TypeNode) {
	_, code, message := input.GetType()
	s.Require().NotZero(code)
	s.Require().NotEmpty(message)
}

func (s *TypesTestSuite) requireDecodeDataTypeError(input DataType) {
	decode := input.GetNode()
	s.Require().Nil(decode)
}

func (s *TypesTestSuite) TestEncodeAndDecodeDataType() {
	s.requireEncodeAndDecodeDataTypeNoError(
		ComposeDataType(DataTypeMajorBool, 0),
		&BoolTypeNode{})
	s.requireEncodeAndDecodeDataTypeNoError(
		ComposeDataType(DataTypeMajorAddress, 0),
		&AddressTypeNode{})
	s.requireEncodeAndDecodeDataTypeNoError(
		ComposeDataType(DataTypeMajorInt, 1),
		&IntTypeNode{Size: 16})
	s.requireEncodeAndDecodeDataTypeNoError(
		ComposeDataType(DataTypeMajorUint, 2),
		&IntTypeNode{Unsigned: true, Size: 24})
	s.requireEncodeAndDecodeDataTypeNoError(
		ComposeDataType(DataTypeMajorFixedBytes, 3),
		&FixedBytesTypeNode{Size: 4})
	s.requireEncodeAndDecodeDataTypeNoError(
		ComposeDataType(DataTypeMajorDynamicBytes, 0),
		&DynamicBytesTypeNode{})
	s.requireEncodeAndDecodeDataTypeNoError(
		ComposeDataType(DataTypeMajorFixed, 1),
		&FixedTypeNode{Size: 8, FractionalDigits: 1})
	s.requireEncodeAndDecodeDataTypeNoError(
		ComposeDataType(DataTypeMajorUfixed+1, 2),
		&FixedTypeNode{Unsigned: true, Size: 16, FractionalDigits: 2})
}

func (s *TypesTestSuite) TestEncodeDataTypeError() {
	s.requireEncodeDataTypeError(&IntTypeNode{Size: 1})
	s.requireEncodeDataTypeError(&IntTypeNode{Size: 257})
	s.requireEncodeDataTypeError(&FixedBytesTypeNode{Size: 0})
	s.requireEncodeDataTypeError(&FixedBytesTypeNode{Size: 257})
	s.requireEncodeDataTypeError(&FixedTypeNode{Size: 1, FractionalDigits: 0})
	s.requireEncodeDataTypeError(&FixedTypeNode{Size: 257, FractionalDigits: 0})
	s.requireEncodeDataTypeError(&FixedTypeNode{Size: 8, FractionalDigits: 81})
}

func (s *TypesTestSuite) TestDecodeDataTypeError() {
	s.requireDecodeDataTypeError(DataTypePending)
	s.requireDecodeDataTypeError(DataTypeBad)
	s.requireDecodeDataTypeError(ComposeDataType(DataTypeMajorInt, 0x20))
	s.requireDecodeDataTypeError(ComposeDataType(DataTypeMajorUint, 0x20))
	s.requireDecodeDataTypeError(ComposeDataType(DataTypeMajorFixedBytes, 0x20))
	s.requireDecodeDataTypeError(ComposeDataType(DataTypeMajorFixed, 81))
	s.requireDecodeDataTypeError(ComposeDataType(DataTypeMajorUfixed, 81))
	s.requireDecodeDataTypeError(ComposeDataType(DataTypeMajorUfixed+0x20, 80))
}

func (s *TypesTestSuite) TestEncodeAndDecodeDecimal() {
	pos := decimal.New(15, 1)
	zero := decimal.Zero
	neg := decimal.New(-150, -1)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorInt, 2),
		pos,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorInt, 2),
		zero,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorInt, 2),
		neg,
		3)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorUint, 2),
		pos,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorUint, 2),
		zero,
		3)

	pos = decimal.New(15, -2)
	neg = decimal.New(-15, -2)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorFixed+2, 2),
		pos,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorFixed+2, 2),
		zero,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorFixed+2, 2),
		neg,
		3)

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorUfixed+2, 2),
		pos,
		3)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorUfixed+2, 2),
		zero,
		3)
}

func (s *TypesTestSuite) TestDataTypeGetMinMax() {
	decAddressMax := decimal.Two.Pow(decimal.New(common.AddressLength*8, 0)).Sub(decimal.One)
	testcases := []struct {
		Name     string
		In       DataType
		Min, Max decimal.Decimal
		Ok       bool
	}{
		{"Bool", ComposeDataType(DataTypeMajorBool, 0), dec.False, dec.True, true},
		{"Address", ComposeDataType(DataTypeMajorAddress, 0), decimal.Zero, decAddressMax, true},
		{"Int8", ComposeDataType(DataTypeMajorInt, 0), decimal.New(-128, 0), decimal.New(127, 0), true},
		{"Int16", ComposeDataType(DataTypeMajorInt, 1), decimal.New(-32768, 0), decimal.New(32767, 0), true},
		{"UInt8", ComposeDataType(DataTypeMajorUint, 0), decimal.Zero, decimal.New(255, 0), true},
		{"UInt16", ComposeDataType(DataTypeMajorUint, 1), decimal.Zero, decimal.New(65535, 0), true},
		{"Bytes1", ComposeDataType(DataTypeMajorFixedBytes, 0), decimal.Zero, decimal.New(255, 0), true},
		{"Bytes2", ComposeDataType(DataTypeMajorFixedBytes, 1), decimal.Zero, decimal.New(65535, 0), true},
		{"Dynamic Bytes", ComposeDataType(DataTypeMajorDynamicBytes, 0), decimal.Zero, decimal.Zero, false},
	}

	var (
		min, max decimal.Decimal
		ok       bool
	)
	for _, t := range testcases {
		min, max, ok = t.In.GetMinMax()
		s.Require().Equal(t.Ok, ok, "Case: %v. Ok not equal: %v != %v", t.Name, t.Ok, ok)
		if !ok {
			continue
		}

		s.Require().True(t.Min.Equal(min), "Case: %v. Min not equal: %v != %v", t.Name, t.Min, min)
		s.Require().True(t.Max.Equal(max), "Case: %v. Max not equal: %v != %v", t.Name, t.Max, max)
	}
}

func (s *TypesTestSuite) TestDataTypeSize() {
	testcases := []struct {
		Name string
		Dt   DataType
		Size uint8
	}{
		{"Bool", ComposeDataType(DataTypeMajorBool, 0), 1},
		{"Address", ComposeDataType(DataTypeMajorAddress, 0), 20},
		{"Int8", ComposeDataType(DataTypeMajorInt, 0), 1},
		{"Int16", ComposeDataType(DataTypeMajorInt, 1), 2},
		{"UInt8", ComposeDataType(DataTypeMajorUint, 0), 1},
		{"UInt16", ComposeDataType(DataTypeMajorUint, 1), 2},
		{"Bytes1", ComposeDataType(DataTypeMajorFixedBytes, 0), 1},
		{"Bytes2", ComposeDataType(DataTypeMajorFixedBytes, 1), 2},
		{"Dynamic Bytes", ComposeDataType(DataTypeMajorDynamicBytes, 0), 32},
	}
	for _, t := range testcases {
		s.Require().Equalf(t.Size, t.Dt.Size(), "Testcase %v", t.Name)
	}
}

func (s *TypesTestSuite) TestDecimalToUint64() {
	pos := decimal.New(15, 1)
	zero := decimal.Zero
	neg := decimal.New(-150, -1)
	posSmall := decimal.New(15, -2)
	negSmall := decimal.New(-15, -2)

	testcases := []struct {
		d   decimal.Decimal
		u   uint64
		err error
	}{
		{pos, 150, nil},
		{zero, 0, nil},
		{neg, 0, errors.ErrorCodeNegDecimalToUint64},
		{posSmall, 0, nil},
		{negSmall, 0, errors.ErrorCodeNegDecimalToUint64},
	}
	for i, t := range testcases {
		u, err := DecimalToUint64(t.d)
		s.Require().Equalf(t.err, err, "err not match. testcase: %v", i)
		s.Require().Equalf(t.u, u, "result not match. testcase: %v", i)
	}
}

func (s *TypesTestSuite) TestBoolValueValidity() {
	var v BoolValue
	s.Require().False(v.Valid())
	s.Require().Panics(func() { _ = v.String() })
	s.Require().Panics(func() { _ = v.NullBool() })
	v = BoolValue(1)
	s.Require().True(v.Valid())
	s.Require().Equal("TRUE", v.String())
	s.Require().Equal(sql.NullBool{Valid: true, Bool: true}, v.NullBool())
	v = BoolValue(4)
	s.Require().False(v.Valid())
	s.Require().Panics(func() { _ = v.String() })
	s.Require().Panics(func() { _ = v.NullBool() })
}

func (s *TypesTestSuite) TestBoolValueOperations() {
	and := func(v, v2 BoolValue) BoolValue {
		if v == BoolValueFalse || v2 == BoolValueFalse {
			return BoolValueFalse
		}
		if v == BoolValueUnknown || v2 == BoolValueUnknown {
			return BoolValueUnknown
		}
		// v is true.
		return v2
	}
	or := func(v, v2 BoolValue) BoolValue {
		if v == BoolValueTrue || v2 == BoolValueTrue {
			return BoolValueTrue
		}
		if v == BoolValueUnknown || v2 == BoolValueUnknown {
			return BoolValueUnknown
		}
		// v is false.
		return v2
	}
	not := func(v BoolValue) BoolValue {
		switch v {
		case BoolValueTrue:
			return BoolValueFalse
		case BoolValueFalse:
			return BoolValueTrue
		case BoolValueUnknown:
			return BoolValueUnknown
		}
		// v is invalid.
		return v
	}
	values := [3]BoolValue{BoolValueTrue, BoolValueFalse, BoolValueUnknown}
	for _, v := range values {
		for _, v2 := range values {
			expected := and(v, v2)
			actual := v.And(v2)
			s.Require().Equalf(expected, actual,
				"%v AND %v = %v, but %v is returned", v, v2, expected, actual)
		}
	}
	for _, v := range values {
		for _, v2 := range values {
			expected := or(v, v2)
			actual := v.Or(v2)
			s.Require().Equalf(expected, actual,
				"%v OR %v = %v, but %v is returned", v, v2, expected, actual)
		}
	}
	for _, v := range values {
		expected := not(v)
		actual := v.Not()
		s.Require().Equalf(expected, actual,
			"NOT %v = %v, but %v is returned", v, expected, actual)
	}
}

func TestTypes(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}
