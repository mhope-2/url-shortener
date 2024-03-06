// Package testrepository defines a repository to used for testing
package testrepository

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/mhope-2/url_shortener/database/mongo"
	"github.com/mhope-2/url_shortener/database/redis"
	"github.com/mhope-2/url_shortener/repository"
	"log"
	"os"
)

func NewRepo() *repository.Repository {
	db, err := mongo.New(&mongo.Config{
		MongodbURI:  os.Getenv("TEST_MONGODB_URI"),
		MongodbName: os.Getenv("TEST_MONGODB_NAME"),
	})

	if err != nil {
		log.Fatal("Failed to connect to Mongo database", err)
	}

	mr, _ := miniredis.Run()
	r := redis.New(&redis.Config{Addr: mr.Addr(), Password: "", DB: 0})

	return &repository.Repository{
		DB:    db,
		Cache: r,
	}

}
