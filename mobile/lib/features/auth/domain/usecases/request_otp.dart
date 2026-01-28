import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/phone_number.dart';
import '../entities/session.dart';
import '../repositories/auth_repository.dart';

part 'request_otp.freezed.dart';

/// Parameters for requesting OTP
@freezed
class RequestOtpParams with _$RequestOtpParams {
  const factory RequestOtpParams({
    required PhoneNumber phone,
    @Default(OtpPurpose.login) OtpPurpose purpose,
  }) = _RequestOtpParams;
}

/// Use case to request OTP for phone verification
@injectable
class RequestOtp implements UseCase<OtpSession, RequestOtpParams> {
  final AuthRepository _repository;

  RequestOtp(this._repository);

  @override
  Future<Either<Failure, OtpSession>> call(RequestOtpParams params) async {
    // Validate phone number
    if (!params.phone.isValid()) {
      return left(
        const Failure.invalidInput('phone', 'Please enter a valid phone number'),
      );
    }

    return _repository.requestOtp(
      phone: params.phone,
      purpose: params.purpose,
    );
  }
}
