package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/idgen"
	"github.com/studyguides-com/study-guides-parser/core/ontology"
	"github.com/studyguides-com/study-guides-parser/core/processor"
)

type ParseRequest struct {
	Content     string `json:"content" binding:"required"`
	ContextType string `json:"context_type"`
}

type ParseResponse struct {
	Success bool                        `json:"success"`
	AST     interface{}                 `json:"ast,omitempty"`
	Errors  []processor.ProcessingError `json:"errors,omitempty"`
	Tree    interface{}                 `json:"tree,omitempty"`
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
	result, err := processor.Lex(lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lexing error: " + err.Error()})
		return
	}

	response := ParseResponse{
		Success: result.Success,
		Errors:  result.Errors,
	}
	if result.Success {
		response.AST = result.Tokens
	}

	c.JSON(http.StatusOK, response)
}

func handlePreparse(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	lines := strings.Split(req.Content, "\n")
	result, err := processor.Preparse(lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Preparsing error: " + err.Error()})
		return
	}

	response := ParseResponse{
		Success: result.Success,
		Errors:  result.Errors,
	}
	if result.Success {
		response.AST = result.Tokens
	}

	c.JSON(http.StatusOK, response)
}

func handleParse(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	lines := strings.Split(req.Content, "\n")
	metadata := config.NewMetaData("parse")
	
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

	result, err := processor.Parse(lines, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Parsing error: " + err.Error()})
		return
	}

	response := ParseResponse{
		Success: result.Success,
		Errors:  result.Errors,
	}
	if result.Success && result.AST != nil {
		response.AST = result.AST
	}

	c.JSON(http.StatusOK, response)
}

func handleBuild(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	lines := strings.Split(req.Content, "\n")
	metadata := config.NewMetaData("build")
	
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

	result, err := processor.Build(lines, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Build error: " + err.Error()})
		return
	}

	response := ParseResponse{
		Success: result.Success,
		Errors:  result.Errors,
	}
	if result.Success && result.Tree != nil {
		response.Tree = result.Tree
	}

	c.JSON(http.StatusOK, response)
}

func handleHash(c *gin.Context) {
	var req HashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	hash := idgen.HashFrom(req.Value)
	c.JSON(200, HashResponse{Hash: hash})
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
		"None",
	}
	
	for _, validType := range validTypes {
		if contextType == validType {
			return true
		}
	}
	return false
}
