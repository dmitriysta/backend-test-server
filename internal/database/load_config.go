package database

import (
	"back-end/logs"
	"encoding/json"
	"os"
)

type ServerConfig struct {
	Port string `json:"port"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	User     string `json:"user"`
}

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

func LoadConfig(confFile string) Config {
	var config Config
	openFile, err := os.Open(confFile)
	if err != nil {
		logs.LogError("Open config file", err.Error())
	}
	defer openFile.Close()

	parsing := json.NewDecoder(openFile)
	parsing.Decode(&config)

	return config
}
