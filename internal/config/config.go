package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbDsn              string `json:"db_dsn"`
	ServerAddr         string `json:"server_address"`
	TimeToLiveLink     int    `json:"time_to_live_link"`
	UserEmailFilesPath string `json:"users_email_files_path"`
}

func ReadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)

	return &cfg, nil
}
