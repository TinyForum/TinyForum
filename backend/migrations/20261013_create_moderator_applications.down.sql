-- 回滚：删除 moderator_applications 表（注意：会同时删除申请记录，请谨慎）
-- 2026-10-13

DROP INDEX IF EXISTS idx_moderator_applications_reviewer_id;
DROP INDEX IF EXISTS idx_moderator_applications_status;
DROP INDEX IF EXISTS idx_moderator_applications_board_id;
DROP INDEX IF EXISTS idx_moderator_applications_user_id;
DROP TABLE IF EXISTS moderator_applications;