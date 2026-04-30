// components/reviewer/ReviewSettings.tsx
"use client";

import { useTranslations } from "next-intl";

export function ReviewSettings() {
  const t = useTranslations("Reviewer");
  return (
    <div className="space-y-6">
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="text-lg font-bold">{t("auto_review_settings")}</h3>

          <div className="form-control">
            <label className="cursor-pointer label">
              <span className="label-text">{t("auto_approve_trusted")}</span>
              <input type="checkbox" className="toggle toggle-primary" />
            </label>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("min_reputation")}</span>
            </label>
            <input
              type="number"
              className="input input-bordered"
              placeholder="100"
            />
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("keywords_blocklist")}</span>
            </label>
            <textarea
              className="textarea textarea-bordered"
              rows={4}
              placeholder={t("keywords_placeholder")}
            />
          </div>

          <div className="mt-4">
            <button className="btn btn-primary">{t("save_settings")}</button>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="text-lg font-bold">{t("notification_settings")}</h3>

          <div className="form-control">
            <label className="cursor-pointer label">
              <span className="label-text">{t("notify_on_new_report")}</span>
              <input
                type="checkbox"
                className="toggle toggle-primary"
                defaultChecked
              />
            </label>
          </div>

          <div className="form-control">
            <label className="cursor-pointer label">
              <span className="label-text">{t("daily_summary")}</span>
              <input type="checkbox" className="toggle toggle-primary" />
            </label>
          </div>
        </div>
      </div>
    </div>
  );
}
