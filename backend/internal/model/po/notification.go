package po

type NotificationType string
type Notification struct {
	BaseModel
	UserID     uint             `gorm:"not null;index" json:"user_id"`
	SenderID   *uint            `gorm:"index" json:"sender_id"`
	Type       NotificationType `gorm:"type:varchar(30)" json:"type"`
	Content    string           `gorm:"size:500" json:"content"`
	TargetID   *uint            `json:"target_id"`
	TargetType string           `gorm:"size:50" json:"target_type"`
	IsRead     bool             `gorm:"default:false" json:"is_read"`

	User   User  `gorm:"foreignKey:UserID" json:"-"`
	Sender *User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

const (
	// ==================== 1. 系统类 ====================
	NotifySystem       NotificationType = "system"       // 系统通知
	NotifyAnnouncement NotificationType = "announcement" // 全局公告
	NotifyMaintenance  NotificationType = "maintenance"  // 系统维护通知
	NotifyUpgrade      NotificationType = "upgrade"      // 版本升级通知
	NotifySecurity     NotificationType = "security"     // 安全提醒（异地登录等）
	NotifyVerify       NotificationType = "verify"       // 验证通知（邮箱/手机）
	NotifyBind         NotificationType = "bind"         // 绑定成功通知
	NotifyUnbind       NotificationType = "unbind"       // 解绑通知

	// ==================== 2. 互动类 ====================
	NotifyComment  NotificationType = "comment"  // 评论我的帖子
	NotifyLike     NotificationType = "like"     // 点赞我的内容
	NotifyReply    NotificationType = "reply"    // 回复我的评论
	NotifyMention  NotificationType = "mention"  // @提及我
	NotifyShare    NotificationType = "share"    // 分享了内容
	NotifyFavorite NotificationType = "favorite" // 收藏了我的内容
	NotifyReward   NotificationType = "reward"   // 打赏/赞赏
	NotifyReport   NotificationType = "report"   // 举报处理结果

	// ==================== 3. 管理类 ====================
	NotifyWarning    NotificationType = "warning"     // 违规警告
	NotifyMute       NotificationType = "mute"        // 禁言通知
	NotifyBan        NotificationType = "ban"         // 封禁通知
	NotifyDelete     NotificationType = "delete"      // 内容被删除通知
	NotifyMove       NotificationType = "move"        // 帖子被移动
	NotifyLock       NotificationType = "lock"        // 帖子被锁定
	NotifyUnlock     NotificationType = "unlock"      // 帖子被解锁
	NotifyStick      NotificationType = "stick"       // 帖子被置顶
	NotifyUnstick    NotificationType = "unstick"     // 取消置顶
	NotifyPromote    NotificationType = "promote"     // 内容被加精/推荐
	NotifyDemote     NotificationType = "demote"      // 取消加精
	NotifyAssignRole NotificationType = "assign_role" // 授予管理角色

	// ==================== 4. 内容类 ====================
	NotifyPostApprove    NotificationType = "post_approve"    // 帖子审核通过
	NotifyPostReject     NotificationType = "post_reject"     // 帖子审核拒绝
	NotifyPostSchedule   NotificationType = "post_schedule"   // 定时发布成功
	NotifyPostExpire     NotificationType = "post_expire"     // 帖子到期
	NotifyCommentApprove NotificationType = "comment_approve" // 评论审核通过
	NotifyCommentReject  NotificationType = "comment_reject"  // 评论审核拒绝
	NotifyDraftSaved     NotificationType = "draft_saved"     // 草稿保存成功

	// ==================== 5. 问答类 ====================
	NotifyQuestion       NotificationType = "question"        // 新问题
	NotifyAnswer         NotificationType = "answer"          // 回答了问题
	NotifyAcceptAnswer   NotificationType = "accept_answer"   // 采纳了我的回答
	NotifyAcceptCancel   NotificationType = "accept_cancel"   // 取消采纳
	NotifyQuestionUpdate NotificationType = "question_update" // 问题被修改
	NotifyQuestionClose  NotificationType = "question_close"  // 问题被关闭
	NotifyQuestionReopen NotificationType = "question_reopen" // 问题被重开
	NotifyBestAnswer     NotificationType = "best_answer"     // 被选为最佳答案
	NotifyInviteAnswer   NotificationType = "invite_answer"   // 邀请我回答问题

	// ==================== 6. 订阅类 ====================
	NotifySubscribePost    NotificationType = "subscribe_post"    // 订阅的帖子有新回复
	NotifySubscribeNode    NotificationType = "subscribe_node"    // 订阅的节点有新帖
	NotifySubscribeTag     NotificationType = "subscribe_tag"     // 订阅的标签有新内容
	NotifySubscribeUser    NotificationType = "subscribe_user"    // 订阅的用户有新动态
	NotifySubscribeComment NotificationType = "subscribe_comment" // 订阅的评论有新回复
	NotifySubscribeTopic   NotificationType = "subscribe_topic"   // 订阅的话题更新
	NotifySubscribeDigest  NotificationType = "subscribe_digest"  // 订阅的摘要/周刊

	// ==================== 7. 关注事件时间线类 ====================
	// 用户动态
	NotifyFollow         NotificationType = "follow"          // 关注了我
	NotifyUnfollow       NotificationType = "unfollow"        // 取消关注（极少用）
	NotifyFollowPost     NotificationType = "follow_post"     // 关注的人发帖
	NotifyFollowComment  NotificationType = "follow_comment"  // 关注的人评论
	NotifyFollowLike     NotificationType = "follow_like"     // 关注的人点赞
	NotifyFollowFavorite NotificationType = "follow_favorite" // 关注的人收藏
	NotifyFollowShare    NotificationType = "follow_share"    // 关注的人分享

	// 内容时间线
	NotifyTimelinePost    NotificationType = "timeline_post"    // 关注的用户发新帖
	NotifyTimelineComment NotificationType = "timeline_comment" // 关注的用户发表评论
	NotifyTimelineLike    NotificationType = "timeline_like"    // 关注的用户点赞
	NotifyTimelineAnswer  NotificationType = "timeline_answer"  // 关注的用户回答问题
	NotifyTimelineFollow  NotificationType = "timeline_follow"  // 关注的用户关注了别人

	// 热门/推荐
	NotifyTrendingPost    NotificationType = "trending_post"    // 帖子成为热门
	NotifyTrendingComment NotificationType = "trending_comment" // 评论成为热门
	NotifyRecommendPost   NotificationType = "recommend_post"   // 推荐帖子给我

	// ==================== 8. 成长/成就类 ====================
	NotifyLevelUp     NotificationType = "level_up"    // 等级提升
	NotifyBadge       NotificationType = "badge"       // 获得勋章
	NotifyAchievement NotificationType = "achievement" // 解锁成就
	NotifyMilestone   NotificationType = "milestone"   // 里程碑达成（如发帖100篇）
	NotifyPoints      NotificationType = "points"      // 积分变动
	NotifyExp         NotificationType = "exp"         // 经验值变动

	// ==================== 9. 社交/关系类 ====================
	NotifyFriendRequest  NotificationType = "friend_request"  // 好友申请
	NotifyFriendAccept   NotificationType = "friend_accept"   // 好友申请通过
	NotifyFriendOnline   NotificationType = "friend_online"   // 好友上线
	NotifyPrivateMessage NotificationType = "private_message" // 私信（通常独立，但可归入）
	NotifyGroupInvite    NotificationType = "group_invite"    // 邀请加入群组
	NotifyGroupJoin      NotificationType = "group_join"      // 加入群组

	// ==================== 10. 交易/经济类 ====================
	NotifyCredit         NotificationType = "credit"          // 积分到账
	NotifyDebit          NotificationType = "debit"           // 积分扣除
	NotifyGift           NotificationType = "gift"            // 收到礼物
	NotifyVipExpire      NotificationType = "vip_expire"      // VIP到期提醒
	NotifyVipRenew       NotificationType = "vip_renew"       // VIP续费成功
	NotifyPaymentSuccess NotificationType = "payment_success" // 支付成功
	NotifyPaymentFailed  NotificationType = "payment_failed"  // 支付失败
	NotifyRefund         NotificationType = "refund"          // 退款通知

	// ==================== 11. 活动/运营类 ====================
	NotifyEvent        NotificationType = "event"         // 活动通知
	NotifyEventRemind  NotificationType = "event_remind"  // 活动开始提醒
	NotifyEventEnd     NotificationType = "event_end"     // 活动结束
	NotifyLottery      NotificationType = "lottery"       // 中奖通知
	NotifyCoupon       NotificationType = "coupon"        // 优惠券发放
	NotifyCouponExpire NotificationType = "coupon_expire" // 优惠券即将过期

	// ==================== 12. 其他 ====================
	NotifyReminder NotificationType = "reminder" // 通用提醒
	NotifyDigest   NotificationType = "digest"   // 每日/每周摘要
	NotifyBirthday NotificationType = "birthday" // 生日祝福
	NotifyFeedback NotificationType = "feedback" // 反馈回复
	NotifySurvey   NotificationType = "survey"   // 问卷调查邀请
)
