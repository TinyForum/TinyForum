// hooks/admin/useAdminReports.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  ReportResponse,
  ListReportsParams,
} from "@/shared/api/modules/admin/report";
import { adminReportsApi } from "@/shared/api/modules/admin/report";
import { postApi } from "@/shared/api/modules/posts";
import { commentApi } from "@/shared/api/modules/comments";
import { adminUsersApi } from "@/shared/api/modules/admin/user";
import { PageData } from "@/shared/api/types/basic.model";
import { toast } from "react-hot-toast";
import { useTranslations } from "next-intl";

const adminReportsKeys = {
  all: ["admin", "reports"] as const,
  lists: () => [...adminReportsKeys.all, "list"] as const,
  list: (params: ListReportsParams) =>
    [...adminReportsKeys.lists(), params] as const,
  detail: (id: number) => [...adminReportsKeys.all, "detail", id] as const,
  pending: () => [...adminReportsKeys.all, "pending"] as const,
};

// ========== 获取所有举报列表（支持筛选/分页）==========
export function useAdminGetReports(params: ListReportsParams = {}) {
  const t = useTranslations("Admin");
  return useQuery({
    queryKey: adminReportsKeys.list(params),
    queryFn: async () => {
      const res = await adminReportsApi.listReports(params);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("get_report_list_failed"));
      }
      return res.data.data as PageData<ReportResponse>;
    },
  });
}

// ========== 获取单个举报详情 ==========
export function useAdminGetReport(id: number) {
  const t = useTranslations("Admin");
  return useQuery({
    queryKey: adminReportsKeys.detail(id),
    queryFn: async () => {
      const res = await adminReportsApi.getReport(id);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("get_report_detail_failed"));
      }
      return res.data.data as ReportResponse;
    },
    enabled: !!id, // 仅在 id 有效时请求
  });
}

// ========== 处理举报（通过/驳回）==========
export function useAdminHandleReport() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: HandleReportRequest }) =>
      adminReportsApi.handleReport(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: adminReportsKeys.all });
      queryClient.invalidateQueries({
        queryKey: adminReportsKeys.detail(variables.id),
      });
      toast.success(t("report_handled"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("report_handle_failed"));
    },
  });
}

// ========== 删除被举报内容（按 target_type 分派到对应删除接口）==========
export function useAdminDeleteReportedContent() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: async ({
      targetType,
      targetId,
    }: {
      targetType: string;
      targetId: number;
    }) => {
      switch (targetType) {
        case "post":
          return postApi.delete(targetId);
        case "comment":
          return commentApi.delete(targetId);
        case "user":
          return adminUsersApi.deleteUser(targetId);
        default:
          throw new Error(t("delete_content_unsupported"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminReportsKeys.all });
      toast.success(t("delete_content_success"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("delete_content_failed"));
    },
  });
}

// 如果需要“批量驳回”或“批量通过”，可额外定义，这里先保持单条处理

// 导出请求体中使用的类型（方便调用方）
export interface HandleReportRequest {
  status: "resolved" | "rejected";
  handle_note: string;
  reject_reason?: string; // 驳回时可额外传原因，后端可存入 handle_note 单独字段
}
