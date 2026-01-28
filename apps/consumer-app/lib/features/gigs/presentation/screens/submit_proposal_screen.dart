import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';
import '../../../../core/utils/currency_helper.dart';
import '../../../../core/widgets/buttons.dart';
import '../../../../core/widgets/inputs.dart';
import '../../../../core/widgets/loaders.dart';
import '../../data/models/gig_model.dart';
import '../../data/repositories/gig_repository.dart';
import '../providers/gigs_provider.dart';

/// =============================================================================
/// SUBMIT PROPOSAL SCREEN
/// =============================================================================

class SubmitProposalScreen extends ConsumerStatefulWidget {
  final String gigId;

  const SubmitProposalScreen({
    super.key,
    required this.gigId,
  });

  @override
  ConsumerState<SubmitProposalScreen> createState() => _SubmitProposalScreenState();
}

class _SubmitProposalScreenState extends ConsumerState<SubmitProposalScreen> {
  final _formKey = GlobalKey<FormState>();
  final _coverLetterController = TextEditingController();
  final _amountController = TextEditingController();
  final _deliveryDaysController = TextEditingController();

  bool _isSubmitting = false;
  int _coverLetterLength = 0;

  @override
  void dispose() {
    _coverLetterController.dispose();
    _amountController.dispose();
    _deliveryDaysController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final gigAsync = ref.watch(gigDetailProvider(widget.gigId));

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Submit Proposal'),
        backgroundColor: AppColors.background,
        elevation: 0,
      ),
      body: gigAsync.when(
        loading: () => const Center(child: AppLoader()),
        error: (error, _) => _buildError(error.toString()),
        data: (gig) {
          if (gig == null) {
            return _buildError('Gig not found');
          }
          return _buildForm(gig);
        },
      ),
    );
  }

  Widget _buildError(String message) {
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
              text: 'Go Back',
              onPressed: () => context.pop(),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildForm(Gig gig) {
    return Column(
      children: [
        Expanded(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Gig summary card
                  _buildGigSummary(gig),

                  const SizedBox(height: 24),

                  // Cover letter section
                  _buildCoverLetterSection(),

                  const SizedBox(height: 24),

                  // Pricing section
                  _buildPricingSection(gig),

                  const SizedBox(height: 24),

                  // Delivery time section
                  _buildDeliverySection(gig),

                  const SizedBox(height: 24),

                  // Tips section
                  _buildTipsCard(),

                  const SizedBox(height: 24),
                ],
              ),
            ),
          ),
        ),

        // Submit button
        _buildBottomBar(gig),
      ],
    );
  }

  Widget _buildGigSummary(Gig gig) {
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
          Row(
            children: [
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: AppColors.primary.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(6),
                ),
                child: Text(
                  gig.category,
                  style: AppTypography.labelSmall.copyWith(
                    color: AppColors.primary,
                  ),
                ),
              ),
              const Spacer(),
              if (gig.isRemote)
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: AppColors.success.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      const Icon(
                        Icons.wifi,
                        size: 12,
                        color: AppColors.success,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        'Remote',
                        style: AppTypography.labelSmall.copyWith(
                          color: AppColors.success,
                        ),
                      ),
                    ],
                  ),
                ),
            ],
          ),
          const SizedBox(height: 12),
          Text(
            gig.title,
            style: AppTypography.titleMedium.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              const Icon(
                Icons.attach_money,
                size: 18,
                color: AppColors.textSecondary,
              ),
              const SizedBox(width: 4),
              Text(
                '${CurrencyHelper.formatNaira(gig.budgetMin)} - ${CurrencyHelper.formatNaira(gig.budgetMax)}',
                style: AppTypography.bodyMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
              const SizedBox(width: 16),
              const Icon(
                Icons.schedule,
                size: 18,
                color: AppColors.textSecondary,
              ),
              const SizedBox(width: 4),
              Text(
                '${gig.duration} days',
                style: AppTypography.bodyMedium.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildCoverLetterSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Cover Letter',
              style: AppTypography.titleSmall.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            Text(
              '$_coverLetterLength/1000',
              style: AppTypography.bodySmall.copyWith(
                color: _coverLetterLength < 100 
                    ? AppColors.error 
                    : AppColors.textSecondary,
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Text(
          'Introduce yourself and explain why you\'re the best fit for this gig.',
          style: AppTypography.bodySmall.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: 12),
        AppTextField(
          controller: _coverLetterController,
          hintText: 'Write your cover letter...\n\n'
              '• Highlight relevant experience\n'
              '• Explain your approach\n'
              '• Ask clarifying questions',
          maxLines: 8,
          maxLength: 1000,
          buildCounter: (_, {required currentLength, required isFocused, maxLength}) {
            return const SizedBox.shrink();
          },
          onChanged: (value) {
            setState(() => _coverLetterLength = value.length);
          },
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please write a cover letter';
            }
            if (value.length < 100) {
              return 'Cover letter must be at least 100 characters';
            }
            return null;
          },
        ),
        const SizedBox(height: 8),
        if (_coverLetterLength < 100)
          Text(
            'Minimum 100 characters required',
            style: AppTypography.bodySmall.copyWith(
              color: AppColors.error,
            ),
          ),
      ],
    );
  }

  Widget _buildPricingSection(Gig gig) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Your Bid',
          style: AppTypography.titleSmall.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Client\'s budget: ${CurrencyHelper.formatNaira(gig.budgetMin)} - ${CurrencyHelper.formatNaira(gig.budgetMax)}',
          style: AppTypography.bodySmall.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: 12),
        AppTextField(
          controller: _amountController,
          hintText: 'Enter your bid amount',
          keyboardType: TextInputType.number,
          prefixIcon: const Icon(Icons.attach_money),
          inputFormatters: [
            FilteringTextInputFormatter.digitsOnly,
          ],
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter your bid amount';
            }
            final amount = double.tryParse(value);
            if (amount == null) {
              return 'Please enter a valid amount';
            }
            if (amount < 1000) {
              return 'Minimum bid is ₦1,000';
            }
            return null;
          },
        ),

        const SizedBox(height: 12),

        // Quick amount buttons
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: [
            _buildQuickAmountChip(gig.budgetMin),
            _buildQuickAmountChip((gig.budgetMin + gig.budgetMax) / 2),
            _buildQuickAmountChip(gig.budgetMax),
          ],
        ),

        const SizedBox(height: 12),

        // Fee breakdown
        _buildFeeBreakdown(),
      ],
    );
  }

  Widget _buildQuickAmountChip(double amount) {
    final roundedAmount = (amount / 1000).round() * 1000;
    return ActionChip(
      label: Text(CurrencyHelper.formatNaira(roundedAmount.toDouble())),
      onPressed: () {
        _amountController.text = roundedAmount.toString();
        setState(() {});
      },
      backgroundColor: AppColors.surface,
      side: const BorderSide(color: AppColors.border),
    );
  }

  Widget _buildFeeBreakdown() {
    final amount = double.tryParse(_amountController.text) ?? 0;
    final platformFee = amount * 0.1; // 10% fee
    final youReceive = amount - platformFee;

    if (amount <= 0) return const SizedBox.shrink();

    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.info.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        children: [
          _buildFeeRow('Your bid', CurrencyHelper.formatNaira(amount)),
          _buildFeeRow('Platform fee (10%)', '-${CurrencyHelper.formatNaira(platformFee)}'),
          const Divider(height: 16),
          _buildFeeRow(
            'You\'ll receive',
            CurrencyHelper.formatNaira(youReceive),
            isBold: true,
          ),
        ],
      ),
    );
  }

  Widget _buildFeeRow(String label, String value, {bool isBold = false}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: isBold
                ? AppTypography.bodyMedium.copyWith(fontWeight: FontWeight.w600)
                : AppTypography.bodySmall.copyWith(color: AppColors.textSecondary),
          ),
          Text(
            value,
            style: isBold
                ? AppTypography.bodyMedium.copyWith(fontWeight: FontWeight.w600)
                : AppTypography.bodySmall,
          ),
        ],
      ),
    );
  }

  Widget _buildDeliverySection(Gig gig) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Delivery Time',
          style: AppTypography.titleSmall.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Client expects delivery within ${gig.duration} days',
          style: AppTypography.bodySmall.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: 12),
        AppTextField(
          controller: _deliveryDaysController,
          hintText: 'Days to deliver',
          keyboardType: TextInputType.number,
          prefixIcon: const Icon(Icons.schedule),
          inputFormatters: [
            FilteringTextInputFormatter.digitsOnly,
            LengthLimitingTextInputFormatter(3),
          ],
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter delivery time';
            }
            final days = int.tryParse(value);
            if (days == null || days < 1) {
              return 'Please enter a valid number of days';
            }
            if (days > 365) {
              return 'Maximum delivery time is 365 days';
            }
            return null;
          },
        ),

        const SizedBox(height: 12),

        // Quick delivery buttons
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: [
            _buildQuickDeliveryChip(gig.duration ~/ 2, 'Fast'),
            _buildQuickDeliveryChip(gig.duration, 'Standard'),
            _buildQuickDeliveryChip((gig.duration * 1.5).round(), 'Flexible'),
          ],
        ),
      ],
    );
  }

  Widget _buildQuickDeliveryChip(int days, String label) {
    return ActionChip(
      label: Text('$days days ($label)'),
      onPressed: () {
        _deliveryDaysController.text = days.toString();
        setState(() {});
      },
      backgroundColor: AppColors.surface,
      side: const BorderSide(color: AppColors.border),
    );
  }

  Widget _buildTipsCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.warning.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: AppColors.warning.withOpacity(0.3),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(
                Icons.lightbulb_outline,
                color: AppColors.warning,
                size: 20,
              ),
              const SizedBox(width: 8),
              Text(
                'Tips for a winning proposal',
                style: AppTypography.titleSmall.copyWith(
                  fontWeight: FontWeight.w600,
                  color: AppColors.warning,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          _buildTipItem('Read the gig description carefully'),
          _buildTipItem('Highlight relevant skills and experience'),
          _buildTipItem('Be specific about your approach'),
          _buildTipItem('Set a competitive but fair price'),
          _buildTipItem('Be realistic about delivery time'),
        ],
      ),
    );
  }

  Widget _buildTipItem(String tip) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(
            Icons.check_circle,
            size: 16,
            color: AppColors.warning,
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              tip,
              style: AppTypography.bodySmall.copyWith(
                color: AppColors.textPrimary,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomBar(Gig gig) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // Summary row
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Your Bid',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                    Text(
                      _amountController.text.isNotEmpty
                          ? CurrencyHelper.formatNaira(
                              double.tryParse(_amountController.text) ?? 0)
                          : '₦0',
                      style: AppTypography.titleMedium.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      'Delivery',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                    Text(
                      _deliveryDaysController.text.isNotEmpty
                          ? '${_deliveryDaysController.text} days'
                          : '-- days',
                      style: AppTypography.titleMedium.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 16),
            AppButton(
              text: _isSubmitting ? 'Submitting...' : 'Submit Proposal',
              onPressed: _isSubmitting ? null : () => _submitProposal(gig),
              isLoading: _isSubmitting,
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _submitProposal(Gig gig) async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isSubmitting = true);

    try {
      final repository = ref.read(gigRepositoryProvider);
      final result = await repository.submitProposal(
        gig.id,
        SubmitProposalRequest(
          coverLetter: _coverLetterController.text,
          amount: double.parse(_amountController.text),
          deliveryDays: int.parse(_deliveryDaysController.text),
        ),
      );

      result.when(
        success: (proposal) {
          // Invalidate relevant providers
          ref.invalidate(myGigsProvider);
          ref.invalidate(gigProposalsProvider(gig.id));

          if (mounted) {
            _showSuccessDialog();
          }
        },
        failure: (message, _) {
          setState(() => _isSubmitting = false);
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(message),
                backgroundColor: AppColors.error,
              ),
            );
          }
        },
      );
    } catch (e) {
      setState(() => _isSubmitting = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Failed to submit proposal'),
            backgroundColor: AppColors.error,
          ),
        );
      }
    }
  }

  void _showSuccessDialog() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 80,
              height: 80,
              decoration: const BoxDecoration(
                color: AppColors.success,
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.check,
                color: Colors.white,
                size: 48,
              ),
            ),
            const SizedBox(height: 24),
            Text(
              'Proposal Submitted!',
              style: AppTypography.titleLarge.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Your proposal has been sent to the client. '
              'You\'ll be notified when they respond.',
              style: AppTypography.bodyMedium.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            AppButton(
              text: 'View My Proposals',
              onPressed: () {
                context.pop(); // Close dialog
                context.pop(); // Go back to gig details
                context.go('/gigs?tab=proposals');
              },
            ),
            const SizedBox(height: 12),
            TextButton(
              onPressed: () {
                context.pop(); // Close dialog
                context.pop(); // Go back to gig details
              },
              child: const Text('Done'),
            ),
          ],
        ),
      ),
    );
  }
}
