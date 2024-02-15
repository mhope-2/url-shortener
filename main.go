package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload" // autoload .env file
	"github.com/mhope-2/url_shortener/database/mongo"
	"github.com/mhope-2/url_shortener/database/redis"
	"github.com/mhope-2/url_shortener/handler"
	"github.com/mhope-2/url_shortener/server"
)

func main() {

	db, err := mongo.New(&mongo.Config{
		MongodbURI:  os.Getenv("MONGODB_URI"),
		MongodbName: os.Getenv("MONGODB_NAME"),
	})
	
	if err != nil {
		log.Fatal("Failed to connect to Mongo database", err)
	}

	r := redis.New(&redis.Config{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	s := server.New()
	h := handler.New(db, r)

	routes := s.Group("")
	h.Register(routes)

	server.Start(&s, &server.Config{
		Port: fmt.Sprintf(":%s", os.Getenv("PORT")),
	})
}
