package processor

import (
	"os"
	"testing"

	"github.com/StudyGuides-com/study-guides-parser/core/parser"
)

func TestParseLines(t *testing.T) {
	tests := []struct {
		name       string
		lines      []string
		parserType ParserType
		wantErr    bool
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
			parserType: Colleges,
			wantErr:    false,
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
			parserType: APExams,
			wantErr:    false,
		},
		{
			name: "missing file header",
			lines: []string{
				"Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations",
				"1. What is x? - A variable",
			},
			parserType: Colleges,
			wantErr:    true,
		},
		{
			name:       "empty lines",
			lines:      []string{},
			parserType: Colleges,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := ParseLines(tt.lines, tt.parserType)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseLines() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("ParseLines() unexpected error: %v", err)
				return
			}
			
			if ast == nil {
				t.Errorf("ParseLines() returned nil AST")
				return
			}
			
			// Basic validation of the AST
			if ast.Root == nil {
				t.Errorf("ParseLines() returned AST with nil root")
			}
			
			if ast.ParserType != parser.ParserType(tt.parserType) {
				t.Errorf("ParseLines() returned AST with wrong parser type: got %s, want %s", ast.ParserType, tt.parserType)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	// Create a temporary test file
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

	// Test parsing the file
	ast, err := ParseFile(tmpFile.Name(), Colleges)
	if err != nil {
		t.Errorf("ParseFile() unexpected error: %v", err)
		return
	}

	if ast == nil {
		t.Errorf("ParseFile() returned nil AST")
		return
	}

	// Basic validation
	if ast.Root == nil {
		t.Errorf("ParseFile() returned AST with nil root")
	}

	if ast.ParserType != "colleges" {
		t.Errorf("ParseFile() returned AST with wrong parser type: got %s, want colleges", ast.ParserType)
	}
}

func TestParseTokens(t *testing.T) {
	// This test would require creating lexer.LineInfo tokens
	// For now, we'll test that the function exists and can be called
	// A more comprehensive test would be added when we have token conversion utilities
	
	t.Run("function exists", func(t *testing.T) {
		// This is a placeholder test to ensure the function exists
		// In a real implementation, we would test with actual tokens
		_ = ParseTokens
	})
}

func TestParserTypeConstants(t *testing.T) {
	tests := []struct {
		name       string
		parserType ParserType
		expected   string
	}{
		{"Colleges", Colleges, "colleges"},
		{"APExams", APExams, "ap_exams"},
		{"Certifications", Certifications, "certifications"},
		{"DOD", DOD, "dod"},
		{"EntranceExams", EntranceExams, "entrance_exams"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.parserType) != tt.expected {
				t.Errorf("ParserType %s = %s, want %s", tt.name, tt.parserType, tt.expected)
			}
		})
	}
} 