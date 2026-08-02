// src/app/[locale]/auth/reset-password/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import ResetPasswordClient from "./ResetPasswordClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("resetPasswordTitle"),
    description: t("resetPasswordDescription"),
  };
}

export default async function ResetPasswordPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<ResetPasswordSkeleton />}>
      <ResetPasswordClient />
    </Suspense>
  );
}

function ResetPasswordSkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-80 w-full max-w-md mx-auto rounded-2xl" />
    </div>
  );
}
