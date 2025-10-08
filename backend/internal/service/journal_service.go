package service

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/parkernilson/micro-journal/gen/journal/v1"
)

// JournalService implements the JournalServiceServer interface
type JournalService struct {
	pb.UnimplementedJournalServiceServer
}

// NewJournalService creates a new instance of JournalService
func NewJournalService() *JournalService {
	return &JournalService{}
}

// CreateJournalEntry creates a new journal entry
func (s *JournalService) CreateJournalEntry(ctx context.Context, req *pb.CreateJournalEntryRequest) (*pb.CreateJournalEntryResponse, error) {
	log.Printf("CreateJournalEntry called with title: %s", req.Title)

	// Stub implementation - return a mock entry
	entry := &pb.JournalEntry{
		Id:        "mock-entry-id-123",
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	return &pb.CreateJournalEntryResponse{
		Entry: entry,
	}, nil
}

// UpdateJournalEntry updates an existing journal entry
func (s *JournalService) UpdateJournalEntry(ctx context.Context, req *pb.UpdateJournalEntryRequest) (*pb.UpdateJournalEntryResponse, error) {
	log.Printf("UpdateJournalEntry called for entry ID: %s", req.Id)

	// Stub implementation - return a mock updated entry
	entry := &pb.JournalEntry{
		Id:        req.Id,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	return &pb.UpdateJournalEntryResponse{
		Entry: entry,
	}, nil
}

// DeleteJournalEntry deletes a journal entry
func (s *JournalService) DeleteJournalEntry(ctx context.Context, req *pb.DeleteJournalEntryRequest) (*pb.DeleteJournalEntryResponse, error) {
	log.Printf("DeleteJournalEntry called for entry ID: %s", req.Id)

	// Stub implementation - return success
	return &pb.DeleteJournalEntryResponse{
		Success: true,
	}, nil
}

// ListJournalEntries returns paginated journal entries sorted by date descending
func (s *JournalService) ListJournalEntries(ctx context.Context, req *pb.ListJournalEntriesRequest) (*pb.ListJournalEntriesResponse, error) {
	log.Printf("ListJournalEntries called with page_size: %d, page_token: %s", req.PageSize, req.PageToken)

	// Stub implementation - return mock entries
	entries := []*pb.JournalEntry{
		{
			Id:        "entry-1",
			Title:     "Mock Entry 1",
			Content:   "This is a mock journal entry",
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
		{
			Id:        "entry-2",
			Title:     "Mock Entry 2",
			Content:   "This is another mock journal entry",
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
	}

	return &pb.ListJournalEntriesResponse{
		Entries:       entries,
		NextPageToken: "",
		TotalCount:    2,
	}, nil
}
