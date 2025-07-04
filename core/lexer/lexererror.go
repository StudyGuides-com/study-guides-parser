package lexer

// ErrorCode represents a service error code
type ErrorCode string

const (
	CodeInvalidToken ErrorCode = "INVALID_TOKEN"
	// Question validation errors
	CodeMissingAnswerDelimiter ErrorCode = "MISSING_ANSWER_DELIMITER"
	// Binary content errors
	CodeBinaryContent ErrorCode = "BINARY_CONTENT"
	// File header errors
	CodeMissingFileHeader ErrorCode = "MISSING_FILE_HEADER"
)

// GeneralError is a base struct for all error types
type LexerError struct {
	Message  string
	Metadata map[string]string
	Code     ErrorCode
	LineInfo LineInfo
}

// Error implements the error interface
func (e *LexerError) Error() string {
	return e.Message
}

// NewParsingError creates a new parsing error
func NewLexerError(code ErrorCode, message string, lineInfo LineInfo) *LexerError {
	return &LexerError{
		Message:  message,
		Metadata: make(map[string]string),
		Code:     code,
		LineInfo: lineInfo,
	}
}
