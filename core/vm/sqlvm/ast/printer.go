package ast

import (
	"fmt"
	"reflect"
)

// PrintAST prints ast to stdout.
func PrintAST(n interface{}, indent string) {
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

	if op, ok := n.(Optional); ok {
		fmt.Printf("%s%s", indent, name)
		m := op.GetOption()
		if m == nil {
			fmt.Printf("\n")
			return
		}
		fmt.Printf(":\n")
		for k := range m {
			fmt.Printf("%s  %s:\n", indent, k)
			PrintAST(m[k], indent+"    ")
		}
		return
	}
	if op, ok := n.(UnaryOperator); ok {
		fmt.Printf("%s%s:\n", indent, name)
		fmt.Printf("%s  Target:\n", indent)
		PrintAST(op.GetTarget(), indent+"    ")
		return
	}
	if op, ok := n.(BinaryOperator); ok {
		fmt.Printf("%s%s:\n", indent, name)
		fmt.Printf("%s  Object:\n", indent)
		PrintAST(op.GetObject(), indent+"    ")
		fmt.Printf("%s  Subject:\n", indent)
		PrintAST(op.GetSubject(), indent+"    ")
		return
	}
	if arr, ok := n.([]interface{}); ok {
		fmt.Printf("%s[\n", indent)
		for idx := range arr {
			PrintAST(arr[idx], indent+"  ")
		}
		fmt.Printf("%s]\n", indent)
		return
	}
	if stringer, ok := n.(fmt.Stringer); ok {
		fmt.Printf("%s%s: %s\n", indent, name, stringer.String())
		return
	}
	fmt.Printf("%s%s: %+v\n", indent, name, valueOf.Interface())
}
