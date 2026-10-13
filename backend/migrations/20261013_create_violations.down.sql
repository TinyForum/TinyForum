-- 回滚：删除 violations 表（注意：会同时删除违规历史，请谨慎）
-- 2026-10-13

DROP INDEX IF EXISTS idx_violations_deleted_at;
DROP INDEX IF EXISTS idx_violations_status;
DROP INDEX IF EXISTS idx_violations_violation_type;
DROP INDEX IF EXISTS idx_user_violation;
DROP TABLE IF EXISTS violations;