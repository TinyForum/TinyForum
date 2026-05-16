package do

// ── 角色定义 ────────────────────────────────────────────────────────────────

type UserRole string

// 有效的用户角色
var validuserRole = map[UserRole]bool{
	RoleGuest:            true,
	RoleUser:             true,
	RoleMember:           true,
	RoleModerator:        true,
	RoleReviewer:         true,
	RoleAdmin:            true,
	RoleSuperAdmin:       true,
	RoleSystemMaintainer: true,
}

const (
	RoleGuest            UserRole = "guest"             // 游客：浏览（只读，受限访问）
	RoleUser             UserRole = "user"              // 普通用户：浏览、评论、发帖等（可读写：受限访问，受限读写）
	RoleMember           UserRole = "member"            // 会员：无广告、自定义表情、创建投票等（可读写：受限访问，受限读写）
	RoleModerator        UserRole = "moderator"         // 版主：管理版块、管理帖子、管理评论等（可读写：受限访问，受限读写）
	RoleReviewer         UserRole = "reviewer"          // 审核员：审核帖子、评论等（可读写：受限访问，受限读写）
	RoleAdmin            UserRole = "admin"             // 管理员：管理用户、管理帖子、管理评论等（可读写：完全访问，完全读写）
	RoleSuperAdmin       UserRole = "super_admin"       // 超级管理员：最高权限（可读写：完全访问，完全读写）
	RoleBot              UserRole = "bot"               // 系统机器人：自动回复、定时任务等（受限访问，受限读写）
	RoleSystemMaintainer UserRole = "system_maintainer" // 系统维护者：系统维护、数据备份等（受限读写，受限访问
)

// enum [guest, user, member, moderator, reviewer, admin, super_admin, bot, system_maintainer]
