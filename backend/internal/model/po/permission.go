package po

// ── 角色定义 ────────────────────────────────────────────────────────────────

type UserRole string

const (
	RoleGuest      UserRole = "guest"
	RoleUser       UserRole = "user"
	RoleMember     UserRole = "member"
	RoleModerator  UserRole = "moderator"
	RoleReviewer   UserRole = "reviewer"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"
	RoleBot        UserRole = "bot"
)

// ── 权限定义 ────────────────────────────────────────────────────────────────

type Permission string

const (
	// 帖子权限
	PermCreatePost    Permission = "post.create"
	PermEditOwnPost   Permission = "post.edit.own"
	PermEditAnyPost   Permission = "post.edit.any"
	PermDeleteOwnPost Permission = "post.delete.own"
	PermDeleteAnyPost Permission = "post.delete.any"
	PermPinPost       Permission = "post.pin"
	PermUnpinPost     Permission = "post.unpin"

	// 评论权限
	PermCreateComment    Permission = "comment.create"
	PermDeleteOwnComment Permission = "comment.delete.own"
	PermDeleteAnyComment Permission = "comment.delete.any"

	// 用户权限
	PermViewUserInfo   Permission = "user.view"
	PermEditOwnProfile Permission = "user.edit.own"
	PermEditAnyProfile Permission = "user.edit.any"
	PermBanUser        Permission = "user.ban"

	// 版块权限
	PermManageBoard Permission = "board.manage"
	PermCreateBoard Permission = "board.create"
	PermDeleteBoard Permission = "board.delete"

	// 系统权限
	PermViewLogs       Permission = "system.logs.view"
	PermModifySettings Permission = "system.settings.modify"
	PermManageAdmin    Permission = "system.admin.manage"

	// 会员专属权限
	PermUploadFile     Permission = "member.upload.file"
	PermNoAds          Permission = "member.no.ads"
	PermUseCustomEmoji Permission = "member.custom.emoji"
	PermCreatePoll     Permission = "member.create.poll"

	// 审核员专属权限
	PermApproveContent Permission = "reviewer.approve"
	PermRejectContent  Permission = "reviewer.reject"
	PermViewPending    Permission = "reviewer.view.pending"

	// 角色管理权限（细粒度）
	PermAssignRoleUser       Permission = "role.assign.user"       // 分配普通用户角色
	PermAssignRoleMember     Permission = "role.assign.member"     // 分配会员角色
	PermAssignRoleModerator  Permission = "role.assign.moderator"  // 分配版主角色
	PermAssignRoleReviewer   Permission = "role.assign.reviewer"   // 分配审核员角色
	PermAssignRoleBot        Permission = "role.assign.bot"        // 分配机器人角色
	PermAssignRoleAdmin      Permission = "role.assign.admin"      // 分配管理员角色
	PermAssignRoleSuperAdmin Permission = "role.assign.superadmin" // 分配超级管理员角色
)

// ── 角色→权限映射 ───────────────────────────────────────────────────────────

// basePerms 普通用户基础权限集，供继承
var basePerms = []Permission{
	PermViewUserInfo,
	PermCreatePost,
	PermEditOwnPost,
	PermDeleteOwnPost,
	PermCreateComment,
	PermDeleteOwnComment,
	PermEditOwnProfile,
}

// moderatorPerms 版主权限集（含基础权限）
var moderatorPerms = append(basePerms,
	PermEditAnyPost,
	PermDeleteAnyPost,
	PermPinPost,
	PermUnpinPost,
	PermDeleteAnyComment,
	PermManageBoard,
	PermBanUser,
	PermAssignRoleUser,
	PermAssignRoleMember,
	PermAssignRoleModerator,
	PermAssignRoleReviewer,
)

// adminPerms 管理员权限集（含版主权限）
var adminPerms = append(moderatorPerms,
	PermEditAnyProfile,
	PermCreateBoard,
	PermDeleteBoard,
	PermViewLogs,
	PermModifySettings,
	PermAssignRoleBot,
	PermAssignRoleAdmin,
)

var rolePermissions = map[UserRole][]Permission{
	RoleGuest: {
		PermViewUserInfo,
	},

	RoleUser: basePerms,

	RoleMember: append(basePerms,
		PermUploadFile,
		PermNoAds,
		PermUseCustomEmoji,
		PermCreatePoll,
	),

	RoleModerator: moderatorPerms,

	RoleReviewer: {
		PermViewUserInfo,
		PermEditAnyPost,
		PermDeleteAnyPost,
		PermDeleteAnyComment,
		PermApproveContent,
		PermRejectContent,
		PermViewPending,
	},

	RoleAdmin: adminPerms,

	RoleSuperAdmin: append(adminPerms,
		PermManageAdmin,
		PermAssignRoleSuperAdmin,
	),

	RoleBot: {
		PermCreatePost,
		PermCreateComment,
		PermViewUserInfo,
	},
}

// ── 角色→可操作目标角色矩阵 ─────────────────────────────────────────────────
//
// 定义"哪些角色的用户，其角色可以被某操作者修改"。
// key = 操作者角色，value = 该操作者可以将目标修改成的角色集合。
// 注意：目标用户当前角色是否可被操作，通过 targetRoleChangeableBy 矩阵控制。

// assignableRoles 操作者可以将目标设置成的角色
var assignableRoles = map[UserRole][]UserRole{
	RoleModerator: {
		RoleUser, RoleMember, RoleModerator, RoleReviewer,
	},
	RoleAdmin: {
		RoleUser, RoleMember, RoleModerator, RoleReviewer, RoleBot, RoleAdmin,
	},
	RoleSuperAdmin: {
		RoleUser, RoleMember, RoleModerator, RoleReviewer, RoleBot, RoleAdmin, RoleSuperAdmin,
	},
}

// operableTargetRoles 操作者可以操作（修改角色）的目标用户当前角色
// 即：目标当前是 X 角色，哪些操作者能动他
var operableTargetRoles = map[UserRole][]UserRole{
	RoleModerator: {
		RoleGuest, RoleUser, RoleMember, // 版主只能动低于自己的角色
	},
	RoleAdmin: {
		RoleGuest, RoleUser, RoleMember, RoleModerator, RoleReviewer, RoleBot,
	},
	RoleSuperAdmin: {
		RoleGuest, RoleUser, RoleMember, RoleModerator, RoleReviewer, RoleBot, RoleAdmin,
	},
}

// ── 辅助查询函数 ─────────────────────────────────────────────────────────────

// GetRolePermissions 获取角色的所有权限
func GetRolePermissions(role UserRole) []Permission {
	if perms, ok := rolePermissions[role]; ok {
		return perms
	}
	return nil
}

// HasPermission 检查角色是否拥有某个权限
func HasPermission(role UserRole, perm Permission) bool {
	for _, p := range rolePermissions[role] {
		if p == perm {
			return true
		}
	}
	return false
}

// HasAllPermissions 检查角色是否同时拥有全部权限
func HasAllPermissions(role UserRole, perms ...Permission) bool {
	for _, perm := range perms {
		if !HasPermission(role, perm) {
			return false
		}
	}
	return true
}

// HasAnyPermission 检查角色是否拥有任意一个权限
func HasAnyPermission(role UserRole, perms ...Permission) bool {
	for _, perm := range perms {
		if HasPermission(role, perm) {
			return true
		}
	}
	return false
}

// CanAssignRole 检查 operator 是否可以将目标设置为 targetRole
func CanAssignRole(operatorRole UserRole, targetRole UserRole) bool {
	for _, r := range assignableRoles[operatorRole] {
		if r == targetRole {
			return true
		}
	}
	return false
}

// CanOperateTarget 检查 operator 是否可以操作当前角色为 currentRole 的用户
func CanOperateTarget(operatorRole UserRole, currentTargetRole UserRole) bool {
	for _, r := range operableTargetRoles[operatorRole] {
		if r == currentTargetRole {
			return true
		}
	}
	return false
}

// IsValidRole 检查角色字符串是否合法（不含 guest/super_admin，不允许外部直接设置）
func IsValidRole(role UserRole) bool {
	switch role {
	case RoleUser, RoleMember, RoleModerator, RoleReviewer, RoleAdmin, RoleBot, RoleSuperAdmin:
		return true
	}
	return false
}

// ── 角色优先级 ───────────────────────────────────────────────────────────────

var rolePriority = map[UserRole]int{
	RoleGuest:      0,
	RoleUser:       1,
	RoleMember:     2,
	RoleReviewer:   3,
	RoleBot:        3,
	RoleModerator:  4,
	RoleAdmin:      5,
	RoleSuperAdmin: 6,
}

func GetRolePriority(role UserRole) int {
	return rolePriority[role]
}

func IsRoleAtLeast(role, target UserRole) bool {
	return rolePriority[role] >= rolePriority[target]
}
