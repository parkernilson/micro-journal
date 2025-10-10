package service

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/parkernilson/micro-journal/gen/journal/v1"
	"github.com/parkernilson/micro-journal/internal/domain"
	"github.com/parkernilson/micro-journal/internal/manager"
)

// mockJournalManager is a mock implementation of JournalManager for testing.
type mockJournalManager struct {
	createEntryFunc func(ctx context.Context, title, content string) (*domain.JournalEntry, error)
	getEntryFunc    func(ctx context.Context, id int64) (*domain.JournalEntry, error)
	updateEntryFunc func(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error)
	deleteEntryFunc func(ctx context.Context, id int64) error
	listEntriesFunc func(ctx context.Context, pageSize int32, pageToken string) (*manager.ListEntriesResult, error)
}

func (m *mockJournalManager) CreateEntry(ctx context.Context, title, content string) (*domain.JournalEntry, error) {
	if m.createEntryFunc != nil {
		return m.createEntryFunc(ctx, title, content)
	}
	return nil, errors.New("not implemented")
}

func (m *mockJournalManager) GetEntry(ctx context.Context, id int64) (*domain.JournalEntry, error) {
	if m.getEntryFunc != nil {
		return m.getEntryFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockJournalManager) UpdateEntry(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error) {
	if m.updateEntryFunc != nil {
		return m.updateEntryFunc(ctx, id, title, content)
	}
	return nil, errors.New("not implemented")
}

func (m *mockJournalManager) DeleteEntry(ctx context.Context, id int64) error {
	if m.deleteEntryFunc != nil {
		return m.deleteEntryFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockJournalManager) ListEntries(ctx context.Context, pageSize int32, pageToken string) (*manager.ListEntriesResult, error) {
	if m.listEntriesFunc != nil {
		return m.listEntriesFunc(ctx, pageSize, pageToken)
	}
	return nil, errors.New("not implemented")
}

func TestJournalService_CreateJournalEntry(t *testing.T) {
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		mockManager := &mockJournalManager{
			createEntryFunc: func(ctx context.Context, title, content string) (*domain.JournalEntry, error) {
				return &domain.JournalEntry{
					ID:        1,
					Title:     title,
					Content:   content,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		service := NewJournalService(mockManager)
		req := &pb.CreateJournalEntryRequest{
			Title:   "Test Title",
			Content: "Test Content",
		}

		resp, err := service.CreateJournalEntry(ctx, req)
		if err != nil {
			t.Fatalf("CreateJournalEntry failed: %v", err)
		}
		if resp.Entry.Id != "1" {
			t.Errorf("Expected ID '1', got '%s'", resp.Entry.Id)
		}
		if resp.Entry.Title != "Test Title" {
			t.Errorf("Expected title 'Test Title', got '%s'", resp.Entry.Title)
		}
	})

	t.Run("manager error", func(t *testing.T) {
		mockManager := &mockJournalManager{
			createEntryFunc: func(ctx context.Context, title, content string) (*domain.JournalEntry, error) {
				return nil, errors.New("validation error")
			},
		}

		service := NewJournalService(mockManager)
		req := &pb.CreateJournalEntryRequest{
			Title:   "",
			Content: "Test Content",
		}

		_, err := service.CreateJournalEntry(ctx, req)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestJournalService_UpdateJournalEntry(t *testing.T) {
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		mockManager := &mockJournalManager{
			updateEntryFunc: func(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error) {
				return &domain.JournalEntry{
					ID:        id,
					Title:     title,
					Content:   content,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		service := NewJournalService(mockManager)
		req := &pb.UpdateJournalEntryRequest{
			Id:      "1",
			Title:   "Updated Title",
			Content: "Updated Content",
		}

		resp, err := service.UpdateJournalEntry(ctx, req)
		if err != nil {
			t.Fatalf("UpdateJournalEntry failed: %v", err)
		}
		if resp.Entry.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", resp.Entry.Title)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockManager := &mockJournalManager{}
		service := NewJournalService(mockManager)
		req := &pb.UpdateJournalEntryRequest{
			Id:      "invalid",
			Title:   "Title",
			Content: "Content",
		}

		_, err := service.UpdateJournalEntry(ctx, req)
		if err == nil {
			t.Error("Expected error for invalid ID, got nil")
		}
	})
}

func TestJournalService_DeleteJournalEntry(t *testing.T) {
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		mockManager := &mockJournalManager{
			deleteEntryFunc: func(ctx context.Context, id int64) error {
				return nil
			},
		}

		service := NewJournalService(mockManager)
		req := &pb.DeleteJournalEntryRequest{
			Id: "1",
		}

		resp, err := service.DeleteJournalEntry(ctx, req)
		if err != nil {
			t.Fatalf("DeleteJournalEntry failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockManager := &mockJournalManager{}
		service := NewJournalService(mockManager)
		req := &pb.DeleteJournalEntryRequest{
			Id: "invalid",
		}

		_, err := service.DeleteJournalEntry(ctx, req)
		if err == nil {
			t.Error("Expected error for invalid ID, got nil")
		}
	})
}

func TestJournalService_ListJournalEntries(t *testing.T) {
	ctx := context.Background()

	t.Run("successful list", func(t *testing.T) {
		mockManager := &mockJournalManager{
			listEntriesFunc: func(ctx context.Context, pageSize int32, pageToken string) (*manager.ListEntriesResult, error) {
				entries := []*domain.JournalEntry{
					{
						ID:        1,
						Title:     "Entry 1",
						Content:   "Content 1",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        2,
						Title:     "Entry 2",
						Content:   "Content 2",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				return &manager.ListEntriesResult{
					Entries:       entries,
					NextPageToken: "next-token",
					TotalCount:    10,
				}, nil
			},
		}

		service := NewJournalService(mockManager)
		req := &pb.ListJournalEntriesRequest{
			PageSize:  10,
			PageToken: "",
		}

		resp, err := service.ListJournalEntries(ctx, req)
		if err != nil {
			t.Fatalf("ListJournalEntries failed: %v", err)
		}
		if len(resp.Entries) != 2 {
			t.Errorf("Expected 2 entries, got %d", len(resp.Entries))
		}
		if resp.NextPageToken != "next-token" {
			t.Errorf("Expected next page token 'next-token', got '%s'", resp.NextPageToken)
		}
		if resp.TotalCount != 10 {
			t.Errorf("Expected total count 10, got %d", resp.TotalCount)
		}
	})

	t.Run("manager error", func(t *testing.T) {
		mockManager := &mockJournalManager{
			listEntriesFunc: func(ctx context.Context, pageSize int32, pageToken string) (*manager.ListEntriesResult, error) {
				return nil, errors.New("database error")
			},
		}

		service := NewJournalService(mockManager)
		req := &pb.ListJournalEntriesRequest{
			PageSize:  10,
			PageToken: "",
		}

		_, err := service.ListJournalEntries(ctx, req)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
