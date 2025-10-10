# Micro Journal

A gRPC-based journaling application with a Go backend and Swift frontend.

## Prerequisites

- Go 1.25 or later
- Protocol Buffers compiler (`protoc`)
- gRPC and protobuf Go plugins:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

## Getting Started

### 1. Install Dependencies

```bash
cd backend
go mod tidy
```

### 2. Generate Protobuf Code

```bash
./script/generate-proto.sh
```

### 3. Run the Server

```bash
cd backend
go run cmd/server/main.go
```

The server will start on port 50051.

### 4. Test the Server

You can test the server using `grpcurl`:

```bash
# List available services
grpcurl -plaintext localhost:50051 list

# List methods for JournalService
grpcurl -plaintext localhost:50051 list journal.v1.JournalService

# Create a journal entry
grpcurl -plaintext -d '{"title": "My First Entry", "content": "Hello, World!"}' \
  localhost:50051 journal.v1.JournalService/CreateJournalEntry

# List journal entries
grpcurl -plaintext -d '{"page_size": 10}' \
  localhost:50051 journal.v1.JournalService/ListJournalEntries
```

## Development

### Running Tests

Run all tests from the backend directory:

```bash
cd backend
go test ./...
```

Run tests with verbose output:

```bash
cd backend
go test -v ./...
```

Run tests with coverage:

```bash
cd backend
go test -cover ./...
```

### Regenerating Protobuf Code

After modifying `.proto` files, regenerate the code:

```bash
./scripts/generate-proto.sh
```

### Adding New Service Methods

1. Update `proto/journal/v1/journal.proto`
2. Run `./scripts/generate-proto.sh`
3. Implement the new methods in `backend/internal/service/journal_service.go`

## Next Steps

- Add authentication and authorization
- Implement the Swift frontend application
- Add comprehensive error handling
- Add unit and integration tests

## License

MIT
