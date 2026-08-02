// app/[locale]/questions/[id]/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import QuestionDetailClient from "./QuestionDetailClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("questionDetailTitle"),
    description: t("questionDetailDescription"),
  };
}

export default async function QuestionDetailPage({
  params,
}: {
  params: Promise<{ locale: string; id: string }>;
}) {
  const { id } = await params;

  return (
    <Suspense fallback={<QuestionDetailSkeleton />}>
      <QuestionDetailClient id={id} />
    </Suspense>
  );
}

function QuestionDetailSkeleton() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="animate-pulse">
          <div className="h-5 w-24 bg-base-200 rounded mb-6" />
          <div className="bg-base-100 rounded-2xl shadow-sm p-6 mb-6">
            <div className="h-7 bg-base-200 rounded-lg w-3/4 mb-4" />
            <div className="flex items-center gap-4 mb-4">
              <div className="h-4 w-24 bg-base-200 rounded" />
              <div className="h-4 w-32 bg-base-200 rounded" />
            </div>
            <div className="space-y-2">
              <div className="h-4 bg-base-200 rounded w-full" />
              <div className="h-4 bg-base-200 rounded w-full" />
              <div className="h-4 bg-base-200 rounded w-2/3" />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
