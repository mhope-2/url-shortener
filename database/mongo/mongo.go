package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// Client func
func Client() *mongo.Client {

	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set.")
	}

	opts := options.Client().ApplyURI(mongodbURI)

	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()

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

// GetCollection is a  function makes a connection with a collection in the database
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	mongodbName := os.Getenv("MONGODB_NAME")
	if mongodbName == "" {
		log.Fatal("MONGODB_NAME environment variable is not set.")
	}

	var collection *mongo.Collection = client.Database(mongodbName).Collection(collectionName)

	return collection
}
