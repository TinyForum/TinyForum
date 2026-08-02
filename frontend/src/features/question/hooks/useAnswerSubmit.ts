// hooks/useAnswerSubmit.ts
import { useMutation } from "@tanstack/react-query";
import { questionApi } from "@/shared/api/modules/questions";

// 提交问题回答
export function useAnswerSubmit(questionId: number) {
  return useMutation({
    mutationFn: async (content: string) => {
      const res = await questionApi.createAnswer(questionId, { content });
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "发布失败");
      }
      return res.data.data;
    },
  });
}
