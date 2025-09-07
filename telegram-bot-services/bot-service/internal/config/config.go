package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server  ServerConfig `json:"server"`
	Storage Storage      `json:"storage"`
	Search  SearchConfig `json:"search"`
	Bot     BotConfig    `json:"bot"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

// Storage defines the configuration for storage services.
type Storage struct {
	PocketBaseURL    string `json:"pocketBaseURL"`
	PocketBaseToken  string `json:"pocketBaseToken"`
	MeilisearchURL   string `json:"meilisearchURL"`
	MeilisearchToken string `json:"meilisearchToken"`
}

// Search defines the configuration for search services.
type SearchConfig struct {
	MeilisearchURL       string `json:"meilisearchURL"`
	MeilisearchKey       string `json:"meilisearchKey"`
	ManagementServiceURL string `json:"managementServiceURL"`
	IndexName            string `json:"indexName"`
}

type BotConfig struct {
	Token                  string   `json:"token"`
	WebhookURL             string   `json:"webhookURL"`
	APIEndpoint            string   `json:"apiEndpoint"`
	ManagementServiceURL   string   `json:"managementServiceURL"`
	ManagementServiceToken string   `json:"managementServiceToken"`
	ReviewChannel          string   `json:"reviewChannel"`
	ReviewBotToken         string   `json:"reviewBotToken"`
	BotTokens              []string `json:"bot_tokens"`
	TokenRotationDuration  int      `json:"token_rotation_duration"`
}

func LoadConfig(env string) (*Config, error) {
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