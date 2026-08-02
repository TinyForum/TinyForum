// hooks/useUpload.ts
import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  uploadApi,
  type FileInfoResponse,
} from "@/shared/api/modules/uploads";
import { pluginApi, PluginListParams } from "@/shared/api/modules/plugin/plugins";
import { PluginMeta } from "@/shared/api/types/plugin.model";
import { userFileKeys } from "@/features/upload/hooks/useUploadKeys";

type UploadType = "post" | "comment" | "plugin";

interface UseUploadReturn {
  isUploading: boolean;
  error: string | null;
  uploadPostFile: (
    postId: string | number,
    file: File,
  ) => Promise<string | null>;
  uploadCommentFile: (
    commentId: string | number,
    file: File,
  ) => Promise<string | null>;
  getUserFiles: (params?: {
    page?: number;
    page_size?: number;
  }) => Promise<FileInfoResponse[]>;
  getFileInfo: (
    type: UploadType,
    fileId: string,
  ) => Promise<FileInfoResponse | null>;
  getUserPluginsList: (params: PluginListParams) => Promise<PluginMeta[]>;
  deleteFile: (type: UploadType, fileId: string) => Promise<boolean>;
  serveFile: (fileId: string) => Promise<Blob | null>;
  resetError: () => void;
}

export function useUpload(): UseUploadReturn {
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const queryClient = useQueryClient();

  const resetError = () => setError(null);

  const setErrorAndReturnNull = (err: unknown) => {
    const message = err instanceof Error ? err.message : "操作失败";
    setError(message);
    return null;
  };

  const uploadPostFile = async (
    postId: string | number,
    file: File,
  ): Promise<string | null> => {
    setIsUploading(true);
    setError(null);
    try {
      const res = await uploadApi.uploadPostFile(postId, file);
      return res.data.data?.url ?? null;
    } catch (err) {
      return setErrorAndReturnNull(err);
    } finally {
      setIsUploading(false);
    }
  };

  const uploadCommentFile = async (
    commentId: string | number,
    file: File,
  ): Promise<string | null> => {
    setIsUploading(true);
    setError(null);
    try {
      const res = await uploadApi.uploadCommentFile(commentId, file);
      return res.data.data?.url ?? null;
    } catch (err) {
      return setErrorAndReturnNull(err);
    } finally {
      setIsUploading(false);
    }
  };

  // 获取用户文件列表
  const getUserFiles = async (params?: {
    page?: number;
    page_size?: number;
  }): Promise<FileInfoResponse[]> => {
    setIsUploading(true);
    setError(null);
    try {
      return await queryClient.fetchQuery({
        queryKey: userFileKeys.filesList(params ?? {}),
        queryFn: async () => {
          const res = await uploadApi.getUserFiles(params);
          return res.data.data?.list ?? [];
        },
      });
    } catch (err) {
      setErrorAndReturnNull(err);
      throw err;
    } finally {
      setIsUploading(false);
    }
  };

  // 获取单个文件信息
  const getFileInfo = async (
    type: UploadType,
    fileId: string,
  ): Promise<FileInfoResponse | null> => {
    setIsUploading(true);
    setError(null);
    try {
      return await queryClient.fetchQuery({
        queryKey: userFileKeys.detail(type, fileId),
        queryFn: async () => {
          if (type === "post") {
            const res = await uploadApi.getPostFile(fileId);
            return res.data.data ?? null;
          }
          if (type === "comment") {
            const res = await uploadApi.getCommentFile(fileId);
            return res.data.data ?? null;
          }
          const res = await uploadApi.getPluginFile(fileId);
          return res.data.data ?? null;
        },
      });
    } catch (err) {
      setErrorAndReturnNull(err);
      throw err;
    } finally {
      setIsUploading(false);
    }
  };

  // 获取用户插件列表
  const getUserPluginsList = async (
    params: PluginListParams,
  ): Promise<PluginMeta[]> => {
    setIsUploading(true);
    setError(null);
    try {
      return await queryClient.fetchQuery({
        queryKey: userFileKeys.pluginsList(params),
        queryFn: async () => {
          const res = await pluginApi.list(params);
          return res.data.data?.list ?? [];
        },
      });
    } catch (err) {
      setErrorAndReturnNull(err);
      throw err;
    } finally {
      setIsUploading(false);
    }
  };

  // 变更：删除文件
  const deleteMutation = useMutation({
    mutationFn: async ({
      type,
      fileId,
    }: {
      type: UploadType;
      fileId: string;
    }) => {
      if (type === "post") await uploadApi.deletePostFile(fileId);
      else if (type === "comment") await uploadApi.deleteCommentFile(fileId);
      else await uploadApi.deletePluginFile(fileId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userFileKeys.all });
    },
  });

  const deleteFile = async (
    type: UploadType,
    fileId: string,
  ): Promise<boolean> => {
    setIsUploading(true);
    setError(null);
    try {
      await deleteMutation.mutateAsync({ type, fileId });
      return true;
    } catch (err) {
      setErrorAndReturnNull(err);
      return false;
    } finally {
      setIsUploading(false);
    }
  };

  const serveFile = async (fileId: string): Promise<Blob | null> => {
    setIsUploading(true);
    setError(null);
    try {
      const res = await uploadApi.serveFile(fileId);
      return res.data;
    } catch (err) {
      setErrorAndReturnNull(err);
      return null;
    } finally {
      setIsUploading(false);
    }
  };

  return {
    isUploading,
    error,
    uploadPostFile,
    uploadCommentFile,

    getUserFiles,
    getFileInfo,
    getUserPluginsList,
    deleteFile,
    serveFile,
    resetError,
  };
}
