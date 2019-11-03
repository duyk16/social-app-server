package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	c "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/duyk16/social-app-server/model"
	"github.com/duyk16/social-app-server/storage"
	"github.com/duyk16/social-app-server/util"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	token := c.Get(r, "token").(model.Token)

	err, imagePath := util.UploadFileAnDeleteOld(r, "static/post", "p-*.png", "")

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  err.Error(),
		})
		return
	}

	var likes []primitive.ObjectID
	post := model.Post{
		ID:        primitive.NewObjectID(),
		Content:   r.FormValue("content"),
		Owner:     token.ID,
		Image:     imagePath,
		Likes:     likes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = storage.Post.InsertOne(context.Background(), &post)

	if err != nil {
		util.JSON(w, 500, util.T{
			"status": 2,
			"error":  err.Error(),
		})
		return
	}

	util.JSON(w, 201, util.T{
		"status": 0,
		"post": util.T{
			"id":    post.ID.Hex(),
			"image": post.Image,
		},
	})
	return

}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	page, limit := util.PaginateList(r)

	countChan := make(chan int64)
	postsChan := make(chan []model.Post)

	go func() {
		count, _ := storage.Post.CountDocuments(context.Background(), bson.M{})
		countChan <- count
	}()

	go func() {
		findOption := options.Find()
		findOption.SetSort(bson.M{"createdAt": -1})
		findOption.SetLimit(limit)
		findOption.SetSkip(page * limit)
		cur, err := storage.Post.Find(context.Background(), bson.M{}, findOption)

		if err != nil {
			util.JSON(w, 500, util.T{
				"status": 1,
				"error":  err.Error(),
			})
			return
		}

		defer cur.Close(context.Background())

		var posts []model.Post

		for cur.Next(context.Background()) {
			var post model.Post
			err = cur.Decode(&post)

			if err != nil {
				util.JSON(w, 500, util.T{
					"status": 2,
					"error":  "Decode fail",
				})
				return
			}

			posts = append(posts, post)
		}

		postsChan <- posts
	}()

	util.JSON(w, 200, util.T{
		"status": 0,
		"posts":  <-postsChan,
		"page":   page,
		"limit":  limit,
		"count":  <-countChan,
	})
}

func GetPostById(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  "id is not valid",
		})
	}

	var post model.Post

	err = storage.Post.FindOne(context.Background(), bson.M{"_id": id}).Decode(&post)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  err.Error(),
		})
		return
	}

	util.JSON(w, 200, util.T{
		"status": 0,
		"post":   post,
	})
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  "ID is not valid",
		})
		return
	}

	var postPost model.Post
	err = json.NewDecoder(r.Body).Decode(&postPost)
	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 2,
			"error":  "Body is not valid",
		})
		return
	}

	var post model.Post
	token := c.Get(r, "token").(model.Token)
	err = storage.Post.FindOne(
		context.Background(),
		bson.M{"_id": id, "owner": token.ID},
	).Decode(&post)
	if err != nil {
		util.JSON(w, 500, util.T{
			"status": 3,
			"error":  err.Error(),
		})
		return
	}

	err = storage.Post.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"content":   postPost.Content,
				"updatedAt": time.Now(),
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&post)

	if err != nil {
		util.JSON(w, 500, util.T{
			"status": 5,
			"error":  err.Error(),
		})
		return
	}

	util.JSON(w, 200, util.T{
		"status": 0,
		"post":   post,
	})
	return
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	token := c.Get(r, "token").(model.Token)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  "ID is not valid",
		})
		return
	}

	_, err = storage.Post.FindOneAndDelete(
		context.Background(),
		bson.M{"_id": id, "owner": token.ID},
	).DecodeBytes()

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 2,
			"error":  "Not found",
		})
		return
	}

	util.JSON(w, 203, util.T{
		"status": 2,
		"error":  "Not found",
	})
	return
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	token := c.Get(r, "token").(model.Token)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  "id is not valid",
		})
		return
	}

	var post model.Post

	err = storage.Post.FindOneAndUpdate(
		context.Background(),
		bson.M{
			"id": id,
		},
		bson.M{
			"$addToSet": bson.M{"likes": token.ID},
			"$set":      bson.M{"updatedAt": time.Now()},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(post)

	if err != nil {
		util.JSON(w, 500, util.T{
			"status": 2,
			"error":  err.Error(),
		})
		return
	}

	util.JSON(w, http.StatusAccepted, util.T{
		"status": 0,
		"post": util.T{
			"id":    post.ID.Hex(),
			"likes": post.Likes,
		},
	})
	return
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	token := c.Get(r, "token").(model.Token)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  "id is not valid",
		})
		return
	}

	var post model.Post

	err = storage.Post.FindOneAndUpdate(
		context.Background(),
		bson.M{
			"id": id,
		},
		bson.M{
			"$pull": bson.M{"likes": token.ID},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(post)

	if err != nil {
		util.JSON(w, 500, util.T{
			"status": 2,
			"error":  err.Error(),
		})
		return
	}

	util.JSON(w, http.StatusAccepted, util.T{
		"status": 0,
		"post": util.T{
			"id":    post.ID.Hex(),
			"likes": post.Likes,
		},
	})
	return
}
