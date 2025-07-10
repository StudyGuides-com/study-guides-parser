package tree

import (
	"github.com/studyguides-com/study-guides-parser/core/ontology"
)

// TreeQAble defines the interface for types that can be QA'd
// The visitor function is called for each tag with its depth
// Traverse must visit all tags in the tree
// GetQAResults/SetQAResults manage QA results
type TreeQAble interface {
	Traverse(visitor func(TagQATarget, int))
	GetQAResults() QAResults
	SetQAResults(results QAResults)
}

// TreeTraverser defines the interface for traversing tree structures
type TreeTraverser interface {
	// Traverse performs a depth-first traversal of the tree
	// The visitor function is called for each tag with its depth
	Traverse(visitor func(TagTypeAssignable, int))
	
	// TraverseWithContext performs traversal with additional context
	TraverseWithContext(visitor func(TagTypeAssignable, int, ontology.ContextType))
}

// TagTypeAssigner defines the interface for assigning tag types
type TagTypeAssigner interface {
	AssignTagTypes(contextType ontology.ContextType) error
}

// TraverseForTagTypes implements TreeTraverser interface for tag type assignment
func (t *Tree) TraverseForTagTypes(visitor func(TagTypeAssignable, int)) {
	if t.Root == nil {
		return
	}
	
	var traverse func(*Tag, int)
	traverse = func(tag *Tag, depth int) {
		if tag == nil {
			return
		}
		
		// Visit current tag
		visitor(tag, depth)
		
		// Recursively visit children
		for _, child := range tag.ChildTags {
			traverse(child, depth+1)
		}
	}
	
	// Start traversal from root-level tags
	for _, child := range t.Root.ChildTags {
		traverse(child, 1)
	}
}

// Traverse implements TreeQAble interface for QA
func (t *Tree) Traverse(visitor func(TagQATarget, int)) {
	if t.Root == nil {
		return
	}
	var traverse func(*Tag, int)
	traverse = func(tag *Tag, depth int) {
		if tag == nil {
			return
		}
		visitor(tag, depth)
		for _, child := range tag.ChildTags {
			traverse(child, depth+1)
		}
	}
	for _, child := range t.Root.ChildTags {
		traverse(child, 1)
	}
}

// TraverseWithContext implements TreeTraverser interface with context
func (t *Tree) TraverseWithContext(visitor func(TagTypeAssignable, int, ontology.ContextType)) {
	if t.Root == nil {
		return
	}
	
	// Get context from metadata
	contextType := ontology.ContextTypeNone
	if t.Metadata != nil {
		contextType = t.Metadata.ContextType
	}
	
	var traverse func(*Tag, int)
	traverse = func(tag *Tag, depth int) {
		if tag == nil {
			return
		}
		
		// Visit current tag with context
		visitor(tag, depth, contextType)
		
		// Recursively visit children
		for _, child := range tag.ChildTags {
			traverse(child, depth+1)
		}
	}
	
	// Start traversal from root-level tags
	for _, child := range t.Root.ChildTags {
		traverse(child, 1)
	}
} 