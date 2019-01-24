package parser

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/shopspring/decimal"
)

// Parser was generated with pigeon v1.0.0-99-gbb0192c.
//go:generate pigeon -no-recover -o grammar.go grammar.peg
//go:generate sh -c "sed -f grammar.sed grammar.go > grammar_new.go"
//go:generate mv grammar_new.go grammar.go
//go:generate goimports -w grammar.go

func prepend(x interface{}, xs interface{}) []interface{} {
	return append([]interface{}{x}, toSlice(xs)...)
}

func toSlice(x interface{}) []interface{} {
	if x == nil {
		return nil
	}
	return x.([]interface{})
}

// TODO(wmin0): finish it.
func isAddress(h []byte) bool {
	return false
}

func hexToInteger(h []byte) interface{} {
	d := decimal.Zero
	l := len(h)
	base := decimal.New(16, 0)
	for idx, b := range h {
		i, err := strconv.ParseInt(string([]byte{b}), 16, 32)
		if err != nil {
			panic(fmt.Sprintf("invalid hex digit %s: %v", []byte{b}, err))
		}
		d = d.Add(
			decimal.New(i, 0).
				Mul(base.Pow(decimal.New(int64(l-idx-1), 0))),
		)
	}
	return ast.IntegerValueNode{V: d, IsAddress: isAddress(h)}
}

func hexToBytes(h []byte) []byte {
	bs := make([]byte, hex.DecodedLen(len(h)))
	_, err := hex.Decode(bs, h)
	if err != nil {
		panic(fmt.Sprintf("invalid hex string %s: %v", h, err))
	}
	return bs
}

func convertNumError(err error) errors.ErrorCode {
	if err == nil {
		return errors.ErrorCodeNil
	}
	switch err.(*strconv.NumError).Err {
	case strconv.ErrSyntax:
		return errors.ErrorCodeSyntax
	case strconv.ErrRange:
		return errors.ErrorCodeIntegerRange
	}
	panic(fmt.Sprintf("unknown NumError: %v", err))
}

func convertDecimalError(err error) errors.ErrorCode {
	if err == nil {
		return errors.ErrorCodeNil
	}
	errStr := err.Error()
	if strings.HasSuffix(errStr, "decimal: fractional part too long") {
		return errors.ErrorCodeFractionalPartTooLong
	} else if strings.HasSuffix(errStr, "decimal: exponent is not numeric") {
		return errors.ErrorCodeSyntax
	} else if strings.HasSuffix(errStr, "decimal: too many .s") {
		return errors.ErrorCodeSyntax
	}
	panic(fmt.Sprintf("unknown decimal error: %v", err))
}

func toUint(b []byte) (uint32, errors.ErrorCode) {
	i, err := strconv.ParseUint(string(b), 10, 32)
	return uint32(i), convertNumError(err)
}

func toDecimal(b []byte) (decimal.Decimal, errors.ErrorCode) {
	d, err := decimal.NewFromString(string(b))
	return d, convertDecimalError(err)
}

func toLower(b []byte) []byte {
	return bytes.ToLower(b)
}

func joinBytes(x interface{}) []byte {
	xs := toSlice(x)
	bs := []byte{}
	for _, b := range xs {
		bs = append(bs, b.([]byte)...)
	}
	return bs
}

func opSetSubject(op interface{}, s interface{}) interface{} {
	x := op.(ast.BinaryOperator)
	x.SetSubject(s)
	return x
}

func opSetObject(op interface{}, o interface{}) interface{} {
	x := op.(ast.BinaryOperator)
	x.SetObject(o)
	return x
}

func opSetTarget(op interface{}, t interface{}) interface{} {
	x := op.(ast.UnaryOperator)
	x.SetTarget(t)
	return x
}

func joinOperator(x interface{}, o interface{}) {
	if op, ok := x.(ast.UnaryOperator); ok {
		joinOperator(op.GetTarget(), o)
		return
	}
	if op, ok := x.(ast.BinaryOperator); ok {
		op.SetObject(o)
		return
	}
}

func rightJoinOperators(o interface{}, x interface{}) interface{} {
	xs := toSlice(x)
	if len(xs) == 0 {
		return o
	}
	l := len(xs)
	for idx := 0; idx < l-1; idx++ {
		joinOperator(xs[idx+1], xs[idx])
	}
	joinOperator(xs[0], o)
	return xs[l-1]
}

// TODO(wmin0): finish it.
func resolveString(s []byte) []byte {
	return s
}

// ParseString parses input string to AST.
func ParseString(s string) (interface{}, error) {
	root, pigeonErr := ParseReader("", strings.NewReader(s))
	if pigeonErr == nil {
		return root, pigeonErr
	}
	pigeonErrList := pigeonErr.(errList)
	sqlvmErrList := make(errors.ErrorList, len(pigeonErrList))
	for i := range pigeonErrList {
		parserErr := pigeonErrList[i].(*parserError)
		if sqlvmErr, ok := parserErr.Inner.(errors.Error); ok {
			sqlvmErrList[i] = sqlvmErr
		} else {
			sqlvmErrList[i] = errors.Error{
				Position: uint32(parserErr.pos.offset),
				Category: errors.ErrorCategoryGrammar,
				Code:     errors.ErrorCodeParser,
				Token:    "",
				Prefix:   parserErr.prefix,
				Message:  parserErr.Inner.Error(),
			}
		}
	}
	return root, sqlvmErrList
}
