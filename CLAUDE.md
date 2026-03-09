# Study Guides Parser - Project Documentation

## Overview
A Go-based parser system for processing educational study guide content into structured hierarchical data (tags, questions, passages).

## Core Architecture

### Processing Pipeline
1. **Lexer** (`core/lexer/`) - Tokenizes raw text into typed tokens
2. **Preparser** (`core/preparser/`) - Parses tokens into structured data
3. **Parser** (`core/parser/`) - Builds Abstract Syntax Tree (AST)
4. **Builder** (`core/builder/`) - Constructs hierarchical tree from AST
5. **Tree** (`core/tree/`) - Final hierarchical data structure

### Key Components

#### Tag System (`core/tree/tag.go`)
- Tags form hierarchical structure (Category → Volume → Range → Topic)
- Each tag has unique hash for database deduplication
- InsertID is pre-generated CUID for database insertion

#### Hash Generation (`core/idgen/`)
- Category tags: `HashFrom(title)` - consistent across files
- Nested tags: `HashFrom(parentTitle + title)` - ensures uniqueness under parent
- Questions/Passages: Hash from content for deduplication

## Recent Fixes

### Category Tag Duplication Bug (Fixed v0.2.1)
**Issue**: Category tags were getting different hashes when imported from different files
**Root Cause**: Hash was computed using file header as parent: `HashFrom(fileHeader + "Encyclopedia")`
**Solution**: Modified `builder.go:175-190` to detect root-level tags and hash without parent

## Development Workflow

### Branches
- `main` - Production releases
- `dev` - Active development

### Testing
```bash
go test ./core/builder/... -v  # Test builder logic
go test ./...                   # Run all tests
```

### Release Process
1. Commit fixes to dev branch
2. Merge dev → main (if needed)
3. Tag from main: `git tag -a v0.X.Y -m "Description"`
4. Push tag: `git push origin v0.X.Y`

## API Endpoints
Server runs on `:8000` with endpoints:
- `POST /lex` - Tokenize text
- `POST /preparse` - Parse tokens
- `POST /parse` - Build AST
- `POST /build` - Build complete tree
- `POST /hash` - Generate hash for value

## Context Types
Valid context types for tag assignment:
- Colleges
- Certifications
- EntranceExams
- APExams
- UserGeneratedContent
- DoD
- None (default)

## Important Implementation Details

### Tag Type Assignment
Tags are assigned types based on depth and context via `tree.AssignTagTypes()`

### QA System
QA runner validates:
- All tags have proper TagType
- All tags have valid ContextType
- Content ratings and descriptors (when enhanced)

### Hash Stability
Critical for database: hashes must be consistent for deduplication to work