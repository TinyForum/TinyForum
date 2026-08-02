// src/app/[locale]/boards/[slug]/apply/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import BoardApplyClient from "./BoardApplyClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("boardApplyTitle"),
    description: t("boardApplyDescription"),
  };
}

export default async function BoardApplyPage({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { slug } = await params;

  return (
    <Suspense fallback={<BoardApplySkeleton />}>
      <BoardApplyClient slug={slug} />
    </Suspense>
  );
}

function BoardApplySkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-96 w-full max-w-2xl mx-auto rounded-2xl" />
    </div>
  );
}
