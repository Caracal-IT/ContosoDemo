package dbsetup

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	once       sync.Once
	mongoURI   string
	database   string = "contoso"
	collection string = "players"

	pgOnce sync.Once
	pgDB   *sql.DB
)

func getMongoURI() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://root:example@localhost:27017"
	}
	return uri
}

func GetMongoCollection() *mongo.Collection {
	once.Do(func() {
		mongoURI = getMongoURI()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			panic("failed to connect to MongoDB: " + err.Error())
		}
	})
	return client.Database(database).Collection(collection)
}

func GetPostgresDB() *sql.DB {
	pgOnce.Do(func() {
		pgURL := os.Getenv("POSTGRES_URL")
		if pgURL == "" {
			pgURL = "postgres://postgres:example@localhost:5432/contoso?sslmode=disable"
		}
		var err error
		pgDB, err = sql.Open("postgres", pgURL)
		if err != nil {
			panic("failed to connect to Postgres: " + err.Error())
		}
		// Migration: create the player table if it doesn't exist
		_, err = pgDB.Exec(`
			CREATE TABLE IF NOT EXISTS players (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				surname TEXT NOT NULL,
				balance DOUBLE PRECISION NOT NULL
			)
		`)
		if err != nil {
			panic("failed to create players table: " + err.Error())
		}
	})
	return pgDB
}
