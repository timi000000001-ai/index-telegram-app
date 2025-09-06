package repository

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

// SearchRepositoryImpl implements the SearchRepository interface.
type searchRepositoryImpl struct {
	client         *resty.Client
	meilisearchURL string
	meilisearchKey string
}

// NewSearchRepository creates a new SearchRepository.
func NewSearchRepository(meilisearchURL, meilisearchKey string) SearchRepository {
	if meilisearchURL == "" || meilisearchKey == "" {
		log.Printf("ERROR: MeiliSearch URL or Key is empty")
		panic("MeiliSearch URL or Key cannot be empty")
	}
	return &searchRepositoryImpl{
		client:         resty.New(),
		meilisearchURL: meilisearchURL,
		meilisearchKey: meilisearchKey,
	}
}

// Search performs a search query against MeiliSearch.
func (s *searchRepositoryImpl) Search(query string, page int, limit int, filter string) ([]byte, error) {
	requestBody := map[string]interface{}{
		"q":           query,
		"page":        page,
		"hitsPerPage": limit,
		"sort":        []string{"MEMBERS_COUNT:desc"},
	}

	var meiliFilter string
	switch filter {
	case "all":
		// No filter, do nothing
	case "group":
		meiliFilter = "TYPE IN [group, supergroup]"
	case "channel":
		meiliFilter = "TYPE=channel"
	case "bot":
		meiliFilter = "TYPE=bot"
	case "message":
		meiliFilter = "MESSAGE_ID EXISTS"
	default:
		log.Printf("WARN: unknown filter type: %s", filter)
	}

	if meiliFilter != "" {
		requestBody["filter"] = meiliFilter
	}

	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+s.meilisearchKey).
		SetBody(requestBody).
		Post(s.meilisearchURL + "/indexes/telegram_index/search")

	if err != nil {
		log.Printf("ERROR: failed to send search request to MeiliSearch: %v", err)
		return nil, fmt.Errorf("failed to send search request to MeiliSearch: %w", err)
	}

	if resp.IsError() {
		log.Printf("ERROR: MeiliSearch returned an error: %s", resp.String())
		return nil, fmt.Errorf("MeiliSearch returned an error: %s", resp.String())
	}

	return resp.Body(), nil
}

// DeleteDocument deletes a document from MeiliSearch.
func (s *searchRepositoryImpl) DeleteDocument(docID string) error {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+s.meilisearchKey).
		Delete(s.meilisearchURL + "/indexes/telegram_index/documents/" + docID)

	if err != nil {
		log.Printf("ERROR: failed to send delete request to MeiliSearch: %v", err)
		return fmt.Errorf("failed to send delete request to MeiliSearch: %w", err)
	}

	if resp.IsError() {
		log.Printf("ERROR: MeiliSearch returned an error on delete: %s", resp.String())
		return fmt.Errorf("MeiliSearch returned an error on delete: %s", resp.String())
	}

	log.Printf("INFO: Document %s deleted successfully from MeiliSearch", docID)
	return nil
}
