package checkers

import (
	"fmt"
	"sort"

	"github.com/dexon-foundation/decimal"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

// In addition to the convention mentioned in utils.go, we have these variable
// names in this file:
//
//   ftd -> foreign table descriptor
//   ftn -> foreign table name
//   fcd -> foreign column descriptor
//   fcn -> foreign column name
//   fid -> foreign index descriptor
//   fin -> foreign index name
//
//   fmid -> first matching index descriptor
//   fmir -> first matching index reference
//   fmin -> first matching index name

// findFirstMatchingIndex finds the first index in 'haystack' matching the
// declaration of 'needle' with attributes specified in 'attrDontCare' ignored.
// This function is considered as a part of the interface, so it have to work
// deterministically.
func findFirstMatchingIndex(haystack []schema.Index, needle schema.Index,
	attrDontCare schema.IndexAttr) (schema.IndexRef, bool) {

	compareAttr := func(a1, a2 schema.IndexAttr) bool {
		a1 = a1.GetDeclaredFlags() | attrDontCare
		a2 = a2.GetDeclaredFlags() | attrDontCare
		return a1 == a2
	}

	compareColumns := func(c1, c2 []schema.ColumnRef) bool {
		if len(c1) != len(c2) {
			return false
		}
		for ci := range c1 {
			if c1[ci] != c2[ci] {
				return false
			}
		}
		return true
	}

	for ii := range haystack {
		if compareAttr(haystack[ii].Attr, needle.Attr) &&
			compareColumns(haystack[ii].Columns, needle.Columns) {
			ir := schema.IndexRef(ii)
			return ir, true
		}
	}
	return 0, false
}

func checkCreateTableStmt(n *ast.CreateTableStmtNode, s *schema.Schema,
	o CheckOptions, c *schemaCache, el *errors.ErrorList) {

	fn := "CheckCreateTableStmt"
	hasError := false

	if c.Begin() != 0 {
		panic("schema cache must not have any open scope")
	}
	defer func() {
		if hasError {
			c.Rollback()
			return
		}
		c.Commit()
	}()

	// Return early if there are too many tables. We cannot ignore this error
	// because it will overflow schema.TableRef, which is used as a part of
	// column key in schemaCache.
	if len(*s) > schema.MaxTableRef {
		el.Append(errors.Error{
			Position: n.GetPosition(),
			Length:   n.GetLength(),
			Category: errors.ErrorCategoryLimit,
			Code:     errors.ErrorCodeTooManyTables,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message: fmt.Sprintf("cannot have more than %d tables",
				schema.MaxTableRef+1),
		}, &hasError)
		return
	}

	table := schema.Table{}
	tr := schema.TableRef(len(*s))
	td := schema.TableDescriptor{Table: tr}

	if len(n.Table.Name) == 0 {
		el.Append(errors.Error{
			Position: n.Table.GetPosition(),
			Length:   n.Table.GetLength(),
			Category: errors.ErrorCategorySemantic,
			Code:     errors.ErrorCodeEmptyTableName,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message:  "cannot create a table with an empty name",
		}, &hasError)
	}

	tn := n.Table.Name
	if !c.AddTable(string(tn), td) {
		el.Append(errors.Error{
			Position: n.Table.GetPosition(),
			Length:   n.Table.GetLength(),
			Category: errors.ErrorCategorySemantic,
			Code:     errors.ErrorCodeDuplicateTableName,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message: fmt.Sprintf("table %s already exists",
				ast.QuoteIdentifier(tn)),
		}, &hasError)
	}
	table.Name = n.Table.Name
	table.Columns = make([]schema.Column, 0, len(n.Column))

	// Handle the primary key index.
	pk := []schema.ColumnRef{}
	// Handle sequences.
	seq := 0
	// Handle indices for unique constraints.
	type localIndex struct {
		index schema.Index
		node  ast.Node
	}
	localIndices := []localIndex{}
	// Handle indices for foreign key constraints.
	type foreignNewIndex struct {
		table schema.TableDescriptor
		index schema.Index
		node  ast.Node
	}
	foreignNewIndices := []foreignNewIndex{}
	type foreignExistingIndex struct {
		index schema.IndexDescriptor
		node  ast.Node
	}
	foreignExistingIndices := []foreignExistingIndex{}

	for ci := range n.Column {
		if len(table.Columns) > schema.MaxColumnRef {
			el.Append(errors.Error{
				Position: n.Column[ci].GetPosition(),
				Length:   n.Column[ci].GetLength(),
				Category: errors.ErrorCategoryLimit,
				Code:     errors.ErrorCodeTooManyColumns,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf("cannot have more than %d columns",
					schema.MaxColumnRef+1),
			}, &hasError)
			return
		}

		column := schema.Column{}
		ok := func() (ok bool) {
			innerHasError := false
			defer func() { ok = !innerHasError }()

			// Block access to the outer hasError variable.
			hasError := struct{}{}
			// Suppress “declared and not used” error.
			_ = hasError

			c.Begin()
			defer func() {
				if innerHasError {
					c.Rollback()
					return
				}
				c.Commit()
			}()

			cr := schema.ColumnRef(len(table.Columns))
			cd := schema.ColumnDescriptor{Table: tr, Column: cr}

			if len(n.Column[ci].Column.Name) == 0 {
				el.Append(errors.Error{
					Position: n.Column[ci].Column.GetPosition(),
					Length:   n.Column[ci].Column.GetLength(),
					Category: errors.ErrorCategorySemantic,
					Code:     errors.ErrorCodeEmptyColumnName,
					Severity: errors.ErrorSeverityError,
					Prefix:   fn,
					Message:  "cannot declare a column with an empty name",
				}, &innerHasError)
			}

			cn := n.Column[ci].Column.Name
			if !c.AddColumn(string(cn), cd) {
				el.Append(errors.Error{
					Position: n.Column[ci].Column.GetPosition(),
					Length:   n.Column[ci].Column.GetLength(),
					Category: errors.ErrorCategorySemantic,
					Code:     errors.ErrorCodeDuplicateColumnName,
					Severity: errors.ErrorSeverityError,
					Prefix:   fn,
					Message: fmt.Sprintf("column %s already exists",
						ast.QuoteIdentifier(cn)),
				}, &innerHasError)
			} else {
				column.Name = n.Column[ci].Column.Name
			}

			dt, code, message := n.Column[ci].DataType.GetType()
			if code == errors.ErrorCodeNil {
				if !dt.ValidColumn() {
					el.Append(errors.Error{
						Position: n.Column[ci].DataType.GetPosition(),
						Length:   n.Column[ci].DataType.GetLength(),
						Category: errors.ErrorCategorySemantic,
						Code:     errors.ErrorCodeInvalidColumnDataType,
						Severity: errors.ErrorSeverityError,
						Prefix:   fn,
						Message: fmt.Sprintf(
							"cannot declare a column with type %s", dt.String()),
					}, &innerHasError)
				}
			} else {
				el.Append(errors.Error{
					Position: n.Column[ci].DataType.GetPosition(),
					Length:   n.Column[ci].DataType.GetLength(),
					Category: errors.ErrorCategorySemantic,
					Code:     code,
					Severity: errors.ErrorSeverityError,
					Prefix:   fn,
					Message:  message,
				}, &innerHasError)
			}
			column.Type = dt

			// Backup lengths of slices in case we have to rollback. We don't
			// have to copy slice headers or data stored in underlying arrays
			// because we always append data at the end.
			defer func(LPK, SEQ, LLI, LFNI, LFEI int) {
				if innerHasError {
					pk = pk[:LPK]
					seq = SEQ
					localIndices = localIndices[:LLI]
					foreignNewIndices = foreignNewIndices[:LFNI]
					foreignExistingIndices = foreignExistingIndices[:LFEI]
				}
			}(
				len(pk), seq, len(localIndices), len(foreignNewIndices),
				len(foreignExistingIndices),
			)

			// cs  -> constraint
			// csi -> constraint index
			for csi := range n.Column[ci].Constraint {
				// Cases are sorted in the same order as internal/grammar.peg.
			cs:
				switch cs := n.Column[ci].Constraint[csi].(type) {
				case *ast.PrimaryOptionNode:
					pk = append(pk, cr)
					column.Attr |= schema.ColumnAttrPrimaryKey

				case *ast.NotNullOptionNode:
					column.Attr |= schema.ColumnAttrNotNull

				case *ast.UniqueOptionNode:
					if (column.Attr & schema.ColumnAttrUnique) != 0 {
						// Don't create duplicate indices on a column.
						break cs
					}
					column.Attr |= schema.ColumnAttrUnique
					indexName := fmt.Sprintf("%s_%s_unique",
						table.Name, column.Name)
					idx := schema.Index{
						Name:    []byte(indexName),
						Attr:    schema.IndexAttrUnique,
						Columns: []schema.ColumnRef{cr},
					}
					localIndices = append(localIndices, localIndex{
						index: idx,
						node:  cs,
					})

				case *ast.DefaultOptionNode:
					if (column.Attr & schema.ColumnAttrHasDefault) != 0 {
						el.Append(errors.Error{
							Position: cs.GetPosition(),
							Length:   cs.GetLength(),
							Category: errors.ErrorCategorySemantic,
							Code:     errors.ErrorCodeMultipleDefaultValues,
							Severity: errors.ErrorSeverityError,
							Prefix:   fn,
							Message:  "cannot have multiple default values",
						}, &innerHasError)
						break cs
					}
					column.Attr |= schema.ColumnAttrHasDefault

					value := checkExpr(cs.Value, *s, o|CheckWithConstantOnly,
						c, el, 0, newTypeActionAssign(column.Type))
					if value == nil {
						innerHasError = true
						break cs
					}
					cs.Value = value

					switch v := cs.Value.(ast.Valuer).(type) {
					case *ast.BoolValueNode:
						sb := v.V.NullBool()
						if !sb.Valid {
							el.Append(errors.Error{
								Position: cs.GetPosition(),
								Length:   cs.GetLength(),
								Category: errors.ErrorCategorySemantic,
								Code:     errors.ErrorCodeNullDefaultValue,
								Severity: errors.ErrorSeverityError,
								Prefix:   fn,
								Message:  "default value must not be NULL",
							}, &innerHasError)
							break cs
						}
						column.Default = sb.Bool

					case *ast.AddressValueNode:
						column.Default = v.V

					case *ast.IntegerValueNode:
						column.Default = v.V

					case *ast.DecimalValueNode:
						column.Default = v.V

					case *ast.BytesValueNode:
						column.Default = v.V

					case *ast.NullValueNode:
						el.Append(errors.Error{
							Position: cs.GetPosition(),
							Length:   cs.GetLength(),
							Category: errors.ErrorCategorySemantic,
							Code:     errors.ErrorCodeNullDefaultValue,
							Severity: errors.ErrorSeverityError,
							Prefix:   fn,
							Message:  "default value must not be NULL",
						}, &innerHasError)
						break cs

					default:
						panic(fmt.Sprintf("unknown constant value type %T", n))
					}

				case *ast.ForeignOptionNode:
					if len(column.ForeignKeys) > schema.MaxForeignKeys {
						el.Append(errors.Error{
							Position: cs.GetPosition(),
							Length:   cs.GetLength(),
							Category: errors.ErrorCategoryLimit,
							Code:     errors.ErrorCodeTooManyForeignKeys,
							Severity: errors.ErrorSeverityError,
							Prefix:   fn,
							Message: fmt.Sprintf(
								"cannot have more than %d foreign key "+
									"constraints in a column",
								schema.MaxForeignKeys+1),
						}, &innerHasError)
						break cs
					}
					column.Attr |= schema.ColumnAttrHasForeignKey
					ftn := cs.Table.Name
					ftd, found := c.FindTableInBase(string(ftn))
					if !found {
						el.Append(errors.Error{
							Position: cs.Table.GetPosition(),
							Length:   cs.Table.GetLength(),
							Category: errors.ErrorCategorySemantic,
							Code:     errors.ErrorCodeTableNotFound,
							Severity: errors.ErrorSeverityError,
							Prefix:   fn,
							Message: fmt.Sprintf(
								"foreign table %s does not exist",
								ast.QuoteIdentifier(ftn)),
						}, &innerHasError)
						break cs
					}
					fcn := cs.Column.Name
					fcd, found := c.FindColumnInBase(ftd.Table, string(fcn))
					if !found {
						el.Append(errors.Error{
							Position: cs.Column.GetPosition(),
							Length:   cs.Column.GetLength(),
							Category: errors.ErrorCategorySemantic,
							Code:     errors.ErrorCodeColumnNotFound,
							Severity: errors.ErrorSeverityError,
							Prefix:   fn,
							Message: fmt.Sprintf(
								"column %s does not exist in foreign table %s",
								ast.QuoteIdentifier(fcn),
								ast.QuoteIdentifier(ftn)),
						}, &innerHasError)
						break cs
					}
					foreignType := (*s)[fcd.Table].Columns[fcd.Column].Type
					if !foreignType.Equal(column.Type) {
						el.Append(errors.Error{
							Position: n.Column[ci].DataType.GetPosition(),
							Length:   n.Column[ci].DataType.GetLength(),
							Category: errors.ErrorCategorySemantic,
							Code:     errors.ErrorCodeForeignKeyDataTypeMismatch,
							Severity: errors.ErrorSeverityError,
							Prefix:   fn,
							Message: fmt.Sprintf(
								"foreign column has type %s (%04x), but "+
									"this column has type %s (%04x)",
								foreignType.String(), uint16(foreignType),
								column.Type.String(), uint16(column.Type)),
						}, &innerHasError)
						break cs
					}

					idx := schema.Index{
						Attr:    schema.IndexAttrReferenced,
						Columns: []schema.ColumnRef{fcd.Column},
					}
					fmir, found := findFirstMatchingIndex(
						(*s)[ftd.Table].Indices, idx, schema.IndexAttrUnique)
					if found {
						fmid := schema.IndexDescriptor{
							Table: ftd.Table,
							Index: fmir,
						}
						foreignExistingIndices = append(
							foreignExistingIndices, foreignExistingIndex{
								index: fmid,
								node:  cs,
							})
					} else {
						if len(column.ForeignKeys) > 0 {
							idx.Name = []byte(fmt.Sprintf("%s_%s_foreign_key_%d",
								table.Name, column.Name, len(column.ForeignKeys)))
						} else {
							idx.Name = []byte(fmt.Sprintf("%s_%s_foreign_key",
								table.Name, column.Name))
						}
						foreignNewIndices = append(
							foreignNewIndices, foreignNewIndex{
								table: ftd,
								index: idx,
								node:  cs,
							})
					}
					column.ForeignKeys = append(column.ForeignKeys, fcd)

				case *ast.AutoIncrementOptionNode:
					if (column.Attr & schema.ColumnAttrHasSequence) != 0 {
						// Don't process AUTOINCREMENT twice.
						break cs
					}
					// We set the flag regardless of the error because we
					// don't want to produce duplicate errors.
					column.Attr |= schema.ColumnAttrHasSequence
					if seq > schema.MaxSequenceRef {
						el.Append(errors.Error{
							Position: cs.GetPosition(),
							Length:   cs.GetLength(),
							Category: errors.ErrorCategoryLimit,
							Code:     errors.ErrorCodeTooManySequences,
							Severity: errors.ErrorSeverityError,
							Prefix:   fn,
							Message: fmt.Sprintf(
								"cannot have more than %d sequences",
								schema.MaxSequenceRef+1),
						}, &innerHasError)
						break cs
					}
					major, _ := ast.DecomposeDataType(column.Type)
					switch major {
					case ast.DataTypeMajorInt, ast.DataTypeMajorUint:
					default:
						el.Append(errors.Error{
							Position: cs.GetPosition(),
							Length:   cs.GetLength(),
							Category: errors.ErrorCategorySemantic,
							Code:     errors.ErrorCodeInvalidAutoIncrementDataType,
							Prefix:   fn,
							Message: fmt.Sprintf(
								"AUTOINCREMENT is only supported on "+
									"INT and UINT types, but this column "+
									"has type %s (%04x)",
								column.Type.String(), uint16(column.Type)),
						}, &innerHasError)
						break cs
					}
					column.Sequence = schema.SequenceRef(seq)
					seq++

				default:
					panic(fmt.Sprintf("unknown constraint type %T", c))
				}
			}

			// The return value will be set by the first defer function.
			return
		}()

		// If an error occurs in the function, stop here and continue
		// processing the next column.
		if !ok {
			hasError = true
			continue
		}

		// Commit the column.
		table.Columns = append(table.Columns, column)
	}

	// Return early if there is any error.
	if hasError {
		return
	}

	mustAddIndex := func(name *[]byte, id schema.IndexDescriptor) {
		for !c.AddIndex(string(*name), id, true) {
			*name = append(*name, '_')
		}
	}

	// Create the primary key index. This is the first index on the table, so
	// it is not possible to exceed the limit on the number of indices.
	ir := schema.IndexRef(len(table.Indices))
	if len(pk) > 0 {
		idx := schema.Index{
			Name:    []byte(fmt.Sprintf("%s_primary_key", table.Name)),
			Attr:    schema.IndexAttrUnique,
			Columns: pk,
		}
		id := schema.IndexDescriptor{Table: tr, Index: ir}
		mustAddIndex(&idx.Name, id)
		table.Indices = append(table.Indices, idx)
	}

	// Create indices for the current table.
	for ii := range localIndices {
		if len(table.Indices) > schema.MaxIndexRef {
			el.Append(errors.Error{
				Position: localIndices[ii].node.GetPosition(),
				Length:   localIndices[ii].node.GetLength(),
				Category: errors.ErrorCategoryLimit,
				Code:     errors.ErrorCodeTooManyIndices,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf("cannot have more than %d indices",
					schema.MaxIndexRef+1),
			}, &hasError)
			return
		}
		idx := localIndices[ii].index
		ir := schema.IndexRef(len(table.Indices))
		id := schema.IndexDescriptor{Table: tr, Index: ir}
		mustAddIndex(&idx.Name, id)
		table.Indices = append(table.Indices, idx)
	}

	// Create indices for foreign tables.
	for ii := range foreignNewIndices {
		ftd := foreignNewIndices[ii].table
		if len((*s)[ftd.Table].Indices) > schema.MaxIndexRef {
			el.Append(errors.Error{
				Position: foreignNewIndices[ii].node.GetPosition(),
				Length:   foreignNewIndices[ii].node.GetLength(),
				Category: errors.ErrorCategoryLimit,
				Code:     errors.ErrorCodeTooManyIndices,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"table %s already has %d indices",
					ast.QuoteIdentifier((*s)[ftd.Table].Name),
					schema.MaxIndexRef+1),
			}, &hasError)
			return
		}
		idx := foreignNewIndices[ii].index
		ir := schema.IndexRef(len((*s)[ftd.Table].Indices))
		id := schema.IndexDescriptor{Table: ftd.Table, Index: ir}
		mustAddIndex(&idx.Name, id)
		(*s)[ftd.Table].Indices = append((*s)[ftd.Table].Indices, idx)
		defer func(tr schema.TableRef, length schema.IndexRef) {
			if hasError {
				(*s)[tr].Indices = (*s)[tr].Indices[:ir]
			}
		}(ftd.Table, ir)
	}

	// Mark existing indices as referenced.
	for ii := range foreignExistingIndices {
		fid := foreignExistingIndices[ii].index
		(*s)[fid.Table].Indices[fid.Index].Attr |= schema.IndexAttrReferenced
	}

	// Finally, we can commit the table definition to the schema.
	*s = append(*s, table)
}

func checkCreateIndexStmt(n *ast.CreateIndexStmtNode, s *schema.Schema,
	o CheckOptions, c *schemaCache, el *errors.ErrorList) {

	fn := "CheckCreateIndexStmt"
	hasError := false

	if c.Begin() != 0 {
		panic("schema cache must not have any open scope")
	}
	defer func() {
		if hasError {
			c.Rollback()
			return
		}
		c.Commit()
	}()

	tn := n.Table.Name
	td, found := c.FindTableInBase(string(tn))
	if !found {
		el.Append(errors.Error{
			Position: n.Table.GetPosition(),
			Length:   n.Table.GetLength(),
			Category: errors.ErrorCategorySemantic,
			Code:     errors.ErrorCodeTableNotFound,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message: fmt.Sprintf(
				"index table %s does not exist",
				ast.QuoteIdentifier(tn)),
		}, &hasError)
		return
	}

	if len(n.Column) > schema.MaxColumnRef {
		begin := n.Column[0].GetPosition()
		last := len(n.Column) - 1
		end := n.Column[last].GetPosition() + n.Column[last].GetLength()
		el.Append(errors.Error{
			Position: begin,
			Length:   end - begin,
			Category: errors.ErrorCategoryLimit,
			Code:     errors.ErrorCodeTooManyColumns,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message: fmt.Sprintf(
				"cannot create an index on more than %d columns",
				schema.MaxColumnRef+1),
		}, &hasError)
		return
	}

	columnRefs := newColumnRefSlice(uint8(len(n.Column)))
	for ci := range n.Column {
		cn := n.Column[ci].Name
		cd, found := c.FindColumnInBase(td.Table, string(cn))
		if !found {
			el.Append(errors.Error{
				Position: n.Column[ci].GetPosition(),
				Length:   n.Column[ci].GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeColumnNotFound,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"column %s does not exist in index table %s",
					ast.QuoteIdentifier(cn),
					ast.QuoteIdentifier(tn)),
			}, &hasError)
			continue
		}
		columnRefs.Append(cd.Column, uint8(ci))
	}
	if hasError {
		return
	}

	sort.Stable(columnRefs)
	for ci := 1; ci < len(n.Column); ci++ {
		if columnRefs.columns[ci] == columnRefs.columns[ci-1] {
			el.Append(errors.Error{
				Position: n.Column[columnRefs.nodes[ci]].GetPosition(),
				Length:   n.Column[columnRefs.nodes[ci]].GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeDuplicateIndexColumn,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"column %s already exists in the column list",
					ast.QuoteIdentifier(n.Column[columnRefs.nodes[ci]].Name)),
			}, &hasError)
			return
		}
	}

	index := schema.Index{}
	index.Columns = columnRefs.columns

	if len((*s)[td.Table].Indices) > schema.MaxIndexRef {
		el.Append(errors.Error{
			Position: n.GetPosition(),
			Length:   n.GetLength(),
			Category: errors.ErrorCategoryLimit,
			Code:     errors.ErrorCodeTooManyIndices,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message: fmt.Sprintf(
				"cannot have more than %d indices in table %s",
				schema.MaxIndexRef+1,
				ast.QuoteIdentifier(tn)),
		}, &hasError)
		return
	}

	ir := schema.IndexRef(len((*s)[td.Table].Indices))
	id := schema.IndexDescriptor{Table: td.Table, Index: ir}
	if n.Unique != nil {
		index.Attr |= schema.IndexAttrUnique
	}

	if len(n.Index.Name) == 0 {
		el.Append(errors.Error{
			Position: n.Index.GetPosition(),
			Length:   n.Table.GetLength(),
			Category: errors.ErrorCategorySemantic,
			Code:     errors.ErrorCodeEmptyIndexName,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message:  "cannot create an index with an empty name",
		}, &hasError)
		return
	}

	// If there is an existing index that is automatically created, rename it
	// instead of creating a new one.
	rename := false
	fmir, found := findFirstMatchingIndex((*s)[id.Table].Indices, index, 0)
	if found {
		fmid := schema.IndexDescriptor{Table: id.Table, Index: fmir}
		fmin := (*s)[id.Table].Indices[fmir].Name
		fminString := string(fmin)
		fmidCache, auto, found := c.FindIndexInBase(fminString)
		if !found {
			panic(fmt.Sprintf("index %s exists in the schema, "+
				"but it cannot be found in the schema cache",
				ast.QuoteIdentifier(fmin)))
		}
		if fmidCache != fmid {
			panic(fmt.Sprintf("index %s has descriptor %+v, "+
				"but the schema cache records it as %+v",
				ast.QuoteIdentifier(fmin), fmid, fmidCache))
		}
		if auto {
			if !c.DeleteIndex(fminString) {
				panic(fmt.Sprintf("unable to mark index %s for deletion",
					ast.QuoteIdentifier(fmin)))
			}
			rename = true
			id = fmid
			ir = id.Index
		}
	}

	in := n.Index.Name
	if !c.AddIndex(string(in), id, false) {
		el.Append(errors.Error{
			Position: n.Index.GetPosition(),
			Length:   n.Index.GetLength(),
			Category: errors.ErrorCategorySemantic,
			Code:     errors.ErrorCodeDuplicateIndexName,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message: fmt.Sprintf("index %s already exists",
				ast.QuoteIdentifier(in)),
		}, &hasError)
		return
	}

	// Commit the change into the schema.
	if rename {
		(*s)[id.Table].Indices[id.Index].Name = n.Index.Name
	} else {
		index.Name = n.Index.Name
		(*s)[id.Table].Indices = append((*s)[id.Table].Indices, index)
	}
}

func checkSelectStmt(n *ast.SelectStmtNode, s schema.Schema,
	o CheckOptions, c *schemaCache, el *errors.ErrorList) {
}

func checkUpdateStmt(n *ast.UpdateStmtNode, s schema.Schema,
	o CheckOptions, c *schemaCache, el *errors.ErrorList) {
}

func checkDeleteStmt(n *ast.DeleteStmtNode, s schema.Schema,
	o CheckOptions, c *schemaCache, el *errors.ErrorList) {
}

func checkInsertStmt(n *ast.InsertStmtNode, s schema.Schema,
	o CheckOptions, c *schemaCache, el *errors.ErrorList) {
}

func checkExpr(n ast.ExprNode,
	s schema.Schema, o CheckOptions, c *schemaCache, el *errors.ErrorList,
	tr schema.TableRef, ta typeAction) ast.ExprNode {

	switch n := n.(type) {
	case *ast.IdentifierNode:
		return checkVariable(n, s, o, c, el, tr, ta)

	case *ast.BoolValueNode:
		return checkBoolValue(n, o, el, ta)

	case *ast.AddressValueNode:
		return checkAddressValue(n, o, el, ta)

	case *ast.IntegerValueNode:
		return checkIntegerValue(n, o, el, ta)

	case *ast.DecimalValueNode:
		return checkDecimalValue(n, o, el, ta)

	case *ast.BytesValueNode:
		return checkBytesValue(n, o, el, ta)

	case *ast.NullValueNode:
		return checkNullValue(n, o, el, ta)

	case *ast.PosOperatorNode:
		return checkPosOperator(n, s, o, c, el, tr, ta)

	case *ast.NegOperatorNode:
		return n

	case *ast.NotOperatorNode:
		return n

	case *ast.ParenOperatorNode:
		return n

	case *ast.AndOperatorNode:
		return n

	case *ast.OrOperatorNode:
		return n

	case *ast.GreaterOrEqualOperatorNode:
		return n

	case *ast.LessOrEqualOperatorNode:
		return n

	case *ast.NotEqualOperatorNode:
		return n

	case *ast.EqualOperatorNode:
		return n

	case *ast.GreaterOperatorNode:
		return n

	case *ast.LessOperatorNode:
		return n

	case *ast.ConcatOperatorNode:
		return n

	case *ast.AddOperatorNode:
		return n

	case *ast.SubOperatorNode:
		return n

	case *ast.MulOperatorNode:
		return n

	case *ast.DivOperatorNode:
		return n

	case *ast.ModOperatorNode:
		return n

	case *ast.IsOperatorNode:
		return n

	case *ast.LikeOperatorNode:
		return n

	case *ast.CastOperatorNode:
		return n

	case *ast.InOperatorNode:
		return n

	case *ast.FunctionOperatorNode:
		return n

	default:
		panic(fmt.Sprintf("unknown expression type %T", n))
	}
}

func checkVariable(n *ast.IdentifierNode,
	s schema.Schema, o CheckOptions, c *schemaCache, el *errors.ErrorList,
	tr schema.TableRef, ta typeAction) ast.ExprNode {

	fn := "CheckVariable"

	if (o & CheckWithConstantOnly) != 0 {
		el.Append(errors.Error{
			Position: n.GetPosition(),
			Length:   n.GetLength(),
			Category: errors.ErrorCategorySemantic,
			Code:     errors.ErrorCodeNonConstantExpression,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message: fmt.Sprintf("%s is not a constant",
				ast.QuoteIdentifier(n.Name)),
		}, nil)
		return nil
	}

	cn := string(n.Name)
	cd, found := c.FindColumnInBase(tr, cn)
	if !found {
		el.Append(errors.Error{
			Position: n.GetPosition(),
			Length:   n.GetLength(),
			Category: errors.ErrorCategorySemantic,
			Code:     errors.ErrorCodeColumnNotFound,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message: fmt.Sprintf(
				"cannot find column %s in table %s",
				ast.QuoteIdentifier(n.Name),
				ast.QuoteIdentifier(s[tr].Name)),
		}, nil)
		return nil
	}

	cr := cd.Column
	dt := s[tr].Columns[cr].Type
	switch a := ta.(type) {
	case typeActionInferDefault:
	case typeActionInferWithSize:
	case typeActionAssign:
		if !dt.Equal(a.dt) {
			el.Append(errors.Error{
				Position: n.GetPosition(),
				Length:   n.GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeTypeError,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"expect %s (%04x), but %s (%04x) is given",
					a.dt.String(), uint16(a.dt), dt.String(), uint16(dt)),
			}, nil)
			return nil
		}
	}

	n.SetType(dt)
	n.Desc = cd
	return n
}

func checkBoolValue(n *ast.BoolValueNode,
	o CheckOptions, el *errors.ErrorList, ta typeAction) ast.ExprNode {

	fn := "CheckBoolValue"

	switch a := ta.(type) {
	case typeActionInferDefault:
	case typeActionInferWithSize:
	case typeActionAssign:
		major, _ := ast.DecomposeDataType(a.dt)
		if major != ast.DataTypeMajorBool {
			el.Append(errors.Error{
				Position: n.GetPosition(),
				Length:   n.GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeTypeError,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"expect %s (%04x), but boolean value is given",
					a.dt.String(), uint16(a.dt)),
			}, nil)
			return nil
		}
	}
	return n
}

func checkAddressValue(n *ast.AddressValueNode,
	o CheckOptions, el *errors.ErrorList, ta typeAction) ast.ExprNode {

	fn := "CheckAddressValue"

	switch a := ta.(type) {
	case typeActionInferDefault:
	case typeActionInferWithSize:
	case typeActionAssign:
		major, _ := ast.DecomposeDataType(a.dt)
		if major != ast.DataTypeMajorAddress {
			el.Append(errors.Error{
				Position: n.GetPosition(),
				Length:   n.GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeTypeError,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"expect %s (%04x), but address value is given",
					a.dt.String(), uint16(a.dt)),
			}, nil)
			return nil
		}
	}
	return n
}

func mustGetMinMax(dt ast.DataType) (decimal.Decimal, decimal.Decimal) {
	min, max, ok := dt.GetMinMax()
	if !ok {
		panic(fmt.Sprintf("GetMinMax does not handle %v", dt))
	}
	return min, max
}

func mustDecimalEncode(dt ast.DataType, d decimal.Decimal) []byte {
	b, ok := ast.DecimalEncode(dt, d)
	if !ok {
		panic(fmt.Sprintf("DecimalEncode does not handle %v", dt))
	}
	return b
}

func mustDecimalDecode(dt ast.DataType, b []byte) decimal.Decimal {
	d, ok := ast.DecimalDecode(dt, b)
	if !ok {
		panic(fmt.Sprintf("DecimalDecode does not handle %v", dt))
	}
	return d
}

func cropDecimal(dt ast.DataType, d decimal.Decimal) decimal.Decimal {
	b := mustDecimalEncode(dt, d)
	return mustDecimalDecode(dt, b)
}

func elAppendConstantTooLongError(el *errors.ErrorList, n ast.Node,
	fn string, v decimal.Decimal) {
	el.Append(errors.Error{
		Position: n.GetPosition(),
		Length:   n.GetLength(),
		Category: errors.ErrorCategorySemantic,
		Code:     errors.ErrorCodeConstantTooLong,
		Severity: errors.ErrorSeverityError,
		Prefix:   fn,
		Message: fmt.Sprintf(
			"constant expression %s has more than %d digits",
			ast.QuoteString(n.GetToken()), MaxIntegerPartDigits),
	}, nil)
}

func elAppendOverflowError(el *errors.ErrorList, n ast.Node,
	fn string, dt ast.DataType, v, min, max decimal.Decimal) {
	el.Append(errors.Error{
		Position: n.GetPosition(),
		Length:   n.GetLength(),
		Category: errors.ErrorCategorySemantic,
		Code:     errors.ErrorCodeOverflow,
		Severity: errors.ErrorSeverityError,
		Prefix:   fn,
		Message: fmt.Sprintf(
			"number %s (%s) overflow %s (%04x)",
			v.String(), ast.QuoteString(n.GetToken()),
			dt.String(), uint16(dt)),
	}, nil)
	el.Append(errors.Error{
		Position: n.GetPosition(),
		Length:   n.GetLength(),
		Category: 0,
		Code:     0,
		Severity: errors.ErrorSeverityNote,
		Prefix:   fn,
		Message: fmt.Sprintf(
			"the range of %s is [%s, %s]",
			dt.String(), min.String(), max.String()),
	}, nil)
}

func elAppendOverflowWarning(el *errors.ErrorList, n ast.Node,
	fn string, dt ast.DataType, from, to decimal.Decimal) {
	el.Append(errors.Error{
		Position: n.GetPosition(),
		Length:   n.GetLength(),
		Category: 0,
		Code:     0,
		Severity: errors.ErrorSeverityWarning,
		Prefix:   fn,
		Message: fmt.Sprintf(
			"number %s (%s) overflow %s (%04x), converted to %s",
			from.String(), ast.QuoteString(n.GetToken()),
			dt.String(), uint16(dt), to.String()),
	}, nil)
}

func checkIntegerValue(n *ast.IntegerValueNode,
	o CheckOptions, el *errors.ErrorList, ta typeAction) ast.ExprNode {

	fn := "CheckIntegerValue"

	normalizeDecimal(&n.V)
	if !safeDecimalRange(n.V) {
		elAppendConstantTooLongError(el, n, fn, n.V)
		return nil
	}

	infer := func(size int) (ast.DataType, bool) {
		// The first case: assume V fits in the signed integer.
		minor := ast.DataTypeMinor(size - 1)
		dt := ast.ComposeDataType(ast.DataTypeMajorInt, minor)
		min, max := mustGetMinMax(dt)
		// Return if V < min. V must be negative so it must be signed.
		if n.V.LessThan(min) {
			if (o & CheckWithSafeMath) != 0 {
				elAppendOverflowError(el, n, fn, dt, n.V, min, max)
				return dt, false
			}
			cropped := cropDecimal(dt, n.V)
			elAppendOverflowWarning(el, n, fn, dt, n.V, cropped)
			normalizeDecimal(&cropped)
			n.V = cropped
			return dt, true
		}
		// We are done if V fits in the signed integer.
		if n.V.LessThanOrEqual(max) {
			return dt, true
		}

		// The second case: V is a non-negative integer, but it does not fit
		// in the signed integer. Test whether the unsigned integer works.
		dt = ast.ComposeDataType(ast.DataTypeMajorUint, minor)
		min, max = mustGetMinMax(dt)
		// Return if V > max. We don't have to test whether V < min because min
		// is always zero and we already know V is non-negative.
		if n.V.GreaterThan(max) {
			if (o & CheckWithSafeMath) != 0 {
				elAppendOverflowError(el, n, fn, dt, n.V, min, max)
				return dt, false
			}
			cropped := cropDecimal(dt, n.V)
			elAppendOverflowWarning(el, n, fn, dt, n.V, cropped)
			normalizeDecimal(&cropped)
			n.V = cropped
			return dt, true
		}
		return dt, true
	}

	dt := ast.DataTypePending
	switch a := ta.(type) {
	case typeActionInferDefault:
		// Default to int256 or uint256.
		var ok bool
		dt, ok = infer(256 / 8)
		if !ok {
			return nil
		}

	case typeActionInferWithSize:
		var ok bool
		dt, ok = infer(a.size)
		if !ok {
			return nil
		}

	case typeActionAssign:
		dt = a.dt
		major, _ := ast.DecomposeDataType(dt)
		switch {
		case major == ast.DataTypeMajorAddress:
			if !n.IsAddress {
				el.Append(errors.Error{
					Position: n.GetPosition(),
					Length:   n.GetLength(),
					Category: errors.ErrorCategorySemantic,
					Code:     errors.ErrorCodeInvalidAddressChecksum,
					Severity: errors.ErrorSeverityError,
					Prefix:   fn,
					Message: fmt.Sprintf(
						"expect %s (%04x), but %s has invalid checksum",
						dt.String(), uint16(dt), n.GetToken()),
				}, nil)
				return nil
			}
			// Redirect to checkAddressValue if it becomes an address.
			addrNode := &ast.AddressValueNode{}
			addrNode.SetPosition(addrNode.GetPosition())
			addrNode.SetLength(addrNode.GetLength())
			addrNode.SetToken(addrNode.GetToken())
			addrNode.V = mustDecimalEncode(ast.ComposeDataType(
				ast.DataTypeMajorUint, ast.DataTypeMinor(160/8-1)), n.V)
			return checkAddressValue(addrNode, o, el, ta)

		case major == ast.DataTypeMajorInt,
			major == ast.DataTypeMajorUint,
			major.IsFixedRange(),
			major.IsUfixedRange():
			min, max := mustGetMinMax(dt)
			if n.V.LessThan(min) || n.V.GreaterThan(max) {
				if (o & CheckWithSafeMath) != 0 {
					elAppendOverflowError(el, n, fn, dt, n.V, min, max)
					return nil
				}
				cropped := cropDecimal(dt, n.V)
				elAppendOverflowWarning(el, n, fn, dt, n.V, cropped)
				normalizeDecimal(&cropped)
				n.V = cropped
			}

		default:
			el.Append(errors.Error{
				Position: n.GetPosition(),
				Length:   n.GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeTypeError,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"expect %s (%04x), but number value is given",
					dt.String(), uint16(dt)),
			}, nil)
			return nil
		}
	}

	if dt != ast.DataTypePending {
		n.SetType(dt)
	}
	return n
}

func checkDecimalValue(n *ast.DecimalValueNode,
	o CheckOptions, el *errors.ErrorList, ta typeAction) ast.ExprNode {

	fn := "CheckDecimalValue"

	normalizeDecimal(&n.V)
	if !safeDecimalRange(n.V) {
		elAppendConstantTooLongError(el, n, fn, n.V)
		return nil
	}

	// Redirect to checkIntegerValue if the value is an integer.
	if intPart := n.V.Truncate(0); n.V.Equal(intPart) {
		intNode := &ast.IntegerValueNode{}
		intNode.SetPosition(n.GetPosition())
		intNode.SetLength(n.GetLength())
		intNode.SetToken(n.GetToken())
		intNode.IsAddress = false
		intNode.V = n.V
		return checkIntegerValue(intNode, o, el, ta)
	}

	infer := func(size, fractionalDigits int) (ast.DataType, bool) {
		major := ast.DataTypeMajorFixed + ast.DataTypeMajor(size-1)
		minor := ast.DataTypeMinor(fractionalDigits)
		dt := ast.ComposeDataType(major, minor)
		min, max := mustGetMinMax(dt)
		if n.V.LessThan(min) {
			if (o & CheckWithSafeMath) != 0 {
				elAppendOverflowError(el, n, fn, dt, n.V, min, max)
				return dt, false
			}
			cropped := cropDecimal(dt, n.V)
			elAppendOverflowWarning(el, n, fn, dt, n.V, cropped)
			normalizeDecimal(&cropped)
			n.V = cropped
			return dt, false
		}
		if n.V.LessThanOrEqual(max) {
			return dt, true
		}

		major = ast.DataTypeMajorUfixed + ast.DataTypeMajor(size-1)
		minor = ast.DataTypeMinor(fractionalDigits)
		dt = ast.ComposeDataType(major, minor)
		min, max = mustGetMinMax(dt)
		if n.V.GreaterThan(max) {
			if (o & CheckWithSafeMath) != 0 {
				elAppendOverflowError(el, n, fn, dt, n.V, min, max)
				return dt, false
			}
			cropped := cropDecimal(dt, n.V)
			elAppendOverflowWarning(el, n, fn, dt, n.V, cropped)
			normalizeDecimal(&cropped)
			n.V = cropped
			return dt, true
		}
		return dt, true
	}

	// Now we are sure the number we are dealing has fractional part.
	dt := ast.DataTypePending
	switch a := ta.(type) {
	case typeActionInferDefault:
		// Default to fixed128x18 and ufixed128x18.
		var ok bool
		dt, ok = infer(128/8, 18)
		if !ok {
			return nil
		}

	case typeActionInferWithSize:
		// It is unclear that what the size hint means for fixed-point numbers,
		// so we just ignore it and do the same thing as the above case.
		var ok bool
		dt, ok = infer(128/8, 18)
		if !ok {
			return nil
		}

	case typeActionAssign:
		dt = a.dt
		major, _ := ast.DecomposeDataType(dt)
		switch {
		case major.IsFixedRange(),
			major.IsUfixedRange():
			min, max := mustGetMinMax(dt)
			if n.V.LessThan(min) || n.V.GreaterThan(max) {
				if (o & CheckWithSafeMath) != 0 {
					elAppendOverflowError(el, n, fn, dt, n.V, min, max)
					return nil
				}
				cropped := cropDecimal(dt, n.V)
				elAppendOverflowWarning(el, n, fn, dt, n.V, cropped)
				normalizeDecimal(&cropped)
				n.V = cropped
			}

		case major == ast.DataTypeMajorInt,
			major == ast.DataTypeMajorUint:
			el.Append(errors.Error{
				Position: n.GetPosition(),
				Length:   n.GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeTypeError,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"expect %s (%04x), but the number %s has fractional part",
					dt.String(), uint16(dt), n.V.String()),
			}, nil)
			return nil

		default:
			el.Append(errors.Error{
				Position: n.GetPosition(),
				Length:   n.GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeTypeError,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"expect %s (%04x), but number value is given",
					dt.String(), uint16(dt)),
			}, nil)
			return nil
		}
	}

	if dt != ast.DataTypePending {
		n.SetType(dt)
		_, minor := ast.DecomposeDataType(dt)
		fractionalDigits := int32(minor)
		n.V = n.V.Round(fractionalDigits)
		normalizeDecimal(&n.V)
	}
	return n
}

func checkBytesValue(n *ast.BytesValueNode,
	o CheckOptions, el *errors.ErrorList, ta typeAction) ast.ExprNode {

	fn := "CheckBytesValue"

	dt := ast.DataTypePending
	switch a := ta.(type) {
	case typeActionInferDefault:
		// Default to bytes.
		major := ast.DataTypeMajorDynamicBytes
		minor := ast.DataTypeMinorDontCare
		dt = ast.ComposeDataType(major, minor)

	case typeActionInferWithSize:
		// Still default to bytes. The size hint does not matter at all.
		major := ast.DataTypeMajorDynamicBytes
		minor := ast.DataTypeMinorDontCare
		dt = ast.ComposeDataType(major, minor)

	case typeActionAssign:
		dt = a.dt
		major, minor := ast.DecomposeDataType(dt)
		switch major {
		case ast.DataTypeMajorDynamicBytes:
			// Do nothing because it is always valid.

		case ast.DataTypeMajorFixedBytes:
			sizeGiven := len(n.V)
			sizeExpected := int(minor) + 1
			if sizeGiven != sizeExpected {
				el.Append(errors.Error{
					Position: n.GetPosition(),
					Length:   n.GetLength(),
					Category: errors.ErrorCategorySemantic,
					Code:     errors.ErrorCodeTypeError,
					Severity: errors.ErrorSeverityError,
					Prefix:   fn,
					Message: fmt.Sprintf(
						"expect %s (%04x), but %s has %d bytes",
						dt.String(), uint16(dt),
						ast.QuoteString(n.V), sizeGiven),
				}, nil)
				return nil
			}

		default:
			el.Append(errors.Error{
				Position: n.GetPosition(),
				Length:   n.GetLength(),
				Category: errors.ErrorCategorySemantic,
				Code:     errors.ErrorCodeTypeError,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"expect %s (%04x), but bytes value is given",
					dt.String(), uint16(dt)),
			}, nil)
			return nil
		}
	}

	if dt != ast.DataTypePending {
		n.SetType(dt)
	}
	return n
}

func checkNullValue(n *ast.NullValueNode,
	o CheckOptions, el *errors.ErrorList, ta typeAction) ast.ExprNode {

	dt := ast.DataTypePending
	switch a := ta.(type) {
	case typeActionInferDefault:
		dt = ast.DataTypeNull
	case typeActionInferWithSize:
		dt = ast.DataTypeNull
	case typeActionAssign:
		dt = a.dt
	}

	if dt != ast.DataTypePending {
		n.SetType(dt)
	}
	return n
}

func checkPosOperator(n *ast.PosOperatorNode,
	s schema.Schema, o CheckOptions, c *schemaCache, el *errors.ErrorList,
	tr schema.TableRef, ta typeAction) ast.ExprNode {

	/*
		fn := "CheckPosOperator"

		switch a := ta.(type) {
		case typeActionInferDefault:
		case typeActionInferWithSize:
		case typeActionAssign:
		}
	*/
	return n
}
