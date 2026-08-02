// src/app/[locale]/boards/applications/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import BoardApplicationsClient from "./BoardApplicationsClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("applicationsTitle"),
    description: t("applicationsDescription"),
  };
}

export default async function BoardApplicationsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<BoardApplicationsSkeleton />}>
      <BoardApplicationsClient />
    </Suspense>
  );
}

function BoardApplicationsSkeleton() {
  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10 space-y-4">
      <div className="skeleton h-10 w-40 rounded-xl" />
      <div className="space-y-4">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="skeleton h-40 w-full rounded-xl" />
        ))}
      </div>
    </div>
  );
}
