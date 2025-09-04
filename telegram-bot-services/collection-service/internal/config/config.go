package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Telegram TelegramConfig `json:"telegram"`
	Storage  StorageConfig  `json:"storage"`
	Search   SearchConfig   `json:"search"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

type TelegramConfig struct {
	AppID   int    `json:"app_id"`
	AppHash string `json:"app_hash"`
}

type StorageConfig struct {
	PocketBaseURL string `json:"pocketBaseURL"`
}

type SearchConfig struct {
	MeilisearchURL  string `json:"meilisearch_url"`
	MeilisearchToken string `json:"meilisearch_token"`
	MessageLimit     int    `json:"message_limit"`
}

func Load(env string) (*Config, error) {
	path := fmt.Sprintf("configs/%s.json", env)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	cfg := &Config{}
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}