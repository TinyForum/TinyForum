export const statsKeys = {
  all: ["stats"] as const,
  totals: (start: string, end: string) =>
    [...statsKeys.all, "total", start, end] as const,
  trend: (metric: "users" | "posts", start: string, end: string) =>
    [...statsKeys.all, "trend", metric, start, end] as const,
};
