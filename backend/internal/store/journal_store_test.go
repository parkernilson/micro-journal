package store

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory SQLite database with the schema initialized.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	// Create the schema
	schema := `
		CREATE TABLE journal_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX idx_journal_entries_created_at ON journal_entries(created_at DESC);
	`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestJournalStore_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	entry, err := store.Create(ctx, "Test Title", "Test Content")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if entry.ID == 0 {
		t.Error("Expected non-zero ID")
	}
	if entry.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", entry.Title)
	}
	if entry.Content != "Test Content" {
		t.Errorf("Expected content 'Test Content', got '%s'", entry.Content)
	}
	if entry.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt")
	}
	if entry.UpdatedAt.IsZero() {
		t.Error("Expected non-zero UpdatedAt")
	}
}

func TestJournalStore_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	// Create an entry first
	created, err := store.Create(ctx, "Test Title", "Test Content")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Retrieve it
	retrieved, err := store.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
	}
	if retrieved.Title != created.Title {
		t.Errorf("Expected title '%s', got '%s'", created.Title, retrieved.Title)
	}
	if retrieved.Content != created.Content {
		t.Errorf("Expected content '%s', got '%s'", created.Content, retrieved.Content)
	}
}

func TestJournalStore_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	_, err := store.GetByID(ctx, 999)
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
}

func TestJournalStore_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	// Create an entry first
	created, err := store.Create(ctx, "Original Title", "Original Content")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update it
	updated, err := store.Update(ctx, created.ID, "Updated Title", "Updated Content")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, updated.ID)
	}
	if updated.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", updated.Title)
	}
	if updated.Content != "Updated Content" {
		t.Errorf("Expected content 'Updated Content', got '%s'", updated.Content)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) && !updated.UpdatedAt.Equal(created.UpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestJournalStore_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	_, err := store.Update(ctx, 999, "Title", "Content")
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
}

func TestJournalStore_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	// Create an entry first
	created, err := store.Create(ctx, "Test Title", "Test Content")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete it
	err = store.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it's gone
	_, err = store.GetByID(ctx, created.ID)
	if err == nil {
		t.Error("Expected error when retrieving deleted entry, got nil")
	}
}

func TestJournalStore_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	err := store.Delete(ctx, 999)
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
}

func TestJournalStore_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	// Create multiple entries
	for i := 1; i <= 5; i++ {
		_, err := store.Create(ctx, "Title "+string(rune('0'+i)), "Content "+string(rune('0'+i)))
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Test pagination - first page
	entries, total, err := store.List(ctx, 2, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total count 5, got %d", total)
	}
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// Test pagination - second page
	entries, total, err = store.List(ctx, 2, 2)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total count 5, got %d", total)
	}
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// Test pagination - last page
	entries, total, err = store.List(ctx, 2, 4)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total count 5, got %d", total)
	}
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
}

func TestJournalStore_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	entries, total, err := store.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if total != 0 {
		t.Errorf("Expected total count 0, got %d", total)
	}
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}

func TestJournalStore_List_OrderedByCreatedAtDesc(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewJournalStore(db)
	ctx := context.Background()

	// Create entries in order
	entry1, _ := store.Create(ctx, "First", "Content 1")
	entry2, _ := store.Create(ctx, "Second", "Content 2")
	entry3, _ := store.Create(ctx, "Third", "Content 3")

	// List should return in reverse order (newest first)
	entries, _, err := store.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Check order (newest to oldest)
	if entries[0].ID != entry3.ID {
		t.Errorf("Expected first entry to be entry3 (ID %d), got ID %d", entry3.ID, entries[0].ID)
	}
	if entries[1].ID != entry2.ID {
		t.Errorf("Expected second entry to be entry2 (ID %d), got ID %d", entry2.ID, entries[1].ID)
	}
	if entries[2].ID != entry1.ID {
		t.Errorf("Expected third entry to be entry1 (ID %d), got ID %d", entry1.ID, entries[2].ID)
	}
}
