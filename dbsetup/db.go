package dbsetup

import (
	"context"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	once       sync.Once
	mongoURI   string
	database   string = "contoso"
	collection string = "players"
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
