"use client";

import { useState, useEffect, useCallback } from "react";
import { announcementApi } from "@/shared/api";
import { toast } from "react-hot-toast";
import {
  MegaphoneIcon,
  PinIcon,
  FileTextIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from "lucide-react";
import { AnnouncementCard } from "@/features/announcements/components/AnnouncementCard";
import {
  AnnouncementDO,
  AnnouncementStatus,
} from "@/shared/api/types/announcement.model";

// 分页组件
function Pagination({
  currentPage,
  totalPages,
  onPageChange,
}: {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}) {
  const getPageNumbers = () => {
    const pages: number[] = [];
    const maxVisible = 5;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) pages.push(i);
    } else {
      if (currentPage <= 3) {
        for (let i = 1; i <= maxVisible; i++) pages.push(i);
      } else if (currentPage >= totalPages - 2) {
        for (let i = totalPages - maxVisible + 1; i <= totalPages; i++)
          pages.push(i);
      } else {
        for (let i = currentPage - 2; i <= currentPage + 2; i++) pages.push(i);
      }
    }
    return pages;
  };

  return (
    <div className="flex justify-center items-center gap-2 mt-10">
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="btn btn-ghost btn-sm gap-1"
      >
        <ChevronLeftIcon className="w-4 h-4" />
        上一页
      </button>

      <div className="flex gap-1.5 mx-2">
        {getPageNumbers().map((pageNum) => (
          <button
            key={pageNum}
            onClick={() => onPageChange(pageNum)}
            className={`btn btn-sm min-w-[2.5rem] ${
              currentPage === pageNum ? "btn-primary" : "btn-ghost"
            }`}
          >
            {pageNum}
          </button>
        ))}
      </div>

      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage >= totalPages}
        className="btn btn-ghost btn-sm gap-1"
      >
        下一页
        <ChevronRightIcon className="w-4 h-4" />
      </button>
    </div>
  );
}

// 加载骨架屏
function LoadingSkeleton() {
  return (
    <div className="space-y-4">
      {[1, 2, 3].map((i) => (
        <div key={i} className="bg-base-100 rounded-xl p-5 animate-pulse">
          <div className="flex items-center gap-2 mb-3">
            <div className="h-4 w-4 bg-base-200 rounded" />
            <div className="h-4 w-16 bg-base-200 rounded-full" />
          </div>
          <div className="h-6 bg-base-200 rounded w-3/4 mb-3" />
          <div className="h-4 bg-base-200 rounded w-full mb-2" />
          <div className="h-4 bg-base-200 rounded w-2/3" />
        </div>
      ))}
    </div>
  );
}

// 空状态组件
function EmptyState() {
  return (
    <div className="bg-base-100 rounded-xl p-16 text-center">
      <div className="text-7xl mb-5 opacity-50">📢</div>
      <h3 className="text-xl font-semibold text-base-content mb-2">暂无公告</h3>
      <p className="text-base-content/60">暂时没有公告，稍后再来看看吧</p>
    </div>
  );
}

export default function AnnouncementsPage() {
  const [announcements, setAnnouncements] = useState<AnnouncementDO[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 12;

  const loadAnnouncements = useCallback(async () => {
    setLoading(true);
    try {
      const response = await announcementApi.list({
        page,
        page_size: pageSize,
        status: AnnouncementStatus.Published,
      });

      if (response.data.code === 0) {
        if (response.data.data) {
          setAnnouncements(response.data.data.list || []);
          setTotal(response.data.data.total || 0);
        }
      } else {
        toast.error(response.data.message || "加载失败");
      }
    } catch {
      toast.error("加载失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    loadAnnouncements();
  }, [page, loadAnnouncements]);

  const totalPages = Math.ceil(total / pageSize);

  const pinnedAnnouncements = announcements.filter((a) => a.is_pinned);
  const normalAnnouncements = announcements.filter((a) => !a.is_pinned);

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      {/* 头部横幅 */}
      <div className="relative overflow-hidden bg-gradient-to-r from-primary/10 via-secondary/5 to-accent/10 border-b border-base-200">
        <div className="max-w-5xl mx-auto px-4 py-12 md:py-16">
          <div className="text-center md:text-left md:flex md:items-center md:justify-between">
            <div className="space-y-3">
              <div className="flex items-center justify-center md:justify-start gap-3">
                <div className="p-3 bg-primary/10 rounded-2xl">
                  <MegaphoneIcon className="w-8 h-8 text-primary" />
                </div>
                <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                  公告中心
                </h1>
              </div>
              <p className="text-base-content/70 text-lg">
                了解最新的平台动态和重要通知
              </p>
            </div>

            {/* 统计卡片 */}
            {!loading && announcements.length > 0 && (
              <div className="mt-6 md:mt-0 flex gap-4">
                <div className="text-center px-6 py-3 bg-base-100/50 backdrop-blur-sm rounded-2xl">
                  <div className="text-2xl font-bold text-primary">{total}</div>
                  <div className="text-xs text-base-content/50">总公告数</div>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* 装饰性背景圆点 */}
        <div className="absolute top-0 right-0 -translate-y-1/2 translate-x-1/2 w-64 h-64 bg-primary/5 rounded-full blur-3xl" />
        <div className="absolute bottom-0 left-0 translate-y-1/2 -translate-x-1/2 w-48 h-48 bg-secondary/5 rounded-full blur-3xl" />
      </div>

      {/* 内容区域 */}
      <div className="max-w-4xl mx-auto px-4 py-10">
        {loading ? (
          <LoadingSkeleton />
        ) : announcements.length === 0 ? (
          <EmptyState />
        ) : (
          <>
            <div className="space-y-4">
              {/* 置顶公告区域 */}
              {pinnedAnnouncements.length > 0 && (
                <div className="mb-8">
                  <div className="flex items-center gap-2.5 mb-4">
                    <div className="w-1 h-6 bg-primary rounded-full" />
                    <PinIcon className="w-4 h-4 text-warning" />
                    <h2 className="text-sm font-semibold text-base-content/70 uppercase tracking-wide">
                      置顶公告
                    </h2>
                    <div className="flex-1 h-px bg-gradient-to-r from-base-200 to-transparent" />
                  </div>
                  <div className="space-y-3">
                    {pinnedAnnouncements.map((announcement) => (
                      <AnnouncementCard
                        key={announcement.id}
                        announcement={announcement}
                      />
                    ))}
                  </div>
                </div>
              )}

              {/* 普通公告区域 */}
              {normalAnnouncements.length > 0 && (
                <div>
                  {pinnedAnnouncements.length > 0 && (
                    <div className="flex items-center gap-2.5 mb-4">
                      <div className="w-1 h-6 bg-secondary rounded-full" />
                      <FileTextIcon className="w-4 h-4 text-secondary" />
                      <h2 className="text-sm font-semibold text-base-content/70 uppercase tracking-wide">
                        最新公告
                      </h2>
                      <div className="flex-1 h-px bg-gradient-to-r from-base-200 to-transparent" />
                    </div>
                  )}
                  <div className="space-y-3">
                    {normalAnnouncements.map((announcement) => (
                      <AnnouncementCard
                        key={announcement.id}
                        announcement={announcement}
                      />
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* 分页 */}
            {totalPages > 1 && (
              <Pagination
                currentPage={page}
                totalPages={totalPages}
                onPageChange={setPage}
              />
            )}

            {/* 显示总数 */}
            <div className="text-center text-xs text-base-content/40 mt-6">
              共 {total} 条公告
            </div>
          </>
        )}
      </div>
    </div>
  );
}
