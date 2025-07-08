package processor

import (
	"os"
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/types"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		metadata *types.Metadata
		wantErr  bool
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
			metadata: types.NewMetadata("test_parser"),
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
			metadata: types.NewMetadata("ap_test_parser"),
			wantErr:  false,
		},
		{
			name: "missing file header",
			lines: []string{
				"Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations",
				"1. What is x? - A variable",
			},
			metadata: types.NewMetadata("test_parser"),
			wantErr:  true,
		},
		{
			name:     "empty lines",
			lines:    []string{},
			metadata: types.NewMetadata("test_parser"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.lines, tt.metadata)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}
			if ast == nil {
				t.Errorf("Parse() returned nil AST")
				return
			}
			if ast.Root == nil {
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
	ast, err := ParseFile(tmpFile.Name(), types.NewMetadata("test_parser"))
	if err != nil {
		t.Errorf("ParseFile() unexpected error: %v", err)
		return
	}
	if ast == nil {
		t.Errorf("ParseFile() returned nil AST")
		return
	}
	if ast.Root == nil {
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
					metadata := types.NewMetadata(tt.parserType)
		if metadata.Type != tt.expected {
			t.Errorf("Metadata Type = %s, want %s", metadata.Type, tt.expected)
		}
		})
	}
}

func TestMetadataWithOption(t *testing.T) {
	metadata := types.NewMetadata("test_parser").WithOption("strict", "true").WithOption("debug", "false")
	
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