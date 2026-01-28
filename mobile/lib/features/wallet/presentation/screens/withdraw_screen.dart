import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/utils/helpers.dart';
import '../../../../core/widgets/buttons.dart';
import '../../../../core/widgets/inputs.dart';
import '../../data/models/wallet_model.dart';
import '../providers/wallet_provider.dart';

class WithdrawScreen extends ConsumerStatefulWidget {
  const WithdrawScreen({super.key});

  @override
  ConsumerState<WithdrawScreen> createState() => _WithdrawScreenState();
}

class _WithdrawScreenState extends ConsumerState<WithdrawScreen> {
  final _formKey = GlobalKey<FormState>();
  final _amountController = TextEditingController();
  final _pinController = TextEditingController();
  final _narrationController = TextEditingController();
  
  BankAccount? _selectedAccount;
  bool _showPinInput = false;

  @override
  void initState() {
    super.initState();
    // Load bank accounts
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(walletProvider.notifier).loadBankAccounts();
    });
  }

  @override
  void dispose() {
    _amountController.dispose();
    _pinController.dispose();
    _narrationController.dispose();
    super.dispose();
  }

  double get _amount {
    final text = _amountController.text.replaceAll(',', '');
    return double.tryParse(text) ?? 0;
  }

  void _proceedToPin() {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedAccount == null) {
      _showError('Please select a bank account');
      return;
    }

    setState(() => _showPinInput = true);
  }

  Future<void> _processWithdrawal() async {
    if (_pinController.text.length != 4) {
      _showError('Please enter your 4-digit PIN');
      return;
    }

    final transferNotifier = ref.read(transferProvider.notifier);
    
    final request = WithdrawalRequest(
      amount: _amount,
      bankAccountId: _selectedAccount!.id,
      pin: _pinController.text,
      narration: _narrationController.text.isEmpty 
          ? null 
          : _narrationController.text,
    );

    final success = await transferNotifier.withdraw(request);

    if (success && mounted) {
      // Refresh wallet balance
      ref.read(walletProvider.notifier).refreshBalance();
      _showSuccess();
    } else if (mounted) {
      final error = ref.read(transferProvider).error;
      _showError(error ?? 'Withdrawal failed');
      _pinController.clear();
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
              'Withdrawal Initiated!',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '${CurrencyUtils.formatNaira(_amount)} is being sent to ${_selectedAccount!.accountName}',
              textAlign: TextAlign.center,
              style: const TextStyle(color: AppColors.textSecondary),
            ),
            const SizedBox(height: 8),
            Text(
              _selectedAccount!.maskedAccountNumber,
              style: const TextStyle(
                fontWeight: FontWeight.w600,
                color: AppColors.textSecondary,
              ),
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

  void _showAddBankAccountSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => const _AddBankAccountSheet(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final walletState = ref.watch(walletProvider);
    final transferState = ref.watch(transferProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Withdraw'),
        centerTitle: true,
      ),
      body: SafeArea(
        child: _showPinInput
            ? _buildPinInputView(transferState)
            : _buildMainForm(walletState),
      ),
    );
  }

  Widget _buildMainForm(WalletState walletState) {
    return Form(
      key: _formKey,
      child: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          // Current balance card
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: AppColors.surface,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: AppColors.border),
            ),
            child: Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: AppColors.primary.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: const Icon(
                    Icons.account_balance_wallet,
                    color: AppColors.primary,
                  ),
                ),
                const SizedBox(width: 16),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Available Balance',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      CurrencyUtils.formatNaira(walletState.availableBalance),
                      style: AppTypography.titleLarge.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),

          const SizedBox(height: 32),

          // Amount input
          Text(
            'Amount to Withdraw',
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
                return 'Minimum withdrawal is ₦100';
              }
              if (amount > walletState.availableBalance) {
                return 'Insufficient balance';
              }
              return null;
            },
          ),

          const SizedBox(height: 24),

          // Bank account selection
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Select Bank Account',
                style: AppTypography.titleMedium.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              TextButton.icon(
                onPressed: _showAddBankAccountSheet,
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Add New'),
              ),
            ],
          ),
          const SizedBox(height: 12),

          if (walletState.bankAccounts.isEmpty)
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: AppColors.surface,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: AppColors.border),
              ),
              child: Column(
                children: [
                  const Icon(
                    Icons.account_balance,
                    size: 48,
                    color: AppColors.textSecondary,
                  ),
                  const SizedBox(height: 12),
                  const Text(
                    'No bank accounts added',
                    style: TextStyle(color: AppColors.textSecondary),
                  ),
                  const SizedBox(height: 12),
                  SecondaryButton(
                    text: 'Add Bank Account',
                    onPressed: _showAddBankAccountSheet,
                  ),
                ],
              ),
            )
          else
            ...walletState.bankAccounts.map((account) => Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: _BankAccountTile(
                account: account,
                isSelected: _selectedAccount?.id == account.id,
                onTap: () => setState(() => _selectedAccount = account),
              ),
            )),

          const SizedBox(height: 24),

          // Narration (optional)
          Text(
            'Narration (Optional)',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          AppTextField(
            controller: _narrationController,
            hintText: 'What is this withdrawal for?',
            maxLines: 2,
          ),

          const SizedBox(height: 16),

          // Fee info
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: AppColors.warning.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppColors.warning.withOpacity(0.3)),
            ),
            child: Row(
              children: [
                const Icon(Icons.info_outline, color: AppColors.warning),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    'A fee of ₦50 will be charged for this withdrawal.',
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.warning,
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
            onPressed: _proceedToPin,
          ),
        ],
      ),
    );
  }

  Widget _buildPinInputView(TransferState transferState) {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          const SizedBox(height: 40),
          
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: AppColors.primary.withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.lock,
              size: 40,
              color: AppColors.primary,
            ),
          ),
          
          const SizedBox(height: 24),
          
          Text(
            'Enter PIN',
            style: AppTypography.headlineSmall.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          
          const SizedBox(height: 8),
          
          Text(
            'Enter your 4-digit PIN to confirm withdrawal',
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
            textAlign: TextAlign.center,
          ),

          const SizedBox(height: 16),

          // Transaction summary
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: AppColors.surface,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppColors.border),
            ),
            child: Column(
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Amount'),
                    Text(
                      CurrencyUtils.formatNaira(_amount),
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Fee'),
                    Text(
                      CurrencyUtils.formatNaira(50),
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                  ],
                ),
                const Divider(height: 24),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('To', style: TextStyle(color: AppColors.textSecondary)),
                    Expanded(
                      child: Text(
                        '${_selectedAccount!.bankName}\n${_selectedAccount!.maskedAccountNumber}',
                        textAlign: TextAlign.right,
                        style: const TextStyle(fontWeight: FontWeight.w600),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),

          const SizedBox(height: 32),
          
          SizedBox(
            width: 200,
            child: AppTextField(
              controller: _pinController,
              hintText: '••••',
              keyboardType: TextInputType.number,
              maxLength: 4,
              obscureText: true,
              textAlign: TextAlign.center,
              style: AppTypography.headlineMedium.copyWith(
                letterSpacing: 16,
              ),
            ),
          ),
          
          const Spacer(),
          
          Row(
            children: [
              Expanded(
                child: SecondaryButton(
                  text: 'Back',
                  onPressed: () {
                    setState(() => _showPinInput = false);
                    _pinController.clear();
                  },
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: PrimaryButton(
                  text: 'Confirm',
                  onPressed: _processWithdrawal,
                  isLoading: transferState.isProcessing,
                ),
              ),
            ],
          ),
          
          const SizedBox(height: 20),
        ],
      ),
    );
  }
}

class _BankAccountTile extends StatelessWidget {
  final BankAccount account;
  final bool isSelected;
  final VoidCallback onTap;

  const _BankAccountTile({
    required this.account,
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
                color: AppColors.backgroundDark,
                borderRadius: BorderRadius.circular(10),
              ),
              child: const Icon(
                Icons.account_balance,
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          account.accountName,
                          style: AppTypography.bodyMedium.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      if (account.isDefault)
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 8,
                            vertical: 2,
                          ),
                          decoration: BoxDecoration(
                            color: AppColors.primary.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: Text(
                            'Default',
                            style: AppTypography.labelSmall.copyWith(
                              color: AppColors.primary,
                            ),
                          ),
                        ),
                    ],
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '${account.bankName} • ${account.maskedAccountNumber}',
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

class _AddBankAccountSheet extends ConsumerStatefulWidget {
  const _AddBankAccountSheet();

  @override
  ConsumerState<_AddBankAccountSheet> createState() => _AddBankAccountSheetState();
}

class _AddBankAccountSheetState extends ConsumerState<_AddBankAccountSheet> {
  final _formKey = GlobalKey<FormState>();
  final _accountNumberController = TextEditingController();
  Bank? _selectedBank;

  @override
  void dispose() {
    _accountNumberController.dispose();
    super.dispose();
  }

  Future<void> _verifyAndAdd() async {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedBank == null) return;

    final verificationNotifier = ref.read(bankVerificationProvider.notifier);
    
    final verified = await verificationNotifier.verifyAccount(
      bankCode: _selectedBank!.code,
      accountNumber: _accountNumberController.text,
    );

    if (!verified || !mounted) return;

    final verification = ref.read(bankVerificationProvider).verification;
    if (verification == null) return;

    // Show confirmation
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Confirm Account'),
        content: Text('Is this account name correct?\n\n${verification.accountName}'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('No'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Yes, Add'),
          ),
        ],
      ),
    );

    if (confirmed != true) return;

    // Add account
    final added = await ref.read(walletProvider.notifier).addBankAccount(
      bankCode: _selectedBank!.code,
      accountNumber: _accountNumberController.text,
    );

    if (added && mounted) {
      Navigator.pop(context);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Bank account added successfully'),
          backgroundColor: AppColors.success,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final banks = ref.watch(banksProvider);
    final verificationState = ref.watch(bankVerificationProvider);

    return Padding(
      padding: EdgeInsets.only(
        left: 20,
        right: 20,
        top: 20,
        bottom: MediaQuery.of(context).viewInsets.bottom + 20,
      ),
      child: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Center(
              child: Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: AppColors.border,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
            const SizedBox(height: 20),
            
            Text(
              'Add Bank Account',
              style: AppTypography.titleLarge.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            
            const SizedBox(height: 24),
            
            // Bank selection
            Text(
              'Select Bank',
              style: AppTypography.bodyMedium.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            banks.when(
              data: (bankList) => AppDropdown<Bank>(
                value: _selectedBank,
                items: bankList,
                itemBuilder: (bank) => Text(bank.name),
                onChanged: (bank) => setState(() => _selectedBank = bank),
                hintText: 'Select a bank',
              ),
              loading: () => const LinearProgressIndicator(),
              error: (_, __) => const Text('Failed to load banks'),
            ),
            
            const SizedBox(height: 16),
            
            // Account number
            Text(
              'Account Number',
              style: AppTypography.bodyMedium.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            AppTextField(
              controller: _accountNumberController,
              hintText: '0123456789',
              keyboardType: TextInputType.number,
              maxLength: 10,
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
            
            if (verificationState.error != null) ...[
              const SizedBox(height: 8),
              Text(
                verificationState.error!,
                style: const TextStyle(color: AppColors.error),
              ),
            ],
            
            const SizedBox(height: 24),
            
            PrimaryButton(
              text: 'Verify & Add',
              onPressed: _verifyAndAdd,
              isLoading: verificationState.isVerifying,
            ),
          ],
        ),
      ),
    );
  }
}
