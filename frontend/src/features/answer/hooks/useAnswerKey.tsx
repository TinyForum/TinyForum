// hooks/useAnswerKeys.ts
export const answerKeys = {
  all: ["answers"] as const,
  lists: () => [...answerKeys.all, "list"] as const,
  list: (filters?: object) => [...answerKeys.lists(), filters] as const,
  details: () => [...answerKeys.all, "detail"] as const,
  detail: (id: number) => [...answerKeys.details(), id] as const,
  voteStatus: (id: number) => [...answerKeys.detail(id), "voteStatus"] as const,
};
