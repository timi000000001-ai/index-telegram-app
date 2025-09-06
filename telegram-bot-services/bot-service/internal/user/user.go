package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bot-service/internal/config"

	"gopkg.in/telebot.v4"
)

// UserSaver is responsible for saving user data to the management service.
type UserSaver struct {
	ManagementServiceURL   string
	ManagementServiceToken string
	Client                 *http.Client
}

// NewUserSaver creates a new UserSaver.
func NewUserSaver(url, token string) *UserSaver {
	return &UserSaver{
		ManagementServiceURL:   url,
		ManagementServiceToken: token,
		Client:                 &http.Client{Timeout: 10 * time.Second},
	}
}

// SaveUser saves or updates a user's information in the management service.
func SaveUser(user *telebot.User) error {
	c, err := config.LoadConfig("development")
	if err != nil {
		return nil
	}
	collectionURL := fmt.Sprintf("%s/api/collections/tele_user/records", c.Bot.ManagementServiceURL)

	// Convert user ID to string for query and payload to avoid precision loss
	userIDStr := strconv.FormatInt(user.ID, 10)
	findURL := fmt.Sprintf("%s?filter=(tg_user_id='%s')", collectionURL, userIDStr)

	req, err := http.NewRequest(http.MethodGet, findURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create find request: %w", err)
	}
	req.Header.Set("Authorization", c.Bot.ManagementServiceToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to query user: %w", err)
	}
	defer resp.Body.Close()

	var findResult struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&findResult); err != nil {
		return fmt.Errorf("failed to decode find user response: %w", err)
	}

	// Use a map to construct the JSON payload to avoid struct marshalling issues.
	rawJSON, _ := json.Marshal(user)
	userData := map[string]interface{}{
		"tg_user_id":      userIDStr,
		"is_bot":          user.IsBot,
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
		"username":        user.Username,
		"language_code":   user.LanguageCode,
		"is_premium":      user.IsPremium,
		"can_join_groups": user.CanJoinGroups,
		"raw_json":        string(rawJSON),
	}

	if len(findResult.Items) > 0 {
		// User exists, update it
		userData["update_time"] = time.Now().UTC().Format(time.RFC3339)
		jsonData, err := json.Marshal(userData)
		if err != nil {
			return fmt.Errorf("failed to marshal user data for update: %w", err)
		}

		recordID := findResult.Items[0].ID
		updateURL := fmt.Sprintf("%s/%s", collectionURL, recordID)
		req, err := http.NewRequest(http.MethodPatch, updateURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create update request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", c.Bot.ManagementServiceToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to update user, status: %s", resp.Status)
		}
	} else {
		// User does not exist, create it
		userData["create_time"] = time.Now().UTC().Format(time.RFC3339)
		jsonData, err := json.Marshal(userData)
		if err != nil {
			return fmt.Errorf("failed to marshal user data for create: %w", err)
		}
		req, err := http.NewRequest(http.MethodPost, collectionURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", c.Bot.ManagementServiceToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create user, status: %s", resp.Status)
		}
	}

	return nil
}