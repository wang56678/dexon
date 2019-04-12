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
	"github.com/dexon-foundation/dexon/crypto"
	"github.com/dexon-foundation/dexon/ethdb"
	"github.com/dexon-foundation/dexon/rlp"
)

type StorageTestSuite struct{ suite.Suite }

func (s *StorageTestSuite) TestUint64ToBytes() {
	testcases := []uint64{1, 65535, math.MaxUint64}
	for _, i := range testcases {
		s.Require().Equal(i, bytesToUint64(uint64ToBytes(i)))
	}
}

func (s *StorageTestSuite) TestGetRowAddress() {
	id := uint64(555666)
	table := []byte("TABLE_A")
	key := [][]byte{
		[]byte("tables"),
		table,
		[]byte("primary"),
		uint64ToBytes(id),
	}
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, key)
	bytes := hw.Sum(nil)
	storage := &Storage{}
	result := storage.GetRowPathHash(table, id)
	s.Require().Equal(bytes, result[:])
}

type decodeTestCase struct {
	name     string
	slotData common.Hash
	result   []byte
}

func (s *StorageTestSuite) TestDecodeDByte() {
	db := ethdb.NewMemDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	storage := NewStorage(state)
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
	SetDataToStorage(head, storage, address, testcase)
	for i, t := range testcase {
		slot := storage.ShiftHashUint64(head, uint64(i))
		result := storage.DecodeDByteBySlot(address, slot)
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
	db := ethdb.NewMemDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	storage := NewStorage(state)

	contractA := common.BytesToAddress([]byte("I'm sad."))
	ownerA := common.BytesToAddress([]byte{5, 5, 6, 6})
	contractB := common.BytesToAddress([]byte{9, 5, 2, 7})
	ownerB := common.BytesToAddress([]byte("Tong Pak-Fu"))

	storage.StoreOwner(contractA, ownerA)
	storage.StoreOwner(contractB, ownerB)
	s.Require().Equal(ownerA, storage.LoadOwner(contractA))
	s.Require().Equal(ownerB, storage.LoadOwner(contractB))

	storage.StoreOwner(contractA, ownerB)
	s.Require().Equal(ownerB, storage.LoadOwner(contractA))
}

func (s *StorageTestSuite) TestTableWriter() {
	db := ethdb.NewMemDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	storage := NewStorage(state)

	table1 := []byte("table1")
	table2 := []byte("table2")
	contractA := common.BytesToAddress([]byte("A"))
	contractB := common.BytesToAddress([]byte("B"))
	addrs := []common.Address{
		common.BytesToAddress([]byte("addr1")),
		common.BytesToAddress([]byte("addr2")),
		common.BytesToAddress([]byte("addr3")),
	}

	// Genesis.
	s.Require().Len(storage.LoadTableWriters(contractA, table1), 0)
	s.Require().Len(storage.LoadTableWriters(contractB, table1), 0)

	// Check writer list.
	storage.InsertTableWriter(contractA, table1, addrs[0])
	storage.InsertTableWriter(contractA, table1, addrs[1])
	storage.InsertTableWriter(contractA, table1, addrs[2])
	storage.InsertTableWriter(contractB, table2, addrs[0])
	s.Require().Equal(addrs, storage.LoadTableWriters(contractA, table1))
	s.Require().Len(storage.LoadTableWriters(contractA, table2), 0)
	s.Require().Len(storage.LoadTableWriters(contractB, table1), 0)
	s.Require().Equal([]common.Address{addrs[0]},
		storage.LoadTableWriters(contractB, table2))

	// Insert duplicate.
	storage.InsertTableWriter(contractA, table1, addrs[0])
	storage.InsertTableWriter(contractA, table1, addrs[1])
	storage.InsertTableWriter(contractA, table1, addrs[2])
	s.Require().Equal(addrs, storage.LoadTableWriters(contractA, table1))

	// Delete some writer.
	storage.DeleteTableWriter(contractA, table1, addrs[0])
	storage.DeleteTableWriter(contractA, table2, addrs[0])
	storage.DeleteTableWriter(contractB, table2, addrs[0])
	s.Require().Equal([]common.Address{addrs[2], addrs[1]},
		storage.LoadTableWriters(contractA, table1))
	s.Require().Len(storage.LoadTableWriters(contractA, table2), 0)
	s.Require().Len(storage.LoadTableWriters(contractB, table1), 0)
	s.Require().Len(storage.LoadTableWriters(contractB, table2), 0)

	// Delete again.
	storage.DeleteTableWriter(contractA, table1, addrs[2])
	s.Require().Equal([]common.Address{addrs[1]},
		storage.LoadTableWriters(contractA, table1))

	// Check writer.
	s.Require().False(storage.IsTableWriter(contractA, table1, addrs[0]))
	s.Require().True(storage.IsTableWriter(contractA, table1, addrs[1]))
	s.Require().False(storage.IsTableWriter(contractA, table1, addrs[2]))
	s.Require().False(storage.IsTableWriter(contractA, table2, addrs[0]))
	s.Require().False(storage.IsTableWriter(contractB, table2, addrs[0]))
}

func (s *StorageTestSuite) TestSequence() {
	db := ethdb.NewMemDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	storage := NewStorage(state)

	table1 := []byte("table1")
	table2 := []byte("table2")
	contract := common.BytesToAddress([]byte("A"))

	s.Require().Equal(uint64(0), storage.IncSequence(contract, table1, 0, 2))
	s.Require().Equal(uint64(2), storage.IncSequence(contract, table1, 0, 1))
	s.Require().Equal(uint64(3), storage.IncSequence(contract, table1, 0, 1))
	// Repeat on another sequence.
	s.Require().Equal(uint64(0), storage.IncSequence(contract, table1, 1, 1))
	s.Require().Equal(uint64(1), storage.IncSequence(contract, table1, 1, 2))
	s.Require().Equal(uint64(3), storage.IncSequence(contract, table1, 1, 3))
	// Repeat on another table.
	s.Require().Equal(uint64(0), storage.IncSequence(contract, table2, 0, 3))
	s.Require().Equal(uint64(3), storage.IncSequence(contract, table2, 0, 4))
	s.Require().Equal(uint64(7), storage.IncSequence(contract, table2, 0, 5))
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
