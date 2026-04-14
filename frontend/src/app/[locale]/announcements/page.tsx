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
  FileTextIcon,
} from 'lucide-react';
import type { Announcement } from '@/lib/api/modules/announcements';

// 公告类型配置
const TYPE_CONFIG: Record<string, { color: string; label: string }> = {
  normal: { color: 'bg-blue-100 text-blue-700', label: '普通' },
  important: { color: 'bg-orange-100 text-orange-700', label: '重要' },
  emergency: { color: 'bg-red-100 text-red-700', label: '紧急' },
  event: { color: 'bg-green-100 text-green-700', label: '活动' },
};

// 获取公告类型样式
function getTypeConfig(type: string) {
  return TYPE_CONFIG[type] || TYPE_CONFIG.normal;
}

// 公告卡片组件
function AnnouncementCard({ announcement }: { announcement: Announcement }) {
  const typeConfig = getTypeConfig(announcement.type);
  
  // 格式化时间
  const formatDate = (dateStr: string | null) => {
    if (!dateStr) return '待发布';
    return new Date(dateStr).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    });
  };

  return (
    <Link href={`/announcements/${announcement.id}`}>
      <div className="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer border border-gray-100">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2 flex-wrap">
              <MegaphoneIcon className="w-4 h-4 text-indigo-500" />
              <span className={`text-xs px-2 py-0.5 rounded-full ${typeConfig.color}`}>
                {typeConfig.label}
              </span>
              {announcement.is_pinned && (
                <span className="text-xs px-2 py-0.5 rounded-full bg-orange-100 text-orange-700 flex items-center gap-1">
                  <PinIcon className="w-3 h-3" />
                  置顶
                </span>
              )}
              {announcement.status === 'draft' && (
                <span className="text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600">
                  草稿
                </span>
              )}
            </div>
            <h3 className="font-semibold text-gray-900 mb-2 line-clamp-1">
              {announcement.title}
            </h3>
            <p className="text-sm text-gray-500 line-clamp-2 mb-3">
              {announcement.summary || announcement.content?.replace(/<[^>]*>/g, '').slice(0, 150)}
            </p>
            <div className="flex items-center gap-4 text-xs text-gray-400">
              <div className="flex items-center gap-1">
                <CalendarIcon className="w-3 h-3" />
                {formatDate(announcement.published_at || announcement.created_at)}
              </div>
              <div className="flex items-center gap-1">
                <EyeIcon className="w-3 h-3" />
                {announcement.view_count || 0} 次浏览
              </div>
              {announcement.board && (
                <div className="flex items-center gap-1">
                  <FileTextIcon className="w-3 h-3" />
                  {announcement.board.name}
                </div>
              )}
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
      const response = await announcementApi.list({ 
        page, 
        page_size: pageSize,
        status: 'published', // 只显示已发布的公告
      });
      
      // 后端返回格式: { code: 0, message: "success", data: { list, total, page, page_size } }
      if (response.data.code === 0) {
        setAnnouncements(response.data.data.list || []);
        setTotal(response.data.data.total || 0);
      } else {
        toast.error(response.data.message || '加载失败');
      }
    } catch (error: any) {
      console.error('Failed to load announcements:', error);
      toast.error(error.response?.data?.message || '加载失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAnnouncements();
  }, [page]);

  const totalPages = Math.ceil(total / pageSize);
  
  // 分离置顶和普通公告
  const pinnedAnnouncements = announcements.filter(a => a.is_pinned);
  const normalAnnouncements = announcements.filter(a => !a.is_pinned);

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
                <div className="flex items-center gap-2 mb-2">
                  <div className="h-4 w-4 bg-gray-200 rounded" />
                  <div className="h-4 w-16 bg-gray-200 rounded" />
                </div>
                <div className="h-5 bg-gray-200 rounded w-3/4 mb-2" />
                <div className="h-4 bg-gray-200 rounded w-full mb-2" />
                <div className="h-4 bg-gray-200 rounded w-2/3" />
              </div>
            ))}
          </div>
        ) : announcements.length === 0 ? (
          <div className="bg-white rounded-lg shadow-sm p-12 text-center">
            <div className="text-6xl mb-4">📢</div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">暂无公告</h3>
            <p className="text-gray-500">暂时没有公告，稍后再来看看吧</p>
          </div>
        ) : (
          <>
            <div className="space-y-3">
              {/* 置顶公告 */}
              {pinnedAnnouncements.length > 0 && (
                <div className="mb-4">
                  <div className="flex items-center gap-2 mb-3">
                    <PinIcon className="w-4 h-4 text-orange-500" />
                    <h2 className="text-sm font-medium text-gray-600">置顶公告</h2>
                  </div>
                  {pinnedAnnouncements.map((announcement) => (
                    <AnnouncementCard key={announcement.id} announcement={announcement} />
                  ))}
                </div>
              )}
              
              {/* 普通公告 */}
              {normalAnnouncements.length > 0 && (
                <div>
                  {pinnedAnnouncements.length > 0 && (
                    <div className="flex items-center gap-2 mb-3">
                      <FileTextIcon className="w-4 h-4 text-gray-400" />
                      <h2 className="text-sm font-medium text-gray-600">最新公告</h2>
                    </div>
                  )}
                  {normalAnnouncements.map((announcement) => (
                    <AnnouncementCard key={announcement.id} announcement={announcement} />
                  ))}
                </div>
              )}
            </div>

            {/* 分页 */}
            {totalPages > 1 && (
              <div className="flex justify-center gap-2 mt-8">
                <button
                  onClick={() => setPage(p => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="px-4 py-2 border rounded-md text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
                >
                  上一页
                </button>
                <div className="flex items-center gap-1">
                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    let pageNum: number;
                    if (totalPages <= 5) {
                      pageNum = i + 1;
                    } else if (page <= 3) {
                      pageNum = i + 1;
                    } else if (page >= totalPages - 2) {
                      pageNum = totalPages - 4 + i;
                    } else {
                      pageNum = page - 2 + i;
                    }
                    
                    return (
                      <button
                        key={pageNum}
                        onClick={() => setPage(pageNum)}
                        className={`w-8 h-8 rounded-md text-sm transition-colors ${
                          page === pageNum
                            ? 'bg-indigo-500 text-white'
                            : 'border text-gray-600 hover:bg-gray-50'
                        }`}
                      >
                        {pageNum}
                      </button>
                    );
                  })}
                </div>
                <button
                  onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                  disabled={page >= totalPages}
                  className="px-4 py-2 border rounded-md text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
                >
                  下一页
                </button>
              </div>
            )}
            
            {/* 显示总数 */}
            <div className="text-center text-xs text-gray-400 mt-4">
              共 {total} 条公告
            </div>
          </>
        )}
      </div>
    </div>
  );
}