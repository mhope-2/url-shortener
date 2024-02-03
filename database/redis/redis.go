package redis

import (
	"github.com/go-redis/redis"
	"log"
	"net/url"
	"os"
)

type Config struct {
	Addr     string
	Password string
	DB       int
	DBurl    string
}

func New(config *Config) *redis.Client {
	if os.Getenv("ENV") == "staging" || os.Getenv("ENV") == "production" {
		parsedURL, _ := url.Parse(config.DBurl)
		password, _ := parsedURL.User.Password()

		return redis.NewClient(&redis.Options{
			Addr:     parsedURL.Host,
			Password: password,
		})
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Error pinging Redis server: %v", err)
	}

	log.Println("Connected to Redis server!")

	return client
}
