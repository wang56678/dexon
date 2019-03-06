package codegen

import (
	"fmt"
	"reflect"

	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	dec "github.com/dexon-foundation/dexon/core/vm/sqlvm/common/decimal"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/planner"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/runtime"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

type codegenCtx struct {
	table        *schema.Table
	tableRef     *schema.TableRef
	colMap       map[schema.ColumnRef]uint
	instructions []runtime.Instruction
	emptyReg     uint
}

func newCodegenCtx(table *schema.Table) *codegenCtx {
	// Preserve len(table.Columns) reg for column raw data.
	i := len(table.Columns)
	return &codegenCtx{
		table:    table,
		colMap:   make(map[schema.ColumnRef]uint),
		emptyReg: uint(i),
	}
}

func (c *codegenCtx) nextReg() uint {
	reg := c.emptyReg
	c.emptyReg++
	return reg
}

func (c *codegenCtx) setEmptyReg(r uint) {
	c.emptyReg = r
}

func (c *codegenCtx) addInstruction(i runtime.Instruction) {
	c.instructions = append(c.instructions, i)
}

func (c *codegenCtx) loadColumnSet(colset planner.ColumnSet, pkReg uint) {
	for _, col := range colset {
		_, ok := c.colMap[col.Column]
		if !ok {
			c.loadColumn(col, pkReg)
		}
	}
}

func (c *codegenCtx) loadColumn(col *schema.ColumnDescriptor, pkReg uint) {
	if c.tableRef == nil || col.Table != *c.tableRef {
		err := fmt.Errorf(
			"operate on different table ctx(%v), req(%v)", c.tableRef, col.Table)
		panic(err)
	}
	tableArg := newTableRefImm(*c.tableRef)
	colArg := newColumnRefImm(col.Column)
	output := uint(col.Column)
	c.addInstruction(runtime.Instruction{
		Op:     runtime.LOAD,
		Input:  []*runtime.Operand{tableArg, newReg(pkReg), colArg},
		Output: output,
		// TODO(wmin0): position.
	})
	c.colMap[col.Column] = output
}

func (c *codegenCtx) clearColumn() {
	c.colMap = make(map[schema.ColumnRef]uint)
}

func (c *codegenCtx) genASTIdentifierNode(
	n *ast.IdentifierNode,
	pkReg *uint,
	colMap map[schema.ColumnRef]uint) *runtime.Operand {
	// TODO(wmin0): correct it.
	return newReg(0)
}

func (c *codegenCtx) genASTValuer(
	n ast.Valuer,
	pkReg *uint,
	colMap map[schema.ColumnRef]uint) *runtime.Operand {
	switch node := n.(type) {
	case *ast.BoolValueNode:
		switch node.V {
		case ast.BoolValueTrue:
			return newImmColDecimal(node.GetType(), dec.True)
		case ast.BoolValueFalse:
			return newImmColDecimal(node.GetType(), dec.False)
		default:
			// case ast.BoolValueUnknown:
			panic("unsupport BoolValueUnknown")
		}
	case *ast.IntegerValueNode:
		return newImmColDecimal(node.GetType(), node.V)
	case *ast.DecimalValueNode:
		return newImmColDecimal(node.GetType(), node.V)
	case *ast.BytesValueNode:
		return newImmColBytes(node.GetType(), node.V)
	case *ast.AnyValueNode:
		return newImmColBytes(node.GetType(), nil)
	case *ast.NullValueNode:
		panic("unsupport NullValueNode")
	default:
		panic(fmt.Sprintf("unknown ast valuer %s", reflect.TypeOf(n)))
	}
}

func (c *codegenCtx) genASTUnaryOperator(
	n ast.UnaryOperator,
	pkReg *uint,
	colMap map[schema.ColumnRef]uint) *runtime.Operand {
	target := c.genASTExprNode(n.GetTarget(), pkReg, colMap)
	output := c.nextReg()

	switch node := n.(type) {
	case *ast.NegOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.NEG,
			Input:    []*runtime.Operand{target},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.NotOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.NOT,
			Input:    []*runtime.Operand{target},
			Output:   output,
			Position: n.GetPosition(),
		})
	default:
		panic(fmt.Sprintf("unknown ast unary operator %s", reflect.TypeOf(n)))
	}

	c.setEmptyReg(output + 1)
	return newReg(output)
}

func (c *codegenCtx) genASTBinaryOperator(
	n ast.BinaryOperator,
	pkReg *uint,
	colMap map[schema.ColumnRef]uint) *runtime.Operand {
	output := c.nextReg()
	object := c.genASTExprNode(n.GetObject(), pkReg, colMap)
	subject := c.genASTExprNode(n.GetSubject(), pkReg, colMap)

	switch node := n.(type) {
	case *ast.AndOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.AND,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.OrOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.OR,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.GreaterOrEqualOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.GT,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
		output1 := c.nextReg()
		c.addInstruction(runtime.Instruction{
			Op:       runtime.EQ,
			Input:    []*runtime.Operand{object, subject},
			Output:   output1,
			Position: n.GetPosition(),
		})
		c.addInstruction(runtime.Instruction{
			Op:       runtime.OR,
			Input:    []*runtime.Operand{newReg(output), newReg(output1)},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.LessOrEqualOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.LT,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
		output1 := c.nextReg()
		c.addInstruction(runtime.Instruction{
			Op:       runtime.EQ,
			Input:    []*runtime.Operand{object, subject},
			Output:   output1,
			Position: n.GetPosition(),
		})
		c.addInstruction(runtime.Instruction{
			Op:       runtime.OR,
			Input:    []*runtime.Operand{newReg(output), newReg(output1)},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.NotEqualOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.EQ,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
		c.addInstruction(runtime.Instruction{
			Op:       runtime.NOT,
			Input:    []*runtime.Operand{newReg(output)},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.EqualOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.EQ,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.GreaterOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.GT,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.LessOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.LT,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.ConcatOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.CONCAT,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.AddOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.ADD,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.SubOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.SUB,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.MulOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.MUL,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.DivOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.DIV,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.ModOperatorNode:
		c.addInstruction(runtime.Instruction{
			Op:       runtime.MOD,
			Input:    []*runtime.Operand{object, subject},
			Output:   output,
			Position: n.GetPosition(),
		})
	case *ast.IsOperatorNode:
		panic("not support IsOperatorNode")
	case *ast.LikeOperatorNode:
		args := make([]*runtime.Operand, 0, 3)
		args = append(args, object, subject)
		if node.Escape != nil {
			args = append(args, c.genASTExprNode(node.Escape, pkReg, colMap))
		}
		c.addInstruction(runtime.Instruction{
			Op:       runtime.LIKE,
			Input:    args,
			Output:   output,
			Position: n.GetPosition(),
		})
	default:
		panic(fmt.Sprintf("unknown ast binary operator %s", reflect.TypeOf(n)))
	}

	c.setEmptyReg(output + 1)
	return newReg(output)
}

func (c *codegenCtx) genASTCastOperatorNode(
	n *ast.CastOperatorNode,
	pkReg *uint,
	colMap map[schema.ColumnRef]uint) *runtime.Operand {
	output := c.nextReg()
	source := c.genASTExprNode(n.SourceExpr, pkReg, colMap)
	target := newImmType(n.GetType())

	c.addInstruction(runtime.Instruction{
		Op:       runtime.CAST,
		Input:    []*runtime.Operand{source, target},
		Output:   output,
		Position: n.GetPosition(),
	})

	c.setEmptyReg(output + 1)
	return newReg(output)
}

func (c *codegenCtx) genASTInOperatorNode(
	n *ast.InOperatorNode,
	pkReg *uint,
	colMap map[schema.ColumnRef]uint) *runtime.Operand {
	output := c.nextReg()
	left := c.genASTExprNode(n.Left, pkReg, colMap)
	right := make([]*runtime.Operand, 0, len(n.Right))
	output1 := c.nextReg()

	for idx, r := range n.Right {
		c.setEmptyReg(output1 + 1)
		right = append(right, c.genASTExprNode(r, pkReg, colMap))
		if idx == 0 {
			c.addInstruction(runtime.Instruction{
				Op:       runtime.EQ,
				Input:    []*runtime.Operand{left, right[idx]},
				Output:   output,
				Position: n.GetPosition(),
			})
		} else {
			c.addInstruction(runtime.Instruction{
				Op:       runtime.EQ,
				Input:    []*runtime.Operand{left, right[idx]},
				Output:   output1,
				Position: n.GetPosition(),
			})
			c.addInstruction(runtime.Instruction{
				Op:       runtime.OR,
				Input:    []*runtime.Operand{newReg(output), newReg(output1)},
				Output:   output,
				Position: n.GetPosition(),
			})
		}
	}

	c.setEmptyReg(output + 1)
	return newReg(output)
}

func (c *codegenCtx) genASTFunctionOperatorNode(
	n *ast.FunctionOperatorNode,
	pkReg *uint,
	colMap map[schema.ColumnRef]uint) *runtime.Operand {
	output := c.nextReg()
	args := make([]*runtime.Operand, 0, len(n.Args)+2)

	if pkReg != nil {
		args = append(args, newReg(*pkReg))
	} else {
		pkArg := newImmColDecimal(
			ast.ComposeDataType(ast.DataTypeMajorUint, 0),
			decimal.New(1, 0),
		)
		args = append(args, pkArg)
	}

	funcArg := newImmColBytes(
		ast.ComposeDataType(
			ast.DataTypeMajorDynamicBytes, ast.DataTypeMinorDontCare),
		n.Name.Name,
	)
	args = append(args, funcArg)

	for _, arg := range n.Args {
		args = append(args, c.genASTExprNode(arg, pkReg, colMap))
	}

	c.addInstruction(runtime.Instruction{
		Op:       runtime.SOLFUNC,
		Input:    args,
		Output:   output,
		Position: n.GetPosition(),
	})

	c.setEmptyReg(output + 1)
	return newReg(output)
}

func (c *codegenCtx) genASTExprNode(
	n ast.ExprNode,
	pkReg *uint,
	colMap map[schema.ColumnRef]uint) *runtime.Operand {
	switch node := n.(type) {
	case *ast.IdentifierNode:
		return c.genASTIdentifierNode(node, pkReg, colMap)
	case ast.Valuer:
		return c.genASTValuer(node, pkReg, colMap)
	case ast.UnaryOperator:
		return c.genASTUnaryOperator(node, pkReg, colMap)
	case ast.BinaryOperator:
		return c.genASTBinaryOperator(node, pkReg, colMap)
	case *ast.CastOperatorNode:
		return c.genASTCastOperatorNode(node, pkReg, colMap)
	case *ast.InOperatorNode:
		return c.genASTInOperatorNode(node, pkReg, colMap)
	case *ast.FunctionOperatorNode:
		return c.genASTFunctionOperatorNode(node, pkReg, colMap)
	default:
		panic(fmt.Sprintf("unknown ast type %s", reflect.TypeOf(n)))
	}
}

func (c *codegenCtx) genUnionStep(
	step *planner.UnionStep, inputs []*runtime.Operand, output uint) {
	c.addInstruction(runtime.Instruction{
		Op:     runtime.UNION,
		Input:  inputs,
		Output: output,
		// TODO(wmin0): position.
	})
}

func (c *codegenCtx) genIntersectStep(
	step *planner.IntersectStep, inputs []*runtime.Operand, output uint) {
	c.addInstruction(runtime.Instruction{
		Op:     runtime.INTXN,
		Input:  inputs,
		Output: output,
		// TODO(wmin0): position.
	})
}

func (c *codegenCtx) filterColWithPk(
	colset planner.ColumnSet,
	pkReg uint,
	output uint,
	cond ast.ExprNode) {
	c.loadColumnSet(colset, pkReg)

	result := c.genASTExprNode(cond, &pkReg, c.colMap)
	c.addInstruction(runtime.Instruction{
		Op:     runtime.FILTER,
		Input:  []*runtime.Operand{newReg(pkReg), result},
		Output: output,
		// TODO(wmin0): position.
	})

	// Maintain colMap data align with pk.
	for col, reg := range c.colMap {
		c.addInstruction(runtime.Instruction{
			Op:     runtime.FILTER,
			Input:  []*runtime.Operand{newReg(reg), result},
			Output: reg,
			// TODO(wmin0): position.
		})
	}
}

func (c *codegenCtx) genFilterStep(
	step *planner.FilterStep, inputs []*runtime.Operand, output uint) {
	c.filterColWithPk(
		step.ColumnSet,
		inputs[0].RegisterIndex,
		output,
		step.Condition,
	)
}

func (c *codegenCtx) genScanIndices(
	step *planner.ScanIndices, inputs []*runtime.Operand, output uint) {
	if c.tableRef == nil || step.Table != *c.tableRef {
		err := fmt.Errorf(
			"operate on different table ctx(%v), req(%v)", c.tableRef, step.Table)
		panic(err)
	}

	rowOp := make([]*runtime.Operand, 0, len(step.Values))
	for _, row := range step.Values {
		colOp := make([]*runtime.Operand, 0, len(row))
		for _, v := range row {
			colOp = append(colOp, c.genASTValuer(v, nil, nil))
		}
		rowOp = append(rowOp, expandImmCol(colOp))
	}
	condArg := expandImmRow(rowOp)
	idxArg := newIndexRefImm(*c.tableRef, step.Index)
	c.addInstruction(runtime.Instruction{
		Op:     runtime.REPEATIDX,
		Input:  []*runtime.Operand{idxArg, condArg},
		Output: output,
		// TODO(wmin0): position.
	})
}

func (c *codegenCtx) genScanIndexValues(
	step *planner.ScanIndexValues, inputs []*runtime.Operand, output uint) {
	if c.tableRef == nil || step.Table != *c.tableRef {
		err := fmt.Errorf(
			"operate on different table ctx(%v), req(%v)", c.tableRef, step.Table)
		panic(err)
	}

	valueReg := c.nextReg()
	idxArg := newIndexRefImm(*c.tableRef, step.Index)
	c.addInstruction(runtime.Instruction{
		Op:     runtime.REPEATIDXV,
		Input:  []*runtime.Operand{idxArg},
		Output: valueReg,
		// TODO(wmin0): position.
	})

	index := c.table.Indices[step.Index]
	colMap := make(map[schema.ColumnRef]uint, len(index.Columns))
	for idx, col := range index.Columns {
		colReg := c.nextReg()
		colMap[col] = colReg
		fieldArg := newColumnRefImm(schema.ColumnRef(idx))
		c.addInstruction(runtime.Instruction{
			Op:     runtime.FIELD,
			Input:  []*runtime.Operand{newReg(valueReg), fieldArg},
			Output: colReg,
			// TODO(wmin0): position.
		})
	}

	result := c.genASTExprNode(step.Condition, nil, colMap)
	c.addInstruction(runtime.Instruction{
		Op:     runtime.FILTER,
		Input:  []*runtime.Operand{newReg(valueReg), result},
		Output: valueReg,
		// TODO(wmin0): position.
	})
	c.addInstruction(runtime.Instruction{
		Op:     runtime.REPEATIDX,
		Input:  []*runtime.Operand{idxArg, newReg(valueReg)},
		Output: output,
		// TODO(wmin0): position.
	})
}

func (c *codegenCtx) genScanTable(
	step *planner.ScanTable, inputs []*runtime.Operand, output uint) {
	if c.tableRef == nil || step.Table != *c.tableRef {
		err := fmt.Errorf(
			"operate on different table ctx(%v), req(%v)", c.tableRef, step.Table)
		panic(err)
	}

	pkReg := c.nextReg()
	tableArg := newTableRefImm(*c.tableRef)
	c.addInstruction(runtime.Instruction{
		Op:     runtime.REPEATPK,
		Input:  []*runtime.Operand{tableArg},
		Output: output,
		// TODO(wmin0): position.
	})

	c.filterColWithPk(
		step.ColumnSet,
		pkReg,
		output,
		step.Condition,
	)
}

func (c *codegenCtx) genInsertStep(
	step *planner.InsertStep, inputs []*runtime.Operand, output uint) {
	if c.tableRef == nil || step.Table != *c.tableRef {
		err := fmt.Errorf(
			"operate on different table ctx(%v), req(%v)", c.tableRef, step.Table)
		panic(err)
	}

	tableArg := newTableRefImm(*c.tableRef)
	colArg := newColumnRefsImm(step.Columns)

	for _, row := range step.Values {
		c.setEmptyReg(output + 1)
		args := make([]*runtime.Operand, 0, 2+len(row))
		args = append(args, tableArg, colArg)
		for _, v := range row {
			args = append(args, c.genASTExprNode(v, nil, nil))
		}
		c.addInstruction(runtime.Instruction{
			Op:     runtime.INSERT,
			Input:  args,
			Output: output,
			// TODO(wmin0): position.
		})
	}
	result := newImmColDecimal(
		ast.ComposeDataType(ast.DataTypeMajorUint, 7),
		decimal.New(int64(len(step.Values)), 0),
	)
	c.addInstruction(runtime.Instruction{
		Op:     runtime.ZIP,
		Input:  []*runtime.Operand{result},
		Output: output,
		// TODO(wmin0): position.
	})
}

func (c *codegenCtx) maybeRangeOnOutput(
	output uint, offset *decimal.Decimal, limit *decimal.Decimal) {
	if offset == nil && limit == nil {
		return
	}

	args := make([]*runtime.Operand, 0, 3)
	args = append(args, newReg(output))
	// Must have offset.
	if offset != nil {
		offsetArg := newImmColDecimal(
			ast.ComposeDataType(ast.DataTypeMajorUint, 7),
			*offset,
		)
		args = append(args, offsetArg)
	} else {
		offsetArg := newImmColDecimal(
			ast.ComposeDataType(ast.DataTypeMajorUint, 7),
			decimal.Zero,
		)
		args = append(args, offsetArg)
	}
	if limit != nil {
		limitArg := newImmColDecimal(
			ast.ComposeDataType(ast.DataTypeMajorUint, 7),
			*limit,
		)
		args = append(args, limitArg)
	}
	c.addInstruction(runtime.Instruction{
		Op:     runtime.RANGE,
		Input:  args,
		Output: output,
		// TODO(wmin0): position.
	})
}

func (c *codegenCtx) genSelectStep(
	step *planner.SelectStep, inputs []*runtime.Operand, output uint) {
	pk := inputs[0]
	c.loadColumnSet(step.ColumnSet, pk.RegisterIndex)

	cols := make([]*runtime.Operand, 0, len(step.Columns)+len(step.Order))
	for _, col := range step.Columns {
		cols = append(cols, c.genASTExprNode(col, &pk.RegisterIndex, c.colMap))
	}

	orders := make([]*runtime.Operand, 0, len(step.Order))
	for _, o := range step.Order {
		field := 0
		i, ok := o.Expr.(*ast.IdentifierNode)
		if ok {
			// TODO(wmin0): correct it.
			field = 0
		} else {
			field = len(cols)
			cols = append(
				cols, c.genASTExprNode(o.Expr, &pk.RegisterIndex, c.colMap))
		}
		orders = append(orders, newOrderImm(field, o.Desc))
	}
	c.addInstruction(runtime.Instruction{
		Op:     runtime.ZIP,
		Input:  cols,
		Output: output,
		// TODO(wmin0): position.
	})

	if len(orders) != 0 {
		c.addInstruction(runtime.Instruction{
			Op:     runtime.SORT,
			Input:  []*runtime.Operand{newReg(output), expandImmRow(orders)},
			Output: output,
			// TODO(wmin0): position.
		})
	}

	// XXX(wmin0): unsupport null first option.
	if len(cols) > len(step.Columns) {
		cutArg := newImmColDecimal(
			ast.ComposeDataType(ast.DataTypeMajorUint, 1),
			decimal.New(int64(len(step.Columns)), 0),
		)
		c.addInstruction(runtime.Instruction{
			Op:     runtime.CUT,
			Input:  []*runtime.Operand{newReg(output), cutArg},
			Output: output,
			// TODO(wmin0): position.
		})
	}

	c.maybeRangeOnOutput(output, step.Offset, step.Limit)
}

func (c *codegenCtx) genSelectWithoutTable(
	step *planner.SelectWithoutTable, inputs []*runtime.Operand, output uint) {
	cols := make([]*runtime.Operand, 0, len(step.Columns))
	for _, col := range step.Columns {
		cols = append(cols, c.genASTExprNode(col, nil, nil))
	}
	c.addInstruction(runtime.Instruction{
		Op:     runtime.ZIP,
		Input:  cols,
		Output: output,
		// TODO(wmin0): position.
	})

	if step.Condition != nil {
		result := c.genASTExprNode(step.Condition, nil, nil)
		c.addInstruction(runtime.Instruction{
			Op:     runtime.FILTER,
			Input:  []*runtime.Operand{newReg(output), result},
			Output: output,
			// TODO(wmin0): position.
		})
	}

	c.maybeRangeOnOutput(output, step.Offset, step.Limit)
}

func (c *codegenCtx) genUpdateStep(
	step *planner.UpdateStep, inputs []*runtime.Operand, output uint) {
	if c.tableRef == nil || step.Table != *c.tableRef {
		err := fmt.Errorf(
			"operate on different table ctx(%v), req(%v)", c.tableRef, step.Table)
		panic(err)
	}

	pk := inputs[0]
	c.loadColumnSet(step.ColumnSet, pk.RegisterIndex)
	tableArg := newTableRefImm(*c.tableRef)
	colArg := newColumnRefsImm(step.Columns)

	args := make([]*runtime.Operand, 0, 2+len(step.Columns))
	args = append(args, tableArg, colArg)
	for _, v := range step.Values {
		args = append(args, c.genASTExprNode(v, &pk.RegisterIndex, c.colMap))
	}
	c.addInstruction(runtime.Instruction{
		Op:     runtime.UPDATE,
		Input:  args,
		Output: output,
		// TODO(wmin0): position.
	})
}

func (c *codegenCtx) genDeleteStep(
	step *planner.DeleteStep, inputs []*runtime.Operand, output uint) {
	c.addInstruction(runtime.Instruction{
		Op:     runtime.DELETE,
		Input:  inputs,
		Output: output,
		// TODO(wmin0): position.
	})
}

func (c *codegenCtx) genOperands(
	step planner.PlanStep, clear bool) []*runtime.Operand {
	operands := step.GetOperands()
	inputs := make([]*runtime.Operand, 0, len(operands))
	for _, op := range operands {
		inputs = append(inputs, c.codegen(op))
		if clear {
			c.clearColumn()
		}
	}
	return inputs
}

func (c *codegenCtx) codegen(s interface{}) *runtime.Operand {
	output := c.nextReg()
	switch step := s.(type) {
	case *planner.UnionStep:
		inputs := c.genOperands(step, true)
		c.genUnionStep(step, inputs, output)
	case *planner.IntersectStep:
		inputs := c.genOperands(step, true)
		c.genIntersectStep(step, inputs, output)
	case *planner.FilterStep:
		inputs := c.genOperands(step, false)
		c.genFilterStep(step, inputs, output)
	case *planner.ScanIndices:
		inputs := c.genOperands(step, false)
		c.genScanIndices(step, inputs, output)
	case *planner.ScanIndexValues:
		inputs := c.genOperands(step, false)
		c.genScanIndexValues(step, inputs, output)
	case *planner.ScanTable:
		inputs := c.genOperands(step, false)
		c.genScanTable(step, inputs, output)
	default:
		panic(fmt.Sprintf("unknown step type %s", reflect.TypeOf(s)))
	}
	c.setEmptyReg(output + 1)
	return newReg(output)
}

func (c *codegenCtx) codegenRoot(s interface{}) *runtime.Operand {
	output := c.nextReg()
	switch step := s.(type) {
	case *planner.InsertStep:
		c.tableRef = &step.Table
		inputs := c.genOperands(step, false)
		c.genInsertStep(step, inputs, output)
	case *planner.SelectStep:
		c.tableRef = &step.Table
		inputs := c.genOperands(step, false)
		c.genSelectStep(step, inputs, output)
	case *planner.SelectWithoutTable:
		inputs := c.genOperands(step, false)
		c.genSelectWithoutTable(step, inputs, output)
	case *planner.UpdateStep:
		c.tableRef = &step.Table
		inputs := c.genOperands(step, false)
		c.genUpdateStep(step, inputs, output)
	case *planner.DeleteStep:
		c.tableRef = &step.Table
		inputs := c.genOperands(step, false)
		c.genDeleteStep(step, inputs, output)
	default:
		panic(fmt.Sprintf("unknown root step type %s", reflect.TypeOf(s)))
	}
	c.setEmptyReg(output + 1)
	return newReg(output)
}

func Codegen(
	table *schema.Table,
	root planner.PlanStep) []runtime.Instruction {
	ctx := newCodegenCtx(table)
	ctx.codegenRoot(root)
	return ctx.instructions
}
