import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/utils/currency_utils.dart';
import '../../../../core/utils/date_utils.dart';
import '../../data/models/wallet_models.dart';

// Filter providers
final transactionFilterProvider = StateProvider<TransactionCategory?>((ref) => null);
final transactionTypeFilterProvider = StateProvider<TransactionType?>((ref) => null);
final transactionSearchProvider = StateProvider<String>((ref) => '');

// Mock transactions provider
final transactionsProvider = FutureProvider<List<Transaction>>((ref) async {
  await Future.delayed(const Duration(milliseconds: 500));
  
  return [
    Transaction(
      id: 'txn-1',
      userId: 'user-1',
      type: TransactionType.credit,
      category: TransactionCategory.gigPayment,
      amount: 75000,
      balanceBefore: 125000,
      balanceAfter: 200000,
      reference: 'TXN-2024-001234',
      description: 'Payment for Logo Design',
      status: TransactionStatus.completed,
      createdAt: DateTime.now().subtract(const Duration(hours: 2)),
      completedAt: DateTime.now().subtract(const Duration(hours: 2)),
    ),
    Transaction(
      id: 'txn-2',
      userId: 'user-1',
      type: TransactionType.debit,
      category: TransactionCategory.circleContribution,
      amount: 10000,
      balanceBefore: 200000,
      balanceAfter: 190000,
      reference: 'TXN-2024-001235',
      description: 'Weekly contribution - Tech Hustlers',
      status: TransactionStatus.completed,
      createdAt: DateTime.now().subtract(const Duration(days: 1)),
      completedAt: DateTime.now().subtract(const Duration(days: 1)),
    ),
    Transaction(
      id: 'txn-3',
      userId: 'user-1',
      type: TransactionType.credit,
      category: TransactionCategory.funding,
      amount: 50000,
      balanceBefore: 140000,
      balanceAfter: 190000,
      reference: 'TXN-2024-001236',
      description: 'Wallet funding via Card',
      status: TransactionStatus.completed,
      createdAt: DateTime.now().subtract(const Duration(days: 2)),
      completedAt: DateTime.now().subtract(const Duration(days: 2)),
    ),
    Transaction(
      id: 'txn-4',
      userId: 'user-1',
      type: TransactionType.debit,
      category: TransactionCategory.transfer,
      amount: 25000,
      balanceBefore: 165000,
      balanceAfter: 140000,
      reference: 'TXN-2024-001237',
      description: 'Transfer to @johndoe',
      status: TransactionStatus.completed,
      createdAt: DateTime.now().subtract(const Duration(days: 3)),
      completedAt: DateTime.now().subtract(const Duration(days: 3)),
    ),
    Transaction(
      id: 'txn-5',
      userId: 'user-1',
      type: TransactionType.credit,
      category: TransactionCategory.circlePayout,
      amount: 120000,
      balanceBefore: 45000,
      balanceAfter: 165000,
      reference: 'TXN-2024-001238',
      description: 'Circle payout - Women in Tech',
      status: TransactionStatus.completed,
      createdAt: DateTime.now().subtract(const Duration(days: 5)),
      completedAt: DateTime.now().subtract(const Duration(days: 5)),
    ),
    Transaction(
      id: 'txn-6',
      userId: 'user-1',
      type: TransactionType.debit,
      category: TransactionCategory.loanRepayment,
      amount: 15000,
      balanceBefore: 60000,
      balanceAfter: 45000,
      reference: 'TXN-2024-001239',
      description: 'Loan repayment - Week 2',
      status: TransactionStatus.completed,
      createdAt: DateTime.now().subtract(const Duration(days: 7)),
      completedAt: DateTime.now().subtract(const Duration(days: 7)),
    ),
    Transaction(
      id: 'txn-7',
      userId: 'user-1',
      type: TransactionType.credit,
      category: TransactionCategory.loanDisbursement,
      amount: 100000,
      balanceBefore: 0,
      balanceAfter: 100000,
      reference: 'TXN-2024-001240',
      description: 'Loan disbursement',
      status: TransactionStatus.completed,
      createdAt: DateTime.now().subtract(const Duration(days: 14)),
      completedAt: DateTime.now().subtract(const Duration(days: 14)),
    ),
    Transaction(
      id: 'txn-8',
      userId: 'user-1',
      type: TransactionType.debit,
      category: TransactionCategory.withdrawal,
      amount: 35000,
      balanceBefore: 35000,
      balanceAfter: 0,
      reference: 'TXN-2024-001241',
      description: 'Withdrawal to Access Bank',
      status: TransactionStatus.completed,
      createdAt: DateTime.now().subtract(const Duration(days: 15)),
      completedAt: DateTime.now().subtract(const Duration(days: 15)),
    ),
  ];
});

// Filtered transactions provider
final filteredTransactionsProvider = Provider<AsyncValue<List<Transaction>>>((ref) {
  final transactionsAsync = ref.watch(transactionsProvider);
  final categoryFilter = ref.watch(transactionFilterProvider);
  final typeFilter = ref.watch(transactionTypeFilterProvider);
  final searchQuery = ref.watch(transactionSearchProvider).toLowerCase();

  return transactionsAsync.whenData((transactions) {
    var filtered = transactions;

    if (categoryFilter != null) {
      filtered = filtered.where((t) => t.category == categoryFilter).toList();
    }

    if (typeFilter != null) {
      filtered = filtered.where((t) => t.type == typeFilter).toList();
    }

    if (searchQuery.isNotEmpty) {
      filtered = filtered.where((t) =>
        t.description.toLowerCase().contains(searchQuery) ||
        t.reference.toLowerCase().contains(searchQuery)
      ).toList();
    }

    return filtered;
  });
});

class TransactionsScreen extends ConsumerStatefulWidget {
  const TransactionsScreen({super.key});

  @override
  ConsumerState<TransactionsScreen> createState() => _TransactionsScreenState();
}

class _TransactionsScreenState extends ConsumerState<TransactionsScreen> {
  final _searchController = TextEditingController();
  bool _showSearch = false;

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final transactionsAsync = ref.watch(filteredTransactionsProvider);
    final categoryFilter = ref.watch(transactionFilterProvider);
    final typeFilter = ref.watch(transactionTypeFilterProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: _showSearch
            ? TextField(
                controller: _searchController,
                autofocus: true,
                decoration: const InputDecoration(
                  hintText: 'Search transactions...',
                  border: InputBorder.none,
                ),
                onChanged: (value) {
                  ref.read(transactionSearchProvider.notifier).state = value;
                },
              )
            : const Text('Transactions'),
        actions: [
          IconButton(
            icon: Icon(_showSearch ? Icons.close : Icons.search),
            onPressed: () {
              setState(() {
                _showSearch = !_showSearch;
                if (!_showSearch) {
                  _searchController.clear();
                  ref.read(transactionSearchProvider.notifier).state = '';
                }
              });
            },
          ),
          IconButton(
            icon: const Icon(Icons.filter_list),
            onPressed: () => _showFilterSheet(context),
          ),
        ],
      ),
      body: Column(
        children: [
          // Active filters display
          if (categoryFilter != null || typeFilter != null)
            _buildActiveFilters(categoryFilter, typeFilter),
          
          // Transactions list
          Expanded(
            child: transactionsAsync.when(
              data: (transactions) {
                if (transactions.isEmpty) {
                  return _buildEmptyState();
                }
                return _buildTransactionsList(transactions);
              },
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (e, _) => Center(child: Text('Error: $e')),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActiveFilters(
    TransactionCategory? categoryFilter,
    TransactionType? typeFilter,
  ) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: AppColors.surface,
      child: Row(
        children: [
          Text(
            'Filters: ',
            style: AppTypography.bodySmall.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          if (typeFilter != null)
            _buildFilterChip(
              typeFilter == TransactionType.credit ? 'Credits' : 'Debits',
              () => ref.read(transactionTypeFilterProvider.notifier).state = null,
            ),
          if (categoryFilter != null)
            _buildFilterChip(
              _getCategoryLabel(categoryFilter),
              () => ref.read(transactionFilterProvider.notifier).state = null,
            ),
          const Spacer(),
          TextButton(
            onPressed: () {
              ref.read(transactionFilterProvider.notifier).state = null;
              ref.read(transactionTypeFilterProvider.notifier).state = null;
            },
            child: const Text('Clear All'),
          ),
        ],
      ),
    );
  }

  Widget _buildFilterChip(String label, VoidCallback onRemove) {
    return Container(
      margin: const EdgeInsets.only(right: 8),
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: AppColors.primary.withOpacity(0.1),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            label,
            style: AppTypography.labelSmall.copyWith(
              color: AppColors.primary,
            ),
          ),
          const SizedBox(width: 4),
          GestureDetector(
            onTap: onRemove,
            child: Icon(
              Icons.close,
              size: 14,
              color: AppColors.primary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.receipt_long_outlined,
              size: 80,
              color: AppColors.textSecondary.withOpacity(0.5),
            ),
            const SizedBox(height: 16),
            Text(
              'No Transactions',
              style: AppTypography.titleLarge.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Your transactions will appear here',
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTransactionsList(List<Transaction> transactions) {
    // Group transactions by date
    final groupedTransactions = <String, List<Transaction>>{};
    
    for (final transaction in transactions) {
      final dateKey = AppDateUtils.formatDate(transaction.createdAt);
      groupedTransactions.putIfAbsent(dateKey, () => []).add(transaction);
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: groupedTransactions.length,
      itemBuilder: (context, index) {
        final date = groupedTransactions.keys.elementAt(index);
        final dayTransactions = groupedTransactions[date]!;

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 8),
              child: Text(
                date,
                style: AppTypography.labelMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
            ),
            ...dayTransactions.map((transaction) => _TransactionTile(
              transaction: transaction,
              onTap: () => context.push('/wallet/transactions/${transaction.id}'),
            )),
            if (index < groupedTransactions.length - 1)
              const SizedBox(height: 8),
          ],
        );
      },
    );
  }

  void _showFilterSheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _FilterSheet(),
    );
  }

  String _getCategoryLabel(TransactionCategory category) {
    switch (category) {
      case TransactionCategory.funding:
        return 'Funding';
      case TransactionCategory.withdrawal:
        return 'Withdrawal';
      case TransactionCategory.transfer:
        return 'Transfer';
      case TransactionCategory.gigPayment:
        return 'Gig Payment';
      case TransactionCategory.circleContribution:
        return 'Circle Contribution';
      case TransactionCategory.circlePayout:
        return 'Circle Payout';
      case TransactionCategory.loanDisbursement:
        return 'Loan Disbursement';
      case TransactionCategory.loanRepayment:
        return 'Loan Repayment';
      case TransactionCategory.fee:
        return 'Fee';
      case TransactionCategory.refund:
        return 'Refund';
      case TransactionCategory.other:
        return 'Other';
    }
  }
}

class _TransactionTile extends StatelessWidget {
  final Transaction transaction;
  final VoidCallback onTap;

  const _TransactionTile({
    required this.transaction,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final isCredit = transaction.type == TransactionType.credit;
    final amountColor = isCredit ? AppColors.success : AppColors.error;
    final prefix = isCredit ? '+' : '-';

    return GestureDetector(
      onTap: onTap,
      child: Container(
        margin: const EdgeInsets.only(bottom: 8),
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: AppColors.border),
        ),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: _getCategoryColor(transaction.category).withOpacity(0.1),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                _getCategoryIcon(transaction.category),
                color: _getCategoryColor(transaction.category),
                size: 20,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    transaction.description,
                    style: AppTypography.titleSmall.copyWith(
                      fontWeight: FontWeight.w500,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 2),
                  Text(
                    AppDateUtils.formatTime(transaction.createdAt),
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '$prefix${CurrencyUtils.formatNaira(transaction.amount)}',
                  style: AppTypography.titleSmall.copyWith(
                    color: amountColor,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 2),
                _buildStatusDot(transaction.status),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatusDot(TransactionStatus status) {
    Color color;
    switch (status) {
      case TransactionStatus.completed:
        color = AppColors.success;
        break;
      case TransactionStatus.pending:
        color = AppColors.warning;
        break;
      case TransactionStatus.failed:
        color = AppColors.error;
        break;
      case TransactionStatus.reversed:
        color = AppColors.textSecondary;
        break;
    }

    return Container(
      width: 8,
      height: 8,
      decoration: BoxDecoration(
        color: color,
        shape: BoxShape.circle,
      ),
    );
  }

  IconData _getCategoryIcon(TransactionCategory category) {
    switch (category) {
      case TransactionCategory.funding:
        return Icons.add_circle;
      case TransactionCategory.withdrawal:
        return Icons.arrow_downward;
      case TransactionCategory.transfer:
        return Icons.swap_horiz;
      case TransactionCategory.gigPayment:
        return Icons.work;
      case TransactionCategory.circleContribution:
        return Icons.people;
      case TransactionCategory.circlePayout:
        return Icons.celebration;
      case TransactionCategory.loanDisbursement:
        return Icons.account_balance;
      case TransactionCategory.loanRepayment:
        return Icons.payments;
      case TransactionCategory.fee:
        return Icons.receipt;
      case TransactionCategory.refund:
        return Icons.undo;
      case TransactionCategory.other:
        return Icons.receipt_long;
    }
  }

  Color _getCategoryColor(TransactionCategory category) {
    switch (category) {
      case TransactionCategory.funding:
        return AppColors.success;
      case TransactionCategory.withdrawal:
        return AppColors.warning;
      case TransactionCategory.transfer:
        return AppColors.info;
      case TransactionCategory.gigPayment:
        return AppColors.primary;
      case TransactionCategory.circleContribution:
        return AppColors.secondary;
      case TransactionCategory.circlePayout:
        return AppColors.success;
      case TransactionCategory.loanDisbursement:
        return AppColors.info;
      case TransactionCategory.loanRepayment:
        return AppColors.warning;
      case TransactionCategory.fee:
        return AppColors.error;
      case TransactionCategory.refund:
        return AppColors.textSecondary;
      case TransactionCategory.other:
        return AppColors.textSecondary;
    }
  }
}

class _FilterSheet extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final typeFilter = ref.watch(transactionTypeFilterProvider);
    final categoryFilter = ref.watch(transactionFilterProvider);

    return Container(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Filter Transactions',
                style: AppTypography.titleLarge.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              TextButton(
                onPressed: () {
                  ref.read(transactionFilterProvider.notifier).state = null;
                  ref.read(transactionTypeFilterProvider.notifier).state = null;
                  Navigator.pop(context);
                },
                child: const Text('Reset'),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // Transaction Type
          Text(
            'Transaction Type',
            style: AppTypography.titleSmall.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              Expanded(
                child: _FilterOption(
                  label: 'All',
                  isSelected: typeFilter == null,
                  onTap: () {
                    ref.read(transactionTypeFilterProvider.notifier).state = null;
                  },
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: _FilterOption(
                  label: 'Credits',
                  isSelected: typeFilter == TransactionType.credit,
                  onTap: () {
                    ref.read(transactionTypeFilterProvider.notifier).state =
                        TransactionType.credit;
                  },
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: _FilterOption(
                  label: 'Debits',
                  isSelected: typeFilter == TransactionType.debit,
                  onTap: () {
                    ref.read(transactionTypeFilterProvider.notifier).state =
                        TransactionType.debit;
                  },
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // Category
          Text(
            'Category',
            style: AppTypography.titleSmall.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: TransactionCategory.values.map((category) {
              final isSelected = categoryFilter == category;
              return GestureDetector(
                onTap: () {
                  ref.read(transactionFilterProvider.notifier).state =
                      isSelected ? null : category;
                },
                child: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 8,
                  ),
                  decoration: BoxDecoration(
                    color: isSelected
                        ? AppColors.primary.withOpacity(0.1)
                        : AppColors.background,
                    borderRadius: BorderRadius.circular(20),
                    border: Border.all(
                      color: isSelected ? AppColors.primary : AppColors.border,
                    ),
                  ),
                  child: Text(
                    _getCategoryLabel(category),
                    style: AppTypography.labelMedium.copyWith(
                      color: isSelected
                          ? AppColors.primary
                          : AppColors.textSecondary,
                    ),
                  ),
                ),
              );
            }).toList(),
          ),
          const SizedBox(height: 24),

          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Apply Filters'),
            ),
          ),
        ],
      ),
    );
  }

  String _getCategoryLabel(TransactionCategory category) {
    switch (category) {
      case TransactionCategory.funding:
        return 'Funding';
      case TransactionCategory.withdrawal:
        return 'Withdrawal';
      case TransactionCategory.transfer:
        return 'Transfer';
      case TransactionCategory.gigPayment:
        return 'Gig Payment';
      case TransactionCategory.circleContribution:
        return 'Circle';
      case TransactionCategory.circlePayout:
        return 'Payout';
      case TransactionCategory.loanDisbursement:
        return 'Loan';
      case TransactionCategory.loanRepayment:
        return 'Repayment';
      case TransactionCategory.fee:
        return 'Fee';
      case TransactionCategory.refund:
        return 'Refund';
      case TransactionCategory.other:
        return 'Other';
    }
  }
}

class _FilterOption extends StatelessWidget {
  final String label;
  final bool isSelected;
  final VoidCallback onTap;

  const _FilterOption({
    required this.label,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: isSelected
              ? AppColors.primary.withOpacity(0.1)
              : AppColors.background,
          borderRadius: BorderRadius.circular(8),
          border: Border.all(
            color: isSelected ? AppColors.primary : AppColors.border,
          ),
        ),
        alignment: Alignment.center,
        child: Text(
          label,
          style: AppTypography.labelMedium.copyWith(
            color: isSelected ? AppColors.primary : AppColors.textSecondary,
            fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
          ),
        ),
      ),
    );
  }
}
