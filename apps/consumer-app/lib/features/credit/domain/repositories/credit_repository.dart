import 'package:dartz/dartz.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../../../../core/domain/value_objects/pin.dart';
import '../entities/credit_score.dart';
import '../entities/loan.dart';
import '../entities/loan_repayment.dart';
import '../entities/loan_offer.dart';

/// Filter parameters for loan listings
class LoanFilter {
  final LoanStatus? status;
  final DateTime? startDate;
  final DateTime? endDate;
  final int page;
  final int limit;

  const LoanFilter({
    this.status,
    this.startDate,
    this.endDate,
    this.page = 1,
    this.limit = 20,
  });

  LoanFilter copyWith({
    LoanStatus? status,
    DateTime? startDate,
    DateTime? endDate,
    int? page,
    int? limit,
  }) {
    return LoanFilter(
      status: status ?? this.status,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      page: page ?? this.page,
      limit: limit ?? this.limit,
    );
  }
}

/// Paginated result for loans
class PaginatedLoans {
  final List<Loan> loans;
  final int total;
  final int page;
  final int limit;
  final bool hasMore;

  const PaginatedLoans({
    required this.loans,
    required this.total,
    required this.page,
    required this.limit,
    required this.hasMore,
  });
}

/// Abstract repository interface for credit feature
abstract class CreditRepository {
  // ============ CREDIT SCORE OPERATIONS ============

  /// Get current user's credit score
  Future<Either<Failure, CreditScore>> getCreditScore();

  /// Watch credit score for real-time updates
  Stream<Either<Failure, CreditScore>> watchCreditScore();

  /// Get credit score history
  Future<Either<Failure, List<CreditScoreHistory>>> getCreditScoreHistory({
    int page = 1,
    int limit = 20,
  });

  /// Get tips to improve credit score
  Future<Either<Failure, List<CreditTip>>> getCreditTips();

  /// Refresh credit score (trigger recalculation)
  Future<Either<Failure, CreditScore>> refreshCreditScore();

  // ============ LOAN ELIGIBILITY ============

  /// Check loan eligibility
  Future<Either<Failure, LoanEligibility>> checkEligibility();

  /// Get available loan offers
  Future<Either<Failure, List<LoanOffer>>> getLoanOffers();

  /// Get specific loan offer details
  Future<Either<Failure, LoanOffer>> getLoanOffer(String offerId);

  // ============ LOAN OPERATIONS ============

  /// Get user's loans
  Future<Either<Failure, PaginatedLoans>> getLoans(LoanFilter filter);

  /// Get active loan (if any)
  Future<Either<Failure, Loan?>> getActiveLoan();

  /// Get a single loan by ID
  Future<Either<Failure, Loan>> getLoan(String loanId);

  /// Watch a loan for real-time updates
  Stream<Either<Failure, Loan>> watchLoan(String loanId);

  /// Apply for a loan
  Future<Either<Failure, Loan>> applyForLoan({
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
  });

  /// Cancel a pending loan application
  Future<Either<Failure, Unit>> cancelLoanApplication(String loanId);

  /// Accept approved loan (trigger disbursement)
  Future<Either<Failure, Loan>> acceptLoanOffer(String loanId);

  // ============ REPAYMENT OPERATIONS ============

  /// Get repayments for a loan
  Future<Either<Failure, List<LoanRepayment>>> getLoanRepayments(String loanId);

  /// Get pending repayments for current user
  Future<Either<Failure, List<LoanRepayment>>> getPendingRepayments();

  /// Get next due repayment
  Future<Either<Failure, LoanRepayment?>> getNextDueRepayment(String loanId);

  /// Make a loan repayment
  Future<Either<Failure, LoanRepayment>> makeRepayment({
    required String loanId,
    required Money amount,
    required Pin pin,
    String? repaymentId,
  });

  /// Make full loan payoff
  Future<Either<Failure, Loan>> payOffLoan({
    required String loanId,
    required Pin pin,
  });

  /// Set up auto-debit for loan repayments
  Future<Either<Failure, Unit>> setupAutoDebit({
    required String loanId,
    required bool enabled,
  });

  // ============ LOAN HISTORY ============

  /// Get loan repayment history
  Future<Either<Failure, List<LoanRepayment>>> getRepaymentHistory({
    String? loanId,
    int page = 1,
    int limit = 20,
  });

  /// Get loan statistics
  Future<Either<Failure, LoanStats>> getLoanStats();
}

/// Loan statistics
class LoanStats {
  final int totalLoans;
  final int activeLoans;
  final int completedLoans;
  final int defaultedLoans;
  final Money totalBorrowed;
  final Money totalRepaid;
  final Money currentOutstanding;
  final double onTimePaymentRate;

  const LoanStats({
    required this.totalLoans,
    required this.activeLoans,
    required this.completedLoans,
    required this.defaultedLoans,
    required this.totalBorrowed,
    required this.totalRepaid,
    required this.currentOutstanding,
    required this.onTimePaymentRate,
  });

  bool get hasActiveLoan => activeLoans > 0;
  bool get hasGoodHistory => onTimePaymentRate >= 80;
}
