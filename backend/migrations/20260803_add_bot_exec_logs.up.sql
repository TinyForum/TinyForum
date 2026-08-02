-- 为 bots 表添加执行日志列
ALTER TABLE bots ADD COLUMN IF NOT EXISTS last_exec_logs JSONB DEFAULT '[]';
ALTER TABLE bots ADD COLUMN IF NOT EXISTS last_exec_duration_ms BIGINT DEFAULT 0;
