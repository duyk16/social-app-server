package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	Content   string               `json:"content"`
	Owner     primitive.ObjectID   `json:"owner"`
	Likes     []primitive.ObjectID `json:"likes"`
	Image     string               `json:"image"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time            `json:"updateAt" bson:"updatedAt"`
}

type PostHome struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	Content   string               `json:"content"`
	Owner     Owner                `json:"owner"`
	Likes     []primitive.ObjectID `json:"likes"`
	Image     string               `json:"image"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time            `json:"updateAt" bson:"updatedAt"`
}

type Owner struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	Avatar    string             `json:"avatar"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	Avatar    string             `json:"avatar"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// type Avatar struct {
// 	Path     string `json:"path"`
// 	FileName string `json:"fileName"`
// }

type Token struct {
	ID    primitive.ObjectID `json:"userId"`
	Email string             `json:"email"`
	jwt.StandardClaims
}
