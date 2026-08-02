// hooks/useDeletionStatus.ts
import { useQuery } from "@tanstack/react-query";
import { authApi } from "@/shared/api/modules/auth";

// 账户删除状态
export interface DeletionStatus {
  is_deleted: boolean;
  deleted_at?: string;
  can_restore: boolean;
  remaining_days?: number;
}

// 账户删除状态 Hook（查询 key 固定为 ["auth", "deletion-status"]）
export function useDeletionStatus(enabled: boolean) {
  return useQuery<DeletionStatus>({
    queryKey: ["auth", "deletion-status"],
    queryFn: async () => {
      const res = await authApi.getDeletionStatus();
      const status = res.data.data;
      if (!status) throw new Error("获取删除状态失败");
      return status;
    },
    enabled,
    retry: false,
  });
}
