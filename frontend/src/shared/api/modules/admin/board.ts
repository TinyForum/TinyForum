import apiClient from "../../client";
import { ApiResponse, PageData, Board } from "../../types";

export const adminBoardsApi = {
  // ── 板块管理 ──────────────────────────────────────────────────────────────
  /** 获取板块列表（分页） */
  listBoards: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Board>>>("/admin/boards", { params }),
};
