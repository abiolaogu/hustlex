import 'package:intl/intl.dart';

/// Currency formatting utilities
class CurrencyUtils {
  static final _nairaFormat = NumberFormat.currency(
    locale: 'en_NG',
    symbol: '₦',
    decimalDigits: 2,
  );

  static final _nairaFormatNoDecimal = NumberFormat.currency(
    locale: 'en_NG',
    symbol: '₦',
    decimalDigits: 0,
  );

  /// Format amount to Naira with 2 decimal places
  static String formatNaira(double amount) {
    return _nairaFormat.format(amount);
  }

  /// Format amount to Naira without decimal places
  static String formatNairaWhole(double amount) {
    return _nairaFormatNoDecimal.format(amount);
  }

  /// Format amount with sign (+ or -)
  static String formatWithSign(double amount) {
    final prefix = amount > 0 ? '+' : '';
    return '$prefix${_nairaFormat.format(amount)}';
  }

  /// Parse formatted currency string to double
  static double parseNaira(String value) {
    final cleaned = value.replaceAll(RegExp(r'[₦,\s]'), '');
    return double.tryParse(cleaned) ?? 0;
  }

  /// Format large numbers (e.g., 1.5M, 500K)
  static String formatCompact(double amount) {
    if (amount >= 1000000) {
      return '₦${(amount / 1000000).toStringAsFixed(1)}M';
    } else if (amount >= 1000) {
      return '₦${(amount / 1000).toStringAsFixed(0)}K';
    }
    return formatNairaWhole(amount);
  }
}

/// Date/Time formatting utilities
class DateUtils {
  /// Format date to readable string (e.g., "Jan 15, 2026")
  static String formatDate(DateTime date) {
    return DateFormat('MMM d, yyyy').format(date);
  }

  /// Format date to short string (e.g., "Jan 15")
  static String formatDateShort(DateTime date) {
    return DateFormat('MMM d').format(date);
  }

  /// Format time (e.g., "2:30 PM")
  static String formatTime(DateTime date) {
    return DateFormat('h:mm a').format(date);
  }

  /// Format date and time (e.g., "Jan 15, 2026 at 2:30 PM")
  static String formatDateTime(DateTime date) {
    return '${formatDate(date)} at ${formatTime(date)}';
  }

  /// Format relative time (e.g., "2 hours ago", "Yesterday")
  static String formatRelative(DateTime date) {
    final now = DateTime.now();
    final difference = now.difference(date);

    if (difference.inSeconds < 60) {
      return 'Just now';
    } else if (difference.inMinutes < 60) {
      return '${difference.inMinutes} min ago';
    } else if (difference.inHours < 24) {
      return '${difference.inHours} hour${difference.inHours > 1 ? 's' : ''} ago';
    } else if (difference.inDays == 1) {
      return 'Yesterday';
    } else if (difference.inDays < 7) {
      return '${difference.inDays} days ago';
    } else if (difference.inDays < 30) {
      final weeks = (difference.inDays / 7).floor();
      return '$weeks week${weeks > 1 ? 's' : ''} ago';
    } else {
      return formatDate(date);
    }
  }

  /// Get greeting based on time of day
  static String getGreeting() {
    final hour = DateTime.now().hour;
    if (hour < 12) return 'Good morning';
    if (hour < 17) return 'Good afternoon';
    return 'Good evening';
  }

  /// Check if date is today
  static bool isToday(DateTime date) {
    final now = DateTime.now();
    return date.year == now.year &&
        date.month == now.month &&
        date.day == now.day;
  }

  /// Check if date is in the past
  static bool isPast(DateTime date) {
    return date.isBefore(DateTime.now());
  }

  /// Get days until date
  static int daysUntil(DateTime date) {
    return date.difference(DateTime.now()).inDays;
  }
}

/// Phone number utilities
class PhoneUtils {
  /// Format Nigerian phone number
  static String formatNigerian(String phone) {
    String digits = phone.replaceAll(RegExp(r'\D'), '');

    // Convert to international format
    if (digits.startsWith('0') && digits.length == 11) {
      digits = '234${digits.substring(1)}';
    }

    if (digits.length == 13 && digits.startsWith('234')) {
      return '+$digits';
    }

    return phone;
  }

  /// Mask phone number (e.g., +234 *** *** 1234)
  static String mask(String phone) {
    if (phone.length < 8) return phone;
    final start = phone.substring(0, 4);
    final end = phone.substring(phone.length - 4);
    return '$start *** *** $end';
  }

  /// Validate Nigerian phone number
  static bool isValidNigerian(String phone) {
    final digits = phone.replaceAll(RegExp(r'\D'), '');
    return (digits.length == 11 && digits.startsWith('0')) ||
        (digits.length == 10) ||
        (digits.length == 13 && digits.startsWith('234'));
  }
}

/// String utilities
class StringUtils {
  /// Capitalize first letter of each word
  static String capitalize(String text) {
    if (text.isEmpty) return text;
    return text.split(' ').map((word) {
      if (word.isEmpty) return word;
      return word[0].toUpperCase() + word.substring(1).toLowerCase();
    }).join(' ');
  }

  /// Get initials from name
  static String getInitials(String name, {int count = 2}) {
    final parts = name.trim().split(' ').where((p) => p.isNotEmpty).toList();
    if (parts.isEmpty) return '';
    if (parts.length == 1) {
      return parts[0][0].toUpperCase();
    }
    return parts.take(count).map((p) => p[0].toUpperCase()).join();
  }

  /// Truncate text with ellipsis
  static String truncate(String text, int maxLength) {
    if (text.length <= maxLength) return text;
    return '${text.substring(0, maxLength)}...';
  }

  /// Generate reference number
  static String generateReference({String prefix = 'HX'}) {
    final timestamp = DateTime.now().millisecondsSinceEpoch;
    return '$prefix$timestamp';
  }
}

/// Validation utilities
class ValidationUtils {
  /// Validate email
  static String? validateEmail(String? value) {
    if (value == null || value.isEmpty) {
      return 'Email is required';
    }
    final emailRegex = RegExp(
      r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
    );
    if (!emailRegex.hasMatch(value)) {
      return 'Please enter a valid email';
    }
    return null;
  }

  /// Validate phone
  static String? validatePhone(String? value) {
    if (value == null || value.isEmpty) {
      return 'Phone number is required';
    }
    if (!PhoneUtils.isValidNigerian(value)) {
      return 'Please enter a valid phone number';
    }
    return null;
  }

  /// Validate name
  static String? validateName(String? value, {String fieldName = 'Name'}) {
    if (value == null || value.trim().isEmpty) {
      return '$fieldName is required';
    }
    if (value.trim().length < 2) {
      return '$fieldName must be at least 2 characters';
    }
    return null;
  }

  /// Validate amount
  static String? validateAmount(
    String? value, {
    double min = 0,
    double? max,
  }) {
    if (value == null || value.isEmpty) {
      return 'Amount is required';
    }
    final amount = CurrencyUtils.parseNaira(value);
    if (amount <= 0) {
      return 'Please enter a valid amount';
    }
    if (amount < min) {
      return 'Minimum amount is ${CurrencyUtils.formatNairaWhole(min)}';
    }
    if (max != null && amount > max) {
      return 'Maximum amount is ${CurrencyUtils.formatNairaWhole(max)}';
    }
    return null;
  }

  /// Validate PIN
  static String? validatePin(String? value, {int length = 4}) {
    if (value == null || value.isEmpty) {
      return 'PIN is required';
    }
    if (value.length != length) {
      return 'PIN must be $length digits';
    }
    if (!RegExp(r'^\d+$').hasMatch(value)) {
      return 'PIN must contain only numbers';
    }
    // Check for simple patterns
    if (value == '0000' || value == '1234' || value == '1111') {
      return 'Please choose a stronger PIN';
    }
    return null;
  }
}

/// Number utilities
class NumberUtils {
  /// Format percentage
  static String formatPercent(double value, {int decimals = 0}) {
    return '${value.toStringAsFixed(decimals)}%';
  }

  /// Clamp value between min and max
  static double clamp(double value, double min, double max) {
    if (value < min) return min;
    if (value > max) return max;
    return value;
  }

  /// Linear interpolation
  static double lerp(double start, double end, double t) {
    return start + (end - start) * t;
  }
}
