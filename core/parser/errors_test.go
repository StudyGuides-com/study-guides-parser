//go:build !prod

package parser

import (
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/preparser"
)

func TestParserError_Error(t *testing.T) {
	err := &ParserError{
		Message:  "parse failed",
		Code:     CodeValidation,
		Metadata: map[string]string{"foo": "bar"},
		LineInfo: preparser.ParsedLineInfo{
			Number: 42,
			Text:   "some text",
		},
	}
	expected := "parse failed (line: 42, text: some text)"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestNewParserError(t *testing.T) {
	line := preparser.ParsedLineInfo{Number: 1, Text: "foo"}
	err := NewParserError(CodeProcessing, "msg", line)
	if err.Message != "msg" {
		t.Errorf("Message = %v, want %v", err.Message, "msg")
	}
	if err.Code != CodeProcessing {
		t.Errorf("Code = %v, want %v", err.Code, CodeProcessing)
	}
	if err.LineInfo != line {
		t.Errorf("LineInfo = %v, want %v", err.LineInfo, line)
	}
	if err.Metadata == nil || len(err.Metadata) != 0 {
		t.Errorf("Metadata should be empty map, got %v", err.Metadata)
	}
}

func TestNewGeneralError(t *testing.T) {
	err := NewGeneralError(CodeValidation, "general fail")
	if err.Message != "general fail" {
		t.Errorf("Message = %v, want %v", err.Message, "general fail")
	}
	if err.Code != CodeValidation {
		t.Errorf("Code = %v, want %v", err.Code, CodeValidation)
	}
	if err.Metadata == nil || len(err.Metadata) != 0 {
		t.Errorf("Metadata should be empty map, got %v", err.Metadata)
	}
}
