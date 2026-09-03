// hooks/useQuestionForm.ts
import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { useTranslations } from "next-intl";
import { questionApi } from "@/shared/api/modules/questions";
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
  const t = useTranslations("Questions");
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
        toast.error(t("select_board_error"));
        return false;
      }

      if (!content.trim()) {
        toast.error(t("content_required"));
        return false;
      }

      if (data.title.length < 5) {
        toast.error(t("title_min_length"));
        return false;
      }

      if (data.title.length > 100) {
        toast.error(t("title_max_length"));
        return false;
      }

      if (data.reward_score > (user?.score || 0)) {
        toast.error(
          t("reward_exceeds_current", { score: user?.score || 0 }),
        );
        return false;
      }

      if (data.reward_score > 100) {
        toast.error(t("reward_exceeds_max"));
        return false;
      }

      return true;
    },
    [content, user?.score, t],
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

          toast.success(t("question_published"));

          if (onSuccess && questionId) {
            onSuccess(questionId);
          } else if (questionId) {
            router.push(`/questions/${questionId}`);
          } else {
            router.push("/questions");
          }
        } else {
          const errorMessage =
            response.data.message || t("publish_failed_retry");
          toast.error(errorMessage);
        }
      } catch (err: unknown) {
        console.error("发布问题错误:", err);

        const error = err as ErrorResponse;

        if (error.response) {
          const errorData = error.response.data;
          const statusCode = error.response.status;

          const errorMessages: Record<number, string> = {
            400: errorData?.message || t("bad_request_error"),
            401: t("please_login"),
            403: errorData?.message || t("insufficient_permission"),
          };

          toast.error(
            errorMessages[statusCode || 0] ||
              t("publish_failed_status", { code: statusCode }),
          );

          if (statusCode === 401) {
            router.push("/login?redirect=/questions/ask");
          }
        } else if (error.request) {
          toast.error(t("network_error"));
        } else {
          toast.error(error.message || t("publish_failed_retry"));
        }
      } finally {
        setLoading(false);
      }
    },
    [validateForm, buildRequestData, router, onSuccess, t],
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
