# 插件

> 注意：
>
> 插件功能目前仅作为测试使用，请勿在生产环境中使用

在 Tiny Forum 中，支持使用插件进行程序的扩展

## 一、插件的运行机制

```bash
浏览器加载页面
  └─ Providers.tsx 挂载
       └─ PluginProvider 初始化
            └─ fetchEnabledPlugins() → GET /api/v1/plugins?enabled=true
                 ├─ [成功] 返回插件列表
                 │    └─ loadPlugins(metas) 遍历每个插件
                 │         └─ loadPlugin(meta)
                 │              ├─ pluginRegistry.registerPlugin(meta)  → status: "loading"
                 │              ├─ loadPluginScript(scriptUrl, id)
                 │              │    ├─ 创建 <script src="..."> 插入 <head>
                 │              │    ├─ 等待 onload
                 │              │    └─ 读取 window.__plugin_<id>__
                 │              ├─ createPluginAPI(...)  创建沙箱
                 │              ├─ entryFn(api)  执行插件入口
                 │              │    └─ api.registerSlot("sidebar-top", MyWidget)
                 │              │         └─ pluginRegistry.registerSlotComponent(...)
                 │              └─ status → "active"
                 └─ [失败] plugins = []，页面正常渲染，插槽为空
                 
页面渲染
  └─ <PluginSlot name="sidebar-top" />
       └─ useSyncExternalStore 读取 pluginRegistry
            └─ 渲染所有注册到该插槽的组件
```

