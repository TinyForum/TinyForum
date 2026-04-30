"use client";

import { User } from "@/shared/api/types";
import Avatar from "./Avatar";
import { useTranslations } from "next-intl";
import { formatDate } from "@/shared/lib/utils";

interface UserInfoCardProps {
  user: User | null;
  isAdmin: boolean;
  isModerator: boolean;
}

export function UserInfoCard({
  user,
  isAdmin,
  isModerator,
}: UserInfoCardProps) {
  const t = useTranslations("User");
  return (
    <div className="card bg-base-100 border border-base-300">
      <div className="card-body">
        <div className="flex items-center gap-4">
          <div className="avatar">
            <div className="w-20 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
              <Avatar avatarUrl={user?.avatar} size={"full"} />
            </div>
          </div>
          <div>
            <h2 className="text-2xl font-bold">{user?.username}</h2>
            <p className="text-base-content/60">{user?.email}</p>
            <div className="flex gap-2 mt-2">
              {isAdmin && (
                <span className="badge badge-primary">{t("admin")}</span>
              )}
              {isModerator && (
                <span className="badge badge-secondary">{t("moderator")}</span>
              )}
              <span className="badge badge-outline">
                {user?.created_at &&
                  `${t("registered_at") + ": " + formatDate(user?.created_at)}`}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
