// src/app/[locale]/posts/new/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import NewPostClient from "./NewPostClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("newPostTitle"),
    description: t("newPostDescription"),
  };
}

export default async function NewPostPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<NewPostSkeleton />}>
      <NewPostClient />
    </Suspense>
  );
}

function NewPostSkeleton() {
  return (
    <div className="max-w-7xl mx-auto px-4 py-6 space-y-4">
      <div className="flex items-center gap-3">
        <div className="skeleton h-8 w-8 rounded-xl" />
        <div className="skeleton h-8 w-48 rounded-xl" />
      </div>
      <div className="flex flex-col lg:flex-row gap-6">
        <div className="lg:w-80 flex-shrink-0 space-y-4">
          <div className="skeleton h-72 w-full rounded-xl" />
        </div>
        <div className="flex-1 min-w-0 space-y-4">
          <div className="skeleton h-64 w-full rounded-xl" />
          <div className="skeleton h-10 w-32 rounded-xl" />
        </div>
      </div>
    </div>
  );
}
