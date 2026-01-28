import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../core/providers/auth_provider.dart';
import '../features/auth/presentation/screens/login_screen.dart';
import '../features/auth/presentation/screens/onboarding_screen.dart';
import '../features/auth/presentation/screens/otp_screen.dart';
import '../features/auth/presentation/screens/pin_setup_screen.dart';
import '../features/auth/presentation/screens/register_screen.dart';
import '../features/credit/presentation/screens/credit_screen.dart';
import '../features/gigs/presentation/screens/gigs_screen.dart';
import '../features/home/presentation/screens/home_screen.dart';
import '../features/home/presentation/screens/main_shell.dart';
import '../features/notifications/presentation/screens/notifications_screen.dart';
import '../features/profile/presentation/screens/profile_screen.dart';
import '../features/savings/presentation/screens/savings_screen.dart';
import '../features/wallet/presentation/screens/wallet_screen.dart';

/// Route names as constants
class AppRoutes {
  static const String splash = '/';
  static const String onboarding = '/onboarding';
  static const String login = '/login';
  static const String otp = '/otp';
  static const String register = '/register';
  static const String pinSetup = '/pin-setup';
  static const String home = '/home';
  static const String gigs = '/gigs';
  static const String gigDetails = '/gigs/:id';
  static const String createGig = '/gigs/create';
  static const String savings = '/savings';
  static const String circleDetails = '/savings/:id';
  static const String createCircle = '/savings/create';
  static const String wallet = '/wallet';
  static const String deposit = '/wallet/deposit';
  static const String withdraw = '/wallet/withdraw';
  static const String transfer = '/wallet/transfer';
  static const String transactions = '/wallet/transactions';
  static const String credit = '/credit';
  static const String loanApplication = '/credit/apply';
  static const String profile = '/profile';
  static const String editProfile = '/profile/edit';
  static const String settings = '/profile/settings';
  static const String notifications = '/notifications';
}

/// Navigator key for shell routes
final _rootNavigatorKey = GlobalKey<NavigatorState>();
final _shellNavigatorKey = GlobalKey<NavigatorState>();

/// Router provider
final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: AppRoutes.splash,
    debugLogDiagnostics: true,
    redirect: (context, state) {
      final isLoading = authState.isLoading;
      final auth = authState.valueOrNull;
      final isAuthenticated = auth?.isAuthenticated ?? false;
      final hasCompletedOnboarding = auth?.hasCompletedOnboarding ?? false;
      final hasPinSet = auth?.hasPinSet ?? false;

      final location = state.matchedLocation;
      final isAuthRoute = location == AppRoutes.login ||
          location == AppRoutes.otp ||
          location == AppRoutes.register ||
          location == AppRoutes.pinSetup;
      final isOnboardingRoute = location == AppRoutes.onboarding;
      final isSplashRoute = location == AppRoutes.splash;

      // Still loading auth state
      if (isLoading) {
        return isSplashRoute ? null : AppRoutes.splash;
      }

      // Not completed onboarding
      if (!hasCompletedOnboarding && !isOnboardingRoute && !isSplashRoute) {
        return AppRoutes.onboarding;
      }

      // Not authenticated
      if (!isAuthenticated) {
        if (isAuthRoute || isOnboardingRoute) return null;
        return AppRoutes.login;
      }

      // Authenticated but no PIN set
      if (!hasPinSet && location != AppRoutes.pinSetup) {
        return AppRoutes.pinSetup;
      }

      // Authenticated, redirect away from auth routes
      if (isAuthRoute || isOnboardingRoute || isSplashRoute) {
        return AppRoutes.home;
      }

      return null;
    },
    routes: [
      // Splash
      GoRoute(
        path: AppRoutes.splash,
        builder: (context, state) => const SplashScreen(),
      ),

      // Onboarding
      GoRoute(
        path: AppRoutes.onboarding,
        builder: (context, state) => const OnboardingScreen(),
      ),

      // Auth routes
      GoRoute(
        path: AppRoutes.login,
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: AppRoutes.otp,
        builder: (context, state) {
          final phone = state.extra as String? ?? '';
          return OtpScreen(phoneNumber: phone);
        },
      ),
      GoRoute(
        path: AppRoutes.register,
        builder: (context, state) {
          final phone = state.extra as String? ?? '';
          return RegisterScreen(phoneNumber: phone);
        },
      ),
      GoRoute(
        path: AppRoutes.pinSetup,
        builder: (context, state) => const PinSetupScreen(),
      ),

      // Main app shell with bottom navigation
      ShellRoute(
        navigatorKey: _shellNavigatorKey,
        builder: (context, state, child) => MainShell(child: child),
        routes: [
          // Home
          GoRoute(
            path: AppRoutes.home,
            pageBuilder: (context, state) => const NoTransitionPage(
              child: HomeScreen(),
            ),
          ),

          // Gigs
          GoRoute(
            path: AppRoutes.gigs,
            pageBuilder: (context, state) => const NoTransitionPage(
              child: GigsScreen(),
            ),
            routes: [
              GoRoute(
                path: 'create',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) => const CreateGigScreen(),
              ),
              GoRoute(
                path: ':id',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) {
                  final id = state.pathParameters['id'] ?? '';
                  return GigDetailsScreen(gigId: id);
                },
              ),
            ],
          ),

          // Savings
          GoRoute(
            path: AppRoutes.savings,
            pageBuilder: (context, state) => const NoTransitionPage(
              child: SavingsScreen(),
            ),
            routes: [
              GoRoute(
                path: 'create',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) => const CreateCircleScreen(),
              ),
              GoRoute(
                path: ':id',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) {
                  final id = state.pathParameters['id'] ?? '';
                  return CircleDetailsScreen(circleId: id);
                },
              ),
            ],
          ),

          // Wallet
          GoRoute(
            path: AppRoutes.wallet,
            pageBuilder: (context, state) => const NoTransitionPage(
              child: WalletScreen(),
            ),
            routes: [
              GoRoute(
                path: 'deposit',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) => const DepositScreen(),
              ),
              GoRoute(
                path: 'withdraw',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) => const WithdrawScreen(),
              ),
              GoRoute(
                path: 'transfer',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) => const TransferScreen(),
              ),
            ],
          ),

          // Profile
          GoRoute(
            path: AppRoutes.profile,
            pageBuilder: (context, state) => const NoTransitionPage(
              child: ProfileScreen(),
            ),
          ),
        ],
      ),

      // Routes outside shell (full screen)
      GoRoute(
        path: AppRoutes.credit,
        builder: (context, state) => const CreditScreen(),
        routes: [
          GoRoute(
            path: 'apply',
            builder: (context, state) => const LoanApplicationScreen(),
          ),
        ],
      ),
      GoRoute(
        path: AppRoutes.notifications,
        builder: (context, state) => const NotificationsScreen(),
      ),
    ],
    errorBuilder: (context, state) => ErrorScreen(error: state.error),
  );
});

/// Splash screen
class SplashScreen extends ConsumerWidget {
  const SplashScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF6C63FF), Color(0xFF5A52E0)],
          ),
        ),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 100,
                height: 100,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(24),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.1),
                      blurRadius: 20,
                      offset: const Offset(0, 10),
                    ),
                  ],
                ),
                child: const Icon(
                  Icons.rocket_launch_rounded,
                  size: 50,
                  color: Color(0xFF6C63FF),
                ),
              ),
              const SizedBox(height: 24),
              const Text(
                'HustleX',
                style: TextStyle(
                  fontSize: 32,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Hustle Smarter',
                style: TextStyle(
                  fontSize: 16,
                  color: Colors.white.withOpacity(0.8),
                ),
              ),
              const SizedBox(height: 48),
              const CircularProgressIndicator(
                valueColor: AlwaysStoppedAnimation(Colors.white),
                strokeWidth: 2,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Error screen
class ErrorScreen extends StatelessWidget {
  final Exception? error;

  const ErrorScreen({super.key, this.error});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(
                Icons.error_outline_rounded,
                size: 64,
                color: Colors.red,
              ),
              const SizedBox(height: 16),
              const Text(
                'Oops! Something went wrong',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 8),
              Text(
                error?.toString() ?? 'Page not found',
                style: const TextStyle(color: Colors.grey),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: () => context.go(AppRoutes.home),
                child: const Text('Go Home'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
