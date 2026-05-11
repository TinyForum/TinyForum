--[[
  欢迎新用户机器人
  ─────────────────────────────────
  触发器：on_user_register
  权限：send:message, read:user

  event 字段：
    event.user_id     -- 新用户 ID (number)
    event.username    -- 用户名 (string)
    event.created_at  -- 注册时间 Unix (number)

  config 字段（在 ConfigValues 中配置）：
    config.welcome_message  -- 欢迎语（可选）
    config.guide_post_id    -- 新手指引帖子 ID（可选，>0 时自动回复）
--]]

function main()
    local uid      = event.user_id
    local username = event.username or "新朋友"

    log("欢迎机器人触发，用户: " .. username .. " (id=" .. tostring(uid) .. ")")

    local welcome = config.welcome_message or "欢迎加入我们的社区！🎉"
    local reg_time = util.format_time(event.created_at or util.now(), "2006-01-02 15:04")

    local msg = string.format(
        "亲爱的 %s，\n\n%s\n\n您的注册时间：%s\n\n如有问题，欢迎在新手专区提问。",
        username, welcome, reg_time
    )

    -- 发送欢迎私信
    local ok, err = forum.sendMessage(uid, msg)
    if not ok then
        log("发送私信失败: " .. (err or "unknown"))
        return { success = false, error = err }
    end

    -- 在新手指引帖留言（可选）
    local guide_id = config.guide_post_id or 0
    if guide_id > 0 then
        local reply = string.format("👋 欢迎 @%s 加入！", username)
        local _, cerr = forum.replyPost(guide_id, reply)
        if cerr then
            log("回复指引帖失败: " .. cerr)
        end
    end

    log("欢迎流程完成")
    return { success = true, user = username }
end