package repository

// SearchRepository defines the interface for searching data.
type SearchRepository interface {
	Search(query string, page int, limit int, filter string) ([]byte, error)
	DeleteDocument(docID string) error
}
