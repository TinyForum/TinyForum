// hooks/useQuestionForm.ts
import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { questionApi } from "@/shared/api";
import { useAuthStore } from "@/store/auth";
import { toast } from "react-hot-toast";
import {
  CreateQuestionPayload,
  CreateQuestionResponse,
} from "@/shared/api/modules/questions";
import { ApiResponse } from "@/shared/api/types/basic.model";

export interface AskFormData {
  title: string;
  content: string;
  summary: string;
  reward_score: number;
}

interface UseQuestionFormProps {
  onSuccess?: (questionId: number) => void;
}

interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
    status?: number;
  };
  request?: unknown;
  message?: string;
}

export function useQuestionForm({ onSuccess }: UseQuestionFormProps = {}) {
  const router = useRouter();
  const { user } = useAuthStore();
  const [loading, setLoading] = useState<boolean>(false);
  const [content, setContent] = useState<string>("");

  const form = useForm<AskFormData>({
    defaultValues: {
      title: "",
      summary: "",
      reward_score: 0,
    },
  });

  const validateForm = useCallback(
    (data: AskFormData, selectedTags: number[], boardId: number): boolean => {
      if (!boardId || boardId === 0) {
        toast.error("请选择板块");
        return false;
      }

      if (!content.trim()) {
        toast.error("请输入问题内容");
        return false;
      }

      if (data.title.length < 5) {
        toast.error("标题至少需要5个字符");
        return false;
      }

      if (data.title.length > 100) {
        toast.error("标题不能超过100个字符");
        return false;
      }

      if (data.reward_score > (user?.score || 0)) {
        toast.error(`悬赏积分不能超过当前积分（${user?.score || 0}）`);
        return false;
      }

      if (data.reward_score > 100) {
        toast.error("悬赏积分不能超过100");
        return false;
      }

      return true;
    },
    [content, user?.score],
  );

  const filterValidTags = useCallback((tagIds: number[]): number[] => {
    return tagIds.filter((id: number): boolean => id > 0); // 修复：直接返回 id > 0 的布尔值
  }, []);

  const buildRequestData = useCallback(
    (
      data: AskFormData,
      selectedTags: number[],
      boardId: number,
    ): CreateQuestionPayload => {
      const validTagIds = filterValidTags(selectedTags);

      return {
        title: data.title.trim(),
        content: content.trim(),
        summary: data.summary?.trim() || "",
        cover: "",
        board_id: boardId,
        tag_ids: validTagIds,
        reward_score: Number(data.reward_score),
      };
    },
    [content, filterValidTags],
  );

  const handleSubmit = useCallback(
    async (
      data: AskFormData,
      selectedTags: number[],
      boardId: number,
    ): Promise<void> => {
      if (!validateForm(data, selectedTags, boardId)) {
        return;
      }

      const requestData = buildRequestData(data, selectedTags, boardId);
      console.log("发送请求:", requestData);

      setLoading(true);

      try {
        const response: { data: ApiResponse<CreateQuestionResponse> } =
          await questionApi.create(requestData);
        console.log("响应:", response.data);

        // 修复：统一使用 code === 0
        const isSuccess = response.data.code === 0;

        if (isSuccess && response.data.data) {
          const questionId = response.data.data.id;

          toast.success("问题发布成功！");

          if (onSuccess && questionId) {
            onSuccess(questionId);
          } else if (questionId) {
            router.push(`/questions/${questionId}`);
          } else {
            router.push("/questions");
          }
        } else {
          const errorMessage = response.data.message || "发布失败，请稍后重试";
          toast.error(errorMessage);
        }
      } catch (err: unknown) {
        console.error("发布问题错误:", err);

        const error = err as ErrorResponse;

        if (error.response) {
          const errorData = error.response.data;
          const statusCode = error.response.status;

          const errorMessages: Record<number, string> = {
            400: errorData?.message || "请求参数错误，请检查输入",
            401: "请先登录",
            403: errorData?.message || "积分不足或权限不足",
          };

          toast.error(
            errorMessages[statusCode || 0] || `发布失败 (${statusCode})`,
          );

          if (statusCode === 401) {
            router.push("/login?redirect=/questions/ask");
          }
        } else if (error.request) {
          toast.error("网络错误，请检查网络连接");
        } else {
          toast.error(error.message || "发布失败，请稍后重试");
        }
      } finally {
        setLoading(false);
      }
    },
    [validateForm, buildRequestData, router, onSuccess],
  );

  return {
    form,
    content,
    setContent,
    loading,
    handleSubmit: (selectedTags: number[], boardId: number): Promise<void> =>
      handleSubmit(form.getValues(), selectedTags, boardId),
  };
}
