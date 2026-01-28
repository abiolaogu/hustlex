import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_constants.dart';

class WalletScreen extends ConsumerStatefulWidget {
  const WalletScreen({super.key});

  @override
  ConsumerState<WalletScreen> createState() => _WalletScreenState();
}

class _WalletScreenState extends ConsumerState<WalletScreen> {
  bool _isBalanceVisible = true;
  String _selectedFilter = 'All';
  
  final _currencyFormat = NumberFormat.currency(
    locale: 'en_NG',
    symbol: '₦',
    decimalDigits: 2,
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Wallet'),
        actions: [
          IconButton(
            onPressed: () {},
            icon: const Icon(Icons.history_rounded),
          ),
          IconButton(
            onPressed: () {},
            icon: const Icon(Icons.more_vert_rounded),
          ),
        ],
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            // Balance card
            Container(
              margin: const EdgeInsets.all(16),
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                gradient: AppColors.primaryGradient,
                borderRadius: BorderRadius.circular(24),
                boxShadow: AppShadows.colored,
              ),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Available Balance',
                        style: AppTypography.bodyMedium.copyWith(
                          color: Colors.white.withOpacity(0.8),
                        ),
                      ),
                      IconButton(
                        onPressed: () {
                          setState(() => _isBalanceVisible = !_isBalanceVisible);
                        },
                        icon: Icon(
                          _isBalanceVisible
                              ? Icons.visibility_outlined
                              : Icons.visibility_off_outlined,
                          color: Colors.white,
                          size: 20,
                        ),
                      ),
                    ],
                  ),
                  Text(
                    _isBalanceVisible
                        ? _currencyFormat.format(125750.50)
                        : '₦ ••••••',
                    style: AppTypography.displaySmall.copyWith(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      _buildBalanceItem('Escrow', 45000),
                      Container(
                        width: 1,
                        height: 30,
                        color: Colors.white.withOpacity(0.3),
                        margin: const EdgeInsets.symmetric(horizontal: 24),
                      ),
                      _buildBalanceItem('Savings', 75000),
                    ],
                  ),
                  const SizedBox(height: 24),
                  Row(
                    children: [
                      Expanded(
                        child: _buildWalletAction(
                          Icons.add_rounded,
                          'Add Money',
                          () => context.push('/wallet/deposit'),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildWalletAction(
                          Icons.send_rounded,
                          'Send',
                          () => context.push('/wallet/transfer'),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: _buildWalletAction(
                          Icons.account_balance_outlined,
                          'Withdraw',
                          () => context.push('/wallet/withdraw'),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),

            // Quick actions
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                children: [
                  Expanded(
                    child: _buildQuickActionCard(
                      Icons.receipt_long_outlined,
                      'Pay Bills',
                      AppColors.warning,
                      () {},
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildQuickActionCard(
                      Icons.phone_android_outlined,
                      'Buy Airtime',
                      AppColors.success,
                      () {},
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildQuickActionCard(
                      Icons.account_balance_wallet_outlined,
                      'Bank Transfer',
                      AppColors.info,
                      () {},
                    ),
                  ),
                ],
              ),
            ),

            const SizedBox(height: 24),

            // Transaction history
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Transaction History',
                    style: AppTypography.titleMedium.copyWith(
                      color: AppColors.textPrimary,
                    ),
                  ),
                  TextButton(
                    onPressed: () {},
                    child: const Text('See all'),
                  ),
                ],
              ),
            ),

            // Filter chips
            SizedBox(
              height: 40,
              child: ListView(
                scrollDirection: Axis.horizontal,
                padding: const EdgeInsets.symmetric(horizontal: 16),
                children: ['All', 'Credits', 'Debits', 'Pending']
                    .map((filter) => Padding(
                          padding: const EdgeInsets.only(right: 8),
                          child: FilterChip(
                            label: Text(filter),
                            selected: _selectedFilter == filter,
                            onSelected: (selected) {
                              setState(() => _selectedFilter = filter);
                            },
                            selectedColor: AppColors.primary.withOpacity(0.2),
                            backgroundColor: AppColors.surfaceVariant,
                            labelStyle: AppTypography.labelMedium.copyWith(
                              color: _selectedFilter == filter
                                  ? AppColors.primary
                                  : AppColors.textSecondary,
                            ),
                          ),
                        ))
                    .toList(),
              ),
            ),

            const SizedBox(height: 12),

            // Transactions list
            ListView.builder(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: 10,
              itemBuilder: (context, index) {
                return _buildTransactionItem(
                  title: [
                    'Deposit from Bank',
                    'Transfer to @johndoe',
                    'Gig Payment - Mobile App',
                    'Savings Contribution',
                    'Withdrawal to GTBank',
                    'Airtime Purchase',
                    'Refund - Cancelled Gig',
                    'Loan Repayment',
                    'Escrow Release',
                    'Transfer from @mary',
                  ][index % 10],
                  amount: [50000, -15000, 75000, -25000, -30000, -500, 10000, -5000, 45000, 20000][index % 10].toDouble(),
                  type: [
                    'deposit',
                    'transfer',
                    'gig',
                    'savings',
                    'withdrawal',
                    'bill',
                    'refund',
                    'loan',
                    'escrow',
                    'transfer'
                  ][index % 10],
                  time: '${index + 1} hour${index > 0 ? 's' : ''} ago',
                  status: index == 4 ? 'pending' : 'completed',
                );
              },
            ),

            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Widget _buildBalanceItem(String label, double amount) {
    return Column(
      children: [
        Text(
          label,
          style: AppTypography.labelSmall.copyWith(
            color: Colors.white.withOpacity(0.7),
          ),
        ),
        const SizedBox(height: 4),
        Text(
          _isBalanceVisible
              ? _currencyFormat.format(amount)
              : '₦ ••••',
          style: AppTypography.titleSmall.copyWith(
            color: Colors.white,
          ),
        ),
      ],
    );
  }

  Widget _buildWalletAction(IconData icon, String label, VoidCallback onTap) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.2),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Icon(icon, color: Colors.white, size: 24),
            const SizedBox(height: 4),
            Text(
              label,
              style: AppTypography.labelSmall.copyWith(
                color: Colors.white,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildQuickActionCard(
    IconData icon,
    String label,
    Color color,
    VoidCallback onTap,
  ) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: color.withOpacity(0.1),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Icon(icon, color: color, size: 24),
            const SizedBox(height: 4),
            Text(
              label,
              style: AppTypography.labelSmall.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTransactionItem({
    required String title,
    required double amount,
    required String type,
    required String time,
    required String status,
  }) {
    final isCredit = amount > 0;
    IconData icon;
    Color iconColor;

    switch (type) {
      case 'deposit':
        icon = Icons.add_rounded;
        iconColor = AppColors.success;
        break;
      case 'transfer':
        icon = Icons.swap_horiz_rounded;
        iconColor = isCredit ? AppColors.success : AppColors.info;
        break;
      case 'gig':
        icon = Icons.work_rounded;
        iconColor = AppColors.primary;
        break;
      case 'savings':
        icon = Icons.savings_rounded;
        iconColor = AppColors.secondary;
        break;
      case 'withdrawal':
        icon = Icons.account_balance_rounded;
        iconColor = AppColors.warning;
        break;
      case 'bill':
        icon = Icons.receipt_rounded;
        iconColor = AppColors.error;
        break;
      case 'refund':
        icon = Icons.replay_rounded;
        iconColor = AppColors.info;
        break;
      case 'loan':
        icon = Icons.credit_score_rounded;
        iconColor = AppColors.warning;
        break;
      case 'escrow':
        icon = Icons.lock_open_rounded;
        iconColor = AppColors.success;
        break;
      default:
        icon = Icons.monetization_on_rounded;
        iconColor = AppColors.textSecondary;
    }

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: iconColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(icon, color: iconColor, size: 22),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: AppTypography.bodyMedium.copyWith(
                    color: AppColors.textPrimary,
                  ),
                ),
                Row(
                  children: [
                    Text(
                      time,
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.textTertiary,
                      ),
                    ),
                    if (status == 'pending') ...[
                      const SizedBox(width: 8),
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 6,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: AppColors.warningLight,
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          'Pending',
                          style: AppTypography.labelSmall.copyWith(
                            color: AppColors.warning,
                            fontSize: 10,
                          ),
                        ),
                      ),
                    ],
                  ],
                ),
              ],
            ),
          ),
          Text(
            '${isCredit ? '+' : ''}${_currencyFormat.format(amount)}',
            style: AppTypography.titleSmall.copyWith(
              color: isCredit ? AppColors.success : AppColors.error,
            ),
          ),
        ],
      ),
    );
  }
}

// Stub screens
class DepositScreen extends StatelessWidget {
  const DepositScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Add Money')),
      body: const Center(child: Text('Deposit Screen')),
    );
  }
}

class WithdrawScreen extends StatelessWidget {
  const WithdrawScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Withdraw')),
      body: const Center(child: Text('Withdraw Screen')),
    );
  }
}

class TransferScreen extends StatelessWidget {
  const TransferScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Send Money')),
      body: const Center(child: Text('Transfer Screen')),
    );
  }
}
