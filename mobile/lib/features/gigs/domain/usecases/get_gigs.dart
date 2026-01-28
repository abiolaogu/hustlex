import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/gig.dart';
import '../repositories/gigs_repository.dart';

/// Use case to get paginated list of gigs
@injectable
class GetGigs implements UseCase<PaginatedGigs, GigFilter> {
  final GigsRepository _repository;

  GetGigs(this._repository);

  @override
  Future<Either<Failure, PaginatedGigs>> call(GigFilter params) {
    return _repository.getGigs(params);
  }
}

/// Use case to get a single gig by ID
@injectable
class GetGig implements UseCase<Gig, String> {
  final GigsRepository _repository;

  GetGig(this._repository);

  @override
  Future<Either<Failure, Gig>> call(String gigId) {
    return _repository.getGig(gigId);
  }
}

/// Use case to watch a gig for real-time updates
@injectable
class WatchGig implements StreamUseCase<Gig, String> {
  final GigsRepository _repository;

  WatchGig(this._repository);

  @override
  Stream<Either<Failure, Gig>> call(String gigId) {
    return _repository.watchGig(gigId);
  }
}

/// Use case to get current user's gigs (as client)
@injectable
class GetMyGigs implements UseCase<PaginatedGigs, GetMyGigsParams> {
  final GigsRepository _repository;

  GetMyGigs(this._repository);

  @override
  Future<Either<Failure, PaginatedGigs>> call(GetMyGigsParams params) {
    return _repository.getMyGigs(
      status: params.status,
      page: params.page,
      limit: params.limit,
    );
  }
}

/// Parameters for GetMyGigs use case
class GetMyGigsParams {
  final GigStatus? status;
  final int page;
  final int limit;

  const GetMyGigsParams({
    this.status,
    this.page = 1,
    this.limit = 20,
  });
}
