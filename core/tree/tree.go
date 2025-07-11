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
