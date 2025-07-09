package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var homeTemplate *template.Template

func main() {
	// Load HTML template from file
	tmplPath := filepath.Join("cmd", "server", "template.html")
	var err error
	homeTemplate, err = template.ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/lex", handleLex)
	http.HandleFunc("/preparse", handlePreparse)
	http.HandleFunc("/parse", handleParse)
	http.HandleFunc("/build", handleBuild)

	port := ":8000"
	log.Printf("Starting development server on http://localhost%s\n", port)
	log.Println("This server is for development/testing only - do not expose to clients!")
	log.Fatal(http.ListenAndServe(port, nil))
} 