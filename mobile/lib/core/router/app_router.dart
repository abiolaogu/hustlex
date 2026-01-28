import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

// Auth screens
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/otp_screen.dart';
import '../../features/auth/presentation/screens/register_screen.dart';
import '../../features/auth/presentation/screens/pin_setup_screen.dart';
import '../../features/auth/presentation/screens/onboarding_screen.dart';

// Main screens
import '../../features/home/presentation/screens/main_shell.dart';
import '../../features/home/presentation/screens/home_screen.dart';

// Gig screens
import '../../features/gigs/presentation/screens/gigs_screen.dart';
import '../../features/gigs/presentation/screens/gig_details_screen.dart';
import '../../features/gigs/presentation/screens/create_gig_screen.dart';
import '../../features/gigs/presentation/screens/submit_proposal_screen.dart';
import '../../features/gigs/presentation/screens/my_gigs_screen.dart';

// Savings screens
import '../../features/savings/presentation/screens/savings_screen.dart';
import '../../features/savings/presentation/screens/circle_details_screen.dart';
import '../../features/savings/presentation/screens/create_circle_screen.dart';
import '../../features/savings/presentation/screens/join_circle_screen.dart';

// Wallet screens
import '../../features/wallet/presentation/screens/wallet_screen.dart';
import '../../features/wallet/presentation/screens/deposit_screen.dart';
import '../../features/wallet/presentation/screens/withdraw_screen.dart';
import '../../features/wallet/presentation/screens/transfer_screen.dart';
import '../../features/wallet/presentation/screens/transactions_screen.dart';
import '../../features/wallet/presentation/screens/transaction_details_screen.dart';
import '../../features/wallet/presentation/screens/bank_accounts_screen.dart';
import '../../features/wallet/presentation/screens/add_bank_account_screen.dart';

// Credit screens
import '../../features/credit/presentation/screens/credit_screen.dart';
import '../../features/credit/presentation/screens/loan_application_screen.dart';
import '../../features/credit/presentation/screens/loan_details_screen.dart';

// Profile screens
import '../../features/profile/presentation/screens/profile_screen.dart';
import '../../features/profile/presentation/screens/edit_profile_screen.dart';
import '../../features/profile/presentation/screens/change_pin_screen.dart';
import '../../features/profile/presentation/screens/settings_screen.dart';

// Notifications
import '../../features/notifications/presentation/screens/notifications_screen.dart';

// Route names
class AppRoutes {
  AppRoutes._();

  // Initial routes
  static const String splash = '/';
  static const String onboarding = '/onboarding';

  // Auth routes
  static const String login = '/auth/login';
  static const String otpVerification = '/auth/otp-verification';
  static const String registration = '/auth/registration';
  static const String pinSetup = '/auth/pin-setup';

  // Main app routes (with bottom nav)
  static const String home = '/home';
  static const String gigs = '/gigs';
  static const String savings = '/savings';
  static const String wallet = '/wallet';
  static const String profile = '/profile';

  // Gig routes
  static const String gigDetails = '/gigs/:id';
  static const String createGig = '/gigs/create';
  static const String myGigs = '/gigs/my';
  static const String submitProposal = '/gigs/:id/propose';

  // Savings routes
  static const String circleDetails = '/savings/circles/:id';
  static const String createCircle = '/savings/circles/create';
  static const String joinCircle = '/savings/circles/join';

  // Wallet routes
  static const String deposit = '/wallet/deposit';
  static const String withdraw = '/wallet/withdraw';
  static const String transfer = '/wallet/transfer';
  static const String transactions = '/wallet/transactions';
  static const String transactionDetails = '/wallet/transactions/:id';
  static const String bankAccounts = '/wallet/banks';
  static const String addBankAccount = '/wallet/add-bank';

  // Credit routes
  static const String credit = '/credit';
  static const String loanApplication = '/credit/apply';
  static const String loanDetails = '/credit/loans/:id';

  // Profile routes
  static const String editProfile = '/profile/edit';
  static const String changePin = '/profile/change-pin';
  static const String settings = '/settings';
  static const String notifications = '/notifications';
}

// Router provider
final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: AppRoutes.onboarding,
    debugLogDiagnostics: true,
    routes: [
      // Onboarding/Splash
      GoRoute(
        path: AppRoutes.splash,
        name: 'splash',
        builder: (context, state) => const OnboardingScreen(),
      ),
      GoRoute(
        path: AppRoutes.onboarding,
        name: 'onboarding',
        builder: (context, state) => const OnboardingScreen(),
      ),

      // Auth routes
      GoRoute(
        path: AppRoutes.login,
        name: 'login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: AppRoutes.otpVerification,
        name: 'otp-verification',
        builder: (context, state) {
          final extra = state.extra as Map<String, dynamic>? ?? {};
          return OtpScreen(
            phoneNumber: extra['phoneNumber'] ?? '',
            isRegistration: extra['isRegistration'] ?? false,
          );
        },
      ),
      GoRoute(
        path: AppRoutes.registration,
        name: 'registration',
        builder: (context, state) {
          final phoneNumber = state.extra as String? ?? '';
          return RegisterScreen(phoneNumber: phoneNumber);
        },
      ),
      GoRoute(
        path: AppRoutes.pinSetup,
        name: 'pin-setup',
        builder: (context, state) => const PinSetupScreen(),
      ),

      // Main app shell with bottom navigation
      ShellRoute(
        builder: (context, state, child) {
          return MainShell(child: child);
        },
        routes: [
          // Home
          GoRoute(
            path: AppRoutes.home,
            name: 'home',
            pageBuilder: (context, state) => const NoTransitionPage(
              child: HomeScreen(),
            ),
          ),

          // Gigs
          GoRoute(
            path: AppRoutes.gigs,
            name: 'gigs',
            pageBuilder: (context, state) => const NoTransitionPage(
              child: GigsScreen(),
            ),
          ),

          // Savings
          GoRoute(
            path: AppRoutes.savings,
            name: 'savings',
            pageBuilder: (context, state) => const NoTransitionPage(
              child: SavingsScreen(),
            ),
          ),

          // Wallet
          GoRoute(
            path: AppRoutes.wallet,
            name: 'wallet',
            pageBuilder: (context, state) => const NoTransitionPage(
              child: WalletScreen(),
            ),
          ),

          // Profile
          GoRoute(
            path: AppRoutes.profile,
            name: 'profile',
            pageBuilder: (context, state) => const NoTransitionPage(
              child: ProfileScreen(),
            ),
          ),
        ],
      ),

      // Gig detail routes (outside shell for full screen)
      GoRoute(
        path: '/gigs/create',
        name: 'create-gig',
        builder: (context, state) => const CreateGigScreen(),
      ),
      GoRoute(
        path: '/gigs/my',
        name: 'my-gigs',
        builder: (context, state) => const MyGigsScreen(),
      ),
      GoRoute(
        path: '/gigs/:id',
        name: 'gig-details',
        builder: (context, state) {
          final gigId = state.pathParameters['id']!;
          return GigDetailsScreen(gigId: gigId);
        },
        routes: [
          GoRoute(
            path: 'propose',
            name: 'submit-proposal',
            builder: (context, state) {
              final gigId = state.pathParameters['id']!;
              return SubmitProposalScreen(gigId: gigId);
            },
          ),
        ],
      ),

      // Savings detail routes
      GoRoute(
        path: '/savings/circles/create',
        name: 'create-circle',
        builder: (context, state) => const CreateCircleScreen(),
      ),
      GoRoute(
        path: '/savings/circles/join',
        name: 'join-circle',
        builder: (context, state) => const JoinCircleScreen(),
      ),
      GoRoute(
        path: '/savings/circles/:id',
        name: 'circle-details',
        builder: (context, state) {
          final circleId = state.pathParameters['id']!;
          return CircleDetailsScreen(circleId: circleId);
        },
      ),

      // Wallet routes
      GoRoute(
        path: AppRoutes.deposit,
        name: 'deposit',
        builder: (context, state) => const DepositScreen(),
      ),
      GoRoute(
        path: AppRoutes.withdraw,
        name: 'withdraw',
        builder: (context, state) => const WithdrawScreen(),
      ),
      GoRoute(
        path: AppRoutes.transfer,
        name: 'transfer',
        builder: (context, state) => const TransferScreen(),
      ),
      GoRoute(
        path: AppRoutes.transactions,
        name: 'transactions',
        builder: (context, state) => const TransactionsScreen(),
      ),
      GoRoute(
        path: '/wallet/transactions/:id',
        name: 'transaction-details',
        builder: (context, state) {
          final transactionId = state.pathParameters['id']!;
          return TransactionDetailsScreen(transactionId: transactionId);
        },
      ),
      GoRoute(
        path: AppRoutes.bankAccounts,
        name: 'bank-accounts',
        builder: (context, state) => const BankAccountsScreen(),
      ),
      GoRoute(
        path: AppRoutes.addBankAccount,
        name: 'add-bank-account',
        builder: (context, state) => const AddBankAccountScreen(),
      ),

      // Credit routes
      GoRoute(
        path: AppRoutes.credit,
        name: 'credit',
        builder: (context, state) => const CreditScreen(),
      ),
      GoRoute(
        path: AppRoutes.loanApplication,
        name: 'loan-application',
        builder: (context, state) {
          final extra = state.extra as Map<String, dynamic>? ?? {};
          return LoanApplicationScreen(
            creditLimit: extra['creditLimit']?.toDouble() ?? 100000,
            creditScore: extra['creditScore'] ?? 650,
          );
        },
      ),
      GoRoute(
        path: '/credit/loans/:id',
        name: 'loan-details',
        builder: (context, state) {
          final loanId = state.pathParameters['id']!;
          return LoanDetailsScreen(loanId: loanId);
        },
      ),

      // Profile routes
      GoRoute(
        path: AppRoutes.editProfile,
        name: 'edit-profile',
        builder: (context, state) => const EditProfileScreen(),
      ),
      GoRoute(
        path: AppRoutes.changePin,
        name: 'change-pin',
        builder: (context, state) => const ChangePinScreen(),
      ),

      // Settings & Notifications
      GoRoute(
        path: AppRoutes.settings,
        name: 'settings',
        builder: (context, state) => const SettingsScreen(),
      ),
      GoRoute(
        path: AppRoutes.notifications,
        name: 'notifications',
        builder: (context, state) => const NotificationsScreen(),
      ),
    ],

    // Error handling
    errorBuilder: (context, state) => Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        title: const Text('Page Not Found'),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.error_outline,
                size: 80,
                color: Colors.grey.shade400,
              ),
              const SizedBox(height: 24),
              Text(
                '404',
                style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: Colors.grey.shade600,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Page not found',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  color: Colors.grey.shade600,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                state.uri.toString(),
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: Colors.grey.shade500,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 32),
              ElevatedButton.icon(
                onPressed: () => context.go(AppRoutes.home),
                icon: const Icon(Icons.home),
                label: const Text('Go Home'),
              ),
            ],
          ),
        ),
      ),
    ),

    // Redirect logic
    redirect: (context, state) {
      // TODO: Add authentication state checks
      // final isLoggedIn = ref.read(authProvider).isLoggedIn;
      // final isOnAuthRoute = state.matchedLocation.startsWith('/auth');
      
      // if (!isLoggedIn && !isOnAuthRoute && state.matchedLocation != AppRoutes.onboarding) {
      //   return AppRoutes.login;
      // }
      
      // if (isLoggedIn && isOnAuthRoute) {
      //   return AppRoutes.home;
      // }
      
      return null;
    },
  );
});
