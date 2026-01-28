import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'contribution.freezed.dart';

/// Contribution status
enum ContributionStatus {
  pending,
  paid,
  overdue,
  missed,
}

/// Extension for ContributionStatus
extension ContributionStatusX on ContributionStatus {
  String get displayName {
    switch (this) {
      case ContributionStatus.pending:
        return 'Pending';
      case ContributionStatus.paid:
        return 'Paid';
      case ContributionStatus.overdue:
        return 'Overdue';
      case ContributionStatus.missed:
        return 'Missed';
    }
  }

  bool get requiresAction =>
      this == ContributionStatus.pending || this == ContributionStatus.overdue;
}

/// Contribution entity representing a member's payment into the circle
@freezed
class Contribution with _$Contribution implements Entity {
  const factory Contribution({
    required String id,
    required String circleId,
    required String memberId,
    String? memberName,
    required Money amount,
    required int cycleNumber,
    required ContributionStatus status,
    required DateTime dueDate,
    DateTime? paidAt,
    String? transactionId,
    String? paymentMethod,
    required DateTime createdAt,
  }) = _Contribution;

  const Contribution._();

  /// Check if contribution is paid
  bool get isPaid => status == ContributionStatus.paid;

  /// Check if contribution is overdue
  bool get isOverdue => status == ContributionStatus.overdue;

  /// Check if contribution was missed
  bool get isMissed => status == ContributionStatus.missed;

  /// Check if contribution needs attention
  bool get needsAttention => status.requiresAction;

  /// Days until due (negative if overdue)
  int get daysUntilDue {
    if (isPaid) return 0;
    return dueDate.difference(DateTime.now()).inDays;
  }

  /// Check if contribution is due today
  bool get isDueToday {
    final now = DateTime.now();
    return dueDate.year == now.year &&
        dueDate.month == now.month &&
        dueDate.day == now.day;
  }
}
