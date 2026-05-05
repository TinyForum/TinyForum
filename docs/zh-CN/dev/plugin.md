按步骤说明，照着做即可。

---

## 第一步：复制文件

将输出的文件按路径放入你的项目：

```
# 新增文件（直接放入）
src/shared/plugin/           ← 6个文件全部放入
src/shared/api/modules/plugins.ts
src/features/admin/hooks/useAdminPlugins.ts
src/features/admin/components/PluginManager.tsx
src/features/admin/components/PluginUploadForm.tsx

# 替换已有文件
src/layout/layout/Providers.tsx
src/features/admin/components/AdminTabs.tsx
```

---

## 第二步：更新 admin.types.ts

```ts
// src/shared/type/admin.types.ts
// 原来是：
export type TabType = "users" | "posts";

// 改为：
export type TabType = "users" | "posts" | "plugins";
```

---

## 第三步：在 Admin 页面挂载 PluginManager

找到你的 `src/app/[locale]/dashboard/admin/page.tsx`，加入插件 Tab 的渲染逻辑：

```tsx
import { PluginManager } from "@/features/admin/components/PluginManager";

// 在渲染 Tab 内容的地方加一个分支：
{activeTab === "users" && <UserManagement ... />}
{activeTab === "posts" && <PostManagement ... />}
{activeTab === "plugins" && <PluginManager t={t} />}  // ← 新增这行
```

---

## 第四步：注册 API 出口

```ts
// src/shared/api/index.ts  末尾加一行：
export { pluginApi } from "./modules/plugins";
```

---

## 第五步：在页面预埋插槽（可选但推荐）

在你想让插件注入 UI 的地方加 `<PluginSlot>`：

```tsx
import { PluginSlot } from "@/shared/plugin/PluginSlot";

// src/layout/home/LeftSidebar.tsx
<PluginSlot name="sidebar-top" />
<PluginSlot name="sidebar-bottom" />

// src/layout/home/mid/PostList.tsx
<PluginSlot name="post-list-top" />
```

---

## 第六步：后端实现 `/api/v1/plugins` 接口

前端调用以下接口，后端需实现：

| 方法   | 路径                         | 说明                            |
| ------ | ---------------------------- | ------------------------------- |
| GET    | `/api/v1/plugins`            | 列表，支持 `?enabled=true` 筛选 |
| POST   | `/api/v1/plugins`            | 安装插件                        |
| PUT    | `/api/v1/plugins/:id`        | 更新配置                        |
| PATCH  | `/api/v1/plugins/:id/toggle` | 启用/禁用                       |
| DELETE | `/api/v1/plugins/:id`        | 删除                            |

返回格式与项目其他接口保持一致：`{ data: ..., total: ... }`

---

## 第七步：编写插件（给插件开发者）

插件是一个 JS 文件，暴露到 `window.__plugin_<id>__`：

```js
// my-plugin.js（打包后上传到 CDN）
window.__plugin_my_plugin__ = function (api) {
  // 注册一个 React 组件到侧边栏顶部插槽
  api.registerSlot("sidebar-top", function MyWidget() {
    return React.createElement(
      "div",
      {
        className: "card bg-base-100 border border-base-300 p-3 text-sm",
      },
      "👋 Hello from My Plugin!",
    );
  });

  // 监听事件
  api.on("post:view", function (data) {
    api.log("info", "User viewed a post: " + JSON.stringify(data));
  });
};
```

---

## 验证一切正常

1. 启动项目，打开 Admin → 应看到 **Plugins** Tab
2. 点击 **Install Plugin**，填入一个测试 JS URL
3. 刷新页面，该插件的组件应出现在对应插槽位置
