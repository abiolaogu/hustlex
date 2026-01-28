import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'milestone.freezed.dart';

/// Milestone status enum
enum MilestoneStatus {
  pending,
  inProgress,
  submitted,
  approved,
  revisionRequested,
  paid,
}

/// Extension for MilestoneStatus
extension MilestoneStatusX on MilestoneStatus {
  String get displayName {
    switch (this) {
      case MilestoneStatus.pending:
        return 'Pending';
      case MilestoneStatus.inProgress:
        return 'In Progress';
      case MilestoneStatus.submitted:
        return 'Submitted';
      case MilestoneStatus.approved:
        return 'Approved';
      case MilestoneStatus.revisionRequested:
        return 'Revision Requested';
      case MilestoneStatus.paid:
        return 'Paid';
    }
  }

  bool get isComplete =>
      this == MilestoneStatus.approved || this == MilestoneStatus.paid;
}

/// Milestone entity representing a deliverable within a contract
@freezed
class Milestone with _$Milestone implements Entity {
  const factory Milestone({
    required String id,
    required String contractId,
    required String title,
    String? description,
    required Money amount,
    required int order,
    required MilestoneStatus status,
    DateTime? dueDate,
    DateTime? submittedAt,
    DateTime? approvedAt,
    DateTime? paidAt,
    @Default([]) List<String> deliverables,
  }) = _Milestone;

  const Milestone._();

  /// Check if milestone is pending
  bool get isPending => status == MilestoneStatus.pending;

  /// Check if milestone is in progress
  bool get isInProgress => status == MilestoneStatus.inProgress;

  /// Check if milestone needs revision
  bool get needsRevision => status == MilestoneStatus.revisionRequested;

  /// Check if milestone is submitted and awaiting approval
  bool get awaitingApproval => status == MilestoneStatus.submitted;

  /// Check if milestone is completed (approved or paid)
  bool get isCompleted => status.isComplete;

  /// Check if milestone is overdue
  bool get isOverdue {
    if (dueDate == null || isCompleted) return false;
    return DateTime.now().isAfter(dueDate!);
  }
}
