package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/studyguides-com/study-guides-parser/core/config"
	"github.com/studyguides-com/study-guides-parser/core/processor"
)

type ParseRequest struct {
	Content string `json:"content"`
}

type ParseResponse struct {
	Success bool                        `json:"success"`
	AST     interface{}                 `json:"ast,omitempty"`
	Errors  []processor.ProcessingError `json:"errors,omitempty"`
	Tree    interface{}                 `json:"tree,omitempty"`
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := homeTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleLex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	lines := strings.Split(req.Content, "\n")
	result, err := processor.Lex(lines)
	if err != nil {
		http.Error(w, fmt.Sprintf("Lexing error: %v", err), http.StatusInternalServerError)
		return
	}
	response := ParseResponse{
		Success: result.Success,
		Errors:  result.Errors,
	}
	if result.Success {
		response.AST = result.Tokens
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handlePreparse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	lines := strings.Split(req.Content, "\n")
	result, err := processor.Preparse(lines)
	if err != nil {
		http.Error(w, fmt.Sprintf("Preparsing error: %v", err), http.StatusInternalServerError)
		return
	}
	response := ParseResponse{
		Success: result.Success,
		Errors:  result.Errors,
	}
	if result.Success {
		response.AST = result.Tokens
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	lines := strings.Split(req.Content, "\n")
	metadata := config.NewMetaData("colleges")
	result, err := processor.Parse(lines, metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parsing error: %v", err), http.StatusInternalServerError)
		return
	}
	response := ParseResponse{
		Success: result.Success,
		Errors:  result.Errors,
	}
	if result.Success && result.AST != nil {
		response.AST = result.AST
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleBuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	lines := strings.Split(req.Content, "\n")
	metadata := config.NewMetaData("colleges")
	result, err := processor.Build(lines, metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Build error: %v", err), http.StatusInternalServerError)
		return
	}
	response := ParseResponse{
		Success: result.Success,
		Errors:  result.Errors,
	}
	if result.Success && result.Tree != nil {
		response.Tree = result.Tree
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 