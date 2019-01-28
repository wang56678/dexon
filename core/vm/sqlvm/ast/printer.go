package ast

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func formatBytes(b []byte) string {
	if utf8.Valid(b) {
		for r, i, size := rune(0), 0, 0; i < len(b); i += size {
			r, size = utf8.DecodeRune(b[i:])
			if !unicode.IsPrint(r) {
				return strconv.Quote(string(b))
			}
		}
		return string(b)
	}
	return fmt.Sprintf("%v", b)
}

func formatString(s string) string {
	if utf8.ValidString(s) {
		for _, r := range s {
			if !unicode.IsPrint(r) {
				return strconv.Quote(s)
			}
		}
		return s
	}
	return fmt.Sprintf("%v", []byte(s))
}

// PrintAST prints ast to stdout.
func PrintAST(n interface{}, indent string, detail bool) {
	if n == nil {
		fmt.Printf("%snil\n", indent)
		return
	}
	typeOf := reflect.TypeOf(n)
	valueOf := reflect.ValueOf(n)
	name := ""
	if typeOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			fmt.Printf("%snil\n", indent)
			return
		}
		name = "*"
		valueOf = valueOf.Elem()
		typeOf = typeOf.Elem()
	}
	name = name + typeOf.Name()

	if op, ok := n.(UnaryOperator); ok {
		fmt.Printf("%s%s:\n", indent, name)
		fmt.Printf("%s  Target:\n", indent)
		PrintAST(op.GetTarget(), indent+"    ", detail)
		return
	}
	if op, ok := n.(BinaryOperator); ok {
		fmt.Printf("%s%s:\n", indent, name)
		fmt.Printf("%s  Object:\n", indent)
		PrintAST(op.GetObject(), indent+"    ", detail)
		fmt.Printf("%s  Subject:\n", indent)
		PrintAST(op.GetSubject(), indent+"    ", detail)
		return
	}
	if arr, ok := n.([]interface{}); ok {
		if len(arr) == 0 {
			fmt.Printf("%s[]\n", indent)
			return
		}
		fmt.Printf("%s[\n", indent)
		for idx := range arr {
			PrintAST(arr[idx], indent+"  ", detail)
		}
		fmt.Printf("%s]\n", indent)
		return
	}
	if stringer, ok := n.(fmt.Stringer); ok {
		s := stringer.String()
		fmt.Printf("%s%s: %s\n", indent, name, formatString(s))
		return
	}
	if typeOf.Kind() == reflect.Struct {
		fmt.Printf("%s%s", indent, name)
		l := typeOf.NumField()
		if l == 0 {
			fmt.Printf(" {}\n")
			return
		}
		fmt.Printf(" {\n")
		for i := 0; i < l; i++ {
			if !detail && typeOf.Field(i).Tag.Get("print") == "-" {
				continue
			}
			fmt.Printf("%s  %s:\n", indent, typeOf.Field(i).Name)
			PrintAST(valueOf.Field(i).Interface(), indent+"    ", detail)
		}
		fmt.Printf("%s}\n", indent)
		return
	}
	if bs, ok := n.([]byte); ok {
		fmt.Printf("%s%s\n", indent, formatBytes(bs))
		return
	}
	fmt.Printf("%s%+v\n", indent, valueOf.Interface())
}
