// src/app/[locale]/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import HomeClient from "./HomeClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("homeTitle"),
    description: t("homeDescription"),
  };
}

export default async function HomePage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<HomeSkeleton />}>
      <HomeClient />
    </Suspense>
  );
}

function HomeSkeleton() {
  return (
    <div className="container mx-auto max-w-7xl px-4 h-full">
      <div className="flex gap-6 h-full py-6">
        <div className="lg:w-64 xl:w-72 flex-none space-y-4">
          <div className="skeleton h-10 w-full rounded-xl" />
          <div className="skeleton h-10 w-full rounded-xl" />
          <div className="skeleton h-10 w-full rounded-xl" />
          <div className="skeleton h-10 w-full rounded-xl" />
        </div>
        <div className="flex-1 min-w-0 space-y-4">
          <div className="skeleton h-12 w-full rounded-xl" />
          <div className="skeleton h-32 w-full rounded-xl" />
          <div className="skeleton h-32 w-full rounded-xl" />
          <div className="skeleton h-32 w-full rounded-xl" />
        </div>
        <div className="lg:w-64 xl:w-72 flex-none space-y-4">
          <div className="skeleton h-32 w-full rounded-xl" />
          <div className="skeleton h-40 w-full rounded-xl" />
          <div className="skeleton h-24 w-full rounded-xl" />
        </div>
      </div>
    </div>
  );
}
