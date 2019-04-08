// Code generated - DO NOT EDIT.

package runtime

import (
	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
)

func (s *instructionSuite) TestOpAdd() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: ADD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-2)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(10)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}},
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(3)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(4)}},
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(20)}},
					{&Raw{Value: decimal.NewFromFloat(-20)}, &Raw{Value: decimal.NewFromFloat(13)}},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: ADD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(11)}, &Raw{Value: decimal.NewFromFloat(8)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(-9)}, &Raw{Value: decimal.NewFromFloat(-12)}, &Raw{Value: decimal.NewFromFloat(-20)}},
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-10)}},
				},
			),
			nil,
		},
		{
			"Immediate 2",
			Instruction{
				Op: ADD,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(11)}, &Raw{Value: decimal.NewFromFloat(8)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(-9)}, &Raw{Value: decimal.NewFromFloat(-12)}, &Raw{Value: decimal.NewFromFloat(-20)}},
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-10)}},
				},
			),
			nil,
		},
		{
			"Overflow - Immediate",
			Instruction{
				Op: ADD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(127)}},
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeOverflow,
		},
		{
			"Overflow None Immediate",
			Instruction{
				Op: ADD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(126)}},
							{&Raw{Value: decimal.NewFromFloat(126)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeOverflow,
		},
		{
			"Underflow - Immediate",
			Instruction{
				Op: ADD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-128)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeUnderflow,
		},
		{
			"Underflow None Immediate",
			Instruction{
				Op: ADD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-127)}},
							{&Raw{Value: decimal.NewFromFloat(-127)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(-2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeUnderflow,
		},
	}

	s.run(testcases, opAdd)
}

func (s *instructionSuite) TestOpSub() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: SUB,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-2)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(10)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}},
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(3)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-4)}},
					{&Raw{Value: decimal.NewFromFloat(20)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(-20)}, &Raw{Value: decimal.NewFromFloat(7)}},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: SUB,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(9)}, &Raw{Value: decimal.NewFromFloat(12)}, &Raw{Value: decimal.NewFromFloat(20)}},
					{&Raw{Value: decimal.NewFromFloat(-11)}, &Raw{Value: decimal.NewFromFloat(-8)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(4)}, &Raw{Value: decimal.NewFromFloat(10)}},
				},
			),
			nil,
		},
		{
			"Immediate 2",
			Instruction{
				Op: SUB,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(-9)}, &Raw{Value: decimal.NewFromFloat(-12)}, &Raw{Value: decimal.NewFromFloat(-20)}},
					{&Raw{Value: decimal.NewFromFloat(11)}, &Raw{Value: decimal.NewFromFloat(8)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-4)}, &Raw{Value: decimal.NewFromFloat(-10)}},
				},
			),
			nil,
		},
		{
			"Overflow - Immediate",
			Instruction{
				Op: SUB,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(127)}},
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeOverflow,
		},
		{
			"Overflow None Immediate",
			Instruction{
				Op: SUB,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(126)}},
							{&Raw{Value: decimal.NewFromFloat(126)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(-2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeOverflow,
		},
		{
			"Underflow - Immediate",
			Instruction{
				Op: SUB,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-128)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeUnderflow,
		},
		{
			"Underflow None Immediate",
			Instruction{
				Op: SUB,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-127)}},
							{&Raw{Value: decimal.NewFromFloat(-127)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeUnderflow,
		},
	}

	s.run(testcases, opSub)
}

func (s *instructionSuite) TestOpMul() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: MUL,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}},
					{&Raw{Value: decimal.NewFromFloat(4)}, &Raw{Value: decimal.NewFromFloat(-1)}},
					{&Raw{Value: decimal.NewFromFloat(-4)}, &Raw{Value: decimal.NewFromFloat(-100)}},
					{&Raw{Value: decimal.NewFromFloat(100)}, &Raw{Value: decimal.NewFromFloat(100)}},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: MUL,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-20)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(20)}, &Raw{Value: decimal.NewFromFloat(0)}},
				},
			),
			nil,
		},
		{
			"Immediate - 2",
			Instruction{
				Op: MUL,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-20)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(20)}, &Raw{Value: decimal.NewFromFloat(0)}},
				},
			),
			nil,
		},
		{
			"Overflow - Immediate",
			Instruction{
				Op: MUL,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(127)}},
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeOverflow,
		},
		{
			"Overflow None Immediate",
			Instruction{
				Op: MUL,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(126)}},
							{&Raw{Value: decimal.NewFromFloat(126)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeOverflow,
		},
		{
			"Underflow - Immediate",
			Instruction{
				Op: MUL,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-128)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeUnderflow,
		},
		{
			"Underflow None Immediate",
			Instruction{
				Op: MUL,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-127)}},
							{&Raw{Value: decimal.NewFromFloat(-127)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeUnderflow,
		},
	}

	s.run(testcases, opMul)
}

func (s *instructionSuite) TestOpDiv() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: DIV,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(-2)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}},
					{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
					{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
					{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: DIV,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(13)}, &Raw{Value: decimal.NewFromFloat(13)}},
							{&Raw{Value: decimal.NewFromFloat(-13)}, &Raw{Value: decimal.NewFromFloat(-13)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}},
					{&Raw{Value: decimal.NewFromFloat(-5)}, &Raw{Value: decimal.NewFromFloat(5)}},
					{&Raw{Value: decimal.NewFromFloat(6)}, &Raw{Value: decimal.NewFromFloat(-6)}},
					{&Raw{Value: decimal.NewFromFloat(-6)}, &Raw{Value: decimal.NewFromFloat(6)}},
				},
			),
			nil,
		},
		{
			"Immediate 2",
			Instruction{
				Op: DIV,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(50)}, &Raw{Value: decimal.NewFromFloat(-50)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(9)}, &Raw{Value: decimal.NewFromFloat(9)}},
							{&Raw{Value: decimal.NewFromFloat(-9)}, &Raw{Value: decimal.NewFromFloat(-9)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}},
					{&Raw{Value: decimal.NewFromFloat(-5)}, &Raw{Value: decimal.NewFromFloat(5)}},
					{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}},
					{&Raw{Value: decimal.NewFromFloat(-5)}, &Raw{Value: decimal.NewFromFloat(5)}},
				},
			),
			nil,
		},
		{
			"DivideByZero Immediate",
			Instruction{
				Op: DIV,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeDividedByZero,
		},
		{
			"DivideByZero None Immediate",
			Instruction{
				Op: DIV,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeDividedByZero,
		},
		{
			"Overflow - Immediate",
			Instruction{
				Op: DIV,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(-128)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeOverflow,
		},
		{
			"Overflow None Immediate",
			Instruction{
				Op: DIV,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-128)}},
							{&Raw{Value: decimal.NewFromFloat(-128)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(-2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeOverflow,
		},
	}

	s.run(testcases, opDiv)
}

func (s *instructionSuite) TestOpMod() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: MOD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}},
							{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(3)}, &Raw{Value: decimal.NewFromFloat(3)}},
							{&Raw{Value: decimal.NewFromFloat(-3)}, &Raw{Value: decimal.NewFromFloat(-3)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}},
					{&Raw{Value: decimal.NewFromFloat(2)}, &Raw{Value: decimal.NewFromFloat(-2)}},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: MOD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(13)}, &Raw{Value: decimal.NewFromFloat(13)}},
							{&Raw{Value: decimal.NewFromFloat(-13)}, &Raw{Value: decimal.NewFromFloat(-13)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(3)}, &Raw{Value: decimal.NewFromFloat(-3)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
					{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
					{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
					{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
					{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
				},
			),
			nil,
		},
		{
			"Immediate - 2",
			Instruction{
				Op: MOD,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(31)}, &Raw{Value: decimal.NewFromFloat(-31)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}, &Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(-10)}, &Raw{Value: decimal.NewFromFloat(-10)}},
							{&Raw{Value: decimal.NewFromFloat(13)}, &Raw{Value: decimal.NewFromFloat(13)}},
							{&Raw{Value: decimal.NewFromFloat(-13)}, &Raw{Value: decimal.NewFromFloat(-13)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
					{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
					{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}},
					{&Raw{Value: decimal.NewFromFloat(5)}, &Raw{Value: decimal.NewFromFloat(-5)}},
				},
			),
			nil,
		},
		{
			"ModideByZero Immediate",
			Instruction{
				Op: MOD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeDividedByZero,
		},
		{
			"ModideByZero None Immediate",
			Instruction{
				Op: MOD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(10)}},
							{&Raw{Value: decimal.NewFromFloat(10)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeDividedByZero,
		},
	}

	s.run(testcases, opMod)
}

func (s *instructionSuite) TestOpLt() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: LT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawTrue, rawTrue},
					{rawFalse, rawFalse, rawTrue},
					{rawFalse, rawFalse, rawFalse},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: LT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawFalse, rawTrue},
				},
			),
			nil,
		},
		{
			"Immediate - 2",
			Instruction{
				Op: LT,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawTrue, rawFalse},
				},
			),
			nil,
		},
	}

	s.run(testcases, opLt)
}

func (s *instructionSuite) TestOpGt() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: GT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
							{&Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawFalse, rawFalse},
					{rawTrue, rawFalse, rawFalse},
					{rawTrue, rawTrue, rawFalse},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: GT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawTrue, rawFalse},
				},
			),
			nil,
		},
		{
			"Immediate - 2",
			Instruction{
				Op: GT,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawFalse, rawTrue},
				},
			),
			nil,
		},
	}

	s.run(testcases, opGt)
}

func (s *instructionSuite) TestOpEq() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: EQ,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(-1)}},
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawTrue, rawTrue},
					{rawTrue, rawFalse, rawFalse},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: EQ,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(-1)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawTrue, rawTrue},
					{rawTrue, rawFalse, rawFalse},
				},
			),
			nil,
		},
	}

	s.run(testcases, opEq)
}

func (s *instructionSuite) TestOpAnd() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: AND,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawTrue},
							{rawFalse, rawFalse},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawFalse},
					{rawFalse, rawFalse},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: AND,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawTrue},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawFalse},
					{rawFalse, rawTrue},
				},
			),
			nil,
		},
		{
			"Immediate - 2",
			Instruction{
				Op: AND,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawTrue},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawFalse},
					{rawFalse, rawTrue},
				},
			),
			nil,
		},
		{
			"Invalid Data Type",
			Instruction{
				Op: AND,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeInvalidDataType,
		},
	}

	s.run(testcases, opAnd)
}

func (s *instructionSuite) TestOpOr() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: OR,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawTrue},
							{rawFalse, rawFalse},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawTrue},
					{rawFalse, rawTrue},
				},
			),
			nil,
		},
		{
			"Immediate",
			Instruction{
				Op: OR,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawTrue},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawTrue},
					{rawTrue, rawTrue},
				},
			),
			nil,
		},
		{
			"Immediate - 2",
			Instruction{
				Op: OR,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawTrue},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawTrue},
					{rawTrue, rawTrue},
				},
			),
			nil,
		},
		{
			"Invalid Data Type",
			Instruction{
				Op: OR,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeInvalidDataType,
		},
	}

	s.run(testcases, opOr)
}

func (s *instructionSuite) TestOpNot() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: NOT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawTrue},
					{rawTrue, rawFalse},
				},
			),
			nil,
		},
		{
			"Errors Invalid Data Type",
			Instruction{
				Op: NOT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeInvalidDataType,
		},
	}

	s.run(testcases, opNot)
}

func (s *instructionSuite) TestOpUnion() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: UNION,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawTrue},
							{rawFalse, rawFalse},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawFalse},
					{rawFalse, rawTrue},
					{rawTrue, rawFalse},
					{rawTrue, rawTrue},
				},
			),
			nil,
		},
	}

	s.run(testcases, opUnion)
}

func (s *instructionSuite) TestOpIntxn() {
	testcases := []opTestcase{
		{
			"None Immediate",
			Instruction{
				Op: INTXN,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
							{rawTrue, rawTrue},
							{rawFalse, rawFalse},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawTrue},
							{rawFalse, rawFalse},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawFalse, rawFalse},
					{rawTrue, rawTrue},
				},
			),
			nil,
		},
	}

	s.run(testcases, opIntxn)
}

func (s *instructionSuite) TestOpLike() {
	testcases := []opTestcase{
		{
			"Like %\\%b% escape \\",
			Instruction{
				Op: LIKE,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("a%bcdefg")}, &Raw{Bytes: []byte("gfedcba")}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("%\\%b%")}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("\\")}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawFalse},
				},
			),
			nil,
		},
		{
			"Like t1 escape t2",
			Instruction{
				Op: LIKE,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("a%bcdefg")}},
							{&Raw{Bytes: []byte("gfedcba")}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("%\\%b%")}},
							{&Raw{Bytes: []byte("_fed%")}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("\\")}},
							{&Raw{Bytes: []byte("")}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue},
					{rawTrue},
				},
			),
			nil,
		},
		{
			"Like with valid and invalid UTF8",
			Instruction{
				Op: LIKE,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte{226, 40, 161, 228, 189, 160, 229, 165, 189}}, &Raw{Bytes: []byte("gfedcba")}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte{37, 228, 189, 160, 37}}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawFalse},
				},
			),
			nil,
		},
	}

	s.run(testcases, opLike)
}

func (s *instructionSuite) TestOpZip() {
	testcases := []opTestcase{
		{
			"Zip two array",
			Instruction{
				Op: ZIP,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
					{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
				},
			),
			nil,
		},
		{
			"Zip immediate",
			Instruction{
				Op: ZIP,
				Input: []*Operand{
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
					{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
				},
			),
			nil,
		},
	}

	s.run(testcases, opZip)
}

func (s *instructionSuite) TestOpField() {
	testcases := []opTestcase{
		{
			"Retrieve 2nd,3rd column",
			Instruction{
				Op: FIELD,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}},
					{&Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}},
				},
			),
			nil,
		},
	}

	s.run(testcases, opField)
}

func (s *instructionSuite) TestOpPrune() {
	testcases := []opTestcase{
		{
			"Prune 2nd,4th,5th column",
			Instruction{
				Op: PRUNE,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse, rawTrue},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(3)}, &Raw{Value: decimal.NewFromFloat(4)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Value: decimal.NewFromFloat(1)}},
					{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Value: decimal.NewFromFloat(2)}},
				},
			),
			nil,
		},
	}

	s.run(testcases, opPrune)
}

func (s *instructionSuite) TestOpCut() {
	testcases := []opTestcase{
		{
			"Cut 2nd to 4th columns",
			Instruction{
				Op: CUT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse, rawTrue},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(3)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("abcdefg-1")}, rawTrue},
					{&Raw{Bytes: []byte("abcdefg-2")}, rawFalse},
				},
			),
			nil,
		},
		{
			"Cut 1st column",
			Instruction{
				Op: CUT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse, rawTrue},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse, rawTrue},
					{&Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue, rawFalse},
				},
			),
			nil,
		},
		{
			"Cut since 2nd column",
			Instruction{
				Op: CUT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("abcdefg-1")}},
					{&Raw{Bytes: []byte("abcdefg-2")}},
				},
			),
			nil,
		},
		{
			"Cut all columns",
			Instruction{
				Op: CUT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{
					{},
					{},
				},
			),
			nil,
		},
		{
			"Cut error range - 1",
			Instruction{
				Op: CUT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(5)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeIndexOutOfRange,
		},
		{
			"Cut error range - 2",
			Instruction{
				Op: CUT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("abcdefg-1")}, &Raw{Bytes: []byte("gfedcba-1")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Bytes: []byte("abcdefg-2")}, &Raw{Bytes: []byte("gfedcba-2")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(15)}, &Raw{Value: decimal.NewFromFloat(17)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeIndexOutOfRange,
		},
	}

	s.run(testcases, opCut)
}

func (s *instructionSuite) TestOpFilter() {
	testcases := []opTestcase{
		{
			"Filter first 2 rows",
			Instruction{
				Op: FILTER,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue, rawFalse},
							{rawFalse, rawTrue},
							{rawTrue, rawTrue},
							{rawFalse, rawFalse},
						},
					),
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue},
							{rawTrue},
							{rawFalse},
							{rawFalse},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawFalse},
					{rawFalse, rawTrue},
				},
			),
			nil,
		},
	}

	s.run(testcases, opFilter)
}

func (s *instructionSuite) TestOpCast() {
	testcases := []opTestcase{
		{
			"None Immediate - int",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 1), ast.ComposeDataType(ast.DataTypeMajorInt, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(127)}, &Raw{Value: decimal.NewFromFloat(127)}},
							{&Raw{Value: decimal.NewFromFloat(-128)}, &Raw{Value: decimal.NewFromFloat(-128)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 2),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 2),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(127)}, &Raw{Value: decimal.NewFromFloat(127)}},
					{&Raw{Value: decimal.NewFromFloat(-128)}, &Raw{Value: decimal.NewFromFloat(-128)}},
				},
			),
			nil,
		},
		{
			"None Immediate - int2",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 1), ast.ComposeDataType(ast.DataTypeMajorInt, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(-32768)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(32768)}},
				},
			),
			nil,
		},
		{
			"None Immediate - int3",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 1), ast.ComposeDataType(ast.DataTypeMajorInt, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(-32768)}},
							{&Raw{Value: decimal.NewFromFloat(0)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawTrue},
					{rawFalse, rawFalse},
				},
			),
			nil,
		},
		{
			"None Immediate - int4",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 1), ast.ComposeDataType(ast.DataTypeMajorInt, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(-32768)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1), ast.ComposeDataType(ast.DataTypeMajorAddress, 0),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1), ast.ComposeDataType(ast.DataTypeMajorAddress, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte{0x7f, 0xff}}, &Raw{Bytes: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0x80, 0x00}}},
				},
			),
			nil,
		},
		{
			"None Immediate - uint",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(128)}, &Raw{Value: decimal.NewFromFloat(128)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 0), ast.ComposeDataType(ast.DataTypeMajorUint, 2),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorUint, 0), ast.ComposeDataType(ast.DataTypeMajorUint, 2),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(128)}, &Raw{Value: decimal.NewFromFloat(128)}},
				},
			),
			nil,
		},
		{
			"None Immediate - uint2",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(32768)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 1), ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 1), ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Bytes: []byte{0x80, 0x00}}},
				},
			),
			nil,
		},
		{
			"None Immediate - uint3",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue, rawFalse},
				},
			),
			nil,
		},
		{
			"None Immediate - uint4",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1), ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1), ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1),
				},
				[]Tuple{
					{&Raw{Bytes: []byte{0x7f, 0xff}}, &Raw{Bytes: []byte{0x00, 0x00}}},
				},
			),
			nil,
		},
		{
			"None Immediate - uint5",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(32767)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorAddress, 1),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorAddress, 1),
				},
				[]Tuple{
					{&Raw{Bytes: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x7f, 0xff}}},
				},
			),
			nil,
		},
		{
			"None Immediate - bytes",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1), ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1),
						},
						[]Tuple{
							{&Raw{Bytes: []byte{0xff, 0xff}}, &Raw{Bytes: []byte{0xff, 0xff}}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 0), ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 2),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 0), ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 2),
				},
				[]Tuple{
					{&Raw{Bytes: []byte{0xff}}, &Raw{Bytes: []byte{0xff, 0xff, 0x00}}},
				},
			),
			nil,
		},
		{
			"None Immediate - bytes2",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1), ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1),
						},
						[]Tuple{
							{&Raw{Bytes: []byte{0x7f, 0xff}}, &Raw{Bytes: []byte{0x80, 0x00}}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(32767)}, &Raw{Value: decimal.NewFromFloat(32768)}},
				},
			),
			nil,
		},
		{
			"None Immediate - bytes3",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorFixedBytes, 1),
						},
						[]Tuple{
							{&Raw{Bytes: []byte{0x7f, 0xff}}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 1),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 1),
				},
				[]Tuple{
					{&Raw{Bytes: []byte{0x7f, 0xff}}},
				},
			),
			nil,
		},
		{
			"Same type",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{rawTrue},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{rawTrue},
				},
			),
			nil,
		},
		{
			"Error Invalid Type",
			Instruction{
				Op: CAST,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 2),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(-32768)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0),
						},
						[]Tuple{},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{},
				[]Tuple{},
			),
			errors.ErrorCodeInvalidCastType,
		},
	}

	s.run(testcases, opCast)
}

func (s *instructionSuite) TestOpSort() {
	testcases := []opTestcase{
		{
			"Multi-column sorting",
			Instruction{
				Op: SORT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue},
							{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawTrue},
							{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse},
							{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
							{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
							{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
							{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{rawFalse, &Raw{Value: decimal.NewFromFloat(1)}},
							{rawTrue, &Raw{Value: decimal.NewFromFloat(2)}},
							{rawFalse, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
					{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
					{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
					{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawTrue},
					{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
					{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue},
					{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse},
					{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
				},
			),
			nil,
		},
		{
			"Multi-column sorting - 2",
			Instruction{
				Op: SORT,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
						},
						[]Tuple{
							{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
							{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue},
							{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawTrue},
							{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse},
							{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
							{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
							{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
							{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorBool, 0), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{rawTrue, &Raw{Value: decimal.NewFromFloat(0)}},
							{rawTrue, &Raw{Value: decimal.NewFromFloat(1)}},
							{rawFalse, &Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorDynamicBytes, 0), ast.ComposeDataType(ast.DataTypeMajorInt, 0), ast.ComposeDataType(ast.DataTypeMajorBool, 0),
				},
				[]Tuple{
					{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(1)}, rawFalse},
					{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawTrue},
					{&Raw{Bytes: []byte("a")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
					{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawTrue},
					{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(2)}, rawFalse},
					{&Raw{Bytes: []byte("b")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
					{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(1)}, rawTrue},
					{&Raw{Bytes: []byte("c")}, &Raw{Value: decimal.NewFromFloat(3)}, rawFalse},
				},
			),
			nil,
		},
	}

	s.run(testcases, opSort)
}

func (s *instructionSuite) TestOpRange() {
	testcases := []opTestcase{
		{
			"Range test limit 2 offset 1",
			Instruction{
				Op: RANGE,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}},
							{&Raw{Value: decimal.NewFromFloat(3)}},
							{&Raw{Value: decimal.NewFromFloat(4)}},
							{&Raw{Value: decimal.NewFromFloat(5)}},
							{&Raw{Value: decimal.NewFromFloat(6)}},
							{&Raw{Value: decimal.NewFromFloat(7)}},
							{&Raw{Value: decimal.NewFromFloat(8)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{
					{&Raw{Value: decimal.NewFromFloat(2)}},
					{&Raw{Value: decimal.NewFromFloat(3)}},
				},
			),
			nil,
		},
		{
			"Range test limit 0 offset 1",
			Instruction{
				Op: RANGE,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}, &Raw{Value: decimal.NewFromFloat(0)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{},
			),
			nil,
		},
		{
			"Range test offset 20",
			Instruction{
				Op: RANGE,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(20)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{},
			),
			nil,
		},
		{
			"Range test limit 10 offset 20",
			Instruction{
				Op: RANGE,
				Input: []*Operand{
					makeOperand(
						false,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorInt, 0),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(1)}},
							{&Raw{Value: decimal.NewFromFloat(2)}},
						},
					),
					makeOperand(
						true,
						[]ast.DataType{
							ast.ComposeDataType(ast.DataTypeMajorUint, 1), ast.ComposeDataType(ast.DataTypeMajorUint, 1),
						},
						[]Tuple{
							{&Raw{Value: decimal.NewFromFloat(20)}, &Raw{Value: decimal.NewFromFloat(10)}},
						},
					),
				},
				Output: 0,
			},
			makeOperand(
				false,
				[]ast.DataType{
					ast.ComposeDataType(ast.DataTypeMajorInt, 0),
				},
				[]Tuple{},
			),
			nil,
		},
	}

	s.run(testcases, opRange)
}
