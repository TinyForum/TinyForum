// src/app/[locale]/posts/[id]/edit/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import EditPostClient from "./EditPostClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("editPostTitle"),
    description: t("editPostDescription"),
  };
}

export default async function EditPostPage({
  params,
}: {
  params: Promise<{ locale: string; id: string }>;
}) {
  const { id } = await params;
  const postId = Number(id);

  return (
    <Suspense fallback={<EditPostSkeleton />}>
      <EditPostClient postId={postId} />
    </Suspense>
  );
}

function EditPostSkeleton() {
  return (
    <div className="max-w-3xl mx-auto space-y-4">
      <div className="skeleton h-8 w-40 rounded-xl" />
      <div className="skeleton h-64 w-full rounded-xl" />
      <div className="skeleton h-40 w-full rounded-xl" />
      <div className="skeleton h-10 w-32 rounded-xl" />
    </div>
  );
}
