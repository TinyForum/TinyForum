import { attachmentApi } from "@/shared/api/modules/attachments";
import {
  UseMutationOptions,
  useQueryClient,
  useMutation,
} from "@tanstack/react-query";
import { attachmentKeys } from "./retrieve";

/**
 * 删除文件 Mutation
 */
export function useDeleteFile(
  options?: UseMutationOptions<{ message: string }, Error, { fileId: string }>,
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ fileId }) => {
      const res = await attachmentApi.deleteFile(fileId);
      if (!res.data.data) {
        throw new Error("Delete failed");
      }
      return res.data.data;
    },
    onSuccess: (data, variables, onMutateResult, context) => {
      queryClient.invalidateQueries({ queryKey: attachmentKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: attachmentKeys.detail(variables.fileId),
      });
      queryClient.removeQueries({
        queryKey: attachmentKeys.file(variables.fileId),
      });
      options?.onSuccess?.(data, variables, onMutateResult, context);
    },
    ...options,
  });
}
