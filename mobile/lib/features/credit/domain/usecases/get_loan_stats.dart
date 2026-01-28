import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/loan_repayment.dart';
import '../repositories/credit_repository.dart';

/// Use case to get loan statistics
@injectable
class GetLoanStats implements UseCaseNoParams<LoanStats> {
  final CreditRepository _repository;

  GetLoanStats(this._repository);

  @override
  Future<Either<Failure, LoanStats>> call() {
    return _repository.getLoanStats();
  }
}

/// Parameters for getting repayment history
class GetRepaymentHistoryParams {
  final String? loanId;
  final int page;
  final int limit;

  const GetRepaymentHistoryParams({
    this.loanId,
    this.page = 1,
    this.limit = 20,
  });
}

/// Use case to get loan repayment history
@injectable
class GetRepaymentHistory
    implements UseCase<List<LoanRepayment>, GetRepaymentHistoryParams> {
  final CreditRepository _repository;

  GetRepaymentHistory(this._repository);

  @override
  Future<Either<Failure, List<LoanRepayment>>> call(
      GetRepaymentHistoryParams params) {
    return _repository.getRepaymentHistory(
      loanId: params.loanId,
      page: params.page,
      limit: params.limit,
    );
  }
}
