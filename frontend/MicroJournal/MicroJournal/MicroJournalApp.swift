//
//  MicroJournalApp.swift
//  MicroJournal
//
//  Created by Parker Nilson on 10/10/25.
//

import SwiftUI

@main
struct MicroJournalApp: App {
    // MARK: - Configuration

    /// App configuration
    enum AppConfiguration {
        case development
        case production

        var journalClient: JournalClient {
            switch self {
            case .development:
                return MockJournalClient.withSampleData
            case .production:
                return createGRPCJournalClient(host: "localhost", port: 50051)
            }
        }
    }

    private let config: AppConfiguration = .production

    // MARK: - App Body

    var body: some Scene {
        WindowGroup {
            ContentView()
                .journalClient(config.journalClient)
        }
    }
}
