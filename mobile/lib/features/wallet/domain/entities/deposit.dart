import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'deposit.freezed.dart';

/// Deposit initiation result from payment gateway
@freezed
class DepositInitiation with _$DepositInitiation {
  const factory DepositInitiation({
    required String reference,
    required String authorizationUrl,
    required String accessCode,
    required Money amount,
  }) = _DepositInitiation;
}

/// Deposit verification result
@freezed
class DepositVerification with _$DepositVerification {
  const factory DepositVerification({
    required String reference,
    required bool success,
    required Money amount,
    String? message,
  }) = _DepositVerification;
}
