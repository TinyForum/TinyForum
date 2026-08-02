export const adminKeys = {
  all: ["admin"] as const,
  users: (page?: number, keyword?: string) =>
    [...adminKeys.all, "users", page, keyword] as const,
  posts: (page?: number, keyword?: string) =>
    [...adminKeys.all, "posts", page, keyword] as const,
  pendingPosts: (page?: number, keyword?: string) =>
    [...adminKeys.all, "pending-posts", page, keyword] as const,
  applications: (params?: unknown) =>
    [...adminKeys.all, "applications", params] as const,
  plugins: (page?: number) => [...adminKeys.all, "plugins", page] as const,
};
