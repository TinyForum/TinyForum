// components/user/ViolationPanel.tsx
"use client";

import { useTranslations } from "next-intl";
import { useState, useEffect } from "react";

// 违规记录类型定义
interface ViolationRecord {
  id: string;
  reason: string; // 违规原因
  date: string; // 违规时间
  status: "pending" | "appealing" | "resolved" | "rejected";
  punishment: string; // 处罚措施（如“禁言3天”、“警告”）
  appealDeadline?: string; // 申诉截止时间
}

// Mock 数据 - 实际应从 API 获取
const mockViolations: ViolationRecord[] = [
  // TODO: 替换为真实 API 调用
  {
    id: "1",
    reason: "发布广告信息",
    date: "2025-03-15",
    status: "pending",
    punishment: "禁言 7 天",
    appealDeadline: "2025-03-22",
  },
  {
    id: "2",
    reason: "人身攻击",
    date: "2025-04-01",
    status: "appealing",
    punishment: "警告一次",
    appealDeadline: "2025-04-08",
  },
  {
    id: "3",
    reason: "刷屏",
    date: "2025-02-10",
    status: "resolved",
    punishment: "禁言 1 天",
  },
];

export function ViolationPanel() {
  const t = useTranslations("Violation");
  const [violations, setViolations] = useState<ViolationRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [appealModal, setAppealModal] = useState<{
    open: boolean;
    violationId: string;
    reason: string;
  }>({ open: false, violationId: "", reason: "" });
  const [appealText, setAppealText] = useState("");

  // 获取违规记录
  useEffect(() => {
    // 替换为真实 API 调用
    const fetchViolations = async () => {
      setLoading(true);
      try {
        // const res = await userApi.getViolations();
        // setViolations(res.data);
        await new Promise((resolve) => setTimeout(resolve, 500));
        setViolations(mockViolations);
      } finally {
        setLoading(false);
      }
    };
    fetchViolations();
  }, []);

  // 提交申诉
  const handleAppealSubmit = async () => {
    if (!appealText.trim()) return;
    try {
      // await userApi.appealViolation(appealModal.violationId, appealText);
      // 更新本地状态
      setViolations((prev) =>
        prev.map((v) =>
          v.id === appealModal.violationId ? { ...v, status: "appealing" } : v,
        ),
      );
      setAppealModal({ open: false, violationId: "", reason: "" });
      setAppealText("");
      alert(t("appeal_submitted")); // 替换为 toast
    } catch {
      alert(t("appeal_failed"));
    }
  };

  const getStatusBadge = (status: ViolationRecord["status"]) => {
    switch (status) {
      case "pending":
        return (
          <span className="badge badge-warning">{t("status_pending")}</span>
        );
      case "appealing":
        return (
          <span className="badge badge-info">{t("status_appealing")}</span>
        );
      case "resolved":
        return (
          <span className="badge badge-success">{t("status_resolved")}</span>
        );
      case "rejected":
        return (
          <span className="badge badge-error">{t("status_rejected")}</span>
        );
    }
  };

  const canAppeal = (record: ViolationRecord) => {
    if (record.status !== "pending") return false;
    if (record.appealDeadline && new Date(record.appealDeadline) < new Date())
      return false;
    return true;
  };

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 违规记录卡片 */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="text-lg font-bold">{t("violation_history")}</h3>
          {violations.length === 0 ? (
            <div className="text-center py-8 text-base-content/60">
              {t("no_violations")}
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>{t("reason")}</th>
                    <th>{t("date")}</th>
                    <th>{t("punishment")}</th>
                    <th>{t("status")}</th>
                    <th>{t("action")}</th>
                  </tr>
                </thead>
                <tbody>
                  {violations.map((record) => (
                    <tr key={record.id}>
                      <td>{record.reason}</td>
                      <td>{record.date}</td>
                      <td>{record.punishment}</td>
                      <td>{getStatusBadge(record.status)}</td>
                      <td>
                        {canAppeal(record) ? (
                          <button
                            className="btn btn-sm btn-outline btn-primary"
                            onClick={() =>
                              setAppealModal({
                                open: true,
                                violationId: record.id,
                                reason: record.reason,
                              })
                            }
                          >
                            {t("appeal")}
                          </button>
                        ) : (
                          <span className="text-sm text-base-content/40">
                            {record.status === "appealing"
                              ? t("appealing_in_progress")
                              : record.status === "resolved"
                                ? t("case_closed")
                                : t("cannot_appeal")}
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>

      {/* 违规处理说明 */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <h3 className="text-lg font-bold">{t("violation_policy")}</h3>
          <ul className="list-disc pl-5 space-y-1 text-sm text-base-content/80">
            <li>{t("policy_1")}</li>
            <li>{t("policy_2")}</li>
            <li>{t("policy_3")}</li>
          </ul>
        </div>
      </div>

      {/* 申诉弹窗 */}
      {appealModal.open && (
        <dialog className="modal modal-open" open>
          <div className="modal-box">
            <h3 className="font-bold text-lg">{t("appeal_title")}</h3>
            <p className="py-2">
              {t("appeal_for")}: {appealModal.reason}
            </p>
            <textarea
              className="textarea textarea-bordered w-full mt-2"
              rows={4}
              placeholder={t("appeal_placeholder")}
              value={appealText}
              onChange={(e) => setAppealText(e.target.value)}
            />
            <div className="modal-action">
              <button className="btn btn-primary" onClick={handleAppealSubmit}>
                {t("submit_appeal")}
              </button>
              <button
                className="btn btn-ghost"
                onClick={() =>
                  setAppealModal({ open: false, violationId: "", reason: "" })
                }
              >
                {t("cancel")}
              </button>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button
              onClick={() =>
                setAppealModal({ open: false, violationId: "", reason: "" })
              }
            >
              close
            </button>
          </form>
        </dialog>
      )}
    </div>
  );
}
