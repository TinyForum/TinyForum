"use client";

import Link from "next/link";
import {
  Flame,
  Clock,
  PenSquare,
  Award,
  RefreshCw,
  Megaphone,
} from "lucide-react";
import { useTranslations } from "next-intl";
import { SortBy } from "@/shared/type/posts.types";
import { useMediaQuery } from "@/features/common/hooks/useMediaQuery";
import { DesktopSortView } from "./DesktopSortView";
import { MobileSortView } from "./MobileSortView";

interface PostFilterBarProps {
  sortBy: SortBy;
  onSortChange: (sortBy: SortBy) => void;
  isAuthenticated: boolean;
  onRefetch: () => void;
  isLoading?: boolean;
  totalCount?: number;
}

export default function PostFilterBar({
  sortBy,
  onSortChange,
  isAuthenticated,
  onRefetch,
  isLoading = false,
  totalCount,
}: PostFilterBarProps) {
  const t = useTranslations("Post");
  const isDesktop = useMediaQuery("(min-width: 1024px)");

  const sortOptions: { value: SortBy; label: string; icon: React.ReactNode }[] =
    [
      {
        value: "random",
        label: t("random"),
        icon: <Clock className="w-4 h-4" />,
      },
      { value: "hot", label: t("hot"), icon: <Flame className="w-4 h-4" /> },
      { value: "like", label: t("like"), icon: <Award className="w-4 h-4" /> },
      {
        value: "latest",
        label: t("latest"),
        icon: <Megaphone className="w-4 h-4" />,
      },
    ];

  return (
    <div className="relative mb-4 flex items-center gap-3 rounded-box border border-base-300 bg-base-100 p-3 shadow-sm">
      {/* 排序和刷新区域 */}
      <div className="flex flex-1 items-center gap-2 min-w-0">
        {/* 排序组件 */}
        {isDesktop ? (
          <DesktopSortView
            sortOptions={sortOptions}
            sortBy={sortBy}
            onSortChange={onSortChange}
          />
        ) : (
          <MobileSortView
            sortOptions={sortOptions}
            sortBy={sortBy}
            onSortChange={onSortChange}
          />
        )}

        {/* 刷新按钮 */}
        <button
          className="btn btn-sm btn-ghost gap-1 shrink-0"
          onClick={() => onRefetch()}
          disabled={isLoading}
          aria-label={t("refresh")}
        >
          <RefreshCw className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`} />
          <span className="hidden sm:inline">{t("refresh")}</span>
        </button>

        {/* 统计信息 */}
        {totalCount !== undefined && totalCount > 0 && (
          <div className="text-xs text-base-content/60 whitespace-nowrap hidden sm:block">
            {isDesktop ? t("total_posts", { count: totalCount }) : totalCount}
          </div>
        )}
      </div>

      {/* 发帖按钮 - 固定宽度 1/5 */}
      {isAuthenticated && (
        <div className="w-1/5 shrink-0 min-w-[80px]">
          <Link
            href="/posts/new"
            className="btn btn-primary btn-sm gap-1 w-full"
          >
            <PenSquare className="h-4 w-4" />
            <span className="hidden sm:inline">{t("create")}</span>
          </Link>
        </div>
      )}
    </div>
  );
}
