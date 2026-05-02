// hooks/admin/useAdminBoards.ts
import { useQuery } from "@tanstack/react-query";
import { adminBoardsApi } from "@/shared/api/modules/admin/board";
import { PageData } from "@/shared/api/types/basic.model";
import { Board } from "@/shared/api/types/board.model";

const adminBoardsKeys = {
  all: ["admin", "boards"] as const,
  lists: () => [...adminBoardsKeys.all, "list"] as const,
  list: (params: object) => [...adminBoardsKeys.lists(), params] as const,
};

// ========== 获取板块列表 ==========
export function useAdminBoards(params?: { page?: number; page_size?: number }) {
  return useQuery({
    queryKey: adminBoardsKeys.list(params || {}),
    queryFn: async () => {
      const res = await adminBoardsApi.listBoards(params);
      if (res.data.code !== 0)
        throw new Error(res.data.message || "获取板块列表失败");
      return res.data.data as PageData<Board>;
    },
  });
}
