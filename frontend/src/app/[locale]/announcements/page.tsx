// app/[locale]/announcements/page.tsx
'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { announcementApi } from '@/lib/api';
import { toast } from 'react-hot-toast';
import {
  MegaphoneIcon,
  CalendarIcon,
  EyeIcon,
  PinIcon,
} from 'lucide-react';
import type { Announcement } from '@/lib/api/types';

// 公告卡片组件
function AnnouncementCard({ announcement }: { announcement: Announcement }) {
  const getTypeColor = (type: string) => {
    switch (type) {
      case 'system':
        return 'bg-blue-100 text-blue-700';
      case 'feature':
        return 'bg-green-100 text-green-700';
      case 'maintenance':
        return 'bg-yellow-100 text-yellow-700';
      case 'policy':
        return 'bg-purple-100 text-purple-700';
      default:
        return 'bg-gray-100 text-gray-700';
    }
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'system':
        return '系统';
      case 'feature':
        return '功能';
      case 'maintenance':
        return '维护';
      case 'policy':
        return '政策';
      default:
        return '公告';
    }
  };

  return (
    <Link href={`/announcements/${announcement.id}`}>
      <div className="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <MegaphoneIcon className="w-4 h-4 text-indigo-500" />
              <span className={`text-xs px-2 py-0.5 rounded-full ${getTypeColor(announcement.type)}`}>
                {getTypeLabel(announcement.type)}
              </span>
              {announcement.is_pinned && (
                <span className="text-xs px-2 py-0.5 rounded-full bg-orange-100 text-orange-700 flex items-center gap-1">
                  <PinIcon className="w-3 h-3" />
                  置顶
                </span>
              )}
            </div>
            <h3 className="font-semibold text-gray-900 mb-2 line-clamp-1">
              {announcement.title}
            </h3>
            <p className="text-sm text-gray-500 line-clamp-2 mb-3">
              {announcement.content?.replace(/<[^>]*>/g, '')}
            </p>
            <div className="flex items-center gap-4 text-xs text-gray-400">
              <div className="flex items-center gap-1">
                <CalendarIcon className="w-3 h-3" />
                {new Date(announcement.created_at).toLocaleDateString()}
              </div>
              <div className="flex items-center gap-1">
                <EyeIcon className="w-3 h-3" />
                {announcement.view_count || 0}
              </div>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}

export default function AnnouncementsPage() {
  const [announcements, setAnnouncements] = useState<Announcement[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 20;

  const loadAnnouncements = async () => {
    setLoading(true);
    try {
      const response = await announcementApi.list({ page, page_size: pageSize });
      if (response.data.code === 200 || response.data.code === 0) {
        setAnnouncements(response.data.data.items || []);
        setTotal(response.data.data.total || 0);
      } else {
        toast.error(response.data.message || '加载失败');
      }
    } catch (error: any) {
      console.error('Failed to load announcements:', error);
      toast.error(error.response?.data?.message || '加载失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAnnouncements();
  }, [page]);

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-3xl mx-auto px-4">
        {/* 头部 */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
          <div className="flex items-center gap-3 mb-2">
            <MegaphoneIcon className="w-8 h-8 text-indigo-500" />
            <h1 className="text-2xl font-bold text-gray-900">公告</h1>
          </div>
          <p className="text-gray-500">了解最新的平台动态和重要通知</p>
        </div>

        {/* 公告列表 */}
        {loading ? (
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-white rounded-lg shadow-sm p-4 animate-pulse">
                <div className="h-5 bg-gray-200 rounded w-1/3 mb-2" />
                <div className="h-4 bg-gray-200 rounded w-full mb-2" />
                <div className="h-4 bg-gray-200 rounded w-2/3" />
              </div>
            ))}
          </div>
        ) : announcements.length === 0 ? (
          <div className="bg-white rounded-lg shadow-sm p-12 text-center">
            <div className="text-6xl mb-4">📢</div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">暂无公告</h3>
            <p className="text-gray-500">稍后再来看看吧</p>
          </div>
        ) : (
          <>
            <div className="space-y-3">
              {/* 置顶公告优先显示 */}
              {announcements
                .filter(a => a.is_pinned)
                .map((announcement) => (
                  <AnnouncementCard key={announcement.id} announcement={announcement} />
                ))}
              {announcements
                .filter(a => !a.is_pinned)
                .map((announcement) => (
                  <AnnouncementCard key={announcement.id} announcement={announcement} />
                ))}
            </div>

            {/* 分页 */}
            {totalPages > 1 && (
              <div className="flex justify-center gap-2 mt-6">
                <button
                  onClick={() => setPage(p => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
                >
                  上一页
                </button>
                <span className="px-3 py-1 text-gray-600">
                  第 {page} / {totalPages} 页
                </span>
                <button
                  onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                  disabled={page >= totalPages}
                  className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
                >
                  下一页
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}