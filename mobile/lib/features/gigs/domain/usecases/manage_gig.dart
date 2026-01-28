import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../entities/gig.dart';
import '../repositories/gigs_repository.dart';

part 'manage_gig.freezed.dart';

/// Parameters for creating a gig
@freezed
class CreateGigParams with _$CreateGigParams {
  const factory CreateGigParams({
    required String title,
    required String description,
    required GigCategory category,
    required Money budgetMin,
    required Money budgetMax,
    required int durationDays,
    @Default(true) bool isRemote,
    String? location,
    DateTime? deadline,
    @Default([]) List<String> skills,
    @Default([]) List<String> attachments,
  }) = _CreateGigParams;
}

/// Use case to create a new gig
@injectable
class CreateGig implements UseCase<Gig, CreateGigParams> {
  final GigsRepository _repository;

  CreateGig(this._repository);

  @override
  Future<Either<Failure, Gig>> call(CreateGigParams params) async {
    // Validate title
    if (params.title.trim().isEmpty) {
      return left(
        const Failure.invalidInput('title', 'Title cannot be empty'),
      );
    }
    if (params.title.trim().length < 10) {
      return left(
        const Failure.invalidInput(
          'title',
          'Title must be at least 10 characters',
        ),
      );
    }

    // Validate description
    if (params.description.trim().isEmpty) {
      return left(
        const Failure.invalidInput('description', 'Description cannot be empty'),
      );
    }
    if (params.description.trim().length < 50) {
      return left(
        const Failure.invalidInput(
          'description',
          'Description must be at least 50 characters',
        ),
      );
    }

    // Validate budget
    if (params.budgetMin.amount <= 0) {
      return left(
        const Failure.invalidInput('budgetMin', 'Minimum budget must be greater than 0'),
      );
    }
    if (params.budgetMax.amount < params.budgetMin.amount) {
      return left(
        const Failure.invalidInput(
          'budgetMax',
          'Maximum budget cannot be less than minimum budget',
        ),
      );
    }

    // Validate duration
    if (params.durationDays <= 0) {
      return left(
        const Failure.invalidInput('durationDays', 'Duration must be at least 1 day'),
      );
    }

    // Validate skills
    if (params.skills.isEmpty) {
      return left(
        const Failure.invalidInput('skills', 'At least one skill is required'),
      );
    }

    // Validate location if not remote
    if (!params.isRemote && (params.location == null || params.location!.isEmpty)) {
      return left(
        const Failure.invalidInput('location', 'Location is required for non-remote gigs'),
      );
    }

    return _repository.createGig(
      title: params.title.trim(),
      description: params.description.trim(),
      category: params.category,
      budgetMin: params.budgetMin,
      budgetMax: params.budgetMax,
      durationDays: params.durationDays,
      isRemote: params.isRemote,
      location: params.location,
      deadline: params.deadline,
      skills: params.skills,
      attachments: params.attachments,
    );
  }
}

/// Parameters for updating a gig
@freezed
class UpdateGigParams with _$UpdateGigParams {
  const factory UpdateGigParams({
    required String gigId,
    String? title,
    String? description,
    GigCategory? category,
    Money? budgetMin,
    Money? budgetMax,
    int? durationDays,
    bool? isRemote,
    String? location,
    DateTime? deadline,
    List<String>? skills,
    List<String>? attachments,
  }) = _UpdateGigParams;
}

/// Use case to update an existing gig
@injectable
class UpdateGig implements UseCase<Gig, UpdateGigParams> {
  final GigsRepository _repository;

  UpdateGig(this._repository);

  @override
  Future<Either<Failure, Gig>> call(UpdateGigParams params) async {
    // Validate title if provided
    if (params.title != null && params.title!.trim().isEmpty) {
      return left(
        const Failure.invalidInput('title', 'Title cannot be empty'),
      );
    }
    if (params.title != null && params.title!.trim().length < 10) {
      return left(
        const Failure.invalidInput(
          'title',
          'Title must be at least 10 characters',
        ),
      );
    }

    // Validate description if provided
    if (params.description != null && params.description!.trim().isEmpty) {
      return left(
        const Failure.invalidInput('description', 'Description cannot be empty'),
      );
    }
    if (params.description != null && params.description!.trim().length < 50) {
      return left(
        const Failure.invalidInput(
          'description',
          'Description must be at least 50 characters',
        ),
      );
    }

    // Validate budget if provided
    if (params.budgetMin != null && params.budgetMin!.amount <= 0) {
      return left(
        const Failure.invalidInput('budgetMin', 'Minimum budget must be greater than 0'),
      );
    }

    // Validate duration if provided
    if (params.durationDays != null && params.durationDays! <= 0) {
      return left(
        const Failure.invalidInput('durationDays', 'Duration must be at least 1 day'),
      );
    }

    return _repository.updateGig(
      gigId: params.gigId,
      title: params.title?.trim(),
      description: params.description?.trim(),
      category: params.category,
      budgetMin: params.budgetMin,
      budgetMax: params.budgetMax,
      durationDays: params.durationDays,
      isRemote: params.isRemote,
      location: params.location,
      deadline: params.deadline,
      skills: params.skills,
      attachments: params.attachments,
    );
  }
}

/// Use case to publish a draft gig
@injectable
class PublishGig implements UseCase<Gig, String> {
  final GigsRepository _repository;

  PublishGig(this._repository);

  @override
  Future<Either<Failure, Gig>> call(String gigId) {
    return _repository.publishGig(gigId);
  }
}

/// Use case to cancel a gig
@injectable
class CancelGig implements UseCase<Unit, String> {
  final GigsRepository _repository;

  CancelGig(this._repository);

  @override
  Future<Either<Failure, Unit>> call(String gigId) {
    return _repository.cancelGig(gigId);
  }
}
