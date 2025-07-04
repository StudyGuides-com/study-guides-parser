# Study Guides Parser

A Go library for parsing study guide text files into structured data.

## Installation

```bash
go get github.com/StudyGuides-com/study-guides-parser
```

## Usage

```go
package main

import (
    "os"
    "log"

    "github.com/StudyGuides-com/study-guides-parser"
)

func main() {
    parser := studyguidesparser.New()

    file, err := os.Open("study-guide.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    sections, err := parser.ParseFile(file)
    if err != nil {
        log.Fatal(err)
    }

    // Process the parsed sections
    for _, section := range sections {
        fmt.Printf("Title: %s\n", section.Title)
        fmt.Printf("Content: %s\n", section.Content)
    }
}
```

## Development

Run tests:

```bash
go test
```

## License

[Add your license here]
