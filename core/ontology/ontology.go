package ontology

// FindTagOntology finds the appropriate TagOntology entry for a given context and depth
func FindTagOntology(contextType ContextType, depth int) *TagOntology {
	for _, ontology := range tagOntology {
		if ontology.ContextType == contextType && ontology.HeaderLength == depth {
			return &ontology
		}
	}
	return nil
}

// AssignTagType assigns the appropriate tag type based on context and depth
func AssignTagType(tag interface{}, contextType ContextType, depth int) {
	// This function will be used by the builder package
	// The actual implementation will depend on the tag interface
	ontology := FindTagOntology(contextType, depth)
	if ontology != nil && depth <= len(ontology.TagTypes) {
		// The builder package will implement the actual tag type assignment
		// This is just a helper function to find the ontology
	}
} 