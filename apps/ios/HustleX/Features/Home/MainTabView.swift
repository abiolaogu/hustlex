import SwiftUI

struct MainTabView: View {
    @State private var selectedTab = 0

    var body: some View {
        TabView(selection: $selectedTab) {
            HomeView()
                .tabItem {
                    Image(systemName: "house.fill")
                    Text("Home")
                }
                .tag(0)

            ServicesView()
                .tabItem {
                    Image(systemName: "square.grid.2x2.fill")
                    Text("Services")
                }
                .tag(1)

            WalletView()
                .tabItem {
                    Image(systemName: "wallet.pass.fill")
                    Text("Wallet")
                }
                .tag(2)

            RemittanceView()
                .tabItem {
                    Image(systemName: "paperplane.fill")
                    Text("Remit")
                }
                .tag(3)

            ProfileView()
                .tabItem {
                    Image(systemName: "person.fill")
                    Text("Profile")
                }
                .tag(4)
        }
        .accentColor(Theme.Colors.primaryFallback)
    }
}

struct MainTabView_Previews: PreviewProvider {
    static var previews: some View {
        MainTabView()
            .environmentObject(AppState())
    }
}
