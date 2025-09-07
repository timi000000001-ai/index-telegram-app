package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/meilisearch/meilisearch-go"
)

const (
	meilisearchHost     = "http://127.0.0.1:7700"
	meilisearchAPIKey   = "timigogogo"
	sourceIndexName     = "telegram_index"
	suggestionIndexName = "suggestions"
)

// MeiliSearchRecord matches the structure from the source index
type MeiliSearchRecord struct {
	ID          int64  `json:"id"`
	Title       string `json:"TITLE"`
	Description string `json:"DESCRIPTION"`
}

// SuggestionRecord is the structure for our new suggestions index
type SuggestionRecord struct {
	Query string `json:"query"`
}

func main() {
	client := meilisearch.New(meilisearchHost, meilisearch.WithAPIKey(meilisearchAPIKey))

	log.Println("Starting suggestion index population...")

	// 1. Get all documents from the source index
	documents, err := getAllDocuments(client, sourceIndexName)
	if err != nil {
		log.Fatalf("FATAL: Could not retrieve documents from '%s': %v", sourceIndexName, err)
	}
	log.Printf("Found %d documents in index '%s'", len(documents), sourceIndexName)

	// 2. Extract keywords
	keywords := extractKeywords(documents)
	log.Printf("Extracted %d unique keywords", len(keywords))

	// 3. Populate the suggestions index
	if err := populateSuggestions(client, suggestionIndexName, keywords); err != nil {
		log.Fatalf("FATAL: Could not populate suggestions index: %v", err)
	}

	log.Println("Suggestion index population finished successfully!")
}

func getAllDocuments(client *meilisearch.Client, indexName string) ([]MeiliSearchRecord, error) {
	var allDocuments []MeiliSearchRecord
	limit := int64(1000)
	offset := int64(0)

	index := client.Index(indexName)

	for {
		var documents []MeiliSearchRecord
		if err := index.GetDocuments(&meilisearch.DocumentsQuery{
			Limit:  limit,
			Offset: offset,
		}, &documents); err != nil {
			return nil, err
		}

		if len(documents) == 0 {
			break
		}

		allDocuments = append(allDocuments, documents...)
		offset += limit
	}

	return allDocuments, nil
}

func extractKeywords(documents []MeiliSearchRecord) map[string]bool {
	keywords := make(map[string]bool)
	// Simple regex to split by spaces and punctuation
	re := regexp.MustCompile(`[\s,.;:!?()\[\]{}]+`)

	for _, doc := range documents {
		// Extract from title
		for _, word := range re.Split(doc.Title, -1) {
			if len(word) > 2 { // Basic filtering for very short words
				keywords[strings.ToLower(word)] = true
			}
		}
		// Extract from description
		for _, word := range re.Split(doc.Description, -1) {
			if len(word) > 2 {
				keywords[strings.ToLower(word)] = true
			}
		}
	}
	return keywords
}

func populateSuggestions(client *meilisearch.Client, indexName string, keywords map[string]bool) error {
	var suggestionRecords []SuggestionRecord
	for keyword := range keywords {
		suggestionRecords = append(suggestionRecords, SuggestionRecord{Query: keyword})
	}

	index := client.Index(indexName)

	// Clear the index first
	task, err := index.DeleteAllDocuments()
	if err != nil {
		var meiliErr *meilisearch.Error
		if errors.As(err, &meiliErr) && meiliErr.ErrCode == meilisearch.ErrCodeIndexNotFound {
			// Index doesn't exist, which is fine, we will create it later
		} else {
			return fmt.Errorf("failed to clear old documents: %w", err)
		}
	} else {
		// Wait for the deletion task to complete
		if _, err := client.WaitForTask(task.TaskUID); err != nil {
			return fmt.Errorf("failed while waiting for document deletion: %w", err)
		}
	}

	// Add new documents
	task, err = index.AddDocuments(suggestionRecords, "query")
	if err != nil {
		return fmt.Errorf("failed to add documents: %w", err)
	}

	log.Printf("Task UID %d to add %d documents has been enqueued.", task.TaskUID, len(suggestionRecords))

	// Configure settings for the suggestions index
	settings := meilisearch.Settings{
		SearchableAttributes: []string{"query"},
		SortableAttributes:   []string{},
		FilterableAttributes: []string{},
		RankingRules: []string{
			"words",
			"typo",
			"proximity",
			"attribute",
			"sort",
			"exactness",
		},
	}

	task, err = index.UpdateSettings(&settings)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	log.Printf("Task UID %d to update settings has been enqueued.", task.TaskUID)
	return nil
}