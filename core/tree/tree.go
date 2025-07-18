package tree

import (
	"github.com/studyguides-com/study-guides-parser/core/config"
)

type Tree struct {
	Root     *Root            `json:"root"`
	Metadata *config.Metadata `json:"metadata"`
}

func NewTree(metadata *config.Metadata) *Tree {
	return &Tree{
		Metadata: metadata,
		Root:     NewRoot(),
	}
}

func (t *Tree) LeafNodes() []*Tag {
	leafNodes := []*Tag{}
	
	// Use the existing Traverse method to walk the tree
	t.Traverse(func(tag TagQATarget, depth int) {
		// Type assert to *Tag to access ChildTags
		if tagTag, ok := tag.(*Tag); ok {
			// Check if this is a leaf node (no children)
			if len(tagTag.ChildTags) == 0 {
				leafNodes = append(leafNodes, tagTag)
			}
		}
	})
	
	return leafNodes
}
