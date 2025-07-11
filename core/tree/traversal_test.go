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
