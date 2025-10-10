package manager

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/parkernilson/micro-journal/internal/domain"
)

// mockJournalStore is a mock implementation of JournalStore for testing.
type mockJournalStore struct {
	createFunc  func(ctx context.Context, title, content string) (*domain.JournalEntry, error)
	getByIDFunc func(ctx context.Context, id int64) (*domain.JournalEntry, error)
	updateFunc  func(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error)
	deleteFunc  func(ctx context.Context, id int64) error
	listFunc    func(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error)
}

func (m *mockJournalStore) Create(ctx context.Context, title, content string) (*domain.JournalEntry, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, title, content)
	}
	return nil, errors.New("not implemented")
}

func (m *mockJournalStore) GetByID(ctx context.Context, id int64) (*domain.JournalEntry, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockJournalStore) Update(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, title, content)
	}
	return nil, errors.New("not implemented")
}

func (m *mockJournalStore) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockJournalStore) List(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return nil, 0, errors.New("not implemented")
}

func TestJournalManager_CreateEntry(t *testing.T) {
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		mockStore := &mockJournalStore{
			createFunc: func(ctx context.Context, title, content string) (*domain.JournalEntry, error) {
				return &domain.JournalEntry{
					ID:        1,
					Title:     title,
					Content:   content,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		manager := NewJournalManager(mockStore)
		entry, err := manager.CreateEntry(ctx, "Test Title", "Test Content")

		if err != nil {
			t.Fatalf("CreateEntry failed: %v", err)
		}
		if entry.Title != "Test Title" {
			t.Errorf("Expected title 'Test Title', got '%s'", entry.Title)
		}
	})

	t.Run("empty title", func(t *testing.T) {
		mockStore := &mockJournalStore{}
		manager := NewJournalManager(mockStore)
		_, err := manager.CreateEntry(ctx, "", "Test Content")

		if err == nil {
			t.Error("Expected error for empty title, got nil")
		}
	})

	t.Run("empty content", func(t *testing.T) {
		mockStore := &mockJournalStore{}
		manager := NewJournalManager(mockStore)
		_, err := manager.CreateEntry(ctx, "Test Title", "")

		if err == nil {
			t.Error("Expected error for empty content, got nil")
		}
	})
}

func TestJournalManager_GetEntry(t *testing.T) {
	ctx := context.Background()

	mockStore := &mockJournalStore{
		getByIDFunc: func(ctx context.Context, id int64) (*domain.JournalEntry, error) {
			if id == 1 {
				return &domain.JournalEntry{
					ID:        1,
					Title:     "Test Title",
					Content:   "Test Content",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			}
			return nil, errors.New("not found")
		},
	}

	manager := NewJournalManager(mockStore)

	t.Run("successful get", func(t *testing.T) {
		entry, err := manager.GetEntry(ctx, 1)
		if err != nil {
			t.Fatalf("GetEntry failed: %v", err)
		}
		if entry.ID != 1 {
			t.Errorf("Expected ID 1, got %d", entry.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := manager.GetEntry(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existent entry, got nil")
		}
	})
}

func TestJournalManager_UpdateEntry(t *testing.T) {
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		mockStore := &mockJournalStore{
			updateFunc: func(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error) {
				return &domain.JournalEntry{
					ID:        id,
					Title:     title,
					Content:   content,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		manager := NewJournalManager(mockStore)
		entry, err := manager.UpdateEntry(ctx, 1, "Updated Title", "Updated Content")

		if err != nil {
			t.Fatalf("UpdateEntry failed: %v", err)
		}
		if entry.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", entry.Title)
		}
	})

	t.Run("empty title", func(t *testing.T) {
		mockStore := &mockJournalStore{}
		manager := NewJournalManager(mockStore)
		_, err := manager.UpdateEntry(ctx, 1, "", "Test Content")

		if err == nil {
			t.Error("Expected error for empty title, got nil")
		}
	})

	t.Run("empty content", func(t *testing.T) {
		mockStore := &mockJournalStore{}
		manager := NewJournalManager(mockStore)
		_, err := manager.UpdateEntry(ctx, 1, "Test Title", "")

		if err == nil {
			t.Error("Expected error for empty content, got nil")
		}
	})
}

func TestJournalManager_DeleteEntry(t *testing.T) {
	ctx := context.Background()

	mockStore := &mockJournalStore{
		deleteFunc: func(ctx context.Context, id int64) error {
			if id == 1 {
				return nil
			}
			return errors.New("not found")
		},
	}

	manager := NewJournalManager(mockStore)

	t.Run("successful delete", func(t *testing.T) {
		err := manager.DeleteEntry(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteEntry failed: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := manager.DeleteEntry(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existent entry, got nil")
		}
	})
}

func TestJournalManager_ListEntries(t *testing.T) {
	ctx := context.Background()

	createMockEntries := func(count int) []*domain.JournalEntry {
		entries := make([]*domain.JournalEntry, count)
		for i := 0; i < count; i++ {
			entries[i] = &domain.JournalEntry{
				ID:        int64(i + 1),
				Title:     "Title " + string(rune('0'+i)),
				Content:   "Content " + string(rune('0'+i)),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
		return entries
	}

	t.Run("first page", func(t *testing.T) {
		mockStore := &mockJournalStore{
			listFunc: func(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error) {
				if limit == 10 && offset == 0 {
					return createMockEntries(10), 25, nil
				}
				return nil, 0, errors.New("unexpected parameters")
			},
		}

		manager := NewJournalManager(mockStore)
		result, err := manager.ListEntries(ctx, 10, "")

		if err != nil {
			t.Fatalf("ListEntries failed: %v", err)
		}
		if len(result.Entries) != 10 {
			t.Errorf("Expected 10 entries, got %d", len(result.Entries))
		}
		if result.TotalCount != 25 {
			t.Errorf("Expected total count 25, got %d", result.TotalCount)
		}
		if result.NextPageToken == "" {
			t.Error("Expected next page token, got empty string")
		}
	})

	t.Run("second page", func(t *testing.T) {
		mockStore := &mockJournalStore{
			listFunc: func(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error) {
				if limit == 10 && offset == 10 {
					return createMockEntries(10), 25, nil
				}
				return nil, 0, errors.New("unexpected parameters")
			},
		}

		manager := NewJournalManager(mockStore)
		// "10" encoded in base64 is "MTA="
		pageToken := "MTA="
		result, err := manager.ListEntries(ctx, 10, pageToken)

		if err != nil {
			t.Fatalf("ListEntries failed: %v", err)
		}
		if len(result.Entries) != 10 {
			t.Errorf("Expected 10 entries, got %d", len(result.Entries))
		}
	})

	t.Run("last page", func(t *testing.T) {
		mockStore := &mockJournalStore{
			listFunc: func(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error) {
				if limit == 10 && offset == 20 {
					return createMockEntries(5), 25, nil
				}
				return nil, 0, errors.New("unexpected parameters")
			},
		}

		manager := NewJournalManager(mockStore)
		// "20" encoded in base64 is "MjA="
		pageToken := "MjA="
		result, err := manager.ListEntries(ctx, 10, pageToken)

		if err != nil {
			t.Fatalf("ListEntries failed: %v", err)
		}
		if len(result.Entries) != 5 {
			t.Errorf("Expected 5 entries, got %d", len(result.Entries))
		}
		if result.NextPageToken != "" {
			t.Errorf("Expected empty next page token, got '%s'", result.NextPageToken)
		}
	})

	t.Run("default page size", func(t *testing.T) {
		mockStore := &mockJournalStore{
			listFunc: func(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error) {
				if limit == 10 && offset == 0 {
					return createMockEntries(10), 10, nil
				}
				return nil, 0, errors.New("unexpected parameters")
			},
		}

		manager := NewJournalManager(mockStore)
		_, err := manager.ListEntries(ctx, 0, "")

		if err != nil {
			t.Fatalf("ListEntries failed: %v", err)
		}
	})

	t.Run("max page size", func(t *testing.T) {
		mockStore := &mockJournalStore{
			listFunc: func(ctx context.Context, limit, offset int) ([]*domain.JournalEntry, int64, error) {
				if limit == 100 && offset == 0 {
					return createMockEntries(100), 200, nil
				}
				return nil, 0, errors.New("unexpected parameters")
			},
		}

		manager := NewJournalManager(mockStore)
		_, err := manager.ListEntries(ctx, 200, "")

		if err != nil {
			t.Fatalf("ListEntries failed: %v", err)
		}
	})

	t.Run("invalid page token", func(t *testing.T) {
		mockStore := &mockJournalStore{}
		manager := NewJournalManager(mockStore)
		_, err := manager.ListEntries(ctx, 10, "invalid-token")

		if err == nil {
			t.Error("Expected error for invalid page token, got nil")
		}
	})
}
