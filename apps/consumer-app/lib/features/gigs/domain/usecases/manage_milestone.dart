import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../entities/milestone.dart';
import '../repositories/gigs_repository.dart';

part 'manage_milestone.freezed.dart';

/// Use case to get milestones for a contract
@injectable
class GetMilestones implements UseCase<List<Milestone>, String> {
  final GigsRepository _repository;

  GetMilestones(this._repository);

  @override
  Future<Either<Failure, List<Milestone>>> call(String contractId) {
    return _repository.getMilestones(contractId);
  }
}

/// Parameters for creating a milestone
@freezed
class CreateMilestoneParams with _$CreateMilestoneParams {
  const factory CreateMilestoneParams({
    required String contractId,
    required String title,
    String? description,
    required Money amount,
    required int order,
    DateTime? dueDate,
  }) = _CreateMilestoneParams;
}

/// Use case to create a milestone
@injectable
class CreateMilestone implements UseCase<Milestone, CreateMilestoneParams> {
  final GigsRepository _repository;

  CreateMilestone(this._repository);

  @override
  Future<Either<Failure, Milestone>> call(CreateMilestoneParams params) async {
    // Validate title
    if (params.title.trim().isEmpty) {
      return left(
        const Failure.invalidInput('title', 'Milestone title is required'),
      );
    }

    // Validate amount
    if (params.amount.amount <= 0) {
      return left(
        const Failure.invalidInput('amount', 'Milestone amount must be greater than 0'),
      );
    }

    // Validate order
    if (params.order < 0) {
      return left(
        const Failure.invalidInput('order', 'Milestone order cannot be negative'),
      );
    }

    return _repository.createMilestone(
      contractId: params.contractId,
      title: params.title.trim(),
      description: params.description?.trim(),
      amount: params.amount,
      order: params.order,
      dueDate: params.dueDate,
    );
  }
}

/// Use case to start working on a milestone
@injectable
class StartMilestone implements UseCase<Milestone, String> {
  final GigsRepository _repository;

  StartMilestone(this._repository);

  @override
  Future<Either<Failure, Milestone>> call(String milestoneId) {
    return _repository.startMilestone(milestoneId);
  }
}

/// Parameters for submitting a milestone
@freezed
class SubmitMilestoneParams with _$SubmitMilestoneParams {
  const factory SubmitMilestoneParams({
    required String milestoneId,
    required List<String> deliverables,
  }) = _SubmitMilestoneParams;
}

/// Use case to submit a milestone for review
@injectable
class SubmitMilestone implements UseCase<Milestone, SubmitMilestoneParams> {
  final GigsRepository _repository;

  SubmitMilestone(this._repository);

  @override
  Future<Either<Failure, Milestone>> call(SubmitMilestoneParams params) async {
    // Validate deliverables
    if (params.deliverables.isEmpty) {
      return left(
        const Failure.invalidInput('deliverables', 'At least one deliverable is required'),
      );
    }

    return _repository.submitMilestone(
      milestoneId: params.milestoneId,
      deliverables: params.deliverables,
    );
  }
}

/// Use case to approve a submitted milestone
@injectable
class ApproveMilestone implements UseCase<Milestone, String> {
  final GigsRepository _repository;

  ApproveMilestone(this._repository);

  @override
  Future<Either<Failure, Milestone>> call(String milestoneId) {
    return _repository.approveMilestone(milestoneId);
  }
}

/// Parameters for requesting revision
@freezed
class RequestRevisionParams with _$RequestRevisionParams {
  const factory RequestRevisionParams({
    required String milestoneId,
    required String feedback,
  }) = _RequestRevisionParams;
}

/// Use case to request revision on a milestone
@injectable
class RequestMilestoneRevision implements UseCase<Milestone, RequestRevisionParams> {
  final GigsRepository _repository;

  RequestMilestoneRevision(this._repository);

  @override
  Future<Either<Failure, Milestone>> call(RequestRevisionParams params) async {
    // Validate feedback
    if (params.feedback.trim().isEmpty) {
      return left(
        const Failure.invalidInput('feedback', 'Revision feedback is required'),
      );
    }
    if (params.feedback.trim().length < 20) {
      return left(
        const Failure.invalidInput(
          'feedback',
          'Please provide more detailed feedback (at least 20 characters)',
        ),
      );
    }

    return _repository.requestRevision(
      milestoneId: params.milestoneId,
      feedback: params.feedback.trim(),
    );
  }
}

/// Use case to release payment for an approved milestone
@injectable
class ReleaseMilestonePayment implements UseCase<Milestone, String> {
  final GigsRepository _repository;

  ReleaseMilestonePayment(this._repository);

  @override
  Future<Either<Failure, Milestone>> call(String milestoneId) {
    return _repository.releaseMilestonePayment(milestoneId);
  }
}
