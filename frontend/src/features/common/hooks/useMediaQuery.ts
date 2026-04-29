// hooks/useMediaQuery.ts
import { useState, useEffect } from "react";

/**
 * 自定义 Hook：用于监听媒体查询并返回匹配结果
 * @param query - 媒体查询字符串，如 "(min-width: 768px)"
 * @returns 返回当前媒体查询是否匹配的布尔值
 */
export function useMediaQuery(query: string): boolean {
  // 使用 useState 初始化匹配状态，SSR（服务器端渲染）环境下默认为 false
  const [matches, setMatches] = useState(false); // SSR 默认 false

  // 使用 useEffect 监听媒体查询变化
  useEffect(() => {
    // 创建媒体查询对象
    const mql = window.matchMedia(query);
    // 设置当前匹配状态
    setMatches(mql.matches);

    // 定义处理媒体查询变化的回调函数
    const handler = (e: MediaQueryListEvent) => setMatches(e.matches);
    // 添加变化事件监听器
    mql.addEventListener("change", handler);
    // 清理函数：组件卸载时移除事件监听器
    return () => mql.removeEventListener("change", handler);
  }, [query]); // 依赖项为 query，当 query 变化时重新执行

  return matches;
}
