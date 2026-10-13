-- 创建 violations 表（用户违规记录）
-- 2026-10-13

CREATE TABLE IF NOT EXISTS violations (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    user_id          BIGINT NOT NULL,
    violation_type   VARCHAR(30) NOT NULL,
    reason           VARCHAR(255) NOT NULL,
    content_snapshot TEXT,
    evidence_url     VARCHAR(512),
    source           VARCHAR(20) NOT NULL DEFAULT 'system',
    status           VARCHAR(20) NOT NULL DEFAULT 'pending',
    operator_id      BIGINT NOT NULL DEFAULT 0,
    punish_type      VARCHAR(30) NOT NULL DEFAULT '',
    punish_expire_at TIMESTAMPTZ,
    appeal_status    VARCHAR(20) NOT NULL DEFAULT 'none',
    appeal_reason    VARCHAR(500),
    appeal_time      TIMESTAMPTZ,
    appeal_result    VARCHAR(500)
);

CREATE INDEX IF NOT EXISTS idx_user_violation ON violations (user_id);
CREATE INDEX IF NOT EXISTS idx_violations_violation_type ON violations (violation_type);
CREATE INDEX IF NOT EXISTS idx_violations_status ON violations (status);
CREATE INDEX IF NOT EXISTS idx_violations_deleted_at ON violations (deleted_at);