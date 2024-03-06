// Package models maintains db table field types
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type URL struct {
	ID        primitive.ObjectID `bson:"_id"`
	URL       string             `json:"url" validate:"required,"`
	Slug      string             `json:"slug" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
