package tree

import (
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
)

func TestTreeLeafNodes(t *testing.T) {
	// Create metadata
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeCertifications,
	}
	tree := NewTree(metadata)

	// Create a complex tree structure with multiple branches
	// Structure:
	// Category
	// ├── Agency1
	// │   ├── Certification1
	// │   │   └── Topic1 (leaf)
	// │   └── Topic2 (leaf)
	// └── Agency2
	//     └── Topic3 (leaf)

	category := NewTag("Category")
	agency1 := NewTag("Agency1")
	agency2 := NewTag("Agency2")
	certification1 := NewTag("Certification1")
	topic1 := NewTag("Topic1")
	topic2 := NewTag("Topic2")
	topic3 := NewTag("Topic3")

	// Build the hierarchy
	category.AddChildTag(agency1)
	category.AddChildTag(agency2)
	agency1.AddChildTag(certification1)
	agency1.AddChildTag(topic2)
	certification1.AddChildTag(topic1)
	agency2.AddChildTag(topic3)

	tree.Root.AddChildTag(category)

	// Get leaf nodes
	leafNodes := tree.LeafNodes()

	// Expected leaf nodes: Topic1, Topic2, Topic3
	expectedLeafTitles := []string{"Topic1", "Topic2", "Topic3"}

	// Check the number of leaf nodes
	if len(leafNodes) != len(expectedLeafTitles) {
		t.Errorf("Expected %d leaf nodes, got %d", len(expectedLeafTitles), len(leafNodes))
	}

	// Check that all returned nodes are actually leaf nodes (no children)
	for i, leaf := range leafNodes {
		if len(leaf.ChildTags) != 0 {
			t.Errorf("Leaf node %d (%s) has %d children, expected 0", i, leaf.Title, len(leaf.ChildTags))
		}
	}

	// Check that we got the expected leaf node titles
	foundTitles := make(map[string]bool)
	for _, leaf := range leafNodes {
		foundTitles[leaf.Title] = true
	}

	for _, expectedTitle := range expectedLeafTitles {
		if !foundTitles[expectedTitle] {
			t.Errorf("Expected leaf node with title '%s' not found", expectedTitle)
		}
	}

	// Verify that non-leaf nodes are not included
	nonLeafTitles := []string{"Category", "Agency1", "Agency2", "Certification1"}
	for _, nonLeafTitle := range nonLeafTitles {
		if foundTitles[nonLeafTitle] {
			t.Errorf("Non-leaf node with title '%s' was incorrectly included in leaf nodes", nonLeafTitle)
		}
	}
}

func TestTreeLeafNodesEmptyTree(t *testing.T) {
	// Test with an empty tree
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeCertifications,
	}
	tree := NewTree(metadata)

	leafNodes := tree.LeafNodes()

	// Empty tree should return empty slice
	if len(leafNodes) != 0 {
		t.Errorf("Expected 0 leaf nodes for empty tree, got %d", len(leafNodes))
	}
}

func TestTreeLeafNodesSingleNode(t *testing.T) {
	// Test with a tree containing only one node (which should be a leaf)
	metadata := &config.Metadata{
		Type:        "test",
		ContextType: ontology.ContextTypeCertifications,
	}
	tree := NewTree(metadata)

	singleTag := NewTag("SingleTag")
	tree.Root.AddChildTag(singleTag)

	leafNodes := tree.LeafNodes()

	// Should have exactly one leaf node
	if len(leafNodes) != 1 {
		t.Errorf("Expected 1 leaf node, got %d", len(leafNodes))
	}

	if leafNodes[0].Title != "SingleTag" {
		t.Errorf("Expected leaf node title 'SingleTag', got '%s'", leafNodes[0].Title)
	}
}
