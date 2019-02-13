package planner

import (
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

// Collect common plan step constructors here, so that the cost calculation
// can be viewed at once.

func (planner *planner) newScanTable(
	tableRef schema.TableRef,
	condition *clause,
) PlanStep {
	p := &ScanTable{}
	p.Table = planner.tableRef
	if condition != nil {
		p.ColumnSet = condition.ColumnSet
		p.Condition = condition.OriginAst
	}
	p.Cost = 10 // TODO(yenlin): calculate cost.
	return p
}

func (planner *planner) newScanIndices(
	tableRef schema.TableRef,
	indexRef schema.IndexRef,
	hashKeys [][]ast.Valuer,
) PlanStep {
	p := &ScanIndices{}
	p.Cost = 0 // TODO(yenlin): calculate cost.
	p.Table = tableRef
	p.Index = indexRef
	p.Values = hashKeys
	return p
}

func (planner *planner) newScanIndexValues(
	tableRef schema.TableRef,
	indexRef schema.IndexRef,
	condition *clause,
) PlanStep {
	index := planner.table.Indices[indexRef]
	p := &ScanIndexValues{}
	p.Table = tableRef
	p.Index = indexRef
	p.Condition = condition.OriginAst
	// TODO(yenlin): calculate cost.
	p.Cost = uint64(len(index.Columns) - len(condition.ColumnSet))
	return p
}

func (planner *planner) newFilterStep(
	condition *clause,
	tableRef schema.TableRef,
	source PlanStep,
) PlanStep {
	p := &FilterStep{
		Table:     tableRef,
		ColumnSet: condition.ColumnSet,
		Condition: condition.OriginAst,
	}
	p.Operands = []PlanStep{source}
	p.Cost = source.GetCost() // TODO(yenlin): calculate cost.
	return p
}

func (planner *planner) newUnionStep(
	subPlans []PlanStep,
) PlanStep {
	p := &UnionStep{}
	p.Cost = 1 // TODO(yenlin): calculate cost.
	for _, subPlan := range subPlans {
		p.Cost += subPlan.GetCost()
	}
	p.Operands = subPlans
	return p
}
