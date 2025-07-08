package processor

import (
	"fmt"
	"os"
	"strings"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/lexer"
	"github.com/studyguides-com/study-guides-parser/core/parser"
	"github.com/studyguides-com/study-guides-parser/core/preparser"
)



type LexerOutput struct {
	Filename string        `json:"filename"`
	Tokens   []lexer.LineInfo `json:"tokens"`
	Errors   []string      `json:"errors"`
	Success  bool          `json:"success"`
}

type PreparserOutput struct {
	Filename string                      `json:"filename"`
	Tokens   []preparser.ParsedLineInfo `json:"tokens"`
	Errors   []string                   `json:"errors"`
	Success  bool                       `json:"success"`
}

// ParseFile reads a file and parses it into an Abstract Syntax Tree
func ParseFile(filename string, metadata *config.Metadata) (*parser.AbstractSyntaxTree, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	lines := strings.Split(string(content), "\n")
	return Parse(lines, metadata)
}

// Parse parses a slice of strings into an Abstract Syntax Tree
func Parse(lines []string, metadata *config.Metadata) (*parser.AbstractSyntaxTree, error) {
	preOut, err := Preparse(lines)
	if err != nil {
		return nil, fmt.Errorf("preparser error: %w", err)
	}
	if !preOut.Success {
		return nil, fmt.Errorf("preparser failed: %v", preOut.Errors)
	}

	p := parser.NewParser(preOut.Tokens)
	ast, parserErr := p.Parse(metadata)
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

	// Convert errors to strings for JSON serialization
	errorStrings := make([]string, len(errors))
	for i, err := range errors {
		errorStrings[i] = err.Error()
	}

	return LexerOutput{
		Tokens:  tokens,
		Errors:  errorStrings,
		Success: len(errors) == 0,
	}, nil
}

func Preparse(lines []string) (PreparserOutput, error) {
	// Step 1: Run lexer and collect all lexer errors
	lexOut, err := Lex(lines)
	if err != nil {
		// If there's a critical error with the lexer itself, return it
		return PreparserOutput{
			Errors:  []string{err.Error()},
			Success: false,
		}, err
	}

	// If lexer failed, return immediately with lexer errors
	if !lexOut.Success {
		return PreparserOutput{
			Tokens:  nil,
			Errors:  lexOut.Errors,
			Success: false,
		}, nil
	}

	// Step 2: Run preparser only if lexer succeeded
	pre := preparser.NewPreparser(lexOut.Tokens, "")
	parsed, prepErrors := pre.Parse()
	
	// Add all preparser errors if any, including line numbers
	var allErrors []string
	for _, prepErr := range prepErrors {
		allErrors = append(allErrors, fmt.Sprintf("line %d: %s", prepErr.LineInfo.Number, prepErr.Error()))
	}
	
	return PreparserOutput{
		Tokens:  parsed,
		Errors:  allErrors,
		Success: len(allErrors) == 0,
	}, nil
}