package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/idgen"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
	"github.com/studyguides-com/study-guides-parser/core/processor"
	"github.com/studyguides-com/study-guides-parser/core/schema"
)

type ParseRequest struct {
	Content     string `json:"content" binding:"required"`
	ContextType string `json:"context_type"`
}

type HashRequest struct {
	Value string `json:"value" binding:"required"`
}

type HashResponse struct {
	Hash string `json:"hash"`
}

func handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "template.html", gin.H{})
}

func handleLex(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	lines := strings.Split(req.Content, "\n")
	metadata := config.NewMetadata("lex")
	result, err := processor.LexWithSchema(lines, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lexing error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func handlePreparse(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	lines := strings.Split(req.Content, "\n")
	metadata := config.NewMetadata("preparse")
	result, err := processor.PreparseWithSchema(lines, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Preparsing error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func handleParse(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	lines := strings.Split(req.Content, "\n")
	metadata := config.NewMetadata("parse")

	// Set context type if provided
	if req.ContextType != "" {
		// Validate context type
		if !isValidContextType(req.ContextType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid context type: " + req.ContextType})
			return
		}
		contextType := ontology.ContextType(req.ContextType)
		metadata.ContextType = contextType
	}

	result, err := processor.ParseWithSchema(lines, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Parsing error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func handleBuild(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	lines := strings.Split(req.Content, "\n")
	metadata := config.NewMetadata("build")

	// Set context type if provided
	if req.ContextType != "" {
		// Validate context type
		if !isValidContextType(req.ContextType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid context type: " + req.ContextType})
			return
		}
		contextType := ontology.ContextType(req.ContextType)
		metadata.ContextType = contextType
	}

	result, err := processor.BuildWithSchema(lines, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Build error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func handleHash(c *gin.Context) {
	var req HashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	hash := idgen.HashFrom(req.Value)
	response := schema.NewEnvelope(schema.SchemaTypeHash, HashResponse{Hash: hash})
	c.JSON(http.StatusOK, response)
}

// isValidContextType validates that the provided context type is valid
func isValidContextType(contextType string) bool {
	validTypes := []string{
		"Colleges",
		"Certifications",
		"EntranceExams",
		"APExams",
		"UserGeneratedContent",
		"DoD",
		"Encyclopedia",
		"General",
		"None",
	}

	for _, validType := range validTypes {
		if contextType == validType {
			return true
		}
	}
	return false
}
