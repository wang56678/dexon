package main

import (
	"flag"
	"fmt"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/parser"
)

func main() {
	var detail bool
	flag.BoolVar(&detail, "detail", false, "print struct detail")

	flag.Parse()

	n, err := parser.ParseString(flag.Arg(0))
	fmt.Printf("detail: %t\n", detail)
	fmt.Printf("err: %+v\n", err)
	if err == nil {
		ast.PrintAST(n, "", detail)
	}
}
