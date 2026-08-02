// hooks/attachment.ts
import {
  UploadFileType,
  UserFileListResponse,
  attachmentApi,
  FileInfoResponse,
} from "@/shared/api/modules/attachments";
import { useQuery, UseQueryOptions } from "@tanstack/react-query";

// ========== 查询键（Query Keys） ==========
export const attachmentKeys = {
  all: ["attachments"] as const,
  lists: () => [...attachmentKeys.all, "list"] as const,
  list: (params: {
    page?: number;
    pageSize?: number;
    fileType?: UploadFileType;
  }) => [...attachmentKeys.lists(), params] as const,
  details: () => [...attachmentKeys.all, "detail"] as const,
  detail: (fileId: string) => [...attachmentKeys.details(), fileId] as const,
  file: (fileId: string) => [...attachmentKeys.all, "file", fileId] as const,
};

// ========== 查询 Hooks ==========

/**
 * 获取当前用户的文件列表（分页）
 */
export function useListMyFiles(
  params?: { page?: number; pageSize?: number; fileType?: UploadFileType },
  options?: Omit<
    UseQueryOptions<UserFileListResponse, Error>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery({
    queryKey: attachmentKeys.list({
      page: params?.page || 1,
      pageSize: params?.pageSize || 20,
      fileType: params?.fileType,
    }),
    queryFn: async () => {
      const res = await attachmentApi.listMyFiles({
        page: params?.page,
        page_size: params?.pageSize,
        file_type: params?.fileType,
      });
      if (!res.data.data) {
        throw new Error("Failed to fetch file list");
      }
      return res.data.data;
    },
    ...options,
  });
}

/**
 * 获取文件元信息
 */
export function useGetFile(
  fileId: string,
  options?: Omit<
    UseQueryOptions<FileInfoResponse, Error>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery({
    queryKey: attachmentKeys.detail(fileId),
    queryFn: async () => {
      const res = await attachmentApi.getFile(fileId);
      if (!res.data.data) {
        throw new Error("File not found");
      }
      return res.data.data;
    },
    enabled: !!fileId,
    ...options,
  });
}

/**
 * 获取文件二进制流（公开访问）
 */
export function useServeFile(
  fileId: string,
  options?: Omit<UseQueryOptions<Blob, Error>, "queryKey" | "queryFn">,
) {
  return useQuery({
    queryKey: attachmentKeys.file(fileId),
    queryFn: async () => {
      const res = await attachmentApi.serveFile(fileId);
      const blob = res.data as Blob;
      if (!blob) {
        throw new Error("No file content");
      }
      return blob;
    },
    enabled: !!fileId,
    staleTime: Infinity,
    ...options,
  });
}
