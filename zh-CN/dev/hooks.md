TanStack Query（前身是 React Query）是一个强大的**服务器状态管理库**。它并非简单的数据获取工具，而是旨在解决服务端状态管理的复杂性问题。

它与`useState`这类管理**客户端状态**的工具不同，主要解决的是**服务端状态**的管理问题，其核心理念是**自动帮你处理获取、缓存、同步和更新服务器状态**这些繁琐的工作。

### 为什么需要 TanStack Query？

手动管理服务端状态非常困难，需要处理大量复杂问题：

- 缓存、去重请求
- 在后台更新过期数据
- 判断数据何时过期
- 分页和懒加载的性能优化
- 内存和垃圾回收管理

TanStack Query 可以将你从这些复杂逻辑中解放出来，让代码更简洁、可维护，并带给用户更快的体验。

### TanStack Query 的核心机制

它通过一套清晰的机制来管理工作流：

1.  **获取 (Fetch)**：通过`queryFn`获取数据，并自动处理重试、取消和去重。
2.  **共享 (Share)**：多个组件使用相同`queryKey`时，会共享同一份缓存数据，避免重复请求。
3.  **重新验证 (Revalidate)**：过期数据可继续显示，同时在后台静默刷新。
4.  **收集 (Collect)**：不用的数据会在内存中保留一段时间，之后被垃圾回收。

其核心是**查询键（Query Key）**，它就像是数据的“身份证”。缓存完全基于查询键管理，键必须是一个**数组**，并且要能**唯一描述**查询函数返回的数据。

### 最佳实践指南

#### 1. 查询键 (Query Keys)

- **使用工厂模式**：创建集中的查询键工厂，确保缓存管理的一致性。
  ```typescript
  // ✅ 推荐
  export const postKeys = {
    all: ["posts"] as const,
    lists: () => [...postKeys.all, "list"] as const,
    list: (params) => [...postKeys.lists(), params] as const,
    details: () => [...postKeys.all, "detail"] as const,
    detail: (id) => [...postKeys.details(), id] as const,
  };
  // ❌ 避免分散在各处写 ['posts', 'list', { page: 1 }]
  ```
- **包含所有依赖**：查询函数中用到的变量，都必须包含在查询键中。
  ```typescript
  // ✅ 正确：todoId 是依赖
  useQuery({
    queryKey: ["todo", todoId],
    queryFn: () => fetchTodoById(todoId),
  });
  ```
- **理解序列化规则**：
  - 对象键顺序不影响缓存命中，`['todos', {status, page}]` 和 `['todos', {page, status}]` 是同一个键。
  - 数组项顺序**会影响**，`['todos', status, page]` 和 `['todos', page, status]` 是不同的键。

#### 2. 查询 (Queries) 配置

- **设置 `staleTime`**：这是最重要的优化之一，用于控制数据多久后被视为“过期”。应根据数据变化频率设置：
  - **静态数据**：`staleTime: 10 * 60 * 1000` (10分钟)
  - **动态数据**：`staleTime: 30 * 1000` (30秒)
  - **实时数据**：`staleTime: 0`
- **使用 `enabled` 控制查询**：当查询依赖的参数不存在时，应禁用查询。
  ```typescript
  useQuery({
    queryKey: ["project", projectId],
    queryFn: () => fetchProject(projectId),
    enabled: !!projectId, // projectId 存在时才发起请求
  });
  ```

#### 3. 变更 (Mutations) 处理

- **变更后失效相关查询**：数据更新后，必须让相关的旧缓存失效。
  ```typescript
  const mutation = useMutation({
    mutationFn: createPost,
    onSuccess: () => {
      // 让帖子列表缓存失效，触发重新获取
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
    },
  });
  ```
- **实现乐观更新**：对用户触发的变更（如点赞、删除），应使用乐观更新提升体验。
  ```typescript
  const mutation = useMutation({
    mutationFn: deletePost,
    onMutate: async (postId) => {
      // 1. 取消正在进行的查询
      await queryClient.cancelQueries({ queryKey: postKeys.detail(postId) });
      // 2. 保存旧数据快照
      const previousData = queryClient.getQueryData(postKeys.detail(postId));
      // 3. 乐观地更新缓存
      queryClient.setQueryData(postKeys.detail(postId), null);
      // 4. 返回包含快照的上下文，用于回滚
      return { previousData };
    },
    onError: (err, postId, context) => {
      // 5. 如果出错，用快照回滚
      queryClient.setQueryData(postKeys.detail(postId), context.previousData);
    },
  });
  ```

#### 4. 错误处理

- **查询错误**：在`queryFn`中抛出错误，TanStack Query 会自动捕获并暴露在`error`状态中。
- **变更错误**：在`useMutation`的`onError`回调中处理错误。

#### 5. 开发者工具 (Devtools)

务必使用官方开发者工具，它可以帮助你可视化地查看缓存、查询键、状态等，极大提升调试效率。

### 总结

总的来说，使用 TanStack Query 的最佳实践可以概括为：

1.  **拥抱声明式**：用`useQuery`和`useMutation`声明数据依赖，而不是在`useEffect`中命令式获取。
2.  **设计好查询键**：建立清晰的工厂模式，这是高效缓存的基础。
3.  **配置好过期时间**：为不同数据设置合理的`staleTime`，减少不必要的请求。
4.  **完善变更处理**：总是使相关查询失效，并对关键操作使用乐观更新。

建议从[官方文档](https://tanstack.com/query/latest)开始学习，并结合 [TanStack Query 的 ESLint 插件](https://tanstack.com/query/latest/docs/eslint/exhaustive-deps) 来强制执行最佳实践。
