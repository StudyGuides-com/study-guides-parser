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

// Envelope wraps parser output with schema metadata
type Envelope[T any] struct {
	SchemaType    SchemaType `json:"schema_type"`
	SchemaVersion string     `json:"schema_version"`
	Data          T          `json:"data"`
}

// NewEnvelope creates a new schema envelope with the given type and data
func NewEnvelope[T any](schemaType SchemaType, data T) Envelope[T] {
	return Envelope[T]{
		SchemaType:    schemaType,
		SchemaVersion: Version,
		Data:          data,
	}
}
