// app/[locale]/timeline/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import TimelineClient from "./TimelineClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("timelineTitle"),
    description: t("timelineDescription"),
  };
}

export default async function TimelinePage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<TimelineSkeleton />}>
      <TimelineClient />
    </Suspense>
  );
}

function TimelineSkeleton() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-3xl mx-auto px-4 py-8 md:py-12">
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="card bg-base-100 shadow-sm p-6 animate-pulse">
              <div className="flex gap-4">
                <div className="w-12 h-12 bg-base-200 rounded-full" />
                <div className="flex-1 space-y-3">
                  <div className="h-4 bg-base-200 rounded w-1/4" />
                  <div className="h-3 bg-base-200 rounded w-3/4" />
                  <div className="h-20 bg-base-200 rounded-lg" />
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
