package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/parser"
)

func main() {
	var detail bool
	flag.BoolVar(&detail, "detail", false, "print struct detail")

	flag.Parse()

	fmt.Fprintf(os.Stderr, "detail: %t\n", detail)
	s := []byte(flag.Arg(0))
	n, parseErr := parser.Parse(s)
	b, printErr := ast.PrintAST(os.Stdout, n, "  ", detail)
	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "Parse error:\n%+v\n", parseErr)
	}
	if printErr != nil {
		fmt.Fprintf(os.Stderr, "Print error:\n%+v\n", printErr)
	}
	fmt.Fprintf(os.Stderr, "Output size: %d bytes\n", b)
	if parseErr != nil || printErr != nil {
		os.Exit(1)
	}
}
