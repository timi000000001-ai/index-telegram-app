package index
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"bot-service/internal/config"
)

// SaveTelegramIndex saves the telegram index data to the management service.
func SaveTelegramIndex(data map[string]interface{}) error {
	cfg, err := config.LoadConfig("development")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

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

	// Index to MeiliSearch
	meiliURL := cfg.Storage.MeilisearchURL + "/indexes/telegram_index/documents"
	docData := make(map[string]interface{})
	for k, v := range data {
		docData[k] = v
	}
	docData["id"] = chatID
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

	if meiliResp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("MeiliSearch indexing failed with status: %d", meiliResp.StatusCode)
	}

	return nil
}