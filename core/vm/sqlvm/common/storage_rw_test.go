package common

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/state"
	"github.com/dexon-foundation/dexon/ethdb"
)

type StorageRWTestSuite struct{ suite.Suite }

func (s *StorageRWTestSuite) TestRW() {
	db := ethdb.NewMemDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	storage := NewStorage(state)
	contract := common.BytesToAddress([]byte("contract"))
	start := common.BytesToHash([]byte("start"))

	// Data to write.
	payload := []byte("What is Lorem Ipsum? Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.")
	// Length used in rw. (Sum should equal len(paylaod))
	writeLens := []int{5, 10, 32, 9, len(payload) - 56}
	readLens := []int{9, 32, 5, 10, len(payload) - 56}

	// Write.
	writer := NewStorageWriter(storage, contract, start)
	cursor := 0
	for _, v := range writeLens {
		n, err := writer.Write(payload[cursor:(cursor + v)])
		s.Require().Nil(err)
		s.Require().Equal(v, n)
		cursor += v
	}
	storage.Commit(false)

	// Read and check.
	reader := NewStorageReader(storage, contract, start)
	cursor = 0
	for _, v := range readLens {
		payloadRead := make([]byte, v)
		n, err := reader.Read(payloadRead)
		s.Require().Nil(err)
		s.Require().Equal(v, n)
		s.Require().Equal(payload[cursor:(cursor+v)], payloadRead)
		cursor += v
	}
	// Check if the remaining data is all zero.
	zeroCheck := make([]byte, 128)
	n, err := reader.Read(zeroCheck)
	s.Require().Nil(err)
	s.Require().Equal(len(zeroCheck), n)
	for _, v := range zeroCheck {
		s.Require().Zero(v)
	}
}

func TestStorageRW(t *testing.T) {
	suite.Run(t, new(StorageRWTestSuite))
}
