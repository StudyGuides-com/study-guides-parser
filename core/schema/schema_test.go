package schema

import (
	"encoding/json"
	"testing"
)

func TestNewEnvelope(t *testing.T) {
	type TestData struct {
		Value string `json:"value"`
	}

	data := TestData{Value: "test"}
	envelope := NewEnvelope(SchemaTypeLexer, data)

	if envelope.SchemaType != SchemaTypeLexer {
		t.Errorf("expected schema type %s, got %s", SchemaTypeLexer, envelope.SchemaType)
	}

	if envelope.SchemaVersion != Version {
		t.Errorf("expected schema version %s, got %s", Version, envelope.SchemaVersion)
	}

	if envelope.Data.Value != "test" {
		t.Errorf("expected data value 'test', got '%s'", envelope.Data.Value)
	}
}

func TestEnvelopeJSONSerialization(t *testing.T) {
	type TestData struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	data := TestData{Success: true, Message: "hello"}
	envelope := NewEnvelope(SchemaTypeBuilder, data)

	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("failed to marshal envelope: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if result["schema_type"] != "builder" {
		t.Errorf("expected schema_type 'builder', got '%v'", result["schema_type"])
	}

	if result["schema_version"] != "1.0.0" {
		t.Errorf("expected schema_version '1.0.0', got '%v'", result["schema_version"])
	}

	dataMap, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be an object")
	}

	if dataMap["success"] != true {
		t.Errorf("expected data.success to be true, got %v", dataMap["success"])
	}

	if dataMap["message"] != "hello" {
		t.Errorf("expected data.message to be 'hello', got %v", dataMap["message"])
	}
}

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

			envelope := NewEnvelope(tt.schemaType, struct{}{})
			if envelope.SchemaType != tt.schemaType {
				t.Errorf("envelope schema type mismatch: expected %s, got %s", tt.schemaType, envelope.SchemaType)
			}
		})
	}
}

func TestVersionConstant(t *testing.T) {
	if Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", Version)
	}
}
