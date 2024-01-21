package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/mhope-2/url_shortener/repository"
)

type Handler struct {
	Repo  *repository.Repository
	Cache *redis.Client
}

func New(cache *redis.Client) *Handler {
	repo := repository.New(cache)

	return &Handler{
		Repo:  repo,
		Cache: cache,
	}
}

func (h *Handler) Register(v1 *gin.RouterGroup) {
	v1.POST("/short-link", h.CreateShortLink)
	v1.GET("/:slug", h.RedirectToOriginalUrl)

}
