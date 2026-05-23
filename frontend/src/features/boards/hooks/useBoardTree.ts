// hooks/useBoardTree.ts
import { useQuery } from "@tanstack/react-query";
import { boardApi } from "@/shared/api/modules/boards";
import { Board } from "@/shared/api/types/board.model";

export const useBoardTree = () => {
  return useQuery({
    queryKey: ["boards-tree"],
    queryFn: async () => {
      const res = await boardApi.getTree();
      return (res.data.data as Board[]) || [];
    },
  });
};
