# JournalClient Dependency Injection System

This directory contains a complete dependency injection system for the JournalClient, allowing you to interact with the gRPC JournalService throughout your SwiftUI app.

## Architecture Overview

The system uses **SwiftUI's Environment** for dependency injection - a simple, native approach perfect for small apps. No external DI frameworks needed!

## Files Created

### 1. `JournalClient.swift`
- **Protocol** defining the interface for journal operations
- **Swift models** for journal entries and requests/responses
- **Error types** for handling failures
- Methods:
  - `createEntry(_:)` - Create a new journal entry
  - `updateEntry(_:)` - Update an existing entry
  - `deleteEntry(id:)` - Delete an entry
  - `listEntries(pageSize:pageToken:)` - List entries with pagination

### 2. `GRPCJournalClient.swift`
- **Production implementation** using gRPC
- Contains TODO comments showing where to integrate gRPC-Swift
- Currently throws errors until gRPC-Swift is set up
- Includes implementation notes and next steps

### 3. `MockJournalClient.swift`
- **Fully functional in-memory implementation**
- Perfect for development, testing, and SwiftUI previews
- Features:
  - Simulated network delay (configurable)
  - Error simulation for testing error handling
  - Sample data factory (`MockJournalClient.withSampleData`)
  - In-memory storage with proper pagination

### 4. `JournalClientEnvironment.swift`
- **SwiftUI Environment setup** for dependency injection
- Provides:
  - `@Environment(\.journalClient)` for accessing the client in views
  - `.journalClient(_:)` modifier for injecting the client
  - `.withMockJournalClient()` helper for previews

### 5. `MicroJournalApp.swift` (Updated)
- **Root injection** of the JournalClient
- Configuration enum to switch between development/production
- Currently set to `.production` (uses GRPCJournalClient)

### 6. `ExampleJournalListView.swift`
- **Example implementation** showing best practices
- Demonstrates:
  - How to access the injected client
  - Async/await usage with SwiftUI
  - Error handling
  - Loading states
  - Multiple preview configurations

## Usage

### In Your Views

```swift
struct MyJournalView: View {
    @Environment(\.journalClient) var journalClient
    @State private var entries: [JournalEntry] = []

    var body: some View {
        List(entries) { entry in
            Text(entry.title)
        }
        .task {
            do {
                let response = try await journalClient.listEntries(
                    pageSize: 20,
                    pageToken: nil
                )
                entries = response.entries
            } catch {
                print("Error: \(error)")
            }
        }
    }
}
```

### In SwiftUI Previews

```swift
#Preview {
    MyJournalView()
        .withMockJournalClient()
}
```

### Creating Entries

```swift
let request = CreateJournalEntryRequest(
    title: "My First Entry",
    content: "This is the content of my journal entry."
)

do {
    let entry = try await journalClient.createEntry(request)
    print("Created entry with ID: \(entry.id)")
} catch {
    print("Error creating entry: \(error)")
}
```

### Updating Entries

```swift
let request = UpdateJournalEntryRequest(
    id: "123",
    title: "Updated Title",
    content: "Updated content"
)

do {
    let entry = try await journalClient.updateEntry(request)
    print("Updated entry: \(entry.title)")
} catch {
    print("Error updating entry: \(error)")
}
```

### Deleting Entries

```swift
do {
    let success = try await journalClient.deleteEntry(id: "123")
    if success {
        print("Entry deleted successfully")
    }
} catch {
    print("Error deleting entry: \(error)")
}
```

## Switching Between Mock and Production

In `MicroJournalApp.swift`, change the configuration:

```swift
// For development with mock data
private let config: AppConfiguration = .development

// For production with real gRPC backend
private let config: AppConfiguration = .production
```

## Setting Up gRPC (For Production)

The gRPC implementation is now **fully functional** and ready to use!

### What's Already Set Up

1. ✅ **gRPC-Swift 2** dependency added to Xcode project
2. ✅ **SwiftProtobuf** dependency added
3. ✅ **Code generation script** (`script/generate-proto.sh`) configured
4. ✅ **Generated Swift files** created in `MicroJournal/Generated/journal/v1/`
5. ✅ **GRPCJournalClient** fully implemented with:
   - HTTP/2 transport layer
   - All four RPC methods (create, update, delete, list)
   - Proto ↔ Swift model conversion
   - Comprehensive error handling

### Requirements

- **iOS 18.0+** or **macOS 15.0+** (required by gRPC Swift 2)
- For older iOS versions, the app will automatically use `MockJournalClient`

### Regenerating Proto Code

If you modify the proto files, regenerate the Swift code:

```bash
./script/generate-proto.sh
```

This generates both Go and Swift code from your proto definitions.

### Using the gRPC Client

The client is configured through `MicroJournalApp.swift`. Simply change the configuration:

```swift
// Switch from development (mock) to production (gRPC)
private let config: AppConfiguration = .production
```

You can also customize the server connection:

```swift
return createGRPCJournalClient(host: "192.168.1.100", port: 8080)
```

## Testing

The MockJournalClient is perfect for testing:

```swift
func testCreateEntry() async throws {
    let client = MockJournalClient()

    let request = CreateJournalEntryRequest(
        title: "Test",
        content: "Test content"
    )

    let entry = try await client.createEntry(request)

    XCTAssertEqual(entry.title, "Test")
    XCTAssertEqual(client.entryCount, 1)
}
```

## Benefits of This Approach

✅ **Native SwiftUI** - Uses built-in Environment system
✅ **Simple** - No external DI frameworks needed
✅ **Testable** - Easy to swap implementations
✅ **Type-safe** - Compiler-checked protocol conformance
✅ **Async/await** - Modern Swift concurrency
✅ **Preview-friendly** - Mock client works great in SwiftUI previews

## Next Steps

1. Build your UI - access the client via `@Environment(\.journalClient)`
2. Test with your gRPC backend (currently configured to use production mode)
3. Switch to `.development` mode in `MicroJournalApp.swift` for testing with mock data

## Examples

See `ExampleJournalListView.swift` for a complete working example that demonstrates:
- Loading entries
- Creating entries
- Deleting entries
- Error handling
- Loading states
- Empty states
- Multiple preview configurations

## Questions?

The system is designed to be simple and intuitive. If you need to:
- **Add a new method** → Update the `JournalClient` protocol and all implementations
- **Change the backend URL** → Modify the `GRPCJournalClient` initializer
- **Add test data** → Use `MockJournalClient.withSampleData` or create your own
- **Mock errors** → Set `client.shouldSimulateError = true` on MockJournalClient
