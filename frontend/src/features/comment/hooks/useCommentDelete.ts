// hooks/useCommentDelete.ts
import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { commentApi } from "@/shared/api/modules/comments";

// 删除评论（按 postId 使评论缓存失效）
export function useCommentDelete(postId: number) {
  const queryClient = useQueryClient();
  const t = useTranslations("Comment");

  return useMutation({
    mutationFn: (commentId: number) => commentApi.delete(commentId),
    onSuccess: () => {
      toast.success(t("deleted"));
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
    },
    onError: () => toast.error(t("delete_failed")),
  });
}
