package repository

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DB    *mongo.Client
	Cache *redis.Client
}

func New(db *mongo.Client, cache *redis.Client) *Repository {
	return &Repository{db, cache}
}
