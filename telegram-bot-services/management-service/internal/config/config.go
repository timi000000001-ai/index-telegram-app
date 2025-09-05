package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// ServerConfig 定义了服务器的配置
type ServerConfig struct {
	Port string `json:"port" envconfig:"SERVER_PORT"`
}

// Config 定义了应用程序的配置结构
type Config struct {
	Server         ServerConfig `json:"server"`
	MeilisearchURL string       `json:"meilisearch_url" envconfig:"MEILISEARCH_URL"`
	MeilisearchKey string       `json:"meilisearch_key" envconfig:"MEILISEARCH_KEY"`
	PocketBaseURL  string       `json:"pocketbase_url" envconfig:"POCKETBASE_URL"`
	BotServiceURL  string       `json:"bot_service_url" envconfig:"BOT_SERVICE_URL"`
}

var (
	cfg  *Config
	once sync.Once
)

// Load 从文件和环境变量加载配置
func Load(configPath string) (*Config, error) {
	var err error
	once.Do(func() {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "development"
		}

		configFilePath := filepath.Join(configPath, fmt.Sprintf("%s.json", env))

		var c Config
		// 从文件加载配置
		configData, err := os.ReadFile(configFilePath)
		if err != nil {
			err = fmt.Errorf("failed to read config file: %w", err)
			return
		}

		if err = json.Unmarshal(configData, &c); err != nil {
			err = fmt.Errorf("failed to unmarshal config data: %w", err)
			return
		}

		// 从环境变量覆盖配置
		if err = envconfig.Process("", &c); err != nil {
			err = fmt.Errorf("failed to process env config: %w", err)
			return
		}

		cfg = &c
	})

	return cfg, err
}

// Get 返回已加载的配置实例
func Get() *Config {
	if cfg == nil {
		panic("config not loaded")
	}
	return cfg
}