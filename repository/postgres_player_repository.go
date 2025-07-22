package repository

import (
	"contoso/models"
	"database/sql"
	"errors"
	"strconv"
)

type PostgresPlayerRepository struct {
	db *sql.DB
}

func NewPostgresPlayerRepository(db *sql.DB) *PostgresPlayerRepository {
	return &PostgresPlayerRepository{db: db}
}

func (r *PostgresPlayerRepository) CreatePlayer(player *models.Player) (*models.Player, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO players (name, surname, balance) VALUES ($1, $2, $3) RETURNING id",
		player.Name, player.Surname, player.Balance,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	player.ID = strconv.Itoa(id)
	return player, nil
}

func (r *PostgresPlayerRepository) GetPlayers() ([]models.Player, error) {
	rows, err := r.db.Query("SELECT id, name, surname, balance FROM players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []models.Player
	for rows.Next() {
		var p models.Player
		var id int
		if err := rows.Scan(&id, &p.Name, &p.Surname, &p.Balance); err == nil {
			p.ID = strconv.Itoa(id)
			players = append(players, p)
		}
	}
	return players, nil
}

func (r *PostgresPlayerRepository) GetPlayer(id string) (*models.Player, error) {
	var p models.Player
	var intID int
	err := r.db.QueryRow("SELECT id, name, surname, balance FROM players WHERE id = $1", id).
		Scan(&intID, &p.Name, &p.Surname, &p.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("player not found")
		}
		return nil, err
	}
	p.ID = strconv.Itoa(intID)
	return &p, nil
}

func (r *PostgresPlayerRepository) UpdatePlayer(id string, input *models.Player) (*models.Player, error) {
	_, err := r.db.Exec(
		"UPDATE players SET name = $1, surname = $2, balance = $3 WHERE id = $4",
		input.Name, input.Surname, input.Balance, id,
	)
	if err != nil {
		return nil, err
	}
	input.ID = id
	return input, nil
}

func (r *PostgresPlayerRepository) DeletePlayer(id string) error {
	_, err := r.db.Exec("DELETE FROM players WHERE id = $1", id)
	return err
}
