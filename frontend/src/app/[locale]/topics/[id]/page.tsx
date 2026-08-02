// app/[locale]/topics/[id]/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import TopicDetailClient from "./TopicDetailClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("topicDetailTitle"),
    description: t("topicDetailDescription"),
  };
}

export default async function TopicDetailPage({
  params,
}: {
  params: Promise<{ locale: string; id: string }>;
}) {
  const { id } = await params;

  return (
    <Suspense fallback={<TopicDetailSkeleton />}>
      <TopicDetailClient id={id} />
    </Suspense>
  );
}

function TopicDetailSkeleton() {
  return (
    <div className="min-h-screen bg-base-200 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <div className="space-y-4">
          <div className="skeleton h-8 w-32 rounded-lg" />
          <div className="skeleton h-64 w-full rounded-2xl" />
          <div className="skeleton h-12 w-full rounded-2xl" />
          <div className="skeleton h-40 w-full rounded-2xl" />
        </div>
      </div>
    </div>
  );
}
