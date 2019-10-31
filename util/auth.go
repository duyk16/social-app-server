package util

import (
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/duyk16/social-app-server/config"
	"github.com/duyk16/social-app-server/model"
)

func HashAndSaltPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println("Hash password fail", err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func GenerateToken(userId primitive.ObjectID, email string) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), model.Token{
		ID:    userId,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})
	tokenString, _ := token.SignedString([]byte(config.ServerConfig.JWTKey))
	return tokenString
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuthRoutes := []string{
			"/api/auth",
		}

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range noAuthRoutes {
			if strings.Contains(r.URL.Path, value) {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			JSON(w, 400, T{
				"status":  "error",
				"message": "Missing token",
			})
			return
		}

		// `Bearer {token-body}`
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			JSON(w, 400, T{
				"status":  "error",
				"message": "Invalid/Malformed auth token",
			})
			return
		}

		tokenString := splitted[1]
		token := model.Token{}

		jwtToken, err := jwt.ParseWithClaims(tokenString, &token, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.ServerConfig.JWTKey), nil
		})

		if err != nil {
			JSON(w, http.StatusForbidden, T{
				"status":  "error",
				"message": "Malformed authentication token",
			})
			return
		}

		if !jwtToken.Valid {
			JSON(w, http.StatusForbidden, T{
				"status":  "error",
				"message": "Token is not valid.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
