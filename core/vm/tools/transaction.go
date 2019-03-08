package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"

	"golang.org/x/crypto/sha3"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/common/hexutil"
	"github.com/dexon-foundation/dexon/crypto"
	"github.com/dexon-foundation/dexon/rlp"
)

//go:generate gencodec -type txdata -field-override txdataMarshaling -out gen_tx_json.go

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
)

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

type Transaction struct {
	Data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type txdataMarshaling struct {
	AccountNonce hexutil.Uint64
	Price        *hexutil.Big
	GasLimit     hexutil.Uint64
	Amount       *hexutil.Big
	Payload      hexutil.Bytes
	V            *hexutil.Big
	R            *hexutil.Big
	S            *hexutil.Big
}

// Protected returns whether the transaction is protected from replay protection.
func (tx *Transaction) Protected() bool {
	return isProtectedV(tx.Data.V)
}
func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28
	}
	// anything not 27 or 28 is considered protected
	return true
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.Data)
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&tx.Data)
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		fmt.Println("tools load")
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.Data)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// MarshalJSON encodes the web3 RPC transaction format.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	hash := tx.Hash()
	data := tx.Data
	data.Hash = &hash
	return data.MarshalJSON()
}

// MarshalJSON marshals as JSON.
func (t txdata) MarshalJSON() ([]byte, error) {
	type txdata struct {
		AccountNonce hexutil.Uint64  `json:"nonce"    gencodec:"required"`
		Price        *hexutil.Big    `json:"gasPrice" gencodec:"required"`
		GasLimit     hexutil.Uint64  `json:"gas"      gencodec:"required"`
		Recipient    *common.Address `json:"to"       rlp:"nil"`
		Amount       *hexutil.Big    `json:"value"    gencodec:"required"`
		Payload      hexutil.Bytes   `json:"input"    gencodec:"required"`
		V            *hexutil.Big    `json:"v" gencodec:"required"`
		R            *hexutil.Big    `json:"r" gencodec:"required"`
		S            *hexutil.Big    `json:"s" gencodec:"required"`
		Hash         *common.Hash    `json:"hash" rlp:"-"`
	}
	var enc txdata
	enc.AccountNonce = hexutil.Uint64(t.AccountNonce)
	enc.Price = (*hexutil.Big)(t.Price)
	enc.GasLimit = hexutil.Uint64(t.GasLimit)
	enc.Recipient = t.Recipient
	enc.Amount = (*hexutil.Big)(t.Amount)
	enc.Payload = t.Payload
	enc.V = (*hexutil.Big)(t.V)
	enc.R = (*hexutil.Big)(t.R)
	enc.S = (*hexutil.Big)(t.S)
	enc.Hash = t.Hash
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (t *txdata) UnmarshalJSON(input []byte) error {
	type txdata struct {
		AccountNonce *hexutil.Uint64 `json:"nonce"    gencodec:"required"`
		Price        *hexutil.Big    `json:"gasPrice" gencodec:"required"`
		GasLimit     *hexutil.Uint64 `json:"gas"      gencodec:"required"`
		Recipient    *common.Address `json:"to"       rlp:"nil"`
		Amount       *hexutil.Big    `json:"value"    gencodec:"required"`
		Payload      *hexutil.Bytes  `json:"input"    gencodec:"required"`
		V            *hexutil.Big    `json:"v" gencodec:"required"`
		R            *hexutil.Big    `json:"r" gencodec:"required"`
		S            *hexutil.Big    `json:"s" gencodec:"required"`
		Hash         *common.Hash    `json:"hash" rlp:"-"`
	}
	var dec txdata
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.AccountNonce == nil {
		return errors.New("missing required field 'nonce' for txdata")
	}
	t.AccountNonce = uint64(*dec.AccountNonce)
	if dec.Price == nil {
		return errors.New("missing required field 'gasPrice' for txdata")
	}
	t.Price = (*big.Int)(dec.Price)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gas' for txdata")
	}
	t.GasLimit = uint64(*dec.GasLimit)
	if dec.Recipient != nil {
		t.Recipient = dec.Recipient
	}
	if dec.Amount == nil {
		return errors.New("missing required field 'value' for txdata")
	}
	t.Amount = (*big.Int)(dec.Amount)
	if dec.Payload == nil {
		return errors.New("missing required field 'input' for txdata")
	}
	t.Payload = *dec.Payload
	if dec.V == nil {
		return errors.New("missing required field 'v' for txdata")
	}
	t.V = (*big.Int)(dec.V)
	if dec.R == nil {
		return errors.New("missing required field 'r' for txdata")
	}
	t.R = (*big.Int)(dec.R)
	if dec.S == nil {
		return errors.New("missing required field 's' for txdata")
	}
	t.S = (*big.Int)(dec.S)
	if dec.Hash != nil {
		t.Hash = dec.Hash
	}
	return nil
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec txdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}

	withSignature := dec.V.Sign() != 0 || dec.R.Sign() != 0 || dec.S.Sign() != 0
	if withSignature {
		var V byte
		if isProtectedV(dec.V) {
			chainID := deriveChainId(dec.V).Uint64()
			V = byte(dec.V.Uint64() - 35 - 2*chainID)
		} else {
			V = byte(dec.V.Uint64() - 27)
		}
		if !crypto.ValidateSignatureValues(V, dec.R, dec.S, false) {
			return ErrInvalidSig
		}
	}

	*tx = Transaction{Data: dec}
	return nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}

// Hash hashes the RLP encoding of tx.
// It uniquely identifies the transaction.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}
func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
func (tx *Transaction) SetFrom(from atomic.Value) {
	tx.from.Store(from.Load())
}
func (tx *Transaction) Gas() uint64        { return tx.Data.GasLimit }
func (tx *Transaction) GasPrice() *big.Int { return new(big.Int).Set(tx.Data.Price) }
func (tx *Transaction) Value() *big.Int    { return new(big.Int).Set(tx.Data.Amount) }
func (tx *Transaction) Nonce() uint64      { return tx.Data.AccountNonce }
func (tx *Transaction) CheckNonce() bool   { return true }
