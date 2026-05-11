// features/upload/hooks/useUpload.ts
import { uploadApi } from "@/shared/api/modules/uploads";
import { useState, useCallback } from "react";

// 文件信息类型（根据后端实际结构定义，此处示例）
export interface UserFile {
  id: string;
  filename: string;
  url: string;
  size: number;
  mime_type: string;
  created_at: string;
}

export interface PluginFile extends UserFile {
  file_type: string;
}

// 上传 Hook（适配多种响应格式）
export function useUpload() {
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  const resetUpload = useCallback(() => {
    setIsUploading(false);
    setUploadError(null);
  }, []);

  // 通用上传处理（适用于返回 { data: string } 的接口）
  const handleUpload = useCallback(
    async <T>(uploadFn: () => Promise<T>): Promise<string | null> => {
      setIsUploading(true);
      setUploadError(null);
      try {
        const response = await uploadFn();
        const fileUrl = (response as any).data?.data;
        if (fileUrl && typeof fileUrl === "string") {
          return fileUrl;
        }
        throw new Error("上传响应无效");
      } catch (err: any) {
        const errorMsg =
          err.response?.data?.message || err.message || "上传失败";
        setUploadError(errorMsg);
        return null;
      } finally {
        setIsUploading(false);
      }
    },
    [],
  );

  // 上传帖子文件
  const uploadPostFile = useCallback(
    async (postId: string | number, file: File): Promise<string | null> =>
      handleUpload(() => uploadApi.uploadPostFile(postId, file)),
    [handleUpload],
  );

  // 上传评论文件
  const uploadCommentFile = useCallback(
    async (commentId: string | number, file: File): Promise<string | null> =>
      handleUpload(() => uploadApi.uploadCommentFile(commentId, file)),
    [handleUpload],
  );

  // 上传插件文件
  const uploadPluginFile = useCallback(
    async (file: File, fileType: string = "plugin"): Promise<string | null> =>
      handleUpload(() => uploadApi.uploadPluginFile(file, fileType)),
    [handleUpload],
  );

  // 上传头像（单独处理，因为响应结构不同）
  const uploadAvatar = useCallback(
    async (file: File): Promise<string | null> => {
      setIsUploading(true);
      setUploadError(null);
      try {
        const response = await uploadApi.uploadAvatar(file);
        if (response.status === 200 && response.data.code === 0) {
          const url = response.data.data?.url;
          if (url && typeof url === "string") {
            return url;
          }
          throw new Error("返回的URL无效");
        } else {
          throw new Error(response.data.message || "上传失败");
        }
      } catch (err: any) {
        const errorMsg =
          err.response?.data?.message || err.message || "上传失败";
        setUploadError(errorMsg);
        return null;
      } finally {
        setIsUploading(false);
      }
    },
    [],
  );

  return {
    isUploading,
    uploadError,
    resetUpload,
    uploadPostFile,
    uploadCommentFile,
    uploadPluginFile,
    uploadAvatar,
  };
}

// 用户文件列表管理 Hook
export function useUserFiles() {
  const [files, setFiles] = useState<UserFile[]>([]);
  const [plugins, setPlugins] = useState<PluginFile[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    page: 1,
    pageSize: 20,
    total: 0,
  });

  const fetchUserFiles = useCallback(async (page = 1, pageSize = 20) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await uploadApi.getUserFiles({
        page,
        page_size: pageSize,
      });
      const data = (response as any).data;
      if (response.status === 200) {
        setFiles(data?.items || []);
        setPagination((prev) => ({
          ...prev,
          page,
          pageSize,
          total: data?.total || 0,
        }));
      } else {
        throw new Error(data?.message || "获取文件列表失败");
      }
    } catch (err: any) {
      setError(err.message || "获取文件列表失败");
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchUserPlugins = useCallback(async (page = 1, pageSize = 20) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await uploadApi.getUserPlugins({
        page,
        page_size: pageSize,
      });
      const data = (response as any).data;
      if (response.status === 200) {
        setPlugins(data?.items || []);
        setPagination((prev) => ({
          ...prev,
          page,
          pageSize,
          total: data?.total || 0,
        }));
      } else {
        throw new Error(data?.message || "获取插件列表失败");
      }
    } catch (err: any) {
      setError(err.message || "获取插件列表失败");
    } finally {
      setIsLoading(false);
    }
  }, []);

  const deleteFile = useCallback(
    async (
      fileId: string,
      type: "post" | "comment" | "plugin",
    ): Promise<boolean> => {
      setIsLoading(true);
      setError(null);
      try {
        let response;
        if (type === "post") response = await uploadApi.deletePostFile(fileId);
        else if (type === "comment")
          response = await uploadApi.deleteCommentFile(fileId);
        else response = await uploadApi.deletePluginFile(fileId);

        if (response.status === 200) {
          if (type === "plugin") {
            setPlugins((prev) => prev.filter((p) => p.id !== fileId));
          } else {
            setFiles((prev) => prev.filter((f) => f.id !== fileId));
          }
          return true;
        } else {
          throw new Error((response as any).data?.message || "删除失败");
        }
      } catch (err: any) {
        setError(err.message || "删除失败");
        return false;
      } finally {
        setIsLoading(false);
      }
    },
    [],
  );

  const getFileInfo = useCallback(
    async (
      fileId: string,
      type: "post" | "comment" | "plugin",
    ): Promise<any | null> => {
      setIsLoading(true);
      setError(null);
      try {
        let response;
        if (type === "post") response = await uploadApi.getPostFile(fileId);
        else if (type === "comment")
          response = await uploadApi.getCommentFile(fileId);
        else response = await uploadApi.getPluginFile(fileId);

        if (response.status === 200) {
          return (response as any).data;
        } else {
          throw new Error(
            (response as any).data?.message || "获取文件信息失败",
          );
        }
      } catch (err: any) {
        setError(err.message || "获取文件信息失败");
        return null;
      } finally {
        setIsLoading(false);
      }
    },
    [],
  );

  return {
    files,
    plugins,
    isLoading,
    error,
    pagination,
    fetchUserFiles,
    fetchUserPlugins,
    deleteFile,
    getFileInfo,
  };
}

// 文件预览/下载 Hook（非认证）
export function useFileServe() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const serveFile = useCallback(
    async (
      fileId: string,
      options?: { download?: boolean; filename?: string },
    ) => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await uploadApi.serveFile(fileId);
        const blob = response.data as Blob;
        const url = URL.createObjectURL(blob);
        if (options?.download) {
          const a = document.createElement("a");
          a.href = url;
          a.download = options.filename || fileId;
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          URL.revokeObjectURL(url);
        } else {
          window.open(url, "_blank");
          setTimeout(() => URL.revokeObjectURL(url), 1000);
        }
      } catch (err: any) {
        setError(err.message || "获取文件失败");
      } finally {
        setIsLoading(false);
      }
    },
    [],
  );

  return { serveFile, isLoading, error };
}
