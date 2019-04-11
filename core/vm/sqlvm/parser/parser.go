package parser

import (
	"bytes"
	"fmt"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/parser/internal"
)

type visitor func(ast.Node, []ast.Node)

func walkSelfFirst(n ast.Node, v visitor) bool {
	return walkSelfFirstWithDepth(n, v, 0)
}

func walkSelfFirstWithDepth(n ast.Node, v visitor, d int) bool {
	if d >= ast.DepthLimit {
		return false
	}
	c := n.GetChildren()
	r := true
	v(n, c)
	for i := range c {
		r = r && walkSelfFirstWithDepth(c[i], v, d+1)
	}
	return r
}

func walkChildrenFirst(n ast.Node, v visitor) bool {
	return walkChildrenFirstWithDepth(n, v, 0)
}

func walkChildrenFirstWithDepth(n ast.Node, v visitor, d int) bool {
	if d >= ast.DepthLimit {
		return false
	}
	c := n.GetChildren()
	r := true
	for i := range c {
		r = r && walkChildrenFirstWithDepth(c[i], v, d+1)
	}
	v(n, c)
	return r
}

// Parse parses SQL commands text and return an AST.
func Parse(b []byte) ([]ast.StmtNode, error) {
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

	// Copy the input text. We will put references to the source code on AST
	// nodes, so we have to make our own copy to prevent the AST from being
	// broken by the caller if the input byte slice was modified afterwards.
	b = append([]byte{}, b...)

	// Process the AST.
	var stmts []ast.StmtNode
	if root != nil {
		stmts = root.([]ast.StmtNode)
	}
	for i := range stmts {
		if stmts[i] == nil {
			continue
		}
		r := true
		r = r && walkChildrenFirst(stmts[i], func(n ast.Node, c []ast.Node) {
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
		r = r && walkSelfFirst(stmts[i], func(n ast.Node, _ []ast.Node) {
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
			n.SetToken(b[fixedBegin:fixedEnd])
		})
		if !r {
			return nil, errors.ErrorList{
				errors.Error{
					Position: 0,
					Length:   0,
					Category: errors.ErrorCategoryLimit,
					Code:     errors.ErrorCodeDepthLimitReached,
					Severity: errors.ErrorSeverityError,
					Prefix:   "",
					Message: fmt.Sprintf("reach syntax tree depth limit %d",
						ast.DepthLimit),
				},
			}
		}
	}
	if pigeonErr == nil {
		return stmts, nil
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
		begin := sqlvmErrList[i].Position
		end := begin + sqlvmErrList[i].Length
		fixedBegin, ok := encMap[begin]
		if !ok {
			panic(fmt.Sprintf("cannot fix error position byte offset %d", begin))
		}
		fixedEnd, ok := encMap[end]
		if !ok {
			panic(fmt.Sprintf("cannot fix error position byte offset %d", end))
		}
		sqlvmErrList[i].Position = fixedBegin
		sqlvmErrList[i].Length = fixedEnd - fixedBegin
	}
	return stmts, sqlvmErrList
}
