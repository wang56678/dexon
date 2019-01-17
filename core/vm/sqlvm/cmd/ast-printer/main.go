package main

import (
	"fmt"
	"os"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm"
)

func main() {
	n, err := sqlvm.ParseString(os.Args[1])
	fmt.Printf("err: %+v\n", err)
	if err == nil {
		sqlvm.PrintAST(n, "")
	}
}
