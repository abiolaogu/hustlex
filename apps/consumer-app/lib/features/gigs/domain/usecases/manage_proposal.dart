import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../entities/proposal.dart';
import '../entities/contract.dart';
import '../repositories/gigs_repository.dart';

part 'manage_proposal.freezed.dart';

/// Parameters for getting proposals
class GetProposalsParams {
  final String gigId;
  final ProposalStatus? status;
  final int page;
  final int limit;

  const GetProposalsParams({
    required this.gigId,
    this.status,
    this.page = 1,
    this.limit = 20,
  });
}

/// Use case to get proposals for a gig
@injectable
class GetProposals implements UseCase<PaginatedProposals, GetProposalsParams> {
  final GigsRepository _repository;

  GetProposals(this._repository);

  @override
  Future<Either<Failure, PaginatedProposals>> call(GetProposalsParams params) {
    return _repository.getProposals(
      gigId: params.gigId,
      status: params.status,
      page: params.page,
      limit: params.limit,
    );
  }
}

/// Parameters for getting my proposals
class GetMyProposalsParams {
  final ProposalStatus? status;
  final int page;
  final int limit;

  const GetMyProposalsParams({
    this.status,
    this.page = 1,
    this.limit = 20,
  });
}

/// Use case to get current user's proposals (as freelancer)
@injectable
class GetMyProposals implements UseCase<PaginatedProposals, GetMyProposalsParams> {
  final GigsRepository _repository;

  GetMyProposals(this._repository);

  @override
  Future<Either<Failure, PaginatedProposals>> call(GetMyProposalsParams params) {
    return _repository.getMyProposals(
      status: params.status,
      page: params.page,
      limit: params.limit,
    );
  }
}

/// Parameters for submitting a proposal
@freezed
class SubmitProposalParams with _$SubmitProposalParams {
  const factory SubmitProposalParams({
    required String gigId,
    required Money proposedAmount,
    required int deliveryDays,
    required String coverLetter,
    @Default([]) List<String> attachments,
  }) = _SubmitProposalParams;
}

/// Use case to submit a proposal for a gig
@injectable
class SubmitProposal implements UseCase<Proposal, SubmitProposalParams> {
  final GigsRepository _repository;

  SubmitProposal(this._repository);

  @override
  Future<Either<Failure, Proposal>> call(SubmitProposalParams params) async {
    // Validate proposed amount
    if (params.proposedAmount.amount <= 0) {
      return left(
        const Failure.invalidInput('proposedAmount', 'Proposed amount must be greater than 0'),
      );
    }

    // Validate delivery days
    if (params.deliveryDays <= 0) {
      return left(
        const Failure.invalidInput('deliveryDays', 'Delivery days must be at least 1'),
      );
    }

    // Validate cover letter
    if (params.coverLetter.trim().isEmpty) {
      return left(
        const Failure.invalidInput('coverLetter', 'Cover letter cannot be empty'),
      );
    }
    if (params.coverLetter.trim().length < 50) {
      return left(
        const Failure.invalidInput(
          'coverLetter',
          'Cover letter must be at least 50 characters',
        ),
      );
    }

    return _repository.submitProposal(
      gigId: params.gigId,
      proposedAmount: params.proposedAmount,
      deliveryDays: params.deliveryDays,
      coverLetter: params.coverLetter.trim(),
      attachments: params.attachments,
    );
  }
}

/// Parameters for updating a proposal
@freezed
class UpdateProposalParams with _$UpdateProposalParams {
  const factory UpdateProposalParams({
    required String proposalId,
    Money? proposedAmount,
    int? deliveryDays,
    String? coverLetter,
    List<String>? attachments,
  }) = _UpdateProposalParams;
}

/// Use case to update an existing proposal
@injectable
class UpdateProposal implements UseCase<Proposal, UpdateProposalParams> {
  final GigsRepository _repository;

  UpdateProposal(this._repository);

  @override
  Future<Either<Failure, Proposal>> call(UpdateProposalParams params) async {
    // Validate proposed amount if provided
    if (params.proposedAmount != null && params.proposedAmount!.amount <= 0) {
      return left(
        const Failure.invalidInput('proposedAmount', 'Proposed amount must be greater than 0'),
      );
    }

    // Validate delivery days if provided
    if (params.deliveryDays != null && params.deliveryDays! <= 0) {
      return left(
        const Failure.invalidInput('deliveryDays', 'Delivery days must be at least 1'),
      );
    }

    // Validate cover letter if provided
    if (params.coverLetter != null && params.coverLetter!.trim().isEmpty) {
      return left(
        const Failure.invalidInput('coverLetter', 'Cover letter cannot be empty'),
      );
    }
    if (params.coverLetter != null && params.coverLetter!.trim().length < 50) {
      return left(
        const Failure.invalidInput(
          'coverLetter',
          'Cover letter must be at least 50 characters',
        ),
      );
    }

    return _repository.updateProposal(
      proposalId: params.proposalId,
      proposedAmount: params.proposedAmount,
      deliveryDays: params.deliveryDays,
      coverLetter: params.coverLetter?.trim(),
      attachments: params.attachments,
    );
  }
}

/// Use case to withdraw a proposal
@injectable
class WithdrawProposal implements UseCase<Unit, String> {
  final GigsRepository _repository;

  WithdrawProposal(this._repository);

  @override
  Future<Either<Failure, Unit>> call(String proposalId) {
    return _repository.withdrawProposal(proposalId);
  }
}

/// Use case to accept a proposal (creates a contract)
@injectable
class AcceptProposal implements UseCase<Contract, String> {
  final GigsRepository _repository;

  AcceptProposal(this._repository);

  @override
  Future<Either<Failure, Contract>> call(String proposalId) {
    return _repository.acceptProposal(proposalId);
  }
}

/// Use case to reject a proposal
@injectable
class RejectProposal implements UseCase<Unit, String> {
  final GigsRepository _repository;

  RejectProposal(this._repository);

  @override
  Future<Either<Failure, Unit>> call(String proposalId) {
    return _repository.rejectProposal(proposalId);
  }
}
