// hooks/useBots.ts
import { useState, useEffect, useCallback } from "react";
import { toast } from "react-hot-toast";
import { botApi } from "@/shared/api/modules/bot";
import {
  BotVO,
  BotListResponse,
  CreateBotRequest,
  UpdateBotRequest,
} from "@/shared/api/types/bot.model";

// 工具函数：从未知错误中提取消息
const getErrorMessage = (err: unknown): string => {
  if (err instanceof Error) return err.message;
  if (typeof err === "string") return err;
  // 检查 Axios 错误结构
  if (err && typeof err === "object" && "response" in err) {
    const response = (err as { response?: { data?: { message?: string } } })
      .response;
    if (response?.data?.message) return response.data.message;
  }
  return "发生未知错误";
};

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
  const { autoLoad = true, page = 1, pageSize = 20 } = options;

  const [bots, setBots] = useState<BotVO[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState<number>(0);

  const loadBots = useCallback(async (): Promise<void> => {
    setLoading(true);
    setError(null);

    try {
      const response = await botApi.list({ page, pageSize });
      if (response.data.code === 0) {
        const data = response.data.data as BotListResponse;
        setBots(data.list || []);
        setTotal(data.total || 0);
      } else {
        throw new Error(response.data.message || "加载机器人列表失败");
      }
    } catch (err: unknown) {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize]);

  const refresh = useCallback(() => loadBots(), [loadBots]);

  useEffect(() => {
    if (autoLoad) {
      loadBots();
    }
  }, [autoLoad, loadBots]);

  return { bots, loading, error, total, loadBots, refresh };
}

// hooks/useMyBots.ts
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
  const { autoLoad = true, page = 1, pageSize = 20 } = options;

  const [bots, setBots] = useState<BotVO[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState<number>(0);

  const loadMyBots = useCallback(async (): Promise<void> => {
    setLoading(true);
    setError(null);

    try {
      const response = await botApi.listMy({ page, pageSize });
      if (response.data.code === 0) {
        const data = response.data.data as BotListResponse;
        setBots(data.list || []);
        setTotal(data.total || 0);
      } else {
        throw new Error(response.data.message || "加载我的机器人失败");
      }
    } catch (err: unknown) {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize]);

  const refresh = useCallback(() => loadMyBots(), [loadMyBots]);

  useEffect(() => {
    if (autoLoad) {
      loadMyBots();
    }
  }, [autoLoad, loadMyBots]);

  return { bots, loading, error, total, loadMyBots, refresh };
}

// hooks/useBotDetail.ts
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
  const { autoLoad = false } = options;

  const [bot, setBot] = useState<BotVO | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [currentId, setCurrentId] = useState<number | null>(null);

  const loadBot = useCallback(async (id: number): Promise<void> => {
    setLoading(true);
    setError(null);
    setCurrentId(id);

    try {
      const response = await botApi.get(id);
      if (response.data.code === 0) {
        setBot(response.data.data as BotVO);
      } else {
        throw new Error(response.data.message || "获取机器人详情失败");
      }
    } catch (err: unknown) {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
      setBot(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const refresh = useCallback(async (): Promise<void> => {
    if (currentId !== null) {
      await loadBot(currentId);
    }
  }, [loadBot, currentId]);

  const clear = useCallback(() => {
    setBot(null);
    setError(null);
    setCurrentId(null);
  }, []);

  useEffect(() => {
    if (autoLoad && currentId !== null) {
      loadBot(currentId);
    }
  }, [autoLoad, currentId, loadBot]);

  return { bot, loading, error, loadBot, refresh, clear };
}

// hooks/useBotActions.ts
interface UseBotActionsReturn {
  createBot: (data: CreateBotRequest) => Promise<number | null>;
  updateBot: (id: number, data: UpdateBotRequest) => Promise<boolean>;
  deleteBot: (id: number) => Promise<boolean>;
  runBot: (id: number, eventData?: Record<string, unknown>) => Promise<boolean>;
  loading: boolean;
  error: string | null;
}

export function useBotActions(): UseBotActionsReturn {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const createBot = useCallback(
    async (data: CreateBotRequest): Promise<number | null> => {
      setLoading(true);
      setError(null);
      try {
        const response = await botApi.create(data);
        if (response.data.code === 0) {
          const id = (response.data.data as { id: number }).id;
          toast.success("机器人创建成功");
          return id;
        } else {
          throw new Error(response.data.message || "创建机器人失败");
        }
      } catch (err: unknown) {
        const errorMsg = getErrorMessage(err);
        setError(errorMsg);
        toast.error(errorMsg);
        return null;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const updateBot = useCallback(
    async (id: number, data: UpdateBotRequest): Promise<boolean> => {
      setLoading(true);
      setError(null);
      try {
        const response = await botApi.update(id, data);
        if (response.data.code === 0) {
          toast.success("机器人更新成功");
          return true;
        } else {
          throw new Error(response.data.message || "更新机器人失败");
        }
      } catch (err: unknown) {
        const errorMsg = getErrorMessage(err);
        setError(errorMsg);
        toast.error(errorMsg);
        return false;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const deleteBot = useCallback(async (id: number): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      const response = await botApi.delete(id);
      if (response.data.code === 0) {
        toast.success("机器人删除成功");
        return true;
      } else {
        throw new Error(response.data.message || "删除机器人失败");
      }
    } catch (err: unknown) {
      const errorMsg = getErrorMessage(err);
      setError(errorMsg);
      toast.error(errorMsg);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const runBot = useCallback(
    async (
      id: number,
      eventData?: Record<string, unknown>,
    ): Promise<boolean> => {
      setLoading(true);
      setError(null);
      try {
        const response = await botApi.runNow(id, eventData);
        if (response.data.code === 0) {
          toast.success("机器人已触发执行");
          return true;
        } else {
          throw new Error(response.data.message || "触发机器人失败");
        }
      } catch (err: unknown) {
        const errorMsg = getErrorMessage(err);
        setError(errorMsg);
        toast.error(errorMsg);
        return false;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return { createBot, updateBot, deleteBot, runBot, loading, error };
}
