import SwiftUI

struct LoginView: View {
    @EnvironmentObject var appState: AppState
    @State private var phone = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var showPassword = false

    var body: some View {
        NavigationView {
            VStack(spacing: Theme.Spacing.lg) {
                Spacer()

                // Logo and Title
                VStack(spacing: Theme.Spacing.sm) {
                    Text("HustleX")
                        .font(Theme.Typography.displaySmall)
                        .foregroundColor(Theme.Colors.primaryFallback)

                    Text("Sign in to continue")
                        .font(Theme.Typography.bodyLarge)
                        .foregroundColor(.secondary)
                }

                Spacer()

                // Form
                VStack(spacing: Theme.Spacing.md) {
                    TextField("Phone Number", text: $phone)
                        .keyboardType(.phonePad)
                        .textFieldStyle(.roundedBorder)

                    HStack {
                        if showPassword {
                            TextField("Password", text: $password)
                        } else {
                            SecureField("Password", text: $password)
                        }

                        Button(action: { showPassword.toggle() }) {
                            Image(systemName: showPassword ? "eye.slash" : "eye")
                                .foregroundColor(.secondary)
                        }
                    }
                    .textFieldStyle(.roundedBorder)

                    Button(action: login) {
                        if isLoading {
                            ProgressView()
                                .progressViewStyle(CircularProgressViewStyle(tint: .white))
                        } else {
                            Text("Sign In")
                                .font(Theme.Typography.labelLarge)
                        }
                    }
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(Theme.Colors.primaryFallback)
                    .foregroundColor(.white)
                    .cornerRadius(Theme.CornerRadius.md)
                    .disabled(phone.isEmpty || password.isEmpty || isLoading)

                    Button("Forgot Password?") {
                        // TODO: Forgot password
                    }
                    .font(Theme.Typography.labelMedium)
                }
                .padding(.horizontal, Theme.Spacing.lg)

                Spacer()

                // Sign Up Link
                HStack {
                    Text("Don't have an account?")
                        .foregroundColor(.secondary)
                    Button("Sign Up") {
                        // TODO: Navigate to register
                    }
                    .foregroundColor(Theme.Colors.primaryFallback)
                }
                .font(Theme.Typography.bodyMedium)
                .padding(.bottom, Theme.Spacing.lg)
            }
            .navigationBarHidden(true)
        }
    }

    private func login() {
        isLoading = true
        // TODO: Implement actual login
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            appState.isAuthenticated = true
            isLoading = false
        }
    }
}

struct LoginView_Previews: PreviewProvider {
    static var previews: some View {
        LoginView()
            .environmentObject(AppState())
    }
}
