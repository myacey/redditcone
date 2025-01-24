package mongorepo

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConfigureMongoClient() (*mongo.Client, error) {
	mongoClientOpts := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	mongoClient, err := mongo.Connect(context.TODO(), mongoClientOpts)
	if err != nil {
		return nil, fmt.Errorf("cant connect to mongo: %v", err)
	}

	if err = mongoClient.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("cant connect to mongo: %v", err)
	}

	return mongoClient, nil
}
