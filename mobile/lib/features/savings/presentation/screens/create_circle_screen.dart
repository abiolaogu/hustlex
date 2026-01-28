import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/widgets/widgets.dart';
import '../../data/models/savings_models.dart';
import '../providers/savings_provider.dart';

class CreateCircleScreen extends ConsumerStatefulWidget {
  const CreateCircleScreen({super.key});

  @override
  ConsumerState<CreateCircleScreen> createState() => _CreateCircleScreenState();
}

class _CreateCircleScreenState extends ConsumerState<CreateCircleScreen> {
  final _formKey = GlobalKey<FormState>();
  final _pageController = PageController();
  
  int _currentStep = 0;
  bool _isLoading = false;

  // Form Controllers
  final _nameController = TextEditingController();
  final _contributionController = TextEditingController();
  final _maxMembersController = TextEditingController(text: '10');
  final _targetAmountController = TextEditingController();

  // Form State
  CircleType _selectedType = CircleType.rotational;
  ContributionFrequency _selectedFrequency = ContributionFrequency.weekly;

  @override
  void dispose() {
    _pageController.dispose();
    _nameController.dispose();
    _contributionController.dispose();
    _maxMembersController.dispose();
    _targetAmountController.dispose();
    super.dispose();
  }

  void _nextStep() {
    if (_currentStep == 0 && _nameController.text.isEmpty) {
      _showError('Please enter a circle name');
      return;
    }
    if (_currentStep == 1 && _contributionController.text.isEmpty) {
      _showError('Please enter contribution amount');
      return;
    }

    if (_currentStep < 2) {
      setState(() => _currentStep++);
      _pageController.nextPage(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
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

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: AppColors.error),
    );
  }

  Future<void> _createCircle() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    try {
      final contribution = double.tryParse(
        _contributionController.text.replaceAll(',', ''),
      ) ?? 0;
      final maxMembers = int.tryParse(_maxMembersController.text) ?? 10;
      
      double? targetAmount;
      if (_selectedType == CircleType.fixedTarget && 
          _targetAmountController.text.isNotEmpty) {
        targetAmount = double.tryParse(
          _targetAmountController.text.replaceAll(',', ''),
        );
      }

      final request = CreateCircleRequest(
        name: _nameController.text,
        type: _selectedType,
        contributionAmount: contribution,
        frequency: _selectedFrequency,
        maxMembers: maxMembers,
        targetAmount: targetAmount,
      );

      await ref.read(savingsProvider.notifier).createCircle(request);

      if (mounted) {
        _showSuccessDialog();
      }
    } catch (e) {
      if (mounted) {
        _showError(e.toString().replaceAll('Exception: ', ''));
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
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
              padding: const EdgeInsets.all(16),
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
            const SizedBox(height: 16),
            const Text(
              'Circle Created!',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            const Text(
              'Share the invite code with friends to start saving together.',
              textAlign: TextAlign.center,
              style: TextStyle(color: AppColors.textSecondary),
            ),
          ],
        ),
        actions: [
          SizedBox(
            width: double.infinity,
            child: PrimaryButton(
              text: 'View My Circles',
              onPressed: () {
                Navigator.of(context).pop();
                context.go('/savings');
              },
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Create Savings Circle'),
        centerTitle: true,
      ),
      body: Form(
        key: _formKey,
        child: Column(
          children: [
            // Progress Indicator
            _buildProgressIndicator(),
            
            // Form Pages
            Expanded(
              child: PageView(
                controller: _pageController,
                physics: const NeverScrollableScrollPhysics(),
                children: [
                  _buildTypeStep(),
                  _buildDetailsStep(),
                  _buildReviewStep(),
                ],
              ),
            ),
            
            // Navigation Buttons
            _buildNavigationButtons(),
          ],
        ),
      ),
    );
  }

  Widget _buildProgressIndicator() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md, vertical: 12),
      child: Row(
        children: List.generate(3, (index) {
          final isActive = index <= _currentStep;
          return Expanded(
            child: Row(
              children: [
                Expanded(
                  child: Container(
                    height: 4,
                    decoration: BoxDecoration(
                      color: isActive ? AppColors.secondary : AppColors.border,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                ),
                if (index < 2) const SizedBox(width: 4),
              ],
            ),
          );
        }),
      ),
    );
  }

  Widget _buildTypeStep() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Choose Circle Type',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          const Text(
            'Select how your savings circle will work.',
            style: TextStyle(color: AppColors.textSecondary),
          ),
          const SizedBox(height: 24),

          // Circle Name
          AppTextField(
            controller: _nameController,
            label: 'Circle Name',
            hintText: 'e.g., Family Savings, Office Ajo',
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter a name';
              }
              return null;
            },
          ),
          const SizedBox(height: 24),

          // Type Selection
          const Text(
            'Circle Type',
            style: TextStyle(fontWeight: FontWeight.w600),
          ),
          const SizedBox(height: 12),

          // Rotational (Ajo/Esusu)
          _buildTypeCard(
            type: CircleType.rotational,
            icon: Icons.refresh,
            title: 'Rotational (Ajo/Esusu)',
            description: 'Members take turns receiving the full pot. '
                'Traditional Nigerian thrift savings.',
            features: [
              'One member receives full pot each round',
              'Position determines payout order',
              'Everyone contributes, everyone receives',
            ],
          ),
          const SizedBox(height: 12),

          // Fixed Target
          _buildTypeCard(
            type: CircleType.fixedTarget,
            icon: Icons.savings,
            title: 'Fixed Target',
            description: 'Save together towards a common goal. '
                'Funds distributed when target is reached.',
            features: [
              'Set a savings goal for the group',
              'Everyone saves until target is met',
              'Equal distribution at completion',
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildTypeCard({
    required CircleType type,
    required IconData icon,
    required String title,
    required String description,
    required List<String> features,
  }) {
    final isSelected = _selectedType == type;

    return InkWell(
      onTap: () => setState(() => _selectedType = type),
      borderRadius: BorderRadius.circular(16),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isSelected
              ? AppColors.secondary.withOpacity(0.1)
              : AppColors.surface,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: isSelected ? AppColors.secondary : AppColors.border,
            width: isSelected ? 2 : 1,
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(
                    color: (isSelected ? AppColors.secondary : AppColors.textSecondary)
                        .withOpacity(0.1),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: Icon(
                    icon,
                    color: isSelected ? AppColors.secondary : AppColors.textSecondary,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    title,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: isSelected ? AppColors.secondary : AppColors.textPrimary,
                    ),
                  ),
                ),
                Radio<CircleType>(
                  value: type,
                  groupValue: _selectedType,
                  onChanged: (value) => setState(() => _selectedType = value!),
                  activeColor: AppColors.secondary,
                ),
              ],
            ),
            const SizedBox(height: 12),
            Text(
              description,
              style: const TextStyle(
                fontSize: 13,
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: 12),
            ...features.map((feature) => Padding(
              padding: const EdgeInsets.only(bottom: 4),
              child: Row(
                children: [
                  Icon(
                    Icons.check_circle,
                    size: 16,
                    color: isSelected ? AppColors.secondary : AppColors.textSecondary,
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      feature,
                      style: const TextStyle(fontSize: 12),
                    ),
                  ),
                ],
              ),
            )),
          ],
        ),
      ),
    );
  }

  Widget _buildDetailsStep() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Circle Details',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          const Text(
            'Set up how contributions will work.',
            style: TextStyle(color: AppColors.textSecondary),
          ),
          const SizedBox(height: 24),

          // Contribution Amount
          AmountTextField(
            controller: _contributionController,
            label: 'Contribution Amount',
            hintText: '5,000',
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please enter amount';
              }
              final amount = double.tryParse(value.replaceAll(',', ''));
              if (amount == null || amount < 100) {
                return 'Minimum is ₦100';
              }
              return null;
            },
          ),
          const SizedBox(height: 8),
          const Text(
            'Amount each member contributes per cycle.',
            style: TextStyle(fontSize: 12, color: AppColors.textSecondary),
          ),

          const SizedBox(height: 24),

          // Frequency
          const Text(
            'Contribution Frequency',
            style: TextStyle(fontWeight: FontWeight.w600),
          ),
          const SizedBox(height: 8),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: ContributionFrequency.values.map((freq) {
              final isSelected = _selectedFrequency == freq;
              return ChoiceChip(
                label: Text(freq.displayName),
                selected: isSelected,
                onSelected: (selected) {
                  if (selected) setState(() => _selectedFrequency = freq);
                },
                selectedColor: AppColors.secondary.withOpacity(0.2),
              );
            }).toList(),
          ),

          const SizedBox(height: 24),

          // Max Members
          AppTextField(
            controller: _maxMembersController,
            label: 'Maximum Members',
            hintText: '10',
            keyboardType: TextInputType.number,
            validator: (value) {
              if (value == null || value.isEmpty) return 'Required';
              final members = int.tryParse(value);
              if (members == null || members < 2) return 'Minimum 2 members';
              if (members > 50) return 'Maximum 50 members';
              return null;
            },
          ),
          const SizedBox(height: 8),
          const Text(
            'How many people can join this circle (2-50).',
            style: TextStyle(fontSize: 12, color: AppColors.textSecondary),
          ),

          // Target Amount (only for fixed target)
          if (_selectedType == CircleType.fixedTarget) ...[
            const SizedBox(height: 24),
            AmountTextField(
              controller: _targetAmountController,
              label: 'Target Amount',
              hintText: '500,000',
              validator: (value) {
                if (_selectedType == CircleType.fixedTarget) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter target amount';
                  }
                }
                return null;
              },
            ),
            const SizedBox(height: 8),
            const Text(
              'The total amount you want to save as a group.',
              style: TextStyle(fontSize: 12, color: AppColors.textSecondary),
            ),
          ],

          const SizedBox(height: 32),

          // Estimated Payout Info
          _buildEstimatedInfo(),
        ],
      ),
    );
  }

  Widget _buildEstimatedInfo() {
    final contribution = double.tryParse(
      _contributionController.text.replaceAll(',', ''),
    ) ?? 0;
    final maxMembers = int.tryParse(_maxMembersController.text) ?? 10;

    if (contribution <= 0) return const SizedBox.shrink();

    final totalPot = contribution * maxMembers;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.info.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.calculate, color: AppColors.info),
              const SizedBox(width: 8),
              const Text(
                'Estimated',
                style: TextStyle(
                  fontWeight: FontWeight.w600,
                  color: AppColors.info,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          if (_selectedType == CircleType.rotational) ...[
            _buildEstimateRow(
              'Pot per round',
              '₦${_formatAmount(totalPot)}',
            ),
            _buildEstimateRow(
              'Total rounds',
              '$maxMembers',
            ),
            _buildEstimateRow(
              'Total cycle',
              '$maxMembers ${_selectedFrequency.displayName.toLowerCase()}s',
            ),
          ] else ...[
            final targetAmount = double.tryParse(
              _targetAmountController.text.replaceAll(',', ''),
            ) ?? totalPot;
            final roundsNeeded = (targetAmount / totalPot).ceil();
            _buildEstimateRow(
              'Per member contribution',
              '₦${_formatAmount(contribution)} ${_selectedFrequency.displayName.toLowerCase()}',
            ),
            _buildEstimateRow(
              'Group savings per cycle',
              '₦${_formatAmount(totalPot)}',
            ),
            _buildEstimateRow(
              'Estimated time to target',
              '~$roundsNeeded ${_selectedFrequency.displayName.toLowerCase()}${roundsNeeded > 1 ? 's' : ''}',
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildEstimateRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(fontSize: 13)),
          Text(
            value,
            style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13),
          ),
        ],
      ),
    );
  }

  Widget _buildReviewStep() {
    final contribution = double.tryParse(
      _contributionController.text.replaceAll(',', ''),
    ) ?? 0;
    final maxMembers = int.tryParse(_maxMembersController.text) ?? 10;
    final targetAmount = _selectedType == CircleType.fixedTarget
        ? double.tryParse(_targetAmountController.text.replaceAll(',', ''))
        : contribution * maxMembers;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Review & Create',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          const Text(
            'Review your circle settings before creating.',
            style: TextStyle(color: AppColors.textSecondary),
          ),
          const SizedBox(height: 24),

          // Preview Card
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [
                  AppColors.secondary,
                  AppColors.secondary.withOpacity(0.8),
                ],
              ),
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Container(
                      padding: const EdgeInsets.all(10),
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.2),
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: Icon(
                        _selectedType == CircleType.rotational
                            ? Icons.refresh
                            : Icons.savings,
                        color: Colors.white,
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            _nameController.text.isEmpty
                                ? 'Untitled Circle'
                                : _nameController.text,
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          Text(
                            _selectedType.displayName,
                            style: TextStyle(
                              color: Colors.white.withOpacity(0.8),
                              fontSize: 13,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 20),
                const Divider(color: Colors.white24),
                const SizedBox(height: 16),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    _buildPreviewStat(
                      'Contribution',
                      '₦${_formatAmount(contribution)}',
                    ),
                    _buildPreviewStat(
                      'Frequency',
                      _selectedFrequency.displayName,
                    ),
                    _buildPreviewStat(
                      'Max Members',
                      maxMembers.toString(),
                    ),
                  ],
                ),
                if (targetAmount != null && targetAmount > 0) ...[
                  const SizedBox(height: 16),
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.15),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          _selectedType == CircleType.rotational
                              ? 'Pot per round: '
                              : 'Target: ',
                          style: TextStyle(
                            color: Colors.white.withOpacity(0.8),
                          ),
                        ),
                        Text(
                          '₦${_formatAmount(targetAmount)}',
                          style: const TextStyle(
                            color: Colors.white,
                            fontWeight: FontWeight.bold,
                            fontSize: 18,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ],
            ),
          ),

          const SizedBox(height: 24),

          // Rules Summary
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: AppColors.surface,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: AppColors.border),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Circle Rules',
                  style: TextStyle(fontWeight: FontWeight.w600),
                ),
                const SizedBox(height: 12),
                _buildRuleItem(
                  'You\'ll be the admin of this circle',
                ),
                _buildRuleItem(
                  'Members must contribute on time to stay in good standing',
                ),
                _buildRuleItem(
                  'Late payments may affect credit score',
                ),
                _buildRuleItem(
                  'You can invite members using a unique code',
                ),
              ],
            ),
          ),

          const SizedBox(height: 16),

          // Terms Notice
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: AppColors.warning.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              children: [
                const Icon(Icons.info_outline, color: AppColors.warning, size: 20),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    'By creating this circle, you agree to our savings circle terms.',
                    style: TextStyle(fontSize: 12, color: AppColors.textSecondary),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPreviewStat(String label, String value) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.bold,
            fontSize: 16,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          label,
          style: TextStyle(
            color: Colors.white.withOpacity(0.7),
            fontSize: 11,
          ),
        ),
      ],
    );
  }

  Widget _buildRuleItem(String text) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(Icons.check, size: 16, color: AppColors.secondary),
          const SizedBox(width: 8),
          Expanded(
            child: Text(text, style: const TextStyle(fontSize: 13)),
          ),
        ],
      ),
    );
  }

  Widget _buildNavigationButtons() {
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
        child: Row(
          children: [
            if (_currentStep > 0)
              Expanded(
                child: SecondaryButton(
                  text: 'Back',
                  onPressed: _previousStep,
                ),
              ),
            if (_currentStep > 0) const SizedBox(width: 12),
            Expanded(
              flex: _currentStep > 0 ? 2 : 1,
              child: PrimaryButton(
                text: _currentStep < 2 ? 'Continue' : 'Create Circle',
                onPressed: _isLoading
                    ? null
                    : (_currentStep < 2 ? _nextStep : _createCircle),
                isLoading: _isLoading,
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatAmount(double amount) {
    return amount.toStringAsFixed(0).replaceAllMapped(
          RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},',
        );
  }
}
