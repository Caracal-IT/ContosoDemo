package main

import (
	"contoso/dbsetup"
	"contoso/elasticlog"
	"contoso/repository"
	"contoso/routes"
	"github.com/gofiber/fiber/v2"
	logger2 "github.com/gofiber/fiber/v2/middleware/logger"
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
	n, err = os.Stderr.Write(p)
	w.logger.Error(string(p), nil)
	return n, err
}

func (w *elasticInfoWriter) Write(p []byte) (n int, err error) {
	n, err = os.Stderr.Write(p)
	w.logger.Info(string(p), nil)
	return n, err
}

func startBackgroundService(logger *elasticlog.Logger) {
	go func() {
		for {
			logger.Info("Background service heartbeat", map[string]interface{}{
				"event": "heartbeat",
				"time":  time.Now().Format(time.RFC3339),
			})
			time.Sleep(1 * time.Minute)
		}
	}()
}

func main() {
	// Setup logger
	logLevel := elasticlog.ParseLogLevel(os.Getenv("LOG_LEVEL"))
	logger := elasticlog.NewLogger(
		logLevel,
		"contoso-", // index
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)

	// Start background service
	startBackgroundService(logger)

	logger.Info("Contoso backend started", map[string]interface{}{
		"event": "startup",
	})

	// Choose repository based on environment variable
	var playerRepo repository.PlayerRepository
	dbType := os.Getenv("DB_TYPE")
	if dbType == "postgres" {
		playerRepo = repository.NewPostgresPlayerRepository(dbsetup.GetPostgresDB())
		logger.Info("Using Postgres repository", nil)
	} else {
		playerRepo = repository.NewMongoPlayerRepository(dbsetup.GetMongoCollection())
		logger.Info("Using MongoDB repository", nil)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error("Fiber error", map[string]interface{}{
				"error":  err.Error(),
				"path":   c.Path(),
				"method": c.Method(),
			})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Add Fiber's logger middleware for endpoint and info logging, logging to both console and elastic
	app.Use(logger2.New(logger2.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: time.RFC3339,
		Output:     &elasticInfoWriter{logger: logger}, // log to both console and elastic
	}))

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		status := c.Response().StatusCode()
		entry := map[string]interface{}{
			"method":  c.Method(),
			"path":    c.Path(),
			"status":  status,
			"latency": latency.String(),
			"client":  c.IP(),
		}
		switch {
		case status >= 500:
			logger.Error("HTTP request", entry)
		case status >= 400:
			logger.Warn("HTTP request", entry)
		default:
			logger.Info("HTTP request", entry)
		}
		return err
	})

	// Pass the repository to the routes/controllers
	routes.RegisterRoutesFiber(app, playerRepo)

	// Serve static files for frontend
	publicDir := "./public"
	app.Static("/assets", filepath.Join(publicDir, "assets"))

	// Serve index.html for non-API routes (SPA fallback)
	app.Use(func(c *fiber.Ctx) error {
		if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
		}
		return c.SendFile(filepath.Join(publicDir, "index.html"))
	})

	err := app.Listen(":8080")
	if err != nil {
		logger.Error("Failed to start server", map[string]interface{}{"error": err.Error()})
	}
}
