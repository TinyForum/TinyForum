"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { AlertTriangle } from "lucide-react";

export function ViolationStatus() {
  const t = useTranslations("Sidebar");

  return (
    <div className="rounded-lg border bg-card">
      <div className="p-3 border-b">
        <h3 className="font-semibold flex items-center gap-2">
          <AlertTriangle className="w-4 h-4 text-yellow-500" />
          {t("violation_status")}
        </h3>
      </div>
      <div className="p-3 space-y-2 text-sm">
        <div className="flex justify-between">
          <span className="text-muted-foreground">{t("violation_count")}</span>
          <span className="font-medium text-green-600">0</span>
        </div>
        <div className="flex justify-between">
          <span className="text-muted-foreground">{t("account_status")}</span>
          <span className="font-medium text-green-600">{t("normal")}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-muted-foreground">{t("warning_level")}</span>
          <span className="font-medium text-green-600">{t("none")}</span>
        </div>
      </div>
    </div>
  );
}
