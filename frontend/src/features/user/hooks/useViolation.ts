import { useCallback } from "react"
import {
  useQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query"
import {
  userViolationApi,
  ViolationVO,
} from "@/shared/api/modules/user/violation"
import { userKeys } from "./useUserKeys"

interface UseViolationReturn {
  violations: ViolationVO[] // ✅ 始终为数组，不是 null
  loadViolations: () => Promise<void>
  isLoading: boolean
  error: string | null
  fetchViolationDetail: (id: string) => Promise<ViolationVO | null> // ✅ 可能返回 null
  submitAppeal: (id: string, reason: string) => Promise<boolean>
  isAppealing: boolean
  appealError: string | null
}

export function useUserViolation(): UseViolationReturn {
  const queryClient = useQueryClient()

  // 查询：拉取当前用户的违规列表
  const { data, isLoading, error, refetch } = useQuery<ViolationVO[]>({
    queryKey: userKeys.violations(),
    queryFn: async () => {
      const res = await userViolationApi.listUserViolations()
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取违规列表失败")
      }
      return res.data.data ?? []
    },
  })

  // 命令式入口：重新拉取
  const loadViolations = useCallback(async (): Promise<void> => {
    await refetch()
  }, [refetch])

  // 查询详情：通过 fetchQuery 按需拉取
  const fetchViolationDetail = useCallback(
    async (id: string): Promise<ViolationVO | null> => {
      try {
        return await queryClient.fetchQuery({
          queryKey: userKeys.violationDetail(id),
          queryFn: async () => {
            const res = await userViolationApi.getViolationDetail(id)
            if (res.data.code !== 0) {
              throw new Error(res.data.message || "获取详情失败")
            }
            if (!res.data.data) {
              throw new Error("获取详情失败")
            }
            return res.data.data
          },
        })
      } catch {
        return null // ✅ 明确返回 null 表示失败
      }
    },
    [queryClient],
  )

  // 变更：提交申诉
  const submitAppealMutation = useMutation<
    boolean,
    Error,
    { id: string; reason: string }
  >({
    mutationFn: async ({ id, reason }) => {
      const res = await userViolationApi.appeal(id, reason)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "申诉提交失败")
      }
      return true
    },
    onSuccess: () => {
      // 申诉成功后刷新违规列表
      queryClient.invalidateQueries({ queryKey: userKeys.violations() })
    },
  })

  const submitAppeal = useCallback(
    async (id: string, reason: string): Promise<boolean> => {
      try {
        await submitAppealMutation.mutateAsync({ id, reason })
        return true
      } catch {
        return false
      }
    },
    [submitAppealMutation],
  )

  return {
    violations: data ?? [],
    loadViolations,
    isLoading,
    error: error?.message ?? null,
    fetchViolationDetail,
    submitAppeal,
    isAppealing: submitAppealMutation.isPending,
    appealError: submitAppealMutation.error?.message ?? null,
  }
}
