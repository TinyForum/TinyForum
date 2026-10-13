-- 创建 board_bans 表（板块封禁记录）
-- 唯一索引为部分索引（deleted_at IS NULL），保证"封禁 → 解封(软删) → 再次封禁"可行
-- 2026-10-13

CREATE TABLE IF NOT EXISTS board_bans (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id    BIGINT NOT NULL,
    board_id   BIGINT NOT NULL,
    banned_by  BIGINT NOT NULL,
    reason     TEXT,
    expires_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_active_ban
    ON board_bans (user_id, board_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_board_bans_deleted_at ON board_bans (deleted_at);
CREATE INDEX IF NOT EXISTS idx_board_bans_expires_at ON board_bans (expires_at);