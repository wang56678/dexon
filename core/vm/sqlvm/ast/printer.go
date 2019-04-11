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

func printAST(w io.Writer, n interface{}, prefix string, detail bool,
	depth int) (int, error) {

	indent := strings.Repeat(prefix, depth)
	indentLong := strings.Repeat(prefix, depth+1)
	if n == nil {
		return fmt.Fprintf(w, "%snil\n", indent)
	}
	typeOf := reflect.TypeOf(n)
	valueOf := reflect.ValueOf(n)
	kind := typeOf.Kind()
	if kind == reflect.Ptr {
		if valueOf.IsNil() {
			return fmt.Fprintf(w, "%snil\n", indent)
		}
		valueOf = valueOf.Elem()
		typeOf = typeOf.Elem()
		kind = typeOf.Kind()
	}
	name := typeOf.Name()

	if stringer, ok := n.(fmt.Stringer); ok {
		s := stringer.String()
		return fmt.Fprintf(w, "%s%s\n", indent, formatString(s))
	}
	if s, ok := n.(string); ok {
		return fmt.Fprintf(w, "%s%s\n", indent, formatString(s))
	}
	if bs, ok := n.([]byte); ok {
		return fmt.Fprintf(w, "%s%s\n", indent, formatBytes(bs))
	}
	if kind == reflect.Slice {
		l := valueOf.Len()
		if l == 0 {
			return fmt.Fprintf(w, "%s[]\n", indent)
		}

		var bytesWritten int
		b, err := fmt.Fprintf(w, "%s[\n", indent)
		bytesWritten += b
		if err != nil {
			return bytesWritten, err
		}
		for i := 0; i < l; i++ {
			v := valueOf.Index(i)
			b, err = printAST(w, v.Interface(), prefix, detail, depth+1)
			bytesWritten += b
			if err != nil {
				return bytesWritten, err
			}
		}
		b, err = fmt.Fprintf(w, "%s]\n", indent)
		bytesWritten += b
		return bytesWritten, err
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

		var position string
		if node, ok := n.(Node); ok {
			begin := node.GetPosition()
			length := node.GetLength()
			if node.HasPosition() {
				end := begin + length - 1
				token := node.GetToken()
				position = fmt.Sprintf("%d-%d %s",
					begin, end, strconv.Quote(string(token)))
			} else {
				position = "no position info"
			}
		}

		var bytesWritten int
		b, err := fmt.Fprintf(w, "%s%s", indent, name)
		bytesWritten += b
		if err != nil {
			return bytesWritten, err
		}
		if len(fields) == 0 {
			b, err = fmt.Fprintf(w, " {}  // %s\n", position)
			bytesWritten += b
			return bytesWritten, err
		}
		b, err = fmt.Fprintf(w, " {  // %s\n", position)
		bytesWritten += b
		if err != nil {
			return bytesWritten, err
		}
		for i := 0; i < len(fields); i++ {
			b, err = fmt.Fprintf(w, "%s%s:\n", indentLong, fields[i].name)
			bytesWritten += b
			if err != nil {
				return bytesWritten, err
			}
			b, err = printAST(w, fields[i].value, prefix, detail, depth+2)
			bytesWritten += b
			if err != nil {
				return bytesWritten, err
			}
		}
		b, err = fmt.Fprintf(w, "%s}\n", indent)
		bytesWritten += b
		return bytesWritten, err
	}
	return fmt.Fprintf(w, "%s%+v\n", indent, valueOf.Interface())
}

// PrintAST prints AST for debugging.
func PrintAST(output io.Writer, node interface{}, indent string, detail bool) (
	int, error) {

	return printAST(output, node, indent, detail, 0)
}
