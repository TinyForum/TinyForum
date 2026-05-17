/**
 * api/modules/bots.ts
 * 机器人管理相关 API（全部需登录）
 */

import { NocodeMetadata } from "@/features/bot/noco.type";
import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";
import {
  BotListResponse,
  BotVO,
  CreateBotRequest,
  Flow,
  RunEventData,
  UpdateBotRequest,
} from "../types/bot.model";

// ========== API 方法 ==========
export const botApi = {
  // ---------- 列表查询（注意顺序：/user/me 必须在 /:id 之前）----------
  /** 获取所有机器人列表（机器人市场） */
  list: (params?: { page?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse<BotListResponse>>("/bots", { params }),

  /** 获取当前用户创建的机器人列表（我的机器人） */
  listMy: (params?: { page?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse<BotListResponse>>("/bots/user/me", { params }),

  // ---------- CRUD ----------
  /** 创建机器人 */
  create: (data: CreateBotRequest) =>
    apiClient.post<ApiResponse<{ id: number }>>("/bots", data),

  /** 获取机器人详情 */
  get: (id: number) => apiClient.get<ApiResponse<BotVO>>(`/bots/${id}`),

  /** 更新机器人（只能操作自己的机器人） */
  update: (id: number, data: UpdateBotRequest) =>
    apiClient.put<ApiResponse<null>>(`/bots/${id}`, data),

  /** 删除机器人（只能删除自己的机器人） */
  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/bots/${id}`),

  // ---------- 执行 ----------
  /** 手动触发机器人执行（可附带事件数据） */
  runNow: (id: number, eventData?: RunEventData) =>
    apiClient.post<ApiResponse<{ message: string }>>(
      `/bots/${id}/run`,
      eventData,
    ),

  // ---------- 零代码支持 ----------
  nocode: {
    getMetadata: () =>
      apiClient.get<ApiResponse<NocodeMetadata>>("/bots/nocode/metadata"),

    validateFlow: (data: Flow) =>
      apiClient.post<ApiResponse<{ valid: boolean; errors?: string[] }>>(
        "/bots/nocode/validate",
        data,
      ),
  },
};
