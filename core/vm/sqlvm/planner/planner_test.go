package planner

import (
	"fmt"
	"testing"

	"github.com/dexon-foundation/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

type PlannerTestSuite struct{ suite.Suite }

// utility functions.
// cleanPlanCost clean up cost and out num recursively.
func cleanPlanCost(plan PlanStep) {
	plan.SetCost(0)
	plan.SetOutNum(0)
	for _, p := range plan.GetOperands() {
		cleanPlanCost(p)
	}
}

func (s *PlannerTestSuite) assertPlanEqual(
	expected, actual PlanStep,
	compareCost bool,
) {
	if !compareCost {
		cleanPlanCost(expected)
		cleanPlanCost(actual)
	}
	s.Require().Equal(expected, actual)
}

// createTable with given name, column names, column ids in every indices.
func (s *PlannerTestSuite) createTable(
	name []byte,
	columns [][]byte,
	indices [][]schema.ColumnRef,
) schema.Table {
	table := schema.Table{
		Name: name,
	}
	table.Columns = make([]schema.Column, len(columns))
	for i, cname := range columns {
		col := schema.Column{}
		col.Name = cname
		table.Columns[i] = col
	}
	table.Indices = make([]schema.Index, len(indices))
	for i, cols := range indices {
		table.Indices[i] = schema.Index{
			Name:    []byte(fmt.Sprintf("index%d", i)),
			Columns: cols,
		}
	}
	return table
}

func (s *PlannerTestSuite) TestBasic() {
	var tables schema.Schema = []schema.Table{
		s.createTable(
			[]byte("table1"),
			[][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
			},
			[][]schema.ColumnRef{
				[]schema.ColumnRef{0},
				[]schema.ColumnRef{1},
				[]schema.ColumnRef{0, 1, 2},
			},
		),
	}
	{
		stmt := &ast.InsertStmtNode{
			Table: &ast.IdentifierNode{
				Name: []byte("table1"),
				Desc: &schema.TableDescriptor{Table: 0},
			},
			Insert: &ast.InsertWithColumnOptionNode{
				Column: []*ast.IdentifierNode{
					&ast.IdentifierNode{
						Name: []byte("a"),
						Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
					},
					&ast.IdentifierNode{
						Name: []byte("b"),
						Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
					},
				},
				Value: [][]ast.ExprNode{
					[]ast.ExprNode{
						&ast.BytesValueNode{V: []byte("valA")},
						&ast.BytesValueNode{V: []byte("valB")},
					},
				},
			},
		}
		expectedPlan := &InsertStep{
			Table:   0,
			Columns: []schema.ColumnRef{0, 1},
			Values: [][]ast.ExprNode{
				[]ast.ExprNode{
					&ast.BytesValueNode{V: []byte("valA")},
					&ast.BytesValueNode{V: []byte("valB")},
				},
			},
		}
		plan, err := Plan(nil, tables, stmt)
		s.Require().Nil(err)
		s.Require().NotNil(plan)
		s.assertPlanEqual(expectedPlan, plan, false)
	}
	{
		expr := &ast.AddOperatorNode{
			BinaryOperatorNode: ast.BinaryOperatorNode{
				Object: &ast.IdentifierNode{
					Name: []byte("b"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
				},
				Subject: &ast.IdentifierNode{
					Name: []byte("b"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
				},
			},
		}
		stmt := &ast.UpdateStmtNode{
			Table: &ast.IdentifierNode{
				Name: []byte("table1"),
				Desc: &schema.TableDescriptor{Table: 0},
			},
			Assignment: []*ast.AssignOperatorNode{
				&ast.AssignOperatorNode{
					Column: &ast.IdentifierNode{
						Name: []byte("a"),
						Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
					},
					Expr: expr,
				},
			},
		}
		expectedPlan := &UpdateStep{
			Table: 0,
			ColumnSet: ColumnSet{
				&schema.ColumnDescriptor{Table: 0, Column: 1},
			},
			Columns: []schema.ColumnRef{0},
			Values:  []ast.ExprNode{expr},
			// Children.
			PlanStepBase: PlanStepBase{
				Operands: []PlanStep{
					&ScanTable{
						Table:     0,
						Condition: nil,
					},
				},
			},
		}
		plan, err := Plan(nil, tables, stmt)
		s.Require().Nil(err)
		s.Require().NotNil(plan)
		s.assertPlanEqual(expectedPlan, plan, false)
	}
	{
		// Test a index scan select.
		// TODO(yenlin): test more on planning.
		limit := decimal.New(10, 0)
		offset := decimal.New(20, 0)
		stmt := &ast.SelectStmtNode{
			Column: []ast.ExprNode{
				&ast.IdentifierNode{
					Name: []byte("a"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
				},
				&ast.IdentifierNode{
					Name: []byte("b"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
				},
			},
			Table: &ast.IdentifierNode{
				Name: []byte("table1"),
				Desc: &schema.TableDescriptor{Table: 0},
			},
			Where: &ast.WhereOptionNode{
				Condition: &ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("b"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
						},
						Subject: &ast.BytesValueNode{V: []byte("valB")},
					},
				},
			},
			Order: []*ast.OrderOptionNode{
				&ast.OrderOptionNode{
					Expr: &ast.IdentifierNode{
						Name: []byte("c"),
						Desc: &schema.ColumnDescriptor{Table: 0, Column: 2},
					},
				},
			},
			Limit: &ast.LimitOptionNode{
				Value: &ast.IntegerValueNode{
					V: limit,
				},
			},
			Offset: &ast.OffsetOptionNode{
				Value: &ast.IntegerValueNode{
					V: offset,
				},
			},
		}
		expectedPlan := &SelectStep{
			Table: 0,
			ColumnSet: ColumnSet{
				&schema.ColumnDescriptor{Table: 0, Column: 0},
				&schema.ColumnDescriptor{Table: 0, Column: 1},
				&schema.ColumnDescriptor{Table: 0, Column: 2},
			},
			Columns: []ast.ExprNode{
				&ast.IdentifierNode{
					Name: []byte("a"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
				},
				&ast.IdentifierNode{
					Name: []byte("b"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
				},
			},
			Order: []*ast.OrderOptionNode{
				&ast.OrderOptionNode{
					Expr: &ast.IdentifierNode{
						Name: []byte("c"),
						Desc: &schema.ColumnDescriptor{Table: 0, Column: 2},
					},
				},
			},
			Limit:  &limit,
			Offset: &offset,

			// Children.
			PlanStepBase: PlanStepBase{
				Operands: []PlanStep{
					&ScanIndices{
						Table: 0,
						Index: 1,
						Values: [][]ast.Valuer{
							[]ast.Valuer{
								&ast.BytesValueNode{V: []byte("valB")},
							},
						},
					},
				},
			},
		}
		plan, err := Plan(nil, tables, stmt)
		s.Require().Nil(err)
		s.Require().NotNil(plan)
		s.assertPlanEqual(expectedPlan, plan, false)
	}
	{
		// Test select without table.
		limit := decimal.New(10, 0)
		offset := decimal.New(20, 0)
		stmt := &ast.SelectStmtNode{
			Column: []ast.ExprNode{
				&ast.FunctionOperatorNode{
					Name: &ast.IdentifierNode{
						Name: []byte("rand"),
					},
				},
				&ast.BytesValueNode{V: []byte("valB")},
			},
			Where: &ast.WhereOptionNode{
				Condition: &ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.FunctionOperatorNode{
							Name: &ast.IdentifierNode{
								Name: []byte("rand"),
							},
						},
						Subject: &ast.DecimalValueNode{V: decimal.New(1, 0)},
					},
				},
			},
			Limit: &ast.LimitOptionNode{
				Value: &ast.IntegerValueNode{
					V: limit,
				},
			},
			Offset: &ast.OffsetOptionNode{
				Value: &ast.IntegerValueNode{
					V: offset,
				},
			},
		}
		expectedPlan := &SelectWithoutTable{
			Columns: []ast.ExprNode{
				&ast.FunctionOperatorNode{
					Name: &ast.IdentifierNode{
						Name: []byte("rand"),
					},
				},
				&ast.BytesValueNode{V: []byte("valB")},
			},
			Condition: &ast.EqualOperatorNode{
				BinaryOperatorNode: ast.BinaryOperatorNode{
					Object: &ast.FunctionOperatorNode{
						Name: &ast.IdentifierNode{
							Name: []byte("rand"),
						},
					},
					Subject: &ast.DecimalValueNode{V: decimal.New(1, 0)},
				},
			},
			Limit:  &limit,
			Offset: &offset,
		}
		plan, err := Plan(nil, tables, stmt)
		s.Require().Nil(err)
		s.Require().NotNil(plan)
		s.assertPlanEqual(expectedPlan, plan, false)
	}
	{
		stmt := &ast.DeleteStmtNode{
			Table: &ast.IdentifierNode{
				Name: []byte("table1"),
				Desc: &schema.TableDescriptor{Table: 0},
			},
			Where: &ast.WhereOptionNode{
				Condition: &ast.GreaterOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("b"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
						},
						Subject: &ast.BytesValueNode{V: []byte("valB")},
					},
				},
			},
		}
		expectedPlan := &DeleteStep{
			Table: 0,
			// Children.
			PlanStepBase: PlanStepBase{
				Operands: []PlanStep{
					&ScanIndexValues{
						Table: 0,
						Index: 1,
						Condition: &ast.GreaterOperatorNode{
							BinaryOperatorNode: ast.BinaryOperatorNode{
								Object: &ast.IdentifierNode{
									Name: []byte("b"),
									Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
								},
								Subject: &ast.BytesValueNode{V: []byte("valB")},
							},
						},
					},
				},
			},
		}
		plan, err := Plan(nil, tables, stmt)
		s.Require().Nil(err)
		s.Require().NotNil(plan)
		s.assertPlanEqual(expectedPlan, plan, false)
	}
}

func (s *PlannerTestSuite) TestHashKeys() {
	var tableSchema schema.Schema = []schema.Table{
		s.createTable(
			[]byte("table1"),
			[][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
			},
			[][]schema.ColumnRef{
				[]schema.ColumnRef{0},
				[]schema.ColumnRef{1},
				[]schema.ColumnRef{0, 1, 2},
			},
		),
	}
	planner := planner{
		schema: tableSchema,
	}
	// Pre-define possible values in every columns.
	consts := [][]ast.Valuer{
		[]ast.Valuer{
			&ast.BytesValueNode{V: []byte("val1")},
			&ast.BytesValueNode{V: []byte("val2")},
			&ast.BytesValueNode{V: []byte("val3")},
		},
		[]ast.Valuer{
			&ast.BytesValueNode{V: []byte("val4")},
			&ast.BytesValueNode{V: []byte("val5")},
			&ast.BytesValueNode{V: []byte("val6")},
		},
		[]ast.Valuer{
			&ast.BytesValueNode{V: []byte("val7")},
			&ast.BytesValueNode{V: []byte("val8")},
			&ast.BytesValueNode{V: []byte("val9")},
		},
	}
	genCombination := func(columns []schema.ColumnRef, values []int) []ast.Valuer {
		s.Require().Equal(len(columns), len(values))
		s.Require().True(len(columns) <= len(consts))

		ret := make([]ast.Valuer, len(columns))
		for i := range ret {
			colIdx := columns[i]
			valIdx := values[i]
			s.Require().True(int(colIdx) < len(consts))
			s.Require().True(valIdx < len(consts[colIdx]))
			ret[i] = consts[colIdx][valIdx]
		}

		return ret
	}

	{
		// Test AND merge.
		cl := &clause{
			ColumnSet: ColumnSet{
				&schema.ColumnDescriptor{Table: 0, Column: 0},
				&schema.ColumnDescriptor{Table: 0, Column: 1},
				&schema.ColumnDescriptor{Table: 0, Column: 2},
			},
			Attr: clauseAttrAnd | clauseAttrEnumerable,
			SubCls: []*clause{
				&clause{
					ColumnSet: ColumnSet{
						&schema.ColumnDescriptor{Table: 0, Column: 0},
						&schema.ColumnDescriptor{Table: 0, Column: 2},
					},
					Attr: clauseAttrEnumerable,
					HashKeys: [][]ast.Valuer{
						genCombination([]schema.ColumnRef{0, 2}, []int{0, 0}),
						genCombination([]schema.ColumnRef{0, 2}, []int{1, 1}),
						genCombination([]schema.ColumnRef{0, 2}, []int{1, 2}),
					},
				},
				&clause{
					ColumnSet: ColumnSet{
						&schema.ColumnDescriptor{Table: 0, Column: 1},
					},
					Attr: clauseAttrEnumerable,
					HashKeys: [][]ast.Valuer{
						genCombination([]schema.ColumnRef{1}, []int{0}),
						genCombination([]schema.ColumnRef{1}, []int{1}),
					},
				},
			},
		}
		planner.mergeAndHashKeys(cl)
		expectedKeys := [][]ast.Valuer{
			genCombination([]schema.ColumnRef{0, 1, 2}, []int{0, 0, 0}),
			genCombination([]schema.ColumnRef{0, 1, 2}, []int{1, 0, 1}),
			genCombination([]schema.ColumnRef{0, 1, 2}, []int{1, 0, 2}),
			genCombination([]schema.ColumnRef{0, 1, 2}, []int{0, 1, 0}),
			genCombination([]schema.ColumnRef{0, 1, 2}, []int{1, 1, 1}),
			genCombination([]schema.ColumnRef{0, 1, 2}, []int{1, 1, 2}),
		}
		s.Require().Equal(expectedKeys, cl.HashKeys)
	}
	{
		// Test AND merge boundary with one size empty.
		cl := &clause{
			ColumnSet: ColumnSet{
				&schema.ColumnDescriptor{Table: 0, Column: 0},
				&schema.ColumnDescriptor{Table: 0, Column: 1},
				&schema.ColumnDescriptor{Table: 0, Column: 2},
			},
			Attr: clauseAttrAnd | clauseAttrEnumerable,
			SubCls: []*clause{
				&clause{
					ColumnSet: ColumnSet{
						&schema.ColumnDescriptor{Table: 0, Column: 0},
						&schema.ColumnDescriptor{Table: 0, Column: 2},
					},
					Attr:     clauseAttrEnumerable,
					HashKeys: [][]ast.Valuer{},
				},
				&clause{
					ColumnSet: ColumnSet{
						&schema.ColumnDescriptor{Table: 0, Column: 1},
					},
					Attr: clauseAttrEnumerable,
					HashKeys: [][]ast.Valuer{
						genCombination([]schema.ColumnRef{1}, []int{0}),
						genCombination([]schema.ColumnRef{1}, []int{1}),
					},
				},
			},
		}
		planner.mergeAndHashKeys(cl)
		expectedKeys := [][]ast.Valuer{}
		s.Require().Equal(expectedKeys, cl.HashKeys)
		cl = &clause{
			ColumnSet: ColumnSet{
				&schema.ColumnDescriptor{Table: 0, Column: 0},
				&schema.ColumnDescriptor{Table: 0, Column: 1},
				&schema.ColumnDescriptor{Table: 0, Column: 2},
			},
			Attr: clauseAttrAnd | clauseAttrEnumerable,
			SubCls: []*clause{
				&clause{
					ColumnSet: ColumnSet{
						&schema.ColumnDescriptor{Table: 0, Column: 0},
						&schema.ColumnDescriptor{Table: 0, Column: 2},
					},
					Attr: clauseAttrEnumerable,
					HashKeys: [][]ast.Valuer{
						genCombination([]schema.ColumnRef{0, 2}, []int{0, 0}),
						genCombination([]schema.ColumnRef{0, 2}, []int{1, 1}),
						genCombination([]schema.ColumnRef{0, 2}, []int{1, 2}),
					},
				},
				&clause{
					ColumnSet: ColumnSet{
						&schema.ColumnDescriptor{Table: 0, Column: 1},
					},
					Attr:     clauseAttrEnumerable,
					HashKeys: [][]ast.Valuer{},
				},
			},
		}
		planner.mergeAndHashKeys(cl)
		expectedKeys = [][]ast.Valuer{}
		s.Require().Equal(expectedKeys, cl.HashKeys)
	}
	{
		// Test OR merge.
		cl := &clause{
			ColumnSet: ColumnSet{
				&schema.ColumnDescriptor{Table: 0, Column: 0},
				&schema.ColumnDescriptor{Table: 0, Column: 2},
			},
			Attr: clauseAttrOr | clauseAttrEnumerable,
			SubCls: []*clause{
				&clause{
					ColumnSet: ColumnSet{
						&schema.ColumnDescriptor{Table: 0, Column: 0},
						&schema.ColumnDescriptor{Table: 0, Column: 2},
					},
					Attr: clauseAttrEnumerable,
					HashKeys: [][]ast.Valuer{
						genCombination([]schema.ColumnRef{0, 2}, []int{0, 0}),
						genCombination([]schema.ColumnRef{0, 2}, []int{1, 1}),
					},
				},
				&clause{
					ColumnSet: ColumnSet{
						&schema.ColumnDescriptor{Table: 0, Column: 0},
						&schema.ColumnDescriptor{Table: 0, Column: 2},
					},
					Attr: clauseAttrEnumerable,
					HashKeys: [][]ast.Valuer{
						genCombination([]schema.ColumnRef{0, 2}, []int{1, 2}),
						genCombination([]schema.ColumnRef{0, 2}, []int{2, 1}),
						genCombination([]schema.ColumnRef{0, 2}, []int{2, 2}),
					},
				},
			},
		}
		planner.mergeOrHashKeys(cl)
		expectedKeys := [][]ast.Valuer{
			genCombination([]schema.ColumnRef{0, 2}, []int{0, 0}),
			genCombination([]schema.ColumnRef{0, 2}, []int{1, 1}),
			genCombination([]schema.ColumnRef{0, 2}, []int{1, 2}),
			genCombination([]schema.ColumnRef{0, 2}, []int{2, 1}),
			genCombination([]schema.ColumnRef{0, 2}, []int{2, 2}),
		}
		s.Require().Equal(expectedKeys, cl.HashKeys)
	}
}

func (s *PlannerTestSuite) TestClauseAttr() {
	var dbSchema schema.Schema = []schema.Table{
		s.createTable(
			[]byte("table1"),
			[][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
			},
			[][]schema.ColumnRef{
				[]schema.ColumnRef{0},
				[]schema.ColumnRef{1},
				[]schema.ColumnRef{0, 1, 2},
			},
		),
	}
	planner := planner{
		schema:   dbSchema,
		table:    &dbSchema[0],
		tableRef: 0,
	}

	{
		node := &ast.IdentifierNode{
			Name: []byte("a"),
			Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Equal(clauseAttrColumn, cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
		}, cl.ColumnSet)
		s.Require().Zero(cl.HashKeys)
	}
	{
		node := &ast.BytesValueNode{V: []byte("valA")}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Equal(clauseAttrConst, cl.Attr)
		s.Require().Zero(cl.ColumnSet)
		s.Require().Equal([][]ast.Valuer{[]ast.Valuer{node}}, cl.HashKeys)
	}
	{
		valA := &ast.BytesValueNode{V: []byte("valA")}
		node := &ast.EqualOperatorNode{
			BinaryOperatorNode: ast.BinaryOperatorNode{
				Object: &ast.IdentifierNode{
					Name: []byte("a"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
				},
				Subject: valA,
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Equal(clauseAttrEnumerable, cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
		}, cl.ColumnSet)
		s.Require().Equal([][]ast.Valuer{[]ast.Valuer{valA}}, cl.HashKeys)
	}
	{
		node := &ast.GreaterOrEqualOperatorNode{
			BinaryOperatorNode: ast.BinaryOperatorNode{
				Object: &ast.IdentifierNode{
					Name: []byte("a"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
				},
				Subject: &ast.BytesValueNode{V: []byte("valA")},
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Zero(cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
		}, cl.ColumnSet)
		s.Require().Zero(cl.HashKeys)
	}
	{
		valA := &ast.BytesValueNode{V: []byte("valA")}
		valB := &ast.BytesValueNode{V: []byte("valB")}
		node := &ast.AndOperatorNode{
			BinaryOperatorNode: ast.BinaryOperatorNode{
				Object: &ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("b"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
						},
						Subject: valB,
					},
				},
				Subject: &ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("a"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
						},
						Subject: valA,
					},
				},
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Equal(clauseAttrEnumerable|clauseAttrAnd, cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
			&schema.ColumnDescriptor{Table: 0, Column: 1},
		}, cl.ColumnSet)
		s.Require().Equal([][]ast.Valuer{
			[]ast.Valuer{valA, valB},
		}, cl.HashKeys)
	}
	{
		node := &ast.OrOperatorNode{
			BinaryOperatorNode: ast.BinaryOperatorNode{
				Object: &ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("a"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
						},
						Subject: &ast.BytesValueNode{V: []byte("valA")},
					},
				},
				Subject: &ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("b"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
						},
						Subject: &ast.BytesValueNode{V: []byte("valB")},
					},
				},
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Equal(clauseAttrOr, cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
			&schema.ColumnDescriptor{Table: 0, Column: 1},
		}, cl.ColumnSet)
		s.Require().Zero(cl.HashKeys)
	}
	{
		valA := &ast.BytesValueNode{V: []byte("valA")}
		valB := &ast.BytesValueNode{V: []byte("valB")}
		node := &ast.OrOperatorNode{
			BinaryOperatorNode: ast.BinaryOperatorNode{
				Object: &ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("b"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
						},
						Subject: valA,
					},
				},
				Subject: &ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("b"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
						},
						Subject: valB,
					},
				},
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Equal(clauseAttrOr|clauseAttrEnumerable, cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 1},
		}, cl.ColumnSet)
		s.Require().Equal([][]ast.Valuer{
			[]ast.Valuer{valA},
			[]ast.Valuer{valB},
		}, cl.HashKeys)
	}
	{
		valA := &ast.BytesValueNode{V: []byte("valA")}
		valB := &ast.BytesValueNode{V: []byte("valB")}
		node := &ast.InOperatorNode{
			Left: &ast.EqualOperatorNode{
				BinaryOperatorNode: ast.BinaryOperatorNode{
					Object: &ast.IdentifierNode{
						Name: []byte("a"),
						Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
					},
					Subject: valA,
				},
			},
			Right: []ast.ExprNode{
				&ast.EqualOperatorNode{
					BinaryOperatorNode: ast.BinaryOperatorNode{
						Object: &ast.IdentifierNode{
							Name: []byte("b"),
							Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
						},
						Subject: valB,
					},
				},
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Zero(cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
			&schema.ColumnDescriptor{Table: 0, Column: 1},
		}, cl.ColumnSet)
		s.Require().Zero(cl.HashKeys)
	}
	{
		valA := &ast.BytesValueNode{V: []byte("valA")}
		valB := &ast.BytesValueNode{V: []byte("valB")}
		node := &ast.InOperatorNode{
			Left: &ast.IdentifierNode{
				Name: []byte("a"),
				Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
			},
			Right: []ast.ExprNode{
				valA, valB,
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Equal(clauseAttrEnumerable, cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
		}, cl.ColumnSet)
		s.Require().Equal([][]ast.Valuer{
			[]ast.Valuer{valA},
			[]ast.Valuer{valB},
		}, cl.HashKeys)
	}
	{
		node := &ast.CastOperatorNode{
			SourceExpr: &ast.AddOperatorNode{
				BinaryOperatorNode: ast.BinaryOperatorNode{
					Object: &ast.IdentifierNode{
						Name: []byte("a"),
						Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
					},
					Subject: &ast.IdentifierNode{
						Name: []byte("b"),
						Desc: &schema.ColumnDescriptor{Table: 0, Column: 1},
					},
				},
			},
			TargetType: &ast.IntTypeNode{
				Unsigned: true,
				Size:     32,
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Zero(cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
			&schema.ColumnDescriptor{Table: 0, Column: 1},
		}, cl.ColumnSet)
		s.Require().Zero(cl.HashKeys)
	}
	{
		node := &ast.EqualOperatorNode{
			BinaryOperatorNode: ast.BinaryOperatorNode{
				Object: &ast.IdentifierNode{
					Name: []byte("a"),
					Desc: &schema.ColumnDescriptor{Table: 0, Column: 0},
				},
				Subject: &ast.FunctionOperatorNode{
					Name: &ast.IdentifierNode{Name: []byte("RAND")},
				},
			},
		}
		cl, err := planner.parseClause(node)
		s.Require().Nil(err)
		s.Require().Equal(clauseAttrForceScan, cl.Attr)
		s.Require().Equal(ColumnSet{
			&schema.ColumnDescriptor{Table: 0, Column: 0},
		}, cl.ColumnSet)
		s.Require().Zero(cl.HashKeys)
	}
}

func TestPlanner(t *testing.T) {
	suite.Run(t, new(PlannerTestSuite))
}
