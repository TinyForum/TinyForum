// hooks/useCoverUpload.ts
import { useCallback } from "react";
import { uploadApi } from "@/shared/api/modules/uploads";
import { UploadResult } from "@/shared/ui/upload/upload.types";

export function useCoverUpload() {
  const uploadCover = useCallback(
    async (file: File): Promise<UploadResult> => {
      const res = await uploadApi.uploadPluginFile(file, "post_cover");
      return { url: res.data.data?.url ?? "", file_id: res.data.data?.file_id };
    },
    [],
  );
  return { uploadCover };
}
