// hooks/useBoardCreate.ts
import { useMutation } from "@tanstack/react-query";
import { boardApi } from "@/shared/api/modules/boards";
import { Board, CreateBoardPayload } from "@/shared/api/types/board.model";

// 创建板块
export function useBoardCreate() {
  return useMutation({
    mutationFn: async (data: CreateBoardPayload): Promise<Board | null> => {
      const res = await boardApi.create(data);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "创建失败");
      }
      return res.data.data ?? null;
    },
  });
}
