package schema

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/rlp"
)

type SchemaTestSuite struct{ suite.Suite }

func (s *SchemaTestSuite) normalizeEmptySlice(i interface{}) {
	var process func(reflect.Type, reflect.Value)
	process = func(t reflect.Type, v reflect.Value) {
		switch t.Kind() {
		case reflect.Ptr:
			process(t.Elem(), v.Elem())
		case reflect.Array:
			l := v.Len()
			for i := 0; i < l; i++ {
				process(t.Elem(), v.Index(i))
			}
		case reflect.Slice:
			if v.IsNil() {
				s := reflect.MakeSlice(t, 0, 0)
				v.Set(s)
			} else {
				l := v.Len()
				for i := 0; i < l; i++ {
					process(t.Elem(), v.Index(i))
				}
			}
		case reflect.Struct:
			l := t.NumField()
			for i := 0; i < l; i++ {
				ft := t.Field(i).Type
				fv := v.Field(i)
				process(ft, fv)
			}
		}
	}
	process(reflect.TypeOf(i), reflect.ValueOf(i))
}

func (s *SchemaTestSuite) requireEncodeAndDecodeColumnNoError(c Column) {
	buffer := bytes.Buffer{}
	w := bufio.NewWriter(&buffer)
	s.Require().NoError(rlp.Encode(w, c))
	w.Flush()

	c2 := Column{}
	r := ioutil.NopCloser(bufio.NewReader(&buffer))
	s.Require().NoError(rlp.Decode(r, &c2))

	s.normalizeEmptySlice(&c.column)
	s.normalizeEmptySlice(&c2.column)
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
		Table{
			Name: []byte("test"),
			Columns: []Column{
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
			Indices: []Index{
				{
					Name:    []byte("idx"),
					Attr:    IndexAttrUnique,
					Columns: []ColumnRef{0},
				},
			},
		},
		Table{
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
			s.Require().Equal(column, column2)
		}

		for j := 0; j < len(table.Indices); j++ {
			index := table.Indices[j]
			index2 := table2.Indices[j]
			s.Require().Equal(index, index2)
		}
	}
}

func (s *SchemaTestSuite) TestGetFieldType() {
	type testCase struct {
		fields        []uint8
		table         *Table
		expectedTypes []ast.DataType
		expectedLenth int
		expectedErr   error
	}
	testCases := []testCase{
		{
			fields: []uint8{uint8(1), uint8(0)},
			table: &Table{
				Name: []byte("Table_A"),
				Columns: []Column{
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorUint, 7),
						},
						nil,
					},
				},
			},
			expectedTypes: []ast.DataType{
				ast.ComposeDataType(ast.DataTypeMajorUint, 7),
				ast.ComposeDataType(ast.DataTypeMajorBool, 0),
			},
			expectedLenth: 2,
			expectedErr:   nil,
		},
		{
			fields: []uint8{uint8(8)},
			table: &Table{
				Name: []byte("Table_B"),
			},
			expectedLenth: 0,
			expectedErr:   errors.ErrorCodeIndexOutOfRange,
		},
	}
	for _, t := range testCases {
		length := t.expectedLenth
		expectedErr := t.expectedErr
		types, err := t.table.GetFieldType(t.fields)
		s.Require().Equal(length, len(types))
		s.Require().Equal(expectedErr, err)
		for i, tt := range types {
			s.Require().Equal(t.expectedTypes[i], tt)
		}
	}
}

func (s *SchemaTestSuite) TestSetupColumnOffset() {
	type testCase struct {
		name               string
		table              *Table
		expectedSlotOffest []uint8
		expectedByteOffset []uint8
	}
	testCases := []testCase{
		{
			name: "Table_A",
			table: &Table{
				Columns: []Column{
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorUint, 7),
						},
						nil,
					},
				},
			},
			expectedByteOffset: []uint8{0, 1},
			expectedSlotOffest: []uint8{0, 0},
		},
		{
			name: "Table_B",
			table: &Table{
				Columns: []Column{
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						nil,
					},
				},
			},
			expectedByteOffset: []uint8{0, 0, 0, 0},
			expectedSlotOffest: []uint8{0, 1, 2, 3},
		},
		{
			name: "Table_C",
			table: &Table{
				Columns: []Column{
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorUint, 7),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorUint, 7),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorUint, 15),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 30),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						nil,
					},
					{
						column{
							Type: ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 31),
						},
						nil,
					},
				},
			},
			expectedByteOffset: []uint8{0, 0, 8, 16, 0, 31, 0},
			expectedSlotOffest: []uint8{0, 1, 1, 1, 2, 2, 3},
		},
	}
	for i, t := range testCases {
		testCases[i].table.SetupColumnOffset()
		shift := 0
		for _, c := range testCases[i].table.Columns {
			s.Require().Equalf(t.expectedSlotOffest[shift],
				c.SlotOffset, "slotOffset not match. Name: %v, shift: %v", t.name, shift)
			s.Require().Equalf(t.expectedByteOffset[shift],
				c.ByteOffset, "byteOffset not match: Name: %v, shift: %v", t.name, shift)
			shift++
		}
	}
}

func TestSchema(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}
