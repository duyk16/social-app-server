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

	// Serve static folder in server
	r.Handle("/static", http.FileServer(http.Dir("./static")))

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
	r.HandleFunc("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	r.HandleFunc("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("PUT")
	r.HandleFunc("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("DELETE")

	r.HandleFunc("/api/post/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	r.HandleFunc("/api/post/", func(w http.ResponseWriter, r *http.Request) {}).Methods("POST")
	r.HandleFunc("/api/post/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	r.HandleFunc("/api/post/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("PUT")
	r.HandleFunc("/api/post/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("DELETE")
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
