package repository_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mhope-2/url_shortener/repository/test_repository"
	"github.com/stretchr/testify/assert"
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

func TestUrlRepo(t *testing.T) {
	t.Run("crud", crud)
	t.Run("slug", slug)
}

// =============================================================================

func crud(t *testing.T) {
	assert.Equal(t, 1, 1)

	repo := test_repository.NewRepo()

	// test repo create url with slug
	url, err := repo.CreateUrl("https://youtube.com/", "A1B2C3D4", "127.0.0.1")
	if err != nil {
		t.Fatalf("Error creating url: %v", err)
	}

	assert.NoError(t, err, "Failed to create url")
	assert.Equal(t, url.Url, "https://youtube.com/")
	assert.Equal(t, url.Slug, "A1B2C3D4")

	// ------------------------------------------------------------------------

	// test repo url retrieve
	url, err = repo.GetUrl("A1B2C3D4", "127.0.0.1")
	if err != nil {
		t.Fatalf("Error retriving url: %v", err)
	}

	assert.NoError(t, err, "Failed to retrieve url")
	assert.Equal(t, url.Url, "https://youtube.com/")
	assert.Equal(t, url.Slug, "A1B2C3D4")
}

// =============================================================================

func slug(t *testing.T) {

	repo := test_repository.NewRepo()

	// test repo generate slug
	slug1 := repo.GenerateSlug("https://go.dev/", 25, 10_000)
	assert.NotNil(t, slug1)
}
