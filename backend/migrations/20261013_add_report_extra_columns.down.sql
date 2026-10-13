-- 回滚：删除 reports 表新增的举报扩展字段（先备份再删除，避免不可逆数据丢失）
-- 2026-10-13

CREATE TABLE IF NOT EXISTS reports_extra_backup AS
SELECT id, content_snapshot, reporter_ip, is_anonymous, priority
FROM reports;

ALTER TABLE reports
    DROP COLUMN IF EXISTS content_snapshot,
    DROP COLUMN IF EXISTS reporter_ip,
    DROP COLUMN IF EXISTS is_anonymous,
    DROP COLUMN IF EXISTS priority;