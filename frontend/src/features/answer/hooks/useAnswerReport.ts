import { useMutation } from "@tanstack/react-query";
import { toast } from "react-hot-toast";
import { useTranslations } from "next-intl";
import {
  reportApi,
  CreateReportRequest,
} from "@/shared/api/modules/reports";

// 提交回答举报（回答即评论，target_type 固定为 comment）
export function useCreateAnswerReport() {
  const t = useTranslations("Questions");
  return useMutation({
    mutationFn: (data: CreateReportRequest) => reportApi.create(data),
    onSuccess: () => {
      toast.success(t("report_submitted"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("report_failed"));
    },
  });
}