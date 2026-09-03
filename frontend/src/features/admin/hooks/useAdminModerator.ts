// hooks/useAdminModerator.ts
import { adminModeratorApi } from "@/shared/api/modules/admin/moderator";
import { moderatorApi } from "@/shared/api/modules/moderator";
import { ReviewApplicationRequest } from "@/shared/api/types/moderator.model";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { adminKeys } from "./useAdminKeys";
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
    queryKey: adminKeys.applications(params),
    queryFn: async () => {
      const res = await adminModeratorApi.listApplications(params);
      return res.data.data;
    },
    enabled: !!params?.status,
  });
};

export const useReviewApplication = () => {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  return useMutation({
    mutationFn: ({
      applicationId,
      data,
    }: {
      applicationId: number;
      data: ReviewApplicationRequest;
    }) => adminModeratorApi.reviewApplication(applicationId, data),
    onSuccess: () => {
      toast.success(t("application_approved"));
      queryClient.invalidateQueries({ queryKey: adminKeys.applications() });
      queryClient.invalidateQueries({
        queryKey: ["boards", "moderators", "apply"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || t("application_review_failed"));
    },
  });
};

// 管理员 - 版主任命与管理
export const useAddModerator = (boardId: number) => {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  return useMutation({
    mutationFn: (data: AddModeratorData) =>
      moderatorApi.addModerator(boardId, data),
    onSuccess: () => {
      toast.success(t("moderator_appointed"));
      queryClient.invalidateQueries({
        queryKey: ["boards", boardId, "moderators"],
      });
      queryClient.invalidateQueries({ queryKey: ["moderator", "my-boards"] });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || t("moderator_appoint_failed"));
    },
  });
};

export const useRemoveModerator = (boardId: number) => {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  return useMutation({
    mutationFn: (userId: number) =>
      moderatorApi.removeModerator(boardId, userId),
    onSuccess: () => {
      toast.success(t("moderator_removed_success"));
      queryClient.invalidateQueries({
        queryKey: ["boards", boardId, "moderators"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || t("moderator_remove_failed"));
    },
  });
};

export const useUpdateModeratorPermissions = (boardId: number) => {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  return useMutation({
    mutationFn: ({ userId, data }: UpdatePermissionsParams) =>
      moderatorApi.updateModeratorPermissions(boardId, userId, data),
    onSuccess: () => {
      toast.success(t("permissions_updated_success"));
      queryClient.invalidateQueries({
        queryKey: ["boards", boardId, "moderators"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || t("permissions_update_failed"));
    },
  });
};

// 管理员 - 全局内容管理
export const useAdminDeletePost = () => {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  return useMutation({
    mutationFn: ({ boardId, postId }: { boardId: number; postId: number }) =>
      moderatorApi.deletePost(boardId, postId),
    onSuccess: () => {
      toast.success(t("post_deleted_success"));
      queryClient.invalidateQueries({ queryKey: ["moderator", "boards"] });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || t("post_delete_failed"));
    },
  });
};

export const useAdminPinPost = () => {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  return useMutation({
    mutationFn: ({ boardId, postId, pinInBoard }: PinPostParams) =>
      moderatorApi.pinPost(boardId, postId, pinInBoard),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["moderator", "boards"] });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || t("operation_failed"));
    },
  });
};

// 管理员 - 禁言管理
export const useAdminBanUser = (boardId: number) => {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  return useMutation({
    mutationFn: (data: BanUserData) => moderatorApi.banUser(boardId, data),
    onSuccess: () => {
      toast.success(t("user_banned_success"));
      queryClient.invalidateQueries({
        queryKey: ["moderator", "boards", boardId, "bans"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || t("user_ban_failed"));
    },
  });
};

export const useAdminUnbanUser = (boardId: number) => {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  return useMutation({
    mutationFn: (userId: number) => moderatorApi.unbanUser(boardId, userId),
    onSuccess: () => {
      toast.success(t("user_unbanned_success"));
      queryClient.invalidateQueries({
        queryKey: ["moderator", "boards", boardId, "bans"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err.response?.data?.message || t("user_unban_failed"));
    },
  });
};

/**
 * 管理员 - 获取板块版主列表
 * @param boardId
 * @returns
 */
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
