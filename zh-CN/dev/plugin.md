# 插件系统运行机制与开发指南

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

该插件系统采用**动态加载 + 注册表模式**，运行时流程如下：

1. **初始化加载**
   - 应用启动时，`PluginProvider`（位于 `PluginContext.tsx`）调用 `fetchEnabledPlugins()` 从后端 `/api/v1/plugins?enabled=true` 获取所有启用的插件元数据（`PluginMeta[]`）。
   - 随后调用 `loadPlugins(metas, options)` 批量加载。

2. **脚本加载**
   - 对每个插件，`loadPlugin()` 会动态创建 `<script>` 标签，`src` 指向元数据中的 `scriptUrl`（通常为 CDN 或内部托管地址）。
   - 设置超时（10秒）和错误处理，加载完成后检查全局对象 `window.__plugin_<slug>__` 是否为函数，该函数即为**插件入口函数**。

3. **插件初始化**
   - 通过 `createPluginAPI()` 构建每个插件的独立 API 沙箱，包含：
     - `registerSlot(slotName, component, options)`：注册 React 组件到指定插槽。
     - `on(event, handler)` / `off()`：订阅/取消系统事件。
     - `getUser()`：获取当前用户信息。
     - `log()`：带前缀的日志输出。
     - `getConfig()`：获取插件的配置数据（存储在注册表中）。
   - 执行入口函数，传入 API 对象，插件在此完成注册操作。

4. **渲染插槽**
   - UI 中使用 `<PluginSlot name="...">` 组件定义插槽位置。
   - `PluginSlot` 内部通过 `useSyncExternalStore` 订阅 `pluginRegistry` 的变化，从注册表中获取该插槽对应的组件列表，按 `order` 排序后渲染。
   - 每个组件被 `<PluginComponentWrapper>` 错误边界包裹，防止单个插件崩溃影响整体 UI。

5. **状态管理与更新**
   - `pluginRegistry` 是单例，管理插件元数据、插槽组件、事件监听器。
   - 任何注册/卸载变更都会触发 `notify()`，通知所有订阅者（如 `PluginSlot`）重新渲染。
   - 管理端操作（启用/禁用/删除）通过 `useAdminPlugins` 调用后端 API，成功后由 `PluginProvider` 重新加载插件列表，实现动态更新。

---

### 二、如何编写一个插件

#### 1. 准备插件代码

插件本质上是一个**独立的 JavaScript 模块**（推荐 UMD 或 IIFE 格式），需在全局暴露一个入口函数。例如：

```javascript
// my-plugin.js
(function () {
  // 入口函数，必须命名为 window.__plugin_my_plugin__
  window.__plugin_my_plugin__ = function (api) {
    // 1. 注册插槽组件
    api.registerSlot(
      "sidebar",
      function ({ user }) {
        return React.createElement(
          "div",
          null,
          `Hello, ${user?.username || "Guest"}!`,
        );
      },
      { order: 10 },
    );

    // 2. 监听事件
    api.on("post:created", (postData) => {
      api.log("info", "New post created: " + postData.title);
    });

    // 3. 使用 API 能力
    const user = api.getUser();
    console.log("Plugin loaded for user:", user);
  };
})();
```

#### 2. 打包与部署

- 将代码打包为单个 JS 文件（可使用 Webpack/Rollup），确保入口函数暴露到 `window`。
- 将 JS 文件上传到可公开访问的 URL（或通过后端上传接口存储），记录该 URL。

#### 3. 创建插件元数据

在管理端通过“安装插件”表单或 API 提交插件信息，必须包含：

- `name`、`version`、`author`、`scriptUrl`（即 JS 文件 URL）
- 可选 `description`、`iconUrl`、`slots`（声明使用的插槽名称，便于管理界面展示）

#### 4. 插件注册与启用

- 管理端调用 `pluginApi.create()` 将元数据保存到数据库，默认 `enabled` 为 true。
- 若为文件上传方式（ZIP），后端需解包并提取元数据（`plugin.json`），然后自动创建记录。

#### 5. 调试与日志

- 使用 `api.log('info', '...')` 输出日志，可在“插件日志”Tab 中查看（需后端收集）。
- 开发时可直接在浏览器控制台观察加载状态和错误信息。

<!--
### 三、如何完善基础设施

当前实现已具备基础骨架，但为进一步生产就绪，建议从以下方面增强：

#### 1. 安全与隔离

- **沙箱隔离**：当前插件直接运行在页面上下文，存在安全风险。可考虑使用 **Web Workers** 或 **Iframe + postMessage** 隔离 DOM 操作。
- **权限控制**：为插件声明所需权限（如 `"read:user"`、`"write:post"`），加载时校验，防止越权。
- **内容安全策略（CSP）**：严格限制 `script-src`，只允许可信 CDN 或内部域名。

#### 2. 开发体验

- **热更新模式**：增加开发环境监听，插件代码变更后自动重新加载，无需刷新页面。
- **TypeScript 类型支持**：发布 `@types/plugin-api` 包，让插件开发者获得类型提示。
- **插件脚手架**：提供 CLI 工具快速生成插件模板，包含构建配置和示例代码。

#### 3. 性能优化

- **代码分割与懒加载**：`PluginSlot` 仅在插槽可见时才加载对应组件（使用 `React.lazy`）。
- **缓存策略**：对已加载的插件脚本缓存至 localStorage 或 IndexedDB，减少重复请求。
- **加载超时重试**：增加重试机制，并在失败时提供用户友好的降级提示。

#### 4. 插件管理与生态

- **依赖管理**：允许插件声明依赖的其他插件或库版本，确保兼容性。
- **插件市场集成**：完善 `PluginMarketTab`，支持浏览远程插件列表、一键安装（需后端市场服务）。
- **配置界面**：`PluginConfigTab` 应动态加载每个插件的配置表单（插件提供配置 Schema），实现可视化配置。
- **版本升级**：支持插件版本更新，提供迁移脚本。

#### 5. 监控与日志

- **日志收集**：实现 `api.log()` 将日志发送至后端，支持按插件、级别、时间检索。
- **性能监控**：记录插件加载耗时、组件渲染性能，上报至监控平台。
- **错误追踪**：集成 Sentry 等工具，捕获插件内部错误并关联插件 ID。

#### 6. 国际化（i18n）

- 在 `PluginAPI` 中注入 `getLocale()` 方法，并传递当前语言，插件可据此渲染多语言内容。
- 管理界面插件描述支持多语言字段。

#### 7. 测试与文档

- 编写单元测试覆盖 `PluginRegistry`、`PluginLoader` 等核心模块。
- 提供完整的插件开发文档，包含 API 参考、示例项目、常见问题。

#### 8. 扩展 API 能力

- 暴露更多系统能力，如路由跳转、弹窗、表单验证、数据请求等。
- 提供 `registerRoute()` 方法，允许插件在特定路径下渲染自定义页面。 -->
