// hooks/useUpload.ts
import { ListPluginRequest, uploadApi } from "@/shared/api/modules/attachments";
import { ApiResponse, PageData } from "@/shared/api/types/basic.model";
import { useState } from "react";
import { PluginMeta } from "@/shared/type/plugin.type";
type UploadType = "post" | "comment" | "plugin";

interface FileInfo {
  file_id: string;
  user_id: number;
  post_id?: number;
  original_name: string;
  stored_name: string;
  stored_path: string;
  size: number;
  mime_type: string;
  file_type: string;
  ext: string;
  status: number;
  created_at: string;
}

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
  }) => Promise<FileInfo>;
  getFileInfo: (type: UploadType, fileId: string) => Promise<FileInfo>;
  getUserPluginsList: (params: ListPluginRequest) => Promise<PluginMeta[]>;
  deleteFile: (type: UploadType, fileId: string) => Promise<boolean>;
  serveFile: (fileId: string) => Promise<Blob | null>;
  resetError: () => void;
}

export function useUpload(): UseUploadReturn {
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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
      return res.data.data;
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
      return res.data.data;
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
  }) => {
    setIsUploading(true);
    setError(null);
    try {
      const res = await uploadApi.getUserFiles(params);
      return res.data;
    } catch (err) {
      setErrorAndReturnNull(err);
      throw err;
    } finally {
      setIsUploading(false);
    }
  };

  const getFileInfo = async (type: UploadType, fileId: string) => {
    setIsUploading(true);
    setError(null);
    try {
      let res;
      if (type === "post") res = await uploadApi.getPostFile(fileId);
      else if (type === "comment") res = await uploadApi.getCommentFile(fileId);
      else res = await uploadApi.getPluginFile(fileId);
      return res.data;
    } catch (err) {
      setErrorAndReturnNull(err);
      throw err;
    } finally {
      setIsUploading(false);
    }
  };

  // 获取用户插件列表
  const getUserPluginsList = async (params: ListPluginRequest) => {
    setIsUploading(true);
    setError(null);
    params.type = "plugin";

    try {
      const res = await uploadApi.getUserPlugins(params);
      return res.data.data?.list || [];
    } catch (err) {
      setErrorAndReturnNull(err);
      throw err;
    } finally {
      setIsUploading(false);
    }
  };

  const deleteFile = async (
    type: UploadType,
    fileId: string,
  ): Promise<boolean> => {
    setIsUploading(true);
    setError(null);
    try {
      if (type === "post") await uploadApi.deletePostFile(fileId);
      else if (type === "comment") await uploadApi.deleteCommentFile(fileId);
      else await uploadApi.deletePluginFile(fileId);
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
