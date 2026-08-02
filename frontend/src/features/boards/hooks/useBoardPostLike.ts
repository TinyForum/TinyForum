import { useMutation, useQueryClient } from "@tanstack/react-query";
import { postApi } from "@/shared/api/modules/posts";
import { postKeys } from "@/features/post/hooks/usePosts";

export function useBoardPostLike() {
  const queryClient = useQueryClient();

  return {
    like: useMutation({
      mutationFn: (id: number) => postApi.like(id),
      onSettled: (_data, _error, variables) => {
        queryClient.invalidateQueries({ queryKey: postKeys.lists() });
        queryClient.invalidateQueries({ queryKey: postKeys.detail(variables) });
      },
    }),
    unlike: useMutation({
      mutationFn: (id: number) => postApi.unlike(id),
      onSettled: (_data, _error, variables) => {
        queryClient.invalidateQueries({ queryKey: postKeys.lists() });
        queryClient.invalidateQueries({ queryKey: postKeys.detail(variables) });
      },
    }),
  };
}
