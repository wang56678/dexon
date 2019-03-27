package ast

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/dexon-foundation/dexon/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

type TypesTestSuite struct{ suite.Suite }

func (s *TypesTestSuite) requireEncodeAndDecodeDecimalNoError(
	d DataType, t decimal.Decimal, bs int) {
	encode, err := DecimalEncode(d, t)
	s.Require().NoError(err)
	s.Require().Len(encode, bs)
	decode, err := DecimalDecode(d, encode)
	s.Require().NoError(err)
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
	s.requireDecodeDataTypeError(DataTypeUnknown)
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

	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorAddress, 0),
		pos,
		20)
	s.requireEncodeAndDecodeDecimalNoError(
		ComposeDataType(DataTypeMajorAddress, 0),
		zero,
		20)

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

func (s *TypesTestSuite) TestGetMinMax() {
	decAddressMax := decimal.New(2, 0).Pow(decimal.New(common.AddressLength*8, 0)).Sub(dec.One)
	testcases := []struct {
		Name     string
		In       DataType
		Min, Max decimal.Decimal
		Err      error
	}{
		{"Bool", ComposeDataType(DataTypeMajorBool, 0), dec.False, dec.True, nil},
		{"Address", ComposeDataType(DataTypeMajorAddress, 0), decimal.Zero, decAddressMax, nil},
		{"Int8", ComposeDataType(DataTypeMajorInt, 0), decimal.New(-128, 0), decimal.New(127, 0), nil},
		{"Int16", ComposeDataType(DataTypeMajorInt, 1), decimal.New(-32768, 0), decimal.New(32767, 0), nil},
		{"UInt8", ComposeDataType(DataTypeMajorUint, 0), decimal.Zero, decimal.New(255, 0), nil},
		{"UInt16", ComposeDataType(DataTypeMajorUint, 1), decimal.Zero, decimal.New(65535, 0), nil},
		{"Bytes1", ComposeDataType(DataTypeMajorFixedBytes, 0), decimal.Zero, decimal.New(255, 0), nil},
		{"Bytes2", ComposeDataType(DataTypeMajorFixedBytes, 1), decimal.Zero, decimal.New(65535, 0), nil},
		{"Dynamic Bytes", ComposeDataType(DataTypeMajorDynamicBytes, 0), decimal.Zero, decimal.Zero, errors.ErrorCodeGetMinMax},
	}

	var (
		min, max decimal.Decimal
		err      error
	)
	for _, t := range testcases {
		min, max, err = GetMinMax(t.In)
		s.Require().Equal(t.Err, err, "Case: %v. Error not equal: %v != %v", t.Name, t.Err, err)
		if t.Err != nil {
			continue
		}

		s.Require().True(t.Min.Equal(min), "Case: %v. Min not equal: %v != %v", t.Name, t.Min, min)
		s.Require().True(t.Max.Equal(max), "Case: %v. Max not equal: %v != %v", t.Name, t.Max, max)
	}
}

func (s *TypesTestSuite) TestSize() {
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

func TestTypes(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}
