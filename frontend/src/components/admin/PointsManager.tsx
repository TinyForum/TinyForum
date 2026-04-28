import { useState } from "react";
import { useScoreData } from "@/hooks/admin/useScoreData";
import {
  Wallet,
  TrendingUp,
  Gift,
  Search,
  Plus,
  Minus,
  RefreshCw,
} from "lucide-react";
import Image from "next/image";
import { ApiResponse } from "@/lib/api/types";
import { useTranslations } from "next-intl";

// 类型定义
interface UserScoreRecord {
  id: number;
  username: string;
  avatar: string;
  score: number;
  created_at?: string;
  updated_at?: string;
}



// 操作类型
type OperationType = "add" | "subtract" | "set";

// 错误响应类型
interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

// ==================== 管理员管理用户积分组件 ====================
export function PointsManager() {
  const {
    scoreRecords,
    addScoreAsync,
    subtractScoreAsync,
    setScoreAsync,
    isAddingScore,
    isSubtractingScore,
    isSettingScore,
  } = useScoreData();

  const t  = useTranslations("Common");
  // 表单状态
  const [userId, setUserId] = useState<string>("");
  const [pointsAmount, setPointsAmount] = useState<number>(0);
  const [reason, setReason] = useState<string>("");
  const [operationType, setOperationType] = useState<OperationType>("add");
  const [searchKeyword, setSearchKeyword] = useState<string>("");

  // 安全地获取积分记录数组
  const getRecordsArray = (): UserScoreRecord[] => {
    if (!scoreRecords) return [];

    // 如果是 ApiResponse 对象且有 data 属性
    if (typeof scoreRecords === "object" && "data" in scoreRecords) {
      const data = (scoreRecords as ApiResponse<UserScoreRecord[]>).data;
      return Array.isArray(data) ? data : [];
    }

    // 如果直接是数组
    if (Array.isArray(scoreRecords)) {
      return scoreRecords;
    }

    return [];
  };

  const records = getRecordsArray();

  // 计算统计数据
  const stats = {
    totalPoints:
      records.length > 0
        ? records.reduce(
            (sum: number, record: UserScoreRecord) => sum + (record.score || 0),
            0,
          )
        : 0,
    todayAwarded: 0,
    exchangeRate: 100,
  };

  // 筛选积分记录
  const filteredRecords = records.filter(
    (record: UserScoreRecord) =>
      searchKeyword === "" ||
      record.username?.toLowerCase().includes(searchKeyword.toLowerCase()) ||
      record.id?.toString().includes(searchKeyword),
  );

  // 显示错误提示
  const showError = (message: string) => {
    alert(message);
  };

  // 显示成功提示
  const showSuccess = (message: string) => {
    alert(message);
  };

  // 处理增加积分
  const handleAddPoints = async (): Promise<void> => {
    if (!userId || pointsAmount <= 0) {
      showError(t("please_fill_correct_info"));
      return;
    }
    if (!reason) {
      showError(t("please_enter_reason"));
      return;
    }

    try {
      await addScoreAsync({
        userId: Number(userId),
        increment: pointsAmount,
        reason: reason,
      });
      // 清空表单
      setUserId("");
      setPointsAmount(0);
      setReason("");
      showSuccess(t("points_awarded_success"));
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      showError(error?.response?.data?.message || t("operation_failed"));
    }
  };

  // 处理扣除积分
  const handleDeductPoints = async (): Promise<void> => {
    if (!userId || pointsAmount <= 0) {
      showError(t("please_fill_correct_info"));
      return;
    }
    if (!reason) {
      showError(t("please_enter_reason"));
      return;
    }

    try {
      await subtractScoreAsync({
        userId: Number(userId),
        decrement: pointsAmount,
        reason: reason,
      });
      // 清空表单
      setUserId("");
      setPointsAmount(0);
      setReason("");
      showSuccess(t("points_deducted_success"));
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      showError(error?.response?.data?.message || t("operation_failed"));
    }
  };

  // 处理设置积分
  const handleSetPoints = async (): Promise<void> => {
    if (!userId || pointsAmount < 0) {
      showError(t("please_fill_correct_info"));
      return;
    }
    if (!reason) {
      showError(t("please_enter_reason"));
      return;
    }

    try {
      await setScoreAsync({
        userId: Number(userId),
        score: pointsAmount,
        reason: reason,
      });
      // 清空表单
      setUserId("");
      setPointsAmount(0);
      setReason("");
      showSuccess(t("points_set_success"));
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      showError(error?.response?.data?.message || t("operation_failed"));
    }
  };

  return (
    <div className="space-y-6">
      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-primary">
            <Wallet className="w-6 h-6" />
          </div>
          <div className="stat-title text-base-content/60">{t("total_points_circulation")}</div>
          <div className="stat-value text-primary text-3xl font-bold">
            {stats.totalPoints.toLocaleString()}
          </div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-secondary">
            <TrendingUp className="w-6 h-6" />
          </div>
          <div className="stat-title text-base-content/60">{t("today_awarded")}</div>
          <div className="stat-value text-secondary text-3xl font-bold">
            {stats.todayAwarded.toLocaleString()}
          </div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-accent">
            <Gift className="w-6 h-6" />
          </div>
          <div className="stat-title text-base-content/60">{t("exchange_rate")}</div>
          <div className="stat-value text-accent text-3xl font-bold">{stats.exchangeRate}</div>
        </div>
      </div>

      {/* 操作面板 */}
      <div className="card bg-base-100 border border-base-300 shadow-sm">
        <div className="card-body p-6">
          <h3 className="font-semibold text-lg mb-4">{t("points_operations")}</h3>

          {/* 操作类型选择 */}
          <div className="flex gap-2 mb-4 flex-wrap">
            <button
              className={`btn btn-sm ${operationType === "add" ? "btn-success" : "btn-ghost"}`}
              onClick={() => setOperationType("add")}
            >
              <Plus className="w-4 h-4" /> {t("award_points")}
            </button>
            <button
              className={`btn btn-sm ${operationType === "subtract" ? "btn-warning" : "btn-ghost"}`}
              onClick={() => setOperationType("subtract")}
            >
              <Minus className="w-4 h-4" /> {t("deduct_points")}
            </button>
            <button
              className={`btn btn-sm ${operationType === "set" ? "btn-primary" : "btn-ghost"}`}
              onClick={() => setOperationType("set")}
            >
              <RefreshCw className="w-4 h-4" /> {t("set_points")}
            </button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">{t("user_id_or_username")}</span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={userId}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => 
                  setUserId(e.target.value)
                }
                placeholder={t("enter_user_id_or_username")}
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">{t("points_amount")}</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={pointsAmount}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => 
                  setPointsAmount(Number(e.target.value))
                }
                placeholder={t("enter_points_amount")}
                min={0}
              />
            </div>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">{t("reason")}</span>
            </label>
            <input
              type="text"
              className="input input-bordered"
              value={reason}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => 
                setReason(e.target.value)
              }
              placeholder={t("enter_operation_reason")}
            />
          </div>

          <div className="flex gap-2 mt-4">
            {operationType === "add" && (
              <button
                className="btn btn-success"
                onClick={handleAddPoints}
                disabled={isAddingScore}
              >
                {isAddingScore ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : null}
                {isAddingScore ? t("processing") : t("award_points")}
              </button>
            )}
            {operationType === "subtract" && (
              <button
                className="btn btn-warning"
                onClick={handleDeductPoints}
                disabled={isSubtractingScore}
              >
                {isSubtractingScore ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : null}
                {isSubtractingScore ? t("processing") : t("deduct_points")}
              </button>
            )}
            {operationType === "set" && (
              <button
                className="btn btn-primary"
                onClick={handleSetPoints}
                disabled={isSettingScore}
              >
                {isSettingScore ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : null}
                {isSettingScore ? t("processing") : t("set_points")}
              </button>
            )}
          </div>
        </div>
      </div>

      {/* 积分记录 */}
      <div className="card bg-base-100 border border-base-300 shadow-sm">
        <div className="card-body p-6">
          <div className="flex justify-between items-center mb-4 flex-wrap gap-2">
            <h3 className="font-semibold text-lg">{t("recent_points_records")}</h3>
            <div className="form-control">
              <div className="flex gap-2">
                <input
                  type="text"
                  className="input input-bordered input-sm"
                  placeholder={t("search_user")}
                  value={searchKeyword}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => 
                    setSearchKeyword(e.target.value)
                  }
                />
                <Search className="w-4 h-4 text-base-content/40" />
              </div>
            </div>
          </div>
          <div className="overflow-x-auto">
            <table className="table table-sm">
              <thead>
                <tr>
                  <th>{t("user")}</th>
                  <th>{t("score")}</th>
                  <th>{t("operation_type")}</th>
                  <th>{t("operation_at")}</th>
                </tr>
              </thead>
              <tbody>
                {filteredRecords.length > 0 ? (
                  filteredRecords.map((record: UserScoreRecord) => (
                    <tr key={record.id}>
                      <td>
                        <div className="flex items-center gap-2">
                          {record.avatar && (
                            <div className="relative w-6 h-6 rounded-full overflow-hidden">
                              <Image
                                src={record.avatar}
                                alt={record.username}
                                fill
                                className="object-cover"
                                sizes="24px"
                              />
                            </div>
                          )}
                          <span>{record.username}</span>
                          <span className="text-xs text-base-content/50">
                            #{record.id}
                          </span>
                        </div>
                      </td>
                      <td className="font-semibold text-primary">
                        {record.score}
                      </td>
                      <td>
                        <span className="badge badge-ghost badge-sm">
                          {t("points_management")}
                        </span>
                      </td>
                      <td className="text-xs text-base-content/50">
                        {record.created_at 
                          ? new Date(record.created_at).toLocaleDateString()
                          : t("unknown")}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td
                      colSpan={4}
                      className="text-center py-8 text-base-content/50"
                    >
                      {t("no_records")}
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}