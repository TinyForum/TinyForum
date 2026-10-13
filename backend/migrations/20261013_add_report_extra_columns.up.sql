-- reports 表新增举报扩展字段（被举报内容快照、举报人IP、是否匿名、优先级）
-- 2026-10-13

ALTER TABLE reports
    ADD COLUMN IF NOT EXISTS content_snapshot TEXT,
    ADD COLUMN IF NOT EXISTS reporter_ip     VARCHAR(45),
    ADD COLUMN IF NOT EXISTS is_anonymous    BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS priority        SMALLINT NOT NULL DEFAULT 2;