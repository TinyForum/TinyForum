// src/app/[locale]/explore/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import ExploreClient from "./ExploreClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("exploreTitle"),
    description: t("exploreDescription"),
  };
}

export default async function ExplorePage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<ExploreSkeleton />}>
      <ExploreClient />
    </Suspense>
  );
}

function ExploreSkeleton() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-4">
      <div className="flex flex-col items-center space-y-3">
        <div className="skeleton h-16 w-16 rounded-2xl" />
        <div className="skeleton h-8 w-40 rounded-xl" />
        <div className="skeleton h-4 w-64 rounded-xl" />
      </div>
      <div className="skeleton h-14 w-full rounded-xl" />
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-4">
          <div className="skeleton h-24 w-full rounded-xl" />
          <div className="skeleton h-24 w-full rounded-xl" />
          <div className="skeleton h-24 w-full rounded-xl" />
        </div>
        <div className="space-y-4">
          <div className="skeleton h-40 w-full rounded-xl" />
          <div className="skeleton h-40 w-full rounded-xl" />
        </div>
      </div>
    </div>
  );
}
