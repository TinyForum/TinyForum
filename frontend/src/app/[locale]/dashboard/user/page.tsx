// src/app/[locale]/dashboard/user/page.tsx (服务端组件)
import { Suspense } from "react";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import UserDashboardClient from "./UserDashboardClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "Site" });
  return {
    title: t("userDashboardTitle"),
    description: t("userDashboardDescription"),
  };
}

export default async function UserDashboardPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  await params;

  return (
    <Suspense fallback={<UserDashboardSkeleton />}>
      <UserDashboardClient />
    </Suspense>
  );
}

function UserDashboardSkeleton() {
  return (
    <div className="max-w-7xl mx-auto px-4 py-8 space-y-4">
      <div className="skeleton h-16 w-64 rounded-xl" />
      <div className="skeleton h-12 w-full rounded-xl" />
      <div className="skeleton h-32 w-full rounded-xl" />
      <div className="skeleton h-32 w-full rounded-xl" />
    </div>
  );
}
