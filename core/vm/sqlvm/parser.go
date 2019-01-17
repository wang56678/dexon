package sqlvm

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// Parser was generated with pigeon v1.0.0-99-gbb0192c.
//go:generate pigeon -no-recover -o grammar.go grammar.peg
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
		i, _ := strconv.ParseInt(string([]byte{b}), 16, 32)
		d = d.Add(
			decimal.New(i, 0).
				Mul(base.Pow(decimal.New(int64(l-idx-1), 0))),
		)
	}
	return integerValueNode{v: d, address: isAddress(h)}
}

func hexToBytes(h []byte) []byte {
	bs, _ := hex.DecodeString(string(h))
	return bs
}

func toInt(b []byte) int32 {
	i, _ := strconv.ParseInt(string(b), 10, 32)
	return int32(i)
}

func toDecimal(b []byte) decimal.Decimal {
	return decimal.RequireFromString(string(b))
}

func toLower(b []byte) []byte {
	return []byte(strings.ToLower(string(b)))
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
	x := op.(binaryOperator)
	x.setSubject(s)
	return x
}

func opSetObject(op interface{}, o interface{}) interface{} {
	x := op.(binaryOperator)
	x.setObject(o)
	return x
}

func opSetTarget(op interface{}, t interface{}) interface{} {
	x := op.(unaryOperator)
	x.setTarget(t)
	return x
}

func joinOperator(x interface{}, o interface{}) {
	if op, ok := x.(unaryOperator); ok {
		joinOperator(op.getTarget(), o)
		return
	}
	if op, ok := x.(binaryOperator); ok {
		op.setObject(o)
		return
	}
}

func rightJoinOperators(o interface{}, x interface{}) interface{} {
	xs := toSlice(x)
	if len(xs) == 0 {
		return o
	}
	l := len(xs)
	for idx := range xs {
		if idx == l-1 {
			break
		}
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
	return ParseReader("parser", strings.NewReader(s))
}
