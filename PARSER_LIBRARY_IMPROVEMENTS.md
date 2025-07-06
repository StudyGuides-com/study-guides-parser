# Study Guides Parser Library Improvements

This document outlines suggestions for improving the `github.com/StudyGuides-com/study-guides-parser` library to make it easier to use and reduce the heavy lifting currently required when working with the library.

## Current Pain Points

The current library requires significant manual conversion between different token formats:

- **Lexer**: `LineInfo` with basic fields
- **Scanner**: `ParsedLineInfo` with `ParsedValue`
- **Parser**: Expects `preparser.ParsedLineInfo`

This forces users to write conversion code like:

```go
// Heavy lifting example from our implementation
func convertTokenType(typeStr string) lexer.TokenType {
    switch typeStr {
    case "file_header":
        return lexer.TokenTypeFileHeader
    // ... many more cases
    }
}

func convertParsedValue(rawValue interface{}, tokenType string) preparser.ParsedValue {
    // Complex conversion logic for each token type
}
```

## Suggested Improvements

### 1. **Standardize Token Formats Across Packages**

**Problem**: Each package has different token formats, requiring manual conversion.

**Solution**: Create a unified token format that all packages can use:

```go
// In a shared package like core/types
type Token struct {
    Number      int         `json:"number"`
    Text        string      `json:"text"`
    Type        TokenType   `json:"type"`
    ParsedValue interface{} `json:"parsed_value,omitempty"`
}

// All packages (lexer, scanner, parser) use this same format
```

### 2. **Add Convenience Functions for Common Workflows**

**Problem**: Users have to manually handle the entire pipeline and format conversions.

**Solution**: Add high-level functions for common workflows:

```go
// In core/processor or similar
func ProcessStudyGuide(lines []string, parserType ParserType) (*AbstractSyntaxTree, error) {
    // Handle the entire pipeline: lexer → scanner → parser
    // Returns the final AST
}

func ProcessTokens(tokens []Token, parserType ParserType) (*AbstractSyntaxTree, error) {
    // Handle scanner output → parser
    // No manual conversion needed
}

func ProcessFile(filename string, parserType ParserType) (*AbstractSyntaxTree, error) {
    // Read file and process in one call
}
```

### 3. **Improve Type Safety and Reduce Interface{} Usage**

**Problem**: Heavy use of `interface{}` makes it hard to work with parsed values.

**Solution**: Use concrete types and generics:

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

// Instead of interface{}, use concrete types
type Token struct {
    Number      int         `json:"number"`
    Text        string      `json:"text"`
    Type        TokenType   `json:"type"`
    ParsedValue ParsedValue `json:"parsed_value,omitempty"`
}
```

### 4. **Add JSON Marshaling/Unmarshaling Support**

**Problem**: Manual conversion between JSON and structs.

**Solution**: Add proper JSON tags and helper functions:

```go
// Add JSON tags to all structs
type ParsedLineInfo struct {
    Number      int         `json:"number"`
    Text        string      `json:"text"`
    Type        TokenType   `json:"type"`
    ParsedValue ParsedValue `json:"parsed_value,omitempty"`
}

// Add helper functions for common conversions
func UnmarshalTokens(data []byte) ([]Token, error) {
    var tokens []Token
    err := json.Unmarshal(data, &tokens)
    return tokens, err
}

func MarshalAST(ast *AbstractSyntaxTree) ([]byte, error) {
    return json.MarshalIndent(ast, "", "  ")
}

// Handle the scanner → parser conversion automatically
func TokensToParsedLineInfo(tokens []Token) []ParsedLineInfo {
    // Automatic conversion without manual mapping
}
```

### 5. **Create a Builder Pattern for Parser Configuration**

**Problem**: Parser instantiation is basic and doesn't allow for configuration.

**Solution**: Add a builder pattern:

```go
parser := NewParserBuilder().
    WithType(ParserType("Colleges")).
    WithTokens(tokens).
    WithOptions(ParserOptions{
        StrictMode: true,
        ValidateStructure: true,
        SkipEmptyLines: false,
    }).
    WithErrorHandler(func(err error) {
        log.Printf("Parser error: %v", err)
    }).
    Build()

ast, err := parser.Parse()
```

### 6. **Add Better Error Handling and Validation**

**Problem**: Errors are basic and don't provide enough context.

**Solution**: Enhanced error types with more context:

```go
type ParsingError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Line    int       `json:"line"`
    Token   string    `json:"token"`
    Context string    `json:"context"`
    Stack   string    `json:"stack,omitempty"`
}

type ValidationError struct {
    Field   string `json:"field"`
    Value   string `json:"value"`
    Rule    string `json:"rule"`
    Message string `json:"message"`
}

// Better error handling
if err != nil {
    if parseErr, ok := err.(*ParsingError); ok {
        fmt.Printf("Error at line %d: %s\n", parseErr.Line, parseErr.Message)
    }
}
```

### 7. **Provide Examples and Documentation**

**Problem**: Limited examples of how to use the library.

**Solution**: Add comprehensive examples:

```go
// examples/basic_usage.go
func ExampleBasicUsage() {
    lines := []string{
        "Study Guide: Mathematics",
        "Subject: Algebra",
        "1. What is x? - A variable",
    }

    ast, err := ProcessStudyGuide(lines, ParserType("Colleges"))
    if err != nil {
        log.Fatal(err)
    }

    jsonData, _ := json.Marshal(ast)
    fmt.Println(string(jsonData))
}

// examples/advanced_usage.go
func ExampleAdvancedUsage() {
    // Custom parser configuration
    parser := NewParserBuilder().
        WithType(ParserType("APExams")).
        WithStrictMode(true).
        Build()

    // Process tokens from scanner
    tokens := getTokensFromScanner()
    ast, err := parser.ParseTokens(tokens)
}
```

### 8. **Add Version Compatibility**

**Problem**: No clear versioning strategy for breaking changes.

**Solution**: Add version compatibility helpers:

```go
func IsCompatibleVersion(version string) bool {
    // Check if the version is compatible
    return strings.HasPrefix(version, "v0.")
}

func MigrateTokens(tokens []Token, fromVersion, toVersion string) ([]Token, error) {
    // Handle migration between versions
    // Useful for backward compatibility
}

// Version-aware processing
func ProcessWithVersion(lines []string, parserType ParserType, version string) (*AbstractSyntaxTree, error) {
    if !IsCompatibleVersion(version) {
        return nil, fmt.Errorf("incompatible version: %s", version)
    }
    return ProcessStudyGuide(lines, parserType)
}
```

### 9. **Create a Simple CLI Tool**

**Problem**: No easy way to test the library.

**Solution**: Add a CLI tool for testing:

```bash
# Process a file directly
study-guides-parser process input.txt --parser-type Colleges --output ast.json

# Validate a file
study-guides-parser validate input.txt

# Convert between formats
study-guides-parser convert --from lexer --to parser input.json

# Interactive mode
study-guides-parser interactive
```

### 10. **Add Streaming Support**

**Problem**: Only supports batch processing.

**Solution**: Add streaming for large files:

```go
func ProcessStream(reader io.Reader, parserType ParserType) (<-chan *AbstractSyntaxTree, <-chan error) {
    astChan := make(chan *AbstractSyntaxTree)
    errChan := make(chan error)

    go func() {
        defer close(astChan)
        defer close(errChan)

        // Process in chunks
        // Send ASTs as they're completed
    }()

    return astChan, errChan
}

// Usage
astChan, errChan := ProcessStream(file, ParserType("Colleges"))
for ast := range astChan {
    // Process each AST
}
```

## Priority Recommendations

### High Priority (Implement First)

1. **Standardize token formats** (#1) - Eliminates most conversion code
2. **Add convenience functions** (#2) - Simplifies common workflows
3. **Improve type safety** (#3) - Reduces interface{} usage

### Medium Priority

4. **Add JSON support** (#4) - Better serialization
5. **Builder pattern** (#5) - More flexible configuration
6. **Enhanced error handling** (#6) - Better debugging

### Low Priority

7. **Examples and documentation** (#7) - Better developer experience
8. **Version compatibility** (#8) - Future-proofing
9. **CLI tool** (#9) - Testing and validation
10. **Streaming support** (#10) - Performance for large files

## Implementation Strategy

1. **Phase 1**: Implement unified token format and convenience functions
2. **Phase 2**: Add type safety improvements and JSON helpers
3. **Phase 3**: Add builder pattern and error handling
4. **Phase 4**: Documentation, examples, and CLI tool

## Expected Benefits

After implementing these improvements:

- **90% reduction** in boilerplate conversion code
- **Simplified API** that's easier to understand and use
- **Better type safety** with fewer runtime errors
- **Improved developer experience** with examples and tools
- **Backward compatibility** through version management

## Example Usage After Improvements

```go
// Before (current heavy lifting)
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

// After (proposed improvements)
tokens := getScannerOutput()
ast, err := ProcessTokens(tokens, ParserType("Colleges"))
```

This would make the library much more developer-friendly and reduce the cognitive load when working with study guide parsing.
