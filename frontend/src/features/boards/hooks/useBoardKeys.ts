// hooks/useBoardKeys.ts
export const boardKeys = {
  all: ["boards"] as const,
  lists: () => [...boardKeys.all, "list"] as const,
  list: (params?: object) => [...boardKeys.lists(), params] as const,
  details: () => [...boardKeys.all, "detail"] as const,
  detail: (id: number) => [...boardKeys.details(), id] as const,
  tree: () => [...boardKeys.all, "tree"] as const,
  slug: (slug: string) => [...boardKeys.all, "slug", slug] as const,
  postsBySlug: (slug: string, page: number) =>
    [...boardKeys.all, slug, "posts", page] as const,
};
