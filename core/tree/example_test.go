package tree

import (
	"fmt"
	"strings"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
)

// Example demonstrates the improved interface design
func Example() {
	// Create metadata with context type
	metadata := &config.Metadata{
		Type:        "example",
		ContextType: ontology.ContextTypeCertifications,
	}

	// Create a tree
	tree := NewTree(metadata)

	// Add some tags
	tag1 := NewTag("Category")
	tag2 := NewTag("Agency")
	tag3 := NewTag("Certification")
	tag4 := NewTag("Topic")

	// Build hierarchy
	tag1.AddChildTag(tag2)
	tag2.AddChildTag(tag3)
	tag3.AddChildTag(tag4)
	tree.Root.AddChildTag(tag1)

	// Traverse the tree and print tag information
	fmt.Println("Tree structure:")
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		indent := strings.Repeat("  ", depth-1)
		fmt.Printf("%s- %s (Type: %s, Context: %s)\n",
			indent, tag.(*Tag).Title, tag.GetTagType(), tag.GetContext())
	})

	// Print the tree as JSON
	fmt.Println("\nTree as JSON:")
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		indent := strings.Repeat("  ", depth-1)
		fmt.Printf("%s- %s (Type: %s, Context: %s)\n",
			indent, tag.(*Tag).Title, tag.GetTagType(), tag.GetContext())
	})

	// Demonstrate the TagTypeAssignable interface
	fmt.Println("Before assignment:")
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		fmt.Printf("  Tag: %s, Depth: %d, Type: %s, Context: %s\n",
			tag.(*Tag).Title, depth, tag.GetTagType(), tag.GetContext())
	})

	// Assign tag types
	tree.AssignTagTypes(ontology.ContextTypeCertifications)

	fmt.Println("\nAfter assignment:")
	tree.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		fmt.Printf("  Tag: %s, Depth: %d, Type: %s, Context: %s\n",
			tag.(*Tag).Title, depth, tag.GetTagType(), tag.GetContext())
	})

	// Demonstrate manual tag type assignment
	fmt.Println("\nManual assignment example:")
	tag1.SetTagType(ontology.TagTypeCategory)
	tag1.SetContext(ontology.ContextTypeAPExams)
	fmt.Printf("  Tag1 Type: %s, Context: %s\n", tag1.GetTagType(), tag1.GetContext())
}

// ExampleTagTypeAssignable demonstrates the TagTypeAssignable interface
func ExampleTagTypeAssignable() {
	// Create a tag
	tag := NewTag("Example")

	// Use the interface methods
	tag.SetTagType(ontology.TagTypeTopic)
	tag.SetContext(ontology.ContextTypeCertifications)

	fmt.Printf("Tag Type: %s\n", tag.GetTagType())
	fmt.Printf("Tag Context: %s\n", tag.GetContext())

	// The interface provides encapsulation and consistency
	var assignable TagTypeAssignable = tag
	assignable.SetTagType(ontology.TagTypeModule)
	fmt.Printf("Updated Tag Type: %s\n", assignable.GetTagType())
}
