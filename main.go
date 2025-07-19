package main

import (
	"contoso/dbsetup"
	"contoso/elasticlog"
	"contoso/routes"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type elasticErrorWriter struct {
	logger *elasticlog.Logger
}

type elasticInfoWriter struct {
	logger *elasticlog.Logger
}

func (w *elasticErrorWriter) Write(p []byte) (n int, err error) {
	// Write to stderr as usual
	n, err = os.Stderr.Write(p)
	// Also log to elastic as error
	w.logger.Error(string(p), nil)
	return n, err
}

func (w *elasticInfoWriter) Write(p []byte) (n int, err error) {
	// Write to stderr as usual
	n, err = os.Stderr.Write(p)
	// Also log to elastic as error
	w.logger.Info(string(p), nil)
	return n, err
}

func main() {
	// Ensure DB is initialized on start
	dbsetup.GetMongoCollection()

	// Setup logger
	logLevel := elasticlog.ParseLogLevel(os.Getenv("LOG_LEVEL"))
	logger := elasticlog.NewLogger(
		logLevel,
		"contoso-", // index
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)

	logger.Info("Contoso backend started", map[string]interface{}{
		"event": "startup",
	})

	// Gin logger middleware
	gin.DefaultWriter = &elasticInfoWriter{logger: logger}
	gin.DefaultErrorWriter = &elasticErrorWriter{logger: logger}
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		entry := map[string]interface{}{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  status,
			"latency": latency.String(),
			"client":  c.ClientIP(),
		}
		switch {
		case status >= 500:
			logger.Error("HTTP request", entry)
		case status >= 400:
			logger.Warn("HTTP request", entry)
		default:
			logger.Info("HTTP request", entry)
		}
	})

	// API routes
	routes.RegisterRoutes(r)

	// Serve static files for frontend
	publicDir := "./public"
	r.Static("/assets", filepath.Join(publicDir, "assets"))

	// Serve index.html for non-API routes (SPA fallback)
	r.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.File(filepath.Join(publicDir, "index.html"))
	})

	err := r.Run(":8080")
	if err != nil {
		logger.Error("Failed to start server", map[string]interface{}{"error": err.Error()})
	}
}
