import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/utils/currency_utils.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_text_field.dart';
import '../../data/models/savings_models.dart';

// Search query provider
final circleSearchQueryProvider = StateProvider<String>((ref) => '');

// Mock available circles
final availableCirclesProvider = FutureProvider<List<SavingsCircle>>((ref) async {
  await Future.delayed(const Duration(milliseconds: 500));
  
  return [
    SavingsCircle(
      id: 'circle-1',
      name: 'Tech Professionals Ajo',
      description: 'Monthly savings for tech workers. Build wealth together!',
      type: CircleType.rotating,
      frequency: ContributionFrequency.monthly,
      contributionAmount: 50000,
      totalSlots: 12,
      filledSlots: 8,
      currentRound: 3,
      totalRounds: 12,
      status: CircleStatus.active,
      startDate: DateTime(2024, 1, 1),
      nextContributionDate: DateTime.now().add(const Duration(days: 5)),
      createdBy: 'user-admin',
      members: [],
      contributions: [],
      createdAt: DateTime(2024, 1, 1),
      updatedAt: DateTime.now(),
    ),
    SavingsCircle(
      id: 'circle-2',
      name: 'Women in Business',
      description: 'Empowering women entrepreneurs through collective savings',
      type: CircleType.rotating,
      frequency: ContributionFrequency.weekly,
      contributionAmount: 10000,
      totalSlots: 10,
      filledSlots: 7,
      currentRound: 5,
      totalRounds: 10,
      status: CircleStatus.active,
      startDate: DateTime(2024, 2, 1),
      nextContributionDate: DateTime.now().add(const Duration(days: 2)),
      createdBy: 'user-admin',
      members: [],
      contributions: [],
      createdAt: DateTime(2024, 2, 1),
      updatedAt: DateTime.now(),
    ),
    SavingsCircle(
      id: 'circle-3',
      name: 'Young Hustlers',
      description: 'For young professionals building their financial future',
      type: CircleType.fixed,
      frequency: ContributionFrequency.weekly,
      contributionAmount: 5000,
      totalSlots: 20,
      filledSlots: 15,
      targetAmount: 500000,
      currentRound: 10,
      totalRounds: 20,
      status: CircleStatus.active,
      startDate: DateTime(2024, 3, 1),
      nextContributionDate: DateTime.now().add(const Duration(days: 1)),
      createdBy: 'user-admin',
      members: [],
      contributions: [],
      createdAt: DateTime(2024, 3, 1),
      updatedAt: DateTime.now(),
    ),
    SavingsCircle(
      id: 'circle-4',
      name: 'Landlord Goals',
      description: 'Saving towards property ownership',
      type: CircleType.fixed,
      frequency: ContributionFrequency.monthly,
      contributionAmount: 100000,
      totalSlots: 6,
      filledSlots: 4,
      targetAmount: 3600000,
      currentRound: 2,
      totalRounds: 6,
      status: CircleStatus.active,
      startDate: DateTime(2024, 1, 15),
      nextContributionDate: DateTime.now().add(const Duration(days: 10)),
      createdBy: 'user-admin',
      members: [],
      contributions: [],
      createdAt: DateTime(2024, 1, 15),
      updatedAt: DateTime.now(),
    ),
  ];
});

// Filtered circles based on search
final filteredCirclesProvider = Provider<AsyncValue<List<SavingsCircle>>>((ref) {
  final query = ref.watch(circleSearchQueryProvider).toLowerCase();
  final circlesAsync = ref.watch(availableCirclesProvider);
  
  return circlesAsync.whenData((circles) {
    if (query.isEmpty) return circles;
    return circles.where((circle) =>
      circle.name.toLowerCase().contains(query) ||
      circle.description.toLowerCase().contains(query)
    ).toList();
  });
});

class JoinCircleScreen extends ConsumerStatefulWidget {
  const JoinCircleScreen({super.key});

  @override
  ConsumerState<JoinCircleScreen> createState() => _JoinCircleScreenState();
}

class _JoinCircleScreenState extends ConsumerState<JoinCircleScreen> {
  final _searchController = TextEditingController();
  String? _selectedFilter;

  final _filters = ['All', 'Weekly', 'Monthly', 'Rotating', 'Fixed'];

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final filteredCircles = ref.watch(filteredCirclesProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Join a Circle'),
      ),
      body: Column(
        children: [
          _buildSearchSection(),
          _buildFilterChips(),
          Expanded(
            child: filteredCircles.when(
              data: (circles) => circles.isEmpty
                  ? _buildEmptyState()
                  : _buildCirclesList(circles),
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (e, _) => Center(
                child: Text('Error: $e'),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSearchSection() {
    return Container(
      padding: const EdgeInsets.all(16),
      color: AppColors.surface,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          AppTextField(
            controller: _searchController,
            hintText: 'Search circles...',
            prefixIcon: Icons.search,
            onChanged: (value) {
              ref.read(circleSearchQueryProvider.notifier).state = value;
            },
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              const Icon(
                Icons.info_outline,
                size: 16,
                color: AppColors.textSecondary,
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Text(
                  'Find and join savings circles that match your goals',
                  style: AppTypography.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildFilterChips() {
    return Container(
      height: 50,
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: 16),
        itemCount: _filters.length,
        separatorBuilder: (_, __) => const SizedBox(width: 8),
        itemBuilder: (context, index) {
          final filter = _filters[index];
          final isSelected = _selectedFilter == filter ||
              (_selectedFilter == null && filter == 'All');
          
          return FilterChip(
            label: Text(filter),
            selected: isSelected,
            onSelected: (selected) {
              setState(() {
                _selectedFilter = selected ? filter : null;
              });
            },
            selectedColor: AppColors.primary.withOpacity(0.2),
            checkmarkColor: AppColors.primary,
          );
        },
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
              Icons.search_off,
              size: 80,
              color: AppColors.textSecondary.withOpacity(0.5),
            ),
            const SizedBox(height: 16),
            Text(
              'No circles found',
              style: AppTypography.titleLarge.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Try adjusting your search or create your own circle',
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            AppButton(
              text: 'Create Circle',
              onPressed: () => context.push('/savings/circles/create'),
              variant: ButtonVariant.outlined,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCirclesList(List<SavingsCircle> circles) {
    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: circles.length,
      separatorBuilder: (_, __) => const SizedBox(height: 12),
      itemBuilder: (context, index) {
        return _CircleCard(
          circle: circles[index],
          onJoin: () => _showJoinDialog(circles[index]),
        );
      },
    );
  }

  void _showJoinDialog(SavingsCircle circle) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _JoinCircleSheet(circle: circle),
    );
  }
}

class _CircleCard extends StatelessWidget {
  final SavingsCircle circle;
  final VoidCallback onJoin;

  const _CircleCard({
    required this.circle,
    required this.onJoin,
  });

  @override
  Widget build(BuildContext context) {
    final availableSlots = circle.totalSlots - circle.filledSlots;
    final isAlmostFull = availableSlots <= 2;

    return Container(
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
                    Container(
                      padding: const EdgeInsets.all(10),
                      decoration: BoxDecoration(
                        color: AppColors.primary.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: Icon(
                        circle.type == CircleType.rotating
                            ? Icons.sync
                            : Icons.savings,
                        color: AppColors.primary,
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
                            style: AppTypography.titleMedium.copyWith(
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                          const SizedBox(height: 2),
                          Row(
                            children: [
                              _buildTag(
                                circle.type == CircleType.rotating
                                    ? 'Rotating'
                                    : 'Fixed',
                                AppColors.primary,
                              ),
                              const SizedBox(width: 6),
                              _buildTag(
                                _getFrequencyLabel(circle.frequency),
                                AppColors.secondary,
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Text(
                  circle.description,
                  style: AppTypography.bodyMedium.copyWith(
                    color: AppColors.textSecondary,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: _buildInfoItem(
                        'Contribution',
                        CurrencyUtils.formatNaira(circle.contributionAmount),
                      ),
                    ),
                    Container(
                      width: 1,
                      height: 30,
                      color: AppColors.border,
                    ),
                    Expanded(
                      child: _buildInfoItem(
                        'Available Slots',
                        '$availableSlots of ${circle.totalSlots}',
                        highlight: isAlmostFull,
                      ),
                    ),
                    Container(
                      width: 1,
                      height: 30,
                      color: AppColors.border,
                    ),
                    Expanded(
                      child: _buildInfoItem(
                        'Round',
                        '${circle.currentRound}/${circle.totalRounds}',
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          Container(
            decoration: BoxDecoration(
              color: AppColors.background,
              borderRadius: const BorderRadius.vertical(
                bottom: Radius.circular(12),
              ),
            ),
            padding: const EdgeInsets.all(12),
            child: Row(
              children: [
                if (isAlmostFull)
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: AppColors.warning.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(
                          Icons.warning_amber,
                          size: 14,
                          color: AppColors.warning,
                        ),
                        const SizedBox(width: 4),
                        Text(
                          'Almost full!',
                          style: AppTypography.labelSmall.copyWith(
                            color: AppColors.warning,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  ),
                const Spacer(),
                TextButton(
                  onPressed: () {},
                  child: const Text('View Details'),
                ),
                const SizedBox(width: 8),
                ElevatedButton(
                  onPressed: onJoin,
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 20,
                      vertical: 10,
                    ),
                  ),
                  child: const Text('Join'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTag(String label, Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
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

  Widget _buildInfoItem(String label, String value, {bool highlight = false}) {
    return Column(
      children: [
        Text(
          value,
          style: AppTypography.titleSmall.copyWith(
            fontWeight: FontWeight.w600,
            color: highlight ? AppColors.warning : AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          label,
          style: AppTypography.labelSmall.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  String _getFrequencyLabel(ContributionFrequency frequency) {
    switch (frequency) {
      case ContributionFrequency.daily:
        return 'Daily';
      case ContributionFrequency.weekly:
        return 'Weekly';
      case ContributionFrequency.biweekly:
        return 'Bi-weekly';
      case ContributionFrequency.monthly:
        return 'Monthly';
    }
  }
}

class _JoinCircleSheet extends ConsumerStatefulWidget {
  final SavingsCircle circle;

  const _JoinCircleSheet({required this.circle});

  @override
  ConsumerState<_JoinCircleSheet> createState() => _JoinCircleSheetState();
}

class _JoinCircleSheetState extends ConsumerState<_JoinCircleSheet> {
  int? _selectedSlot;
  bool _agreedToTerms = false;
  bool _isLoading = false;

  @override
  Widget build(BuildContext context) {
    final availableSlots = widget.circle.totalSlots - widget.circle.filledSlots;
    final slotsAvailable = List.generate(
      availableSlots,
      (i) => widget.circle.filledSlots + i + 1,
    );

    return Padding(
      padding: EdgeInsets.only(
        left: 16,
        right: 16,
        top: 16,
        bottom: MediaQuery.of(context).viewInsets.bottom + 16,
      ),
      child: SingleChildScrollView(
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
              'Join ${widget.circle.name}',
              style: AppTypography.titleLarge.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 24),
            
            // Circle Summary
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: AppColors.primary.withOpacity(0.05),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: AppColors.primary.withOpacity(0.2)),
              ),
              child: Column(
                children: [
                  _buildSummaryRow(
                    'Contribution Amount',
                    CurrencyUtils.formatNaira(widget.circle.contributionAmount),
                  ),
                  const Divider(height: 16),
                  _buildSummaryRow(
                    'Frequency',
                    _getFrequencyLabel(widget.circle.frequency),
                  ),
                  const Divider(height: 16),
                  _buildSummaryRow(
                    'Total Payout',
                    CurrencyUtils.formatNaira(
                      widget.circle.contributionAmount * widget.circle.totalSlots,
                    ),
                  ),
                  if (widget.circle.type == CircleType.rotating) ...[
                    const Divider(height: 16),
                    _buildSummaryRow(
                      'Your Turn',
                      'Round ${_selectedSlot ?? '?'}',
                    ),
                  ],
                ],
              ),
            ),
            const SizedBox(height: 24),

            // Slot Selection (for rotating circles)
            if (widget.circle.type == CircleType.rotating) ...[
              Text(
                'Select Your Payout Position',
                style: AppTypography.titleMedium.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Choose when you want to receive your payout',
                style: AppTypography.bodySmall.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
              const SizedBox(height: 12),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: slotsAvailable.map((slot) {
                  final isSelected = _selectedSlot == slot;
                  return GestureDetector(
                    onTap: () => setState(() => _selectedSlot = slot),
                    child: Container(
                      width: 50,
                      height: 50,
                      decoration: BoxDecoration(
                        color: isSelected
                            ? AppColors.primary
                            : AppColors.surface,
                        borderRadius: BorderRadius.circular(8),
                        border: Border.all(
                          color: isSelected
                              ? AppColors.primary
                              : AppColors.border,
                        ),
                      ),
                      alignment: Alignment.center,
                      child: Text(
                        '$slot',
                        style: AppTypography.titleMedium.copyWith(
                          color: isSelected
                              ? Colors.white
                              : AppColors.textPrimary,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  );
                }).toList(),
              ),
              const SizedBox(height: 24),
            ],

            // Terms agreement
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Checkbox(
                  value: _agreedToTerms,
                  onChanged: (value) {
                    setState(() => _agreedToTerms = value ?? false);
                  },
                ),
                Expanded(
                  child: GestureDetector(
                    onTap: () {
                      setState(() => _agreedToTerms = !_agreedToTerms);
                    },
                    child: Padding(
                      padding: const EdgeInsets.only(top: 12),
                      child: Text.rich(
                        TextSpan(
                          text: 'I agree to the ',
                          style: AppTypography.bodyMedium,
                          children: [
                            TextSpan(
                              text: 'Circle Terms & Conditions',
                              style: AppTypography.bodyMedium.copyWith(
                                color: AppColors.primary,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                            const TextSpan(
                              text: ' and understand that contributions are binding',
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),

            // Warning notice
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: AppColors.warning.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                children: [
                  Icon(
                    Icons.info_outline,
                    color: AppColors.warning,
                    size: 20,
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Text(
                      'First contribution of ${CurrencyUtils.formatNaira(widget.circle.contributionAmount)} will be deducted from your wallet upon joining.',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.warning,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // Join button
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _canJoin() ? _handleJoin : null,
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
                        'Join Circle (${CurrencyUtils.formatNaira(widget.circle.contributionAmount)})',
                      ),
              ),
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  Widget _buildSummaryRow(String label, String value) {
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
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }

  String _getFrequencyLabel(ContributionFrequency frequency) {
    switch (frequency) {
      case ContributionFrequency.daily:
        return 'Daily';
      case ContributionFrequency.weekly:
        return 'Weekly';
      case ContributionFrequency.biweekly:
        return 'Bi-weekly';
      case ContributionFrequency.monthly:
        return 'Monthly';
    }
  }

  bool _canJoin() {
    if (!_agreedToTerms) return false;
    if (widget.circle.type == CircleType.rotating && _selectedSlot == null) {
      return false;
    }
    return true;
  }

  void _handleJoin() async {
    setState(() => _isLoading = true);

    // Simulate API call
    await Future.delayed(const Duration(seconds: 2));

    if (mounted) {
      Navigator.pop(context);
      
      // Show success dialog
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
              const Text('Welcome!'),
            ],
          ),
          content: Text(
            'You have successfully joined ${widget.circle.name}. Your first contribution has been deducted.',
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(context);
                context.push('/savings/circles/${widget.circle.id}');
              },
              child: const Text('View Circle'),
            ),
          ],
        ),
      );
    }
  }
}
