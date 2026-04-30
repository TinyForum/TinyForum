// hooks/user/useUserInfo.ts
import { useState, useCallback } from "react";
import { UpdateProfilePayload } from "@/shared/api/modules/user";
import type { User } from "@/shared/api/types";
import { toast } from "react-hot-toast";
import { userApi } from "@/shared/api/modules/user";

export interface ErrorResponse {
  response?: { data?: { message?: string } };
  message?: string;
}

// ========== 用户资料 ==========
interface UseProfileReturn {
  user: User | null;
  loading: boolean;
  error: string | null;
  loadProfile: (id: number) => Promise<void>;
  updateProfile: (data: UpdateProfilePayload) => Promise<boolean>;
}

export function useUserProfile(): UseProfileReturn {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const loadProfile = useCallback(async (id: number): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const response = await userApi.getProfile(id);
      if (response.data.code === 0 && response.data.data) {
        setUser(response.data.data);
      } else {
        throw new Error(response.data.message || "获取用户信息失败");
      }
    } catch (err: unknown) {
      const errorObj = err as ErrorResponse;
      const errorMsg =
        errorObj.response?.data?.message ||
        errorObj.message ||
        "获取用户信息失败";
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  }, []);

  const updateProfile = useCallback(
    async (data: UpdateProfilePayload): Promise<boolean> => {
      setLoading(true);
      try {
        const response = await userApi.updateProfile(data);
        if (response.data.code === 0 && response.data.data) {
          setUser(response.data.data);
          toast.success("资料更新成功");
          return true;
        } else {
          throw new Error(response.data.message || "更新失败");
        }
      } catch (err: unknown) {
        const errorObj = err as ErrorResponse;
        const errorMsg =
          errorObj.response?.data?.message || errorObj.message || "更新失败";
        toast.error(errorMsg);
        return false;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return { user, loading, error, loadProfile, updateProfile };
}
