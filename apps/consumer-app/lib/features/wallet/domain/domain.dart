/// Wallet domain layer barrel file
library;

// Entities
export 'entities/wallet.dart';
export 'entities/transaction.dart';
export 'entities/bank_account.dart';
export 'entities/deposit.dart';

// Repository interface
export 'repositories/wallet_repository.dart';

// Use cases
export 'usecases/get_wallet.dart';
export 'usecases/get_transactions.dart';
export 'usecases/transfer_funds.dart';
export 'usecases/initiate_deposit.dart';
export 'usecases/initiate_withdrawal.dart';
export 'usecases/manage_bank_accounts.dart';
