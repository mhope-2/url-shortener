package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var (
	AllowedSlugLength = 6
	LowerBound        = 100
	UpperBound        = 1_000_000
)

type CreateShortUrlRequest struct {
	Url  string `json:"url" validate:"required"`
	Slug string `json:"slug,omitempty"` // optional
}

// CreateShortLink returns a shortened url
func (h *Handler) CreateShortLink(c *gin.Context) {
	var data CreateShortUrlRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "failed to parse request",
		})
	}

	if data.Slug == "" {
		data.Slug = h.Repo.GenerateSlug(data.Url, LowerBound, UpperBound)
	}

	// slug length should be at least 6
	if data.Slug != "" && len(data.Slug) < AllowedSlugLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "slug length should be at least 6",
		})
		return
	} else {
		// If user has a preferred slug, check if it's available
		existingUrl, err := h.Repo.GetUrlBySlug(data.Slug)

		if err != nil {
			log.Println("Error shortening url: ", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "Request failed",
			})
			return
		}

		if existingUrl != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "Slug is not available, try again",
			})
			return
		}
	}

	// check cache layer first

	url, err := h.Repo.CreateUrl(data.Url, data.Slug)

	if err != nil {
		log.Println("Error shortening url: ", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "Request failed",
		})
		return
	}

	shortened_url := fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), url.Slug)

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"shortened_url": shortened_url,
		},
	})
}

// RedirectToOriginalUrl redirects to the stored url for the given slug
func (h *Handler) RedirectToOriginalUrl(c *gin.Context) {
	slug := c.Param("slug")

	// get from cache first

	url, err := h.Repo.GetUrlBySlug(slug)

	if err != nil {
		log.Println("Error retrieving url: ", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "Request failed",
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url.Url)
}
