package mongoDb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"health/config"
)

func NewMongoClient() (*mongo.Client, *mongo.Database, error) {
	cfg := config.Load()
	clientOptions := options.Client().ApplyURI(cfg.MongoURI).SetAuth(options.Credential{Username: "root", Password: "example"})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	db := client.Database(cfg.MongoDBName)
	return client, db, nil
}
