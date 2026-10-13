-- 为 creations 表新增 locked 字段，用于锁定帖子，锁定后禁止编辑/回复
-- 2026-10-12

ALTER TABLE creations
    ADD COLUMN locked boolean NOT NULL DEFAULT false;