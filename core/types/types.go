package types

// Metadata contains configuration and metadata for parsing
type Metadata struct {
	Type    string            `json:"type"`
	Options map[string]string `json:"options,omitempty"`
}



// NewMetadata creates a new Metadata struct with the given type
func NewMetadata(typeName string) *Metadata {
	return &Metadata{
		Type:    typeName,
		Options: make(map[string]string),
	}
}

// WithOption adds an option to the metadata
func (m *Metadata) WithOption(key, value string) *Metadata {
	m.Options[key] = value
	return m
} 