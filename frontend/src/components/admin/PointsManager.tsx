import { useState } from "react";
import { useScoreData } from "@/hooks/admin/useScoreData";
import { Wallet, TrendingUp, Gift, Search, Plus, Minus, RefreshCw } from "lucide-react";

// 类型定义
interface UserScoreRecord {
  id: number;
  username: string;
  avatar: string;
  score: number;
}

// ==================== 管理员管理用户积分组件 ====================
export function PointsManager({ t }: { t: (key: string) => string }) {
  const { 
    scoreRecords, 
    myScore,
    addScore,
    subtractScore,
    setScore,
    addScoreAsync,
    subtractScoreAsync,
    setScoreAsync,
    isAddingScore,
    isSubtractingScore,
    isSettingScore
  } = useScoreData();

  // 表单状态
  const [userId, setUserId] = useState("");
  const [pointsAmount, setPointsAmount] = useState<number>(0);
  const [reason, setReason] = useState("");
  const [operationType, setOperationType] = useState<"add" | "subtract" | "set">("add");
  const [searchKeyword, setSearchKeyword] = useState("");

  // 修复：安全地获取积分记录数组
  // scoreRecords 可能是 ApiResponse 对象，需要提取 data 字段
  const getRecordsArray = (): UserScoreRecord[] => {
    if (!scoreRecords) return [];
    
    // 如果是 ApiResponse 对象且有 data 属性
    if (typeof scoreRecords === 'object' && 'data' in scoreRecords) {
      const data = (scoreRecords as any).data;
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
    totalPoints: records.length > 0 
      ? records.reduce((sum: number, record: UserScoreRecord) => sum + (record.score || 0), 0)
      : 0,
    todayAwarded: 0, // 需要后端接口支持
    exchangeRate: 100,
  };

  // 筛选积分记录
  const filteredRecords = records.filter((record: UserScoreRecord) => 
    searchKeyword === "" || 
    record.username?.toLowerCase().includes(searchKeyword.toLowerCase()) ||
    record.id?.toString().includes(searchKeyword)
  );

  // 处理增加积分
  const handleAddPoints = async () => {
    if (!userId || pointsAmount <= 0) {
      alert(t("please_fill_correct_info"));
      return;
    }
    if (!reason) {
      alert(t("please_enter_reason"));
      return;
    }
    
    try {
      await addScoreAsync({
        userId: Number(userId),
        increment: pointsAmount,
        reason: reason
      });
      // 清空表单
      setUserId("");
      setPointsAmount(0);
      setReason("");
      alert(t("points_awarded_success"));
    } catch (error) {
      alert(t("operation_failed"));
    }
  };

  // 处理扣除积分
  const handleDeductPoints = async () => {
    if (!userId || pointsAmount <= 0) {
      alert(t("please_fill_correct_info"));
      return;
    }
    if (!reason) {
      alert(t("please_enter_reason"));
      return;
    }
    
    try {
      await subtractScoreAsync({
        userId: Number(userId),
        decrement: pointsAmount,
        reason: reason
      });
      // 清空表单
      setUserId("");
      setPointsAmount(0);
      setReason("");
      alert(t("points_deducted_success"));
    } catch (error) {
      alert(t("operation_failed"));
    }
  };

  // 处理设置积分
  const handleSetPoints = async () => {
    if (!userId || pointsAmount < 0) {
      alert(t("please_fill_correct_info"));
      return;
    }
    if (!reason) {
      alert(t("please_enter_reason"));
      return;
    }
    
    try {
      await setScoreAsync({
        userId: Number(userId),
        score: pointsAmount,
        reason: reason
      });
      // 清空表单
      setUserId("");
      setPointsAmount(0);
      setReason("");
      alert(t("points_set_success"));
    } catch (error) {
      alert(t("operation_failed"));
    }
  };

  return (
    <div className="space-y-6">
      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-primary">
            <Wallet className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("total_points_circulation")}</div>
          <div className="stat-value text-primary">{stats.totalPoints.toLocaleString()}</div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-secondary">
            <TrendingUp className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("today_awarded")}</div>
          <div className="stat-value text-secondary">{stats.todayAwarded.toLocaleString()}</div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-accent">
            <Gift className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("exchange_rate")}</div>
          <div className="stat-value text-accent">{stats.exchangeRate}</div>
        </div>
      </div>

      {/* 操作面板 */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="font-semibold mb-4">{t("points_operations")}</h3>
          
          {/* 操作类型选择 */}
          <div className="flex gap-2 mb-4">
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
                <span className="label-text">{t("user_id_or_username")}</span>
              </label>
              <input 
                type="text" 
                className="input input-bordered" 
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                placeholder={t("enter_user_id_or_username")}
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t("points_amount")}</span>
              </label>
              <input 
                type="number" 
                className="input input-bordered" 
                value={pointsAmount}
                onChange={(e) => setPointsAmount(Number(e.target.value))}
                placeholder={t("enter_points_amount")}
              />
            </div>
          </div>
          
          <div className="form-control">
            <label className="label">
              <span className="label-text">{t("reason")}</span>
            </label>
            <input 
              type="text" 
              className="input input-bordered" 
              value={reason}
              onChange={(e) => setReason(e.target.value)}
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
                {isAddingScore ? t("processing") : t("award_points")}
              </button>
            )}
            {operationType === "subtract" && (
              <button 
                className="btn btn-warning" 
                onClick={handleDeductPoints}
                disabled={isSubtractingScore}
              >
                {isSubtractingScore ? t("processing") : t("deduct_points")}
              </button>
            )}
            {operationType === "set" && (
              <button 
                className="btn btn-primary" 
                onClick={handleSetPoints}
                disabled={isSettingScore}
              >
                {isSettingScore ? t("processing") : t("set_points")}
              </button>
            )}
          </div>
        </div>
      </div>

      {/* 积分记录 */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h3 className="font-semibold">{t("recent_points_records")}</h3>
            <div className="form-control">
              <div className="flex gap-2">
                <input 
                  type="text" 
                  className="input input-bordered input-sm" 
                  placeholder={t("search_user")}
                  value={searchKeyword}
                  onChange={(e) => setSearchKeyword(e.target.value)}
                />
                <button className="btn btn-ghost btn-sm">
                  <Search className="w-4 h-4" />
                </button>
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
                            <img src={record.avatar} alt={record.username} className="w-6 h-6 rounded-full" />
                          )}
                          <span>{record.username}</span>
                          <span className="text-xs text-base-content/50">#{record.id}</span>
                        </div>
                      </td>
                      <td className="font-semibold">{record.score}</td>
                      <td>
                        <span className="text-xs text-base-content/50">
                          {t("points_management")}
                        </span>
                      </td>
                      <td className="text-xs">
                        {new Date().toLocaleDateString()}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={4} className="text-center py-8 text-base-content/50">
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