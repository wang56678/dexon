package common

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/sha3"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/state"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
	"github.com/dexon-foundation/dexon/crypto"
	"github.com/dexon-foundation/dexon/ethdb"
	"github.com/dexon-foundation/dexon/rlp"
)

type StorageTestSuite struct {
	suite.Suite
	storage *Storage
	address common.Address
}

func (s *StorageTestSuite) SetupTest() {
	db := ethdb.NewMemDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	s.storage = NewStorage(state)
	s.address = common.BytesToAddress([]byte("5566"))
}

func (s *StorageTestSuite) TestUint64ToBytes() {
	testcases := []uint64{1, 65535, math.MaxUint64}
	for _, i := range testcases {
		s.Require().Equal(i, bytesToUint64(uint64ToBytes(i)))
	}
}

func (s *StorageTestSuite) TestGetRowAddress() {
	id := uint64(555666)
	table := schema.TableRef(1)
	key := [][]byte{
		[]byte("tables"),
		{uint8(table)},
		[]byte("primary"),
		uint64ToBytes(id),
	}
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, key)
	bytes := hw.Sum(nil)
	result := s.storage.GetRowPathHash(table, id)
	s.Require().Equal(bytes, result[:])
}

type decodeTestCase struct {
	name     string
	slotData common.Hash
	result   []byte
}

func (s *StorageTestSuite) TestDecodeDByte() {
	address := common.BytesToAddress([]byte("123"))
	head := common.HexToHash("0x5566")
	testcase := []decodeTestCase{
		{
			name:     "small size",
			slotData: common.HexToHash("0x48656c6c6f2c20776f726c64210000000000000000000000000000000000001a"),
			result:   common.FromHex("0x48656c6c6f2c20776f726c6421"),
		},
		{
			name:     "32 byte case",
			slotData: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000041"),
			result:   []byte("Hello world. Hello DEXON, SQLVM."),
		},
		{
			name:     "large size",
			slotData: common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000047D"),
			result:   []byte("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."),
		},
		{
			name:     "empty",
			slotData: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			result:   []byte(""),
		},
	}
	SetDataToStorage(head, s.storage, address, testcase)
	for i, t := range testcase {
		slot := s.storage.ShiftHashUint64(head, uint64(i))
		result := s.storage.DecodeDByteBySlot(address, slot)
		s.Require().Truef(bytes.Equal(result, t.result), fmt.Sprintf("name %v", t.name))
	}
}

func SetDataToStorage(head common.Hash, storage *Storage, addr common.Address,
	testcase []decodeTestCase) {
	for i, t := range testcase {
		slot := storage.ShiftHashUint64(head, uint64(i))
		storage.SetState(addr, slot, t.slotData)
		b := t.slotData.Bytes()
		if b[len(b)-1]&0x1 != 0 {
			length := len(t.result)
			slotNum := (length-1)/common.HashLength + 1
			ptr := crypto.Keccak256Hash(slot.Bytes())
			for s := 0; s < slotNum; s++ {
				start := s * common.HashLength
				end := (s + 1) * common.HashLength
				if end > len(t.result) {
					end = len(t.result)
				}
				hash := common.Hash{}
				copy(hash[:], t.result[start:end])
				storage.SetState(addr, ptr, hash)
				ptr = storage.ShiftHashUint64(ptr, 1)
			}
		}
	}
}

func (s *StorageTestSuite) TestOwner() {
	contractA := common.BytesToAddress([]byte("I'm sad."))
	ownerA := common.BytesToAddress([]byte{5, 5, 6, 6})
	contractB := common.BytesToAddress([]byte{9, 5, 2, 7})
	ownerB := common.BytesToAddress([]byte("Tong Pak-Fu"))

	s.storage.StoreOwner(contractA, ownerA)
	s.storage.StoreOwner(contractB, ownerB)
	s.Require().Equal(ownerA, s.storage.LoadOwner(contractA))
	s.Require().Equal(ownerB, s.storage.LoadOwner(contractB))

	s.storage.StoreOwner(contractA, ownerB)
	s.Require().Equal(ownerB, s.storage.LoadOwner(contractA))
}

func (s *StorageTestSuite) TestTableWriter() {
	table1 := schema.TableRef(0)
	table2 := schema.TableRef(1)
	contractA := common.BytesToAddress([]byte("A"))
	contractB := common.BytesToAddress([]byte("B"))
	addrs := []common.Address{
		common.BytesToAddress([]byte("addr1")),
		common.BytesToAddress([]byte("addr2")),
		common.BytesToAddress([]byte("addr3")),
	}

	// Genesis.
	s.Require().Len(s.storage.LoadTableWriters(contractA, table1), 0)
	s.Require().Len(s.storage.LoadTableWriters(contractB, table1), 0)

	// Check writer list.
	s.storage.InsertTableWriter(contractA, table1, addrs[0])
	s.storage.InsertTableWriter(contractA, table1, addrs[1])
	s.storage.InsertTableWriter(contractA, table1, addrs[2])
	s.storage.InsertTableWriter(contractB, table2, addrs[0])
	s.Require().Equal(addrs, s.storage.LoadTableWriters(contractA, table1))
	s.Require().Len(s.storage.LoadTableWriters(contractA, table2), 0)
	s.Require().Len(s.storage.LoadTableWriters(contractB, table1), 0)
	s.Require().Equal([]common.Address{addrs[0]},
		s.storage.LoadTableWriters(contractB, table2))

	// Insert duplicate.
	s.storage.InsertTableWriter(contractA, table1, addrs[0])
	s.storage.InsertTableWriter(contractA, table1, addrs[1])
	s.storage.InsertTableWriter(contractA, table1, addrs[2])
	s.Require().Equal(addrs, s.storage.LoadTableWriters(contractA, table1))

	// Delete some writer.
	s.storage.DeleteTableWriter(contractA, table1, addrs[0])
	s.storage.DeleteTableWriter(contractA, table2, addrs[0])
	s.storage.DeleteTableWriter(contractB, table2, addrs[0])
	s.Require().Equal([]common.Address{addrs[2], addrs[1]},
		s.storage.LoadTableWriters(contractA, table1))
	s.Require().Len(s.storage.LoadTableWriters(contractA, table2), 0)
	s.Require().Len(s.storage.LoadTableWriters(contractB, table1), 0)
	s.Require().Len(s.storage.LoadTableWriters(contractB, table2), 0)

	// Delete again.
	s.storage.DeleteTableWriter(contractA, table1, addrs[2])
	s.Require().Equal([]common.Address{addrs[1]},
		s.storage.LoadTableWriters(contractA, table1))

	// Check writer.
	s.Require().False(s.storage.IsTableWriter(contractA, table1, addrs[0]))
	s.Require().True(s.storage.IsTableWriter(contractA, table1, addrs[1]))
	s.Require().False(s.storage.IsTableWriter(contractA, table1, addrs[2]))
	s.Require().False(s.storage.IsTableWriter(contractA, table2, addrs[0]))
	s.Require().False(s.storage.IsTableWriter(contractB, table2, addrs[0]))
}

func (s *StorageTestSuite) TestSequence() {
	table1 := schema.TableRef(0)
	table2 := schema.TableRef(1)
	contract := common.BytesToAddress([]byte("A"))

	s.Require().Equal(uint64(0), s.storage.IncSequence(contract, table1, 0, 2))
	s.Require().Equal(uint64(2), s.storage.IncSequence(contract, table1, 0, 1))
	s.Require().Equal(uint64(3), s.storage.IncSequence(contract, table1, 0, 1))
	// Repeat on another sequence.
	s.Require().Equal(uint64(0), s.storage.IncSequence(contract, table1, 1, 1))
	s.Require().Equal(uint64(1), s.storage.IncSequence(contract, table1, 1, 2))
	s.Require().Equal(uint64(3), s.storage.IncSequence(contract, table1, 1, 3))
	// Repeat on another table.
	s.Require().Equal(uint64(0), s.storage.IncSequence(contract, table2, 0, 3))
	s.Require().Equal(uint64(3), s.storage.IncSequence(contract, table2, 0, 4))
	s.Require().Equal(uint64(7), s.storage.IncSequence(contract, table2, 0, 5))
}

func (s *StorageTestSuite) TestPKHeaderEncodeDecode() {
	lastRowID := uint64(5566)
	rowCount := uint64(6655)
	bm := bitMap{}
	bm.encodeHeader(lastRowID, rowCount)
	newLastRowID, newRowCount := bm.decodeHeader()
	s.Require().Equal(lastRowID, newLastRowID)
	s.Require().Equal(rowCount, newRowCount)
}

func (s *StorageTestSuite) TestRepeatPK() {
	type testCase struct {
		address   common.Address
		tableRef  schema.TableRef
		expectIDs []uint64
	}
	testCases := []testCase{
		{
			address:   common.BytesToAddress([]byte("0")),
			tableRef:  schema.TableRef(0),
			expectIDs: []uint64{0, 1, 2},
		},
		{
			address:   common.BytesToAddress([]byte("1")),
			tableRef:  schema.TableRef(1),
			expectIDs: []uint64{1234, 5566},
		},
		{
			address:   common.BytesToAddress([]byte("2")),
			tableRef:  schema.TableRef(2),
			expectIDs: []uint64{0, 128, 256, 512, 1024},
		},
	}
	for i, t := range testCases {
		headerSlot := s.storage.GetPrimaryPathHash(t.tableRef)
		s.storage.SetPK(t.address, headerSlot, t.expectIDs)
		IDs := s.storage.RepeatPK(t.address, t.tableRef)
		s.Require().Equalf(t.expectIDs, IDs, "testCase #%v\n", i)
	}
}

func (s *StorageTestSuite) TestBitMapIncreasePK() {
	type testCase struct {
		tableRef schema.TableRef
		IDs      []uint64
	}
	testCases := []testCase{
		{
			tableRef: schema.TableRef(0),
			IDs:      []uint64{0, 1, 2},
		},
		{
			tableRef: schema.TableRef(1),
			IDs:      []uint64{1234, 5566},
		},
		{
			tableRef: schema.TableRef(2),
			IDs:      []uint64{0, 128, 256, 512, 1024},
		},
	}
	for i, t := range testCases {
		hash := s.storage.GetPrimaryPathHash(t.tableRef)
		s.storage.SetPK(s.address, hash, t.IDs)
		bm := newBitMap(hash, s.address, s.storage)
		newID := bm.increasePK()

		t.IDs = append(t.IDs, newID)
		IDs := s.storage.RepeatPK(s.address, t.tableRef)
		s.Require().Equalf(t.IDs, IDs, "testCase #%v\n", i)
	}
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
