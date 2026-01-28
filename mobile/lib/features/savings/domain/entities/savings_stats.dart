import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'savings_stats.freezed.dart';

/// Statistics for user's savings activity
@freezed
class SavingsStats with _$SavingsStats {
  const factory SavingsStats({
    required Money totalSaved,
    @Default(0) int activeCircles,
    @Default(0) int completedCircles,
    required Money totalEarned,
    required Money totalPayoutsReceived,
    @Default(100.0) double contributionRate,
    @Default(0) int totalContributions,
    @Default(0) int missedContributions,
  }) = _SavingsStats;

  const SavingsStats._();

  /// Total number of circles user has joined
  int get totalCircles => activeCircles + completedCircles;

  /// Check if user has good contribution history
  bool get hasGoodHistory => contributionRate >= 80;

  /// Number of successful contributions
  int get successfulContributions => totalContributions - missedContributions;
}
