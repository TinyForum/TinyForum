"use client";

import { useTranslations } from "next-intl";

export default function SiteInfo() {
  const t = useTranslations("Post");

  return (
    <div className="card bg-base-100 border border-base-300 shadow-sm">
      <div className="card-body p-4 text-xs text-base-content/50 space-y-1">
        <p className="font-medium text-base-content/70">{t("about")}</p>
        <p>{t("description")}</p>
        <p className="pt-1">
          © {new Date().getFullYear()} {t("copyright")}
        </p>
      </div>
    </div>
  );
}
