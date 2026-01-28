import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/api/api_client.dart';
import '../../../../core/exceptions/api_exception.dart';

/// Credit API Service
/// Handles all credit scoring and loan-related API calls
class CreditService {
  final ApiClient _apiClient;

  CreditService(this._apiClient);

  /// Get user's credit profile (score, limit, history)
  Future<CreditProfile> getCreditProfile() async {
    try {
      final response = await _apiClient.get('/api/v1/credit/profile');
      return CreditProfile.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get credit score breakdown and factors
  Future<CreditScoreDetails> getCreditScoreDetails() async {
    try {
      final response = await _apiClient.get('/api/v1/credit/score');
      return CreditScoreDetails.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get credit score history
  Future<CreditScoreHistoryResponse> getCreditScoreHistory({
    int months = 12,
  }) async {
    try {
      final response = await _apiClient.get(
        '/api/v1/credit/score/history',
        queryParameters: {'months': months},
      );
      return CreditScoreHistoryResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Apply for a loan
  Future<LoanApplication> applyForLoan({
    required double amount,
    required int tenureMonths,
    required String purpose,
    String? additionalInfo,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/credit/loans/apply', data: {
        'amount': amount,
        'tenure_months': tenureMonths,
        'purpose': purpose,
        if (additionalInfo != null) 'additional_info': additionalInfo,
      });
      return LoanApplication.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get loan eligibility check (pre-qualification)
  Future<LoanEligibility> checkLoanEligibility({
    required double amount,
    required int tenureMonths,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/credit/loans/eligibility', data: {
        'amount': amount,
        'tenure_months': tenureMonths,
      });
      return LoanEligibility.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get all user's loans
  Future<LoansListResponse> getLoans({
    int page = 1,
    int perPage = 20,
    String? status, // pending, active, completed, defaulted
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (status != null) queryParams['status'] = status;

      final response = await _apiClient.get(
        '/api/v1/credit/loans',
        queryParameters: queryParams,
      );
      return LoansListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get single loan details
  Future<LoanDetail> getLoanDetails(String loanId) async {
    try {
      final response = await _apiClient.get('/api/v1/credit/loans/$loanId');
      return LoanDetail.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get loan repayment schedule
  Future<RepaymentScheduleResponse> getRepaymentSchedule(String loanId) async {
    try {
      final response = await _apiClient.get('/api/v1/credit/loans/$loanId/schedule');
      return RepaymentScheduleResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Make loan repayment
  Future<RepaymentResponse> makeRepayment(String loanId, {
    required double amount,
    String? pin,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/credit/loans/$loanId/repay', data: {
        'amount': amount,
        if (pin != null) 'pin': pin,
      });
      return RepaymentResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get loan repayment history
  Future<RepaymentsListResponse> getRepaymentHistory(String loanId, {
    int page = 1,
    int perPage = 20,
  }) async {
    try {
      final response = await _apiClient.get(
        '/api/v1/credit/loans/$loanId/repayments',
        queryParameters: {'page': page, 'per_page': perPage},
      );
      return RepaymentsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get loan statement (download URL)
  Future<LoanStatementResponse> getLoanStatement(String loanId, {
    String format = 'pdf', // pdf, csv
  }) async {
    try {
      final response = await _apiClient.get(
        '/api/v1/credit/loans/$loanId/statement',
        queryParameters: {'format': format},
      );
      return LoanStatementResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get credit tips/recommendations
  Future<CreditTipsResponse> getCreditTips() async {
    try {
      final response = await _apiClient.get('/api/v1/credit/tips');
      return CreditTipsResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get loan offers/products
  Future<LoanOffersResponse> getLoanOffers() async {
    try {
      final response = await _apiClient.get('/api/v1/credit/offers');
      return LoanOffersResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  ApiException _handleError(dynamic error) {
    if (error is ApiException) return error;
    return ApiException(message: error.toString());
  }
}

// Response models
class CreditProfile {
  final int creditScore;
  final String scoreRating; // Excellent, Good, Fair, Poor
  final double creditLimit;
  final double availableCredit;
  final double usedCredit;
  final int activeLoans;
  final int completedLoans;
  final double totalBorrowed;
  final double totalRepaid;
  final bool isEligibleForLoan;
  final String? ineligibleReason;
  final DateTime lastUpdated;

  CreditProfile({
    required this.creditScore,
    required this.scoreRating,
    required this.creditLimit,
    required this.availableCredit,
    required this.usedCredit,
    required this.activeLoans,
    required this.completedLoans,
    required this.totalBorrowed,
    required this.totalRepaid,
    required this.isEligibleForLoan,
    this.ineligibleReason,
    required this.lastUpdated,
  });

  factory CreditProfile.fromJson(Map<String, dynamic> json) {
    return CreditProfile(
      creditScore: json['credit_score'] ?? 0,
      scoreRating: json['score_rating'] ?? 'Unknown',
      creditLimit: (json['credit_limit'] ?? 0).toDouble(),
      availableCredit: (json['available_credit'] ?? 0).toDouble(),
      usedCredit: (json['used_credit'] ?? 0).toDouble(),
      activeLoans: json['active_loans'] ?? 0,
      completedLoans: json['completed_loans'] ?? 0,
      totalBorrowed: (json['total_borrowed'] ?? 0).toDouble(),
      totalRepaid: (json['total_repaid'] ?? 0).toDouble(),
      isEligibleForLoan: json['is_eligible_for_loan'] ?? false,
      ineligibleReason: json['ineligible_reason'],
      lastUpdated: DateTime.tryParse(json['last_updated'] ?? '') ?? DateTime.now(),
    );
  }

  double get utilizationPercentage => creditLimit > 0 ? (usedCredit / creditLimit) * 100 : 0;
}

class CreditScoreDetails {
  final int score;
  final String rating;
  final List<ScoreFactor> factors;
  final List<ScoreImpact> recentImpacts;
  final int maxScore;
  final int minScore;

  CreditScoreDetails({
    required this.score,
    required this.rating,
    required this.factors,
    required this.recentImpacts,
    this.maxScore = 850,
    this.minScore = 300,
  });

  factory CreditScoreDetails.fromJson(Map<String, dynamic> json) {
    return CreditScoreDetails(
      score: json['score'] ?? 0,
      rating: json['rating'] ?? 'Unknown',
      factors: (json['factors'] as List? ?? [])
          .map((f) => ScoreFactor.fromJson(f))
          .toList(),
      recentImpacts: (json['recent_impacts'] as List? ?? [])
          .map((i) => ScoreImpact.fromJson(i))
          .toList(),
      maxScore: json['max_score'] ?? 850,
      minScore: json['min_score'] ?? 300,
    );
  }

  double get scorePercentage => ((score - minScore) / (maxScore - minScore)) * 100;
}

class ScoreFactor {
  final String name;
  final String category; // payment_history, credit_utilization, account_age, etc.
  final double weight; // Percentage impact on score
  final String impact; // positive, negative, neutral
  final String description;
  final int? contribution; // Points contributed

  ScoreFactor({
    required this.name,
    required this.category,
    required this.weight,
    required this.impact,
    required this.description,
    this.contribution,
  });

  factory ScoreFactor.fromJson(Map<String, dynamic> json) {
    return ScoreFactor(
      name: json['name'] ?? '',
      category: json['category'] ?? '',
      weight: (json['weight'] ?? 0).toDouble(),
      impact: json['impact'] ?? 'neutral',
      description: json['description'] ?? '',
      contribution: json['contribution'],
    );
  }
}

class ScoreImpact {
  final String event;
  final int pointsChange;
  final String reason;
  final DateTime date;

  ScoreImpact({
    required this.event,
    required this.pointsChange,
    required this.reason,
    required this.date,
  });

  factory ScoreImpact.fromJson(Map<String, dynamic> json) {
    return ScoreImpact(
      event: json['event'] ?? '',
      pointsChange: json['points_change'] ?? 0,
      reason: json['reason'] ?? '',
      date: DateTime.tryParse(json['date'] ?? '') ?? DateTime.now(),
    );
  }
}

class CreditScoreHistoryResponse {
  final List<ScoreHistoryPoint> history;

  CreditScoreHistoryResponse({required this.history});

  factory CreditScoreHistoryResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return CreditScoreHistoryResponse(
      history: data.map((h) => ScoreHistoryPoint.fromJson(h)).toList(),
    );
  }
}

class ScoreHistoryPoint {
  final int score;
  final String rating;
  final DateTime date;
  final int? changeFromPrevious;

  ScoreHistoryPoint({
    required this.score,
    required this.rating,
    required this.date,
    this.changeFromPrevious,
  });

  factory ScoreHistoryPoint.fromJson(Map<String, dynamic> json) {
    return ScoreHistoryPoint(
      score: json['score'] ?? 0,
      rating: json['rating'] ?? '',
      date: DateTime.tryParse(json['date'] ?? '') ?? DateTime.now(),
      changeFromPrevious: json['change_from_previous'],
    );
  }
}

class LoanEligibility {
  final bool eligible;
  final double maxAmount;
  final double minAmount;
  final int maxTenureMonths;
  final int minTenureMonths;
  final double interestRate;
  final double estimatedMonthlyPayment;
  final List<String>? requirements;
  final String? rejectionReason;

  LoanEligibility({
    required this.eligible,
    required this.maxAmount,
    required this.minAmount,
    required this.maxTenureMonths,
    required this.minTenureMonths,
    required this.interestRate,
    required this.estimatedMonthlyPayment,
    this.requirements,
    this.rejectionReason,
  });

  factory LoanEligibility.fromJson(Map<String, dynamic> json) {
    return LoanEligibility(
      eligible: json['eligible'] ?? false,
      maxAmount: (json['max_amount'] ?? 0).toDouble(),
      minAmount: (json['min_amount'] ?? 0).toDouble(),
      maxTenureMonths: json['max_tenure_months'] ?? 12,
      minTenureMonths: json['min_tenure_months'] ?? 1,
      interestRate: (json['interest_rate'] ?? 0).toDouble(),
      estimatedMonthlyPayment: (json['estimated_monthly_payment'] ?? 0).toDouble(),
      requirements: json['requirements'] != null
          ? List<String>.from(json['requirements'])
          : null,
      rejectionReason: json['rejection_reason'],
    );
  }
}

class LoanApplication {
  final String id;
  final double amount;
  final int tenureMonths;
  final String purpose;
  final String status; // pending, approved, rejected, cancelled
  final double? approvedAmount;
  final double? interestRate;
  final String? rejectionReason;
  final DateTime createdAt;
  final DateTime? decidedAt;

  LoanApplication({
    required this.id,
    required this.amount,
    required this.tenureMonths,
    required this.purpose,
    required this.status,
    this.approvedAmount,
    this.interestRate,
    this.rejectionReason,
    required this.createdAt,
    this.decidedAt,
  });

  factory LoanApplication.fromJson(Map<String, dynamic> json) {
    return LoanApplication(
      id: json['id'] ?? '',
      amount: (json['amount'] ?? 0).toDouble(),
      tenureMonths: json['tenure_months'] ?? 0,
      purpose: json['purpose'] ?? '',
      status: json['status'] ?? 'pending',
      approvedAmount: json['approved_amount']?.toDouble(),
      interestRate: json['interest_rate']?.toDouble(),
      rejectionReason: json['rejection_reason'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      decidedAt: json['decided_at'] != null ? DateTime.tryParse(json['decided_at']) : null,
    );
  }
}

class LoansListResponse {
  final List<LoanSummary> loans;
  final PaginationMeta meta;

  LoansListResponse({required this.loans, required this.meta});

  factory LoansListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return LoansListResponse(
      loans: data.map((l) => LoanSummary.fromJson(l)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class LoanSummary {
  final String id;
  final double principalAmount;
  final double totalAmount;
  final double interestRate;
  final int tenureMonths;
  final String status; // active, completed, defaulted
  final double amountRepaid;
  final double amountRemaining;
  final DateTime? nextDueDate;
  final double? nextPaymentAmount;
  final int daysOverdue;
  final DateTime disbursedAt;
  final DateTime createdAt;

  LoanSummary({
    required this.id,
    required this.principalAmount,
    required this.totalAmount,
    required this.interestRate,
    required this.tenureMonths,
    required this.status,
    required this.amountRepaid,
    required this.amountRemaining,
    this.nextDueDate,
    this.nextPaymentAmount,
    required this.daysOverdue,
    required this.disbursedAt,
    required this.createdAt,
  });

  factory LoanSummary.fromJson(Map<String, dynamic> json) {
    return LoanSummary(
      id: json['id'] ?? '',
      principalAmount: (json['principal_amount'] ?? 0).toDouble(),
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      interestRate: (json['interest_rate'] ?? 0).toDouble(),
      tenureMonths: json['tenure_months'] ?? 0,
      status: json['status'] ?? 'active',
      amountRepaid: (json['amount_repaid'] ?? 0).toDouble(),
      amountRemaining: (json['amount_remaining'] ?? 0).toDouble(),
      nextDueDate: json['next_due_date'] != null ? DateTime.tryParse(json['next_due_date']) : null,
      nextPaymentAmount: json['next_payment_amount']?.toDouble(),
      daysOverdue: json['days_overdue'] ?? 0,
      disbursedAt: DateTime.tryParse(json['disbursed_at'] ?? '') ?? DateTime.now(),
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }

  double get repaymentProgress => totalAmount > 0 ? (amountRepaid / totalAmount) * 100 : 0;
  bool get isOverdue => daysOverdue > 0;
}

class LoanDetail extends LoanSummary {
  final String purpose;
  final List<ScheduledPayment> schedule;
  final List<Repayment> repayments;
  final double interestAmount;
  final double penaltyAmount;
  final double totalPaid;

  LoanDetail({
    required super.id,
    required super.principalAmount,
    required super.totalAmount,
    required super.interestRate,
    required super.tenureMonths,
    required super.status,
    required super.amountRepaid,
    required super.amountRemaining,
    super.nextDueDate,
    super.nextPaymentAmount,
    required super.daysOverdue,
    required super.disbursedAt,
    required super.createdAt,
    required this.purpose,
    required this.schedule,
    required this.repayments,
    required this.interestAmount,
    required this.penaltyAmount,
    required this.totalPaid,
  });

  factory LoanDetail.fromJson(Map<String, dynamic> json) {
    return LoanDetail(
      id: json['id'] ?? '',
      principalAmount: (json['principal_amount'] ?? 0).toDouble(),
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      interestRate: (json['interest_rate'] ?? 0).toDouble(),
      tenureMonths: json['tenure_months'] ?? 0,
      status: json['status'] ?? 'active',
      amountRepaid: (json['amount_repaid'] ?? 0).toDouble(),
      amountRemaining: (json['amount_remaining'] ?? 0).toDouble(),
      nextDueDate: json['next_due_date'] != null ? DateTime.tryParse(json['next_due_date']) : null,
      nextPaymentAmount: json['next_payment_amount']?.toDouble(),
      daysOverdue: json['days_overdue'] ?? 0,
      disbursedAt: DateTime.tryParse(json['disbursed_at'] ?? '') ?? DateTime.now(),
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      purpose: json['purpose'] ?? '',
      schedule: (json['schedule'] as List? ?? [])
          .map((s) => ScheduledPayment.fromJson(s))
          .toList(),
      repayments: (json['repayments'] as List? ?? [])
          .map((r) => Repayment.fromJson(r))
          .toList(),
      interestAmount: (json['interest_amount'] ?? 0).toDouble(),
      penaltyAmount: (json['penalty_amount'] ?? 0).toDouble(),
      totalPaid: (json['total_paid'] ?? 0).toDouble(),
    );
  }
}

class ScheduledPayment {
  final int installmentNumber;
  final double principalAmount;
  final double interestAmount;
  final double totalAmount;
  final DateTime dueDate;
  final String status; // pending, paid, overdue, partially_paid
  final double? amountPaid;
  final DateTime? paidAt;

  ScheduledPayment({
    required this.installmentNumber,
    required this.principalAmount,
    required this.interestAmount,
    required this.totalAmount,
    required this.dueDate,
    required this.status,
    this.amountPaid,
    this.paidAt,
  });

  factory ScheduledPayment.fromJson(Map<String, dynamic> json) {
    return ScheduledPayment(
      installmentNumber: json['installment_number'] ?? 0,
      principalAmount: (json['principal_amount'] ?? 0).toDouble(),
      interestAmount: (json['interest_amount'] ?? 0).toDouble(),
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      dueDate: DateTime.tryParse(json['due_date'] ?? '') ?? DateTime.now(),
      status: json['status'] ?? 'pending',
      amountPaid: json['amount_paid']?.toDouble(),
      paidAt: json['paid_at'] != null ? DateTime.tryParse(json['paid_at']) : null,
    );
  }
}

class RepaymentScheduleResponse {
  final List<ScheduledPayment> schedule;
  final double totalAmount;
  final double totalPrincipal;
  final double totalInterest;

  RepaymentScheduleResponse({
    required this.schedule,
    required this.totalAmount,
    required this.totalPrincipal,
    required this.totalInterest,
  });

  factory RepaymentScheduleResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return RepaymentScheduleResponse(
      schedule: data.map((s) => ScheduledPayment.fromJson(s)).toList(),
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      totalPrincipal: (json['total_principal'] ?? 0).toDouble(),
      totalInterest: (json['total_interest'] ?? 0).toDouble(),
    );
  }
}

class Repayment {
  final String id;
  final double amount;
  final String paymentMethod;
  final String status; // completed, pending, failed
  final String? transactionRef;
  final DateTime createdAt;

  Repayment({
    required this.id,
    required this.amount,
    required this.paymentMethod,
    required this.status,
    this.transactionRef,
    required this.createdAt,
  });

  factory Repayment.fromJson(Map<String, dynamic> json) {
    return Repayment(
      id: json['id'] ?? '',
      amount: (json['amount'] ?? 0).toDouble(),
      paymentMethod: json['payment_method'] ?? 'wallet',
      status: json['status'] ?? 'completed',
      transactionRef: json['transaction_ref'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class RepaymentResponse {
  final bool success;
  final String message;
  final Repayment? repayment;
  final double? newBalance;

  RepaymentResponse({
    required this.success,
    required this.message,
    this.repayment,
    this.newBalance,
  });

  factory RepaymentResponse.fromJson(Map<String, dynamic> json) {
    return RepaymentResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      repayment: json['data'] != null ? Repayment.fromJson(json['data']) : null,
      newBalance: json['new_balance']?.toDouble(),
    );
  }
}

class RepaymentsListResponse {
  final List<Repayment> repayments;
  final PaginationMeta meta;

  RepaymentsListResponse({required this.repayments, required this.meta});

  factory RepaymentsListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return RepaymentsListResponse(
      repayments: data.map((r) => Repayment.fromJson(r)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class LoanStatementResponse {
  final String downloadUrl;
  final String filename;
  final DateTime generatedAt;

  LoanStatementResponse({
    required this.downloadUrl,
    required this.filename,
    required this.generatedAt,
  });

  factory LoanStatementResponse.fromJson(Map<String, dynamic> json) {
    return LoanStatementResponse(
      downloadUrl: json['download_url'] ?? '',
      filename: json['filename'] ?? 'loan_statement.pdf',
      generatedAt: DateTime.tryParse(json['generated_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class CreditTipsResponse {
  final List<CreditTip> tips;

  CreditTipsResponse({required this.tips});

  factory CreditTipsResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return CreditTipsResponse(
      tips: data.map((t) => CreditTip.fromJson(t)).toList(),
    );
  }
}

class CreditTip {
  final String id;
  final String title;
  final String description;
  final String category;
  final int potentialImpact; // Points improvement possible
  final String? actionUrl;
  final bool isDismissed;

  CreditTip({
    required this.id,
    required this.title,
    required this.description,
    required this.category,
    required this.potentialImpact,
    this.actionUrl,
    this.isDismissed = false,
  });

  factory CreditTip.fromJson(Map<String, dynamic> json) {
    return CreditTip(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      description: json['description'] ?? '',
      category: json['category'] ?? '',
      potentialImpact: json['potential_impact'] ?? 0,
      actionUrl: json['action_url'],
      isDismissed: json['is_dismissed'] ?? false,
    );
  }
}

class LoanOffersResponse {
  final List<LoanOffer> offers;

  LoanOffersResponse({required this.offers});

  factory LoanOffersResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return LoanOffersResponse(
      offers: data.map((o) => LoanOffer.fromJson(o)).toList(),
    );
  }
}

class LoanOffer {
  final String id;
  final String name;
  final String description;
  final double minAmount;
  final double maxAmount;
  final int minTenure;
  final int maxTenure;
  final double interestRate;
  final double processingFee;
  final List<String> requirements;
  final bool isEligible;

  LoanOffer({
    required this.id,
    required this.name,
    required this.description,
    required this.minAmount,
    required this.maxAmount,
    required this.minTenure,
    required this.maxTenure,
    required this.interestRate,
    required this.processingFee,
    required this.requirements,
    required this.isEligible,
  });

  factory LoanOffer.fromJson(Map<String, dynamic> json) {
    return LoanOffer(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      minAmount: (json['min_amount'] ?? 0).toDouble(),
      maxAmount: (json['max_amount'] ?? 0).toDouble(),
      minTenure: json['min_tenure'] ?? 1,
      maxTenure: json['max_tenure'] ?? 12,
      interestRate: (json['interest_rate'] ?? 0).toDouble(),
      processingFee: (json['processing_fee'] ?? 0).toDouble(),
      requirements: List<String>.from(json['requirements'] ?? []),
      isEligible: json['is_eligible'] ?? false,
    );
  }
}

class PaginationMeta {
  final int currentPage;
  final int lastPage;
  final int perPage;
  final int total;

  PaginationMeta({
    required this.currentPage,
    required this.lastPage,
    required this.perPage,
    required this.total,
  });

  factory PaginationMeta.fromJson(Map<String, dynamic> json) {
    return PaginationMeta(
      currentPage: json['current_page'] ?? 1,
      lastPage: json['last_page'] ?? 1,
      perPage: json['per_page'] ?? 20,
      total: json['total'] ?? 0,
    );
  }

  bool get hasNextPage => currentPage < lastPage;
}

// Provider
final creditServiceProvider = Provider<CreditService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return CreditService(apiClient);
});
