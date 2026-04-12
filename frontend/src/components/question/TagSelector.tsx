// components/question/TagSelector.tsx
'use client';

import { TagIcon } from '@heroicons/react/24/outline';
import { Tag } from '@/lib/api/types';

interface TagSelectorProps {
  tags: Tag[];
  selectedTags: number[];
  onToggleTag: (tagId: number) => void;
  loading?: boolean;
}

export function TagSelector({ tags, selectedTags, onToggleTag, loading }: TagSelectorProps) {
  if (loading) {
    return (
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          <TagIcon className="w-4 h-4 inline mr-1" />
          标签
        </label>
        <div className="flex flex-wrap gap-2">
          {[1, 2, 3, 4].map(i => (
            <div key={i} className="w-20 h-8 bg-gray-200 rounded-full animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        <TagIcon className="w-4 h-4 inline mr-1" />
        标签
      </label>
      <div className="flex flex-wrap gap-2">
        {tags.filter(tag => tag.id > 0).map((tag) => (
          <button
            key={tag.id}
            type="button"
            onClick={() => onToggleTag(tag.id)}
            className={`px-3 py-1.5 rounded-full text-sm transition-colors ${
              selectedTags.includes(tag.id)
                ? 'bg-indigo-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
            style={{
              borderLeft: selectedTags.includes(tag.id) 
                ? undefined 
                : `3px solid ${tag.color || '#6366f1'}`
            }}
          >
            {tag.name}
          </button>
        ))}
      </div>
      <p className="mt-2 text-sm text-gray-400">选择相关的标签，帮助其他人更快找到你的问题</p>
    </div>
  );
}