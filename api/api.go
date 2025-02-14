package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/duyk16/social-app-server/config"
	"github.com/duyk16/social-app-server/handler"
	"github.com/duyk16/social-app-server/util"
)

func Init() {
	r := mux.NewRouter()

	// Serve static file
	go func() {
		http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
		http.ListenAndServe(":"+config.ServerConfig.StaticPort, nil)
	}()

	// Middleware
	r.Use(util.JwtAuthentication)

	// API router
	setRouter(r)

	listen(r, config.ServerConfig.Port)
}

func setRouter(r *mux.Router) {
	r.HandleFunc("/api/auth/login", handler.LoginUser).Methods("POST")
	r.HandleFunc("/api/auth/register", handler.CreateUser).Methods("POST")

	// r.HandleFunc("/api/user/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	// r.HandleFunc("/api/user/", func(w http.ResponseWriter, r *http.Request) {}).Methods("POST")
	r.HandleFunc("/api/user/{id}", handler.GetUserByID).Methods("GET")
	r.HandleFunc("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("PUT")
	r.HandleFunc("/api/user/{id}/avatar", handler.UpdateAvatar).Methods("PUT")
	// r.HandleFunc("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("DELETE")

	r.HandleFunc("/api/post", handler.GetPosts).Methods("GET")
	r.HandleFunc("/api/post", handler.CreatePost).Methods("POST")
	r.HandleFunc("/api/post/{id}", handler.GetPostById).Methods("GET")
	r.HandleFunc("/api/post/{id}", handler.UpdatePost).Methods("PUT")
	r.HandleFunc("/api/post/{id}", handler.DeletePost).Methods("DELETE")
	r.HandleFunc("/api/post/{id}/like", handler.LikePost).Methods("PUT")
	r.HandleFunc("/api/post/{id}/unlike", handler.UnlikePost).Methods("PUT")
}

func listen(r *mux.Router, port string) {
	log.Printf("Server listening on port %v...\n", port)

	err := http.ListenAndServe(":"+port, r)

	if err != nil {
		log.Println("Serve server fail", err)
	}
}

func test(w http.ResponseWriter, r *http.Request) {

}
