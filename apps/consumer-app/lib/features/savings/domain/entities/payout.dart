import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';
import 'savings_circle.dart';

part 'payout.freezed.dart';

/// Payout status
enum PayoutStatus {
  pending,
  processing,
  completed,
  failed,
}

/// Extension for PayoutStatus
extension PayoutStatusX on PayoutStatus {
  String get displayName {
    switch (this) {
      case PayoutStatus.pending:
        return 'Pending';
      case PayoutStatus.processing:
        return 'Processing';
      case PayoutStatus.completed:
        return 'Completed';
      case PayoutStatus.failed:
        return 'Failed';
    }
  }

  bool get isTerminal =>
      this == PayoutStatus.completed || this == PayoutStatus.failed;
}

/// Payout entity representing a disbursement to a circle member
@freezed
class Payout with _$Payout implements Entity {
  const factory Payout({
    required String id,
    required String circleId,
    required String memberId,
    String? memberName,
    required Money amount,
    required int cycleNumber,
    required PayoutStatus status,
    required DateTime scheduledDate,
    DateTime? processedAt,
    String? transactionId,
    String? failureReason,
    CircleMemberInfo? recipient,
    required DateTime createdAt,
  }) = _Payout;

  const Payout._();

  /// Check if payout is completed
  bool get isCompleted => status == PayoutStatus.completed;

  /// Check if payout is pending
  bool get isPending => status == PayoutStatus.pending;

  /// Check if payout failed
  bool get isFailed => status == PayoutStatus.failed;

  /// Check if payout is being processed
  bool get isProcessing => status == PayoutStatus.processing;

  /// Days until scheduled (negative if past due)
  int get daysUntilScheduled {
    if (isCompleted) return 0;
    return scheduledDate.difference(DateTime.now()).inDays;
  }
}
