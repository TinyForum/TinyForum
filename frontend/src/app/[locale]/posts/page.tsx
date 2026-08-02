// src/app/[locale]/posts/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import PostsClient from "./PostsClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("postsTitle"),
    description: t("postsDescription"),
  };
}

export default async function PostsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<PostsSkeleton />}>
      <PostsClient />
    </Suspense>
  );
}

function PostsSkeleton() {
  return (
    <div className="max-w-3xl mx-auto space-y-4">
      <div className="flex items-center gap-3">
        <div className="skeleton h-6 w-6 rounded-xl" />
        <div className="skeleton h-8 w-48 rounded-xl" />
      </div>
      <div className="skeleton h-28 w-full rounded-xl" />
      <div className="skeleton h-28 w-full rounded-xl" />
      <div className="skeleton h-28 w-full rounded-xl" />
      <div className="skeleton h-28 w-full rounded-xl" />
      <div className="skeleton h-28 w-full rounded-xl" />
    </div>
  );
}
