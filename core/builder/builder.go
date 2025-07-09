package builder

import (
	"strings"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/lexer"
	"github.com/studyguides-com/study-guides-parser/core/parser"
)

func Build(ast *parser.AbstractSyntaxTree, metadata *config.Metadata) *Tree {
	tree := NewTree(metadata)
	
	if ast.Root == nil {
		return tree
	}
	
	// Walk through the AST and build the tree
	buildTree(ast.Root, tree.Root)
	
	return tree
}

func buildTree(node *parser.Node, currentTag *Tag) {
	if node == nil {
		return
	}
	
	switch node.Type {
	case lexer.TokenTypeFileHeader:
		// File header contains the overall title
		if fileHeader := node.Data.GetFileHeader(); fileHeader != nil {
			currentTag.Title = fileHeader.Title
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
			q := &Question{
				Prompt: question.QuestionText,
				Answer: question.AnswerText,
			}
			// Check for learn more content in children
			for _, child := range node.Children {
				if child.Type == lexer.TokenTypeLearnMore {
					if learnMore := child.Data.GetLearnMore(); learnMore != nil {
						q.LearnMore = learnMore.Text
					}
				}
			}
			currentTag.Questions = append(currentTag.Questions, q)
		}
		
	case lexer.TokenTypePassage:
		// Passage creates a new passage structure
		if passage := node.Data.GetPassage(); passage != nil {
			p := &Passage{
				Title: passage.Text,
			}
			// Process children (content and questions) and add them to the passage
			var contentLines []string
			for _, child := range node.Children {
				if child.Type == lexer.TokenTypeQuestion {
					if question := child.Data.GetQuestion(); question != nil {
						q := &Question{
							Prompt: question.QuestionText,
							Answer: question.AnswerText,
						}
						// Check for learn more content in question children
						for _, grandChild := range child.Children {
							if grandChild.Type == lexer.TokenTypeLearnMore {
								if learnMore := grandChild.Data.GetLearnMore(); learnMore != nil {
									q.LearnMore = learnMore.Text
								}
							}
						}
						p.Questions = append(p.Questions, q)
					}
				} else if child.Type == lexer.TokenTypeContent {
					if content := child.Data.GetContent(); content != nil {
						contentLines = append(contentLines, content.Text)
					}
				}
			}
			// Concatenate content lines with newlines
			if len(contentLines) > 0 {
				p.Content = strings.Join(contentLines, "\n")
			}
			// Add the passage to the current tag's Passages
			currentTag.Passages = append(currentTag.Passages, p)
		}
		
	default:
		// For other node types, just process children
		for _, child := range node.Children {
			buildTree(child, currentTag)
		}
	}
}

func buildTagHierarchy(parentTag *Tag, headerParts []string) *Tag {
	if len(headerParts) == 0 {
		return parentTag
	}
	
	// Find or create the tag for the first part
	var currentTag *Tag
	for _, child := range parentTag.ChildTags {
		if child.Title == headerParts[0] {
			currentTag = child
			break
		}
	}
	
	if currentTag == nil {
		// Create new tag
		currentTag = NewTag(headerParts[0])
		parentTag.ChildTags = append(parentTag.ChildTags, currentTag)
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
