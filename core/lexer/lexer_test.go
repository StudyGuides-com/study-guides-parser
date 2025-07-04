//go:build !prod

package lexer

import (
	"testing"
)

func TestNewLexerBasic(t *testing.T) {
	lexer := NewLexer()

	if lexer == nil {
		t.Fatal("NewLexer returned nil")
	}

	if len(lexer.classifiers) == 0 {
		t.Error("no classifiers were initialized")
	}

	// Verify classifier order
	expectedOrder := []string{
		"isBinary",
		"isFileHeader",
		"isComment",
		"isQuestion",
		"isHeader",
		"isPassage",
		"isLearnMore",
		"isEmpty",
	}

	if len(lexer.classifiers) != len(expectedOrder) {
		t.Errorf("expected %d classifiers, got %d", len(expectedOrder), len(lexer.classifiers))
	}
}

func TestProcessLineBasic(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "process valid question",
			line:        "1. What is Go? - A programming language",
			lineNum:     2,
			wantType:    TokenTypeQuestion,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "process valid header",
			line:        "Subject: Topic: Subtopic",
			lineNum:     2,
			wantType:    TokenTypeHeader,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "process valid comment",
			line:        "# This is a comment",
			lineNum:     2,
			wantType:    TokenTypeComment,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "process empty line",
			line:        "",
			lineNum:     2,
			wantType:    TokenTypeEmpty,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "process file header on first line",
			line:        "Study Guide: Go Programming",
			lineNum:     1,
			wantType:    TokenTypeFileHeader,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "process valid passage",
			line:        "Passage: This is a passage",
			lineNum:     2,
			wantType:    TokenTypePassage,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "process learn more line",
			line:        "Learn More: Additional information",
			lineNum:     2,
			wantType:    TokenTypeLearnMore,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "process binary content",
			line:        "Normal text\x00with null byte",
			lineNum:     2,
			wantType:    TokenTypeBinary,
			wantErrCode: CodeBinaryContent,
			wantErrMsg:  "contains binary or non-printable characters",
		},
		{
			name:        "process content with no specific type",
			line:        "This is regular content",
			lineNum:     2,
			wantType:    TokenTypeContent,
			wantErrCode: "",
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer()

			got, err := lexer.ProcessLine(tt.line, tt.lineNum)

			// Check line info
			if got.Type != tt.wantType {
				t.Errorf("ProcessLine() got type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Number != tt.lineNum {
				t.Errorf("ProcessLine() got line number = %v, want %v", got.Number, tt.lineNum)
			}
			if got.Text != tt.line {
				t.Errorf("ProcessLine() got text = %v, want %v", got.Text, tt.line)
			}

			// Check error
			if tt.wantErrMsg != "" {
				if err == nil {
					t.Errorf("ProcessLine() expected error with message %q, got nil", tt.wantErrMsg)
				} else if err.Error() != tt.wantErrMsg {
					t.Errorf("ProcessLine() got error = %v, want %v", err, tt.wantErrMsg)
				}
			} else if err != nil {
				t.Errorf("ProcessLine() unexpected error: %v", err)
			}
		})
	}
}

func TestProcessLineIntegrationBasic(t *testing.T) {
	lexer := NewLexer()

	// Test processing multiple lines in sequence
	lines := []struct {
		line    string
		lineNum int
		want    TokenType
	}{
		{"Study Guide: Go Programming", 1, TokenTypeFileHeader},
		{"# This is a comment", 2, TokenTypeComment},
		{"", 3, TokenTypeEmpty},
		{"Subject: Topic: Subtopic", 4, TokenTypeHeader},
		{"1. What is Go? - A programming language", 5, TokenTypeQuestion},
		{"Passage: This is a passage", 6, TokenTypePassage},
		{"Learn More: Additional information", 7, TokenTypeLearnMore},
		{"This is regular content", 8, TokenTypeContent},
	}

	for _, tt := range lines {
		t.Run(tt.line, func(t *testing.T) {
			got, err := lexer.ProcessLine(tt.line, tt.lineNum)
			if err != nil {
				t.Errorf("ProcessLine() unexpected error: %v", err)
			}
			if got.Type != tt.want {
				t.Errorf("ProcessLine() got type = %v, want %v", got.Type, tt.want)
			}
		})
	}
}
