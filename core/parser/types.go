package parser

import (
	"github.com/studyguides-com/study-guides-parser/core/lexer"
)

// ParserType represents the type of parser being used
type ParserType string

// Node represents a single node in the parser syntax tree
type Node struct {
	Type     lexer.TokenType `json:"type"`
	Data     interface{}     `json:"data,omitempty"` // nullable
	Children []*Node         `json:"children,omitempty"`
	Parent   *Node           `json:"-"` // already nullable
}

// AbstractSyntaxTree represents the output of a parser tree
type AbstractSyntaxTree struct {
	ParserType ParserType `json:"parser_type"`
	Timestamp  string     `json:"timestamp"`
	Root       *Node      `json:"root"`
}
