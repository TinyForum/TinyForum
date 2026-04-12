// app/[locale]/announcements/[id]/page.tsx
'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuthStore } from '@/store/auth';
import { toast } from 'react-hot-toast';
import {
  ArrowLeftIcon,
  CalendarIcon,
  EyeIcon,
  UserCircleIcon,
  MegaphoneIcon,
  DocumentTextIcon,
  ShareIcon,
  FlagIcon,
} from '@heroicons/react/24/outline';
import { announcementApi } from '@/lib/api';
import type { Announcement } from '@/lib/api/types';

// 公告详情组件
function AnnouncementDetail({ announcement }: { announcement: Announcement }) {
  const [showShareMenu, setShowShareMenu] = useState(false);

  const handleShare = async () => {
    try {
      await navigator.clipboard.writeText(window.location.href);
      toast.success('链接已复制到剪贴板');
      setShowShareMenu(false);
    } catch (error) {
      toast.error('复制失败');
    }
  };

  const getAnnouncementTypeIcon = (type: string) => {
    switch (type) {
      case 'system':
        return '🔧';
      case 'feature':
        return '✨';
      case 'maintenance':
        return '🔨';
      case 'policy':
        return '📋';
      default:
        return '📢';
    }
  };

  const getAnnouncementTypeColor = (type: string) => {
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

  const getAnnouncementTypeLabel = (type: string) => {
    switch (type) {
      case 'system':
        return '系统公告';
      case 'feature':
        return '功能更新';
      case 'maintenance':
        return '维护通知';
      case 'policy':
        return '政策变更';
      default:
        return '公告';
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm overflow-hidden">
      {/* 头部 */}
      <div className="p-6 border-b">
        <div className="flex items-center gap-3 mb-4 flex-wrap">
          <MegaphoneIcon className="w-8 h-8 text-indigo-500" />
          <div className={`px-3 py-1 rounded-full text-sm font-medium ${getAnnouncementTypeColor(announcement.type)}`}>
            {getAnnouncementTypeIcon(announcement.type)} {getAnnouncementTypeLabel(announcement.type)}
          </div>
          {announcement.is_pinned && (
            <div className="px-3 py-1 rounded-full text-sm font-medium bg-orange-100 text-orange-700">
              📌 置顶
            </div>
          )}
        </div>

        <h1 className="text-2xl font-bold text-gray-900 mb-4">
          {announcement.title}
        </h1>

        <div className="flex flex-wrap items-center gap-4 text-sm text-gray-500">
          <div className="flex items-center gap-1">
            <UserCircleIcon className="w-4 h-4" />
            <Link href={`/users/${announcement.author_id}`} className="hover:text-indigo-600">
              {announcement.author?.username || `用户${announcement.author_id}`}
            </Link>
          </div>
          <div className="flex items-center gap-1">
            <CalendarIcon className="w-4 h-4" />
            {new Date(announcement.created_at).toLocaleDateString('zh-CN', {
              year: 'numeric',
              month: 'long',
              day: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
            })}
          </div>
          <div className="flex items-center gap-1">
            <EyeIcon className="w-4 h-4" />
            {announcement.view_count || 0} 次阅读
          </div>
          <div className="flex items-center gap-1">
            <DocumentTextIcon className="w-4 h-4" />
            {announcement.content?.length || 0} 字
          </div>
        </div>
      </div>

      {/* 内容 */}
      <div className="p-6">
        <div
          className="prose max-w-none"
          dangerouslySetInnerHTML={{ __html: announcement.content }}
        />
      </div>

      {/* 操作按钮 */}
      <div className="px-6 pb-6 flex gap-2">
        <div className="relative">
          <button
            onClick={() => setShowShareMenu(!showShareMenu)}
            className="flex items-center gap-1 px-3 py-1 text-gray-500 hover:text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <ShareIcon className="w-4 h-4" />
            分享
          </button>
          {showShareMenu && (
            <>
              <div
                className="fixed inset-0 z-40"
                onClick={() => setShowShareMenu(false)}
              />
              <div className="absolute top-full left-0 mt-1 bg-white rounded-lg shadow-lg border p-2 z-50 min-w-[120px]">
                <button
                  onClick={handleShare}
                  className="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-md w-full"
                >
                  📋 复制链接
                </button>
              </div>
            </>
          )}
        </div>
        <button className="flex items-center gap-1 px-3 py-1 text-gray-500 hover:text-red-500 rounded-lg hover:bg-gray-100 transition-colors">
          <FlagIcon className="w-4 h-4" />
          举报
        </button>
      </div>
    </div>
  );
}

// 加载骨架屏
function LoadingSkeleton() {
  return (
    <div className="bg-white rounded-lg shadow-sm overflow-hidden animate-pulse">
      <div className="p-6 border-b">
        <div className="flex items-center gap-3 mb-4">
          <div className="w-8 h-8 bg-gray-200 rounded-full" />
          <div className="w-20 h-6 bg-gray-200 rounded-full" />
        </div>
        <div className="h-8 bg-gray-200 rounded w-3/4 mb-4" />
        <div className="flex gap-4">
          <div className="h-4 bg-gray-200 rounded w-24" />
          <div className="h-4 bg-gray-200 rounded w-32" />
          <div className="h-4 bg-gray-200 rounded w-20" />
        </div>
      </div>
      <div className="p-6 space-y-3">
        <div className="h-4 bg-gray-200 rounded w-full" />
        <div className="h-4 bg-gray-200 rounded w-11/12" />
        <div className="h-4 bg-gray-200 rounded w-10/12" />
        <div className="h-4 bg-gray-200 rounded w-9/12" />
      </div>
    </div>
  );
}

// 错误状态组件
function ErrorState({ message, onRetry }: { message: string; onRetry: () => void }) {
  return (
    <div className="bg-white rounded-lg shadow-sm p-12 text-center">
      <div className="text-6xl mb-4">⚠️</div>
      <h3 className="text-lg font-medium text-gray-900 mb-2">加载失败</h3>
      <p className="text-gray-500 mb-4">{message}</p>
      <button
        onClick={onRetry}
        className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
      >
        重新加载
      </button>
    </div>
  );
}

// 公告不存在状态
function NotFoundState() {
  return (
    <div className="bg-white rounded-lg shadow-sm p-12 text-center">
      <div className="text-6xl mb-4">📭</div>
      <h3 className="text-lg font-medium text-gray-900 mb-2">公告不存在</h3>
      <p className="text-gray-500 mb-4">该公告可能已被删除或不存在</p>
      <Link
        href="/announcements"
        className="inline-block px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
      >
        返回公告列表
      </Link>
    </div>
  );
}

export default function AnnouncementDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [announcement, setAnnouncement] = useState<Announcement | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const id = Number(params.id);

  // 加载公告详情
  const loadAnnouncement = async () => {
    if (!id || isNaN(id)) {
      setError('无效的公告ID');
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);
    
    try {
      const response = await announcementApi.getById(id);
      
      if (response.data.code === 200 || response.data.code === 0) {
        setAnnouncement(response.data.data);
      } else if (response.data.code === 404) {
        setAnnouncement(null);
      } else {
        throw new Error(response.data.message || '加载失败');
      }
    } catch (error: any) {
      console.error('Failed to load announcement:', error);
      const errorMsg = error.response?.data?.message || error.message || '加载公告失败';
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAnnouncement();
  }, [id]);

  // 页面标题
  useEffect(() => {
    if (announcement) {
      document.title = `${announcement.title} - 公告`;
    } else if (!loading) {
      document.title = '公告详情';
    }
  }, [announcement, loading]);

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* 返回按钮 */}
        <div className="mb-4">
          <Link
            href="/announcements"
            className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            返回公告列表
          </Link>
        </div>

        {/* 内容区域 */}
        {loading ? (
          <LoadingSkeleton />
        ) : error ? (
          <ErrorState message={error} onRetry={loadAnnouncement} />
        ) : !announcement ? (
          <NotFoundState />
        ) : (
          <AnnouncementDetail announcement={announcement} />
        )}
      </div>
    </div>
  );
}