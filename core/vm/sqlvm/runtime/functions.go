package runtime

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	se "github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/crypto"
)

// function identifier
const (
	BLOCKHASH uint16 = iota
	BLOCKNUMBER
	BLOCKTIMESTAMP
	BLOCKCOINBASE
	BLOCKGASLIMIT
	MSGSENDER
	MSGDATA
	TXORIGIN
	NOW
	RAND
	BITAND
	BITOR
	BITXOR
	BITNOT
)

type fn func(*common.Context, []*Operand, uint64) (*Operand, error)

var (
	fnTable = []fn{
		BLOCKHASH:      fnBlockHash,
		BLOCKNUMBER:    fnBlockNumber,
		BLOCKTIMESTAMP: fnBlockTimestamp,
		NOW:            fnBlockTimestamp,
		BLOCKCOINBASE:  fnBlockCoinBase,
		BLOCKGASLIMIT:  fnBlockGasLimit,
		MSGSENDER:      fnMsgSender,
		MSGDATA:        fnMsgData,
		TXORIGIN:       fnTxOrigin,
		RAND:           fnRand,
		BITAND:         fnBitAnd,
		BITOR:          fnBitOr,
		BITXOR:         fnBitXor,
		BITNOT:         fnBitNot,
	}
)

func assignFuncResult(meta []ast.DataType, fn func() *Raw, length uint64) (result *Operand) {
	result = &Operand{Meta: meta, Data: make([]Tuple, length)}
	for i := uint64(0); i < length; i++ {
		result.Data[i] = Tuple{fn()}
	}
	return
}

func evalBlockHash(ctx *common.Context, num, cur decimal.Decimal) (r *Raw, err error) {
	r = &Raw{Bytes: make([]byte, 32)}

	cNum := cur.Sub(dec.Dec257)
	if num.Cmp(cNum) > 0 && num.Cmp(cur) < 0 {
		var num64 uint64
		num64, err = ast.DecimalToUint64(num)
		if err != nil {
			return
		}
		r.Bytes = ctx.GetHash(num64).Bytes()
	}
	return
}

func fnBlockHash(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	if len(ops) != 1 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}

	meta := []ast.DataType{ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 3)}
	cNum := decimal.NewFromBigInt(ctx.BlockNumber, 0)

	if ops[0].IsImmediate {
		var r *Raw
		r, err = evalBlockHash(ctx, ops[0].Data[0][0].Value, cNum)
		if err != nil {
			return
		}
		result = assignFuncResult(meta, r.clone, length)
	} else {
		result = &Operand{Meta: meta, Data: make([]Tuple, length)}
		for i := uint64(0); i < length; i++ {
			var r *Raw
			r, err = evalBlockHash(ctx, ops[0].Data[i][0].Value, cNum)
			if err != nil {
				return
			}
			result.Data[i] = Tuple{r}
		}
	}
	return
}

func fnBlockNumber(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	r := &Raw{Value: decimal.NewFromBigInt(ctx.BlockNumber, 0)}
	result = assignFuncResult(
		[]ast.DataType{ast.ComposeDataType(ast.DataTypeMajorUint, 31)},
		r.clone, length,
	)
	return
}

func fnBlockTimestamp(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	r := &Raw{Value: decimal.NewFromBigInt(ctx.Time, 0)}
	result = assignFuncResult(
		[]ast.DataType{ast.ComposeDataType(ast.DataTypeMajorUint, 31)},
		r.clone, length,
	)
	return
}

func fnBlockCoinBase(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	r := &Raw{Bytes: ctx.Coinbase.Bytes()}
	result = assignFuncResult(
		[]ast.DataType{ast.ComposeDataType(ast.DataTypeMajorAddress, 0)},
		r.clone, length,
	)
	return
}

func fnBlockGasLimit(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	r := &Raw{}
	if ctx.GasLimit > uint64(math.MaxInt64) {
		r.Value, err = decimal.NewFromString(fmt.Sprint(ctx.GasLimit))
		if err != nil {
			return
		}
	} else {
		r.Value = decimal.New(int64(ctx.GasLimit), 0)
	}
	result = assignFuncResult(
		[]ast.DataType{ast.ComposeDataType(ast.DataTypeMajorUint, 7)},
		r.clone, length,
	)
	return
}

func fnMsgSender(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	r := &Raw{Bytes: ctx.Contract.CallerAddress.Bytes()}
	result = assignFuncResult(
		[]ast.DataType{ast.ComposeDataType(ast.DataTypeMajorAddress, 0)},
		r.clone, length,
	)
	return
}

func fnMsgData(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	r := &Raw{Bytes: ctx.Contract.Input}
	result = assignFuncResult(
		[]ast.DataType{ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0)},
		r.clone, length,
	)
	return
}

func fnTxOrigin(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	r := &Raw{Bytes: ctx.Origin.Bytes()}
	result = assignFuncResult(
		[]ast.DataType{ast.ComposeDataType(ast.DataTypeMajorAddress, 0)},
		r.clone, length,
	)
	return
}

func fnRand(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	binaryOriginNonce := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(binaryOriginNonce, ctx.Storage.GetNonce(ctx.Origin))

	binaryUsedIndex := make([]byte, binary.MaxVarintLen64)
	vType := ast.ComposeDataType(ast.DataTypeMajorUint, 31)

	fn := func() *Raw {
		binary.PutUvarint(binaryUsedIndex, ctx.RandCallIndex)
		ctx.RandCallIndex++

		hash := crypto.Keccak256(
			ctx.Randomness,
			ctx.Origin.Bytes(),
			binaryOriginNonce,
			binaryUsedIndex)

		v, _ := ast.DecimalDecode(vType, hash)
		return &Raw{Value: v}
	}

	result = assignFuncResult(
		[]ast.DataType{ast.ComposeDataType(ast.DataTypeMajorUint, 31)},
		fn, length,
	)
	return
}

func metaBitOp(dType ast.DataType) bool {
	dMajor, _ := ast.DecomposeDataType(dType)
	switch dMajor {
	case ast.DataTypeMajorUint,
		ast.DataTypeMajorInt,
		ast.DataTypeMajorFixedBytes:
		return true
	}
	return false
}

func metaAllBitOp(op *Operand) bool { return metaAll(op, metaBitOp) }

func extractOps(ops []*Operand) (n int, op1, op2 *Operand, err error) {
	if len(ops) < 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}

	n, err = findMaxDataLength(ops)
	if err != nil {
		return
	}

	op1, op2 = ops[0], ops[1]

	if !metaAllEq(op1, op2) || !metaAllBitOp(op1) {
		err = se.ErrorCodeInvalidDataType
	}
	return
}

func (r *Raw) toBytes(dType ast.DataType) []byte {
	dMajor, _ := ast.DecomposeDataType(dType)
	switch dMajor {
	case ast.DataTypeMajorFixedBytes,
		ast.DataTypeMajorAddress,
		ast.DataTypeMajorDynamicBytes:
		return r.Bytes
	case ast.DataTypeMajorUint,
		ast.DataTypeMajorInt,
		ast.DataTypeMajorFixed:
		bytes, err := ast.DecimalEncode(dType, r.Value)
		if err != nil {
			panic(err)
		}
		return bytes
	default:
		panic(fmt.Errorf("unrecognized data type: %v", dType))
	}
}

func (r *Raw) fromBytes(bytes []byte, dType ast.DataType) {
	dMajor, _ := ast.DecomposeDataType(dType)
	switch dMajor {
	case ast.DataTypeMajorFixedBytes,
		ast.DataTypeMajorAddress,
		ast.DataTypeMajorDynamicBytes:
		r.Bytes = bytes
	case ast.DataTypeMajorUint,
		ast.DataTypeMajorInt,
		ast.DataTypeMajorFixed:
		var err error
		r.Value, err = ast.DecimalDecode(dType, bytes)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("unrecognized data type: %v", dType))
	}
}

type bitBinFunc func(b1, b2 byte) byte

func fnBitAnd(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	n, op1, op2, err := extractOps(ops)
	if err != nil {
		return
	}

	result = op1.clone(true)
	result.Data = make([]Tuple, n)
	for i := 0; i < n; i++ {
		result.Data[i] = op1.Data[i].bitBinOp(
			op2.Data[i],
			op1.Meta,
			func(b1, b2 byte) byte { return b1 & b2 },
		)
	}
	return
}

func (t Tuple) bitBinOp(t2 Tuple, meta []ast.DataType, bFn bitBinFunc) (t3 Tuple) {
	t3 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t3[i] = t[i].bitBinOp(t2[i], meta[i], bFn)
	}
	return
}

func (r *Raw) bitBinOp(r2 *Raw, dType ast.DataType, bFn bitBinFunc) (r3 *Raw) {
	bytes1, bytes2 := r.toBytes(dType), r2.toBytes(dType)

	if len(bytes1) != len(bytes2) {
		panic(fmt.Errorf("bitwise operand on differnt length bits"))
	}

	n := len(bytes1)
	bytes3 := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes3[i] = bFn(bytes1[i], bytes2[i])
	}

	r3 = &Raw{}
	r3.fromBytes(bytes3, dType)
	return
}

func fnBitOr(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	n, op1, op2, err := extractOps(ops)
	if err != nil {
		return
	}

	result = op1.clone(true)
	result.Data = make([]Tuple, n)
	for i := 0; i < n; i++ {
		result.Data[i] = op1.Data[i].bitBinOp(
			op2.Data[i],
			op1.Meta,
			func(b1, b2 byte) byte { return b1 | b2 },
		)
	}
	return
}

func fnBitXor(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	n, op1, op2, err := extractOps(ops)
	if err != nil {
		return
	}

	result = op1.clone(true)
	result.Data = make([]Tuple, n)
	for i := 0; i < n; i++ {
		result.Data[i] = op1.Data[i].bitBinOp(
			op2.Data[i],
			op1.Meta,
			func(b1, b2 byte) byte { return b1 ^ b2 },
		)
	}
	return
}

type bitUnFunc func(b byte) byte

func fnBitNot(ctx *common.Context, ops []*Operand, length uint64) (result *Operand, err error) {
	if len(ops) < 2 {
		err = se.ErrorCodeInvalidOperandNum
		return
	}

	op := ops[1]
	if !metaAllBitOp(op) {
		err = se.ErrorCodeInvalidDataType
		return
	}

	result = op.clone(true)
	result.Data = make([]Tuple, len(op.Data))
	for i := 0; i < len(op.Data); i++ {
		result.Data[i] = op.Data[i].bitUnOp(
			op.Meta,
			func(b byte) byte { return ^b },
		)
	}
	return
}

func (t Tuple) bitUnOp(meta []ast.DataType, bFn bitUnFunc) (t2 Tuple) {
	t2 = make(Tuple, len(t))
	for i := 0; i < len(t); i++ {
		t2[i] = t[i].bitUnOp(meta[i], bFn)
	}
	return
}

func (r *Raw) bitUnOp(dType ast.DataType, bFn bitUnFunc) (r2 *Raw) {
	bytes := r.toBytes(dType)

	n := len(bytes)
	bytes2 := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes2[i] = bFn(bytes[i])
	}

	r2 = &Raw{}
	r2.fromBytes(bytes2, dType)
	return
}
