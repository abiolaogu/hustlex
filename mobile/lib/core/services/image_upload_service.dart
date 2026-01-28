import 'dart:io';

import 'package:dio/dio.dart';
import 'package:flutter_image_compress/flutter_image_compress.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';
import 'package:path/path.dart' as path;
import 'package:path_provider/path_provider.dart';

import '../api/api_client.dart';
import '../../core/di/providers.dart';

/// Image source for picking
enum ImageSource {
  camera,
  gallery,
}

/// Compression quality presets
enum ImageQuality {
  low,     // 30% quality, max 500px
  medium,  // 50% quality, max 800px
  high,    // 70% quality, max 1200px
  original, // No compression
}

extension ImageQualityX on ImageQuality {
  int get quality {
    switch (this) {
      case ImageQuality.low:
        return 30;
      case ImageQuality.medium:
        return 50;
      case ImageQuality.high:
        return 70;
      case ImageQuality.original:
        return 100;
    }
  }

  int get maxSize {
    switch (this) {
      case ImageQuality.low:
        return 500;
      case ImageQuality.medium:
        return 800;
      case ImageQuality.high:
        return 1200;
      case ImageQuality.original:
        return 4096; // No practical limit
    }
  }
}

/// Upload progress callback
typedef UploadProgressCallback = void Function(int sent, int total);

/// Upload result
class UploadResult {
  final bool success;
  final String? url;
  final String? id;
  final String? error;
  final Map<String, dynamic>? metadata;

  const UploadResult({
    required this.success,
    this.url,
    this.id,
    this.error,
    this.metadata,
  });

  factory UploadResult.success({
    required String url,
    String? id,
    Map<String, dynamic>? metadata,
  }) {
    return UploadResult(
      success: true,
      url: url,
      id: id,
      metadata: metadata,
    );
  }

  factory UploadResult.failed({String? error}) {
    return UploadResult(
      success: false,
      error: error ?? 'Upload failed',
    );
  }
}

/// Image upload service
class ImageUploadService {
  final ApiClient _apiClient;
  final ImagePicker _picker;

  ImageUploadService({
    required ApiClient apiClient,
    ImagePicker? picker,
  })  : _apiClient = apiClient,
        _picker = picker ?? ImagePicker();

  /// Pick image from source
  Future<File?> pickImage({
    ImageSource source = ImageSource.gallery,
    double? maxWidth,
    double? maxHeight,
    int? imageQuality,
  }) async {
    try {
      final XFile? pickedFile = await _picker.pickImage(
        source: source == ImageSource.camera 
            ? ImageSource.camera as image_picker.ImageSource
            : ImageSource.gallery as image_picker.ImageSource,
        maxWidth: maxWidth,
        maxHeight: maxHeight,
        imageQuality: imageQuality,
      );

      if (pickedFile != null) {
        return File(pickedFile.path);
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  /// Pick multiple images from gallery
  Future<List<File>> pickMultipleImages({
    int? maxImages,
    double? maxWidth,
    double? maxHeight,
    int? imageQuality,
  }) async {
    try {
      final List<XFile> pickedFiles = await _picker.pickMultiImage(
        maxWidth: maxWidth,
        maxHeight: maxHeight,
        imageQuality: imageQuality,
        limit: maxImages,
      );

      return pickedFiles.map((xFile) => File(xFile.path)).toList();
    } catch (e) {
      return [];
    }
  }

  /// Compress image
  Future<File?> compressImage(
    File file, {
    ImageQuality quality = ImageQuality.medium,
    String? targetPath,
  }) async {
    try {
      final dir = await getTemporaryDirectory();
      final targetFilePath = targetPath ??
          path.join(
            dir.path,
            '${DateTime.now().millisecondsSinceEpoch}_compressed.jpg',
          );

      final result = await FlutterImageCompress.compressAndGetFile(
        file.absolute.path,
        targetFilePath,
        quality: quality.quality,
        minWidth: quality.maxSize,
        minHeight: quality.maxSize,
        format: CompressFormat.jpeg,
      );

      return result != null ? File(result.path) : null;
    } catch (e) {
      return null;
    }
  }

  /// Pick and compress image
  Future<File?> pickAndCompressImage({
    ImageSource source = ImageSource.gallery,
    ImageQuality quality = ImageQuality.medium,
  }) async {
    final file = await pickImage(source: source);
    if (file == null) return null;

    if (quality == ImageQuality.original) {
      return file;
    }

    return await compressImage(file, quality: quality);
  }

  /// Upload profile photo
  Future<UploadResult> uploadProfilePhoto(
    File file, {
    UploadProgressCallback? onProgress,
    ImageQuality quality = ImageQuality.medium,
  }) async {
    // Compress before upload
    final compressedFile = quality != ImageQuality.original
        ? await compressImage(file, quality: quality) ?? file
        : file;

    return _uploadFile(
      file: compressedFile,
      endpoint: '/profile/photo',
      fieldName: 'photo',
      onProgress: onProgress,
    );
  }

  /// Upload KYC document
  Future<UploadResult> uploadKycDocument(
    File file, {
    required String documentType, // 'id_card', 'proof_of_address', 'selfie'
    UploadProgressCallback? onProgress,
    ImageQuality quality = ImageQuality.high,
  }) async {
    final compressedFile = quality != ImageQuality.original
        ? await compressImage(file, quality: quality) ?? file
        : file;

    return _uploadFile(
      file: compressedFile,
      endpoint: '/profile/kyc/documents',
      fieldName: documentType,
      additionalFields: {'document_type': documentType},
      onProgress: onProgress,
    );
  }

  /// Upload gig attachment
  Future<UploadResult> uploadGigAttachment(
    File file, {
    String? gigId,
    UploadProgressCallback? onProgress,
    ImageQuality quality = ImageQuality.medium,
  }) async {
    final compressedFile = quality != ImageQuality.original
        ? await compressImage(file, quality: quality) ?? file
        : file;

    return _uploadFile(
      file: compressedFile,
      endpoint: '/gigs/attachments',
      fieldName: 'attachment',
      additionalFields: gigId != null ? {'gig_id': gigId} : null,
      onProgress: onProgress,
    );
  }

  /// Upload chat/message attachment
  Future<UploadResult> uploadMessageAttachment(
    File file, {
    required String conversationId,
    UploadProgressCallback? onProgress,
    ImageQuality quality = ImageQuality.medium,
  }) async {
    final compressedFile = quality != ImageQuality.original
        ? await compressImage(file, quality: quality) ?? file
        : file;

    return _uploadFile(
      file: compressedFile,
      endpoint: '/messages/attachments',
      fieldName: 'attachment',
      additionalFields: {'conversation_id': conversationId},
      onProgress: onProgress,
    );
  }

  /// Upload feedback attachment
  Future<UploadResult> uploadFeedbackAttachment(
    File file, {
    UploadProgressCallback? onProgress,
    ImageQuality quality = ImageQuality.medium,
  }) async {
    final compressedFile = quality != ImageQuality.original
        ? await compressImage(file, quality: quality) ?? file
        : file;

    return _uploadFile(
      file: compressedFile,
      endpoint: '/feedback/attachments',
      fieldName: 'attachment',
      onProgress: onProgress,
    );
  }

  /// Generic file upload
  Future<UploadResult> _uploadFile({
    required File file,
    required String endpoint,
    required String fieldName,
    Map<String, dynamic>? additionalFields,
    UploadProgressCallback? onProgress,
  }) async {
    try {
      final fileName = path.basename(file.path);
      
      final formData = FormData.fromMap({
        fieldName: await MultipartFile.fromFile(
          file.path,
          filename: fileName,
        ),
        if (additionalFields != null) ...additionalFields,
      });

      final response = await _apiClient.dio.post(
        endpoint,
        data: formData,
        onSendProgress: onProgress,
        options: Options(
          headers: {'Content-Type': 'multipart/form-data'},
        ),
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = response.data['data'];
        return UploadResult.success(
          url: data['url'] ?? data['photo_url'] ?? data['attachment_url'],
          id: data['id'],
          metadata: data is Map<String, dynamic> ? data : null,
        );
      } else {
        return UploadResult.failed(
          error: response.data['message'] ?? 'Upload failed',
        );
      }
    } on DioException catch (e) {
      return UploadResult.failed(
        error: e.response?.data?['message'] ?? e.message ?? 'Upload failed',
      );
    } catch (e) {
      return UploadResult.failed(error: e.toString());
    }
  }

  /// Get file size in MB
  double getFileSizeMB(File file) {
    final bytes = file.lengthSync();
    return bytes / (1024 * 1024);
  }

  /// Check if file is within size limit
  bool isFileSizeValid(File file, {double maxSizeMB = 5.0}) {
    return getFileSizeMB(file) <= maxSizeMB;
  }

  /// Get file extension
  String getFileExtension(File file) {
    return path.extension(file.path).toLowerCase();
  }

  /// Check if file is a valid image
  bool isValidImage(File file) {
    final ext = getFileExtension(file);
    return ['.jpg', '.jpeg', '.png', '.gif', '.webp'].contains(ext);
  }

  /// Delete temp files
  Future<void> cleanupTempFiles() async {
    try {
      final dir = await getTemporaryDirectory();
      final files = dir.listSync();
      
      for (final file in files) {
        if (file is File && file.path.contains('_compressed')) {
          await file.delete();
        }
      }
    } catch (e) {
      // Ignore cleanup errors
    }
  }
}

/// Image upload service provider
final imageUploadServiceProvider = Provider<ImageUploadService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return ImageUploadService(apiClient: apiClient);
});

/// Upload state for UI
class UploadState {
  final bool isUploading;
  final double progress; // 0.0 to 1.0
  final UploadResult? result;
  final File? selectedFile;

  const UploadState({
    this.isUploading = false,
    this.progress = 0.0,
    this.result,
    this.selectedFile,
  });

  UploadState copyWith({
    bool? isUploading,
    double? progress,
    UploadResult? result,
    File? selectedFile,
  }) {
    return UploadState(
      isUploading: isUploading ?? this.isUploading,
      progress: progress ?? this.progress,
      result: result,
      selectedFile: selectedFile ?? this.selectedFile,
    );
  }

  bool get isSuccess => result?.success ?? false;
  bool get isFailed => result != null && !result!.success;
  int get progressPercent => (progress * 100).round();
}

/// Upload state notifier
class UploadNotifier extends StateNotifier<UploadState> {
  final ImageUploadService _uploadService;

  UploadNotifier(this._uploadService) : super(const UploadState());

  /// Pick image
  Future<bool> pickImage({
    ImageSource source = ImageSource.gallery,
    ImageQuality quality = ImageQuality.medium,
  }) async {
    final file = await _uploadService.pickAndCompressImage(
      source: source,
      quality: quality,
    );

    if (file != null) {
      state = state.copyWith(selectedFile: file);
      return true;
    }
    return false;
  }

  /// Upload profile photo
  Future<UploadResult> uploadProfilePhoto({
    ImageQuality quality = ImageQuality.medium,
  }) async {
    if (state.selectedFile == null) {
      return UploadResult.failed(error: 'No file selected');
    }

    state = state.copyWith(isUploading: true, progress: 0.0);

    final result = await _uploadService.uploadProfilePhoto(
      state.selectedFile!,
      quality: quality,
      onProgress: (sent, total) {
        state = state.copyWith(progress: sent / total);
      },
    );

    state = state.copyWith(
      isUploading: false,
      progress: result.success ? 1.0 : 0.0,
      result: result,
    );

    return result;
  }

  /// Clear selected file
  void clearSelection() {
    state = state.copyWith(selectedFile: null, result: null);
  }

  /// Reset state
  void reset() {
    state = const UploadState();
  }
}

/// Upload notifier provider
final uploadProvider = StateNotifierProvider<UploadNotifier, UploadState>((ref) {
  final service = ref.watch(imageUploadServiceProvider);
  return UploadNotifier(service);
});
