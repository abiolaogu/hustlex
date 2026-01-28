import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:pin_code_fields/pin_code_fields.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/providers/auth_provider.dart';
import '../../../../core/exceptions/api_exception.dart';
import '../../../../router/app_router.dart';

class PinSetupScreen extends ConsumerStatefulWidget {
  const PinSetupScreen({super.key});

  @override
  ConsumerState<PinSetupScreen> createState() => _PinSetupScreenState();
}

class _PinSetupScreenState extends ConsumerState<PinSetupScreen> {
  final _pinController = TextEditingController();
  final _confirmPinController = TextEditingController();
  
  bool _isLoading = false;
  bool _isConfirmStep = false;
  String? _errorMessage;
  String _firstPin = '';

  @override
  void dispose() {
    _pinController.dispose();
    _confirmPinController.dispose();
    super.dispose();
  }

  void _onPinEntered(String pin) {
    if (!_isConfirmStep) {
      // First PIN entry
      setState(() {
        _firstPin = pin;
        _isConfirmStep = true;
        _errorMessage = null;
      });
    } else {
      // Confirm PIN entry
      if (pin == _firstPin) {
        _setPin(pin);
      } else {
        setState(() {
          _errorMessage = 'PINs do not match. Please try again.';
          _confirmPinController.clear();
        });
      }
    }
  }

  Future<void> _setPin(String pin) async {
    setState(() {
      _isLoading = true;
      _errorMessage = null;
    });

    try {
      await ref.read(authStateProvider.notifier).setTransactionPin(pin);

      if (mounted) {
        // Show success message and navigate
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Transaction PIN set successfully!'),
            backgroundColor: AppColors.success,
          ),
        );
        context.go(AppRoutes.home);
      }
    } on ApiException catch (e) {
      setState(() {
        _errorMessage = e.message;
        _confirmPinController.clear();
      });
    } catch (e) {
      setState(() {
        _errorMessage = 'Failed to set PIN. Please try again.';
        _confirmPinController.clear();
      });
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  void _goBack() {
    if (_isConfirmStep) {
      setState(() {
        _isConfirmStep = false;
        _firstPin = '';
        _errorMessage = null;
        _pinController.clear();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: _isConfirmStep
            ? IconButton(
                icon: const Icon(Icons.arrow_back_ios_new_rounded),
                onPressed: _goBack,
              )
            : null,
        automaticallyImplyLeading: false,
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const SizedBox(height: 20),

              // Icon
              Center(
                child: Container(
                  width: 80,
                  height: 80,
                  decoration: BoxDecoration(
                    color: AppColors.primaryLight.withOpacity(0.2),
                    shape: BoxShape.circle,
                  ),
                  child: Icon(
                    _isConfirmStep
                        ? Icons.verified_user_outlined
                        : Icons.lock_outline_rounded,
                    size: 40,
                    color: AppColors.primary,
                  ),
                ),
              ),

              const SizedBox(height: 32),

              // Header
              Center(
                child: Column(
                  children: [
                    Text(
                      _isConfirmStep
                          ? 'Confirm your PIN'
                          : 'Set your transaction PIN',
                      style: AppTypography.headlineSmall.copyWith(
                        color: AppColors.textPrimary,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      _isConfirmStep
                          ? 'Enter your PIN again to confirm'
                          : 'This PIN will be used to authorize transactions',
                      style: AppTypography.bodyMedium.copyWith(
                        color: AppColors.textSecondary,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 48),

              // PIN Input
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 40),
                child: PinCodeTextField(
                  appContext: context,
                  length: AppConstants.pinLength,
                  controller:
                      _isConfirmStep ? _confirmPinController : _pinController,
                  autoFocus: true,
                  obscureText: true,
                  obscuringCharacter: '●',
                  keyboardType: TextInputType.number,
                  animationType: AnimationType.scale,
                  enableActiveFill: true,
                  cursorColor: AppColors.primary,
                  textStyle: AppTypography.headlineMedium.copyWith(
                    color: AppColors.textPrimary,
                  ),
                  pinTheme: PinTheme(
                    shape: PinCodeFieldShape.box,
                    borderRadius: BorderRadius.circular(16),
                    fieldHeight: 64,
                    fieldWidth: 56,
                    activeColor: AppColors.primary,
                    selectedColor: AppColors.primary,
                    inactiveColor: AppColors.border,
                    activeFillColor: AppColors.white,
                    selectedFillColor: AppColors.white,
                    inactiveFillColor: AppColors.surfaceVariant,
                    errorBorderColor: AppColors.error,
                  ),
                  animationDuration: AppConstants.shortAnimationDuration,
                  onCompleted: _onPinEntered,
                  onChanged: (_) {
                    if (_errorMessage != null) {
                      setState(() => _errorMessage = null);
                    }
                  },
                ),
              ),

              // Error message
              if (_errorMessage != null) ...[
                const SizedBox(height: 16),
                Center(
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 12,
                    ),
                    decoration: BoxDecoration(
                      color: AppColors.errorLight,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(
                          Icons.error_outline,
                          color: AppColors.error,
                          size: 20,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          _errorMessage!,
                          style: AppTypography.bodySmall.copyWith(
                            color: AppColors.error,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],

              // Loading indicator
              if (_isLoading) ...[
                const SizedBox(height: 24),
                const Center(
                  child: CircularProgressIndicator(),
                ),
              ],

              const Spacer(),

              // Security tips
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: AppColors.warningLight,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(
                          Icons.security_rounded,
                          color: AppColors.warning,
                          size: 20,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          'Security Tips',
                          style: AppTypography.titleSmall.copyWith(
                            color: AppColors.textPrimary,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Text(
                      '• Don\'t use easily guessable PINs like 1234 or 0000\n'
                      '• Don\'t share your PIN with anyone\n'
                      '• Change your PIN regularly for security',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                        height: 1.5,
                      ),
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 24),
            ],
          ),
        ),
      ),
    );
  }
}
