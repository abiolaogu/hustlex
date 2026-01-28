import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import 'gig.dart';

part 'review.freezed.dart';

/// Review entity for gig completion feedback
@freezed
class Review with _$Review implements Entity {
  const factory Review({
    required String id,
    required String contractId,
    required String reviewerId,
    required String revieweeId,
    required double rating,
    String? comment,
    GigParticipant? reviewer,
    GigParticipant? reviewee,
    required DateTime createdAt,
  }) = _Review;

  const Review._();

  /// Check if rating is positive (4 or above)
  bool get isPositive => rating >= 4.0;

  /// Get star rating (1-5)
  int get starRating => rating.round().clamp(1, 5);
}
