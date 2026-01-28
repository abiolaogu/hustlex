import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/base_entity.dart';
import '../../../../core/domain/value_objects/email.dart';
import '../../../../core/domain/value_objects/phone_number.dart';

part 'user.freezed.dart';

/// User domain entity
@freezed
class User with _$User implements Entity {
  const User._();

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

  /// Get user's full name
  String get fullName {
    if (firstName == null && lastName == null) return 'User';
    return '${firstName ?? ''} ${lastName ?? ''}'.trim();
  }

  /// Get user's initials for avatar placeholder
  String get initials {
    final first = firstName?.isNotEmpty == true ? firstName![0] : '';
    final last = lastName?.isNotEmpty == true ? lastName![0] : '';
    if (first.isEmpty && last.isEmpty) return 'U';
    return '$first$last'.toUpperCase();
  }

  /// Get masked phone for display
  String get maskedPhone {
    if (phone.length < 8) return phone;
    return '${phone.substring(0, 4)}****${phone.substring(phone.length - 4)}';
  }

  /// Get PhoneNumber value object
  PhoneNumber get phoneNumber => PhoneNumber(phone);

  /// Get Email value object (if email exists)
  Email? get emailVO => email != null ? Email(email!) : null;

  /// Check if profile is complete
  bool get isProfileComplete =>
      firstName != null &&
      lastName != null &&
      firstName!.isNotEmpty &&
      lastName!.isNotEmpty;

  /// Check if user can perform financial transactions
  bool get canTransact => isVerified && hasSetPin;
}
