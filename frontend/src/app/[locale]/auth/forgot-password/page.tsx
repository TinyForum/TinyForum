// src/app/[locale]/auth/forgot-password/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import ForgotPasswordClient from "./ForgotPasswordClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("forgotPasswordTitle"),
    description: t("forgotPasswordDescription"),
  };
}

export default async function ForgotPasswordPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<ForgotPasswordSkeleton />}>
      <ForgotPasswordClient />
    </Suspense>
  );
}

function ForgotPasswordSkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-72 w-full max-w-md mx-auto rounded-2xl" />
    </div>
  );
}
