//
//  MockJournalClient.swift
//  MicroJournal
//
//  Mock implementation of JournalClient for testing and development
//

import Foundation

/// Mock implementation of JournalClient with in-memory storage
/// Perfect for development, testing, and SwiftUI previews
class MockJournalClient: JournalClient {
    private var entries: [String: JournalEntry] = [:]
    private var nextId: Int = 1

    /// Simulated network delay (in seconds) for more realistic testing
    var simulatedDelay: TimeInterval = 0.5

    /// Flag to simulate errors for testing error handling
    var shouldSimulateError: Bool = false

    init(preloadedEntries: [JournalEntry] = []) {
        for entry in preloadedEntries {
            entries[entry.id] = entry
            if let idInt = Int(entry.id), idInt >= nextId {
                nextId = idInt + 1
            }
        }
    }

    func createEntry(_ request: CreateJournalEntryRequest) async throws -> JournalEntry {
        try await simulateDelay()

        if shouldSimulateError {
            throw JournalClientError.serverError("Simulated error")
        }

        let now = Date()
        let id = String(nextId)
        nextId += 1

        let entry = JournalEntry(
            id: id,
            title: request.title,
            content: request.content,
            createdAt: now,
            updatedAt: now
        )

        entries[id] = entry
        return entry
    }

    func updateEntry(_ request: UpdateJournalEntryRequest) async throws -> JournalEntry {
        try await simulateDelay()

        if shouldSimulateError {
            throw JournalClientError.serverError("Simulated error")
        }

        guard let existingEntry = entries[request.id] else {
            throw JournalClientError.notFound
        }

        let updatedEntry = JournalEntry(
            id: existingEntry.id,
            title: request.title,
            content: request.content,
            createdAt: existingEntry.createdAt,
            updatedAt: Date()
        )

        entries[request.id] = updatedEntry
        return updatedEntry
    }

    func deleteEntry(id: String) async throws -> Bool {
        try await simulateDelay()

        if shouldSimulateError {
            throw JournalClientError.serverError("Simulated error")
        }

        guard entries[id] != nil else {
            throw JournalClientError.notFound
        }

        entries.removeValue(forKey: id)
        return true
    }

    func listEntries(pageSize: Int, pageToken: String?) async throws -> ListJournalEntriesResponse {
        try await simulateDelay()

        if shouldSimulateError {
            throw JournalClientError.serverError("Simulated error")
        }

        // Sort entries by creation date descending (newest first)
        let sortedEntries = entries.values.sorted { $0.createdAt > $1.createdAt }

        // Calculate pagination
        let startIndex = pageToken.flatMap(Int.init) ?? 0
        let endIndex = min(startIndex + pageSize, sortedEntries.count)
        let pageEntries = Array(sortedEntries[startIndex..<endIndex])

        let nextPageToken: String? = endIndex < sortedEntries.count ? String(endIndex) : nil

        return ListJournalEntriesResponse(
            entries: pageEntries,
            nextPageToken: nextPageToken,
            totalCount: sortedEntries.count
        )
    }

    // MARK: - Helper Methods

    /// Simulates network delay for more realistic testing
    private func simulateDelay() async throws {
        if simulatedDelay > 0 {
            try await Task.sleep(nanoseconds: UInt64(simulatedDelay * 1_000_000_000))
        }
    }

    /// Clears all entries (useful for testing)
    func clearAllEntries() {
        entries.removeAll()
        nextId = 1
    }

    /// Returns the count of entries (useful for testing)
    var entryCount: Int {
        entries.count
    }
}

// MARK: - Preview Helper

extension MockJournalClient {
    /// Creates a mock client with sample data for SwiftUI previews
    static var withSampleData: MockJournalClient {
        let sampleEntries = [
            JournalEntry(
                id: "1",
                title: "First Entry",
                content: "This is my first journal entry. I'm excited to start journaling!",
                createdAt: Date().addingTimeInterval(-86400 * 7), // 7 days ago
                updatedAt: Date().addingTimeInterval(-86400 * 7)
            ),
            JournalEntry(
                id: "2",
                title: "Great Day",
                content: "Had an amazing day today. The weather was perfect and I accomplished a lot.",
                createdAt: Date().addingTimeInterval(-86400 * 3), // 3 days ago
                updatedAt: Date().addingTimeInterval(-86400 * 3)
            ),
            JournalEntry(
                id: "3",
                title: "Reflection",
                content: "Taking time to reflect on my progress. I'm learning a lot about myself through this journal.",
                createdAt: Date().addingTimeInterval(-86400), // Yesterday
                updatedAt: Date().addingTimeInterval(-86400)
            ),
            JournalEntry(
                id: "4",
                title: "Today's Thoughts",
                content: "Just writing down my thoughts for today. It's been a productive day.",
                createdAt: Date(), // Today
                updatedAt: Date()
            )
        ]

        let client = MockJournalClient(preloadedEntries: sampleEntries)
        client.simulatedDelay = 0.3 // Shorter delay for previews
        return client
    }
}
