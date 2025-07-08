//go:build !prod

package lexer

import (
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/cleanstring"
)

// errorMatches checks if a LexerError matches the expected error code and message
func errorMatches(got *LexerError, wantCode ErrorCode, wantMessage string) bool {
	if got == nil {
		return wantCode == "" && wantMessage == ""
	}
	return got.Code == wantCode && got.Message == wantMessage
}

func TestIsQuestion(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "valid numbered question",
			line:        "1. What is Go? - A programming language",
			lineNum:     1,
			wantType:    TokenTypeQuestion,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid bullet question",
			line:        "* What is Go? - A programming language",
			lineNum:     1,
			wantType:    TokenTypeQuestion,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:       "valid dash question",
			line:       "- What is Go? - A programming language",
			lineNum:    1,
			wantType:   TokenTypeQuestion,
			wantErrMsg: "",
		},
		{
			name:        "missing answer delimiter",
			line:        "1. What is Go?",
			lineNum:     1,
			wantType:    TokenTypeQuestion,
			wantErrCode: CodeMissingAnswerDelimiter,
			wantErrMsg:  "missing answer delimiter ' - '",
		},
		{
			name:        "no list prefix",
			line:        "What is Go? - A programming language",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotErr := isQuestion(tt.line, tt.lineNum)
			if gotType != tt.wantType {
				t.Errorf("isQuestion() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !errorMatches(gotErr, tt.wantErrCode, tt.wantErrMsg) {
				t.Errorf("isQuestion() gotErr = %v, want code=%v msg=%v", gotErr, tt.wantErrCode, tt.wantErrMsg)
			}
		})
	}
}

func TestIsHeader(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "valid header with two colons",
			line:        "Subject: Topic: Subtopic",
			lineNum:     1,
			wantType:    TokenTypeHeader,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid header with three colons",
			line:        "Subject: Topic: Subtopic: Detail",
			lineNum:     1,
			wantType:    TokenTypeHeader,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not a header - only one colon",
			line:        "Subject: Topic",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not a header - is a passage",
			line:        "Passage: Some passage",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not a header - is a question",
			line:        "1. Question - Answer",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not a header - is learn more",
			line:        "Learn More: Some info",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotErr := isHeader(tt.line, tt.lineNum)
			if gotType != tt.wantType {
				t.Errorf("isHeader() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !errorMatches(gotErr, tt.wantErrCode, tt.wantErrMsg) {
				t.Errorf("isHeader() gotErr = %v, want code=%v msg=%v", gotErr, tt.wantErrCode, tt.wantErrMsg)
			}
		})
	}
}

func TestIsComment(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "valid single hash comment",
			line:        "# This is a comment",
			lineNum:     1,
			wantType:    TokenTypeComment,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid comment with leading spaces",
			line:        "   # This is a comment",
			lineNum:     1,
			wantType:    TokenTypeComment,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not a comment - double hash",
			line:        "## This is not a comment",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not a comment - no hash",
			line:        "This is not a comment",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotErr := isComment(tt.line, tt.lineNum)
			if gotType != tt.wantType {
				t.Errorf("isComment() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !errorMatches(gotErr, tt.wantErrCode, tt.wantErrMsg) {
				t.Errorf("isComment() gotErr = %v, want code=%v msg=%v", gotErr, tt.wantErrCode, tt.wantErrMsg)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "empty string",
			line:        "",
			lineNum:     1,
			wantType:    TokenTypeEmpty,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "whitespace only",
			line:        "   ",
			lineNum:     1,
			wantType:    TokenTypeEmpty,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "tabs only",
			line:        "\t\t",
			lineNum:     1,
			wantType:    TokenTypeEmpty,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "mixed whitespace",
			line:        " \t \n \r ",
			lineNum:     1,
			wantType:    TokenTypeEmpty,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not empty - has content",
			line:        "  Some content  ",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate how the lexer calls the classifier with cleaned line
			cleaned := cleanstring.New(tt.line).Clean()
			gotType, gotErr := isEmpty(cleaned, tt.lineNum)
			if gotType != tt.wantType {
				t.Errorf("isEmpty() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !errorMatches(gotErr, tt.wantErrCode, tt.wantErrMsg) {
				t.Errorf("isEmpty() gotErr = %v, want code=%v msg=%v", gotErr, tt.wantErrCode, tt.wantErrMsg)
			}
		})
	}
}

func TestIsFileHeader(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "valid file header on first line",
			line:        "Study Guide: Go Programming",
			lineNum:     1,
			wantType:    TokenTypeFileHeader,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid file header with more content",
			line:        "Advanced Placement (AP): AP Computer Science A Exam",
			lineNum:     1,
			wantType:    TokenTypeFileHeader,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not first line - should not be file header",
			line:        "Study Guide: Go Programming",
			lineNum:     2,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "empty first line - should error",
			line:        "",
			lineNum:     1,
			wantType:    TokenTypeEmpty,
			wantErrCode: CodeMissingFileHeader,
			wantErrMsg:  "first line cannot be empty",
		},
		{
			name:        "regular header on first line - should error",
			line:        "Subject: Topic: Subtopic",
			lineNum:     1,
			wantType:    TokenTypeHeader,
			wantErrCode: CodeMissingFileHeader,
			wantErrMsg:  "first line must be a file header, not a regular header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotErr := isFileHeader(tt.line, tt.lineNum)
			if gotType != tt.wantType {
				t.Errorf("isFileHeader() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !errorMatches(gotErr, tt.wantErrCode, tt.wantErrMsg) {
				t.Errorf("isFileHeader() gotErr = %v, want code=%v msg=%v", gotErr, tt.wantErrCode, tt.wantErrMsg)
			}
		})
	}
}

func TestIsPassage(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "valid passage with single hash",
			line:        "### Passage: This is a passage",
			lineNum:     1,
			wantType:    TokenTypePassage,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid passage with double hash",
			line:        "#### Passage: This is a passage",
			lineNum:     1,
			wantType:    TokenTypePassage,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid passage without hash",
			line:        "Passage: This is a passage",
			lineNum:     1,
			wantType:    TokenTypePassage,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid passage with leading spaces",
			line:        "   Passage: This is a passage",
			lineNum:     1,
			wantType:    TokenTypePassage,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not a passage - different prefix",
			line:        "Section: This is not a passage",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not a passage - no prefix",
			line:        "This is not a passage",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotErr := isPassage(tt.line, tt.lineNum)
			if gotType != tt.wantType {
				t.Errorf("isPassage() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !errorMatches(gotErr, tt.wantErrCode, tt.wantErrMsg) {
				t.Errorf("isPassage() gotErr = %v, want code=%v msg=%v", gotErr, tt.wantErrCode, tt.wantErrMsg)
			}
		})
	}
}

func TestIsLearnMore(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "valid learn more line",
			line:        "Learn More: Additional information",
			lineNum:     1,
			wantType:    TokenTypeLearnMore,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid learn more with lowercase m",
			line:        "Learn more: Additional information",
			lineNum:     1,
			wantType:    TokenTypeLearnMore,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "valid learn more with leading spaces",
			line:        "   Learn More: Additional information",
			lineNum:     1,
			wantType:    TokenTypeLearnMore,
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not learn more - different prefix",
			line:        "Read More: Additional information",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "not learn more - no prefix",
			line:        "Additional information",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotErr := isLearnMore(tt.line, tt.lineNum)
			if gotType != tt.wantType {
				t.Errorf("isLearnMore() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !errorMatches(gotErr, tt.wantErrCode, tt.wantErrMsg) {
				t.Errorf("isLearnMore() gotErr = %v, want code=%v msg=%v", gotErr, tt.wantErrCode, tt.wantErrMsg)
			}
		})
	}
}

func TestIsBinary(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		lineNum     int
		wantType    TokenType
		wantErrCode ErrorCode
		wantErrMsg  string
	}{
		{
			name:        "contains null byte",
			line:        "Normal text\x00with null byte",
			lineNum:     1,
			wantType:    TokenTypeBinary,
			wantErrCode: CodeBinaryContent,
			wantErrMsg:  "contains binary or non-printable characters",
		},
		{
			name:        "contains control characters",
			line:        "Text with\x01control\x02chars",
			lineNum:     1,
			wantType:    TokenTypeBinary,
			wantErrCode: CodeBinaryContent,
			wantErrMsg:  "contains binary or non-printable characters",
		},
		{
			name:        "normal text - should not be binary",
			line:        "This is normal text",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
		{
			name:        "text with tabs and newlines - should not be binary",
			line:        "Text\twith\nnewlines\r",
			lineNum:     1,
			wantType:    "",
			wantErrCode: "",
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotErr := isBinary(tt.line, tt.lineNum)
			if gotType != tt.wantType {
				t.Errorf("isBinary() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !errorMatches(gotErr, tt.wantErrCode, tt.wantErrMsg) {
				t.Errorf("isBinary() gotErr = %v, want code=%v msg=%v", gotErr, tt.wantErrCode, tt.wantErrMsg)
			}
		})
	}
}
