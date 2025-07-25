package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Ping godoc
// @Summary Health check
// @Description Returns pong if the server is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/ping [get]
func Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "pong"})
}
