package runtime

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
)

type tmplData struct {
	BinOpCollections []*tmplTestCollection
}

func (d *tmplData) process() (err error) {
	for _, c := range d.BinOpCollections {
		err = c.process()
		if err != nil {
			return
		}
	}
	return
}

type tmplTestCollection struct {
	TestName string
	Cases    []*tmplTestCase
	OpFunc   string
}

func (c *tmplTestCollection) process() (err error) {
	for _, c := range c.Cases {
		err = c.process()
		if err != nil {
			return
		}
	}
	return
}

type tmplTestCase struct {
	Name   string
	OpCode string
	Inputs []*tmplOp
	Output *tmplOp
	Error  string
}

func (c *tmplTestCase) process() (err error) {
	for _, ic := range c.Inputs {
		err = ic.process()
		if err != nil {
			return
		}
	}

	err = c.Output.process()
	return
}

type tmplOp struct {
	Im    bool
	Metas []*tmplOpMeta
	Data  []string
}

func (o *tmplOp) process() (err error) {
	for i, r := range o.Data {
		o.Data[i] = processRaw(r)
	}
	return
}

type tmplOpMeta struct {
	Major string
	Minor ast.DataTypeMinor
}

// RenderOpTest render op test to test file.
func RenderOpTest(output string) (err error) {
	binOpT, err := template.New("binOp").Parse(binOpTmplStr)
	if err != nil {
		return
	}

	b := new(bytes.Buffer)

	err = testData.process()
	if err != nil {
		fmt.Printf("data process error: %v\n", err)
		return
	}

	err = binOpT.Execute(b, testData)
	if err != nil {
		fmt.Printf("template render error: %v\n", err)
		return
	}

	src, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Printf(
			`!!!!!
format source error: %v
note: render source for debugging
!!!!!`, err)
		src = b.Bytes()
	}

	f, err := os.Create(output)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write(src)
	return
}

func processRaw(raw string) (dsRaw string) {
	var (
		b    = &strings.Builder{}
		size = 1

		err error
	)

	for i := 0; i < len(raw); i += size {
		switch raw[i] {
		case 'V':
			size, err = writeV(b, raw, i)
		case 'B':
			size, err = writeB(b, raw, i)
		case 'T':
			_, err = b.WriteString("rawTrue")
			size = 1
		case 'F':
			_, err = b.WriteString("rawFalse")
			size = 1
		default:
			err = b.WriteByte(raw[i])
			size = 1
		}

		if err != nil {
			panic(err)
		}
	}
	dsRaw = b.String()
	return
}

func trim(raw string, i int) (size, j int) {
	j = i
	t := true
	for t {
		switch raw[j] {
		case ' ', '\t', '\n':
		case ':':
			t = false
		default:
			panic(fmt.Errorf("trim '%v' fail. char: '%c'", raw, raw[j]))
		}
		j++
		size++
	}
	return

}

func writeV(b *strings.Builder, raw string, i int) (size int, err error) {
	_, err = b.WriteString("&Raw{Value: decimal.NewFromFloat(")
	if err != nil {
		return
	}

	size, i = trim(raw, i+1)

VLOOP:
	for j := i; j < len(raw); j++ {
		size++
		switch raw[j] {
		case ',', '}':
			break VLOOP
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '.', ' ':
			b.WriteByte(raw[j])
		}
	}

	_, err = b.WriteString(")}")
	return
}

func writeB(b *strings.Builder, raw string, i int) (size int, err error) {
	_, err = b.WriteString("&Raw{Bytes: []byte")
	if err != nil {
		return
	}

	size, i = trim(raw, i+1)

	str := false
	level := 0

BLOOP:
	for j := i; j < len(raw); j++ {
		size++
		switch raw[j] {
		case ',':
			if level <= 0 {
				break BLOOP
			}
		case '{':
			level++
		case '}':
			level--
			if level <= 0 {
				break BLOOP
			}
		case '"':
			if !str {
				b.WriteByte(byte('('))
				str = true
			}
		}
		b.WriteByte(raw[j])
	}

	if str {
		b.WriteByte(byte(')'))
	}

	err = b.WriteByte(byte('}'))
	return
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
							{{range .Metas}}ast.ComposeDataType(ast.DataTypeMajor{{.Major}}, {{.Minor}}),{{end}}
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
					{{range .Output.Metas}}ast.ComposeDataType(ast.DataTypeMajor{{.Major}}, {{.Minor}}),{{end}}
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
