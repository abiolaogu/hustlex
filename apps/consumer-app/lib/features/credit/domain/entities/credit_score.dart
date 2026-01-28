import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'credit_score.freezed.dart';

/// Credit tier based on score ranges
enum CreditTier {
  poor, // 300-499
  fair, // 500-649
  good, // 650-749
  excellent, // 750-850
}

/// Extension for CreditTier
extension CreditTierX on CreditTier {
  String get displayName {
    switch (this) {
      case CreditTier.poor:
        return 'Poor';
      case CreditTier.fair:
        return 'Fair';
      case CreditTier.good:
        return 'Good';
      case CreditTier.excellent:
        return 'Excellent';
    }
  }

  int get minScore {
    switch (this) {
      case CreditTier.poor:
        return 300;
      case CreditTier.fair:
        return 500;
      case CreditTier.good:
        return 650;
      case CreditTier.excellent:
        return 750;
    }
  }

  int get maxScore {
    switch (this) {
      case CreditTier.poor:
        return 499;
      case CreditTier.fair:
        return 649;
      case CreditTier.good:
        return 749;
      case CreditTier.excellent:
        return 850;
    }
  }

  /// Maximum loan amount multiplier based on tier
  double get maxLoanMultiplier {
    switch (this) {
      case CreditTier.poor:
        return 0.5;
      case CreditTier.fair:
        return 1.0;
      case CreditTier.good:
        return 2.0;
      case CreditTier.excellent:
        return 3.0;
    }
  }

  /// Get credit tier from score
  static CreditTier fromScore(int score) {
    if (score >= 750) return CreditTier.excellent;
    if (score >= 650) return CreditTier.good;
    if (score >= 500) return CreditTier.fair;
    return CreditTier.poor;
  }
}

/// Credit score factors breakdown
@freezed
class CreditScoreFactors with _$CreditScoreFactors {
  const factory CreditScoreFactors({
    @Default(0) int paymentHistory, // 0-100
    @Default(0) int savingsConsistency, // 0-100
    @Default(0) int gigPerformance, // 0-100
    @Default(0) int accountAge, // 0-100
    @Default(0) int walletActivity, // 0-100
    @Default(0) int communityTrust, // 0-100
  }) = _CreditScoreFactors;

  const CreditScoreFactors._();

  /// Get all factors as a list
  List<CreditFactor> toList() => [
        CreditFactor(
          name: 'Payment History',
          description: 'Your track record of paying bills and loans on time',
          score: paymentHistory,
          weight: 0.25,
        ),
        CreditFactor(
          name: 'Savings Consistency',
          description: 'How regularly you contribute to savings circles',
          score: savingsConsistency,
          weight: 0.20,
        ),
        CreditFactor(
          name: 'Gig Performance',
          description: 'Your ratings and completion rate on gigs',
          score: gigPerformance,
          weight: 0.20,
        ),
        CreditFactor(
          name: 'Account Age',
          description: 'How long you have been using HustleX',
          score: accountAge,
          weight: 0.10,
        ),
        CreditFactor(
          name: 'Wallet Activity',
          description: 'Your transaction history and wallet usage',
          score: walletActivity,
          weight: 0.15,
        ),
        CreditFactor(
          name: 'Community Trust',
          description: 'Trust signals from savings circle members',
          score: communityTrust,
          weight: 0.10,
        ),
      ];
}

/// Individual credit factor
@freezed
class CreditFactor with _$CreditFactor {
  const factory CreditFactor({
    required String name,
    required String description,
    required int score,
    required double weight,
  }) = _CreditFactor;

  const CreditFactor._();

  /// Get rating label based on score
  String get rating {
    if (score >= 80) return 'Excellent';
    if (score >= 60) return 'Good';
    if (score >= 40) return 'Fair';
    return 'Needs Work';
  }

  /// Progress value (0.0 - 1.0)
  double get progress => score / 100.0;
}

/// Credit score history entry
@freezed
class CreditScoreHistory with _$CreditScoreHistory {
  const factory CreditScoreHistory({
    required int score,
    required CreditTier tier,
    required DateTime date,
    int? delta,
    String? reason,
  }) = _CreditScoreHistory;
}

/// Credit score entity
@freezed
class CreditScore with _$CreditScore {
  const factory CreditScore({
    required String userId,
    required int score,
    required CreditTier tier,
    @Default(0) int scoreDelta,
    required Money maxLoanAmount,
    required CreditScoreFactors factors,
    @Default([]) List<CreditScoreHistory> history,
    DateTime? lastUpdated,
    DateTime? nextUpdate,
  }) = _CreditScore;

  const CreditScore._();

  /// Score progress (0.0 - 1.0 based on 300-850 range)
  double get scoreProgress => (score - 300) / 550;

  /// Score percentage (0 - 100)
  int get scorePercent => (scoreProgress * 100).round();

  /// Check if score is improving
  bool get isImproving => scoreDelta > 0;

  /// Check if score is declining
  bool get isDeclining => scoreDelta < 0;

  /// Check if user is eligible for loans
  bool get isLoanEligible => score >= CreditTier.fair.minScore;
}

/// Tips to improve credit score
@freezed
class CreditTip with _$CreditTip {
  const factory CreditTip({
    required String title,
    required String description,
    int? potentialScoreIncrease,
  }) = _CreditTip;
}
