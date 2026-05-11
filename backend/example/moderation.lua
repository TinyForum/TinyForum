-- --[[
--   内容审核机器人
--   触发器：on_keyword（关键词触发）
--   权限：manage:content, send:message, read:user

--   event 结构：
--     event.post_id      -- 触发帖子 ID（若为帖子）
--     event.comment_id   -- 触发评论 ID（若为评论）
--     event.user_id      -- 作者 ID
--     event.content_type -- "post" | "comment"
--     event.matched_kw   -- 匹配到的关键词
-- --]]

-- -- 违规级别配置
-- local LEVEL = {
--     WARNING = 1,  -- 警告
--     HIDE    = 2,  -- 隐藏内容
--     DELETE  = 3,  -- 删除 + 警告
--     BAN     = 4,  -- 封禁
-- }

-- -- 从 config 读取关键词映射
-- -- config.rules = [{ keywords: ["xxx"], level: 2 }]
-- local function get_level(keyword)
--     local rules = config.rules or {}
--     for _, rule in ipairs(rules) do
--         for _, kw in ipairs(rule.keywords or {}) do
--             if util.lower(kw) == util.lower(keyword) then
--                 return rule.level or LEVEL.WARNING
--             end
--         end
--     end
--     return LEVEL.WARNING
-- end

-- function main()
--     local user_id      = event.user_id
--     local content_type = event.content_type or "post"
--     local matched_kw   = event.matched_kw or ""
--     local post_id      = event.post_id
--     local comment_id   = event.comment_id

--     logf("审核触发: 用户=%d 类型=%s 关键词=%s", user_id, content_type, matched_kw)

--     local level = get_level(matched_kw)
--     logf("违规级别: %d", level)

--     -- 获取用户信息
--     local user, uerr = forum.getUser(user_id)
--     if uerr then
--         log("获取用户失败: " .. uerr)
--         return { success = false }
--     end

--     -- 管理员/版主免疫
--     if user.role == "admin" or user.role == "moderator" then
--         log("用户为管理员，跳过审核")
--         return { success = true, skipped = true }
--     end

--     local action_taken = {}

--     if level >= LEVEL.HIDE then
--         -- 隐藏内容
--         if content_type == "post" and post_id then
--             local ok, err = forum.moderatePost(post_id, "hide", "包含违禁词: " .. matched_kw)
--             if ok then
--                 table.insert(action_taken, "post_hidden")
--             else
--                 log("隐藏帖子失败: " .. (err or ""))
--             end
--         elseif content_type == "comment" and comment_id then
--             local ok, err = forum.deleteComment(comment_id)
--             if ok then
--                 table.insert(action_taken, "comment_deleted")
--             else
--                 log("删除评论失败: " .. (err or ""))
--             end
--         end
--     end

--     if level >= LEVEL.DELETE and content_type == "post" and post_id then
--         local ok, err = forum.deletePost(post_id)
--         if ok then
--             table.insert(action_taken, "post_deleted")
--         end
--     end

--     -- 发送警告私信
--     if level >= LEVEL.WARNING then
--         local warn_msg = string.format(
--             "您好 %s，您的内容因包含违禁词「%s」已被处理。\n请遵守社区规范，感谢配合。",
--             user.username, matched_kw
--         )
--         forum.sendMessage(user_id, warn_msg)
--         table.insert(action_taken, "warned")
--     end

--     if level >= LEVEL.BAN then
--         local ban_dur = config.ban_duration_sec or 86400
--         local ok, err = forum.banUser(user_id, "发布违禁内容: " .. matched_kw, ban_dur)
--         if ok then
--             table.insert(action_taken, "banned")
--         else
--             log("封禁失败: " .. (err or ""))
--         end
--     end

--     log("审核完成，执行操作: " .. table.concat(action_taken, ", "))
--     return { success = true, actions = action_taken, level = level }
-- end