"use client";

import Link from "next/link";
import { Flame, Clock, PenSquare, TrendingUp, Award, RefreshCw, Megaphone } from "lucide-react";
import { useTranslations } from "next-intl";
import { SortBy } from "@/types";

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
  const t = useTranslations("post");

  const sortOptions: { value: SortBy; label: string; icon: React.ReactNode }[] = [
    { value: "random", label: t("random"), icon: <Clock className="w-4 h-4" /> },
    { value: "hot", label: t("hot"), icon: <Flame className="w-4 h-4" /> },
    { value: "like", label: t("like"), icon: <Award className="w-4 h-4" /> },
    { value: "latest", label: t("latest"), icon: <Megaphone className="w-4 h-4" /> },
  ];

  return (
    <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-4 bg-base-100 rounded-xl p-3 border border-base-300">
      {/* 左侧：排序选项 */}
      <div className="flex flex-wrap items-center gap-2">
        {sortOptions.map((option) => (
          <button
            key={option.value}
            className={`btn btn-sm gap-1 transition-all duration-200 ${
              sortBy === option.value 
                ? "btn-primary shadow-md" 
                : "btn-ghost hover:bg-base-200"
            }`}
            onClick={() => onSortChange(option.value)}
          >
            {option.icon}
            <span className="hidden sm:inline">{option.label}</span>
          </button>
        ))}
        
        {/* 刷新按钮 */}
        <button
          className="btn btn-sm btn-ghost gap-1"
          onClick={() => onRefetch()}
          disabled={isLoading}
        >
          <RefreshCw className={`w-4 h-4 ${isLoading ? "animate-spin" : ""}`} />
          <span className="hidden sm:inline">{t("refresh")}</span>
        </button>
      </div>

      {/* 右侧：发帖按钮和统计 */}
      <div className="flex items-center gap-3 w-full sm:w-auto">
        {totalCount !== undefined && totalCount > 0 && (
          <div className="text-xs text-muted-foreground whitespace-nowrap">
            {t("total_posts", { count: totalCount })}
          </div>
        )}
        
        {isAuthenticated && (
          <Link href="/posts/new" className="btn btn-primary btn-sm gap-1 w-full sm:w-auto">
            <PenSquare className="w-4 h-4" /> 
            <span>{t("create")}</span>
          </Link>
        )}
      </div>
    </div>
  );
}