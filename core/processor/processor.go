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

type LexerOutput struct {
	Filename string
	Tokens []lexer.LineInfo
	Errors    []error
	Success   bool
}

type PreparserOutput struct {
	Filename string
	Tokens []preparser.ParsedLineInfo
	Errors      []error
	Success     bool
}

// ParseFile reads a file and parses it into an Abstract Syntax Tree
func ParseFile(filename string, parserType ParserType) (*parser.AbstractSyntaxTree, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	lines := strings.Split(string(content), "\n")
	return Parse(lines, parserType)
}

// Parse parses a slice of strings into an Abstract Syntax Tree
func Parse(lines []string, parserType ParserType) (*parser.AbstractSyntaxTree, error) {
	preOut, err := Preparse(lines)
	if err != nil {
		return nil, fmt.Errorf("preparser error: %w", err)
	}
	if !preOut.Success {
		return nil, fmt.Errorf("preparser failed: %v", preOut.Errors)
	}

	p := parser.NewParser(parser.ParserType(parserType), preOut.Tokens)
	ast, parserErr := p.Parse()
	if parserErr != nil {
		return nil, fmt.Errorf("parser error: %w", parserErr)
	}
	return ast, nil
}

func Lex(lines []string) (LexerOutput, error) {
	lex := lexer.NewLexer()
	var tokens []lexer.LineInfo
	var errors []error

	for i, line := range lines {
		lineInfo, err := lex.ProcessLine(line, i+1)
		if err != nil {
			errors = append(errors, fmt.Errorf("line %d: %w", i+1, err))
		}
		tokens = append(tokens, lineInfo)
	}

	return LexerOutput{
		Tokens:  tokens,
		Errors:  errors,
		Success: len(errors) == 0,
	}, nil
}

func Preparse(lines []string) (PreparserOutput, error) {
	lexOut, err := Lex(lines)
	if err != nil || !lexOut.Success {
		return PreparserOutput{
			Errors:  append([]error{}, err),
			Success: false,
		}, err
	}

	pre := preparser.NewPreparser(lexOut.Tokens, "")
	parsed, prepErr := pre.Parse()
	var errors []error
	if prepErr != nil {
		errors = append(errors, prepErr)
	}
	return PreparserOutput{
		Tokens:  parsed,
		Errors:  errors,
		Success: len(errors) == 0,
	}, nil
}