// hooks/useProfileSettings.ts
import { useCallback } from "react";
import { userApi } from "@/shared/api/modules/user";
import { UpdateProfilePayload } from "@/shared/api/types/user.model";

// 个人资料保存 Hook
export function useProfileSettings() {
  // 更新当前用户资料（仅调用 API，错误由调用方统一处理）
  const updateProfile = useCallback(
    async (data: UpdateProfilePayload): Promise<void> => {
      await userApi.updateProfile(data);
    },
    [],
  );

  return { updateProfile };
}
