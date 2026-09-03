"use client";

import { useTranslations } from "next-intl";

export default function NotificationSettings() {
  const t = useTranslations("Settings");
  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">{t("notification_settings")}</h1>
      <div className="card bg-base-100 border border-base-200 shadow-sm p-6">
        <p>{t("notification_settings_development")}</p>
      </div>
    </div>
  );
}
