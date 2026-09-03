// hooks/useTopicCoverUpload.ts
import { useCallback } from "react";
import { uploadApi } from "@/shared/api/modules/uploads";
import { toast } from "react-hot-toast";
import { useTranslations } from "next-intl";

// 话题封面上传 Hook（使用通用附件上传接口）
export function useTopicCoverUpload() {
  const t = useTranslations("Topic");
  const uploadCover = useCallback(
    async (file: File): Promise<{ url: string }> => {
      try {
        const res = await uploadApi.uploadPluginFile(file, "topic_cover");
        return { url: res.data.data?.url ?? "" };
      } catch (error) {
        toast.error(t("cover_upload_failed"));
        throw error;
      }
    },
    [t],
  );
  return { uploadCover };
}
