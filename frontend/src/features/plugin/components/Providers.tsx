"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "react-hot-toast";
import { useState, useCallback } from "react";
import { useAuthStore } from "@/store";
import { useLocale } from "next-intl";
import { PluginProvider } from "./PluginContext";

function InnerProviders({ children }: { children: React.ReactNode }) {
  const { user } = useAuthStore();
  const locale = useLocale();

  const getUser = useCallback(() => {
    if (!user) return null;
    return {
      id: String(user.id),
      username: user.username ?? "",
      role: user.role ?? "user",
    };
  }, [user]);

  const getLocale = useCallback(() => locale, [locale]);

  return (
    <PluginProvider getUser={getUser} getLocale={getLocale}>
      {children}
    </PluginProvider>
  );
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
