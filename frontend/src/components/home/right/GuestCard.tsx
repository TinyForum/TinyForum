"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { User } from "lucide-react";

export function GuestCard() {
  const t = useTranslations("Sidebar");

  return (
    <div className="rounded-lg border bg-card p-4 text-center">
      <User className="w-12 h-12 mx-auto text-muted-foreground mb-2" />
      <p className="text-sm text-muted-foreground mb-3">{t("not_logged_in")}</p>
      <Link
        href="/auth/login"
        className="inline-block w-full bg-primary text-primary-foreground rounded-lg px-4 py-2 text-sm hover:bg-primary/90 transition-colors"
      >
        {t("login")}
      </Link>
      <Link
        href="/auth/register"
        className="inline-block w-full mt-2 text-sm text-muted-foreground hover:text-primary transition-colors"
      >
        {t("register")}
      </Link>
    </div>
  );
}