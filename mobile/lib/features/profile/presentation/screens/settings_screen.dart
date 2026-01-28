import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_typography.dart';

// Settings state providers
final notificationEnabledProvider = StateProvider<bool>((ref) => true);
final pushNotificationsProvider = StateProvider<bool>((ref) => true);
final emailNotificationsProvider = StateProvider<bool>((ref) => true);
final smsNotificationsProvider = StateProvider<bool>((ref) => false);
final transactionAlertsProvider = StateProvider<bool>((ref) => true);
final savingsRemindersProvider = StateProvider<bool>((ref) => true);
final gigUpdatesProvider = StateProvider<bool>((ref) => true);
final marketingProvider = StateProvider<bool>((ref) => false);

final biometricEnabledProvider = StateProvider<bool>((ref) => false);
final darkModeProvider = StateProvider<bool>((ref) => false);
final selectedLanguageProvider = StateProvider<String>((ref) => 'English');
final selectedCurrencyProvider = StateProvider<String>((ref) => 'NGN');

class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
      ),
      body: ListView(
        children: [
          _buildSection(
            'Account',
            [
              _SettingsTile(
                icon: Icons.person_outline,
                title: 'Edit Profile',
                subtitle: 'Update your personal information',
                onTap: () => context.push('/profile/edit'),
              ),
              _SettingsTile(
                icon: Icons.lock_outline,
                title: 'Change PIN',
                subtitle: 'Update your transaction PIN',
                onTap: () => context.push('/profile/change-pin'),
              ),
              _SettingsTile(
                icon: Icons.account_balance_outlined,
                title: 'Bank Accounts',
                subtitle: 'Manage linked bank accounts',
                onTap: () => context.push('/wallet/banks'),
              ),
              Consumer(
                builder: (context, ref, child) {
                  final biometricEnabled = ref.watch(biometricEnabledProvider);
                  return _SettingsSwitchTile(
                    icon: Icons.fingerprint,
                    title: 'Biometric Login',
                    subtitle: 'Use fingerprint or face ID',
                    value: biometricEnabled,
                    onChanged: (value) {
                      ref.read(biometricEnabledProvider.notifier).state = value;
                    },
                  );
                },
              ),
            ],
          ),
          _buildSection(
            'Notifications',
            [
              Consumer(
                builder: (context, ref, child) {
                  final pushEnabled = ref.watch(pushNotificationsProvider);
                  return _SettingsSwitchTile(
                    icon: Icons.notifications_outlined,
                    title: 'Push Notifications',
                    subtitle: 'Receive app notifications',
                    value: pushEnabled,
                    onChanged: (value) {
                      ref.read(pushNotificationsProvider.notifier).state = value;
                    },
                  );
                },
              ),
              Consumer(
                builder: (context, ref, child) {
                  final emailEnabled = ref.watch(emailNotificationsProvider);
                  return _SettingsSwitchTile(
                    icon: Icons.email_outlined,
                    title: 'Email Notifications',
                    subtitle: 'Receive email updates',
                    value: emailEnabled,
                    onChanged: (value) {
                      ref.read(emailNotificationsProvider.notifier).state = value;
                    },
                  );
                },
              ),
              Consumer(
                builder: (context, ref, child) {
                  final smsEnabled = ref.watch(smsNotificationsProvider);
                  return _SettingsSwitchTile(
                    icon: Icons.sms_outlined,
                    title: 'SMS Notifications',
                    subtitle: 'Receive text messages',
                    value: smsEnabled,
                    onChanged: (value) {
                      ref.read(smsNotificationsProvider.notifier).state = value;
                    },
                  );
                },
              ),
              _SettingsTile(
                icon: Icons.tune,
                title: 'Notification Preferences',
                subtitle: 'Customize notification types',
                onTap: () => _showNotificationPreferences(context, ref),
              ),
            ],
          ),
          _buildSection(
            'Appearance',
            [
              Consumer(
                builder: (context, ref, child) {
                  final darkMode = ref.watch(darkModeProvider);
                  return _SettingsSwitchTile(
                    icon: Icons.dark_mode_outlined,
                    title: 'Dark Mode',
                    subtitle: 'Use dark theme',
                    value: darkMode,
                    onChanged: (value) {
                      ref.read(darkModeProvider.notifier).state = value;
                    },
                  );
                },
              ),
              Consumer(
                builder: (context, ref, child) {
                  final language = ref.watch(selectedLanguageProvider);
                  return _SettingsTile(
                    icon: Icons.language,
                    title: 'Language',
                    subtitle: language,
                    onTap: () => _showLanguageSelector(context, ref),
                  );
                },
              ),
              Consumer(
                builder: (context, ref, child) {
                  final currency = ref.watch(selectedCurrencyProvider);
                  return _SettingsTile(
                    icon: Icons.attach_money,
                    title: 'Currency',
                    subtitle: currency,
                    onTap: () => _showCurrencySelector(context, ref),
                  );
                },
              ),
            ],
          ),
          _buildSection(
            'Privacy & Security',
            [
              _SettingsTile(
                icon: Icons.privacy_tip_outlined,
                title: 'Privacy Policy',
                subtitle: 'Read our privacy policy',
                onTap: () => _showWebView(context, 'Privacy Policy', 'https://hustlex.app/privacy'),
              ),
              _SettingsTile(
                icon: Icons.description_outlined,
                title: 'Terms of Service',
                subtitle: 'Read our terms',
                onTap: () => _showWebView(context, 'Terms of Service', 'https://hustlex.app/terms'),
              ),
              _SettingsTile(
                icon: Icons.security,
                title: 'Security Tips',
                subtitle: 'Keep your account safe',
                onTap: () => _showSecurityTips(context),
              ),
            ],
          ),
          _buildSection(
            'Support',
            [
              _SettingsTile(
                icon: Icons.help_outline,
                title: 'Help Center',
                subtitle: 'Get help with HustleX',
                onTap: () {},
              ),
              _SettingsTile(
                icon: Icons.chat_bubble_outline,
                title: 'Contact Support',
                subtitle: 'Chat with our team',
                onTap: () => _showContactOptions(context),
              ),
              _SettingsTile(
                icon: Icons.bug_report_outlined,
                title: 'Report a Bug',
                subtitle: 'Help us improve',
                onTap: () {},
              ),
            ],
          ),
          _buildSection(
            'About',
            [
              _SettingsTile(
                icon: Icons.info_outline,
                title: 'About HustleX',
                subtitle: 'Version 1.0.0',
                onTap: () => _showAboutDialog(context),
              ),
              _SettingsTile(
                icon: Icons.star_outline,
                title: 'Rate App',
                subtitle: 'Rate us on the app store',
                onTap: () {},
              ),
              _SettingsTile(
                icon: Icons.share_outlined,
                title: 'Share App',
                subtitle: 'Tell friends about HustleX',
                onTap: () {},
              ),
            ],
          ),
          const SizedBox(height: 16),
          _buildLogoutButton(context),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSection(String title, List<Widget> children) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 24, 16, 8),
          child: Text(
            title,
            style: AppTypography.titleSmall.copyWith(
              color: AppColors.textSecondary,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
        ...children,
      ],
    );
  }

  Widget _buildLogoutButton(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: OutlinedButton(
        onPressed: () => _showLogoutDialog(context),
        style: OutlinedButton.styleFrom(
          foregroundColor: AppColors.error,
          side: const BorderSide(color: AppColors.error),
          padding: const EdgeInsets.symmetric(vertical: 16),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.logout),
            const SizedBox(width: 8),
            const Text('Log Out'),
          ],
        ),
      ),
    );
  }

  void _showNotificationPreferences(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.6,
        maxChildSize: 0.9,
        minChildSize: 0.4,
        expand: false,
        builder: (context, scrollController) => Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: AppColors.border,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Notification Preferences',
                    style: AppTypography.titleMedium.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ),
            Expanded(
              child: ListView(
                controller: scrollController,
                children: [
                  Consumer(
                    builder: (context, ref, child) {
                      final enabled = ref.watch(transactionAlertsProvider);
                      return _SettingsSwitchTile(
                        icon: Icons.receipt_long_outlined,
                        title: 'Transaction Alerts',
                        subtitle: 'Notifications for deposits, withdrawals, and transfers',
                        value: enabled,
                        onChanged: (value) {
                          ref.read(transactionAlertsProvider.notifier).state = value;
                        },
                      );
                    },
                  ),
                  Consumer(
                    builder: (context, ref, child) {
                      final enabled = ref.watch(savingsRemindersProvider);
                      return _SettingsSwitchTile(
                        icon: Icons.savings_outlined,
                        title: 'Savings Reminders',
                        subtitle: 'Reminders for circle contributions',
                        value: enabled,
                        onChanged: (value) {
                          ref.read(savingsRemindersProvider.notifier).state = value;
                        },
                      );
                    },
                  ),
                  Consumer(
                    builder: (context, ref, child) {
                      final enabled = ref.watch(gigUpdatesProvider);
                      return _SettingsSwitchTile(
                        icon: Icons.work_outline,
                        title: 'Gig Updates',
                        subtitle: 'New gigs and proposal updates',
                        value: enabled,
                        onChanged: (value) {
                          ref.read(gigUpdatesProvider.notifier).state = value;
                        },
                      );
                    },
                  ),
                  Consumer(
                    builder: (context, ref, child) {
                      final enabled = ref.watch(marketingProvider);
                      return _SettingsSwitchTile(
                        icon: Icons.campaign_outlined,
                        title: 'Marketing',
                        subtitle: 'Promotions, offers, and announcements',
                        value: enabled,
                        onChanged: (value) {
                          ref.read(marketingProvider.notifier).state = value;
                        },
                      );
                    },
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showLanguageSelector(BuildContext context, WidgetRef ref) {
    final languages = ['English', 'Yoruba', 'Igbo', 'Hausa', 'Pidgin'];
    final currentLanguage = ref.read(selectedLanguageProvider);

    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: AppColors.border,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Select Language',
                    style: AppTypography.titleMedium.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ),
            ...languages.map((language) {
              final isSelected = currentLanguage == language;
              return ListTile(
                title: Text(language),
                trailing: isSelected
                    ? const Icon(Icons.check, color: AppColors.primary)
                    : null,
                onTap: () {
                  ref.read(selectedLanguageProvider.notifier).state = language;
                  Navigator.pop(context);
                },
              );
            }),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  void _showCurrencySelector(BuildContext context, WidgetRef ref) {
    final currencies = [
      {'code': 'NGN', 'name': 'Nigerian Naira', 'symbol': '₦'},
      {'code': 'USD', 'name': 'US Dollar', 'symbol': '\$'},
      {'code': 'GBP', 'name': 'British Pound', 'symbol': '£'},
      {'code': 'EUR', 'name': 'Euro', 'symbol': '€'},
    ];
    final currentCurrency = ref.read(selectedCurrencyProvider);

    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: AppColors.border,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Select Currency',
                    style: AppTypography.titleMedium.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ),
            ...currencies.map((currency) {
              final isSelected = currentCurrency == currency['code'];
              return ListTile(
                leading: Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: AppColors.surfaceVariant,
                    shape: BoxShape.circle,
                  ),
                  child: Center(
                    child: Text(
                      currency['symbol']!,
                      style: AppTypography.titleMedium.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ),
                title: Text(currency['code']!),
                subtitle: Text(currency['name']!),
                trailing: isSelected
                    ? const Icon(Icons.check, color: AppColors.primary)
                    : null,
                onTap: () {
                  ref.read(selectedCurrencyProvider.notifier).state = currency['code']!;
                  Navigator.pop(context);
                },
              );
            }),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  void _showWebView(BuildContext context, String title, String url) {
    // TODO: Implement WebView
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Opening $title...')),
    );
  }

  void _showSecurityTips(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.7,
        maxChildSize: 0.9,
        minChildSize: 0.5,
        expand: false,
        builder: (context, scrollController) => Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: AppColors.border,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Security Tips',
                    style: AppTypography.titleMedium.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ),
            Expanded(
              child: ListView(
                controller: scrollController,
                padding: const EdgeInsets.all(16),
                children: [
                  _SecurityTipCard(
                    icon: Icons.lock,
                    color: AppColors.primary,
                    title: 'Keep Your PIN Secret',
                    description: 'Never share your PIN with anyone, including HustleX staff.',
                  ),
                  _SecurityTipCard(
                    icon: Icons.phone_android,
                    color: AppColors.secondary,
                    title: 'Secure Your Device',
                    description: 'Use screen lock and keep your phone software updated.',
                  ),
                  _SecurityTipCard(
                    icon: Icons.wifi_off,
                    color: AppColors.warning,
                    title: 'Avoid Public WiFi',
                    description: 'Don\'t access your wallet on public or unsecured networks.',
                  ),
                  _SecurityTipCard(
                    icon: Icons.report_problem,
                    color: AppColors.error,
                    title: 'Report Suspicious Activity',
                    description: 'Contact support immediately if you notice unauthorized transactions.',
                  ),
                  _SecurityTipCard(
                    icon: Icons.verified_user,
                    color: AppColors.success,
                    title: 'Verify Official Communications',
                    description: 'HustleX will never ask for your password via email or phone.',
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showContactOptions(BuildContext context) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: AppColors.border,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 24),
              Text(
                'Contact Support',
                style: AppTypography.titleMedium.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 24),
              _ContactOption(
                icon: Icons.chat,
                title: 'Live Chat',
                subtitle: 'Chat with support',
                onTap: () {
                  Navigator.pop(context);
                  // TODO: Open live chat
                },
              ),
              _ContactOption(
                icon: Icons.email,
                title: 'Email',
                subtitle: 'support@hustlex.app',
                onTap: () {
                  Navigator.pop(context);
                  // TODO: Open email
                },
              ),
              _ContactOption(
                icon: Icons.phone,
                title: 'Phone',
                subtitle: '+234 800 HUSTLEX',
                onTap: () {
                  Navigator.pop(context);
                  // TODO: Make phone call
                },
              ),
              _ContactOption(
                icon: Icons.message,
                title: 'WhatsApp',
                subtitle: '+234 800 000 0000',
                onTap: () {
                  Navigator.pop(context);
                  // TODO: Open WhatsApp
                },
              ),
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }

  void _showAboutDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: AppColors.primary.withOpacity(0.1),
                shape: BoxShape.circle,
              ),
              child: Icon(
                Icons.bolt,
                color: AppColors.primary,
                size: 48,
              ),
            ),
            const SizedBox(height: 16),
            Text(
              'HustleX',
              style: AppTypography.headlineSmall.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              'Version 1.0.0 (Build 1)',
              style: AppTypography.bodySmall.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: 16),
            Text(
              'Your all-in-one platform for gigs, savings, and credit building.',
              style: AppTypography.bodyMedium,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            Text(
              '© 2024 HustleX. All rights reserved.',
              style: AppTypography.bodySmall.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  void _showLogoutDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: const Text('Log Out'),
        content: const Text('Are you sure you want to log out?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              // TODO: Perform logout
              context.go('/auth/login');
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.error,
              foregroundColor: Colors.white,
            ),
            child: const Text('Log Out'),
          ),
        ],
      ),
    );
  }
}

class _SettingsTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback? onTap;

  const _SettingsTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: AppColors.surfaceVariant,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Icon(icon, color: AppColors.textSecondary, size: 22),
      ),
      title: Text(
        title,
        style: AppTypography.bodyMedium.copyWith(
          fontWeight: FontWeight.w500,
        ),
      ),
      subtitle: Text(
        subtitle,
        style: AppTypography.bodySmall.copyWith(
          color: AppColors.textSecondary,
        ),
      ),
      trailing: const Icon(Icons.chevron_right, color: AppColors.textTertiary),
      onTap: onTap,
    );
  }
}

class _SettingsSwitchTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final bool value;
  final ValueChanged<bool> onChanged;

  const _SettingsSwitchTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.value,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return SwitchListTile(
      secondary: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: AppColors.surfaceVariant,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Icon(icon, color: AppColors.textSecondary, size: 22),
      ),
      title: Text(
        title,
        style: AppTypography.bodyMedium.copyWith(
          fontWeight: FontWeight.w500,
        ),
      ),
      subtitle: Text(
        subtitle,
        style: AppTypography.bodySmall.copyWith(
          color: AppColors.textSecondary,
        ),
      ),
      value: value,
      onChanged: onChanged,
      activeColor: AppColors.primary,
    );
  }
}

class _SecurityTipCard extends StatelessWidget {
  final IconData icon;
  final Color color;
  final String title;
  final String description;

  const _SecurityTipCard({
    required this.icon,
    required this.color,
    required this.title,
    required this.description,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withOpacity(0.3)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: color.withOpacity(0.2),
              shape: BoxShape.circle,
            ),
            child: Icon(icon, color: color, size: 20),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: AppTypography.titleSmall.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  description,
                  style: AppTypography.bodySmall.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _ContactOption extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  const _ContactOption({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Container(
        padding: const EdgeInsets.all(10),
        decoration: BoxDecoration(
          color: AppColors.primary.withOpacity(0.1),
          shape: BoxShape.circle,
        ),
        child: Icon(icon, color: AppColors.primary),
      ),
      title: Text(title),
      subtitle: Text(subtitle),
      onTap: onTap,
    );
  }
}
