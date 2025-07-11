package tree

import (
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
)

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
