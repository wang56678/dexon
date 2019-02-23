package runtime

import "github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"

type tmplData struct {
	BinOpCollections []tmplTestCollection
}

type tmplTestCollection struct {
	TestName string
	Cases    []tmplTestCase
	OpFunc   string
}

type tmplTestCase struct {
	Name   string
	OpCode string
	Inputs []tmplOp
	Output tmplOp
	Error  string
}

type tmplOp struct {
	Im    bool
	Metas []tmplOpMeta
	Data  []string
}

type tmplOpMeta struct {
	Major ast.DataTypeMajor
	Minor ast.DataTypeMinor
}

const binOpTmplStr = `
// Code generated - DO NOT EDIT.

package runtime

import (
	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

{{range .BinOpCollections}}
func (s *instructionSuite) Test{{.TestName}}() {
	testcases := []opTestcase{ {{range .Cases}}
		{
			"{{.Name}}",
			Instruction{
				Op: {{.OpCode}},
				Input: []*Operand{ {{range .Inputs}}
					makeOperand(
						{{.Im}},
						[]ast.DataType{
							{{range .Metas}}ast.ComposeDataType({{.Major}}, {{.Minor}}),{{end}}
						},
						[]Tuple{ {{range .Data}}
							{{.}},{{end}}
						},
					),{{end}}
				},
				Output: 0,
			},
			makeOperand(
				{{.Output.Im}},
				[]ast.DataType{
					{{range .Output.Metas}}ast.ComposeDataType({{.Major}}, {{.Minor}}),{{end}}
				},
				[]Tuple{ {{range .Output.Data}}
					{{.}},{{end}}
				},
			),
			{{.Error}},
		},{{end}}
	}

	s.run(testcases, {{.OpFunc}})
}
{{end}}
`
