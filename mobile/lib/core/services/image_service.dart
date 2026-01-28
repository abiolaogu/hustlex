import 'dart:async';
import 'dart:io';
import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_cropper/image_cropper.dart';
import 'package:image_picker/image_picker.dart';
import 'package:logger/logger.dart';
import 'package:permission_handler/permission_handler.dart';

import '../api/api_client.dart';
import '../constants/app_colors.dart';

/// Image source selection
enum ImageSourceOption {
  camera,
  gallery,
}

/// Image upload result
class ImageUploadResult {
  final bool success;
  final String? imageUrl;
  final String? message;

  ImageUploadResult({
    required this.success,
    this.imageUrl,
    this.message,
  });

  factory ImageUploadResult.success(String imageUrl) {
    return ImageUploadResult(
      success: true,
      imageUrl: imageUrl,
    );
  }

  factory ImageUploadResult.failure(String message) {
    return ImageUploadResult(
      success: false,
      message: message,
    );
  }

  factory ImageUploadResult.cancelled() {
    return ImageUploadResult(
      success: false,
      message: 'Image selection cancelled',
    );
  }
}

/// Image pick result
class ImagePickResult {
  final bool success;
  final File? file;
  final Uint8List? bytes;
  final String? message;

  ImagePickResult({
    required this.success,
    this.file,
    this.bytes,
    this.message,
  });

  factory ImagePickResult.success(File file, [Uint8List? bytes]) {
    return ImagePickResult(
      success: true,
      file: file,
      bytes: bytes,
    );
  }

  factory ImagePickResult.failure(String message) {
    return ImagePickResult(
      success: false,
      message: message,
    );
  }

  factory ImagePickResult.cancelled() {
    return ImagePickResult(
      success: false,
      message: 'Selection cancelled',
    );
  }
}

/// Image service for picking, cropping, and uploading images
class ImageService {
  final Logger _logger = Logger();
  final ImagePicker _picker;
  final ApiClient _apiClient;

  ImageService({
    required ApiClient apiClient,
    ImagePicker? picker,
  })  : _apiClient = apiClient,
        _picker = picker ?? ImagePicker();

  /// Pick image from camera or gallery
  Future<ImagePickResult> pickImage({
    required ImageSourceOption source,
    int maxWidth = 1024,
    int maxHeight = 1024,
    int quality = 85,
  }) async {
    try {
      // Check permissions
      final hasPermission = await _checkPermission(source);
      if (!hasPermission) {
        return ImagePickResult.failure(
          source == ImageSourceOption.camera
              ? 'Camera permission denied'
              : 'Gallery permission denied',
        );
      }

      // Pick image
      final XFile? pickedFile = await _picker.pickImage(
        source: source == ImageSourceOption.camera
            ? ImageSource.camera
            : ImageSource.gallery,
        maxWidth: maxWidth.toDouble(),
        maxHeight: maxHeight.toDouble(),
        imageQuality: quality,
        preferredCameraDevice: CameraDevice.front,
      );

      if (pickedFile == null) {
        return ImagePickResult.cancelled();
      }

      final file = File(pickedFile.path);
      final bytes = await file.readAsBytes();

      _logger.i('Image picked: ${pickedFile.path}');
      return ImagePickResult.success(file, bytes);
    } catch (e) {
      _logger.e('Error picking image', error: e);
      return ImagePickResult.failure('Failed to pick image: $e');
    }
  }

  /// Check and request permission
  Future<bool> _checkPermission(ImageSourceOption source) async {
    Permission permission;

    if (source == ImageSourceOption.camera) {
      permission = Permission.camera;
    } else {
      // For gallery, check photos permission
      if (Platform.isIOS) {
        permission = Permission.photos;
      } else {
        // Android 13+ uses different permissions
        final androidVersion = int.tryParse(
          Platform.version.split('.').first,
        );
        if (androidVersion != null && androidVersion >= 33) {
          permission = Permission.photos;
        } else {
          permission = Permission.storage;
        }
      }
    }

    final status = await permission.status;

    if (status.isGranted) {
      return true;
    }

    if (status.isDenied) {
      final result = await permission.request();
      return result.isGranted;
    }

    if (status.isPermanentlyDenied) {
      // Open app settings
      await openAppSettings();
      return false;
    }

    return false;
  }

  /// Crop image with circular mask (for profile photos)
  Future<ImagePickResult> cropImageCircular({
    required File imageFile,
    String? title,
  }) async {
    try {
      final croppedFile = await ImageCropper().cropImage(
        sourcePath: imageFile.path,
        cropStyle: CropStyle.circle,
        aspectRatio: const CropAspectRatio(ratioX: 1, ratioY: 1),
        compressQuality: 85,
        maxWidth: 512,
        maxHeight: 512,
        uiSettings: [
          AndroidUiSettings(
            toolbarTitle: title ?? 'Crop Image',
            toolbarColor: AppColors.primary,
            toolbarWidgetColor: Colors.white,
            initAspectRatio: CropAspectRatioPreset.square,
            lockAspectRatio: true,
            hideBottomControls: false,
          ),
          IOSUiSettings(
            title: title ?? 'Crop Image',
            aspectRatioLockEnabled: true,
            resetAspectRatioEnabled: false,
            rotateButtonsHidden: true,
            aspectRatioPickerButtonHidden: true,
          ),
        ],
      );

      if (croppedFile == null) {
        return ImagePickResult.cancelled();
      }

      final file = File(croppedFile.path);
      final bytes = await file.readAsBytes();

      _logger.i('Image cropped: ${croppedFile.path}');
      return ImagePickResult.success(file, bytes);
    } catch (e) {
      _logger.e('Error cropping image', error: e);
      return ImagePickResult.failure('Failed to crop image: $e');
    }
  }

  /// Crop image with rectangle (for documents)
  Future<ImagePickResult> cropImageRectangle({
    required File imageFile,
    CropAspectRatio? aspectRatio,
    String? title,
  }) async {
    try {
      final croppedFile = await ImageCropper().cropImage(
        sourcePath: imageFile.path,
        cropStyle: CropStyle.rectangle,
        aspectRatio: aspectRatio,
        compressQuality: 85,
        maxWidth: 1920,
        maxHeight: 1920,
        uiSettings: [
          AndroidUiSettings(
            toolbarTitle: title ?? 'Crop Document',
            toolbarColor: AppColors.primary,
            toolbarWidgetColor: Colors.white,
            initAspectRatio: aspectRatio != null
                ? CropAspectRatioPreset.original
                : CropAspectRatioPreset.original,
            lockAspectRatio: aspectRatio != null,
          ),
          IOSUiSettings(
            title: title ?? 'Crop Document',
            aspectRatioLockEnabled: aspectRatio != null,
            resetAspectRatioEnabled: aspectRatio == null,
          ),
        ],
      );

      if (croppedFile == null) {
        return ImagePickResult.cancelled();
      }

      final file = File(croppedFile.path);
      final bytes = await file.readAsBytes();

      return ImagePickResult.success(file, bytes);
    } catch (e) {
      _logger.e('Error cropping document', error: e);
      return ImagePickResult.failure('Failed to crop document: $e');
    }
  }

  /// Pick and crop profile image
  Future<ImagePickResult> pickProfileImage({
    required ImageSourceOption source,
  }) async {
    // Pick image
    final pickResult = await pickImage(
      source: source,
      maxWidth: 1024,
      maxHeight: 1024,
      quality: 90,
    );

    if (!pickResult.success || pickResult.file == null) {
      return pickResult;
    }

    // Crop to circular
    return await cropImageCircular(
      imageFile: pickResult.file!,
      title: 'Crop Profile Photo',
    );
  }

  /// Pick document image (ID card, license, etc.)
  Future<ImagePickResult> pickDocumentImage({
    required ImageSourceOption source,
    CropAspectRatio? aspectRatio,
    String? documentType,
  }) async {
    final pickResult = await pickImage(
      source: source,
      maxWidth: 1920,
      maxHeight: 1920,
      quality: 95,
    );

    if (!pickResult.success || pickResult.file == null) {
      return pickResult;
    }

    return await cropImageRectangle(
      imageFile: pickResult.file!,
      aspectRatio: aspectRatio,
      title: documentType != null ? 'Crop $documentType' : 'Crop Document',
    );
  }

  /// Upload profile image
  Future<ImageUploadResult> uploadProfileImage(File file) async {
    try {
      _logger.i('Uploading profile image...');

      final response = await _apiClient.uploadFile(
        endpoint: '/users/profile/photo',
        file: file,
        fieldName: 'photo',
      );

      final imageUrl = response['data']?['url'] as String?;
      
      if (imageUrl != null) {
        _logger.i('Profile image uploaded: $imageUrl');
        return ImageUploadResult.success(imageUrl);
      }

      return ImageUploadResult.failure('Failed to get image URL');
    } catch (e) {
      _logger.e('Error uploading profile image', error: e);
      return ImageUploadResult.failure('Upload failed: $e');
    }
  }

  /// Upload KYC document
  Future<ImageUploadResult> uploadKycDocument({
    required File file,
    required String documentType,
  }) async {
    try {
      _logger.i('Uploading KYC document: $documentType');

      final response = await _apiClient.uploadFile(
        endpoint: '/kyc/documents',
        file: file,
        fieldName: 'document',
        additionalFields: {'type': documentType},
      );

      final imageUrl = response['data']?['url'] as String?;

      if (imageUrl != null) {
        _logger.i('KYC document uploaded: $imageUrl');
        return ImageUploadResult.success(imageUrl);
      }

      return ImageUploadResult.failure('Failed to get document URL');
    } catch (e) {
      _logger.e('Error uploading KYC document', error: e);
      return ImageUploadResult.failure('Upload failed: $e');
    }
  }

  /// Upload gig attachment
  Future<ImageUploadResult> uploadGigAttachment({
    required File file,
    required String gigId,
  }) async {
    try {
      _logger.i('Uploading gig attachment for: $gigId');

      final response = await _apiClient.uploadFile(
        endpoint: '/gigs/$gigId/attachments',
        file: file,
        fieldName: 'attachment',
      );

      final imageUrl = response['data']?['url'] as String?;

      if (imageUrl != null) {
        return ImageUploadResult.success(imageUrl);
      }

      return ImageUploadResult.failure('Failed to get attachment URL');
    } catch (e) {
      _logger.e('Error uploading gig attachment', error: e);
      return ImageUploadResult.failure('Upload failed: $e');
    }
  }

  /// Show image source selection dialog
  static Future<ImageSourceOption?> showSourcePicker(BuildContext context) async {
    return await showModalBottomSheet<ImageSourceOption>(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey[300],
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 20),
              const Text(
                'Select Image Source',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 20),
              ListTile(
                leading: const CircleAvatar(
                  backgroundColor: AppColors.primary,
                  child: Icon(Icons.camera_alt, color: Colors.white),
                ),
                title: const Text('Camera'),
                subtitle: const Text('Take a new photo'),
                onTap: () => Navigator.pop(context, ImageSourceOption.camera),
              ),
              ListTile(
                leading: CircleAvatar(
                  backgroundColor: AppColors.primary.withOpacity(0.8),
                  child: const Icon(Icons.photo_library, color: Colors.white),
                ),
                title: const Text('Gallery'),
                subtitle: const Text('Choose from gallery'),
                onTap: () => Navigator.pop(context, ImageSourceOption.gallery),
              ),
              const SizedBox(height: 10),
            ],
          ),
        ),
      ),
    );
  }
}

/// Image upload state
class ImageUploadState {
  final bool isPicking;
  final bool isUploading;
  final double uploadProgress;
  final ImagePickResult? pickResult;
  final ImageUploadResult? uploadResult;
  final String? error;

  const ImageUploadState({
    this.isPicking = false,
    this.isUploading = false,
    this.uploadProgress = 0,
    this.pickResult,
    this.uploadResult,
    this.error,
  });

  ImageUploadState copyWith({
    bool? isPicking,
    bool? isUploading,
    double? uploadProgress,
    ImagePickResult? pickResult,
    ImageUploadResult? uploadResult,
    String? error,
  }) {
    return ImageUploadState(
      isPicking: isPicking ?? this.isPicking,
      isUploading: isUploading ?? this.isUploading,
      uploadProgress: uploadProgress ?? this.uploadProgress,
      pickResult: pickResult,
      uploadResult: uploadResult,
      error: error,
    );
  }

  bool get hasImage => pickResult?.file != null;
  File? get imageFile => pickResult?.file;
}

/// Image upload notifier
class ImageUploadNotifier extends StateNotifier<ImageUploadState> {
  final ImageService _service;

  ImageUploadNotifier(this._service) : super(const ImageUploadState());

  /// Pick profile image
  Future<void> pickProfileImage(ImageSourceOption source) async {
    state = state.copyWith(isPicking: true, error: null);

    final result = await _service.pickProfileImage(source: source);

    state = state.copyWith(
      isPicking: false,
      pickResult: result,
      error: result.success ? null : result.message,
    );
  }

  /// Upload profile image
  Future<ImageUploadResult> uploadProfileImage() async {
    if (state.imageFile == null) {
      return ImageUploadResult.failure('No image selected');
    }

    state = state.copyWith(isUploading: true, uploadProgress: 0, error: null);

    final result = await _service.uploadProfileImage(state.imageFile!);

    state = state.copyWith(
      isUploading: false,
      uploadProgress: result.success ? 1.0 : 0,
      uploadResult: result,
      error: result.success ? null : result.message,
    );

    return result;
  }

  /// Pick and upload profile image in one step
  Future<ImageUploadResult?> pickAndUploadProfileImage(
    ImageSourceOption source,
  ) async {
    await pickProfileImage(source);

    if (!state.hasImage) {
      return null;
    }

    return await uploadProfileImage();
  }

  /// Reset state
  void reset() {
    state = const ImageUploadState();
  }

  /// Clear error
  void clearError() {
    state = state.copyWith(error: null);
  }
}

/// Image service provider
final imageServiceProvider = Provider<ImageService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return ImageService(apiClient: apiClient);
});

/// API client provider (duplicate for standalone usage)
final apiClientProvider = Provider<ApiClient>((ref) {
  throw UnimplementedError('Must be overridden in ProviderScope');
});

/// Image upload provider
final imageUploadProvider =
    StateNotifierProvider<ImageUploadNotifier, ImageUploadState>((ref) {
  final service = ref.watch(imageServiceProvider);
  return ImageUploadNotifier(service);
});

/// Profile image upload provider (separate instance)
final profileImageProvider =
    StateNotifierProvider<ImageUploadNotifier, ImageUploadState>((ref) {
  final service = ref.watch(imageServiceProvider);
  return ImageUploadNotifier(service);
});

/// KYC document upload provider
final kycDocumentProvider =
    StateNotifierProvider<ImageUploadNotifier, ImageUploadState>((ref) {
  final service = ref.watch(imageServiceProvider);
  return ImageUploadNotifier(service);
});
