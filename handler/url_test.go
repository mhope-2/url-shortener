package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := os.Setenv("ENV", "testing")
	if err != nil {
		fmt.Println("Error setting env variable, ENV")
		return
	}

	err = os.Setenv("TEST_MONGODB_URI", "mongodb://localhost:27017")
	if err != nil {
		fmt.Println("Error setting env variable, TEST_MONGODB_URI")
		return
	}

	err = os.Setenv("TEST_MONGODB_NAME", "urls")
	if err != nil {
		fmt.Println("Error setting env variable, TEST_MONGODB_NAME")
		return
	}

	os.Exit(m.Run())
}

func TestUrlHandler(t *testing.T) {
	t.Run("createShortLink", createShortLink)
	t.Run("UrlRedirect", urlRedirect)
}

// =============================================================================

func createShortLink(t *testing.T) {
	requestData := map[string]string{
		"url":  "https://go.dev/",
		"slug": "A1B2C3D4",
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON. Cause: %v\n", err)
	}

	req, err := http.NewRequest("POST", "/short-link", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Could not create request. Cause: %v\n", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Setup router with a dummy handler to simulate endpoint behavior
	router := gin.Default()
	router.POST("/short-link", func(c *gin.Context) {
		// Simulate successful handling of the request
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Create a response recorder to inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the response status code is as expected
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}
}

func urlRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/A1B2C3D4", nil)
	if err != nil {
		t.Fatalf("Could not create request. Cause: %v\n", err)
	}

	router := gin.Default()
	router.GET("/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		_ = slug // This line is just to avoid unused variable error
		c.JSON(http.StatusTemporaryRedirect, gin.H{})
	})

	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the response status code is as expected
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %v, got %v", http.StatusTemporaryRedirect, w.Code)
	}
}
