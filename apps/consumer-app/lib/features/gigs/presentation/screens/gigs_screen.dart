import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '../../../../core/constants/app_constants.dart';

class GigsScreen extends ConsumerStatefulWidget {
  const GigsScreen({super.key});

  @override
  ConsumerState<GigsScreen> createState() => _GigsScreenState();
}

class _GigsScreenState extends ConsumerState<GigsScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final _searchController = TextEditingController();
  String _selectedCategory = 'All';
  
  final _currencyFormat = NumberFormat.currency(
    locale: 'en_NG',
    symbol: '₦',
    decimalDigits: 0,
  );

  final List<String> _categories = [
    'All',
    'Tech',
    'Design',
    'Writing',
    'Marketing',
    'Video',
    'Music',
    'Business',
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Gig Marketplace'),
        actions: [
          IconButton(
            onPressed: () {},
            icon: const Icon(Icons.filter_list_rounded),
          ),
        ],
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(110),
          child: Column(
            children: [
              // Search bar
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Search gigs...',
                    prefixIcon: const Icon(Icons.search, size: 20),
                    suffixIcon: IconButton(
                      onPressed: () => _searchController.clear(),
                      icon: const Icon(Icons.close, size: 20),
                    ),
                    filled: true,
                    fillColor: AppColors.surfaceVariant,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide.none,
                    ),
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 12,
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 12),
              // Tab bar
              TabBar(
                controller: _tabController,
                tabs: const [
                  Tab(text: 'Browse'),
                  Tab(text: 'My Gigs'),
                  Tab(text: 'Proposals'),
                ],
                labelColor: AppColors.primary,
                unselectedLabelColor: AppColors.textSecondary,
                indicatorColor: AppColors.primary,
                indicatorSize: TabBarIndicatorSize.label,
              ),
            ],
          ),
        ),
      ),
      body: Column(
        children: [
          // Category chips
          Container(
            height: 50,
            padding: const EdgeInsets.symmetric(vertical: 8),
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: _categories.length,
              itemBuilder: (context, index) {
                final category = _categories[index];
                final isSelected = _selectedCategory == category;
                return Padding(
                  padding: const EdgeInsets.only(right: 8),
                  child: FilterChip(
                    label: Text(category),
                    selected: isSelected,
                    onSelected: (selected) {
                      setState(() => _selectedCategory = category);
                    },
                    selectedColor: AppColors.primary.withOpacity(0.2),
                    backgroundColor: AppColors.surfaceVariant,
                    labelStyle: AppTypography.labelMedium.copyWith(
                      color: isSelected
                          ? AppColors.primary
                          : AppColors.textSecondary,
                    ),
                    checkmarkColor: AppColors.primary,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(20),
                    ),
                  ),
                );
              },
            ),
          ),
          // Content
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                _buildBrowseTab(),
                _buildMyGigsTab(),
                _buildProposalsTab(),
              ],
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => context.push('/gigs/create'),
        backgroundColor: AppColors.primary,
        icon: const Icon(Icons.add_rounded),
        label: const Text('Post a Gig'),
      ),
    );
  }

  Widget _buildBrowseTab() {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: 10,
      itemBuilder: (context, index) {
        return _buildGigCard(
          title: [
            'Mobile App Development',
            'Logo Design for Startup',
            'Content Writing - Blog Posts',
            'Social Media Management',
            'Video Editing - YouTube',
            'WordPress Website',
            'Data Entry Specialist',
            'Voice Over - Narration',
            'Translation Services',
            'Virtual Assistant',
          ][index % 10],
          budget: [150000, 25000, 35000, 80000, 45000, 100000, 20000, 30000, 50000, 40000][index % 10],
          category: ['Tech', 'Design', 'Writing', 'Marketing', 'Video', 'Tech', 'Business', 'Music', 'Writing', 'Business'][index % 10],
          proposals: index * 3 + 2,
          isRemote: index % 2 == 0,
          postedAt: 'Posted ${index + 1} days ago',
          onTap: () => context.push('/gigs/gig-${index + 1}'),
        );
      },
    );
  }

  Widget _buildMyGigsTab() {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: 3,
      itemBuilder: (context, index) {
        return _buildGigCard(
          title: ['Build an E-commerce App', 'Create Marketing Materials', 'Write Product Descriptions'][index],
          budget: [200000, 50000, 15000][index],
          category: ['Tech', 'Design', 'Writing'][index],
          proposals: [5, 12, 3][index],
          isRemote: true,
          postedAt: ['Active', 'In Progress', 'Completed'][index],
          status: ['active', 'in_progress', 'completed'][index],
          onTap: () {},
        );
      },
    );
  }

  Widget _buildProposalsTab() {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: 4,
      itemBuilder: (context, index) {
        return _buildProposalCard(
          gigTitle: ['Mobile App Dev', 'Logo Design', 'Content Writing', 'Website Dev'][index],
          proposedAmount: [140000, 22000, 30000, 90000][index],
          status: ['pending', 'accepted', 'rejected', 'pending'][index],
          submittedAt: '${index + 1} days ago',
        );
      },
    );
  }

  Widget _buildGigCard({
    required String title,
    required int budget,
    required String category,
    required int proposals,
    required bool isRemote,
    required String postedAt,
    String? status,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: AppColors.border),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 8,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: AppColors.primary.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Text(
                    category,
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.primary,
                    ),
                  ),
                ),
                const Spacer(),
                if (isRemote)
                  Row(
                    children: [
                      const Icon(
                        Icons.location_off_outlined,
                        size: 14,
                        color: AppColors.textTertiary,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        'Remote',
                        style: AppTypography.labelSmall.copyWith(
                          color: AppColors.textTertiary,
                        ),
                      ),
                    ],
                  ),
              ],
            ),
            const SizedBox(height: 12),
            Text(
              title,
              style: AppTypography.titleMedium.copyWith(
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Text(
                  _currencyFormat.format(budget),
                  style: AppTypography.titleSmall.copyWith(
                    color: AppColors.secondary,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const Spacer(),
                if (status != null)
                  _buildStatusBadge(status)
                else
                  Text(
                    '$proposals proposals',
                    style: AppTypography.bodySmall.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              postedAt,
              style: AppTypography.labelSmall.copyWith(
                color: AppColors.textTertiary,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildProposalCard({
    required String gigTitle,
    required int proposedAmount,
    required String status,
    required String submittedAt,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(
                child: Text(
                  gigTitle,
                  style: AppTypography.titleSmall.copyWith(
                    color: AppColors.textPrimary,
                  ),
                ),
              ),
              _buildStatusBadge(status),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            'Your bid: ${_currencyFormat.format(proposedAmount)}',
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.secondary,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            'Submitted $submittedAt',
            style: AppTypography.labelSmall.copyWith(
              color: AppColors.textTertiary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusBadge(String status) {
    Color color;
    String text;

    switch (status) {
      case 'active':
        color = AppColors.success;
        text = 'Active';
        break;
      case 'in_progress':
        color = AppColors.info;
        text = 'In Progress';
        break;
      case 'completed':
        color = AppColors.secondary;
        text = 'Completed';
        break;
      case 'accepted':
        color = AppColors.success;
        text = 'Accepted';
        break;
      case 'rejected':
        color = AppColors.error;
        text = 'Rejected';
        break;
      case 'pending':
      default:
        color = AppColors.warning;
        text = 'Pending';
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        text,
        style: AppTypography.labelSmall.copyWith(color: color),
      ),
    );
  }
}

// Stub screens for navigation
class GigDetailsScreen extends StatelessWidget {
  final String gigId;
  const GigDetailsScreen({super.key, required this.gigId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Gig Details')),
      body: Center(child: Text('Gig ID: $gigId')),
    );
  }
}

class CreateGigScreen extends StatelessWidget {
  const CreateGigScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Post a Gig')),
      body: const Center(child: Text('Create Gig Form')),
    );
  }
}
