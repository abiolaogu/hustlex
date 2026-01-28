import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';

part 'circle_invite.freezed.dart';

/// Circle invite entity
@freezed
class CircleInvite with _$CircleInvite implements Entity {
  const factory CircleInvite({
    required String id,
    required String circleId,
    required String circleName,
    required String inviterId,
    required String inviterName,
    required String inviteePhone,
    @Default(false) bool isAccepted,
    @Default(false) bool isExpired,
    DateTime? expiresAt,
    required DateTime createdAt,
  }) = _CircleInvite;

  const CircleInvite._();

  /// Check if invite is still valid
  bool get isValid {
    if (isAccepted || isExpired) return false;
    if (expiresAt == null) return true;
    return DateTime.now().isBefore(expiresAt!);
  }

  /// Days until expiry
  int get daysUntilExpiry {
    if (expiresAt == null) return -1;
    return expiresAt!.difference(DateTime.now()).inDays;
  }
}
