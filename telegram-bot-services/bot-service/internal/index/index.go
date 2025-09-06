package index

import (
	"bot-service/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// SaveTelegramIndex saves the telegram index data to the management service.
func SaveTelegramIndex(cfg *config.Config, data map[string]interface{}) error {
	baseURL := cfg.Bot.ManagementServiceURL + "/api/collections/telegram_index/records"
	chatID, ok := data["chat_id"].(string)
	if !ok {
		return fmt.Errorf("chat_id not found or not a string")
	}

	// 查询是否存在
	queryURL := baseURL + "?filter=(chat_id='" + chatID + "')"
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create query request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Bot.ManagementServiceToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("query failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode query response: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	var method, targetURL string
	if len(result.Items) > 0 {
		// 更新
		method = "PATCH"
		targetURL = baseURL + "/" + result.Items[0].ID
	} else {
		// 插入
		method = "POST"
		targetURL = baseURL
	}

	req, err = http.NewRequest(method, targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Bot.ManagementServiceToken))

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform %s: %w", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%s failed with status: %d", method, resp.StatusCode)
	}

	// Check if index exists, create if not with primaryKey 'id'
	indexURL := cfg.Storage.MeilisearchURL + "/indexes/telegram_index"
	indexReq, err := http.NewRequest("GET", indexURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create index check request: %w", err)
	}
	indexReq.Header.Set("Authorization", "Bearer " + cfg.Storage.MeilisearchToken)

	indexResp, err := client.Do(indexReq)
	if err != nil {
		return fmt.Errorf("failed to check index: %w", err)
	}
	defer indexResp.Body.Close()

	if indexResp.StatusCode == http.StatusNotFound {
		// Create index with primaryKey 'id'
		createPayload := map[string]string{"uid": "telegram_index", "primaryKey": "id"}
		createJSON, err := json.Marshal(createPayload)
		if err != nil {
		return fmt.Errorf("failed to marshal create index payload: %w", err)
		}
		createReq, err := http.NewRequest("POST", cfg.Storage.MeilisearchURL + "/indexes", bytes.NewBuffer(createJSON))
		if err != nil {
		return fmt.Errorf("failed to create index creation request: %w", err)
		}
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("Authorization", "Bearer " + cfg.Storage.MeilisearchToken)

		createResp, err := client.Do(createReq)
		if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
		}
		defer createResp.Body.Close()

		if createResp.StatusCode != http.StatusAccepted {
			body, _ := io.ReadAll(createResp.Body)
			return fmt.Errorf("index creation failed with status: %d, body: %s", createResp.StatusCode, string(body))
		}
		// Wait for task to complete (simplified, in production use task status)
		time.Sleep(2 * time.Second)

		// Update index settings
		if err := updateIndexSettings(cfg); err != nil {
			// Log the error but don't fail the whole process, as settings can be applied later
			log.Printf("failed to update meilisearch settings: %v", err)
		}
	}

	// Index to MeiliSearch
	meiliURL := cfg.Storage.MeilisearchURL + "/indexes/telegram_index/documents"
	docData := make(map[string]interface{})
	for k, v := range data {
		docData[k] = v
	}
	docData["id"] = chatID
	delete(docData, "chat_id")
	docs := []map[string]interface{}{docData}
	docsJSON, err := json.Marshal(docs)
	if err != nil {
		return fmt.Errorf("failed to marshal MeiliSearch docs: %w", err)
	}
	meiliReq, err := http.NewRequest("POST", meiliURL, bytes.NewBuffer(docsJSON))
	if err != nil {
		return fmt.Errorf("failed to create MeiliSearch request: %w", err)
	}
	meiliReq.Header.Set("Content-Type", "application/json")
	meiliReq.Header.Set("Authorization", "Bearer " + cfg.Storage.MeilisearchToken)

	meiliResp, err := client.Do(meiliReq)
	if err != nil {
		return fmt.Errorf("failed to index to MeiliSearch: %w", err)
	}
	defer meiliResp.Body.Close()

	body, _ := io.ReadAll(meiliResp.Body)
	log.Printf("MeiliSearch response: status=%d, body=%s", meiliResp.StatusCode, string(body))

	if meiliResp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("MeiliSearch indexing failed with status: %d, body: %s", meiliResp.StatusCode, string(body))
	}

	return nil
}

// updateIndexSettings updates the MeiliSearch index settings.
func updateIndexSettings(cfg *config.Config) error {
	settings := map[string]interface{}{
		"searchableAttributes": []string{
			"title",
			"username",
			"description",
			"first_name",
			"last_name",
		},
		"filterableAttributes": []string{
			"type",
			"is_verified",
			"is_restricted",
			"is_scam",
			"is_fake",
			"language_code",
			"tags",
			"content_types",
			"members_count",
			"sender_is_bot",
		},
		"rankingRules": []string{
			"words",
			"typo",
			"proximity",
			"attribute",
			"sort",
			"exactness",
			"members_count:desc",
		},
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal index settings: %w", err)
	}

	settingsURL := cfg.Storage.MeilisearchURL + "/indexes/telegram_index/settings"
	req, err := http.NewRequest("POST", settingsURL, bytes.NewBuffer(settingsJSON))
	if err != nil {
		return fmt.Errorf("failed to create settings update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + cfg.Storage.MeilisearchToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update index settings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("index settings update failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Println("MeiliSearch index settings updated successfully.")
	return nil
}