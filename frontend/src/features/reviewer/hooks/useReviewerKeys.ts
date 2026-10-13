// hooks/useReviewerKeys.ts
export const reviewerKeys = {
  all: ["reviewer"] as const,
  tasks: () => [...reviewerKeys.all, "tasks"] as const,
  taskList: (params?: object) => [...reviewerKeys.tasks(), params] as const,
  history: () => [...reviewerKeys.all, "history"] as const,
  historyList: (params?: object) => [...reviewerKeys.history(), params] as const,
};
