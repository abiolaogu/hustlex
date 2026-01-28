import '../../../core/api/api_client.dart';
import '../../../core/repositories/base_repository.dart';
import '../models/wallet_model.dart';

class WalletRepository extends BaseRepository {
  final ApiClient _apiClient;

  WalletRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  // Wallet

  /// Get current user's wallet
  Future<Result<Wallet>> getWallet() {
    return safeCall(() async {
      final response = await _apiClient.get('/wallet');
      return Wallet.fromJson(response.data['data']);
    });
  }

  /// Get wallet balance only (lighter call)
  Future<Result<double>> getBalance() {
    return safeCall(() async {
      final response = await _apiClient.get('/wallet/balance');
      return (response.data['data']['available_balance'] as num).toDouble();
    });
  }

  /// Get wallet statistics
  Future<Result<WalletStats>> getWalletStats({
    DateTime? startDate,
    DateTime? endDate,
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (startDate != null) {
        queryParams['start_date'] = startDate.toIso8601String();
      }
      if (endDate != null) {
        queryParams['end_date'] = endDate.toIso8601String();
      }

      final response = await _apiClient.get('/wallet/stats', queryParameters: queryParams);
      return WalletStats.fromJson(response.data['data']);
    });
  }

  // Transactions

  /// Get paginated transactions with filters
  Future<Result<PaginatedTransactions>> getTransactions(TransactionFilter filter) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{
        'page': filter.page,
        'limit': filter.limit,
        'sort_by': filter.sortBy,
        'sort_order': filter.sortOrder,
      };

      if (filter.type != null) {
        queryParams['type'] = filter.type!.name;
      }
      if (filter.status != null) {
        queryParams['status'] = filter.status!.name;
      }
      if (filter.startDate != null) {
        queryParams['start_date'] = filter.startDate!.toIso8601String();
      }
      if (filter.endDate != null) {
        queryParams['end_date'] = filter.endDate!.toIso8601String();
      }
      if (filter.minAmount != null) {
        queryParams['min_amount'] = filter.minAmount;
      }
      if (filter.maxAmount != null) {
        queryParams['max_amount'] = filter.maxAmount;
      }
      if (filter.search != null && filter.search!.isNotEmpty) {
        queryParams['search'] = filter.search;
      }

      final response = await _apiClient.get('/transactions', queryParameters: queryParams);
      return PaginatedTransactions.fromJson(response.data['data']);
    });
  }

  /// Get recent transactions (for home screen)
  Future<Result<List<Transaction>>> getRecentTransactions({int limit = 5}) {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/transactions',
        queryParameters: {'limit': limit, 'sort_order': 'desc'},
      );
      final data = response.data['data']['transactions'] as List;
      return data.map((e) => Transaction.fromJson(e)).toList();
    });
  }

  /// Get a single transaction by ID
  Future<Result<Transaction>> getTransaction(String transactionId) {
    return safeCall(() async {
      final response = await _apiClient.get('/transactions/$transactionId');
      return Transaction.fromJson(response.data['data']);
    });
  }

  // Deposits

  /// Initialize deposit with Paystack
  Future<Result<PaystackInitResponse>> initializeDeposit(DepositRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/wallet/deposit/init', data: {
        'amount': request.amount,
        'method': request.method.name,
        if (request.metadata != null) 'metadata': request.metadata,
      });
      return PaystackInitResponse.fromJson(response.data['data']);
    });
  }

  /// Verify deposit after Paystack callback
  Future<Result<Transaction>> verifyDeposit(String reference) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/wallet/deposit/verify',
        data: {'reference': reference},
      );
      return Transaction.fromJson(response.data['data']);
    });
  }

  // Withdrawals

  /// Initiate withdrawal to bank account
  Future<Result<Transaction>> withdraw(WithdrawalRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/wallet/withdraw', data: {
        'amount': request.amount,
        'bank_account_id': request.bankAccountId,
        'pin': request.pin,
        if (request.narration != null) 'narration': request.narration,
      });
      return Transaction.fromJson(response.data['data']);
    });
  }

  // Transfers

  /// Transfer to another HustleX user
  Future<Result<Transaction>> transfer(TransferRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/wallet/transfer', data: {
        'amount': request.amount,
        'recipient_phone': request.recipientPhone,
        'pin': request.pin,
        if (request.narration != null) 'narration': request.narration,
      });
      return Transaction.fromJson(response.data['data']);
    });
  }

  /// Transfer to bank account
  Future<Result<Transaction>> bankTransfer(BankTransferRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/wallet/bank-transfer', data: {
        'amount': request.amount,
        'bank_code': request.bankCode,
        'account_number': request.accountNumber,
        'account_name': request.accountName,
        'pin': request.pin,
        if (request.narration != null) 'narration': request.narration,
      });
      return Transaction.fromJson(response.data['data']);
    });
  }

  // Bank Accounts

  /// Get user's saved bank accounts
  Future<Result<List<BankAccount>>> getBankAccounts() {
    return safeCall(() async {
      final response = await _apiClient.get('/bank-accounts');
      final data = response.data['data'] as List;
      return data.map((e) => BankAccount.fromJson(e)).toList();
    });
  }

  /// Add a new bank account
  Future<Result<BankAccount>> addBankAccount({
    required String bankCode,
    required String accountNumber,
  }) {
    return safeCall(() async {
      final response = await _apiClient.post('/bank-accounts', data: {
        'bank_code': bankCode,
        'account_number': accountNumber,
      });
      return BankAccount.fromJson(response.data['data']);
    });
  }

  /// Remove a bank account
  Future<Result<void>> removeBankAccount(String accountId) {
    return safeVoidCall(() async {
      await _apiClient.delete('/bank-accounts/$accountId');
    });
  }

  /// Set default bank account
  Future<Result<void>> setDefaultBankAccount(String accountId) {
    return safeVoidCall(() async {
      await _apiClient.post('/bank-accounts/$accountId/default');
    });
  }

  /// Verify bank account (name enquiry)
  Future<Result<AccountVerification>> verifyBankAccount({
    required String bankCode,
    required String accountNumber,
  }) {
    return safeCall(() async {
      final response = await _apiClient.post('/bank-accounts/verify', data: {
        'bank_code': bankCode,
        'account_number': accountNumber,
      });
      return AccountVerification.fromJson(response.data['data']);
    });
  }

  /// Get list of supported banks
  Future<Result<List<Bank>>> getBanks() {
    return safeCall(() async {
      final response = await _apiClient.get('/banks');
      final data = response.data['data'] as List;
      return data.map((e) => Bank.fromJson(e)).toList();
    });
  }

  // Bill Payments

  /// Pay a bill
  Future<Result<Transaction>> payBill(BillPaymentRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/bills/pay', data: {
        'category': request.category.name,
        'provider_id': request.providerId,
        'customer_id': request.customerId,
        'amount': request.amount,
        'pin': request.pin,
        if (request.phone != null) 'phone': request.phone,
        if (request.metadata != null) 'metadata': request.metadata,
      });
      return Transaction.fromJson(response.data['data']);
    });
  }

  /// Get bill providers for a category
  Future<Result<List<Map<String, dynamic>>>> getBillProviders(BillCategory category) {
    return safeCall(() async {
      final response = await _apiClient.get('/bills/providers/${category.name}');
      return (response.data['data'] as List).cast<Map<String, dynamic>>();
    });
  }

  /// Validate bill customer ID
  Future<Result<Map<String, dynamic>>> validateBillCustomer({
    required BillCategory category,
    required String providerId,
    required String customerId,
  }) {
    return safeCall(() async {
      final response = await _apiClient.post('/bills/validate', data: {
        'category': category.name,
        'provider_id': providerId,
        'customer_id': customerId,
      });
      return response.data['data'] as Map<String, dynamic>;
    });
  }

  // Airtime

  /// Buy airtime
  Future<Result<Transaction>> buyAirtime(AirtimeRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/airtime/buy', data: {
        'phone': request.phone,
        'amount': request.amount,
        'network': request.network,
        'pin': request.pin,
      });
      return Transaction.fromJson(response.data['data']);
    });
  }

  /// Get available networks
  Future<Result<List<String>>> getNetworks() {
    return safeCall(() async {
      final response = await _apiClient.get('/airtime/networks');
      return (response.data['data'] as List).cast<String>();
    });
  }

  /// Get data plans for a network
  Future<Result<List<Map<String, dynamic>>>> getDataPlans(String network) {
    return safeCall(() async {
      final response = await _apiClient.get('/data/plans/$network');
      return (response.data['data'] as List).cast<Map<String, dynamic>>();
    });
  }

  /// Buy data
  Future<Result<Transaction>> buyData({
    required String phone,
    required String network,
    required String planId,
    required String pin,
  }) {
    return safeCall(() async {
      final response = await _apiClient.post('/data/buy', data: {
        'phone': phone,
        'network': network,
        'plan_id': planId,
        'pin': pin,
      });
      return Transaction.fromJson(response.data['data']);
    });
  }
}
