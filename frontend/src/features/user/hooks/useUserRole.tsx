import { userApi } from "@/shared/api/modules/user";
import { RoleResponse } from "@/shared/api/modules/users";
import { useState, useCallback } from "react";
import { ErrorResponse } from "./useUserProfile";

// ========== 用户角色 ==========
interface UseUserRoleReturn {
  role: RoleResponse | null;
  isLoading: boolean;
  error: string | null;
  loadRole: () => Promise<void>;
  isAdmin: boolean;
  isModerator: boolean;
  isUser: boolean;
  isMember: boolean;
  isSuperAdmin: boolean;
  isReviewer: boolean;
  isSystemMaintainer: boolean;
}

export function useUserRole(): UseUserRoleReturn {
  const [role, setRole] = useState<RoleResponse | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const loadRole = useCallback(async (): Promise<void> => {
    setIsLoading(true);
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
      setIsLoading(false);
    }
  }, []);

  return {
    role,
    isLoading,
    error,
    loadRole,
    isAdmin: role?.role === "admin",
    isModerator: role?.role === "moderator",
    isUser: role?.role === "user",
    isMember: role?.role === "member",
    isSuperAdmin: role?.role === "super_admin",
    isReviewer: role?.role === "reviewer",
    isSystemMaintainer: role?.role === "system_maintainer",
  };
}
