package schema

// SchemaType identifies the type of parser output
type SchemaType string

const (
	SchemaTypeLexer     SchemaType = "lexer"
	SchemaTypePreparser SchemaType = "preparser"
	SchemaTypeParser    SchemaType = "parser"
	SchemaTypeBuilder   SchemaType = "builder"
	SchemaTypeHash      SchemaType = "hash"
)

// Version is the current schema version
const Version = "1.0.0"
