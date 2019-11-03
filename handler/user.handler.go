package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/duyk16/social-app-server/config"
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
		Avatar:    "",
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
		"userId": user.ID,
	})
	return
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, _ := primitive.ObjectIDFromHex(params["id"])

	var user model.User
	err := storage.User.FindOne(context.Background(), bson.M{"_id": userId}).Decode(&user)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  err.Error(),
		})
		return
	}

	user.Password = ""

	util.JSON(w, 200, util.T{
		"statusr": 0,
		"user": util.T{
			"id":        user.ID.Hex(),
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"avatar":    user.Avatar,
		},
	})
	return
}

func UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, _ := primitive.ObjectIDFromHex(params["id"])

	var user model.User
	err := storage.User.FindOne(context.Background(), bson.M{"_id": userId}).Decode(&user)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  err.Error(),
		})
		return
	}

	err, path := util.UploadFileAnDeleteOld(r, "static/avatar", "u-*.png", user.Avatar)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 2,
			"error":  err.Error(),
		})

		return
	}

	err = storage.User.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"avatar": path,
			},
		},
		options.
			FindOneAndUpdate().
			SetReturnDocument(options.After),
	).Decode(&user)

	if err != nil {
		util.JSON(w, 500, util.T{
			"status": 4,
			"error":  "",
		})
		return
	}

	util.JSON(w, 200, util.T{
		"status": 0,
		"user": util.T{
			"id":     user.ID.Hex(),
			"avatar": "http://" + config.ServerConfig.ServerIP + config.ServerConfig.StaticPort + "/" + path,
		},
	})
}
