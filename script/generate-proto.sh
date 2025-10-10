#!/bin/bash

# Script to generate Go and Swift code from protobuf definitions

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Generating protobuf code...${NC}"

# Navigate to project root
cd "$(dirname "$0")/.."

# Create output directories if they don't exist
mkdir -p backend/gen/proto/journal/v1
mkdir -p frontend/MicroJournal/MicroJournal/Generated

# Generate Go code from proto files
echo -e "${YELLOW}Generating Go code...${NC}"
protoc \
  --go_out=backend/gen \
  --go_opt=paths=source_relative \
  --go-grpc_out=backend/gen \
  --go-grpc_opt=paths=source_relative \
  --proto_path=proto \
  proto/journal/v1/journal.proto

# Generate Swift code from proto files
echo -e "${YELLOW}Generating Swift code...${NC}"
protoc \
  --swift_out=frontend/MicroJournal/MicroJournal/Generated \
  --grpc-swift_out=frontend/MicroJournal/MicroJournal/Generated \
  --proto_path=proto \
  proto/journal/v1/journal.proto

echo -e "${GREEN}Protobuf code generated successfully!${NC}"
echo -e "Generated files:"
echo -e "  Go:"
echo -e "    - backend/gen/proto/journal/v1/journal.pb.go"
echo -e "    - backend/gen/proto/journal/v1/journal_grpc.pb.go"
echo -e "  Swift:"
echo -e "    - frontend/MicroJournal/MicroJournal/Generated/journal.pb.swift"
echo -e "    - frontend/MicroJournal/MicroJournal/Generated/journal.grpc.swift"
