import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/widgets/widgets.dart';
import '../../data/models/savings_models.dart';
import '../providers/savings_provider.dart';

class CircleDetailsScreen extends ConsumerStatefulWidget {
  final String circleId;

  const CircleDetailsScreen({super.key, required this.circleId});

  @override
  ConsumerState<CircleDetailsScreen> createState() => _CircleDetailsScreenState();
}

class _CircleDetailsScreenState extends ConsumerState<CircleDetailsScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(savingsProvider.notifier).fetchCircleDetails(widget.circleId);
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final savingsState = ref.watch(savingsProvider);
    final circle = savingsState.selectedCircle;

    if (savingsState.isLoading && circle == null) {
      return const Scaffold(
        body: LoadingScreen(message: 'Loading circle details...'),
      );
    }

    if (circle == null) {
      return Scaffold(
        appBar: AppBar(),
        body: const EmptyStateCard(
          icon: Icons.error_outline,
          title: 'Circle Not Found',
          subtitle: 'This savings circle may have been removed.',
        ),
      );
    }

    return Scaffold(
      body: NestedScrollView(
        headerSliverBuilder: (context, innerBoxIsScrolled) => [
          SliverAppBar(
            expandedHeight: 200,
            pinned: true,
            flexibleSpace: FlexibleSpaceBar(
              background: _buildHeader(circle),
            ),
            actions: [
              if (circle.isAdmin)
                PopupMenuButton<String>(
                  icon: const Icon(Icons.more_vert, color: Colors.white),
                  onSelected: (value) => _handleMenuAction(value, circle),
                  itemBuilder: (context) => [
                    const PopupMenuItem(
                      value: 'edit',
                      child: Row(
                        children: [
                          Icon(Icons.edit, size: 20),
                          SizedBox(width: 8),
                          Text('Edit Circle'),
                        ],
                      ),
                    ),
                    const PopupMenuItem(
                      value: 'invite',
                      child: Row(
                        children: [
                          Icon(Icons.person_add, size: 20),
                          SizedBox(width: 8),
                          Text('Invite Members'),
                        ],
                      ),
                    ),
                  ],
                ),
            ],
          ),
          SliverPersistentHeader(
            pinned: true,
            delegate: _SliverTabBarDelegate(
              TabBar(
                controller: _tabController,
                labelColor: AppColors.primary,
                unselectedLabelColor: AppColors.textSecondary,
                indicatorColor: AppColors.primary,
                tabs: const [
                  Tab(text: 'Overview'),
                  Tab(text: 'Members'),
                  Tab(text: 'Payments'),
                ],
              ),
            ),
          ),
        ],
        body: TabBarView(
          controller: _tabController,
          children: [
            _buildOverviewTab(circle),
            _buildMembersTab(circle),
            _buildPaymentsTab(circle),
          ],
        ),
      ),
      bottomNavigationBar: _buildBottomBar(circle),
    );
  }

  Widget _buildHeader(SavingsCircle circle) {
    return Container(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            AppColors.secondary,
            AppColors.secondary.withOpacity(0.8),
          ],
        ),
      ),
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.md),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.end,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      circle.type == CircleType.rotational
                          ? Icons.refresh
                          : Icons.savings,
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
                          circle.name,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        Text(
                          circle.type.displayName,
                          style: TextStyle(
                            color: Colors.white.withOpacity(0.8),
                            fontSize: 14,
                          ),
                        ),
                      ],
                    ),
                  ),
                  StatusBadge(
                    label: circle.status.displayName,
                    color: _getStatusColor(circle.status),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              // Progress
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        '₦${_formatAmount(circle.currentAmount)}',
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Text(
                        'of ₦${_formatAmount(circle.targetAmount)}',
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.8),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  LinearProgressIndicator(
                    value: circle.progress,
                    backgroundColor: Colors.white.withOpacity(0.3),
                    valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
                    borderRadius: BorderRadius.circular(4),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildOverviewTab(SavingsCircle circle) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Key Stats
          Row(
            children: [
              Expanded(
                child: _buildStatCard(
                  icon: Icons.payments_outlined,
                  label: 'Contribution',
                  value: '₦${_formatAmount(circle.contributionAmount)}',
                  subtitle: circle.frequency.displayName,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildStatCard(
                  icon: Icons.people_outline,
                  label: 'Members',
                  value: '${circle.memberCount}/${circle.maxMembers}',
                  subtitle: 'Joined',
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(
                child: _buildStatCard(
                  icon: Icons.loop,
                  label: 'Current Round',
                  value: '${circle.currentRound}/${circle.totalRounds}',
                  subtitle: 'Progress',
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildStatCard(
                  icon: Icons.calendar_today,
                  label: 'Next Due',
                  value: circle.nextDueDate != null
                      ? DateFormat('MMM d').format(circle.nextDueDate!)
                      : 'N/A',
                  subtitle: circle.nextDueDate != null
                      ? _getDaysUntil(circle.nextDueDate!)
                      : '',
                ),
              ),
            ],
          ),

          const SizedBox(height: 24),

          // How It Works
          const Text(
            'How It Works',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: AppColors.surface,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppColors.border),
            ),
            child: Column(
              children: [
                if (circle.type == CircleType.rotational) ...[
                  _buildInfoRow(
                    Icons.people,
                    'Each member contributes ${circle.frequency.displayName.toLowerCase()}',
                  ),
                  _buildInfoRow(
                    Icons.account_balance_wallet,
                    'One member receives the full pot each round',
                  ),
                  _buildInfoRow(
                    Icons.repeat,
                    'Continues until everyone has received',
                  ),
                ] else ...[
                  _buildInfoRow(
                    Icons.savings,
                    'Everyone saves together towards a goal',
                  ),
                  _buildInfoRow(
                    Icons.payments,
                    'Contributions pool until target is reached',
                  ),
                  _buildInfoRow(
                    Icons.celebration,
                    'Funds distributed equally at completion',
                  ),
                ],
              ],
            ),
          ),

          const SizedBox(height: 24),

          // Your Status
          const Text(
            'Your Status',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: AppColors.primary.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Your Position'),
                    Text(
                      '#${circle.myPosition ?? '-'}',
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                  ],
                ),
                const Divider(height: 24),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Total Contributed'),
                    Text(
                      '₦${_formatAmount(circle.myTotalContributed ?? 0)}',
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                  ],
                ),
                if (circle.type == CircleType.rotational) ...[
                  const Divider(height: 24),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Payout Status'),
                      StatusBadge(
                        label: circle.hasReceivedPayout ? 'Received' : 'Pending',
                        color: circle.hasReceivedPayout
                            ? AppColors.success
                            : AppColors.warning,
                      ),
                    ],
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMembersTab(SavingsCircle circle) {
    final members = circle.members ?? [];

    if (members.isEmpty) {
      return const EmptyStateCard(
        icon: Icons.people_outline,
        title: 'No Members Yet',
        subtitle: 'Invite friends to join your savings circle.',
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(AppSpacing.md),
      itemCount: members.length,
      itemBuilder: (context, index) {
        final member = members[index];
        return _buildMemberCard(member, circle);
      },
    );
  }

  Widget _buildMemberCard(CircleMember member, SavingsCircle circle) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Stack(
            children: [
              CircleAvatar(
                radius: 24,
                backgroundColor: AppColors.primary.withOpacity(0.1),
                child: Text(
                  member.user?.initials ?? '#${member.position}',
                  style: const TextStyle(
                    color: AppColors.primary,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
              if (circle.type == CircleType.rotational)
                Positioned(
                  bottom: 0,
                  right: 0,
                  child: Container(
                    padding: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                      color: member.hasReceivedPayout
                          ? AppColors.success
                          : AppColors.warning,
                      shape: BoxShape.circle,
                      border: Border.all(color: Colors.white, width: 2),
                    ),
                    child: Icon(
                      member.hasReceivedPayout ? Icons.check : Icons.schedule,
                      size: 10,
                      color: Colors.white,
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Text(
                      member.user?.fullName ?? 'Member #${member.position}',
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                    if (member.isAdmin) ...[
                      const SizedBox(width: 6),
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 6,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: AppColors.primary.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: const Text(
                          'Admin',
                          style: TextStyle(
                            fontSize: 10,
                            color: AppColors.primary,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ),
                    ],
                  ],
                ),
                const SizedBox(height: 2),
                Text(
                  'Contributed: ₦${_formatAmount(member.totalContributed)}',
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          if (circle.type == CircleType.rotational)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(
                color: AppColors.surface,
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: AppColors.border),
              ),
              child: Text(
                'Position #${member.position}',
                style: const TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildPaymentsTab(SavingsCircle circle) {
    final contributions = circle.contributions ?? [];
    final payouts = circle.payouts ?? [];

    return DefaultTabController(
      length: 2,
      child: Column(
        children: [
          Container(
            margin: const EdgeInsets.all(AppSpacing.md),
            decoration: BoxDecoration(
              color: AppColors.surface,
              borderRadius: BorderRadius.circular(8),
            ),
            child: TabBar(
              labelColor: AppColors.primary,
              unselectedLabelColor: AppColors.textSecondary,
              indicatorSize: TabBarIndicatorSize.tab,
              dividerColor: Colors.transparent,
              tabs: const [
                Tab(text: 'Contributions'),
                Tab(text: 'Payouts'),
              ],
            ),
          ),
          Expanded(
            child: TabBarView(
              children: [
                // Contributions
                contributions.isEmpty
                    ? const EmptyStateCard(
                        icon: Icons.payments_outlined,
                        title: 'No Contributions Yet',
                        subtitle: 'Your contribution history will appear here.',
                      )
                    : ListView.builder(
                        padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
                        itemCount: contributions.length,
                        itemBuilder: (context, index) {
                          final contribution = contributions[index];
                          return _buildContributionItem(contribution);
                        },
                      ),
                // Payouts
                payouts.isEmpty
                    ? const EmptyStateCard(
                        icon: Icons.account_balance_wallet,
                        title: 'No Payouts Yet',
                        subtitle: 'Payout history will appear here.',
                      )
                    : ListView.builder(
                        padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
                        itemCount: payouts.length,
                        itemBuilder: (context, index) {
                          final payout = payouts[index];
                          return _buildPayoutItem(payout);
                        },
                      ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildContributionItem(Contribution contribution) {
    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: contribution.isPaid
                  ? AppColors.success.withOpacity(0.1)
                  : contribution.isOverdue
                      ? AppColors.error.withOpacity(0.1)
                      : AppColors.warning.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Icon(
              contribution.isPaid
                  ? Icons.check_circle
                  : contribution.isOverdue
                      ? Icons.error
                      : Icons.schedule,
              color: contribution.isPaid
                  ? AppColors.success
                  : contribution.isOverdue
                      ? AppColors.error
                      : AppColors.warning,
              size: 20,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Round ${contribution.round}',
                  style: const TextStyle(fontWeight: FontWeight.w600),
                ),
                Text(
                  contribution.isPaid
                      ? 'Paid ${DateFormat('MMM d, yyyy').format(contribution.paidAt!)}'
                      : 'Due ${DateFormat('MMM d, yyyy').format(contribution.dueDate)}',
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          Text(
            '₦${_formatAmount(contribution.amount)}',
            style: const TextStyle(fontWeight: FontWeight.w600),
          ),
        ],
      ),
    );
  }

  Widget _buildPayoutItem(Payout payout) {
    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: AppColors.success.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: const Icon(
              Icons.account_balance_wallet,
              color: AppColors.success,
              size: 20,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Round ${payout.round} Payout',
                  style: const TextStyle(fontWeight: FontWeight.w600),
                ),
                Text(
                  DateFormat('MMM d, yyyy').format(payout.scheduledDate),
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          Text(
            '₦${_formatAmount(payout.amount)}',
            style: const TextStyle(
              fontWeight: FontWeight.w600,
              color: AppColors.success,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatCard({
    required IconData icon,
    required String label,
    required String value,
    required String subtitle,
  }) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, color: AppColors.textSecondary, size: 20),
          const SizedBox(height: 8),
          Text(
            value,
            style: const TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 2),
          Text(
            label,
            style: const TextStyle(
              fontSize: 12,
              color: AppColors.textSecondary,
            ),
          ),
          if (subtitle.isNotEmpty)
            Text(
              subtitle,
              style: const TextStyle(
                fontSize: 11,
                color: AppColors.textSecondary,
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildInfoRow(IconData icon, String text) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        children: [
          Icon(icon, size: 18, color: AppColors.secondary),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              text,
              style: const TextStyle(fontSize: 13),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomBar(SavingsCircle circle) {
    if (circle.status != CircleStatus.active) return const SizedBox.shrink();

    final hasUpcomingPayment = circle.nextDueDate != null &&
        circle.nextDueDate!.isBefore(DateTime.now().add(const Duration(days: 7)));

    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.background,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, -5),
          ),
        ],
      ),
      child: SafeArea(
        child: hasUpcomingPayment
            ? PrimaryButton(
                text: 'Make Contribution',
                onPressed: () => _makeContribution(circle),
              )
            : SecondaryButton(
                text: 'View Schedule',
                onPressed: () {
                  _tabController.animateTo(2);
                },
              ),
      ),
    );
  }

  void _makeContribution(SavingsCircle circle) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _ContributionSheet(circle: circle),
    );
  }

  void _handleMenuAction(String action, SavingsCircle circle) {
    switch (action) {
      case 'edit':
        // Navigate to edit screen
        break;
      case 'invite':
        _showInviteDialog(circle);
        break;
    }
  }

  void _showInviteDialog(SavingsCircle circle) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Invite Members'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('Share this code with friends to join your circle:'),
            const SizedBox(height: 16),
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: AppColors.surface,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    circle.inviteCode ?? 'INVITE123',
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                      letterSpacing: 4,
                    ),
                  ),
                  const SizedBox(width: 12),
                  IconButton(
                    icon: const Icon(Icons.copy),
                    onPressed: () {
                      // Copy to clipboard
                      Navigator.pop(context);
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Code copied!')),
                      );
                    },
                  ),
                ],
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
          ElevatedButton.icon(
            icon: const Icon(Icons.share),
            label: const Text('Share'),
            onPressed: () {
              // Share functionality
              Navigator.pop(context);
            },
          ),
        ],
      ),
    );
  }

  Color _getStatusColor(CircleStatus status) {
    switch (status) {
      case CircleStatus.pending:
        return AppColors.warning;
      case CircleStatus.active:
        return AppColors.success;
      case CircleStatus.completed:
        return AppColors.primary;
      case CircleStatus.cancelled:
        return AppColors.error;
    }
  }

  String _formatAmount(double amount) {
    return amount.toStringAsFixed(0).replaceAllMapped(
          RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},',
        );
  }

  String _getDaysUntil(DateTime date) {
    final days = date.difference(DateTime.now()).inDays;
    if (days == 0) return 'Today';
    if (days == 1) return 'Tomorrow';
    if (days < 0) return '${-days} days ago';
    return 'in $days days';
  }
}

// Sliver Tab Bar Delegate
class _SliverTabBarDelegate extends SliverPersistentHeaderDelegate {
  final TabBar tabBar;

  _SliverTabBarDelegate(this.tabBar);

  @override
  double get minExtent => tabBar.preferredSize.height;

  @override
  double get maxExtent => tabBar.preferredSize.height;

  @override
  Widget build(BuildContext context, double shrinkOffset, bool overlapsContent) {
    return Container(
      color: AppColors.background,
      child: tabBar,
    );
  }

  @override
  bool shouldRebuild(_SliverTabBarDelegate oldDelegate) => false;
}

// Contribution Bottom Sheet
class _ContributionSheet extends ConsumerStatefulWidget {
  final SavingsCircle circle;

  const _ContributionSheet({required this.circle});

  @override
  ConsumerState<_ContributionSheet> createState() => _ContributionSheetState();
}

class _ContributionSheetState extends ConsumerState<_ContributionSheet> {
  bool _isLoading = false;

  Future<void> _submitContribution() async {
    setState(() => _isLoading = true);

    try {
      await ref.read(savingsProvider.notifier).makeContribution(
            widget.circle.id,
            widget.circle.contributionAmount,
          );

      if (mounted) {
        Navigator.pop(context);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Contribution successful!'),
            backgroundColor: AppColors.success,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(e.toString()),
            backgroundColor: AppColors.error,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
      ),
      child: Container(
        padding: const EdgeInsets.all(AppSpacing.lg),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: AppColors.border,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const SizedBox(height: 24),
            const Text(
              'Make Contribution',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Contribution to ${widget.circle.name}',
              style: const TextStyle(color: AppColors.textSecondary),
            ),
            const SizedBox(height: 24),
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: AppColors.surface,
                borderRadius: BorderRadius.circular(16),
              ),
              child: Column(
                children: [
                  const Text(
                    'Amount',
                    style: TextStyle(color: AppColors.textSecondary),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    '₦${widget.circle.contributionAmount.toStringAsFixed(0).replaceAllMapped(
                          RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
                          (Match m) => '${m[1]},',
                        )}',
                    style: const TextStyle(
                      fontSize: 32,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),
            PrimaryButton(
              text: 'Pay from Wallet',
              onPressed: _isLoading ? null : _submitContribution,
              isLoading: _isLoading,
            ),
            const SizedBox(height: 12),
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancel'),
            ),
          ],
        ),
      ),
    );
  }
}
