// src/components/home/LeftSidebar.tsx
"use client";

import { useState } from "react";
import Link from "next/link";
import { 
  LayoutDashboard, 
  Hash, 
  MessageSquare, 
  FileText, 
  TrendingUp,
  HelpCircle,
  Megaphone,
  ChevronDown,
  ChevronRight,
  Tag,
  Sparkles,
  PenIcon
} from "lucide-react";
import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";
import { PostType } from "@/lib/api";
// import { PostType } from "@/types";

export type FilterType = "all" | PostType;

interface Board {
  id: number;
  name: string;
  slug: string;
  icon: string;
  description: string;
  post_count: number;
  children?: Board[];
}

interface LeftSidebarProps {
  boards: Board[];
  tags: any[];
  selectedBoard: number | null;
  selectedTag: number | null;
    filterType: FilterType;  
  onBoardChange: (boardId: number | null) => void;
  onTagChange: (tagId: number | null) => void;
  onPostTypeChange: (type: FilterType) => void;
}

export default function LeftSidebar({
  boards,
  tags,
  selectedBoard,
  selectedTag,
  filterType,
  onBoardChange,
  onTagChange,
  onPostTypeChange,
}: LeftSidebarProps) {
  const t = useTranslations("sidebar");
  const [expandedBoards, setExpandedBoards] = useState<Set<number>>(new Set());

  const toggleBoard = (boardId: number) => {
    const newExpanded = new Set(expandedBoards);
    if (newExpanded.has(boardId)) {
      newExpanded.delete(boardId);
    } else {
      newExpanded.add(boardId);
    }
    setExpandedBoards(newExpanded);
  };

  const renderBoardTree = (boardList: Board[], level = 0) => {
    return boardList.map((board) => (
      <div key={board.id} style={{ marginLeft: level * 12 }}>
        <button
          onClick={() => {
            if (board.children && board.children.length > 0) {
              toggleBoard(board.id);
            }
            onBoardChange(board.id);
          }}
          className={cn(
            "w-full flex items-center justify-between px-3 py-2 rounded-lg text-sm transition-colors",
            selectedBoard === board.id
              ? "bg-primary/10 text-primary font-medium"
              : "hover:bg-muted text-muted-foreground hover:text-foreground"
          )}
        >
          <div className="flex items-center gap-2 flex-1 min-w-0">
            {board.children && board.children.length > 0 && (
              <span className="flex-shrink-0">
                {expandedBoards.has(board.id) ? (
                  <ChevronDown className="w-3 h-3" />
                ) : (
                  <ChevronRight className="w-3 h-3" />
                )}
              </span>
            )}
            <span className="flex-shrink-0 w-4 h-4">
              {board.icon || <Hash className="w-4 h-4" />}
            </span>
            <span className="truncate">{board.name}</span>
          </div>
          {board.post_count > 0 && (
            <span className="text-xs text-muted-foreground">{board.post_count}</span>
          )}
        </button>
        {board.children && expandedBoards.has(board.id) && (
          <div className="mt-1">
            {renderBoardTree(board.children, level + 1)}
          </div>
        )}
      </div>
    ));
  };

  return (
    <aside className="space-y-4">
      {/* 导航菜单 */}
      <div className="rounded-lg border bg-card">
        <div className="p-3 border-b">
          <h3 className="font-semibold flex items-center gap-2">
            <LayoutDashboard className="w-4 h-4" />
            {t("navigation")}
          </h3>
        </div>
        <div className="p-2 space-y-1">
          {/* 所有 */}
          <button
            onClick={() => {
              onBoardChange(null);
              onTagChange(null);
              onPostTypeChange("all");
            }}
            className={cn(
              "w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors",
              !selectedBoard && !selectedTag && filterType === "all"
                ? "bg-primary/10 text-primary font-medium"
                : "hover:bg-muted text-muted-foreground hover:text-foreground"
            )}
          >
            <TrendingUp className="w-4 h-4" />
            {t("all")}
          </button>
          
          {/* 问答 */}
          <button
            onClick={() => {
              onBoardChange(null);
              onTagChange(null);
              onPostTypeChange("question");
            }}
            className={cn(
              "w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors",
              filterType === "question"
                ? "bg-primary/10 text-primary font-medium"
                : "hover:bg-muted text-muted-foreground hover:text-foreground"
            )}
          >
            <HelpCircle className="w-4 h-4" />
            {t("questions")}
            <span className="ml-auto text-xs text-muted-foreground">问答</span>
          </button>
          
          {/* 文章 */}
          <button
            onClick={() => {
              onBoardChange(null);
              onTagChange(null);
              onPostTypeChange("article");
            }}
            className={cn(
              "w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors",
              filterType === "article"
                ? "bg-primary/10 text-primary font-medium"
                : "hover:bg-muted text-muted-foreground hover:text-foreground"
            )}
          >
            <FileText className="w-4 h-4" />
            {t("articles")}
            <span className="ml-auto text-xs text-muted-foreground">文章</span>
          </button>
          
          {/* 帖子 */}
          <button
            onClick={() => {
              onBoardChange(null);
              onTagChange(null);
              onPostTypeChange("post");
            }}
            className={cn(
              "w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors",
              filterType === "post"
                ? "bg-primary/10 text-primary font-medium"
                : "hover:bg-muted text-muted-foreground hover:text-foreground"
            )}
          >
            <PenIcon className="w-4 h-4" />
            {t("posts")}
            <span className="ml-auto text-xs text-muted-foreground">帖子</span>
          </button>
        </div>
      </div>

      {/* 公告栏 */}
      <div className="rounded-lg border bg-card">
        <div className="p-3 border-b">
          <h3 className="font-semibold flex items-center gap-2">
            <Megaphone className="w-4 h-4" />
            {t("announcements")}
          </h3>
        </div>
        <div className="p-3 space-y-2">
          <Link 
            href="/announcements/1" 
            className="block text-sm hover:text-primary transition-colors"
          >
            • 社区新功能上线通知
          </Link>
          <Link 
            href="/announcements/2" 
            className="block text-sm hover:text-primary transition-colors"
          >
            • 积分规则调整公告
          </Link>
          <Link 
            href="/announcements/3" 
            className="block text-sm hover:text-primary transition-colors"
          >
            • 版主招募中
          </Link>
          <Link 
            href="/announcements" 
            className="block text-xs text-muted-foreground hover:text-primary mt-2"
          >
            查看全部 →
          </Link>
        </div>
      </div>

      {/* 板块列表 */}
      {boards && boards.length > 0 && (
        <div className="rounded-lg border bg-card">
          <div className="p-3 border-b">
            <h3 className="font-semibold flex items-center gap-2">
              <Hash className="w-4 h-4" />
              {t("boards")}
            </h3>
          </div>
          <div className="p-2 space-y-1 max-h-[400px] overflow-y-auto">
            {renderBoardTree(boards)}
          </div>
        </div>
      )}

      {/* 热门标签 */}
      {tags && tags.length > 0 && (
        <div className="rounded-lg border bg-card">
          <div className="p-3 border-b">
            <h3 className="font-semibold flex items-center gap-2">
              <Tag className="w-4 h-4" />
              {t("hot_tags")}
            </h3>
          </div>
          <div className="p-3 flex flex-wrap gap-2">
            {tags.slice(0, 15).map((tag) => (
              <button
                key={tag.id}
                onClick={() => onTagChange(tag.id)}
                className={cn(
                  "px-2 py-1 rounded-md text-xs transition-colors",
                  selectedTag === tag.id
                    ? "bg-primary text-primary-foreground"
                    : "bg-muted hover:bg-muted/80 text-muted-foreground"
                )}
              >
                {tag.name}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* 社区统计 */}
      <div className="rounded-lg border bg-card">
        <div className="p-3 border-b">
          <h3 className="font-semibold flex items-center gap-2">
            <Sparkles className="w-4 h-4" />
            {t("community_stats")}
          </h3>
        </div>
        <div className="p-3 space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-muted-foreground">今日发帖</span>
            <span className="font-medium">128</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">昨日活跃</span>
            <span className="font-medium">1,234</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">总用户数</span>
            <span className="font-medium">5,678</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">总帖子数</span>
            <span className="font-medium">12,345</span>
          </div>
        </div>
      </div>
    </aside>
  );
}