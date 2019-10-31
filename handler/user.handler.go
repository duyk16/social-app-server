package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/duyk16/social-app-server/model"
	"github.com/duyk16/social-app-server/storage"
	"github.com/duyk16/social-app-server/util"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var postUser *model.User
	err := json.NewDecoder(r.Body).Decode(&postUser)

	if err != nil || postUser.Email == "" || postUser.Password == "" || postUser.FirstName == "" || postUser.LastName == "" {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  "Body is not valid",
		})
		return
	}

	user := model.User{
		ID:        primitive.NewObjectID(),
		FirstName: postUser.FirstName,
		LastName:  postUser.LastName,
		Email:     postUser.Email,
		Password:  util.HashAndSaltPassword(postUser.Password),
		Avatar:    model.Avatar{Path: "", FileName: ""},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = storage.User.InsertOne(context.Background(), user)

	if err != nil {
		mongoErr, ok := err.(mongo.WriteException)

		if ok && mongoErr.WriteErrors[0].Code == 11000 {
			util.JSON(w, 400, util.T{
				"status": 2,
				"code":   11000,
				"error":  "Email was used before",
			})
			return
		}

		util.JSON(w, 500, util.T{
			"status":  3,
			"message": "Insert user fail",
		})
		return
	}

	util.JSON(w, 201, util.T{
		"status": 0,
		"user":   user,
	})
	return
}

func GetUserByID() {

}
