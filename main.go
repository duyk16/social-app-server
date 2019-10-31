package main

import (
	"github.com/duyk16/social-app-server/api"
	"github.com/duyk16/social-app-server/config"
	"github.com/duyk16/social-app-server/storage"
)

func main() {
	config.Init()
	storage.Init()
	api.Init()
}
