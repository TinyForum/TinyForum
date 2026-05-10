/**
 * api/modules/bots.ts
 * 机器人管理相关 API
 */

import apiClient from "../client";
import { ApiResponse, PageData } from "../types/basic.model";
import {
  BotListResponse,
  BotVO,
  CreateBotRequest,
  UpdateBotRequest,
} from "../types/bot.model";

/** 手动触发时的事件数据（任意对象） */
export type RunEventData = Record<string, any>;

// ========== API 方法 ==========
export const botApi = {
  // ---------- 创建/更新/删除 ----------
  /** 创建机器人（需登录） */
  create: (data: CreateBotRequest) =>
    apiClient.post<ApiResponse<{ id: number }>>("/bots", data),

  /** 更新机器人（只能操作自己的机器人） */
  update: (id: number, data: UpdateBotRequest) =>
    apiClient.put<ApiResponse<null>>(`/bots/${id}`, data),

  /** 删除机器人（只能删除自己的机器人） */
  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/bots/${id}`),

  /** 获取机器人详情（需要认证） */
  get: (id: number) => apiClient.get<ApiResponse<BotVO>>(`/bots/${id}`),

  /** 手动触发机器人执行（可附带事件数据） */
  runNow: (id: number, eventData?: RunEventData) =>
    apiClient.post<ApiResponse<{ message: string }>>(
      `/bots/${id}/run`,
      eventData,
    ),

  // ---------- 列表查询 ----------
  /** 获取所有机器人列表（机器人市场）*/
  list: (params?: { page?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse<BotListResponse>>("/bots", { params }),

  /** 获取当前用户创建的机器人列表（我的机器人）*/
  listMy: (params?: { page?: number; pageSize?: number }) =>
    apiClient.get<ApiResponse<BotListResponse>>("/bots/user/me", { params }),
};
