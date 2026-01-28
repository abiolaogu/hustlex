import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/phone_number.dart';
import '../../../../core/domain/value_objects/email.dart';
import '../entities/session.dart';
import '../repositories/auth_repository.dart';

part 'register_user.freezed.dart';

/// Parameters for user registration
@freezed
class RegisterParams with _$RegisterParams {
  const factory RegisterParams({
    required PhoneNumber phone,
    required String firstName,
    required String lastName,
    Email? email,
    String? referralCode,
  }) = _RegisterParams;
}

/// Use case to register a new user
@injectable
class RegisterUser implements UseCase<Session, RegisterParams> {
  final AuthRepository _repository;

  RegisterUser(this._repository);

  @override
  Future<Either<Failure, Session>> call(RegisterParams params) async {
    // Validate phone number
    if (!params.phone.isValid()) {
      return left(
        const Failure.invalidInput('phone', 'Invalid phone number'),
      );
    }

    // Validate first name
    if (params.firstName.trim().isEmpty) {
      return left(
        const Failure.invalidInput('firstName', 'First name is required'),
      );
    }

    if (params.firstName.trim().length < 2) {
      return left(
        const Failure.invalidInput(
          'firstName',
          'First name must be at least 2 characters',
        ),
      );
    }

    // Validate last name
    if (params.lastName.trim().isEmpty) {
      return left(
        const Failure.invalidInput('lastName', 'Last name is required'),
      );
    }

    if (params.lastName.trim().length < 2) {
      return left(
        const Failure.invalidInput(
          'lastName',
          'Last name must be at least 2 characters',
        ),
      );
    }

    // Validate email if provided
    if (params.email != null && !params.email!.isValid()) {
      return left(
        const Failure.invalidInput('email', 'Invalid email address'),
      );
    }

    return _repository.register(
      phone: params.phone,
      firstName: params.firstName.trim(),
      lastName: params.lastName.trim(),
      email: params.email?.getOrCrash(),
      referralCode: params.referralCode?.trim(),
    );
  }
}
