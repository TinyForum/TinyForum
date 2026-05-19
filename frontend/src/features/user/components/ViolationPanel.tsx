// components/user/ViolationPanel.tsx
"use client";

import { useTranslations } from "next-intl";
import { useState, useEffect, useMemo } from "react";
import { getViolationStatusBadge } from "@/shared/lib/utils/violation";
import { ViolationRecord } from "@/shared/api/types/violation.model";
import { useUserViolation } from "../hooks/useViolation";
import type { ViolationVO } from "@/shared/api/modules/user/violation";
import toast from "react-hot-toast";

// 安全的数据转换函数（防御式）
function toViolationRecord(vo: ViolationVO): ViolationRecord {
  // 映射 UI 状态
  let uiStatus: "pending" | "appealing" | "resolved";
  if (vo.appeal_status && vo.appeal_status !== "none") {
    uiStatus = "appealing";
  } else if (vo.status === "resolved" || vo.status === "closed") {
    uiStatus = "resolved";
  } else {
    uiStatus = "pending";
  }

  // 生成可读的处罚内容
  let punishmentText = "";
  if (vo.punish_type) {
    if (vo.punish_expire_at) {
      punishmentText = `${vo.punish_type} 至 ${new Date(vo.punish_expire_at).toLocaleDateString()}`;
    } else {
      punishmentText = vo.punish_type;
    }
  }

  return {
    id: vo.id,
    reason: vo.reason,
    date: new Date(vo.created_at).toLocaleDateString(),
    status: uiStatus,
    punishment: punishmentText,
    appealDeadline: undefined,
  };
}

export function ViolationPanel() {
  const t = useTranslations("Violation");
  const {
    violations: rawViolations,
    loadViolations,
    isLoading,
    submitAppeal,
    isAppealing,
    appealError,
    error: loadError,
  } = useUserViolation();

  // 安全转换后的展示数据
  const [violations, setViolations] = useState<ViolationRecord[]>([]);
  const [appealModal, setAppealModal] = useState({
    open: false,
    violationId: "",
    violationReason: "",
  });
  const [appealText, setAppealText] = useState("");

  // 始终确保 rawViolations 是数组，如果不是则记录警告并转为空数组
  const safeViolations = useMemo(() => {
    if (Array.isArray(rawViolations)) {
      return rawViolations;
    }
    console.warn("useUserViolation 返回的 violations 不是数组:", rawViolations);
    return [];
  }, [rawViolations]);

  // 将安全的违规数据转换为 UI 数据
  useEffect(() => {
    setViolations(safeViolations.map(toViolationRecord));
  }, [safeViolations]);

  // 首次加载数据
  useEffect(() => {
    loadViolations();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const canAppeal = (record: ViolationRecord) => {
    if (record.status !== "pending") return false;
    if (record.appealDeadline && new Date(record.appealDeadline) < new Date())
      return false;
    return true;
  };

  const handleAppealSubmit = async () => {
    if (!appealText.trim()) return;
    const success = await submitAppeal(appealModal.violationId, appealText);
    if (success) {
      setAppealModal({ open: false, violationId: "", violationReason: "" });
      setAppealText("");
      // 使用更友好的提示，实际项目可替换为 toast
      toast.success(t("appeal_submitted"));
    } else {
      toast.error(appealError || t("appeal_failed"));
    }
  };

  // 加载中
  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  // 加载错误
  if (loadError) {
    return (
      <div className="card bg-base-100 border border-error/30">
        <div className="card-body text-center text-error">
          <p>{loadError}</p>
          <button
            className="btn btn-sm btn-outline mt-2"
            onClick={loadViolations}
          >
            {t("retry")}
          </button>
        </div>
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
                      <td>{getViolationStatusBadge(record.status, t)}</td>
                      <td>
                        {canAppeal(record) ? (
                          <button
                            className="btn btn-sm btn-outline btn-primary"
                            onClick={() =>
                              setAppealModal({
                                open: true,
                                violationId: record.id,
                                violationReason: record.reason,
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
              {t("appeal_for")}: {appealModal.violationReason}
            </p>
            <textarea
              className="textarea textarea-bordered w-full mt-2"
              rows={4}
              placeholder={t("appeal_placeholder")}
              value={appealText}
              onChange={(e) => setAppealText(e.target.value)}
              disabled={isAppealing}
            />
            <div className="modal-action">
              <button
                className="btn btn-primary"
                onClick={handleAppealSubmit}
                disabled={isAppealing || !appealText.trim()}
              >
                {isAppealing ? (
                  <span className="loading loading-spinner loading-xs"></span>
                ) : (
                  t("submit_appeal")
                )}
              </button>
              <button
                className="btn btn-ghost"
                onClick={() =>
                  setAppealModal({
                    open: false,
                    violationId: "",
                    violationReason: "",
                  })
                }
                disabled={isAppealing}
              >
                {t("cancel")}
              </button>
            </div>
            {appealError && (
              <p className="text-error text-sm mt-2">{appealError}</p>
            )}
          </div>
          <form method="dialog" className="modal-backdrop">
            <button
              onClick={() =>
                setAppealModal({
                  open: false,
                  violationId: "",
                  violationReason: "",
                })
              }
              disabled={isAppealing}
            >
              close
            </button>
          </form>
        </dialog>
      )}
    </div>
  );
}
