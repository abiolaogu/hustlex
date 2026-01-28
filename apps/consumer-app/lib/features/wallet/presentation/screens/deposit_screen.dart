import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/utils/helpers.dart';
import '../../../../core/widgets/buttons.dart';
import '../../../../core/widgets/inputs.dart';
import '../../data/models/wallet_model.dart';
import '../providers/wallet_provider.dart';

class DepositScreen extends ConsumerStatefulWidget {
  const DepositScreen({super.key});

  @override
  ConsumerState<DepositScreen> createState() => _DepositScreenState();
}

class _DepositScreenState extends ConsumerState<DepositScreen> {
  final _formKey = GlobalKey<FormState>();
  final _amountController = TextEditingController();
  PaymentMethod _selectedMethod = PaymentMethod.card;

  final List<int> _quickAmounts = [1000, 2000, 5000, 10000, 20000, 50000];

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  double get _amount {
    final text = _amountController.text.replaceAll(',', '');
    return double.tryParse(text) ?? 0;
  }

  Future<void> _processDeposit() async {
    if (!_formKey.currentState!.validate()) return;

    final depositNotifier = ref.read(depositProvider.notifier);
    
    final request = DepositRequest(
      amount: _amount,
      method: _selectedMethod,
    );

    final success = await depositNotifier.initializeDeposit(request);

    if (success && mounted) {
      final initResponse = ref.read(depositProvider).initResponse;
      if (initResponse != null) {
        // Open Paystack payment page
        _openPaystackPayment(initResponse);
      }
    } else if (mounted) {
      final error = ref.read(depositProvider).error;
      _showError(error ?? 'Failed to initialize deposit');
    }
  }

  void _openPaystackPayment(PaystackInitResponse response) {
    // In a real implementation, use paystack_flutter package
    // For now, show a success dialog
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: const Row(
          children: [
            Icon(Icons.payment, color: AppColors.primary),
            SizedBox(width: 12),
            Text('Payment'),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Amount: ${CurrencyUtils.formatNaira(_amount)}'),
            const SizedBox(height: 8),
            Text('Reference: ${response.reference}'),
            const SizedBox(height: 16),
            const Text(
              'In production, this would open the Paystack payment page.',
              style: TextStyle(color: AppColors.textSecondary, fontSize: 12),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              _verifyPayment(response.reference);
            },
            child: const Text('Simulate Success'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              _showError('Payment cancelled');
            },
            child: const Text('Cancel'),
          ),
        ],
      ),
    );
  }

  Future<void> _verifyPayment(String reference) async {
    final depositNotifier = ref.read(depositProvider.notifier);
    final success = await depositNotifier.verifyDeposit(reference);

    if (success && mounted) {
      // Refresh wallet balance
      ref.read(walletProvider.notifier).refreshBalance();
      _showSuccess();
    } else if (mounted) {
      final error = ref.read(depositProvider).error;
      _showError(error ?? 'Payment verification failed');
    }
  }

  void _showSuccess() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const CircleAvatar(
              radius: 40,
              backgroundColor: AppColors.success,
              child: Icon(Icons.check, color: Colors.white, size: 40),
            ),
            const SizedBox(height: 24),
            const Text(
              'Deposit Successful!',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '${CurrencyUtils.formatNaira(_amount)} has been added to your wallet',
              textAlign: TextAlign.center,
              style: const TextStyle(color: AppColors.textSecondary),
            ),
          ],
        ),
        actions: [
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: () {
                Navigator.pop(context);
                context.pop();
              },
              child: const Text('Done'),
            ),
          ),
        ],
      ),
    );
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: AppColors.error,
        behavior: SnackBarBehavior.floating,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final depositState = ref.watch(depositProvider);
    final walletState = ref.watch(walletProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Add Money'),
        centerTitle: true,
      ),
      body: SafeArea(
        child: Form(
          key: _formKey,
          child: ListView(
            padding: const EdgeInsets.all(20),
            children: [
              // Current balance card
              Container(
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  gradient: AppColors.primaryGradient,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Available Balance',
                      style: AppTypography.bodySmall.copyWith(
                        color: Colors.white.withOpacity(0.8),
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      CurrencyUtils.formatNaira(walletState.availableBalance),
                      style: AppTypography.headlineMedium.copyWith(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 32),

              // Amount input
              Text(
                'Enter Amount',
                style: AppTypography.titleMedium.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 12),
              AmountTextField(
                controller: _amountController,
                hintText: '0.00',
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter amount';
                  }
                  final amount = double.tryParse(value.replaceAll(',', ''));
                  if (amount == null || amount < 100) {
                    return 'Minimum deposit is ₦100';
                  }
                  if (amount > 5000000) {
                    return 'Maximum deposit is ₦5,000,000';
                  }
                  return null;
                },
              ),

              const SizedBox(height: 16),

              // Quick amounts
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: _quickAmounts.map((amount) {
                  return ActionChip(
                    label: Text(CurrencyUtils.formatNairaWhole(amount.toDouble())),
                    onPressed: () {
                      _amountController.text = NumberFormat('#,###').format(amount);
                    },
                    backgroundColor: AppColors.surface,
                    side: const BorderSide(color: AppColors.border),
                  );
                }).toList(),
              ),

              const SizedBox(height: 32),

              // Payment method
              Text(
                'Payment Method',
                style: AppTypography.titleMedium.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 12),

              _PaymentMethodTile(
                icon: Icons.credit_card,
                title: 'Debit/Credit Card',
                subtitle: 'Visa, Mastercard, Verve',
                isSelected: _selectedMethod == PaymentMethod.card,
                onTap: () => setState(() => _selectedMethod = PaymentMethod.card),
              ),
              const SizedBox(height: 8),
              _PaymentMethodTile(
                icon: Icons.account_balance,
                title: 'Bank Transfer',
                subtitle: 'Direct bank transfer',
                isSelected: _selectedMethod == PaymentMethod.bankTransfer,
                onTap: () => setState(() => _selectedMethod = PaymentMethod.bankTransfer),
              ),
              const SizedBox(height: 8),
              _PaymentMethodTile(
                icon: Icons.qr_code,
                title: 'USSD',
                subtitle: 'Pay with USSD code',
                isSelected: _selectedMethod == PaymentMethod.ussd,
                onTap: () => setState(() => _selectedMethod = PaymentMethod.ussd),
              ),

              const SizedBox(height: 32),

              // Info card
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: AppColors.info.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: AppColors.info.withOpacity(0.3)),
                ),
                child: Row(
                  children: [
                    const Icon(Icons.info_outline, color: AppColors.info),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Text(
                        'Deposits are instant and secure. Powered by Paystack.',
                        style: AppTypography.bodySmall.copyWith(
                          color: AppColors.info,
                        ),
                      ),
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 32),

              // Continue button
              PrimaryButton(
                text: 'Continue',
                onPressed: _processDeposit,
                isLoading: depositState.isProcessing,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _PaymentMethodTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final bool isSelected;
  final VoidCallback onTap;

  const _PaymentMethodTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected ? AppColors.primary.withOpacity(0.1) : AppColors.surface,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isSelected ? AppColors.primary : AppColors.border,
            width: isSelected ? 2 : 1,
          ),
        ),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: isSelected
                    ? AppColors.primary.withOpacity(0.2)
                    : AppColors.backgroundDark,
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                icon,
                color: isSelected ? AppColors.primary : AppColors.textSecondary,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: AppTypography.bodyMedium.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    subtitle,
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
            Radio<bool>(
              value: true,
              groupValue: isSelected,
              onChanged: (_) => onTap(),
              activeColor: AppColors.primary,
            ),
          ],
        ),
      ),
    );
  }
}

// Temporary NumberFormat class - in real app, use intl package
class NumberFormat {
  final String pattern;
  NumberFormat(this.pattern);
  
  String format(num number) {
    return number.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (Match m) => '${m[1]},',
    );
  }
}
