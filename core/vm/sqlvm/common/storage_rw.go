package common

import (
	"io"

	"github.com/dexon-foundation/dexon/common"
)

// StorageReader implements io.Reader on Storage.
// Notice that we have cache in Reader, so it become invalid after any writes
// to Storage or StateDB.
type StorageReader struct {
	storage  *Storage
	contract common.Address
	cursor   common.Hash
	buffer   []byte
}

// assert io.Reader interface is implemented.
var _ io.Reader = (*StorageReader)(nil)

// NewStorageReader create a Reader on Storage.
func NewStorageReader(
	storage *Storage,
	contract common.Address,
	startPos common.Hash,
) *StorageReader {
	return &StorageReader{
		storage:  storage,
		contract: contract,
		cursor:   startPos,
	}
}

// Read implements the function defined in io.Reader interface.
func (r *StorageReader) Read(p []byte) (n int, err error) {
	copyBuffer := func() {
		lenCopy := len(p)
		if lenCopy > len(r.buffer) {
			lenCopy = len(r.buffer)
		}
		copy(p[:lenCopy], r.buffer[:lenCopy])
		p = p[lenCopy:]
		r.buffer = r.buffer[lenCopy:]
		n += lenCopy
	}
	// Flush old buffer first.
	copyBuffer()
	for len(p) > 0 {
		// Read slot by slot.
		r.buffer = r.storage.GetState(r.contract, r.cursor).Bytes()
		r.cursor = r.storage.ShiftHashUint64(r.cursor, 1)
		copyBuffer()
	}
	return
}

// StorageWriter implements io.Writer on Storage.
type StorageWriter struct {
	storage   *Storage
	contract  common.Address
	cursor    common.Hash
	byteShift int // bytes already written in last slot. value: 0 ~ 31
}

// assert io.Writer interface is implemented.
var _ io.Writer = (*StorageWriter)(nil)

// NewStorageWriter create a Writer on Storage.
func NewStorageWriter(
	storage *Storage,
	contract common.Address,
	startPos common.Hash,
) *StorageWriter {
	return &StorageWriter{
		storage:  storage,
		contract: contract,
		cursor:   startPos,
	}
}

// Write implements the function defined in io.Writer interface.
func (w *StorageWriter) Write(p []byte) (n int, err error) {
	var payload common.Hash
	for len(p) > 0 {
		// Setup common.Hash to write.
		remain := common.HashLength - w.byteShift
		lenCopy := remain
		if lenCopy > len(p) {
			lenCopy = len(p)
		}
		if lenCopy != common.HashLength {
			// Not writing an entire slot, need load first.
			payload = w.storage.GetState(w.contract, w.cursor)
		}
		b := payload.Bytes()
		start := w.byteShift
		end := w.byteShift + lenCopy
		copy(b[start:end], p[:lenCopy])
		payload.SetBytes(b)

		w.storage.SetState(w.contract, w.cursor, payload)

		// Update state.
		p = p[lenCopy:]
		n += lenCopy
		w.byteShift += lenCopy
		if w.byteShift == common.HashLength {
			w.byteShift = 0
			w.cursor = w.storage.ShiftHashUint64(w.cursor, 1)
		}
	}
	return
}
