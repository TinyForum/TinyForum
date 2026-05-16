#!/usr/bin/env bash
# mock_data.sh - PostgreSQL Mock 数据生成脚本
# 用法: ./mock_data.sh [选项]
# 示例: USERS=100 POSTS=500 ./mock_data.sh --execute
# USERS=200 POSTS=1000 COMMENTS=3000 ./mock_data.sh --execute
# 查看所有用户 SELECT rolname FROM pg_roles;
set -euo pipefail

# ==================== 默认配置（可通过环境变量覆盖）====================
: "${DB_HOST:=localhost}"
: "${DB_PORT:=5432}"
: "${DB_NAME:=tiny_forum}"
: "${DB_USER:=simons}"
: "${DB_PASS:=tf-password}"              # 空则使用 .pgpass 或 trust 认证

# 各表数据量（设为 0 则跳过）
: "${USERS:=500}"
: "${BOARDS:=100}"
: "${TAGS:=200}"
: "${POSTS:=2000}"
: "${COMMENTS:=5000}"
: "${LIKES:=3000}"
: "${FOLLOWS:=1000}"
: "${NOTIFICATIONS:=1000}"
: "${TOPICS:=150}"
: "${TOPIC_POSTS:=500}"
: "${TOPIC_FOLLOWS:=800}"
: "${ANNOUNCEMENTS:=100}"
: "${QUESTIONS:=500}"
: "${REPORTS:=200}"
: "${SIGN_INS:=2000}"
: "${MODERATORS:=150}"
: "${MODERATOR_APPLICATIONS:=200}"
: "${AUDIT_LOGS:=500}"
: "${REFRESH_TOKENS:=1000}"
: "${ATTACHMENTS:=500}"
: "${PLUGINS:=100}"
: "${BOTS:=100}"
: "${BLOCKED_IPS:=100}"
: "${IP_RISK_RECORDS:=200}"
: "${USER_RISK_RECORDS:=200}"
: "${CONTENT_AUDIT_TASKS:=300}"
: "${TIMELINE_EVENTS:=1000}"
: "${TIMELINE_SUBSCRIPTIONS:=800}"
: "${FAVORITES:=1000}"
: "${VIOLATIONS:=200}"
: "${VOTES:=1000}"
: "${ANSWER_VOTES:=500}"

# 执行模式：默认仅打印 SQL（dry-run），加 --execute 才真正执行
EXECUTE=false
if [[ "${1:-}" == "--execute" ]]; then
    EXECUTE=true
fi



# ==================== 工具函数 ====================

# 生成随机时间戳（过去 N 天内）
random_timestamp() {
    local days="${1:-365}"
    local offset=$((RANDOM % (days * 86400)))
    local ts=$(($(date +%s) - offset))
    if [[ "$(uname)" == "Darwin" ]]; then
        date -u -r "$ts" '+%Y-%m-%d %H:%M:%S'
    else
        date -u -d "@$ts" '+%Y-%m-%d %H:%M:%S'
    fi
}

# 生成随机 UUID
random_uuid() {
    if command -v uuidgen >/dev/null 2>&1; then
        uuidgen | tr '[:upper:]' '[:lower:]'
    elif [[ -f /proc/sys/kernel/random/uuid ]]; then
        cat /proc/sys/kernel/random/uuid
    else
        # fallback: 生成简单的随机十六进制
        openssl rand -hex 16 | sed 's/\(..\)/\1-/g; s/-$//'
    fi
}

# 随机选择数组元素
random_pick() {
    local arr=("$@")
    echo "${arr[$RANDOM % ${#arr[@]}]}"
}

# 随机整数 [min, max]
random_int() {
    local min=$1
    local max=$2
    echo $((min + RANDOM % (max - min + 1)))
}

# 随机 IP
random_ip() {
    echo "$((RANDOM % 256)).$((RANDOM % 256)).$((RANDOM % 256)).$((RANDOM % 256))"
}

# 生成随机字符串
random_string() {
    local len=$1
    openssl rand -hex $((len / 2 + 1)) | head -c "$len"
}

# 生成中文标题（模拟）
random_title() {
    local titles=(
        # ---------- 按性别 & 日常兴趣 ----------
        "女生必看的3个护肤小技巧"          # 女性常见话题
        "男生穿搭避坑指南"                  # 男性常见话题
        "给闺蜜的一封信"                    # 女性情感
        "兄弟深夜撸串局"                    # 男性社交

        # ---------- 按年龄阶段 ----------
        "20岁，花时间搞懂这5件事"           # 年轻人成长
        "三十而已：30岁后我学会的松弛感"     # 30+ 感悟
        "四十不惑，开始享受慢生活"           # 40+ 生活
        "退休两年，我的旅行笔记"             # 老年生活

        # ---------- 按职业人群 ----------
        "程序员的桌面装备分享"               # IT
        "教师日常：那些让你又气又笑的瞬间"   # 教育
        "医生的手机里存了哪些App"            # 医疗
        "新手会计入职第一周体验"             # 财务
        "设计师的灵感网站清单"               # 设计
        "律师：普通人最该留意的3个法律风险"  # 法律
        "自由职业者如何安排一天"             # 自由职业

        # ---------- 通用生活 & 情感 ----------
        "今日份小确幸☀️"                     # 所有人
        "记录一下最近的胡思乱想"             # 情感
        "我妈说：家里干净，运气就好"          # 家庭
        "最近循环播放的一首歌🎵"             # 娱乐

        # ---------- 健康 & 生活方式 ----------
        "坚持早起30天，变化有多大？"         # 自律
        "减脂期外卖怎么点？附点单模板"       # 健康
        "周末给自己做一顿brunch"             # 精致生活

        # ---------- 学习 & 自我提升 ----------
        "碎片时间学英语，我用这几个方法"     # 学习
        "下班后2小时，我考下了证书"          # 技能提升
    )
    # 随机选一个标题，后面可加随机后缀（保留原函数风格）
    echo "${titles[$RANDOM % ${#titles[@]}]} $(random_string 4)"
}

# 生成中文内容
random_content() {
    local contents=(
        # 生活日常
        "今天尝试了一道新菜，家人说很好吃，幸福感满满！🥘"
        "周末带孩子去公园放风筝，阳光正好，心情也变好了。☀️"
        "整理房间时翻到老照片，满满都是回忆呀～📸"
        "坚持晨跑一周，感觉精力更充沛了，继续加油！🏃‍♀️"
        # 职场与学习
        "工作中遇到难题，请教了同事后豁然开朗，团队合作真棒！💼"
        "终于考完试了，不管结果如何，努力过就不后悔。📚"
        "加班到很晚，但项目上线成功，一切值得。🚀"
        "参加行业分享会，学到了新思路，保持学习很重要。🎓"
        # 亲子与家庭
        "宝宝第一次叫妈妈，心都融化了！👶"
        "陪父母去医院体检，身体健康就是最大的福气。❤️"
        "和另一半一起做饭、追剧，平凡的日子也很甜。🍲"
        # 健康与养生
        "开始注意饮食搭配，少油少盐，身体感觉轻松多了。🥗"
        "老年人也要多活动，今天和老姐妹去跳广场舞了。💃"
        "瑜伽打卡第30天，柔韧性和心态都变好了。🧘"
        # 旅行与户外
        "山里的空气真好，远离城市喧嚣，治愈了疲惫。🏞️"
        "海边日出太美了，推荐大家去看看！🌊"
        # 兴趣与娱乐
        "学吉他两个月，终于能弹唱一首完整的歌了。🎸"
        "追完一部好剧，结局好感动，有没有同好？📺"
        # 实用分享
        "发现一个整理收纳的小技巧，空间瞬间变大～🧺"
        "手机这些设置可以省电又护眼，亲测有效。📱"
    )
    echo "${contents[$RANDOM % ${#contents[@]}]} $(random_string 8)"
}

# 生成随机创建时间
random_created_at() {
    local days="${1:-365}"
    local offset=$((RANDOM % (days * 86400)))
    local ts=$(($(date +%s) - offset))
    if [[ "$(uname)" == "Darwin" ]]; then
        date -u -r "$ts" '+%Y-%m-%d %H:%M:%S'
    else
        date -u -d "@$ts" '+%Y-%m-%d %H:%M:%S'
    fi
}
# 生成随机文章状态
random_status() {
    local statuses=("draft" "published" "deleted")
    echo "${statuses[$RANDOM % ${#statuses[@]}]}"
}

# 执行 SQL（dry-run 或实际执行）
run_sql() {
    local sql="$1"
    if $EXECUTE; then
        if [[ -n "$DB_PASS" ]]; then
            PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -q -c "$sql"
        else
            psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -q -c "$sql"
        fi
    else
        echo "$sql"
        echo ""
    fi
}

# 批量插入的 VALUES 缓冲
BATCH_SIZE=1000
VALUES_BUFFER=""
BUFFER_COUNT=0
BUFFER_TABLE=""

# 开始批量插入
begin_batch() {
    BUFFER_TABLE="$1"
    VALUES_BUFFER=""
    BUFFER_COUNT=0
}

# 添加一条记录到缓冲
add_batch() {
    local values="$1"
    if [[ $BUFFER_COUNT -gt 0 ]]; then
        VALUES_BUFFER+=", "
    fi
    VALUES_BUFFER+="$values"
    ((BUFFER_COUNT++))
    
    if [[ $BUFFER_COUNT -ge $BATCH_SIZE ]]; then
        flush_batch
    fi
}

# 刷新缓冲到数据库
flush_batch() {
    if [[ $BUFFER_COUNT -gt 0 ]]; then
        local sql="INSERT INTO $BUFFER_TABLE VALUES $VALUES_BUFFER ON CONFLICT DO NOTHING;"
        run_sql "$sql"
        VALUES_BUFFER=""
        BUFFER_COUNT=0
    fi
}

# ==================== 数据生成 ====================

echo "-- Mock Data Generation Script"
echo "-- Database: $DB_NAME"
echo "-- Execute mode: $EXECUTE"
echo ""

# 1. users（基础表，无依赖）

if [[ $USERS -gt 0 ]]; then
    echo "-- Generating $USERS users..."
    begin_batch "public.users ( created_at, updated_at, username, email, password, avatar, bio, role, score, is_active, is_blocked, last_login, invited_by_id)"
    
    for i in $(seq 1 $USERS); do
        username="user_${i}_$(random_string 4)"
        email="user${i}@example.com"
        password="\$2a\$10\$$(random_string 60)"  # bcrypt hash 占位
        avatar="https://api.dicebear.com/7.x/avataaars/svg?seed=${username}"
        bio="用户 ${i} 的个人简介"
        role=$(random_pick "user" "user" "user" "moderator" "admin")
        score=$(random_int 0 10000)
        is_active=$(random_pick "true" "true" "true" "false")
        is_blocked=$(random_pick "false" "false" "false" "true")
        last_login="'"$(random_timestamp 30)"'"
        invited_by_id="NULL"
        [[ $i -gt 1 && $(random_int 1 3) -eq 1 ]] && invited_by_id=$(random_int 1 $((i-1)))
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '$username', '$email', '$password', '$avatar', '$bio', '$role', $score, $is_active, $is_blocked, $last_login, $invited_by_id)"
    done
    flush_batch
fi

# 2. boards（基础表）
if [[ $BOARDS -gt 0 ]]; then
    echo "-- Generating $BOARDS boards..."
    begin_batch "public.boards ( created_at, updated_at, name, slug, description, icon, cover, parent_id, sort_order, view_role, post_role, reply_role)"
    
    for i in $(seq 1 $BOARDS); do
        name="板块_${i}_$(random_string 3)"
        slug="board-$(random_string 6)"
        description="这是第 ${i} 个板块的描述"
        parent_id="NULL"
        [[ $i -gt 3 && $(random_int 1 3) -eq 1 ]] && parent_id=$(random_int 1 $((i-1)))
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '$name', '$slug', '$description', 'icon-$i', 'cover-$i.jpg', $parent_id, $i, 'user', 'user', 'user')"
    done
    flush_batch
fi

# 3. tags（基础表）
if [[ $TAGS -gt 0 ]]; then
    echo "-- Generating $TAGS tags..."
    begin_batch "public.tags ( created_at, updated_at, name, description, color, post_count)"
    
    tag_names=("技术" "生活" "职场" "前端" "后端" "数据库" "云原生" "AI" "开源" "教程" "面试" "吐槽" "分享" "提问" "Bug" "优化" "安全" "测试" "运维" "架构")
    
    for i in $(seq 1 $TAGS); do
        name="${tag_names[$((i-1)) % ${#tag_names[@]}]}-$i"
        color=$(random_pick "#6366f1" "#ef4444" "#10b981" "#f59e0b" "#8b5cf6" "#ec4899")
        post_count=$(random_int 0 50)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '$name', '标签描述 $i', '$color', $post_count)"
    done
    flush_batch
fi

# 4. topics（依赖 users）
if [[ $TOPICS -gt 0 ]]; then
    echo "-- Generating $TOPICS topics..."
    begin_batch "public.topics ( created_at, updated_at, title, description, cover, creator_id, is_public, post_count, follower_count)"
    
    for i in $(seq 1 $TOPICS); do
        creator_id=$(random_int 1 $USERS)
        is_public=$(random_pick "true" "true" "false")
        post_count=$(random_int 0 20)
        follower_count=$(random_int 0 100)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '话题标题 $i', '话题描述 $i', 'cover-$i.jpg', $creator_id, $is_public, $post_count, $follower_count)"
    done
    flush_batch
fi

# 5. posts（依赖 users, boards）
if [[ $POSTS -gt 0 ]]; then
    echo "-- Generating $POSTS posts..."
    begin_batch "public.posts ( created_at, updated_at, title, content, summary, cover, type, post_status, moderation_status, author_id, view_count, like_count, pin_top, board_id, pin_in_board)"
    
    for i in $(seq 1 $POSTS); do
        author_id=$(random_int 1 $USERS)
        board_id="NULL"
        [[ $BOARDS -gt 0 ]] && board_id=$(random_int 1 $BOARDS)
        type=$(random_pick "post" "post" "article" "topic" "question")
        status=$(random_pick "published" "published" "draft" "archived")
        mod_status=$(random_pick "approved" "approved" "pending" "rejected")
        view_count=$(random_int 0 10000)
        like_count=$(random_int 0 500)
        pin_top=$(random_pick "false" "false" "false" "true")
        pin_in_board=$(random_pick "false" "false" "true")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '$(random_title)', '$(random_content)', '摘要 $i', 'cover-$i.jpg', '$type', '$status', '$mod_status', $author_id, $view_count, $like_count, $pin_top, $board_id, $pin_in_board)"
    done
    flush_batch
fi

# 6. post_tags（关联表，依赖 posts, tags）
if [[ $POSTS -gt 0 && $TAGS -gt 0 ]]; then
    echo "-- Generating post_tags..."
    
    # 创建临时文件存储所有随机组合
    tmpfile=$(mktemp /tmp/post_tags.XXXXXX)
    max_pairs=$((POSTS * 2))
    [[ $max_pairs -gt 500 ]] && max_pairs=500
    
    i=0
    while [[ $i -lt $max_pairs ]]; do
        post_id=$(random_int 1 $POSTS)
        tag_id=$(random_int 1 $TAGS)
        echo "$post_id,$tag_id" >> "$tmpfile"
        ((i++))
    done
    
    # 去重并生成 SQL
    unique_pairs=$(sort -u "$tmpfile" | while IFS=, read p t; do
        echo "($p, $t)"
    done | paste -sd, -)
    
    rm -f "$tmpfile"
    
    if [[ -n "$unique_pairs" ]]; then
        sql="INSERT INTO public.post_tags (post_id, tag_id) VALUES $unique_pairs ON CONFLICT DO NOTHING;"
        run_sql "$sql"
    fi
fi

# 7. comments（依赖 users, posts, comments）
if [[ $COMMENTS -gt 0 ]]; then
    echo "-- Generating $COMMENTS comments..."
    begin_batch "public.comments ( created_at, updated_at, content, post_id, author_id, parent_id, like_count, status, is_answer, is_accepted, vote_count)"
    
    for i in $(seq 1 $COMMENTS); do
        post_id=$(random_int 1 $POSTS)
        author_id=$(random_int 1 $USERS)
        parent_id="NULL"
        [[ $i -gt 10 && $(random_int 1 3) -eq 1 ]] && parent_id=$(random_int 1 $((i-1)))
        is_answer=$(random_pick "false" "false" "true")
        is_accepted="false"
        [[ "$is_answer" == "true" && $(random_int 1 4) -eq 1 ]] && is_accepted="true"
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '评论内容 $i: $(random_content)', $post_id, $author_id, $parent_id, $(random_int 0 100), 'visible', $is_answer, $is_accepted, $(random_int 0 50))"
    done
    flush_batch
fi

# 8. questions（依赖 posts, comments）
if [[ $QUESTIONS -gt 0 && $POSTS -gt 0 ]]; then
    echo "-- Generating $QUESTIONS questions..."
    begin_batch "public.questions ( created_at, updated_at, post_id, accepted_answer_id, reward_score, answer_count, view_count)"
    
    for i in $(seq 1 $QUESTIONS); do
        post_id=$((i % POSTS + 1))
        accepted_answer_id="NULL"
        [[ $COMMENTS -gt 0 && $(random_int 1 3) -eq 1 ]] && accepted_answer_id=$(random_int 1 $COMMENTS)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $post_id, $accepted_answer_id, $(random_int 0 100), $(random_int 0 20), $(random_int 0 1000))"
    done
    flush_batch
fi

# 9. likes（依赖 users, posts, comments）
if [[ $LIKES -gt 0 ]]; then
    echo "-- Generating $LIKES likes..."
    begin_batch "public.likes ( created_at, updated_at, user_id, post_id, comment_id)"
    
    for i in $(seq 1 $LIKES); do
        user_id=$(random_int 1 $USERS)
        post_id="NULL"
        comment_id="NULL"
        if [[ $(random_int 1 2) -eq 1 && $POSTS -gt 0 ]]; then
            post_id=$(random_int 1 $POSTS)
        elif [[ $COMMENTS -gt 0 ]]; then
            comment_id=$(random_int 1 $COMMENTS)
        fi
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $post_id, $comment_id)"
    done
    flush_batch
fi

# 10. follows（依赖 users）
if [[ $FOLLOWS -gt 0 ]]; then
    echo "-- Generating $FOLLOWS follows..."
    begin_batch "public.follows ( created_at, updated_at, follower_id, following_id)"
    
    for i in $(seq 1 $FOLLOWS); do
        follower_id=$(random_int 1 $USERS)
        following_id=$(random_int 1 $USERS)
        while [[ $follower_id -eq $following_id ]]; do
            following_id=$(random_int 1 $USERS)
        done
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $follower_id, $following_id)"
    done
    flush_batch
fi

# 11. notifications（依赖 users）
if [[ $NOTIFICATIONS -gt 0 ]]; then
    echo "-- Generating $NOTIFICATIONS notifications..."
    begin_batch "public.notifications ( created_at, updated_at, user_id, sender_id, type, content, target_id, target_type, is_read)"
    
    for i in $(seq 1 $NOTIFICATIONS); do
        user_id=$(random_int 1 $USERS)
        sender_id=$(random_int 1 $USERS)
        type=$(random_pick "like" "comment" "follow" "mention" "system")
        content="通知内容 $i"
        target_id=$(random_int 1 $POSTS)
        target_type=$(random_pick "post" "comment" "user")
        is_read=$(random_pick "true" "false")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $sender_id, '$type', '$content', $target_id, '$target_type', $is_read)"
    done
    flush_batch
fi

# 12. topic_posts（依赖 topics, posts, users）
if [[ $TOPIC_POSTS -gt 0 && $TOPICS -gt 0 && $POSTS -gt 0 ]]; then
    echo "-- Generating $TOPIC_POSTS topic_posts..."
    begin_batch "public.topic_posts ( created_at, updated_at, topic_id, post_id, sort_order, added_by)"
    
    for i in $(seq 1 $TOPIC_POSTS); do
        topic_id=$(random_int 1 $TOPICS)
        post_id=$(random_int 1 $POSTS)
        added_by=$(random_int 1 $USERS)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $topic_id, $post_id, $i, $added_by)"
    done
    flush_batch
fi

# 13. topic_follows（依赖 users, topics）
if [[ $TOPIC_FOLLOWS -gt 0 && $USERS -gt 0 && $TOPICS -gt 0 ]]; then
    echo "-- Generating $TOPIC_FOLLOWS topic_follows..."
    begin_batch "public.topic_follows ( created_at, updated_at, user_id, topic_id)"
    
    for i in $(seq 1 $TOPIC_FOLLOWS); do
        user_id=$(random_int 1 $USERS)
        topic_id=$(random_int 1 $TOPICS)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $topic_id)"
    done
    flush_batch
fi

# 14. announcements（依赖 users, boards）
if [[ $ANNOUNCEMENTS -gt 0 ]]; then
    echo "-- Generating $ANNOUNCEMENTS announcements..."
    begin_batch "public.announcements ( created_at, updated_at, title, content, summary, cover, type, status, is_pinned, is_global, board_id, published_at, expired_at, view_count, created_by, updated_by)"
    
    for i in $(seq 1 $ANNOUNCEMENTS); do
        created_by=$(random_int 1 $USERS)
        board_id="NULL"
        [[ $BOARDS -gt 0 && $(random_int 1 3) -eq 1 ]] && board_id=$(random_int 1 $BOARDS)
        is_global=$(random_pick "true" "true" "false")
        is_pinned=$(random_pick "false" "false" "true")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '公告标题 $i', '公告内容 $i', '摘要 $i', 'cover-$i.jpg', $(random_int 0 3), $(random_int 0 2), $is_pinned, $is_global, $board_id, '$(random_timestamp)', '$(random_timestamp 30)', $(random_int 0 10000), $created_by, $created_by)"
    done
    flush_batch
fi

# 15. reports（依赖 users）
if [[ $REPORTS -gt 0 ]]; then
    echo "-- Generating $REPORTS reports..."
    begin_batch "public.reports ( created_at, updated_at, reporter_id, target_id, target_type, type, reason, status, handler_id, handle_note, handle_at, content_snapshot, reporter_ip, is_anonymous, priority)"
    
    for i in $(seq 1 $REPORTS); do
        reporter_id=$(random_int 1 $USERS)
        target_id=$(random_int 1 $POSTS)
        target_type=$(random_pick "post" "comment" "user")
        status=$(random_pick "pending" "pending" "resolved" "rejected")
        handler_id="NULL"
        [[ "$status" != "pending" ]] && handler_id=$(random_int 1 $USERS)
        is_anonymous=$(random_pick "false" "false" "true")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $reporter_id, $target_id, '$target_type', 'other', '举报理由 $i', '$status', $handler_id, '处理备注 $i', '$(random_timestamp)', '内容快照 $i', '$(random_ip)', $is_anonymous, $(random_int 1 3))"
    done
    flush_batch
fi

# 16. sign_ins（依赖 users）
if [[ $SIGN_INS -gt 0 ]]; then
    echo "-- Generating $SIGN_INS sign_ins..."
    begin_batch "public.sign_ins ( created_at, updated_at, user_id, sign_date, score, continued)"
    
    for i in $(seq 1 $SIGN_INS); do
        user_id=$(random_int 1 $USERS)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, '$(random_timestamp 7)', $(random_int 5 20), $(random_int 1 30))"
    done
    flush_batch
fi

# 17. moderators（依赖 users, boards）
if [[ $MODERATORS -gt 0 && $USERS -gt 0 && $BOARDS -gt 0 ]]; then
    echo "-- Generating $MODERATORS moderators..."
    begin_batch "public.moderators ( created_at, updated_at, user_id, board_id, permissions)"
    
    for i in $(seq 1 $MODERATORS); do
        user_id=$(random_int 1 $USERS)
        board_id=$(random_int 1 $BOARDS)
        permissions="'{\"delete_post\": true, \"pin_post\": true}'"
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $board_id, $permissions)"
    done
    flush_batch
fi

# 18. moderator_applications（依赖 users, boards）
if [[ $MODERATOR_APPLICATIONS -gt 0 ]]; then
    echo "-- Generating $MODERATOR_APPLICATIONS moderator_applications..."
    begin_batch "public.moderator_applications ( created_at, updated_at, user_id, board_id, reason, status, reviewer_id, review_note, req_delete_post, req_pin_post, req_edit_any_post, req_manage_moderator, req_ban_user)"
    
    for i in $(seq 1 $MODERATOR_APPLICATIONS); do
        user_id=$(random_int 1 $USERS)
        board_id=$(random_int 1 $BOARDS)
        status=$(random_pick "pending" "approved" "rejected")
        reviewer_id="NULL"
        [[ "$status" != "pending" ]] && reviewer_id=$(random_int 1 $USERS)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $board_id, '申请理由 $i', '$status', $reviewer_id, '审核备注 $i', true, true, false, false, true)"
    done
    flush_batch
fi

# 19. audit_logs（依赖 users）
if [[ $AUDIT_LOGS -gt 0 ]]; then
    echo "-- Generating $AUDIT_LOGS audit_logs..."
    begin_batch "public.audit_logs ( created_at, updated_at, operator_id, operator_ip, action, target_type, target_id, before, after, reason, ip)"
    
    for i in $(seq 1 $AUDIT_LOGS); do
        operator_id=$(random_int 1 $USERS)
        action=$(random_pick "create" "update" "delete" "login" "logout")
        target_type=$(random_pick "post" "user" "comment" "board")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $operator_id, '$(random_ip)', '$action', '$target_type', $i, '{\"old\": \"value\"}', '{\"new\": \"value\"}', '操作原因 $i', '$(random_ip)')"
    done
    flush_batch
fi

# 20. refresh_tokens（依赖 users）
if [[ $REFRESH_TOKENS -gt 0 ]]; then
    echo "-- Generating $REFRESH_TOKENS refresh_tokens..."
    begin_batch "public.refresh_tokens ( created_at, updated_at, user_id, token, jti, user_agent, ip, expires_at, is_used)"
    
    for i in $(seq 1 $REFRESH_TOKENS); do
        user_id=$(random_int 1 $USERS)
        token="token_$(random_string 32)_$i"
        jti="$(random_uuid)"
        is_used=$(random_pick "false" "false" "true")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, '$token', '$jti', 'Mozilla/5.0', '$(random_ip)', '$(random_timestamp 30)', $is_used)"
    done
    flush_batch
fi

# 21. attachments（依赖 users）
if [[ $ATTACHMENTS -gt 0 ]]; then
    echo "-- Generating $ATTACHMENTS attachments..."
    begin_batch "public.attachments ( created_at, updated_at, file_id, user_id, plugin_id, post_id, reply_id, original_name, stored_name, stored_path, size, file_type, mime_type, mime_major, ext, width, height, status, upload_ip, plugin_meta, file_hash)"
    
    for i in $(seq 1 $ATTACHMENTS); do
        user_id=$(random_int 1 $USERS)
        post_id=$(random_int 1 $POSTS)
        ext=$(random_pick "jpg" "png" "pdf" "mp4")
        mime_major=$(random_pick "image" "image" "application" "video")
        mime_type="$mime_major/$ext"
        [[ "$ext" == "jpg" ]] && mime_type="image/jpeg"
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '$(random_uuid)', $user_id, NULL, $post_id, 0, 'file_$i.$ext', 'stored_$i.$ext', 'uploads/2024/$i.$ext', $(random_int 1024 10485760), 'attachment', '$mime_type', '$mime_major', '$ext', $(random_int 100 1920), $(random_int 100 1080), 1, '$(random_ip)', NULL, '$(random_string 64)')"
    done
    flush_batch
fi

# 22. plugins（基础表）
if [[ $PLUGINS -gt 0 ]]; then
    echo "-- Generating $PLUGINS plugins..."
    begin_batch "public.plugins ( created_at, updated_at, name, slug, version, description, summary, icon_url, screenshots, homepage_url, type, category, tags, author_id, author_email, author_url, script_url, server_entry, slots, routes, pricing, compatibility, permissions, enabled, status, install_count, rating, config_schema, config)"
    
    for i in $(seq 1 $PLUGINS); do
        name="插件_$i"
        slug="plugin-$(random_string 6)"
        enabled=$(random_pick "true" "false")
        status=$([[ "$enabled" == "true" ]] && echo "active" || echo "inactive")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '$name', '$slug', '1.0.0', '插件描述 $i', '摘要 $i', 'icon.png', '[]', 'https://example.com/$slug', 'widget', 'tool', '[\"tag1\", \"tag2\"]', $(random_int 1 $USERS), 'author@example.com', 'https://author.com', 'https://cdn.example.com/$slug.js', 'server.js', '[]', '[]', '{}', '{}', '{}', $enabled, '$status', $(random_int 0 1000), $(random_int 10 50)/10.0, '{}', '{}')"
    done
    flush_batch
fi

# 23. bots（依赖 users）
if [[ $BOTS -gt 0 ]]; then
    echo "-- Generating $BOTS bots..."
    begin_batch "public.bots ( created_at, updated_at, name, version, description, summary, avatar_url, screenshots, homepage_url, type, tags, creator_id, creator_name, script_code, script_url, trigger_type, cron_expr, event_filter, timeout_sec, retry_times, env_vars, resource_limit, pricing, permissions, enabled, status, exec_count, last_exec_at, error_msg, config_schema, config_values)"
    
    for i in $(seq 1 $BOTS); do
        creator_id=$(random_int 1 $USERS)
        enabled=$(random_pick "true" "false")
        status=$([[ "$enabled" == "true" ]] && echo "active" || echo "inactive")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', 'Bot_$i', '1.0.0', '机器人描述 $i', '摘要 $i', 'avatar.png', '[]', 'https://example.com/bot$i', 'cron', '[\"auto\"]', $creator_id, 'creator_$creator_id', 'console.log(\"hello\")', 'https://cdn.example.com/bot$i.js', 'cron', '0 0 * * *', '', 10, 3, '{}', '{}', '{}', '{}', $enabled, '$status', $(random_int 0 1000), '$(random_timestamp)', '', '{}', '{}')"
    done
    flush_batch
fi

# 24. blocked_ips
if [[ $BLOCKED_IPS -gt 0 ]]; then
    echo "-- Generating $BLOCKED_IPS blocked_ips..."
    begin_batch "public.blocked_ips ( ip, reason, operator_id, expire_at, created_at, updated_at)"
    
    for i in $(seq 1 $BLOCKED_IPS); do
        operator_id=$(random_int 1 $USERS)
        
        add_batch "( '$(random_ip)', '封禁原因 $i', $operator_id, '$(random_timestamp 30)', '$(random_timestamp)', '$(random_timestamp)')"
    done
    flush_batch
fi

# 25. ip_risk_records
if [[ $IP_RISK_RECORDS -gt 0 ]]; then
    echo "-- Generating $IP_RISK_RECORDS ip_risk_records..."
    begin_batch "public.ip_risk_records ( ip, event_type, event_detail, expire_at, created_at)"
    
    for i in $(seq 1 $IP_RISK_RECORDS); do
        event_type=$(random_pick "brute_force" "spam" "abuse" "scan")
        
        add_batch "( '$(random_ip)', '$event_type', '事件详情 $i', '$(random_timestamp 30)', '$(random_timestamp)')"
    done
    flush_batch
fi

# 26. user_risk_records（依赖 users）
if [[ $USER_RISK_RECORDS -gt 0 ]]; then
    echo "-- Generating $USER_RISK_RECORDS user_risk_records..."
    begin_batch "public.user_risk_records ( created_at, updated_at, user_id, event_type, event_detail, expire_at)"
    
    for i in $(seq 1 $USER_RISK_RECORDS); do
        user_id=$(random_int 1 $USERS)
        event_type=$(random_pick "spam" "abuse" "cheat" "harass")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, '$event_type', '事件详情 $i', '$(random_timestamp 30)')"
    done
    flush_batch
fi

# 27. content_audit_tasks
if [[ $CONTENT_AUDIT_TASKS -gt 0 ]]; then
    echo "-- Generating $CONTENT_AUDIT_TASKS content_audit_tasks..."
    begin_batch "public.content_audit_tasks ( created_at, updated_at, target_type, target_id, trigger_type, trigger_meta, status, reviewer_id, review_note, reviewed_at)"
    
    for i in $(seq 1 $CONTENT_AUDIT_TASKS); do
        target_type=$(random_pick "post" "comment" "user")
        target_id=$(random_int 1 $POSTS)
        status=$(random_pick "pending" "approved" "rejected")
        reviewer_id="NULL"
        [[ "$status" != "pending" ]] && reviewer_id=$(random_int 1 $USERS)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', '$target_type', $target_id, 'auto', '{}', '$status', $reviewer_id, '审核备注 $i', '$(random_timestamp)')"
    done
    flush_batch
fi

# 28. timeline_events（依赖 users）
if [[ $TIMELINE_EVENTS -gt 0 ]]; then
    echo "-- Generating $TIMELINE_EVENTS timeline_events..."
    begin_batch "public.timeline_events ( created_at, updated_at, user_id, actor_id, action, target_id, target_type, payload, score)"
    
    for i in $(seq 1 $TIMELINE_EVENTS); do
        user_id=$(random_int 1 $USERS)
        actor_id=$(random_int 1 $USERS)
        action=$(random_pick "post" "like" "follow" "comment")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $actor_id, '$action', $(random_int 1 $POSTS), 'post', '{}', $(random_int 0 100))"
    done
    flush_batch
fi

# 29. timeline_subscriptions（依赖 users）
if [[ $TIMELINE_SUBSCRIPTIONS -gt 0 ]]; then
    echo "-- Generating $TIMELINE_SUBSCRIPTIONS timeline_subscriptions..."
    begin_batch "public.timeline_subscriptions ( created_at, updated_at, subscriber_id, target_user_id, target_type, target_id, is_active)"
    
    for i in $(seq 1 $TIMELINE_SUBSCRIPTIONS); do
        subscriber_id=$(random_int 1 $USERS)
        target_user_id=$(random_int 1 $USERS)
        while [[ $subscriber_id -eq $target_user_id ]]; do
            target_user_id=$(random_int 1 $USERS)
        done
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $subscriber_id, $target_user_id, 'user', $target_user_id, true)"
    done
    flush_batch
fi

# 30. favorites（依赖 users）
if [[ $FAVORITES -gt 0 ]]; then
    echo "-- Generating $FAVORITES favorites..."
    begin_batch "public.favorites ( created_at, updated_at, user_id, target_id, target_type, group_id, status)"
    
    for i in $(seq 1 $FAVORITES); do
        user_id=$(random_int 1 $USERS)
        target_id=$(random_int 1 $POSTS)
        target_type=$(random_pick "post" "comment")
        group_id=$(random_int 1 5)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $target_id, '$target_type', $group_id, 1)"
    done
    flush_batch
fi

# 31. violations（依赖 users）
if [[ $VIOLATIONS -gt 0 ]]; then
    echo "-- Generating $VIOLATIONS violations..."
    begin_batch "public.violations ( created_at, updated_at, user_id, violation_type, reason, content_snapshot, evidence_url, source, status, operator_id, punish_type, punish_expire_at, appeal_status, appeal_reason, appeal_time, appeal_result)"
    
    for i in $(seq 1 $VIOLATIONS); do
        user_id=$(random_int 1 $USERS)
        status=$(random_int 1 3)
        operator_id="NULL"
        [[ $status -gt 1 ]] && operator_id=$(random_int 1 $USERS)
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $(random_int 1 10), '违规原因 $i', '内容快照 $i', 'https://evidence.com/$i.jpg', 1, $status, $operator_id, $(random_int 0 3), '$(random_timestamp 30)', 0, '', '$(random_timestamp)', '')"
    done
    flush_batch
fi

# 32. votes（依赖 users, comments）
if [[ $VOTES -gt 0 && $USERS -gt 0 && $COMMENTS -gt 0 ]]; then
    echo "-- Generating $VOTES votes..."
    begin_batch "public.votes ( user_id, comment_id, value, created_at, updated_at)"
    
    for i in $(seq 1 $VOTES); do
        user_id=$(random_int 1 $USERS)
        comment_id=$(random_int 1 $COMMENTS)
        value=$(random_pick "1" "-1")
        
        add_batch "( $user_id, $comment_id, $value, '$(random_timestamp)', '$(random_timestamp)')"
    done
    flush_batch
fi

# 33. answer_votes（依赖 users, comments）
if [[ $ANSWER_VOTES -gt 0 && $USERS -gt 0 && $COMMENTS -gt 0 ]]; then
    echo "-- Generating $ANSWER_VOTES answer_votes..."
    begin_batch "public.answer_votes ( created_at, updated_at, user_id, comment_id, vote_type)"
    
    for i in $(seq 1 $ANSWER_VOTES); do
        user_id=$(random_int 1 $USERS)
        comment_id=$(random_int 1 $COMMENTS)
        vote_type=$(random_pick "up" "down")
        
        add_batch "( '$(random_timestamp)', '$(random_timestamp)', $user_id, $comment_id, '$vote_type')"
    done
    flush_batch
fi

# 34. casbin_rule（基础数据）
# if [[ $(random_int 1 2) -eq 1 ]]; then
#     echo "-- Generating casbin rules..."
#     run_sql "INSERT INTO public.casbin_rule ( ptype, v0, v1, v2, v3, v4, v5) VALUES 
#         (1, 'p', 'admin', '*', '*', '', '', ''),
#         (2, 'p', 'moderator', 'board', 'moderate', '', '', ''),
#         (3, 'p', 'user', 'post', 'create', '', '', ''),
#         (4, 'g', 'moderator', 'user', '', '', '', '')
#         ON CONFLICT DO NOTHING;"
# fi

echo ""
echo "-- Mock data generation complete!"
if ! $EXECUTE; then
    echo "-- This was a DRY RUN. Add --execute to actually insert data."
    echo "-- Example: USERS=100 POSTS=500 ./mock_data.sh --execute"
fi