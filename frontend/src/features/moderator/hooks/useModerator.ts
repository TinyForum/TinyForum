// hooks/useModerator.ts
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  moderatorApi,
  ApplyModeratorForm,
  BanUserRequest,
  ModeratorBoard,
} from "@/shared/api/modules/moderator";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { useAuthStore } from "@/store";
import { ApiResponse } from "@/shared/api/types/basic.model";

// ========== 类型定义 ==========
interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

interface PinPostParams {
  postId: number;
  pinInBoard: boolean;
}

// ========== 申请相关 ==========
export function useMyApplications(params?: {
  page?: number;
  page_size?: number;
}) {
  return useQuery({
    queryKey: ["moderator", "my-applications", params],
    queryFn: () =>
      moderatorApi
        .getMyApplications(params)
        .then((r: { data: ApiResponse<unknown> }) => r.data.data),
    staleTime: 60 * 1000,
  });
}

export function useApplyModerator(boardId: number) {
  const queryClient = useQueryClient();
  const t = useTranslations("moderator");

  return useMutation({
    mutationFn: (data: ApplyModeratorForm) =>
      moderatorApi.applyModerator(boardId, data),
    onSuccess: () => {
      toast.success(t("application_submitted"));
      queryClient.invalidateQueries({
        queryKey: ["moderator", "my-applications"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err?.response?.data?.message || t("application_failed"));
    },
  });
}

export function useCancelApplication() {
  const queryClient = useQueryClient();
  const t = useTranslations("moderator");

  return useMutation({
    mutationFn: (applicationId: number) =>
      moderatorApi.cancelApplication(applicationId),
    onSuccess: () => {
      toast.success(t("application_canceled"));
      queryClient.invalidateQueries({
        queryKey: ["moderator", "my-applications"],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err?.response?.data?.message || t("cancel_failed"));
    },
  });
}

// ========== 版主操作 ==========
export function useModerators(boardId: number, enabled: boolean = true) {
  return useQuery({
    queryKey: ["moderator", "list", boardId],
    queryFn: () =>
      moderatorApi
        .getModerators(boardId)
        .then((r: { data: ApiResponse<unknown> }) => r.data.data),
    enabled: !!boardId && enabled,
  });
}

export function useBanUser(boardId: number) {
  const queryClient = useQueryClient();
  const t = useTranslations("moderator");

  return useMutation({
    mutationFn: (data: BanUserRequest) => moderatorApi.banUser(boardId, data),
    onSuccess: () => {
      toast.success(t("user_banned"));
      queryClient.invalidateQueries({
        queryKey: ["moderator", "bans", boardId],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err?.response?.data?.message || t("ban_failed"));
    },
  });
}

export function useUnbanUser(boardId: number) {
  const queryClient = useQueryClient();
  const t = useTranslations("moderator");

  return useMutation({
    mutationFn: (userId: number) => moderatorApi.unbanUser(boardId, userId),
    onSuccess: () => {
      toast.success(t("user_unbanned"));
      queryClient.invalidateQueries({
        queryKey: ["moderator", "bans", boardId],
      });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err?.response?.data?.message || t("unban_failed"));
    },
  });
}

export function useDeletePost(boardId: number) {
  const queryClient = useQueryClient();
  const t = useTranslations("moderator");

  return useMutation({
    mutationFn: (postId: number) => moderatorApi.deletePost(boardId, postId),
    onSuccess: () => {
      toast.success(t("post_deleted"));
      queryClient.invalidateQueries({ queryKey: ["posts", boardId] });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err?.response?.data?.message || t("delete_post_failed"));
    },
  });
}

export function usePinPost(boardId: number) {
  const queryClient = useQueryClient();
  const t = useTranslations("moderator");

  return useMutation({
    mutationFn: ({ postId, pinInBoard }: PinPostParams) =>
      moderatorApi.pinPost(boardId, postId, pinInBoard),
    onSuccess: (_: unknown, variables: PinPostParams) => {
      toast.success(
        variables.pinInBoard ? t("post_pinned") : t("post_unpinned"),
      );
      queryClient.invalidateQueries({ queryKey: ["posts", boardId] });
      queryClient.invalidateQueries({ queryKey: ["post", variables.postId] });
    },
    onError: (error: unknown) => {
      const err = error as ErrorResponse;
      toast.error(err?.response?.data?.message || t("pin_failed"));
    },
  });
}

// ========== 权限检查 ==========
export function useModeratorPermissions(boardId: number) {
  const { data: moderators, isLoading } = useModerators(boardId);
  const { user } = useAuthStore();

  // 确保 moderators 是数组
  const moderatorsArray = Array.isArray(moderators) ? moderators : [];

  const currentModerator = moderatorsArray.find(
    (m: { user_id: number }) => m.user_id === user?.id,
  );

  return {
    isModerator: !!currentModerator,
    canDeletePost: currentModerator?.permissions?.can_delete_post || false,
    canPinPost: currentModerator?.permissions?.can_pin_post || false,
    canEditAnyPost: currentModerator?.permissions?.can_edit_any_post || false,
    canManageModerator:
      currentModerator?.permissions?.can_manage_moderator || false,
    canBanUser: currentModerator?.permissions?.can_ban_user || false,
    isLoading, // 可选：返回加载状态
  };
}

export function useMyModeratorBoards() {
  const { user, isHydrated } = useAuthStore();

  return useQuery<ModeratorBoard[]>({
    queryKey: ["moderator", "my-boards", user?.id],
    queryFn: async (): Promise<ModeratorBoard[]> => {
      const response = await moderatorApi.getMyModeratorBoards();
      const data = response.data?.data;
      return Array.isArray(data) ? data : [];
    },
    enabled: !!user?.id && isHydrated,
  });
}
