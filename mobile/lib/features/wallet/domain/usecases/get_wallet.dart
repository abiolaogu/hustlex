import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/wallet.dart';
import '../repositories/wallet_repository.dart';

/// Use case to get the current user's wallet
@injectable
class GetWallet implements UseCaseNoParams<Wallet> {
  final WalletRepository _repository;

  GetWallet(this._repository);

  @override
  Future<Either<Failure, Wallet>> call() {
    return _repository.getWallet();
  }
}

/// Use case to watch wallet changes in real-time
@injectable
class WatchWallet implements StreamUseCaseNoParams<Wallet> {
  final WalletRepository _repository;

  WatchWallet(this._repository);

  @override
  Stream<Either<Failure, Wallet>> call() {
    return _repository.watchWallet();
  }
}
