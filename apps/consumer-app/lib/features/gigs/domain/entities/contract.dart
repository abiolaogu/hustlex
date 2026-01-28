import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';
import 'gig.dart';

part 'contract.freezed.dart';

/// Contract status enum
enum ContractStatus {
  pending,
  active,
  completed,
  cancelled,
  disputed,
}

/// Extension for ContractStatus
extension ContractStatusX on ContractStatus {
  String get displayName {
    switch (this) {
      case ContractStatus.pending:
        return 'Pending';
      case ContractStatus.active:
        return 'Active';
      case ContractStatus.completed:
        return 'Completed';
      case ContractStatus.cancelled:
        return 'Cancelled';
      case ContractStatus.disputed:
        return 'Disputed';
    }
  }

  bool get isTerminal =>
      this == ContractStatus.completed || this == ContractStatus.cancelled;
}

/// Contract entity representing an agreement between client and freelancer
@freezed
class Contract with _$Contract implements Entity {
  const factory Contract({
    required String id,
    required String gigId,
    required String clientId,
    required String freelancerId,
    required Money amount,
    required Money escrowAmount,
    required Money platformFee,
    required int deliveryDays,
    required DateTime startDate,
    required DateTime dueDate,
    required ContractStatus status,
    DateTime? completedAt,
    GigParticipant? client,
    GigParticipant? freelancer,
    Gig? gig,
    required DateTime createdAt,
  }) = _Contract;

  const Contract._();

  /// Check if contract is active
  bool get isActive => status == ContractStatus.active;

  /// Check if contract is completed
  bool get isCompleted => status == ContractStatus.completed;

  /// Check if contract is overdue
  bool get isOverdue =>
      status == ContractStatus.active && DateTime.now().isAfter(dueDate);

  /// Get days remaining until due date
  int get daysRemaining {
    if (status != ContractStatus.active) return 0;
    return dueDate.difference(DateTime.now()).inDays;
  }

  /// Net amount after platform fee
  Money get netAmount => amount - platformFee;
}
