package runtime

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"

	"github.com/dexon-foundation/decimal"
	"github.com/stretchr/testify/suite"

	dexCommon "github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/state"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
	"github.com/dexon-foundation/dexon/crypto"
	"github.com/dexon-foundation/dexon/ethdb"
)

type opLoadSuite struct {
	suite.Suite
	ctx      *common.Context
	headHash dexCommon.Hash
	address  dexCommon.Address
	slotHash []dexCommon.Hash
	raws     []*raw
}

type raw struct {
	Raw
	slotShift uint8
	byteShift uint8
	major     ast.DataTypeMajor
	minor     ast.DataTypeMinor
}

func createSchema(storage *common.Storage, raws []*raw) {
	storage.Schema = schema.Schema{
		schema.Table{
			Name: []byte("Table_A"),
		},
		schema.Table{
			Name:    []byte("Table_B"),
			Columns: make([]schema.Column, len(raws)),
		},
		schema.Table{
			Name: []byte("Table_C"),
		},
	}
	for i := range raws {
		storage.Schema[1].Columns[i] = schema.NewColumn(
			[]byte{byte(i)},
			ast.ComposeDataType(raws[i].major, raws[i].minor),
			0, nil, 0,
		)
	}
	storage.Schema.SetupColumnOffset()
}

// setSlotDataInStateDB store data in StateDB, and
// return corresponding slot hash and raw slice.
func setSlotDataInStateDB(head dexCommon.Hash, addr dexCommon.Address,
	storage *common.Storage) ([]dexCommon.Hash, []*raw) {

	hash := dexCommon.Hash{}
	var b []byte
	slotHash := []string{
		"0123112233445566778800000000000000000000000000000000000000000000",
		"48656c6c6f2c20776f726c64210000000000000000000000000000000000001a",
		"3132333435363738393000000000000000000000000000000000000000000000",
		"53514c564d2069732075736566756c2100000000000000000000000000000020",
		"0000000000000000000000000000000000000000000000000000000000000041",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	}
	uInt256Dt := ast.ComposeDataType(ast.DataTypeMajorUint, ast.DataTypeMinor(31))

	raws := []*raw{
		{
			Raw: Raw{
				Value: decimal.New(0x0123, 0),
				Bytes: nil,
			},
			slotShift: 0,
			byteShift: 0,
			major:     ast.DataTypeMajorUint,
			minor:     ast.DataTypeMinor(1),
		},
		{
			Raw: Raw{
				Value: decimal.New(0x1122334455667788, 0),
				Bytes: nil,
			},
			slotShift: 0,
			byteShift: 2,
			major:     ast.DataTypeMajorUint,
			minor:     ast.DataTypeMinor(7),
		},
		{
			Raw: Raw{
				Value: dec.False,
				Bytes: nil,
			},
			slotShift: 0,
			byteShift: 10,
			major:     ast.DataTypeMajorBool,
			minor:     ast.DataTypeMinor(0),
		},
		{
			Raw: Raw{
				Bytes: []byte("Hello, world!"),
			},
			slotShift: 1,
			byteShift: 0,
			major:     ast.DataTypeMajorDynamicBytes,
			minor:     ast.DataTypeMinor(0),
		},
		{
			Raw: Raw{
				Bytes: hexToBytes(slotHash[2][:20]),
			},
			slotShift: 2,
			byteShift: 0,
			major:     ast.DataTypeMajorFixedBytes,
			minor:     ast.DataTypeMinor(9),
		},
		{
			Raw: Raw{
				Bytes: []byte("SQLVM is useful!"),
			},
			slotShift: 3,
			byteShift: 0,
			major:     ast.DataTypeMajorDynamicBytes,
			minor:     ast.DataTypeMinor(0),
		},
		{
			Raw: Raw{
				Bytes: []byte("Hello world. Hello DEXON, SQLVM."),
			},
			slotShift: 4,
			byteShift: 0,
			major:     ast.DataTypeMajorDynamicBytes,
			minor:     ast.DataTypeMinor(0),
		},
		{
			Raw: Raw{
				Value: hexToDec(slotHash[5], uInt256Dt),
				Bytes: nil,
			},
			slotShift: 5,
			byteShift: 0,
			major:     ast.DataTypeMajorUint,
			minor:     ast.DataTypeMinor(31),
		},
	}

	// set slot hash
	hData := make([]dexCommon.Hash, len(slotHash))
	ptr := head
	for i, s := range slotHash {
		b, _ = hex.DecodeString(s)
		hData[i].SetBytes(b)
		storage.SetState(addr, ptr, hData[i])
		ptr = storage.ShiftHashUint64(ptr, uint64(1))
	}

	// set dynamic bytes data
	longDBytesLoc := 6
	longRaw := raws[longDBytesLoc]
	hash.SetBytes(longRaw.Bytes)
	ptr = storage.ShiftHashUint64(head, uint64(longRaw.slotShift))
	ptr = crypto.Keccak256Hash(ptr.Bytes())
	storage.SetState(addr, ptr, hash)

	return hData, raws
}

func hexToDec(s string, dt ast.DataType) decimal.Decimal {
	b, _ := hex.DecodeString(s)
	d, _ := ast.DecimalDecode(dt, b)
	return d
}

func hexToBytes(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

type opLoadTestCase struct {
	title          string
	outputIdx      uint
	expectedOutput *Operand
	expectedErr    error
	ids            []uint64
	fields         []uint8
	tableRef       schema.TableRef
}

func (s *opLoadSuite) SetupTest() {
	s.ctx = &common.Context{}
	s.ctx.Storage = newStorage()
	targetTableRef := schema.TableRef(1)
	s.headHash = s.ctx.Storage.GetRowPathHash(targetTableRef, uint64(123456))
	s.address = dexCommon.HexToAddress("0x6655")
	s.ctx.Storage.CreateAccount(s.address)
	s.ctx.Contract = vm.NewContract(vm.AccountRef(s.address),
		vm.AccountRef(s.address), new(big.Int), 0)
	s.slotHash, s.raws = setSlotDataInStateDB(s.headHash, s.address, s.ctx.Storage)
	createSchema(s.ctx.Storage, s.raws)
	s.setColData(targetTableRef, 654321)
}

func (s *opLoadSuite) setColData(tableRef schema.TableRef, id uint64) {
	h := s.ctx.Storage.GetRowPathHash(tableRef, id)
	setSlotDataInStateDB(h, s.address, s.ctx.Storage)
}

func (s *opLoadSuite) getOpLoadTestCases(raws []*raw) []opLoadTestCase {
	testCases := []opLoadTestCase{
		{
			title:          "NIL_RESULT",
			outputIdx:      0,
			expectedOutput: &Operand{Meta: make([]ast.DataType, 0), Data: make([]Tuple, 0)},
			expectedErr:    nil,
			ids:            nil,
			fields:         nil,
			tableRef:       0,
		},
		{
			title:          "NOT_EXIST_TABLE",
			outputIdx:      0,
			expectedOutput: nil,
			expectedErr:    errors.ErrorCodeIndexOutOfRange,
			ids:            nil,
			fields:         nil,
			tableRef:       13,
		},
		{
			title:          "OK_CASE",
			outputIdx:      0,
			expectedOutput: s.getOKCaseOutput(raws),
			expectedErr:    nil,
			ids:            []uint64{123456, 654321},
			fields:         s.getOKCaseFields(raws),
			tableRef:       1,
		},
	}
	return testCases
}

func (s *opLoadSuite) getOKCaseOutput(raws []*raw) *Operand {
	rValue := &Operand{}
	size := len(raws)
	rValue.Meta = make([]ast.DataType, size)
	rValue.Data = make([]Tuple, 2)
	for j := range rValue.Data {
		rValue.Data[j] = make([]*Raw, size)
		for i, raw := range raws {
			rValue.Meta[i] = ast.ComposeDataType(raw.major, raw.minor)
			rValue.Data[j][i] = &raw.Raw
		}
	}
	return rValue
}

func (s *opLoadSuite) getOKCaseFields(raws []*raw) []uint8 {
	rValue := make([]uint8, len(raws))
	for i := range raws {
		rValue[i] = uint8(i)
	}
	return rValue
}

type decodeTestCase struct {
	dt             ast.DataType
	expectData     *Raw
	expectSlotHash dexCommon.Hash
	shift          uint64
	inputBytes     []byte
	dBytes         []byte
}

func (s *opLoadSuite) getDecodeTestCases(headHash dexCommon.Hash,
	address dexCommon.Address, storage *common.Storage) []decodeTestCase {

	slotHash, raws := setSlotDataInStateDB(headHash, address, storage)
	createSchema(storage, raws)
	testCases := make([]decodeTestCase, len(raws))

	for i := range testCases {
		r := raws[i]
		testCases[i].dt = ast.ComposeDataType(r.major, r.minor)
		testCases[i].shift = uint64(r.slotShift)
		testCases[i].expectSlotHash = slotHash[r.slotShift]
		testCases[i].expectData = &r.Raw
		slot := slotHash[r.slotShift]
		start := r.byteShift
		end := r.byteShift + testCases[i].dt.Size()
		testCases[i].inputBytes = slot.Bytes()[start:end]
	}
	return testCases
}

func (s *opLoadSuite) newRegisters(tableRef schema.TableRef, ids []uint64, fields []uint8) []*Operand {
	o := make([]*Operand, 4)
	o[1] = newTableRefOperand(tableRef)
	o[2] = newIDsOperand(ids)
	o[3] = newFieldsOperand(fields)
	return o
}

func newInput(nums []int) []*Operand {
	o := make([]*Operand, len(nums))
	for i, n := range nums {
		o[i] = &Operand{
			IsImmediate:   false,
			RegisterIndex: uint(n),
		}
	}
	return o
}

func newTableRefOperand(tableRef schema.TableRef) *Operand {
	if tableRef < 0 {
		return nil
	}
	o := &Operand{
		Meta: []ast.DataType{
			ast.ComposeDataType(ast.DataTypeMajorUint, 0),
		},
		Data: []Tuple{
			[]*Raw{
				{
					Value: decimal.New(int64(tableRef), 0),
				},
			},
		},
	}
	return o
}

func newIDsOperand(ids []uint64) *Operand {
	o := &Operand{
		Meta: []ast.DataType{
			ast.ComposeDataType(ast.DataTypeMajorUint, 7),
		},
	}
	o.Data = make([]Tuple, len(ids))
	for i := range o.Data {
		o.Data[i] = make([]*Raw, 1)
		o.Data[i][0] = &Raw{
			Value: decimal.New(int64(ids[i]), 0),
		}
	}
	return o
}

func newFieldsOperand(fields []uint8) *Operand {
	o := &Operand{
		Meta: []ast.DataType{
			ast.ComposeDataType(ast.DataTypeMajorUint, 0),
		},
	}
	o.Data = make([]Tuple, len(fields))
	for i := range o.Data {
		o.Data[i] = make([]*Raw, 1)
		o.Data[i][0] = &Raw{
			Value: decimal.New(int64(fields[i]), 0),
		}
	}
	return o
}

func newStorage() *common.Storage {
	db := ethdb.NewMemDatabase()
	state, _ := state.New(dexCommon.Hash{}, state.NewDatabase(db))
	storage := common.NewStorage(state)
	return storage
}

func (s *opLoadSuite) TestDecode() {
	testCases := s.getDecodeTestCases(s.headHash, s.address, s.ctx.Storage)
	for _, tt := range testCases {
		M, _ := ast.DecomposeDataType(tt.dt)
		slot := s.ctx.Storage.ShiftHashUint64(s.headHash, tt.shift)
		slotHash := s.ctx.Storage.GetState(s.address, slot)
		s.Require().Equal(tt.expectSlotHash, slotHash)

		data, err := decode(s.ctx, tt.dt, slot, tt.inputBytes)
		s.Require().Nil(err)

		if M == ast.DataTypeMajorDynamicBytes {
			s.Require().Equal(tt.expectData.Bytes, data.Bytes)
		} else {
			s.Require().True(tt.expectData.Value.Equal(data.Value))
		}
	}
}

func (s *opLoadSuite) TestOpLoad() {
	testCases := s.getOpLoadTestCases(s.raws)
	for _, t := range testCases {
		input := newInput([]int{1, 2, 3})
		reg := s.newRegisters(t.tableRef, t.ids, t.fields)

		loadRegister(input, reg)
		err := opLoad(s.ctx, input, reg, t.outputIdx)

		s.Require().Equalf(t.expectedErr, err, "testcase: [%v]", t.title)
		s.Require().Truef(reflect.DeepEqual(t.expectedOutput, reg[t.outputIdx]),
			"testcase: [%v], expect: %+v, result: %+v",
			t.title, t.expectedOutput, reg[t.outputIdx],
		)
	}
}

type opRepeatPKSuite struct{ suite.Suite }

type repeatPKTestCase struct {
	tableRef    schema.TableRef
	address     dexCommon.Address
	title       string
	expectedErr error
	expectedIDs []uint64
}

func (s opRepeatPKSuite) newInput(tableRef schema.TableRef) []*Operand {
	o := make([]*Operand, 1)
	o[0] = newTableRefOperand(tableRef)
	return o
}

func (s opRepeatPKSuite) newRegisters(tableRef schema.TableRef) []*Operand {
	o := make([]*Operand, 2)
	o[1] = newTableRefOperand(tableRef)
	return o
}

func (s opRepeatPKSuite) getTestCases(storage *common.Storage) []repeatPKTestCase {
	testCases := []repeatPKTestCase{
		{
			tableRef:    0,
			address:     dexCommon.BytesToAddress([]byte("0")),
			title:       "no IDs",
			expectedErr: nil,
			expectedIDs: []uint64{},
		},
		{
			tableRef:    1,
			address:     dexCommon.BytesToAddress([]byte("1")),
			title:       "ok case",
			expectedErr: nil,
			expectedIDs: []uint64{1, 2, 3, 4},
		},
	}
	for _, t := range testCases {
		headerSlot := storage.GetPrimaryPathHash(t.tableRef)
		storage.SetPK(t.address, headerSlot, t.expectedIDs)
	}
	return testCases
}

func (s opRepeatPKSuite) TestRepeatPK() {
	ctx := &common.Context{}
	ctx.Storage = newStorage()
	testCases := s.getTestCases(ctx.Storage)
	for _, t := range testCases {
		address := t.address
		ctx.Contract = vm.NewContract(vm.AccountRef(address),
			vm.AccountRef(address), new(big.Int), uint64(0))
		reg := s.newRegisters(t.tableRef)
		input := newInput([]int{1})
		loadRegister(input, reg)
		err := opRepeatPK(ctx, input, reg, 0)
		s.Require().Equalf(t.expectedErr, err, "testcase: [%v]", t.title)
		result, _ := reg[0].toUint64()
		s.Require().Equalf(t.expectedIDs, result, "testcase: [%v]", t.title)
	}
}

func TestInstructions(t *testing.T) {
	suite.Run(t, new(opLoadSuite))
	suite.Run(t, new(opRepeatPKSuite))
	suite.Run(t, new(instructionSuite))
}

func makeOperand(im bool, meta []ast.DataType, pTuple []Tuple) (op *Operand) {
	op = &Operand{IsImmediate: im, Meta: meta, Data: pTuple}
	return
}

func loadRegister(input, registers []*Operand) {
	for i, operand := range input {
		if operand != nil && !operand.IsImmediate {
			input[i] = registers[operand.RegisterIndex]
		}
	}
}

type opTestcase struct {
	Name   string
	In     Instruction
	Output *Operand
	Err    error
}

type instructionSuite struct {
	suite.Suite
}

func (s *instructionSuite) run(testcases []opTestcase, opfunc OpFunction) {
	for idx, c := range testcases {
		registers := make([]*Operand, len(c.In.Input))

		for i, j := 0, 0; i < len(c.In.Input); i++ {
			if !c.In.Input[i].IsImmediate {
				registers[j] = c.In.Input[i]
				j++
			}
		}
		err := opfunc(
			&common.Context{Opt: common.Option{SafeMath: true}},
			c.In.Input, registers, c.In.Output)
		s.Require().Equal(
			c.Err, err,
			"idx: %v, op: %v, case: %v\nerror not equal: %v != %v",
			idx, c.In.Op, c.Name, c.Err, err,
		)
		if c.Err != nil {
			continue
		}

		result := registers[0]
		s.Require().True(
			c.Output.Equal(result),
			"idx: %v, op: %v, case: %v\noutput not equal.\nExpect: %v\nResult: %v\n",
			idx, c.In.Op, c.Name, c.Output, result,
		)
	}
}
