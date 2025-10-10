package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/parkernilson/micro-journal/internal/domain"
)

// JournalStore handles data access operations for journal entries.
type JournalStore struct {
	db *sql.DB
}

// NewJournalStore creates a new instance of JournalStore.
func NewJournalStore(db *sql.DB) *JournalStore {
	return &JournalStore{db: db}
}

// Create inserts a new journal entry into the database.
func (s *JournalStore) Create(ctx context.Context, title, content string) (*domain.JournalEntry, error) {
	query := `
		INSERT INTO journal_entries (title, content, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := s.db.ExecContext(ctx, query, title, content)
	if err != nil {
		return nil, fmt.Errorf("failed to insert journal entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Fetch the created entry to get accurate timestamps
	return s.GetByID(ctx, id)
}

// GetByID retrieves a journal entry by its ID.
func (s *JournalStore) GetByID(ctx context.Context, id int64) (*domain.JournalEntry, error) {
	query := `
		SELECT id, title, content, created_at, updated_at
		FROM journal_entries
		WHERE id = ?
	`

	entry := &domain.JournalEntry{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID,
		&entry.Title,
		&entry.Content,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("journal entry not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entry: %w", err)
	}

	return entry, nil
}

// Update modifies an existing journal entry.
func (s *JournalStore) Update(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error) {
	query := `
		UPDATE journal_entries
		SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query, title, content, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update journal entry: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil, fmt.Errorf("journal entry not found: %d", id)
	}

	return s.GetByID(ctx, id)
}

// Delete removes a journal entry from the database.
func (s *JournalStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM journal_entries WHERE id = ?`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete journal entry: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("journal entry not found: %d", id)
	}

	return nil
}

// List retrieves journal entries with pagination.
// Returns the entries and the total count of all entries.
func (s *JournalStore) List(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error) {
	// Get total count
	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM journal_entries`
	err := s.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count journal entries: %w", err)
	}

	// Get paginated entries
	query := `
		SELECT id, title, content, created_at, updated_at
		FROM journal_entries
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query journal entries: %w", err)
	}
	defer rows.Close()

	var entries []*domain.JournalEntry
	for rows.Next() {
		entry := &domain.JournalEntry{}
		err := rows.Scan(
			&entry.ID,
			&entry.Title,
			&entry.Content,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan journal entry: %w", err)
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return entries, totalCount, nil
}
