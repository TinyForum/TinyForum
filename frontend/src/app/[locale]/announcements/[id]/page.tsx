'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { announcementApi } from '@/lib/api';
import { toast } from 'react-hot-toast';
import {
  MegaphoneIcon,
  CalendarIcon,
  EyeIcon,
  PinIcon,
  ArrowLeftIcon,
  TagIcon,
} from 'lucide-react';
import type { Announcement } from '@/lib/api/modules/announcements';

// 公告类型配置
const TYPE_CONFIG: Record<string, { color: string; label: string; bgColor: string }> = {
  normal: { color: 'text-blue-700', label: '普通公告', bgColor: 'bg-blue-50' },
  important: { color: 'text-orange-700', label: '重要公告', bgColor: 'bg-orange-50' },
  emergency: { color: 'text-red-700', label: '紧急公告', bgColor: 'bg-red-50' },
  event: { color: 'text-green-700', label: '活动公告', bgColor: 'bg-green-50' },
};

export default function AnnouncementDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [announcement, setAnnouncement] = useState<Announcement | null>(null);
  const [loading, setLoading] = useState(true);

  const loadAnnouncement = async () => {
    const id = parseInt(params.id as string);
    if (isNaN(id)) {
      toast.error('无效的公告ID');
      router.push('/announcements');
      return;
    }

    setLoading(true);
    try {
      const response = await announcementApi.getById(id);
      if (response.data.code === 0) {
        setAnnouncement(response.data.data);
      } else {
        toast.error(response.data.message || '公告不存在');
        router.push('/announcements');
      }
    } catch (error: any) {
      console.error('Failed to load announcement:', error);
      if (error.response?.status === 404) {
        toast.error('公告不存在');
        router.push('/announcements');
      } else {
        toast.error('加载失败，请稍后重试');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAnnouncement();
  }, [params.id]);

  const formatDate = (dateStr: string | null) => {
    if (!dateStr) return '待发布';
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const typeConfig = announcement ? TYPE_CONFIG[announcement.type] || TYPE_CONFIG.normal : null;

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-3xl mx-auto px-4">
          <div className="bg-white rounded-lg shadow-sm p-8 animate-pulse">
            <div className="h-8 bg-gray-200 rounded w-3/4 mb-4" />
            <div className="h-4 bg-gray-200 rounded w-1/2 mb-6" />
            <div className="space-y-3">
              <div className="h-4 bg-gray-200 rounded w-full" />
              <div className="h-4 bg-gray-200 rounded w-full" />
              <div className="h-4 bg-gray-200 rounded w-2/3" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!announcement) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-3xl mx-auto px-4">
        {/* 返回按钮 */}
        <Link
          href="/announcements"
          className="inline-flex items-center gap-2 text-gray-500 hover:text-gray-700 mb-4 transition-colors"
        >
          <ArrowLeftIcon className="w-4 h-4" />
          返回公告列表
        </Link>

        {/* 公告内容 */}
        <article className="bg-white rounded-lg shadow-sm overflow-hidden">
          {/* 头部 */}
          <div className="p-6 border-b border-gray-100">
            <div className="flex items-center gap-2 mb-4 flex-wrap">
              {typeConfig && (
                <span className={`text-xs px-2 py-0.5 rounded-full ${typeConfig.bgColor} ${typeConfig.color}`}>
                  {typeConfig.label}
                </span>
              )}
              {announcement.is_pinned && (
                <span className="text-xs px-2 py-0.5 rounded-full bg-orange-100 text-orange-700 flex items-center gap-1">
                  <PinIcon className="w-3 h-3" />
                  置顶
                </span>
              )}
              {announcement.board && (
                <span className="text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 flex items-center gap-1">
                  <TagIcon className="w-3 h-3" />
                  {announcement.board.name}
                </span>
              )}
            </div>
            
            <h1 className="text-2xl font-bold text-gray-900 mb-4">
              {announcement.title}
            </h1>
            
            <div className="flex items-center gap-4 text-sm text-gray-400">
              <div className="flex items-center gap-1">
                <CalendarIcon className="w-4 h-4" />
                发布时间：{formatDate(announcement.published_at || announcement.created_at)}
              </div>
              <div className="flex items-center gap-1">
                <EyeIcon className="w-4 h-4" />
                浏览 {announcement.view_count || 0} 次
              </div>
            </div>
          </div>

          {/* 内容 */}
          <div className="p-6">
            <div 
              className="prose prose-gray max-w-none"
              dangerouslySetInnerHTML={{ __html: announcement.content }}
            />
          </div>

          {/* 页脚 */}
          <div className="p-6 bg-gray-50 border-t border-gray-100">
            <div className="text-xs text-gray-400 text-center">
              {announcement.created_by && (
                <p>发布者 ID: {announcement.created_by}</p>
              )}
              {announcement.expired_at && new Date(announcement.expired_at) < new Date() && (
                <p className="text-orange-500 mt-1">此公告已过期</p>
              )}
            </div>
          </div>
        </article>
      </div>
    </div>
  );
}