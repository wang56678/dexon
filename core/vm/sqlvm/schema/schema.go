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

// ColumnAttr enums.
const (
	ColumnAttrHasDefault ColumnAttr = 1 << iota
	ColumnAttrNotNull
	ColumnAttrHasSequence
	ColumnAttrHasForeignKey
)

// TableRef defines the type for table index in Schema.
type TableRef uint8

// ColumnRef defines the type for column index in Table.Columns.
type ColumnRef uint8

// IndexRef defines the type for array index of Column.Indices.
type IndexRef uint8

// SequenceRef defines the type for sequence index in Table.
type SequenceRef uint8

// IndexAttr defines bit flags for describing index attribute.
type IndexAttr uint16

// IndexAttr enums.
const (
	IndexAttrUnique IndexAttr = 1 << iota
)

// Schema defines sqlvm schema struct.
type Schema []*Table

// Table defiens sqlvm table struct.
type Table struct {
	Name    []byte
	Columns []*Column
	Indices []*Index
}

// Index defines sqlvm index struct.
type Index struct {
	Name    []byte
	Attr    IndexAttr
	Columns []ColumnRef
}

type column struct {
	Name          []byte
	Type          ast.DataType
	Attr          ColumnAttr
	Sequence      SequenceRef
	ForeignTable  TableRef
	ForeignColumn ColumnRef
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
