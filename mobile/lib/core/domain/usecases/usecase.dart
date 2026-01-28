import 'package:dartz/dartz.dart';
import '../failures/failure.dart';

/// Base class for use cases that take parameters.
///
/// Use cases encapsulate a single business operation.
/// They receive parameters and return either a [Failure] or a result [Type].
///
/// Example:
/// ```dart
/// class GetUser implements UseCase<User, GetUserParams> {
///   final UserRepository repository;
///
///   GetUser(this.repository);
///
///   @override
///   Future<Either<Failure, User>> call(GetUserParams params) {
///     return repository.getUser(params.userId);
///   }
/// }
/// ```
abstract class UseCase<Type, Params> {
  Future<Either<Failure, Type>> call(Params params);
}

/// Base class for use cases without parameters.
///
/// Example:
/// ```dart
/// class GetCurrentUser implements UseCaseNoParams<User> {
///   final AuthRepository repository;
///
///   GetCurrentUser(this.repository);
///
///   @override
///   Future<Either<Failure, User>> call() {
///     return repository.getCurrentUser();
///   }
/// }
/// ```
abstract class UseCaseNoParams<Type> {
  Future<Either<Failure, Type>> call();
}

/// Base class for stream-based use cases.
///
/// Use for operations that emit multiple values over time.
///
/// Example:
/// ```dart
/// class WatchWalletBalance implements StreamUseCase<Money, String> {
///   final WalletRepository repository;
///
///   WatchWalletBalance(this.repository);
///
///   @override
///   Stream<Either<Failure, Money>> call(String walletId) {
///     return repository.watchBalance(walletId);
///   }
/// }
/// ```
abstract class StreamUseCase<Type, Params> {
  Stream<Either<Failure, Type>> call(Params params);
}

/// Stream use case without parameters.
abstract class StreamUseCaseNoParams<Type> {
  Stream<Either<Failure, Type>> call();
}

/// Empty params class for use cases that don't need parameters.
///
/// Use this when you want to keep consistency with [UseCase] but
/// don't need any parameters.
class NoParams {
  const NoParams();
}

/// Unit type for use cases that don't return a value.
///
/// Example: DeleteUser returns Either<Failure, Unit>
typedef Unit = void;
const unit = null;
