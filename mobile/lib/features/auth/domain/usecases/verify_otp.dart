import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/phone_number.dart';
import '../entities/session.dart';
import '../repositories/auth_repository.dart';

part 'verify_otp.freezed.dart';

/// Parameters for verifying OTP
@freezed
class VerifyOtpParams with _$VerifyOtpParams {
  const factory VerifyOtpParams({
    required PhoneNumber phone,
    required String code,
  }) = _VerifyOtpParams;
}

/// Use case to verify OTP code
@injectable
class VerifyOtp implements UseCase<Session, VerifyOtpParams> {
  final AuthRepository _repository;

  VerifyOtp(this._repository);

  @override
  Future<Either<Failure, Session>> call(VerifyOtpParams params) async {
    // Validate phone number
    if (!params.phone.isValid()) {
      return left(
        const Failure.invalidInput('phone', 'Invalid phone number'),
      );
    }

    // Validate OTP code format (6 digits)
    if (params.code.isEmpty) {
      return left(
        const Failure.invalidInput('code', 'Please enter the OTP code'),
      );
    }

    if (params.code.length != 6 || !RegExp(r'^\d{6}$').hasMatch(params.code)) {
      return left(
        const Failure.invalidInput('code', 'OTP must be 6 digits'),
      );
    }

    return _repository.verifyOtp(
      phone: params.phone,
      code: params.code,
    );
  }
}
