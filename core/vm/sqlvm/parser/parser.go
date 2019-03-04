package parser

import (
	"bytes"
	"fmt"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/parser/internal"
)

func walkSelfFirst(n ast.Node, v func(ast.Node, []ast.Node)) {
	c := n.GetChildren()
	v(n, c)
	for i := range c {
		walkSelfFirst(c[i], v)
	}
}

func walkChildrenFirst(n ast.Node, v func(ast.Node, []ast.Node)) {
	c := n.GetChildren()
	for i := range c {
		walkChildrenFirst(c[i], v)
	}
	v(n, c)
}

// Parse parses SQL commands text and return an AST.
func Parse(b []byte) ([]ast.Node, error) {
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
	options := []internal.Option{internal.Recover(false)}
	root, pigeonErr := internal.Parse("", eb, options...)

	// Process the AST.
	var stmts []ast.Node
	if root != nil {
		stmts = root.([]ast.Node)
	}
	for i := range stmts {
		if stmts[i] == nil {
			continue
		}
		walkChildrenFirst(stmts[i], func(n ast.Node, c []ast.Node) {
			minBegin := uint32(len(eb))
			maxEnd := uint32(0)
			for _, cn := range append(c, n) {
				if cn.HasPosition() {
					begin := cn.GetPosition()
					end := begin + cn.GetLength()
					if begin < minBegin {
						minBegin = begin
					}
					if end > maxEnd {
						maxEnd = end
					}
				}
			}
			n.SetPosition(minBegin)
			n.SetLength(maxEnd - minBegin)
		})
		walkSelfFirst(stmts[i], func(n ast.Node, _ []ast.Node) {
			begin := n.GetPosition()
			end := begin + n.GetLength()
			fixedBegin, ok := encMap[begin]
			if !ok {
				panic(fmt.Sprintf("cannot fix node begin byte offset %d", begin))
			}
			fixedEnd, ok := encMap[end]
			if !ok {
				panic(fmt.Sprintf("cannot fix node end byte offset %d", end))
			}
			n.SetPosition(fixedBegin)
			n.SetLength(fixedEnd - fixedBegin)
		})
	}
	if pigeonErr == nil {
		return stmts, pigeonErr
	}

	// Process errors.
	pigeonErrList := pigeonErr.(internal.ErrList)
	sqlvmErrList := make(errors.ErrorList, len(pigeonErrList))
	for i := range pigeonErrList {
		parserErr := pigeonErrList[i].(*internal.ParserError)
		if sqlvmErr, ok := parserErr.Inner.(errors.Error); ok {
			sqlvmErrList[i] = sqlvmErr
		} else {
			sqlvmErrList[i] = parserErr.SQLVMError()
		}
		sqlvmErrList[i].Token =
			string(internal.DecodeString([]byte(sqlvmErrList[i].Token)))
		if offset, ok := encMap[sqlvmErrList[i].Position]; ok {
			sqlvmErrList[i].Position = offset
		} else {
			panic(fmt.Sprintf("cannot fix error position byte offset %d",
				sqlvmErrList[i].Position))
		}
	}
	return stmts, sqlvmErrList
}
