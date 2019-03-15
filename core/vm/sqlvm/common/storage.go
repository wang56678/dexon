package common

import (
	"math/big"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/state"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
	"github.com/dexon-foundation/dexon/crypto"
	"github.com/dexon-foundation/dexon/rlp"
)

// Storage holds SQLVM required data and method.
type Storage struct {
	state.StateDB
	Schema schema.Schema
}

// NewStorage return Storage instance.
func NewStorage(state *state.StateDB) Storage {
	s := Storage{*state, schema.Schema{}}
	return s
}

func convertIDtoBytes(id uint64) []byte {
	bigIntID := new(big.Int).SetUint64(id)
	decimalID := decimal.NewFromBigInt(bigIntID, 0)
	dt := ast.ComposeDataType(ast.DataTypeMajorUint, 7)
	byteID, _ := ast.DecimalEncode(dt, decimalID)
	return byteID
}

// GetPrimaryKeyHash return primary key hash.
func (s Storage) GetPrimaryKeyHash(tableName []byte, id uint64) (h common.Hash) {
	key := [][]byte{
		[]byte("tables"),
		tableName,
		[]byte("primary"),
		convertIDtoBytes(id),
	}
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, key)
	// length of common.Hash is 256bit,
	// so it can properly match the size of hw.Sum
	hw.Sum(h[:0])
	return
}

// ShiftHashUint64 shift hash in uint64.
func (s Storage) ShiftHashUint64(hash common.Hash, shift uint64) common.Hash {
	bigIntOffset := new(big.Int)
	bigIntOffset.SetUint64(shift)
	return s.ShiftHashBigInt(hash, bigIntOffset)
}

// ShiftHashBigInt shift hash in big.Int
func (s Storage) ShiftHashBigInt(hash common.Hash, shift *big.Int) common.Hash {
	head := hash.Big()
	head.Add(head, shift)
	return common.BytesToHash(head.Bytes())
}

func getDByteSize(data common.Hash) uint64 {
	bytes := data.Bytes()
	lastByte := bytes[len(bytes)-1]
	if lastByte&0x1 == 0 {
		return uint64(lastByte / 2)
	}
	return new(big.Int).Div(new(big.Int).Sub(
		data.Big(), big.NewInt(1)), big.NewInt(2)).Uint64()
}

// DecodeDByteBySlot given contract address and slot return the dynamic bytes data.
func (s Storage) DecodeDByteBySlot(address common.Address, slot common.Hash) []byte {
	data := s.GetState(address, slot)
	length := getDByteSize(data)
	if length < 32 {
		return data[:length]
	}
	ptr := crypto.Keccak256Hash(slot.Bytes())
	slotNum := (length-1)/32 + 1
	rVal := make([]byte, slotNum*32)
	for i := uint64(0); i < slotNum; i++ {
		start := i * 32
		copy(rVal[start:start+32], s.GetState(address, ptr).Bytes())
		ptr = s.ShiftHashUint64(ptr, 1)
	}
	return rVal[:length]
}
