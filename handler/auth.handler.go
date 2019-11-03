package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/duyk16/social-app-server/model"
	"github.com/duyk16/social-app-server/storage"
	"github.com/duyk16/social-app-server/util"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var postUser model.User
	err := json.NewDecoder(r.Body).Decode(&postUser)

	if err != nil || postUser.Email == "" || postUser.Password == "" {
		util.JSON(w, 400, util.T{
			"status": 1,
			"error":  "Body is not valid",
		})
		return
	}

	var user model.User
	err = storage.User.FindOne(context.Background(), bson.M{"email": postUser.Email}).Decode(&user)

	if err != nil {
		util.JSON(w, 400, util.T{
			"status": 3,
			"error":  "Not found email",
		})
		return
	}

	if !util.ComparePasswords(user.Password, postUser.Password) {
		util.JSON(w, 400, util.T{
			"status": 4,
			"error":  "Password is incorrect",
		})
		return
	}

	token := util.GenerateToken(user.ID, user.Email)

	util.JSON(w, 200, util.T{
		"status": 0,
		"token":  token,
		"userId": user.ID,
	})

}
