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

class TransferScreen extends ConsumerStatefulWidget {
  const TransferScreen({super.key});

  @override
  ConsumerState<TransferScreen> createState() => _TransferScreenState();
}

class _TransferScreenState extends ConsumerState<TransferScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  
  final _formKey = GlobalKey<FormState>();
  final _amountController = TextEditingController();
  final _phoneController = TextEditingController();
  final _accountNumberController = TextEditingController();
  final _narrationController = TextEditingController();
  final _pinController = TextEditingController();
  
  Bank? _selectedBank;
  bool _showPinInput = false;
  String? _recipientName;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _tabController.addListener(() {
      if (_tabController.indexIsChanging) {
        setState(() {
          _showPinInput = false;
          _recipientName = null;
          _pinController.clear();
        });
      }
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    _amountController.dispose();
    _phoneController.dispose();
    _accountNumberController.dispose();
    _narrationController.dispose();
    _pinController.dispose();
    super.dispose();
  }

  double get _amount {
    final text = _amountController.text.replaceAll(',', '');
    return double.tryParse(text) ?? 0;
  }

  bool get _isHustleXTransfer => _tabController.index == 0;

  void _proceedToPin() {
    if (!_formKey.currentState!.validate()) return;

    if (_isHustleXTransfer) {
      // For HustleX transfers, verify recipient exists
      _verifyHustleXRecipient();
    } else {
      // For bank transfers, verify bank account
      _verifyBankAccount();
    }
  }

  Future<void> _verifyHustleXRecipient() async {
    // Simulate recipient lookup
    // In real app, call API to verify recipient
    setState(() {
      _recipientName = 'John Doe'; // Would come from API
      _showPinInput = true;
    });
  }

  Future<void> _verifyBankAccount() async {
    if (_selectedBank == null) {
      _showError('Please select a bank');
      return;
    }

    final verificationNotifier = ref.read(bankVerificationProvider.notifier);
    
    final verified = await verificationNotifier.verifyAccount(
      bankCode: _selectedBank!.code,
      accountNumber: _accountNumberController.text,
    );

    if (verified && mounted) {
      final verification = ref.read(bankVerificationProvider).verification;
      setState(() {
        _recipientName = verification?.accountName ?? 'Unknown';
        _showPinInput = true;
      });
    } else if (mounted) {
      final error = ref.read(bankVerificationProvider).error;
      _showError(error ?? 'Failed to verify account');
    }
  }

  Future<void> _processTransfer() async {
    if (_pinController.text.length != 4) {
      _showError('Please enter your 4-digit PIN');
      return;
    }

    final transferNotifier = ref.read(transferProvider.notifier);
    bool success;

    if (_isHustleXTransfer) {
      final request = TransferRequest(
        amount: _amount,
        recipientPhone: PhoneUtils.formatNigerian(_phoneController.text),
        pin: _pinController.text,
        narration: _narrationController.text.isEmpty 
            ? null 
            : _narrationController.text,
      );
      success = await transferNotifier.transfer(request);
    } else {
      final request = BankTransferRequest(
        amount: _amount,
        bankCode: _selectedBank!.code,
        accountNumber: _accountNumberController.text,
        accountName: _recipientName!,
        pin: _pinController.text,
        narration: _narrationController.text.isEmpty 
            ? null 
            : _narrationController.text,
      );
      success = await transferNotifier.bankTransfer(request);
    }

    if (success && mounted) {
      ref.read(walletProvider.notifier).refreshBalance();
      _showSuccess();
    } else if (mounted) {
      final error = ref.read(transferProvider).error;
      _showError(error ?? 'Transfer failed');
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
              'Transfer Successful!',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '${CurrencyUtils.formatNaira(_amount)} sent to $_recipientName',
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
    final walletState = ref.watch(walletProvider);
    final transferState = ref.watch(transferProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Transfer'),
        centerTitle: true,
        bottom: _showPinInput
            ? null
            : TabBar(
                controller: _tabController,
                labelColor: AppColors.primary,
                unselectedLabelColor: AppColors.textSecondary,
                indicatorColor: AppColors.primary,
                tabs: const [
                  Tab(text: 'HustleX User'),
                  Tab(text: 'Bank Account'),
                ],
              ),
      ),
      body: SafeArea(
        child: _showPinInput
            ? _buildPinInputView(transferState)
            : TabBarView(
                controller: _tabController,
                children: [
                  _buildHustleXTransferForm(walletState),
                  _buildBankTransferForm(walletState),
                ],
              ),
      ),
    );
  }

  Widget _buildHustleXTransferForm(WalletState walletState) {
    return Form(
      key: _formKey,
      child: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          // Info banner
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              gradient: AppColors.primaryGradient,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              children: [
                const Icon(Icons.bolt, color: Colors.white),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Instant & Free',
                        style: AppTypography.bodyMedium.copyWith(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Text(
                        'Transfer to other HustleX users instantly with no fees',
                        style: AppTypography.bodySmall.copyWith(
                          color: Colors.white.withOpacity(0.8),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),

          const SizedBox(height: 24),

          // Recipient phone
          Text(
            'Recipient Phone Number',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          PhoneTextField(
            controller: _phoneController,
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter phone number';
              }
              if (!PhoneUtils.isValidNigerian(value)) {
                return 'Please enter a valid Nigerian phone number';
              }
              return null;
            },
          ),

          const SizedBox(height: 24),

          // Amount
          Text(
            'Amount',
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
              if (amount == null || amount < 10) {
                return 'Minimum transfer is ₦10';
              }
              if (amount > walletState.availableBalance) {
                return 'Insufficient balance';
              }
              return null;
            },
          ),
          const SizedBox(height: 8),
          Text(
            'Available: ${CurrencyUtils.formatNaira(walletState.availableBalance)}',
            style: AppTypography.bodySmall.copyWith(
              color: AppColors.textSecondary,
            ),
          ),

          const SizedBox(height: 24),

          // Narration
          Text(
            'Note (Optional)',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          AppTextField(
            controller: _narrationController,
            hintText: 'What is this transfer for?',
            maxLines: 2,
          ),

          const SizedBox(height: 32),

          PrimaryButton(
            text: 'Continue',
            onPressed: _proceedToPin,
          ),
        ],
      ),
    );
  }

  Widget _buildBankTransferForm(WalletState walletState) {
    final banks = ref.watch(banksProvider);
    final verificationState = ref.watch(bankVerificationProvider);

    return Form(
      key: _formKey,
      child: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          // Bank selection
          Text(
            'Select Bank',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
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

          const SizedBox(height: 24),

          // Account number
          Text(
            'Account Number',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
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

          const SizedBox(height: 24),

          // Amount
          Text(
            'Amount',
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
                return 'Minimum transfer is ₦100';
              }
              if (amount > walletState.availableBalance) {
                return 'Insufficient balance';
              }
              return null;
            },
          ),
          const SizedBox(height: 8),
          Text(
            'Available: ${CurrencyUtils.formatNaira(walletState.availableBalance)}',
            style: AppTypography.bodySmall.copyWith(
              color: AppColors.textSecondary,
            ),
          ),

          const SizedBox(height: 24),

          // Narration
          Text(
            'Narration (Optional)',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          AppTextField(
            controller: _narrationController,
            hintText: 'What is this transfer for?',
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
                    'A fee of ₦25 will be charged for bank transfers.',
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.warning,
                    ),
                  ),
                ),
              ],
            ),
          ),

          if (verificationState.error != null) ...[
            const SizedBox(height: 8),
            Text(
              verificationState.error!,
              style: const TextStyle(color: AppColors.error),
            ),
          ],

          const SizedBox(height: 32),

          PrimaryButton(
            text: 'Continue',
            onPressed: _proceedToPin,
            isLoading: verificationState.isVerifying,
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
              Icons.send_rounded,
              size: 40,
              color: AppColors.primary,
            ),
          ),
          
          const SizedBox(height: 24),
          
          Text(
            'Confirm Transfer',
            style: AppTypography.headlineSmall.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          
          const SizedBox(height: 8),
          
          Text(
            'Enter your PIN to send money',
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),

          const SizedBox(height: 24),

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
                // Recipient
                Row(
                  children: [
                    CircleAvatar(
                      backgroundColor: AppColors.primary.withOpacity(0.1),
                      child: Text(
                        StringUtils.getInitials(_recipientName ?? 'U'),
                        style: const TextStyle(
                          color: AppColors.primary,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            _recipientName ?? 'Unknown',
                            style: const TextStyle(fontWeight: FontWeight.w600),
                          ),
                          Text(
                            _isHustleXTransfer
                                ? PhoneUtils.mask(_phoneController.text)
                                : '${_selectedBank?.name ?? ''} • ****${_accountNumberController.text.substring(6)}',
                            style: AppTypography.bodySmall.copyWith(
                              color: AppColors.textSecondary,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                const Divider(height: 24),
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
                if (!_isHustleXTransfer) ...[
                  const SizedBox(height: 8),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Fee'),
                      Text(
                        CurrencyUtils.formatNaira(25),
                        style: const TextStyle(fontWeight: FontWeight.w600),
                      ),
                    ],
                  ),
                ],
                const Divider(height: 24),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Total', style: TextStyle(fontWeight: FontWeight.bold)),
                    Text(
                      CurrencyUtils.formatNaira(_amount + (_isHustleXTransfer ? 0 : 25)),
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        color: AppColors.primary,
                        fontSize: 18,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),

          const SizedBox(height: 32),

          Text(
            'Enter PIN',
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 12),
          
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
                  text: 'Send',
                  onPressed: _processTransfer,
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
