-- 创建 moderator_applications 表（版主申请记录）
-- 2026-10-13

CREATE TABLE IF NOT EXISTS moderator_applications (
    id                     BIGSERIAL PRIMARY KEY,
    created_at             TIMESTAMPTZ,
    updated_at             TIMESTAMPTZ,
    deleted_at             TIMESTAMPTZ,
    user_id                BIGINT NOT NULL,
    board_id               BIGINT NOT NULL,
    reason                 VARCHAR(500),
    requested_permissions  JSON,
    status                 VARCHAR(20) NOT NULL DEFAULT 'pending',
    reviewer_id            BIGINT,
    review_note            VARCHAR(500),
    reviewed_at            TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_moderator_applications_user_id ON moderator_applications (user_id);
CREATE INDEX IF NOT EXISTS idx_moderator_applications_board_id ON moderator_applications (board_id);
CREATE INDEX IF NOT EXISTS idx_moderator_applications_status ON moderator_applications (status);
CREATE INDEX IF NOT EXISTS idx_moderator_applications_reviewer_id ON moderator_applications (reviewer_id);