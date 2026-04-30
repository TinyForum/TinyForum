// components/moderator/ReportedContentTable.tsx

import { useTranslations } from "next-intl";

export interface Report {
  id: number;
  target_id: number;
  target_title: string;
  reporter_id: number;
  reporter_name: string;
  reason: string;
  target_type: string;
  status: string;
  created_at: string;
}

interface ModeratorPermissions {
  canDeletePost: boolean;
  canPinPost: boolean;
  canEditAnyPost: boolean;
  canManageModerator: boolean;
  canBanUser: boolean;
}

interface ReportedContentTableProps {
  reports: Report[];
  onDeletePost?: (postId: number) => void;
  onBanUser?: (data: { userId: number; reason: string }) => void;
  isDeleting: boolean;
  isBanning: boolean;
  permissions: ModeratorPermissions;
  t: (key: string) => string;
}

export function ReportedContentTable({
  reports,
  onDeletePost,
  onBanUser,
  isDeleting,
  isBanning,
  permissions,
}: ReportedContentTableProps) {
  const t = useTranslations("Moderator");
  if (!reports?.length) {
    return (
      <div className="text-center py-8 text-base-content/50">
        {t("no_reports")}
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="table">
        <thead>
          <tr>
            <th>{t("content")}</th>
            <th>{t("reporter")}</th>
            <th>{t("reason")}</th>
            <th>{t("actions")}</th>
          </tr>
        </thead>
        <tbody>
          {reports.map((report: Report) => (
            <tr key={report.id}>
              <td>{report.target_title}</td>
              <td>{report.reporter_name}</td>
              <td>{report.reason}</td>
              <td className="flex gap-2">
                {permissions.canDeletePost && onDeletePost && (
                  <button
                    className="btn btn-xs btn-error"
                    onClick={() => onDeletePost(report.target_id)}
                    disabled={isDeleting}
                  >
                    {t("delete_post")}
                  </button>
                )}
                {permissions.canBanUser && onBanUser && (
                  <button
                    className="btn btn-xs btn-warning"
                    onClick={() =>
                      onBanUser({
                        userId: report.reporter_id,
                        reason: report.reason,
                      })
                    }
                    disabled={isBanning}
                  >
                    {t("ban_user")}
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
