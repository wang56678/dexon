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
	Length   uint32 // Length is the length in bytes of the error token.
	Category ErrorCategory
	Code     ErrorCode

	// These keys are only used for debugging purposes and not included in ABI.
	// Values stored in these fields are not guaranteed to be stable, so they
	// MUST NOT be returned to the contract caller.
	Severity ErrorSeverity
	Prefix   string // Prefix identified the cause of the error.
	Message  string // Message provides detailed the error message.
}

func (e Error) Error() string {
	b := strings.Builder{}
	// It is possible for an error to have zero length because not all errors
	// correspond to tokens. The parser can report an error with no length when
	// it encounters an unexpected token.
	if e.Position > 0 || e.Length > 0 {
		b.WriteString(fmt.Sprintf("offset %d", e.Position))
		if e.Length > 0 {
			b.WriteString(fmt.Sprintf(", length %d", e.Length))
		}
	} else {
		b.WriteString("no location")
	}
	if e.Category > 0 {
		b.WriteString(fmt.Sprintf(", category %d (%s)", e.Category, e.Category))
	}
	if e.Code > 0 {
		b.WriteString(fmt.Sprintf(", code %d (%s)", e.Code, e.Code))
	}
	if e.Prefix != "" {
		b.WriteString(", prefix ")
		b.WriteString(strconv.Quote(e.Prefix))
	}
	if e.Message != "" {
		b.WriteString(", ")
		b.WriteString(e.Severity.String())
		b.WriteString(": ")
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
	ErrorCategoryNil ErrorCategory = iota
	ErrorCategoryLimit
	ErrorCategoryGrammar
	ErrorCategorySemantic
	ErrorCategoryRuntime
)

var errorCategoryMap = [...]string{
	ErrorCategoryLimit:    "limit",
	ErrorCategoryGrammar:  "grammar",
	ErrorCategorySemantic: "semantic",
	ErrorCategoryRuntime:  "runtime",
}

func (c ErrorCategory) Error() string {
	return errorCategoryMap[c]
}

// ErrorCode describes the reason of the error.
type ErrorCode uint16

// Error code starts from 1. Zero value is invalid.
const (
	ErrorCodeNil ErrorCode = iota
	ErrorCodeDepthLimitReached
	ErrorCodeParser
	ErrorCodeInvalidIntegerSyntax
	ErrorCodeInvalidNumberSyntax
	ErrorCodeIntegerOutOfRange
	ErrorCodeNumberOutOfRange
	ErrorCodeFractionalPartTooLong
	ErrorCodeEscapeSequenceTooShort
	ErrorCodeInvalidUnicodeCodePoint
	ErrorCodeUnknownEscapeSequence
	ErrorCodeInvalidBytesSize
	ErrorCodeInvalidIntSize
	ErrorCodeInvalidUintSize
	ErrorCodeInvalidFixedSize
	ErrorCodeInvalidUfixedSize
	ErrorCodeInvalidFixedFractionalDigits
	ErrorCodeInvalidUfixedFractionalDigits

	// Runtime Error
	ErrorCodeInvalidOperandNum
	ErrorCodeInvalidDataType
	ErrorCodeOverflow
	ErrorCodeIndexOutOfRange
	ErrorCodeInvalidCastType
	ErrorCodeDividedByZero
	ErrorCodeNegDecimalToUint64
	ErrorCodeDataLengthNotMatch
	ErrorCodeMultipleEscapeByte
	ErrorCodePendingEscapeByte
	ErrorCodeNoSuchFunction
)

var errorCodeMap = [...]string{
	ErrorCodeDepthLimitReached:             "depth limit reached",
	ErrorCodeParser:                        "parser error",
	ErrorCodeInvalidIntegerSyntax:          "invalid integer syntax",
	ErrorCodeInvalidNumberSyntax:           "invalid number syntax",
	ErrorCodeIntegerOutOfRange:             "integer out of range",
	ErrorCodeNumberOutOfRange:              "number out of range",
	ErrorCodeFractionalPartTooLong:         "fractional part too long",
	ErrorCodeEscapeSequenceTooShort:        "escape sequence too short",
	ErrorCodeInvalidUnicodeCodePoint:       "invalid unicode code point",
	ErrorCodeUnknownEscapeSequence:         "unknown escape sequence",
	ErrorCodeInvalidBytesSize:              "invalid bytes size",
	ErrorCodeInvalidIntSize:                "invalid int size",
	ErrorCodeInvalidUintSize:               "invalid uint size",
	ErrorCodeInvalidFixedSize:              "invalid fixed size",
	ErrorCodeInvalidUfixedSize:             "invalid ufixed size",
	ErrorCodeInvalidFixedFractionalDigits:  "invalid fixed fractional digits",
	ErrorCodeInvalidUfixedFractionalDigits: "invalid ufixed fractional digits",
	// Runtime Error
	ErrorCodeInvalidOperandNum:  "invalid operand number",
	ErrorCodeInvalidDataType:    "invalid data type",
	ErrorCodeOverflow:           "overflow",
	ErrorCodeIndexOutOfRange:    "index out of range",
	ErrorCodeInvalidCastType:    "invalid cast type",
	ErrorCodeDividedByZero:      "divide by zero",
	ErrorCodeNegDecimalToUint64: "negative deciaml to uint64",
	ErrorCodeDataLengthNotMatch: "data length not match",
	ErrorCodeMultipleEscapeByte: "multiple escape byte",
	ErrorCodePendingEscapeByte:  "pending escape byte",
	ErrorCodeNoSuchFunction:     "no such function",
}

func (c ErrorCode) Error() string {
	return errorCodeMap[c]
}

// ErrorSeverity describes the severity of the error.
type ErrorSeverity uint8

// Error severity starts from 0. Zero value indicates an error which causes an
// operation to be aborted. Other values are used for messages which are just
// informative and do not affect operations.
const (
	ErrorSeverityError ErrorSeverity = iota
	ErrorSeverityWarning
	ErrorSeverityNote
)

var errorSeverityMap = [...]string{
	ErrorSeverityError:   "error",
	ErrorSeverityWarning: "warning",
	ErrorSeverityNote:    "note",
}

func (s ErrorSeverity) String() string {
	return errorSeverityMap[s]
}
