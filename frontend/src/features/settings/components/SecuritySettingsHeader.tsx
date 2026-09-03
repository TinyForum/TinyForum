"use client";

import { useTranslations } from "next-intl";

export function SecuritySettingsHeader() {
  const t = useTranslations("Settings");
  return (
    <div className="mb-6">
      <h1 className="text-2xl font-bold">{t("security_settings")}</h1>
      <p className="text-sm text-base-content/60 mt-1">{t("security_settings_desc")}</p>
    </div>
  );
}
