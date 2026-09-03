"use client";

import { useTranslations } from "next-intl";
import { Smartphone } from "lucide-react";

export function TwoFactorAuthCard() {
  const t = useTranslations("Settings");
  return (
    <div className="card bg-primary/5 border border-primary/20">
      <div className="card-body p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/20 rounded-lg">
              <Smartphone className="w-5 h-5 text-primary" />
            </div>
            <div>
              <p className="font-medium">{t("two_factor_auth")}</p>
              <p className="text-xs text-base-content/60">
                {t("two_factor_auth_desc")}
              </p>
            </div>
          </div>
          <button className="btn btn-sm btn-outline btn-primary">{t("enable")}</button>
        </div>
      </div>
    </div>
  );
}
