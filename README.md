# Micro Journal

A gRPC-based journaling application with a Go backend and Swift frontend.

## Project Structure

```
micro-journal/
├── proto/                      # Protobuf definitions
│   └── journal/
│       └── v1/
│           └── journal.proto
├── backend/                    # Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # Server entry point
│   ├── internal/
│   │   └── service/
│   │       └── journal_service.go  # Service implementation
│   └── gen/                   # Generated protobuf code
├── frontend/                   # Swift frontend (to be implemented)
├── scripts/
│   └── generate-proto.sh      # Protobuf code generation script
└── README.md
```

## Features

The journal service supports the following operations:

1. **Create Journal Entry** - Create a new journal entry with title and content
2. **Update Journal Entry** - Edit an existing journal entry
3. **Delete Journal Entry** - Remove a journal entry
4. **List Journal Entries** - Get paginated journal entries sorted by date (descending)

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
./scripts/generate-proto.sh
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

## API Reference

### CreateJournalEntry

Creates a new journal entry.

**Request:**
```protobuf
message CreateJournalEntryRequest {
  string title = 1;
  string content = 2;
}
```

**Response:**
```protobuf
message CreateJournalEntryResponse {
  JournalEntry entry = 1;
}
```

### UpdateJournalEntry

Updates an existing journal entry.

**Request:**
```protobuf
message UpdateJournalEntryRequest {
  string id = 1;
  string title = 2;
  string content = 3;
}
```

**Response:**
```protobuf
message UpdateJournalEntryResponse {
  JournalEntry entry = 1;
}
```

### DeleteJournalEntry

Deletes a journal entry.

**Request:**
```protobuf
message DeleteJournalEntryRequest {
  string id = 1;
}
```

**Response:**
```protobuf
message DeleteJournalEntryResponse {
  bool success = 1;
}
```

### ListJournalEntries

Returns paginated journal entries sorted by date descending.

**Request:**
```protobuf
message ListJournalEntriesRequest {
  int32 page_size = 1;
  string page_token = 2;
}
```

**Response:**
```protobuf
message ListJournalEntriesResponse {
  repeated JournalEntry entries = 1;
  string next_page_token = 2;
  int32 total_count = 3;
}
```

## Development

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

- Implement database persistence (currently using stub implementations)
- Add authentication and authorization
- Implement the Swift frontend application
- Add comprehensive error handling
- Add unit and integration tests

## License

MIT
