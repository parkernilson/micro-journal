package service

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/parkernilson/micro-journal/gen/journal/v1"
	"github.com/parkernilson/micro-journal/internal/domain"
	"github.com/parkernilson/micro-journal/internal/manager"
)

// JournalManager defines the interface for the manager layer.
type JournalManager interface {
	CreateEntry(ctx context.Context, title, content string) (*domain.JournalEntry, error)
	GetEntry(ctx context.Context, id int64) (*domain.JournalEntry, error)
	UpdateEntry(ctx context.Context, id int64, title, content string) (*domain.JournalEntry, error)
	DeleteEntry(ctx context.Context, id int64) error
	ListEntries(ctx context.Context, pageSize int32, pageToken string) (*manager.ListEntriesResult, error)
}

// JournalService implements the JournalServiceServer interface
type JournalService struct {
	pb.UnimplementedJournalServiceServer
	manager JournalManager
}

// NewJournalService creates a new instance of JournalService
func NewJournalService(manager JournalManager) *JournalService {
	return &JournalService{manager: manager}
}

// CreateJournalEntry creates a new journal entry
func (s *JournalService) CreateJournalEntry(ctx context.Context, req *pb.CreateJournalEntryRequest) (*pb.CreateJournalEntryResponse, error) {
	log.Printf("CreateJournalEntry called with title: %s", req.Title)

	entry, err := s.manager.CreateEntry(ctx, req.Title, req.Content)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create entry: %v", err)
	}

	return &pb.CreateJournalEntryResponse{
		Entry: domainToProto(entry),
	}, nil
}

// UpdateJournalEntry updates an existing journal entry
func (s *JournalService) UpdateJournalEntry(ctx context.Context, req *pb.UpdateJournalEntryRequest) (*pb.UpdateJournalEntryResponse, error) {
	log.Printf("UpdateJournalEntry called for entry ID: %s", req.Id)

	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid entry ID: %v", err)
	}

	entry, err := s.manager.UpdateEntry(ctx, id, req.Title, req.Content)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update entry: %v", err)
	}

	return &pb.UpdateJournalEntryResponse{
		Entry: domainToProto(entry),
	}, nil
}

// DeleteJournalEntry deletes a journal entry
func (s *JournalService) DeleteJournalEntry(ctx context.Context, req *pb.DeleteJournalEntryRequest) (*pb.DeleteJournalEntryResponse, error) {
	log.Printf("DeleteJournalEntry called for entry ID: %s", req.Id)

	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid entry ID: %v", err)
	}

	err = s.manager.DeleteEntry(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete entry: %v", err)
	}

	return &pb.DeleteJournalEntryResponse{
		Success: true,
	}, nil
}

// ListJournalEntries returns paginated journal entries sorted by date descending
func (s *JournalService) ListJournalEntries(ctx context.Context, req *pb.ListJournalEntriesRequest) (*pb.ListJournalEntriesResponse, error) {
	log.Printf("ListJournalEntries called with page_size: %d, page_token: %s", req.PageSize, req.PageToken)

	result, err := s.manager.ListEntries(ctx, req.PageSize, req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to list entries: %v", err)
	}

	// Convert domain entries to protobuf entries
	protoEntries := make([]*pb.JournalEntry, len(result.Entries))
	for i, entry := range result.Entries {
		protoEntries[i] = domainToProto(entry)
	}

	return &pb.ListJournalEntriesResponse{
		Entries:       protoEntries,
		NextPageToken: result.NextPageToken,
		TotalCount:    int32(result.TotalCount),
	}, nil
}

// domainToProto converts a domain JournalEntry to a protobuf JournalEntry
func domainToProto(entry *domain.JournalEntry) *pb.JournalEntry {
	return &pb.JournalEntry{
		Id:        fmt.Sprintf("%d", entry.ID),
		Title:     entry.Title,
		Content:   entry.Content,
		CreatedAt: timestamppb.New(entry.CreatedAt),
		UpdatedAt: timestamppb.New(entry.UpdatedAt),
	}
}
