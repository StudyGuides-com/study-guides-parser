package tree

import (
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
)

func TestTreeTraversal(t *testing.T) {
	// Create a simple tree structure
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeCertifications,
	}
	tree := NewTree(metadata)

	// Add some test tags
	tag1 := NewTag("Category")
	tag2 := NewTag("Agency")
	tag3 := NewTag("Certification")
	tag4 := NewTag("Topic")

	// Build hierarchy: Category -> Agency -> Certification -> Topic
	tag1.AddChildTag(tag2)
	tag2.AddChildTag(tag3)
	tag3.AddChildTag(tag4)
	tree.Root.AddChildTag(tag1)

	// Test traversal
	var visited []string
	var depths []int
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		visited = append(visited, tag.(*Tag).Title)
		depths = append(depths, depth)
	})

	expectedTags := []string{"Category", "Agency", "Certification", "Topic"}
	expectedDepths := []int{1, 2, 3, 4}

	if len(visited) != len(expectedTags) {
		t.Errorf("Expected %d visited tags, got %d", len(expectedTags), len(visited))
	}

	for i, tag := range expectedTags {
		if visited[i] != tag {
			t.Errorf("Expected tag %s at position %d, got %s", tag, i, visited[i])
		}
	}

	for i, depth := range expectedDepths {
		if depths[i] != depth {
			t.Errorf("Expected depth %d at position %d, got %d", depth, i, depths[i])
		}
	}
}

func TestTagTypeAssignment(t *testing.T) {
	// Create a tree with certification context
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeCertifications,
	}
	tree := NewTree(metadata)

	// Add tags with depth 4 structure (matching ontology data)
	tag1 := NewTag("Category")
	tag2 := NewTag("Agency")
	tag3 := NewTag("Certification")
	tag4 := NewTag("Topic")

	// Build hierarchy: Category -> Agency -> Certification -> Topic
	tag1.AddChildTag(tag2)
	tag2.AddChildTag(tag3)
	tag3.AddChildTag(tag4)
	tree.Root.AddChildTag(tag1)

	// Debug: Print before assignment
	t.Log("Before assignment:")
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		t.Logf("Tag: %s, Depth: %d, Type: %s, Context: %s", tag.(*Tag).Title, depth, tag.GetTagType(), tag.GetContext())
	})

	// Assign tag types using the total depth of 4
	if err := tree.AssignTagTypes(ontology.ContextTypeCertifications); err != nil {
		t.Fatalf("Failed to assign tag types: %v", err)
	}

	// Debug: Print after assignment
	t.Log("After assignment:")
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		t.Logf("Tag: %s, Depth: %d, Type: %s, Context: %s", tag.(*Tag).Title, depth, tag.GetTagType(), tag.GetContext())
	})

	// Verify tag types were assigned correctly
	// For depth 4, the ontology should assign: Category, Certifying_Agency, Certification, Topic
	expectedTypes := []ontology.TagType{
		ontology.TagTypeCategory,
		ontology.TagTypeCertifyingAgency,
		ontology.TagTypeCertification,
		ontology.TagTypeTopic,
	}

	// Collect actual tag types in order
	var actualTypes []ontology.TagType
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		actualTypes = append(actualTypes, tag.GetTagType())
	})

	if len(actualTypes) != len(expectedTypes) {
		t.Errorf("Expected %d tag types, got %d", len(expectedTypes), len(actualTypes))
		return
	}

	for i, expectedType := range expectedTypes {
		if actualTypes[i] != expectedType {
			t.Errorf("Expected tag type %s at position %d, got %s", expectedType, i, actualTypes[i])
		}
	}

	// Verify context was set
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		if tag.GetContext() != ontology.ContextTypeCertifications {
			t.Errorf("Expected context %s, got %s", ontology.ContextTypeCertifications, tag.GetContext())
		}
	})
}

func TestTraverseWithContext(t *testing.T) {
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeAPExams,
	}
	tree := NewTree(metadata)

	tag1 := NewTag("Category")
	tag2 := NewTag("AP Exam")
	tag3 := NewTag("Topic")

	tag1.AddChildTag(tag2)
	tag2.AddChildTag(tag3)
	tree.Root.AddChildTag(tag1)

	var contexts []ontology.ContextType
	tree.TraverseWithContext(func(tag TagTypeAssignable, depth int, context ontology.ContextType) {
		contexts = append(contexts, context)
	})

	// All tags should have the same context
	for i, context := range contexts {
		if context != ontology.ContextTypeAPExams {
			t.Errorf("Expected context %s at position %d, got %s", ontology.ContextTypeAPExams, i, context)
		}
	}
}

func TestAssignTagTypesError(t *testing.T) {
	// Create a tree with a depth that doesn't exist in the ontology
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeCertifications,
	}
	tree := NewTree(metadata)

	// Add tags with depth 2 structure (which doesn't exist for Certifications)
	tag1 := NewTag("Category")
	tag2 := NewTag("Topic")

	tag1.AddChildTag(tag2)
	tree.Root.AddChildTag(tag1)

	// This should return an error because there's no ontology for Certifications with depth 2
	err := tree.AssignTagTypes(ontology.ContextTypeCertifications)
	if err == nil {
		t.Error("Expected error when no ontology is found, but got nil")
	}

	expectedError := "no ontology found for context type 'Certifications' with depth 2"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
} 