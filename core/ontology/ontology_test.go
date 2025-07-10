package ontology

import "testing"

func TestFindTagOntology(t *testing.T) {
	tests := []struct {
		name        string
		contextType ContextType
		depth       int
		expected    *TagOntology
	}{
		{
			name:        "AP Exams depth 3",
			contextType: ContextTypeAPExams,
			depth:       3,
			expected: &TagOntology{
				ContextType:  ContextTypeAPExams,
				HeaderLength: 3,
				TagTypes:     []TagType{TagTypeCategory, TagTypeAPExam, TagTypeTopic},
			},
		},
		{
			name:        "Certifications depth 4",
			contextType: ContextTypeCertifications,
			depth:       4,
			expected: &TagOntology{
				ContextType:  ContextTypeCertifications,
				HeaderLength: 4,
				TagTypes:     []TagType{TagTypeCategory, TagTypeCertifyingAgency, TagTypeCertification, TagTypeTopic},
			},
		},
		{
			name:        "Non-existent combination",
			contextType: ContextTypeColleges,
			depth:       2,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindTagOntology(tt.contextType, tt.depth)
			
			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %+v", result)
				}
				return
			}
			
			if result == nil {
				t.Errorf("Expected %+v, got nil", tt.expected)
				return
			}
			
			if result.ContextType != tt.expected.ContextType {
				t.Errorf("ContextType mismatch: expected %s, got %s", tt.expected.ContextType, result.ContextType)
			}
			
			if result.HeaderLength != tt.expected.HeaderLength {
				t.Errorf("HeaderLength mismatch: expected %d, got %d", tt.expected.HeaderLength, result.HeaderLength)
			}
			
			if len(result.TagTypes) != len(tt.expected.TagTypes) {
				t.Errorf("TagTypes length mismatch: expected %d, got %d", len(tt.expected.TagTypes), len(result.TagTypes))
				return
			}
			
			for i, tagType := range tt.expected.TagTypes {
				if result.TagTypes[i] != tagType {
					t.Errorf("TagTypes[%d] mismatch: expected %s, got %s", i, tagType, result.TagTypes[i])
				}
			}
		})
	}
} 