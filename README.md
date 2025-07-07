# Study Guides Parser Library

A Go library for parsing study guides into structured Abstract Syntax Trees (ASTs). Currently supports multiple formats including college study guides, AP exams, certifications, and more.

## Current State

The library provides a robust parsing pipeline but requires significant manual work from users to handle format conversions and pipeline orchestration.

### Current Pain Points

- **Multiple incompatible token formats**: Users must manually convert between `lexer.LineInfo`, `preparser.ParsedLineInfo`, and different token types
- **Heavy boilerplate**: Users need to write conversion functions like `convertTokenType()` and `convertParsedValue()`
- **Complex pipeline management**: Users must manually orchestrate lexer → scanner → parser flow
- **Poor type safety**: Heavy use of `interface{}` makes it hard to work with parsed values

### Current Usage (Complex)

```go
// Current heavy lifting required
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
parser := parser.NewParser(lines)
ast, err := parser.Parse(parserType)
```

## Improvements Status

### ✅ Phase 1 Complete: Simple API (Implemented)

We've successfully implemented the simple API! Users can now parse study guides with just one function call instead of the previous 20+ lines of boilerplate.

**New Simple Usage:**

```go
// Super simple - just one function call
ast, err := processor.ParseFile("study_guide.txt", processor.Colleges)
if err != nil {
    log.Fatal(err)
}

// Or from strings
lines := []string{"Mathematics Study Guide", "Colleges: Virginia: ODU: MATH 101: Linear Equations", "1. What is x? - A variable"}
ast, err := processor.ParseLines(lines, processor.Colleges)
```

**Available Functions:**

- `processor.ParseFile(filename, parserType)` - Parse directly from a file
- `processor.ParseLines(lines, parserType)` - Parse from string slices
- `processor.ParseTokens(tokens, parserType)` - Parse from pre-processed tokens

**Supported Parser Types:**

- `processor.Colleges` - College study guides
- `processor.APExams` - AP exam study guides
- `processor.Certifications` - Certification study guides
- `processor.DOD` - Department of Defense study guides
- `processor.EntranceExams` - Entrance exam study guides

### 🔄 Remaining Phases

We're working to make this library even better. The remaining phases will:

### Target Usage (Simple)

```go
// Super simple - just one function call
ast, err := processor.ParseFile("study_guide.txt", processor.Colleges)
if err != nil {
    log.Fatal(err)
}

// Output as JSON
jsonData, _ := json.MarshalIndent(ast, "", "  ")
fmt.Println(string(jsonData))
```

## Improvement Plan

### Phase 1: Create Simple API (Week 1) - **HIGHEST PRIORITY**

Create high-level functions that handle everything internally:

```go
// In core/processor/processor.go
func ParseFile(filename string, parserType ParserType) (*AbstractSyntaxTree, error)
func ParseLines(lines []string, parserType ParserType) (*AbstractSyntaxTree, error)
func ParseTokens(tokens []Token, parserType ParserType) (*AbstractSyntaxTree, error)
```

**Benefits:**

- 90% reduction in boilerplate code
- Single function call for most use cases
- Automatic internal conversions

### Phase 2: Unify Token Types (Week 2)

Create a unified token format that works across all packages:

```go
// In core/types/token.go
type Token struct {
    Number      int         `json:"number"`
    Text        string      `json:"text"`
    Type        TokenType   `json:"type"`
    ParsedValue ParsedValue `json:"parsed_value,omitempty"`
}

type TokenType string

const (
    TokenTypeFileHeader TokenType = "file_header"
    TokenTypeHeader     TokenType = "header"
    TokenTypeQuestion   TokenType = "question"
    TokenTypeContent    TokenType = "content"
    // ... etc
)
```

**Benefits:**

- No more manual token type conversions
- Consistent format across all packages
- Better JSON serialization

### Phase 3: Improve Type Safety (Week 3)

Replace `interface{}` with concrete types:

```go
type ParsedValue struct {
    FileHeader *FileHeaderResult `json:"file_header,omitempty"`
    Header     *HeaderResult     `json:"header,omitempty"`
    Question   *QuestionResult   `json:"question,omitempty"`
    Comment    *CommentResult    `json:"comment,omitempty"`
    Passage    *PassageResult    `json:"passage,omitempty"`
    LearnMore  *LearnMoreResult  `json:"learn_more,omitempty"`
    Content    *ContentResult    `json:"content,omitempty"`
    Empty      *EmptyLineResult  `json:"empty,omitempty"`
    Binary     *BinaryResult     `json:"binary,omitempty"`
}
```

**Benefits:**

- Better type safety
- Easier to work with parsed values
- Better IDE support and autocomplete

### Phase 4: Add Convenience Features (Week 4)

- **Builder pattern** for advanced configuration
- **Better error handling** with context and line numbers
- **CLI tool** for testing and validation
- **JSON helpers** for easy serialization

## Supported Parser Types

```go
type ParserType string

const (
    Colleges        ParserType = "colleges"
    APExams         ParserType = "ap_exams"
    Certifications  ParserType = "certifications"
    DOD             ParserType = "dod"
    EntranceExams   ParserType = "entrance_exams"
)
```

## Example Study Guide Format

```
Study Guide: Mathematics
Subject: Algebra
Topic: Linear Equations

1. What is a linear equation?
   A linear equation is an equation where the highest power of the variable is 1.

2. How do you solve 2x + 3 = 7?
   Subtract 3 from both sides: 2x = 4
   Divide both sides by 2: x = 2

Learn More: See Khan Academy's linear equations course.
```

## Current Library Structure

```
core/
├── lexer/          # Token classification
├── preparser/      # Token parsing and value extraction
├── parser/         # AST construction
├── cleanstring/    # Text cleaning utilities
├── constants/      # Shared constants
├── regexes/        # Regular expression patterns
└── utils/          # Utility functions
```

## Development

### Prerequisites

- Go 1.20 or higher

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

### Running Examples

```bash
# Process a study guide file
go run examples/basic_usage.go input.txt

# Validate a file
go run examples/validate.go input.txt
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Priorities

1. **Phase 1**: Implement simple API functions
2. **Phase 2**: Unify token types across packages
3. **Phase 3**: Improve type safety
4. **Phase 4**: Add convenience features

## License

[Add your license information here]

## Support

For questions, issues, or contributions, please:

1. Check the [Issues](https://github.com/StudyGuides-com/study-guides-parser/issues) page
2. Create a new issue with a clear description
3. Include example input and expected output

---

**Note**: This library is actively being improved to provide a super simple API. The current version requires manual pipeline management, but future versions will provide one-line parsing capabilities.
