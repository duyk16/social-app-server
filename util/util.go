package util

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type T map[string]interface{}

func SetResponseHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

func JSON(w http.ResponseWriter, status int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func PaginateList(r *http.Request) (page, limit int64) {
	var err error

	limit, err = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}

	page, err = strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 0
	}

	return page, limit
}
