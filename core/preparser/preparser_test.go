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

func TestGetParserForType(t *testing.T) {
	tests := []struct {
		name     string
		lineType TokenType
		wantErr  bool
	}{
		{
			name:     "valid question type",
			lineType: TokenTypeQuestion,
			wantErr:  false,
		},
		{
			name:     "valid header type",
			lineType: TokenTypeHeader,
			wantErr:  false,
		},
		{
			name:     "valid comment type",
			lineType: TokenTypeComment,
			wantErr:  false,
		},
		{
			name:     "valid empty type",
			lineType: TokenTypeEmpty,
			wantErr:  false,
		},
		{
			name:     "valid file header type",
			lineType: TokenTypeFileHeader,
			wantErr:  false,
		},
		{
			name:     "valid passage type",
			lineType: TokenTypePassage,
			wantErr:  false,
		},
		{
			name:     "valid learn more type",
			lineType: TokenTypeLearnMore,
			wantErr:  false,
		},
		{
			name:     "valid content type",
			lineType: TokenTypeContent,
			wantErr:  false,
		},
		{
			name:     "invalid type",
			lineType: "InvalidType",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lineInfo := LineInfo{Type: tt.lineType}
			parser, err := GetParserForType(tt.lineType, lineInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParserForType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && parser == nil {
				t.Error("GetParserForType() returned nil parser for valid type")
			}
		})
	}
}
