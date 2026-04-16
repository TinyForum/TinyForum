// components/moderator/ReportedContentTable.tsx
interface ReportedContentTableProps {
  reports: any[];
  onDeletePost?: (postId: number) => void;
  onBanUser?: (data: { userId: number; reason: string }) => void;
  isDeleting: boolean;
  isBanning: boolean;
  permissions: any;
  t: (key: string) => string;
}

export function ReportedContentTable({ reports, onDeletePost, onBanUser, isDeleting, isBanning, permissions, t }: ReportedContentTableProps) {
  if (!reports?.length) {
    return <div className="text-center py-8 text-base-content/50">{t("no_reports")}</div>;
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
          {reports.map((report) => (
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
                    onClick={() => onBanUser({ userId: report.reporter_id, reason: report.reason })}
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