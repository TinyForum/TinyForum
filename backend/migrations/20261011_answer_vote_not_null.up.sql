-- 回答投票：将 vote_type 设为 NOT NULL
-- 根因：代码中直接解引用 *VoteType，历史脏数据（NULL）会导致 panic。
-- 先回填存量 NULL 行（默认视为 up），再收紧约束。
-- 2026-10-11

UPDATE answer_votes SET vote_type = 'up' WHERE vote_type IS NULL;

ALTER TABLE answer_votes
    ALTER COLUMN vote_type SET NOT NULL;