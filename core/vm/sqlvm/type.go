package sqlvm

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
)

type identifierNode string

type valuer interface {
	value() interface{}
}

type boolValueNode bool

func (n boolValueNode) value() interface{} { return bool(n) }

type integerValueNode struct {
	v       decimal.Decimal
	address bool
}

func (n integerValueNode) value() interface{} { return n.v }
func (n integerValueNode) String() string     { return n.v.String() }

type decimalValueNode decimal.Decimal

func (n decimalValueNode) value() interface{} { return decimal.Decimal(n) }
func (n decimalValueNode) String() string     { return decimal.Decimal(n).String() }

type bytesValueNode []byte

func (n bytesValueNode) value() interface{} { return []byte(n) }
func (n bytesValueNode) String() string     { return string(n) }

type anyValueNode struct{}

func (n anyValueNode) value() interface{} { return n }

type defaultValueNode struct{}

func (n defaultValueNode) value() interface{} { return n }

type nullValueNode struct{}

func (n nullValueNode) value() interface{} { return n }

type intTypeNode struct {
	size     int32
	unsigned bool
}

type fixedTypeNode struct {
	integerSize    int32
	fractionalSize int32
	unsigned       bool
}

type dynamicBytesTypeNode struct{}

type fixedBytesTypeNode struct {
	size int32
}

type addressTypeNode struct{}
type boolTypeNode struct{}

type unaryOperator interface {
	getTarget() interface{}
	setTarget(interface{})
}

type binaryOperator interface {
	getObject() interface{}
	getSubject() interface{}
	setObject(interface{})
	setSubject(interface{})
}

type unaryOperatorNode struct {
	target interface{}
}

func (n unaryOperatorNode) getTarget() interface{} {
	return n.target
}

func (n *unaryOperatorNode) setTarget(t interface{}) {
	n.target = t
}

type binaryOperatorNode struct {
	object  interface{}
	subject interface{}
}

func (n binaryOperatorNode) getObject() interface{} {
	return n.object
}

func (n binaryOperatorNode) getSubject() interface{} {
	return n.subject
}

func (n *binaryOperatorNode) setObject(o interface{}) {
	n.object = o
}

func (n *binaryOperatorNode) setSubject(s interface{}) {
	n.subject = s
}

type posOperatorNode struct{ unaryOperatorNode }
type negOperatorNode struct{ unaryOperatorNode }
type notOperatorNode struct{ unaryOperatorNode }
type andOperatorNode struct{ binaryOperatorNode }
type orOperatorNode struct{ binaryOperatorNode }

type greaterOrEqualOperatorNode struct{ binaryOperatorNode }
type lessOrEqualOperatorNode struct{ binaryOperatorNode }
type notEqualOperatorNode struct{ binaryOperatorNode }
type equalOperatorNode struct{ binaryOperatorNode }
type greaterOperatorNode struct{ binaryOperatorNode }
type lessOperatorNode struct{ binaryOperatorNode }

type concatOperatorNode struct{ binaryOperatorNode }
type addOperatorNode struct{ binaryOperatorNode }
type subOperatorNode struct{ binaryOperatorNode }
type mulOperatorNode struct{ binaryOperatorNode }
type divOperatorNode struct{ binaryOperatorNode }
type modOperatorNode struct{ binaryOperatorNode }

type inOperatorNode struct{ binaryOperatorNode }
type isOperatorNode struct{ binaryOperatorNode }
type likeOperatorNode struct{ binaryOperatorNode }

type castOperatorNode struct{ binaryOperatorNode }
type assignOperatorNode struct{ binaryOperatorNode }
type functionOperatorNode struct{ binaryOperatorNode }

type optional interface {
	getOption() map[string]interface{}
}

type nilOptionNode struct{}

func (n nilOptionNode) getOption() map[string]interface{} { return nil }

type whereOptionNode struct {
	condition interface{}
}

func (n whereOptionNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"condition": n.condition,
	}
}

type orderOptionNode struct {
	expr      interface{}
	desc      bool
	nullfirst bool
}

func (n orderOptionNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"expr":      n.expr,
		"desc":      n.desc,
		"nullfirst": n.nullfirst,
	}
}

type groupOptionNode struct {
	expr interface{}
}

func (n groupOptionNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"expr": n.expr,
	}
}

type offsetOptionNode struct {
	value integerValueNode
}

func (n offsetOptionNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"value": n.value,
	}
}

type limitOptionNode struct {
	value integerValueNode
}

func (n limitOptionNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"value": n.value,
	}
}

type insertWithColumnOptionNode struct {
	column []interface{}
	value  []interface{}
}

func (n insertWithColumnOptionNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"column": n.column,
		"value":  n.value,
	}
}

type insertWithDefaultOptionNode struct{ nilOptionNode }
type primaryOptionNode struct{ nilOptionNode }
type notNullOptionNode struct{ nilOptionNode }
type uniqueOptionNode struct{ nilOptionNode }
type autoincrementOptionNode struct{ nilOptionNode }

type defaultOptionNode struct {
	value interface{}
}

func (n defaultOptionNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"value": n.value,
	}
}

type foreignOptionNode struct {
	table  identifierNode
	column identifierNode
}

func (n foreignOptionNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"table":  n.table,
		"column": n.column,
	}
}

type selectStmtNode struct {
	column []interface{}
	table  *identifierNode
	where  *whereOptionNode
	group  []interface{}
	order  []interface{}
	limit  *limitOptionNode
	offset *offsetOptionNode
}

func (n selectStmtNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"column": n.column,
		"table":  n.table,
		"where":  n.where,
		"group":  n.group,
		"order":  n.order,
		"limit":  n.limit,
		"offset": n.offset,
	}
}

type updateStmtNode struct {
	table      identifierNode
	assignment []interface{}
	where      *whereOptionNode
}

func (n updateStmtNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"table":      n.table,
		"assignment": n.assignment,
		"where":      n.where,
	}
}

type deleteStmtNode struct {
	table identifierNode
	where *whereOptionNode
}

func (n deleteStmtNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"table": n.table,
		"where": n.where,
	}
}

type insertStmtNode struct {
	table  identifierNode
	insert interface{}
}

func (n insertStmtNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"table":  n.table,
		"insert": n.insert,
	}
}

type createTableStmtNode struct {
	table  identifierNode
	column []interface{}
}

func (n createTableStmtNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"table":  n.table,
		"column": n.column,
	}
}

type columnSchemaNode struct {
	column     identifierNode
	dataType   interface{}
	constraint []interface{}
}

func (n columnSchemaNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"column":     n.column,
		"data_type":  n.dataType,
		"constraint": n.constraint,
	}
}

type createIndexStmtNode struct {
	index  identifierNode
	table  identifierNode
	column []interface{}
	unique *uniqueOptionNode
}

func (n createIndexStmtNode) getOption() map[string]interface{} {
	return map[string]interface{}{
		"index":  n.index,
		"table":  n.table,
		"column": n.column,
		"unique": n.unique,
	}
}

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

	if op, ok := n.(optional); ok {
		fmt.Printf("%s%s", indent, name)
		m := op.getOption()
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
	if op, ok := n.(unaryOperator); ok {
		fmt.Printf("%s%s:\n", indent, name)
		fmt.Printf("%s  target:\n", indent)
		PrintAST(op.getTarget(), indent+"    ")
		return
	}
	if op, ok := n.(binaryOperator); ok {
		fmt.Printf("%s%s:\n", indent, name)
		fmt.Printf("%s  object:\n", indent)
		PrintAST(op.getObject(), indent+"    ")
		fmt.Printf("%s  subject:\n", indent)
		PrintAST(op.getSubject(), indent+"    ")
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
