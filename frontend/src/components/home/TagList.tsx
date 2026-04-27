"use client";

import { TagIcon } from "lucide-react";
import { useTranslations } from "next-intl";

interface TagListProps {
  tags: any[];
  selectedTag: number | null;
  onTagChange: (tagId: number | null) => void;
}

export default function TagList({
  tags,
  selectedTag,
  onTagChange,
}: TagListProps) {
  const t = useTranslations("post");

  return (
    <div className="card bg-base-100 border border-base-300 shadow-sm">
      <div className="card-body p-4">
        <h3 className="font-bold flex items-center gap-2 mb-3">
          <TagIcon className="w-4 h-4 text-primary" /> {t("hot_tags")}
        </h3>
        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => onTagChange(null)}
            className={`badge badge-lg cursor-pointer ${!selectedTag ? "badge-primary" : "badge-ghost hover:badge-primary"}`}
          >
            {t("all")}
          </button>
          {tags.slice(0, 12).map((tag) => (
            <TagButton
              key={tag.id}
              tag={tag}
              isSelected={selectedTag === tag.id}
              onSelect={() =>
                onTagChange(selectedTag === tag.id ? null : tag.id)
              }
            />
          ))}
        </div>
      </div>
    </div>
  );
}

function TagButton({
  tag,
  isSelected,
  onSelect,
}: {
  tag: any;
  isSelected: boolean;
  onSelect: () => void;
}) {
  return (
    <button
      onClick={onSelect}
      className="badge badge-lg cursor-pointer hover:opacity-80 transition-opacity"
      style={{
        backgroundColor: isSelected ? tag.color : tag.color + "20",
        color: tag.color,
        borderColor: tag.color + "40",
      }}
    >
      {tag.name}
      <span className="ml-1 opacity-60 text-xs">({tag.post_count})</span>
    </button>
  );
}
