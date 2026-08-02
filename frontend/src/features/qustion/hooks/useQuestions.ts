// hooks/useQuestions.ts
import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
  UseMutationOptions,
} from "@tanstack/react-query";

import { AxiosError } from "axios";
import { questionKeys } from "./useQuestionKeys";
import {
  questionApi,
  QuestionListParams,
} from "@/shared/api/modules/questions";
import { PageData } from "@/shared/api/types/basic.model";
import {
  CommentResponse,
  CreateCommentPayload,
} from "@/shared/api/types/comment.model";
import {
  QuestionResponse,
  QuestionSimple,
} from "@/shared/api/types/question.model";
import { AnswerResponse } from "@/shared/api/types/answer.model";

// ---------- 查询：问题详情 ----------
export function useQuestion(
  questionId: number,
  options?: Omit<
    UseQueryOptions<QuestionResponse, AxiosError>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery<QuestionResponse, AxiosError>({
    queryKey: questionKeys.detail(questionId),
    queryFn: async () => {
      const res = await questionApi.getDetail(questionId);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取问题详情失败");
      }
      if (!res.data.data) {
        throw new Error("问题详情数据为空");
      }
      return res.data.data;
    },
    enabled: !!questionId,
    ...options,
  });
}

// ---------- 查询：回答列表（分页） ----------
export function useAnswers(
  questionId: number,
  params?: { page?: number; page_size?: number },
  options?: Omit<
    UseQueryOptions<PageData<AnswerResponse>, AxiosError>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery<PageData<AnswerResponse>, AxiosError>({
    queryKey: questionKeys.answers(questionId, params),
    queryFn: async () => {
      const res = await questionApi.getAnswers(questionId, params);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取回答列表失败");
      }
      if (!res.data.data) {
        throw new Error("回答列表数据为空");
      }
      return res.data.data; // 假设返回 PageData<CommentResponse>
    },
    enabled: !!questionId,
    ...options,
  });
}

// ---------- 变更：创建回答 ----------
export function useCreateAnswer(
  questionId: number,
  options?: UseMutationOptions<
    CommentResponse,
    AxiosError,
    Omit<CreateCommentPayload, "post_id">
  >,
) {
  const queryClient = useQueryClient();

  return useMutation<
    CommentResponse,
    AxiosError,
    Omit<CreateCommentPayload, "post_id">
  >({
    mutationFn: async (payload) => {
      // 内部组装 post_id
      const fullPayload: CreateCommentPayload = {
        ...payload,
        post_id: questionId,
      };
      const res = await questionApi.createAnswer(questionId, fullPayload);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "创建回答失败");
      }
      if (!res.data.data) {
        throw new Error("创建回答返回数据为空");
      }
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: questionKeys.answers(questionId),
      });
    },
    ...options,
  });
}

// ---------- 变更：删除回答（如果存在） ----------
// 注意：如果 questionApi 没有 deleteAnswer，此方法需移除或改用 commentApi.delete
// 这里假设 questionApi.deleteAnswer 存在，若没有则注释掉
/*
export function useDeleteAnswer(
  questionId: number,
  options?: UseMutationOptions<void, AxiosError, number>
) {
  const queryClient = useQueryClient();

  return useMutation<void, AxiosError, number>({
    mutationFn: async (answerId) => {
      const res = await questionApi.deleteAnswer(answerId);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '删除回答失败');
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: questionKeys.answers(questionId),
      });
    },
    ...options,
  });
}
*/

// ---------- 变更：点赞问题（如果存在） ----------
// 如果 questionApi 没有 toggleLike，可改为使用 commentApi 的投票接口（如果有）
// 这里暂不实现，保留注释
/*
export function useToggleQuestionLike(
  options?: UseMutationOptions<{ liked: boolean; likeCount: number }, AxiosError, number>
) {
  const queryClient = useQueryClient();

  return useMutation<{ liked: boolean; likeCount: number }, AxiosError, number>({
    mutationFn: async (questionId) => {
      const res = await questionApi.toggleLike(questionId);
      if (res.data.code !== 0) throw new Error(res.data.message || '操作失败');
      if (!res.data.data) throw new Error('返回数据为空');
      return res.data.data;
    },
    onSuccess: (data, questionId) => {
      queryClient.setQueryData<QuestionResponse>(
        questionKeys.detail(questionId),
        (old) => {
          if (!old) return old;
          return {
            ...old,
            liked: data.liked,
            likeCount: data.likeCount,
          };
        }
      );
    },
    ...options,
  });
}
*/

// ---------- 兼容原 useQuestionDetail 的组合 Hook ----------
export function useQuestionDetail(questionId: number) {
  const {
    data: question,
    isLoading: questionLoading,
    error: questionError,
    refetch: refetchQuestion,
  } = useQuestion(questionId);

  const {
    data: answersData,
    isLoading: answersLoading,
    error: answersError,
    refetch: refetchAnswers,
  } = useAnswers(questionId, { page: 1, page_size: 100 });

  const loading = questionLoading || answersLoading;
  const error = questionError?.message || answersError?.message || null;

  const refresh = () => {
    refetchQuestion();
    refetchAnswers();
  };

  const answers = answersData?.list || [];
  const answersTotal = answersData?.total || 0;
  const liked = question?.liked ?? false;

  return {
    question,
    answers,
    answersTotal,
    liked,
    loading,
    error,
    refresh,
    // 保留 setter 占位（实际应通过 mutation 修改）
    setLiked: () => {
      console.warn(
        "setLiked 已废弃，请使用 useToggleQuestionLike 或相应 mutation",
      );
    },
    setAnswers: () => {
      console.warn("setAnswers 已废弃，数据由 TanStack Query 管理");
    },
    setAnswersTotal: () => {
      console.warn("setAnswersTotal 已废弃，数据由 TanStack Query 管理");
    },
  };
}

// ---------- 查询：问题列表（精简版，用于首页） ----------
export function useQuestionList(
  params?: QuestionListParams,
  options?: Omit<
    UseQueryOptions<PageData<QuestionSimple>, AxiosError>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery<PageData<QuestionSimple>, AxiosError>({
    queryKey: [...questionKeys.list(), params],
    queryFn: async () => {
      const res = await questionApi.getSimple(params);
      return unwrap(res);
    },
    enabled: !!params?.page, // 或者根据是否有参数决定，但默认启用
    ...options,
  });
}

// 辅助解包函数
function unwrap<T>(res: {
  data: { code: number; message?: string; data?: T };
}): T {
  if (res.data.code !== 0) {
    throw new Error(res.data.message || "请求失败");
  }
  if (!res.data.data) {
    throw new Error("返回数据为空");
  }
  return res.data.data;
}
