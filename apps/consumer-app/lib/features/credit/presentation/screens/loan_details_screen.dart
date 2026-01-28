import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/utils/currency_utils.dart';
import '../../../../core/utils/date_utils.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../data/models/credit_models.dart';

class LoanDetailsScreen extends ConsumerWidget {
  final String loanId;

  const LoanDetailsScreen({
    super.key,
    required this.loanId,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // TODO: Fetch from provider
    final loan = Loan(
      id: loanId,
      userId: 'user-1',
      amount: 75000,
      interestRate: 10,
      interestAmount: 7500,
      totalAmount: 82500,
      tenureDays: 30,
      purpose: LoanPurpose.business,
      status: LoanStatus.active,
      disbursementStatus: DisbursementStatus.completed,
      repaymentStatus: RepaymentStatus.onTrack,
      amountRepaid: 27500,
      applicationDate: DateTime.now().subtract(const Duration(days: 15)),
      approvalDate: DateTime.now().subtract(const Duration(days: 14)),
      disbursementDate: DateTime.now().subtract(const Duration(days: 14)),
      dueDate: DateTime.now().add(const Duration(days: 16)),
      bankAccountId: 'bank-1',
      createdAt: DateTime.now().subtract(const Duration(days: 15)),
      updatedAt: DateTime.now(),
    );

    final repaymentSchedule = [
      _RepaymentItem(
        amount: 27500,
        dueDate: DateTime.now().subtract(const Duration(days: 7)),
        status: 'paid',
        paidDate: DateTime.now().subtract(const Duration(days: 7)),
      ),
      _RepaymentItem(
        amount: 27500,
        dueDate: DateTime.now().add(const Duration(days: 1)),
        status: 'upcoming',
      ),
      _RepaymentItem(
        amount: 27500,
        dueDate: DateTime.now().add(const Duration(days: 8)),
        status: 'pending',
      ),
    ];

    return Scaffold(
      backgroundColor: AppColors.background,
      body: CustomScrollView(
        slivers: [
          _buildAppBar(context, loan),
          SliverToBoxAdapter(
            child: Column(
              children: [
                _buildStatusCard(loan),
                const SizedBox(height: 16),
                _buildProgressSection(loan),
                const SizedBox(height: 16),
                _buildDetailsCard(loan),
                const SizedBox(height: 16),
                _buildRepaymentSchedule(repaymentSchedule),
                const SizedBox(height: 16),
                _buildActionsSection(context, loan),
                const SizedBox(height: 32),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAppBar(BuildContext context, Loan loan) {
    return SliverAppBar(
      expandedHeight: 200,
      pinned: true,
      flexibleSpace: FlexibleSpaceBar(
        background: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.primary, AppColors.primaryDark],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
          child: SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  Text(
                    CurrencyUtils.formatNaira(loan.totalAmount),
                    style: AppTypography.headlineLarge.copyWith(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'Total Loan Amount',
                    style: AppTypography.bodyMedium.copyWith(
                      color: Colors.white.withOpacity(0.8),
                    ),
                  ),
                  const SizedBox(height: 16),
                  _buildStatusBadge(loan.status),
                ],
              ),
            ),
          ),
        ),
      ),
      actions: [
        IconButton(
          icon: const Icon(Icons.help_outline),
          onPressed: () => _showLoanHelp(context),
        ),
      ],
    );
  }

  Widget _buildStatusBadge(LoanStatus status) {
    Color bgColor;
    Color textColor;
    String label;
    IconData icon;

    switch (status) {
      case LoanStatus.pending:
        bgColor = Colors.white;
        textColor = AppColors.warning;
        label = 'Pending Approval';
        icon = Icons.hourglass_empty;
        break;
      case LoanStatus.approved:
        bgColor = Colors.white;
        textColor = AppColors.info;
        label = 'Approved';
        icon = Icons.check;
        break;
      case LoanStatus.active:
        bgColor = Colors.white;
        textColor = AppColors.success;
        label = 'Active';
        icon = Icons.trending_up;
        break;
      case LoanStatus.overdue:
        bgColor = Colors.white;
        textColor = AppColors.error;
        label = 'Overdue';
        icon = Icons.warning;
        break;
      case LoanStatus.completed:
        bgColor = Colors.white;
        textColor = AppColors.success;
        label = 'Completed';
        icon = Icons.check_circle;
        break;
      case LoanStatus.defaulted:
        bgColor = Colors.white;
        textColor = AppColors.error;
        label = 'Defaulted';
        icon = Icons.cancel;
        break;
      case LoanStatus.rejected:
        bgColor = Colors.white;
        textColor = AppColors.error;
        label = 'Rejected';
        icon = Icons.block;
        break;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 16, color: textColor),
          const SizedBox(width: 6),
          Text(
            label,
            style: AppTypography.labelMedium.copyWith(
              color: textColor,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusCard(Loan loan) {
    final daysRemaining = loan.dueDate.difference(DateTime.now()).inDays;
    final isOverdue = daysRemaining < 0;

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          Expanded(
            child: _buildStatusItem(
              icon: Icons.calendar_today,
              label: 'Due Date',
              value: AppDateUtils.formatDate(loan.dueDate),
              color: isOverdue ? AppColors.error : AppColors.textPrimary,
            ),
          ),
          Container(
            width: 1,
            height: 50,
            color: AppColors.border,
          ),
          Expanded(
            child: _buildStatusItem(
              icon: Icons.timer,
              label: isOverdue ? 'Days Overdue' : 'Days Left',
              value: '${daysRemaining.abs()}',
              color: isOverdue ? AppColors.error : AppColors.success,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusItem({
    required IconData icon,
    required String label,
    required String value,
    Color? color,
  }) {
    return Column(
      children: [
        Icon(icon, color: color ?? AppColors.textSecondary, size: 24),
        const SizedBox(height: 8),
        Text(
          value,
          style: AppTypography.titleLarge.copyWith(
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: AppTypography.bodySmall.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildProgressSection(Loan loan) {
    final progress = loan.amountRepaid / loan.totalAmount;
    final remaining = loan.totalAmount - loan.amountRepaid;

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Repayment Progress',
                style: AppTypography.titleMedium.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              Text(
                '${(progress * 100).toInt()}%',
                style: AppTypography.titleMedium.copyWith(
                  fontWeight: FontWeight.bold,
                  color: AppColors.primary,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          ClipRRect(
            borderRadius: BorderRadius.circular(4),
            child: LinearProgressIndicator(
              value: progress,
              minHeight: 8,
              backgroundColor: AppColors.border,
              valueColor: const AlwaysStoppedAnimation<Color>(AppColors.primary),
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Amount Repaid',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      CurrencyUtils.formatNaira(loan.amountRepaid),
                      style: AppTypography.titleMedium.copyWith(
                        fontWeight: FontWeight.w600,
                        color: AppColors.success,
                      ),
                    ),
                  ],
                ),
              ),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      'Remaining',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      CurrencyUtils.formatNaira(remaining),
                      style: AppTypography.titleMedium.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildDetailsCard(Loan loan) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Loan Details',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 16),
          _buildDetailRow('Principal Amount', CurrencyUtils.formatNaira(loan.amount)),
          const Divider(height: 24),
          _buildDetailRow('Interest Rate', '${loan.interestRate}%'),
          const Divider(height: 24),
          _buildDetailRow('Interest Amount', CurrencyUtils.formatNaira(loan.interestAmount)),
          const Divider(height: 24),
          _buildDetailRow('Tenure', '${loan.tenureDays} days'),
          const Divider(height: 24),
          _buildDetailRow('Purpose', _getPurposeLabel(loan.purpose)),
          const Divider(height: 24),
          _buildDetailRow(
            'Disbursed On',
            loan.disbursementDate != null
                ? AppDateUtils.formatDate(loan.disbursementDate!)
                : 'Pending',
          ),
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Row(
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
          style: AppTypography.bodyMedium.copyWith(
            fontWeight: FontWeight.w500,
          ),
        ),
      ],
    );
  }

  Widget _buildRepaymentSchedule(List<_RepaymentItem> schedule) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Repayment Schedule',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 16),
          ...schedule.asMap().entries.map((entry) {
            final index = entry.key;
            final item = entry.value;
            final isLast = index == schedule.length - 1;
            
            return _buildScheduleItem(item, isLast);
          }),
        ],
      ),
    );
  }

  Widget _buildScheduleItem(_RepaymentItem item, bool isLast) {
    IconData icon;
    Color iconColor;
    Color bgColor;

    switch (item.status) {
      case 'paid':
        icon = Icons.check_circle;
        iconColor = AppColors.success;
        bgColor = AppColors.success.withOpacity(0.1);
        break;
      case 'upcoming':
        icon = Icons.schedule;
        iconColor = AppColors.warning;
        bgColor = AppColors.warning.withOpacity(0.1);
        break;
      default:
        icon = Icons.circle_outlined;
        iconColor = AppColors.textSecondary;
        bgColor = AppColors.background;
    }

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Column(
          children: [
            Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: bgColor,
                shape: BoxShape.circle,
              ),
              child: Icon(icon, color: iconColor, size: 20),
            ),
            if (!isLast)
              Container(
                width: 2,
                height: 40,
                color: AppColors.border,
              ),
          ],
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Padding(
            padding: EdgeInsets.only(bottom: isLast ? 0 : 16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      CurrencyUtils.formatNaira(item.amount),
                      style: AppTypography.titleSmall.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    _buildScheduleStatusBadge(item.status),
                  ],
                ),
                const SizedBox(height: 4),
                Text(
                  'Due: ${AppDateUtils.formatDate(item.dueDate)}',
                  style: AppTypography.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
                if (item.paidDate != null) ...[
                  const SizedBox(height: 2),
                  Text(
                    'Paid: ${AppDateUtils.formatDate(item.paidDate!)}',
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.success,
                    ),
                  ),
                ],
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildScheduleStatusBadge(String status) {
    Color color;
    String label;

    switch (status) {
      case 'paid':
        color = AppColors.success;
        label = 'Paid';
        break;
      case 'upcoming':
        color = AppColors.warning;
        label = 'Upcoming';
        break;
      default:
        color = AppColors.textSecondary;
        label = 'Pending';
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        label,
        style: AppTypography.labelSmall.copyWith(
          color: color,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Widget _buildActionsSection(BuildContext context, Loan loan) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        children: [
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: loan.status == LoanStatus.active
                  ? () => _showRepaymentDialog(context, loan)
                  : null,
              icon: const Icon(Icons.payment),
              label: const Text('Make Repayment'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 14),
              ),
            ),
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(
                child: OutlinedButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.receipt_long, size: 18),
                  label: const Text('Statement'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: OutlinedButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.support_agent, size: 18),
                  label: const Text('Get Help'),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _getPurposeLabel(LoanPurpose purpose) {
    switch (purpose) {
      case LoanPurpose.emergency:
        return 'Emergency';
      case LoanPurpose.business:
        return 'Business';
      case LoanPurpose.education:
        return 'Education';
      case LoanPurpose.medical:
        return 'Medical';
      case LoanPurpose.personal:
        return 'Personal';
      case LoanPurpose.rent:
        return 'Rent';
      case LoanPurpose.other:
        return 'Other';
    }
  }

  void _showRepaymentDialog(BuildContext context, Loan loan) {
    final remaining = loan.totalAmount - loan.amountRepaid;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _RepaymentSheet(
        remainingAmount: remaining,
        loanId: loan.id,
      ),
    );
  }

  void _showLoanHelp(BuildContext context) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Container(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Loan FAQ',
              style: AppTypography.titleLarge.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            _buildFAQItem(
              'How do I make a repayment?',
              'Tap the "Make Repayment" button and choose to pay the full amount or a partial payment from your wallet.',
            ),
            _buildFAQItem(
              'What happens if I miss a payment?',
              'Late payments may affect your credit score and incur additional charges. Contact support if you\'re having difficulties.',
            ),
            _buildFAQItem(
              'Can I pay off early?',
              'Yes! You can pay off your loan early at any time with no prepayment penalties.',
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  Widget _buildFAQItem(String question, String answer) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            question,
            style: AppTypography.titleSmall.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            answer,
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}

class _RepaymentItem {
  final double amount;
  final DateTime dueDate;
  final String status;
  final DateTime? paidDate;

  _RepaymentItem({
    required this.amount,
    required this.dueDate,
    required this.status,
    this.paidDate,
  });
}

class _RepaymentSheet extends StatefulWidget {
  final double remainingAmount;
  final String loanId;

  const _RepaymentSheet({
    required this.remainingAmount,
    required this.loanId,
  });

  @override
  State<_RepaymentSheet> createState() => _RepaymentSheetState();
}

class _RepaymentSheetState extends State<_RepaymentSheet> {
  bool _payFullAmount = true;
  final _amountController = TextEditingController();
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _amountController.text = widget.remainingAmount.toStringAsFixed(0);
  }

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(
        left: 16,
        right: 16,
        top: 16,
        bottom: MediaQuery.of(context).viewInsets.bottom + 16,
      ),
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
          const SizedBox(height: 16),
          Text(
            'Make Repayment',
            style: AppTypography.titleLarge.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Outstanding balance: ${CurrencyUtils.formatNaira(widget.remainingAmount)}',
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 24),

          // Payment options
          _buildPaymentOption(
            'Pay Full Amount',
            CurrencyUtils.formatNaira(widget.remainingAmount),
            _payFullAmount,
            () => setState(() {
              _payFullAmount = true;
              _amountController.text = widget.remainingAmount.toStringAsFixed(0);
            }),
          ),
          const SizedBox(height: 12),
          _buildPaymentOption(
            'Pay Custom Amount',
            'Enter amount',
            !_payFullAmount,
            () => setState(() => _payFullAmount = false),
          ),

          if (!_payFullAmount) ...[
            const SizedBox(height: 16),
            TextField(
              controller: _amountController,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(
                labelText: 'Amount',
                prefixText: '₦ ',
                border: const OutlineInputBorder(),
                helperText: 'Min: ₦1,000 | Max: ${CurrencyUtils.formatNaira(widget.remainingAmount)}',
              ),
            ),
          ],

          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _isLoading ? null : _handleRepayment,
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: _isLoading
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    )
                  : Text(
                      'Pay ${_payFullAmount ? CurrencyUtils.formatNaira(widget.remainingAmount) : ''}',
                    ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPaymentOption(
    String title,
    String subtitle,
    bool isSelected,
    VoidCallback onTap,
  ) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected
              ? AppColors.primary.withOpacity(0.05)
              : AppColors.surface,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isSelected ? AppColors.primary : AppColors.border,
          ),
        ),
        child: Row(
          children: [
            Icon(
              isSelected ? Icons.radio_button_checked : Icons.radio_button_off,
              color: isSelected ? AppColors.primary : AppColors.textSecondary,
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: AppTypography.titleSmall.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    subtitle,
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _handleRepayment() async {
    final amount = double.tryParse(_amountController.text) ?? 0;
    
    if (amount < 1000) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Minimum amount is ₦1,000')),
      );
      return;
    }

    if (amount > widget.remainingAmount) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Amount cannot exceed ${CurrencyUtils.formatNaira(widget.remainingAmount)}'),
        ),
      );
      return;
    }

    setState(() => _isLoading = true);

    // Simulate API call
    await Future.delayed(const Duration(seconds: 2));

    if (mounted) {
      Navigator.pop(context);
      
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: AppColors.success.withOpacity(0.1),
                  shape: BoxShape.circle,
                ),
                child: const Icon(
                  Icons.check_circle,
                  color: AppColors.success,
                  size: 24,
                ),
              ),
              const SizedBox(width: 12),
              const Text('Payment Successful'),
            ],
          ),
          content: Text(
            'Your payment of ${CurrencyUtils.formatNaira(amount)} has been processed successfully.',
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Done'),
            ),
          ],
        ),
      );
    }
  }
}
