package internal

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

// Parser was generated with pigeon v1.0.0-99-gbb0192c.
//go:generate pigeon -no-recover -o grammar.go grammar.peg

func prepend(x interface{}, xs []interface{}) []interface{} {
	return append([]interface{}{x}, xs...)
}

func assertSlice(x interface{}) []interface{} {
	if x == nil {
		return nil
	}
	return x.([]interface{})
}

func assertExprSlice(x interface{}) []ast.ExprNode {
	xs := assertSlice(x)
	es := make([]ast.ExprNode, len(xs))
	for i := 0; i < len(xs); i++ {
		if xs[i] != nil {
			es[i] = xs[i].(ast.ExprNode)
		}
	}
	return es
}

func isAddress(h []byte) bool {
	ma, err := common.NewMixedcaseAddressFromString(string(h))
	if err != nil {
		return false
	}
	return ma.ValidChecksum()
}

func hexToInteger(h []byte) *ast.IntegerValueNode {
	d := decimal.Zero
	x := h[2:]
	l := len(x)
	base := decimal.New(16, 0)
	for idx, b := range x {
		i, err := strconv.ParseInt(string([]byte{b}), 16, 32)
		if err != nil {
			panic(fmt.Sprintf("invalid hex digit %s: %v", []byte{b}, err))
		}
		d = d.Add(
			decimal.New(i, 0).
				Mul(base.Pow(decimal.New(int64(l-idx-1), 0))),
		)
	}
	node := &ast.IntegerValueNode{}
	node.IsAddress = isAddress(h)
	node.V = d
	return node
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
		return errors.ErrorCodeInvalidIntegerSyntax
	case strconv.ErrRange:
		return errors.ErrorCodeIntegerOutOfRange
	}
	panic(fmt.Sprintf("unknown NumError: %v", err))
}

func convertDecimalError(err error) errors.ErrorCode {
	if err == nil {
		return errors.ErrorCodeNil
	}
	switch err.(type) {
	case *decimal.ErrorExponentLimit:
		return errors.ErrorCodeFractionalPartTooLong
	case *decimal.ErrorInvalidFormat:
		return errors.ErrorCodeInvalidNumberSyntax
	}
	panic(fmt.Sprintf("unknown decimal error: %v", err))
}

func toUint(b []byte) (uint32, errors.ErrorCode) {
	i, err := strconv.ParseUint(string(b), 10, 32)
	return uint32(i), convertNumError(err)
}

func toDecimal(b []byte) (decimal.Decimal, errors.ErrorCode) {
	if len(b) > 0 && b[0] == byte('.') {
		b = append([]byte{'0'}, b...)
	}
	d, err := decimal.NewFromString(string(b))
	return d, convertDecimalError(err)
}

func joinBytes(x []interface{}) []byte {
	bs := []byte{}
	for _, b := range x {
		if b != nil {
			bs = append(bs, b.([]byte)...)
		}
	}
	return bs
}

func opSetSubject(op ast.BinaryOperator, s ast.ExprNode) ast.BinaryOperator {
	op.SetSubject(s)
	return op
}

func opSetObject(op ast.BinaryOperator, o ast.ExprNode) ast.BinaryOperator {
	op.SetObject(o)
	return op
}

func opSetTarget(op ast.UnaryOperator, t ast.ExprNode) ast.UnaryOperator {
	op.SetTarget(t)
	return op
}

func joinOperator(x ast.ExprNode, o ast.ExprNode) {
	switch op := x.(type) {
	case ast.UnaryOperator:
		joinOperator(op.GetTarget(), o)
	case ast.BinaryOperator:
		opSetObject(op, o)
	case *ast.InOperatorNode:
		op.Left = o
	default:
		_ = "suppress go-sumtype error"
		panic(fmt.Sprintf("unable to join operators %T and %T", x, o))
	}
}

func rightJoinOperators(o ast.ExprNode, x []ast.ExprNode) ast.ExprNode {
	if len(x) == 0 {
		return o
	}
	l := len(x)
	joinOperator(x[0], o)
	for idx := 0; idx < l-1; idx++ {
		joinOperator(x[idx+1], x[idx])
	}
	return x[l-1]
}

func sanitizeBadEscape(s []byte) []byte {
	o := bytes.Buffer{}
	for _, b := range s {
		if b >= 0x20 && b <= 0x7e && b != '\'' {
			o.WriteByte(b)
		} else {
			o.WriteString(fmt.Sprintf("<%02X>", b))
		}
	}
	return o.Bytes()
}

func decodeString(s []byte) []byte {
	o := bytes.Buffer{}
	for r, i, size := rune(0), 0, 0; i < len(s); i += size {
		r, size = utf8.DecodeRune(s[i:])
		if r > 0xff {
			panic(fmt.Sprintf("invalid encoded rune U+%04X", r))
		}
		o.WriteByte(byte(r))
	}
	return o.Bytes()
}

func resolveString(s []byte) ([]byte, []byte, errors.ErrorCode) {
	s = decodeString(s)
	o := bytes.Buffer{}
	for i, size := 0, 0; i < len(s); i += size {
		if s[i] == '\\' {
			if i+1 >= len(s) {
				panic("trailing backslash in string literal")
			}
			switch s[i+1] {
			case '\n':
				size = 2

			case '\\':
				o.WriteByte('\\')
				size = 2
			case '\'':
				o.WriteByte('\'')
				size = 2
			case '"':
				o.WriteByte('"')
				size = 2
			case 'b':
				o.WriteByte('\b')
				size = 2
			case 'f':
				o.WriteByte('\f')
				size = 2
			case 'n':
				o.WriteByte('\n')
				size = 2
			case 'r':
				o.WriteByte('\r')
				size = 2
			case 't':
				o.WriteByte('\t')
				size = 2
			case 'v':
				o.WriteByte('\v')
				size = 2

			case 'x':
				if i+3 >= len(s) {
					return nil, s[i:], errors.ErrorCodeEscapeSequenceTooShort
				}
				b, err := strconv.ParseUint(string(s[i+2:i+4]), 16, 8)
				if err != nil {
					return nil, s[i : i+4], convertNumError(err)
				}
				o.WriteByte(uint8(b))
				size = 4

			case 'u':
				if i+5 >= len(s) {
					return nil, s[i:], errors.ErrorCodeEscapeSequenceTooShort
				}
				u, err := strconv.ParseUint(string(s[i+2:i+6]), 16, 16)
				if err != nil {
					return nil, s[i : i+6], convertNumError(err)
				}
				if u >= 0xd800 && u <= 0xdfff {
					return nil, s[i : i+6], errors.ErrorCodeInvalidUnicodeCodePoint
				}
				o.WriteRune(rune(u))
				size = 6

			case 'U':
				if i+9 >= len(s) {
					return nil, s[i:], errors.ErrorCodeEscapeSequenceTooShort
				}
				r, err := strconv.ParseUint(string(s[i+2:i+10]), 16, 32)
				if err != nil {
					return nil, s[i : i+10], convertNumError(err)
				}
				if r > 0x10ffff || (r >= 0xd800 && r <= 0xdfff) {
					return nil, s[i : i+10], errors.ErrorCodeInvalidUnicodeCodePoint
				}
				o.WriteRune(rune(r))
				size = 10

			default:
				return nil, s[i : i+2], errors.ErrorCodeUnknownEscapeSequence
			}
		} else {
			o.WriteByte(s[i])
			size = 1
		}
	}
	return o.Bytes(), nil, errors.ErrorCodeNil
}

// ErrList exports pigeon errList for the parser package.
type ErrList = errList

// ParserError exports pigeon parserError for parser package.
type ParserError = parserError

// SQLVMError converts a pigeon parserError to a SQLVM error.
func (p *parserError) SQLVMError() errors.Error {
	return errors.Error{
		Position: uint32(p.pos.offset),
		Length:   0,
		Category: errors.ErrorCategoryGrammar,
		Code:     errors.ErrorCodeParser,
		Severity: errors.ErrorSeverityError,
		Prefix:   p.prefix,
		Message:  p.Inner.Error(),
	}
}
