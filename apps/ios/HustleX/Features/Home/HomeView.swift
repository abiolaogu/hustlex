import SwiftUI

struct HomeView: View {
    @EnvironmentObject var appState: AppState

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: Theme.Spacing.lg) {
                    // Balance Card
                    WalletCard(balance: "₦250,000.00")

                    // Quick Actions
                    VStack(alignment: .leading, spacing: Theme.Spacing.md) {
                        Text("Quick Actions")
                            .font(Theme.Typography.titleMedium)

                        HStack(spacing: Theme.Spacing.md) {
                            QuickActionButton(icon: "arrow.up.right", label: "Transfer")
                            QuickActionButton(icon: "airplane.departure", label: "Remit")
                            QuickActionButton(icon: "phone.fill", label: "Airtime")
                            QuickActionButton(icon: "ellipsis", label: "More")
                        }
                    }
                    .padding(.horizontal)

                    // Services Section
                    VStack(alignment: .leading, spacing: Theme.Spacing.md) {
                        HStack {
                            Text("Services")
                                .font(Theme.Typography.titleMedium)
                            Spacer()
                            Button("See All") {}
                                .font(Theme.Typography.labelMedium)
                                .foregroundColor(Theme.Colors.primaryFallback)
                        }

                        ScrollView(.horizontal, showsIndicators: false) {
                            HStack(spacing: Theme.Spacing.md) {
                                ServiceCard(name: "Cleaning", icon: "sparkles")
                                ServiceCard(name: "Plumbing", icon: "wrench.fill")
                                ServiceCard(name: "Electrical", icon: "bolt.fill")
                                ServiceCard(name: "Beauty", icon: "face.smiling.fill")
                            }
                        }
                    }
                    .padding(.horizontal)

                    // Recent Transactions
                    VStack(alignment: .leading, spacing: Theme.Spacing.md) {
                        Text("Recent Transactions")
                            .font(Theme.Typography.titleMedium)

                        TransactionRow(title: "House Cleaning", date: "Yesterday", amount: "-₦15,000")
                        TransactionRow(title: "Airtime Purchase", date: "2 days ago", amount: "-₦2,000")
                        TransactionRow(title: "Wallet Top-up", date: "3 days ago", amount: "+₦50,000", isCredit: true)
                    }
                    .padding(.horizontal)
                }
                .padding(.vertical)
            }
            .navigationTitle("Home")
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: {}) {
                        Image(systemName: "bell.fill")
                    }
                }
            }
        }
    }
}

struct WalletCard: View {
    let balance: String

    var body: some View {
        VStack(alignment: .leading, spacing: Theme.Spacing.sm) {
            Text("Available Balance")
                .font(Theme.Typography.bodyMedium)
                .foregroundColor(.white.opacity(0.8))

            Text(balance)
                .font(Theme.Typography.headlineLarge)
                .foregroundColor(.white)

            HStack {
                Button(action: {}) {
                    HStack {
                        Image(systemName: "plus")
                        Text("Add Money")
                    }
                }
                .foregroundColor(.white)

                Spacer()

                Button(action: {}) {
                    HStack {
                        Image(systemName: "arrow.down")
                        Text("Withdraw")
                    }
                }
                .foregroundColor(.white)
            }
            .font(Theme.Typography.labelMedium)
        }
        .padding()
        .background(
            LinearGradient(
                gradient: Gradient(colors: [Theme.Colors.primaryFallback, Theme.Colors.primaryFallback.opacity(0.8)]),
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
        )
        .cornerRadius(Theme.CornerRadius.lg)
        .padding(.horizontal)
    }
}

struct QuickActionButton: View {
    let icon: String
    let label: String

    var body: some View {
        VStack(spacing: Theme.Spacing.xs) {
            Image(systemName: icon)
                .font(.title2)
                .frame(width: 56, height: 56)
                .background(Theme.Colors.primaryFallback.opacity(0.1))
                .foregroundColor(Theme.Colors.primaryFallback)
                .clipShape(Circle())

            Text(label)
                .font(Theme.Typography.labelSmall)
        }
        .frame(maxWidth: .infinity)
    }
}

struct ServiceCard: View {
    let name: String
    let icon: String

    var body: some View {
        VStack(spacing: Theme.Spacing.sm) {
            Image(systemName: icon)
                .font(.title2)
                .foregroundColor(Theme.Colors.primaryFallback)

            Text(name)
                .font(Theme.Typography.labelMedium)
        }
        .frame(width: 100, height: 80)
        .background(Color(.systemGray6))
        .cornerRadius(Theme.CornerRadius.md)
    }
}

struct TransactionRow: View {
    let title: String
    let date: String
    let amount: String
    var isCredit: Bool = false

    var body: some View {
        HStack {
            VStack(alignment: .leading) {
                Text(title)
                    .font(Theme.Typography.bodyMedium)
                Text(date)
                    .font(Theme.Typography.bodySmall)
                    .foregroundColor(.secondary)
            }

            Spacer()

            Text(amount)
                .font(Theme.Typography.bodyMedium)
                .foregroundColor(isCredit ? Theme.Colors.secondaryFallback : .primary)
        }
        .padding(.vertical, Theme.Spacing.xs)
    }
}

struct HomeView_Previews: PreviewProvider {
    static var previews: some View {
        HomeView()
            .environmentObject(AppState())
    }
}
