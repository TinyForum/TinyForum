// hooks/useBoardApply.ts
import { useCallback } from "react";
import { useMutation } from "@tanstack/react-query";
import { boardApi } from "@/shared/api/modules/boards";
import { moderatorApi } from "@/shared/api/modules/moderator";
import { adminModeratorApi } from "@/shared/api/modules/admin/moderator";
import { Board } from "@/shared/api/types/board.model";
import {
  ApplyModeratorForm,
  ModeratorApplication,
} from "@/shared/api/types/moderator.model";

// 申请版主页面：封装板块加载、申请状态查询与提交
export function useBoardApply() {
  // 通过 slug 加载板块
  const fetchBoardBySlug = useCallback(
    async (slug: string): Promise<Board | null> => {
      const res = await boardApi.getBySlug(slug);
      return res.data.data ?? null;
    },
    [],
  );

  // 查询某板块的全部申请（用于定位当前用户申请）
  const fetchBoardApplications = useCallback(
    async (boardId: number): Promise<ModeratorApplication[]> => {
      const res = await adminModeratorApi.listApplications({
        board_id: boardId,
        page: 1,
        page_size: 100,
      });
      return res.data.data?.list || [];
    },
    [],
  );

  // 提交版主申请
  const applyModeratorMutation = useMutation({
    mutationFn: ({
      boardId,
      data,
    }: {
      boardId: number;
      data: ApplyModeratorForm;
    }) => moderatorApi.applyModerator(boardId, data),
  });

  return {
    fetchBoardBySlug,
    fetchBoardApplications,
    applyModeratorMutation,
  };
}
