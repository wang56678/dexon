package main

import (
	"flag"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/runtime"
)

func main() {
	var output string
	flag.StringVar(
		&output, "o", "./runtime/instructions_op_test.go",
		"the output path of generated testcases",
	)
	flag.Parse()

	err := runtime.RenderOpTest(output)
	if err != nil {
		panic(err)
	}
}
