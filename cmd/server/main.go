package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode for production-like behavior
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("cmd/server/template.html")

	// Routes
	r.GET("/", handleHome)
	r.POST("/lex", handleLex)
	r.POST("/preparse", handlePreparse)
	r.POST("/parse", handleParse)
	r.POST("/build", handleBuild)
	r.POST("/hash", handleHash)

	port := ":8000"
	log.Printf("Starting development server on http://localhost%s", port)
	log.Println("This server is for development/testing only - do not expose to clients!")

	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
