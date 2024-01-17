package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Url struct {
	ID        primitive.ObjectID `bson:"_id"`
	Url       string             `json:"url" validate:"required,"`
	Slug      string             `json:"slug" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
