// hooks/useBoard.ts
import { useState, useEffect, useCallback } from "react";
import { boardApi } from "@/lib/api/modules/boards";
import type { Board, ApiResponse } from "@/lib/api/types";
import { toast } from "react-hot-toast";

interface UseBoardOptions {
  autoLoad?: boolean;
  page?: number;
  pageSize?: number;
}

interface BoardListResponse {
  list: Board[];
  total: number;
  page: number;
  page_size: number;
}

interface UseBoardReturn {
  boards: Board[];
  loading: boolean;
  error: string | null;
  total: number;
  loadBoards: () => Promise<void>;
  getBoardById: (id: number) => Board | undefined;
  getDefaultBoard: () => Board | null;
  refresh: () => Promise<void>;
}

interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

export function useBoard(options: UseBoardOptions = {}): UseBoardReturn {
  const { autoLoad = true, page = 1, pageSize = 100 } = options;

  const [boards, setBoards] = useState<Board[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState<number>(0);

  const loadBoards = useCallback(async (): Promise<void> => {
    setLoading(true);
    setError(null);

    try {
      const response: { data: ApiResponse<BoardListResponse | Board[]> } = 
        await boardApi.list({ page, page_size: pageSize });

      // 统一使用 code === 0
      if (response.data.code === 0) {
        const responseData = response.data.data;

        // 优先检查分页数据格式 { list: [], total, page, page_size }
        if (responseData && "list" in responseData && Array.isArray(responseData.list)) {
          setBoards(responseData.list);
          setTotal(responseData.total || 0);
          console.log("加载板块成功:", responseData.list.length, "个板块");
        }
        // 处理直接返回数组的情况
        else if (Array.isArray(responseData)) {
          setBoards(responseData);
          setTotal(responseData.length);
          console.log("加载板块成功(数组):", responseData.length, "个板块");
        }
        // 处理空数据
        else if (!responseData) {
          setBoards([]);
          setTotal(0);
          console.log("没有板块数据");
        } else {
          console.warn("未知的数据格式:", responseData);
          setBoards([]);
          setTotal(0);
        }
      } else {
        throw new Error(response.data.message || "加载板块失败");
      }
    } catch (err: unknown) {
      const errorObj = err as ErrorResponse;
      const errorMsg = errorObj.response?.data?.message || errorObj.message || "加载板块失败";
      setError(errorMsg);
      console.error("加载板块失败:", err);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize]);

  const refresh = useCallback(async (): Promise<void> => {
    await loadBoards();
  }, [loadBoards]);

  const getBoardById = useCallback(
    (id: number): Board | undefined => {
      return boards.find((board: Board): boolean => board.id === id);
    },
    [boards],
  );

  const getDefaultBoard = useCallback((): Board | null => {
    if (boards.length === 0) return null;
    return boards[0];
  }, [boards]);

   
  useEffect((): void => {
    if (autoLoad) {
      loadBoards();
    }
  }, [autoLoad, loadBoards]);

  return {
    boards,
    loading,
    error,
    total,
    loadBoards,
    getBoardById,
    getDefaultBoard,
    refresh,
  };
}