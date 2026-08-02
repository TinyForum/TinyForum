import {
  UploadFileType,
  attachmentApi,
  UploadFileResponse,
} from "@/shared/api/modules/attachments";
import { useMutation, useQueryClient, UseMutationOptions } from "@tanstack/react-query";
import { attachmentKeys } from "./retrieve";
// ========== 变更 Hooks ==========

/**
 * 通用上传文件 Mutation
 */
export function useUploadFile(
  options?: UseMutationOptions<
    UploadFileResponse,
    Error,
    { file: File; type: UploadFileType; post_id?: number; reply_id?: number }
  >,
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (params) => {
      const res = await attachmentApi.uploadFile(params);
      if (!res.data.data) {
        throw new Error("Upload failed");
      }
      return res.data.data;
    },
    onSuccess: (data, variables, onMutateResult, context) => {
      queryClient.invalidateQueries({ queryKey: attachmentKeys.lists() });
      options?.onSuccess?.(data, variables, onMutateResult, context);
    },
    ...options,
  });
}

/**
 * 便捷上传：帖子图片
 */
export function useUploadPostImage(
  options?: UseMutationOptions<
    UploadFileResponse,
    Error,
    { file: File; postId: number }
  >,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ file, postId }) => {
      const res = await attachmentApi.uploadFile({
        file,
        type: "post_image",
        post_id: postId,
      });
      if (!res.data.data) {
        throw new Error("Upload failed");
      }
      return res.data.data;
    },
    onSuccess: (data, variables, onMutateResult, context) => {
      queryClient.invalidateQueries({ queryKey: attachmentKeys.lists() });
      options?.onSuccess?.(data, variables, onMutateResult, context);
    },
    ...options,
  });
}

/**
 * 便捷上传：评论附件
 */
export function useUploadCommentFile(
  options?: UseMutationOptions<
    UploadFileResponse,
    Error,
    { file: File; replyId: number }
  >,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ file, replyId }) => {
      const res = await attachmentApi.uploadFile({
        file,
        type: "comment_file",
        reply_id: replyId,
      });
      if (!res.data.data) {
        throw new Error("Upload failed");
      }
      return res.data.data;
    },
    onSuccess: (data, variables, onMutateResult, context) => {
      queryClient.invalidateQueries({ queryKey: attachmentKeys.lists() });
      options?.onSuccess?.(data, variables, onMutateResult, context);
    },
    ...options,
  });
}

/**
 * 便捷上传：头像
 */
export function useUploadAvatar(
  options?: UseMutationOptions<UploadFileResponse, Error, { file: File }>,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ file }) => {
      const res = await attachmentApi.uploadFile({ file, type: "avatar" });
      if (!res.data.data) {
        throw new Error("Upload failed");
      }
      return res.data.data;
    },
    onSuccess: (data, variables, onMutateResult, context) => {
      queryClient.invalidateQueries({ queryKey: attachmentKeys.lists() });
      options?.onSuccess?.(data, variables, onMutateResult, context);
    },
    ...options,
  });
}

/**
 * 便捷上传：插件
 */
export function useUploadPlugin(
  options?: UseMutationOptions<UploadFileResponse, Error, { file: File }>,
) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ file }) => {
      const res = await attachmentApi.uploadFile({ file, type: "plugin" });
      if (!res.data.data) {
        throw new Error("Upload failed");
      }
      return res.data.data;
    },
    onSuccess: (data, variables, onMutateResult, context) => {
      queryClient.invalidateQueries({ queryKey: attachmentKeys.lists() });
      options?.onSuccess?.(data, variables, onMutateResult, context);
    },
    ...options,
  });
}
