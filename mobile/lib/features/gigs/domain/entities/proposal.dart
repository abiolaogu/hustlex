import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';
import 'gig.dart';

part 'proposal.freezed.dart';

/// Proposal status enum
enum ProposalStatus {
  pending,
  accepted,
  rejected,
  withdrawn,
}

/// Extension for ProposalStatus
extension ProposalStatusX on ProposalStatus {
  String get displayName {
    switch (this) {
      case ProposalStatus.pending:
        return 'Pending';
      case ProposalStatus.accepted:
        return 'Accepted';
      case ProposalStatus.rejected:
        return 'Rejected';
      case ProposalStatus.withdrawn:
        return 'Withdrawn';
    }
  }

  bool get isTerminal =>
      this == ProposalStatus.accepted ||
      this == ProposalStatus.rejected ||
      this == ProposalStatus.withdrawn;
}

/// Proposal entity representing a freelancer's bid on a gig
@freezed
class Proposal with _$Proposal implements Entity {
  const factory Proposal({
    required String id,
    required String gigId,
    required String freelancerId,
    required Money proposedAmount,
    required int deliveryDays,
    required String coverLetter,
    required ProposalStatus status,
    @Default([]) List<String> attachments,
    GigParticipant? freelancer,
    Gig? gig,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _Proposal;

  const Proposal._();

  /// Check if proposal is pending
  bool get isPending => status == ProposalStatus.pending;

  /// Check if proposal is accepted
  bool get isAccepted => status == ProposalStatus.accepted;

  /// Check if proposal can be withdrawn
  bool get canWithdraw => status == ProposalStatus.pending;

  /// Check if proposal can be modified
  bool get canModify => status == ProposalStatus.pending;
}
