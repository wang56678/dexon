package schema

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/common"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	se "github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/rlp"
)

// Error defines for encode and decode.
var (
	ErrEncodeUnexpectedDataType    = errors.New("encode unexpected data type")
	ErrEncodeUnexpectedDefaultType = errors.New("encode unexpected default type")
	ErrDecodeUnexpectedDataType    = errors.New("decode unexpected data type")
	ErrDecodeUnexpectedDefaultType = errors.New("decode unexpected default type")
)

// ColumnAttr defines bit flags for describing column attribute.
type ColumnAttr uint16

const (
	// ColumnAttrPrimaryKey is a no-op. Primary key constraints are converted
	// to a unique index during contract creation.
	ColumnAttrPrimaryKey ColumnAttr = 1 << iota
	// ColumnAttrNotNull is a no-op. We have not supported NULL values so all
	// columns are implicitly non-null.
	ColumnAttrNotNull
	// ColumnAttrUnique is a no-op. Unique constraints are converted to unique
	// indices during contract creation.
	ColumnAttrUnique
	// ColumnAttrHasDefault indicates whether a column has a default value. The
	// default value does not affect the starting value of AUTOINCREMENT.
	ColumnAttrHasDefault
	// ColumnAttrHasForeignKey indicates whether a column references a column
	// on a different table.
	ColumnAttrHasForeignKey
	// ColumnAttrHasSequence indicates whether a column is declared with
	// AUTOINCREMENT. It is only valid on integer fields.
	ColumnAttrHasSequence
)

// GetDeclaredFlags returns flags which can be mapped to the source code tokens.
func (a ColumnAttr) GetDeclaredFlags() ColumnAttr {
	mask := ColumnAttrPrimaryKey |
		ColumnAttrNotNull |
		ColumnAttrUnique |
		ColumnAttrHasDefault |
		ColumnAttrHasForeignKey |
		ColumnAttrHasSequence
	return a & mask

}

// GetDerivedFlags returns flags which are not declared in the source code but
// can be derived from it.
func (a ColumnAttr) GetDerivedFlags() ColumnAttr {
	mask := ColumnAttr(0)
	return a & mask
}

// FunctionRef defines the type for number of builtin function.
type FunctionRef uint16

// MaxFunctionRef is the maximum value of FunctionRef.
const MaxFunctionRef = math.MaxUint16

// TableRef defines the type for table index in Schema.
type TableRef uint8

// MaxTableRef is the maximum value of TableRef.
const MaxTableRef = math.MaxUint8

// ColumnRef defines the type for column index in Table.Columns.
type ColumnRef uint8

// MaxColumnRef is the maximum value of ColumnRef.
const MaxColumnRef = math.MaxUint8

// IndexRef defines the type for array index of Column.Indices.
type IndexRef uint8

// MaxIndexRef is the maximum value of IndexRef.
const MaxIndexRef = math.MaxUint8

// SequenceRef defines the type for sequence index in Table.
type SequenceRef uint8

// MaxSequenceRef is the maximum value of SequenceRef.
const MaxSequenceRef = math.MaxUint8

// SelectColumnRef defines the type for column index in SelectStmtNode.Column.
type SelectColumnRef uint16

// MaxSelectColumnRef is the maximum value of SelectColumnRef.
const MaxSelectColumnRef = math.MaxUint16

// IndexAttr defines bit flags for describing index attribute.
type IndexAttr uint16

const (
	// IndexAttrUnique indicates whether an index is unique.
	IndexAttrUnique IndexAttr = 1 << iota
	// IndexAttrReferenced indicates whether an index is referenced by columns
	// with foreign key constraints. This attribute cannot be specified by
	// users. It is computed automatically during contract creation.
	IndexAttrReferenced
)

// GetDeclaredFlags returns flags which can be mapped to the source code tokens.
func (a IndexAttr) GetDeclaredFlags() IndexAttr {
	mask := IndexAttrUnique
	return a & mask
}

// GetDerivedFlags returns flags which are not declared in the source code but
// can be derived from it.
func (a IndexAttr) GetDerivedFlags() IndexAttr {
	mask := IndexAttrReferenced
	return a & mask
}

// Schema defines sqlvm schema struct.
type Schema []Table

// SetupColumnOffset set all tables' column offset.
func (s Schema) SetupColumnOffset() {
	for i := range s {
		s[i].SetupColumnOffset()
	}
}

func (s Schema) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf(
		"-- DEXON SQLVM database schema dump (%d tables)\n", len(s)))
	for _, t := range s {
		b.WriteString("CREATE TABLE ")
		b.Write(ast.QuoteIdentifierOptional(t.Name))
		b.WriteString(" (\n")
		for ci, c := range t.Columns {
			b.WriteString("    ")
			b.Write(ast.QuoteIdentifierOptional(c.Name))
			b.WriteByte(' ')
			b.WriteString(c.Type.String())
			comments := []string{
				fmt.Sprintf("slot %d", c.SlotOffset),
				fmt.Sprintf("byte %d", c.ByteOffset),
			}
			if (c.Attr & ColumnAttrPrimaryKey) != 0 {
				b.WriteString(" PRIMARY KEY")
			}
			if (c.Attr & ColumnAttrNotNull) != 0 {
				b.WriteString(" NOT NULL")
			}
			if (c.Attr & ColumnAttrUnique) != 0 {
				b.WriteString(" UNIQUE")
			}
			if (c.Attr & ColumnAttrHasDefault) != 0 {
				b.WriteString(" DEFAULT ")
				switch v := c.Default.(type) {
				case nil:
					b.WriteString("NULL")
				case bool:
					if v {
						b.WriteString("TRUE")
					} else {
						b.WriteString("FALSE")
					}
				case []byte:
					major, _ := ast.DecomposeDataType(c.Type)
					if major == ast.DataTypeMajorAddress {
						b.WriteString(common.BytesToAddress(v).String())
						break
					}
					b.Write(ast.QuoteString(v))
				case decimal.Decimal:
					b.WriteString(v.String())
				default:
					b.WriteString("<?>")
				}
			}
			if (c.Attr & ColumnAttrHasForeignKey) != 0 {
				for _, fk := range c.ForeignKeys {
					b.WriteString(" REFERENCES ")
					b.Write(ast.QuoteIdentifierOptional(
						s[fk.Table].Name))
					b.WriteByte('(')
					b.Write(ast.QuoteIdentifierOptional(
						s[fk.Table].Columns[fk.Column].Name))
					b.WriteByte(')')
				}
			}
			if (c.Attr & ColumnAttrHasSequence) != 0 {
				b.WriteString(" AUTOINCREMENT")
				comments = append(comments,
					fmt.Sprintf("sequence %d", c.Sequence))
			}
			if ci < len(t.Columns)-1 {
				b.WriteByte(',')
			}
			b.WriteString(" -- ")
			b.WriteString(strings.Join(comments, ", "))
			b.WriteByte('\n')
		}
		b.WriteString(");\n")
		for _, i := range t.Indices {
			comments := []string{}
			b.WriteString("CREATE")
			if (i.Attr & IndexAttrUnique) != 0 {
				b.WriteString(" UNIQUE")
			}
			if (i.Attr & IndexAttrReferenced) != 0 {
				comments = append(comments, "referenced")
			}
			b.WriteString(" INDEX ")
			b.Write(ast.QuoteIdentifierOptional(i.Name))
			b.WriteString(" ON ")
			b.Write(ast.QuoteIdentifierOptional(t.Name))
			b.WriteByte('(')
			for ci, c := range i.Columns {
				b.Write(ast.QuoteIdentifierOptional(t.Columns[c].Name))
				if ci < len(i.Columns)-1 {
					b.WriteString(", ")
				}
			}
			b.WriteString(");")
			if len(comments) > 0 {
				b.WriteString(" -- ")
				b.WriteString(strings.Join(comments, ", "))
			}
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// Table defiens sqlvm table struct.
type Table struct {
	Name    []byte
	Columns []Column
	Indices []Index
}

// GetFieldType return fields' data type.
func (t *Table) GetFieldType(fields []uint8) ([]ast.DataType, error) {
	types := make([]ast.DataType, len(fields))
	columns := t.Columns
	for i, f := range fields {
		if int(f) >= len(columns) {
			return nil, se.ErrorCodeIndexOutOfRange
		}
		types[i] = columns[f].Type
	}
	return types, nil
}

// SetupColumnOffset set columns' slot and byte offset.
func (t *Table) SetupColumnOffset() {
	slotOffset := uint8(0)
	byteOffset := uint8(0)
	for i, col := range t.Columns {
		size := col.Type.Size()
		if size+byteOffset > common.HashLength {
			slotOffset++
			byteOffset = 0
		}
		t.Columns[i].SlotOffset = slotOffset
		t.Columns[i].ByteOffset = byteOffset
		byteOffset += size
	}
}

// Index defines sqlvm index struct.
type Index struct {
	Name    []byte
	Attr    IndexAttr
	Columns []ColumnRef // Columns must be sorted in ascending order.
}

type column struct {
	Name        []byte
	Type        ast.DataType
	Attr        ColumnAttr
	ForeignKeys []ColumnDescriptor
	Sequence    SequenceRef
	SlotOffset  uint8 `rlp:"-"`
	ByteOffset  uint8 `rlp:"-"`
	// Rest is a special field reserved for use in EncodeRLP. The value stored
	// in it will be overwritten every time EncodeRLP is called.
	Rest interface{}
}

// MaxForeignKeys is the maximum number of foreign key constraints which can be
// defined on a column.
const MaxForeignKeys = math.MaxUint8

// Column defines sqlvm index struct.
type Column struct {
	column
	Default interface{} // decimal.Decimal, bool, []byte
}

// NewColumn return a Column instance.
func NewColumn(Name []byte, Type ast.DataType, Attr ColumnAttr,
	ForeignKeys []ColumnDescriptor, Sequence SequenceRef) Column {
	c := column{
		Name:        Name,
		Type:        Type,
		Attr:        Attr,
		ForeignKeys: ForeignKeys,
		Sequence:    Sequence,
	}

	return Column{
		column: c,
	}
}

var _ rlp.Decoder = (*Column)(nil)
var _ rlp.Encoder = Column{}

// EncodeRLP encodes column with rlp encode.
func (c Column) EncodeRLP(w io.Writer) error {
	if c.Default != nil {
		switch d := c.Default.(type) {
		case bool:
			v := byte(0)
			if d {
				v = byte(1)
			}
			c.Rest = []byte{v}
		case []byte:
			c.Rest = d
		case decimal.Decimal:
			var ok bool
			c.Rest, ok = ast.DecimalEncode(c.Type, d)
			if !ok {
				return ErrEncodeUnexpectedDataType
			}
		default:
			return ErrEncodeUnexpectedDefaultType
		}
	} else {
		c.Rest = nil
	}

	return rlp.Encode(w, c.column)
}

// DecodeRLP decodes column with rlp decode.
func (c *Column) DecodeRLP(s *rlp.Stream) error {
	defer func() { c.Rest = nil }()

	err := s.Decode(&c.column)
	if err != nil {
		return err
	}

	switch rest := c.Rest.(type) {
	case []interface{}:
		// nil is converted to empty list by encoder, while empty list is
		// converted to []interface{} by decoder.
		// So we view this case as nil and skip it.
		c.Default = nil
	case []byte:
		major, _ := ast.DecomposeDataType(c.Type)
		switch major {
		case ast.DataTypeMajorBool:
			if rest[0] == 1 {
				c.Default = true
			} else {
				c.Default = false
			}
		case ast.DataTypeMajorAddress,
			ast.DataTypeMajorFixedBytes,
			ast.DataTypeMajorDynamicBytes:
			c.Default = rest
		default:
			d, ok := ast.DecimalDecode(c.Type, rest)
			if !ok {
				return ErrDecodeUnexpectedDataType
			}
			c.Default = d
		}
	default:
		return ErrDecodeUnexpectedDefaultType
	}

	return nil
}

// FunctionDescriptor identifies a function.
type FunctionDescriptor struct {
	Function FunctionRef
}

var _ ast.IdentifierDescriptor = (*FunctionDescriptor)(nil)

// GetDescriptor is a useless function to satisfy the interface.
func (d FunctionDescriptor) GetDescriptor() uint32 {
	return uint32(0)<<24 | uint32(d.Function)<<8
}

// TableDescriptor identifies a table in a schema by an array index.
type TableDescriptor struct {
	Table TableRef
}

var _ ast.IdentifierDescriptor = (*TableDescriptor)(nil)

// GetDescriptor is a useless function to satisfy the interface.
func (d TableDescriptor) GetDescriptor() uint32 {
	return uint32(1)<<24 | uint32(d.Table)<<16
}

// ColumnDescriptor identifies a column in a schema by array indices.
type ColumnDescriptor struct {
	Table  TableRef
	Column ColumnRef
}

var _ ast.IdentifierDescriptor = (*ColumnDescriptor)(nil)

// GetDescriptor is a useless function to satisfy the interface.
func (d ColumnDescriptor) GetDescriptor() uint32 {
	return uint32(2)<<24 | uint32(d.Table)<<16 | uint32(d.Column)<<8
}

// IndexDescriptor identifies a index in a schema by array indices.
type IndexDescriptor struct {
	Table TableRef
	Index IndexRef
}

var _ ast.IdentifierDescriptor = (*IndexDescriptor)(nil)

// GetDescriptor is a useless function to satisfy the interface.
func (d IndexDescriptor) GetDescriptor() uint32 {
	return uint32(3)<<24 | uint32(d.Table)<<16 | uint32(d.Index)<<8
}

// SelectColumnDescriptor identifies a column specified in a select command by
// an array index.
type SelectColumnDescriptor struct {
	SelectColumn SelectColumnRef
}

var _ ast.IdentifierDescriptor = (*SelectColumnDescriptor)(nil)

// GetDescriptor is a useless function to satisfy the interface.
func (d SelectColumnDescriptor) GetDescriptor() uint32 {
	return uint32(4)<<24 | uint32(d.SelectColumn)<<8
}
