package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client func
func Client() *mongo.Client {

	mongodbURI := fmt.Sprintf("%s", os.Getenv("MONGODB_URI"))

	if mongodbURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set.")
	}

	opts := options.Client().ApplyURI(mongodbURI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")

	return client
}

// OpenCollection is a  function makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database(os.Getenv("MONGODB_NAME")).Collection(collectionName)

	return collection
}
