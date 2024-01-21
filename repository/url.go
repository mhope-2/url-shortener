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
	localmongo "github.com/mhope-2/url_shortener/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	collection                      = "url"
	urlCollection *mongo.Collection = localmongo.GetCollection(localmongo.Client(), collection)
)

// CreateUrl creates a url object, stores it in the db and the caches it
func (r *Repository) CreateUrl(originalUrl string, slug string, clientIP string) (*models.Url, error) {

	var url models.Url

	// check if the slug already exists
	existingUrl, err := r.GetUrl(slug, clientIP)

	if err != nil {
		return nil, err
	}

	if existingUrl != nil {
		// a URL with the given slug already exists
		return existingUrl, nil
	}

	url = models.Url{ID: primitive.NewObjectID(), Url: originalUrl, Slug: slug, CreatedAt: time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = urlCollection.InsertOne(ctx, url)

	if err != nil {
		return nil, err
	}

	// Cache the url
	if err = r.CacheUrl(&shared.Url{Url: url.Url, Slug: url.Slug}, clientIP); err != nil {
		return nil, err
	}

	return &url, nil
}

// GetUrl returns matching url objects for the given slug
func (r *Repository) GetUrl(slug string, clientIP string) (*models.Url, error) {

	// Attempt to get the URL from the cache
	cachedUrl, err := r.GetUrlFromCache(slug, clientIP)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		log.Printf("Error getting url from cache: %v", err)
		return nil, err
	}

	if cachedUrl != nil {
		return &models.Url{Url: cachedUrl.Url, Slug: cachedUrl.Slug}, nil
	}

	var url models.Url

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

// TODO: Add client ip to cache key
func (r *Repository) CacheUrl(url *shared.Url, clientIP string) error {
	stringifiedUrl, err := json.Marshal(url)
	if err != nil {
		return err
	}

	// set slug cache key
	key := fmt.Sprintf("%s-%s", url.Slug, clientIP)
	_, err = r.Cache.Set(key, stringifiedUrl, 0).Result()

	// set original url cache key
	key = fmt.Sprintf("%s-%s", url.Url, clientIP)
	_, err = r.Cache.Set(key, stringifiedUrl, 0).Result()

	if err != nil {
		log.Printf("Error caching url: %v", err)
		return err
	}
	return nil
}

// GetUrlFromCache returns the url from the cache using the given slug as the key
func (r *Repository) GetUrlFromCache(cacheKey string, clientIP string) (*shared.Url, error) {
	var url shared.Url

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
// TODO: Review this approach to scale, i.e. reduce rate of possible collisions
func (r *Repository) GenerateSlug(url string, min, max int) string {
	var existingURL models.Url

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
