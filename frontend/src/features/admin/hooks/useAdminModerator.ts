// hooks/useAdminModerator.ts
import { adminModeratorApi } from "@/shared/api/modules/admin/moderator";
import {
  moderatorApi,
  ReviewApplicationRequest,
} from "@/shared/api/modules/moderator";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
// import { ApiResponse } from "@/shared/api/types";

// ============ 类型定义 ============
interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

interface ApplicationParams {
  board_id?: number;
  status?: "pending" | "approved" | "rejected";
  page?: number;
  page_size?: number;
}

interface AddModeratorData {
  user_id: number;
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

interface UpdatePermissionsData {
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

interface UpdatePermissionsParams {
  userId: number;
  data: UpdatePermissionsData;
}

interface BanUserData {
  user_id: number;
  reason: string;
  expires_at?: string;
}

interface PinPostParams {
  boardId: number;
  postId: number;
  pinInBoard: boolean;
}

// 管理员 - 版主申请管理
export const useAdminApplications = (params?: ApplicationParams) => {
  return useQuery({
    queryKey: ["admin", "applications", params],
    queryFn: async () => {
      const res = await adminModeratorApi.listApplications(params);
      return res.data.data;
    },
    enabled: !!params?.status,
  });
};

export const useReviewApplication = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      applicationId,
      data,
    }: {
      applicationId: number;
      data: ReviewApplicationRequest;
    }) => adminModeratorApi.reviewApplication(applicationId, data),
    onSuccess: () => {
      toast.success("审批成功");
      queryClient.invalidateQueries({ queryKey: ["admin", "applications"] });
      queryClient.invalidateQueries({
        queryKey: ["boards", "moderators", "apply"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || "审批失败");
    },
  });
};

// 管理员 - 版主任命与管理
export const useAddModerator = (boardId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AddModeratorData) =>
      moderatorApi.addModerator(boardId, data),
    onSuccess: () => {
      toast.success("任命版主成功");
      queryClient.invalidateQueries({
        queryKey: ["boards", boardId, "moderators"],
      });
      queryClient.invalidateQueries({ queryKey: ["moderator", "my-boards"] });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || "任命失败");
    },
  });
};

export const useRemoveModerator = (boardId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userId: number) =>
      moderatorApi.removeModerator(boardId, userId),
    onSuccess: () => {
      toast.success("移除版主成功");
      queryClient.invalidateQueries({
        queryKey: ["boards", boardId, "moderators"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || "移除失败");
    },
  });
};

export const useUpdateModeratorPermissions = (boardId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userId, data }: UpdatePermissionsParams) =>
      moderatorApi.updateModeratorPermissions(boardId, userId, data),
    onSuccess: () => {
      toast.success("更新权限成功");
      queryClient.invalidateQueries({
        queryKey: ["boards", boardId, "moderators"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || "更新权限失败");
    },
  });
};

// 管理员 - 全局内容管理
export const useAdminDeletePost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ boardId, postId }: { boardId: number; postId: number }) =>
      moderatorApi.deletePost(boardId, postId),
    onSuccess: () => {
      toast.success("删除帖子成功");
      queryClient.invalidateQueries({ queryKey: ["moderator", "boards"] });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || "删除失败");
    },
  });
};

export const useAdminPinPost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ boardId, postId, pinInBoard }: PinPostParams) =>
      moderatorApi.pinPost(boardId, postId, pinInBoard),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["moderator", "boards"] });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || "操作失败");
    },
  });
};

// 管理员 - 禁言管理
export const useAdminBanUser = (boardId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BanUserData) => moderatorApi.banUser(boardId, data),
    onSuccess: () => {
      toast.success("禁言用户成功");
      queryClient.invalidateQueries({
        queryKey: ["moderator", "boards", boardId, "bans"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || "禁言失败");
    },
  });
};

export const useAdminUnbanUser = (boardId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userId: number) => moderatorApi.unbanUser(boardId, userId),
    onSuccess: () => {
      toast.success("解除禁言成功");
      queryClient.invalidateQueries({
        queryKey: ["moderator", "boards", boardId, "bans"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || "解除禁言失败");
    },
  });
};

// 管理员 - 获取板块版主列表
export const useAdminModeratorList = (boardId?: number) => {
  return useQuery({
    queryKey: ["boards", boardId, "moderators"],
    queryFn: async () => {
      if (!boardId) return null;
      const res = await moderatorApi.getModerators(boardId);
      return res.data.data;
    },
    enabled: !!boardId,
  });
};

// 管理员 - 获取板块禁言列表
export const useAdminBannedUsers = (
  boardId: number,
  params?: { page?: number; page_size?: number },
) => {
  return useQuery({
    queryKey: ["moderator", "boards", boardId, "bans", params],
    queryFn: async () => {
      const res = await moderatorApi.getBoardBannedUsers(boardId, params);
      return res.data.data;
    },
    enabled: !!boardId,
  });
};

// 管理员 - 获取板块举报列表
export const useAdminReports = (
  boardId: number,
  params?: { page?: number; page_size?: number; status?: string },
) => {
  return useQuery({
    queryKey: ["moderator", "boards", boardId, "reports", params],
    queryFn: async () => {
      const res = await moderatorApi.getBoardReports(boardId, params);
      return res.data.data;
    },
    enabled: !!boardId,
  });
};

// 管理员 - 获取板块帖子列表
export const useAdminBoardPosts = (
  boardId: number,
  params?: {
    page?: number;
    page_size?: number;
    keyword?: string;
    status?: string;
  },
) => {
  return useQuery({
    queryKey: ["moderator", "boards", boardId, "posts", params],
    queryFn: async () => {
      const res = await moderatorApi.getBoardPosts(boardId, params);
      return res.data.data;
    },
    enabled: !!boardId,
  });
};
