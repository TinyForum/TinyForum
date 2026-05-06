// hooks/useUpload.ts
import { uploadApi } from "@/shared/api/modules/uploads";
import { useState } from "react";

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
  uploadPluginFile: (file: File) => Promise<string | null>;
  getUserFiles: (params?: {
    page?: number;
    page_size?: number;
  }) => Promise<any>;
  getFileInfo: (type: UploadType, fileId: string) => Promise<any>;
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

  const uploadPluginFile = async (file: File): Promise<string | null> => {
    setIsUploading(true);
    setError(null);
    try {
      const res = await uploadApi.uploadPluginFile(file);
      return res.data.data;
    } catch (err) {
      return setErrorAndReturnNull(err);
    } finally {
      setIsUploading(false);
    }
  };

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
    uploadPluginFile,
    getUserFiles,
    getFileInfo,
    deleteFile,
    serveFile,
    resetError,
  };
}
