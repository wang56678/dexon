package main

import (
	"fmt"
	"os"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/parser"
)

func main() {
	n, err := parser.ParseString(os.Args[1])
	fmt.Printf("err: %+v\n", err)
	if err == nil {
		ast.PrintAST(n, "")
	}
}
