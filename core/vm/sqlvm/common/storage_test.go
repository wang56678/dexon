package common

import (
	"bytes"
	"fmt"
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

func (s *StorageTestSuite) TestGetPrimaryKeyHash() {
	id := uint64(555666)
	table := []byte("TABLE_A")
	key := [][]byte{
		[]byte("tables"),
		table,
		[]byte("primary"),
		convertIDtoBytes(id),
	}
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, key)
	bytes := hw.Sum(nil)
	storage := Storage{}
	result := storage.GetPrimaryKeyHash(table, id)
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
	SetDataToStateDB(head, storage, address, testcase)
	for i, t := range testcase {
		slot := storage.ShiftHashUint64(head, uint64(i))
		result := storage.DecodeDByteBySlot(address, slot)
		s.Require().Truef(bytes.Equal(result, t.result), fmt.Sprintf("name %v", t.name))
	}
}

func SetDataToStateDB(head common.Hash, storage Storage, addr common.Address,
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
	storage.Commit(false)
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
