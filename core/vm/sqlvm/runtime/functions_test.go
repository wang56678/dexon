package runtime

import (
	"math"
	"math/big"
	"testing"

	"github.com/dexon-foundation/decimal"
	"github.com/stretchr/testify/suite"

	dexCommon "github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
)

func TestFunction(t *testing.T) {
	suite.Run(t, new(FunctionSuite))
}

type FunctionSuite struct {
	suite.Suite
}

var (
	hash1   = dexCommon.BigToHash(big.NewInt(1))
	hash255 = dexCommon.BigToHash(big.NewInt(255))
)

var mockNumberHashTable = map[uint64]dexCommon.Hash{1: hash1, 255: hash255}

func mockGetHashFunc(u uint64) dexCommon.Hash { return mockNumberHashTable[u] }

func (s *FunctionSuite) TestFnBlockHash() {
	type blockHashCase struct {
		Name   string
		Ops    []*Operand
		Length uint64
		Res    [][]byte
		Cur    *big.Int
		Err    error
	}

	testcases := []blockHashCase{
		{"Immediate OP", []*Operand{
			{IsImmediate: true, Meta: nil, Data: []Tuple{{&Raw{Value: decimal.New(1, 0)}}}},
		}, 2, [][]byte{hash1.Bytes(), hash1.Bytes()}, big.NewInt(255), nil},
		{"OP", []*Operand{
			{IsImmediate: false, Meta: nil, Data: []Tuple{
				{&Raw{Value: decimal.New(255, 0)}},
				{&Raw{Value: decimal.New(515, 0)}},
			}},
		}, 2, [][]byte{hash255.Bytes(), make([]byte, 32)}, big.NewInt(256), nil},
		{"Older than 257 block", []*Operand{
			{IsImmediate: false, Meta: nil, Data: []Tuple{
				{&Raw{Value: decimal.New(1, 0)}},
			}},
		}, 1, [][]byte{make([]byte, 32)}, big.NewInt(512), nil},
	}

	callFn := func(c blockHashCase) (*Operand, error) {
		return fnBlockHash(
			&common.Context{
				Context: vm.Context{
					GetHash:     mockGetHashFunc,
					BlockNumber: c.Cur,
				},
			},
			c.Ops,
			c.Length,
		)
	}

	meta := []ast.DataType{ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 3)}

	for idx, tCase := range testcases {
		r, err := callFn(tCase)
		s.Require().Equal(
			tCase.Err, err,
			"Index: %v. Error not expected: %v != %v", idx, tCase.Err, err)
		s.Require().Equal(
			meta, r.Meta,
			"Index: %v. Meta not equal: %v != %v", idx, meta, r.Meta)
		s.Require().Equal(
			uint64(len(r.Data)), tCase.Length,
			"Index: %v. Length not equal: %v != %v", idx, len(r.Data), tCase.Length)

		for i := 0; i < len(r.Data); i++ {
			s.Require().Equal(
				tCase.Res[i], r.Data[i][0].Bytes,
				"TestCase Index: %v. Data Index: %v. Value not equal: %v != %v",
				idx, i, tCase.Res[i], r.Data[i][0].Bytes)
		}
	}
}

func (s *FunctionSuite) TestFnBlockNumber() {
	type blockNumberCase struct {
		Name       string
		RawNum     *big.Int
		Length     uint64
		ResNum     decimal.Decimal
		Err        error
		AsserPanic bool
	}

	testcases := []blockNumberCase{
		{"number 1 with length 1", big.NewInt(1), 1, decimal.New(1, 0), nil, false},
		{"number 10 with length 10", big.NewInt(10), 10, decimal.New(10, 0), nil, false},
		{"number 1 with length 0", big.NewInt(1), 0, decimal.New(1, 0), nil, false},
		{"panic on invalid context", nil, 0, decimal.New(1, 0), nil, true},
	}

	callFn := func(c blockNumberCase) (*Operand, error) {
		return fnBlockNumber(
			&common.Context{
				Context: vm.Context{BlockNumber: c.RawNum},
			},
			nil,
			c.Length)
	}

	meta := []ast.DataType{ast.ComposeDataType(ast.DataTypeMajorUint, 31)}

	for idx, tCase := range testcases {
		if tCase.AsserPanic {
			s.Require().Panicsf(
				func() { callFn(tCase) },
				"Index: %v. Not Panic on '%v'", idx, tCase.Name,
			)
		} else {
			r, err := callFn(tCase)
			s.Require().Equal(
				tCase.Err, err,
				"Index: %v. Error not expected: %v != %v", idx, tCase.Err, err)
			s.Require().Equal(
				meta, r.Meta,
				"Index: %v. Meta not equal: %v != %v", idx, meta, r.Meta)
			s.Require().Equal(
				uint64(len(r.Data)), tCase.Length,
				"Index: %v. Length not equal: %v != %v", idx, len(r.Data), tCase.Length)

			for i := 0; i < len(r.Data); i++ {
				s.Require().True(
					tCase.ResNum.Equal(r.Data[i][0].Value),
					"Index: %v Data index: %v. Value not equal: %v != %v",
					idx, i, tCase.ResNum, r.Data[i][0].Value)
			}
		}
	}
}

func (s *FunctionSuite) TestFnBlockTimestamp() {
	type blockTimestampCase struct {
		Name       string
		Timestamp  *big.Int
		Length     uint64
		Res        decimal.Decimal
		Err        error
		AsserPanic bool
	}

	testcases := []blockTimestampCase{
		{"number 1 with length 1", big.NewInt(1), 1, decimal.New(1, 0), nil, false},
		{"number 10 with length 10", big.NewInt(10), 10, decimal.New(10, 0), nil, false},
		{"number 1 with length 0", big.NewInt(1), 0, decimal.New(1, 0), nil, false},
		{"panic on invalid context", nil, 0, decimal.New(1, 0), nil, true},
	}

	callFn := func(c blockTimestampCase) (*Operand, error) {
		return fnBlockTimestamp(
			&common.Context{
				Context: vm.Context{Time: c.Timestamp},
			},
			nil,
			c.Length)
	}

	meta := []ast.DataType{ast.ComposeDataType(ast.DataTypeMajorUint, 31)}

	for idx, tCase := range testcases {
		if tCase.AsserPanic {
			s.Require().Panicsf(
				func() { callFn(tCase) },
				"Index: %v. Not Panic on '%v'", idx, tCase.Name,
			)
		} else {
			r, err := callFn(tCase)
			s.Require().Equal(
				tCase.Err, err,
				"Index: %v. Error not expected: %v != %v", idx, tCase.Err, err)
			s.Require().Equal(
				meta, r.Meta,
				"Index: %v. Meta not equal: %v != %v", idx, meta, r.Meta)
			s.Require().Equal(
				uint64(len(r.Data)), tCase.Length,
				"Index: %v. Length not equal: %v != %v", idx, len(r.Data), tCase.Length)

			for i := 0; i < len(r.Data); i++ {
				s.Require().True(
					tCase.Res.Equal(r.Data[i][0].Value),
					"Index: %v. Data Index: %v. Value not equal: %v != %v",
					idx, i, tCase.Res, r.Data[i][0].Value)
			}
		}
	}
}

func (s *FunctionSuite) TestFnCoinBase() {
	type blockCoinBaseCase struct {
		Name    string
		Address dexCommon.Address
		Length  uint64
		Res     []byte
		Err     error
	}

	res := []byte{
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67}
	address := dexCommon.BytesToAddress(res)

	testcases := []blockCoinBaseCase{
		{"address with length 1", address, 1, res, nil},
		{"address with length 10", address, 10, res, nil},
		{"address with length 0", address, 0, res, nil},
	}

	callFn := func(c blockCoinBaseCase) (*Operand, error) {
		return fnBlockCoinBase(
			&common.Context{
				Context: vm.Context{Coinbase: c.Address},
			},
			nil,
			c.Length)
	}

	meta := []ast.DataType{ast.ComposeDataType(ast.DataTypeMajorAddress, 0)}

	for idx, tCase := range testcases {
		r, err := callFn(tCase)
		s.Require().Equal(
			tCase.Err, err,
			"Index: %v. Error not expected: %v != %v", idx, tCase.Err, err)
		s.Require().Equal(
			meta, r.Meta,
			"Index: %v. Meta not equal: %v != %v", idx, meta, r.Meta)
		s.Require().Equal(
			uint64(len(r.Data)), tCase.Length,
			"Index: %v. Length not equal: %v != %v", idx, len(r.Data), tCase.Length)

		for i := 0; i < len(r.Data); i++ {
			s.Require().Equal(
				tCase.Res, r.Data[i][0].Bytes,
				"Index: %v. Data Index: %v. Value not equal: %v != %v",
				idx, i, tCase.Res, r.Data[i][0].Bytes)
		}
	}
}

func (s *FunctionSuite) TestFnGasLimit() {
	type blockGasLimitCase struct {
		Name   string
		Limit  uint64
		Length uint64
		Res    decimal.Decimal
		Err    error
	}
	testcases := []blockGasLimitCase{
		{"max int64 with length 1", uint64(math.MaxInt64), 1, dec.MaxInt64, nil},
		{"max uint64 with length 1", math.MaxUint64, 10, dec.MaxUint64, nil},
		{"address with length 0", math.MaxUint64, 0, decimal.Zero, nil},
	}

	callFn := func(c blockGasLimitCase) (*Operand, error) {
		return fnBlockGasLimit(
			&common.Context{
				Context: vm.Context{GasLimit: c.Limit},
			},
			nil,
			c.Length)
	}

	meta := []ast.DataType{ast.ComposeDataType(ast.DataTypeMajorUint, 7)}

	for idx, tCase := range testcases {
		r, err := callFn(tCase)
		s.Require().Equal(
			tCase.Err, err,
			"Index: %v. Error not expected: %v != %v", idx, tCase.Err, err)
		s.Require().Equal(
			meta, r.Meta,
			"Index: %v. Meta not equal: %v != %v", idx, meta, r.Meta)
		s.Require().Equal(
			uint64(len(r.Data)), tCase.Length,
			"Index: %v. Length not equal: %v != %v", idx, len(r.Data), tCase.Length)

		for i := 0; i < len(r.Data); i++ {
			s.Require().True(
				tCase.Res.Equal(r.Data[i][0].Value),
				"Index: %v. Data Index: %v. Value not equal: %v != %v",
				idx, i, tCase.Res, r.Data[i][0].Value)
		}
	}
}

func (s *FunctionSuite) TestFnMsgSender() {
	type txMsgSenderCase struct {
		Name    string
		Address dexCommon.Address
		Length  uint64
		Res     []byte
		Err     error
	}

	res := []byte{
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67}
	address := dexCommon.BytesToAddress(res)

	testcases := []txMsgSenderCase{
		{"address with length 1", address, 1, res, nil},
		{"address with length 10", address, 10, res, nil},
		{"address with length 0", address, 0, res, nil},
	}

	callFn := func(c txMsgSenderCase) (*Operand, error) {
		return fnMsgSender(
			&common.Context{
				Contract: &vm.Contract{CallerAddress: c.Address},
			},
			nil,
			c.Length)
	}

	meta := []ast.DataType{ast.ComposeDataType(ast.DataTypeMajorAddress, 0)}

	for idx, tCase := range testcases {
		r, err := callFn(tCase)
		s.Require().Equal(
			tCase.Err, err,
			"Index: %v. Error not expected: %v != %v", idx, tCase.Err, err)
		s.Require().Equal(
			meta, r.Meta,
			"Index: %v. Meta not equal: %v != %v", idx, meta, r.Meta)
		s.Require().Equal(
			uint64(len(r.Data)), tCase.Length,
			"Index: %v. Length not equal: %v != %v", idx, len(r.Data), tCase.Length)

		for i := 0; i < len(r.Data); i++ {
			s.Require().Equal(
				tCase.Res, r.Data[i][0].Bytes,
				"Index: %v. Data Index: %v. Value not equal: %v != %v",
				idx, i, tCase.Res, r.Data[i][0].Bytes)
		}
	}
}

