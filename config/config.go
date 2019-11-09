package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Threads    int     `json:"threads"`
	Name       string  `json:"name"`
	ServerIP   string  `json:"server_ip"`
	Port       string  `json:"port"`
	StaticPort string  `json:"file_static_port"`
	Storage    Storage `json:"storage"`
	JWTKey     string  `json:"jwt_key"`
	JWTExpire  int     `json:"jwt_expire"`
}

type Storage struct {
	Uri  string `json:"uri"`
	Name string `json:"name"`
}

var ServerConfig Config

func Init() {
	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&ServerConfig); err != nil {
		log.Fatal("Config error: ", err.Error())
	}

	// Set the number of operating system threads
	if ServerConfig.Threads > 0 {
		runtime.GOMAXPROCS(ServerConfig.Threads)
		log.Printf("Running with %v threads", ServerConfig.Threads)
	}
}
