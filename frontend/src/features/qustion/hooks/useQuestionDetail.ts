// hooks/useQuestionDetail.ts (完整版)
import { CommentResponse } from "@/shared/api/types/comment.model";
import { useQuestion, useAnswers } from "./useQuestions";

export function useQuestionDetail(questionId: number) {
  // 问题详情查询
  const {
    data: question,
    isLoading: questionLoading,
    error: questionError,
    refetch: refetchQuestion,
  } = useQuestion(questionId);

  // 回答列表查询（默认不分页，获取全部）
  const {
    data: answersData,
    isLoading: answersLoading,
    error: answersError,
    refetch: refetchAnswers,
  } = useAnswers(questionId, { page: 1, page_size: 100 }); // 按需调整默认值

  // 合并 loading 和 error 状态（与原 Hook 一致）
  const loading = questionLoading || answersLoading;
  const error = questionError?.message || answersError?.message || null;

  // 刷新方法：同时刷新两个查询
  const refresh = () => {
    refetchQuestion();
    refetchAnswers();
  };

  // 回答数据解构
  const answers = answersData?.list || [];
  const answersTotal = answersData?.total || 0;

  // 从问题数据中获取点赞状态（假设 question 中有 liked 字段）
  const liked = question?.liked ?? false;

  // 返回与原 Hook 完全一致的接口
  return {
    question,
    answers,
    answersTotal,
    liked,
    loading,
    error,
    refresh,
    // 保留 set 方法以便与现有组件兼容（但建议使用 mutation 更新）
    setLiked: (newLiked: boolean) => {
      // 注意：此处不直接修改，而是通过 mutation 实现，这里仅留作占位
      void newLiked;
      console.warn("setLiked 已废弃，请使用 useLikeMutation 替代");
    },
    setAnswers: (newAnswers: CommentResponse[]) => {
      void newAnswers;
      console.warn("setAnswers 已废弃，数据由 TanStack Query 管理");
    },
    setAnswersTotal: (newTotal: number) => {
      void newTotal;
      console.warn("setAnswersTotal 已废弃，数据由 TanStack Query 管理");
    },
  };
}
