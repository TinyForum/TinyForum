# 清空所有表
truncate_tables() {
    if $EXECUTE; then
        echo "Clearing existing data..."
        run_sql "TRUNCATE TABLE 
            users, boards, tags, topics, posts, post_tags, comments, questions, 
            likes, follows, notifications, topic_posts, topic_follows, announcements,
            reports, sign_ins, moderators, moderator_applications, audit_logs,
            refresh_tokens, attachments, plugins, bots, blocked_ips, ip_risk_records,
            user_risk_records, content_audit_tasks, timeline_events, timeline_subscriptions,
            favorites, violations, votes, answer_votes, casbin_rule
        RESTART IDENTITY CASCADE;"
    fi
}
truncate_tables