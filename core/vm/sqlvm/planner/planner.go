package planner

import (
	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/common"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

const (
	maxHashKeySize = 64
	maxPlanDepth   = 5
)

type planner struct {
	context  *common.Context
	schema   schema.Schema
	tableRef schema.TableRef
	table    *schema.Table
}

func idNodeToTable(node *ast.IdentifierNode) schema.TableRef {
	return node.Desc.(*schema.TableDescriptor).Table
}

func idNodeToColumn(node *ast.IdentifierNode) (schema.TableRef, schema.ColumnRef) {
	desc := node.Desc.(*schema.ColumnDescriptor)
	return desc.Table, desc.Column
}

func (planner *planner) planInsert(stmt *ast.InsertStmtNode) (PlanStep, error) {
	planner.tableRef = idNodeToTable(stmt.Table)
	planner.table = &planner.schema[planner.tableRef]

	plan := &InsertStep{}
	plan.Table = planner.tableRef

	switch node := stmt.Insert.(type) {
	case *ast.InsertWithColumnOptionNode:
		plan.Columns = make([]schema.ColumnRef, len(node.Column))
		for i, node := range node.Column {
			_, plan.Columns[i] = idNodeToColumn(node)
		}
		plan.Values = node.Value
	case *ast.InsertWithDefaultOptionNode:
		plan.Values = [][]ast.ExprNode{[]ast.ExprNode{}}
	}
	return plan, nil
}

func (planner *planner) mergeAndHashKeys(cl *clause) {
	// Merge hash keys from sub-clause in disjoint AND.
	colNum := len(cl.ColumnSet)
	leftLen := len(cl.SubCls[0].HashKeys)
	rightLen := len(cl.SubCls[1].HashKeys)
	totalLen := leftLen * rightLen
	// Ignore enumerable if the hask key size is too big.
	if totalLen > maxHashKeySize {
		cl.Attr = cl.Attr &^ clauseAttrEnumerable
		return
	}
	// Prepare columns idx mapping.
	leftCols := cl.SubCls[0].ColumnSet
	leftColMap := make([]int, len(cl.SubCls[0].ColumnSet))
	rightCols := cl.SubCls[1].ColumnSet
	rightColMap := make([]int, len(cl.SubCls[1].ColumnSet))
	for i, j, k := 0, 0, 0; k < colNum; {
		switch {
		case i < len(leftCols) &&
			compareColumn(leftCols[i], cl.ColumnSet[k]) == 0:
			leftColMap[i] = k
			i++
		case j < len(rightCols) &&
			compareColumn(rightCols[j], cl.ColumnSet[k]) == 0:
			rightColMap[j] = k
			j++
		default:
			panic("Merging hash keys of invalid clauses.")
		}
		k++
	}
	cl.HashKeys = make([][]ast.Valuer, totalLen)
	for i := 0; i < totalLen; i++ {
		row := make([]ast.Valuer, colNum)
		for j, colIdx := range leftColMap {
			row[colIdx] = cl.SubCls[0].HashKeys[i%leftLen][j]
		}
		for j, colIdx := range rightColMap {
			row[colIdx] = cl.SubCls[1].HashKeys[i/leftLen][j]
		}
		cl.HashKeys[i] = row
	}
}

func (planner *planner) mergeOrHashKeys(cl *clause) {
	// Merge hash keys from sub-clause with same ColumnSet in OR.
	totalLen := len(cl.SubCls[0].HashKeys)
	totalLen += len(cl.SubCls[1].HashKeys)
	// Ignore enumerable if the hask key size is too big.
	if totalLen > maxHashKeySize {
		cl.Attr = cl.Attr &^ clauseAttrEnumerable
		return
	}
	cl.HashKeys = make([][]ast.Valuer, totalLen)
	i := 0
	for _, keys := range cl.SubCls[0].HashKeys {
		cl.HashKeys[i] = keys
		i++
	}
	for _, keys := range cl.SubCls[1].HashKeys {
		cl.HashKeys[i] = keys
		i++
	}
}

func (planner *planner) parseClause(node ast.ExprNode) (*clause, error) {
	// General parsing.
	cl := &clause{
		OriginAst: node,
	}
	children := node.GetChildren()
	// Parse sub expressions.
	for _, c := range children {
		child, ok := c.(ast.ExprNode)
		if !ok {
			continue
		}
		subCl, err := planner.parseClause(child)
		if err != nil {
			return nil, err
		}
		cl.SubCls = append(cl.SubCls, subCl)
		cl.ColumnSet = cl.ColumnSet.Join(subCl.ColumnSet)
		// General attributes.
		cl.Attr |= subCl.Attr & clauseAttrForceScan
	}

	// Special attributes.
	switch op := node.(type) {
	case *ast.IdentifierNode:
		// Column.
		switch desc := op.Desc.(type) {
		case *schema.ColumnDescriptor:
			cl.ColumnSet = []*schema.ColumnDescriptor{desc}
		default:
			// This is function name.
			// TODO(yenlin): distinguish function name from column names.
			break
		}
		// TODO(yenlin): bool column is directly enumerable.
		cl.Attr |= clauseAttrColumn
	case ast.Valuer:
		// Constant value.
		if node.IsConstant() {
			cl.Attr |= clauseAttrConst
			cl.HashKeys = [][]ast.Valuer{[]ast.Valuer{op}}
		}
	case ast.UnaryOperator:
		// Unary op.
		// TODO(yenlin): negation of bool column is still enumerable.
	case ast.BinaryOperator:
		// Binary op.
		if cl.Attr&clauseAttrForceScan != 0 {
			// No need of optimization.
			break
		}
		// Set optimization attributes by detailed type.
		switch node.(type) {
		case *ast.AndOperatorNode:
			cl.Attr |= clauseAttrAnd
			if cl.SubCls[0].ColumnSet.IsDisjoint(cl.SubCls[1].ColumnSet) &&
				(cl.SubCls[0].Attr&clauseAttrEnumerable != 0) &&
				(cl.SubCls[1].Attr&clauseAttrEnumerable != 0) {
				cl.Attr |= clauseAttrEnumerable
				planner.mergeAndHashKeys(cl)
			}
		case *ast.OrOperatorNode:
			cl.Attr |= clauseAttrOr
			if cl.SubCls[0].ColumnSet.Equal(cl.SubCls[1].ColumnSet) &&
				(cl.SubCls[0].Attr&clauseAttrEnumerable != 0) &&
				(cl.SubCls[1].Attr&clauseAttrEnumerable != 0) {
				cl.Attr |= clauseAttrEnumerable
				planner.mergeOrHashKeys(cl)
			}
		case *ast.EqualOperatorNode:
			hashableAttr := clauseAttrConst | clauseAttrColumn
			attrs := (cl.SubCls[0].Attr | cl.SubCls[1].Attr)
			attrs &= hashableAttr
			if attrs == hashableAttr {
				cl.Attr |= clauseAttrEnumerable
				if cl.SubCls[0].Attr&clauseAttrConst != 0 {
					cl.HashKeys = cl.SubCls[0].HashKeys
				} else {
					cl.HashKeys = cl.SubCls[1].HashKeys
				}
			}
		}
	case *ast.InOperatorNode:
		if cl.Attr&clauseAttrForceScan != 0 {
			// No need of optimization.
			break
		}
		left := cl.SubCls[0]
		enumerable := (left.Attr & clauseAttrColumn) != 0
		hashKeys := make([][]ast.Valuer, len(cl.SubCls)-1)
		for i, right := range cl.SubCls[1:] {
			if !enumerable {
				break
			}
			if right.Attr&clauseAttrConst == 0 {
				enumerable = false
				break
			}
			hashKeys[i] = right.HashKeys[0]
		}
		if enumerable {
			cl.Attr |= clauseAttrEnumerable
			cl.HashKeys = hashKeys
		}
	case *ast.CastOperatorNode:
		// No optimization can be done.
	case *ast.FunctionOperatorNode:
		// TODO(yenlin): enumerate the force scan function calls. (e.g. rand)
		cl.Attr |= clauseAttrForceScan
	default:
		err := errors.Error{
			Position: node.GetPosition(),
			Category: errors.ErrorCategoryPlanner,
			Code:     errors.ErrorCodePlanner,
			Message:  "Unsupported node type",
		}
		return nil, err
	}
	return cl, nil
}

func (planner *planner) parseWhere(
	whereNode *ast.WhereOptionNode,
) (*clause, error) {
	if whereNode == nil {
		return nil, nil
	}
	return planner.parseClause(whereNode.Condition)
}

func (planner *planner) planWhereclause(
	clause *clause,
	depth int,
) (PlanStep, error) {
	// TOOD(yenlin): recursion depth limit.
	var curPlan PlanStep
	checkPlan := func(newPlan PlanStep) {
		if curPlan == nil {
			curPlan = newPlan
		}
		if newPlan != nil && curPlan.GetCost() > newPlan.GetCost() {
			// Found a better plan, replace the old one.
			curPlan = newPlan
		}
	}

	// Basic brute force plan.
	p := planner.newScanTable(planner.tableRef, clause)
	checkPlan(p)

	if clause == nil {
		// No where specified, return ScanTable directly.
		return curPlan, nil
	}

	// If not forced scan.
	if clause.Attr&clauseAttrForceScan == 0 {
		// Plan on the whole clause.
		for i, index := range planner.table.Indices {
			var plan PlanStep

			var columnSet ColumnSet
			columnSet = make([]*schema.ColumnDescriptor, len(index.Columns))
			for i := range columnSet {
				columnSet[i] = &schema.ColumnDescriptor{
					Table:  planner.tableRef,
					Column: index.Columns[i],
				}
			}
			if clause.Attr&clauseAttrEnumerable != 0 &&
				columnSet.Equal(clause.ColumnSet) {
				// Values are known for hash.
				plan = planner.newScanIndices(planner.tableRef, schema.IndexRef(i),
					clause.HashKeys)
			} else if columnSet.Contains(clause.ColumnSet) {
				plan = planner.newScanIndexValues(planner.tableRef,
					schema.IndexRef(i), clause)
			}
			checkPlan(plan)
		}
	}

	if depth >= maxPlanDepth {
		// Maximum plan depth reached, just current best plan.
		return curPlan, nil
	}
	// Plan on sub-clauses.
	switch {
	case clause.Attr&clauseAttrAnd != 0:
		// It's an AND condition. Try to find a better plan by index on one
		// sub-clause and filter the other sub-clause.
		for i, subCl := range clause.SubCls {
			plan, err := planner.planWhereclause(subCl, depth+1)
			if err != nil {
				return nil, err
			}
			if plan == nil {
				continue
			}
			for j, filterCl := range clause.SubCls {
				if j != i {
					plan = planner.newFilterStep(
						filterCl, planner.tableRef, plan)
				}
			}
			checkPlan(plan)
		}
	case clause.Attr&clauseAttrOr != 0:
		// It's an OR condition. Try union from sub-clauses.
		subPlans := make([]PlanStep, len(clause.SubCls))
		for i, subCl := range clause.SubCls {
			subPlan, err := planner.planWhereclause(subCl, depth+1)
			if err != nil {
				return nil, err
			}
			subPlans[i] = subPlan
		}
		plan := planner.newUnionStep(subPlans)
		checkPlan(plan)
	}
	// Cleanup SubCls as they are no longer needed.
	return curPlan, nil
}

func (planner *planner) planWhere(
	whereNode *ast.WhereOptionNode,
) (PlanStep, error) {
	whereclause, err := planner.parseWhere(whereNode)
	if err != nil {
		return nil, err
	}
	return planner.planWhereclause(whereclause, 0)
}

func (planner *planner) planSelectWithoutTable(stmt *ast.SelectStmtNode) (
	PlanStep, error) {

	plan := &SelectWithoutTable{
		Columns: stmt.Column,
	}
	if stmt.Where != nil {
		plan.Condition = stmt.Where.Condition
	}
	if stmt.Offset != nil {
		plan.Offset = &decimal.Decimal{}
		*plan.Offset = stmt.Offset.Value.Value().(decimal.Decimal)
	}
	if stmt.Limit != nil {
		plan.Limit = &decimal.Decimal{}
		*plan.Limit = stmt.Limit.Value.Value().(decimal.Decimal)
	}
	return plan, nil
}

func (planner *planner) planSelect(stmt *ast.SelectStmtNode) (PlanStep, error) {
	if stmt.Table == nil {
		return planner.planSelectWithoutTable(stmt)
	}
	planner.tableRef = idNodeToTable(stmt.Table)
	planner.table = &planner.schema[planner.tableRef]

	wherePlan, err := planner.planWhere(stmt.Where)
	if err != nil {
		return nil, err
	}
	plan := &SelectStep{}
	plan.Table = planner.tableRef
	if stmt.Offset != nil {
		plan.Offset = &decimal.Decimal{}
		*plan.Offset = stmt.Offset.Value.Value().(decimal.Decimal)
	}
	if stmt.Limit != nil {
		plan.Limit = &decimal.Decimal{}
		*plan.Limit = stmt.Limit.Value.Value().(decimal.Decimal)
	}
	// TODO(yenlin): we may expect there's a columnset struct in every node.
	//               use clause for current development.
	plan.Columns = stmt.Column
	for _, col := range plan.Columns {
		cl, err := planner.parseClause(col)
		if err != nil {
			return nil, err
		}
		plan.ColumnSet = plan.ColumnSet.Join(cl.ColumnSet)
	}
	plan.Order = stmt.Order
	for i := range plan.Order {
		cl, err := planner.parseClause(plan.Order[i].Expr)
		if err != nil {
			return nil, err
		}
		plan.ColumnSet = plan.ColumnSet.Join(cl.ColumnSet)
	}
	plan.Operands = []PlanStep{wherePlan}
	return plan, nil
}

func (planner *planner) planUpdate(stmt *ast.UpdateStmtNode) (PlanStep, error) {
	planner.tableRef = idNodeToTable(stmt.Table)
	planner.table = &planner.schema[planner.tableRef]

	wherePlan, err := planner.planWhere(stmt.Where)
	if err != nil {
		return nil, err
	}
	plan := &UpdateStep{}
	plan.Table = planner.tableRef
	// Parse assignment.
	plan.Columns = make([]schema.ColumnRef, len(stmt.Assignment))
	plan.Values = make([]ast.ExprNode, len(stmt.Assignment))
	for i, as := range stmt.Assignment {
		cl, err := planner.parseClause(as.Expr)
		if err != nil {
			return nil, err
		}
		_, plan.Columns[i] = idNodeToColumn(as.Column)
		plan.Values[i] = as.Expr
		plan.ColumnSet = plan.ColumnSet.Join(cl.ColumnSet)
	}
	plan.Operands = []PlanStep{wherePlan}
	return plan, nil
}

func (planner *planner) planDelete(stmt *ast.DeleteStmtNode) (PlanStep, error) {
	planner.tableRef = idNodeToTable(stmt.Table)
	planner.table = &planner.schema[planner.tableRef]

	wherePlan, err := planner.planWhere(stmt.Where)
	if err != nil {
		return nil, err
	}
	plan := &DeleteStep{}
	plan.Table = planner.tableRef
	plan.Operands = []PlanStep{wherePlan}
	return plan, nil
}

// Plan the steps to achieve the sql statement in stmt.
func Plan(
	context *common.Context,
	schema schema.Schema,
	stmt ast.Node,
) (PlanStep, error) {
	planner := planner{
		context: context,
		schema:  schema,
	}

	var plan PlanStep
	var err error
	// Handle the action.
	switch st := stmt.(type) {
	case *ast.InsertStmtNode:
		plan, err = planner.planInsert(st)
	case *ast.SelectStmtNode:
		plan, err = planner.planSelect(st)
	case *ast.UpdateStmtNode:
		plan, err = planner.planUpdate(st)
	case *ast.DeleteStmtNode:
		plan, err = planner.planDelete(st)
	default:
		err = errors.Error{
			Position: stmt.GetPosition(),
			Category: errors.ErrorCategoryPlanner,
			Code:     errors.ErrorCodePlanner,
			Message:  "Unsupported statement type",
		}
	}

	return plan, err
}
