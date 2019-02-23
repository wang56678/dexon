package runtime

import (
	"bytes"
	"go/format"
	"os"
	"text/template"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
)

// RenderOpTest render op test to test file.
func RenderOpTest(output string) (err error) {
	binOpT, err := template.New("binOp").Parse(binOpTmplStr)
	if err != nil {
		return
	}

	b := new(bytes.Buffer)

	err = binOpT.Execute(b, testData)
	if err != nil {
		return
	}

	src, err := format.Source(b.Bytes())
	if err != nil {
		return
	}

	f, err := os.Create(output)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write(src)
	return
}

var testData = tmplData{
	BinOpCollections: []tmplTestCollection{
		{
			TestName: "OpAdd", OpFunc: "opAdd",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "ADD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-2)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(3)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(4)}}",
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(20)}}",
							"{&Raw{Value: decimal.NewFromFloat(-20)}, &Raw{Value: decimal.NewFromFloat(13)}}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "ADD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(11)}, &Raw{Value: decimal.NewFromFloat(8)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(-9)}, &Raw{Value: decimal.NewFromFloat(-12)}, &Raw{Value: decimal.NewFromFloat(-20)}}",
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
						},
					},
				},
				{
					Name:  "Immediate 2",
					Error: "nil", OpCode: "ADD",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(11)}, &Raw{Value: decimal.NewFromFloat(8)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(-9)}, &Raw{Value: decimal.NewFromFloat(-12)}, &Raw{Value: decimal.NewFromFloat(-20)}}",
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
						},
					},
				},
				{
					Name:  "Overflow - Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "ADD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(127)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Overflow None Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "ADD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(126)}}",
								"{&Raw{Value: decimal.NewFromFloat(126)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Underflow - Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "ADD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-128)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Underflow None Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "ADD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-127)}}",
								"{&Raw{Value: decimal.NewFromFloat(-127)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
			},
		},
		// -- end of ADD
		{
			TestName: "OpSub", OpFunc: "opSub",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "SUB",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-2)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(3)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-4)}}",
							"{&Raw{Value: decimal.NewFromFloat(20)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(-20)}, &Raw{Value: decimal.NewFromFloat(7)}}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "SUB",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(9)}, &Raw{Value: decimal.NewFromFloat(12)}, &Raw{Value: decimal.NewFromFloat(20)}}",
							"{&Raw{Value: decimal.NewFromFloat(-11)}, &Raw{Value: decimal.NewFromFloat(-8)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(4)}, &Raw{Value: decimal.NewFromFloat(10)}}",
						},
					},
				},
				{
					Name:  "Immediate 2",
					Error: "nil", OpCode: "SUB",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(-9)}, &Raw{Value: decimal.NewFromFloat(-12)}, &Raw{Value: decimal.NewFromFloat(-20)}}",
							"{&Raw{Value: decimal.NewFromFloat(11)}, &Raw{Value: decimal.NewFromFloat(8)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-4)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
						},
					},
				},
				{
					Name:  "Overflow - Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "SUB",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(127)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Overflow None Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "SUB",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(126)}}",
								"{&Raw{Value: decimal.NewFromFloat(126)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Underflow - Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "SUB",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-128)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Underflow None Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "SUB",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-127)}}",
								"{&Raw{Value: decimal.NewFromFloat(-127)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
			},
		},
		// -- end of SUB
		{
			TestName: "OpMul", OpFunc: "opMul",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "MUL",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}}",
							"{&Raw{Value: decimal.NewFromFloat(4)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							"{&Raw{Value: decimal.NewFromFloat(-4)}, &Raw{Value: decimal.NewFromFloat(-100)}}",
							"{&Raw{Value: decimal.NewFromFloat(100)}, &Raw{Value: decimal.NewFromFloat(100)}}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "MUL",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-20)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(20)}, &Raw{Value: decimal.NewFromFloat(0)}}",
						},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "MUL",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-20)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(20)}, &Raw{Value: decimal.NewFromFloat(0)}}",
						},
					},
				},
				{
					Name:  "Overflow - Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "MUL",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(127)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Overflow None Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "MUL",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(126)}}",
								"{&Raw{Value: decimal.NewFromFloat(126)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Underflow - Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "MUL",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-128)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Underflow None Immediate",
					Error: "errors.ErrorCodeUnderflow", OpCode: "MUL",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-127)}}",
								"{&Raw{Value: decimal.NewFromFloat(-127)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
			},
		},
		// -- end of MUL
		{
			TestName: "OpDiv", OpFunc: "opDiv",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "DIV",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}}",
							"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "DIV",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(13)}, &Raw{Value: decimal.NewFromFloat(13)}}",
								"{&Raw{Value: decimal.NewFromFloat(-13)}, &Raw{Value: decimal.NewFromFloat(-13)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}}",
							"{&Raw{Value: decimal.NewFromFloat(-5)}, &Raw{Value: decimal.NewFromFloat(5)}}",
							"{&Raw{Value: decimal.NewFromFloat(6)}, &Raw{Value: decimal.NewFromFloat(-6)}}",
							"{&Raw{Value: decimal.NewFromFloat(-6)}, &Raw{Value: decimal.NewFromFloat(6)}}",
						},
					},
				},
				{
					Name:  "Immediate 2",
					Error: "nil", OpCode: "DIV",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(50)}, &Raw{Value: decimal.NewFromFloat(-50)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(9)}, &Raw{Value: decimal.NewFromFloat(9)}}",
								"{&Raw{Value: decimal.NewFromFloat(-9)}, &Raw{Value: decimal.NewFromFloat(-9)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}}",
							"{&Raw{Value: decimal.NewFromFloat(-5)}, &Raw{Value: decimal.NewFromFloat(5)}}",
							"{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}}",
							"{&Raw{Value: decimal.NewFromFloat(-5)}, &Raw{Value: decimal.NewFromFloat(5)}}",
						},
					},
				},
				{
					Name:  "DivideByZero Immediate",
					Error: "errors.ErrorCodeDividedByZero", OpCode: "DIV",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "DivideByZero None Immediate",
					Error: "errors.ErrorCodeDividedByZero", OpCode: "DIV",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Overflow - Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "DIV",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-128)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "Overflow None Immediate",
					Error: "errors.ErrorCodeOverflow", OpCode: "DIV",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-128)}}",
								"{&Raw{Value: decimal.NewFromFloat(-128)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-2)}}",
							},
						},
					},
					Output: tmplOp{},
				},
			},
		},
		// -- end of DIV
		{
			TestName: "OpMod", OpFunc: "opMod",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "MOD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(3)}, &Raw{Value: decimal.NewFromFloat(3)}}",
								"{&Raw{Value: decimal.NewFromFloat(-3)}, &Raw{Value: decimal.NewFromFloat(-3)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}}",
							"{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "MOD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(13)}, &Raw{Value: decimal.NewFromFloat(13)}}",
								"{&Raw{Value: decimal.NewFromFloat(-13)}, &Raw{Value: decimal.NewFromFloat(-13)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(3)}, &Raw{Value: decimal.NewFromFloat(-3)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
							"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
							"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
						},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "MOD",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(31)}, &Raw{Value: decimal.NewFromFloat(-31)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}}",
								"{&Raw{Value: decimal.NewFromFloat(13)}, &Raw{Value: decimal.NewFromFloat(13)}}",
								"{&Raw{Value: decimal.NewFromFloat(-13)}, &Raw{Value: decimal.NewFromFloat(-13)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							"{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}}",
							"{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}}",
						},
					},
				},
				{
					Name:  "ModideByZero Immediate",
					Error: "errors.ErrorCodeDividedByZero", OpCode: "MOD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{},
				},
				{
					Name:  "ModideByZero None Immediate",
					Error: "errors.ErrorCodeDividedByZero", OpCode: "MOD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
								"{&Raw{Value: decimal.NewFromFloat(10)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{},
				},
			},
		},
		// -- end of MOD
		{
			TestName: "OpLt", OpFunc: "opLt",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "LT",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawTrue, rawTrue}",
							"{rawFalse, rawFalse, rawTrue}",
							"{rawFalse, rawFalse, rawFalse}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "LT",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawFalse, rawTrue}",
						},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "LT",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawTrue, rawFalse}",
						},
					},
				},
			},
		},
		// -- end of LT
		{
			TestName: "OpGt", OpFunc: "opGt",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "GT",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
								"{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawFalse, rawFalse}",
							"{rawTrue, rawFalse, rawFalse}",
							"{rawTrue, rawTrue, rawFalse}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "GT",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawTrue, rawFalse}",
						},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "GT",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawFalse, rawTrue}",
						},
					},
				},
			},
		},
		// -- end of GT
		{
			TestName: "OpEq", OpFunc: "opEq",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "EQ",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawTrue, rawTrue}",
							"{rawTrue, rawFalse, rawFalse}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "EQ",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawTrue, rawTrue}",
							"{rawTrue, rawFalse, rawFalse}",
						},
					},
				},
			},
		},
		// -- end of EQ
		{
			TestName: "OpAnd", OpFunc: "opAnd",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "AND",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawTrue}",
								"{rawFalse, rawFalse}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawFalse}",
							"{rawFalse, rawFalse}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "AND",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawTrue}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawFalse}",
							"{rawFalse, rawTrue}",
						},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "AND",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawTrue}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawFalse}",
							"{rawFalse, rawTrue}",
						},
					},
				},
				{
					Name:  "Invalid Data Type",
					Error: "errors.ErrorCodeInvalidDataType", OpCode: "AND",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
					},
					Output: tmplOp{},
				},
			},
		},
		// -- end of AND
		{
			TestName: "OpOr", OpFunc: "opOr",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "OR",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawTrue}",
								"{rawFalse, rawFalse}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawTrue}",
							"{rawFalse, rawTrue}",
						},
					},
				},
				{
					Name:  "Immediate",
					Error: "nil", OpCode: "OR",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawTrue}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawTrue}",
							"{rawTrue, rawTrue}",
						},
					},
				},
				{
					Name:  "Immediate - 2",
					Error: "nil", OpCode: "OR",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawTrue}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawTrue}",
							"{rawTrue, rawTrue}",
						},
					},
				},
				{
					Name:  "Invalid Data Type",
					Error: "errors.ErrorCodeInvalidDataType", OpCode: "OR",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
					},
					Output: tmplOp{},
				},
			},
		},
		// -- end of OR
		{
			TestName: "OpNot", OpFunc: "opNot",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "NOT",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawTrue}",
							"{rawTrue, rawFalse}",
						},
					},
				},
				{
					Name:  "Errors Invalid Data Type",
					Error: "errors.ErrorCodeInvalidDataType", OpCode: "NOT",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}}",
							},
						},
					},
					Output: tmplOp{},
				},
			},
		},
		// -- end of NOT
		{
			TestName: "OpUnion", OpFunc: "opUnion",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "UNION",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawTrue}",
								"{rawFalse, rawFalse}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawFalse}",
							"{rawFalse, rawTrue}",
							"{rawTrue, rawFalse}",
							"{rawTrue, rawTrue}",
						},
					},
				},
			},
		},
		// -- end of UNION
		{
			TestName: "OpIntxn", OpFunc: "opIntxn",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate",
					Error: "nil", OpCode: "INTXN",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
								"{rawTrue, rawTrue}",
								"{rawFalse, rawFalse}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawTrue}",
								"{rawFalse, rawFalse}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawFalse, rawFalse}",
							"{rawTrue, rawTrue}",
						},
					},
				},
			},
		},
		// -- end of INTXN
		{
			TestName: "OpLike", OpFunc: "opLike",
			Cases: []tmplTestCase{
				{
					Name:  `Like %\\%b% escape \\`, // \\ is raw string escape for \
					Error: "nil", OpCode: "LIKE",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("a%bcdefg")}, &Raw{Bytes: []byte("gfedcba")}}`,
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{`{&Raw{Bytes: []byte("%\\%b%")}}`},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{`{&Raw{Bytes: []byte("\\")}}`},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawFalse}",
						},
					},
				},
				{
					Name:  `Like t1 escape t2`,
					Error: "nil", OpCode: "LIKE",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("a%bcdefg")}}`,
								`{&Raw{Bytes: []byte("gfedcba")}}`,
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("%\\%b%")}}`,
								`{&Raw{Bytes: []byte("_fed%")}}`,
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("\\")}}`,
								`{&Raw{Bytes: []byte("")}}`,
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue}",
							"{rawTrue}",
						},
					},
				},
				{
					Name:  `Like with valid and invalid UTF8`,
					Error: "nil", OpCode: "LIKE",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte{226, 40, 161, 228, 189, 160, 229, 165, 189}}, &Raw{Bytes: []byte("gfedcba")}}`,
								// "\xe2(\xa1你好"
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte{37, 228, 189, 160, 37}}}`,
								// "%你%"
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawFalse}",
						},
					},
				},
			},
		},
		// -- end of LIKE
		{
			TestName: "OpZip", OpFunc: "opZip",
			Cases: []tmplTestCase{
				{
					Name:  "Zip two array",
					Error: "nil", OpCode: "ZIP",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}}`,
								`{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}}`,
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, rawTrue}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, rawFalse}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							`{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue}`,
							`{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse}`,
						},
					},
				},
				{
					Name:  "Zip immediate",
					Error: "nil", OpCode: "ZIP",
					Inputs: []tmplOp{
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}}`,
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, rawTrue}",
								"{&Raw{Value: decimal.NewFromFloat(2)}, rawFalse}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							`{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue}`,
							`{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse}`,
						},
					},
				},
			},
		},
		// -- end of ZIP
		{
			TestName: "OpField", OpFunc: "opField",
			Cases: []tmplTestCase{
				{
					Name:  "Retrieve 2nd,3rd column",
					Error: "nil", OpCode: "FIELD",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue}`,
								`{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse}`,
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							`{&Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}}`,
							`{&Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}}`,
						},
					},
				},
			},
		},
		// -- end of FIELD
		{
			TestName: "OpPrune", OpFunc: "opPrune",
			Cases: []tmplTestCase{
				{
					Name:  "Prune 2nd,4th,5th column",
					Error: "nil", OpCode: "PRUNE",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse, rawTrue}`,
								`{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue, rawFalse}`,
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(3)}, &Raw{Value: decimal.NewFromFloat(4)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
						},
						Data: []string{
							`{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Value: decimal.NewFromFloat(1)}}`,
							`{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Value: decimal.NewFromFloat(2)}}`,
						},
					},
				},
			},
		},
		// -- end of PRUNE
		{
			TestName: "OpFilter", OpFunc: "opFilter",
			Cases: []tmplTestCase{
				{
					Name:  "Filter first 2 rows",
					Error: "nil", OpCode: "FILTER",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue, rawFalse}",
								"{rawFalse, rawTrue}",
								"{rawTrue, rawTrue}",
								"{rawFalse, rawFalse}",
							},
						},
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue}",
								"{rawTrue}",
								"{rawFalse}",
								"{rawFalse}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawFalse}",
							"{rawFalse, rawTrue}",
						},
					},
				},
			},
		},
		// -- end of FILTER
		{
			TestName: "OpCast", OpFunc: "opCast",
			Cases: []tmplTestCase{
				{
					Name:  "None Immediate - int",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 1}, // int16 -> int8
								{Major: ast.DataTypeMajorInt, Minor: 1}, // int16 -> int24
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(127)}, &Raw{Value: decimal.NewFromFloat(127)}}",
								"{&Raw{Value: decimal.NewFromFloat(-128)}, &Raw{Value: decimal.NewFromFloat(-128)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 2},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 2},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(127)}, &Raw{Value: decimal.NewFromFloat(127)}}",
							"{&Raw{Value: decimal.NewFromFloat(-128)}, &Raw{Value: decimal.NewFromFloat(-128)}}",
						},
					},
				},
				{
					Name:  "None Immediate - int2",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 1}, // int16 -> uint16
								{Major: ast.DataTypeMajorInt, Minor: 1}, // int16 -> uint16
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(-32768)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorUint, Minor: 1},
								{Major: ast.DataTypeMajorUint, Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorUint, Minor: 1},
							{Major: ast.DataTypeMajorUint, Minor: 1},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(32768)}}",
						},
					},
				},
				{
					Name:  "None Immediate - int3",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 1}, // int16 -> bool
								{Major: ast.DataTypeMajorInt, Minor: 1}, // int16 -> bool
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(-32768)}}",
								"{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawTrue}", "{rawFalse, rawFalse}",
						},
					},
				},
				{
					Name:  "None Immediate - int4",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 1}, // int16 -> bytes16
								{Major: ast.DataTypeMajorInt, Minor: 1}, // int16 -> address
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(-32768)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1},
								{Major: ast.DataTypeMajorAddress, Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorFixedBytes, Minor: 1},
							{Major: ast.DataTypeMajorAddress, Minor: 0},
						},
						Data: []string{
							"{&Raw{Bytes: []byte{0x7f, 0xff}}, &Raw{Bytes: []byte{255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,0x80,0x00}}}",
						},
					},
				},
				{
					Name:  "None Immediate - uint",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> uint8
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> uint24
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(128)}, &Raw{Value: decimal.NewFromFloat(128)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorUint, Minor: 0},
								{Major: ast.DataTypeMajorUint, Minor: 2},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorUint, Minor: 0},
							{Major: ast.DataTypeMajorUint, Minor: 2},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(128)}, &Raw{Value: decimal.NewFromFloat(128)}}",
						},
					},
				},
				{
					Name:  "None Immediate - uint2",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> int16
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> byte16
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(32768)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 1},
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 1},
							{Major: ast.DataTypeMajorFixedBytes, Minor: 1},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Bytes: []byte{0x80,0x00}}}",
						},
					},
				},
				{
					Name:  "None Immediate - uint3",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> bool
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> bool
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue, rawFalse}",
						},
					},
				},
				{
					Name:  "None Immediate - uint4",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> bytes
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> bytes
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1},
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorFixedBytes, Minor: 1},
							{Major: ast.DataTypeMajorFixedBytes, Minor: 1},
						},
						Data: []string{
							"{&Raw{Bytes: []byte{0x7f, 0xff}}, &Raw{Bytes: []byte{0x00, 0x00}}}",
						},
					},
				},
				{
					Name:  "None Immediate - uint5",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorUint, Minor: 1}, // uint16 -> address
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(32767)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorAddress, Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorAddress, Minor: 1},
						},
						Data: []string{
							"{&Raw{Bytes: []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0x7f,0xff}}}",
						},
					},
				},
				{
					Name:  "None Immediate - bytes",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1}, // byte16 -> byte8
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1}, // byte16 -> byte24
							},
							Data: []string{
								"{&Raw{Bytes: []byte{0xff, 0xff}}, &Raw{Bytes: []byte{0xff, 0xff}}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorFixedBytes, Minor: 0},
								{Major: ast.DataTypeMajorFixedBytes, Minor: 2},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorFixedBytes, Minor: 0},
							{Major: ast.DataTypeMajorFixedBytes, Minor: 2},
						},
						Data: []string{
							"{&Raw{Bytes: []byte{0xff}}, &Raw{Bytes: []byte{0xff, 0xff, 0x00}}}",
						},
					},
				},
				{
					Name:  "None Immediate - bytes2",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1}, // byte16 -> int16
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1}, // byte16 -> uint16
							},
							Data: []string{
								"{&Raw{Bytes: []byte{0x7f, 0xff}}, &Raw{Bytes: []byte{0x80, 0x00}}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 1},
								{Major: ast.DataTypeMajorUint, Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorInt, Minor: 1},
							{Major: ast.DataTypeMajorUint, Minor: 1},
						},
						Data: []string{
							"{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(32768)}}",
						},
					},
				},
				{
					Name:  "None Immediate - bytes3",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorFixedBytes, Minor: 1}, // byte16 -> dyn
							},
							Data: []string{
								"{&Raw{Bytes: []byte{0x7f, 0xff}}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 1},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 1},
						},
						Data: []string{
							"{&Raw{Bytes: []byte{0x7f, 0xff}}}",
						},
					},
				},
				{
					Name:  "Same type",
					Error: "nil", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								"{rawTrue}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							"{rawTrue}",
						},
					},
				},
				{
					Name:  "Error Invalid Type",
					Error: "errors.ErrorCodeInvalidCastType", OpCode: "CAST",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorInt, Minor: 2},
							},
							Data: []string{
								"{&Raw{Value: decimal.NewFromFloat(-32768)}}",
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							},
							Data: []string{},
						},
					},
					Output: tmplOp{
						Im:    false,
						Metas: []tmplOpMeta{},
						Data:  []string{},
					},
				},
			},
		},
		// -- end of CAST
		{
			TestName: "OpSort", OpFunc: "opSort",
			Cases: []tmplTestCase{
				{
					Name:  "Multi-column sorting",
					Error: "nil", OpCode: "SORT",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue}`,
								`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue}`,
								`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawTrue}`,
								`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse}`,
								`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse}`,
								`{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
								`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
								`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{rawFalse, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{rawTrue, &Raw{Value: decimal.NewFromFloat(2)}}",
								"{rawFalse, &Raw{Value: decimal.NewFromFloat(0)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							`{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
							`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
							`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
							`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawTrue}`,
							`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse}`,
							`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue}`,
							`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse}`,
							`{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue}`,
						},
					},
				},
				{
					Name:  "Multi-column sorting - 2",
					Error: "nil", OpCode: "SORT",
					Inputs: []tmplOp{
						{
							Im: false,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
								{Major: ast.DataTypeMajorBool, Minor: 0},
							},
							Data: []string{
								`{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue}`,
								`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue}`,
								`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawTrue}`,
								`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse}`,
								`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse}`,
								`{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
								`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
								`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
							},
						},
						{
							Im: true,
							Metas: []tmplOpMeta{
								{Major: ast.DataTypeMajorBool, Minor: 0},
								{Major: ast.DataTypeMajorInt, Minor: 0},
							},
							Data: []string{
								"{rawTrue, &Raw{Value: decimal.NewFromFloat(0)}}",
								"{rawTrue, &Raw{Value: decimal.NewFromFloat(1)}}",
								"{rawFalse, &Raw{Value: decimal.NewFromFloat(2)}}",
							},
						},
					},
					Output: tmplOp{
						Im: false,
						Metas: []tmplOpMeta{
							{Major: ast.DataTypeMajorDynamicBytes, Minor: 0},
							{Major: ast.DataTypeMajorInt, Minor: 0},
							{Major: ast.DataTypeMajorBool, Minor: 0},
						},
						Data: []string{
							`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse}`,
							`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawTrue}`,
							`{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
							`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue}`,
							`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse}`,
							`{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
							`{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue}`,
							`{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse}`,
						},
					},
				},
			},
		},
		// -- end of SORT
	},
}
