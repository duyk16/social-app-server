package main

import (
	"github.com/duyk16/secure-rest-api/storage"
	"github.com/duyk16/social-app-server/api"
	"github.com/duyk16/social-app-server/config"
)

func main() {
	config.Init()
	storage.Init()
	api.Init()
}
