import { usePointsData } from "@/hooks/admin/usePointsData";
import { Wallet, TrendingUp, Gift } from "lucide-react";

// ==================== 积分管理组件 ====================
export function PointsManager({ t }: { t: (key: string) => string }) {
  const { pointsRecords, stats, awardPoints, deductPoints, setExchangeRate } = usePointsData();

  return (
    <div className="space-y-6">
      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-primary">
            <Wallet className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("total_points_circulation")}</div>
          <div className="stat-value text-primary">{stats?.totalPoints || 0}</div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-secondary">
            <TrendingUp className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("today_awarded")}</div>
          <div className="stat-value text-secondary">{stats?.todayAwarded || 0}</div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-accent">
            <Gift className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("exchange_rate")}</div>
          <div className="stat-value text-accent">{stats?.exchangeRate || 100}</div>
        </div>
      </div>

      {/* 操作面板 */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="font-semibold mb-4">{t("points_operations")}</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t("user_id_or_username")}</span>
              </label>
              <input type="text" className="input input-bordered" />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t("points_amount")}</span>
              </label>
              <input type="number" className="input input-bordered" />
            </div>
          </div>
          <div className="flex gap-2 mt-4">
            <button className="btn btn-success btn-sm">{t("award_points")}</button>
            <button className="btn btn-warning btn-sm">{t("deduct_points")}</button>
          </div>
        </div>
      </div>

      {/* 积分记录 */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="font-semibold mb-4">{t("recent_points_records")}</h3>
          <div className="overflow-x-auto">
            <table className="table table-sm">
              <thead>
                <tr>
                  <th>{t("user")}</th>
                  <th>{t("operation")}</th>
                  <th>{t("points_change")}</th>
                  <th>{t("reason")}</th>
                  <th>{t("time")}</th>
                </tr>
              </thead>
              <tbody>
                {pointsRecords?.map((record) => (
                  <tr key={record.id}>
                    <td>{record.username}</td>
                    <td>
                      <span className={`badge badge-sm ${record.type === "award" ? "badge-success" : "badge-error"}`}>
                        {record.type === "award" ? t("award") : t("deduct")}
                      </span>
                    </td>
                    <td className={record.type === "award" ? "text-success" : "text-error"}>
                      {record.type === "award" ? "+" : "-"}{record.amount}
                    </td>
                    <td>{record.reason}</td>
                    <td className="text-xs">{record.created_at}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}