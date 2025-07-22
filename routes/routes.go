package routes

import (
	"contoso/controllers"
	"contoso/repository"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutesFiber registers API routes on the provided Fiber app
func RegisterRoutesFiber(app *fiber.App, playerRepo repository.PlayerRepository) {
	api := app.Group("/api")
	api.Get("/ping", controllers.Ping)
	// Player CRUD routes
	api.Get("/players", controllers.GetPlayers(playerRepo))
	api.Get("/players/:id", controllers.GetPlayer(playerRepo))
	api.Post("/players", controllers.CreatePlayer(playerRepo))
	api.Put("/players/:id", controllers.UpdatePlayer(playerRepo))
	api.Delete("/players/:id", controllers.DeletePlayer(playerRepo))
}
