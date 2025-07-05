package parser

import (
	"fmt"

	"github.com/StudyGuides-com/study-guides-parser/core/preparser"
)

// ErrorCode represents a service error code
type ErrorCode string

const (
	CodeValidation ErrorCode = "VALIDATION"
	CodeProcessing ErrorCode = "PROCESSING"
)

// ParserError represents a parsing error with context
type ParserError struct {
	Message  string
	Metadata map[string]string
	Code     ErrorCode
	LineInfo preparser.ParsedLineInfo
}

// Error implements the error interface
func (e *ParserError) Error() string {
	return fmt.Sprintf("%s (line: %d, text: %s)", e.Message, e.LineInfo.Number, e.LineInfo.Text)
}

// NewParserError creates a new parser error with the given code, message, and line information.
// This function is used to create consistent error instances throughout the parser package.
func NewParserError(code ErrorCode, message string, lineInfo preparser.ParsedLineInfo) *ParserError {
	return &ParserError{
		Message:  message,
		Metadata: map[string]string{},
		Code:     code,
		LineInfo: lineInfo,
	}
}

type GeneralError struct {
	Message  string
	Metadata map[string]string
	Code     ErrorCode
}

// NewGeneralError creates a new general error with the given code and message.
// This function is used for errors that are not specific to a particular line.
func NewGeneralError(code ErrorCode, message string) *GeneralError {
	return &GeneralError{
		Message:  message,
		Metadata: map[string]string{},
		Code:     code,
	}
}
