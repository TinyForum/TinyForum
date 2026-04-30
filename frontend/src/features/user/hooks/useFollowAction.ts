import { userApi } from "@/shared/api/modules/user";
import { useState, useCallback } from "react";
import toast from "react-hot-toast";
import { ErrorResponse } from "./useUserProfile";

// ========== 关注/取消关注 ==========
interface UseFollowReturn {
  following: boolean;
  loading: boolean;
  follow: (userId: number) => Promise<boolean>;
  unfollow: (userId: number) => Promise<boolean>;
  checkFollowStatus: (userId: number) => Promise<boolean>;
}

export function useFollowAction(): UseFollowReturn {
  const [loading, setLoading] = useState<boolean>(false);
  const [following, setFollowing] = useState<boolean>(false);

  const follow = useCallback(async (userId: number): Promise<boolean> => {
    setLoading(true);
    try {
      const response = await userApi.follow(userId);
      if (response.data.code === 0) {
        setFollowing(true);
        toast.success("关注成功");
        return true;
      } else {
        throw new Error(response.data.message || "关注失败");
      }
    } catch (err: unknown) {
      const errorObj = err as ErrorResponse;
      const errorMsg =
        errorObj.response?.data?.message || errorObj.message || "关注失败";
      toast.error(errorMsg);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const unfollow = useCallback(async (userId: number): Promise<boolean> => {
    setLoading(true);
    try {
      const response = await userApi.unfollow(userId);
      if (response.data.code === 0) {
        setFollowing(false);
        toast.success("取消关注成功");
        return true;
      } else {
        throw new Error(response.data.message || "取消关注失败");
      }
    } catch (err: unknown) {
      const errorObj = err as ErrorResponse;
      const errorMsg =
        errorObj.response?.data?.message || errorObj.message || "取消关注失败";
      toast.error(errorMsg);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const checkFollowStatus = useCallback(
    async (userId: number): Promise<boolean> => {
      // 需要真实接口时可调用 userApi.getFollowing 检查
      console.log("check follow status", userId, following);
      return following;
    },
    [following],
  );

  return { following, loading, follow, unfollow, checkFollowStatus };
}
