import SwiftUI

struct ProfileView: View {
    @EnvironmentObject var appState: AppState

    var body: some View {
        NavigationView {
            List {
                Section {
                    HStack {
                        Image(systemName: "person.circle.fill")
                            .font(.system(size: 60))
                            .foregroundColor(Theme.Colors.primaryFallback)

                        VStack(alignment: .leading) {
                            Text("John Doe")
                                .font(Theme.Typography.titleMedium)
                            Text("+234 800 123 4567")
                                .font(Theme.Typography.bodySmall)
                                .foregroundColor(.secondary)
                        }
                    }
                    .padding(.vertical, Theme.Spacing.sm)
                }

                Section("Account") {
                    NavigationLink(destination: Text("Edit Profile")) {
                        Label("Edit Profile", systemImage: "person.fill")
                    }
                    NavigationLink(destination: Text("KYC")) {
                        Label("KYC Verification", systemImage: "checkmark.shield.fill")
                    }
                    NavigationLink(destination: Text("Security")) {
                        Label("Security", systemImage: "lock.fill")
                    }
                }

                Section("Diaspora") {
                    NavigationLink(destination: Text("Beneficiaries")) {
                        Label("Beneficiaries", systemImage: "person.2.fill")
                    }
                    NavigationLink(destination: Text("Transfer History")) {
                        Label("Transfer History", systemImage: "clock.fill")
                    }
                }

                Section("Preferences") {
                    NavigationLink(destination: Text("Notifications")) {
                        Label("Notifications", systemImage: "bell.fill")
                    }
                    NavigationLink(destination: Text("Currency")) {
                        Label("Currency", systemImage: "dollarsign.circle.fill")
                    }
                }

                Section {
                    Button(action: logout) {
                        Label("Sign Out", systemImage: "rectangle.portrait.and.arrow.right")
                            .foregroundColor(.red)
                    }
                }
            }
            .navigationTitle("Profile")
        }
    }

    private func logout() {
        appState.isAuthenticated = false
        appState.user = nil
    }
}

struct ProfileView_Previews: PreviewProvider {
    static var previews: some View {
        ProfileView()
            .environmentObject(AppState())
    }
}
