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

// ProcessingError represents a structured error with line information
type ProcessingError struct {
	LineNumber int    `json:"line_number"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Text       string `json:"text,omitempty"`
	Type       string `json:"type,omitempty"`
}

type LexerOutput struct {
	Filename string        `json:"filename"`
	Tokens   []lexer.LineInfo `json:"tokens"`
	Errors   []ProcessingError `json:"errors"`
	Success  bool          `json:"success"`
}

type PreparserOutput struct {
	Filename string                      `json:"filename"`
	Tokens   []preparser.ParsedLineInfo `json:"tokens"`
	Errors   []ProcessingError           `json:"errors"`
	Success  bool                       `json:"success"`
}

// ParseResult represents the result of parsing with structured errors
type ParseResult struct {
	AST    *parser.AbstractSyntaxTree `json:"ast,omitempty"`
	Errors []ProcessingError          `json:"errors,omitempty"`
	Success bool                      `json:"success"`
}

// ParseFile reads a file and parses it into an Abstract Syntax Tree
func ParseFile(filename string, metadata *config.Metadata) (*ParseResult, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	lines := strings.Split(string(content), "\n")
	return Parse(lines, metadata)
}

// Parse parses a slice of strings into an Abstract Syntax Tree
func Parse(lines []string, metadata *config.Metadata) (*ParseResult, error) {
	preOut, err := Preparse(lines)
	if err != nil {
		return nil, fmt.Errorf("preparser error: %w", err)
	}
	if !preOut.Success {
		return &ParseResult{
			Errors:  preOut.Errors,
			Success: false,
		}, nil
	}

	p := parser.NewParser(preOut.Tokens)
	ast, parserErr := p.Parse(metadata)
	if parserErr != nil {
		// Convert parser error to ProcessingError format
		parserError := ProcessingError{
			LineNumber: parserErr.LineInfo.Number,
			Message:    parserErr.Message,
			Code:       string(parserErr.Code),
			Text:       parserErr.LineInfo.Text,
			Type:       string(parserErr.LineInfo.Type),
		}
		return &ParseResult{
			Errors:  []ProcessingError{parserError},
			Success: false,
		}, nil
	}
	return &ParseResult{
		AST:     ast,
		Success: true,
	}, nil
}

func Lex(lines []string) (LexerOutput, error) {
	lex := lexer.NewLexer()
	var tokens []lexer.LineInfo
	var errors []*lexer.LexerError

	for i, line := range lines {
		lineInfo, err := lex.ProcessLine(line, i+1)
		if err != nil {
			errors = append(errors, err)
		}
		tokens = append(tokens, lineInfo)
	}

	// Convert lexer errors to ProcessingError structs for JSON serialization
	processingErrors := make([]ProcessingError, len(errors))
	for i, err := range errors {
		processingErrors[i] = ProcessingError{
			LineNumber: err.LineInfo.Number,
			Message:    err.Message,
			Code:       string(err.Code),
			Text:       err.LineInfo.Text,
			Type:       string(err.LineInfo.Type),
		}
	}

	return LexerOutput{
		Tokens:  tokens,
		Errors:  processingErrors,
		Success: len(errors) == 0,
	}, nil
}

func Preparse(lines []string) (PreparserOutput, error) {
	// Step 1: Run lexer and collect all lexer errors
	lexOut, err := Lex(lines)
	if err != nil {
		// If there's a critical error with the lexer itself, return it
		return PreparserOutput{
			Errors:  []ProcessingError{{LineNumber: 0, Message: err.Error(), Code: "CRITICAL_ERROR"}},
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
	var allErrors []ProcessingError
	for _, prepErr := range prepErrors {
		allErrors = append(allErrors, ProcessingError{
			LineNumber: prepErr.LineInfo.Number,
			Message:    prepErr.Message,
			Code:       string(prepErr.Code),
			Text:       prepErr.LineInfo.Text,
			Type:       string(prepErr.LineInfo.Type),
		})
	}
	
	return PreparserOutput{
		Tokens:  parsed,
		Errors:  allErrors,
		Success: len(allErrors) == 0,
	}, nil
}

// ParseWithErrors parses a slice of strings and returns both the AST and any errors that occurred
func ParseWithErrors(lines []string, metadata *config.Metadata) (*parser.AbstractSyntaxTree, []ProcessingError, error) {
	preOut, err := Preparse(lines)
	if err != nil {
		return nil, nil, fmt.Errorf("preparser error: %w", err)
	}
	if !preOut.Success {
		return nil, preOut.Errors, nil
	}

	p := parser.NewParser(preOut.Tokens)
	ast, parserErr := p.Parse(metadata)
	if parserErr != nil {
		// Convert parser error to ProcessingError format
		parserError := ProcessingError{
			LineNumber: parserErr.LineInfo.Number,
			Message:    parserErr.Message,
			Code:       string(parserErr.Code),
			Text:       parserErr.LineInfo.Text,
			Type:       string(parserErr.LineInfo.Type),
		}
		return nil, []ProcessingError{parserError}, nil
	}
	return ast, nil, nil
}

// ParseFileWithErrors reads a file and parses it, returning both the AST and any errors
func ParseFileWithErrors(filename string, metadata *config.Metadata) (*parser.AbstractSyntaxTree, []ProcessingError, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	lines := strings.Split(string(content), "\n")
	return ParseWithErrors(lines, metadata)
}