package repository

import (
	"context"
	"contoso/models"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPlayerRepository struct {
	collection *mongo.Collection
}

func NewMongoPlayerRepository(col *mongo.Collection) *MongoPlayerRepository {
	return &MongoPlayerRepository{collection: col}
}

func (r *MongoPlayerRepository) CreatePlayer(player *models.Player) (*models.Player, error) {
	player.ID = ""
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := r.collection.InsertOne(ctx, player)
	if err != nil {
		return nil, err
	}
	player.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return player, nil
}

func (r *MongoPlayerRepository) GetPlayers() ([]models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var players []models.Player
	for cursor.Next(ctx) {
		var player models.Player
		if err := cursor.Decode(&player); err == nil {
			if oid, ok := cursor.Current.Lookup("_id").ObjectIDOK(); ok {
				player.ID = oid.Hex()
			}
			players = append(players, player)
		}
	}
	return players, nil
}

func (r *MongoPlayerRepository) GetPlayer(id string) (*models.Player, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var player models.Player
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&player); err != nil {
		return nil, errors.New("player not found")
	}
	player.ID = objID.Hex()
	return &player, nil
}

func (r *MongoPlayerRepository) UpdatePlayer(id string, input *models.Player) (*models.Player, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{
		"$set": bson.M{
			"name":    input.Name,
			"surname": input.Surname,
			"balance": input.Balance,
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}
	input.ID = id
	return input, nil
}

func (r *MongoPlayerRepository) DeletePlayer(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
