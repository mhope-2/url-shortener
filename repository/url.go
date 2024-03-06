package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mhope-2/url_shortener/shared"
	"log"
	"math/rand"
	"time"

	"github.com/mhope-2/url_shortener/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLRepository interface {
	CreateUrl(originalURL string, slug string, clientIP string) (*models.URL, error)
	GetUrl(slug string, clientIP string) (*models.URL, error)
	CacheUrl(url *shared.URL, clientIP string) error
	GetUrlFromCache(cacheKey string, clientIP string) (*shared.URL, error)
	GenerateRandomNumber(min, max int) int
	GenerateSlug(url string, min, max int) string
}

var collection = "url"

// CreateURL creates a url object, stores it in the db and the caches it
func (r *Repository) CreateURL(originalURL string, slug string, clientIP string) (*models.URL, error) {

	urlCollection := r.DB.Collection(collection)

	var url models.URL

	existingURL, err := r.GetURL(slug, clientIP)

	if err != nil {
		return nil, err
	}

	if existingURL != nil {
		return existingURL, nil
	}

	url = models.URL{ID: primitive.NewObjectID(), URL: originalURL, Slug: slug, CreatedAt: time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = urlCollection.InsertOne(ctx, url)

	if err != nil {
		return nil, err
	}

	if err = r.CacheURL(&shared.URL{URL: url.URL, Slug: url.Slug}, clientIP); err != nil {
		return nil, err
	}

	return &url, nil
}

// GetURL returns matching url objects for the given slug
func (r *Repository) GetURL(slug string, clientIP string) (*models.URL, error) {

	urlCollection := r.DB.Collection(collection)

	// Attempt to get the URL from the cache
	cachedURL, err := r.GetURLFromCache(slug, clientIP)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		log.Printf("Error getting url from cache: %v", err)
		return nil, err
	}

	if cachedURL != nil {
		return &models.URL{URL: cachedURL.URL, Slug: cachedURL.URL}, nil
	}

	var url models.URL

	filter := bson.M{"slug": slug}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = urlCollection.FindOne(ctx, filter).Decode(&url)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &url, nil
}

// CacheURL caches the given url
func (r *Repository) CacheURL(url *shared.URL, clientIP string) error {
	stringifiedURL, err := json.Marshal(url)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s-%s", url.Slug, clientIP)
	_, err = r.Cache.Set(key, stringifiedURL, 0).Result()

	if err != nil {
		log.Printf("Error caching URL; key, slug: %v", err)
		return err
	}

	key = fmt.Sprintf("%s-%s", url.URL, clientIP)
	_, err = r.Cache.Set(key, stringifiedURL, 0).Result()

	if err != nil {
		log.Printf("Error caching URL; key url: %v", err)
		return err
	}
	return nil
}

// GetURLFromCache returns the url from the cache using the given slug as the key
func (r *Repository) GetURLFromCache(cacheKey string, clientIP string) (*shared.URL, error) {
	var url shared.URL

	key := fmt.Sprintf("%s-%s", cacheKey, clientIP)

	result, err := r.Cache.Get(key).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &url)

	if err != nil {
		return nil, err
	}
	return &url, nil
}

// GenerateRandomNumber returns as an int a pseudo-random number for the given interval
func (r *Repository) GenerateRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// GenerateSlug returns as a string, an encoded form of the given url + timestamp + a pseudo-random number
// TODO: Review this approach to scale, i.e. reduce frequency of possible collisions
func (r *Repository) GenerateSlug(url string, min, max int) string {

	urlCollection := r.DB.Collection("url")

	var existingURL models.URL

	uniqueStr := fmt.Sprintf("%s+%d+%d", url, r.GenerateRandomNumber(min, max), time.Now().Unix())
	encodedStr := base64.RawURLEncoding.EncodeToString([]byte(uniqueStr))
	slug := encodedStr[len(encodedStr)-8:]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"slug": slug}
	err := urlCollection.FindOne(ctx, filter).Decode(&url)

	// Regenerate the slug if it already exists
	for err == nil {
		uniqueStr = fmt.Sprintf("%d%s%d", time.Now().Unix(), url, r.GenerateRandomNumber(min, max))

		encodedStr = base64.RawURLEncoding.EncodeToString([]byte(uniqueStr))

		slug = encodedStr[len(encodedStr)-8:]

		filter = bson.M{"slug": slug}
		err = urlCollection.FindOne(ctx, filter).Decode(&existingURL)
	}
	return slug
}
