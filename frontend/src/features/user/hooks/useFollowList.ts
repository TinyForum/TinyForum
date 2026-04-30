// hooks/user/useUserFollow.ts
import { useState, useCallback } from "react";
import type { User } from "@/shared/api/types";
import { toast } from "react-hot-toast";
import { userApi } from "@/shared/api/modules/user";

interface ErrorResponse {
  response?: { data?: { message?: string } };
  message?: string;
}

// ========== 粉丝/关注列表（分页） ==========
interface UseFollowListReturn {
  users: User[];
  loading: boolean;
  error: string | null;
  total: number;
  page: number;
  loadFollowers: (
    userId: number,
    pageNum?: number,
    pageSize?: number,
  ) => Promise<void>;
  loadFollowing: (
    userId: number,
    pageNum?: number,
    pageSize?: number,
  ) => Promise<void>;
}

export function useFollowList(): UseFollowListReturn {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);

  const loadFollowers = useCallback(
    async (
      userId: number,
      pageNum: number = 1,
      pageSize: number = 20,
    ): Promise<void> => {
      setLoading(true);
      setError(null);
      try {
        const response = await userApi.getFollowers(userId, {
          page: pageNum,
          page_size: pageSize,
        });
        if (response.data.code === 0 && response.data.data) {
          setUsers(response.data.data.list || []);
          setTotal(response.data.data.total || 0);
          setPage(response.data.data.page || pageNum);
        } else {
          throw new Error(response.data.message || "获取粉丝列表失败");
        }
      } catch (err: unknown) {
        const errorObj = err as ErrorResponse;
        const errorMsg =
          errorObj.response?.data?.message ||
          errorObj.message ||
          "获取粉丝列表失败";
        setError(errorMsg);
        toast.error(errorMsg);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const loadFollowing = useCallback(
    async (
      userId: number,
      pageNum: number = 1,
      pageSize: number = 20,
    ): Promise<void> => {
      setLoading(true);
      setError(null);
      try {
        const response = await userApi.getFollowing(userId, {
          page: pageNum,
          page_size: pageSize,
        });
        if (response.data.code === 0 && response.data.data) {
          setUsers(response.data.data.list || []);
          setTotal(response.data.data.total || 0);
          setPage(response.data.data.page || pageNum);
        } else {
          throw new Error(response.data.message || "获取关注列表失败");
        }
      } catch (err: unknown) {
        const errorObj = err as ErrorResponse;
        const errorMsg =
          errorObj.response?.data?.message ||
          errorObj.message ||
          "获取关注列表失败";
        setError(errorMsg);
        toast.error(errorMsg);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return { users, loading, error, total, page, loadFollowers, loadFollowing };
}
