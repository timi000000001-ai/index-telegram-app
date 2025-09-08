package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/wangbin/jiebago"
)

const (
	meilisearchHost   = "http://127.0.0.1:7700"
	meilisearchAPIKey = "timigogogo"
	sourceIndexName   = "telegram_index"
	destIndexName     = "suggestions"
)

// MeiliSearchRecord defines the structure of documents in the source index.
type MeiliSearchRecord struct {
	ID          float64 `json:"id"`
	Title       string  `json:"TITLE"`
	Description string  `json:"DESCRIPTION"`
	Username    string  `json:"USERNAME"`
}

// SuggestionRecord defines the structure of documents in the suggestions index.
 type SuggestionRecord struct {
	ID    string `json:"id"`
	Query string `json:"query"`
}

// MeiliClientForWait defines an interface for waiting for tasks.
type MeiliClientForWait interface {
	WaitForTask(taskUID int64, interval time.Duration) (*meilisearch.Task, error)
}

// MeiliIndex defines an interface for index operations.
type MeiliIndex interface {
	GetDocuments(req *meilisearch.DocumentsQuery, resp *meilisearch.DocumentsResult) error
	AddDocuments(documents interface{}, primaryKey *string) (*meilisearch.TaskInfo, error)
}

var seg jiebago.Segmenter

func main() {
	// Initialize jiebago
	log.Println("INFO: Loading dictionary...")
	seg.LoadDictionary("dict.txt")
	log.Println("INFO: Dictionary loaded.")

	// Initialize Meilisearch client
	log.Println("INFO: Initializing Meilisearch client...")
	client := meilisearch.New(meilisearchHost, meilisearch.WithAPIKey(meilisearchAPIKey))

	// Delete the index if it exists to reset the primary key
	log.Println("INFO: Deleting existing 'suggestions' index if it exists...")
	task, err := client.DeleteIndex(destIndexName)
	if err == nil {
		log.Printf("INFO: Enqueued task %d to delete index '%s'.", task.TaskUID, destIndexName)
		log.Println("INFO: Waiting for index deletion task to complete...")
		finalTask, err := client.WaitForTask(task.TaskUID, 5*time.Second)
		if err != nil {
			log.Printf("WARN: Failed to wait for index deletion task completion: %v", err)
		} else if finalTask.Error.Message != "" {
			log.Printf("WARN: Meilisearch task to delete index failed: %s (%s)", finalTask.Error.Message, finalTask.Error.Code)
		} else {
			log.Println("INFO: Index deletion task completed.")
		}
	} else {
		log.Printf("WARN: Could not delete index %s (it may not exist or another error occurred): %v", destIndexName, err)
	}

	var sourceIndex MeiliIndex = client.Index(sourceIndexName)
	var destIndex MeiliIndex = client.Index(destIndexName)
	log.Println("INFO: Meilisearch client initialized.")

	log.Println("INFO: Starting to populate suggestions index...")

	// Get all documents from the source index
	documents, err := getAllDocuments(sourceIndex)
	if err != nil {
		log.Fatalf("ERROR: Failed to get documents: %v", err)
	}
	log.Printf("INFO: Found %d documents in index '%s'.", len(documents), sourceIndexName)

	if len(documents) == 0 {
		log.Println("WARN: No documents found in source index. Nothing to process.")
		return
	}

	// Extract keywords
	keywords := extractKeywords(documents)
	log.Printf("INFO: Extracted %d unique keywords.", len(keywords))

	if len(keywords) > 10 {
		log.Printf("INFO: Sample keywords: %v", keywords[:10])
	} else {
		log.Printf("INFO: Sample keywords: %v", keywords)
	}


	// Populate suggestions index
	if err := populateSuggestions(client, destIndex, keywords); err != nil {
		log.Fatalf("ERROR: Failed to populate suggestions: %v", err)
	}

	log.Println("INFO: Successfully populated suggestions index.")
}

// getAllDocuments retrieves all documents from the source index.
func getAllDocuments(index MeiliIndex) ([]MeiliSearchRecord, error) {
	var documents []MeiliSearchRecord
	limit := int64(1000)
	offset := int64(0)

	log.Println("INFO: Fetching documents from source index...")
	for {
		req := &meilisearch.DocumentsQuery{
			Limit:  limit,
			Offset: offset,
			Fields: []string{"id", "TITLE", "DESCRIPTION", "USERNAME"},
		}
		resp := new(meilisearch.DocumentsResult)

		if err := index.GetDocuments(req, resp); err != nil {
			return nil, fmt.Errorf("failed to get documents from Meilisearch at offset %d: %w", offset, err)
		}

		// Manually unmarshal the results
		var docs []MeiliSearchRecord
		jsonBytes, err := json.Marshal(resp.Results)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal document results: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &docs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal documents into MeiliSearchRecord: %w", err)
		}

		if len(docs) == 0 {
			log.Println("INFO: No more documents to fetch.")
			break
		}
		log.Printf("INFO: Fetched %d documents (offset: %d).", len(docs), offset)

		documents = append(documents, docs...)
		offset += limit
	}

	return documents, nil
}

// extractKeywords uses jiebago to cut words from the content of the documents.
func extractKeywords(documents []MeiliSearchRecord) []string {
	log.Printf("INFO: Extracting keywords from %d documents...", len(documents))
	keywordMap := make(map[string]struct{})
	for i, doc := range documents {
		fullContent := doc.Title + " " + doc.Description + " " + doc.Username
		words := seg.CutForSearch(fullContent, true)
		for word := range words {
			trimmedWord := strings.TrimSpace(word)
			if len(trimmedWord) > 1 {
				keywordMap[trimmedWord] = struct{}{}
			}
		}
		if (i+1)%100 == 0 {
			log.Printf("INFO: Processed %d/%d documents for keyword extraction.", i+1, len(documents))
		}
	}

	uniqueKeywords := make([]string, 0, len(keywordMap))
	for k := range keywordMap {
		uniqueKeywords = append(uniqueKeywords, k)
	}
	log.Printf("INFO: Finished keyword extraction. Found %d unique keywords.", len(uniqueKeywords))
	return uniqueKeywords
}

// populateSuggestions adds the extracted keywords to the suggestions index.
func populateSuggestions(client MeiliClientForWait, index MeiliIndex, keywords []string) error {
	if len(keywords) == 0 {
		log.Println("INFO: No keywords to populate.")
		return nil
	}
	log.Printf("INFO: Populating suggestions index with %d keywords...", len(keywords))
	suggestionRecords := make([]SuggestionRecord, len(keywords))
	for i, k := range keywords {
		h := sha1.New()
		h.Write([]byte(k))
		suggestionRecords[i] = SuggestionRecord{
			ID:    hex.EncodeToString(h.Sum(nil)),
			Query: k,
		}
	}

	primaryKey := "id"
	task, err := index.AddDocuments(suggestionRecords, &primaryKey)
	if err != nil {
		return fmt.Errorf("failed to add documents to Meilisearch: %w", err)
	}
	log.Printf("INFO: Enqueued task %d to add documents. Details: %+v", task.TaskUID, task)

	log.Println("INFO: Waiting for task to complete...")
	finalTask, err := client.WaitForTask(task.TaskUID, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for task completion: %w", err)
	}

	log.Printf("INFO: Task %d completed. Status: %s", finalTask.UID, finalTask.Status)
	if finalTask.Error.Message != "" {
		log.Printf("ERROR: Meilisearch task failed: %s (%s)", finalTask.Error.Message, finalTask.Error.Code)
	}

	return nil
}