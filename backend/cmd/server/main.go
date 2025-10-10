package main

import (
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "modernc.org/sqlite"

	pb "github.com/parkernilson/micro-journal/gen/journal/v1"
	"github.com/parkernilson/micro-journal/internal/manager"
	"github.com/parkernilson/micro-journal/internal/service"
	"github.com/parkernilson/micro-journal/internal/store"
)

const (
	port   = ":50051"
	dbPath = "data/micro_journal.db"
)

func main() {
	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Configure connection pool (SQLite works best with limited connections)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Verify database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Printf("Connected to database at %s", dbPath)

	// Create layers: Store -> Manager -> Service
	journalStore := store.NewJournalStore(db)
	journalManager := manager.NewJournalManager(journalStore)
	journalService := service.NewJournalService(journalManager)

	// Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the JournalService
	pb.RegisterJournalServiceServer(grpcServer, journalService)

	// Register reflection service on gRPC server (useful for debugging with grpcurl)
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on port %s", port)
	log.Printf("Server is ready to accept connections")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
