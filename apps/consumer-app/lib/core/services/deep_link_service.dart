import 'dart:async';

import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

/// Deep link types
enum DeepLinkType {
  gig,
  circle,
  wallet,
  transaction,
  profile,
  referral,
  promo,
  verify,
  resetPassword,
  unknown,
}

/// Parsed deep link
class DeepLink {
  final DeepLinkType type;
  final String? id;
  final Map<String, String> params;
  final String rawUri;

  const DeepLink({
    required this.type,
    this.id,
    this.params = const {},
    required this.rawUri,
  });

  @override
  String toString() => 'DeepLink(type: $type, id: $id, params: $params)';
}

/// Deep link service
class DeepLinkService {
  StreamSubscription<Uri?>? _linkSubscription;
  
  // Callback for handling deep links
  void Function(DeepLink link)? onDeepLink;
  
  // Initial link (app opened via deep link)
  DeepLink? _initialLink;
  DeepLink? get initialLink => _initialLink;

  /// Initialize deep link handling
  Future<void> initialize() async {
    // Skip deep link setup on web - use browser URL instead
    if (kIsWeb) {
      // On web, use URL parameters from window.location
      // Deep links are handled via GoRouter's initial route parsing
      return;
    }

    // Handle initial link (cold start) - mobile only
    // Note: uni_links package should be imported conditionally for mobile
    // For now, we skip uni_links initialization as it's not web-compatible
    // TODO: Implement platform-specific deep link handling
  }

  /// Parse URI to DeepLink
  DeepLink _parseUri(Uri uri) {
    final path = uri.path;
    final queryParams = uri.queryParameters;
    
    // Handle hustlex:// scheme
    // Handle https://hustlex.ng/... universal links
    
    DeepLinkType type = DeepLinkType.unknown;
    String? id;
    
    final pathSegments = path.split('/').where((s) => s.isNotEmpty).toList();
    
    if (pathSegments.isEmpty) {
      return DeepLink(
        type: DeepLinkType.unknown,
        params: queryParams,
        rawUri: uri.toString(),
      );
    }

    switch (pathSegments[0]) {
      case 'gig':
      case 'gigs':
        type = DeepLinkType.gig;
        if (pathSegments.length > 1) {
          id = pathSegments[1];
        }
        break;
        
      case 'circle':
      case 'circles':
      case 'savings':
        type = DeepLinkType.circle;
        if (pathSegments.length > 1) {
          id = pathSegments[1];
        }
        break;
        
      case 'wallet':
        type = DeepLinkType.wallet;
        break;
        
      case 'transaction':
      case 'transactions':
        type = DeepLinkType.transaction;
        if (pathSegments.length > 1) {
          id = pathSegments[1];
        }
        break;
        
      case 'profile':
      case 'user':
        type = DeepLinkType.profile;
        if (pathSegments.length > 1) {
          id = pathSegments[1];
        }
        break;
        
      case 'referral':
      case 'ref':
      case 'invite':
        type = DeepLinkType.referral;
        id = queryParams['code'] ?? (pathSegments.length > 1 ? pathSegments[1] : null);
        break;
        
      case 'promo':
      case 'promotion':
        type = DeepLinkType.promo;
        id = queryParams['code'] ?? (pathSegments.length > 1 ? pathSegments[1] : null);
        break;
        
      case 'verify':
      case 'verification':
        type = DeepLinkType.verify;
        id = queryParams['token'];
        break;
        
      case 'reset-password':
      case 'reset':
        type = DeepLinkType.resetPassword;
        id = queryParams['token'];
        break;
    }

    return DeepLink(
      type: type,
      id: id,
      params: queryParams,
      rawUri: uri.toString(),
    );
  }

  /// Get route for deep link
  String getRouteForDeepLink(DeepLink link) {
    switch (link.type) {
      case DeepLinkType.gig:
        if (link.id != null) {
          return '/gigs/${link.id}';
        }
        return '/gigs';
        
      case DeepLinkType.circle:
        if (link.id != null) {
          return '/savings/circle/${link.id}';
        }
        return '/savings';
        
      case DeepLinkType.wallet:
        return '/wallet';
        
      case DeepLinkType.transaction:
        if (link.id != null) {
          return '/wallet/transaction/${link.id}';
        }
        return '/wallet/transactions';
        
      case DeepLinkType.profile:
        return '/profile';
        
      case DeepLinkType.referral:
        // Store referral code and redirect to registration
        return '/auth/register?ref=${link.id ?? ''}';
        
      case DeepLinkType.promo:
        // Handle promo code
        return '/wallet?promo=${link.id ?? ''}';
        
      case DeepLinkType.verify:
        return '/auth/verify?token=${link.id ?? ''}';
        
      case DeepLinkType.resetPassword:
        return '/auth/reset-password?token=${link.id ?? ''}';
        
      case DeepLinkType.unknown:
        return '/';
    }
  }

  /// Handle deep link with GoRouter
  void handleDeepLink(DeepLink link, GoRouter router) {
    final route = getRouteForDeepLink(link);
    router.go(route);
  }

  /// Generate share link for gig
  String generateGigShareLink(String gigId) {
    return 'https://hustlex.ng/gig/$gigId';
  }

  /// Generate share link for circle
  String generateCircleShareLink(String circleId) {
    return 'https://hustlex.ng/circle/$circleId';
  }

  /// Generate referral link
  String generateReferralLink(String referralCode) {
    return 'https://hustlex.ng/invite/$referralCode';
  }

  /// Generate profile share link
  String generateProfileShareLink(String userId) {
    return 'https://hustlex.ng/user/$userId';
  }

  /// Clean up
  void dispose() {
    _linkSubscription?.cancel();
    _linkSubscription = null;
  }
}

/// Deep link service provider
final deepLinkServiceProvider = Provider<DeepLinkService>((ref) {
  final service = DeepLinkService();
  ref.onDispose(() => service.dispose());
  return service;
});

/// Initial deep link provider
final initialDeepLinkProvider = FutureProvider<DeepLink?>((ref) async {
  final service = ref.watch(deepLinkServiceProvider);
  await service.initialize();
  return service.initialLink;
});

/// Deep link stream provider
final deepLinkStreamProvider = StreamProvider<DeepLink>((ref) {
  final service = ref.watch(deepLinkServiceProvider);
  final controller = StreamController<DeepLink>();
  
  service.onDeepLink = (link) {
    controller.add(link);
  };
  
  ref.onDispose(() {
    service.onDeepLink = null;
    controller.close();
  });
  
  return controller.stream;
});

/// Deep link handler mixin for widgets
mixin DeepLinkHandler {
  /// Handle incoming deep link
  void handleDeepLink(DeepLink link, GoRouter router) {
    // Special handling for certain link types
    switch (link.type) {
      case DeepLinkType.referral:
        _handleReferralLink(link, router);
        break;
      case DeepLinkType.promo:
        _handlePromoLink(link, router);
        break;
      default:
        final route = DeepLinkService().getRouteForDeepLink(link);
        router.go(route);
    }
  }

  void _handleReferralLink(DeepLink link, GoRouter router) {
    // Store referral code for later use
    if (link.id != null) {
      // Store in secure storage or preferences
      // Then navigate to registration
      router.go('/auth/register?ref=${link.id}');
    }
  }

  void _handlePromoLink(DeepLink link, GoRouter router) {
    // Store promo code for later use
    if (link.id != null) {
      // Show promo modal or navigate with promo
      router.go('/wallet?promo=${link.id}');
    }
  }
}
