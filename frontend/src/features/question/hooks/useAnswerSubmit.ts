// hooks/useAnswerSubmit.ts
import { useMutation } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { questionApi } from "@/shared/api/modules/questions";

// 提交问题回答
export function useAnswerSubmit(questionId: number) {
  const t = useTranslations("Questions");
  return useMutation({
    mutationFn: async (content: string) => {
      const res = await questionApi.createAnswer(questionId, { content });
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("publish_failed"));
      }
      return res.data.data;
    },
  });
}
