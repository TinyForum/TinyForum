// hooks/useBoardCreationApply.ts
import { useMutation } from "@tanstack/react-query";
import { boardApi } from "@/shared/api/modules/boards";
import { ApplyBoardCreationPayload } from "@/shared/api/types/board.model";

// 提交板块创建申请（普通用户）
export function useBoardCreationApply() {
  return useMutation({
    mutationFn: async (data: ApplyBoardCreationPayload) => {
      const res = await boardApi.applyBoardCreation(data);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "apply_failed");
      }
      return res.data.data;
    },
  });
}
