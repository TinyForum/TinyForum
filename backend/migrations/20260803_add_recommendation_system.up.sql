-- 推荐系统核心表（up）
-- 2026-08-03

-- 1. 用户行为事件表：记录每一次用户与内容的交互
CREATE TABLE IF NOT EXISTS user_behavior_events (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    user_id         INTEGER NOT NULL,
    target_id       INTEGER NOT NULL,
    target_type     VARCHAR(30)  NOT NULL DEFAULT 'creation',
    behavior_type   VARCHAR(30)  NOT NULL,
    value           REAL         NOT NULL DEFAULT 1.0,
    session_id      VARCHAR(64)  NOT NULL DEFAULT '',
    context_json    TEXT         NOT NULL DEFAULT '{}',
    created_ts      BIGINT       NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_ube_user_id       ON user_behavior_events (user_id);
CREATE INDEX IF NOT EXISTS idx_ube_target        ON user_behavior_events (target_id, target_type);
CREATE INDEX IF NOT EXISTS idx_ube_user_behavior ON user_behavior_events (user_id, behavior_type, created_ts);
CREATE INDEX IF NOT EXISTS idx_ube_created_ts    ON user_behavior_events (created_ts);
CREATE INDEX IF NOT EXISTS idx_ube_deleted_at    ON user_behavior_events (deleted_at);

-- 2. 内容特征表：缓存内容的特征向量与元数据
CREATE TABLE IF NOT EXISTS content_features (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    creation_id     INTEGER      NOT NULL UNIQUE,
    tag_ids_json    TEXT         NOT NULL DEFAULT '[]',
    board_id        INTEGER      NOT NULL DEFAULT 0,
    author_id       INTEGER      NOT NULL DEFAULT 0,
    quality_score   REAL         NOT NULL DEFAULT 0,
    hot_score       REAL         NOT NULL DEFAULT 0,
    freshness_score REAL         NOT NULL DEFAULT 1.0,
    feature_vector  TEXT         NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_cf_creation_id ON content_features (creation_id);
CREATE INDEX IF NOT EXISTS idx_cf_hot_score   ON content_features (hot_score DESC);
CREATE INDEX IF NOT EXISTS idx_cf_deleted_at  ON content_features (deleted_at);

-- 3. 推荐反馈表：记录用户看到推荐内容后的行为反馈
CREATE TABLE IF NOT EXISTS recommendation_feedbacks (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    user_id         INTEGER      NOT NULL,
    creation_id     INTEGER      NOT NULL,
    feedback_type   VARCHAR(20)  NOT NULL DEFAULT 'impression',
    source_type     VARCHAR(30)  NOT NULL DEFAULT 'recommend',
    position        INTEGER      NOT NULL DEFAULT 0,
    session_id      VARCHAR(64)  NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_rf_user_id      ON recommendation_feedbacks (user_id);
CREATE INDEX IF NOT EXISTS idx_rf_creation_id  ON recommendation_feedbacks (creation_id);
CREATE INDEX IF NOT EXISTS idx_rf_feedback     ON recommendation_feedbacks (user_id, feedback_type, created_at);
CREATE INDEX IF NOT EXISTS idx_rf_deleted_at   ON recommendation_feedbacks (deleted_at);

-- 4. 用户兴趣画像表：缓存用户近期兴趣向量
CREATE TABLE IF NOT EXISTS user_interest_profiles (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    user_id         INTEGER NOT NULL UNIQUE,
    tag_weights_json TEXT    NOT NULL DEFAULT '{}',
    board_weights_json TEXT  NOT NULL DEFAULT '{}',
    active_tags_json  TEXT   NOT NULL DEFAULT '[]',
    last_updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_uip_user_id      ON user_interest_profiles (user_id);
CREATE INDEX IF NOT EXISTS idx_uip_deleted_at   ON user_interest_profiles (deleted_at);
