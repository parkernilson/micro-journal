package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/parkernilson/micro-journal/gen/journal/v1"
	"github.com/parkernilson/micro-journal/internal/service"
)

const (
	port = ":50051"
)

func main() {
	// Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the JournalService
	journalService := service.NewJournalService()
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
