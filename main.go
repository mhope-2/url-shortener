package main

import (
	"fmt"
	"github.com/mhope-2/url_shortener/database/redis"
	"os"

	_ "github.com/joho/godotenv/autoload" // autoload .env file
	"github.com/mhope-2/url_shortener/handler"
	"github.com/mhope-2/url_shortener/server"
)

func main() {

	r := redis.New(&redis.Config{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	s := server.New()
	h := handler.New(r)

	routes := s.Group("")
	h.Register(routes)

	server.Start(&s, &server.Config{
		Port: fmt.Sprintf(":%s", os.Getenv("PORT")),
	})
}
