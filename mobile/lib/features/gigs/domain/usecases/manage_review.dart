import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/review.dart';
import '../repositories/gigs_repository.dart';

part 'manage_review.freezed.dart';

/// Parameters for submitting a review
@freezed
class SubmitReviewParams with _$SubmitReviewParams {
  const factory SubmitReviewParams({
    required String contractId,
    required double rating,
    String? comment,
  }) = _SubmitReviewParams;
}

/// Use case to submit a review for a completed contract
@injectable
class SubmitReview implements UseCase<Review, SubmitReviewParams> {
  final GigsRepository _repository;

  SubmitReview(this._repository);

  @override
  Future<Either<Failure, Review>> call(SubmitReviewParams params) async {
    // Validate rating
    if (params.rating < 1.0 || params.rating > 5.0) {
      return left(
        const Failure.invalidInput('rating', 'Rating must be between 1 and 5'),
      );
    }

    // Validate comment if provided
    if (params.comment != null && params.comment!.trim().isNotEmpty) {
      if (params.comment!.trim().length < 10) {
        return left(
          const Failure.invalidInput(
            'comment',
            'Review comment must be at least 10 characters',
          ),
        );
      }
    }

    return _repository.submitReview(
      contractId: params.contractId,
      rating: params.rating,
      comment: params.comment?.trim(),
    );
  }
}

/// Use case to get reviews for a user
@injectable
class GetUserReviews implements UseCase<List<Review>, String> {
  final GigsRepository _repository;

  GetUserReviews(this._repository);

  @override
  Future<Either<Failure, List<Review>>> call(String userId) {
    return _repository.getUserReviews(userId);
  }
}

/// Use case to get review for a contract
@injectable
class GetContractReview implements UseCase<Review?, String> {
  final GigsRepository _repository;

  GetContractReview(this._repository);

  @override
  Future<Either<Failure, Review?>> call(String contractId) {
    return _repository.getContractReview(contractId);
  }
}
