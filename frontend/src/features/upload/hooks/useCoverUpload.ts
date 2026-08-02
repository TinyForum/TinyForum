// hooks/useCoverUpload.ts
import { useCallback } from "react";
import { uploadApi } from "@/shared/api/modules/uploads";

// 封面上传 Hook：使用通用附件上传接口（post_cover 类型）
export function useCoverUpload() {
  const uploadCover = useCallback(
    async (file: File): Promise<{ url: string }> => {
      const res = await uploadApi.uploadPluginFile(file, "post_cover");
      return { url: res.data.data?.url ?? "" };
    },
    [],
  );
  return { uploadCover };
}
