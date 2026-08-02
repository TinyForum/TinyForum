// src/app/[locale]/dashboard/system/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import SystemDashboardClient from "./SystemDashboardClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("systemTitle"),
    description: t("systemDescription"),
  };
}

export default async function SystemDashboardPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<SystemDashboardSkeleton />}>
      <SystemDashboardClient />
    </Suspense>
  );
}

function SystemDashboardSkeleton() {
  return (
    <div className="flex h-[calc(100vh-64px)] bg-base-200 overflow-hidden rounded-xl border border-base-300 shadow-sm animate-pulse">
      <div className="w-56 bg-base-300 border-r border-base-300" />
      <div className="flex-1 p-6 space-y-4">
        <div className="skeleton h-6 w-40" />
        <div className="skeleton h-64 w-full rounded-2xl" />
      </div>
    </div>
  );
}
