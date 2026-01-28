import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../repositories/gigs_repository.dart';

/// Use case to get all available gig categories
@injectable
class GetCategories implements UseCaseNoParams<List<GigCategoryInfo>> {
  final GigsRepository _repository;

  GetCategories(this._repository);

  @override
  Future<Either<Failure, List<GigCategoryInfo>>> call() {
    return _repository.getCategories();
  }
}
