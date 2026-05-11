--[[
  每日统计播报机器人
  ─────────────────────────────────
  触发器：on_schedule  cron: "0 9 * * *"（每天 9:00）
  权限：read:stats, write:posts

  config 字段：
    config.report_section_id  -- 发布报告的板块 ID（默认 1）
--]]

function main()
    log("开始生成每日统计报告")

    local stats, err = forum.getStats()
    if err then
        log("获取统计失败: " .. err)
        return { success = false, error = err }
    end

    local today      = util.format_time(util.now(), "2006-01-02")
    local section_id = config.report_section_id or 1

    local lines = {
        "# 📊 论坛日报 · " .. today,
        "",
        "| 指标 | 数值 |",
        "|------|------|",
        string.format("| 总帖子数 | %d |", stats.post_count),
        string.format("| 总用户数 | %d |", stats.user_count),
        string.format("| 总评论数 | %d |", stats.comment_count),
        string.format("| 今日活跃 | %d |", stats.active_today),
        "",
        "_本报告由系统机器人自动生成于 " ..
            util.format_time(util.now(), "2006-01-02 15:04:05") .. "_",
    }

    local post, perr = forum.createPost(
        "论坛日报 · " .. today,
        table.concat(lines, "\n"),
        section_id
    )
    if perr then
        log("发布报告失败: " .. perr)
        return { success = false, error = perr }
    end

    logf("报告发布成功，帖子 ID: %d", post.id)
    return { success = true, post_id = post.id, stats = stats }
end