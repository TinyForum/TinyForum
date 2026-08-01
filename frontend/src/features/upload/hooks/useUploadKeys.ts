// hooks/useUploadKeys.ts
export const userFileKeys = {
  all: ["userFiles"] as const,
  filesLists: () => [...userFileKeys.all, "list", "files"] as const,
  filesList: (params?: object) => [...userFileKeys.filesLists(), params] as const,
  pluginsLists: () => [...userFileKeys.all, "list", "plugins"] as const,
  pluginsList: (params?: object) => [...userFileKeys.pluginsLists(), params] as const,
  details: () => [...userFileKeys.all, "detail"] as const,
  detail: (type: string, fileId: string) =>
    [...userFileKeys.details(), type, fileId] as const,
};

export const fileServeKeys = {
  all: ["fileServe"] as const,
  file: (fileId: string) => [...fileServeKeys.all, fileId] as const,
};
