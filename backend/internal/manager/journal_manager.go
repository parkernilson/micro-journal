package manager

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/parkernilson/micro-journal/internal/domain"
)

// JournalStore defines the interface for the store layer.
// This allows the manager to be tested with a mock store.
type JournalStore interface {
	Create(ctx context.Context, title, content string) (*domain.JournalEntry, error)
	GetByID(ctx context.Context, id int64) (*domain.JournalEntry, error)
	Update(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error)
}

// JournalManager handles business logic for journal entries.
type JournalManager struct {
	store JournalStore
}

// NewJournalManager creates a new instance of JournalManager.
func NewJournalManager(store JournalStore) *JournalManager {
	return &JournalManager{store: store}
}

// CreateEntry creates a new journal entry.
func (m *JournalManager) CreateEntry(ctx context.Context, title, content string) (*domain.JournalEntry, error) {
	// Add any business logic validation here
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	return m.store.Create(ctx, title, content)
}

// GetEntry retrieves a journal entry by ID.
func (m *JournalManager) GetEntry(ctx context.Context, id int64) (*domain.JournalEntry, error) {
	return m.store.GetByID(ctx, id)
}

// UpdateEntry updates an existing journal entry.
func (m *JournalManager) UpdateEntry(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error) {
	// Add any business logic validation here
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	return m.store.Update(ctx, id, title, content)
}

// DeleteEntry deletes a journal entry.
func (m *JournalManager) DeleteEntry(ctx context.Context, id int64) error {
	return m.store.Delete(ctx, id)
}

// ListEntriesResult contains the result of listing journal entries.
type ListEntriesResult struct {
	Entries       []*domain.JournalEntry
	NextPageToken string
	TotalCount    int64
}

// ListEntries retrieves journal entries with pagination.
// pageSize determines how many entries to return per page.
// pageToken is a base64-encoded offset for pagination (empty for first page).
func (m *JournalManager) ListEntries(ctx context.Context, pageSize int32, pageToken string) (*ListEntriesResult, error) {
	// Default page size
	if pageSize <= 0 {
		pageSize = 10
	}

	// Maximum page size
	if pageSize > 100 {
		pageSize = 100
	}

	// Decode page token to get offset
	offset := 0
	if pageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(pageToken)
		if err != nil {
			return nil, fmt.Errorf("invalid page token: %w", err)
		}
		offset, err = strconv.Atoi(string(decoded))
		if err != nil {
			return nil, fmt.Errorf("invalid page token: %w", err)
		}
	}

	// Get entries from store
	entries, totalCount, err := m.store.List(ctx, int(pageSize), offset)
	if err != nil {
		return nil, err
	}

	// Calculate next page token
	nextPageToken := ""
	nextOffset := offset + len(entries)
	if nextOffset < int(totalCount) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(nextOffset)))
	}

	return &ListEntriesResult{
		Entries:       entries,
		NextPageToken: nextPageToken,
		TotalCount:    totalCount,
	}, nil
}
