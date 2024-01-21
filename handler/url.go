package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mhope-2/url_shortener/shared"
	"log"
	"net/http"
	"strings"
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
	clientIP := c.ClientIP()

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "failed to parse request"})
		return
	}

	// Prefix http:// to a given url if it is without one
	// c.Redirect(code int, location string) requires the http prefix
	if !strings.HasPrefix(data.Url, "http://") && !strings.HasPrefix(data.Url, "https://") {
		data.Url = "http://" + data.Url
	}

	if data.Slug != "" {
		if len(data.Slug) < AllowedSlugLength {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "slug length should be at least 6"})
			return
		}

		// Check if the provided slug is available
		existingUrl, err := h.Repo.GetUrl(data.Slug, clientIP)
		if err != nil {
			log.Println("Error shortening URL:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Request failed"})
			return
		}

		if existingUrl != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Slug is not available, try again"})
			return
		}
	}

	if data.Slug == "" {
		existingUrl, err := h.Repo.GetUrl(data.Url, clientIP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Request failed"})
			return
		}

		if existingUrl != nil {
			fmt.Println("EX: ", existingUrl)
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"shortened_url": shared.GetShortenedUrl(existingUrl.Slug),
				},
			})
			return
		}

		data.Slug = h.Repo.GenerateSlug(data.Url, LowerBound, UpperBound)
	}

	c.ClientIP()

	url, err := h.Repo.CreateUrl(data.Url, data.Slug, clientIP)

	if err != nil {
		log.Println("Error shortening url: ", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "Request failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"shortened_url": shared.GetShortenedUrl(url.Slug),
		},
	})
}

// RedirectToOriginalUrl redirects to the stored url for the given slug
func (h *Handler) RedirectToOriginalUrl(c *gin.Context) {
	slug := c.Param("slug")
	clientIP := c.ClientIP()

	url, err := h.Repo.GetUrl(slug, clientIP)

	if err != nil {
		log.Println("Error retrieving url: ", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "Request failed",
		})
		return
	}

	if url == nil {
		log.Println("Invalid slug: ", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "Request failed",
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url.Url)
	c.Abort()
}
