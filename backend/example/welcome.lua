-- --[[
--   欢迎新用户机器人
--   触发器：on_user_register（新用户注册事件）
--   权限：send:message, read:user

--   event 结构：
--     event.user_id     -- 新注册用户 ID
--     event.username    -- 用户名
--     event.email       -- 邮箱
--     event.created_at  -- 注册时间（Unix）
-- --]]

-- function main()
--     local uid = event.user_id
--     local username = event.username or "新用户"

--     log("欢迎机器人触发，新用户: " .. username .. " (id=" .. tostring(uid) .. ")")

--     -- 读取配置
--     local welcome_msg = config.welcome_message or "欢迎加入论坛！"
--     local guide_post_id = config.guide_post_id or 0

--     -- 构建欢迎语
--     local msg = string.format(
--         "亲爱的 %s，%s\n\n" ..
--         "注册时间：%s\n\n" ..
--         "如有问题，请查看新手指引。",
--         username,
--         welcome_msg,
--         util.format_time(event.created_at or util.now(), "2006-01-02 15:04")
--     )

--     -- 发送私信
--     local ok, err = forum.sendMessage(uid, msg)
--     if not ok then
--         log("发送欢迎私信失败: " .. (err or "unknown"))
--         return { success = false, error = err }
--     end

--     -- 如果配置了新手指引帖子，回复该帖
--     if guide_post_id > 0 then
--         local reply_content = string.format("@%s 欢迎来到论坛！🎉", username)
--         local comment, cerr = forum.replyPost(guide_post_id, reply_content)
--         if cerr then
--             log("回复新手指引帖失败: " .. cerr)
--         end
--     end

--     log("欢迎流程完成")
--     return { success = true, user = username }
-- end