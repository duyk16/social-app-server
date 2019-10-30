package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	FirstName       string             `json:"firstName"`
	LastName        string             `json:"lastName"`
	Email           string             `json:"email"`
	Password        string             `json:"password"`
	Avatar          Avatar             `json:"avatar"`
	PermissionLevel int                `json:"permissionLevel"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

type Avatar struct {
	Path     string `json:"path"`
	FileName string `json:"fileName"`
}
