-- 新增板块创建申请表 board_creation_applications
-- 用于用户申请创建新板块，管理员审批通过后自动创建对应板块
-- 2026-10-11

CREATE TABLE IF NOT EXISTS board_creation_applications (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    user_id          BIGINT       NOT NULL,
    name             VARCHAR(50)  NOT NULL,
    slug             VARCHAR(50)  NOT NULL,
    description      VARCHAR(500),
    reason           VARCHAR(500),
    status           VARCHAR(20)  DEFAULT 'pending',
    reviewer_id      BIGINT,
    review_note      VARCHAR(500),
    reviewed_at      TIMESTAMPTZ,
    created_board_id BIGINT
);

CREATE INDEX IF NOT EXISTS idx_board_creation_applications_user_id ON board_creation_applications (user_id);
-- 部分唯一索引：同一 slug 至多一条待审核（未删除）申请，从数据库层杜绝并发重复申请
CREATE UNIQUE INDEX IF NOT EXISTS uq_board_creation_applications_pending_slug
    ON board_creation_applications (slug)
    WHERE status = 'pending' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_board_creation_applications_status ON board_creation_applications (status);
CREATE INDEX IF NOT EXISTS idx_board_creation_applications_deleted_at ON board_creation_applications (deleted_at);
