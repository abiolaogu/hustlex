/// HustleX Core Module
/// 
/// This barrel file exports all core functionality for the app.
/// 
/// Usage:
/// ```dart
/// import 'package:hustlex/core/core.dart';
/// ```

// API
export 'api/api_client.dart';

// Constants
export 'constants/app_colors.dart';
export 'constants/app_constants.dart';
export 'constants/app_theme.dart';
export 'constants/app_typography.dart';

// Dependency Injection
export 'di/providers.dart';

// Exceptions
export 'exceptions/api_exception.dart';

// Repositories
export 'repositories/base_repository.dart';

// Router
export 'router/router.dart';

// Services
export 'services/services.dart';

// Storage
export 'storage/local_cache_service.dart';
export 'storage/secure_storage.dart';

// Utils
export 'utils/utils.dart';
