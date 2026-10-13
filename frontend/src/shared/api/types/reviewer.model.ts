export interface AuditTaskVO {
  id: number;
  created_at: string;
  target_type: string; // post | comment | user
  target_id: number;
  trigger_type: string; // sensitive_word | report_aggregate | manual
  trigger_meta: string;
  status: string; // pending | approved | rejected
  // 审核历史字段
  reviewer_id?: number;
  reviewer_name?: string;
  review_note?: string;
  reviewed_at?: string;
}
