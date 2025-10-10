//
//  JournalClientEnvironment.swift
//  MicroJournal
//
//  SwiftUI Environment setup for JournalClient dependency injection
//

import SwiftUI

// MARK: - Environment Key

/// Environment key for injecting JournalClient into SwiftUI views
private struct JournalClientKey: EnvironmentKey {
    static let defaultValue: JournalClient = MockJournalClient()
}

// MARK: - Environment Values Extension

extension EnvironmentValues {
    /// Access the JournalClient from the environment
    ///
    /// Use this in your views like:
    /// ```swift
    /// @Environment(\.journalClient) var journalClient
    /// ```
    var journalClient: JournalClient {
        get { self[JournalClientKey.self] }
        set { self[JournalClientKey.self] = newValue }
    }
}

// MARK: - View Extension

extension View {
    /// Injects a JournalClient into the environment
    ///
    /// Use this in your app's root view or any parent view:
    /// ```swift
    /// ContentView()
    ///     .journalClient(MockJournalClient.withSampleData)
    /// ```
    func journalClient(_ client: JournalClient) -> some View {
        environment(\.journalClient, client)
    }
}

// MARK: - Preview Helper

#if DEBUG
extension View {
    /// Injects a mock JournalClient with sample data for previews
    ///
    /// Use this in your previews:
    /// ```swift
    /// #Preview {
    ///     MyView()
    ///         .withMockJournalClient()
    /// }
    /// ```
    func withMockJournalClient() -> some View {
        self.journalClient(MockJournalClient.withSampleData)
    }
}
#endif

/*
 USAGE EXAMPLES:

 1. In your app entry point (MicroJournalApp.swift):
    ```swift
    @main
    struct MicroJournalApp: App {
        var body: some Scene {
            WindowGroup {
                ContentView()
                    .journalClient(MockJournalClient.withSampleData) // Use mock for development
                    // OR
                    .journalClient(GRPCJournalClient()) // Use real gRPC client
            }
        }
    }
    ```

 2. In any SwiftUI view:
    ```swift
    struct EntryListView: View {
        @Environment(\.journalClient) var journalClient
        @State private var entries: [JournalEntry] = []

        var body: some View {
            List(entries) { entry in
                Text(entry.title)
            }
            .task {
                do {
                    let response = try await journalClient.listEntries(pageSize: 20, pageToken: nil)
                    entries = response.entries
                } catch {
                    print("Error loading entries: \(error)")
                }
            }
        }
    }
    ```

 3. In SwiftUI previews:
    ```swift
    #Preview {
        EntryListView()
            .withMockJournalClient()
    }
    ```

 4. Switching between implementations:
    ```swift
    // Create an enum to manage which implementation to use
    enum AppConfiguration {
        case development
        case production

        var journalClient: JournalClient {
            switch self {
            case .development:
                return MockJournalClient.withSampleData
            case .production:
                return GRPCJournalClient()
            }
        }
    }

    @main
    struct MicroJournalApp: App {
        let config: AppConfiguration = .production

        var body: some Scene {
            WindowGroup {
                ContentView()
                    .journalClient(config.journalClient)
            }
        }
    }
    ```
 */
