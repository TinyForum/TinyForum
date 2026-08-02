// hooks/useBotKeys.ts
export const botKeys = {
  all: ["bots"] as const,
  lists: () => [...botKeys.all, "list"] as const,
  list: (params?: object) => [...botKeys.lists(), params] as const,
  myLists: () => [...botKeys.all, "list", "my"] as const,
  myList: (params?: object) => [...botKeys.myLists(), params] as const,
  details: () => [...botKeys.all, "detail"] as const,
  detail: (id: number) => [...botKeys.details(), id] as const,
  metadata: () => [...botKeys.all, "metadata"] as const,
};
