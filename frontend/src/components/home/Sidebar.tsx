"use client";

import Link from "next/link";
import { TagIcon, Trophy, ChevronRight } from "lucide-react";
import { useTranslations } from "next-intl";
import TagList from "./TagList";
import LeaderboardList from "./LeaderboardList";
import SiteInfo from "./SiteInfo";

interface SidebarProps {
  tags: any[];
  selectedTag: number | null;
  onTagChange: (tagId: number | null) => void;
  leaderboard: any[];
}

export default function Sidebar({
  tags,
  selectedTag,
  onTagChange,
  leaderboard,
}: SidebarProps) {
  const t = useTranslations("post");

  return (
    <aside className="w-full lg:w-64 xl:w-72 flex-none space-y-4">
      <TagList
        tags={tags}
        selectedTag={selectedTag}
        onTagChange={onTagChange}
      />
      <LeaderboardList leaderboard={leaderboard} />
      <SiteInfo />
    </aside>
  );
}