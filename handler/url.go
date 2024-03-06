package handler

import (
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

type CreateShortURLRequest struct {
	URL  string `json:"url" validate:"required"`
	Slug string `json:"slug,omitempty"` // optional
}

// CreateShortLink returns a shortened url
func (h *Handler) CreateShortLink(c *gin.Context) {
	var data CreateShortURLRequest
	clientIP := c.ClientIP()

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "failed to parse request"})
		return
	}

	// Prefix http:// to a given url if it is without one
	// c.Redirect(code int, location string) requires the http prefix
	if !strings.HasPrefix(data.URL, "http://") && !strings.HasPrefix(data.URL, "https://") {
		data.URL = "http://" + data.URL
	}

	if data.Slug != "" {
		if len(data.Slug) < AllowedSlugLength {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "slug length should be at least 6"})
			return
		}

		// Check if the provided slug is available
		existingURL, err := h.Repo.GetURL(data.Slug, clientIP)
		if err != nil {
			log.Println("Error shortening URL:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Request failed"})
			return
		}

		if existingURL != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Slug is not available, try again"})
			return
		}
	}

	if data.Slug == "" {
		existingURL, err := h.Repo.GetURL(data.URL, clientIP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Request failed"})
			return
		}

		if existingURL != nil {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"shortened_url": shared.GetShortenedURL(existingURL.Slug),
				},
			})
			return
		}

		data.Slug = h.Repo.GenerateSlug(data.URL, LowerBound, UpperBound)
	}

	c.ClientIP()

	url, err := h.Repo.CreateURL(data.URL, data.Slug, clientIP)

	if err != nil {
		log.Println("Error shortening url: ", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "Request failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"shortened_url": shared.GetShortenedURL(url.Slug),
		},
	})
}

// RedirectToOriginalURL redirects to the stored url for the given slug
func (h *Handler) RedirectToOriginalURL(c *gin.Context) {
	slug := c.Param("slug")
	clientIP := c.ClientIP()

	url, err := h.Repo.GetURL(slug, clientIP)

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

	c.Redirect(http.StatusTemporaryRedirect, url.URL)
	c.Abort()
}
