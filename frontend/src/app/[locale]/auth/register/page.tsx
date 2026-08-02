// src/app/[locale]/auth/register/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import RegisterClient from "./RegisterClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("registerTitle"),
    description: t("registerDescription"),
  };
}

export default async function RegisterPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<RegisterSkeleton />}>
      <RegisterClient />
    </Suspense>
  );
}

function RegisterSkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-96 w-full max-w-md mx-auto rounded-2xl" />
    </div>
  );
}
