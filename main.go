package main

import (
	"contoso/dbsetup"
	"contoso/elasticlog"
	"contoso/routes"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func main() {
	// Ensure DB is initialized on start
	dbsetup.GetMongoCollection()

	// Semantic log to Elasticsearch on startup
	elasticlog.LogToElastic("info", "Contoso backend started", map[string]interface{}{
		"event": "startup",
	}, "contoso", "", "", // index, username, password: use env/default
	)

	router := gin.Default()

	// API routes
	routes.RegisterRoutes(router)

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

	err := router.Run(":8080")

	if err != nil {
		return
	}
}
