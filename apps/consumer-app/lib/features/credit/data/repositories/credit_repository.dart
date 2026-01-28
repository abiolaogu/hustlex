import '../../../core/api/api_client.dart';
import '../../../core/repositories/base_repository.dart';
import '../models/credit_model.dart';

class CreditRepository extends BaseRepository {
  final ApiClient _apiClient;

  CreditRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  // Credit Score

  /// Get current user's credit score
  Future<Result<CreditScore>> getCreditScore() {
    return safeCall(() async {
      final response = await _apiClient.get('/credit/score');
      return CreditScore.fromJson(response.data['data']);
    });
  }

  /// Get credit score history
  Future<Result<List<CreditScoreHistory>>> getCreditScoreHistory({
    int months = 6,
  }) {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/credit/score/history',
        queryParameters: {'months': months},
      );
      final data = response.data['data'] as List;
      return data.map((e) => CreditScoreHistory.fromJson(e)).toList();
    });
  }

  /// Get tips to improve credit score
  Future<Result<List<CreditTip>>> getCreditTips() {
    return safeCall(() async {
      final response = await _apiClient.get('/credit/tips');
      final data = response.data['data'] as List;
      return data.map((e) => CreditTip.fromJson(e)).toList();
    });
  }

  /// Refresh credit score (manual refresh)
  Future<Result<CreditScore>> refreshCreditScore() {
    return safeCall(() async {
      final response = await _apiClient.post('/credit/score/refresh');
      return CreditScore.fromJson(response.data['data']);
    });
  }

  // Loans

  /// Get loan offer based on credit score
  Future<Result<LoanOffer>> getLoanOffer() {
    return safeCall(() async {
      final response = await _apiClient.get('/loans/offer');
      return LoanOffer.fromJson(response.data['data']);
    });
  }

  /// Apply for a loan
  Future<Result<Loan>> applyForLoan(LoanApplication application) {
    return safeCall(() async {
      final response = await _apiClient.post('/loans/apply', data: {
        'amount': application.amount,
        'tenor_months': application.tenorMonths,
        'repayment_frequency': application.repaymentFrequency.name,
        'purpose': application.purpose.name,
        if (application.purposeDescription != null)
          'purpose_description': application.purposeDescription,
        if (application.employmentStatus != null)
          'employment_status': application.employmentStatus,
        if (application.monthlyIncome != null)
          'monthly_income': application.monthlyIncome,
        if (application.guarantorName != null)
          'guarantor_name': application.guarantorName,
        if (application.guarantorPhone != null)
          'guarantor_phone': application.guarantorPhone,
      });
      return Loan.fromJson(response.data['data']);
    });
  }

  /// Get all user's loans with filters
  Future<Result<PaginatedLoans>> getLoans(LoanFilter filter) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{
        'page': filter.page,
        'limit': filter.limit,
      };

      if (filter.status != null) {
        queryParams['status'] = filter.status!.name;
      }
      if (filter.startDate != null) {
        queryParams['start_date'] = filter.startDate!.toIso8601String();
      }
      if (filter.endDate != null) {
        queryParams['end_date'] = filter.endDate!.toIso8601String();
      }

      final response = await _apiClient.get('/loans', queryParameters: queryParams);
      return PaginatedLoans.fromJson(response.data['data']);
    });
  }

  /// Get active loans
  Future<Result<List<Loan>>> getActiveLoans() {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/loans',
        queryParameters: {'status': 'active'},
      );
      final data = response.data['data']['loans'] as List;
      return data.map((e) => Loan.fromJson(e)).toList();
    });
  }

  /// Get a single loan by ID
  Future<Result<Loan>> getLoan(String loanId) {
    return safeCall(() async {
      final response = await _apiClient.get('/loans/$loanId');
      return Loan.fromJson(response.data['data']);
    });
  }

  /// Get loan repayment schedule
  Future<Result<List<LoanRepayment>>> getLoanRepayments(String loanId) {
    return safeCall(() async {
      final response = await _apiClient.get('/loans/$loanId/repayments');
      final data = response.data['data'] as List;
      return data.map((e) => LoanRepayment.fromJson(e)).toList();
    });
  }

  /// Make a loan payment
  Future<Result<LoanRepayment>> makeLoanPayment(MakeLoanPaymentRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/loans/${request.loanId}/pay', data: {
        'amount': request.amount,
        'pin': request.pin,
        if (request.repaymentId != null) 'repayment_id': request.repaymentId,
      });
      return LoanRepayment.fromJson(response.data['data']);
    });
  }

  /// Calculate loan details (preview before applying)
  Future<Result<Map<String, dynamic>>> calculateLoan({
    required double amount,
    required int tenorMonths,
    required RepaymentFrequency frequency,
  }) {
    return safeCall(() async {
      final response = await _apiClient.post('/loans/calculate', data: {
        'amount': amount,
        'tenor_months': tenorMonths,
        'repayment_frequency': frequency.name,
      });
      return response.data['data'] as Map<String, dynamic>;
    });
  }

  /// Get loan statistics
  Future<Result<Map<String, dynamic>>> getLoanStats() {
    return safeCall(() async {
      final response = await _apiClient.get('/loans/stats');
      return response.data['data'] as Map<String, dynamic>;
    });
  }

  // Eligibility

  /// Check loan eligibility
  Future<Result<Map<String, dynamic>>> checkEligibility() {
    return safeCall(() async {
      final response = await _apiClient.get('/loans/eligibility');
      return response.data['data'] as Map<String, dynamic>;
    });
  }
}
