// src/app/[locale]/announcements/[id]/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import AnnouncementDetailClient from "./AnnouncementDetailClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("announcementDetailTitle"),
    description: t("announcementDetailDescription"),
  };
}

export default async function AnnouncementDetailPage({
  params,
}: {
  params: Promise<{ locale: string; id: string }>;
}) {
  const { id } = await params;

  return (
    <Suspense fallback={<AnnouncementDetailSkeleton />}>
      <AnnouncementDetailClient id={id} />
    </Suspense>
  );
}

function AnnouncementDetailSkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-5 w-28" />
      <div className="skeleton h-64 w-full rounded-2xl" />
    </div>
  );
}
