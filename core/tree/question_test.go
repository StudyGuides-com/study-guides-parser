package tree

import (
	"encoding/json"
	"testing"
)

func TestNewQuestionWithNilDistractors(t *testing.T) {
	// Test with nil distractors
	question := NewQuestion("What is 1 + 1?", "2", nil, "Simple addition", 1)

	jsonData, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("Failed to marshal question: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify distractors is an empty array instead of null
	if !contains(jsonStr, `"distractors":[]`) {
		t.Errorf("Expected distractors to be empty array, got: %s", jsonStr)
	}

	// Verify we don't have null
	if contains(jsonStr, `"distractors":null`) {
		t.Errorf("Found null for distractors, expected empty array")
	}
}

func TestNewQuestionWithEmptyDistractors(t *testing.T) {
	// Test with empty slice
	question := NewQuestion("What is 1 + 1?", "2", []string{}, "Simple addition", 1)

	jsonData, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("Failed to marshal question: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify distractors is an empty array
	if !contains(jsonStr, `"distractors":[]`) {
		t.Errorf("Expected distractors to be empty array, got: %s", jsonStr)
	}
}

func TestNewQuestionWithDistractors(t *testing.T) {
	// Test with actual distractors
	distractors := []string{"3", "4", "5"}
	question := NewQuestion("What is 1 + 1?", "2", distractors, "Simple addition", 1)

	jsonData, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("Failed to marshal question: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify distractors contains the expected values
	if !contains(jsonStr, `"distractors":["3","4","5"]`) {
		t.Errorf("Expected distractors to contain values, got: %s", jsonStr)
	}
}

func TestNewQuestionOrder(t *testing.T) {
	// Test that order is correctly set
	q1 := NewQuestion("Question 1?", "A", nil, "", 1)
	q2 := NewQuestion("Question 2?", "B", nil, "", 2)
	q3 := NewQuestion("Question 3?", "C", nil, "", 5)

	if q1.Order != 1 {
		t.Errorf("Expected q1.Order to be 1, got %d", q1.Order)
	}
	if q2.Order != 2 {
		t.Errorf("Expected q2.Order to be 2, got %d", q2.Order)
	}
	if q3.Order != 5 {
		t.Errorf("Expected q3.Order to be 5, got %d", q3.Order)
	}
}

func TestQuestionOrderInJSON(t *testing.T) {
	// Test that order is included in JSON serialization
	question := NewQuestion("What is 2 + 2?", "4", nil, "", 3)

	jsonData, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("Failed to marshal question: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify order is in the JSON output
	if !contains(jsonStr, `"order":3`) {
		t.Errorf("Expected order:3 in JSON, got: %s", jsonStr)
	}
}
