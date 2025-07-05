//go:build !prod

package preparser

import (
	"testing"
)

func TestPreparser(t *testing.T) {
	tests := []struct {
		name    string
		lines   []LineInfo
		wantErr bool
	}{
		{
			name: "valid document with all line types",
			lines: []LineInfo{
				{
					Number: 1,
					Type:   TokenTypeFileHeader,
					Text:   "Study Guide",
				},
				{
					Number: 2,
					Type:   TokenTypeHeader,
					Text:   "Section: Chapter 1: Introduction",
				},
				{
					Number: 3,
					Type:   TokenTypeContent,
					Text:   "Some content here",
				},
				{
					Number: 4,
					Type:   TokenTypeQuestion,
					Text:   "1. What is Go? - A programming language",
				},
				{
					Number: 5,
					Type:   TokenTypeComment,
					Text:   "# This is a comment",
				},
				{
					Number: 6,
					Type:   TokenTypePassage,
					Text:   "Passage: This is a passage",
				},
				{
					Number: 7,
					Type:   TokenTypeLearnMore,
					Text:   "Learn More: Additional info",
				},
				{
					Number: 8,
					Type:   TokenTypeEmpty,
					Text:   "",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid document with error in first line",
			lines: []LineInfo{
				{
					Number: 1,
					Type:   TokenTypeFileHeader,
					Text:   "Section: Chapter 1: Introduction", // Invalid file header
				},
			},
			wantErr: true,
		},
		{
			name:    "empty document",
			lines:   []LineInfo{},
			wantErr: false,
		},
		{
			name: "document with only empty lines",
			lines: []LineInfo{
				{
					Number: 1,
					Type:   TokenTypeEmpty,
					Text:   "",
				},
				{
					Number: 2,
					Type:   TokenTypeEmpty,
					Text:   "   \t  ",
				},
			},
			wantErr: false,
		},
		{
			name: "document with binary line",
			lines: []LineInfo{
				{
					Number: 1,
					Type:   TokenTypeBinary,
					Text:   "binary data with \x00\x01\x02 bytes",
				},
			},
			wantErr: false,
		},
		{
			name: "document with unknown token type",
			lines: []LineInfo{
				{
					Number: 1,
					Type:   TokenType("UNKNOWN"),
					Text:   "unknown line type",
				},
			},
			wantErr: true,
		},
		{
			name: "document with mixed types including binary",
			lines: []LineInfo{
				{
					Number: 1,
					Type:   TokenTypeFileHeader,
					Text:   "Study Guide",
				},
				{
					Number: 2,
					Type:   TokenTypeBinary,
					Text:   "binary content \x00\x01",
				},
				{
					Number: 3,
					Type:   TokenTypeContent,
					Text:   "regular content",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewPreparser(tt.lines, "test")
			result, err := parser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Preparser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(result) != len(tt.lines) {
					t.Errorf("Preparser.Parse() returned %d lines, want %d", len(result), len(tt.lines))
				}
			}
		})
	}
}

func TestPreparserParseLineBinary(t *testing.T) {
	parser := NewPreparser([]LineInfo{}, "test")

	line := LineInfo{
		Number: 1,
		Type:   TokenTypeBinary,
		Text:   "binary data \x00\x01\x02\x03",
	}

	result, err := parser.parseLine(line)
	if err != nil {
		t.Errorf("parseLine() for binary type should not return error, got %v", err)
	}

	if !result.IsBinary() {
		t.Error("parseLine() for binary type should return Binary result")
	}

	binaryResult := result.GetBinary()
	if binaryResult == nil {
		t.Error("GetBinary() should not return nil for binary type")
	}

	if binaryResult.Text != line.Text {
		t.Errorf("Binary text = %v, want %v", binaryResult.Text, line.Text)
	}
}

func TestPreparserParseLineUnknownType(t *testing.T) {
	parser := NewPreparser([]LineInfo{}, "test")

	line := LineInfo{
		Number: 1,
		Type:   TokenType("UNKNOWN_TYPE"),
		Text:   "unknown line",
	}

	result, err := parser.parseLine(line)
	if err == nil {
		t.Error("parseLine() for unknown type should return error")
	}

	if result != (ParsedValue{}) {
		t.Error("parseLine() for unknown type should return empty ParsedValue")
	}

	// Check that it's a validation error
	if err.Code != CodeValidation {
		t.Errorf("Error code = %v, want %v", err.Code, CodeValidation)
	}

	// Check that the error message contains the unknown type
	if err.Message == "" {
		t.Error("Error message should not be empty")
	}
}
