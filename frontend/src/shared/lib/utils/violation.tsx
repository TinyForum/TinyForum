import { ViolationRecord } from "@/shared/type/violation.type";

export const getViolationStatusBadge = (
  status: ViolationRecord["status"],
  t: (key: string) => string,
) => {
  switch (status) {
    case "pending":
      return <span className="badge badge-warning">{t("status_pending")}</span>;
    case "appealing":
      return <span className="badge badge-info">{t("status_appealing")}</span>;
    case "resolved":
      return (
        <span className="badge badge-success">{t("status_resolved")}</span>
      );
    case "rejected":
      return <span className="badge badge-error">{t("status_rejected")}</span>;
  }
};
