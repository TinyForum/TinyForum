"use client";

import { Search, X } from "lucide-react";
import { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import { useDebounce } from "@/hooks/useDebounce";
import { postApi } from "@/lib/api";

interface SearchBarProps {
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  onSearch: (e: React.FormEvent) => void;
}

interface SearchSuggestion {
  id: number;
  title: string;
  type: "post" | "user" | "topic";
  avatar?: string;
}

export default function SearchBar({ searchQuery, setSearchQuery, onSearch }: SearchBarProps) {
  const router = useRouter();
  const [isExpanded, setIsExpanded] = useState(false);
  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const searchInputRef = useRef<HTMLInputElement>(null);
  const debouncedQuery = useDebounce(searchQuery, 300);

  // 搜索建议
  useEffect(() => {
    if (debouncedQuery.length >= 2) {
      // 获取搜索建议
      const fetchSuggestions = async () => {
        try {
          const response = await postApi.list({ keyword: debouncedQuery, page_size: 5 });
          const posts = response.data.data?.list || [];
          setSuggestions(posts.map((post: any) => ({
            id: post.id,
            title: post.title,
            type: "post",
          })));
          setShowSuggestions(true);
        } catch (error) {
          console.error("Failed to fetch suggestions:", error);
        }
      };
      fetchSuggestions();
    } else {
      setSuggestions([]);
      setShowSuggestions(false);
    }
  }, [debouncedQuery]);

  const handleSuggestionClick = (suggestion: SearchSuggestion) => {
    if (suggestion.type === "post") {
      router.push(`/posts/${suggestion.id}`);
    }
    setSearchQuery("");
    setShowSuggestions(false);
    setIsExpanded(false);
  };

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
      setSearchQuery("");
      setShowSuggestions(false);
      setIsExpanded(false);
    }
  };

  // 移动端搜索展开
  if (isExpanded) {
    return (
      <form onSubmit={handleSearchSubmit} className="fixed inset-x-0 top-0 z-50 p-4 bg-base-100 shadow-lg animate-slideDown">
        <div className="flex gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
            <input
              ref={searchInputRef}
              type="text"
              placeholder="搜索帖子、用户..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input input-bordered w-full pl-9"
              autoFocus
            />
            {suggestions.length > 0 && (
              <div className="absolute top-full left-0 right-0 mt-2 bg-base-100 rounded-lg shadow-lg border border-base-200 z-50">
                {suggestions.map((suggestion) => (
                  <button
                    key={suggestion.id}
                    onClick={() => handleSuggestionClick(suggestion)}
                    className="w-full px-4 py-2 text-left hover:bg-base-200 transition-colors flex items-center gap-2"
                  >
                    <Search className="w-4 h-4 text-base-content/40" />
                    <span className="text-sm">{suggestion.title}</span>
                  </button>
                ))}
              </div>
            )}
          </div>
          <button
            type="button"
            onClick={() => setIsExpanded(false)}
            className="btn btn-ghost btn-sm"
          >
            取消
          </button>
        </div>
      </form>
    );
  }

  return (
    <>
      {/* 移动端搜索图标 */}
      <button
        onClick={() => setIsExpanded(true)}
        className="btn btn-ghost btn-sm btn-circle lg:hidden"
        aria-label="搜索"
      >
        <Search className="w-5 h-5" />
      </button>

      {/* 桌面端完整搜索框 */}
      <form onSubmit={handleSearchSubmit} className="hidden lg:block relative">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
          <input
            ref={searchInputRef}
            type="text"
            placeholder="搜索帖子、用户、专题..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onFocus={() => searchQuery.length >= 2 && setShowSuggestions(true)}
            onBlur={() => setTimeout(() => setShowSuggestions(false), 200)}
            className="input input-bordered input-sm w-64 focus:w-80 transition-all duration-300 pl-9"
          />
          {searchQuery && (
            <button
              type="button"
              onClick={() => setSearchQuery("")}
              className="absolute right-3 top-1/2 -translate-y-1/2"
            >
              <X className="w-4 h-4 text-base-content/40 hover:text-base-content" />
            </button>
          )}
        </div>

        {/* 搜索建议下拉 */}
        {showSuggestions && suggestions.length > 0 && (
          <div className="absolute top-full left-0 right-0 mt-2 bg-base-100 rounded-lg shadow-lg border border-base-200 z-50 overflow-hidden">
            <div className="px-3 py-2 text-xs text-base-content/50 border-b border-base-200">
              搜索建议
            </div>
            {suggestions.map((suggestion) => (
              <button
                key={suggestion.id}
                onClick={() => handleSuggestionClick(suggestion)}
                className="w-full px-4 py-2 text-left hover:bg-base-200 transition-colors flex items-center gap-2"
              >
                <Search className="w-4 h-4 text-base-content/40" />
                <span className="text-sm">{suggestion.title}</span>
              </button>
            ))}
          </div>
        )}
      </form>
    </>
  );
}