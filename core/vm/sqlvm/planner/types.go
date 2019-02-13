package planner

import (
	"github.com/shopspring/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

// clause types.

// clauseAttr is the attribute bitmap for every clause.
type clauseAttr uint64

// clauseAttr enums.
const (
	clauseAttrConst clauseAttr = 1 << iota
	clauseAttrColumn
	clauseAttrEnumerable
	clauseAttrForceScan
	clauseAttrAnd
	clauseAttrOr
)

// clause contains metadata used by planner about operation nodes.
type clause struct {
	ColumnSet ColumnSet
	Attr      clauseAttr
	HashKeys  [][]ast.Valuer
	OriginAst ast.ExprNode
	SubCls    []*clause `print:"-"`
}

// Plan types.

// PlanStep is the interface for all plan or sub-plans.
type PlanStep interface {
	GetCost() uint64
	SetCost(uint64)
	GetOutNum() uint64
	SetOutNum(uint64)
	GetOperands() []PlanStep
}

// PlanStepBase implements PlanStep interface.
type PlanStepBase struct {
	Cost     uint64
	OutNum   uint64
	Operands []PlanStep
}

// GetCost gets the cost of the plan.
func (b PlanStepBase) GetCost() uint64 {
	return b.Cost
}

// SetCost sets the cost of the plan.
func (b *PlanStepBase) SetCost(c uint64) {
	b.Cost = c
}

// GetOutNum gets the estimated output row count of the plan.
func (b PlanStepBase) GetOutNum() uint64 {
	return b.OutNum
}

// SetOutNum sets the estimated output row count of the plan.
func (b *PlanStepBase) SetOutNum(n uint64) {
	b.OutNum = n
}

// GetOperands gets the sub-plans to generate the operands of this plan.
func (b PlanStepBase) GetOperands() []PlanStep {
	return b.Operands
}

var _ PlanStep = (*PlanStepBase)(nil)

// Plan step types.

// ScanTable means to scan whole table.
type ScanTable struct {
	PlanStepBase
	Table     schema.TableRef
	ColumnSet ColumnSet
	Condition ast.ExprNode
}

// ScanIndices means to scan known hash keys on a index.
type ScanIndices struct {
	PlanStepBase
	Table  schema.TableRef
	Index  schema.IndexRef
	Values [][]ast.Valuer
}

// ScanIndexValues means to scan all possible values on a index.
type ScanIndexValues struct {
	PlanStepBase
	Table     schema.TableRef
	Index     schema.IndexRef
	Condition ast.ExprNode
}

// FilterStep means to further filter the rows in Operands[0].
type FilterStep struct {
	PlanStepBase
	Table     schema.TableRef
	ColumnSet ColumnSet
	Condition ast.ExprNode
}

// UnionStep means to take union the results from Operands.
type UnionStep struct {
	PlanStepBase
}

// IntersectStep means to take intersect of the result from Operands.
type IntersectStep struct {
	PlanStepBase
}

// InsertStep inserts rows (with Columns = Values) into Table.
type InsertStep struct {
	PlanStepBase
	Table   schema.TableRef
	Columns []schema.ColumnRef
	Values  [][]ast.ExprNode
}

// SelectStep select Columns from the rows in the result of Operands[0].
type SelectStep struct {
	PlanStepBase
	Table     schema.TableRef
	ColumnSet ColumnSet
	Columns   []ast.ExprNode
	Order     []*ast.OrderOptionNode
	Offset    *decimal.Decimal
	Limit     *decimal.Decimal
}

// SelectWithoutTable is a special case of select when table is nil.
// In this case, we should output 1 or 0 row data depending on whether the
// condition is true.
type SelectWithoutTable struct {
	PlanStepBase
	Columns   []ast.ExprNode
	Condition ast.ExprNode
	Offset    *decimal.Decimal
	Limit     *decimal.Decimal
}

// UpdateStep does the assignment to the rows in the result of Operands[0].
type UpdateStep struct {
	PlanStepBase
	Table     schema.TableRef
	ColumnSet ColumnSet
	Columns   []schema.ColumnRef
	Values    []ast.ExprNode
}

// DeleteStep deletes the rows in the result of Operands[0] from the table.
type DeleteStep struct {
	PlanStepBase
	Table schema.TableRef
}
