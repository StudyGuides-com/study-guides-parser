//go:build !prod

package parser

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/preparser"
)

func TestParserWithContentInPassage(t *testing.T) {
	// Create the input lines that represent the user's example
	lines := []preparser.ParsedLineInfo{
		{
			Type: preparser.TokenTypeFileHeader,
			ParsedValue: preparser.ParsedValue{
				FileHeader: &preparser.FileHeaderResult{
					Title: "TestFile",
				},
			},
		},
		{
			Type: preparser.TokenTypeHeader,
			ParsedValue: preparser.ParsedValue{
				Header: &preparser.HeaderResult{
					Parts: []string{"TagA", "TagB", "TagC", "TagD"},
				},
			},
		},
		{
			Type: preparser.TokenTypeQuestion,
			ParsedValue: preparser.ParsedValue{
				Question: &preparser.QuestionResult{
					QuestionText: "What is 1 + 1?",
					AnswerText:   "2",
				},
			},
		},
		{
			Type: preparser.TokenTypeLearnMore,
			ParsedValue: preparser.ParsedValue{
				LearnMore: &preparser.LearnMoreResult{
					Text: "This is simple addition",
				},
			},
		},
		{
			Type: preparser.TokenTypeQuestion,
			ParsedValue: preparser.ParsedValue{
				Question: &preparser.QuestionResult{
					QuestionText: "What is 2 - 2?",
					AnswerText:   "0",
				},
			},
		},
		{
			Type: preparser.TokenTypeLearnMore,
			ParsedValue: preparser.ParsedValue{
				LearnMore: &preparser.LearnMoreResult{
					Text: "This is simple subtraction",
				},
			},
		},
		{
			Type: preparser.TokenTypePassage,
			ParsedValue: preparser.ParsedValue{
				Passage: &preparser.PassageResult{
					Text: "Two best friends hanging out.",
				},
			},
		},
		{
			Type: preparser.TokenTypeContent,
			ParsedValue: preparser.ParsedValue{
				Content: &preparser.ContentResult{
					Text: "It's a sunny day so the boys decide to have a snack",
				},
			},
		},
		{
			Type: preparser.TokenTypeContent,
			ParsedValue: preparser.ParsedValue{
				Content: &preparser.ContentResult{
					Text: "Tim had 5 apples and gave Mike 3",
				},
			},
		},
		{
			Type: preparser.TokenTypeQuestion,
			ParsedValue: preparser.ParsedValue{
				Question: &preparser.QuestionResult{
					QuestionText: "How many apples are there?",
					AnswerText:   "5",
				},
			},
		},
		{
			Type: preparser.TokenTypeQuestion,
			ParsedValue: preparser.ParsedValue{
				Question: &preparser.QuestionResult{
					QuestionText: "How many apples does Tim have?",
					AnswerText:   "2",
				},
			},
		},
		{
			Type: preparser.TokenTypeQuestion,
			ParsedValue: preparser.ParsedValue{
				Question: &preparser.QuestionResult{
					QuestionText: "How many apples does Mike have?",
					AnswerText:   "3",
				},
			},
		},
		{
			Type: preparser.TokenTypePassage,
			ParsedValue: preparser.ParsedValue{
				Passage: &preparser.PassageResult{
					Text: "Tim had $10 and gave Mike $5",
				},
			},
		},
		{
			Type: preparser.TokenTypeQuestion,
			ParsedValue: preparser.ParsedValue{
				Question: &preparser.QuestionResult{
					QuestionText: "How many dollars are there?",
					AnswerText:   "$10",
				},
			},
		},
		{
			Type: preparser.TokenTypeQuestion,
			ParsedValue: preparser.ParsedValue{
				Question: &preparser.QuestionResult{
					QuestionText: "How many dollars does Tim have?",
					AnswerText:   "$5",
				},
			},
		},
		{
			Type: preparser.TokenTypeQuestion,
			ParsedValue: preparser.ParsedValue{
				Question: &preparser.QuestionResult{
					QuestionText: "How many dollars does Mike have?",
					AnswerText:   "$5",
				},
			},
		},
	}

	parser := NewParser(lines)
	metadata := &config.Metadata{
		Options: map[string]string{
			"file": "input.txt",
		},
		Type: "info",
	}

	ast, err := parser.Parse(metadata)
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	// Print the AST for inspection
	jsonData, jsonErr := json.MarshalIndent(ast, "", "  ")
	if jsonErr != nil {
		t.Fatalf("Failed to marshal JSON: %v", jsonErr)
	}
	fmt.Println("Generated AST:")
	fmt.Println(string(jsonData))

	// Verify the structure
	if ast.Root == nil {
		t.Fatal("AST root is nil")
	}

	// Navigate to the first passage and check its children
	header := ast.Root.Children[0]
	if len(header.Children) < 3 {
		t.Fatalf("Expected at least 3 children under header, got %d", len(header.Children))
	}

	// Find the first passage
	var firstPassage *Node
	for _, child := range header.Children {
		if child.Type == preparser.TokenTypePassage {
			firstPassage = child
			break
		}
	}

	if firstPassage == nil {
		t.Fatal("First passage not found")
	}

	// Check that the passage has content and questions as children
	contentCount := 0
	questionCount := 0
	for _, child := range firstPassage.Children {
		switch child.Type {
		case preparser.TokenTypeContent:
			contentCount++
		case preparser.TokenTypeQuestion:
			questionCount++
		}
	}

	if contentCount != 2 {
		t.Errorf("Expected 2 content children in first passage, got %d", contentCount)
	}
	if questionCount != 3 {
		t.Errorf("Expected 3 question children in first passage, got %d", questionCount)
	}
}
