import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:logger/logger.dart';

import '../storage/local_cache_service.dart';
import '../storage/secure_storage.dart';
import 'biometric_service.dart';
import 'deep_link_service.dart';
import 'firebase_service.dart';
import 'payment_service.dart';

/// App initialization stages
enum InitStage {
  starting,
  environment,
  storage,
  firebase,
  payments,
  biometrics,
  deepLinks,
  complete,
  error,
}

/// Initialization state
class AppInitState {
  final InitStage stage;
  final double progress;
  final String? error;
  final bool isComplete;

  const AppInitState({
    this.stage = InitStage.starting,
    this.progress = 0,
    this.error,
    this.isComplete = false,
  });

  AppInitState copyWith({
    InitStage? stage,
    double? progress,
    String? error,
    bool? isComplete,
  }) {
    return AppInitState(
      stage: stage ?? this.stage,
      progress: progress ?? this.progress,
      error: error,
      isComplete: isComplete ?? this.isComplete,
    );
  }

  String get stageMessage {
    switch (stage) {
      case InitStage.starting:
        return 'Starting...';
      case InitStage.environment:
        return 'Loading configuration...';
      case InitStage.storage:
        return 'Initializing storage...';
      case InitStage.firebase:
        return 'Setting up notifications...';
      case InitStage.payments:
        return 'Configuring payments...';
      case InitStage.biometrics:
        return 'Setting up security...';
      case InitStage.deepLinks:
        return 'Configuring links...';
      case InitStage.complete:
        return 'Ready!';
      case InitStage.error:
        return 'Error occurred';
    }
  }
}

/// App initialization service
class AppInitService {
  final Logger _logger = Logger(
    printer: PrettyPrinter(
      methodCount: 0,
      errorMethodCount: 5,
      lineLength: 80,
      colors: true,
      printEmojis: true,
      dateTimeFormat: DateTimeFormat.onlyTimeAndSinceStart,
    ),
  );

  /// Initialize the app
  Future<void> initialize({
    required WidgetRef ref,
    void Function(AppInitState)? onStateChange,
  }) async {
    AppInitState state = const AppInitState();
    
    void updateState(InitStage stage, double progress) {
      state = state.copyWith(stage: stage, progress: progress);
      onStateChange?.call(state);
    }

    try {
      _logger.i('🚀 Starting HustleX initialization...');
      
      // Stage 1: Environment
      updateState(InitStage.environment, 0.1);
      await _initEnvironment();

      // Stage 2: Storage
      updateState(InitStage.storage, 0.25);
      await _initStorage();

      // Stage 3: Firebase
      updateState(InitStage.firebase, 0.4);
      await _initFirebase(ref);

      // Stage 4: Payments
      updateState(InitStage.payments, 0.6);
      await _initPayments(ref);

      // Stage 5: Biometrics
      updateState(InitStage.biometrics, 0.75);
      await _initBiometrics(ref);

      // Stage 6: Deep Links
      updateState(InitStage.deepLinks, 0.9);
      await _initDeepLinks(ref);

      // Complete
      updateState(InitStage.complete, 1.0);
      state = state.copyWith(isComplete: true);
      onStateChange?.call(state);

      _logger.i('✅ HustleX initialization complete!');
    } catch (e, stack) {
      _logger.e('❌ Initialization failed', error: e, stackTrace: stack);
      state = state.copyWith(
        stage: InitStage.error,
        error: e.toString(),
      );
      onStateChange?.call(state);
      rethrow;
    }
  }

  /// Initialize environment and configuration
  Future<void> _initEnvironment() async {
    _logger.d('Loading environment configuration...');

    // Load .env file
    try {
      await dotenv.load(fileName: '.env');
      _logger.d('Environment loaded');
    } catch (e) {
      _logger.w('Could not load .env file, using defaults');
    }

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

    _logger.i('✓ Environment initialized');
  }

  /// Initialize storage systems
  Future<void> _initStorage() async {
    _logger.d('Initializing storage systems...');

    // Initialize Hive
    await Hive.initFlutter();

    // Initialize local cache
    final cacheService = LocalCacheService();
    await cacheService.initialize();

    // Verify secure storage
    final secureStorage = SecureStorage();
    await secureStorage.containsKey('test'); // Warm up

    _logger.i('✓ Storage initialized');
  }

  /// Initialize Firebase services
  Future<void> _initFirebase(WidgetRef ref) async {
    _logger.d('Initializing Firebase...');

    try {
      final firebaseService = ref.read(firebaseServiceProvider);
      await firebaseService.initialize();
      _logger.i('✓ Firebase initialized');
    } catch (e) {
      _logger.w('Firebase initialization failed: $e');
      // Non-critical, app can continue without Firebase
    }
  }

  /// Initialize payment services
  Future<void> _initPayments(WidgetRef ref) async {
    _logger.d('Initializing payment services...');

    try {
      final paymentService = ref.read(paymentServiceProvider);
      await paymentService.initialize();
      _logger.i('✓ Payment services initialized');
    } catch (e) {
      _logger.w('Payment service initialization failed: $e');
      // Non-critical in some cases
    }
  }

  /// Initialize biometric services
  Future<void> _initBiometrics(WidgetRef ref) async {
    _logger.d('Checking biometric capabilities...');

    try {
      final biometricService = ref.read(biometricServiceProvider);
      final isAvailable = await biometricService.canCheckBiometrics();
      
      if (isAvailable) {
        final types = await biometricService.getAvailableBiometrics();
        _logger.d('Available biometrics: $types');
      }
      
      _logger.i('✓ Biometrics checked');
    } catch (e) {
      _logger.w('Biometric check failed: $e');
    }
  }

  /// Initialize deep link handling
  Future<void> _initDeepLinks(WidgetRef ref) async {
    _logger.d('Initializing deep link handling...');

    try {
      final deepLinkService = ref.read(deepLinkServiceProvider);
      await deepLinkService.initialize();
      _logger.i('✓ Deep links initialized');
    } catch (e) {
      _logger.w('Deep link initialization failed: $e');
    }
  }
}

/// Quick initialization for essential services only
class QuickInitService {
  static Future<void> initialize() async {
    // Essential initialization only
    WidgetsFlutterBinding.ensureInitialized();

    // Load environment
    try {
      await dotenv.load(fileName: '.env');
    } catch (_) {}

    // Initialize Hive
    await Hive.initFlutter();

    // Set orientations
    await SystemChrome.setPreferredOrientations([
      DeviceOrientation.portraitUp,
      DeviceOrientation.portraitDown,
    ]);
  }
}

/// Initialize app before runApp
Future<void> initializeApp() async {
  await QuickInitService.initialize();
}

/// App initialization service provider
final appInitServiceProvider = Provider<AppInitService>((ref) {
  return AppInitService();
});

/// App init state provider
final appInitStateProvider = StateProvider<AppInitState>((ref) {
  return const AppInitState();
});

/// App ready provider
final appReadyProvider = Provider<bool>((ref) {
  final state = ref.watch(appInitStateProvider);
  return state.isComplete;
});

/// Initialization progress provider
final initProgressProvider = Provider<double>((ref) {
  final state = ref.watch(appInitStateProvider);
  return state.progress;
});

/// Logger configuration for the app
Logger createLogger({
  String? name,
  bool includeTimestamp = true,
}) {
  return Logger(
    printer: PrettyPrinter(
      methodCount: kDebugMode ? 2 : 0,
      errorMethodCount: 8,
      lineLength: 120,
      colors: true,
      printEmojis: true,
      dateTimeFormat: includeTimestamp
          ? DateTimeFormat.onlyTimeAndSinceStart
          : DateTimeFormat.none,
    ),
    filter: kDebugMode ? DevelopmentFilter() : ProductionFilter(),
  );
}

/// App-wide logger instance
final appLogger = createLogger(name: 'HustleX');
