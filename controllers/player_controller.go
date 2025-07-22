package controllers

import (
	"contoso/models"
	"contoso/repository"
	"github.com/gofiber/fiber/v2"
)

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

func GetPlayers(repo repository.PlayerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		players, err := repo.GetPlayers()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(players)
	}
}

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
