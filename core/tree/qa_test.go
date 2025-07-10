package tree

import (
	"testing"
)

func TestNewQAResult(t *testing.T) {
	// Test creating a passing QA result
	passingResult := NewQAResult("Test QA", true)
	
	if passingResult.Name != "Test QA" {
		t.Errorf("Expected name 'Test QA', got '%s'", passingResult.Name)
	}
	
	if !passingResult.Passed {
		t.Errorf("Expected passed to be true, got false")
	}
	
	if passingResult.Warnings == nil {
		t.Errorf("Expected warnings to be an empty slice, got nil")
	}
	
	if len(passingResult.Warnings) != 0 {
		t.Errorf("Expected empty warnings slice, got %d warnings", len(passingResult.Warnings))
	}
	
	// Test creating a failing QA result
	failingResult := NewQAResult("Failing QA", false)
	
	if failingResult.Name != "Failing QA" {
		t.Errorf("Expected name 'Failing QA', got '%s'", failingResult.Name)
	}
	
	if failingResult.Passed {
		t.Errorf("Expected passed to be false, got true")
	}
	
	if failingResult.Warnings == nil {
		t.Errorf("Expected warnings to be an empty slice, got nil")
	}
	
	if len(failingResult.Warnings) != 0 {
		t.Errorf("Expected empty warnings slice, got %d warnings", len(failingResult.Warnings))
	}
}

func TestNewQAResults(t *testing.T) {
	// Test creating new QA results
	qaResults := NewQAResults()
	
	if !qaResults.OverallPassed {
		t.Errorf("Expected overall passed to be true by default, got false")
	}
	
	if qaResults.Results == nil {
		t.Errorf("Expected results to be an empty slice, got nil")
	}
	
	if len(qaResults.Results) != 0 {
		t.Errorf("Expected empty results slice, got %d results", len(qaResults.Results))
	}
}

func TestQAResultWithWarnings(t *testing.T) {
	// Test creating a QA result and then adding warnings
	result := NewQAResult("Warning Test", false)
	
	// Add some warnings
	result.Warnings = append(result.Warnings, "Warning 1")
	result.Warnings = append(result.Warnings, "Warning 2")
	
	if len(result.Warnings) != 2 {
		t.Errorf("Expected 2 warnings, got %d", len(result.Warnings))
	}
	
	if result.Warnings[0] != "Warning 1" {
		t.Errorf("Expected first warning 'Warning 1', got '%s'", result.Warnings[0])
	}
	
	if result.Warnings[1] != "Warning 2" {
		t.Errorf("Expected second warning 'Warning 2', got '%s'", result.Warnings[1])
	}
} 