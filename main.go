package main

import (
	"contoso/dbsetup"
	"contoso/routes"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func main() {
	// Ensure DB is initialized on start
	dbsetup.GetMongoCollection()

	router := gin.Default()

	// API routes
	routes.RegisterRoutes(router) // Uncommented to register API routes

	// Serve static files for frontend
	publicDir := "./public"
	router.Static("/assets", filepath.Join(publicDir, "assets"))

	// Serve index.html for non-API routes (SPA fallback)
	router.NoRoute(func(c *gin.Context) {
		// Only serve index.html if the path does NOT start with /api
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.File(filepath.Join(publicDir, "index.html"))
	})

	router.Run(":8080")
}
