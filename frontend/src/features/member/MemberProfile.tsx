// components/member/MemberProfile.tsx
"use client";

import { useTranslations } from "next-intl";
import { useState } from "react";
import Avatar from "../user/components/Avatar";

export function MemberProfile() {
  const t = useTranslations("Member");
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  // TODO: 从 props 获取用户信息
  // const { member, updateProfile, changePassword, deleteAccount } = useMemberProfile();

  return (
    <div className="space-y-6">
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="text-lg font-bold">{t("profile_info")}</h3>

          <div className="flex items-center gap-6 mb-6">
            <div className="avatar">
              <div className="w-24 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                <Avatar username="Undefined" />
              </div>
            </div>
            <div>
              <button className="btn btn-outline btn-sm">
                {t("change_avatar")}
              </button>
            </div>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("username")}</span>
            </label>
            <input
              type="text"
              className="input input-bordered"
              value="username"
              disabled
            />
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("email")}</span>
            </label>
            <input
              type="email"
              className="input input-bordered"
              value="user@example.com"
              disabled
            />
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("bio")}</span>
            </label>
            <textarea
              className="textarea textarea-bordered"
              rows={3}
              placeholder={t("bio_placeholder")}
            />
          </div>

          <div className="mt-4">
            <button className="btn btn-primary">{t("save_changes")}</button>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="text-lg font-bold">{t("change_password")}</h3>

          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("current_password")}</span>
            </label>
            <input type="password" className="input input-bordered" />
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("new_password")}</span>
            </label>
            <input type="password" className="input input-bordered" />
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("confirm_password")}</span>
            </label>
            <input type="password" className="input input-bordered" />
          </div>

          <div className="mt-4">
            <button className="btn btn-primary">{t("update_password")}</button>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="text-lg font-bold text-error">{t("danger_zone")}</h3>
          <p className="text-sm text-base-content/60">
            {t("delete_account_warning")}
          </p>

          {!showDeleteConfirm ? (
            <button
              onClick={() => setShowDeleteConfirm(true)}
              className="btn btn-error btn-sm w-32"
            >
              {t("delete_account")}
            </button>
          ) : (
            <div className="flex gap-2">
              <button className="btn btn-error btn-sm">
                {t("confirm_delete")}
              </button>
              <button
                onClick={() => setShowDeleteConfirm(false)}
                className="btn btn-ghost btn-sm"
              >
                {t("cancel")}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
