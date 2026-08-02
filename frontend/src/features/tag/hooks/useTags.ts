// hooks/useTags.ts
import { useState, useCallback } from 'react'
import { useQuery } from '@tanstack/react-query'
import { toast } from 'react-hot-toast'
import { tagKeys } from './useTagKeys'
import { tagApi } from '@/shared/api/modules/tags'
import { Tag } from '@/shared/api/types/tag.model'

export function useTags() {
  const [selectedTags, setSelectedTags] = useState<number[]>([])

  // 查询：获取全部标签
  const query = useQuery<Tag[]>({
    queryKey: tagKeys.all,
    queryFn: async () => {
      const res = await tagApi.list()
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '加载标签失败')
      }
      return res.data.data ?? []
    },
  })

  const tags = query.data ?? []

  const toggleTag = useCallback((tagId: number): void => {
    if (tagId === 0) {
      toast.error('无效的标签')
      return
    }

    setSelectedTags((prev: number[]): number[] =>
      prev.includes(tagId)
        ? prev.filter((id: number): boolean => id !== tagId)
        : [...prev, tagId],
    )
  }, [])

  const clearSelectedTags = useCallback((): void => {
    setSelectedTags([])
  }, [])

  return {
    tags,
    selectedTags,
    loading: query.isLoading,
    toggleTag,
    clearSelectedTags,
    reloadTags: query.refetch,
  }
}
