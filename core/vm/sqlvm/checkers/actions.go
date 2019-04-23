package checkers

import (
	"fmt"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/ast"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/errors"
	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

// CheckOptions stores boolean options for Check* functions.
type CheckOptions uint32

const (
	// CheckWithSafeMath enables overflow and underflow checks during expression
	// evaluation. An error will be thrown when the result is out of range.
	CheckWithSafeMath CheckOptions = 1 << iota
	// CheckWithSafeCast enables overflow and underflow checks during casting.
	// An error will be thrown if the value does not fit in the target type.
	CheckWithSafeCast
	// CheckWithConstantOnly restricts the expression to be a constant. An error
	// will be thrown if the expression cannot be folded into a constant.
	CheckWithConstantOnly
)

// CheckCreate runs CREATE commands to generate a database schema. It modifies
// AST in-place during evaluation of expressions.
func CheckCreate(ss []ast.StmtNode, o CheckOptions) (schema.Schema, error) {
	fn := "CheckCreate"
	s := schema.Schema{}
	c := newSchemaCache()
	el := errors.ErrorList{}

	for idx := range ss {
		if ss[idx] == nil {
			continue
		}

		switch n := ss[idx].(type) {
		case *ast.CreateTableStmtNode:
			checkCreateTableStmt(n, &s, o, c, &el)
		case *ast.CreateIndexStmtNode:
			checkCreateIndexStmt(n, &s, o, c, &el)
		default:
			el.Append(errors.Error{
				Position: ss[idx].GetPosition(),
				Length:   ss[idx].GetLength(),
				Category: errors.ErrorCategoryCommand,
				Code:     errors.ErrorCodeDisallowedCommand,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"command %s is not allowed when creating a contract",
					ast.QuoteIdentifier(ss[idx].GetVerb())),
			}, nil)
		}
	}

	if len(s) == 0 && len(el) == 0 {
		el.Append(errors.Error{
			Position: 0,
			Length:   0,
			Category: errors.ErrorCategoryCommand,
			Code:     errors.ErrorCodeNoCommand,
			Severity: errors.ErrorSeverityError,
			Prefix:   fn,
			Message:  "creating a contract without a table is not allowed",
		}, nil)
	}
	if len(el) != 0 {
		return s, el
	}
	return s, nil
}

// CheckQuery checks and modifies SELECT commands with a given database schema.
func CheckQuery(ss []ast.StmtNode, s schema.Schema, o CheckOptions) error {
	fn := "CheckQuery"
	c := newSchemaCache()
	el := errors.ErrorList{}

	for idx := range ss {
		if ss[idx] == nil {
			continue
		}

		switch n := ss[idx].(type) {
		case *ast.SelectStmtNode:
			checkSelectStmt(n, s, o, c, &el)
		default:
			el.Append(errors.Error{
				Position: ss[idx].GetPosition(),
				Length:   ss[idx].GetLength(),
				Category: errors.ErrorCategoryCommand,
				Code:     errors.ErrorCodeDisallowedCommand,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"command %s is not allowed when calling query",
					ast.QuoteIdentifier(ss[idx].GetVerb())),
			}, nil)
		}
	}
	if len(el) != 0 {
		return el
	}
	return nil
}

// CheckExec checks and modifies UPDATE, DELETE, INSERT commands with a given
// database schema.
func CheckExec(ss []ast.StmtNode, s schema.Schema, o CheckOptions) error {
	fn := "CheckExec"
	c := newSchemaCache()
	el := errors.ErrorList{}

	for idx := range ss {
		if ss[idx] == nil {
			continue
		}

		switch n := ss[idx].(type) {
		case *ast.UpdateStmtNode:
			checkUpdateStmt(n, s, o, c, &el)
		case *ast.DeleteStmtNode:
			checkDeleteStmt(n, s, o, c, &el)
		case *ast.InsertStmtNode:
			checkInsertStmt(n, s, o, c, &el)
		default:
			el.Append(errors.Error{
				Position: ss[idx].GetPosition(),
				Length:   ss[idx].GetLength(),
				Category: errors.ErrorCategoryCommand,
				Code:     errors.ErrorCodeDisallowedCommand,
				Severity: errors.ErrorSeverityError,
				Prefix:   fn,
				Message: fmt.Sprintf(
					"command %s is not allowed when calling exec",
					ast.QuoteIdentifier(ss[idx].GetVerb())),
			}, nil)
		}
	}
	if len(el) != 0 {
		return el
	}
	return nil
}
