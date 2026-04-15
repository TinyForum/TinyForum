'use client';

import { useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Notification, notificationApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import { useRouter } from 'next/navigation';
import { timeAgo } from '@/lib/utils';
import Image from 'next/image';
import { Bell, CheckCheck, Heart, MessageSquare, UserPlus, Info } from 'lucide-react';
import toast from 'react-hot-toast';
// import type { Notification } from '@/types';
import Avatar from '@/components/user/Avatar';
import { useTranslations } from 'next-intl';

const NotifIcon = ({ type }: { type: Notification['type'] }) => {
  const cls = 'w-4 h-4';
  switch (type) {
    case 'like': return <Heart className={`${cls} text-error`} />;
    case 'comment': return <MessageSquare className={`${cls} text-info`} />;
    case 'reply': return <MessageSquare className={`${cls} text-success`} />;
    case 'follow': return <UserPlus className={`${cls} text-primary`} />;
    default: return <Info className={`${cls} text-warning`} />;
  }
};

export default function NotificationsPage() {
  const { isAuthenticated } = useAuthStore();
  const router = useRouter();
  const queryClient = useQueryClient();
  const t = useTranslations('notifications');

  useEffect(() => {
    if (!isAuthenticated) router.push('/auth/login');
  }, [isAuthenticated, router]);

  const { data, isLoading } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => notificationApi.list({ page: 1, page_size: 50 }).then((r) => r.data.data),
    enabled: isAuthenticated,
  });

  const markAllMutation = useMutation({
    mutationFn: () => notificationApi.markAllRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
      queryClient.invalidateQueries({ queryKey: ['notifications', 'unread'] });
      toast.success('已全部标记为已读');
    },
  });

  const notifications = data?.list ?? [];
  const unread = notifications.filter((n) => !n.is_read).length;

  return (
    <div className="max-w-2xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold flex items-center gap-2">
          <Bell className="w-6 h-6 text-primary" />
          {t("title")}
          {unread > 0 && (
            <span className="badge badge-error badge-sm">{unread}</span>
          )}
        </h1>
        {unread > 0 && (
          <button
            className="btn btn-ghost btn-sm gap-1"
            onClick={() => markAllMutation.mutate()}
            disabled={markAllMutation.isPending}
          >
            <CheckCheck className="w-4 h-4" /> {t("mark_all_read")}
          </button>
        )}
      </div>

      {isLoading ? (
        <div className="space-y-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="skeleton h-16 w-full rounded-xl" />
          ))}
        </div>
      ) : notifications.length === 0 ? (
        <div className="text-center py-20 text-base-content/40">
          <Bell className="w-12 h-12 mx-auto mb-3 opacity-30" />
          <p>{t("no_notifications")}</p>
        </div>
      ) : (
        <div className="space-y-2">
          {notifications.map((notif) => (
            <div
              key={notif.id}
              className={`card border transition-colors ${
                notif.is_read ? 'bg-base-100 border-base-300' : 'bg-primary/5 border-primary/20'
              }`}
            >
              <div className="card-body p-4">
                <div className="flex items-start gap-3">
                  {notif.sender ? (
                    <div className="avatar flex-none">
                      <div className="w-9 h-9 rounded-full">

                         <Avatar 
  username={notif.sender.username} 
  avatarUrl={notif.sender.avatar}  // 数据库中的头像
  size="md" 
/>
                      </div>
                    </div>
                  ) : (
                    <div className="w-9 h-9 rounded-full bg-base-200 flex items-center justify-center flex-none">
                      <Bell className="w-4 h-4 text-base-content/40" />
                    </div>
                  )}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <NotifIcon type={notif.type} />
                      <p className="text-sm text-base-content/80">{notif.content}</p>
                      {!notif.is_read && (
                        <span className="w-2 h-2 bg-primary rounded-full flex-none ml-auto" />
                      )}
                    </div>
                    <p className="text-xs text-base-content/40 mt-1">{timeAgo(notif.created_at)}</p>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
