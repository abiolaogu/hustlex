import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../../../../core/domain/value_objects/pin.dart';
import '../entities/loan.dart';
import '../entities/loan_repayment.dart';
import '../repositories/credit_repository.dart';

part 'manage_repayment.freezed.dart';

/// Use case to get repayments for a loan
@injectable
class GetLoanRepayments implements UseCase<List<LoanRepayment>, String> {
  final CreditRepository _repository;

  GetLoanRepayments(this._repository);

  @override
  Future<Either<Failure, List<LoanRepayment>>> call(String loanId) {
    return _repository.getLoanRepayments(loanId);
  }
}

/// Use case to get pending repayments
@injectable
class GetPendingRepayments implements UseCaseNoParams<List<LoanRepayment>> {
  final CreditRepository _repository;

  GetPendingRepayments(this._repository);

  @override
  Future<Either<Failure, List<LoanRepayment>>> call() {
    return _repository.getPendingRepayments();
  }
}

/// Use case to get next due repayment
@injectable
class GetNextDueRepayment implements UseCase<LoanRepayment?, String> {
  final CreditRepository _repository;

  GetNextDueRepayment(this._repository);

  @override
  Future<Either<Failure, LoanRepayment?>> call(String loanId) {
    return _repository.getNextDueRepayment(loanId);
  }
}

/// Parameters for making a repayment
@freezed
class MakeRepaymentParams with _$MakeRepaymentParams {
  const factory MakeRepaymentParams({
    required String loanId,
    required Money amount,
    required Pin pin,
    String? repaymentId,
  }) = _MakeRepaymentParams;
}

/// Use case to make a loan repayment
@injectable
class MakeRepayment implements UseCase<LoanRepayment, MakeRepaymentParams> {
  final CreditRepository _repository;

  MakeRepayment(this._repository);

  @override
  Future<Either<Failure, LoanRepayment>> call(MakeRepaymentParams params) async {
    // Validate PIN
    if (!params.pin.isValid()) {
      return left(const Failure.invalidPin());
    }

    // Validate amount
    if (params.amount.amount <= 0) {
      return left(
        const Failure.invalidInput('amount', 'Payment amount must be greater than 0'),
      );
    }

    return _repository.makeRepayment(
      loanId: params.loanId,
      amount: params.amount,
      pin: params.pin,
      repaymentId: params.repaymentId,
    );
  }
}

/// Parameters for paying off a loan
@freezed
class PayOffLoanParams with _$PayOffLoanParams {
  const factory PayOffLoanParams({
    required String loanId,
    required Pin pin,
  }) = _PayOffLoanParams;
}

/// Use case to pay off a loan in full
@injectable
class PayOffLoan implements UseCase<Loan, PayOffLoanParams> {
  final CreditRepository _repository;

  PayOffLoan(this._repository);

  @override
  Future<Either<Failure, Loan>> call(PayOffLoanParams params) async {
    // Validate PIN
    if (!params.pin.isValid()) {
      return left(const Failure.invalidPin());
    }

    return _repository.payOffLoan(
      loanId: params.loanId,
      pin: params.pin,
    );
  }
}

/// Parameters for setting up auto-debit
@freezed
class SetupAutoDebitParams with _$SetupAutoDebitParams {
  const factory SetupAutoDebitParams({
    required String loanId,
    required bool enabled,
  }) = _SetupAutoDebitParams;
}

/// Use case to set up auto-debit for loan repayments
@injectable
class SetupAutoDebit implements UseCase<Unit, SetupAutoDebitParams> {
  final CreditRepository _repository;

  SetupAutoDebit(this._repository);

  @override
  Future<Either<Failure, Unit>> call(SetupAutoDebitParams params) {
    return _repository.setupAutoDebit(
      loanId: params.loanId,
      enabled: params.enabled,
    );
  }
}
