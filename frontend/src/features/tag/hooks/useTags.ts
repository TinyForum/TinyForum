// hooks/useTags.ts
import { useState, useEffect, useCallback } from "react";
import { Tag, ApiResponse } from "@/shared/api/types";
import { toast } from "react-hot-toast";
import { tagApi } from "@/shared/api";

export function useTags() {
  const [tags, setTags] = useState<Tag[]>([]);
  const [selectedTags, setSelectedTags] = useState<number[]>([]);
  const [loading, setLoading] = useState<boolean>(false);

  const loadTags = useCallback(async (): Promise<void> => {
    setLoading(true);
    try {
      const response: { data: ApiResponse<Tag[]> } = await tagApi.list();
      // 修复：后端返回的 code 应该是 0 表示成功，不是 200
      if (response.data.code === 0) {
        setTags(response.data.data || []);
      } else {
        toast.error(response.data.message || "加载标签失败");
      }
    } catch (error: unknown) {
      console.error("Failed to load tags:", error);
      toast.error("加载标签失败");
    } finally {
      setLoading(false);
    }
  }, []);

  const toggleTag = useCallback((tagId: number): void => {
    if (tagId === 0) {
      toast.error("无效的标签");
      return;
    }

    setSelectedTags((prev: number[]): number[] =>
      prev.includes(tagId)
        ? prev.filter((id: number): boolean => id !== tagId)
        : [...prev, tagId],
    );
  }, []);

  const clearSelectedTags = useCallback((): void => {
    setSelectedTags([]);
  }, []);

  useEffect((): void => {
    loadTags();
  }, [loadTags]);

  return {
    tags,
    selectedTags,
    loading,
    toggleTag,
    clearSelectedTags,
    reloadTags: loadTags,
  };
}
