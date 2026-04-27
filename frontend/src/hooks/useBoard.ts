// hooks/useBoard.ts
import { useState, useEffect, useCallback } from "react";
import { boardApi } from "@/lib/api/modules/boards";
import type { Board } from "@/lib/api/types";
import { toast } from "react-hot-toast";

interface UseBoardOptions {
  autoLoad?: boolean; // 是否自动加载
  page?: number; // 页码
  pageSize?: number; // 每页数量
}

interface UseBoardReturn {
  boards: Board[]; // 板块列表
  loading: boolean; // 加载状态
  error: string | null; // 错误信息
  total: number; // 总记录数
  loadBoards: () => Promise<void>; // 手动加载
  getBoardById: (id: number) => Board | undefined; // 根据ID获取板块
  getDefaultBoard: () => Board | null; // 获取默认板块
  refresh: () => Promise<void>; // 刷新数据
}

export function useBoard(options: UseBoardOptions = {}): UseBoardReturn {
  const { autoLoad = true, page = 1, pageSize = 100 } = options;

  const [boards, setBoards] = useState<Board[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState<number>(0);

  const loadBoards = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await boardApi.list({ page, page_size: pageSize });

      // 响应结构: { code: 0, message: "success", data: { list: [], total, page, page_size } }
      if (response.data.code === 0 || response.data.code === 200) {
        const responseData = response.data.data;

        // 优先检查分页数据格式 { list: [], total, page, page_size }
        if (
          responseData &&
          "list" in responseData &&
          Array.isArray(responseData.list)
        ) {
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
    } catch (err: any) {
      const errorMsg =
        err.response?.data?.message || err.message || "加载板块失败";
      setError(errorMsg);
      console.error("加载板块失败:", err);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize]);

  const refresh = useCallback(async () => {
    await loadBoards();
  }, [loadBoards]);

  const getBoardById = useCallback(
    (id: number): Board | undefined => {
      return boards.find((board) => board.id === id);
    },
    [boards],
  );

  const getDefaultBoard = useCallback((): Board | null => {
    if (boards.length === 0) return null;
    // 返回第一个板块作为默认
    return boards[0];
  }, [boards]);

  useEffect(() => {
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
