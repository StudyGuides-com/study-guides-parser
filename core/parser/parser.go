package parser

import (
	"fmt"
	"time"

	"github.com/StudyGuides-com/study-guides-parser/core/lexer"
	"github.com/StudyGuides-com/study-guides-parser/core/preparser"
)

type Parser struct {
	Root    *Node // The root node representing the entire tree (now it will be the file header)
	Current *Node // The current "open" node (e.g., a header or passage)
	Lines   []preparser.ParsedLineInfo
}

func NewParser(lines []preparser.ParsedLineInfo) *Parser {
	return &Parser{
		Lines: lines,
	}
}

// Helper to add node under current node if current is of expected type
func (p *Parser) addUnderCurrent(expected lexer.TokenType, line preparser.ParsedLineInfo) *ParserError {
	if p.Current == nil {
		return NewParserError(CodeValidation, fmt.Sprintf(" unexpected %s under <nil>", line.Type), line)
	}
	if p.Current.Type != expected {
		return NewParserError(CodeValidation, fmt.Sprintf(" unexpected %s under %s", line.Type, p.Current.Type), line)
	}
	node := &Node{
		Type:     line.Type,
		Data:     line.ParsedValue,
		Children: []*Node{},
		Parent:   p.Current,
	}
	p.Current.Children = append(p.Current.Children, node)
	return nil
}

// Helper to walk up the tree to find nearest parent of a certain type
func (p *Parser) findNearest(target lexer.TokenType) *Node {
	node := p.Current
	for node != nil && node.Type != target {
		node = node.Parent
	}
	return node
}

// finalize returns the parsed content as an AbstractSyntaxTree
func (p *Parser) finalize(parserType ParserType) (*AbstractSyntaxTree, *ParserError) {
	if p.Root == nil {
		return nil, NewParserError(CodeValidation, "no root node found", preparser.ParsedLineInfo{})
	}

	output := &AbstractSyntaxTree{
		ParserType: parserType,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Root:       p.Root,
	}

	return output, nil
}

// see internal/services/parser/gramar.ebnf for the grammar
func (p *Parser) Parse(parserType ParserType) (*AbstractSyntaxTree, *ParserError) {
	if len(p.Lines) == 0 {
		return nil, NewParserError(CodeValidation, "no lines to parse", preparser.ParsedLineInfo{})
	}

	firstLine := p.Lines[0]
	if firstLine.Type != lexer.TokenTypeFileHeader {
		return nil, NewParserError(CodeValidation, "first line must be a file header", firstLine)
	}

	// Initialize root node as file header
	p.Root = &Node{
		Type:     lexer.TokenTypeFileHeader,
		Data:     firstLine.ParsedValue, // This is already a *preparser.FileHeaderResult
		Children: []*Node{},
	}
	p.Current = p.Root

	// Process the remaining lines
	for _, line := range p.Lines[1:] {
		switch line.Type {

		// Header
		case lexer.TokenTypeHeader:
			node := &Node{
				Type:     lexer.TokenTypeHeader,
				Data:     line.ParsedValue,
				Children: []*Node{},
			}
			p.Root.Children = append(p.Root.Children, node)
			node.Parent = p.Root
			p.Current = node

		// Passage
		case lexer.TokenTypePassage:
			// Find the nearest header
			parent := p.findNearest(lexer.TokenTypeHeader)
			if parent == nil {
				return nil, NewParserError(CodeValidation, fmt.Sprintf("%s without parent %s", line.Type, lexer.TokenTypeHeader), line)
			}
			node := &Node{
				Type:     lexer.TokenTypePassage,
				Data:     line.ParsedValue,
				Children: []*Node{},
				Parent:   parent,
			}
			parent.Children = append(parent.Children, node)
			p.Current = node

		case lexer.TokenTypeContent:
			// Add the content under the current passage
			if err := p.addUnderCurrent(lexer.TokenTypePassage, line); err != nil {
				return nil, err
			}

		// Question
		case lexer.TokenTypeQuestion:
			// Find the nearest passage or header
			parent := p.findNearest(lexer.TokenTypePassage)
			if parent == nil {
				parent = p.findNearest(lexer.TokenTypeHeader)
			}
			if parent == nil {
				return nil, NewParserError(CodeValidation, "question without valid parent", line)
			}
			node := &Node{
				Type:     lexer.TokenTypeQuestion,
				Data:     line.ParsedValue,
				Children: []*Node{},
				Parent:   parent,
			}
			parent.Children = append(parent.Children, node)
			p.Current = node

		// LearnMore
		case lexer.TokenTypeLearnMore:
			// Add the learn more under the current question
			if err := p.addUnderCurrent(lexer.TokenTypeQuestion, line); err != nil {
				return nil, err
			}

		default:
			continue
		}
	}

	return p.finalize(parserType)
}
