import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/di/providers.dart';
import '../../data/models/wallet_model.dart';
import '../../data/repositories/wallet_repository.dart';

/// Wallet state
class WalletState {
  final Wallet? wallet;
  final WalletStats? stats;
  final List<BankAccount> bankAccounts;
  final bool isLoading;
  final String? error;

  const WalletState({
    this.wallet,
    this.stats,
    this.bankAccounts = const [],
    this.isLoading = false,
    this.error,
  });

  double get availableBalance => wallet?.availableBalance ?? 0;
  double get escrowBalance => wallet?.escrowBalance ?? 0;
  double get savingsBalance => wallet?.savingsBalance ?? 0;
  double get totalBalance => wallet?.totalBalance ?? 0;

  WalletState copyWith({
    Wallet? wallet,
    WalletStats? stats,
    List<BankAccount>? bankAccounts,
    bool? isLoading,
    String? error,
  }) {
    return WalletState(
      wallet: wallet ?? this.wallet,
      stats: stats ?? this.stats,
      bankAccounts: bankAccounts ?? this.bankAccounts,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// Wallet state notifier
class WalletNotifier extends StateNotifier<WalletState> {
  final WalletRepository _repository;

  WalletNotifier(this._repository) : super(const WalletState());

  /// Load wallet data
  Future<void> loadWallet() async {
    state = state.copyWith(isLoading: true, error: null);

    final result = await _repository.getWallet();

    result.when(
      success: (wallet) {
        state = state.copyWith(wallet: wallet, isLoading: false);
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
      },
    );
  }

  /// Load wallet statistics
  Future<void> loadStats({DateTime? startDate, DateTime? endDate}) async {
    final result = await _repository.getWalletStats(
      startDate: startDate,
      endDate: endDate,
    );

    result.when(
      success: (stats) {
        state = state.copyWith(stats: stats);
      },
      failure: (_, __) {},
    );
  }

  /// Load bank accounts
  Future<void> loadBankAccounts() async {
    final result = await _repository.getBankAccounts();

    result.when(
      success: (accounts) {
        state = state.copyWith(bankAccounts: accounts);
      },
      failure: (_, __) {},
    );
  }

  /// Load all wallet data
  Future<void> loadAll() async {
    state = state.copyWith(isLoading: true, error: null);

    await Future.wait([
      loadWallet(),
      loadStats(),
      loadBankAccounts(),
    ]);

    state = state.copyWith(isLoading: false);
  }

  /// Refresh wallet balance only
  Future<void> refreshBalance() async {
    final result = await _repository.getBalance();

    result.when(
      success: (balance) {
        if (state.wallet != null) {
          state = state.copyWith(
            wallet: state.wallet!.copyWith(availableBalance: balance),
          );
        }
      },
      failure: (_, __) {},
    );
  }

  /// Add bank account
  Future<bool> addBankAccount({
    required String bankCode,
    required String accountNumber,
  }) async {
    final result = await _repository.addBankAccount(
      bankCode: bankCode,
      accountNumber: accountNumber,
    );

    return result.when(
      success: (account) {
        state = state.copyWith(
          bankAccounts: [...state.bankAccounts, account],
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Remove bank account
  Future<bool> removeBankAccount(String accountId) async {
    final result = await _repository.removeBankAccount(accountId);

    return result.when(
      success: (_) {
        state = state.copyWith(
          bankAccounts: state.bankAccounts
              .where((a) => a.id != accountId)
              .toList(),
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Set default bank account
  Future<bool> setDefaultBankAccount(String accountId) async {
    final result = await _repository.setDefaultBankAccount(accountId);

    return result.when(
      success: (_) {
        state = state.copyWith(
          bankAccounts: state.bankAccounts.map((a) {
            return a.copyWith(isDefault: a.id == accountId);
          }).toList(),
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }
}

/// Main wallet provider
final walletProvider = StateNotifierProvider<WalletNotifier, WalletState>((ref) {
  final repository = ref.watch(walletRepositoryProvider);
  return WalletNotifier(repository);
});

/// Wallet balance only provider (for widgets)
final walletBalanceProvider = Provider<double>((ref) {
  return ref.watch(walletProvider).availableBalance;
});

/// Transactions state
class TransactionsState {
  final List<Transaction> transactions;
  final bool isLoading;
  final bool hasMore;
  final int currentPage;
  final TransactionFilter filter;
  final String? error;

  const TransactionsState({
    this.transactions = const [],
    this.isLoading = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.filter = const TransactionFilter(),
    this.error,
  });

  TransactionsState copyWith({
    List<Transaction>? transactions,
    bool? isLoading,
    bool? hasMore,
    int? currentPage,
    TransactionFilter? filter,
    String? error,
  }) {
    return TransactionsState(
      transactions: transactions ?? this.transactions,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      filter: filter ?? this.filter,
      error: error,
    );
  }
}

/// Transactions notifier with pagination
class TransactionsNotifier extends StateNotifier<TransactionsState> {
  final WalletRepository _repository;

  TransactionsNotifier(this._repository) : super(const TransactionsState());

  /// Load transactions with pagination
  Future<void> loadTransactions({bool refresh = false}) async {
    if (state.isLoading) return;
    if (!refresh && !state.hasMore) return;

    state = state.copyWith(
      isLoading: true,
      error: null,
      currentPage: refresh ? 1 : state.currentPage,
    );

    final filter = state.filter.copyWith(page: refresh ? 1 : state.currentPage);
    final result = await _repository.getTransactions(filter);

    result.when(
      success: (data) {
        final newTransactions = refresh
            ? data.transactions
            : [...state.transactions, ...data.transactions];
        state = state.copyWith(
          transactions: newTransactions,
          isLoading: false,
          hasMore: data.hasMore,
          currentPage: data.page + 1,
        );
      },
      failure: (message, _) {
        state = state.copyWith(
          isLoading: false,
          error: message,
        );
      },
    );
  }

  /// Refresh transactions
  Future<void> refresh() => loadTransactions(refresh: true);

  /// Update filter and reload
  void updateFilter(TransactionFilter filter) {
    state = state.copyWith(filter: filter);
    loadTransactions(refresh: true);
  }

  /// Filter by type
  void filterByType(TransactionType? type) {
    updateFilter(state.filter.copyWith(type: type));
  }

  /// Filter by status
  void filterByStatus(TransactionStatus? status) {
    updateFilter(state.filter.copyWith(status: status));
  }

  /// Filter by date range
  void filterByDateRange(DateTime? start, DateTime? end) {
    updateFilter(state.filter.copyWith(startDate: start, endDate: end));
  }

  /// Search transactions
  void search(String? query) {
    updateFilter(state.filter.copyWith(search: query));
  }

  /// Clear all filters
  void clearFilters() {
    updateFilter(const TransactionFilter());
  }
}

/// Transactions provider
final transactionsProvider =
    StateNotifierProvider<TransactionsNotifier, TransactionsState>((ref) {
  final repository = ref.watch(walletRepositoryProvider);
  return TransactionsNotifier(repository);
});

/// Recent transactions provider (for home screen)
final recentTransactionsProvider = FutureProvider<List<Transaction>>((ref) async {
  final repository = ref.watch(walletRepositoryProvider);
  final result = await repository.getRecentTransactions(limit: 5);
  return result.data ?? [];
});

/// Transaction detail provider
final transactionDetailProvider =
    FutureProvider.family<Transaction?, String>((ref, transactionId) async {
  final repository = ref.watch(walletRepositoryProvider);
  final result = await repository.getTransaction(transactionId);
  return result.data;
});

/// Bank account verification state
class BankVerificationState {
  final bool isVerifying;
  final AccountVerification? verification;
  final String? error;

  const BankVerificationState({
    this.isVerifying = false,
    this.verification,
    this.error,
  });

  BankVerificationState copyWith({
    bool? isVerifying,
    AccountVerification? verification,
    String? error,
  }) {
    return BankVerificationState(
      isVerifying: isVerifying ?? this.isVerifying,
      verification: verification ?? this.verification,
      error: error,
    );
  }
}

/// Bank verification notifier
class BankVerificationNotifier extends StateNotifier<BankVerificationState> {
  final WalletRepository _repository;

  BankVerificationNotifier(this._repository)
      : super(const BankVerificationState());

  /// Verify bank account (name enquiry)
  Future<bool> verifyAccount({
    required String bankCode,
    required String accountNumber,
  }) async {
    state = state.copyWith(isVerifying: true, error: null, verification: null);

    final result = await _repository.verifyBankAccount(
      bankCode: bankCode,
      accountNumber: accountNumber,
    );

    return result.when(
      success: (verification) {
        state = state.copyWith(
          isVerifying: false,
          verification: verification,
        );
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(
          isVerifying: false,
          error: message,
        );
        return false;
      },
    );
  }

  /// Clear verification
  void clear() {
    state = const BankVerificationState();
  }
}

/// Bank verification provider
final bankVerificationProvider =
    StateNotifierProvider<BankVerificationNotifier, BankVerificationState>((ref) {
  final repository = ref.watch(walletRepositoryProvider);
  return BankVerificationNotifier(repository);
});

/// Banks list provider
final banksProvider = FutureProvider<List<Bank>>((ref) async {
  final repository = ref.watch(walletRepositoryProvider);
  final result = await repository.getBanks();
  return result.data ?? [];
});

/// Networks provider (for airtime)
final networksProvider = FutureProvider<List<String>>((ref) async {
  final repository = ref.watch(walletRepositoryProvider);
  final result = await repository.getNetworks();
  return result.data ?? [];
});

/// Data plans provider
final dataPlansProvider =
    FutureProvider.family<List<Map<String, dynamic>>, String>((ref, network) async {
  final repository = ref.watch(walletRepositoryProvider);
  final result = await repository.getDataPlans(network);
  return result.data ?? [];
});

/// Deposit state
class DepositState {
  final bool isProcessing;
  final PaystackInitResponse? initResponse;
  final Transaction? transaction;
  final String? error;

  const DepositState({
    this.isProcessing = false,
    this.initResponse,
    this.transaction,
    this.error,
  });

  DepositState copyWith({
    bool? isProcessing,
    PaystackInitResponse? initResponse,
    Transaction? transaction,
    String? error,
  }) {
    return DepositState(
      isProcessing: isProcessing ?? this.isProcessing,
      initResponse: initResponse ?? this.initResponse,
      transaction: transaction ?? this.transaction,
      error: error,
    );
  }
}

/// Deposit notifier
class DepositNotifier extends StateNotifier<DepositState> {
  final WalletRepository _repository;

  DepositNotifier(this._repository) : super(const DepositState());

  /// Initialize deposit
  Future<bool> initializeDeposit(DepositRequest request) async {
    state = state.copyWith(isProcessing: true, error: null);

    final result = await _repository.initializeDeposit(request);

    return result.when(
      success: (response) {
        state = state.copyWith(
          isProcessing: false,
          initResponse: response,
        );
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(
          isProcessing: false,
          error: message,
        );
        return false;
      },
    );
  }

  /// Verify deposit after payment
  Future<bool> verifyDeposit(String reference) async {
    state = state.copyWith(isProcessing: true, error: null);

    final result = await _repository.verifyDeposit(reference);

    return result.when(
      success: (transaction) {
        state = state.copyWith(
          isProcessing: false,
          transaction: transaction,
        );
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(
          isProcessing: false,
          error: message,
        );
        return false;
      },
    );
  }

  /// Reset state
  void reset() {
    state = const DepositState();
  }
}

/// Deposit provider
final depositProvider = StateNotifierProvider<DepositNotifier, DepositState>((ref) {
  final repository = ref.watch(walletRepositoryProvider);
  return DepositNotifier(repository);
});

/// Transfer state
class TransferState {
  final bool isProcessing;
  final Transaction? transaction;
  final String? error;

  const TransferState({
    this.isProcessing = false,
    this.transaction,
    this.error,
  });

  TransferState copyWith({
    bool? isProcessing,
    Transaction? transaction,
    String? error,
  }) {
    return TransferState(
      isProcessing: isProcessing ?? this.isProcessing,
      transaction: transaction ?? this.transaction,
      error: error,
    );
  }
}

/// Transfer notifier
class TransferNotifier extends StateNotifier<TransferState> {
  final WalletRepository _repository;

  TransferNotifier(this._repository) : super(const TransferState());

  /// Transfer to HustleX user
  Future<bool> transfer(TransferRequest request) async {
    state = state.copyWith(isProcessing: true, error: null);

    final result = await _repository.transfer(request);

    return result.when(
      success: (transaction) {
        state = state.copyWith(
          isProcessing: false,
          transaction: transaction,
        );
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(
          isProcessing: false,
          error: message,
        );
        return false;
      },
    );
  }

  /// Transfer to bank account
  Future<bool> bankTransfer(BankTransferRequest request) async {
    state = state.copyWith(isProcessing: true, error: null);

    final result = await _repository.bankTransfer(request);

    return result.when(
      success: (transaction) {
        state = state.copyWith(
          isProcessing: false,
          transaction: transaction,
        );
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(
          isProcessing: false,
          error: message,
        );
        return false;
      },
    );
  }

  /// Withdraw to saved bank account
  Future<bool> withdraw(WithdrawalRequest request) async {
    state = state.copyWith(isProcessing: true, error: null);

    final result = await _repository.withdraw(request);

    return result.when(
      success: (transaction) {
        state = state.copyWith(
          isProcessing: false,
          transaction: transaction,
        );
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(
          isProcessing: false,
          error: message,
        );
        return false;
      },
    );
  }

  /// Reset state
  void reset() {
    state = const TransferState();
  }
}

/// Transfer provider
final transferProvider = StateNotifierProvider<TransferNotifier, TransferState>((ref) {
  final repository = ref.watch(walletRepositoryProvider);
  return TransferNotifier(repository);
});
