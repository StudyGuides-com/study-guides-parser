package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/processor"
)

type ParseRequest struct {
	Content string `json:"content" binding:"required"`
}

type ParseResponse struct {
	Success bool                        `json:"success"`
	AST     interface{}                 `json:"ast,omitempty"`
	Errors  []processor.ProcessingError `json:"errors,omitempty"`
	Tree    interface{}                 `json:"tree,omitempty"`
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
	metadata := config.NewMetaData("colleges")
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
	metadata := config.NewMetaData("colleges")
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