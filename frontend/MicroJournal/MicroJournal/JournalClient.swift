//
//  JournalClient.swift
//  MicroJournal
//
//  Protocol defining the interface for interacting with the Journal service
//

import Foundation

// MARK: - Models

/// Represents a journal entry
struct JournalEntry: Identifiable, Codable {
    let id: String
    let title: String
    let content: String
    let createdAt: Date
    let updatedAt: Date
}

/// Request to create a new journal entry
struct CreateJournalEntryRequest {
    let title: String
    let content: String
}

/// Request to update an existing journal entry
struct UpdateJournalEntryRequest {
    let id: String
    let title: String
    let content: String
}

/// Response from listing journal entries
struct ListJournalEntriesResponse {
    let entries: [JournalEntry]
    let nextPageToken: String?
    let totalCount: Int
}

// MARK: - JournalClient Protocol

/// Protocol defining the interface for journal operations
protocol JournalClient {
    /// Creates a new journal entry
    /// - Parameter request: The create request containing title and content
    /// - Returns: The created journal entry
    /// - Throws: Error if the operation fails
    func createEntry(_ request: CreateJournalEntryRequest) async throws -> JournalEntry

    /// Updates an existing journal entry
    /// - Parameter request: The update request containing id, title, and content
    /// - Returns: The updated journal entry
    /// - Throws: Error if the operation fails
    func updateEntry(_ request: UpdateJournalEntryRequest) async throws -> JournalEntry

    /// Deletes a journal entry
    /// - Parameter id: The ID of the entry to delete
    /// - Returns: True if deletion was successful
    /// - Throws: Error if the operation fails
    func deleteEntry(id: String) async throws -> Bool

    /// Lists journal entries with pagination
    /// - Parameters:
    ///   - pageSize: Number of entries to return per page
    ///   - pageToken: Token for fetching the next page (nil for first page)
    /// - Returns: Response containing entries, next page token, and total count
    /// - Throws: Error if the operation fails
    func listEntries(pageSize: Int, pageToken: String?) async throws -> ListJournalEntriesResponse
}

// MARK: - Error Types

enum JournalClientError: LocalizedError {
    case networkError(String)
    case invalidResponse
    case notFound
    case serverError(String)

    var errorDescription: String? {
        switch self {
        case .networkError(let message):
            return "Network error: \(message)"
        case .invalidResponse:
            return "Invalid response from server"
        case .notFound:
            return "Entry not found"
        case .serverError(let message):
            return "Server error: \(message)"
        }
    }
}
