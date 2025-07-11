package builder

import (
	"strings"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/lexer"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
	"github.com/studyguides-com/study-guides-parser/core/parser"
	"github.com/studyguides-com/study-guides-parser/core/qa"
	"github.com/studyguides-com/study-guides-parser/core/tree"
)

func Build(ast *parser.AbstractSyntaxTree, metadata *config.Metadata) *tree.Tree {
	tree := tree.NewTree(metadata)

	if ast.Root == nil {
		return tree
	}

	// Walk through the AST and build the tree
	buildTree(ast.Root, tree.Root)

	// Assign tag types based on context
	if metadata.ContextType != ontology.ContextTypeNone {
		_ = tree.AssignTagTypes(metadata.ContextType)
	}

	// Run QA
	qaRunner := qa.NewTreeQARunner(
		qa.NewTagTypeQA(),
		qa.NewContextTypeQA(),
	)
	qaRunner.RunQAAndUpdate(tree)

	return tree
}

// BuildWithContext builds a tree and requires a context type for tag assignment
func BuildWithContext(ast *parser.AbstractSyntaxTree, metadata *config.Metadata, contextType ontology.ContextType) *tree.Tree {
	tree := tree.NewTree(metadata)

	if ast.Root == nil {
		return tree
	}

	// Walk through the AST and build the tree
	buildTree(ast.Root, tree.Root)

	// Assign tag types based on the provided context
	tree.AssignTagTypes(contextType)

	// Run QA
	qaRunner := qa.NewTreeQARunner(
		qa.NewTagTypeQA(),
		qa.NewContextTypeQA(),
	)
	qaRunner.RunQAAndUpdate(tree)

	return tree
}

func buildTree(node *parser.Node, currentTag tree.TagContainer) {
	if node == nil {
		return
	}

	switch node.Type {
	case lexer.TokenTypeFileHeader:
		// File header contains the overall title
		if fileHeader := node.Data.GetFileHeader(); fileHeader != nil {
			if root, ok := currentTag.(*tree.Root); ok {
				root.Title = fileHeader.Title
			}
		}
		// Process children
		for _, child := range node.Children {
			buildTree(child, currentTag)
		}

	case lexer.TokenTypeHeader:
		// Header creates a new tag structure
		if header := node.Data.GetHeader(); header != nil {
			// Build the tag hierarchy from header parts
			tag := buildTagHierarchy(currentTag, header.Parts)
			// Process children (questions, passages, etc.) and add them to the last tag
			for _, child := range node.Children {
				buildTree(child, tag)
			}
		}

	case lexer.TokenTypeQuestion:
		// Question gets added to the current tag
		if question := node.Data.GetQuestion(); question != nil {
			// Check for learn more content in children
			var learnMoreText string
			for _, child := range node.Children {
				if child.Type == lexer.TokenTypeLearnMore {
					if learnMore := child.Data.GetLearnMore(); learnMore != nil {
						learnMoreText = learnMore.Text
					}
				}
			}
			q := tree.NewQuestion(question.QuestionText, question.AnswerText, nil, learnMoreText)
			if tag, ok := currentTag.(*tree.Tag); ok {
				tag.Questions = append(tag.Questions, q)
			}
		}

	case lexer.TokenTypePassage:
		// Passage creates a new passage structure
		if passage := node.Data.GetPassage(); passage != nil {
			// Process children (content and questions) and collect data
			var contentLines []string
			var questions []*tree.Question
			for _, child := range node.Children {
				if child.Type == lexer.TokenTypeQuestion {
					if question := child.Data.GetQuestion(); question != nil {
						// Check for learn more content in question children
						var learnMoreText string
						for _, grandChild := range child.Children {
							if grandChild.Type == lexer.TokenTypeLearnMore {
								if learnMore := grandChild.Data.GetLearnMore(); learnMore != nil {
									learnMoreText = learnMore.Text
								}
							}
						}
						q := tree.NewQuestion(question.QuestionText, question.AnswerText, nil, learnMoreText)
						questions = append(questions, q)
					}
				} else if child.Type == lexer.TokenTypeContent {
					if content := child.Data.GetContent(); content != nil {
						contentLines = append(contentLines, content.Text)
					}
				}
			}
			// Concatenate content lines with newlines
			var content string
			if len(contentLines) > 0 {
				content = strings.Join(contentLines, "\n")
			}
			// Create passage using NewPassage constructor
			p := tree.NewPassage(passage.Text, content, questions)
			// Add the passage to the current tag's Passages
			if tag, ok := currentTag.(*tree.Tag); ok {
				tag.Passages = append(tag.Passages, p)
			}
		}

	default:
		// For other node types, just process children
		for _, child := range node.Children {
			buildTree(child, currentTag)
		}
	}
}

func buildTagHierarchy(parentTag tree.TagContainer, headerParts []string) *tree.Tag {
	if len(headerParts) == 0 {
		if tag, ok := parentTag.(*tree.Tag); ok {
			return tag
		}
		return nil
	}

	// Find or create the tag for the first part
	var currentTag *tree.Tag
	for _, child := range parentTag.GetChildTags() {
		if child.Title == headerParts[0] {
			currentTag = child
			break
		}
	}

	if currentTag == nil {
		// For the first header part, create as root-level tag (no parent hash)
		// For subsequent parts, create with parent-based hash
		var parentTitle string
		if root, ok := parentTag.(*tree.Root); ok {
			parentTitle = root.Title
		} else if tag, ok := parentTag.(*tree.Tag); ok {
			parentTitle = tag.Title
		}

		if parentTitle == "TestFile" || parentTitle == "Root" {
			currentTag = tree.NewTag(headerParts[0])
		} else {
			currentTag = tree.NewTagWithParent(headerParts[0], parentTitle)
		}
		parentTag.AddChildTag(currentTag)
	}

	// Recursively build the rest of the hierarchy
	if len(headerParts) > 1 {
		return buildTagHierarchy(currentTag, headerParts[1:])
	}

	return currentTag
}

// example text file
/*
TestFile

TagA: TagB: TagC: TagD

1. What is 1 + 1? - 2
Learn More: This is simple addition
2. What is 2 - 2? - 0
Learn More: This is simple subtraction

Passage: Tim had 5 apples and gave Mike 3

1. How many apples are there? - 5
2. How many apples does Tim have? - 2
3. How many apples does Mike have? - 3

Passage: Tim had $10 and gave Mike $5

1. How many dollars are there? - $10
2. How many dollars does Tim have? - $5
3. How many dollars does Mike have? - $5

*/

// example ast
/*
{
  "ast": {
    "metadata": {
      "options": {
        "file": "input.txt"
      },
      "type": "info"
    },
    "root": {
      "children": [
        {
          "children": [
            {
              "children": [
                {
                  "data": {
                    "learn_more": {
                      "Text": "This is simple addition"
                    }
                  },
                  "type": "learn_more"
                }
              ],
              "data": {
                "question": {
                  "AnswerText": "2",
                  "QuestionText": "What is 1 + 1?"
                }
              },
              "type": "question"
            },
            {
              "children": [
                {
                  "data": {
                    "learn_more": {
                      "Text": "This is simple subtraction"
                    }
                  },
                  "type": "learn_more"
                }
              ],
              "data": {
                "question": {
                  "AnswerText": "2? - 0",
                  "QuestionText": "What is 2"
                }
              },
              "type": "question"
            },
            {
              "children": [
                {
                  "data": {
                    "question": {
                      "AnswerText": "5",
                      "QuestionText": "How many apples are there?"
                    }
                  },
                  "type": "question"
                },
                {
                  "data": {
                    "question": {
                      "AnswerText": "2",
                      "QuestionText": "How many apples does Tim have?"
                    }
                  },
                  "type": "question"
                },
                {
                  "data": {
                    "question": {
                      "AnswerText": "3",
                      "QuestionText": "How many apples does Mike have?"
                    }
                  },
                  "type": "question"
                }
              ],
              "data": {
                "passage": {
                  "Text": "Tim had 5 apples and gave Mike 3"
                }
              },
              "type": "passage"
            },
            {
              "children": [
                {
                  "data": {
                    "question": {
                      "AnswerText": "$10",
                      "QuestionText": "How many dollars are there?"
                    }
                  },
                  "type": "question"
                },
                {
                  "data": {
                    "question": {
                      "AnswerText": "$5",
                      "QuestionText": "How many dollars does Tim have?"
                    }
                  },
                  "type": "question"
                },
                {
                  "data": {
                    "question": {
                      "AnswerText": "$5",
                      "QuestionText": "How many dollars does Mike have?"
                    }
                  },
                  "type": "question"
                }
              ],
              "data": {
                "passage": {
                  "Text": "Tim had $10 and gave Mike $5"
                }
              },
              "type": "passage"
            }
          ],
          "data": {
            "header": {
              "Parts": [
                "TagA",
                "TagB",
                "TagC",
                "TagD"
              ]
            }
          },
          "type": "header"
        }
      ],
      "data": {
        "file_header": {
          "Title": "TestFile"
        }
      },
      "type": "file_header"
    },
    "timestamp": "2025-07-09T11:15:18Z"
  },
  "success": true
}
*/
