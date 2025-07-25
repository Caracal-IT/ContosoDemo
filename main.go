package main

import (
	"contoso/dbsetup"
	_ "contoso/docs" // swaggo docs
	"contoso/elasticlog"
	"contoso/repository"
	"contoso/routes"
	"github.com/gofiber/fiber/v2"
	logger2 "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/swaggo/fiber-swagger"
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

	// Serve Swagger UI at /swagger/*
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Serve raw swagger.json for ReDoc and other tools
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join("docs", "swagger.json"), true)
	})

	// Serve ReDoc at /redoc
	app.Get("/redoc", func(c *fiber.Ctx) error {
		html := `<!DOCTYPE html>
<html>
  <head>
    <title>ReDoc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" href="data:,">
    <style>
      @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap');
      html, body, redoc, #redoc-container {
        height: 100%;
        width: 100%;
        margin: 0;
        padding: 0;
        background: #181a1b !important;
        color: #e0e0e0 !important;
        font-family: 'Inter', 'Segoe UI', Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
      }
      /* Headings */
      h1, h2, h3, h4, h5, h6 {
        font-family: 'Inter', 'Segoe UI', Arial, sans-serif;
        font-weight: 700;
        color: #6ca0ff !important;
      }
      /* Make SVG .sc class font and fill white */
      svg .sc {
        color: #fff !important;
        fill: #fff !important;
      }
      /* Tag cards */
      [id^="tag/"] {
        background: #23272a !important;
        border-radius: 12px;
        box-shadow: 0 2px 8px 0 rgba(0,0,0,0.25);
        margin-bottom: 32px !important;
        padding: 24px 32px !important;
        transition: box-shadow 0.2s;
        color: #e0e0e0 !important;
      }
      [id^="tag/"]:hover {
        box-shadow: 0 4px 16px 0 rgba(0,0,0,0.35);
      }
      /* Links and buttons */
      a, button {
        color: #6ca0ff !important;
        border-radius: 6px;
        font-weight: 600;
        text-decoration: none;
        transition: background 0.2s, color 0.2s;
      }
      a:hover, button:hover {
        background: #26324a !important;
        color: #a3c9ff !important;
      }
      /* Code blocks */
      code, pre {
        background: #23272a !important;
        color: #a3c9ff !important;
        border-radius: 6px;
        font-size: 0.97em;
        padding: 2px 8px;
      }
      /* Section backgrounds */
      .sc-eDvSVe, .sc-jrsJWt, .sc-hKwDye, .sc-cPiKLX, .menu-content {
        background: #202225 !important;
        color: #e0e0e0 !important;
      }
      /* Force dark theme for all menu/sidebar and ReDoc UI elements */
      .menu-content, .sc-dkzDqf, .sc-hKwDye, .sc-cPiKLX, .sc-eDvSVe, .sc-jrsJWt, .sc-gEvEer, .sc-ksZaOG, .sc-hBUSln, .sc-bZQynM, .sc-lllmON, .sc-cmTdod, .sc-dcJsrY, .sc-hKwDye, .sc-cPiKLX, .sc-jrsJWt, .sc-fubCfw, .sc-kgflAQ, .sc-lllmON, .sc-cmTdod, .sc-dcJsrY {
        background: #181a1b !important;
        color: #e0e0e0 !important;
        border-color: #23272a !important;
      }
      .menu-content *, .sc-dkzDqf *, .sc-hKwDye *, .sc-cPiKLX *, .sc-eDvSVe *, .sc-jrsJWt *, .sc-gEvEer *, .sc-ksZaOG *, .sc-hBUSln *, .sc-bZQynM *, .sc-lllmON *, .sc-cmTdod *, .sc-dcJsrY *, .sc-fubCfw *, .sc-kgflAQ * {
        color: #e0e0e0 !important;
        background: transparent !important;
      }
      .menu-content a, .menu-content a *, .sc-dkzDqf a, .sc-dkzDqf a *, .sc-hKwDye a, .sc-hKwDye a *, .sc-cPiKLX a, .sc-cPiKLX a * {
        color: #6ca0ff !important;
      }
      .menu-content a:hover, .sc-dkzDqf a:hover, .sc-hKwDye a:hover, .sc-cPiKLX a:hover {
        color: #a3c9ff !important;
        background: #23272a !important;
      }
      /* Make menu SVG arrows lighter for dark theme */
      .menu-content svg, .sc-cBoqAE svg, .sc-dkzDqf svg, .sc-hKwDye svg, .sc-cPiKLX svg, .sc-jrsJWt svg, .sc-eDvSVe svg {
        fill: #b3cfff !important;
        color: #b3cfff !important;
        opacity: 1 !important;
      }

	svg polygon {
  		fill: white !important;
	}
      /* Remove box-shadow from menu for a flat look */
      .menu-content, .sc-dkzDqf {
        box-shadow: none !important;
      }
      /* Fix search bar and input fields */
      input, .sc-hKwDye input, .sc-cPiKLX input {
        background: #23272a !important;
        color: #e0e0e0 !important;
        border: 1px solid #23272a !important;
      }
      input::placeholder {
        color: #888 !important;
      }
      /* Fix scrollbar in menu */
      .menu-content ::-webkit-scrollbar {
        width: 8px;
        background: #23272a;
      }
      .menu-content ::-webkit-scrollbar-thumb {
        background: #181a1b;
        border-radius: 4px;
      }
      /* Application/JSON dark theme fix, but keep response body (second part) light */
      .sc-dkzDqf, .sc-eDvSVe, .sc-hKwDye, .sc-cPiKLX {
        color: #e0e0e0 !important;
        background: #23272a !important;
      }
      .sc-dkzDqf code, .sc-eDvSVe code, .sc-hKwDye code, .sc-cPiKLX code {
        color: #a3c9ff !important;
        background: #23272a !important;
      }
      .sc-dkzDqf pre, .sc-eDvSVe pre, .sc-hKwDye pre, .sc-cPiKLX pre {
        color: #a3c9ff !important;
        background: #23272a !important;
      }
      /* Keep the second part (response body) light */
      .sc-ikZpkk, .sc-ikZpkk * {
        background: #fff !important;
        color: #23272a !important;
      }
      /* Make menu operation verbs colored and menu text lighter */
      .sc-cBoqAE, .sc-cBoqAE * {
        color: #f0f0f0 !important;
      }
      /* HTTP verb colors in menu */
      .sc-cBoqAE span[title="get"], .sc-cBoqAE .http-verb-get {
        color: #61affe !important;
        background: none !important;
      }
      .sc-cBoqAE span[title="post"], .sc-cBoqAE .http-verb-post {
        color: #49cc90 !important;
        background: none !important;
      }
      .sc-cBoqAE span[title="put"], .sc-cBoqAE .http-verb-put {
        color: #fca130 !important;
        background: none !important;
      }
      .sc-cBoqAE span[title="delete"], .sc-cBoqAE .http-verb-delete {
        color: #f93e3e !important;
        background: none !important;
      }
      /* For PATCH and other verbs */
      .sc-cBoqAE span[title="patch"], .sc-cBoqAE .http-verb-patch {
        color: #bada55 !important;
        background: none !important;
      }
      .sc-cBoqAE span[title="options"], .sc-cBoqAE .http-verb-options {
        color: #ebebeb !important;
        background: none !important;
      }
      /* HTTP verb colors for operation-type classes in menu and content */
      .operation-type.get, .sc-ikXwFM.get {
        color: #61affe !important;
      }
      .operation-type.post, .sc-ikXwFM.post {
        color: #49cc90 !important;
      }
      .operation-type.put, .sc-ikXwFM.put {
        color: #fca130 !important;
      }
      .operation-type.delete, .sc-ikXwFM.delete {
        color: #f93e3e !important;
      }
      .operation-type.patch, .sc-ikXwFM.patch {
        color: #bada55 !important;
      }
      .operation-type.options, .sc-ikXwFM.options {
        color: #ebebeb !important;
      }
      /* Highlight selected menu item */
      .sc-cBoqAE .sc-hSdWYo.selected, .sc-cBoqAE .sc-hSdWYo.selected *,
      .sc-cBoqAE .sc-hSdWYo.active, .sc-cBoqAE .sc-hSdWYo.active * {
        background: #23272a !important;
        color: #6ca0ff !important;
        border-radius: 6px;
        font-weight: 700;
        box-shadow: 0 0 0 2px #6ca0ff33;
        transition: background 0.2s, color 0.2s;
      }
      /* Also highlight operation-type in menu if selected */
      .sc-cBoqAE .sc-hSdWYo.selected .operation-type,
      .sc-cBoqAE .sc-hSdWYo.active .operation-type {
        color: #fff !important;
        background: #6ca0ff !important;
        border-radius: 4px;
        padding: 2px 8px;
      }
      @media (max-width: 900px) {
        [id^="tag/"] {
          padding: 12px 8px !important;
        }
        h1 { font-size: 2rem; }
        h2 { font-size: 1.5rem; }
        h3 { font-size: 1.2rem; }
      }
    </style>
  </head>
  <body>
    <redoc spec-url='/swagger/doc.json'></redoc>
    <script src='https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js'></script>
  </body>
</html>`
		return c.Type("html").SendString(html)
	})

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
