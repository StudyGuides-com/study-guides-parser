package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/lexer"
	"github.com/studyguides-com/study-guides-parser/core/parser"
	"github.com/studyguides-com/study-guides-parser/core/preparser"
)

func TestBuildTree(t *testing.T) {
	// Create a test AST based on the new example
	ast := &parser.AbstractSyntaxTree{
		Metadata: &config.Metadata{
			Options: map[string]string{
				"file": "input.txt",
			},
			Type: "info",
		},
		Timestamp: "2025-07-09T11:15:18Z",
		Root: &parser.Node{
			Type: lexer.TokenTypeFileHeader,
			Data: preparser.ParsedValue{
				FileHeader: &preparser.FileHeaderResult{
					Title: "TestFile",
				},
			},
			Children: []*parser.Node{
				{
					Type: lexer.TokenTypeHeader,
					Data: preparser.ParsedValue{
						Header: &preparser.HeaderResult{
							Parts: []string{
								"TagA",
								"TagB", 
								"TagC",
								"TagD",
							},
						},
					},
					Children: []*parser.Node{
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 1 + 1?",
									AnswerText:   "2",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple addition",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 2 - 2?",
									AnswerText:   "0",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple subtraction",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Tim had 5 apples and gave Mike 3",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples are there?",
											AnswerText:   "5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Tim have?",
											AnswerText:   "2",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Mike have?",
											AnswerText:   "3",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Tim had $10 and gave Mike $5",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars are there?",
											AnswerText:   "$10",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Tim have?",
											AnswerText:   "$5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Mike have?",
											AnswerText:   "$5",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Build the tree
	tree := Build(ast, ast.Metadata)

	// Verify the tree structure
	if tree.Root == nil {
		t.Fatal("Tree root should not be nil")
	}

	if tree.Root.Title != "TestFile" {
		t.Errorf("Expected root title 'TestFile', got '%s'", tree.Root.Title)
	}

	if len(tree.Root.ChildTags) != 1 {
		t.Fatalf("Expected 1 child tag, got %d", len(tree.Root.ChildTags))
	}

	// Check the hierarchy
	tagA := tree.Root.ChildTags[0]
	if tagA.Title != "TagA" {
		t.Errorf("Expected 'TagA', got '%s'", tagA.Title)
	}

	tagB := tagA.ChildTags[0]
	if tagB.Title != "TagB" {
		t.Errorf("Expected 'TagB', got '%s'", tagB.Title)
	}

	tagC := tagB.ChildTags[0]
	if tagC.Title != "TagC" {
		t.Errorf("Expected 'TagC', got '%s'", tagC.Title)
	}

	tagD := tagC.ChildTags[0]
	if tagD.Title != "TagD" {
		t.Errorf("Expected 'TagD', got '%s'", tagD.Title)
	}

	// Check that all questions are in the last tag (including passage questions)
	expectedStandaloneQuestions := 2 // Only the two standalone questions
	expectedPassages := 2
	if len(tagD.Questions) != expectedStandaloneQuestions {
		t.Fatalf("Expected %d standalone questions, got %d", expectedStandaloneQuestions, len(tagD.Questions))
	}
	if len(tagD.Passages) != expectedPassages {
		t.Fatalf("Expected %d passages, got %d", expectedPassages, len(tagD.Passages))
	}

	// Check the first standalone question with learn more
	if tagD.Questions[0].Prompt != "What is 1 + 1?" {
		t.Errorf("Expected question 'What is 1 + 1?', got '%s'", tagD.Questions[0].Prompt)
	}
	if tagD.Questions[0].Answer != "2" {
		t.Errorf("Expected answer '2', got '%s'", tagD.Questions[0].Answer)
	}
	if tagD.Questions[0].LearnMore != "This is simple addition" {
		t.Errorf("Expected learn more 'This is simple addition', got '%s'", tagD.Questions[0].LearnMore)
	}

	// Check the second standalone question with learn more
	if tagD.Questions[1].Prompt != "What is 2 - 2?" {
		t.Errorf("Expected question 'What is 2 - 2?', got '%s'", tagD.Questions[1].Prompt)
	}
	if tagD.Questions[1].Answer != "0" {
		t.Errorf("Expected answer '0', got '%s'", tagD.Questions[1].Answer)
	}
	if tagD.Questions[1].LearnMore != "This is simple subtraction" {
		t.Errorf("Expected learn more 'This is simple subtraction', got '%s'", tagD.Questions[1].LearnMore)
	}

	// Check the first passage
	p1 := tagD.Passages[0]
	if p1.Title != "Tim had 5 apples and gave Mike 3" {
		t.Errorf("Expected passage text 'Tim had 5 apples and gave Mike 3', got '%s'", p1.Title)
	}
	if len(p1.Questions) != 3 {
		t.Errorf("Expected 3 questions in first passage, got %d", len(p1.Questions))
	}
	if p1.Questions[0].Prompt != "How many apples are there?" {
		t.Errorf("Expected passage question 'How many apples are there?', got '%s'", p1.Questions[0].Prompt)
	}
	if p1.Questions[0].Answer != "5" {
		t.Errorf("Expected passage answer '5', got '%s'", p1.Questions[0].Answer)
	}

	// Check the second passage
	p2 := tagD.Passages[1]
	if p2.Title != "Tim had $10 and gave Mike $5" {
		t.Errorf("Expected passage text 'Tim had $10 and gave Mike $5', got '%s'", p2.Title)
	}
	if len(p2.Questions) != 3 {
		t.Errorf("Expected 3 questions in second passage, got %d", len(p2.Questions))
	}
	if p2.Questions[0].Prompt != "How many dollars are there?" {
		t.Errorf("Expected passage question 'How many dollars are there?', got '%s'", p2.Questions[0].Prompt)
	}
	if p2.Questions[0].Answer != "$10" {
		t.Errorf("Expected passage answer '$10', got '%s'", p2.Questions[0].Answer)
	}
} 

func TestBuildTree_JSONOutput(t *testing.T) {
	// Use the same AST as in TestBuildTree
	ast := &parser.AbstractSyntaxTree{
		Metadata: &config.Metadata{
			Options: map[string]string{
				"file": "input.txt",
			},
			Type: "info",
		},
		Timestamp: "2025-07-09T11:15:18Z",
		Root: &parser.Node{
			Type: lexer.TokenTypeFileHeader,
			Data: preparser.ParsedValue{
				FileHeader: &preparser.FileHeaderResult{
					Title: "TestFile",
				},
			},
			Children: []*parser.Node{
				{
					Type: lexer.TokenTypeHeader,
					Data: preparser.ParsedValue{
						Header: &preparser.HeaderResult{
							Parts: []string{
								"TagA",
								"TagB", 
								"TagC",
								"TagD",
							},
						},
					},
					Children: []*parser.Node{
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 1 + 1?",
									AnswerText:   "2",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple addition",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 2 - 2?",
									AnswerText:   "0",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple subtraction",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Tim had 5 apples and gave Mike 3",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples are there?",
											AnswerText:   "5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Tim have?",
											AnswerText:   "2",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Mike have?",
											AnswerText:   "3",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Tim had $10 and gave Mike $5",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars are there?",
											AnswerText:   "$10",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Tim have?",
											AnswerText:   "$5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Mike have?",
											AnswerText:   "$5",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tree := Build(ast, ast.Metadata)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(tree.Root); err != nil {
		t.Fatalf("Failed to encode tree to JSON: %v", err)
	}
} 

func TestActualJSONOutput(t *testing.T) {
	// Use the same AST as in TestBuildTree
	ast := &parser.AbstractSyntaxTree{
		Metadata: &config.Metadata{
			Options: map[string]string{
				"file": "input.txt",
			},
			Type: "info",
		},
		Timestamp: "2025-07-09T11:15:18Z",
		Root: &parser.Node{
			Type: lexer.TokenTypeFileHeader,
			Data: preparser.ParsedValue{
				FileHeader: &preparser.FileHeaderResult{
					Title: "TestFile",
				},
			},
			Children: []*parser.Node{
				{
					Type: lexer.TokenTypeHeader,
					Data: preparser.ParsedValue{
						Header: &preparser.HeaderResult{
							Parts: []string{
								"TagA",
								"TagB", 
								"TagC",
								"TagD",
							},
						},
					},
					Children: []*parser.Node{
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 1 + 1?",
									AnswerText:   "2",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple addition",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 2 - 2?",
									AnswerText:   "0",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple subtraction",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Tim had 5 apples and gave Mike 3",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples are there?",
											AnswerText:   "5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Tim have?",
											AnswerText:   "2",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Mike have?",
											AnswerText:   "3",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Tim had $10 and gave Mike $5",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars are there?",
											AnswerText:   "$10",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Tim have?",
											AnswerText:   "$5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Mike have?",
											AnswerText:   "$5",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tree := Build(ast, ast.Metadata)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(tree.Root); err != nil {
		t.Fatalf("Failed to encode tree to JSON: %v", err)
	}
} 

func TestBuildTreeWithContent(t *testing.T) {
	// Create the AST based on the user's actual example
	ast := &parser.AbstractSyntaxTree{
		Root: &parser.Node{
			Type: lexer.TokenTypeFileHeader,
			Data: preparser.ParsedValue{
				FileHeader: &preparser.FileHeaderResult{
					Title: "TestFile",
				},
			},
			Children: []*parser.Node{
				{
					Type: lexer.TokenTypeHeader,
					Data: preparser.ParsedValue{
						Header: &preparser.HeaderResult{
							Parts: []string{"TagA", "TagB", "TagC", "TagD"},
						},
					},
					Children: []*parser.Node{
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 1 + 1?",
									AnswerText:   "2",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple addition",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 2 - 2?",
									AnswerText:   "0",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple subtraction",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Two best friends hanging out.",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeContent,
									Data: preparser.ParsedValue{
										Content: &preparser.ContentResult{
											Text: "It's a sunny day so the boys decide to have a snack",
										},
									},
								},
								{
									Type: lexer.TokenTypeContent,
									Data: preparser.ParsedValue{
										Content: &preparser.ContentResult{
											Text: "Tim had 5 apples and gave Mike 3",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples are there?",
											AnswerText:   "5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Tim have?",
											AnswerText:   "2",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Mike have?",
											AnswerText:   "3",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Tim had $10 and gave Mike $5",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars are there?",
											AnswerText:   "$10",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Tim have?",
											AnswerText:   "$5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Mike have?",
											AnswerText:   "$5",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	metadata := &config.Metadata{
		Options: map[string]string{
			"file": "input.txt",
		},
		Type: "info",
	}

	tree := Build(ast, metadata)

	// Verify the structure
	if tree.Root.Title != "TestFile" {
		t.Errorf("Expected root title 'TestFile', got '%s'", tree.Root.Title)
	}

	// Navigate to TagD
	tagA := tree.Root.ChildTags[0]
	tagB := tagA.ChildTags[0]
	tagC := tagB.ChildTags[0]
	tagD := tagC.ChildTags[0]

	// Check standalone questions
	if len(tagD.Questions) != 2 {
		t.Errorf("Expected 2 standalone questions, got %d", len(tagD.Questions))
	}

	// Check passages
	if len(tagD.Passages) != 2 {
		t.Errorf("Expected 2 passages, got %d", len(tagD.Passages))
	}

	// Check the first passage (with content)
	p1 := tagD.Passages[0]
	if p1.Title != "Two best friends hanging out." {
		t.Errorf("Expected passage title 'Two best friends hanging out.', got '%s'", p1.Title)
	}
	expectedContent := "It's a sunny day so the boys decide to have a snack\nTim had 5 apples and gave Mike 3"
	if p1.Content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, p1.Content)
	}
	if len(p1.Questions) != 3 {
		t.Errorf("Expected 3 questions in first passage, got %d", len(p1.Questions))
	}

	// Check the second passage (without content)
	p2 := tagD.Passages[1]
	if p2.Title != "Tim had $10 and gave Mike $5" {
		t.Errorf("Expected passage title 'Tim had $10 and gave Mike $5', got '%s'", p2.Title)
	}
	if p2.Content != "" {
		t.Errorf("Expected empty content in second passage, got '%s'", p2.Content)
	}
	if len(p2.Questions) != 3 {
		t.Errorf("Expected 3 questions in second passage, got %d", len(p2.Questions))
	}
} 

func TestBuildTreeWithContent_JSONOutput(t *testing.T) {
	// Create the AST based on the user's actual example
	ast := &parser.AbstractSyntaxTree{
		Root: &parser.Node{
			Type: lexer.TokenTypeFileHeader,
			Data: preparser.ParsedValue{
				FileHeader: &preparser.FileHeaderResult{
					Title: "TestFile",
				},
			},
			Children: []*parser.Node{
				{
					Type: lexer.TokenTypeHeader,
					Data: preparser.ParsedValue{
						Header: &preparser.HeaderResult{
							Parts: []string{"TagA", "TagB", "TagC", "TagD"},
						},
					},
					Children: []*parser.Node{
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 1 + 1?",
									AnswerText:   "2",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple addition",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypeQuestion,
							Data: preparser.ParsedValue{
								Question: &preparser.QuestionResult{
									QuestionText: "What is 2 - 2?",
									AnswerText:   "0",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeLearnMore,
									Data: preparser.ParsedValue{
										LearnMore: &preparser.LearnMoreResult{
											Text: "This is simple subtraction",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Two best friends hanging out.",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeContent,
									Data: preparser.ParsedValue{
										Content: &preparser.ContentResult{
											Text: "It's a sunny day so the boys decide to have a snack",
										},
									},
								},
								{
									Type: lexer.TokenTypeContent,
									Data: preparser.ParsedValue{
										Content: &preparser.ContentResult{
											Text: "Tim had 5 apples and gave Mike 3",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples are there?",
											AnswerText:   "5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Tim have?",
											AnswerText:   "2",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many apples does Mike have?",
											AnswerText:   "3",
										},
									},
								},
							},
						},
						{
							Type: lexer.TokenTypePassage,
							Data: preparser.ParsedValue{
								Passage: &preparser.PassageResult{
									Text: "Tim had $10 and gave Mike $5",
								},
							},
							Children: []*parser.Node{
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars are there?",
											AnswerText:   "$10",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Tim have?",
											AnswerText:   "$5",
										},
									},
								},
								{
									Type: lexer.TokenTypeQuestion,
									Data: preparser.ParsedValue{
										Question: &preparser.QuestionResult{
											QuestionText: "How many dollars does Mike have?",
											AnswerText:   "$5",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	metadata := &config.Metadata{
		Options: map[string]string{
			"file": "input.txt",
		},
		Type: "info",
	}

	tree := Build(ast, metadata)

	// Print the JSON output
	jsonData, err := json.MarshalIndent(tree.Root, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonData))
} 