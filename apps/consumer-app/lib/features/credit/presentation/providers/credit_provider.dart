import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/di/providers.dart';
import '../../data/models/credit_model.dart';
import '../../data/repositories/credit_repository.dart';

/// Credit overview state
class CreditState {
  final CreditScore? creditScore;
  final List<CreditScoreHistory> scoreHistory;
  final CreditEligibility? eligibility;
  final List<LoanProduct> loanProducts;
  final bool isLoading;
  final String? error;

  const CreditState({
    this.creditScore,
    this.scoreHistory = const [],
    this.eligibility,
    this.loanProducts = const [],
    this.isLoading = false,
    this.error,
  });

  // Computed properties
  int get score => creditScore?.score ?? 0;
  CreditTier get tier => creditScore?.tier ?? CreditTier.poor;
  double get maxLoanAmount => eligibility?.maxAmount ?? 0;
  bool get isEligible => eligibility?.isEligible ?? false;

  CreditState copyWith({
    CreditScore? creditScore,
    List<CreditScoreHistory>? scoreHistory,
    CreditEligibility? eligibility,
    List<LoanProduct>? loanProducts,
    bool? isLoading,
    String? error,
  }) {
    return CreditState(
      creditScore: creditScore ?? this.creditScore,
      scoreHistory: scoreHistory ?? this.scoreHistory,
      eligibility: eligibility ?? this.eligibility,
      loanProducts: loanProducts ?? this.loanProducts,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// Credit state notifier
class CreditNotifier extends StateNotifier<CreditState> {
  final CreditRepository _repository;

  CreditNotifier(this._repository) : super(const CreditState());

  /// Load credit score
  Future<void> loadCreditScore() async {
    state = state.copyWith(isLoading: true, error: null);

    final result = await _repository.getCreditScore();

    result.when(
      success: (score) {
        state = state.copyWith(creditScore: score, isLoading: false);
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
      },
    );
  }

  /// Load credit score history
  Future<void> loadScoreHistory({int months = 12}) async {
    final result = await _repository.getCreditScoreHistory(months: months);

    result.when(
      success: (history) {
        state = state.copyWith(scoreHistory: history);
      },
      failure: (_, __) {},
    );
  }

  /// Load loan eligibility
  Future<void> loadEligibility() async {
    final result = await _repository.checkEligibility();

    result.when(
      success: (eligibility) {
        state = state.copyWith(eligibility: eligibility);
      },
      failure: (_, __) {},
    );
  }

  /// Load loan products
  Future<void> loadLoanProducts() async {
    final result = await _repository.getLoanProducts();

    result.when(
      success: (products) {
        state = state.copyWith(loanProducts: products);
      },
      failure: (_, __) {},
    );
  }

  /// Load all credit data
  Future<void> loadAll() async {
    state = state.copyWith(isLoading: true, error: null);

    await Future.wait([
      loadCreditScore(),
      loadScoreHistory(),
      loadEligibility(),
      loadLoanProducts(),
    ]);

    state = state.copyWith(isLoading: false);
  }

  /// Refresh
  Future<void> refresh() => loadAll();
}

/// Main credit provider
final creditProvider = StateNotifierProvider<CreditNotifier, CreditState>((ref) {
  final repository = ref.watch(creditRepositoryProvider);
  return CreditNotifier(repository);
});

/// Credit score only provider
final creditScoreProvider = Provider<int>((ref) {
  return ref.watch(creditProvider).score;
});

/// Credit tier provider
final creditTierProvider = Provider<CreditTier>((ref) {
  return ref.watch(creditProvider).tier;
});

/// Loans list state
class LoansState {
  final List<Loan> loans;
  final bool isLoading;
  final String? error;

  const LoansState({
    this.loans = const [],
    this.isLoading = false,
    this.error,
  });

  // Computed properties
  List<Loan> get activeLoans =>
      loans.where((l) => l.status == LoanStatus.active || 
                         l.status == LoanStatus.disbursed).toList();
  
  List<Loan> get pendingLoans =>
      loans.where((l) => l.status == LoanStatus.pending || 
                         l.status == LoanStatus.approved).toList();
  
  List<Loan> get completedLoans =>
      loans.where((l) => l.status == LoanStatus.completed).toList();

  double get totalOutstanding =>
      activeLoans.fold(0, (sum, l) => sum + l.amountRemaining);

  bool get hasActiveLoan => activeLoans.isNotEmpty;

  LoansState copyWith({
    List<Loan>? loans,
    bool? isLoading,
    String? error,
  }) {
    return LoansState(
      loans: loans ?? this.loans,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// Loans notifier
class LoansNotifier extends StateNotifier<LoansState> {
  final CreditRepository _repository;

  LoansNotifier(this._repository) : super(const LoansState());

  /// Load all loans
  Future<void> loadLoans() async {
    state = state.copyWith(isLoading: true, error: null);

    final result = await _repository.getMyLoans();

    result.when(
      success: (loans) {
        state = state.copyWith(loans: loans, isLoading: false);
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
      },
    );
  }

  /// Refresh
  Future<void> refresh() => loadLoans();
}

/// Loans provider
final loansProvider = StateNotifierProvider<LoansNotifier, LoansState>((ref) {
  final repository = ref.watch(creditRepositoryProvider);
  return LoansNotifier(repository);
});

/// Active loan provider (for dashboard)
final activeLoanProvider = Provider<Loan?>((ref) {
  final loans = ref.watch(loansProvider).activeLoans;
  return loans.isNotEmpty ? loans.first : null;
});

/// Loan detail state
class LoanDetailState {
  final Loan? loan;
  final List<LoanRepayment> repayments;
  final LoanRepayment? nextRepayment;
  final bool isLoading;
  final String? error;

  const LoanDetailState({
    this.loan,
    this.repayments = const [],
    this.nextRepayment,
    this.isLoading = false,
    this.error,
  });

  // Computed properties
  double get progress => loan?.progress ?? 0;
  double get amountPaid => loan?.amountPaid ?? 0;
  double get amountRemaining => loan?.amountRemaining ?? 0;
  bool get isOverdue => loan?.isOverdue ?? false;

  LoanDetailState copyWith({
    Loan? loan,
    List<LoanRepayment>? repayments,
    LoanRepayment? nextRepayment,
    bool? isLoading,
    String? error,
  }) {
    return LoanDetailState(
      loan: loan ?? this.loan,
      repayments: repayments ?? this.repayments,
      nextRepayment: nextRepayment ?? this.nextRepayment,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// Loan detail notifier
class LoanDetailNotifier extends StateNotifier<LoanDetailState> {
  final CreditRepository _repository;
  final String loanId;

  LoanDetailNotifier(this._repository, this.loanId)
      : super(const LoanDetailState());

  /// Load loan details
  Future<void> loadLoan() async {
    state = state.copyWith(isLoading: true, error: null);

    final result = await _repository.getLoan(loanId);

    result.when(
      success: (loan) {
        state = state.copyWith(loan: loan, isLoading: false);
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
      },
    );
  }

  /// Load repayment history
  Future<void> loadRepayments() async {
    final result = await _repository.getRepaymentHistory(loanId);

    result.when(
      success: (repayments) {
        state = state.copyWith(repayments: repayments);
      },
      failure: (_, __) {},
    );
  }

  /// Load next repayment
  Future<void> loadNextRepayment() async {
    final result = await _repository.getNextRepayment(loanId);

    result.when(
      success: (repayment) {
        state = state.copyWith(nextRepayment: repayment);
      },
      failure: (_, __) {},
    );
  }

  /// Load all data
  Future<void> loadAll() async {
    state = state.copyWith(isLoading: true, error: null);

    await Future.wait([
      loadLoan(),
      loadRepayments(),
      loadNextRepayment(),
    ]);

    state = state.copyWith(isLoading: false);
  }

  /// Refresh
  Future<void> refresh() => loadAll();

  /// Make repayment
  Future<bool> makeRepayment({
    required double amount,
    required String pin,
  }) async {
    final result = await _repository.makeRepayment(
      loanId: loanId,
      amount: amount,
      pin: pin,
    );

    return result.when(
      success: (repayment) {
        // Refresh loan details after payment
        loadAll();
        return true;
      },
      failure: (_, __) => false,
    );
  }
}

/// Loan detail provider
final loanDetailProvider = StateNotifierProvider.family<LoanDetailNotifier,
    LoanDetailState, String>((ref, loanId) {
  final repository = ref.watch(creditRepositoryProvider);
  return LoanDetailNotifier(repository, loanId);
});

/// Loan application state
class LoanApplicationState {
  final bool isApplying;
  final Loan? loan;
  final String? error;
  
  // Form state
  final double amount;
  final int tenureDays;
  final LoanProduct? selectedProduct;

  const LoanApplicationState({
    this.isApplying = false,
    this.loan,
    this.error,
    this.amount = 0,
    this.tenureDays = 30,
    this.selectedProduct,
  });

  // Computed properties
  double get interestRate => selectedProduct?.interestRate ?? 0;
  double get totalRepayment => amount + (amount * interestRate / 100);
  double get monthlyPayment => totalRepayment / (tenureDays / 30);

  LoanApplicationState copyWith({
    bool? isApplying,
    Loan? loan,
    String? error,
    double? amount,
    int? tenureDays,
    LoanProduct? selectedProduct,
  }) {
    return LoanApplicationState(
      isApplying: isApplying ?? this.isApplying,
      loan: loan ?? this.loan,
      error: error,
      amount: amount ?? this.amount,
      tenureDays: tenureDays ?? this.tenureDays,
      selectedProduct: selectedProduct ?? this.selectedProduct,
    );
  }
}

/// Loan application notifier
class LoanApplicationNotifier extends StateNotifier<LoanApplicationState> {
  final CreditRepository _repository;

  LoanApplicationNotifier(this._repository)
      : super(const LoanApplicationState());

  /// Update amount
  void setAmount(double amount) {
    state = state.copyWith(amount: amount);
  }

  /// Update tenure
  void setTenure(int days) {
    state = state.copyWith(tenureDays: days);
  }

  /// Select product
  void selectProduct(LoanProduct product) {
    state = state.copyWith(
      selectedProduct: product,
      amount: product.minAmount,
      tenureDays: product.minTenureDays,
    );
  }

  /// Apply for loan
  Future<bool> applyForLoan({required String pin}) async {
    if (state.selectedProduct == null) return false;
    if (state.amount <= 0) return false;

    state = state.copyWith(isApplying: true, error: null);

    final result = await _repository.applyForLoan(
      LoanApplicationRequest(
        productId: state.selectedProduct!.id,
        amount: state.amount,
        tenureDays: state.tenureDays,
        pin: pin,
      ),
    );

    return result.when(
      success: (loan) {
        state = state.copyWith(isApplying: false, loan: loan);
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(isApplying: false, error: message);
        return false;
      },
    );
  }

  /// Reset state
  void reset() {
    state = const LoanApplicationState();
  }
}

/// Loan application provider
final loanApplicationProvider =
    StateNotifierProvider<LoanApplicationNotifier, LoanApplicationState>((ref) {
  final repository = ref.watch(creditRepositoryProvider);
  return LoanApplicationNotifier(repository);
});

/// Loan products provider
final loanProductsProvider = FutureProvider<List<LoanProduct>>((ref) async {
  final repository = ref.watch(creditRepositoryProvider);
  final result = await repository.getLoanProducts();
  return result.data ?? [];
});

/// Eligibility provider
final eligibilityProvider = FutureProvider<CreditEligibility?>((ref) async {
  final repository = ref.watch(creditRepositoryProvider);
  final result = await repository.checkEligibility();
  return result.data;
});

/// Credit score history provider
final creditScoreHistoryProvider =
    FutureProvider.family<List<CreditScoreHistory>, int>((ref, months) async {
  final repository = ref.watch(creditRepositoryProvider);
  final result = await repository.getCreditScoreHistory(months: months);
  return result.data ?? [];
});

/// Credit factors provider
final creditFactorsProvider = FutureProvider<CreditFactors?>((ref) async {
  final creditScore = ref.watch(creditProvider).creditScore;
  return creditScore?.factors;
});

/// Credit tips provider
final creditTipsProvider = FutureProvider<List<String>>((ref) async {
  final repository = ref.watch(creditRepositoryProvider);
  final result = await repository.getCreditTips();
  return result.data ?? [];
});
