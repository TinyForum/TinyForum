// app/[locale]/questions/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import QuestionsClient from "./QuestionsClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("questionsTitle"),
    description: t("questionsDescription"),
  };
}

export default async function QuestionsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<QuestionsSkeleton />}>
      <QuestionsClient />
    </Suspense>
  );
}

function QuestionsSkeleton() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-base-200 via-base-100 to-base-200">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="card bg-base-100 shadow-md border border-base-200">
          <div className="card-body p-0">
            <div className="space-y-3 p-4">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="card bg-base-100 border border-base-200">
                  <div className="card-body p-5">
                    <div className="h-5 bg-base-200 rounded w-3/4 mb-2 animate-pulse" />
                    <div className="h-4 bg-base-200 rounded w-1/2 mb-3 animate-pulse" />
                    <div className="flex gap-4">
                      <div className="h-3 bg-base-200 rounded w-16 animate-pulse" />
                      <div className="h-3 bg-base-200 rounded w-16 animate-pulse" />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
