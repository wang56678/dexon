package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/checker"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/parser"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
	"github.com/dexon-foundation/dexon/rlp"
)

func create(sql string, o checker.CheckOptions) int {
	n, parseErr := parser.Parse([]byte(sql))
	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "Parse error:\n%+v\n", parseErr)
	}
	s, checkErr := checker.CheckCreate(n, o)
	if checkErr != nil {
		fmt.Fprintf(os.Stderr, "Check error:\n%+v\n", checkErr)
	}
	b := bytes.Buffer{}
	rlpErr := rlp.Encode(&b, s)
	if rlpErr != nil {
		fmt.Fprintf(os.Stderr, "RLP encode error: %v\n", rlpErr)
		return 1
	}
	fmt.Println(hex.EncodeToString(b.Bytes()))
	if parseErr != nil || checkErr != nil {
		return 1
	}
	return 0
}

func decode(ss string) int {
	b, hexErr := hex.DecodeString(ss)
	if hexErr != nil {
		fmt.Fprintf(os.Stderr, "Hex decode error: %v\n", hexErr)
		return 1
	}
	s := schema.Schema{}
	rlpErr := rlp.Decode(bytes.NewReader(b), &s)
	if rlpErr != nil {
		fmt.Fprintf(os.Stderr, "RLP decode error: %v\n", rlpErr)
		return 1
	}
	s.SetupColumnOffset()
	fmt.Print(s.String())
	return 0
}

func query(ss, sql string, o checker.CheckOptions) int {
	fmt.Fprintln(os.Stderr, "Function not implemented")
	return 1
}

func exec(ss, sql string, o checker.CheckOptions) int {
	fmt.Fprintln(os.Stderr, "Function not implemented")
	return 1
}

func main() {
	var safeMath bool
	var safeCast bool
	var constantOnly bool
	flag.BoolVar(&safeMath, "safe-math", true, "")
	flag.BoolVar(&safeCast, "safe-cast", true, "")
	flag.BoolVar(&constantOnly, "constant-only", false, " (default false)")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] <action> <arguments>\n"+
				"Options:\n"+
				"  -help    Show options\n"+
				"Actions:\n"+
				"  create   (SQL) -> schema\n"+
				"  decode   (schema) -> SQL\n"+
				"  query    (schema, SQL) -> AST\n"+
				"  exec     (schema, SQL) -> AST\n",
			os.Args[0])
		os.Exit(1)
	}

	var o checker.CheckOptions
	if safeMath {
		o |= checker.CheckWithSafeMath
	}
	if safeCast {
		o |= checker.CheckWithSafeCast
	}
	if constantOnly {
		o |= checker.CheckWithConstantOnly
	}

	action := flag.Arg(0)
	switch action {
	case "create":
		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "create needs 1 argument")
			os.Exit(1)
		}
		os.Exit(create(flag.Arg(1), o))
	case "decode":
		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "decode needs 1 argument")
			os.Exit(1)
		}
		os.Exit(decode(flag.Arg(1)))
	case "query":
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "query needs 2 arguments")
			os.Exit(1)
		}
		os.Exit(query(flag.Arg(1), flag.Arg(2), o))
	case "exec":
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "exec needs 2 arguments")
			os.Exit(1)
		}
		os.Exit(exec(flag.Arg(1), flag.Arg(2), o))
	default:
		fmt.Fprintf(os.Stderr, "Invalid action %s\n", action)
		os.Exit(1)
	}
}
