import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/user.dart';
import '../entities/session.dart';
import '../repositories/auth_repository.dart';

/// Use case to get current session
@injectable
class GetCurrentSession implements UseCaseNoParams<Session> {
  final AuthRepository _repository;

  GetCurrentSession(this._repository);

  @override
  Future<Either<Failure, Session>> call() {
    return _repository.getCurrentSession();
  }
}

/// Use case to get current user
@injectable
class GetCurrentUser implements UseCaseNoParams<User> {
  final AuthRepository _repository;

  GetCurrentUser(this._repository);

  @override
  Future<Either<Failure, User>> call() {
    return _repository.getCurrentUser();
  }
}

/// Use case to refresh tokens
@injectable
class RefreshTokens implements UseCaseNoParams<AuthTokens> {
  final AuthRepository _repository;

  RefreshTokens(this._repository);

  @override
  Future<Either<Failure, AuthTokens>> call() {
    return _repository.refreshTokens();
  }
}

/// Use case to logout
@injectable
class Logout implements UseCaseNoParams<Unit> {
  final AuthRepository _repository;

  Logout(this._repository);

  @override
  Future<Either<Failure, Unit>> call() {
    return _repository.logout();
  }
}

/// Use case to check if user is authenticated
@injectable
class IsAuthenticated {
  final AuthRepository _repository;

  IsAuthenticated(this._repository);

  Future<bool> call() {
    return _repository.isAuthenticated();
  }
}
