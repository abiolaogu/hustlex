import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../shared/widgets/buttons.dart';

class ChangePinScreen extends ConsumerStatefulWidget {
  const ChangePinScreen({super.key});

  @override
  ConsumerState<ChangePinScreen> createState() => _ChangePinScreenState();
}

class _ChangePinScreenState extends ConsumerState<ChangePinScreen> {
  final _currentPinController = TextEditingController();
  final _newPinController = TextEditingController();
  final _confirmPinController = TextEditingController();
  
  final _currentPinFocus = FocusNode();
  final _newPinFocus = FocusNode();
  final _confirmPinFocus = FocusNode();
  
  int _currentStep = 0; // 0: current PIN, 1: new PIN, 2: confirm PIN
  bool _isVerifying = false;
  bool _isChanging = false;
  String? _errorMessage;

  static const int _pinLength = 4;

  @override
  void initState() {
    super.initState();
    _currentPinFocus.requestFocus();
  }

  @override
  void dispose() {
    _currentPinController.dispose();
    _newPinController.dispose();
    _confirmPinController.dispose();
    _currentPinFocus.dispose();
    _newPinFocus.dispose();
    _confirmPinFocus.dispose();
    super.dispose();
  }

  Future<void> _verifyCurrentPin() async {
    if (_currentPinController.text.length != _pinLength) return;

    setState(() {
      _isVerifying = true;
      _errorMessage = null;
    });

    try {
      // TODO: Verify current PIN via API
      await Future.delayed(const Duration(seconds: 1));
      
      // Mock verification - assume PIN is "1234"
      if (_currentPinController.text == '1234') {
        setState(() {
          _currentStep = 1;
          _isVerifying = false;
        });
        _newPinFocus.requestFocus();
      } else {
        setState(() {
          _errorMessage = 'Incorrect PIN. Please try again.';
          _isVerifying = false;
        });
        _currentPinController.clear();
      }
    } catch (e) {
      setState(() {
        _errorMessage = 'Verification failed. Please try again.';
        _isVerifying = false;
      });
    }
  }

  void _onNewPinComplete() {
    if (_newPinController.text.length != _pinLength) return;

    // Check for weak PINs
    final weakPins = ['0000', '1111', '2222', '3333', '4444', '5555', '6666', '7777', '8888', '9999', '1234', '4321'];
    if (weakPins.contains(_newPinController.text)) {
      setState(() {
        _errorMessage = 'Please choose a stronger PIN';
      });
      _newPinController.clear();
      return;
    }

    // Check if new PIN is same as current
    if (_newPinController.text == _currentPinController.text) {
      setState(() {
        _errorMessage = 'New PIN must be different from current PIN';
      });
      _newPinController.clear();
      return;
    }

    setState(() {
      _currentStep = 2;
      _errorMessage = null;
    });
    _confirmPinFocus.requestFocus();
  }

  Future<void> _onConfirmPinComplete() async {
    if (_confirmPinController.text.length != _pinLength) return;

    // Check if PINs match
    if (_confirmPinController.text != _newPinController.text) {
      setState(() {
        _errorMessage = 'PINs do not match. Please try again.';
      });
      _confirmPinController.clear();
      return;
    }

    setState(() {
      _isChanging = true;
      _errorMessage = null;
    });

    try {
      // TODO: Change PIN via API
      await Future.delayed(const Duration(seconds: 2));

      if (mounted) {
        _showSuccessDialog();
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _errorMessage = 'Failed to change PIN. Please try again.';
          _isChanging = false;
        });
      }
    }
  }

  void _showSuccessDialog() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: AppColors.success.withOpacity(0.1),
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.check_circle,
                color: AppColors.success,
                size: 48,
              ),
            ),
            const SizedBox(height: 24),
            Text(
              'PIN Changed!',
              style: AppTypography.headlineSmall.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Your transaction PIN has been changed successfully.',
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            PrimaryButton(
              text: 'Done',
              onPressed: () {
                context.pop();
                context.pop();
              },
            ),
          ],
        ),
      ),
    );
  }

  void _goBack() {
    if (_currentStep > 0) {
      setState(() {
        _currentStep--;
        _errorMessage = null;
        if (_currentStep == 0) {
          _currentPinController.clear();
          _newPinController.clear();
          _confirmPinController.clear();
          _currentPinFocus.requestFocus();
        } else if (_currentStep == 1) {
          _newPinController.clear();
          _confirmPinController.clear();
          _newPinFocus.requestFocus();
        }
      });
    } else {
      context.pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Change PIN'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: _goBack,
        ),
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            children: [
              _buildProgressIndicator(),
              const SizedBox(height: 48),
              _buildStepContent(),
              if (_errorMessage != null) ...[
                const SizedBox(height: 24),
                _buildErrorMessage(),
              ],
              const Spacer(),
              _buildSecurityTip(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildProgressIndicator() {
    return Row(
      children: List.generate(3, (index) {
        final isCompleted = index < _currentStep;
        final isCurrent = index == _currentStep;
        return Expanded(
          child: Container(
            margin: EdgeInsets.only(right: index < 2 ? 8 : 0),
            height: 4,
            decoration: BoxDecoration(
              color: isCompleted
                  ? AppColors.success
                  : isCurrent
                      ? AppColors.primary
                      : AppColors.border,
              borderRadius: BorderRadius.circular(2),
            ),
          ),
        );
      }),
    );
  }

  Widget _buildStepContent() {
    switch (_currentStep) {
      case 0:
        return _buildCurrentPinStep();
      case 1:
        return _buildNewPinStep();
      case 2:
        return _buildConfirmPinStep();
      default:
        return const SizedBox.shrink();
    }
  }

  Widget _buildCurrentPinStep() {
    return Column(
      children: [
        Container(
          width: 80,
          height: 80,
          decoration: BoxDecoration(
            color: AppColors.primary.withOpacity(0.1),
            shape: BoxShape.circle,
          ),
          child: Icon(
            Icons.lock_outline,
            color: AppColors.primary,
            size: 40,
          ),
        ),
        const SizedBox(height: 24),
        Text(
          'Enter Current PIN',
          style: AppTypography.headlineSmall.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Enter your current 4-digit PIN to continue',
          style: AppTypography.bodyMedium.copyWith(
            color: AppColors.textSecondary,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 32),
        _buildPinInput(
          controller: _currentPinController,
          focusNode: _currentPinFocus,
          onComplete: _verifyCurrentPin,
          isLoading: _isVerifying,
        ),
        const SizedBox(height: 16),
        TextButton(
          onPressed: () {
            // TODO: Implement forgot PIN flow
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Forgot PIN flow coming soon'),
              ),
            );
          },
          child: Text(
            'Forgot PIN?',
            style: AppTypography.labelMedium.copyWith(
              color: AppColors.primary,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildNewPinStep() {
    return Column(
      children: [
        Container(
          width: 80,
          height: 80,
          decoration: BoxDecoration(
            color: AppColors.secondary.withOpacity(0.1),
            shape: BoxShape.circle,
          ),
          child: Icon(
            Icons.vpn_key_outlined,
            color: AppColors.secondary,
            size: 40,
          ),
        ),
        const SizedBox(height: 24),
        Text(
          'Create New PIN',
          style: AppTypography.headlineSmall.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Enter a new 4-digit PIN you\'ll remember',
          style: AppTypography.bodyMedium.copyWith(
            color: AppColors.textSecondary,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 32),
        _buildPinInput(
          controller: _newPinController,
          focusNode: _newPinFocus,
          onComplete: _onNewPinComplete,
          isLoading: false,
        ),
      ],
    );
  }

  Widget _buildConfirmPinStep() {
    return Column(
      children: [
        Container(
          width: 80,
          height: 80,
          decoration: BoxDecoration(
            color: AppColors.success.withOpacity(0.1),
            shape: BoxShape.circle,
          ),
          child: Icon(
            Icons.check_circle_outline,
            color: AppColors.success,
            size: 40,
          ),
        ),
        const SizedBox(height: 24),
        Text(
          'Confirm New PIN',
          style: AppTypography.headlineSmall.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Re-enter your new PIN to confirm',
          style: AppTypography.bodyMedium.copyWith(
            color: AppColors.textSecondary,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 32),
        _buildPinInput(
          controller: _confirmPinController,
          focusNode: _confirmPinFocus,
          onComplete: _onConfirmPinComplete,
          isLoading: _isChanging,
        ),
      ],
    );
  }

  Widget _buildPinInput({
    required TextEditingController controller,
    required FocusNode focusNode,
    required VoidCallback onComplete,
    required bool isLoading,
  }) {
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: List.generate(_pinLength, (index) {
            final isFilled = controller.text.length > index;
            return Container(
              margin: const EdgeInsets.symmetric(horizontal: 8),
              width: 56,
              height: 56,
              decoration: BoxDecoration(
                color: isFilled
                    ? AppColors.primary.withOpacity(0.1)
                    : AppColors.surfaceVariant,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: isFilled ? AppColors.primary : AppColors.border,
                  width: 2,
                ),
              ),
              child: Center(
                child: isFilled
                    ? Container(
                        width: 16,
                        height: 16,
                        decoration: const BoxDecoration(
                          color: AppColors.primary,
                          shape: BoxShape.circle,
                        ),
                      )
                    : null,
              ),
            );
          }),
        ),
        const SizedBox(height: 16),
        if (isLoading)
          const SizedBox(
            width: 24,
            height: 24,
            child: CircularProgressIndicator(strokeWidth: 2),
          )
        else
          Opacity(
            opacity: 0,
            child: TextField(
              controller: controller,
              focusNode: focusNode,
              keyboardType: TextInputType.number,
              maxLength: _pinLength,
              autofocus: true,
              inputFormatters: [FilteringTextInputFormatter.digitsOnly],
              decoration: const InputDecoration(
                counterText: '',
                border: InputBorder.none,
              ),
              onChanged: (value) {
                setState(() {});
                if (value.length == _pinLength) {
                  onComplete();
                }
              },
            ),
          ),
      ],
    );
  }

  Widget _buildErrorMessage() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.error.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Icon(
            Icons.error_outline,
            color: AppColors.error,
            size: 20,
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              _errorMessage!,
              style: AppTypography.bodySmall.copyWith(
                color: AppColors.error,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSecurityTip() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.warning.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(
            Icons.lightbulb_outline,
            color: AppColors.warning,
            size: 20,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Security Tip',
                  style: AppTypography.labelMedium.copyWith(
                    fontWeight: FontWeight.w600,
                    color: AppColors.warning,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  'Never share your PIN with anyone. HustleX will never ask for your PIN via call, SMS, or email.',
                  style: AppTypography.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
