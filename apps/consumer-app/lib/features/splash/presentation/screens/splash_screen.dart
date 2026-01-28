import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/constants/app_colors.dart';
import '../../core/constants/app_constants.dart';
import '../../features/auth/presentation/providers/auth_provider.dart';

/// Splash screen with app initialization
class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  bool _showLoading = false;
  String _loadingMessage = 'Starting...';
  double _progress = 0;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      duration: const Duration(milliseconds: 800),
      vsync: this,
    );
    _startInitialization();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _startInitialization() async {
    // Start logo animation
    _controller.forward();

    // Show loading after animation
    await Future.delayed(const Duration(milliseconds: 1000));
    
    if (!mounted) return;
    setState(() => _showLoading = true);

    // Simulate initialization stages
    await _initializeApp();
  }

  Future<void> _initializeApp() async {
    if (!mounted) return;

    // Stage 1: Check storage
    _updateProgress(0.2, 'Loading preferences...');
    await Future.delayed(const Duration(milliseconds: 300));

    // Stage 2: Initialize services
    _updateProgress(0.4, 'Setting up services...');
    await Future.delayed(const Duration(milliseconds: 300));

    // Stage 3: Check authentication
    _updateProgress(0.6, 'Checking authentication...');
    await Future.delayed(const Duration(milliseconds: 300));

    // Stage 4: Load user data
    _updateProgress(0.8, 'Loading your data...');
    await Future.delayed(const Duration(milliseconds: 300));

    // Complete
    _updateProgress(1.0, 'Ready!');
    await Future.delayed(const Duration(milliseconds: 200));

    // Navigate based on auth state
    if (!mounted) return;
    _navigateToNextScreen();
  }

  void _updateProgress(double progress, String message) {
    if (!mounted) return;
    setState(() {
      _progress = progress;
      _loadingMessage = message;
    });
  }

  void _navigateToNextScreen() {
    final authState = ref.read(authProvider);

    if (authState.isAuthenticated) {
      // User is logged in, go to home
      context.go('/');
    } else if (authState.hasSeenOnboarding) {
      // User has seen onboarding, go to login
      context.go('/login');
    } else {
      // First time user, show onboarding
      context.go('/onboarding');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.primary,
      body: SafeArea(
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              // Logo
              _buildLogo(),
              
              const SizedBox(height: 32),
              
              // App name
              _buildAppName(),
              
              const SizedBox(height: 48),
              
              // Loading indicator
              if (_showLoading) ...[
                _buildLoadingIndicator(),
                const SizedBox(height: 16),
                _buildLoadingMessage(),
              ],
              
              const Spacer(),
              
              // Bottom tagline
              _buildTagline(),
              
              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildLogo() {
    return Container(
      width: 120,
      height: 120,
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(30),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.2),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Center(
        child: Text(
          'H',
          style: TextStyle(
            fontSize: 64,
            fontWeight: FontWeight.bold,
            color: AppColors.primary,
            height: 1,
          ),
        ),
      ),
    )
        .animate(controller: _controller)
        .scale(
          begin: const Offset(0.5, 0.5),
          end: const Offset(1, 1),
          curve: Curves.elasticOut,
        )
        .fadeIn();
  }

  Widget _buildAppName() {
    return Column(
      children: [
        Text(
          AppConstants.appName,
          style: const TextStyle(
            fontSize: 36,
            fontWeight: FontWeight.bold,
            color: Colors.white,
            letterSpacing: 2,
          ),
        )
            .animate(controller: _controller)
            .fadeIn(delay: const Duration(milliseconds: 200))
            .slideY(begin: 0.3, end: 0),
        const SizedBox(height: 8),
        Text(
          'Hustle Smart, Save Smarter',
          style: TextStyle(
            fontSize: 16,
            color: Colors.white.withOpacity(0.8),
            fontWeight: FontWeight.w500,
          ),
        )
            .animate(controller: _controller)
            .fadeIn(delay: const Duration(milliseconds: 400))
            .slideY(begin: 0.3, end: 0),
      ],
    );
  }

  Widget _buildLoadingIndicator() {
    return SizedBox(
      width: 200,
      child: Column(
        children: [
          // Progress bar
          ClipRRect(
            borderRadius: BorderRadius.circular(4),
            child: LinearProgressIndicator(
              value: _progress,
              backgroundColor: Colors.white.withOpacity(0.3),
              valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
              minHeight: 4,
            ),
          ),
          const SizedBox(height: 8),
          // Percentage
          Text(
            '${(_progress * 100).toInt()}%',
            style: TextStyle(
              fontSize: 12,
              color: Colors.white.withOpacity(0.8),
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    ).animate().fadeIn();
  }

  Widget _buildLoadingMessage() {
    return Text(
      _loadingMessage,
      style: TextStyle(
        fontSize: 14,
        color: Colors.white.withOpacity(0.8),
      ),
    ).animate().fadeIn();
  }

  Widget _buildTagline() {
    return Column(
      children: [
        Text(
          'Your gateway to',
          style: TextStyle(
            fontSize: 12,
            color: Colors.white.withOpacity(0.6),
          ),
        ),
        const SizedBox(height: 4),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            _buildFeatureChip('Gigs'),
            _buildFeatureDot(),
            _buildFeatureChip('Savings'),
            _buildFeatureDot(),
            _buildFeatureChip('Credit'),
          ],
        ),
      ],
    )
        .animate(controller: _controller)
        .fadeIn(delay: const Duration(milliseconds: 600));
  }

  Widget _buildFeatureChip(String label) {
    return Text(
      label,
      style: TextStyle(
        fontSize: 12,
        color: Colors.white.withOpacity(0.9),
        fontWeight: FontWeight.w600,
      ),
    );
  }

  Widget _buildFeatureDot() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 8),
      child: Container(
        width: 4,
        height: 4,
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.5),
          shape: BoxShape.circle,
        ),
      ),
    );
  }
}

/// Minimal splash for quick transitions
class MinimalSplash extends StatelessWidget {
  const MinimalSplash({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.primary,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
              ),
              child: Center(
                child: Text(
                  'H',
                  style: TextStyle(
                    fontSize: 42,
                    fontWeight: FontWeight.bold,
                    color: AppColors.primary,
                  ),
                ),
              ),
            ),
            const SizedBox(height: 24),
            const CircularProgressIndicator(
              valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
              strokeWidth: 2,
            ),
          ],
        ),
      ),
    );
  }
}
