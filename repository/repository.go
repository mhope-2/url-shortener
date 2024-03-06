// Package repository houses the data layer
package repository

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DB    *mongo.Database
	Cache *redis.Client
}

func New(db *mongo.Database, cache *redis.Client) *Repository {
	return &Repository{db, cache}
}
