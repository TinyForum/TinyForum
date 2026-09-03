"use client";

import { useTranslations } from "next-intl";

export default function AppearanceSettings() {
  const t = useTranslations("Settings");
  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">{t("appearance_theme")}</h1>
      <div className="card bg-base-100 border border-base-200 shadow-sm p-6">
        <p>{t("theme_settings_development")}</p>
      </div>
    </div>
  );
}
