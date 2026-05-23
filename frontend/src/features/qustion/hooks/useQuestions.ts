// features/question/hooks/useQuestions.ts
import { useQuery, UseQueryOptions } from "@tanstack/react-query";
import { questionApi } from "@/shared/api/modules/questions";
import { PageData } from "@/shared/api/types/basic.model";
import { QuestionSimple } from "@/shared/api/types/question.model";

export const useQuestions = (
  params?: { page?: number; page_size?: number; board_id?: number },
  options?: Omit<
    UseQueryOptions<PageData<QuestionSimple>>,
    "queryKey" | "queryFn"
  >,
) => {
  return useQuery({
    queryKey: ["questions", params],
    queryFn: async () => {
      const res = await questionApi.getSimple(params);
      if (res.data.code !== 0)
        throw new Error(res.data.message || "加载问答列表失败");
      const data = res.data.data as PageData<QuestionSimple>;
      if (!data) throw new Error("数据为空");
      return data;
    },
    ...options,
  });
};
