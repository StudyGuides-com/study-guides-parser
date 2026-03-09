# Study Guides Parser

A Go library and development server for parsing educational study guide content into structured hierarchical data. Converts plain text documents into Abstract Syntax Trees (ASTs) and tree structures containing tags, questions, and passages.

## Installation

```bash
go get github.com/studyguides-com/study-guides-parser
```

**Requirements:** Go 1.20+

## Quick Start

### Programmatic Usage

```go
import (
    "github.com/studyguides-com/study-guides-parser/core/processor"
    "github.com/studyguides-com/study-guides-parser/core/config"
)

// Parse from file
result, err := processor.ParseFile("study_guide.txt", config.NewMetadata("colleges"))

// Parse from strings
lines := []string{
    "Mathematics Study Guide",
    "Colleges: Virginia: ODU: MATH 101: Linear Equations",
    "1. What is x? - A variable",
}
result, err := processor.Parse(lines, config.NewMetadata("colleges"))

// Full pipeline to Tree (with tags, questions, passages)
result, err := processor.Build(lines, config.NewMetadata("colleges"))
tree := result.Tree
```

### Available Functions

| Function | Description |
|----------|-------------|
| `processor.ParseFile(filename, metadata)` | Parse file to AST |
| `processor.Parse(lines, metadata)` | Parse string slice to AST |
| `processor.Build(lines, metadata)` | Full pipeline to Tree structure |
| `processor.Preparse(lines, metadata)` | Tokenize and parse values |
| `processor.Lex(lines, metadata)` | Lexical analysis only |

## Input Format

Study guides are plain text files with a specific structure:

```
Mathematics Study Guide
Colleges: Virginia: Old Dominion University (ODU): Mathematics (MATH): MATH 101: Linear Equations

1. What is a linear equation? - An equation where the highest power of the variable is 1.
2. How do you solve 2x + 3 = 7? - Subtract 3 from both sides, then divide by 2.

Learn More: See Khan Academy's linear equations course.

Passage: Introduction to Linear Systems

A linear system consists of two or more linear equations.

1. What defines a linear system? - Two or more linear equations
```

### Line Types

| Type | Format | Example |
|------|--------|---------|
| File Header | First line of document | `Mathematics Study Guide` |
| Header | Colon-separated hierarchy | `Colleges: Virginia: ODU: MATH 101` |
| Question | `N. Question? - Answer` | `1. What is x? - A variable` |
| Passage | `Passage: Title` | `Passage: Introduction` |
| Learn More | `Learn More: Text` | `Learn More: See Khan Academy` |
| Content | Body text | Any regular text |
| Comment | Lines starting with `#` | `# This is a comment` |

## Context Types

Create metadata with the appropriate context type:

```go
config.NewMetadata("colleges")       // College study guides
config.NewMetadata("certifications") // Professional certifications
config.NewMetadata("ap_exams")       // AP exam prep
config.NewMetadata("entrance_exams") // College entrance exams
config.NewMetadata("dod")            // Department of Defense materials
```

### Tag Hierarchies by Context

Each context type has its own tag hierarchy:

| Context | Hierarchy |
|---------|-----------|
| Colleges | Category > Region > University > Department > Course > Topic |
| Certifications | Category > Certifying Agency > Certification > Domain > Module |
| AP Exams | Category > AP Exam > Domain > Part > Topic |
| Entrance Exams | Category > Entrance Exam > Section > Topic |
| DoD | Category > Branch > Instruction Type > Instruction Group > Instruction > Chapter |

## Output Structure

### Tree Structure

```go
type Tree struct {
    Root     *Tag
    Metadata *Metadata
}

type Tag struct {
    Title              string
    TagType            TagType       // Category, Topic, Course, etc.
    InsertID           string        // CUID for database insertion
    Hash               string        // SHA256 for deduplication
    Context            ContextType
    ContentRating      ContentRatingType
    ContentDescriptors []string
    MetaTags           []string
    Overview           *Overview
    Questions          []*Question
    Passages           []*Passage
    ChildTags          []*Tag
}

type Question struct {
    InsertID    string
    Hash        string
    Prompt      string
    Answer      string
    Distractors []string
    LearnMore   string
}

type Passage struct {
    InsertID  string
    Hash      string
    Title     string
    Content   string
    Questions []*Question
}
```

## Development Server

A web server is included for testing and development:

```bash
# Start the server
make server
# or
go run cmd/server/main.go
```

Server runs at `http://localhost:8000` with an interactive web UI.

**Note:** Development use only. Do not expose to production.

### API Endpoints

All endpoints accept POST with JSON body:

```json
{
  "content": "Your study guide content here",
  "context_type": "Colleges"
}
```

| Endpoint | Description | Returns |
|----------|-------------|---------|
| `POST /lex` | Tokenize text | Line tokens with types |
| `POST /preparse` | Parse token values | Parsed line info |
| `POST /parse` | Build AST | Abstract Syntax Tree |
| `POST /build` | Full pipeline | Complete Tree structure |
| `POST /hash` | Generate hash | SHA256 hash of input |

### Response Format

```json
{
  "success": true,
  "tree": { ... },
  "errors": []
}
```

Errors include line numbers and context:

```json
{
  "success": false,
  "errors": [
    {
      "line_number": 1,
      "message": "first line must be a file header",
      "text": "Some Content",
      "type": "content"
    }
  ]
}
```

## Processing Pipeline

The parser follows a 5-stage pipeline:

```
Text Input --> Lexer --> Preparser --> Parser --> Builder --> Tree
                |           |            |          |          |
             Tokens     Parsed        AST      Tag Types    Final
                        Values                 Assigned     Output
```

1. **Lexer** - Classifies each line by type (header, question, etc.)
2. **Preparser** - Extracts semantic values from tokens
3. **Parser** - Builds Abstract Syntax Tree with hierarchy
4. **Builder** - Creates domain objects, assigns tag types
5. **Tree** - Final hierarchical structure with hashes and IDs

## Hash Generation

Hashes are used for database deduplication:

- **Category tags**: `HashFrom(title)` - consistent across files
- **Nested tags**: `HashFrom(parentTitle + title)` - unique under parent
- **Questions/Passages**: Hash from content

## Commands

```bash
make fmt      # Format code
make test     # Run tests
make build    # Build binary
make server   # Start dev server

go test ./...              # Run all tests
go test ./core/builder/... # Test specific package
```

## Library Structure

```
core/
├── builder/      # Tree construction from AST
├── config/       # Metadata and configuration
├── idgen/        # Hash and CUID generation
├── lexer/        # Line tokenization
├── ontology/     # Tag types and context types
├── parser/       # AST construction
├── preparser/    # Token value extraction
├── processor/    # High-level API functions
├── qa/           # Validation runner
└── tree/         # Tree data structures
```

## License

[Add license information]
