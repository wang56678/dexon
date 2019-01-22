package ast

import (
	"github.com/shopspring/decimal"
)

// ---------------------------------------------------------------------------
// Identifiers
// ---------------------------------------------------------------------------

// IdentifierNode references table, column, or function.
type IdentifierNode struct {
	Name []byte
}

// ---------------------------------------------------------------------------
// Values
// ---------------------------------------------------------------------------

// Valuer defines the interface of a constant value.
type Valuer interface {
	Value() interface{}
}

// BoolValueNode is a boolean constant.
type BoolValueNode struct {
	V bool
}

// Value returns the value of BoolValueNode.
func (n BoolValueNode) Value() interface{} { return n.V }

// IntegerValueNode is an integer constant.
type IntegerValueNode struct {
	IsAddress bool
	V         decimal.Decimal
}

// Value returns the value of IntegerValueNode.
func (n IntegerValueNode) Value() interface{} { return n.V }
func (n IntegerValueNode) String() string     { return n.V.String() }

// DecimalValueNode is a number constant.
type DecimalValueNode struct {
	V decimal.Decimal
}

// Value returns the value of DecimalValueNode.
func (n DecimalValueNode) Value() interface{} { return n.V }
func (n DecimalValueNode) String() string     { return n.V.String() }

// BytesValueNode is a dynamic or fixed bytes constant.
type BytesValueNode struct {
	V []byte
}

// Value returns the value of BytesValueNode.
func (n BytesValueNode) Value() interface{} { return n.V }
func (n BytesValueNode) String() string     { return string(n.V) }

// AnyValueNode is '*' used in SELECT and function call.
type AnyValueNode struct{}

// Value returns itself.
func (n AnyValueNode) Value() interface{} { return n }

// DefaultValueNode represents the default value used in INSERT and UPDATE.
type DefaultValueNode struct{}

// Value returns itself.
func (n DefaultValueNode) Value() interface{} { return n }

// NullValueNode is NULL.
type NullValueNode struct{}

// Value returns itself.
func (n NullValueNode) Value() interface{} { return n }

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// IntTypeNode represents solidity int{X} and uint{X} types.
type IntTypeNode struct {
	Unsigned bool
	Size     uint32
}

// FixedTypeNode represents solidity fixed{M}x{N} and ufixed{M}x{N} types.
type FixedTypeNode struct {
	Unsigned         bool
	Size             uint32
	FractionalDigits uint32
}

// DynamicBytesTypeNode represents solidity bytes type.
type DynamicBytesTypeNode struct{}

// FixedBytesTypeNode represents solidity bytes{X} type.
type FixedBytesTypeNode struct {
	Size uint32
}

// AddressTypeNode represents solidity address type.
type AddressTypeNode struct{}

// BoolTypeNode represents solidity bool type.
type BoolTypeNode struct{}

// ---------------------------------------------------------------------------
// Operators
// ---------------------------------------------------------------------------

// UnaryOperator defines the interface of a unary operator.
type UnaryOperator interface {
	GetTarget() interface{}
	SetTarget(interface{})
}

// BinaryOperator defines the interface of a binary operator.
type BinaryOperator interface {
	GetObject() interface{}
	GetSubject() interface{}
	SetObject(interface{})
	SetSubject(interface{})
}

// UnaryOperatorNode is a base struct of unary operators.
type UnaryOperatorNode struct {
	Target interface{}
}

// GetTarget gets the target of the operation.
func (n UnaryOperatorNode) GetTarget() interface{} {
	return n.Target
}

// SetTarget sets the target of the operation.
func (n *UnaryOperatorNode) SetTarget(t interface{}) {
	n.Target = t
}

// BinaryOperatorNode is a base struct of binary operators.
type BinaryOperatorNode struct {
	Object  interface{}
	Subject interface{}
}

// GetObject gets the node on which the operation is applied.
func (n BinaryOperatorNode) GetObject() interface{} {
	return n.Object
}

// GetSubject gets the node whose value is applied on the object.
func (n BinaryOperatorNode) GetSubject() interface{} {
	return n.Subject
}

// SetObject sets the object of the operation.
func (n *BinaryOperatorNode) SetObject(o interface{}) {
	n.Object = o
}

// SetSubject sets the subject of the operation.
func (n *BinaryOperatorNode) SetSubject(s interface{}) {
	n.Subject = s
}

// PosOperatorNode is '+'.
type PosOperatorNode struct{ UnaryOperatorNode }

// NegOperatorNode is '-'.
type NegOperatorNode struct{ UnaryOperatorNode }

// NotOperatorNode is 'NOT'.
type NotOperatorNode struct{ UnaryOperatorNode }

// AndOperatorNode is 'AND'.
type AndOperatorNode struct{ BinaryOperatorNode }

// OrOperatorNode is 'OR'.
type OrOperatorNode struct{ BinaryOperatorNode }

// GreaterOrEqualOperatorNode is '>='.
type GreaterOrEqualOperatorNode struct{ BinaryOperatorNode }

// LessOrEqualOperatorNode is '<='.
type LessOrEqualOperatorNode struct{ BinaryOperatorNode }

// NotEqualOperatorNode is '<>'.
type NotEqualOperatorNode struct{ BinaryOperatorNode }

// EqualOperatorNode is '=' used in expressions.
type EqualOperatorNode struct{ BinaryOperatorNode }

// GreaterOperatorNode is '>'.
type GreaterOperatorNode struct{ BinaryOperatorNode }

// LessOperatorNode is '<'.
type LessOperatorNode struct{ BinaryOperatorNode }

// ConcatOperatorNode is '||'.
type ConcatOperatorNode struct{ BinaryOperatorNode }

// AddOperatorNode is '+'.
type AddOperatorNode struct{ BinaryOperatorNode }

// SubOperatorNode is '-'.
type SubOperatorNode struct{ BinaryOperatorNode }

// MulOperatorNode is '*'.
type MulOperatorNode struct{ BinaryOperatorNode }

// DivOperatorNode is '/'.
type DivOperatorNode struct{ BinaryOperatorNode }

// ModOperatorNode is '%'.
type ModOperatorNode struct{ BinaryOperatorNode }

// InOperatorNode is 'IN'.
type InOperatorNode struct{ BinaryOperatorNode }

// IsOperatorNode is 'IS NULL'.
type IsOperatorNode struct{ BinaryOperatorNode }

// LikeOperatorNode is 'LIKE'.
type LikeOperatorNode struct{ BinaryOperatorNode }

// CastOperatorNode is 'CAST(expr AS type)'.
type CastOperatorNode struct{ BinaryOperatorNode }

// AssignOperatorNode is '=' used in UPDATE to set values.
type AssignOperatorNode struct{ BinaryOperatorNode }

// FunctionOperatorNode is a function call.
type FunctionOperatorNode struct{ BinaryOperatorNode }

// ---------------------------------------------------------------------------
// Options
// ---------------------------------------------------------------------------

// Optional defines the interface for printing AST.
type Optional interface {
	GetOption() map[string]interface{}
}

// NilOptionNode is a base struct implementing Optional interface.
type NilOptionNode struct{}

// GetOption returns a value for printing AST.
func (n NilOptionNode) GetOption() map[string]interface{} { return nil }

// WhereOptionNode is 'WHERE' used in SELECT, UPDATE, DELETE.
type WhereOptionNode struct {
	Condition interface{}
}

// GetOption returns a value for printing AST.
func (n WhereOptionNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Condition": n.Condition,
	}
}

// OrderOptionNode is an expression in 'ORDER BY' used in SELECT.
type OrderOptionNode struct {
	Expr       interface{}
	Desc       bool
	NullsFirst bool
}

// GetOption returns a value for printing AST.
func (n OrderOptionNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Expr":       n.Expr,
		"Desc":       n.Desc,
		"NullsFirst": n.NullsFirst,
	}
}

// GroupOptionNode is 'GROUP BY' used in SELECT.
type GroupOptionNode struct {
	Expr interface{}
}

// GetOption returns a value for printing AST.
func (n GroupOptionNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Expr": n.Expr,
	}
}

// OffsetOptionNode is 'OFFSET' used in SELECT.
type OffsetOptionNode struct {
	Value IntegerValueNode
}

// GetOption returns a value for printing AST.
func (n OffsetOptionNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Value": n.Value,
	}
}

// LimitOptionNode is 'LIMIT' used in SELECT.
type LimitOptionNode struct {
	Value IntegerValueNode
}

// GetOption returns a value for printing AST.
func (n LimitOptionNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Value": n.Value,
	}
}

// InsertWithColumnOptionNode stores columns and values used in INSERT.
type InsertWithColumnOptionNode struct {
	Column []interface{}
	Value  []interface{}
}

// GetOption returns a value for printing AST.
func (n InsertWithColumnOptionNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Column": n.Column,
		"Value":  n.Value,
	}
}

// InsertWithDefaultOptionNode is 'DEFAULT VALUES' used in INSERT.
type InsertWithDefaultOptionNode struct{ NilOptionNode }

// PrimaryOptionNode is 'PRIMARY KEY' used in CREATE TABLE.
type PrimaryOptionNode struct{ NilOptionNode }

// NotNullOptionNode is 'NOT NULL' used in CREATE TABLE.
type NotNullOptionNode struct{ NilOptionNode }

// UniqueOptionNode is 'UNIQUE' used in CREATE TABLE and CREATE INDEX.
type UniqueOptionNode struct{ NilOptionNode }

// AutoIncrementOptionNode is 'AUTOINCREMENT' used in CREATE TABLE.
type AutoIncrementOptionNode struct{ NilOptionNode }

// DefaultOptionNode is 'DEFAULT' used in CREATE TABLE.
type DefaultOptionNode struct {
	Value interface{}
}

// GetOption returns a value for printing AST.
func (n DefaultOptionNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Value": n.Value,
	}
}

// ForeignOptionNode is 'REFERENCES' used in CREATE TABLE.
type ForeignOptionNode struct {
	Table  IdentifierNode
	Column IdentifierNode
}

// GetOption returns a value for printing AST.
func (n ForeignOptionNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Table":  n.Table,
		"Column": n.Column,
	}
}

// ---------------------------------------------------------------------------
// Statements
// ---------------------------------------------------------------------------

// SelectStmtNode is SELECT.
type SelectStmtNode struct {
	Column []interface{}
	Table  *IdentifierNode
	Where  *WhereOptionNode
	Group  []interface{}
	Order  []interface{}
	Limit  *LimitOptionNode
	Offset *OffsetOptionNode
}

// GetOption returns a value for printing AST.
func (n SelectStmtNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Column": n.Column,
		"Table":  n.Table,
		"Where":  n.Where,
		"Group":  n.Group,
		"Order":  n.Order,
		"Limit":  n.Limit,
		"Offset": n.Offset,
	}
}

// UpdateStmtNode is UPDATE.
type UpdateStmtNode struct {
	Table      IdentifierNode
	Assignment []interface{}
	Where      *WhereOptionNode
}

// GetOption returns a value for printing AST.
func (n UpdateStmtNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Table":      n.Table,
		"Assignment": n.Assignment,
		"Where":      n.Where,
	}
}

// DeleteStmtNode is DELETE.
type DeleteStmtNode struct {
	Table IdentifierNode
	Where *WhereOptionNode
}

// GetOption returns a value for printing AST.
func (n DeleteStmtNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Table": n.Table,
		"Where": n.Where,
	}
}

// InsertStmtNode is INSERT.
type InsertStmtNode struct {
	Table  IdentifierNode
	Insert interface{}
}

// GetOption returns a value for printing AST.
func (n InsertStmtNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Table":  n.Table,
		"Insert": n.Insert,
	}
}

// CreateTableStmtNode is CREATE TABLE.
type CreateTableStmtNode struct {
	Table  IdentifierNode
	Column []interface{}
}

// GetOption returns a value for printing AST.
func (n CreateTableStmtNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Table":  n.Table,
		"Column": n.Column,
	}
}

// ColumnSchemaNode specifies a column in CREATE TABLE.
type ColumnSchemaNode struct {
	Column     IdentifierNode
	DataType   interface{}
	Constraint []interface{}
}

// GetOption returns a value for printing AST.
func (n ColumnSchemaNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Column":     n.Column,
		"DataYype":   n.DataType,
		"Constraint": n.Constraint,
	}
}

// CreateIndexStmtNode is CREATE INDEX.
type CreateIndexStmtNode struct {
	Index  IdentifierNode
	Table  IdentifierNode
	Column []interface{}
	Unique *UniqueOptionNode
}

// GetOption returns a value for printing AST.
func (n CreateIndexStmtNode) GetOption() map[string]interface{} {
	return map[string]interface{}{
		"Index":  n.Index,
		"Table":  n.Table,
		"Column": n.Column,
		"Unique": n.Unique,
	}
}
