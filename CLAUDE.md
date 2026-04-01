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
6. **Processor** (`core/processor/`) - Orchestrates pipeline stages, used by HTTP handlers

### Key Components

#### Tag System (`core/tree/tag.go`)
- Tags form hierarchical structures based on context type (see Ontology below)
- Each tag has unique hash for database deduplication
- InsertID is pre-generated CUID for database insertion

#### Hash Generation (`core/idgen/`)
- Top-level tags: `HashFrom(title)` - consistent across files
- Nested tags: `HashFrom(parentTitle + title)` - ensures uniqueness under parent
- Questions/Passages: Hash from content for deduplication

#### Ontology System (`core/ontology/`)
Tag types are assigned based on context and tree depth. Each context type defines its own hierarchy:
- **Colleges**: Category → Region → University → Department → Course → Topic
- **Certifications**: Category → CertifyingAgency → Certification → Domain → Module → Topic
- **EntranceExams**: Category → EntranceExam → Module → Topic
- **APExams**: Category → APExam → Module → Topic
- **DoD**: Category → Branch → InstructionType → InstructionGroup → Instruction → Section
- **Encyclopedia**: Category → Volume → Range → Topic
- **General**: Category → SubCategory → Topic

#### Schema Versioning (`core/schema/`)
All API responses include `schema_type` and `schema_version` for client compatibility.

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

All endpoints accept `context_type` in the request body for tag type assignment.

## Context Types
Valid context types for tag assignment:
- Colleges
- Certifications
- EntranceExams
- APExams
- UserGeneratedContent
- DoD
- Encyclopedia
- General
- None (default)

## Important Implementation Details

### Tag Type Assignment
Tags are assigned types based on depth and context via `tree.AssignTagTypes()`. The ontology system (`core/ontology/data.go`) defines the mapping.

### QA System
QA runner (`core/qa/`) validates:
- All tags have proper TagType (not TagTypeNone)
- All tags have valid ContextType (not ContextTypeNone)

### Hash Stability
Critical for database: hashes must be consistent for deduplication to work. Top-level (category) tags are hashed without parent to ensure consistency across files (see `builder.go:178-192`).