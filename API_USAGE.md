# Study Guides Parser API Usage Guide

This guide shows how to use the simplified API for the `github.com/StudyGuides-com/study-guides-parser` library.

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/StudyGuides-com/study-guides-parser/core/processor"
)

func main() {
    lines := []string{
        "Mathematics Study Guide",
        "Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations",
        "",
        "1. What is a linear equation? - An equation where the highest power of the variable is 1.",
    }

    metadata := processor.NewMetadata(processor.Colleges)
    ast, err := processor.Parse(lines, metadata)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Parsed %d children\n", len(ast.Root.Children))
}
```

## API Overview

### Core Functions

- `Parse(lines []string, metadata *Metadata) (*AbstractSyntaxTree, error)` - Parse strings into AST
- `ParseFile(filename string, metadata *Metadata) (*AbstractSyntaxTree, error)` - Parse file into AST
- `Lex(lines []string) (LexerOutput, error)` - Lex strings into tokens
- `Preparse(lines []string) (PreparserOutput, error)` - Preparse strings into parsed tokens

### Metadata Configuration

```go
// Simple metadata
metadata := processor.NewMetadata(processor.Colleges)

// Advanced metadata with options
metadata := processor.NewMetadata(processor.APExams).
    WithOption("strict", "true").
    WithOption("debug", "false").
    WithOption("version", "1.0")
```

### Parser Types

- `processor.Colleges` - College study guides
- `processor.APExams` - Advanced Placement exams
- `processor.Certifications` - Professional certifications
- `processor.DOD` - Department of Defense
- `processor.EntranceExams` - College entrance exams

## Examples

### Basic Usage

```go
// Parse from strings
lines := []string{
    "Study Guide Title",
    "Colleges: State: University: Department: Course: Topic",
    "",
    "1. Question? - Answer.",
}

metadata := processor.NewMetadata(processor.Colleges)
ast, err := processor.Parse(lines, metadata)
if err != nil {
    log.Fatal(err)
}
```

### File Processing

```go
// Parse from file
metadata := processor.NewMetadata(processor.APExams)
ast, err := processor.ParseFile("study_guide.txt", metadata)
if err != nil {
    log.Fatal(err)
}
```

### Advanced Configuration

```go
// Configure with options
metadata := processor.NewMetadata(processor.Colleges).
    WithOption("strict", "true").
    WithOption("debug", "false").
    WithOption("output_format", "json")

ast, err := processor.Parse(lines, metadata)
```

### Error Handling

```go
ast, err := processor.Parse(lines, metadata)
if err != nil {
    if strings.Contains(err.Error(), "preparser failed") {
        fmt.Println("Input format error")
    } else if strings.Contains(err.Error(), "parser error") {
        fmt.Println("Parsing logic error")
    } else {
        fmt.Println("Unexpected error:", err)
    }
    return
}
```

### Working with AST

```go
ast, err := processor.Parse(lines, metadata)
if err != nil {
    log.Fatal(err)
}

// Access AST properties
fmt.Printf("Parser Type: %s\n", ast.ParserType)
fmt.Printf("Timestamp: %s\n", ast.Timestamp)
fmt.Printf("Root Type: %s\n", ast.Root.Type)

// Traverse children
for _, child := range ast.Root.Children {
    fmt.Printf("Child: %s\n", child.Type)
}
```

## Input Format

### College Format

```
Study Guide Title
Colleges: State: University: Department: Course: Topic

1. Question? - Answer.
2. Another question? - Another answer.

Learn More: Additional information.
```

### AP Exam Format

```
AP Subject Study Guide
Advanced Placement (AP): AP Subject: Topic: Subtopic

1. Question? - Answer.
2. Another question? - Another answer.
```

## Output Structure

The parser returns an `AbstractSyntaxTree` with:

```go
type AbstractSyntaxTree struct {
    ParserType string    `json:"parser_type"`
    Timestamp  string    `json:"timestamp"`
    Root       *ASTNode  `json:"root"`
}

type ASTNode struct {
    Type     string                 `json:"type"`
    Data     map[string]interface{} `json:"data"`
    Children []*ASTNode             `json:"children"`
}
```

## Best Practices

1. **Always check errors** - The parser can fail on malformed input
2. **Use appropriate parser type** - Different formats have different conventions
3. **Validate input** - Ensure your study guide follows the expected format
4. **Handle large files** - For very large files, consider processing in chunks
5. **Use metadata options** - Configure behavior for your specific needs

## Migration from Old API

### Before (Old API)

```go
// Manual pipeline management
lex := lexer.NewLexer()
var tokens []lexer.LineInfo
for i, line := range lines {
    token, err := lex.ProcessLine(line, i+1)
    if err != nil {
        return err
    }
    tokens = append(tokens, token)
}

// Manual conversion
pre := preparser.NewPreparser(tokens, "")
parsed, err := pre.Parse()
if err != nil {
    return err
}

// Manual parser setup
p := parser.NewParser(parser.ParserType("colleges"), parsed)
ast, err := p.Parse()
```

### After (New API)

```go
// One function call
metadata := processor.NewMetadata(processor.Colleges)
ast, err := processor.Parse(lines, metadata)
```

## Error Types

- **Preparser errors** - Input format issues, missing headers, etc.
- **Parser errors** - Logic errors during AST construction
- **File errors** - File not found, permission issues, etc.

## Performance Notes

- The parser is designed for typical study guide sizes (hundreds to thousands of lines)
- For very large files (>10MB), consider streaming or chunked processing
- Memory usage scales linearly with input size
- Parsing time is typically O(n) where n is the number of lines

## Future Extensibility

The `Metadata` struct allows for future extensions without breaking the API:

```go
// Future options could include:
metadata := processor.NewMetadata(processor.Colleges).
    WithOption("validation_level", "strict").
    WithOption("output_format", "json").
    WithOption("include_comments", "true").
    WithOption("max_questions", "100")
```

This design ensures backward compatibility while allowing new features to be added seamlessly.

## Study Guide Format

### File Structure

Study guides must follow this format:

```
[File Header]                    # First line: "Study Guide: Title"
[Header]                         # Optional: "Subject: Topic: Subtopic"
[Empty Line]                     # Optional: blank lines for spacing
[Question]                       # "1. Question text? - Answer text"
[Content]                        # Optional: additional content under questions
[Learn More]                     # Optional: "Learn More: Additional info"
```

### Supported Line Types

| Type        | Format                     | Example                                               |
| ----------- | -------------------------- | ----------------------------------------------------- |
| File Header | `Study Guide: Title`       | `Study Guide: Mathematics`                            |
| Header      | `Subject: Topic: Subtopic` | `Colleges: Virginia: ODU: MATH 101: Linear Equations` |
| Question    | `1. Question? - Answer`    | `1. What is x? - A variable`                          |
| Content     | Any text                   | `This is additional content`                          |
| Learn More  | `Learn More: Info`         | `Learn More: See Khan Academy`                        |
| Comment     | `# Comment`                | `# This is a comment`                                 |
| Empty       | Blank line                 | ``                                                    |

### Parser-Specific Formats

#### Colleges Format

```
Study Guide: [Title]
Colleges: [State]: [University]: [Department]: [Course]: [Topic]
[Questions and content...]
```

#### AP Exams Format

```
Study Guide: [Title]
Advanced Placement (AP): [Exam]: [Topic]: [Subtopic]
[Questions and content...]
```

## Output Structure

The parser returns an `AbstractSyntaxTree` with this structure:

```go
type AbstractSyntaxTree struct {
    ParserType string    // "colleges", "ap_exams", etc.
    Timestamp  string    // RFC3339 timestamp
    Root       *Node     // Root node (file header)
}

type Node struct {
    Type     string                 // Token type
    Data     interface{}            // Parsed data
    Children []*Node               // Child nodes
    Parent   *Node                 // Parent node (not in JSON)
}
```

### Example JSON Output

```json
{
  "parser_type": "colleges",
  "timestamp": "2024-01-15T10:30:00Z",
  "root": {
    "type": "file_header",
    "data": {
      "FileHeader": {
        "Title": "Mathematics Study Guide"
      }
    },
    "children": [
      {
        "type": "header",
        "data": {
          "Header": {
            "Parts": [
              "Colleges",
              "Virginia",
              "ODU",
              "MATH 101",
              "Linear Equations"
            ]
          }
        },
        "children": [
          {
            "type": "question",
            "data": {
              "Question": {
                "QuestionText": "What is a linear equation?",
                "AnswerText": "An equation where the highest power of the variable is 1."
              }
            }
          }
        ]
      }
    ]
  }
}
```

## Best Practices

### 1. Always Check Errors

```go
ast, err := processor.Parse(lines, metadata)
if err != nil {
    // Handle error appropriately
    return err
}
```

### 2. Use Appropriate Parser Type

```go
// For college study guides
metadata := processor.NewMetadata(processor.Colleges)

// For AP exams
metadata := processor.NewMetadata(processor.APExams)
```

### 3. Validate Input Format

```go
if len(lines) == 0 {
    return fmt.Errorf("no lines to parse")
}

if !strings.HasPrefix(lines[0], "Study Guide:") {
    return fmt.Errorf("first line must be a file header")
}
```

### 4. Handle Large Files

```go
// For very large files, consider processing in chunks
// or using streaming if available in future versions
```

### 5. Debugging

```go
// Use step-by-step processing for debugging
lexerOutput, err := processor.Lex(lines)
if err != nil {
    log.Printf("Lexer error: %v", err)
    return
}

if !lexerOutput.Success {
    for _, err := range lexerOutput.Errors {
        log.Printf("Lexer error: %v", err)
    }
}
```

## Common Error Messages

| Error                              | Cause                             | Solution                                  |
| ---------------------------------- | --------------------------------- | ----------------------------------------- |
| `first line must be a file header` | Missing or invalid file header    | Ensure first line is "Study Guide: Title" |
| `no lines to parse`                | Empty input                       | Check that input contains lines           |
| `unexpected [type] under [parent]` | Invalid structure                 | Check study guide format                  |
| `question without valid parent`    | Question not under header/passage | Ensure proper nesting                     |

## Migration from Old API

If you were using the old API with manual pipeline management:

**Old (Complex):**

```go
// 20+ lines of boilerplate
tokens := getScannerOutput()
var lines []preparser.ParsedLineInfo
for _, token := range tokens {
    parsedLine := preparser.ParsedLineInfo{
        Number:      token.Number,
        Text:        token.Text,
        Type:        convertTokenType(token.Type),
        ParsedValue: convertParsedValue(token.ParsedValue, token.Type),
    }
    lines = append(lines, parsedLine)
}
parser := parser.NewParser(parserType, lines)
ast, err := parser.Parse()
```

**New (Simple):**

```go
// Just 3 lines!
metadata := processor.NewMetadata(processor.Colleges)
ast, err := processor.Parse(lines, metadata)
if err != nil {
    log.Fatal(err)
}
```

## Support

For questions or issues:

1. Check this guide first
2. Review the test examples in `core/processor/processor_test.go`
3. Check the main README.md for project information
4. Create an issue with example input and expected output
