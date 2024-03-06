// Package mongo houses code for connecting to mongodb
package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Config struct {
	MongodbURI  string
	MongodbName string
}

func New(config *Config) (*mongo.Database, error) {

	var mongodbURI string

	if os.Getenv("ENV") == "testing" {
		mongodbURI = os.Getenv("TEST_MONGODB_URI")
	} else {
		mongodbURI = os.Getenv("MONGODB_URI")
	}

	//mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set.")
	}

	opts := options.Client().ApplyURI(config.MongodbURI)

	// set context with timeout for mongodb connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	var result bson.M

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}

	database := client.Database(config.MongodbName)

	log.Println("Connected to MongoDB!")

	return database, nil
}
