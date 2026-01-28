import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/utils/currency_helper.dart';
import '../../../../core/widgets/buttons.dart';
import '../../../../core/widgets/loaders.dart';
import '../../domain/models/wallet_models.dart';
import '../providers/wallet_provider.dart';

/// =============================================================================
/// TRANSACTION DETAILS SCREEN
/// =============================================================================

class TransactionDetailsScreen extends ConsumerWidget {
  final String transactionId;

  const TransactionDetailsScreen({
    super.key,
    required this.transactionId,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final transactionAsync = ref.watch(transactionDetailProvider(transactionId));

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Transaction Details'),
        backgroundColor: AppColors.background,
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(Icons.share_outlined),
            onPressed: () => _shareTransaction(context),
          ),
        ],
      ),
      body: transactionAsync.when(
        loading: () => const Center(child: AppLoader()),
        error: (error, _) => _buildError(context, error.toString(), ref),
        data: (transaction) {
          if (transaction == null) {
            return _buildError(context, 'Transaction not found', ref);
          }
          return _buildContent(context, transaction);
        },
      ),
    );
  }

  Widget _buildError(BuildContext context, String message, WidgetRef ref) {
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
              message,
              style: AppTypography.bodyLarge,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            AppButton(
              text: 'Try Again',
              onPressed: () => ref.refresh(transactionDetailProvider(transactionId)),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildContent(BuildContext context, Transaction transaction) {
    final isCredit = transaction.type == TransactionType.credit ||
        transaction.type == TransactionType.deposit ||
        transaction.type == TransactionType.payout;

    return SingleChildScrollView(
      child: Column(
        children: [
          // Status Header
          _buildStatusHeader(transaction, isCredit),

          const SizedBox(height: 24),

          // Amount
          _buildAmountSection(transaction, isCredit),

          const SizedBox(height: 24),

          // Details Card
          _buildDetailsCard(transaction),

          const SizedBox(height: 16),

          // Timeline
          if (transaction.timeline != null && transaction.timeline!.isNotEmpty)
            _buildTimelineCard(transaction.timeline!),

          const SizedBox(height: 16),

          // Actions
          _buildActions(context, transaction),

          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildStatusHeader(Transaction transaction, bool isCredit) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(vertical: 32),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            _getStatusColor(transaction.status).withOpacity(0.1),
            AppColors.background,
          ],
        ),
      ),
      child: Column(
        children: [
          // Icon
          Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              color: _getStatusColor(transaction.status).withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: Icon(
              _getStatusIcon(transaction.status),
              size: 40,
              color: _getStatusColor(transaction.status),
            ),
          ),

          const SizedBox(height: 16),

          // Status Badge
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: BoxDecoration(
              color: _getStatusColor(transaction.status).withOpacity(0.1),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Text(
              _getStatusText(transaction.status),
              style: AppTypography.labelMedium.copyWith(
                color: _getStatusColor(transaction.status),
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAmountSection(Transaction transaction, bool isCredit) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24),
      child: Column(
        children: [
          Text(
            '${isCredit ? '+' : '-'}${CurrencyHelper.formatNaira(transaction.amount)}',
            style: AppTypography.displayMedium.copyWith(
              color: isCredit ? AppColors.success : AppColors.error,
              fontWeight: FontWeight.bold,
            ),
          ),

          if (transaction.fee != null && transaction.fee! > 0) ...[
            const SizedBox(height: 8),
            Text(
              'Fee: ${CurrencyHelper.formatNaira(transaction.fee!)}',
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          ],

          const SizedBox(height: 4),
          Text(
            _formatTransactionType(transaction.type),
            style: AppTypography.bodyLarge.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDetailsCard(Transaction transaction) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        children: [
          _buildDetailRow(
            'Reference',
            transaction.reference,
            showCopy: true,
          ),
          const Divider(height: 1),

          _buildDetailRow(
            'Date',
            DateFormat('MMMM d, yyyy').format(transaction.createdAt),
          ),
          const Divider(height: 1),

          _buildDetailRow(
            'Time',
            DateFormat('h:mm a').format(transaction.createdAt),
          ),

          if (transaction.description != null) ...[
            const Divider(height: 1),
            _buildDetailRow('Description', transaction.description!),
          ],

          if (transaction.category != null) ...[
            const Divider(height: 1),
            _buildDetailRow('Category', transaction.category!),
          ],

          if (transaction.counterparty != null) ...[
            const Divider(height: 1),
            _buildDetailRow(
              transaction.type == TransactionType.transfer ? 'Sent to' : 'From',
              transaction.counterparty!,
            ),
          ],

          if (transaction.bankName != null) ...[
            const Divider(height: 1),
            _buildDetailRow('Bank', transaction.bankName!),
          ],

          if (transaction.accountNumber != null) ...[
            const Divider(height: 1),
            _buildDetailRow(
              'Account',
              _maskAccountNumber(transaction.accountNumber!),
            ),
          ],

          if (transaction.balanceAfter != null) ...[
            const Divider(height: 1),
            _buildDetailRow(
              'Balance After',
              CurrencyHelper.formatNaira(transaction.balanceAfter!),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value, {bool showCopy = false}) {
    return Builder(
      builder: (context) => Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SizedBox(
              width: 120,
              child: Text(
                label,
                style: AppTypography.bodyMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
            ),
            Expanded(
              child: Text(
                value,
                style: AppTypography.bodyMedium.copyWith(
                  fontWeight: FontWeight.w500,
                ),
                textAlign: TextAlign.right,
              ),
            ),
            if (showCopy) ...[
              const SizedBox(width: 8),
              GestureDetector(
                onTap: () => _copyToClipboard(context, value),
                child: const Icon(
                  Icons.copy_outlined,
                  size: 18,
                  color: AppColors.primary,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildTimelineCard(List<TransactionEvent> timeline) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Transaction Timeline',
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 16),
          ...timeline.asMap().entries.map((entry) {
            final index = entry.key;
            final event = entry.value;
            final isLast = index == timeline.length - 1;

            return _buildTimelineItem(event, isLast);
          }),
        ],
      ),
    );
  }

  Widget _buildTimelineItem(TransactionEvent event, bool isLast) {
    return IntrinsicHeight(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Timeline dot and line
          Column(
            children: [
              Container(
                width: 12,
                height: 12,
                decoration: BoxDecoration(
                  color: _getEventColor(event.status),
                  shape: BoxShape.circle,
                ),
              ),
              if (!isLast)
                Expanded(
                  child: Container(
                    width: 2,
                    color: AppColors.border,
                  ),
                ),
            ],
          ),

          const SizedBox(width: 12),

          // Event content
          Expanded(
            child: Padding(
              padding: EdgeInsets.only(bottom: isLast ? 0 : 16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    event.title,
                    style: AppTypography.bodyMedium.copyWith(
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    DateFormat('MMM d, h:mm a').format(event.timestamp),
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                  if (event.description != null) ...[
                    const SizedBox(height: 4),
                    Text(
                      event.description!,
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActions(BuildContext context, Transaction transaction) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        children: [
          // Download receipt
          if (transaction.status == TransactionStatus.completed)
            _buildActionButton(
              icon: Icons.receipt_long_outlined,
              label: 'Download Receipt',
              onTap: () => _downloadReceipt(context, transaction),
            ),

          const SizedBox(height: 12),

          // Report issue
          _buildActionButton(
            icon: Icons.report_problem_outlined,
            label: 'Report an Issue',
            onTap: () => _reportIssue(context, transaction),
          ),

          // Retry for failed transactions
          if (transaction.status == TransactionStatus.failed) ...[
            const SizedBox(height: 12),
            _buildActionButton(
              icon: Icons.refresh,
              label: 'Retry Transaction',
              onTap: () => _retryTransaction(context, transaction),
              isPrimary: true,
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildActionButton({
    required IconData icon,
    required String label,
    required VoidCallback onTap,
    bool isPrimary = false,
  }) {
    return Material(
      color: isPrimary ? AppColors.primary : AppColors.surface,
      borderRadius: BorderRadius.circular(12),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 16),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                icon,
                size: 20,
                color: isPrimary ? Colors.white : AppColors.textPrimary,
              ),
              const SizedBox(width: 8),
              Text(
                label,
                style: AppTypography.bodyMedium.copyWith(
                  fontWeight: FontWeight.w500,
                  color: isPrimary ? Colors.white : AppColors.textPrimary,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // Helper methods
  Color _getStatusColor(TransactionStatus status) {
    switch (status) {
      case TransactionStatus.completed:
        return AppColors.success;
      case TransactionStatus.pending:
        return AppColors.warning;
      case TransactionStatus.processing:
        return AppColors.info;
      case TransactionStatus.failed:
        return AppColors.error;
      case TransactionStatus.cancelled:
        return AppColors.textSecondary;
    }
  }

  IconData _getStatusIcon(TransactionStatus status) {
    switch (status) {
      case TransactionStatus.completed:
        return Icons.check_circle;
      case TransactionStatus.pending:
        return Icons.schedule;
      case TransactionStatus.processing:
        return Icons.sync;
      case TransactionStatus.failed:
        return Icons.error;
      case TransactionStatus.cancelled:
        return Icons.cancel;
    }
  }

  String _getStatusText(TransactionStatus status) {
    switch (status) {
      case TransactionStatus.completed:
        return 'Completed';
      case TransactionStatus.pending:
        return 'Pending';
      case TransactionStatus.processing:
        return 'Processing';
      case TransactionStatus.failed:
        return 'Failed';
      case TransactionStatus.cancelled:
        return 'Cancelled';
    }
  }

  Color _getEventColor(String status) {
    switch (status.toLowerCase()) {
      case 'completed':
      case 'success':
        return AppColors.success;
      case 'pending':
        return AppColors.warning;
      case 'failed':
      case 'error':
        return AppColors.error;
      default:
        return AppColors.textSecondary;
    }
  }

  String _formatTransactionType(TransactionType type) {
    switch (type) {
      case TransactionType.deposit:
        return 'Wallet Deposit';
      case TransactionType.withdrawal:
        return 'Withdrawal';
      case TransactionType.transfer:
        return 'Transfer';
      case TransactionType.payment:
        return 'Payment';
      case TransactionType.refund:
        return 'Refund';
      case TransactionType.credit:
        return 'Credit';
      case TransactionType.debit:
        return 'Debit';
      case TransactionType.contribution:
        return 'Circle Contribution';
      case TransactionType.payout:
        return 'Circle Payout';
      case TransactionType.loanDisbursement:
        return 'Loan Disbursement';
      case TransactionType.loanRepayment:
        return 'Loan Repayment';
      case TransactionType.gigPayment:
        return 'Gig Payment';
      case TransactionType.gigEarning:
        return 'Gig Earning';
    }
  }

  String _maskAccountNumber(String accountNumber) {
    if (accountNumber.length <= 4) return accountNumber;
    final masked = '*' * (accountNumber.length - 4);
    return masked + accountNumber.substring(accountNumber.length - 4);
  }

  void _copyToClipboard(BuildContext context, String text) {
    Clipboard.setData(ClipboardData(text: text));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Copied to clipboard'),
        duration: Duration(seconds: 2),
      ),
    );
  }

  void _shareTransaction(BuildContext context) {
    // TODO: Implement share functionality
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Share feature coming soon')),
    );
  }

  void _downloadReceipt(BuildContext context, Transaction transaction) {
    // TODO: Implement receipt download
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Receipt download coming soon')),
    );
  }

  void _reportIssue(BuildContext context, Transaction transaction) {
    // TODO: Navigate to support/report screen
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Report feature coming soon')),
    );
  }

  void _retryTransaction(BuildContext context, Transaction transaction) {
    // TODO: Implement retry logic based on transaction type
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Retry feature coming soon')),
    );
  }
}

/// =============================================================================
/// TRANSACTION EVENT MODEL
/// =============================================================================

class TransactionEvent {
  final String title;
  final String status;
  final DateTime timestamp;
  final String? description;

  TransactionEvent({
    required this.title,
    required this.status,
    required this.timestamp,
    this.description,
  });

  factory TransactionEvent.fromJson(Map<String, dynamic> json) {
    return TransactionEvent(
      title: json['title'] as String,
      status: json['status'] as String,
      timestamp: DateTime.parse(json['timestamp'] as String),
      description: json['description'] as String?,
    );
  }
}
