package processor

import (
	"fmt"
	"os"
	"strings"

	"github.com/StudyGuides-com/study-guides-parser/core/lexer"
	"github.com/StudyGuides-com/study-guides-parser/core/parser"
	"github.com/StudyGuides-com/study-guides-parser/core/preparser"
)

// ParserType represents the type of study guide parser to use
type ParserType string

const (
	Colleges       ParserType = "colleges"
	APExams        ParserType = "ap_exams"
	Certifications ParserType = "certifications"
	DOD            ParserType = "dod"
	EntranceExams  ParserType = "entrance_exams"
)

type LexerResponse struct {
	Filename string
	Tokens []lexer.LineInfo
	Errors    []error
	Success   bool
}

type PreparserResponse struct {
	Filename string
	Tokens []preparser.LineInfo
	Errors      []error
	Success     bool
}


// ParseFile reads a file and parses it into an Abstract Syntax Tree
func ParseFile(filename string, parserType ParserType) (*parser.AbstractSyntaxTree, error) {
	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Split into lines
	lines := strings.Split(string(content), "\n")

	// Parse the lines
	return ParseLines(lines, parserType)
}

// ParseLines parses a slice of strings into an Abstract Syntax Tree
func ParseLines(lines []string, parserType ParserType) (*parser.AbstractSyntaxTree, error) {
	// Step 1: Lexer - classify each line
	lex := lexer.NewLexer()
	var lineInfos []lexer.LineInfo
	
	for i, line := range lines {
		lineInfo, err := lex.ProcessLine(line, i+1)
		if err != nil {
			return nil, fmt.Errorf("lexer error at line %d: %w", i+1, err)
		}
		lineInfos = append(lineInfos, lineInfo)
	}

	// Step 2: Preparser - parse values from each line
	pre := preparser.NewPreparser(lineInfos, string(parserType))
	parsedLines, prepErr := pre.Parse()
	if prepErr != nil {
		return nil, fmt.Errorf("preparser error: %w", prepErr)
	}

	// Step 3: Parser - build the AST
	p := parser.NewParser(parser.ParserType(parserType), parsedLines)
	ast, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("parser error: %w", err)
	}

	return ast, nil
}

// ParseTokens parses pre-processed tokens into an Abstract Syntax Tree
// This is useful when you already have tokens from a scanner or other source
func ParseTokens(tokens []lexer.LineInfo, parserType ParserType) (*parser.AbstractSyntaxTree, error) {
	// Step 1: Preparser - parse values from tokens
	pre := preparser.NewPreparser(tokens, string(parserType))
	parsedLines, prepErr := pre.Parse()
	if prepErr != nil {
		return nil, fmt.Errorf("preparser error: %w", prepErr)
	}

	// Step 2: Parser - build the AST
	p := parser.NewParser(parser.ParserType(parserType), parsedLines)
	ast, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("parser error: %w", err)
	}

	return ast, nil
}

func LexFile(filename string) (LexerResponse, error) {
	return LexerResponse{}, nil
}

func LexLines(lines []string) (LexerResponse, error) {
	return LexerResponse{}, nil
}

func PreparseFile(filename string) (PreparserResponse, error) {
	return PreparserResponse{}, nil
}

func PreparseLines(lines []string) (PreparserResponse, error) {
	return PreparserResponse{}, nil
}