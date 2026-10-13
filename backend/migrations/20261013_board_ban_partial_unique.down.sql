-- 回滚：恢复 board_bans 全量唯一索引（注意：若存在因部分索引而新增的软删重复记录，恢复全量唯一索引前需先清理）
-- 2026-10-13

DROP INDEX IF EXISTS idx_active_ban;

CREATE UNIQUE INDEX idx_active_ban
    ON board_bans (user_id, board_id);