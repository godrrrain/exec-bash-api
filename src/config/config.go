package config

import (
	"encoding/json"
	"os"
)

type Database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type Server struct {
	Port int `json:"port"`
}

type Config struct {
	Database `json:"database"`
	Server   `json:"server"`
}

// ReadConfig считывает конфиг
func ReadConfig() (Config, error) {

	var config Config

	config.Database.Host = "localhost"
	config.Database.Port = 5432
	config.Database.User = "postgres"
	config.Database.Password = "postgres"
	config.Database.Dbname = "postgres"
	config.Server.Port = 8080

	data, err := os.ReadFile("config.json")
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
