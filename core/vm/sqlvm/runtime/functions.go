package runtime

import (
	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	se "github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

// function identifier
const (
	BLOCKHASH   = "BLOCK_HASH"
	BLOCKNUMBER = "BLOCK_NUMBER"
)

type fn func(*common.Context, []*Operand, uint64) (*Operand, error)

var (
	fnTable = map[string]fn{
		BLOCKHASH:   fnBlockHash,
		BLOCKNUMBER: fnBlockNumber,
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
