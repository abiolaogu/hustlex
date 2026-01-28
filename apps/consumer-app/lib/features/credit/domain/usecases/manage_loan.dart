import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../entities/loan.dart';
import '../repositories/credit_repository.dart';

part 'manage_loan.freezed.dart';

/// Use case to get user's loans
@injectable
class GetLoans implements UseCase<PaginatedLoans, LoanFilter> {
  final CreditRepository _repository;

  GetLoans(this._repository);

  @override
  Future<Either<Failure, PaginatedLoans>> call(LoanFilter params) {
    return _repository.getLoans(params);
  }
}

/// Use case to get active loan
@injectable
class GetActiveLoan implements UseCaseNoParams<Loan?> {
  final CreditRepository _repository;

  GetActiveLoan(this._repository);

  @override
  Future<Either<Failure, Loan?>> call() {
    return _repository.getActiveLoan();
  }
}

/// Use case to get a single loan
@injectable
class GetLoan implements UseCase<Loan, String> {
  final CreditRepository _repository;

  GetLoan(this._repository);

  @override
  Future<Either<Failure, Loan>> call(String loanId) {
    return _repository.getLoan(loanId);
  }
}

/// Use case to watch a loan for real-time updates
@injectable
class WatchLoan implements StreamUseCase<Loan, String> {
  final CreditRepository _repository;

  WatchLoan(this._repository);

  @override
  Stream<Either<Failure, Loan>> call(String loanId) {
    return _repository.watchLoan(loanId);
  }
}

/// Parameters for loan application
@freezed
class ApplyForLoanParams with _$ApplyForLoanParams {
  const factory ApplyForLoanParams({
    required String offerId,
    required Money amount,
    required int tenorMonths,
    required RepaymentFrequency repaymentFrequency,
    required LoanPurpose purpose,
    String? purposeDescription,
    String? employmentStatus,
    Money? monthlyIncome,
    String? guarantorName,
    String? guarantorPhone,
  }) = _ApplyForLoanParams;
}

/// Use case to apply for a loan
@injectable
class ApplyForLoan implements UseCase<Loan, ApplyForLoanParams> {
  final CreditRepository _repository;

  ApplyForLoan(this._repository);

  @override
  Future<Either<Failure, Loan>> call(ApplyForLoanParams params) async {
    // Validate amount
    if (params.amount.amount <= 0) {
      return left(
        const Failure.invalidInput('amount', 'Loan amount must be greater than 0'),
      );
    }

    // Minimum loan amount
    if (params.amount.amount < 5000) {
      return left(
        const Failure.invalidInput('amount', 'Minimum loan amount is ₦5,000'),
      );
    }

    // Validate tenor
    if (params.tenorMonths <= 0) {
      return left(
        const Failure.invalidInput('tenorMonths', 'Loan tenor must be at least 1 month'),
      );
    }
    if (params.tenorMonths > 24) {
      return left(
        const Failure.invalidInput('tenorMonths', 'Maximum loan tenor is 24 months'),
      );
    }

    // Validate purpose description for "other" purpose
    if (params.purpose == LoanPurpose.other) {
      if (params.purposeDescription == null ||
          params.purposeDescription!.trim().isEmpty) {
        return left(
          const Failure.invalidInput(
            'purposeDescription',
            'Please describe the purpose of this loan',
          ),
        );
      }
    }

    return _repository.applyForLoan(
      offerId: params.offerId,
      amount: params.amount,
      tenorMonths: params.tenorMonths,
      repaymentFrequency: params.repaymentFrequency,
      purpose: params.purpose,
      purposeDescription: params.purposeDescription?.trim(),
      employmentStatus: params.employmentStatus,
      monthlyIncome: params.monthlyIncome,
      guarantorName: params.guarantorName?.trim(),
      guarantorPhone: params.guarantorPhone?.trim(),
    );
  }
}

/// Use case to cancel a pending loan application
@injectable
class CancelLoanApplication implements UseCase<Unit, String> {
  final CreditRepository _repository;

  CancelLoanApplication(this._repository);

  @override
  Future<Either<Failure, Unit>> call(String loanId) {
    return _repository.cancelLoanApplication(loanId);
  }
}

/// Use case to accept approved loan offer
@injectable
class AcceptLoanOffer implements UseCase<Loan, String> {
  final CreditRepository _repository;

  AcceptLoanOffer(this._repository);

  @override
  Future<Either<Failure, Loan>> call(String loanId) {
    return _repository.acceptLoanOffer(loanId);
  }
}
