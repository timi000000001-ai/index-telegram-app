package management

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"bot-service/internal/api/handler"
	"bot-service/internal/config"
)

// BotConfigResponse matches the structure of the response from the management service.
type BotConfigResponse struct {
	Items []handler.BotConfig `json:"items"`
}

// Client is a client for the management service.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new management service client.
func NewClient(baseURL string, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

// GetBotConfigs fetches bot configurations from the management service.
func GetBotsToken() ([]handler.BotConfig, error) {

	c, err := config.LoadConfig("development")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	url := c.Bot.ManagementServiceURL+"/api/collections/bot_info/records"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Bot.ManagementServiceToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
		} else {
			log.Printf("Failed to get bot configs: status %d, body: %s", resp.StatusCode, string(bodyBytes))
		}
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var botConfigResponse BotConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&botConfigResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return botConfigResponse.Items, nil
}

