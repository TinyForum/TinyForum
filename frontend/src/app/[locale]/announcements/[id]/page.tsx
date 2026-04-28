"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { announcementApi } from "@/shared/api";
import {
  MegaphoneIcon,
  CalendarIcon,
  EyeIcon,
  PinIcon,
  ArrowLeftIcon,
  TagIcon,
  ClockIcon,
  UserIcon,
  AlertTriangleIcon,
} from "lucide-react";
import type { Announcement } from "@/shared/api/modules/announcements";

// 错误响应类型
interface ErrorResponse {
  response?: {
    status?: number;
    data?: {
      message?: string;
    };
  };
  message?: string;
}

// 公告类型配置 - 使用主题色
const TYPE_CONFIG: Record<
  string,
  { color: string; label: string; bgColor: string; icon: typeof MegaphoneIcon }
> = {
  normal: {
    color: "text-primary",
    label: "普通公告",
    bgColor: "bg-primary/10",
    icon: MegaphoneIcon,
  },
  important: {
    color: "text-secondary",
    label: "重要公告",
    bgColor: "bg-secondary/10",
    icon: AlertTriangleIcon,
  },
  emergency: {
    color: "text-error",
    label: "紧急公告",
    bgColor: "bg-error/10",
    icon: AlertTriangleIcon,
  },
  event: {
    color: "text-accent",
    label: "活动公告",
    bgColor: "bg-accent/10",
    icon: CalendarIcon,
  },
};

// 加载骨架屏
function LoadingSkeleton() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="animate-pulse">
          {/* 返回按钮骨架 */}
          <div className="h-5 w-28 bg-base-200 rounded mb-6" />

          {/* 内容卡片骨架 */}
          <div className="bg-base-100 rounded-2xl shadow-sm overflow-hidden">
            <div className="p-8 border-b border-base-200">
              <div className="flex items-center gap-2 mb-5">
                <div className="h-6 w-16 bg-base-200 rounded-full" />
                <div className="h-6 w-12 bg-base-200 rounded-full" />
              </div>
              <div className="h-9 bg-base-200 rounded-lg w-3/4 mb-5" />
              <div className="flex items-center gap-5">
                <div className="h-5 w-32 bg-base-200 rounded" />
                <div className="h-5 w-24 bg-base-200 rounded" />
              </div>
            </div>
            <div className="p-8 space-y-3">
              <div className="h-4 bg-base-200 rounded w-full" />
              <div className="h-4 bg-base-200 rounded w-full" />
              <div className="h-4 bg-base-200 rounded w-2/3" />
              <div className="h-4 bg-base-200 rounded w-full" />
              <div className="h-4 bg-base-200 rounded w-5/6" />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// 错误状态组件
function ErrorState({ message }: { message: string }) {
  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100 flex items-center justify-center">
      <div className="text-center">
        <div className="text-6xl mb-5">😕</div>
        <h3 className="text-xl font-semibold text-base-content mb-2">
          加载失败
        </h3>
        <p className="text-base-content/60 mb-6">{message}</p>
        <Link href="/announcements" className="btn btn-primary btn-sm gap-2">
          <ArrowLeftIcon className="w-4 h-4" />
          返回公告列表
        </Link>
      </div>
    </div>
  );
}

export default function AnnouncementDetailPage() {
  const params = useParams();
  const [announcement, setAnnouncement] = useState<Announcement | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadAnnouncement = useCallback(async () => {
    const id = parseInt(params.id as string);
    if (isNaN(id)) {
      setError("无效的公告ID");
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const response = await announcementApi.getById(id);
      if (response.data.code === 0) {
        if (response.data.data) {
          setAnnouncement(response.data.data);
        }
      } else {
        setError(response.data.message || "公告不存在");
      }
    } catch (err: unknown) {
      console.error("Failed to load announcement:", err);
      const error = err as ErrorResponse;
      if (error.response?.status === 404) {
        setError("公告不存在");
      } else {
        setError("加载失败，请稍后重试");
      }
    } finally {
      setLoading(false);
    }
  }, [params.id]);

  useEffect(() => {
    loadAnnouncement();
  }, [params.id, loadAnnouncement]);

  const formatDate = (dateStr: string | null, withTime: boolean = true) => {
    if (!dateStr) return "待发布";
    const date = new Date(dateStr);
    if (withTime) {
      return date.toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    }
    return date.toLocaleDateString("zh-CN", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const isExpired =
    announcement?.expired_at && new Date(announcement.expired_at) < new Date();
  const typeConfig = announcement
    ? TYPE_CONFIG[announcement.type] || TYPE_CONFIG.normal
    : null;
  const TypeIcon = typeConfig?.icon || MegaphoneIcon;

  if (loading) {
    return <LoadingSkeleton />;
  }

  if (error || !announcement) {
    return <ErrorState message={error || "公告不存在"} />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-4xl mx-auto px-4 py-8 md:py-12">
        {/* 返回按钮 - 优化样式 */}
        <Link
          href="/announcements"
          className="inline-flex items-center gap-2 text-base-content/60 hover:text-primary transition-all duration-200 mb-6 group"
        >
          <ArrowLeftIcon className="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
          <span className="text-sm font-medium">返回公告列表</span>
        </Link>

        {/* 公告卡片 - 使用主题色和现代化设计 */}
        <article className="bg-base-100 rounded-2xl shadow-lg overflow-hidden border border-base-200">
          {/* 头部区域 - 增强视觉层次 */}
          <div className="relative p-8 pb-6 border-b border-base-200">
            {/* 装饰性背景 */}
            <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-primary/5 to-transparent rounded-full blur-2xl" />

            {/* 标签区域 */}
            <div className="flex items-center gap-2 mb-5 flex-wrap relative z-10">
              {typeConfig && (
                <span
                  className={`inline-flex items-center gap-1.5 text-xs font-medium px-3 py-1 rounded-full ${typeConfig.bgColor} ${typeConfig.color}`}
                >
                  <TypeIcon className="w-3.5 h-3.5" />
                  {typeConfig.label}
                </span>
              )}
              {announcement.is_pinned && (
                <span className="inline-flex items-center gap-1.5 text-xs font-medium px-3 py-1 rounded-full bg-warning/10 text-warning">
                  <PinIcon className="w-3.5 h-3.5" />
                  置顶公告
                </span>
              )}
              {isExpired && (
                <span className="inline-flex items-center gap-1.5 text-xs font-medium px-3 py-1 rounded-full bg-error/10 text-error">
                  <ClockIcon className="w-3.5 h-3.5" />
                  已过期
                </span>
              )}
              {announcement.board && (
                <span className="inline-flex items-center gap-1.5 text-xs font-medium px-3 py-1 rounded-full bg-base-200 text-base-content/60">
                  <TagIcon className="w-3.5 h-3.5" />
                  {announcement.board.name}
                </span>
              )}
            </div>

            {/* 标题 */}
            <h1 className="text-2xl md:text-3xl font-bold text-base-content mb-5 leading-tight">
              {announcement.title}
            </h1>

            {/* 元信息 - 优化布局 */}
            <div className="flex flex-wrap items-center gap-x-5 gap-y-2 text-sm text-base-content/50">
              <div className="flex items-center gap-1.5">
                <CalendarIcon className="w-4 h-4" />
                <span>
                  发布于{" "}
                  {formatDate(
                    announcement.published_at || announcement.created_at,
                  )}
                </span>
              </div>
              <div className="flex items-center gap-1.5">
                <EyeIcon className="w-4 h-4" />
                <span>{announcement.view_count || 0} 次阅读</span>
              </div>
              {announcement.created_by && (
                <div className="flex items-center gap-1.5">
                  <UserIcon className="w-4 h-4" />
                  <span>发布者 ID: {announcement.created_by}</span>
                </div>
              )}
            </div>
          </div>

          {/* 内容区域 - 优化阅读体验 */}
          <div className="p-8 md:p-10">
            <div
              className="prose prose-base max-w-none
                prose-headings:text-base-content prose-headings:font-semibold
                prose-h1:text-2xl prose-h2:text-xl prose-h3:text-lg
                prose-p:text-base-content/80 prose-p:leading-relaxed
                prose-a:text-primary prose-a:no-underline hover:prose-a:underline
                prose-strong:text-base-content prose-strong:font-semibold
                prose-code:text-primary prose-code:bg-primary/5 prose-code:px-1 prose-code:rounded
                prose-pre:bg-base-200 prose-pre:text-base-content
                prose-ul:text-base-content/80 prose-ol:text-base-content/80
                prose-li:marker:text-primary
                prose-blockquote:border-l-primary prose-blockquote:bg-primary/5 prose-blockquote:py-2
                prose-img:rounded-lg prose-img:shadow-md
                dark:prose-invert"
              dangerouslySetInnerHTML={{ __html: announcement.content }}
            />
          </div>

          {/* 页脚 - 添加过期提示和附加信息 */}
          {(isExpired || announcement.expired_at) && (
            <div className="p-6 bg-base-200/50 border-t border-base-200">
              {isExpired ? (
                <div className="flex items-center justify-center gap-2 text-sm text-error">
                  <AlertTriangleIcon className="w-4 h-4" />
                  <span>
                    此公告已于 {formatDate(announcement.expired_at, false)} 过期
                  </span>
                </div>
              ) : (
                <div className="flex items-center justify-center gap-2 text-sm text-base-content/40">
                  <ClockIcon className="w-4 h-4" />
                  <span>
                    有效期至 {formatDate(announcement.expired_at, false)}
                  </span>
                </div>
              )}
            </div>
          )}

          {/* 底部操作栏 */}
          <div className="p-6 bg-base-200/30 border-t border-base-200 flex justify-center gap-3">
            <Link href="/announcements" className="btn btn-ghost btn-sm gap-2">
              <MegaphoneIcon className="w-4 h-4" />
              查看全部公告
            </Link>
          </div>
        </article>
      </div>
    </div>
  );
}
