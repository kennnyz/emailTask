package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DB             *DBConfig
	ServerAddr     string `json:"server_addr"`
	TimeToLiveLink string `json:"time_to_live_link"`
}

type DBConfig struct {
	Dsn string
}

type ServerConfig struct {
	serverAddr string
}

func ReadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	// decode json file ServerAddr and StorageType fields
	err = json.NewDecoder(file).Decode(&cfg)

	db := &DBConfig{}
	cfg.DB = db
	cfg.DB.Dsn = os.Getenv("DB_DSN")

	return &cfg, nil
}
