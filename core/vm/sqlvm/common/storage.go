package common

import (
	"math/big"

	"github.com/dexon-foundation/decimal"
	"golang.org/x/crypto/sha3"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
	"github.com/dexon-foundation/dexon/crypto"
	"github.com/dexon-foundation/dexon/rlp"
)

// Constants for path keys.
var (
	pathCompTables         = []byte("tables")
	pathCompPrimary        = []byte("primary")
	pathCompIndices        = []byte("indices")
	pathCompSequence       = []byte("sequence")
	pathCompOwner          = []byte("owner")
	pathCompWriters        = []byte("writers")
	pathCompReverseIndices = []byte("reverse_indices")
)

// Storage holds SQLVM required data and method.
type Storage struct {
	vm.StateDB
	Schema schema.Schema
}

// NewStorage return Storage instance.
func NewStorage(state vm.StateDB) *Storage {
	s := &Storage{state, schema.Schema{}}
	return s
}

// TODO(yenlin): Do we really need to use ast encode/decode here?
func uint64ToBytes(id uint64) []byte {
	bigIntID := new(big.Int).SetUint64(id)
	decimalID := decimal.NewFromBigInt(bigIntID, 0)
	dt := ast.ComposeDataType(ast.DataTypeMajorUint, 7)
	byteID, _ := ast.DecimalEncode(dt, decimalID)
	return byteID
}

func bytesToUint64(b []byte) uint64 {
	dt := ast.ComposeDataType(ast.DataTypeMajorUint, 7)
	d, _ := ast.DecimalDecode(dt, b)
	// TODO(yenlin): Not yet a convenient way to extract uint64 from decimal...
	bigInt := d.Rescale(0).Coefficient()
	return bigInt.Uint64()
}

func hashToAddress(hash common.Hash) common.Address {
	return common.BytesToAddress(hash.Bytes())
}

func addressToHash(addr common.Address) common.Hash {
	return common.BytesToHash(addr.Bytes())
}

func (s *Storage) hashPathKey(key [][]byte) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, key)
	// length of common.Hash is 256bit,
	// so it can properly match the size of hw.Sum
	hw.Sum(h[:0])
	return
}

// GetRowPathHash return primary key hash which points to row data.
func (s *Storage) GetRowPathHash(tableName []byte, rowID uint64) common.Hash {
	// PathKey(["tables", "{table_name}", "primary", uint64({row_id})])
	key := [][]byte{
		pathCompTables,
		tableName,
		pathCompPrimary,
		uint64ToBytes(rowID),
	}
	return s.hashPathKey(key)
}

// GetIndexValuesPathHash return the hash address to IndexValues structure
// which contains all possible values.
func (s *Storage) GetIndexValuesPathHash(
	tableName, indexName []byte,
) common.Hash {
	// PathKey(["tables", "{table_name}", "indices", "{index_name}"])
	key := [][]byte{
		pathCompTables,
		tableName,
		pathCompIndices,
		indexName,
	}
	return s.hashPathKey(key)
}

// GetIndexEntryPathHash return the hash address to IndexEntry structure for a
// given value.
func (s *Storage) GetIndexEntryPathHash(
	tableName, indexName []byte,
	values ...[]byte,
) common.Hash {
	// PathKey(["tables", "{table_name}", "indices", "{index_name}", field_1, field_2, field_3, ...])
	key := make([][]byte, 0, 4+len(values))
	key = append(key, pathCompTables, tableName, pathCompIndices, indexName)
	key = append(key, values...)
	return s.hashPathKey(key)
}

// GetReverseIndexPathHash return the hash address to IndexRev structure for a
// row in a table.
func (s *Storage) GetReverseIndexPathHash(
	tableName []byte,
	rowID uint64,
) common.Hash {
	// PathKey(["tables", "{table_name}", "reverse_indices", "{RowID}"])
	key := [][]byte{
		pathCompTables,
		tableName,
		pathCompReverseIndices,
		uint64ToBytes(rowID),
	}
	return s.hashPathKey(key)
}

// getSequencePathHash return the hash address of a sequence.
func (s *Storage) getSequencePathHash(
	tableName []byte, seqIdx uint8,
) common.Hash {
	// PathKey(["tables", "{table_name}", "sequence", uint8(sequence_idx)])
	key := [][]byte{
		pathCompTables,
		tableName,
		pathCompSequence,
		{seqIdx}, // TODO(yenlin): use some other encode method on uint8?
	}
	return s.hashPathKey(key)
}

func (s *Storage) getOwnerPathHash() common.Hash {
	// PathKey(["owner"])
	key := [][]byte{pathCompOwner}
	return s.hashPathKey(key)
}

func (s *Storage) getTableWritersPathHash(tableName []byte) common.Hash {
	// PathKey(["tables", "{table_name}", "writers"])
	key := [][]byte{
		pathCompTables,
		tableName,
		pathCompWriters,
	}
	return s.hashPathKey(key)
}

func (s *Storage) getTableWriterRevIdxPathHash(
	tableName []byte,
	account common.Address,
) common.Hash {
	// PathKey(["tables", "{table_name}", "writers", "{addr}"])
	key := [][]byte{
		pathCompTables,
		tableName,
		pathCompWriters,
		account.Bytes(),
	}
	return s.hashPathKey(key)
}

// ShiftHashUint64 shift hash in uint64.
func (s *Storage) ShiftHashUint64(hash common.Hash, shift uint64) common.Hash {
	bigIntOffset := new(big.Int)
	bigIntOffset.SetUint64(shift)
	return s.ShiftHashBigInt(hash, bigIntOffset)
}

// ShiftHashBigInt shift hash in big.Int
func (s *Storage) ShiftHashBigInt(hash common.Hash, shift *big.Int) common.Hash {
	head := hash.Big()
	head.Add(head, shift)
	return common.BytesToHash(head.Bytes())
}

// ShiftHashListEntry shift hash from the head of a list to the hash of
// idx-th entry.
func (s *Storage) ShiftHashListEntry(
	base common.Hash,
	headerSize uint64,
	entrySize uint64,
	idx uint64,
) common.Hash {
	// TODO(yenlin): tuning when headerSize+entrySize*idx do not overflow.
	shift := new(big.Int)
	operand := new(big.Int)
	shift.SetUint64(entrySize)
	operand.SetUint64(idx)
	shift.Mul(shift, operand)
	operand.SetUint64(headerSize)
	shift.Add(shift, operand)
	return s.ShiftHashBigInt(base, shift)
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
func (s *Storage) DecodeDByteBySlot(address common.Address, slot common.Hash) []byte {
	data := s.GetState(address, slot)
	length := getDByteSize(data)
	if length < common.HashLength {
		return data[:length]
	}
	ptr := crypto.Keccak256Hash(slot.Bytes())
	slotNum := (length-1)/common.HashLength + 1
	rVal := make([]byte, slotNum*common.HashLength)
	for i := uint64(0); i < slotNum; i++ {
		start := i * common.HashLength
		copy(rVal[start:start+common.HashLength], s.GetState(address, ptr).Bytes())
		ptr = s.ShiftHashUint64(ptr, 1)
	}
	return rVal[:length]
}

// SQLVM metadata structure operations.

// IndexValues contain addresses to all possible values of an index.
type IndexValues struct {
	// Header.
	Length uint64
	// 3 unused uint64 fields here.
	// Contents.
	ValueHashes []common.Hash
}

// IndexEntry contain row ids of a given value in an index.
type IndexEntry struct {
	// Header.
	Length              uint64
	IndexToValuesOffset uint64
	ForeignKeyRefCount  uint64
	// 1 unused uint64 field here.
	// Contents.
	RowIDs []uint64
}

// LoadIndexValues load IndexValues struct of a given index.
func (s *Storage) LoadIndexValues(
	contract common.Address,
	tableName, indexName []byte,
	onlyHeader bool,
) *IndexValues {
	ret := &IndexValues{}
	slot := s.GetIndexValuesPathHash(tableName, indexName)
	data := s.GetState(contract, slot)
	ret.Length = bytesToUint64(data[:8])
	if onlyHeader {
		return ret
	}
	// Load all ValueHashes.
	ret.ValueHashes = make([]common.Hash, ret.Length)
	for i := uint64(0); i < ret.Length; i++ {
		slot = s.ShiftHashUint64(slot, 1)
		ret.ValueHashes[i] = s.GetState(contract, slot)
	}
	return ret
}

// LoadIndexEntry load IndexEntry struct of a given value key on an index.
func (s *Storage) LoadIndexEntry(
	contract common.Address,
	tableName, indexName []byte,
	onlyHeader bool,
	values ...[]byte,
) *IndexEntry {
	ret := &IndexEntry{}
	slot := s.GetIndexEntryPathHash(tableName, indexName, values...)
	data := s.GetState(contract, slot)
	ret.Length = bytesToUint64(data[:8])
	ret.IndexToValuesOffset = bytesToUint64(data[8:16])
	ret.ForeignKeyRefCount = bytesToUint64(data[16:24])

	if onlyHeader {
		return ret
	}
	// Load all RowIDs.
	ret.RowIDs = make([]uint64, 0, ret.Length)
	remain := ret.Length
	for remain > 0 {
		bound := remain
		if bound > 4 {
			bound = 4
		}
		slot = s.ShiftHashUint64(slot, 1)
		data := s.GetState(contract, slot).Bytes()
		for i := uint64(0); i < bound; i++ {
			ret.RowIDs = append(ret.RowIDs, bytesToUint64(data[:8]))
			data = data[8:]
		}
		remain -= bound
	}
	return ret
}

// LoadOwner load the owner of a SQLVM contract from storage.
func (s *Storage) LoadOwner(contract common.Address) common.Address {
	return hashToAddress(s.GetState(contract, s.getOwnerPathHash()))
}

// StoreOwner save the owner of a SQLVM contract to storage.
func (s *Storage) StoreOwner(contract, newOwner common.Address) {
	s.SetState(contract, s.getOwnerPathHash(), addressToHash(newOwner))
}

type tableWriters struct {
	Length uint64
	// 3 unused uint64 in slot 1.
	Writers []common.Address // Each address consumes one slot, right aligned.
}

type tableWriterRevIdx struct {
	IndexToValuesOffset uint64
	// 3 unused uint64 in the slot.
}

func (c *tableWriterRevIdx) Valid() bool {
	return c.IndexToValuesOffset != 0
}

func (s *Storage) loadTableWriterRevIdx(
	contract common.Address,
	path common.Hash,
) *tableWriterRevIdx {
	ret := &tableWriterRevIdx{}
	data := s.GetState(contract, path)
	ret.IndexToValuesOffset = bytesToUint64(data[:8])
	return ret
}

func (s *Storage) storeTableWriterRevIdx(
	contract common.Address,
	path common.Hash,
	rev *tableWriterRevIdx,
) {
	var data common.Hash // One slot.
	copy(data[:8], uint64ToBytes(rev.IndexToValuesOffset))
	s.SetState(contract, path, data)
}

func (s *Storage) loadTableWriters(
	contract common.Address,
	pathHash common.Hash,
	onlyHeader bool,
) *tableWriters {
	ret := &tableWriters{}
	header := s.GetState(contract, pathHash)
	ret.Length = bytesToUint64(header[:8])
	if onlyHeader {
		return ret
	}
	ret.Writers = make([]common.Address, ret.Length)
	for i := uint64(0); i < ret.Length; i++ {
		ret.Writers[i] = s.loadSingleTableWriter(contract, pathHash, i)
	}
	return ret
}

func (s *Storage) storeTableWritersHeader(
	contract common.Address,
	pathHash common.Hash,
	w *tableWriters,
) {
	var header common.Hash
	copy(header[:8], uint64ToBytes(w.Length))
	s.SetState(contract, pathHash, header)
}

func (s *Storage) shiftTableWriterList(
	base common.Hash,
	idx uint64,
) common.Hash {
	return s.ShiftHashListEntry(base, 1, 1, idx)
}

func (s *Storage) loadSingleTableWriter(
	contract common.Address,
	writersPathHash common.Hash,
	idx uint64,
) common.Address {
	slot := s.shiftTableWriterList(writersPathHash, idx)
	acc := s.GetState(contract, slot)
	return hashToAddress(acc)
}

func (s *Storage) storeSingleTableWriter(
	contract common.Address,
	writersPathHash common.Hash,
	idx uint64,
	acc common.Address,
) {
	slot := s.shiftTableWriterList(writersPathHash, idx)
	s.SetState(contract, slot, addressToHash(acc))
}

// IsTableWriter check if an account is writer to the table.
func (s *Storage) IsTableWriter(
	contract common.Address,
	tableName []byte,
	account common.Address,
) bool {
	path := s.getTableWriterRevIdxPathHash(tableName, account)
	rev := s.loadTableWriterRevIdx(contract, path)
	return rev.Valid()
}

// LoadTableWriters load writers of a table.
func (s *Storage) LoadTableWriters(
	contract common.Address,
	tableName []byte,
) (ret []common.Address) {
	path := s.getTableWritersPathHash(tableName)
	writers := s.loadTableWriters(contract, path, false)
	return writers.Writers
}

// InsertTableWriter insert an account into writer list of the table.
func (s *Storage) InsertTableWriter(
	contract common.Address,
	tableName []byte,
	account common.Address,
) {
	revPath := s.getTableWriterRevIdxPathHash(tableName, account)
	rev := s.loadTableWriterRevIdx(contract, revPath)
	if rev.Valid() {
		return
	}
	path := s.getTableWritersPathHash(tableName)
	writers := s.loadTableWriters(contract, path, true)
	// Store modification.
	s.storeSingleTableWriter(contract, path, writers.Length, account)
	writers.Length++
	s.storeTableWritersHeader(contract, path, writers)
	// Notice: IndexToValuesOffset starts from 1.
	s.storeTableWriterRevIdx(contract, revPath, &tableWriterRevIdx{
		IndexToValuesOffset: writers.Length,
	})
}

// DeleteTableWriter delete an account from writer list of the table.
func (s *Storage) DeleteTableWriter(
	contract common.Address,
	tableName []byte,
	account common.Address,
) {
	revPath := s.getTableWriterRevIdxPathHash(tableName, account)
	rev := s.loadTableWriterRevIdx(contract, revPath)
	if !rev.Valid() {
		return
	}
	path := s.getTableWritersPathHash(tableName)
	writers := s.loadTableWriters(contract, path, true)

	// Store modification.
	if rev.IndexToValuesOffset != writers.Length {
		// Move last to deleted slot.
		lastAcc := s.loadSingleTableWriter(contract, path, writers.Length-1)
		s.storeSingleTableWriter(contract, path, rev.IndexToValuesOffset-1,
			lastAcc)
		s.storeTableWriterRevIdx(contract, s.getTableWriterRevIdxPathHash(
			tableName, lastAcc), rev)
	}
	// Delete last.
	writers.Length--
	s.storeTableWritersHeader(contract, path, writers)
	s.storeSingleTableWriter(contract, path, writers.Length, common.Address{})
	s.storeTableWriterRevIdx(contract, revPath, &tableWriterRevIdx{})
}

// IncSequence increment value of sequence by inc and return the old value.
func (s *Storage) IncSequence(
	contract common.Address,
	tableName []byte,
	seqIdx uint8,
	inc uint64,
) uint64 {
	seqPath := s.getSequencePathHash(tableName, seqIdx)
	slot := s.GetState(contract, seqPath)
	val := bytesToUint64(slot.Bytes())
	// TODO(yenlin): Check overflow?
	s.SetState(contract, seqPath, common.BytesToHash(uint64ToBytes(val+inc)))
	return val
}
