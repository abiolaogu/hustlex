import SwiftUI

@main
struct HustleXApp: App {
    @StateObject private var appState = AppState()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appState)
        }
    }
}

class AppState: ObservableObject {
    @Published var isAuthenticated = false
    @Published var user: User?
}

struct User: Identifiable, Codable {
    let id: String
    let email: String?
    let phone: String?
    let firstName: String?
    let lastName: String?
}
