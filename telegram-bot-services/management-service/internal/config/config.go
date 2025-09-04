package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"github.com/kelseyhightower/envconfig"
)

// Config 定义了应用程序的配置结构
type Config struct {
	MeilisearchURL string `json:"meilisearch_url" envconfig:"MEILISEARCH_URL"`
	MeilisearchKey string `json:"meilisearch_key" envconfig:"MEILISEARCH_KEY"`
}

var (
	cfg  *Config
	once sync.Once
)

// Load 从文件和环境变量加载配置
// 它首先会根据 `APP_ENV` 环境变量（默认为 "development"）查找对应的 JSON 配置文件。
// 然后，它会使用环境变量覆盖 JSON 文件中的值。
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
		// 如果配置未加载，可以返回一个默认配置或触发一个 panic
		// 这里我们选择 panic，因为在没有配置的情况下应用程序无法正常工作
		panic("config not loaded")
	}
	return cfg
}