package routes

import (
	"contoso/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers API routes on the provided Gin router
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/ping", controllers.Ping)
		// Player CRUD routes
		api.GET("/players", controllers.GetPlayers)
		api.GET("/players/:id", controllers.GetPlayer)
		api.POST("/players", controllers.CreatePlayer)
		api.PUT("/players/:id", controllers.UpdatePlayer)
		api.DELETE("/players/:id", controllers.DeletePlayer)
	}
}
