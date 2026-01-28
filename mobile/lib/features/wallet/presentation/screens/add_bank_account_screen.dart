import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/widgets/buttons.dart';
import '../../../../core/widgets/inputs.dart';
import '../../../../core/widgets/loaders.dart';
import '../providers/wallet_provider.dart';

/// =============================================================================
/// ADD BANK ACCOUNT SCREEN
/// =============================================================================

class AddBankAccountScreen extends ConsumerStatefulWidget {
  const AddBankAccountScreen({super.key});

  @override
  ConsumerState<AddBankAccountScreen> createState() => _AddBankAccountScreenState();
}

class _AddBankAccountScreenState extends ConsumerState<AddBankAccountScreen> {
  final _formKey = GlobalKey<FormState>();
  final _accountNumberController = TextEditingController();
  final _accountNumberFocus = FocusNode();

  String? _selectedBankCode;
  String? _selectedBankName;
  String? _verifiedAccountName;
  bool _isVerifying = false;
  bool _isSubmitting = false;
  String? _error;

  @override
  void dispose() {
    _accountNumberController.dispose();
    _accountNumberFocus.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final banksAsync = ref.watch(banksProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Add Bank Account'),
        backgroundColor: AppColors.background,
        elevation: 0,
      ),
      body: banksAsync.when(
        loading: () => const Center(child: AppLoader()),
        error: (error, _) => _buildError(error.toString()),
        data: (banks) => _buildForm(banks),
      ),
    );
  }

  Widget _buildError(String message) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.error_outline,
              size: 64,
              color: AppColors.error,
            ),
            const SizedBox(height: 16),
            Text(
              'Failed to load banks',
              style: AppTypography.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              message,
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            AppButton(
              text: 'Try Again',
              onPressed: () => ref.refresh(banksProvider),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildForm(List<Bank> banks) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Info card
            _buildInfoCard(),

            const SizedBox(height: 24),

            // Bank selection
            Text(
              'Select Bank',
              style: AppTypography.labelLarge,
            ),
            const SizedBox(height: 8),
            _buildBankDropdown(banks),

            const SizedBox(height: 20),

            // Account number
            Text(
              'Account Number',
              style: AppTypography.labelLarge,
            ),
            const SizedBox(height: 8),
            AppTextField(
              controller: _accountNumberController,
              focusNode: _accountNumberFocus,
              hintText: 'Enter 10-digit account number',
              keyboardType: TextInputType.number,
              inputFormatters: [
                FilteringTextInputFormatter.digitsOnly,
                LengthLimitingTextInputFormatter(10),
              ],
              onChanged: (value) {
                // Reset verification when account number changes
                if (_verifiedAccountName != null) {
                  setState(() {
                    _verifiedAccountName = null;
                    _error = null;
                  });
                }

                // Auto-verify when 10 digits entered
                if (value.length == 10 && _selectedBankCode != null) {
                  _verifyAccount();
                }
              },
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter account number';
                }
                if (value.length != 10) {
                  return 'Account number must be 10 digits';
                }
                return null;
              },
            ),

            const SizedBox(height: 20),

            // Verification result
            if (_isVerifying) _buildVerifyingIndicator(),
            if (_verifiedAccountName != null) _buildVerifiedAccountCard(),
            if (_error != null) _buildErrorMessage(),

            const SizedBox(height: 32),

            // Submit button
            AppButton(
              text: _isSubmitting ? 'Adding...' : 'Add Bank Account',
              onPressed: _canSubmit ? _submitAccount : null,
              isLoading: _isSubmitting,
            ),

            const SizedBox(height: 16),

            // Security note
            _buildSecurityNote(),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.info.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppColors.info.withOpacity(0.3),
        ),
      ),
      child: Row(
        children: [
          const Icon(
            Icons.info_outline,
            color: AppColors.info,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              'Add a bank account to withdraw funds from your wallet. '
              'Make sure the account name matches your HustleX profile.',
              style: AppTypography.bodySmall.copyWith(
                color: AppColors.info,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBankDropdown(List<Bank> banks) {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: DropdownButtonFormField<String>(
        value: _selectedBankCode,
        decoration: const InputDecoration(
          contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          border: InputBorder.none,
          hintText: 'Select your bank',
        ),
        isExpanded: true,
        icon: const Icon(Icons.keyboard_arrow_down),
        items: banks.map((bank) {
          return DropdownMenuItem(
            value: bank.code,
            child: Row(
              children: [
                if (bank.logoUrl != null) ...[
                  Container(
                    width: 32,
                    height: 32,
                    decoration: BoxDecoration(
                      color: AppColors.background,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Image.network(
                      bank.logoUrl!,
                      errorBuilder: (_, __, ___) => const Icon(
                        Icons.account_balance,
                        size: 20,
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                ],
                Expanded(
                  child: Text(
                    bank.name,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
          );
        }).toList(),
        onChanged: (value) {
          setState(() {
            _selectedBankCode = value;
            _selectedBankName = banks.firstWhere((b) => b.code == value).name;
            _verifiedAccountName = null;
            _error = null;
          });

          // Auto-verify if account number already entered
          if (_accountNumberController.text.length == 10) {
            _verifyAccount();
          }
        },
        validator: (value) {
          if (value == null) {
            return 'Please select a bank';
          }
          return null;
        },
      ),
    );
  }

  Widget _buildVerifyingIndicator() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          const SizedBox(
            width: 20,
            height: 20,
            child: CircularProgressIndicator(strokeWidth: 2),
          ),
          const SizedBox(width: 12),
          Text(
            'Verifying account...',
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildVerifiedAccountCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.success.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppColors.success.withOpacity(0.3),
        ),
      ),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: const BoxDecoration(
              color: AppColors.success,
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.check,
              color: Colors.white,
              size: 24,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Account Verified',
                  style: AppTypography.labelSmall.copyWith(
                    color: AppColors.success,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  _verifiedAccountName!,
                  style: AppTypography.titleSmall.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                Text(
                  '$_selectedBankName • ${_accountNumberController.text}',
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

  Widget _buildErrorMessage() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.error.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppColors.error.withOpacity(0.3),
        ),
      ),
      child: Row(
        children: [
          const Icon(
            Icons.error_outline,
            color: AppColors.error,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              _error!,
              style: AppTypography.bodySmall.copyWith(
                color: AppColors.error,
              ),
            ),
          ),
          TextButton(
            onPressed: _verifyAccount,
            child: const Text('Retry'),
          ),
        ],
      ),
    );
  }

  Widget _buildSecurityNote() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          const Icon(
            Icons.lock_outline,
            size: 16,
            color: AppColors.textSecondary,
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              'Your bank details are encrypted and securely stored',
              style: AppTypography.bodySmall.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          ),
        ],
      ),
    );
  }

  bool get _canSubmit =>
      _selectedBankCode != null &&
      _accountNumberController.text.length == 10 &&
      _verifiedAccountName != null &&
      !_isSubmitting;

  Future<void> _verifyAccount() async {
    if (_selectedBankCode == null || _accountNumberController.text.length != 10) {
      return;
    }

    setState(() {
      _isVerifying = true;
      _verifiedAccountName = null;
      _error = null;
    });

    try {
      final result = await ref.read(walletRepositoryProvider).verifyBankAccount(
            _accountNumberController.text,
            _selectedBankCode!,
          );

      result.when(
        success: (data) {
          setState(() {
            _verifiedAccountName = data.accountName;
            _isVerifying = false;
          });
        },
        failure: (message, _) {
          setState(() {
            _error = message;
            _isVerifying = false;
          });
        },
      );
    } catch (e) {
      setState(() {
        _error = 'Failed to verify account. Please try again.';
        _isVerifying = false;
      });
    }
  }

  Future<void> _submitAccount() async {
    if (!_formKey.currentState!.validate() || !_canSubmit) return;

    setState(() => _isSubmitting = true);

    try {
      final result = await ref.read(walletRepositoryProvider).addBankAccount(
            AddBankAccountRequest(
              accountNumber: _accountNumberController.text,
              bankCode: _selectedBankCode!,
              accountName: _verifiedAccountName!,
            ),
          );

      result.when(
        success: (account) {
          // Refresh bank accounts
          ref.invalidate(bankAccountsProvider);

          // Show success message
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Bank account added successfully'),
                backgroundColor: AppColors.success,
              ),
            );
            context.pop(account);
          }
        },
        failure: (message, _) {
          setState(() => _isSubmitting = false);
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(message),
                backgroundColor: AppColors.error,
              ),
            );
          }
        },
      );
    } catch (e) {
      setState(() => _isSubmitting = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Failed to add bank account'),
            backgroundColor: AppColors.error,
          ),
        );
      }
    }
  }
}

/// =============================================================================
/// BANK MODEL & PROVIDERS (to be moved to wallet models/providers)
/// =============================================================================

class Bank {
  final String code;
  final String name;
  final String? logoUrl;

  Bank({
    required this.code,
    required this.name,
    this.logoUrl,
  });

  factory Bank.fromJson(Map<String, dynamic> json) {
    return Bank(
      code: json['code'] as String,
      name: json['name'] as String,
      logoUrl: json['logo_url'] as String?,
    );
  }
}

class AddBankAccountRequest {
  final String accountNumber;
  final String bankCode;
  final String accountName;

  AddBankAccountRequest({
    required this.accountNumber,
    required this.bankCode,
    required this.accountName,
  });

  Map<String, dynamic> toJson() => {
        'account_number': accountNumber,
        'bank_code': bankCode,
        'account_name': accountName,
      };
}

class BankAccountVerification {
  final String accountNumber;
  final String accountName;
  final String bankCode;
  final String bankName;

  BankAccountVerification({
    required this.accountNumber,
    required this.accountName,
    required this.bankCode,
    required this.bankName,
  });

  factory BankAccountVerification.fromJson(Map<String, dynamic> json) {
    return BankAccountVerification(
      accountNumber: json['account_number'] as String,
      accountName: json['account_name'] as String,
      bankCode: json['bank_code'] as String,
      bankName: json['bank_name'] as String,
    );
  }
}
