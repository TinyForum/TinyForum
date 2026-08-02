"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { FileText, MessageSquare, BookOpen, Hash } from "lucide-react";
import { ViolationStatus } from "./ViolationStatus";
import WorkCard from "@/shared/ui/creation/WorkCard";
import { PostType } from "@/shared/api/types/post.model";
import { useProfileContent } from "../hooks/useProfileContent";

interface ProfileContentProps {
  userId: number;
  isAuthenticated: boolean;
}

const TAB_CONFIG = [
  { key: "image_text" as PostType, label: "the_posts", icon: FileText },
  { key: "article" as PostType, label: "the_articles", icon: BookOpen },
  { key: "question" as PostType, label: "the_questions", icon: MessageSquare },
  { key: "topic" as PostType, label: "the_topics", icon: Hash },
];

export function ProfileContent({
  userId,
  isAuthenticated,
}: ProfileContentProps) {
  const [tab, setTab] = useState<PostType>("image_text");
  const t = useTranslations("Profile");

  const { data: postsData, isLoading } = useProfileContent(userId, tab);

  const posts = postsData?.list ?? [];
  const currentTabConfig = TAB_CONFIG.find((tb) => tb.key === tab);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="skeleton h-12 w-full rounded-xl" />
        <div className="skeleton h-40 w-full rounded-xl" />
        <div className="skeleton h-40 w-full rounded-xl" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {isAuthenticated && <ViolationStatus />}

      <div className="tabs tabs-boxed bg-base-100 border border-base-300 p-1">
        {TAB_CONFIG.map(({ key, label, icon: Icon }) => (
          <button
            key={key}
            className={`tab gap-2 flex-1 sm:flex-initial ${tab === key ? "tab-active" : ""}`}
            onClick={() => setTab(key)}
          >
            <Icon className="w-4 h-4" />
            <span className="hidden sm:inline">{t(label)}</span>
          </button>
        ))}
      </div>

      {posts.length === 0 ? (
        <div className="text-center py-16 bg-base-100 rounded-xl border border-base-200">
          {currentTabConfig && (
            <currentTabConfig.icon className="w-12 h-12 mx-auto mb-3 opacity-30" />
          )}
          <p className="text-base-content/40">
            {tab === "image_text"
              ? t("no_posts")
              : tab === "article"
                ? t("no_articles")
                : tab === "question"
                  ? t("no_questions")
                  : t("no_topics")}
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {posts.map((post) => (
            <WorkCard key={post.id} post={post} />
          ))}
        </div>
      )}
    </div>
  );
}
