package ast

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
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

func printAST(w io.Writer, n interface{}, depth int, base string, detail bool) {
	indent := strings.Repeat(base, depth)
	indentLong := strings.Repeat(base, depth+1)
	if n == nil {
		fmt.Fprintf(w, "%snil\n", indent)
		return
	}
	typeOf := reflect.TypeOf(n)
	valueOf := reflect.ValueOf(n)
	kind := typeOf.Kind()
	if kind == reflect.Ptr {
		if valueOf.IsNil() {
			fmt.Fprintf(w, "%snil\n", indent)
			return
		}
		valueOf = valueOf.Elem()
		typeOf = typeOf.Elem()
		kind = typeOf.Kind()
	}
	name := typeOf.Name()

	if op, ok := n.(UnaryOperator); ok {
		fmt.Fprintf(w, "%s%s:\n", indent, name)
		fmt.Fprintf(w, "%sTarget:\n", indentLong)
		printAST(w, op.GetTarget(), depth+2, base, detail)
		return
	}
	if op, ok := n.(BinaryOperator); ok {
		fmt.Fprintf(w, "%s%s:\n", indent, name)
		fmt.Fprintf(w, "%sObject:\n", indentLong)
		printAST(w, op.GetObject(), depth+2, base, detail)
		fmt.Fprintf(w, "%sSubject:\n", indentLong)
		printAST(w, op.GetSubject(), depth+2, base, detail)
		return
	}
	if stringer, ok := n.(fmt.Stringer); ok {
		s := stringer.String()
		fmt.Fprintf(w, "%s%s\n", indent, formatString(s))
		return
	}
	if s, ok := n.(string); ok {
		fmt.Fprintf(w, "%s%s\n", indent, formatString(s))
		return
	}
	if bs, ok := n.([]byte); ok {
		fmt.Fprintf(w, "%s%s\n", indent, formatBytes(bs))
		return
	}
	if kind == reflect.Slice {
		l := valueOf.Len()
		if l == 0 {
			fmt.Fprintf(w, "%s[]\n", indent)
			return
		}
		fmt.Fprintf(w, "%s[\n", indent)
		for i := 0; i < l; i++ {
			v := valueOf.Index(i)
			printAST(w, v.Interface(), depth+1, base, detail)
		}
		fmt.Fprintf(w, "%s]\n", indent)
		return
	}
	if kind == reflect.Struct {
		type field struct {
			name  string
			value interface{}
		}
		var fields []field
		var collect func(reflect.Type, reflect.Value)
		collect = func(typeOf reflect.Type, valueOf reflect.Value) {
			l := typeOf.NumField()
			for i := 0; i < l; i++ {
				if !detail && typeOf.Field(i).Tag.Get("print") == "-" {
					continue
				}
				if typeOf.Field(i).Anonymous {
					embeddedInterface := valueOf.Field(i).Interface()
					embeddedTypeOf := reflect.TypeOf(embeddedInterface)
					embeddedValueOf := reflect.ValueOf(embeddedInterface)
					collect(embeddedTypeOf, embeddedValueOf)
					continue
				}
				fields = append(fields, field{
					name:  typeOf.Field(i).Name,
					value: valueOf.Field(i).Interface(),
				})
			}
		}
		collect(typeOf, valueOf)
		fmt.Fprintf(w, "%s%s", indent, name)
		if len(fields) == 0 {
			fmt.Fprintf(w, " {}\n")
			return
		}
		fmt.Fprintf(w, " {\n")
		for i := 0; i < len(fields); i++ {
			fmt.Fprintf(w, "%s%s:\n", indentLong, fields[i].name)
			printAST(w, fields[i].value, depth+2, base, detail)
		}
		fmt.Fprintf(w, "%s}\n", indent)
		return
	}
	fmt.Fprintf(w, "%s%+v\n", indent, valueOf.Interface())
}

// PrintAST prints AST for debugging.
func PrintAST(w io.Writer, n interface{}, indent string, detail bool) {
	printAST(w, n, 0, indent, detail)
}
