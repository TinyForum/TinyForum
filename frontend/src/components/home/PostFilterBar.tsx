"use client";

import Link from "next/link";
import { Flame, Clock, PenSquare } from "lucide-react";
import { useTranslations } from "next-intl";

interface PostFilterBarProps {
  sortBy: "" | "hot";
  onSortChange: (sortBy: "" | "hot") => void;
  isAuthenticated: boolean;
}

export default function PostFilterBar({
  sortBy,
  onSortChange,
  isAuthenticated,
}: PostFilterBarProps) {
  const t = useTranslations("post");

  return (
    <div className="flex items-center justify-between mb-4 bg-base-100 rounded-xl p-3 border border-base-300">
      <div className="flex items-center gap-2">
        <button
          className={`btn btn-sm gap-1 ${sortBy === "" ? "btn-primary" : "btn-ghost"}`}
          onClick={() => onSortChange("")}
        >
          <Clock className="w-4 h-4" /> {t("latest_posts")}
        </button>
        <button
          className={`btn btn-sm gap-1 ${sortBy === "hot" ? "btn-primary" : "btn-ghost"}`}
          onClick={() => onSortChange("hot")}
        >
          <Flame className="w-4 h-4" /> {t("hot_posts")}
        </button>
      </div>

      {isAuthenticated && (
        <Link href="/posts/new" className="btn btn-primary btn-sm gap-1">
          <PenSquare className="w-4 h-4" /> {t("create")}
        </Link>
      )}
    </div>
  );
}