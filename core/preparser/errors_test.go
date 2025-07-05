//go:build !prod

package preparser

import (
	"testing"
)

func TestPreParsingError_Error(t *testing.T) {
	tests := []struct {
		name     string
		error    *PreParsingError
		expected string
	}{
		{
			name: "simple error message",
			error: &PreParsingError{
				Message: "validation failed",
				Code:    CodeValidation,
				LineInfo: LineInfo{
					Number: 1,
					Type:   TokenTypeQuestion,
					Text:   "test line",
				},
			},
			expected: "validation failed",
		},
		{
			name: "empty error message",
			error: &PreParsingError{
				Message: "",
				Code:    CodeProcessing,
				LineInfo: LineInfo{
					Number: 5,
					Type:   TokenTypeHeader,
					Text:   "header line",
				},
			},
			expected: "",
		},
		{
			name: "complex error message",
			error: &PreParsingError{
				Message: "This is a very long error message with special characters: !@#$%^&*()",
				Code:    CodeValidation,
				LineInfo: LineInfo{
					Number: 10,
					Type:   TokenTypeContent,
					Text:   "content line",
				},
			},
			expected: "This is a very long error message with special characters: !@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.error.Error()
			if result != tt.expected {
				t.Errorf("Error() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewPreParsingError(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		message  string
		lineInfo LineInfo
		expected *PreParsingError
	}{
		{
			name:    "validation error",
			code:    CodeValidation,
			message: "invalid format",
			lineInfo: LineInfo{
				Number: 1,
				Type:   TokenTypeQuestion,
				Text:   "1. Question - Answer",
			},
			expected: &PreParsingError{
				Message:  "invalid format",
				Code:     CodeValidation,
				Metadata: map[string]string{},
				LineInfo: LineInfo{
					Number: 1,
					Type:   TokenTypeQuestion,
					Text:   "1. Question - Answer",
				},
			},
		},
		{
			name:    "processing error",
			code:    CodeProcessing,
			message: "internal error",
			lineInfo: LineInfo{
				Number: 5,
				Type:   TokenTypeHeader,
				Text:   "Section: Chapter: Title",
			},
			expected: &PreParsingError{
				Message:  "internal error",
				Code:     CodeProcessing,
				Metadata: map[string]string{},
				LineInfo: LineInfo{
					Number: 5,
					Type:   TokenTypeHeader,
					Text:   "Section: Chapter: Title",
				},
			},
		},
		{
			name:    "empty message",
			code:    CodeValidation,
			message: "",
			lineInfo: LineInfo{
				Number: 10,
				Type:   TokenTypeContent,
				Text:   "content",
			},
			expected: &PreParsingError{
				Message:  "",
				Code:     CodeValidation,
				Metadata: map[string]string{},
				LineInfo: LineInfo{
					Number: 10,
					Type:   TokenTypeContent,
					Text:   "content",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPreParsingError(tt.code, tt.message, tt.lineInfo)

			// Check basic fields
			if result.Message != tt.expected.Message {
				t.Errorf("Message = %v, want %v", result.Message, tt.expected.Message)
			}
			if result.Code != tt.expected.Code {
				t.Errorf("Code = %v, want %v", result.Code, tt.expected.Code)
			}
			if result.LineInfo.Number != tt.expected.LineInfo.Number {
				t.Errorf("LineInfo.Number = %v, want %v", result.LineInfo.Number, tt.expected.LineInfo.Number)
			}
			if result.LineInfo.Type != tt.expected.LineInfo.Type {
				t.Errorf("LineInfo.Type = %v, want %v", result.LineInfo.Type, tt.expected.LineInfo.Type)
			}
			if result.LineInfo.Text != tt.expected.LineInfo.Text {
				t.Errorf("LineInfo.Text = %v, want %v", result.LineInfo.Text, tt.expected.LineInfo.Text)
			}

			// Check that Metadata is initialized as empty map
			if result.Metadata == nil {
				t.Error("Metadata should be initialized as empty map, got nil")
			}
			if len(result.Metadata) != 0 {
				t.Errorf("Metadata should be empty, got %d items", len(result.Metadata))
			}
		})
	}
}

func TestNewGeneralError(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		message  string
		expected *GeneralError
	}{
		{
			name:    "validation error",
			code:    CodeValidation,
			message: "general validation failed",
			expected: &GeneralError{
				Message:  "general validation failed",
				Code:     CodeValidation,
				Metadata: map[string]string{},
			},
		},
		{
			name:    "processing error",
			code:    CodeProcessing,
			message: "general processing error",
			expected: &GeneralError{
				Message:  "general processing error",
				Code:     CodeProcessing,
				Metadata: map[string]string{},
			},
		},
		{
			name:    "empty message",
			code:    CodeValidation,
			message: "",
			expected: &GeneralError{
				Message:  "",
				Code:     CodeValidation,
				Metadata: map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewGeneralError(tt.code, tt.message)

			// Check basic fields
			if result.Message != tt.expected.Message {
				t.Errorf("Message = %v, want %v", result.Message, tt.expected.Message)
			}
			if result.Code != tt.expected.Code {
				t.Errorf("Code = %v, want %v", result.Code, tt.expected.Code)
			}

			// Check that Metadata is initialized as empty map
			if result.Metadata == nil {
				t.Error("Metadata should be initialized as empty map, got nil")
			}
			if len(result.Metadata) != 0 {
				t.Errorf("Metadata should be empty, got %d items", len(result.Metadata))
			}
		})
	}
}
