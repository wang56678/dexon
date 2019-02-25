package schema

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/rlp"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type SchemaTestSuite struct{ suite.Suite }

func (s *SchemaTestSuite) requireEncodeAndDecodeColumnNoError(c Column) {
	buffer := bytes.Buffer{}
	w := bufio.NewWriter(&buffer)
	s.Require().NoError(rlp.Encode(w, c))
	w.Flush()

	c2 := Column{}
	r := ioutil.NopCloser(bufio.NewReader(&buffer))
	s.Require().NoError(rlp.Decode(r, &c2))
	s.Require().Equal(c, c2)
}

func (s *SchemaTestSuite) TestEncodeAndDecodeColumn() {
	s.requireEncodeAndDecodeColumnNoError(Column{
		column: column{
			Name:     []byte("a"),
			Type:     ast.ComposeDataType(ast.DataTypeMajorBool, 0),
			Attr:     ColumnAttrHasSequence | ColumnAttrHasDefault,
			Sequence: 1,
		},
		Default: true,
	})

	s.requireEncodeAndDecodeColumnNoError(Column{
		column: column{
			Name: []byte("b"),
			Type: ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 0),
		},
	})

	s.requireEncodeAndDecodeColumnNoError(Column{
		column: column{
			Name: []byte("c"),
			Type: ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
			Attr: ColumnAttrNotNull | ColumnAttrHasDefault,
		},
		Default: []byte{},
	})

	s.requireEncodeAndDecodeColumnNoError(Column{
		column: column{
			Name: []byte("d"),
			Type: ast.ComposeDataType(ast.DataTypeMajorUint, 0),
			Attr: ColumnAttrNotNull | ColumnAttrHasDefault,
		},
		Default: decimal.New(1, 0),
	})
}

func (s *SchemaTestSuite) TestEncodeAndDecodeSchema() {
	schema := Schema{
		&Table{
			Name: []byte("test"),
			Columns: []*Column{
				{
					column: column{
						Name:     []byte("a"),
						Type:     ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						Attr:     ColumnAttrHasSequence | ColumnAttrHasDefault,
						Sequence: 1,
					},
					Default: true,
				},
			},
			Indices: []*Index{
				{
					Name:    []byte("idx"),
					Attr:    IndexAttrUnique,
					Columns: []ColumnRef{0},
				},
			},
		},
		&Table{
			Name: []byte("test2"),
		},
	}
	buffer := bytes.Buffer{}
	w := bufio.NewWriter(&buffer)
	s.Require().NoError(rlp.Encode(w, schema))
	w.Flush()

	schema2 := Schema{}
	r := ioutil.NopCloser(bufio.NewReader(&buffer))
	s.Require().NoError(rlp.Decode(r, &schema2))

	s.Require().Equal(len(schema), len(schema2))

	for i := 0; i < len(schema); i++ {
		table := schema[i]
		table2 := schema2[i]
		s.Require().Equal(table.Name, table2.Name)
		s.Require().Equal(len(table.Columns), len(table2.Columns))
		s.Require().Equal(len(table.Indices), len(table2.Indices))

		for j := 0; j < len(table.Columns); j++ {
			column := table.Columns[j]
			column2 := table.Columns[j]
			s.Require().Equal(*column, *column2)
		}

		for j := 0; j < len(table.Indices); j++ {
			index := table.Indices[j]
			index2 := table2.Indices[j]
			s.Require().Equal(*index, *index2)
		}
	}
}

func TestSchema(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}
