"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "react-hot-toast";
import { useState } from "react";

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
      {children}
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
        containerStyle={{
          top: 80, // 向下偏移 80px，根据你的 header 高度调整
        }}
      />
    </QueryClientProvider>
  );
}
