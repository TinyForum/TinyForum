-- 回滚：恢复 favorites 全量唯一索引（注意：若存在因部分索引而新增的软删重复记录，恢复前需先清理）
-- 2026-10-13

DROP INDEX IF EXISTS idx_fav_unique;

CREATE UNIQUE INDEX idx_fav_unique
    ON favorites (user_id, target_id, target_type);
