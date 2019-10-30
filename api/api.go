package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init() {
	r := mux.NewRouter()
	listen(r, "8000")
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
