import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:timeago/timeago.dart' as timeago;

import '../../../../core/constants/app_constants.dart';

class NotificationsScreen extends ConsumerStatefulWidget {
  const NotificationsScreen({super.key});

  @override
  ConsumerState<NotificationsScreen> createState() =>
      _NotificationsScreenState();
}

class _NotificationsScreenState extends ConsumerState<NotificationsScreen> {
  String _selectedFilter = 'All';

  final List<_Notification> _notifications = [
    _Notification(
      id: '1',
      type: 'gig',
      title: 'New Proposal Received',
      body: 'John Doe submitted a proposal for your "Mobile App Development" gig.',
      time: DateTime.now().subtract(const Duration(minutes: 5)),
      isRead: false,
    ),
    _Notification(
      id: '2',
      type: 'savings',
      title: 'Contribution Reminder',
      body: 'Your ₦25,000 contribution to "Monthly Savings Club" is due tomorrow.',
      time: DateTime.now().subtract(const Duration(hours: 1)),
      isRead: false,
    ),
    _Notification(
      id: '3',
      type: 'wallet',
      title: 'Deposit Successful',
      body: '₦50,000 has been added to your wallet.',
      time: DateTime.now().subtract(const Duration(hours: 3)),
      isRead: true,
    ),
    _Notification(
      id: '4',
      type: 'credit',
      title: 'Credit Score Updated',
      body: 'Your credit score increased by 15 points! Keep up the good work.',
      time: DateTime.now().subtract(const Duration(hours: 5)),
      isRead: true,
    ),
    _Notification(
      id: '5',
      type: 'loan',
      title: 'Loan Payment Due',
      body: 'Your loan repayment of ₦6,875 is due in 3 days.',
      time: DateTime.now().subtract(const Duration(days: 1)),
      isRead: false,
    ),
    _Notification(
      id: '6',
      type: 'gig',
      title: 'Gig Completed',
      body: 'Payment of ₦75,000 for "Logo Design" has been released to your wallet.',
      time: DateTime.now().subtract(const Duration(days: 1)),
      isRead: true,
    ),
    _Notification(
      id: '7',
      type: 'savings',
      title: 'Payout Received! 🎉',
      body: 'You received ₦250,000 payout from "Monthly Savings Club".',
      time: DateTime.now().subtract(const Duration(days: 2)),
      isRead: true,
    ),
    _Notification(
      id: '8',
      type: 'system',
      title: 'Account Verified',
      body: 'Your identity has been verified successfully. You now have full access to all features.',
      time: DateTime.now().subtract(const Duration(days: 3)),
      isRead: true,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final filteredNotifications = _selectedFilter == 'All'
        ? _notifications
        : _selectedFilter == 'Unread'
            ? _notifications.where((n) => !n.isRead).toList()
            : _notifications.where((n) => n.type == _selectedFilter.toLowerCase()).toList();

    final unreadCount = _notifications.where((n) => !n.isRead).length;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Notifications'),
        actions: [
          if (unreadCount > 0)
            TextButton(
              onPressed: _markAllAsRead,
              child: const Text('Mark all read'),
            ),
        ],
      ),
      body: Column(
        children: [
          // Filter chips
          SizedBox(
            height: 50,
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              children: ['All', 'Unread', 'Gig', 'Savings', 'Wallet', 'Credit']
                  .map((filter) => Padding(
                        padding: const EdgeInsets.only(right: 8),
                        child: FilterChip(
                          label: Text(filter),
                          selected: _selectedFilter == filter,
                          onSelected: (selected) {
                            setState(() => _selectedFilter = filter);
                          },
                          selectedColor: AppColors.primary.withOpacity(0.2),
                          backgroundColor: AppColors.surfaceVariant,
                          labelStyle: AppTypography.labelMedium.copyWith(
                            color: _selectedFilter == filter
                                ? AppColors.primary
                                : AppColors.textSecondary,
                          ),
                        ),
                      ))
                  .toList(),
            ),
          ),

          // Notifications list
          Expanded(
            child: filteredNotifications.isEmpty
                ? _buildEmptyState()
                : ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: filteredNotifications.length,
                    itemBuilder: (context, index) {
                      final notification = filteredNotifications[index];
                      return _buildNotificationCard(notification);
                    },
                  ),
          ),
        ],
      ),
    );
  }

  Widget _buildNotificationCard(_Notification notification) {
    return Dismissible(
      key: Key(notification.id),
      direction: DismissDirection.endToStart,
      background: Container(
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: 20),
        color: AppColors.error,
        child: const Icon(Icons.delete_outline, color: Colors.white),
      ),
      onDismissed: (_) {
        setState(() {
          _notifications.removeWhere((n) => n.id == notification.id);
        });
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Notification deleted')),
        );
      },
      child: GestureDetector(
        onTap: () => _handleNotificationTap(notification),
        child: Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: notification.isRead
                ? AppColors.surface
                : AppColors.primary.withOpacity(0.05),
            borderRadius: BorderRadius.circular(16),
            border: Border.all(
              color: notification.isRead
                  ? AppColors.border
                  : AppColors.primary.withOpacity(0.2),
            ),
          ),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildNotificationIcon(notification.type),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            notification.title,
                            style: AppTypography.titleSmall.copyWith(
                              color: AppColors.textPrimary,
                              fontWeight: notification.isRead
                                  ? FontWeight.w500
                                  : FontWeight.w600,
                            ),
                          ),
                        ),
                        if (!notification.isRead)
                          Container(
                            width: 8,
                            height: 8,
                            decoration: const BoxDecoration(
                              color: AppColors.primary,
                              shape: BoxShape.circle,
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text(
                      notification.body,
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      timeago.format(notification.time),
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.textTertiary,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildNotificationIcon(String type) {
    IconData icon;
    Color color;

    switch (type) {
      case 'gig':
        icon = Icons.work_rounded;
        color = AppColors.primary;
        break;
      case 'savings':
        icon = Icons.savings_rounded;
        color = AppColors.secondary;
        break;
      case 'wallet':
        icon = Icons.account_balance_wallet_rounded;
        color = AppColors.info;
        break;
      case 'credit':
        icon = Icons.credit_score_rounded;
        color = AppColors.credit;
        break;
      case 'loan':
        icon = Icons.monetization_on_rounded;
        color = AppColors.warning;
        break;
      case 'system':
      default:
        icon = Icons.notifications_rounded;
        color = AppColors.textSecondary;
    }

    return Container(
      width: 44,
      height: 44,
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Icon(icon, color: color, size: 22),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.notifications_off_outlined,
            size: 64,
            color: AppColors.textTertiary,
          ),
          const SizedBox(height: 16),
          Text(
            'No notifications',
            style: AppTypography.titleMedium.copyWith(
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'You\'re all caught up!',
            style: AppTypography.bodyMedium.copyWith(
              color: AppColors.textTertiary,
            ),
          ),
        ],
      ),
    );
  }

  void _markAllAsRead() {
    setState(() {
      for (var notification in _notifications) {
        notification.isRead = true;
      }
    });
  }

  void _handleNotificationTap(_Notification notification) {
    setState(() {
      notification.isRead = true;
    });

    // Navigate based on type
    // switch (notification.type) {
    //   case 'gig':
    //     context.push('/gigs/${notification.data['gig_id']}');
    //     break;
    //   case 'savings':
    //     context.push('/savings/${notification.data['circle_id']}');
    //     break;
    //   // etc.
    // }
  }
}

class _Notification {
  final String id;
  final String type;
  final String title;
  final String body;
  final DateTime time;
  bool isRead;
  final Map<String, dynamic>? data;

  _Notification({
    required this.id,
    required this.type,
    required this.title,
    required this.body,
    required this.time,
    required this.isRead,
    this.data,
  });
}
