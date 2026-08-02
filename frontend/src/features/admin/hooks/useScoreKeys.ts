export const scoreKeys = {
  all: ["score"] as const,
  list: (userId?: string) => [...scoreKeys.all, "list", userId] as const,
  me: () => [...scoreKeys.all, "me"] as const,
};
