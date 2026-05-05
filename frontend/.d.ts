declare module "js-yaml";

declare interface Window {
  [key: `__plugin_${string}__`]: PluginEntryFn | undefined;
}
