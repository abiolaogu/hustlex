import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/api/api_client.dart';
import '../../../../core/exceptions/api_exception.dart';
import '../models/wallet_models.dart';

/// Wallet API Service
/// Handles all wallet-related API calls to the backend
class WalletService {
  final ApiClient _apiClient;

  WalletService(this._apiClient);

  /// Get wallet details for the authenticated user
  Future<WalletResponse> getWallet() async {
    try {
      final response = await _apiClient.get('/api/v1/wallet');
      return WalletResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get wallet transactions with optional pagination and filters
  Future<TransactionsResponse> getTransactions({
    int page = 1,
    int perPage = 20,
    String? type, // credit, debit
    String? category,
    DateTime? startDate,
    DateTime? endDate,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (type != null) queryParams['type'] = type;
      if (category != null) queryParams['category'] = category;
      if (startDate != null) queryParams['start_date'] = startDate.toIso8601String();
      if (endDate != null) queryParams['end_date'] = endDate.toIso8601String();

      final response = await _apiClient.get(
        '/api/v1/wallet/transactions',
        queryParameters: queryParams,
      );
      return TransactionsResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get single transaction details
  Future<Transaction> getTransactionDetails(String transactionId) async {
    try {
      final response = await _apiClient.get('/api/v1/wallet/transactions/$transactionId');
      return Transaction.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Initialize deposit (get payment link from Paystack)
  Future<DepositInitResponse> initializeDeposit({
    required double amount,
    required String channel, // card, bank_transfer, ussd
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/wallet/deposit', data: {
        'amount': amount,
        'channel': channel,
      });
      return DepositInitResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Verify deposit after Paystack callback
  Future<Transaction> verifyDeposit(String reference) async {
    try {
      final response = await _apiClient.post('/api/v1/wallet/deposit/verify', data: {
        'reference': reference,
      });
      return Transaction.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Withdraw to bank account
  Future<WithdrawalResponse> withdraw({
    required double amount,
    required String bankAccountId,
    String? pin,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/wallet/withdraw', data: {
        'amount': amount,
        'bank_account_id': bankAccountId,
        if (pin != null) 'pin': pin,
      });
      return WithdrawalResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Transfer to another HustleX user
  Future<TransferResponse> transfer({
    required double amount,
    required String recipientPhone,
    String? narration,
    String? pin,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/wallet/transfer', data: {
        'amount': amount,
        'recipient_phone': recipientPhone,
        if (narration != null) 'narration': narration,
        if (pin != null) 'pin': pin,
      });
      return TransferResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Look up recipient by phone number before transfer
  Future<RecipientLookupResponse> lookupRecipient(String phone) async {
    try {
      final response = await _apiClient.get(
        '/api/v1/wallet/transfer/lookup',
        queryParameters: {'phone': phone},
      );
      return RecipientLookupResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get all linked bank accounts
  Future<BankAccountsResponse> getBankAccounts() async {
    try {
      final response = await _apiClient.get('/api/v1/wallet/bank-accounts');
      return BankAccountsResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Add new bank account
  Future<BankAccount> addBankAccount({
    required String bankCode,
    required String accountNumber,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/wallet/bank-accounts', data: {
        'bank_code': bankCode,
        'account_number': accountNumber,
      });
      return BankAccount.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Verify bank account (resolve account name)
  Future<BankAccountVerifyResponse> verifyBankAccount({
    required String bankCode,
    required String accountNumber,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/wallet/bank-accounts/verify', data: {
        'bank_code': bankCode,
        'account_number': accountNumber,
      });
      return BankAccountVerifyResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Delete bank account
  Future<void> deleteBankAccount(String accountId) async {
    try {
      await _apiClient.delete('/api/v1/wallet/bank-accounts/$accountId');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Set bank account as default
  Future<BankAccount> setDefaultBankAccount(String accountId) async {
    try {
      final response = await _apiClient.patch('/api/v1/wallet/bank-accounts/$accountId/default');
      return BankAccount.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get list of supported banks
  Future<BanksListResponse> getSupportedBanks() async {
    try {
      final response = await _apiClient.get('/api/v1/banks');
      return BanksListResponse.fromJson(response.data);
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
class WalletResponse {
  final String id;
  final double balance;
  final double availableBalance;
  final double ledgerBalance;
  final String currency;
  final bool isLocked;
  final DateTime createdAt;
  final DateTime updatedAt;

  WalletResponse({
    required this.id,
    required this.balance,
    required this.availableBalance,
    required this.ledgerBalance,
    required this.currency,
    required this.isLocked,
    required this.createdAt,
    required this.updatedAt,
  });

  factory WalletResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] ?? json;
    return WalletResponse(
      id: data['id'] ?? '',
      balance: (data['balance'] ?? 0).toDouble(),
      availableBalance: (data['available_balance'] ?? data['balance'] ?? 0).toDouble(),
      ledgerBalance: (data['ledger_balance'] ?? data['balance'] ?? 0).toDouble(),
      currency: data['currency'] ?? 'NGN',
      isLocked: data['is_locked'] ?? false,
      createdAt: DateTime.tryParse(data['created_at'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(data['updated_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class TransactionsResponse {
  final List<Transaction> transactions;
  final PaginationMeta meta;

  TransactionsResponse({required this.transactions, required this.meta});

  factory TransactionsResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return TransactionsResponse(
      transactions: data.map((t) => Transaction.fromJson(t)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class Transaction {
  final String id;
  final String type; // credit, debit
  final String category; // deposit, withdrawal, transfer_in, transfer_out, gig_payment, etc.
  final double amount;
  final double balanceBefore;
  final double balanceAfter;
  final String status; // pending, completed, failed
  final String? reference;
  final String? narration;
  final Map<String, dynamic>? metadata;
  final DateTime createdAt;

  Transaction({
    required this.id,
    required this.type,
    required this.category,
    required this.amount,
    required this.balanceBefore,
    required this.balanceAfter,
    required this.status,
    this.reference,
    this.narration,
    this.metadata,
    required this.createdAt,
  });

  factory Transaction.fromJson(Map<String, dynamic> json) {
    return Transaction(
      id: json['id'] ?? '',
      type: json['type'] ?? '',
      category: json['category'] ?? '',
      amount: (json['amount'] ?? 0).toDouble(),
      balanceBefore: (json['balance_before'] ?? 0).toDouble(),
      balanceAfter: (json['balance_after'] ?? 0).toDouble(),
      status: json['status'] ?? 'pending',
      reference: json['reference'],
      narration: json['narration'],
      metadata: json['metadata'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
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

class DepositInitResponse {
  final String reference;
  final String authorizationUrl;
  final String accessCode;

  DepositInitResponse({
    required this.reference,
    required this.authorizationUrl,
    required this.accessCode,
  });

  factory DepositInitResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] ?? json;
    return DepositInitResponse(
      reference: data['reference'] ?? '',
      authorizationUrl: data['authorization_url'] ?? '',
      accessCode: data['access_code'] ?? '',
    );
  }
}

class WithdrawalResponse {
  final bool success;
  final String message;
  final Transaction? transaction;
  final String? reference;

  WithdrawalResponse({
    required this.success,
    required this.message,
    this.transaction,
    this.reference,
  });

  factory WithdrawalResponse.fromJson(Map<String, dynamic> json) {
    return WithdrawalResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      transaction: json['data'] != null ? Transaction.fromJson(json['data']) : null,
      reference: json['reference'],
    );
  }
}

class TransferResponse {
  final bool success;
  final String message;
  final Transaction? transaction;
  final RecipientInfo? recipient;

  TransferResponse({
    required this.success,
    required this.message,
    this.transaction,
    this.recipient,
  });

  factory TransferResponse.fromJson(Map<String, dynamic> json) {
    return TransferResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      transaction: json['data'] != null ? Transaction.fromJson(json['data']) : null,
      recipient: json['recipient'] != null ? RecipientInfo.fromJson(json['recipient']) : null,
    );
  }
}

class RecipientInfo {
  final String id;
  final String name;
  final String phone;
  final String? avatar;

  RecipientInfo({
    required this.id,
    required this.name,
    required this.phone,
    this.avatar,
  });

  factory RecipientInfo.fromJson(Map<String, dynamic> json) {
    return RecipientInfo(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      phone: json['phone'] ?? '',
      avatar: json['avatar'],
    );
  }
}

class RecipientLookupResponse {
  final bool found;
  final RecipientInfo? recipient;
  final String? message;

  RecipientLookupResponse({
    required this.found,
    this.recipient,
    this.message,
  });

  factory RecipientLookupResponse.fromJson(Map<String, dynamic> json) {
    return RecipientLookupResponse(
      found: json['found'] ?? false,
      recipient: json['data'] != null ? RecipientInfo.fromJson(json['data']) : null,
      message: json['message'],
    );
  }
}

class BankAccountsResponse {
  final List<BankAccount> accounts;

  BankAccountsResponse({required this.accounts});

  factory BankAccountsResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return BankAccountsResponse(
      accounts: data.map((a) => BankAccount.fromJson(a)).toList(),
    );
  }
}

class BankAccount {
  final String id;
  final String bankCode;
  final String bankName;
  final String accountNumber;
  final String accountName;
  final bool isDefault;
  final bool isVerified;
  final DateTime createdAt;

  BankAccount({
    required this.id,
    required this.bankCode,
    required this.bankName,
    required this.accountNumber,
    required this.accountName,
    required this.isDefault,
    required this.isVerified,
    required this.createdAt,
  });

  factory BankAccount.fromJson(Map<String, dynamic> json) {
    return BankAccount(
      id: json['id'] ?? '',
      bankCode: json['bank_code'] ?? '',
      bankName: json['bank_name'] ?? '',
      accountNumber: json['account_number'] ?? '',
      accountName: json['account_name'] ?? '',
      isDefault: json['is_default'] ?? false,
      isVerified: json['is_verified'] ?? false,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class BankAccountVerifyResponse {
  final bool verified;
  final String accountName;
  final String accountNumber;
  final String bankCode;

  BankAccountVerifyResponse({
    required this.verified,
    required this.accountName,
    required this.accountNumber,
    required this.bankCode,
  });

  factory BankAccountVerifyResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] ?? json;
    return BankAccountVerifyResponse(
      verified: data['verified'] ?? true,
      accountName: data['account_name'] ?? '',
      accountNumber: data['account_number'] ?? '',
      bankCode: data['bank_code'] ?? '',
    );
  }
}

class BanksListResponse {
  final List<BankInfo> banks;

  BanksListResponse({required this.banks});

  factory BanksListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return BanksListResponse(
      banks: data.map((b) => BankInfo.fromJson(b)).toList(),
    );
  }
}

class BankInfo {
  final String code;
  final String name;
  final String? slug;
  final String? logo;

  BankInfo({
    required this.code,
    required this.name,
    this.slug,
    this.logo,
  });

  factory BankInfo.fromJson(Map<String, dynamic> json) {
    return BankInfo(
      code: json['code'] ?? '',
      name: json['name'] ?? '',
      slug: json['slug'],
      logo: json['logo'],
    );
  }
}

// Provider
final walletServiceProvider = Provider<WalletService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return WalletService(apiClient);
});
