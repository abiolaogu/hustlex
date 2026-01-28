import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/utils/currency_utils.dart';
import '../../../../core/utils/date_utils.dart';
import '../../data/models/gig_models.dart';

// Tab state provider
final myGigsTabProvider = StateProvider<int>((ref) => 0);

// Mock data providers
final myPostedGigsProvider = FutureProvider<List<Gig>>((ref) async {
  await Future.delayed(const Duration(milliseconds: 500));
  return [
    Gig(
      id: 'gig-1',
      clientId: 'user-1',
      title: 'Mobile App UI Design',
      description: 'Need a UI designer for my mobile app',
      category: GigCategory.design,
      budget: 150000,
      budgetType: BudgetType.fixed,
      location: 'Lagos, Nigeria',
      locationType: LocationType.remote,
      status: GigStatus.open,
      skills: ['Figma', 'UI/UX', 'Mobile Design'],
      deadline: DateTime.now().add(const Duration(days: 14)),
      proposalCount: 8,
      viewCount: 45,
      createdAt: DateTime.now().subtract(const Duration(days: 2)),
      updatedAt: DateTime.now(),
    ),
    Gig(
      id: 'gig-2',
      clientId: 'user-1',
      title: 'Backend API Development',
      description: 'Node.js backend for e-commerce platform',
      category: GigCategory.development,
      budget: 300000,
      budgetType: BudgetType.fixed,
      location: 'Remote',
      locationType: LocationType.remote,
      status: GigStatus.inProgress,
      skills: ['Node.js', 'PostgreSQL', 'REST API'],
      deadline: DateTime.now().add(const Duration(days: 30)),
      proposalCount: 15,
      viewCount: 120,
      assignedFreelancerId: 'freelancer-1',
      createdAt: DateTime.now().subtract(const Duration(days: 10)),
      updatedAt: DateTime.now(),
    ),
    Gig(
      id: 'gig-3',
      clientId: 'user-1',
      title: 'Social Media Marketing',
      description: 'Monthly social media management',
      category: GigCategory.marketing,
      budget: 80000,
      budgetType: BudgetType.fixed,
      location: 'Lagos',
      locationType: LocationType.hybrid,
      status: GigStatus.completed,
      skills: ['Social Media', 'Content Creation', 'Analytics'],
      deadline: DateTime.now().subtract(const Duration(days: 5)),
      proposalCount: 22,
      viewCount: 89,
      assignedFreelancerId: 'freelancer-2',
      createdAt: DateTime.now().subtract(const Duration(days: 45)),
      updatedAt: DateTime.now(),
    ),
  ];
});

final myProposalsProvider = FutureProvider<List<GigProposal>>((ref) async {
  await Future.delayed(const Duration(milliseconds: 500));
  return [
    GigProposal(
      id: 'proposal-1',
      gigId: 'gig-10',
      freelancerId: 'user-1',
      coverLetter: 'I am excited to work on this project...',
      proposedAmount: 120000,
      estimatedDuration: '2 weeks',
      status: ProposalStatus.pending,
      createdAt: DateTime.now().subtract(const Duration(hours: 5)),
      updatedAt: DateTime.now(),
    ),
    GigProposal(
      id: 'proposal-2',
      gigId: 'gig-11',
      freelancerId: 'user-1',
      coverLetter: 'With my experience in mobile development...',
      proposedAmount: 250000,
      estimatedDuration: '4 weeks',
      status: ProposalStatus.accepted,
      createdAt: DateTime.now().subtract(const Duration(days: 3)),
      updatedAt: DateTime.now(),
    ),
    GigProposal(
      id: 'proposal-3',
      gigId: 'gig-12',
      freelancerId: 'user-1',
      coverLetter: 'I specialize in this type of work...',
      proposedAmount: 75000,
      estimatedDuration: '1 week',
      status: ProposalStatus.rejected,
      createdAt: DateTime.now().subtract(const Duration(days: 7)),
      updatedAt: DateTime.now(),
    ),
  ];
});

class MyGigsScreen extends ConsumerWidget {
  const MyGigsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final currentTab = ref.watch(myGigsTabProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('My Gigs'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => context.push('/gigs/create'),
          ),
        ],
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(48),
          child: _buildTabBar(ref, currentTab),
        ),
      ),
      body: currentTab == 0
          ? _PostedGigsList()
          : _ProposalsList(),
    );
  }

  Widget _buildTabBar(WidgetRef ref, int currentTab) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: AppColors.background,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Expanded(
            child: _TabButton(
              label: 'Posted Gigs',
              isSelected: currentTab == 0,
              onTap: () => ref.read(myGigsTabProvider.notifier).state = 0,
            ),
          ),
          Expanded(
            child: _TabButton(
              label: 'My Proposals',
              isSelected: currentTab == 1,
              onTap: () => ref.read(myGigsTabProvider.notifier).state = 1,
            ),
          ),
        ],
      ),
    );
  }
}

class _TabButton extends StatelessWidget {
  final String label;
  final bool isSelected;
  final VoidCallback onTap;

  const _TabButton({
    required this.label,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 10),
        decoration: BoxDecoration(
          color: isSelected ? AppColors.primary : Colors.transparent,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text(
          label,
          textAlign: TextAlign.center,
          style: AppTypography.labelLarge.copyWith(
            color: isSelected ? Colors.white : AppColors.textSecondary,
            fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
          ),
        ),
      ),
    );
  }
}

class _PostedGigsList extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final gigsAsync = ref.watch(myPostedGigsProvider);

    return gigsAsync.when(
      data: (gigs) {
        if (gigs.isEmpty) {
          return _buildEmptyState(
            context,
            icon: Icons.work_outline,
            title: 'No Gigs Posted',
            subtitle: 'Start by posting your first gig',
            actionLabel: 'Post a Gig',
            onAction: () => context.push('/gigs/create'),
          );
        }

        return ListView.separated(
          padding: const EdgeInsets.all(16),
          itemCount: gigs.length,
          separatorBuilder: (_, __) => const SizedBox(height: 12),
          itemBuilder: (context, index) {
            return _PostedGigCard(gig: gigs[index]);
          },
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) => Center(child: Text('Error: $e')),
    );
  }

  Widget _buildEmptyState(
    BuildContext context, {
    required IconData icon,
    required String title,
    required String subtitle,
    required String actionLabel,
    required VoidCallback onAction,
  }) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              size: 80,
              color: AppColors.textSecondary.withOpacity(0.5),
            ),
            const SizedBox(height: 16),
            Text(
              title,
              style: AppTypography.titleLarge.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              subtitle,
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: onAction,
              child: Text(actionLabel),
            ),
          ],
        ),
      ),
    );
  }
}

class _PostedGigCard extends StatelessWidget {
  final Gig gig;

  const _PostedGigCard({required this.gig});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () => context.push('/gigs/${gig.id}'),
      child: Container(
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: AppColors.border),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          gig.title,
                          style: AppTypography.titleMedium.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                      _buildStatusBadge(gig.status),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text(
                    gig.description,
                    style: AppTypography.bodyMedium.copyWith(
                      color: AppColors.textSecondary,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      _buildMetric(
                        Icons.description,
                        '${gig.proposalCount} proposals',
                      ),
                      const SizedBox(width: 16),
                      _buildMetric(
                        Icons.visibility,
                        '${gig.viewCount} views',
                      ),
                    ],
                  ),
                ],
              ),
            ),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              decoration: BoxDecoration(
                color: AppColors.background,
                borderRadius: const BorderRadius.vertical(
                  bottom: Radius.circular(12),
                ),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    CurrencyUtils.formatNaira(gig.budget),
                    style: AppTypography.titleSmall.copyWith(
                      color: AppColors.primary,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    'Posted ${AppDateUtils.formatRelative(gig.createdAt)}',
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

  Widget _buildStatusBadge(GigStatus status) {
    Color color;
    String label;

    switch (status) {
      case GigStatus.draft:
        color = AppColors.textSecondary;
        label = 'Draft';
        break;
      case GigStatus.open:
        color = AppColors.success;
        label = 'Open';
        break;
      case GigStatus.inProgress:
        color = AppColors.info;
        label = 'In Progress';
        break;
      case GigStatus.completed:
        color = AppColors.primary;
        label = 'Completed';
        break;
      case GigStatus.cancelled:
        color = AppColors.error;
        label = 'Cancelled';
        break;
      case GigStatus.disputed:
        color = AppColors.warning;
        label = 'Disputed';
        break;
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

  Widget _buildMetric(IconData icon, String label) {
    return Row(
      children: [
        Icon(icon, size: 16, color: AppColors.textSecondary),
        const SizedBox(width: 4),
        Text(
          label,
          style: AppTypography.bodySmall.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }
}

class _ProposalsList extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final proposalsAsync = ref.watch(myProposalsProvider);

    return proposalsAsync.when(
      data: (proposals) {
        if (proposals.isEmpty) {
          return _buildEmptyState(
            context,
            icon: Icons.send_outlined,
            title: 'No Proposals Yet',
            subtitle: 'Browse gigs and submit proposals to get started',
            actionLabel: 'Browse Gigs',
            onAction: () => context.go('/gigs'),
          );
        }

        return ListView.separated(
          padding: const EdgeInsets.all(16),
          itemCount: proposals.length,
          separatorBuilder: (_, __) => const SizedBox(height: 12),
          itemBuilder: (context, index) {
            return _ProposalCard(proposal: proposals[index]);
          },
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) => Center(child: Text('Error: $e')),
    );
  }

  Widget _buildEmptyState(
    BuildContext context, {
    required IconData icon,
    required String title,
    required String subtitle,
    required String actionLabel,
    required VoidCallback onAction,
  }) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              size: 80,
              color: AppColors.textSecondary.withOpacity(0.5),
            ),
            const SizedBox(height: 16),
            Text(
              title,
              style: AppTypography.titleLarge.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              subtitle,
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: onAction,
              child: Text(actionLabel),
            ),
          ],
        ),
      ),
    );
  }
}

class _ProposalCard extends StatelessWidget {
  final GigProposal proposal;

  const _ProposalCard({required this.proposal});

  @override
  Widget build(BuildContext context) {
    // Mock gig data for display
    final mockGigTitles = {
      'gig-10': 'Flutter App Development',
      'gig-11': 'E-commerce Website',
      'gig-12': 'Logo Design Project',
    };

    return GestureDetector(
      onTap: () => context.push('/gigs/${proposal.gigId}'),
      child: Container(
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: AppColors.border),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          mockGigTitles[proposal.gigId] ?? 'Gig Title',
                          style: AppTypography.titleMedium.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                      _buildStatusBadge(proposal.status),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text(
                    proposal.coverLetter,
                    style: AppTypography.bodyMedium.copyWith(
                      color: AppColors.textSecondary,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Icon(
                        Icons.payments,
                        size: 16,
                        color: AppColors.textSecondary,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        'Bid: ${CurrencyUtils.formatNaira(proposal.proposedAmount)}',
                        style: AppTypography.bodySmall.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                      const SizedBox(width: 16),
                      Icon(
                        Icons.schedule,
                        size: 16,
                        color: AppColors.textSecondary,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        proposal.estimatedDuration,
                        style: AppTypography.bodySmall.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              decoration: BoxDecoration(
                color: AppColors.background,
                borderRadius: const BorderRadius.vertical(
                  bottom: Radius.circular(12),
                ),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Submitted ${AppDateUtils.formatRelative(proposal.createdAt)}',
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                  if (proposal.status == ProposalStatus.pending)
                    TextButton(
                      onPressed: () => _showWithdrawDialog(context),
                      style: TextButton.styleFrom(
                        padding: EdgeInsets.zero,
                        minimumSize: const Size(0, 0),
                        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                      ),
                      child: Text(
                        'Withdraw',
                        style: AppTypography.labelMedium.copyWith(
                          color: AppColors.error,
                        ),
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

  Widget _buildStatusBadge(ProposalStatus status) {
    Color color;
    String label;
    IconData icon;

    switch (status) {
      case ProposalStatus.pending:
        color = AppColors.warning;
        label = 'Pending';
        icon = Icons.hourglass_empty;
        break;
      case ProposalStatus.shortlisted:
        color = AppColors.info;
        label = 'Shortlisted';
        icon = Icons.star;
        break;
      case ProposalStatus.accepted:
        color = AppColors.success;
        label = 'Accepted';
        icon = Icons.check_circle;
        break;
      case ProposalStatus.rejected:
        color = AppColors.error;
        label = 'Rejected';
        icon = Icons.cancel;
        break;
      case ProposalStatus.withdrawn:
        color = AppColors.textSecondary;
        label = 'Withdrawn';
        icon = Icons.undo;
        break;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 12, color: color),
          const SizedBox(width: 4),
          Text(
            label,
            style: AppTypography.labelSmall.copyWith(
              color: color,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }

  void _showWithdrawDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Withdraw Proposal?'),
        content: const Text(
          'Are you sure you want to withdraw this proposal? This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Proposal withdrawn'),
                ),
              );
            },
            style: TextButton.styleFrom(
              foregroundColor: AppColors.error,
            ),
            child: const Text('Withdraw'),
          ),
        ],
      ),
    );
  }
}
