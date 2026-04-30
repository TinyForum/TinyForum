import { userApi } from "@/shared/api/modules/user";
import { LeaderboardItemResponse } from "@/shared/api/modules/users";
import { useState, useCallback } from "react";
import toast from "react-hot-toast";
import { ErrorResponse } from "./useUserProfile";

// ========== 排行榜 ==========
interface UseLeaderboardReturn {
  data: LeaderboardItemResponse[];
  loading: boolean;
  error: string | null;
  loadLeaderboard: (simple?: boolean, limit?: number) => Promise<void>;
}

export function useLeaderboard(): UseLeaderboardReturn {
  const [data, setData] = useState<LeaderboardItemResponse[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const loadLeaderboard = useCallback(
    async (simple: boolean = true, limit: number = 100): Promise<void> => {
      setLoading(true);
      setError(null);
      try {
        const response = simple
          ? await userApi.getLeaderboardSimple({ limit })
          : await userApi.getLeaderboardDetail({ limit });
        if (response.data.code === 0 && response.data.data) {
          setData(response.data.data);
        } else {
          throw new Error(response.data.message || "加载排行榜失败");
        }
      } catch (err: unknown) {
        const errorObj = err as ErrorResponse;
        const errorMsg =
          errorObj.response?.data?.message ||
          errorObj.message ||
          "加载排行榜失败";
        setError(errorMsg);
        toast.error(errorMsg);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return { data, loading, error, loadLeaderboard };
}
