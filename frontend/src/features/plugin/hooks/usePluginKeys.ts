export const pluginKeys = {
  all: ["admin", "plugins"] as const,
  list: (page: number) => [...pluginKeys.all, page] as const,
  enabled: () => [...pluginKeys.all, "enabled"] as const,
};
