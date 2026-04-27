// components/question/TagSelector.tsx
"use client";

import { TagIcon, HashtagIcon } from "@heroicons/react/24/outline";
import { Tag } from "@/lib/api/types";

interface TagSelectorProps {
  tags: Tag[];
  selectedTags: number[];
  onToggleTag: (tagId: number) => void;
  loading?: boolean;
  maxTags?: number;
}

export function TagSelector({
  tags,
  selectedTags,
  onToggleTag,
  loading,
  maxTags = 5,
}: TagSelectorProps) {
  const isMaxTagsReached = selectedTags.length >= maxTags;

  if (loading) {
    return (
      <div className="space-y-2">
        <label className="text-sm font-medium text-base-content flex items-center gap-1">
          <TagIcon className="w-4 h-4" />
          标签
        </label>
        <div className="flex flex-wrap gap-2">
          {[1, 2, 3, 4, 5].map((i) => (
            <div
              key={i}
              className="h-8 w-16 bg-base-200 rounded-full animate-pulse"
            />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <label className="text-sm font-medium text-base-content flex items-center gap-1">
          <TagIcon className="w-4 h-4" />
          标签
          <span className="text-xs text-base-content/40 ml-1">
            (可选，最多 {maxTags} 个)
          </span>
        </label>
        {selectedTags.length > 0 && (
          <span className="text-xs text-base-content/40">
            已选择 {selectedTags.length}/{maxTags}
          </span>
        )}
      </div>

      <div className="flex flex-wrap gap-2">
        {tags
          .filter((tag) => tag.id > 0)
          .map((tag) => {
            const isSelected = selectedTags.includes(tag.id);
            const isDisabled = !isSelected && isMaxTagsReached;

            return (
              <button
                key={tag.id}
                type="button"
                onClick={() => !isDisabled && onToggleTag(tag.id)}
                disabled={isDisabled}
                className={`
                relative group px-3 py-1.5 rounded-full text-sm font-medium transition-all duration-200
                ${
                  isSelected
                    ? "bg-primary text-primary-content shadow-md hover:bg-primary-focus"
                    : "bg-base-200 text-base-content hover:bg-base-300"
                }
                ${
                  isDisabled && !isSelected
                    ? "opacity-50 cursor-not-allowed"
                    : "cursor-pointer"
                }
              `}
                style={
                  !isSelected && !isDisabled
                    ? {
                        borderLeft: `3px solid ${tag.color || "#ef4444"}`,
                      }
                    : undefined
                }
              >
                <div className="flex items-center gap-1.5">
                  <HashtagIcon className="w-3 h-3" />
                  <span>{tag.name}</span>
                  {isSelected && <span className="text-xs ml-0.5">✓</span>}
                </div>

                {/* 悬停提示 */}
                {isDisabled && !isSelected && (
                  <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1 px-2 py-0.5 bg-gray-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none z-10">
                    最多选择 {maxTags} 个标签
                  </div>
                )}
              </button>
            );
          })}
      </div>

      {tags.length === 0 && (
        <div className="text-sm text-base-content/40 text-center py-4 bg-base-200/50 rounded-lg">
          <TagIcon className="w-8 h-8 mx-auto mb-2 opacity-30" />
          暂无可用标签
        </div>
      )}

      <p className="text-xs text-base-content/50 flex items-center gap-1">
        <TagIcon className="w-3 h-3" />
        选择相关的标签，帮助其他人更快找到你的问题
      </p>
    </div>
  );
}
