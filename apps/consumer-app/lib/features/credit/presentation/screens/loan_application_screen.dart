import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../shared/widgets/buttons.dart';
import '../../../../shared/widgets/inputs.dart';

class LoanApplicationScreen extends ConsumerStatefulWidget {
  final double? creditLimit;
  final int? creditScore;

  const LoanApplicationScreen({
    super.key,
    this.creditLimit,
    this.creditScore,
  });

  @override
  ConsumerState<LoanApplicationScreen> createState() => _LoanApplicationScreenState();
}

class _LoanApplicationScreenState extends ConsumerState<LoanApplicationScreen> {
  final _pageController = PageController();
  final _formKey = GlobalKey<FormState>();
  
  int _currentStep = 0;
  bool _isSubmitting = false;
  
  // Step 1: Loan Amount
  final _amountController = TextEditingController();
  int _selectedTenure = 30; // days
  String _selectedPurpose = '';
  
  // Step 2: Employment Info
  String _employmentStatus = '';
  final _monthlyIncomeController = TextEditingController();
  final _employerNameController = TextEditingController();
  
  // Step 3: Bank Details
  String _selectedBankId = '';
  final _accountNumberController = TextEditingController();
  bool _agreeToTerms = false;
  
  final List<int> _tenureOptions = [7, 14, 30, 60, 90];
  
  final List<String> _purposeOptions = [
    'Business Expansion',
    'Emergency',
    'Education',
    'Medical',
    'Rent/Bills',
    'Personal',
    'Other',
  ];
  
  final List<String> _employmentOptions = [
    'Employed',
    'Self-Employed',
    'Business Owner',
    'Freelancer',
    'Student',
    'Unemployed',
  ];
  
  // Mock banks
  final List<Map<String, String>> _banks = [
    {'id': '1', 'name': 'Access Bank', 'code': '044'},
    {'id': '2', 'name': 'GTBank', 'code': '058'},
    {'id': '3', 'name': 'First Bank', 'code': '011'},
    {'id': '4', 'name': 'UBA', 'code': '033'},
    {'id': '5', 'name': 'Zenith Bank', 'code': '057'},
  ];

  @override
  void dispose() {
    _pageController.dispose();
    _amountController.dispose();
    _monthlyIncomeController.dispose();
    _employerNameController.dispose();
    _accountNumberController.dispose();
    super.dispose();
  }

  double get _loanAmount {
    final text = _amountController.text.replaceAll(',', '');
    return double.tryParse(text) ?? 0;
  }

  double get _monthlyIncome {
    final text = _monthlyIncomeController.text.replaceAll(',', '');
    return double.tryParse(text) ?? 0;
  }

  double get _interestRate {
    // Interest rate based on tenure
    if (_selectedTenure <= 7) return 0.05;
    if (_selectedTenure <= 14) return 0.08;
    if (_selectedTenure <= 30) return 0.10;
    if (_selectedTenure <= 60) return 0.15;
    return 0.20;
  }

  double get _interestAmount => _loanAmount * _interestRate;
  double get _totalRepayment => _loanAmount + _interestAmount;
  
  DateTime get _dueDate => DateTime.now().add(Duration(days: _selectedTenure));

  void _nextStep() {
    if (_currentStep < 2) {
      if (_validateCurrentStep()) {
        setState(() => _currentStep++);
        _pageController.nextPage(
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeInOut,
        );
      }
    } else {
      _submitApplication();
    }
  }

  void _previousStep() {
    if (_currentStep > 0) {
      setState(() => _currentStep--);
      _pageController.previousPage(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
    }
  }

  bool _validateCurrentStep() {
    switch (_currentStep) {
      case 0:
        if (_loanAmount <= 0) {
          _showError('Please enter loan amount');
          return false;
        }
        if (widget.creditLimit != null && _loanAmount > widget.creditLimit!) {
          _showError('Amount exceeds your credit limit');
          return false;
        }
        if (_selectedPurpose.isEmpty) {
          _showError('Please select loan purpose');
          return false;
        }
        return true;
      case 1:
        if (_employmentStatus.isEmpty) {
          _showError('Please select employment status');
          return false;
        }
        if (_monthlyIncome <= 0) {
          _showError('Please enter monthly income');
          return false;
        }
        return true;
      case 2:
        if (_selectedBankId.isEmpty) {
          _showError('Please select a bank');
          return false;
        }
        if (_accountNumberController.text.length != 10) {
          _showError('Please enter valid account number');
          return false;
        }
        if (!_agreeToTerms) {
          _showError('Please agree to the terms and conditions');
          return false;
        }
        return true;
      default:
        return false;
    }
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: AppColors.error,
      ),
    );
  }

  Future<void> _submitApplication() async {
    if (!_validateCurrentStep()) return;

    setState(() => _isSubmitting = true);

    try {
      // TODO: Submit loan application via API
      await Future.delayed(const Duration(seconds: 2));

      if (mounted) {
        _showSuccessDialog();
      }
    } catch (e) {
      if (mounted) {
        _showError('Failed to submit application: $e');
      }
    } finally {
      if (mounted) {
        setState(() => _isSubmitting = false);
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
              'Application Submitted!',
              style: AppTypography.headlineSmall.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Your loan application is being reviewed. You\'ll receive a decision within 24 hours.',
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: AppColors.surfaceVariant,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                children: [
                  _buildSummaryRow('Loan Amount', '₦${_loanAmount.toStringAsFixed(0)}'),
                  const SizedBox(height: 8),
                  _buildSummaryRow('Repayment', '₦${_totalRepayment.toStringAsFixed(0)}'),
                  const SizedBox(height: 8),
                  _buildSummaryRow('Due Date', _formatDate(_dueDate)),
                ],
              ),
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

  Widget _buildSummaryRow(String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label, style: AppTypography.bodySmall),
        Text(
          value,
          style: AppTypography.bodySmall.copyWith(fontWeight: FontWeight.w600),
        ),
      ],
    );
  }

  String _formatDate(DateTime date) {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return '${months[date.month - 1]} ${date.day}, ${date.year}';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Apply for Loan'),
      ),
      body: Column(
        children: [
          _buildProgressIndicator(),
          Expanded(
            child: PageView(
              controller: _pageController,
              physics: const NeverScrollableScrollPhysics(),
              children: [
                _buildStep1(),
                _buildStep2(),
                _buildStep3(),
              ],
            ),
          ),
        ],
      ),
      bottomNavigationBar: _buildBottomBar(),
    );
  }

  Widget _buildProgressIndicator() {
    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          Row(
            children: List.generate(3, (index) {
              final isCompleted = index < _currentStep;
              final isCurrent = index == _currentStep;
              return Expanded(
                child: Container(
                  margin: EdgeInsets.only(right: index < 2 ? 8 : 0),
                  height: 4,
                  decoration: BoxDecoration(
                    color: isCompleted || isCurrent
                        ? AppColors.primary
                        : AppColors.border,
                    borderRadius: BorderRadius.circular(2),
                  ),
                ),
              );
            }),
          ),
          const SizedBox(height: 8),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                ['Loan Details', 'Employment', 'Bank Details'][_currentStep],
                style: AppTypography.labelMedium.copyWith(
                  color: AppColors.primary,
                  fontWeight: FontWeight.w600,
                ),
              ),
              Text(
                'Step ${_currentStep + 1} of 3',
                style: AppTypography.bodySmall.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildStep1() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (widget.creditLimit != null) _buildCreditLimitCard(),
          const SizedBox(height: 24),
          _buildAmountSection(),
          const SizedBox(height: 24),
          _buildTenureSection(),
          const SizedBox(height: 24),
          _buildPurposeSection(),
          const SizedBox(height: 24),
          if (_loanAmount > 0) _buildLoanSummaryCard(),
        ],
      ),
    );
  }

  Widget _buildCreditLimitCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [AppColors.primary, AppColors.primary.withOpacity(0.8)],
        ),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.account_balance_wallet,
              color: Colors.white,
              size: 24,
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Your Credit Limit',
                  style: AppTypography.bodySmall.copyWith(
                    color: Colors.white.withOpacity(0.8),
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '₦${widget.creditLimit!.toStringAsFixed(0)}',
                  style: AppTypography.headlineMedium.copyWith(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
          if (widget.creditScore != null)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
              ),
              child: Text(
                'Score: ${widget.creditScore}',
                style: AppTypography.labelSmall.copyWith(
                  color: AppColors.primary,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildAmountSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'How much do you need?',
          style: AppTypography.titleMedium.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 12),
        AppTextField(
          controller: _amountController,
          label: 'Loan Amount',
          hint: '0',
          prefix: const Text('₦'),
          keyboardType: TextInputType.number,
          inputFormatters: [
            FilteringTextInputFormatter.digitsOnly,
            _ThousandsSeparatorFormatter(),
          ],
          onChanged: (_) => setState(() {}),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: [5000, 10000, 25000, 50000, 100000].map((amount) {
            return ActionChip(
              label: Text('₦${_formatNumber(amount)}'),
              onPressed: () {
                _amountController.text = amount.toString();
                setState(() {});
              },
            );
          }).toList(),
        ),
      ],
    );
  }

  Widget _buildTenureSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Repayment Period',
          style: AppTypography.titleMedium.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: _tenureOptions.map((days) {
            final isSelected = _selectedTenure == days;
            return ChoiceChip(
              label: Text('$days days'),
              selected: isSelected,
              onSelected: (selected) {
                if (selected) {
                  setState(() => _selectedTenure = days);
                }
              },
              selectedColor: AppColors.primary,
              labelStyle: TextStyle(
                color: isSelected ? Colors.white : AppColors.textPrimary,
              ),
            );
          }).toList(),
        ),
      ],
    );
  }

  Widget _buildPurposeSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'What\'s this loan for?',
          style: AppTypography.titleMedium.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: _purposeOptions.map((purpose) {
            final isSelected = _selectedPurpose == purpose;
            return ChoiceChip(
              label: Text(purpose),
              selected: isSelected,
              onSelected: (selected) {
                if (selected) {
                  setState(() => _selectedPurpose = purpose);
                }
              },
              selectedColor: AppColors.primary,
              labelStyle: TextStyle(
                color: isSelected ? Colors.white : AppColors.textPrimary,
              ),
            );
          }).toList(),
        ),
      ],
    );
  }

  Widget _buildLoanSummaryCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surfaceVariant,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Loan Summary',
            style: AppTypography.titleSmall.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 16),
          _buildDetailRow('Loan Amount', '₦${_loanAmount.toStringAsFixed(0)}'),
          _buildDetailRow('Interest Rate', '${(_interestRate * 100).toStringAsFixed(0)}%'),
          _buildDetailRow('Interest Amount', '₦${_interestAmount.toStringAsFixed(0)}'),
          const Divider(height: 24),
          _buildDetailRow(
            'Total Repayment',
            '₦${_totalRepayment.toStringAsFixed(0)}',
            isBold: true,
          ),
          _buildDetailRow('Due Date', _formatDate(_dueDate)),
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value, {bool isBold = false}) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          Text(
            value,
            style: isBold
                ? AppTypography.titleMedium.copyWith(fontWeight: FontWeight.bold)
                : AppTypography.bodyMedium.copyWith(fontWeight: FontWeight.w500),
          ),
        ],
      ),
    );
  }

  Widget _buildStep2() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Employment Information',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'This helps us verify your ability to repay the loan.',
            style: AppTypography.bodySmall.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 24),
          Text(
            'Employment Status',
            style: AppTypography.labelMedium.copyWith(
              fontWeight: FontWeight.w500,
            ),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: _employmentOptions.map((status) {
              final isSelected = _employmentStatus == status;
              return ChoiceChip(
                label: Text(status),
                selected: isSelected,
                onSelected: (selected) {
                  if (selected) {
                    setState(() => _employmentStatus = status);
                  }
                },
                selectedColor: AppColors.primary,
                labelStyle: TextStyle(
                  color: isSelected ? Colors.white : AppColors.textPrimary,
                ),
              );
            }).toList(),
          ),
          const SizedBox(height: 24),
          AppTextField(
            controller: _monthlyIncomeController,
            label: 'Monthly Income',
            hint: '0',
            prefix: const Text('₦'),
            keyboardType: TextInputType.number,
            inputFormatters: [
              FilteringTextInputFormatter.digitsOnly,
              _ThousandsSeparatorFormatter(),
            ],
            onChanged: (_) => setState(() {}),
          ),
          const SizedBox(height: 16),
          if (_employmentStatus == 'Employed')
            AppTextField(
              controller: _employerNameController,
              label: 'Employer Name',
              hint: 'Enter your employer\'s name',
            ),
          const SizedBox(height: 24),
          _buildIncomeVerificationNote(),
        ],
      ),
    );
  }

  Widget _buildIncomeVerificationNote() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.info.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.info.withOpacity(0.3)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(
            Icons.info_outline,
            color: AppColors.info,
            size: 20,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Income Verification',
                  style: AppTypography.labelMedium.copyWith(
                    fontWeight: FontWeight.w600,
                    color: AppColors.info,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  'We may request bank statements or pay slips to verify your income. Providing accurate information speeds up approval.',
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

  Widget _buildStep3() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Bank Account for Disbursement',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Where should we send your loan?',
            style: AppTypography.bodySmall.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 24),
          Text(
            'Select Bank',
            style: AppTypography.labelMedium.copyWith(
              fontWeight: FontWeight.w500,
            ),
          ),
          const SizedBox(height: 12),
          ..._banks.map((bank) => _buildBankOption(bank)),
          const SizedBox(height: 24),
          AppTextField(
            controller: _accountNumberController,
            label: 'Account Number',
            hint: 'Enter 10-digit account number',
            keyboardType: TextInputType.number,
            inputFormatters: [
              FilteringTextInputFormatter.digitsOnly,
              LengthLimitingTextInputFormatter(10),
            ],
          ),
          const SizedBox(height: 24),
          _buildFinalSummary(),
          const SizedBox(height: 24),
          _buildTermsCheckbox(),
        ],
      ),
    );
  }

  Widget _buildBankOption(Map<String, String> bank) {
    final isSelected = _selectedBankId == bank['id'];
    return GestureDetector(
      onTap: () => setState(() => _selectedBankId = bank['id']!),
      child: Container(
        margin: const EdgeInsets.only(bottom: 8),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected ? AppColors.primary.withOpacity(0.1) : Colors.white,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isSelected ? AppColors.primary : AppColors.border,
            width: isSelected ? 2 : 1,
          ),
        ),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: AppColors.surfaceVariant,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Center(
                child: Text(
                  bank['name']![0],
                  style: AppTypography.titleMedium.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Text(
                bank['name']!,
                style: AppTypography.bodyMedium.copyWith(
                  fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
                ),
              ),
            ),
            if (isSelected)
              const Icon(
                Icons.check_circle,
                color: AppColors.primary,
                size: 24,
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildFinalSummary() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [AppColors.primary, AppColors.primary.withOpacity(0.8)],
        ),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          Text(
            'Loan Application Summary',
            style: AppTypography.titleSmall.copyWith(
              color: Colors.white,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 16),
          _buildWhiteDetailRow('Loan Amount', '₦${_loanAmount.toStringAsFixed(0)}'),
          _buildWhiteDetailRow('Interest', '₦${_interestAmount.toStringAsFixed(0)} (${(_interestRate * 100).toStringAsFixed(0)}%)'),
          _buildWhiteDetailRow('Tenure', '$_selectedTenure days'),
          _buildWhiteDetailRow('Purpose', _selectedPurpose),
          const Divider(color: Colors.white24, height: 24),
          _buildWhiteDetailRow('Total Repayment', '₦${_totalRepayment.toStringAsFixed(0)}', isBold: true),
          _buildWhiteDetailRow('Due Date', _formatDate(_dueDate)),
        ],
      ),
    );
  }

  Widget _buildWhiteDetailRow(String label, String value, {bool isBold = false}) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: AppTypography.bodySmall.copyWith(
              color: Colors.white.withOpacity(0.8),
            ),
          ),
          Text(
            value,
            style: isBold
                ? AppTypography.titleMedium.copyWith(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                  )
                : AppTypography.bodySmall.copyWith(
                    color: Colors.white,
                    fontWeight: FontWeight.w500,
                  ),
          ),
        ],
      ),
    );
  }

  Widget _buildTermsCheckbox() {
    return GestureDetector(
      onTap: () => setState(() => _agreeToTerms = !_agreeToTerms),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Checkbox(
            value: _agreeToTerms,
            onChanged: (value) => setState(() => _agreeToTerms = value ?? false),
            activeColor: AppColors.primary,
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.only(top: 12),
              child: Text.rich(
                TextSpan(
                  text: 'I agree to the ',
                  style: AppTypography.bodySmall,
                  children: [
                    TextSpan(
                      text: 'Loan Agreement',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.primary,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const TextSpan(text: ' and '),
                    TextSpan(
                      text: 'Terms of Service',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.primary,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const TextSpan(
                      text: '. I understand that late repayment may affect my credit score.',
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomBar() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Theme.of(context).scaffoldBackgroundColor,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, -5),
          ),
        ],
      ),
      child: SafeArea(
        child: Row(
          children: [
            if (_currentStep > 0)
              Expanded(
                child: OutlinedButton(
                  onPressed: _previousStep,
                  child: const Text('Back'),
                ),
              ),
            if (_currentStep > 0) const SizedBox(width: 16),
            Expanded(
              flex: _currentStep > 0 ? 2 : 1,
              child: PrimaryButton(
                text: _currentStep < 2 ? 'Continue' : 'Submit Application',
                onPressed: _isSubmitting ? null : _nextStep,
                isLoading: _isSubmitting,
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatNumber(int number) {
    return number.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (Match m) => '${m[1]},',
    );
  }
}

class _ThousandsSeparatorFormatter extends TextInputFormatter {
  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    if (newValue.text.isEmpty) {
      return newValue;
    }

    final number = int.tryParse(newValue.text.replaceAll(',', ''));
    if (number == null) {
      return oldValue;
    }

    final formatted = number.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (Match m) => '${m[1]},',
    );

    return TextEditingValue(
      text: formatted,
      selection: TextSelection.collapsed(offset: formatted.length),
    );
  }
}
