//
//  ExampleJournalListView.swift
//  MicroJournal
//
//  Example view demonstrating how to use the JournalClient with dependency injection
//

import SwiftUI

struct ExampleJournalListView: View {
    // MARK: - Dependencies

    /// Access the injected JournalClient from the environment
    @Environment(\.journalClient) var journalClient

    // MARK: - State

    @State private var entries: [JournalEntry] = []
    @State private var isLoading = false
    @State private var errorMessage: String?

    // MARK: - Body

    var body: some View {
        NavigationStack {
            Group {
                if isLoading {
                    ProgressView("Loading entries...")
                } else if let error = errorMessage {
                    VStack(spacing: 16) {
                        Image(systemName: "exclamationmark.triangle")
                            .font(.largeTitle)
                            .foregroundStyle(.red)
                        Text("Error")
                            .font(.headline)
                        Text(error)
                            .font(.body)
                            .foregroundStyle(.secondary)
                        Button("Retry") {
                            Task { await loadEntries() }
                        }
                        .buttonStyle(.borderedProminent)
                    }
                    .padding()
                } else if entries.isEmpty {
                    ContentUnavailableView {
                        Label("No Entries", systemImage: "book.closed")
                    } description: {
                        Text("Create your first journal entry to get started")
                    } actions: {
                        Button("Create Entry") {
                            Task { await createSampleEntry() }
                        }
                        .buttonStyle(.borderedProminent)
                    }
                } else {
                    List {
                        ForEach(entries) { entry in
                            VStack(alignment: .leading, spacing: 8) {
                                Text(entry.title)
                                    .font(.headline)
                                Text(entry.content)
                                    .font(.body)
                                    .foregroundStyle(.secondary)
                                    .lineLimit(2)
                                Text(entry.createdAt, style: .date)
                                    .font(.caption)
                                    .foregroundStyle(.tertiary)
                            }
                            .padding(.vertical, 4)
                        }
                        .onDelete(perform: deleteEntries)
                    }
                }
            }
            .navigationTitle("Journal Entries")
            .toolbar {
                ToolbarItem(placement: .primaryAction) {
                    Button {
                        Task { await createSampleEntry() }
                    } label: {
                        Image(systemName: "plus")
                    }
                }
                ToolbarItem(placement: .topBarLeading) {
                    Button {
                        Task { await loadEntries() }
                    } label: {
                        Image(systemName: "arrow.clockwise")
                    }
                    .disabled(isLoading)
                }
            }
            .task {
                await loadEntries()
            }
        }
    }

    // MARK: - Methods

    /// Loads journal entries using the injected client
    private func loadEntries() async {
        isLoading = true
        errorMessage = nil

        do {
            let response = try await journalClient.listEntries(pageSize: 50, pageToken: nil)
            entries = response.entries
        } catch {
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    /// Creates a sample entry
    private func createSampleEntry() async {
        let titles = ["Morning Thoughts", "Daily Reflection", "Quick Note", "Today's Goals"]
        let contents = [
            "Feeling energized and ready for the day!",
            "Reflecting on the progress I've made this week.",
            "Just a quick note to capture this moment.",
            "Setting my intentions for a productive day."
        ]

        do {
            let request = CreateJournalEntryRequest(
                title: titles.randomElement()!,
                content: contents.randomElement()!
            )
            let newEntry = try await journalClient.createEntry(request)
            entries.insert(newEntry, at: 0) // Add to beginning (newest first)
        } catch {
            errorMessage = error.localizedDescription
        }
    }

    /// Deletes entries at the specified offsets
    private func deleteEntries(at offsets: IndexSet) {
        for index in offsets {
            let entry = entries[index]
            Task {
                do {
                    _ = try await journalClient.deleteEntry(id: entry.id)
                    await MainActor.run {
                        entries.remove(at: index)
                    }
                } catch {
                    await MainActor.run {
                        errorMessage = error.localizedDescription
                    }
                }
            }
        }
    }
}

// MARK: - Previews

#Preview("With Sample Data") {
    ExampleJournalListView()
        .withMockJournalClient()
}

#Preview("Empty State") {
    ExampleJournalListView()
        .journalClient(MockJournalClient())
}

#Preview("Loading State") {
    ExampleJournalListView()
        .journalClient({
            let client = MockJournalClient()
            client.simulatedDelay = 10 // Long delay to show loading
            return client
        }())
}
