import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/constants/app_constants.dart';
import '../../../../core/providers/auth_provider.dart';
import '../../../../router/app_router.dart';

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(currentUserProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Profile'),
        actions: [
          IconButton(
            onPressed: () {},
            icon: const Icon(Icons.settings_outlined),
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            // Profile header
            _buildProfileHeader(context, user),

            const SizedBox(height: 24),

            // Stats
            _buildStatsRow(),

            const SizedBox(height: 24),

            // Menu sections
            _buildSection(
              'Account',
              [
                _MenuItem(
                  Icons.person_outline_rounded,
                  'Edit Profile',
                  () {},
                ),
                _MenuItem(
                  Icons.lock_outline_rounded,
                  'Change PIN',
                  () {},
                ),
                _MenuItem(
                  Icons.verified_user_outlined,
                  'Verify Identity',
                  () {},
                  badge: 'Pending',
                ),
                _MenuItem(
                  Icons.account_balance_outlined,
                  'Bank Accounts',
                  () {},
                ),
              ],
            ),

            const SizedBox(height: 16),

            _buildSection(
              'Preferences',
              [
                _MenuItem(
                  Icons.notifications_outlined,
                  'Notifications',
                  () => context.push(AppRoutes.notifications),
                ),
                _MenuItem(
                  Icons.language_rounded,
                  'Language',
                  () {},
                  trailing: 'English',
                ),
                _MenuItem(
                  Icons.dark_mode_outlined,
                  'Theme',
                  () {},
                  trailing: 'System',
                ),
                _MenuItem(
                  Icons.fingerprint_rounded,
                  'Biometric Login',
                  () {},
                  isToggle: true,
                  toggleValue: true,
                ),
              ],
            ),

            const SizedBox(height: 16),

            _buildSection(
              'Support',
              [
                _MenuItem(
                  Icons.help_outline_rounded,
                  'Help Center',
                  () {},
                ),
                _MenuItem(
                  Icons.chat_bubble_outline_rounded,
                  'Contact Support',
                  () {},
                ),
                _MenuItem(
                  Icons.description_outlined,
                  'Terms & Conditions',
                  () {},
                ),
                _MenuItem(
                  Icons.privacy_tip_outlined,
                  'Privacy Policy',
                  () {},
                ),
              ],
            ),

            const SizedBox(height: 16),

            _buildSection(
              '',
              [
                _MenuItem(
                  Icons.share_outlined,
                  'Invite Friends',
                  () {},
                  subtitle: 'Earn ₦500 per referral',
                ),
                _MenuItem(
                  Icons.star_outline_rounded,
                  'Rate Us',
                  () {},
                ),
              ],
            ),

            const SizedBox(height: 24),

            // Logout button
            SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: () => _showLogoutDialog(context, ref),
                icon: const Icon(Icons.logout_rounded, color: AppColors.error),
                label: Text(
                  'Log Out',
                  style: AppTypography.button.copyWith(color: AppColors.error),
                ),
                style: OutlinedButton.styleFrom(
                  side: const BorderSide(color: AppColors.error),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
              ),
            ),

            const SizedBox(height: 16),

            // App version
            Text(
              'Version ${AppConstants.appVersion}',
              style: AppTypography.labelSmall.copyWith(
                color: AppColors.textTertiary,
              ),
            ),

            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Widget _buildProfileHeader(BuildContext context, User? user) {
    return Column(
      children: [
        Stack(
          children: [
            CircleAvatar(
              radius: 50,
              backgroundColor: AppColors.primary,
              backgroundImage: user?.avatar != null
                  ? NetworkImage(user!.avatar!)
                  : null,
              child: user?.avatar == null
                  ? Text(
                      user?.firstName.isNotEmpty == true
                          ? user!.firstName[0].toUpperCase()
                          : 'U',
                      style: AppTypography.displaySmall.copyWith(
                        color: Colors.white,
                      ),
                    )
                  : null,
            ),
            Positioned(
              bottom: 0,
              right: 0,
              child: Container(
                width: 32,
                height: 32,
                decoration: BoxDecoration(
                  color: AppColors.secondary,
                  shape: BoxShape.circle,
                  border: Border.all(color: Colors.white, width: 2),
                ),
                child: const Icon(
                  Icons.camera_alt_rounded,
                  color: Colors.white,
                  size: 16,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        Text(
          user?.fullName ?? 'User Name',
          style: AppTypography.headlineSmall.copyWith(
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          user?.phone ?? '+234 XXX XXX XXXX',
          style: AppTypography.bodyMedium.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: 8),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              decoration: BoxDecoration(
                color: AppColors.successLight,
                borderRadius: BorderRadius.circular(20),
              ),
              child: Row(
                children: [
                  const Icon(
                    Icons.verified_rounded,
                    color: AppColors.success,
                    size: 16,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    'Verified',
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.success,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 8),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              decoration: BoxDecoration(
                color: AppColors.credit.withOpacity(0.1),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Row(
                children: [
                  const Icon(
                    Icons.star_rounded,
                    color: AppColors.credit,
                    size: 16,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    'Silver Tier',
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.credit,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildStatsRow() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Expanded(
            child: _buildStatItem('Gigs Completed', '12'),
          ),
          Container(
            width: 1,
            height: 40,
            color: AppColors.border,
          ),
          Expanded(
            child: _buildStatItem('Savings Circles', '3'),
          ),
          Container(
            width: 1,
            height: 40,
            color: AppColors.border,
          ),
          Expanded(
            child: _buildStatItem('Credit Score', '680'),
          ),
        ],
      ),
    );
  }

  Widget _buildStatItem(String label, String value) {
    return Column(
      children: [
        Text(
          value,
          style: AppTypography.titleLarge.copyWith(
            color: AppColors.primary,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: AppTypography.labelSmall.copyWith(
            color: AppColors.textSecondary,
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }

  Widget _buildSection(String title, List<_MenuItem> items) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (title.isNotEmpty) ...[
          Text(
            title,
            style: AppTypography.titleSmall.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 8),
        ],
        Container(
          decoration: BoxDecoration(
            color: AppColors.surface,
            borderRadius: BorderRadius.circular(16),
            border: Border.all(color: AppColors.border),
          ),
          child: Column(
            children: items.asMap().entries.map((entry) {
              final index = entry.key;
              final item = entry.value;
              return Column(
                children: [
                  _buildMenuItem(item),
                  if (index < items.length - 1)
                    const Divider(height: 1, indent: 56),
                ],
              );
            }).toList(),
          ),
        ),
      ],
    );
  }

  Widget _buildMenuItem(_MenuItem item) {
    return ListTile(
      leading: Icon(item.icon, color: AppColors.textSecondary),
      title: Text(
        item.title,
        style: AppTypography.bodyMedium.copyWith(
          color: AppColors.textPrimary,
        ),
      ),
      subtitle: item.subtitle != null
          ? Text(
              item.subtitle!,
              style: AppTypography.labelSmall.copyWith(
                color: AppColors.secondary,
              ),
            )
          : null,
      trailing: item.isToggle
          ? Switch(
              value: item.toggleValue ?? false,
              onChanged: (_) {},
              activeColor: AppColors.primary,
            )
          : item.badge != null
              ? Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 8,
                    vertical: 4,
                  ),
                  decoration: BoxDecoration(
                    color: AppColors.warningLight,
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Text(
                    item.badge!,
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.warning,
                    ),
                  ),
                )
              : item.trailing != null
                  ? Text(
                      item.trailing!,
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textTertiary,
                      ),
                    )
                  : const Icon(
                      Icons.chevron_right_rounded,
                      color: AppColors.textTertiary,
                    ),
      onTap: item.onTap,
    );
  }

  void _showLogoutDialog(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Log Out'),
        content: const Text('Are you sure you want to log out?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(context);
              await ref.read(authStateProvider.notifier).logout();
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.error,
            ),
            child: const Text('Log Out'),
          ),
        ],
      ),
    );
  }
}

class _MenuItem {
  final IconData icon;
  final String title;
  final VoidCallback onTap;
  final String? subtitle;
  final String? badge;
  final String? trailing;
  final bool isToggle;
  final bool? toggleValue;

  _MenuItem(
    this.icon,
    this.title,
    this.onTap, {
    this.subtitle,
    this.badge,
    this.trailing,
    this.isToggle = false,
    this.toggleValue,
  });
}
