import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:timeago/timeago.dart' as timeago;

import 'core/constants/app_constants.dart';
import 'core/constants/app_theme.dart';
import 'core/services/app_init_service.dart';
import 'core/services/deep_link_service.dart';
import 'core/services/firebase_service.dart';
import 'router/app_router.dart';

void main() async {
  // Ensure Flutter bindings are initialized
  WidgetsFlutterBinding.ensureInitialized();

  // Run in a zone to catch errors
  await runZonedGuarded(
    () async {
      // Load environment variables
      await dotenv.load(fileName: '.env').catchError((_) {
        // Use defaults if .env not found
        debugPrint('No .env file found, using defaults');
      });

      // Initialize Hive for local storage
      await Hive.initFlutter();

      // Set preferred orientations
      await SystemChrome.setPreferredOrientations([
        DeviceOrientation.portraitUp,
        DeviceOrientation.portraitDown,
      ]);

      // Set system UI overlay style
      SystemChrome.setSystemUIOverlayStyle(
        const SystemUiOverlayStyle(
          statusBarColor: Colors.transparent,
          statusBarIconBrightness: Brightness.dark,
          statusBarBrightness: Brightness.light,
          systemNavigationBarColor: Colors.white,
          systemNavigationBarIconBrightness: Brightness.dark,
        ),
      );

      // Configure timeago messages
      timeago.setLocaleMessages('en', timeago.EnMessages());

      // Run the app
      runApp(
        const ProviderScope(
          child: HustleXApp(),
        ),
      );
    },
    (error, stackTrace) {
      // Log errors in production
      debugPrint('Uncaught error: $error');
      debugPrint('Stack trace: $stackTrace');
    },
  );
}

class HustleXApp extends ConsumerStatefulWidget {
  const HustleXApp({super.key});

  @override
  ConsumerState<HustleXApp> createState() => _HustleXAppState();
}

class _HustleXAppState extends ConsumerState<HustleXApp> with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _initializeServices();
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  /// Initialize background services after widget tree is built
  Future<void> _initializeServices() async {
    // Wait for widget tree to be built
    await Future.delayed(const Duration(milliseconds: 100));

    if (!mounted) return;

    try {
      // Initialize Firebase
      final firebaseService = ref.read(firebaseServiceProvider);
      await firebaseService.initialize();

      // Initialize deep link service
      final deepLinkService = ref.read(deepLinkServiceProvider);
      await deepLinkService.initialize();

      // Listen for deep links
      deepLinkService.onDeepLink.listen((link) {
        if (!mounted) return;
        final router = ref.read(routerProvider);
        deepLinkService.navigateToDeepLink(router, link);
      });

      debugPrint('✅ Background services initialized');
    } catch (e) {
      debugPrint('❌ Error initializing services: $e');
    }
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);

    switch (state) {
      case AppLifecycleState.resumed:
        // App came to foreground
        _onAppResumed();
        break;
      case AppLifecycleState.paused:
        // App went to background
        _onAppPaused();
        break;
      case AppLifecycleState.inactive:
      case AppLifecycleState.detached:
      case AppLifecycleState.hidden:
        break;
    }
  }

  void _onAppResumed() {
    // Clear notification badge
    try {
      final firebaseService = ref.read(firebaseServiceProvider);
      firebaseService.clearBadge();
    } catch (_) {}
  }

  void _onAppPaused() {
    // Handle app backgrounding if needed
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: AppConstants.appName,
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: ThemeMode.light,
      routerConfig: router,
      builder: (context, child) {
        // Limit text scaling to prevent layout issues
        return MediaQuery(
          data: MediaQuery.of(context).copyWith(
            textScaler: TextScaler.linear(
              MediaQuery.of(context).textScaler.scale(1.0).clamp(0.8, 1.2),
            ),
          ),
          child: _ErrorBoundary(child: child!),
        );
      },
    );
  }
}

/// Error boundary widget to catch rendering errors
class _ErrorBoundary extends StatefulWidget {
  final Widget child;

  const _ErrorBoundary({required this.child});

  @override
  State<_ErrorBoundary> createState() => _ErrorBoundaryState();
}

class _ErrorBoundaryState extends State<_ErrorBoundary> {
  bool _hasError = false;
  FlutterErrorDetails? _errorDetails;

  @override
  void initState() {
    super.initState();
    
    // In debug mode, don't catch errors (let them bubble up for easier debugging)
    if (!kDebugMode) {
      FlutterError.onError = (details) {
        setState(() {
          _hasError = true;
          _errorDetails = details;
        });
      };
    }
  }

  void _resetError() {
    setState(() {
      _hasError = false;
      _errorDetails = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_hasError) {
      return Material(
        child: SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(
                  Icons.error_outline,
                  size: 64,
                  color: Colors.red,
                ),
                const SizedBox(height: 16),
                const Text(
                  'Something went wrong',
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  'An unexpected error occurred. Please try again.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 16,
                    color: Colors.grey[600],
                  ),
                ),
                if (kDebugMode && _errorDetails != null) ...[
                  const SizedBox(height: 16),
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: Colors.red[50],
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      _errorDetails!.exceptionAsString(),
                      style: const TextStyle(
                        fontSize: 12,
                        fontFamily: 'monospace',
                      ),
                    ),
                  ),
                ],
                const SizedBox(height: 24),
                ElevatedButton(
                  onPressed: _resetError,
                  child: const Text('Try Again'),
                ),
              ],
            ),
          ),
        ),
      );
    }

    return widget.child;
  }
}
