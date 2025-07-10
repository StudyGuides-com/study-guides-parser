package tree

import (
	"fmt"

	"github.com/studyguides-com/study-guides-parser/core/ontology"
)

// AssignTagTypes implements TagTypeAssigner interface
func (t *Tree) AssignTagTypes(contextType ontology.ContextType) error {
	// First, determine the maximum depth of the tree
	maxDepth := t.getMaxDepth()
	
	// Find the ontology entry for this total depth
	tagOntology := ontology.FindTagOntology(contextType, maxDepth)
	if tagOntology == nil {
		return fmt.Errorf("no ontology found for context type '%s' with depth %d", contextType, maxDepth)
	}
	
	// Now traverse and assign types based on individual tag depths
	t.TraverseForTagTypes(func(tag TagTypeAssignable, depth int) {
		assignTagTypeFromOntology(tag, contextType, depth, tagOntology)
	})
	
	return nil
}

// getMaxDepth calculates the maximum depth of the tree
func (t *Tree) getMaxDepth() int {
	if t.Root == nil {
		return 0
	}
	
	var maxDepth int
	var traverse func(*Tag, int)
	traverse = func(tag *Tag, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		for _, child := range tag.ChildTags {
			traverse(child, depth+1)
		}
	}
	
	for _, child := range t.Root.ChildTags {
		traverse(child, 1)
	}
	
	return maxDepth
}

// assignTagTypeFromOntology assigns the appropriate tag type based on context and depth within a known ontology
func assignTagTypeFromOntology(tag TagTypeAssignable, contextType ontology.ContextType, depth int, tagOntology *ontology.TagOntology) {
	if depth <= len(tagOntology.TagTypes) {
		tag.SetTagType(tagOntology.TagTypes[depth-1]) // depth is 1-indexed, slice is 0-indexed
		tag.SetContext(contextType)
	}
} 