package service

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// BotInfo represents the structure of a bot information record.
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type BotInfo struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	// Add other fields from your bot_info collection as needed
}

// BotInfoService defines the interface for fetching bot information.
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type BotInfoService interface {
	GetAllBotInfos() ([]BotInfo, error)
}

// botInfoServiceImpl implements the BotInfoService interface.
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type botInfoServiceImpl struct {
	client *resty.Client
	config PocketBaseConfig
}

// PocketBaseConfig holds the configuration for connecting to PocketBase.
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type PocketBaseConfig struct {
	BaseURL string
}

// NewBotInfoService creates a new BotInfoService instance.
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param config The PocketBase configuration.
// @return BotInfoService A new BotInfoService instance.
func NewBotInfoService(config PocketBaseConfig) BotInfoService {
	return &botInfoServiceImpl{
		client: resty.New(),
		config: config,
	}
}

// GetAllBotInfos fetches all bot information records from PocketBase.
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return []BotInfo A slice of bot information records.
// @return error An error if the request fails.
func (s *botInfoServiceImpl) GetAllBotInfos() ([]BotInfo, error) {
	resp, err := s.client.R().
		SetResult(&struct{ Items []BotInfo }{}).
		Get(fmt.Sprintf("%s/api/collections/bot_info/records", s.config.BaseURL))

	if err != nil {
		return nil, fmt.Errorf("failed to get bot infos: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get bot infos: status code %d", resp.StatusCode())
	}

	result := resp.Result().(*struct{ Items []BotInfo })
	return result.Items, nil
}