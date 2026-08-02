export const statsKeys = {
  all: ["stats"] as const,
  day: (params?: unknown) => [...statsKeys.all, "day", params] as const,
  total: (params?: unknown) => [...statsKeys.all, "total", params] as const,
  range: (params?: unknown) => [...statsKeys.all, "range", params] as const,
  totals: (start: string, end: string) =>
    [...statsKeys.all, "total", start, end] as const,
  trend: (metric: "users" | "posts", start: string, end: string) =>
    [...statsKeys.all, "trend", metric, start, end] as const,
};
