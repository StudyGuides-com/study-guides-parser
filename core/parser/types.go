package parser

import (
	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/lexer"
)

// Node represents a single node in the parser syntax tree
type Node struct {
	Type     lexer.TokenType `json:"type"`
	Data     interface{}     `json:"data,omitempty"` // nullable
	Children []*Node         `json:"children,omitempty"`
	Parent   *Node           `json:"-"` // already nullable
}

// AbstractSyntaxTree represents the output of a parser tree
type AbstractSyntaxTree struct {
	Metadata *config.Metadata `json:"metadata"`
	Timestamp string          `json:"timestamp"`
	Root      *Node           `json:"root"`
}
