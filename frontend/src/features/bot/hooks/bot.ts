// hooks/useBots.ts
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { toast } from "react-hot-toast";
import {
  BotVO,
  CreateBotRequest,
  UpdateBotRequest,
  RunEventData,
} from "@/shared/api/types/bot.model";
import { botApi } from "@/shared/api/modules/bot";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { botKeys } from "./useBotKeys";

// 工具函数：从未知错误中提取消息（保持不变）
const getErrorMessage = (err: unknown): string => {
  if (err instanceof Error) return err.message;
  if (typeof err === "string") return err;
  if (err && typeof err === "object" && "response" in err) {
    const response = (err as { response?: { data?: { message?: string } } })
      .response;
    if (response?.data?.message) return response.data.message;
  }
  return "发生未知错误";
};

// ==================== 原有 hooks（保持功能不变，仅调整导入） ====================

interface UseBotsOptions {
  autoLoad?: boolean;
  page?: number;
  pageSize?: number;
}

interface UseBotsReturn {
  bots: BotVO[];
  loading: boolean;
  error: string | null;
  total: number;
  loadBots: () => Promise<void>;
  refresh: () => Promise<void>;
}

export function useBots(options: UseBotsOptions = {}): UseBotsReturn {
  const t = useTranslations("Bot");
  const { autoLoad = true, page = 1, pageSize = 20 } = options;

  // 查询：机器人市场列表（分页）
  const query = useQuery({
    queryKey: botKeys.list({ page, pageSize }),
    queryFn: async () => {
      const res = await botApi.list({ page, pageSize });
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("load_bots_failed"));
      }
      if (!res.data.data) {
        throw new Error(t("bot_list_empty"));
      }
      return res.data.data;
    },
    enabled: autoLoad,
  });

  const bots = query.data?.list ?? [];
  const total = query.data?.total ?? 0;
  const loading = query.isLoading;
  const error = query.error ? getErrorMessage(query.error) : null;

  const loadBots = async (): Promise<void> => {
    await query.refetch();
  };

  const refresh = async (): Promise<void> => {
    await query.refetch();
  };

  return { bots, loading, error, total, loadBots, refresh };
}

interface UseMyBotsOptions {
  autoLoad?: boolean;
  page?: number;
  pageSize?: number;
}

interface UseMyBotsReturn {
  bots: BotVO[];
  loading: boolean;
  error: string | null;
  total: number;
  loadMyBots: () => Promise<void>;
  refresh: () => Promise<void>;
}

export function useMyBots(options: UseMyBotsOptions = {}): UseMyBotsReturn {
  const t = useTranslations("Bot");
  const { autoLoad = true, page = 1, pageSize = 20 } = options;

  // 查询：我的机器人列表（分页）
  const query = useQuery({
    queryKey: botKeys.myList({ page, pageSize }),
    queryFn: async () => {
      const res = await botApi.listMy({ page, pageSize });
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("load_my_bots_failed"));
      }
      if (!res.data.data) {
        throw new Error(t("my_bot_list_empty"));
      }
      return res.data.data;
    },
    enabled: autoLoad,
  });

  const bots = query.data?.list ?? [];
  const total = query.data?.total ?? 0;
  const loading = query.isLoading;
  const error = query.error ? getErrorMessage(query.error) : null;

  const loadMyBots = async (): Promise<void> => {
    await query.refetch();
  };

  const refresh = async (): Promise<void> => {
    await query.refetch();
  };

  return { bots, loading, error, total, loadMyBots, refresh };
}

interface UseBotDetailOptions {
  autoLoad?: boolean;
}

interface UseBotDetailReturn {
  bot: BotVO | null;
  loading: boolean;
  error: string | null;
  loadBot: (id: number) => Promise<void>;
  refresh: () => Promise<void>;
  clear: () => void;
}

export function useBotDetail(
  options: UseBotDetailOptions = {},
): UseBotDetailReturn {
  const t = useTranslations("Bot");
  const { autoLoad = false } = options;
  const queryClient = useQueryClient();
  const [currentId, setCurrentId] = useState<number | null>(null);

  // 查询：机器人详情（由 currentId 驱动）
  const query = useQuery({
    queryKey: botKeys.detail(currentId ?? -1),
    queryFn: async () => {
      if (currentId === null) {
        throw new Error(t("missing_bot_id"));
      }
      const res = await botApi.get(currentId);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("get_bot_detail_failed"));
      }
      if (!res.data.data) {
        throw new Error(t("bot_detail_empty"));
      }
      return res.data.data;
    },
    enabled: autoLoad && currentId !== null,
  });

  const bot = query.data ?? null;
  const loading = query.isLoading;
  const error = query.error ? getErrorMessage(query.error) : null;

  const loadBot = async (id: number): Promise<void> => {
    setCurrentId(id);
  };

  const refresh = async (): Promise<void> => {
    if (currentId !== null) {
      await query.refetch();
    }
  };

  const clear = (): void => {
    setCurrentId(null);
    queryClient.removeQueries({ queryKey: botKeys.details() });
  };

  return { bot, loading, error, loadBot, refresh, clear };
}

interface UseBotActionsReturn {
  createBot: (data: CreateBotRequest) => Promise<number | null>;
  updateBot: (id: number, data: UpdateBotRequest) => Promise<boolean>;
  deleteBot: (id: number) => Promise<boolean>;
  runBot: (id: number, eventData?: RunEventData) => Promise<boolean>;
  loading: boolean;
  error: string | null;
}

export function useBotActions(): UseBotActionsReturn {
  const t = useTranslations("Bot");
  const queryClient = useQueryClient();
  const [error, setError] = useState<string | null>(null);

  // 变更：创建机器人
  const createMutation = useMutation({
    mutationFn: async (data: CreateBotRequest): Promise<{ id: number }> => {
      const res = await botApi.create(data);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("create_bot_failed"));
      }
      if (!res.data.data) {
        throw new Error(t("create_bot_empty"));
      }
      return res.data.data;
    },
    onSuccess: () => {
      toast.success(t("create_bot_success"));
      queryClient.invalidateQueries({ queryKey: botKeys.lists() });
      queryClient.invalidateQueries({ queryKey: botKeys.myLists() });
    },
    onError: (err: Error) => {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
    },
  });

  // 变更：更新机器人
  const updateMutation = useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: number;
      data: UpdateBotRequest;
    }): Promise<void> => {
      const res = await botApi.update(id, data);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("update_bot_failed"));
      }
    },
    onSuccess: (_, variables) => {
      toast.success(t("update_bot_success"));
      queryClient.invalidateQueries({ queryKey: botKeys.lists() });
      queryClient.invalidateQueries({ queryKey: botKeys.myLists() });
      queryClient.invalidateQueries({
        queryKey: botKeys.detail(variables.id),
      });
    },
    onError: (err: Error) => {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
    },
  });

  // 变更：删除机器人
  const deleteMutation = useMutation({
    mutationFn: async (id: number): Promise<void> => {
      const res = await botApi.delete(id);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("delete_bot_failed"));
      }
    },
    onSuccess: (_, id) => {
      toast.success(t("delete_bot_success"));
      queryClient.invalidateQueries({ queryKey: botKeys.lists() });
      queryClient.invalidateQueries({ queryKey: botKeys.myLists() });
      queryClient.invalidateQueries({ queryKey: botKeys.detail(id) });
    },
    onError: (err: Error) => {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
    },
  });

  // 变更：手动触发机器人执行
  const runMutation = useMutation({
    mutationFn: async ({
      id,
      eventData,
    }: {
      id: number;
      eventData?: RunEventData;
    }): Promise<void> => {
      const res = await botApi.runNow(id, eventData);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("run_bot_failed"));
      }
    },
    onSuccess: (_, variables) => {
      toast.success(t("run_bot_success"));
      queryClient.invalidateQueries({
        queryKey: botKeys.detail(variables.id),
      });
    },
    onError: (err: Error) => {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
    },
  });

  const createBot = async (data: CreateBotRequest): Promise<number | null> => {
    setError(null);
    try {
      const result = await createMutation.mutateAsync(data);
      return result?.id ?? null;
    } catch {
      return null;
    }
  };

  const updateBot = async (
    id: number,
    data: UpdateBotRequest,
  ): Promise<boolean> => {
    setError(null);
    try {
      await updateMutation.mutateAsync({ id, data });
      return true;
    } catch {
      return false;
    }
  };

  const deleteBot = async (id: number): Promise<boolean> => {
    setError(null);
    try {
      await deleteMutation.mutateAsync(id);
      return true;
    } catch {
      return false;
    }
  };

  const runBot = async (
    id: number,
    eventData?: RunEventData,
  ): Promise<boolean> => {
    setError(null);
    try {
      await runMutation.mutateAsync({ id, eventData });
      return true;
    } catch {
      return false;
    }
  };

  const loading =
    createMutation.isPending ||
    updateMutation.isPending ||
    deleteMutation.isPending ||
    runMutation.isPending;

  return { createBot, updateBot, deleteBot, runBot, loading, error };
}
