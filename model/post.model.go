package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Content   string               `json:"content"`
	Likes     []primitive.ObjectID `json:"likes"`
	Image     string               `json:"image"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updateAt"`
}

