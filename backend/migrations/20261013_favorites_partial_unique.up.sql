-- 修复 favorites 唯一索引：改为部分唯一索引（仅约束未软删记录）
-- 根因：取消收藏使用 GORM 软删除（deleted_at 非空），旧全量唯一索引不区分软删记录，
-- 导致"取消收藏后再次收藏"同一 (user_id, target_id, target_type) 时触发唯一约束冲突。
-- 2026-10-13

DROP INDEX IF EXISTS idx_fav_unique;

CREATE UNIQUE INDEX idx_fav_unique
    ON favorites (user_id, target_id, target_type)
    WHERE deleted_at IS NULL;
