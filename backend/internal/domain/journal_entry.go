package domain

import "time"

// JournalEntry represents a journal entry in the domain model.
// This is the shared type used across all layers (store, manager, service).
type JournalEntry struct {
	ID        int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
