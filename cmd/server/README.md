# Development Server

This is a development-only web server for testing the study guides parser. **Do not expose this server to clients or production environments.**

## Features

- **Web Interface**: Paste study guide content and see parsing results in real-time
- **Multiple Parser Types**: Test with different parser types (colleges, AP exams, certifications, etc.)
- **Example Content**: Pre-loaded examples for each parser type
- **Two Modes**:
  - **Parse**: Shows the Abstract Syntax Tree (AST)
  - **Build**: Shows the built tree structure
- **Error Display**: Shows detailed parsing errors with line numbers

## Usage

### Starting the Server

```bash
# Option 1: Using make
make server

# Option 2: Direct command
go run cmd/server/main.go
```

The server will start on `http://localhost:8000`

### Using the Web Interface

1. **Select Parser Type**: Choose the appropriate parser type from the dropdown
2. **Load Example** (optional): Click one of the example buttons to load sample content
3. **Enter Content**: Paste your study guide content in the text area
4. **Parse or Build**:
   - Click "Parse Content" to see the AST
   - Click "Build Tree" to see the built tree structure
5. **View Results**: Results are displayed in JSON format below the form

### Example Study Guide Format

```
Study Guide Title
Parser: Type: Category: Subcategory: Topic

1. Question text? - Answer text.
2. Another question? - Another answer.

Learn More: Additional information or resources.
```

## Supported Parser Types

- **Colleges**: College study guides
- **AP Exams**: Advanced Placement exam study guides
- **Certifications**: Professional certification study guides
- **DOD**: Department of Defense study guides
- **Entrance Exams**: College entrance exam study guides

## API Endpoints

The server also provides REST API endpoints for programmatic access:

### POST /parse

Parses content and returns the AST.

**Request:**

```json
{
  "content": "Study guide content...",
  "parser_type": "colleges"
}
```

**Response:**

```json
{
  "success": true,
  "ast": {
    /* Abstract Syntax Tree */
  },
  "errors": []
}
```

### POST /build

Builds content and returns the tree structure.

**Request:**

```json
{
  "content": "Study guide content...",
  "parser_type": "colleges"
}
```

**Response:**

```json
{
  "success": true,
  "tree": {
    /* Built Tree */
  },
  "errors": []
}
```

## Security Notice

⚠️ **This server is for development and testing only.** It has no authentication, rate limiting, or security measures. Do not expose it to the internet or use it in production environments.
