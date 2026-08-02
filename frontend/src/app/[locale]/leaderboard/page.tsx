// src/app/[locale]/leaderboard/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import LeaderboardClient from "./LeaderboardClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("leaderboardTitle"),
    description: t("leaderboardDescription"),
  };
}

export default async function LeaderboardPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<LeaderboardSkeleton />}>
      <LeaderboardClient />
    </Suspense>
  );
}

function LeaderboardSkeleton() {
  return (
    <div className="max-w-3xl mx-auto px-4 py-8 md:py-12 space-y-4">
      <div className="flex flex-col items-center space-y-3">
        <div className="skeleton h-16 w-16 rounded-full" />
        <div className="skeleton h-8 w-48 rounded-xl" />
        <div className="skeleton h-4 w-72 rounded-xl" />
      </div>
      <div className="grid grid-cols-3 gap-4">
        <div className="skeleton h-20 rounded-xl" />
        <div className="skeleton h-20 rounded-xl" />
        <div className="skeleton h-20 rounded-xl" />
      </div>
      <div className="skeleton h-16 w-full rounded-xl" />
      <div className="skeleton h-16 w-full rounded-xl" />
      <div className="skeleton h-16 w-full rounded-xl" />
      <div className="skeleton h-16 w-full rounded-xl" />
    </div>
  );
}
