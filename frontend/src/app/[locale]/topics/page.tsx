// app/[locale]/topics/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import TopicsClient from "./TopicsClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("topicsTitle"),
    description: t("topicsDescription"),
  };
}

export default async function TopicsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<TopicsSkeleton />}>
      <TopicsClient />
    </Suspense>
  );
}

function TopicsSkeleton() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-4xl mx-auto px-4 py-8 md:py-12">
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="card bg-base-100 shadow-sm p-6 animate-pulse">
              <div className="flex gap-4">
                <div className="skeleton h-16 w-16 rounded-xl" />
                <div className="flex-1 space-y-3">
                  <div className="skeleton h-6 w-1/3" />
                  <div className="skeleton h-4 w-2/3" />
                  <div className="flex gap-4">
                    <div className="skeleton h-4 w-20" />
                    <div className="skeleton h-4 w-20" />
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
