package schema

import (
	"errors"
	"io"

	"github.com/shopspring/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/rlp"
)

// Error defines for encode and decode.
var (
	ErrEncodeUnexpectedType = errors.New("encode unexpected type")
	ErrDecodeUnexpectedType = errors.New("decode unexpected type")
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

// TableRef defines the type for table index in Schema.
type TableRef uint8

// ColumnRef defines the type for column index in Table.Columns.
type ColumnRef uint8

// IndexRef defines the type for array index of Column.Indices.
type IndexRef uint8

// SequenceRef defines the type for sequence index in Table.
type SequenceRef uint8

// SelectColumnRef defines the type for column index in SelectStmtNode.Column.
type SelectColumnRef uint16

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

// Table defiens sqlvm table struct.
type Table struct {
	Name    []byte
	Columns []Column
	Indices []Index
}

// Index defines sqlvm index struct.
type Index struct {
	Name    []byte
	Attr    IndexAttr
	Columns []ColumnRef // Columns must be sorted in ascending order.
}

type column struct {
	Name          []byte
	Type          ast.DataType
	Attr          ColumnAttr
	ForeignTable  TableRef
	ForeignColumn ColumnRef
	Sequence      SequenceRef
	Rest          interface{}
}

// Column defines sqlvm index struct.
type Column struct {
	column
	Default interface{} // decimal.Decimal, bool, []byte
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
			var err error
			c.Rest, err = ast.DecimalEncode(c.Type, d)
			if err != nil {
				return err
			}
		default:
			return ErrEncodeUnexpectedType
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
	case []byte:
		major, _ := ast.DecomposeDataType(c.Type)
		switch major {
		case ast.DataTypeMajorBool:
			if rest[0] == 1 {
				c.Default = true
			} else {
				c.Default = false
			}
		case ast.DataTypeMajorFixedBytes, ast.DataTypeMajorDynamicBytes:
			c.Default = rest
		default:
			d, err := ast.DecimalDecode(c.Type, rest)
			if err != nil {
				return err
			}
			c.Default = d
		}
	default:
		return ErrDecodeUnexpectedType
	}

	return nil
}

// TableDescriptor identifies a table in a schema by an array index.
type TableDescriptor struct {
	Table TableRef
}

// ColumnDescriptor identifies a column in a schema by array indices.
type ColumnDescriptor struct {
	Table  TableRef
	Column ColumnRef
}

// IndexDescriptor identifies a index in a schema by array indices.
type IndexDescriptor struct {
	Table TableRef
	Index IndexRef
}

// SelectColumnDescriptor identifies a column specified in a select command by
// an array index.
type SelectColumnDescriptor struct {
	SelectColumn SelectColumnRef
}
