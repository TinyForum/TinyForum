-- KEYS[1]   : Redis key，格式 "rl:{userID}:{action}"
-- ARGV[1]   : limit (配额上限)
-- ARGV[2]   : window_ms (时间窗口长度，毫秒)
-- ARGV[3]   : now_ms (当前毫秒时间戳)
-- ARGV[4]   : member (唯一成员标识，使用纳秒时间戳)
-- 返回值    : { allowed_flag, current, limit, ttl_ms }
--            allowed_flag: 1 表示允许，0 表示拒绝
--            current: 当前窗口内请求数（允许时已包含本次）
--            limit: 配额上限
--            ttl_ms: 距离窗口重置的毫秒数（仅当拒绝时有效，允许时返回0）
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window_ms = tonumber(ARGV[2])
local now_ms = tonumber(ARGV[3])
local member = ARGV[4]

redis.log(redis.LOG_NOTICE, "ratelimit script executed, key=", KEYS[1])
-- 1. 移除窗口外的旧记录
local window_start = now_ms - window_ms
redis.call('ZREMRANGEBYSCORE', key, 0, window_start)

-- 2. 统计当前窗口内的请求数
local current = redis.call('ZCARD', key)

-- 3. 判断是否允许
if current < limit then
    -- 允许：插入本次请求记录
    redis.call('ZADD', key, now_ms, member)
    -- 设置过期时间（窗口长度 + 1分钟缓冲）
    redis.call('EXPIRE', key, (window_ms / 1000) + 60)
    return {1, current + 1, limit, 0}
else
    -- 拒绝：计算还需要多久窗口重置（最早一条记录的时间 + 窗口长度 - 当前时间）
    -- 获取最早记录的 score（即最小时间戳）
    local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
    local reset_ms = 0
    if #oldest >= 2 then
        local earliest_score = tonumber(oldest[2])
        reset_ms = (earliest_score + window_ms) - now_ms
        if reset_ms < 0 then reset_ms = 0 end
    else
        reset_ms = window_ms  -- 降级：如果没有记录（理论上不会发生），使用整个窗口
    end
    return {0, current, limit, reset_ms}
end