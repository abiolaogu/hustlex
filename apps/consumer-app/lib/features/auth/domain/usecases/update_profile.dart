import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/email.dart';
import '../entities/user.dart';
import '../repositories/auth_repository.dart';

part 'update_profile.freezed.dart';

/// Parameters for updating profile
@freezed
class UpdateProfileParams with _$UpdateProfileParams {
  const factory UpdateProfileParams({
    String? firstName,
    String? lastName,
    Email? email,
    String? avatar,
  }) = _UpdateProfileParams;
}

/// Use case to update user profile
@injectable
class UpdateProfile implements UseCase<User, UpdateProfileParams> {
  final AuthRepository _repository;

  UpdateProfile(this._repository);

  @override
  Future<Either<Failure, User>> call(UpdateProfileParams params) async {
    // Validate first name if provided
    if (params.firstName != null && params.firstName!.trim().isEmpty) {
      return left(
        const Failure.invalidInput('firstName', 'First name cannot be empty'),
      );
    }

    if (params.firstName != null && params.firstName!.trim().length < 2) {
      return left(
        const Failure.invalidInput(
          'firstName',
          'First name must be at least 2 characters',
        ),
      );
    }

    // Validate last name if provided
    if (params.lastName != null && params.lastName!.trim().isEmpty) {
      return left(
        const Failure.invalidInput('lastName', 'Last name cannot be empty'),
      );
    }

    if (params.lastName != null && params.lastName!.trim().length < 2) {
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

    return _repository.updateProfile(
      firstName: params.firstName?.trim(),
      lastName: params.lastName?.trim(),
      email: params.email?.getOrCrash(),
      avatar: params.avatar,
    );
  }
}
