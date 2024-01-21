package repository

import (
	"github.com/go-redis/redis"
)

type Repository struct {
	Cache *redis.Client
}

func New(cache *redis.Client) *Repository {
	return &Repository{cache}
}
