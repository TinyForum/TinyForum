-- 回滚：删除 creations.locked 字段前，先将数据备份到临时表，避免不可逆的数据丢失
-- 2026-10-12

CREATE TABLE IF NOT EXISTS creations_locked_backup (
    id     BIGINT PRIMARY KEY,
    locked boolean NOT NULL DEFAULT false
);

INSERT INTO creations_locked_backup (id, locked)
SELECT id, locked FROM creations;

ALTER TABLE creations
    DROP COLUMN locked;