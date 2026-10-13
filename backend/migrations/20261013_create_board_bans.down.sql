-- 回滚：删除 board_bans 表（注意：会同时删除封禁历史，请谨慎）
-- 2026-10-13

DROP INDEX IF EXISTS idx_board_bans_expires_at;
DROP INDEX IF EXISTS idx_board_bans_deleted_at;
DROP INDEX IF EXISTS idx_active_ban;
DROP TABLE IF EXISTS board_bans;