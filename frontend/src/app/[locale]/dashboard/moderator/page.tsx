// src/app/[locale]/dashboard/moderator/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import ModeratorDashboardClient from "./ModeratorDashboardClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("moderatorTitle"),
    description: t("moderatorDescription"),
  };
}

export default async function ModeratorDashboardPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<ModeratorDashboardSkeleton />}>
      <ModeratorDashboardClient />
    </Suspense>
  );
}

function ModeratorDashboardSkeleton() {
  return (
    <div className="flex h-screen bg-base-100 animate-pulse">
      <div className="w-60 bg-base-200 border-r border-base-300" />
      <div className="flex-1 p-6 space-y-4">
        <div className="skeleton h-8 w-40" />
        <div className="skeleton h-64 w-full rounded-2xl" />
      </div>
    </div>
  );
}
