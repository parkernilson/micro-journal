//
//  GRPCJournalClient.swift
//  MicroJournal
//
//  gRPC implementation of JournalClient using generated protobuf code
//

import Foundation
import GRPCCore
import GRPCNIOTransportHTTP2

/// Concrete implementation of JournalClient using gRPC Swift 2
@available(iOS 18.0, macOS 15.0, *)
class GRPCJournalClient: JournalClient {
    private let grpcClient: Journal_V1_JournalService.Client<HTTP2ClientTransport.Posix>
    private let host: String
    private let port: Int

    /// Initializes the gRPC client with server details
    /// - Parameters:
    ///   - host: The server host (e.g., "localhost" or "192.168.1.100")
    ///   - port: The server port (e.g., 50051)
    init(host: String = "localhost", port: Int = 50051) {
        self.host = host
        self.port = port

        // Create the gRPC client with HTTP/2 transport
        let transport = try! HTTP2ClientTransport.Posix(
            target: .ipv4(host: host, port: port),
            config: .defaults(transportSecurity: .plaintext)
        )

        let client = GRPCClient(transport: transport)
        self.grpcClient = Journal_V1_JournalService.Client(wrapping: client)
    }

    // MARK: - JournalClient Implementation

    func createEntry(_ request: CreateJournalEntryRequest) async throws -> JournalEntry {
        // Convert Swift request to protobuf request
        var protoRequest = Journal_V1_CreateJournalEntryRequest()
        protoRequest.title = request.title
        protoRequest.content = request.content

        do {
            // Make the gRPC call
            let response: Journal_V1_CreateJournalEntryResponse = try await grpcClient.createJournalEntry(
                protoRequest
            )

            // Convert protobuf response to Swift model
            guard response.hasEntry else {
                throw JournalClientError.invalidResponse
            }

            return try convertToSwiftModel(response.entry)

        } catch {
            throw mapError(error)
        }
    }

    func updateEntry(_ request: UpdateJournalEntryRequest) async throws -> JournalEntry {
        // Convert Swift request to protobuf request
        var protoRequest = Journal_V1_UpdateJournalEntryRequest()
        protoRequest.id = request.id
        protoRequest.title = request.title
        protoRequest.content = request.content

        do {
            // Make the gRPC call
            let response: Journal_V1_UpdateJournalEntryResponse = try await grpcClient.updateJournalEntry(
                protoRequest
            )

            // Convert protobuf response to Swift model
            guard response.hasEntry else {
                throw JournalClientError.invalidResponse
            }

            return try convertToSwiftModel(response.entry)

        } catch {
            throw mapError(error)
        }
    }

    func deleteEntry(id: String) async throws -> Bool {
        // Convert Swift request to protobuf request
        var protoRequest = Journal_V1_DeleteJournalEntryRequest()
        protoRequest.id = id

        do {
            // Make the gRPC call
            let response: Journal_V1_DeleteJournalEntryResponse = try await grpcClient.deleteJournalEntry(
                protoRequest
            )

            return response.success

        } catch {
            throw mapError(error)
        }
    }

    func listEntries(pageSize: Int, pageToken: String?) async throws -> ListJournalEntriesResponse {
        // Convert Swift request to protobuf request
        var protoRequest = Journal_V1_ListJournalEntriesRequest()
        protoRequest.pageSize = Int32(pageSize)
        if let token = pageToken, !token.isEmpty {
            protoRequest.pageToken = token
        }

        do {
            // Make the gRPC call
            let response: Journal_V1_ListJournalEntriesResponse = try await grpcClient.listJournalEntries(
                protoRequest
            )

            // Convert protobuf response to Swift model
            let entries = try response.entries.map { try convertToSwiftModel($0) }
            let nextToken = response.nextPageToken.isEmpty ? nil : response.nextPageToken

            return ListJournalEntriesResponse(
                entries: entries,
                nextPageToken: nextToken,
                totalCount: Int(response.totalCount)
            )

        } catch {
            throw mapError(error)
        }
    }

    // MARK: - Helper Methods

    /// Converts a proto JournalEntry to Swift model
    /// - Parameter protoEntry: The proto entry
    /// - Returns: Swift JournalEntry model
    /// - Throws: JournalClientError if conversion fails
    private func convertToSwiftModel(_ protoEntry: Journal_V1_JournalEntry) throws -> JournalEntry {
        // Check that timestamps exist
        guard protoEntry.hasCreatedAt, protoEntry.hasUpdatedAt else {
            throw JournalClientError.invalidResponse
        }

        // Convert Google.Protobuf.Timestamp to Swift Date
        let createdAt = Date(
            timeIntervalSince1970: TimeInterval(protoEntry.createdAt.seconds) +
            TimeInterval(protoEntry.createdAt.nanos) / 1_000_000_000
        )

        let updatedAt = Date(
            timeIntervalSince1970: TimeInterval(protoEntry.updatedAt.seconds) +
            TimeInterval(protoEntry.updatedAt.nanos) / 1_000_000_000
        )

        return JournalEntry(
            id: protoEntry.id,
            title: protoEntry.title,
            content: protoEntry.content,
            createdAt: createdAt,
            updatedAt: updatedAt
        )
    }

    /// Maps gRPC errors to JournalClientError
    /// - Parameter error: The original error
    /// - Returns: Mapped JournalClientError
    private func mapError(_ error: Error) -> JournalClientError {
        // Check if it's already a JournalClientError
        if let journalError = error as? JournalClientError {
            return journalError
        }

        // Map gRPC errors
        if let rpcError = error as? RPCError {
            switch rpcError.code {
            case .notFound:
                return .notFound
            case .unavailable:
                return .networkError("Server unavailable")
            case .deadlineExceeded:
                return .networkError("Request timeout")
            case .unauthenticated:
                return .serverError("Authentication required")
            case .permissionDenied:
                return .serverError("Permission denied")
            case .invalidArgument:
                return .serverError("Invalid request")
            default:
                return .serverError("gRPC error: \(rpcError.code)")
            }
        }

        // Default to network error
        return .networkError(error.localizedDescription)
    }
}

// MARK: - Fallback for older iOS versions

/// Mock implementation for iOS versions < 18.0
/// Since gRPC Swift 2 requires iOS 18.0+, we provide a fallback that shows an error
class GRPCJournalClientFallback: JournalClient {
    private let requiredVersion = "iOS 18.0 / macOS 15.0"

    func createEntry(_ request: CreateJournalEntryRequest) async throws -> JournalEntry {
        throw JournalClientError.serverError(
            "gRPC client requires \(requiredVersion) or later. Please use MockJournalClient for testing on older versions."
        )
    }

    func updateEntry(_ request: UpdateJournalEntryRequest) async throws -> JournalEntry {
        throw JournalClientError.serverError(
            "gRPC client requires \(requiredVersion) or later. Please use MockJournalClient for testing on older versions."
        )
    }

    func deleteEntry(id: String) async throws -> Bool {
        throw JournalClientError.serverError(
            "gRPC client requires \(requiredVersion) or later. Please use MockJournalClient for testing on older versions."
        )
    }

    func listEntries(pageSize: Int, pageToken: String?) async throws -> ListJournalEntriesResponse {
        throw JournalClientError.serverError(
            "gRPC client requires \(requiredVersion) or later. Please use MockJournalClient for testing on older versions."
        )
    }
}

// MARK: - Factory Function

/// Creates the appropriate gRPC client based on iOS version
func createGRPCJournalClient(host: String = "localhost", port: Int = 50051) -> JournalClient {
    if #available(iOS 18.0, macOS 15.0, *) {
        return GRPCJournalClient(host: host, port: port)
    } else {
        return GRPCJournalClientFallback()
    }
}
