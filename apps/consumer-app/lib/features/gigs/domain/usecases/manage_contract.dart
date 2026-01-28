import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/contract.dart';
import '../repositories/gigs_repository.dart';

part 'manage_contract.freezed.dart';

/// Parameters for getting contracts
class GetContractsParams {
  final ContractStatus? status;
  final int page;
  final int limit;

  const GetContractsParams({
    this.status,
    this.page = 1,
    this.limit = 20,
  });
}

/// Use case to get contracts for current user
@injectable
class GetContracts implements UseCase<PaginatedContracts, GetContractsParams> {
  final GigsRepository _repository;

  GetContracts(this._repository);

  @override
  Future<Either<Failure, PaginatedContracts>> call(GetContractsParams params) {
    return _repository.getContracts(
      status: params.status,
      page: params.page,
      limit: params.limit,
    );
  }
}

/// Use case to get a single contract
@injectable
class GetContract implements UseCase<Contract, String> {
  final GigsRepository _repository;

  GetContract(this._repository);

  @override
  Future<Either<Failure, Contract>> call(String contractId) {
    return _repository.getContract(contractId);
  }
}

/// Use case to watch a contract for real-time updates
@injectable
class WatchContract implements StreamUseCase<Contract, String> {
  final GigsRepository _repository;

  WatchContract(this._repository);

  @override
  Stream<Either<Failure, Contract>> call(String contractId) {
    return _repository.watchContract(contractId);
  }
}

/// Use case to complete a contract
@injectable
class CompleteContract implements UseCase<Contract, String> {
  final GigsRepository _repository;

  CompleteContract(this._repository);

  @override
  Future<Either<Failure, Contract>> call(String contractId) {
    return _repository.completeContract(contractId);
  }
}

/// Parameters for requesting cancellation
@freezed
class RequestCancellationParams with _$RequestCancellationParams {
  const factory RequestCancellationParams({
    required String contractId,
    required String reason,
  }) = _RequestCancellationParams;
}

/// Use case to request contract cancellation
@injectable
class RequestCancellation implements UseCase<Unit, RequestCancellationParams> {
  final GigsRepository _repository;

  RequestCancellation(this._repository);

  @override
  Future<Either<Failure, Unit>> call(RequestCancellationParams params) async {
    // Validate reason
    if (params.reason.trim().isEmpty) {
      return left(
        const Failure.invalidInput('reason', 'Cancellation reason is required'),
      );
    }
    if (params.reason.trim().length < 20) {
      return left(
        const Failure.invalidInput(
          'reason',
          'Please provide a more detailed reason (at least 20 characters)',
        ),
      );
    }

    return _repository.requestCancellation(
      contractId: params.contractId,
      reason: params.reason.trim(),
    );
  }
}

/// Parameters for raising a dispute
@freezed
class RaiseDisputeParams with _$RaiseDisputeParams {
  const factory RaiseDisputeParams({
    required String contractId,
    required String reason,
    @Default([]) List<String> evidence,
  }) = _RaiseDisputeParams;
}

/// Use case to raise a dispute on a contract
@injectable
class RaiseDispute implements UseCase<Unit, RaiseDisputeParams> {
  final GigsRepository _repository;

  RaiseDispute(this._repository);

  @override
  Future<Either<Failure, Unit>> call(RaiseDisputeParams params) async {
    // Validate reason
    if (params.reason.trim().isEmpty) {
      return left(
        const Failure.invalidInput('reason', 'Dispute reason is required'),
      );
    }
    if (params.reason.trim().length < 50) {
      return left(
        const Failure.invalidInput(
          'reason',
          'Please provide a detailed explanation of the dispute (at least 50 characters)',
        ),
      );
    }

    return _repository.raiseDispute(
      contractId: params.contractId,
      reason: params.reason.trim(),
      evidence: params.evidence,
    );
  }
}
