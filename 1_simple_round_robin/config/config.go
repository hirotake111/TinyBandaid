package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	ServerUrls []string `json:"server_urls"`
}

func LoadConfigFile(fileName string) (*Config, error) {
	var config *Config
	data, err := os.Open(fileName)
	if err != nil {
		return config, err
	}
	defer data.Close()
	byteArray, err := io.ReadAll(data)
	if err != nil {
		return config, err
	}
	json.Unmarshal(byteArray, &config)
	return config, nil
}
