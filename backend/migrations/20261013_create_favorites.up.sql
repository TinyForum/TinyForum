-- 创建 favorite_groups 与 favorites 表
-- favorites 唯一索引为部分索引（deleted_at IS NULL），保证"收藏 → 取消(软删) → 再次收藏"可行
-- 2026-10-13

CREATE TABLE IF NOT EXISTS favorite_groups (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id    BIGINT NOT NULL,
    name       VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    is_private BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_user_id ON favorite_groups (user_id);
CREATE INDEX IF NOT EXISTS idx_user_default ON favorite_groups (user_id, is_default);
CREATE INDEX IF NOT EXISTS idx_favorite_groups_deleted_at ON favorite_groups (deleted_at);

CREATE TABLE IF NOT EXISTS favorites (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    user_id     BIGINT NOT NULL,
    target_id   BIGINT NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    group_id    BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_fav_unique
    ON favorites (user_id, target_id, target_type)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_favorites_group_id ON favorites (group_id);
CREATE INDEX IF NOT EXISTS idx_favorites_deleted_at ON favorites (deleted_at);
