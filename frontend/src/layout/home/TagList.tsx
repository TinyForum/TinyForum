"use client";

import { TagIcon } from "lucide-react";
import { useTranslations } from "next-intl";

interface Tag {
  id: number;
  name: string;
  color: string;
  post_count: number;
  description?: string;
}

interface TagListProps {
  tags: Tag[];
  selectedTag: number | null;
  onTagChange: (tagId: number | null) => void;
}

interface TagButtonProps {
  tag: Tag;
  isSelected: boolean;
  onSelect: () => void;
}

export default function TagList({
  tags,
  selectedTag,
  onTagChange,
}: TagListProps) {
  const t = useTranslations("Post");

  return (
    <div className="card bg-base-100 border border-base-300 shadow-sm">
      <div className="card-body p-4">
        <h3 className="font-bold flex items-center gap-2 mb-3">
          <TagIcon className="w-4 h-4 text-primary" /> {t("hot_tags")}
        </h3>
        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => onTagChange(null)}
            className={`badge badge-lg cursor-pointer transition-all duration-200 ${
              !selectedTag
                ? "badge-primary shadow-sm"
                : "badge-ghost hover:badge-primary"
            }`}
          >
            {t("all")}
          </button>
          {tags.slice(0, 12).map((tag: Tag) => (
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

function TagButton({ tag, isSelected, onSelect }: TagButtonProps) {
  // 计算背景色
  const getBackgroundColor = (): string => {
    if (isSelected) {
      return tag.color;
    }
    // 使用正则验证颜色格式，确保颜色值有效
    if (tag.color && /^#[0-9A-Fa-f]{6}$/.test(tag.color)) {
      return tag.color + "20";
    }
    return "transparent";
  };

  // 计算文字颜色
  const getTextColor = (): string => {
    if (isSelected && tag.color) {
      // 根据背景色决定文字颜色（白色或黑色）
      const hexColor = tag.color.replace("#", "");
      const r = parseInt(hexColor.substring(0, 2), 16);
      const g = parseInt(hexColor.substring(2, 4), 16);
      const b = parseInt(hexColor.substring(4, 6), 16);
      const brightness = (r * 299 + g * 587 + b * 114) / 1000;
      return brightness > 128 ? "#000000" : "#FFFFFF";
    }
    return tag.color || "#666666";
  };

  // 计算边框颜色
  const getBorderColor = (): string => {
    if (tag.color && /^#[0-9A-Fa-f]{6}$/.test(tag.color)) {
      return tag.color + "40";
    }
    return "#e5e7eb";
  };

  return (
    <button
      onClick={onSelect}
      className="badge badge-lg cursor-pointer hover:scale-105 transition-all duration-200 font-medium"
      style={{
        backgroundColor: getBackgroundColor(),
        color: getTextColor(),
        borderColor: getBorderColor(),
        boxShadow: isSelected ? "0 2px 4px rgba(0,0,0,0.1)" : "none",
      }}
    >
      {tag.name}
      <span
        className="ml-1 text-xs opacity-75"
        style={{ color: isSelected ? "inherit" : tag.color }}
      >
        ({tag.post_count})
      </span>
    </button>
  );
}
