package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `file`
	file, handler, err := r.FormFile("file")
	if err != nil {
		// log.Printf("Error Retrieving the File %v", err)
		util.JSON(w, 400, util.T{
			"status": 2,
			"error":  "You must upload file",
		})
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %v\n", handler.Filename)
	// log.Printf("File Size: %+v\n", handler.Size)
	// log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a file
	tempFile, err := ioutil.TempFile("static/avatar", "u-*.png")
	if err != nil {
		log.Println(err)
		util.JSON(w, 500, util.T{
			"status": 3,
			"error":  "",
		})
		return
	}

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		util.JSON(w, 500, util.T{
			"status": 3,
			"error":  "",
		})
		return
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// log.Printf("Successfully Uploaded File\n")

	// remove old avatar
	os.Remove(user.Avatar)

	err = storage.User.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"avatar": tempFile.Name(),
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
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
			"avatar": "http://" + config.ServerConfig.ServerIP + config.ServerConfig.StaticPort + "/" + tempFile.Name(),
		},
	})
}
