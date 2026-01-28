import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/loan_offer.dart';
import '../repositories/credit_repository.dart';

/// Use case to check loan eligibility
@injectable
class CheckLoanEligibility implements UseCaseNoParams<LoanEligibility> {
  final CreditRepository _repository;

  CheckLoanEligibility(this._repository);

  @override
  Future<Either<Failure, LoanEligibility>> call() {
    return _repository.checkEligibility();
  }
}

/// Use case to get available loan offers
@injectable
class GetLoanOffers implements UseCaseNoParams<List<LoanOffer>> {
  final CreditRepository _repository;

  GetLoanOffers(this._repository);

  @override
  Future<Either<Failure, List<LoanOffer>>> call() {
    return _repository.getLoanOffers();
  }
}

/// Use case to get specific loan offer details
@injectable
class GetLoanOffer implements UseCase<LoanOffer, String> {
  final CreditRepository _repository;

  GetLoanOffer(this._repository);

  @override
  Future<Either<Failure, LoanOffer>> call(String offerId) {
    return _repository.getLoanOffer(offerId);
  }
}
