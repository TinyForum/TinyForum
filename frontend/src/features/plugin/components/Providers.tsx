"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "react-hot-toast";
import { useState, useCallback } from "react";
import { useAuthStore } from "@/store";
import { useLocale } from "next-intl";
import { PluginProvider } from "./PluginContext";

/**
 * InnerProviders 组件，用于提供用户和本地化信息给子组件
 * @param children - React 节点，将被包裹在 Provider 中
 */
export function InnerProviders({ children }: { children: React.ReactNode }) {
  // 使用 useAuthStore 获取用户信息
  const { user } = useAuthStore();
  // 使用 useLocale 获取本地化信息
  // const locale = useLocale();

  // 使用 useCallback 缓存 getUser 函数，避免不必要的重新创建
  const getUser = useCallback(() => {
    // 如果用户不存在，返回 null
    if (!user) return null;
    // 返回用户信息对象，包含 id、username 和 role
    return {
      id: String(user.id),
      username: user.username ?? "", // 如果 username 不存在，使用空字符串作为默认值
      role: user.role ?? "user", // 如果 role 不存在，使用 "user" 作为默认值
    };
  }, [user]); // 依赖项为 user，当 user 变化时重新创建函数

  // 使用 useCallback 缓存 getLocale 函数，避免不必要的重新创建
  // const getLocale = useCallback(() => locale, [locale]); // 依赖项为 locale，当 locale 变化时重新创建函数

  return <PluginProvider getUser={getUser}>{children}</PluginProvider>;
}

export default function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 30 * 1000,
            retry: 1,
          },
        },
      }),
  );

  return (
    <QueryClientProvider client={queryClient}>
      <InnerProviders>{children}</InnerProviders>
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 3000,
          style: {
            background: "var(--fallback-b1,oklch(var(--b1)/1))",
            color: "var(--fallback-bc,oklch(var(--bc)/1))",
            border: "1px solid var(--fallback-b3,oklch(var(--b3)/1))",
          },
        }}
        containerStyle={{ top: 80 }}
      />
    </QueryClientProvider>
  );
}
