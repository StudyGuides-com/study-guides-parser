package preparser

import (
	"testing"
)

func TestPreParsingError_Error(t *testing.T) {
	lineInfo := LineInfo{Number: 1, Text: "test line", Type: TokenTypeQuestion}
	err := NewPreParsingError(CodeValidation, "test error message", lineInfo)

	expected := "test error message"
	if err.Error() != expected {
		t.Errorf("PreParsingError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestNewGeneralError(t *testing.T) {
	err := NewGeneralError(CodeProcessing, "general error message")

	if err == nil {
		t.Fatal("NewGeneralError returned nil")
	}

	if err.Code != CodeProcessing {
		t.Errorf("NewGeneralError.Code = %v, want %v", err.Code, CodeProcessing)
	}

	if err.Message != "general error message" {
		t.Errorf("NewGeneralError.Message = %v, want %v", err.Message, "general error message")
	}
}

func TestBinaryParser_Parse(t *testing.T) {
	parser := &BinaryParser{}
	lineInfo := LineInfo{Number: 1, Text: "binary\x00content", Type: TokenTypeBinary}

	result, err := parser.Parse(lineInfo)
	if err != nil {
		t.Errorf("BinaryParser.Parse() unexpected error: %v", err)
	}

	if result.Text != "binary\x00content" {
		t.Errorf("BinaryParser.Parse() returned text = %v, want %v", result.Text, "binary\x00content")
	}
}

func TestGetParserForType_UnknownType(t *testing.T) {
	lineInfo := LineInfo{Number: 1, Text: "test", Type: "unknown_type"}

	parser, err := GetParserForType[*QuestionResult]("unknown_type", lineInfo)
	if err == nil {
		t.Error("GetParserForType() expected error for unknown type, got nil")
	}

	if parser != nil {
		t.Error("GetParserForType() expected nil parser for unknown type")
	}

	expectedMsg := "unknown line type: unknown_type"
	if err.Error() != expectedMsg {
		t.Errorf("GetParserForType() error message = %v, want %v", err.Error(), expectedMsg)
	}
}

func TestPreparser_Parse_ErrorHandling(t *testing.T) {
	// Test with a line that will cause a parsing error
	lines := []LineInfo{
		{Number: 1, Text: "invalid question format", Type: TokenTypeQuestion},
	}

	preparser := NewPreparser(lines, "test")
	result, err := preparser.Parse()

	if err == nil {
		t.Error("Preparser.Parse() expected error for invalid question, got nil")
	}

	if result != nil {
		t.Error("Preparser.Parse() expected nil result when error occurs")
	}
}

func TestPreparser_Parse_GetParserError(t *testing.T) {
	// Test with an unknown token type
	lines := []LineInfo{
		{Number: 1, Text: "test", Type: "unknown_type"},
	}

	preparser := NewPreparser(lines, "test")
	result, err := preparser.Parse()

	if err == nil {
		t.Error("Preparser.Parse() expected error for unknown token type, got nil")
	}

	if result != nil {
		t.Error("Preparser.Parse() expected nil result when error occurs")
	}
}

func TestQuestionParser_Parse_EdgeCases(t *testing.T) {
	parser := &QuestionParser{}

	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
		{
			name:    "no prefix",
			line:    "What is Go? - A programming language",
			wantErr: true,
		},
		{
			name:    "no delimiter",
			line:    "1. What is Go?",
			wantErr: true,
		},
		{
			name:    "multiple delimiters",
			line:    "1. What is Go? - A programming language - extra",
			wantErr: false, // Should work, takes first delimiter
		},
		{
			name:    "only prefix and delimiter",
			line:    "1. - answer",
			wantErr: false, // Edge case: empty question text
		},
		{
			name:    "whitespace only question",
			line:    "1.    - answer",
			wantErr: false, // Should work, whitespace gets trimmed
		},
		{
			name:    "whitespace only answer",
			line:    "1. question -    ",
			wantErr: false, // Should work, whitespace gets trimmed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lineInfo := LineInfo{Number: 1, Text: tt.line, Type: TokenTypeQuestion}
			_, err := parser.Parse(lineInfo)

			if (err != nil) != tt.wantErr {
				t.Errorf("QuestionParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetParserForType_AllTypes(t *testing.T) {
	lineInfo := LineInfo{Number: 1, Text: "test", Type: TokenTypeQuestion}

	testCases := []struct {
		tokenType TokenType
		resultType any
	}{
		{TokenTypeQuestion, &QuestionResult{}},
		{TokenTypeHeader, HeaderResult{}},
		{TokenTypeComment, &CommentResult{}},
		{TokenTypeEmpty, &EmptyLineResult{}},
		{TokenTypeFileHeader, &FileHeaderResult{}},
		{TokenTypePassage, &PassageResult{}},
		{TokenTypeLearnMore, &LearnMoreResult{}},
		{TokenTypeContent, &ContentResult{}},
		{TokenTypeBinary, &BinaryResult{}},
	}

	for _, tc := range testCases {
		t.Run(string(tc.tokenType), func(t *testing.T) {
			lineInfo.Type = tc.tokenType
			
			// Use the appropriate type parameter based on the token type
			switch tc.tokenType {
			case TokenTypeQuestion:
				parser, err := GetParserForType[*QuestionResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			case TokenTypeHeader:
				parser, err := GetParserForType[HeaderResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			case TokenTypeComment:
				parser, err := GetParserForType[*CommentResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			case TokenTypeEmpty:
				parser, err := GetParserForType[*EmptyLineResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			case TokenTypeFileHeader:
				parser, err := GetParserForType[*FileHeaderResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			case TokenTypePassage:
				parser, err := GetParserForType[*PassageResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			case TokenTypeLearnMore:
				parser, err := GetParserForType[*LearnMoreResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			case TokenTypeContent:
				parser, err := GetParserForType[*ContentResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			case TokenTypeBinary:
				parser, err := GetParserForType[*BinaryResult](tc.tokenType, lineInfo)
				if err != nil {
					t.Errorf("GetParserForType() unexpected error for %v: %v", tc.tokenType, err)
				}
				if parser == nil {
					t.Errorf("GetParserForType() returned nil parser for %v", tc.tokenType)
				}
			}
		})
	}
}

func TestContentParser_Parse_EdgeCases(t *testing.T) {
	parser := &ContentParser{}

	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{
			name:    "normal content",
			line:    "This is normal content",
			wantErr: false,
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			line:    "   \t\n\r   ",
			wantErr: true,
		},
		{
			name:    "content with invisible characters",
			line:    "Content with\u200Bzero-width\u200Bspace",
			wantErr: false, // Should work, invisible chars get removed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lineInfo := LineInfo{Number: 1, Text: tt.line, Type: TokenTypeContent}
			_, err := parser.Parse(lineInfo)

			if (err != nil) != tt.wantErr {
				t.Errorf("ContentParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
