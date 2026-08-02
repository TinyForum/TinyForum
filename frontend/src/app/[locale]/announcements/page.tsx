// src/app/[locale]/announcements/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import AnnouncementsClient from "./AnnouncementsClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("announcementsTitle"),
    description: t("announcementsDescription"),
  };
}

export default async function AnnouncementsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<AnnouncementsSkeleton />}>
      <AnnouncementsClient />
    </Suspense>
  );
}

function AnnouncementsSkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-32 w-full rounded-2xl" />
      <div className="skeleton h-24 w-full rounded-xl" />
      <div className="skeleton h-24 w-full rounded-xl" />
    </div>
  );
}
