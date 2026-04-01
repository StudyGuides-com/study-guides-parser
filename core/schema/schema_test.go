package schema

import (
	"testing"
)

func TestAllSchemaTypes(t *testing.T) {
	tests := []struct {
		schemaType SchemaType
		expected   string
	}{
		{SchemaTypeLexer, "lexer"},
		{SchemaTypePreparser, "preparser"},
		{SchemaTypeParser, "parser"},
		{SchemaTypeBuilder, "builder"},
		{SchemaTypeHash, "hash"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.schemaType) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.schemaType)
			}
		})
	}
}

func TestVersionConstant(t *testing.T) {
	if Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", Version)
	}
}
