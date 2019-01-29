package ast

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TypeTestSuite struct{ suite.Suite }

func (s *TypeTestSuite) requireEncodeAndDecodeNoError(
	d DataType, t interface{}) {
	encode, err := DataTypeEncode(t)
	s.Require().NoError(err)
	s.Require().Equal(d, encode)
	decode, err := DataTypeDecode(d)
	s.Require().NoError(err)
	s.Require().Equal(t, decode)
}

func (s *TypeTestSuite) requireEncodeError(input interface{}) {
	_, err := DataTypeEncode(input)
	s.Require().Error(err)
}

func (s *TypeTestSuite) requireDecodeError(input DataType) {
	_, err := DataTypeDecode(input)
	s.Require().Error(err)
}

func (s *TypeTestSuite) TestEncodeAndDecode() {
	s.requireEncodeAndDecodeNoError(
		composeDataType(DataTypeMajorBool, 0),
		BoolTypeNode{})
	s.requireEncodeAndDecodeNoError(
		composeDataType(DataTypeMajorAddress, 0),
		AddressTypeNode{})
	s.requireEncodeAndDecodeNoError(
		composeDataType(DataTypeMajorInt, 1),
		IntTypeNode{Size: 16})
	s.requireEncodeAndDecodeNoError(
		composeDataType(DataTypeMajorUint, 2),
		IntTypeNode{Unsigned: true, Size: 24})
	s.requireEncodeAndDecodeNoError(
		composeDataType(DataTypeMajorFixedBytes, 3),
		FixedBytesTypeNode{Size: 32})
	s.requireEncodeAndDecodeNoError(
		composeDataType(DataTypeMajorDynamicBytes, 0),
		DynamicBytesTypeNode{})
	s.requireEncodeAndDecodeNoError(
		composeDataType(DataTypeMajorFixed, 1),
		FixedTypeNode{Size: 8, FractionalDigits: 1})
	s.requireEncodeAndDecodeNoError(
		composeDataType(DataTypeMajorUfixed+1, 2),
		FixedTypeNode{Unsigned: true, Size: 16, FractionalDigits: 2})
}

func (s *TypeTestSuite) TestEncodeError() {
	s.requireEncodeError(struct{}{})
	s.requireEncodeError(IntTypeNode{Size: 1})
	s.requireEncodeError(IntTypeNode{Size: 257})
	s.requireEncodeError(FixedBytesTypeNode{Size: 1})
	s.requireEncodeError(FixedBytesTypeNode{Size: 257})
	s.requireEncodeError(FixedTypeNode{Size: 1, FractionalDigits: 0})
	s.requireEncodeError(FixedTypeNode{Size: 257, FractionalDigits: 0})
	s.requireEncodeError(FixedTypeNode{Size: 8, FractionalDigits: 81})
}

func (s *TypeTestSuite) TestDecodeError() {
	s.requireDecodeError(DataTypeUnknown)
	s.requireDecodeError(composeDataType(DataTypeMajorBool, 1))
	s.requireDecodeError(composeDataType(DataTypeMajorAddress, 1))
	s.requireDecodeError(composeDataType(DataTypeMajorInt, 0x20))
	s.requireDecodeError(composeDataType(DataTypeMajorUint, 0x20))
	s.requireDecodeError(composeDataType(DataTypeMajorFixedBytes, 0x20))
	s.requireDecodeError(composeDataType(DataTypeMajorDynamicBytes, 1))
	s.requireDecodeError(composeDataType(DataTypeMajorFixed, 81))
	s.requireDecodeError(composeDataType(DataTypeMajorUfixed, 81))
	s.requireDecodeError(composeDataType(DataTypeMajorUfixed+0x20, 80))
}

func TestType(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}
