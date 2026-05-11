
-- --[[
--   每日统计播报机器人
--   触发器：on_schedule  cron: "0 9 * * *"（每天 9:00）
--   权限：read:stats, write:posts

--   config：
--     report_section_id  -- 发布统计帖的板块 ID
--     mention_top_n      -- 展示排名前 N 个（默认 5）
-- --]]

-- function main()
--     log("开始生成每日统计报告")

--     -- 获取论坛统计
--     local stats, err = forum.getStats()
--     if err then
--         log("获取统计失败: " .. err)
--         return { success = false, error = err }
--     end

--     local today = util.format_time(util.now(), "2006-01-02")
--     local top_n = config.mention_top_n or 5
--     local section_id = config.report_section_id or 1

--     -- 构建报告内容
--     local lines = {
--         string.format("# 📊 论坛日报 · %s", today),
--         "",
--         "## 今日数据一览",
--         "",
--         string.format("| 指标       | 数值   |"),
--         string.format("|------------|--------|"),
--         string.format("| 总帖子数   | %d    |", stats.post_count),
--         string.format("| 总用户数   | %d    |", stats.user_count),
--         string.format("| 总评论数   | %d    |", stats.comment_count),
--         string.format("| 今日活跃   | %d    |", stats.active_today),
--         "",
--         string.format("_本报告由系统机器人自动生成于 %s_",
--             util.format_time(util.now(), "2006-01-02 15:04:05")),
--     }

--     local content = table.concat(lines, "\n")

--     -- 发布统计帖
--     local post, perr = forum.createPost(
--         string.format("论坛日报 · %s", today),
--         content,
--         section_id
--     )

--     if perr then
--         log("发布统计帖失败: " .. perr)
--         return { success = false, error = perr }
--     end

--     logf("统计帖发布成功，帖子 ID: %d", post.id)
--     return {
--         success    = true,
--         post_id    = post.id,
--         stats      = stats,
--     }
-- end