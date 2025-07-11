package processor

import (
	"os"
	"strings"
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/config"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		lines         []string
		metadata      *config.Metadata
		wantErr       bool
		wantErrSubstr string
	}{
		{
			name: "valid college study guide",
			lines: []string{
				"Mathematics Study Guide",
				"Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations",
				"",
				"1. What is a linear equation? - An equation where the highest power of the variable is 1.",
				"2. How do you solve 2x + 3 = 7? - Subtract 3 from both sides: 2x = 4, then divide by 2: x = 2.",
				"",
				"Learn More: See Khan Academy's linear equations course.",
			},
			metadata: config.NewMetadata("test_parser"),
			wantErr:  false,
		},
		{
			name: "valid AP exam study guide",
			lines: []string{
				"AP Calculus AB Study Guide",
				"Advanced Placement (AP): AP Calculus AB: Derivatives: Introduction to Derivatives",
				"",
				"1. What is a derivative? - The rate of change of a function.",
				"2. How do you find the derivative of x²? - Use the power rule: 2x.",
			},
			metadata: config.NewMetadata("ap_test_parser"),
			wantErr:  false,
		},
		{
			name: "valid study guide with learn more line containing leading whitespace",
			lines: []string{
				"Military Study Guide",
				"DOD: Navy: Regulations: Leave Policy",
				"",
				"1. What happens if an absence extends over special liberty? - The entire period is charged as leave.",
				"    Learn More: If an unavoidable absence extends over special liberty exceeding 24 hours and return occurs after 0900, the entire period is charged as leave from the day special liberty started through the day of return, as defined in MILPERSMAN 1050-110.",
			},
			metadata: config.NewMetadata("military_test_parser"),
			wantErr:  false,
		},
		{
			name: "missing file header",
			lines: []string{
				"Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations",
				"1. What is x? - A variable",
			},
			metadata:      config.NewMetadata("test_parser"),
			wantErr:       true,
			wantErrSubstr: "file header",
		},
		{
			name:          "empty lines",
			lines:         []string{},
			metadata:      config.NewMetadata("test_parser"),
			wantErr:       true,
			wantErrSubstr: "no lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.lines, tt.metadata)
			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
			}
			if tt.wantErr {
				if result.Success {
					t.Errorf("Parse() expected error but got success")
				} else if tt.wantErrSubstr != "" {
					found := false
					for _, parseErr := range result.Errors {
						if strings.Contains(parseErr.Message, tt.wantErrSubstr) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Parse() error messages = %v, want substring %q", result.Errors, tt.wantErrSubstr)
					}
				}
				return
			}
			if !result.Success {
				t.Errorf("Parse() unexpected errors: %v", result.Errors)
				return
			}
			if result.AST == nil {
				t.Errorf("Parse() returned nil AST")
				return
			}
			if result.AST.Root == nil {
				t.Errorf("Parse() returned AST with nil root")
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	testContent := `Test Mathematics Study Guide
Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations

1. What is a linear equation? - An equation where the highest power of the variable is 1.
2. How do you solve 2x + 3 = 7? - Subtract 3 from both sides: 2x = 4, then divide by 2: x = 2.

Learn More: See Khan Academy's linear equations course.`
	tmpFile, err := os.CreateTemp("", "test_study_guide_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	result, err := ParseFile(tmpFile.Name(), config.NewMetadata("test_parser"))
	if err != nil {
		t.Errorf("ParseFile() unexpected error: %v", err)
		return
	}
	if !result.Success {
		t.Errorf("ParseFile() unexpected errors: %v", result.Errors)
		return
	}
	if result.AST == nil {
		t.Errorf("ParseFile() returned nil AST")
		return
	}
	if result.AST.Root == nil {
		t.Errorf("ParseFile() returned AST with nil root")
	}

}

func TestMetadata(t *testing.T) {
	tests := []struct {
		name       string
		parserType string
		expected   string
	}{
		{"Colleges", "colleges", "colleges"},
		{"APExams", "ap_exams", "ap_exams"},
		{"Certifications", "certifications", "certifications"},
		{"DOD", "dod", "dod"},
		{"EntranceExams", "entrance_exams", "entrance_exams"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := config.NewMetadata(tt.parserType)
			if metadata.Type != tt.expected {
				t.Errorf("Metadata Type = %s, want %s", metadata.Type, tt.expected)
			}
		})
	}
}

func TestMetadataWithOption(t *testing.T) {
	metadata := config.NewMetadata("test_parser").WithOption("strict", "true").WithOption("debug", "false")

	if metadata.Type != "test_parser" {
		t.Errorf("Expected parser type 'test_parser', got %s", metadata.Type)
	}

	if metadata.Options["strict"] != "true" {
		t.Errorf("Expected option 'strict' to be 'true', got %s", metadata.Options["strict"])
	}

	if metadata.Options["debug"] != "false" {
		t.Errorf("Expected option 'debug' to be 'false', got %s", metadata.Options["debug"])
	}
}

func TestPreparseReturnsOnlyLexerErrors(t *testing.T) {
	lines := []string{
		"\x00\x01\x02Invalid binary data", // Should trigger a lexer error
		"Another line",                    // Should not be reached by preparser
	}

	result, err := Preparse(lines, config.NewMetadata("test_parser"))
	if err != nil {
		t.Fatalf("Preparse() returned unexpected error: %v", err)
	}

	if result.Success {
		t.Error("Preparse() should return Success=false when there are lexer errors")
	}

	if len(result.Errors) == 0 {
		t.Error("Preparse() should return lexer errors")
	}

	if len(result.Tokens) != 0 {
		t.Errorf("Preparse() should not return tokens when there are lexer errors, got %d", len(result.Tokens))
	}

	for _, err := range result.Errors {
		if err.LineNumber != 1 {
			t.Errorf("Expected lexer error to be on line 1, got line %d", err.LineNumber)
		}
	}
}

func TestPreparseReturnsOnlyPreparserErrors(t *testing.T) {
	lines := []string{
		"Study Guide",                      // Valid file header
		"Section: Chapter 1",               // Valid header
		"Learn More: ",                     // Invalid - empty after colon
		"1. What is Go? - Answer",          // Valid question
		"Passage: This is a valid passage", // Valid passage
		"Some content here",                // Valid content
	}

	result, err := Preparse(lines, config.NewMetadata("test_parser"))
	if err != nil {
		t.Fatalf("Preparse() returned unexpected error: %v", err)
	}

	if result.Success {
		t.Error("Preparse() should return Success=false when there are preparser errors")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Preparse() should collect one preparser error, got %d: %v", len(result.Errors), result.Errors)
	}

	if !strings.Contains(result.Errors[0].Message, "learn more line must contain text after 'Learn More:'") {
		t.Error("Expected error about learn more line missing text not found")
	}

	if len(result.Tokens) != 5 {
		t.Errorf("Expected 5 valid tokens, got %d", len(result.Tokens))
	}
}

func TestParseFileWithLexerError(t *testing.T) {
	testContent := "\x00\x01\x02Invalid binary data\nAnother line"
	tmpFile, err := os.CreateTemp("", "test_study_guide_lexerr_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	result, err := ParseFile(tmpFile.Name(), config.NewMetadata("test_parser"))
	if err != nil {
		t.Errorf("ParseFile() unexpected error: %v", err)
		return
	}
	if result.Success {
		t.Error("ParseFile() should return Success=false for file with lexer error")
	} else if len(result.Errors) == 0 {
		t.Error("ParseFile() should return lexer errors")
	} else {
		found := false
		for _, parseErr := range result.Errors {
			if parseErr.LineNumber == 1 {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected lexer error to be on line 1")
		}
	}
}
