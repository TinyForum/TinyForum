"use client";

import Link from "next/link";
import { PenSquare, Sparkles, MessageCircleQuestion } from "lucide-react";

interface QuickActionsProps {
  isAuthenticated: boolean;
}

export default function QuickActions({ isAuthenticated }: QuickActionsProps) {
  if (!isAuthenticated) return null;

  return (
    <div className="hidden md:flex items-center gap-1">
      {/* 快速发帖 */}
      <Link
        href="/posts/new"
        className="btn btn-primary btn-sm gap-1 shadow-md hover:shadow-lg transition-all"
      >
        <PenSquare className="w-4 h-4" />
        <span className="hidden sm:inline">写帖子</span>
      </Link>

      {/* 快速提问 */}
      <Link
        href="/questions/ask"
        className="btn btn-outline btn-sm gap-1"
      >
        <MessageCircleQuestion className="w-4 h-4" />
        <span className="hidden sm:inline">提问</span>
      </Link>
    </div>
  );
}