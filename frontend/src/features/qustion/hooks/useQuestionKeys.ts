// hooks/useQuestionKeys.ts
export const questionKeys = {
  all: ["questions"] as const,
  lists: () => [...questionKeys.all, "list"] as const,
  list: (filters?: object) => [...questionKeys.lists(), filters] as const,
  details: () => [...questionKeys.all, "detail"] as const,
  detail: (id: number) => [...questionKeys.details(), id] as const,
  answers: (
    questionId: number,
    params?: { page?: number; page_size?: number },
  ) => [...questionKeys.detail(questionId), "answers", params] as const,
};
