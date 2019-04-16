package ast

import (
	"fmt"

	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

// ---------------------------------------------------------------------------
// Base
// ---------------------------------------------------------------------------

// Node is an interface which should be satisfied by all nodes in AST.
type Node interface {
	HasPosition() bool
	GetPosition() uint32
	SetPosition(uint32)
	GetLength() uint32
	SetLength(uint32)
	GetToken() []byte
	SetToken([]byte)
	GetChildren() []Node
}

// NodeBase is a base struct embedded by structs implementing Node interface.
type NodeBase struct {
	Position uint32 `print:"-"`
	Length   uint32 `print:"-"`
	Token    []byte `print:"-"`
}

// HasPosition returns whether the position is set.
func (n *NodeBase) HasPosition() bool {
	return n.Length > 0
}

// GetPosition returns the offset in bytes where the corresponding token starts.
func (n *NodeBase) GetPosition() uint32 {
	return n.Position
}

// SetPosition sets the offset in bytes where the corresponding token starts.
func (n *NodeBase) SetPosition(position uint32) {
	n.Position = position
}

// GetLength returns the length in bytes of the corresponding token.
func (n *NodeBase) GetLength() uint32 {
	return n.Length
}

// SetLength sets the length in bytes of the corresponding token.
func (n *NodeBase) SetLength(length uint32) {
	n.Length = length
}

// GetToken returns the corresponding token of the node.
func (n *NodeBase) GetToken() []byte {
	return n.Token
}

// SetToken sets the corresponding token of the node.
func (n *NodeBase) SetToken(token []byte) {
	n.Token = token
}

// ---------------------------------------------------------------------------
// Identifiers
// ---------------------------------------------------------------------------

// ExprNode is an interface which should be satisfied all nodes in expressions.
//go-sumtype:decl ExprNode
type ExprNode interface {
	Node
	IsConstant() bool
	GetType() DataType
	SetType(DataType)
	ˉExprNode()
	ˉExprWithAnyNode()
	ˉExprWithDefaultNode()
}

var _ ExprWithAnyNode = (ExprNode)(nil)
var _ ExprWithDefaultNode = (ExprNode)(nil)

// UntaggedExprNodeBase is a base struct embedded by nodes whose types can be
// decided without any context and database schema.
type UntaggedExprNodeBase struct {
	NodeBase
}

func (n *UntaggedExprNodeBase) ˉExprNode()            {}
func (n *UntaggedExprNodeBase) ˉExprWithAnyNode()     {}
func (n *UntaggedExprNodeBase) ˉExprWithDefaultNode() {}

// SetType always panics because it is not reasonable to set data type on nodes
// whose types are already decided.
func (n *UntaggedExprNodeBase) SetType(t DataType) {
	panic("setting type on untagged expression node")
}

// TaggedExprNodeBase is a base struct embedded by nodes whose types depend on
// the context and can only be decided after loading database schemas.
type TaggedExprNodeBase struct {
	NodeBase
	Type DataType `print:"-"`
}

func (n *TaggedExprNodeBase) ˉExprNode()            {}
func (n *TaggedExprNodeBase) ˉExprWithAnyNode()     {}
func (n *TaggedExprNodeBase) ˉExprWithDefaultNode() {}

// GetType gets the data type of the node.
func (n *TaggedExprNodeBase) GetType() DataType {
	return n.Type
}

// SetType sets the data type of the node.
func (n *TaggedExprNodeBase) SetType(t DataType) {
	n.Type = t
}

// IdentifierDescriptor defines the interface of a descriptor. A descriptor
// identifies an object in a SQL statement. This interface is intended to be
// used by IdentifierNode to store the target after the name is resolved.
// FIXME: IdentifierDescriptor is a sum type, but all implementations are in
// different packages.
type IdentifierDescriptor interface {
	GetDescriptor() uint32
}

// IdentifierNode references table, column, or function.
type IdentifierNode struct {
	TaggedExprNodeBase
	Name []byte
	Desc IdentifierDescriptor
}

var _ ExprNode = (*IdentifierNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *IdentifierNode) GetChildren() []Node {
	return nil
}

// IsConstant returns whether a node is a constant.
func (n *IdentifierNode) IsConstant() bool {
	return false
}

// ---------------------------------------------------------------------------
// Values
// ---------------------------------------------------------------------------

// Valuer defines the interface of a constant value.
//go-sumtype:decl Valuer
type Valuer interface {
	ExprNode
	ˉValuer()
}

// BoolValueNode is a boolean constant.
type BoolValueNode struct {
	UntaggedExprNodeBase
	V BoolValue
}

var _ Valuer = (*BoolValueNode)(nil)

func (n *BoolValueNode) ˉValuer() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *BoolValueNode) GetChildren() []Node {
	return nil
}

// IsConstant returns whether a node is a constant.
func (n *BoolValueNode) IsConstant() bool {
	return true
}

// GetType returns the type of 'bool'.
func (n *BoolValueNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// IntegerValueNode is an integer constant.
type IntegerValueNode struct {
	TaggedExprNodeBase
	IsAddress bool
	V         decimal.Decimal
}

var _ Valuer = (*IntegerValueNode)(nil)

func (n *IntegerValueNode) ˉValuer() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *IntegerValueNode) GetChildren() []Node {
	return nil
}

// IsConstant returns whether a node is a constant.
func (n *IntegerValueNode) IsConstant() bool {
	return true
}

// DecimalValueNode is a number constant.
type DecimalValueNode struct {
	TaggedExprNodeBase
	V decimal.Decimal
}

var _ Valuer = (*DecimalValueNode)(nil)

func (n *DecimalValueNode) ˉValuer() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *DecimalValueNode) GetChildren() []Node {
	return nil
}

// IsConstant returns whether a node is a constant.
func (n *DecimalValueNode) IsConstant() bool {
	return true
}

// BytesValueNode is a dynamic or fixed bytes constant.
type BytesValueNode struct {
	TaggedExprNodeBase
	V []byte
}

var _ Valuer = (*BytesValueNode)(nil)

func (n *BytesValueNode) ˉValuer() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *BytesValueNode) GetChildren() []Node {
	return nil
}

// IsConstant returns whether a node is a constant.
func (n *BytesValueNode) IsConstant() bool {
	return true
}

// NullValueNode is NULL.
type NullValueNode struct {
	TaggedExprNodeBase
}

var _ Valuer = (*NullValueNode)(nil)

func (n *NullValueNode) ˉValuer() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *NullValueNode) GetChildren() []Node {
	return nil
}

// IsConstant returns whether a node is a constant.
func (n *NullValueNode) IsConstant() bool {
	return true
}

// ---------------------------------------------------------------------------
// Symbols
// ---------------------------------------------------------------------------

// ExprWithAnyNode is a sum type of ExprNode and AnySymbolNode.
//go-sumtype:decl ExprWithAnyNode
type ExprWithAnyNode interface {
	Node
	ˉExprWithAnyNode()
}

// AnySymbolNode is '*' used in SELECT and function call.
type AnySymbolNode struct {
	NodeBase
}

var _ ExprWithAnyNode = (*AnySymbolNode)(nil)

func (n *AnySymbolNode) ˉExprWithAnyNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *AnySymbolNode) GetChildren() []Node {
	return nil
}

// ExprWithDefaultNode is a sum type of ExprNode and DefaultSymbolNode.
//go-sumtype:decl ExprWithDefaultNode
type ExprWithDefaultNode interface {
	Node
	ˉExprWithDefaultNode()
}

// DefaultSymbolNode represents the default value used in INSERT and UPDATE.
type DefaultSymbolNode struct {
	NodeBase
}

var _ ExprWithDefaultNode = (*DefaultSymbolNode)(nil)

func (n *DefaultSymbolNode) ˉExprWithDefaultNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *DefaultSymbolNode) GetChildren() []Node {
	return nil
}

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// TypeNode is an interface which should be satisfied nodes representing types.
//go-sumtype:decl TypeNode
type TypeNode interface {
	Node
	GetType() (DataType, errors.ErrorCode, string)
	ˉTypeNode()
}

// IntTypeNode represents solidity int{X} and uint{X} types.
type IntTypeNode struct {
	NodeBase
	Unsigned bool
	Size     uint32
}

var _ TypeNode = (*IntTypeNode)(nil)

func (n *IntTypeNode) ˉTypeNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *IntTypeNode) GetChildren() []Node {
	return nil
}

// GetType returns the type represented by the node.
func (n *IntTypeNode) GetType() (DataType, errors.ErrorCode, string) {
	isMultiple := n.Size%8 == 0
	isNotZero := n.Size != 0
	isInRange := n.Size <= 256
	if !isMultiple || !isNotZero || !isInRange {
		name := "int"
		code := errors.ErrorCodeInvalidIntSize
		if n.Unsigned {
			name = "uint"
			code = errors.ErrorCodeInvalidUintSize
		}
		if !isMultiple {
			return DataTypeBad, code, fmt.Sprintf(
				"%s size %d is not a multiple of 8", name, n.Size)
		}
		if !isNotZero {
			return DataTypeBad, code, fmt.Sprintf(
				"%s size cannot be zero", name)
		}
		if !isInRange {
			return DataTypeBad, code, fmt.Sprintf(
				"%s size %d cannot be larger than 256", name, n.Size)
		}
		panic("unreachable")
	}
	var major DataTypeMajor
	var minor DataTypeMinor
	if n.Unsigned {
		major = DataTypeMajorUint
	} else {
		major = DataTypeMajorInt
	}
	minor = DataTypeMinor(n.Size/8 - 1)
	return ComposeDataType(major, minor), errors.ErrorCodeNil, ""
}

// FixedTypeNode represents solidity fixed{M}x{N} and ufixed{M}x{N} types.
type FixedTypeNode struct {
	NodeBase
	Unsigned         bool
	Size             uint32
	FractionalDigits uint32
}

var _ TypeNode = (*FixedTypeNode)(nil)

func (n *FixedTypeNode) ˉTypeNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *FixedTypeNode) GetChildren() []Node {
	return nil
}

// GetType returns the type represented by the node.
func (n *FixedTypeNode) GetType() (DataType, errors.ErrorCode, string) {
	sizeIsMultiple := n.Size%8 == 0
	sizeIsNotZero := n.Size != 0
	sizeIsInRange := n.Size <= 256
	fractionalDigitsInRange := n.FractionalDigits <= 80
	if !sizeIsMultiple || !sizeIsNotZero || !sizeIsInRange ||
		!fractionalDigitsInRange {
		name := "fixed"
		code := errors.ErrorCodeInvalidFixedSize
		if n.Unsigned {
			name = "ufixed"
			code = errors.ErrorCodeInvalidUfixedSize
		}
		if !sizeIsMultiple {
			return DataTypeBad, code, fmt.Sprintf(
				"%s size %d is not a multiple of 8", name, n.Size)
		}
		if !sizeIsNotZero {
			return DataTypeBad, code, fmt.Sprintf(
				"%s size cannot be zero", name)
		}
		if !sizeIsInRange {
			return DataTypeBad, code, fmt.Sprintf(
				"%s size %d cannot be larger than 256", name, n.Size)
		}
		code = errors.ErrorCodeInvalidFixedFractionalDigits
		if n.Unsigned {
			code = errors.ErrorCodeInvalidUfixedFractionalDigits
		}
		if !fractionalDigitsInRange {
			return DataTypeBad, code, fmt.Sprintf(
				"%s fractional digits %d cannot be larger than 80",
				name, n.FractionalDigits)
		}
		panic("unreachable")
	}
	var major DataTypeMajor
	var minor DataTypeMinor
	if n.Unsigned {
		major = DataTypeMajorUfixed
	} else {
		major = DataTypeMajorFixed
	}
	major += DataTypeMajor(n.Size/8 - 1)
	minor = DataTypeMinor(n.FractionalDigits)
	return ComposeDataType(major, minor), errors.ErrorCodeNil, ""
}

// DynamicBytesTypeNode represents solidity bytes type.
type DynamicBytesTypeNode struct {
	NodeBase
}

var _ TypeNode = (*DynamicBytesTypeNode)(nil)

func (n *DynamicBytesTypeNode) ˉTypeNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *DynamicBytesTypeNode) GetChildren() []Node {
	return nil
}

// GetType returns the type represented by the node.
func (n *DynamicBytesTypeNode) GetType() (DataType, errors.ErrorCode, string) {
	return ComposeDataType(DataTypeMajorDynamicBytes, DataTypeMinorDontCare),
		errors.ErrorCodeNil, ""
}

// FixedBytesTypeNode represents solidity bytes{X} type.
type FixedBytesTypeNode struct {
	NodeBase
	Size uint32
}

var _ TypeNode = (*FixedBytesTypeNode)(nil)

func (n *FixedBytesTypeNode) ˉTypeNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *FixedBytesTypeNode) GetChildren() []Node {
	return nil
}

// GetType returns the type represented by the node.
func (n *FixedBytesTypeNode) GetType() (DataType, errors.ErrorCode, string) {
	isNotZero := n.Size != 0
	isInRange := n.Size <= 32
	if !isNotZero || !isInRange {
		code := errors.ErrorCodeInvalidBytesSize
		if !isNotZero {
			return DataTypeBad, code, "bytes size cannot be zero"
		}
		if !isInRange {
			return DataTypeBad, code, fmt.Sprintf(
				"bytes size %d cannot be larger than 32", n.Size)
		}
		panic("unreachable")
	}
	major := DataTypeMajorFixedBytes
	minor := DataTypeMinor(n.Size - 1)
	return ComposeDataType(major, minor), errors.ErrorCodeNil, ""
}

// AddressTypeNode represents solidity address type.
type AddressTypeNode struct {
	NodeBase
}

var _ TypeNode = (*AddressTypeNode)(nil)

func (n *AddressTypeNode) ˉTypeNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *AddressTypeNode) GetChildren() []Node {
	return nil
}

// GetType returns the type represented by the node.
func (n *AddressTypeNode) GetType() (DataType, errors.ErrorCode, string) {
	return ComposeDataType(DataTypeMajorAddress, DataTypeMinorDontCare),
		errors.ErrorCodeNil, ""
}

// BoolTypeNode represents solidity bool type.
type BoolTypeNode struct {
	NodeBase
}

var _ TypeNode = (*BoolTypeNode)(nil)

func (n *BoolTypeNode) ˉTypeNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *BoolTypeNode) GetChildren() []Node {
	return nil
}

// GetType returns the type represented by the node.
func (n *BoolTypeNode) GetType() (DataType, errors.ErrorCode, string) {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare),
		errors.ErrorCodeNil, ""
}

// ---------------------------------------------------------------------------
// Operators
// ---------------------------------------------------------------------------

// UnaryOperator defines the interface of a unary operator.
//go-sumtype:decl UnaryOperator
type UnaryOperator interface {
	ExprNode
	GetTarget() ExprNode
	SetTarget(ExprNode)
	ˉUnaryOperator()
}

// BinaryOperator defines the interface of a binary operator.
//go-sumtype:decl BinaryOperator
type BinaryOperator interface {
	ExprNode
	GetObject() ExprNode
	GetSubject() ExprNode
	SetObject(ExprNode)
	SetSubject(ExprNode)
	ˉBinaryOperator()
}

// UnaryOperatorNode is a base struct of unary operators.
type UnaryOperatorNode struct {
	Target ExprNode
}

func (n *UnaryOperatorNode) ˉUnaryOperator() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *UnaryOperatorNode) GetChildren() []Node {
	return []Node{n.Target}
}

// IsConstant returns whether a node is a constant.
func (n *UnaryOperatorNode) IsConstant() bool {
	return n.Target.IsConstant()
}

// GetTarget gets the target of the operation.
func (n *UnaryOperatorNode) GetTarget() ExprNode {
	return n.Target
}

// SetTarget sets the target of the operation.
func (n *UnaryOperatorNode) SetTarget(t ExprNode) {
	n.Target = t
}

// BinaryOperatorNode is a base struct of binary operators.
type BinaryOperatorNode struct {
	Object  ExprNode
	Subject ExprNode
}

func (n *BinaryOperatorNode) ˉBinaryOperator() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *BinaryOperatorNode) GetChildren() []Node {
	return []Node{n.Object, n.Subject}
}

// IsConstant returns whether a node is a constant.
func (n *BinaryOperatorNode) IsConstant() bool {
	return n.Object.IsConstant() && n.Subject.IsConstant()
}

// GetObject gets the node on which the operation is applied.
func (n *BinaryOperatorNode) GetObject() ExprNode {
	return n.Object
}

// GetSubject gets the node whose value is applied on the object.
func (n *BinaryOperatorNode) GetSubject() ExprNode {
	return n.Subject
}

// SetObject sets the object of the operation.
func (n *BinaryOperatorNode) SetObject(o ExprNode) {
	n.Object = o
}

// SetSubject sets the subject of the operation.
func (n *BinaryOperatorNode) SetSubject(s ExprNode) {
	n.Subject = s
}

// PosOperatorNode is '+'.
type PosOperatorNode struct {
	TaggedExprNodeBase
	UnaryOperatorNode
}

var _ UnaryOperator = (*PosOperatorNode)(nil)

// NegOperatorNode is '-'.
type NegOperatorNode struct {
	TaggedExprNodeBase
	UnaryOperatorNode
}

var _ UnaryOperator = (*NegOperatorNode)(nil)

// NotOperatorNode is 'NOT'.
type NotOperatorNode struct {
	UntaggedExprNodeBase
	UnaryOperatorNode
}

var _ UnaryOperator = (*NotOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *NotOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// ParenOperatorNode is a pair of '(' and ')', representing a parenthesized
// expression.
type ParenOperatorNode struct {
	TaggedExprNodeBase
	UnaryOperatorNode
}

var _ UnaryOperator = (*ParenOperatorNode)(nil)

// AndOperatorNode is 'AND'.
type AndOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*AndOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *AndOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// OrOperatorNode is 'OR'.
type OrOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*OrOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *OrOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// GreaterOrEqualOperatorNode is '>='.
type GreaterOrEqualOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*GreaterOrEqualOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *GreaterOrEqualOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// LessOrEqualOperatorNode is '<='.
type LessOrEqualOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*LessOrEqualOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *LessOrEqualOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// NotEqualOperatorNode is '<>'.
type NotEqualOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*NotEqualOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *NotEqualOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// EqualOperatorNode is '=' used in expressions.
type EqualOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*EqualOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *EqualOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// GreaterOperatorNode is '>'.
type GreaterOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*GreaterOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *GreaterOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// LessOperatorNode is '<'.
type LessOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*LessOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *LessOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// ConcatOperatorNode is '||'.
type ConcatOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*ConcatOperatorNode)(nil)

// GetType returns the type of 'bytes'.
func (n *ConcatOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorDynamicBytes, DataTypeMinorDontCare)
}

// AddOperatorNode is '+'.
type AddOperatorNode struct {
	TaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*AddOperatorNode)(nil)

// SubOperatorNode is '-'.
type SubOperatorNode struct {
	TaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*SubOperatorNode)(nil)

// MulOperatorNode is '*'.
type MulOperatorNode struct {
	TaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*MulOperatorNode)(nil)

// DivOperatorNode is '/'.
type DivOperatorNode struct {
	TaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*DivOperatorNode)(nil)

// ModOperatorNode is '%'.
type ModOperatorNode struct {
	TaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*ModOperatorNode)(nil)

// IsOperatorNode is 'IS NULL'.
type IsOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
}

var _ BinaryOperator = (*IsOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *IsOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// LikeOperatorNode is 'LIKE'.
type LikeOperatorNode struct {
	UntaggedExprNodeBase
	BinaryOperatorNode
	Escape ExprNode
}

var _ BinaryOperator = (*LikeOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *LikeOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// GetChildren returns a list of child nodes used for traversing.
func (n *LikeOperatorNode) GetChildren() []Node {
	size := 2
	if n.Escape != nil {
		size++
	}

	idx := 0
	nodes := make([]Node, size)
	nodes[idx] = n.Object
	idx++
	nodes[idx] = n.Subject
	idx++
	if n.Escape != nil {
		nodes[idx] = n.Escape
		idx++
	}
	return nodes
}

// IsConstant returns whether a node is a constant.
func (n *LikeOperatorNode) IsConstant() bool {
	if !n.Object.IsConstant() {
		return false
	}
	if !n.Subject.IsConstant() {
		return false
	}
	if n.Escape != nil && !n.Escape.IsConstant() {
		return false
	}
	return true
}

// ---------------------------------------------------------------------------
// Cast
// ---------------------------------------------------------------------------

// CastOperatorNode is 'CAST(expr AS type)'.
type CastOperatorNode struct {
	UntaggedExprNodeBase
	SourceExpr ExprNode
	TargetType TypeNode
}

var _ ExprNode = (*CastOperatorNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *CastOperatorNode) GetChildren() []Node {
	return []Node{n.SourceExpr, n.TargetType}
}

// IsConstant returns whether a node is a constant.
func (n *CastOperatorNode) IsConstant() bool {
	return n.SourceExpr.IsConstant()
}

// GetType returns the type of CAST expression, which is always the target type.
func (n *CastOperatorNode) GetType() DataType {
	if dt, code, _ := n.TargetType.GetType(); code == errors.ErrorCodeNil {
		return dt
	}
	return DataTypeBad
}

// ---------------------------------------------------------------------------
// Assignment
// ---------------------------------------------------------------------------

// AssignOperatorNode is '=' used in UPDATE to set values.
type AssignOperatorNode struct {
	NodeBase
	Column *IdentifierNode
	Expr   ExprWithDefaultNode
}

var _ Node = (*AssignOperatorNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *AssignOperatorNode) GetChildren() []Node {
	return []Node{n.Column, n.Expr}
}

// ---------------------------------------------------------------------------
// In
// ---------------------------------------------------------------------------

// InOperatorNode is 'IN'.
type InOperatorNode struct {
	UntaggedExprNodeBase
	Left  ExprNode
	Right []ExprNode
}

var _ ExprNode = (*InOperatorNode)(nil)

// GetType returns the type of 'bool'.
func (n *InOperatorNode) GetType() DataType {
	return ComposeDataType(DataTypeMajorBool, DataTypeMinorDontCare)
}

// GetChildren returns a list of child nodes used for traversing.
func (n *InOperatorNode) GetChildren() []Node {
	nodes := make([]Node, 1+len(n.Right))
	nodes[0] = n.Left
	for i := 0; i < len(n.Right); i++ {
		nodes[i+1] = n.Right[i]
	}
	return nodes
}

// IsConstant returns whether a node is a constant.
func (n *InOperatorNode) IsConstant() bool {
	if !n.Left.IsConstant() {
		return false
	}
	for i := 0; i < len(n.Right); i++ {
		if !n.Right[i].IsConstant() {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Function
// ---------------------------------------------------------------------------

// FunctionOperatorNode is a function call.
type FunctionOperatorNode struct {
	TaggedExprNodeBase
	Name *IdentifierNode
	Args []ExprWithAnyNode
}

var _ ExprNode = (*FunctionOperatorNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *FunctionOperatorNode) GetChildren() []Node {
	nodes := make([]Node, 1+len(n.Args))
	nodes[0] = n.Name
	for i := 0; i < len(n.Args); i++ {
		nodes[i+1] = n.Args[i]
	}
	return nodes
}

// IsConstant returns whether a node is a constant.
func (n *FunctionOperatorNode) IsConstant() bool {
	return false
}

// ---------------------------------------------------------------------------
// Options
// ---------------------------------------------------------------------------

// WhereOptionNode is 'WHERE' used in SELECT, UPDATE, DELETE.
type WhereOptionNode struct {
	NodeBase
	Condition ExprNode
}

var _ Node = (*WhereOptionNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *WhereOptionNode) GetChildren() []Node {
	return []Node{n.Condition}
}

// OrderOptionNode is an expression in 'ORDER BY' used in SELECT.
type OrderOptionNode struct {
	NodeBase
	Expr       ExprNode
	Desc       bool
	NullsFirst bool
}

var _ Node = (*OrderOptionNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *OrderOptionNode) GetChildren() []Node {
	return []Node{n.Expr}
}

// GroupOptionNode is 'GROUP BY' used in SELECT.
type GroupOptionNode struct {
	NodeBase
	Expr ExprNode
}

var _ Node = (*GroupOptionNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *GroupOptionNode) GetChildren() []Node {
	return []Node{n.Expr}
}

// OffsetOptionNode is 'OFFSET' used in SELECT.
type OffsetOptionNode struct {
	NodeBase
	Value *IntegerValueNode
}

var _ Node = (*OffsetOptionNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *OffsetOptionNode) GetChildren() []Node {
	return []Node{n.Value}
}

// LimitOptionNode is 'LIMIT' used in SELECT.
type LimitOptionNode struct {
	NodeBase
	Value *IntegerValueNode
}

var _ Node = (*LimitOptionNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *LimitOptionNode) GetChildren() []Node {
	return []Node{n.Value}
}

// InsertOptionNode is a sum type of all INSERT options.
//go-sumtype:decl InsertOptionNode
type InsertOptionNode interface {
	Node
	ˉInsertOptionNode()
}

// InsertWithColumnOptionNode stores columns and values used in INSERT.
type InsertWithColumnOptionNode struct {
	NodeBase
	Column []*IdentifierNode
	Value  [][]ExprWithDefaultNode
}

var _ InsertOptionNode = (*InsertWithColumnOptionNode)(nil)

func (n *InsertWithColumnOptionNode) ˉInsertOptionNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *InsertWithColumnOptionNode) GetChildren() []Node {
	size := len(n.Column)
	for i := 0; i < len(n.Value); i++ {
		size += len(n.Value[i])
	}

	nodes := make([]Node, size)
	idx := 0
	for i := 0; i < len(n.Column); i, idx = i+1, idx+1 {
		nodes[idx] = n.Column[i]
	}
	for i := 0; i < len(n.Value); i++ {
		for j := 0; j < len(n.Value[i]); j, idx = j+1, idx+1 {
			nodes[idx] = n.Value[i][j]
		}
	}
	return nodes
}

// InsertWithDefaultOptionNode is 'DEFAULT VALUES' used in INSERT.
type InsertWithDefaultOptionNode struct {
	NodeBase
}

var _ InsertOptionNode = (*InsertWithDefaultOptionNode)(nil)

func (n *InsertWithDefaultOptionNode) ˉInsertOptionNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *InsertWithDefaultOptionNode) GetChildren() []Node {
	return nil
}

// ColumnConstraintNode is a sum type of all column constraints.
//go-sumtype:decl ColumnConstraintNode
type ColumnConstraintNode interface {
	Node
	ˉColumnConstraintNode()
}

// PrimaryOptionNode is 'PRIMARY KEY' used in CREATE TABLE.
type PrimaryOptionNode struct {
	NodeBase
}

var _ ColumnConstraintNode = (*PrimaryOptionNode)(nil)

func (n *PrimaryOptionNode) ˉColumnConstraintNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *PrimaryOptionNode) GetChildren() []Node {
	return nil
}

// NotNullOptionNode is 'NOT NULL' used in CREATE TABLE.
type NotNullOptionNode struct {
	NodeBase
}

var _ ColumnConstraintNode = (*NotNullOptionNode)(nil)

func (n *NotNullOptionNode) ˉColumnConstraintNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *NotNullOptionNode) GetChildren() []Node {
	return nil
}

// UniqueOptionNode is 'UNIQUE' used in CREATE TABLE and CREATE INDEX.
type UniqueOptionNode struct {
	NodeBase
}

var _ ColumnConstraintNode = (*UniqueOptionNode)(nil)

func (n *UniqueOptionNode) ˉColumnConstraintNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *UniqueOptionNode) GetChildren() []Node {
	return nil
}

// AutoIncrementOptionNode is 'AUTOINCREMENT' used in CREATE TABLE.
type AutoIncrementOptionNode struct {
	NodeBase
}

var _ ColumnConstraintNode = (*AutoIncrementOptionNode)(nil)

func (n *AutoIncrementOptionNode) ˉColumnConstraintNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *AutoIncrementOptionNode) GetChildren() []Node {
	return nil
}

// DefaultOptionNode is 'DEFAULT' used in CREATE TABLE.
type DefaultOptionNode struct {
	NodeBase
	Value ExprNode
}

var _ ColumnConstraintNode = (*DefaultOptionNode)(nil)

func (n *DefaultOptionNode) ˉColumnConstraintNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *DefaultOptionNode) GetChildren() []Node {
	return []Node{n.Value}
}

// ForeignOptionNode is 'REFERENCES' used in CREATE TABLE.
type ForeignOptionNode struct {
	NodeBase
	Table  *IdentifierNode
	Column *IdentifierNode
}

var _ ColumnConstraintNode = (*ForeignOptionNode)(nil)

func (n *ForeignOptionNode) ˉColumnConstraintNode() {}

// GetChildren returns a list of child nodes used for traversing.
func (n *ForeignOptionNode) GetChildren() []Node {
	return []Node{n.Table, n.Column}
}

// ---------------------------------------------------------------------------
// Statements
// ---------------------------------------------------------------------------

// StmtNode defines the interface of a statement.
//go-sumtype:decl StmtNode
type StmtNode interface {
	Node
	GetVerb() []byte
	SetVerb([]byte)
	ˉStmtNode()
}

// StmtNodeBase is a base struct embedded by statement nodes.
type StmtNodeBase struct {
	Verb []byte `print:"-"`
}

func (n *StmtNodeBase) ˉStmtNode() {}

// GetVerb returns the verb used to identify the statement.
func (n *StmtNodeBase) GetVerb() []byte {
	return n.Verb
}

// SetVerb sets the verb used to identify the statement.
func (n *StmtNodeBase) SetVerb(verb []byte) {
	n.Verb = verb
}

// SelectStmtNode is SELECT.
type SelectStmtNode struct {
	NodeBase
	StmtNodeBase
	Column []ExprWithAnyNode
	Table  *IdentifierNode
	Where  *WhereOptionNode
	Group  []*GroupOptionNode
	Order  []*OrderOptionNode
	Limit  *LimitOptionNode
	Offset *OffsetOptionNode
}

var _ StmtNode = (*SelectStmtNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *SelectStmtNode) GetChildren() []Node {
	size := len(n.Column) + len(n.Group) + len(n.Order)
	if n.Table != nil {
		size++
	}
	if n.Where != nil {
		size++
	}
	if n.Limit != nil {
		size++
	}
	if n.Offset != nil {
		size++
	}

	nodes := make([]Node, size)
	idx := 0
	for i := 0; i < len(n.Column); i, idx = i+1, idx+1 {
		nodes[idx] = n.Column[i]
	}
	if n.Table != nil {
		nodes[idx] = n.Table
		idx++
	}
	if n.Where != nil {
		nodes[idx] = n.Where
		idx++
	}
	for i := 0; i < len(n.Group); i, idx = i+1, idx+1 {
		nodes[idx] = n.Group[i]
	}
	for i := 0; i < len(n.Order); i, idx = i+1, idx+1 {
		nodes[idx] = n.Order[i]
	}
	if n.Limit != nil {
		nodes[idx] = n.Limit
		idx++
	}
	if n.Offset != nil {
		nodes[idx] = n.Offset
		idx++
	}
	return nodes
}

// UpdateStmtNode is UPDATE.
type UpdateStmtNode struct {
	NodeBase
	StmtNodeBase
	Table      *IdentifierNode
	Assignment []*AssignOperatorNode
	Where      *WhereOptionNode
}

var _ StmtNode = (*UpdateStmtNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *UpdateStmtNode) GetChildren() []Node {
	size := 1 + len(n.Assignment)
	if n.Where != nil {
		size++
	}

	nodes := make([]Node, size)
	idx := 0
	nodes[idx] = n.Table
	idx++
	for i := 0; i < len(n.Assignment); i, idx = i+1, idx+1 {
		nodes[idx] = n.Assignment[i]
	}
	if n.Where != nil {
		nodes[idx] = n.Where
		idx++
	}
	return nodes
}

// DeleteStmtNode is DELETE.
type DeleteStmtNode struct {
	NodeBase
	StmtNodeBase
	Table *IdentifierNode
	Where *WhereOptionNode
}

var _ StmtNode = (*DeleteStmtNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *DeleteStmtNode) GetChildren() []Node {
	if n.Where == nil {
		return []Node{n.Table}
	}
	return []Node{n.Table, n.Where}
}

// InsertStmtNode is INSERT.
type InsertStmtNode struct {
	NodeBase
	StmtNodeBase
	Table  *IdentifierNode
	Insert InsertOptionNode
}

var _ StmtNode = (*InsertStmtNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *InsertStmtNode) GetChildren() []Node {
	return []Node{n.Table, n.Insert}
}

// CreateTableStmtNode is CREATE TABLE.
type CreateTableStmtNode struct {
	NodeBase
	StmtNodeBase
	Table  *IdentifierNode
	Column []*ColumnSchemaNode
}

var _ StmtNode = (*CreateTableStmtNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *CreateTableStmtNode) GetChildren() []Node {
	nodes := make([]Node, 1+len(n.Column))
	nodes[0] = n.Table
	for i := 0; i < len(n.Column); i++ {
		nodes[i+1] = n.Column[i]
	}
	return nodes
}

// ColumnSchemaNode specifies a column in CREATE TABLE.
type ColumnSchemaNode struct {
	NodeBase
	Column     *IdentifierNode
	DataType   TypeNode
	Constraint []ColumnConstraintNode
}

var _ Node = (*ColumnSchemaNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *ColumnSchemaNode) GetChildren() []Node {
	nodes := make([]Node, 2+len(n.Constraint))
	nodes[0] = n.Column
	nodes[1] = n.DataType
	for i := 0; i < len(n.Constraint); i++ {
		nodes[i+2] = n.Constraint[i]
	}
	return nodes
}

// CreateIndexStmtNode is CREATE INDEX.
type CreateIndexStmtNode struct {
	NodeBase
	StmtNodeBase
	Index  *IdentifierNode
	Table  *IdentifierNode
	Column []*IdentifierNode
	Unique *UniqueOptionNode
}

var _ StmtNode = (*CreateIndexStmtNode)(nil)

// GetChildren returns a list of child nodes used for traversing.
func (n *CreateIndexStmtNode) GetChildren() []Node {
	size := 2 + len(n.Column)
	if n.Unique != nil {
		size++
	}

	nodes := make([]Node, size)
	idx := 0
	nodes[idx] = n.Index
	idx++
	nodes[idx] = n.Table
	idx++
	for i := 0; i < len(n.Column); i, idx = i+1, idx+1 {
		nodes[idx] = n.Column[i]
	}
	if n.Unique != nil {
		nodes[idx] = n.Unique
		idx++
	}
	return nodes
}
