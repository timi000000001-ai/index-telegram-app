package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// --- Configuration ---
const (
	pocketbaseURL      = "http://127.0.0.1:8090"
	pocketbaseToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb2xsZWN0aW9uSWQiOiJwYmNfMzE0MjYzNTgyMyIsImV4cCI6MTg0NjExNzczMiwiaWQiOiJwaWJvcjZzb3BqaDNrN3MiLCJyZWZyZXNoYWJsZSI6ZmFsc2UsInR5cGUiOiJhdXRoIn0.2jGluOIh9pvcY7hTeTWYctVzrMIBZ2vzvh3_aGSHMGQ"
	meilisearchHost    = "http://127.0.0.1:7700"
	meilisearchAPIKey  = "timigogogo"
	meilisearchIndex   = "telegram_index"
	sqlFilePath        = "/root/新建文件夹/index-telegram-app/sql/tele_index_record.sql"
	pocketbaseCollection = "telegram_index"
)

// --- Structs ---

// PocketBaseRecord matches the structure of the 'telegram_index' collection
type PocketBaseRecord struct {
	ID           string `json:"id,omitempty"`
	Type         string `json:"type"`
	ChatID       int64  `json:"chat_id"`
	Title        string `json:"title"`
	Username     string `json:"username,omitempty"`
	Description  string `json:"description,omitempty"`
	MembersCount int64  `json:"members_count,omitempty"`
	IsVerified   bool   `json:"is_verified,omitempty"`
	// Add other fields from your PocketBase collection structure if needed
}

// MeiliSearchRecord matches the structure for MeiliSearch
type MeiliSearchRecord struct {
	ID           int64  `json:"id"` // Mapped from chat_id
	Title        string `json:"TITLE"`
	Description  string `json:"DESCRIPTION"`
	Username     string `json:"USERNAME"`
	Type         string `json:"TYPE"`
	MembersCount int64  `json:"MEMBERS_COUNT"`
	Link         string `json:"link"`
	// Add other fields from your MeiliSearch structure if needed
}

// PocketBaseListResult is used to decode the list response from PocketBase
type PocketBaseListResult struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}

func main() {
	// Clear existing data first
	// if err := clearMeiliSearchIndex(); err != nil {
	// 	log.Fatalf("FATAL: Could not clear MeiliSearch index: %v", err)
	// }
	// if err := clearPocketBaseCollection(); err != nil {
	// 	log.Fatalf("FATAL: Could not clear PocketBase collection: %v", err)
	// }

	file, err := os.Open(sqlFilePath)
	if err != nil {
		log.Fatalf("FATAL: Cannot open SQL file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxCapacity = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	log.Println("Starting data import...")

	var statement strings.Builder
	insideInsert := false

	const numWorkers = 10
	jobs := make(chan map[string]string, 100) // Buffered channel
	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, &wg)
	}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "INSERT INTO `tele_index_record`") {
			insideInsert = true
		}

		if insideInsert {
			statement.WriteString(line)
			statement.WriteString("\n") // Preserve newlines

			if strings.HasSuffix(strings.TrimSpace(line), ";") {
				// End of statement
				fullStatement := statement.String()
				statement.Reset()
				insideInsert = false

				// Extract column names and values
				valuesIndex := strings.Index(fullStatement, "VALUES")
				if valuesIndex == -1 {
					log.Println("WARN: Found INSERT statement without VALUES clause, skipping.")
					continue
				}

				// Extract column names
				colStr := fullStatement[:valuesIndex]
				startCols := strings.Index(colStr, "(")
				endCols := strings.LastIndex(colStr, ")")
				if startCols == -1 || endCols == -1 {
					log.Println("WARN: Could not parse column names from INSERT statement.")
					continue
				}
				cols := strings.Split(colStr[startCols+1:endCols], ",")
				for i, col := range cols {
					cols[i] = strings.Trim(strings.TrimSpace(col), "`")
				}

				// Extract values string
				valuesStr := fullStatement[valuesIndex+len("VALUES"):]
				valuesStr = strings.TrimSuffix(strings.TrimSpace(valuesStr), ";")

				recordsStr := splitSqlValues(valuesStr)
				for _, recordStr := range recordsStr {
					rawValues := parseSqlValues(recordStr)
					if len(rawValues) != len(cols) {
						log.Printf("WARN: Column/value mismatch. Got %d columns and %d values. Skipping record.", len(cols), len(rawValues))
						continue
					}

					// Create a map of col -> value
					valueMap := make(map[string]string)
					for i, col := range cols {
						valueMap[col] = rawValues[i]
					}
					wg.Add(1)
					jobs <- valueMap
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("FATAL: Error reading SQL file: %v", err)
	}

	close(jobs)
	wg.Wait()

	log.Println("Data import finished.")
}

func worker(id int, jobs <-chan map[string]string, wg *sync.WaitGroup) {
	for j := range jobs {
		processRecord(j)
		wg.Done()
	}
}

func processRecord(valueMap map[string]string) {
	// --- Map SQL values to structs ---

	chatID, _ := strconv.ParseInt(valueMap["chat_id"], 10, 64)
	if chatID == 0 {
		log.Println("WARN: Skipping record with empty chat_id")
		return
	}
	members, _ := strconv.ParseInt(valueMap["members"], 10, 64)
	isVerified := valueMap["isofficial"] == "Y"

	rawType := valueMap["type"]
	var recordType string
	switch strings.TrimSpace(rawType) {
	case "频道":
		recordType = "channel"
	case "群组":
		recordType = "group"
	case "机器人":
		recordType = "bot"
	case "超级群组":
		recordType = "supergroup"
	default:
		log.Printf("WARN: Unknown record type: '%s', defaulting to 'group'", rawType)
		recordType = "group"
	}

	username := valueMap["username"]
	if username == "" {
		link := valueMap["link"]
		if strings.HasPrefix(link, "https://t.me/") {
			path := strings.TrimPrefix(link, "https://t.me/")
			if parts := strings.Split(path, "/"); len(parts) > 0 {
				username = parts[0]
			}
		}
	}

	// PocketBase Record
	pbRecord := PocketBaseRecord{
		Type:         recordType,
		ChatID:       chatID,
		Title:        valueMap["title"],
		Description:  valueMap["description"],
		Username:     username,
		MembersCount: members,
		IsVerified:   isVerified,
	}

	// MeiliSearch Record
	meiliRecord := MeiliSearchRecord{
		ID:           chatID,
		Title:        valueMap["title"],
		Description:  valueMap["description"],
		Username:     username,
		Type:         recordType,
		MembersCount: members,
		Link:         valueMap["link"],
	}

	// --- Send to Services ---
	if err := sendToPocketBase(pbRecord); err != nil {
		log.Printf("ERROR: Failed to send record to PocketBase (ChatID: %d): %v", chatID, err)
	}

	if err := sendToMeiliSearch(meiliRecord); err != nil {
		log.Printf("ERROR: Failed to send record to MeiliSearch (ChatID: %d): %v", chatID, err)
	} else {
		log.Printf("INFO: Successfully sent record to MeiliSearch (ChatID: %d)", chatID)
	}
}

func sendToPocketBase(record PocketBaseRecord) error {
	// 1. Check if record exists by chat_id
	client := &http.Client{Timeout: 10 * time.Second}
	filter := fmt.Sprintf("(chat_id='%d')", record.ChatID)
	queryURL := fmt.Sprintf("%s/api/collections/%s/records?filter=%s", pocketbaseURL, pocketbaseCollection, filter)

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return fmt.Errorf("could not create check request: %w", err)
	}
	req.Header.Set("Authorization", pocketbaseToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("check request returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var listResult PocketBaseListResult
	if err := json.NewDecoder(resp.Body).Decode(&listResult); err != nil {
		return fmt.Errorf("could not decode check response: %w", err)
	}

	// 2. Marshal data for POST or PATCH
	jsonData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("could not marshal record: %w", err)
	}

	var method, url string
	if len(listResult.Items) > 0 {
		// Record exists, update it
		recordID := listResult.Items[0].ID
		method = "PATCH"
		url = fmt.Sprintf("%s/api/collections/%s/records/%s", pocketbaseURL, pocketbaseCollection, recordID)
	} else {
		// Record does not exist, create it
		method = "POST"
		url = fmt.Sprintf("%s/api/collections/%s/records", pocketbaseURL, pocketbaseCollection)
	}

	req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("could not create %s request: %w", method, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", pocketbaseToken)

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("%s request failed: %w", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-2xx status code for %s: %d, body: %s", method, resp.StatusCode, string(bodyBytes))
	}

	if method == "PATCH" {
		log.Printf("INFO: Successfully updated record in PocketBase (ChatID: %d)", record.ChatID)
	} else {
		log.Printf("INFO: Successfully created record in PocketBase (ChatID: %d)", record.ChatID)
	}

	return nil
}

func sendToMeiliSearch(record MeiliSearchRecord) error {
	// MeiliSearch expects a list of documents
	records := []MeiliSearchRecord{record}
	jsonData, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("could not marshal record: %w", err)
	}

	url := fmt.Sprintf("%s/indexes/%s/documents?primaryKey=id", meilisearchHost, meilisearchIndex)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+meilisearchAPIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var bodyBytes []byte
		resp.Body.Read(bodyBytes)
		return fmt.Errorf("received non-2xx status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// splitSqlValues splits a string of multiple SQL value tuples like "(...),(...)"
func splitSqlValues(s string) []string {
	var result []string
	level := 0
	start := 0
	inString := false

	for i, r := range s {
		if r == '\'' {
			inString = !inString
		}
		if !inString {
			if r == '(' {
				if level == 0 {
					start = i
				}
				level++
			} else if r == ')' {
				level--
				if level == 0 {
					result = append(result, s[start:i+1])
				}
			}
		}
	}
	return result
}

// parseSqlValues parses a single SQL value tuple like "('a', 1, NULL)"
func parseSqlValues(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")

	var values []string
	var currentVal strings.Builder
	inString := false
	
	for i := 0; i < len(s); i++ {
		char := s[i]
		if char == '\'' {
			inString = !inString
		} else if char == ',' && !inString {
			values = append(values, strings.TrimSpace(currentVal.String()))
			currentVal.Reset()
		} else {
			currentVal.WriteByte(char)
		}
	}
	values = append(values, strings.TrimSpace(currentVal.String()))

	// Unquote strings and handle NULL
	for i, v := range values {
		if strings.ToUpper(v) == "NULL" {
			values[i] = ""
		} else if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			values[i] = strings.TrimSpace(v[1 : len(v)-1])
		}
	}

	return values
}

// clearPocketBaseCollection deletes all records from the PocketBase collection.
func clearPocketBaseCollection() error {
	log.Println("INFO: Clearing all records from PocketBase collection:", pocketbaseCollection)

	// 1. Get all record IDs
	client := &http.Client{Timeout: 30 * time.Second}
	var allRecords []struct {
		ID string `json:"id"`
	}

	page := 1
	for {
		// Fetch a page of records
		listURL := fmt.Sprintf("%s/api/collections/%s/records?perPage=500&page=%d&fields=id", pocketbaseURL, pocketbaseCollection, page)
		req, err := http.NewRequest("GET", listURL, nil)
		if err != nil {
			return fmt.Errorf("could not create list request for page %d: %w", page, err)
		}
		req.Header.Set("Authorization", pocketbaseToken)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("list request failed for page %d: %w", page, err)
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("list request returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
		}

		var listResult struct {
			Items      []struct{ ID string `json:"id"` } `json:"items"`
			TotalItems int                          `json:"totalItems"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&listResult); err != nil {
			resp.Body.Close()
			return fmt.Errorf("could not decode list response for page %d: %w", page, err)
		}
		resp.Body.Close()

		if len(listResult.Items) == 0 {
			break // No more records
		}

		allRecords = append(allRecords, listResult.Items...)
		page++
	}

	log.Printf("INFO: Found %d records to delete.", len(allRecords))

	// 2. Delete each record
	for _, record := range allRecords {
		deleteURL := fmt.Sprintf("%s/api/collections/%s/records/%s", pocketbaseURL, pocketbaseCollection, record.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		if err != nil {
			log.Printf("WARN: Could not create delete request for record %s: %v", record.ID, err)
			continue // Continue to next record
		}
		req.Header.Set("Authorization", pocketbaseToken)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("WARN: Delete request for record %s failed: %v", record.ID, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			log.Printf("WARN: Delete request for record %s returned status %d", record.ID, resp.StatusCode)
		}
	}

	log.Println("INFO: Finished clearing PocketBase collection.")
	return nil
}

// clearMeiliSearchIndex deletes all documents from the MeiliSearch index.
func clearMeiliSearchIndex() error {
	log.Println("INFO: Clearing all documents from MeiliSearch index:", meilisearchIndex)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/indexes/%s/documents", meilisearchHost, meilisearchIndex), nil)
	if err != nil {
		return fmt.Errorf("could not create delete request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+meilisearchAPIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var bodyBytes []byte
		resp.Body.Read(bodyBytes)
		return fmt.Errorf("received non-2xx status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Println("INFO: Finished clearing PocketBase collection.")
	return nil
}