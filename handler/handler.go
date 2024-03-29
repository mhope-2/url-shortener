// Package handler defines the API and DB interface
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/mhope-2/url_shortener/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	DB    *mongo.Database
	Repo  *repository.Repository
	Cache *redis.Client
}

func New(db *mongo.Database, cache *redis.Client) *Handler {
	repo := repository.New(db, cache)

	return &Handler{
		DB:    db,
		Repo:  repo,
		Cache: cache,
	}
}

func (h *Handler) Register(v1 *gin.RouterGroup) {
	v1.POST("/short-link", h.CreateShortLink)
	v1.GET("/:slug", h.RedirectToOriginalURL)

}
