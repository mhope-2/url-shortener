package repository

import (
	"context"
	"encoding/base64"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"

	"github.com/mhope-2/url_shortener/database/models"
	localmongo "github.com/mhope-2/url_shortener/database/mongo"
)

var (
	collection                      = "url"
	urlCollection *mongo.Collection = localmongo.OpenCollection(localmongo.Client(), collection)
)

// CreateUrl creates a url object and slug in the db
func (r *Repository) CreateUrl(originalUrl string, slug string) (*models.Url, error) {

	var url models.Url

	// check if the slug already exists
	existingUrl, err := r.GetUrlBySlug(slug)

	if err != nil {
		return nil, err
	}

	if existingUrl != nil {
		// a URL with the given slug already exists
		return existingUrl, nil
	}

	url = models.Url{ID: primitive.NewObjectID(), Url: originalUrl, Slug: slug, CreatedAt: time.Now()}

	_, err = urlCollection.InsertOne(context.TODO(), url)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

// GetUrlBySlug returns matching url objects for the given slug
func (r *Repository) GetUrlBySlug(slug string) (*models.Url, error) {

	var url models.Url

	filter := bson.M{"slug": slug}

	err := urlCollection.FindOne(context.TODO(), filter).Decode(&url)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &url, nil
}

// GetSlugByUrl returns matching url objects for the given url
//func (r *Repository) GetSlugByUrl(url string) (*models.Url, error) {
//
//	var result models.Url
//
//	filter := bson.M{"url": url}
//
//	err := urlCollection.FindOne(context.TODO(), filter).Decode(&result)
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return nil, nil
//		}
//		return nil, err
//	}
//
//	return &result, nil
//}

//func (r *Repository) InsertUrlIntoRedis()

// GenerateRandomNumber returns as an int a pseudo-random number for the given interval
func (r *Repository) GenerateRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// GenerateSlug returns as a string, an encoded form of the given url + timestamp + a pseudo-random number
func (r *Repository) GenerateSlug(url string, min, max int) string {
	var existingURL models.Url

	uniqueStr := fmt.Sprintf("%s+%d+%d", url, r.GenerateRandomNumber(min, max), time.Now().Unix())

	encodedStr := base64.URLEncoding.EncodeToString([]byte(uniqueStr))

	// TODO: Review approach to reduce collision rate and ensure only alphanumeric characters
	slug := encodedStr[len(encodedStr)-8:]

	filter := bson.M{"slug": slug}
	err := urlCollection.FindOne(context.TODO(), filter).Decode(&url)

	// Regenerate the slug if it already exists
	for err == nil {
		uniqueStr = fmt.Sprintf("%d%s%d", time.Now().Unix(), url, r.GenerateRandomNumber(min, max))

		encodedStr = base64.URLEncoding.EncodeToString([]byte(uniqueStr))

		slug = encodedStr[len(encodedStr)-8:]

		filter = bson.M{"slug": slug}
		err = urlCollection.FindOne(context.TODO(), filter).Decode(&existingURL)
	}

	return slug
}

func (r *Repository) InsertUrlIntoRedis(slug, shortenedUrl string) error {
	//stringifiedUser, err := json.Marshal(user)
	//if err != nil {
	//	return err
	//}
	//
	//key := fmt.Sprintf("courier_%s", user.Id)
	_, err = r.Store.Set(slug, shortenedUrl, 0).Result()

	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) GetUrlFromRedis(id string) (*shared.User, error) {
	var user shared.User
	key := fmt.Sprintf("courier_%s", id)

	result, err := r.RedisDB.Get(key).Result()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &user)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
