// hooks/useTags.ts
import { useState, useEffect, useCallback } from 'react';
import { tagApi } from '@/lib/api';
import { Tag } from '@/lib/api/types';
import { toast } from 'react-hot-toast';

export function useTags() {
  const [tags, setTags] = useState<Tag[]>([]);
  const [selectedTags, setSelectedTags] = useState<number[]>([]);
  const [loading, setLoading] = useState(false);

  const loadTags = useCallback(async () => {
    setLoading(true);
    try {
      const response = await tagApi.list();
      if (response.data.code === 200) {
        setTags(response.data.data);
      }
    } catch (error) {
      console.error('Failed to load tags:', error);
      toast.error('加载标签失败');
    } finally {
      setLoading(false);
    }
  }, []);

  const toggleTag = useCallback((tagId: number) => {
    if (tagId === 0) {
      toast.error('无效的标签');
      return;
    }
    
    setSelectedTags(prev =>
      prev.includes(tagId)
        ? prev.filter(id => id !== tagId)
        : [...prev, tagId]
    );
  }, []);

  const clearSelectedTags = useCallback(() => {
    setSelectedTags([]);
  }, []);

  useEffect(() => {
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