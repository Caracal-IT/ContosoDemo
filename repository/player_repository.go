package repository

import (
	"contoso/models"
)

// PlayerRepository abstracts player CRUD operations.
type PlayerRepository interface {
	CreatePlayer(player *models.Player) (*models.Player, error)
	GetPlayers() ([]models.Player, error)
	GetPlayer(id string) (*models.Player, error)
	UpdatePlayer(id string, player *models.Player) (*models.Player, error)
	DeletePlayer(id string) error
}
