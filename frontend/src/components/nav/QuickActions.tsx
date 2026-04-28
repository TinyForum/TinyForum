"use client";

import Link from "next/link";
import { PenSquare,  MessageCircleQuestion } from "lucide-react";
import { useTranslations } from "next-intl";

interface QuickActionsProps {
  isAuthenticated: boolean;
}

export default function QuickActions({ isAuthenticated }: QuickActionsProps) {
  const t = useTranslations("Nav");
  if (!isAuthenticated) return null;

  return (
    <div className="hidden md:flex items-center gap-1">
      {/* 快速发帖 */}
      <Link
        href="/posts/new"
        className="btn btn-primary btn-sm gap-1 shadow-md hover:shadow-lg transition-all"
      >
        <PenSquare className="w-4 h-4" />
        <span className="hidden sm:inline">{t("create_post")}</span>
      </Link>

      {/* 快速提问 */}
      <Link href="/questions/ask" className="btn btn-outline btn-sm gap-1">
        <MessageCircleQuestion className="w-4 h-4" />
        <span className="hidden sm:inline">{t("ask_question")}</span>
      </Link>
    </div>
  );
}
