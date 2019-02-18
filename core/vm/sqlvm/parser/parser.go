package parser

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

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
	errStr := err.Error()
	if strings.HasSuffix(errStr, "decimal: fractional part too long") {
		return errors.ErrorCodeFractionalPartTooLong
	} else if strings.HasSuffix(errStr, "decimal: exponent is not numeric") {
		return errors.ErrorCodeInvalidNumberSyntax
	} else if strings.HasSuffix(errStr, "decimal: too many .s") {
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
				if r >= 0x10ffff || (r >= 0xd800 && r <= 0xdfff) {
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

// Parse parses SQL commands text and return an AST.
func Parse(b []byte, o ...Option) (interface{}, error) {
	// The string sent from the caller is not guaranteed to be valid UTF-8.
	// We don't really care non-ASCII characters in the string because all
	// keywords and special symbols are defined in ASCII. Therefore, as long
	// as the encoding is compatible with ASCII, we can process text with
	// unknown encoding.
	//
	// However, pigeon requires input text to be valid UTF-8, throwing an error
	// and exiting early when it cannot decode the input as UTF-8. In order to
	// workaround it, we preprocess the input text by assuming each byte value
	// is a Unicode code point and encoding the input text as UTF-8.
	//
	// This means that the byte offset reported by pigeon is wrong. We have to
	// scan the the error list and the AST to fix positions in these structs
	// before returning them to the caller.

	// Encode the input text.
	encBuf := bytes.Buffer{}
	encMap := map[uint32]uint32{}
	for i, c := range b {
		encMap[uint32(encBuf.Len())] = uint32(i)
		encBuf.WriteRune(rune(c))
	}
	encMap[uint32(encBuf.Len())] = uint32(len(b))

	// Prepare arguments and call the parser.
	eb := encBuf.Bytes()
	options := append([]Option{Recover(false)}, o...)
	root, pigeonErr := parse("", eb, options...)

	// Process the AST.
	if pigeonErr == nil {
		return root, pigeonErr
	}

	// Process errors.
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
		sqlvmErrList[i].Token =
			string(decodeString([]byte(sqlvmErrList[i].Token)))
		if offset, ok := encMap[sqlvmErrList[i].Position]; ok {
			sqlvmErrList[i].Position = offset
		} else {
			panic(fmt.Sprintf(
				"cannot fix byte offset %d", sqlvmErrList[i].Position))
		}
	}
	return root, sqlvmErrList
}
