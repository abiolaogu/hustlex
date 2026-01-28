import 'package:dartz/dartz.dart';
import '../failures/value_failures.dart';
import 'value_object.dart';

/// Nigerian phone number value object.
/// Validates and normalizes Nigerian mobile phone numbers.
class PhoneNumber extends ValueObject<String> {
  @override
  final Either<ValueFailure<String>, String> value;

  factory PhoneNumber(String input) {
    return PhoneNumber._(_validateAndNormalize(input));
  }

  const PhoneNumber._(this.value);

  static Either<ValueFailure<String>, String> _validateAndNormalize(
      String input) {
    if (input.isEmpty) {
      return left(ValueFailure.empty(failedValue: input));
    }

    // Remove all non-digit characters except leading +
    String cleaned = input.replaceAll(RegExp(r'[^\d+]'), '');

    // Remove leading + for processing
    if (cleaned.startsWith('+')) {
      cleaned = cleaned.substring(1);
    }

    String normalized;

    // Handle different formats
    if (cleaned.length == 11 && cleaned.startsWith('0')) {
      // Local format: 08012345678 -> 2348012345678
      normalized = '234${cleaned.substring(1)}';
    } else if (cleaned.length == 13 && cleaned.startsWith('234')) {
      // International without +: 2348012345678
      normalized = cleaned;
    } else if (cleaned.length == 10 && !cleaned.startsWith('0')) {
      // Without leading zero: 8012345678 -> 2348012345678
      normalized = '234$cleaned';
    } else {
      return left(ValueFailure.invalidPhoneNumber(failedValue: input));
    }

    // Validate Nigerian mobile prefixes
    final validPrefixes = <String>{
      // MTN
      '234703',
      '234706',
      '234803',
      '234806',
      '234810',
      '234813',
      '234814',
      '234816',
      '234903',
      '234906',
      '234913',
      '234916',
      // Glo
      '234705',
      '234805',
      '234807',
      '234811',
      '234815',
      '234905',
      '234915',
      // Airtel
      '234701',
      '234708',
      '234802',
      '234808',
      '234812',
      '234901',
      '234902',
      '234904',
      '234907',
      '234912',
      // 9mobile
      '234809',
      '234817',
      '234818',
      '234909',
      '234908',
    };

    final prefix = normalized.substring(0, 6);
    if (!validPrefixes.contains(prefix)) {
      return left(ValueFailure.invalidPhoneNumber(failedValue: input));
    }

    return right(normalized);
  }

  /// Returns phone in local format (e.g., 08012345678)
  String get localFormat {
    final normalized = getOrCrash();
    return '0${normalized.substring(3)}';
  }

  /// Returns phone in international format (e.g., +2348012345678)
  String get internationalFormat {
    return '+${getOrCrash()}';
  }

  /// Returns masked phone for display (e.g., 080****5678)
  String get masked {
    final local = localFormat;
    return '${local.substring(0, 3)}****${local.substring(7)}';
  }

  /// Returns the raw normalized number without + (e.g., 2348012345678)
  String get raw => getOrCrash();
}
