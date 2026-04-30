import { userApi } from "@/shared/api/modules/user";
import { RoleResponse } from "@/shared/api/modules/users";
import { useState, useCallback } from "react";
import { ErrorResponse } from "./useUserProfile";

// ========== 用户角色 ==========
interface UseUserRoleReturn {
  role: RoleResponse | null;
  loading: boolean;
  error: string | null;
  loadRole: () => Promise<void>;
  isAdmin: boolean;
  isModerator: boolean;
  isUser: boolean;
}

export function useUserRole(): UseUserRoleReturn {
  const [role, setRole] = useState<RoleResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const loadRole = useCallback(async (): Promise<void> => {
    setLoading(true);
    setError(null);
    try {
      const response = await userApi.getMeRole();
      if (response.data.code === 0 && response.data.data) {
        setRole(response.data.data);
      } else {
        throw new Error(response.data.message || "获取角色失败");
      }
    } catch (err: unknown) {
      const errorObj = err as ErrorResponse;
      const errorMsg =
        errorObj.response?.data?.message || errorObj.message || "获取角色失败";
      setError(errorMsg);
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    role,
    loading,
    error,
    loadRole,
    isAdmin: role?.role === "admin",
    isModerator: role?.role === "moderator",
    isUser: role?.role === "user",
  };
}
