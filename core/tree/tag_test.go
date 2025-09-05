package tree

import (
	"encoding/json"
	"testing"
)

func TestTagJSONSerialization(t *testing.T) {
	// Test NewTag function
	tag := NewTag("Test Tag")

	jsonData, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to marshal tag: %v", err)
	}

	// Check that the JSON contains empty arrays instead of null
	jsonStr := string(jsonData)

	// Verify content_descriptors is an empty array
	if !contains(jsonStr, `"content_descriptors":[]`) {
		t.Errorf("Expected content_descriptors to be empty array, got: %s", jsonStr)
	}

	// Verify meta_tags is an empty array
	if !contains(jsonStr, `"meta_tags":[]`) {
		t.Errorf("Expected meta_tags to be empty array, got: %s", jsonStr)
	}

	// Verify child_tags is an empty array
	if !contains(jsonStr, `"child_tags":[]`) {
		t.Errorf("Expected child_tags to be empty array, got: %s", jsonStr)
	}
}

func TestTagWithParentJSONSerialization(t *testing.T) {
	// Test NewTagWithParent function
	tag := NewTagWithParent("Test Tag", "Parent Tag")

	jsonData, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to marshal tag: %v", err)
	}

	// Check that the JSON contains empty arrays instead of null
	jsonStr := string(jsonData)

	// Verify content_descriptors is an empty array
	if !contains(jsonStr, `"content_descriptors":[]`) {
		t.Errorf("Expected content_descriptors to be empty array, got: %s", jsonStr)
	}

	// Verify meta_tags is an empty array
	if !contains(jsonStr, `"meta_tags":[]`) {
		t.Errorf("Expected meta_tags to be empty array, got: %s", jsonStr)
	}

	// Verify child_tags is an empty array
	if !contains(jsonStr, `"child_tags":[]`) {
		t.Errorf("Expected child_tags to be empty array, got: %s", jsonStr)
	}
}

func TestRootJSONSerialization(t *testing.T) {
	// Test NewRoot function
	root := NewRoot()

	jsonData, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("Failed to marshal root: %v", err)
	}

	// Check that the JSON contains empty arrays instead of null
	jsonStr := string(jsonData)

	// Verify warnings is an empty array
	if !contains(jsonStr, `"warnings":[]`) {
		t.Errorf("Expected warnings to be empty array, got: %s", jsonStr)
	}

	// Verify child_tags is an empty array
	if !contains(jsonStr, `"child_tags":[]`) {
		t.Errorf("Expected child_tags to be empty array, got: %s", jsonStr)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
