--[[
  内容审核机器人
  ─────────────────────────────────
  触发器：on_keyword
  权限：manage:content, send:message, read:user

  event 字段：
    event.post_id      -- 帖子 ID（若为帖子触发）
    event.comment_id   -- 评论 ID（若为评论触发）
    event.user_id      -- 作者 ID
    event.content_type -- "post" | "comment"
    event.matched_kw   -- 命中的关键词

  config 字段：
    config.rules = [{ keywords: ["xxx"], level: 2 }]
      level: 1=警告, 2=隐藏+警告, 3=删除+警告, 4=封禁
    config.ban_duration_sec -- 封禁时长（秒），默认 86400
--]]

local LEVEL = { WARNING = 1, HIDE = 2, DELETE = 3, BAN = 4 }

local function get_level(keyword)
    for _, rule in ipairs(config.rules or {}) do
        for _, kw in ipairs(rule.keywords or {}) do
            if util.lower(kw) == util.lower(keyword) then
                return rule.level or LEVEL.WARNING
            end
        end
    end
    return LEVEL.WARNING
end

function main()
    local uid          = event.user_id
    local content_type = event.content_type or "post"
    local matched_kw   = event.matched_kw   or ""
    local post_id      = event.post_id
    local comment_id   = event.comment_id

    logf("审核触发: user=%d type=%s kw=%s", uid, content_type, matched_kw)

    -- 获取用户信息，管理员/版主豁免
    local user, uerr = forum.getUser(uid)
    if uerr then
        log("获取用户失败: " .. uerr)
        return { success = false }
    end
    if user.role == "admin" or user.role == "moderator" then
        log("管理员豁免，跳过审核")
        return { success = true, skipped = true }
    end

    local level   = get_level(matched_kw)
    local actions = {}

    -- 隐藏内容
    if level >= LEVEL.HIDE then
        if content_type == "post" and post_id then
            local ok, err = forum.moderatePost(post_id, "hide", "违禁词: " .. matched_kw)
            if ok then table.insert(actions, "post_hidden")
            else log("隐藏帖子失败: " .. (err or "")) end
        elseif content_type == "comment" and comment_id then
            local ok, err = forum.deleteComment(comment_id)
            if ok then table.insert(actions, "comment_deleted")
            else log("删除评论失败: " .. (err or "")) end
        end
    end

    -- 删除帖子
    if level >= LEVEL.DELETE and content_type == "post" and post_id then
        local ok, _ = forum.deletePost(post_id)
        if ok then table.insert(actions, "post_deleted") end
    end

    -- 发送警告私信
    if level >= LEVEL.WARNING then
        local warn = string.format(
            "您好 %s，您的内容因包含违禁词「%s」已被处理，请遵守社区规范。",
            user.username, matched_kw
        )
        forum.sendMessage(uid, warn)
        table.insert(actions, "warned")
    end

    -- 封禁
    if level >= LEVEL.BAN then
        local dur = config.ban_duration_sec or 86400
        local ok, err = forum.banUser(uid, "发布违禁内容: " .. matched_kw, dur)
        if ok then table.insert(actions, "banned")
        else log("封禁失败: " .. (err or "")) end
    end

    logf("审核完成，操作: %s", table.concat(actions, ", "))
    return { success = true, actions = actions, level = level }
end