package controllers

import (
	"contoso/models"
	"contoso/repository"
	"github.com/gofiber/fiber/v2"
)

// CreatePlayer godoc
// @Summary Create a new player
// @Description Create a new player in the system
// @Tags players
// @Accept json
// @Produce json
// @Param player body models.Player true "Player data"
// @Success 201 {object} models.Player
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /players [post]
func CreatePlayer(repo repository.PlayerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var player models.Player
		if err := c.BodyParser(&player); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		created, err := repo.CreatePlayer(&player)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(created)
	}
}

// GetPlayers godoc
// @Summary Get all players
// @Description Get a list of all players
// @Tags players
// @Produce json
// @Success 200 {array} models.Player
// @Failure 500 {object} map[string]string
// @Router /players [get]
func GetPlayers(repo repository.PlayerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		players, err := repo.GetPlayers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(players)
	}
}

// GetPlayer godoc
// @Summary Get a player by ID
// @Description Get details of a player by ID
// @Tags players
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} models.Player
// @Failure 404 {object} map[string]string
// @Router /players/{id} [get]
func GetPlayer(repo repository.PlayerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		player, err := repo.GetPlayer(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(player)
	}
}

// UpdatePlayer godoc
// @Summary Update a player
// @Description Update a player's information
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Param player body models.Player true "Player data"
// @Success 200 {object} models.Player
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /players/{id} [put]
func UpdatePlayer(repo repository.PlayerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var input models.Player
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		updated, err := repo.UpdatePlayer(id, &input)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(updated)
	}
}

// DeletePlayer godoc
// @Summary Delete a player
// @Description Delete a player by ID
// @Tags players
// @Param id path string true "Player ID"
// @Success 204 {string} string "No Content"
// @Failure 500 {object} map[string]string
// @Router /players/{id} [delete]
func DeletePlayer(repo repository.PlayerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		err := repo.DeletePlayer(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
