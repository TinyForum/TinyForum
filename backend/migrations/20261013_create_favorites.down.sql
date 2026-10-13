-- 回滚：删除 favorites 与 favorite_groups 表（注意：会同时删除收藏历史，请谨慎）
-- 2026-10-13

DROP INDEX IF EXISTS idx_favorites_deleted_at;
DROP INDEX IF EXISTS idx_favorites_group_id;
DROP INDEX IF EXISTS idx_fav_unique;
DROP TABLE IF EXISTS favorites;

DROP INDEX IF EXISTS idx_favorite_groups_deleted_at;
DROP INDEX IF EXISTS idx_user_default;
DROP INDEX IF EXISTS idx_user_id;
DROP TABLE IF EXISTS favorite_groups;
