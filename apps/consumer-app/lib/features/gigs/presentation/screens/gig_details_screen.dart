import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:timeago/timeago.dart' as timeago;

import '../../../../core/constants/app_constants.dart';
import '../../../../core/widgets/widgets.dart';
import '../../data/models/gig_models.dart';
import '../providers/gigs_provider.dart';

class GigDetailsScreen extends ConsumerStatefulWidget {
  final String gigId;

  const GigDetailsScreen({super.key, required this.gigId});

  @override
  ConsumerState<GigDetailsScreen> createState() => _GigDetailsScreenState();
}

class _GigDetailsScreenState extends ConsumerState<GigDetailsScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    // Fetch gig details
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(gigsProvider.notifier).fetchGigDetails(widget.gigId);
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final gigsState = ref.watch(gigsProvider);
    final gig = gigsState.selectedGig;

    if (gigsState.isLoading && gig == null) {
      return const Scaffold(
        body: LoadingScreen(message: 'Loading gig details...'),
      );
    }

    if (gig == null) {
      return Scaffold(
        appBar: AppBar(),
        body: const EmptyStateCard(
          icon: Icons.error_outline,
          title: 'Gig Not Found',
          subtitle: 'This gig may have been removed or is no longer available.',
        ),
      );
    }

    return Scaffold(
      body: CustomScrollView(
        slivers: [
          // App Bar
          SliverAppBar(
            expandedHeight: 120,
            pinned: true,
            flexibleSpace: FlexibleSpaceBar(
              title: Text(
                gig.title,
                style: const TextStyle(fontSize: 16),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
              background: Container(
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [
                      AppColors.primary,
                      AppColors.primary.withOpacity(0.8),
                    ],
                  ),
                ),
              ),
            ),
            actions: [
              if (gig.isOwner)
                PopupMenuButton<String>(
                  icon: const Icon(Icons.more_vert, color: Colors.white),
                  onSelected: (value) => _handleMenuAction(value, gig),
                  itemBuilder: (context) => [
                    const PopupMenuItem(
                      value: 'edit',
                      child: Row(
                        children: [
                          Icon(Icons.edit, size: 20),
                          SizedBox(width: 8),
                          Text('Edit Gig'),
                        ],
                      ),
                    ),
                    if (gig.status == GigStatus.open)
                      const PopupMenuItem(
                        value: 'close',
                        child: Row(
                          children: [
                            Icon(Icons.close, size: 20),
                            SizedBox(width: 8),
                            Text('Close Gig'),
                          ],
                        ),
                      ),
                    const PopupMenuItem(
                      value: 'delete',
                      child: Row(
                        children: [
                          Icon(Icons.delete, size: 20, color: AppColors.error),
                          SizedBox(width: 8),
                          Text('Delete', style: TextStyle(color: AppColors.error)),
                        ],
                      ),
                    ),
                  ],
                ),
            ],
          ),

          // Content
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(AppSpacing.md),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Status and Category
                  Row(
                    children: [
                      StatusBadge(
                        label: gig.status.displayName,
                        color: _getStatusColor(gig.status),
                      ),
                      const SizedBox(width: 8),
                      if (gig.category != null)
                        Chip(
                          label: Text(
                            gig.category!.name,
                            style: const TextStyle(fontSize: 12),
                          ),
                          backgroundColor: AppColors.surface,
                          side: const BorderSide(color: AppColors.border),
                          padding: EdgeInsets.zero,
                          materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                        ),
                      const Spacer(),
                      if (gig.isRemote)
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 8,
                            vertical: 4,
                          ),
                          decoration: BoxDecoration(
                            color: AppColors.info.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: const Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(Icons.wifi, size: 14, color: AppColors.info),
                              SizedBox(width: 4),
                              Text(
                                'Remote',
                                style: TextStyle(
                                  fontSize: 12,
                                  color: AppColors.info,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ],
                          ),
                        ),
                    ],
                  ),

                  const SizedBox(height: 16),

                  // Budget
                  _buildInfoSection(
                    icon: Icons.payments_outlined,
                    title: 'Budget',
                    content: gig.budgetMin == gig.budgetMax
                        ? '₦${_formatAmount(gig.budgetMin)}'
                        : '₦${_formatAmount(gig.budgetMin)} - ₦${_formatAmount(gig.budgetMax)}',
                    highlight: true,
                  ),

                  const SizedBox(height: 12),

                  // Duration
                  _buildInfoSection(
                    icon: Icons.schedule,
                    title: 'Duration',
                    content: '${gig.durationDays} day${gig.durationDays > 1 ? 's' : ''}',
                  ),

                  const SizedBox(height: 12),

                  // Posted
                  _buildInfoSection(
                    icon: Icons.access_time,
                    title: 'Posted',
                    content: timeago.format(gig.createdAt),
                  ),

                  const SizedBox(height: 12),

                  // Proposals
                  _buildInfoSection(
                    icon: Icons.people_outline,
                    title: 'Proposals',
                    content: '${gig.proposalCount} received',
                  ),

                  const Divider(height: 32),

                  // Description
                  const Text(
                    'Description',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    gig.description,
                    style: const TextStyle(
                      fontSize: 14,
                      height: 1.6,
                      color: AppColors.textSecondary,
                    ),
                  ),

                  // Skills
                  if (gig.skills.isNotEmpty) ...[
                    const SizedBox(height: 24),
                    const Text(
                      'Required Skills',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Wrap(
                      spacing: 8,
                      runSpacing: 8,
                      children: gig.skills.map((skill) {
                        return Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 12,
                            vertical: 6,
                          ),
                          decoration: BoxDecoration(
                            color: AppColors.primary.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(20),
                          ),
                          child: Text(
                            skill,
                            style: const TextStyle(
                              fontSize: 13,
                              color: AppColors.primary,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                  ],

                  // Client Info (if not owner)
                  if (!gig.isOwner && gig.client != null) ...[
                    const Divider(height: 32),
                    const Text(
                      'Posted By',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 12),
                    _buildClientCard(gig.client!),
                  ],
                ],
              ),
            ),
          ),

          // Proposals Tab (for gig owner)
          if (gig.isOwner) ...[
            SliverPersistentHeader(
              pinned: true,
              delegate: _SliverTabBarDelegate(
                TabBar(
                  controller: _tabController,
                  labelColor: AppColors.primary,
                  unselectedLabelColor: AppColors.textSecondary,
                  indicatorColor: AppColors.primary,
                  tabs: const [
                    Tab(text: 'Details'),
                    Tab(text: 'Proposals'),
                  ],
                ),
              ),
            ),
            SliverFillRemaining(
              child: TabBarView(
                controller: _tabController,
                children: [
                  // Details tab content (already shown above)
                  const SizedBox.shrink(),
                  // Proposals tab
                  _buildProposalsList(gig),
                ],
              ),
            ),
          ],
        ],
      ),
      bottomNavigationBar: !gig.isOwner && gig.status == GigStatus.open
          ? _buildSubmitProposalBar()
          : null,
    );
  }

  Widget _buildInfoSection({
    required IconData icon,
    required String title,
    required String content,
    bool highlight = false,
  }) {
    return Row(
      children: [
        Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: (highlight ? AppColors.primary : AppColors.textSecondary)
                .withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Icon(
            icon,
            size: 20,
            color: highlight ? AppColors.primary : AppColors.textSecondary,
          ),
        ),
        const SizedBox(width: 12),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: const TextStyle(
                fontSize: 12,
                color: AppColors.textSecondary,
              ),
            ),
            Text(
              content,
              style: TextStyle(
                fontSize: 14,
                fontWeight: highlight ? FontWeight.w600 : FontWeight.w500,
                color: highlight ? AppColors.primary : AppColors.textPrimary,
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildClientCard(GigUser client) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          CircleAvatar(
            radius: 24,
            backgroundColor: AppColors.primary.withOpacity(0.1),
            backgroundImage:
                client.avatarUrl != null ? NetworkImage(client.avatarUrl!) : null,
            child: client.avatarUrl == null
                ? Text(
                    client.initials,
                    style: const TextStyle(
                      color: AppColors.primary,
                      fontWeight: FontWeight.w600,
                    ),
                  )
                : null,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  client.fullName,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 2),
                Row(
                  children: [
                    const Icon(Icons.star, size: 14, color: AppColors.warning),
                    const SizedBox(width: 4),
                    Text(
                      client.rating?.toStringAsFixed(1) ?? 'New',
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.textSecondary,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      '${client.completedGigs ?? 0} gigs completed',
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.message_outlined),
            onPressed: () {
              // Open chat with client
            },
          ),
        ],
      ),
    );
  }

  Widget _buildProposalsList(Gig gig) {
    final proposals = gig.proposals ?? [];

    if (proposals.isEmpty) {
      return const EmptyStateCard(
        icon: Icons.inbox_outlined,
        title: 'No Proposals Yet',
        subtitle: 'Proposals from freelancers will appear here.',
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(AppSpacing.md),
      itemCount: proposals.length,
      itemBuilder: (context, index) {
        final proposal = proposals[index];
        return _buildProposalCard(proposal);
      },
    );
  }

  Widget _buildProposalCard(Proposal proposal) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
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
            children: [
              CircleAvatar(
                radius: 20,
                backgroundColor: AppColors.primary.withOpacity(0.1),
                child: Text(
                  proposal.freelancer?.initials ?? 'U',
                  style: const TextStyle(
                    color: AppColors.primary,
                    fontWeight: FontWeight.w600,
                    fontSize: 14,
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      proposal.freelancer?.fullName ?? 'Unknown',
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                    Row(
                      children: [
                        const Icon(Icons.star, size: 12, color: AppColors.warning),
                        const SizedBox(width: 2),
                        Text(
                          proposal.freelancer?.rating?.toStringAsFixed(1) ?? 'New',
                          style: const TextStyle(
                            fontSize: 11,
                            color: AppColors.textSecondary,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              StatusBadge(
                label: proposal.status.displayName,
                color: _getProposalStatusColor(proposal.status),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Text(
            proposal.coverLetter,
            maxLines: 3,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(
              fontSize: 13,
              color: AppColors.textSecondary,
              height: 1.4,
            ),
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              _buildProposalInfo(
                Icons.payments_outlined,
                '₦${_formatAmount(proposal.proposedAmount)}',
              ),
              const SizedBox(width: 16),
              _buildProposalInfo(
                Icons.schedule,
                '${proposal.deliveryDays} days',
              ),
              const Spacer(),
              Text(
                timeago.format(proposal.createdAt),
                style: const TextStyle(
                  fontSize: 11,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
          if (proposal.status == ProposalStatus.pending) ...[
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: () => _rejectProposal(proposal),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppColors.error,
                      side: const BorderSide(color: AppColors.error),
                    ),
                    child: const Text('Decline'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton(
                    onPressed: () => _acceptProposal(proposal),
                    child: const Text('Accept'),
                  ),
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildProposalInfo(IconData icon, String text) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: AppColors.textSecondary),
        const SizedBox(width: 4),
        Text(
          text,
          style: const TextStyle(
            fontSize: 13,
            fontWeight: FontWeight.w500,
          ),
        ),
      ],
    );
  }

  Widget _buildSubmitProposalBar() {
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
        child: PrimaryButton(
          text: 'Submit Proposal',
          onPressed: () {
            context.push('/gigs/${widget.gigId}/submit-proposal');
          },
        ),
      ),
    );
  }

  Color _getStatusColor(GigStatus status) {
    switch (status) {
      case GigStatus.draft:
        return AppColors.textSecondary;
      case GigStatus.open:
        return AppColors.success;
      case GigStatus.inProgress:
        return AppColors.info;
      case GigStatus.completed:
        return AppColors.primary;
      case GigStatus.cancelled:
        return AppColors.error;
      case GigStatus.disputed:
        return AppColors.warning;
    }
  }

  Color _getProposalStatusColor(ProposalStatus status) {
    switch (status) {
      case ProposalStatus.pending:
        return AppColors.warning;
      case ProposalStatus.accepted:
        return AppColors.success;
      case ProposalStatus.rejected:
        return AppColors.error;
      case ProposalStatus.withdrawn:
        return AppColors.textSecondary;
    }
  }

  String _formatAmount(double amount) {
    return amount.toStringAsFixed(0).replaceAllMapped(
          RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},',
        );
  }

  void _handleMenuAction(String action, Gig gig) {
    switch (action) {
      case 'edit':
        context.push('/gigs/${gig.id}/edit');
        break;
      case 'close':
        _showCloseConfirmation(gig);
        break;
      case 'delete':
        _showDeleteConfirmation(gig);
        break;
    }
  }

  void _showCloseConfirmation(Gig gig) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Close Gig?'),
        content: const Text(
          'This will stop accepting new proposals. Existing proposals can still be reviewed.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              ref.read(gigsProvider.notifier).closeGig(gig.id);
            },
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  void _showDeleteConfirmation(Gig gig) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Gig?'),
        content: const Text(
          'This action cannot be undone. All proposals will also be deleted.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              ref.read(gigsProvider.notifier).deleteGig(gig.id);
              context.pop();
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.error,
            ),
            child: const Text('Delete'),
          ),
        ],
      ),
    );
  }

  void _acceptProposal(Proposal proposal) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Accept Proposal?'),
        content: Text(
          'You are about to accept ${proposal.freelancer?.fullName}\'s proposal for ₦${_formatAmount(proposal.proposedAmount)}. '
          'The funds will be held in escrow until the gig is completed.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              ref.read(gigsProvider.notifier).acceptProposal(
                    widget.gigId,
                    proposal.id,
                  );
            },
            child: const Text('Accept'),
          ),
        ],
      ),
    );
  }

  void _rejectProposal(Proposal proposal) {
    ref.read(gigsProvider.notifier).rejectProposal(widget.gigId, proposal.id);
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
  Widget build(
    BuildContext context,
    double shrinkOffset,
    bool overlapsContent,
  ) {
    return Container(
      color: AppColors.background,
      child: tabBar,
    );
  }

  @override
  bool shouldRebuild(_SliverTabBarDelegate oldDelegate) {
    return false;
  }
}
