package errors

import (
	"fmt"
	"strconv"
	"strings"
)

// Error collects error information which should be reported to users.
type Error struct {
	// These keys are parts of SQL VM ABI. Database contract callers can
	// obtain values stored in these fields from function return values.
	Position uint32 // Position is the offset in bytes to the error location.
	Category ErrorCategory
	Code     ErrorCode

	// These keys are only used for debugging purposes and not included in ABI.
	// Values stored in these fields are not guaranteed to be stable, so they
	// MUST NOT be returned to the contract caller.
	Token   string // Token is the source code token where the error occurred.
	Prefix  string // Prefix identified the cause of the error.
	Message string // Message provides detailed the error message.
}

func (e Error) Error() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("offset %d, category %d (%s), code %d (%s)",
		e.Position, e.Category, e.Category, e.Code, e.Code))
	if e.Token != "" {
		b.WriteString(", token ")
		b.WriteString(strconv.Quote(e.Token))
	}
	if e.Prefix != "" {
		b.WriteString(", hint ")
		b.WriteString(strconv.Quote(e.Prefix))
	}
	if e.Message != "" {
		b.WriteString(", message: ")
		b.WriteString(e.Message)
	}
	return b.String()
}

// ErrorList is a list of Error.
type ErrorList []Error

func (e ErrorList) Error() string {
	b := strings.Builder{}
	for i := range e {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(e[i].Error())
	}
	return b.String()
}

// ErrorCategory is used to distinguish errors come from different phases.
type ErrorCategory uint16

// Error category starts from 1. Zero value is invalid.
const (
	ErrorCategorNil ErrorCategory = iota
	ErrorCategoryGrammar
	ErrorCategorySemantic
)

var errorCategoryMap = [...]string{
	ErrorCategoryGrammar:  "grammar",
	ErrorCategorySemantic: "semantic",
}

func (c ErrorCategory) Error() string {
	return errorCategoryMap[c]
}

// ErrorCode describes the reason of the error.
type ErrorCode uint16

// Error code starts from 1. Zero value is invalid.
const (
	ErrorCodeNil ErrorCode = iota
	ErrorCodeParser
	ErrorCodeSyntax
	ErrorCodeIntegerRange
	ErrorCodeFractionalPartTooLong
)

var errorCodeMap = [...]string{
	ErrorCodeParser:                "parser error",
	ErrorCodeSyntax:                "syntax error",
	ErrorCodeIntegerRange:          "integer out of range",
	ErrorCodeFractionalPartTooLong: "fractional part too long",
}

func (c ErrorCode) Error() string {
	return errorCodeMap[c]
}
