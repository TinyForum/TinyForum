// 违规记录类型定义
export interface ViolationRecord {
  id: string;
  reason: string; // 违规原因
  date: string; // 违规时间
  status: "pending" | "appealing" | "resolved" | "rejected";
  punishment: string; // 处罚措施（如“禁言3天”、“警告”）
  appealDeadline?: string; // 申诉截止时间
}
