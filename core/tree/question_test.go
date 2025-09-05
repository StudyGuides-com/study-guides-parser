package tree

import (
	"encoding/json"
	"testing"
)

func TestNewQuestionWithNilDistractors(t *testing.T) {
	// Test with nil distractors
	question := NewQuestion("What is 1 + 1?", "2", nil, "Simple addition")

	jsonData, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("Failed to marshal question: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify distractors is an empty array instead of null
	if !contains(jsonStr, `"distractor":[]`) {
		t.Errorf("Expected distractors to be empty array, got: %s", jsonStr)
	}

	// Verify we don't have null
	if contains(jsonStr, `"distractor":null`) {
		t.Errorf("Found null for distractors, expected empty array")
	}
}

func TestNewQuestionWithEmptyDistractors(t *testing.T) {
	// Test with empty slice
	question := NewQuestion("What is 1 + 1?", "2", []string{}, "Simple addition")

	jsonData, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("Failed to marshal question: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify distractors is an empty array
	if !contains(jsonStr, `"distractor":[]`) {
		t.Errorf("Expected distractors to be empty array, got: %s", jsonStr)
	}
}

func TestNewQuestionWithDistractors(t *testing.T) {
	// Test with actual distractors
	distractors := []string{"3", "4", "5"}
	question := NewQuestion("What is 1 + 1?", "2", distractors, "Simple addition")

	jsonData, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("Failed to marshal question: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify distractors contains the expected values
	if !contains(jsonStr, `"distractor":["3","4","5"]`) {
		t.Errorf("Expected distractors to contain values, got: %s", jsonStr)
	}
}
