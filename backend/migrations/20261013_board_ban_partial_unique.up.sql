-- 修复 board_bans 唯一索引：改为部分唯一索引（仅约束未软删记录）
-- 根因：UnbanUser 使用 GORM 软删除（deleted_at 非空），旧唯一索引不区分软删记录，
-- 导致"解封后再次封禁"同一 (user_id, board_id) 时触发唯一约束冲突。
-- 2026-10-13

-- 先删除旧的全量唯一索引
DROP INDEX IF EXISTS idx_active_ban;

-- 重建为部分唯一索引：仅当 deleted_at IS NULL（未软删）时保证唯一
CREATE UNIQUE INDEX idx_active_ban
    ON board_bans (user_id, board_id)
    WHERE deleted_at IS NULL;