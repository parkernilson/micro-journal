#!/bin/bash

# Start the Go gRPC backend server

cd "$(dirname "$0")/../backend"

go run ./cmd/server
