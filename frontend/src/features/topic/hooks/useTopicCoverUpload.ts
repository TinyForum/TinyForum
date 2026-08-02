// hooks/useTopicCoverUpload.ts
import { useCallback } from "react";
import { uploadApi } from "@/shared/api/modules/uploads";
import { toast } from "react-hot-toast";

// 话题封面上传 Hook（使用通用附件上传接口）
export function useTopicCoverUpload() {
  const uploadCover = useCallback(
    async (file: File): Promise<{ url: string }> => {
      try {
        const res = await uploadApi.uploadPluginFile(file, "topic_cover");
        return { url: res.data.data?.url ?? "" };
      } catch (error) {
        toast.error("封面上传失败");
        throw error;
      }
    },
    [],
  );
  return { uploadCover };
}
