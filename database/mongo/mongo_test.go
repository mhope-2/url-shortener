package mongo

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientConnection(t *testing.T) {
	// Test will fail if the environment variables are not set.
	if os.Getenv("MONGODB_URI") == "" || os.Getenv("MONGODB_NAME") == "" {
		t.Fatal("Environment variables MONGODB_URI and MONGODB_NAME must be set")
	}

	client := Client()

	// Assert client is not nil
	assert.NotNil(t, client, "Client should not be nil")

	err := client.Disconnect(context.TODO())
	assert.Nil(t, err)
}

func TestGetCollection(t *testing.T) {
	client := Client()

	collectionName := "testCollection"
	collection := GetCollection(client, collectionName)

	// Assert collection is not nil and has correct name
	assert.NotNil(t, collection, "Collection should not be nil")
	assert.Equal(t, collectionName, collection.Name(), "Collection name should match")

	err := client.Disconnect(context.TODO())
	assert.Nil(t, err)
}
