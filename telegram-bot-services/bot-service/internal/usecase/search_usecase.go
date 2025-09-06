package usecase

import (
	"bot-service/internal/repository"
	"fmt"
	"log"
)

// SearchUsecase defines the interface for the search use case.
type SearchUsecase interface {
	Search(query string, page int, limit int, filter string) ([]byte, error)
	DeleteDocument(docID string) error
}

// searchUsecaseImpl implements the SearchUsecase interface.
type searchUsecaseImpl struct {
	searchRepo repository.SearchRepository
}

// NewSearchUsecase creates a new SearchUsecase.
func NewSearchUsecase(searchRepo repository.SearchRepository) SearchUsecase {
	if searchRepo == nil {
		log.Printf("ERROR: search repository is nil")
		panic("search repository cannot be nil")
	}
	return &searchUsecaseImpl{
		searchRepo: searchRepo,
	}
}

// Search performs a search using the search repository.
func (s *searchUsecaseImpl) Search(query string, page int, limit int, filter string) ([]byte, error) {
	log.Printf("INFO: Performing search: query='%s', page=%d, limit=%d, filter='%s'", query, page, limit, filter)

	if query == "" {
		log.Printf("ERROR: empty search query")
		return nil, fmt.Errorf("search query cannot be empty")
	}

	if page < 1 {
		log.Printf("ERROR: invalid page number: %d", page)
		return nil, fmt.Errorf("page number must be greater than 0")
	}

	if limit < 1 {
		log.Printf("ERROR: invalid limit: %d", limit)
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	result, err := s.searchRepo.Search(query, page, limit, filter)
	if err != nil {
		log.Printf("ERROR: failed to search: %v", err)
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	return result, nil
}

// DeleteDocument deletes a document from the search repository.
func (s *searchUsecaseImpl) DeleteDocument(docID string) error {
	log.Printf("INFO: Deleting document: docID='%s'", docID)
	if docID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}
	return s.searchRepo.DeleteDocument(docID)
}
