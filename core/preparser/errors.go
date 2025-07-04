package preparser

// ErrorCode represents a service error code
type ErrorCode string

const (
	CodeValidation ErrorCode = "VALIDATION"
	CodeProcessing ErrorCode = "PROCESSING"
)

// GeneralError is a base struct for all error types
type PreParsingError struct {
	Message  string
	Metadata map[string]string
	Code     ErrorCode
	LineInfo LineInfo
}

// Error implements the error interface
func (e *PreParsingError) Error() string {
	return e.Message
}

// NewParsingError creates a new parsing error
func NewPreParsingError(code ErrorCode, message string, lineInfo LineInfo) *PreParsingError {
	return &PreParsingError{
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

func NewGeneralError(code ErrorCode, message string) *GeneralError {
	return &GeneralError{
		Message:  message,
		Metadata: map[string]string{},
		Code:     code,
	}
}
