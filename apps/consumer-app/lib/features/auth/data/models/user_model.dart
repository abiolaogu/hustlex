import 'package:freezed_annotation/freezed_annotation.dart';

part 'user_model.freezed.dart';
part 'user_model.g.dart';

@freezed
class User with _$User {
  const factory User({
    required String id,
    required String phone,
    String? email,
    String? firstName,
    String? lastName,
    String? avatar,
    @Default(false) bool isVerified,
    @Default(false) bool hasSetPin,
    @Default(0) int creditScore,
    String? referralCode,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _User;

  const User._();

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);

  String get fullName {
    if (firstName == null && lastName == null) return 'User';
    return '${firstName ?? ''} ${lastName ?? ''}'.trim();
  }

  String get initials {
    final first = firstName?.isNotEmpty == true ? firstName![0] : '';
    final last = lastName?.isNotEmpty == true ? lastName![0] : '';
    if (first.isEmpty && last.isEmpty) return 'U';
    return '$first$last'.toUpperCase();
  }

  String get maskedPhone {
    if (phone.length < 8) return phone;
    return '${phone.substring(0, 4)}****${phone.substring(phone.length - 4)}';
  }
}

@freezed
class AuthTokens with _$AuthTokens {
  const factory AuthTokens({
    required String accessToken,
    required String refreshToken,
    required DateTime expiresAt,
  }) = _AuthTokens;

  const AuthTokens._();

  factory AuthTokens.fromJson(Map<String, dynamic> json) => _$AuthTokensFromJson(json);

  bool get isExpired => DateTime.now().isAfter(expiresAt);

  bool get shouldRefresh {
    final buffer = const Duration(minutes: 5);
    return DateTime.now().isAfter(expiresAt.subtract(buffer));
  }
}

@freezed
class OtpRequest with _$OtpRequest {
  const factory OtpRequest({
    required String phone,
  }) = _OtpRequest;

  factory OtpRequest.fromJson(Map<String, dynamic> json) => _$OtpRequestFromJson(json);
}

@freezed
class OtpVerify with _$OtpVerify {
  const factory OtpVerify({
    required String phone,
    required String code,
  }) = _OtpVerify;

  factory OtpVerify.fromJson(Map<String, dynamic> json) => _$OtpVerifyFromJson(json);
}

@freezed
class RegisterRequest with _$RegisterRequest {
  const factory RegisterRequest({
    required String phone,
    required String firstName,
    required String lastName,
    String? email,
    String? referralCode,
  }) = _RegisterRequest;

  factory RegisterRequest.fromJson(Map<String, dynamic> json) => _$RegisterRequestFromJson(json);
}

@freezed
class SetPinRequest with _$SetPinRequest {
  const factory SetPinRequest({
    required String pin,
  }) = _SetPinRequest;

  factory SetPinRequest.fromJson(Map<String, dynamic> json) => _$SetPinRequestFromJson(json);
}

@freezed
class AuthResponse with _$AuthResponse {
  const factory AuthResponse({
    required User user,
    required AuthTokens tokens,
    @Default(false) bool isNewUser,
  }) = _AuthResponse;

  factory AuthResponse.fromJson(Map<String, dynamic> json) => _$AuthResponseFromJson(json);
}
