// src/app/[locale]/boards/[slug]/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import BoardDetailClient from "./BoardDetailClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("boardDetailTitle"),
    description: t("boardDetailDescription"),
  };
}

export default async function BoardDetailPage({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { slug } = await params;

  return (
    <Suspense fallback={<BoardDetailSkeleton />}>
      <BoardDetailClient slug={slug} />
    </Suspense>
  );
}

function BoardDetailSkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-32 w-full max-w-5xl mx-auto rounded-2xl" />
      <div className="skeleton h-24 w-full max-w-5xl mx-auto rounded-xl" />
    </div>
  );
}
