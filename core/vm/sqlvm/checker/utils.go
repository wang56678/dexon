package checker

import (
	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

// Variable name convention:
//
//   fn -> function name
//   el -> error list
//
//   td -> table descriptor
//   tr -> table reference
//   ti -> table index
//   tn -< table name
//
//   cd -> column descriptor
//   cr -> column reference
//   ci -> column index
//   cn -> column name
//
//   id -> index descriptor
//   ir -> index reference
//   ii -> index index
//   in -> index name

const (
	MaxIntegerPartDigits    int32 = 200
	MaxFractionalPartDigits int32 = 200
)

var (
	MaxConstant = func() decimal.Decimal {
		max := (decimal.New(1, MaxIntegerPartDigits).
			Sub(decimal.New(1, -MaxFractionalPartDigits)))
		normalizeDecimal(&max)
		return max
	}()
	MinConstant = MaxConstant.Neg()
)

func normalizeDecimal(d *decimal.Decimal) {
	if d.Exponent() != -MaxFractionalPartDigits {
		*d = d.Rescale(-MaxFractionalPartDigits)
	}
}

func safeDecimalRange(d decimal.Decimal) bool {
	return d.GreaterThanOrEqual(MinConstant) && d.LessThanOrEqual(MaxConstant)
}

// schemaCache is a multi-layer symbol table used to support the checker.
// It allows changes to be easily rolled back by keeping modifications in a
// separate layer, providing an experience similar to a database transaction.
type schemaCache struct {
	base   schemaCacheBase
	scopes []schemaCacheScope
}

type schemaCacheIndexValue struct {
	id   schema.IndexDescriptor
	auto bool
}

type schemaCacheColumnKey struct {
	tr schema.TableRef
	n  string
}

type schemaCacheBase struct {
	table  map[string]schema.TableDescriptor
	index  map[string]schemaCacheIndexValue
	column map[schemaCacheColumnKey]schema.ColumnDescriptor
}

func (lower *schemaCacheBase) Merge(upper schemaCacheScope) {
	// Process deletions.
	for n := range upper.tableDeleted {
		delete(lower.table, n)
	}
	for n := range upper.indexDeleted {
		delete(lower.index, n)
	}
	for ck := range upper.columnDeleted {
		delete(lower.column, ck)
	}

	// Process additions.
	for n, td := range upper.table {
		lower.table[n] = td
	}
	for n, iv := range upper.index {
		lower.index[n] = iv
	}
	for ck, cd := range upper.column {
		lower.column[ck] = cd
	}
}

type schemaCacheScope struct {
	table         map[string]schema.TableDescriptor
	tableDeleted  map[string]struct{}
	index         map[string]schemaCacheIndexValue
	indexDeleted  map[string]struct{}
	column        map[schemaCacheColumnKey]schema.ColumnDescriptor
	columnDeleted map[schemaCacheColumnKey]struct{}
}

func (lower *schemaCacheScope) Merge(upper schemaCacheScope) {
	// Process deletions.
	for n := range upper.tableDeleted {
		delete(lower.table, n)
		lower.tableDeleted[n] = struct{}{}
	}
	for n := range upper.indexDeleted {
		delete(lower.index, n)
		lower.indexDeleted[n] = struct{}{}
	}
	for ck := range upper.columnDeleted {
		delete(lower.column, ck)
		lower.columnDeleted[ck] = struct{}{}
	}

	// Process additions.
	for n, td := range upper.table {
		lower.table[n] = td
	}
	for n, iv := range upper.index {
		lower.index[n] = iv
	}
	for ck, cd := range upper.column {
		lower.column[ck] = cd
	}
}

func newSchemaCache() *schemaCache {
	return &schemaCache{
		base: schemaCacheBase{
			table:  map[string]schema.TableDescriptor{},
			index:  map[string]schemaCacheIndexValue{},
			column: map[schemaCacheColumnKey]schema.ColumnDescriptor{},
		},
	}
}

func (c *schemaCache) Begin() int {
	position := len(c.scopes)
	scope := schemaCacheScope{
		table:         map[string]schema.TableDescriptor{},
		tableDeleted:  map[string]struct{}{},
		index:         map[string]schemaCacheIndexValue{},
		indexDeleted:  map[string]struct{}{},
		column:        map[schemaCacheColumnKey]schema.ColumnDescriptor{},
		columnDeleted: map[schemaCacheColumnKey]struct{}{},
	}
	c.scopes = append(c.scopes, scope)
	return position
}

func (c *schemaCache) Rollback() {
	if len(c.scopes) == 0 {
		panic("there is no scope to rollback")
	}
	c.scopes = c.scopes[:len(c.scopes)-1]
}

func (c *schemaCache) RollbackTo(position int) {
	for position <= len(c.scopes) {
		c.Rollback()
	}
}

func (c *schemaCache) Commit() {
	if len(c.scopes) == 0 {
		panic("there is no scope to commit")
	}
	if len(c.scopes) == 1 {
		c.base.Merge(c.scopes[0])
	} else {
		src := len(c.scopes) - 1
		dst := len(c.scopes) - 2
		c.scopes[dst].Merge(c.scopes[src])
	}
	c.scopes = c.scopes[:len(c.scopes)-1]
}

func (c *schemaCache) CommitTo(position int) {
	for position <= len(c.scopes) {
		c.Commit()
	}
}

func (c *schemaCache) FindTableInBase(n string) (
	schema.TableDescriptor, bool) {

	td, exists := c.base.table[n]
	return td, exists
}

func (c *schemaCache) FindTableInScope(n string) (
	schema.TableDescriptor, bool) {

	for si := range c.scopes {
		si = len(c.scopes) - si - 1
		if td, exists := c.scopes[si].table[n]; exists {
			return td, true
		}
		if _, exists := c.scopes[si].tableDeleted[n]; exists {
			return schema.TableDescriptor{}, false
		}
	}
	return c.FindTableInBase(n)
}

func (c *schemaCache) FindTableInBaseWithFallback(n string,
	fallback schema.Schema) (schema.TableDescriptor, bool) {

	if td, found := c.FindTableInBase(n); found {
		return td, true
	}
	if fallback == nil {
		return schema.TableDescriptor{}, false
	}

	s := fallback
	for ti := range s {
		if n == string(s[ti].Name) {
			tr := schema.TableRef(ti)
			td := schema.TableDescriptor{Table: tr}
			c.base.table[n] = td
			return td, true
		}
	}
	return schema.TableDescriptor{}, false
}

func (c *schemaCache) FindIndexInBase(n string) (
	schema.IndexDescriptor, bool, bool) {

	iv, exists := c.base.index[n]
	return iv.id, iv.auto, exists
}

func (c *schemaCache) FindIndexInScope(n string) (
	schema.IndexDescriptor, bool, bool) {

	for si := range c.scopes {
		si = len(c.scopes) - si - 1
		if iv, exists := c.scopes[si].index[n]; exists {
			return iv.id, iv.auto, true
		}
		if _, exists := c.scopes[si].indexDeleted[n]; exists {
			return schema.IndexDescriptor{}, false, false
		}
	}
	return c.FindIndexInBase(n)
}

func (c *schemaCache) FindIndexInBaseWithFallback(n string,
	fallback schema.Schema) (schema.IndexDescriptor, bool, bool) {

	if id, auto, found := c.FindIndexInBase(n); found {
		return id, auto, true
	}
	if fallback == nil {
		return schema.IndexDescriptor{}, false, false
	}

	s := fallback
	for ti := range s {
		for ii := range s[ti].Indices {
			if n == string(s[ti].Indices[ii].Name) {
				tr := schema.TableRef(ti)
				ir := schema.IndexRef(ii)
				id := schema.IndexDescriptor{Table: tr, Index: ir}
				iv := schemaCacheIndexValue{id: id, auto: false}
				c.base.index[n] = iv
				return id, false, true
			}
		}
	}
	return schema.IndexDescriptor{}, false, false
}

func (c *schemaCache) FindColumnInBase(tr schema.TableRef, n string) (
	schema.ColumnDescriptor, bool) {

	cd, exists := c.base.column[schemaCacheColumnKey{tr: tr, n: n}]
	return cd, exists
}

func (c *schemaCache) FindColumnInScope(tr schema.TableRef, n string) (
	schema.ColumnDescriptor, bool) {

	ck := schemaCacheColumnKey{tr: tr, n: n}
	for si := range c.scopes {
		si = len(c.scopes) - si - 1
		if cd, exists := c.scopes[si].column[ck]; exists {
			return cd, true
		}
		if _, exists := c.scopes[si].columnDeleted[ck]; exists {
			return schema.ColumnDescriptor{}, false
		}
	}
	return c.FindColumnInBase(tr, n)
}

func (c *schemaCache) FindColumnInBaseWithFallback(tr schema.TableRef, n string,
	fallback schema.Schema) (schema.ColumnDescriptor, bool) {

	if cd, found := c.FindColumnInBase(tr, n); found {
		return cd, true
	}
	if fallback == nil {
		return schema.ColumnDescriptor{}, false
	}

	s := fallback
	for ci := range s[tr].Columns {
		if n == string(s[tr].Columns[ci].Name) {
			cr := schema.ColumnRef(ci)
			cd := schema.ColumnDescriptor{Table: tr, Column: cr}
			ck := schemaCacheColumnKey{tr: tr, n: n}
			c.base.column[ck] = cd
			return cd, true
		}
	}
	return schema.ColumnDescriptor{}, false
}

func (c *schemaCache) AddTable(n string,
	td schema.TableDescriptor) bool {

	top := len(c.scopes) - 1
	if _, found := c.FindTableInScope(n); found {
		return false
	}

	c.scopes[top].table[n] = td
	return true
}

func (c *schemaCache) AddIndex(n string,
	id schema.IndexDescriptor, auto bool) bool {

	top := len(c.scopes) - 1
	if _, _, found := c.FindIndexInScope(n); found {
		return false
	}

	iv := schemaCacheIndexValue{id: id, auto: auto}
	c.scopes[top].index[n] = iv
	return true
}

func (c *schemaCache) AddColumn(n string,
	cd schema.ColumnDescriptor) bool {

	top := len(c.scopes) - 1
	tr := cd.Table
	if _, found := c.FindColumnInScope(tr, n); found {
		return false
	}

	ck := schemaCacheColumnKey{tr: tr, n: n}
	c.scopes[top].column[ck] = cd
	return true
}

func (c *schemaCache) DeleteTable(n string) bool {
	top := len(c.scopes) - 1
	if _, found := c.FindTableInScope(n); !found {
		return false
	}

	delete(c.scopes[top].table, n)
	c.scopes[top].tableDeleted[n] = struct{}{}
	return true
}

func (c *schemaCache) DeleteIndex(n string) bool {
	top := len(c.scopes) - 1
	if _, _, found := c.FindIndexInScope(n); !found {
		return false
	}

	delete(c.scopes[top].index, n)
	c.scopes[top].indexDeleted[n] = struct{}{}
	return true
}

func (c *schemaCache) DeleteColumn(tr schema.TableRef, n string) bool {
	top := len(c.scopes) - 1
	if _, found := c.FindColumnInScope(tr, n); !found {
		return false
	}

	ck := schemaCacheColumnKey{tr: tr, n: n}
	delete(c.scopes[top].column, ck)
	c.scopes[top].columnDeleted[ck] = struct{}{}
	return true
}

// columnRefSlice implements sort.Interface. It allows sorting a slice of
// schema.ColumnRef while keeping references to AST nodes they originate from.
type columnRefSlice struct {
	columns []schema.ColumnRef
	nodes   []uint8
}

func newColumnRefSlice(c uint8) columnRefSlice {
	return columnRefSlice{
		columns: make([]schema.ColumnRef, 0, c),
		nodes:   make([]uint8, 0, c),
	}
}

func (s *columnRefSlice) Append(c schema.ColumnRef, i uint8) {
	s.columns = append(s.columns, c)
	s.nodes = append(s.nodes, i)
}

func (s columnRefSlice) Len() int {
	return len(s.columns)
}

func (s columnRefSlice) Less(i, j int) bool {
	return s.columns[i] < s.columns[j]
}

func (s columnRefSlice) Swap(i, j int) {
	s.columns[i], s.columns[j] = s.columns[j], s.columns[i]
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

// typeAction represents an action on type inference requested from the parent
// node. An action is usually only applied on a single node. It is seldom
// propagated to child nodes because we want to delay the assignment of types
// until it is necessary, making constant operations easier to use without
// being restricted by data types.
//go-sumtype:decl typeAction
type typeAction interface {
	ˉtypeAction()
}

// typeActionInferDefault requests the node to infer the type using its default
// rule. It usually means that the parent node does not care the data type,
// such as the select list in a SELECT statement. It is an advisory request.
// If the type of the node is already determined, it should ignore the request.
type typeActionInferDefault struct{}

func newTypeActionInferDefaultSize() typeActionInferDefault {
	return typeActionInferDefault{}
}

var _ typeAction = typeActionInferDefault{}

func (typeActionInferDefault) ˉtypeAction() {}

// typeActionInferWithSize requests the node to infer the type with size
// requirement. The size is measured in bytes. It is indented to be used in
// CAST to support conversion between integer and fixed-size bytes types.
// It is an advisory request. If the type is already determined, the request is
// ignored and the parent node should be able to handle the problem by itself.
type typeActionInferWithSize struct {
	size int
}

func newTypeActionInferWithSize(bytes int) typeActionInferWithSize {
	return typeActionInferWithSize{size: bytes}
}

var _ typeAction = typeActionInferWithSize{}

func (typeActionInferWithSize) ˉtypeAction() {}

type typeActionAssign struct {
	dt ast.DataType
}

// newTypeActionAssign requests the node to have a specific type. It is a
// mandatory request. If the node is unable to meet the requirement, it should
// throw an error. It is not allowed to ignore the request.
func newTypeActionAssign(expected ast.DataType) typeActionAssign {
	return typeActionAssign{dt: expected}
}

var _ typeAction = typeActionAssign{}

func (typeActionAssign) ˉtypeAction() {}
