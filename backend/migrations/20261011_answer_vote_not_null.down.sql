-- 回答投票：回滚 vote_type 非空约束
-- 2026-10-11

ALTER TABLE answer_votes
    ALTER COLUMN vote_type DROP NOT NULL;