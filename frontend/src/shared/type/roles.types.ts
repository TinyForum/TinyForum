// constants/roles.ts

export const UserRole = {
  Guest: "guest",
  User: "user",
  Member: "member",
  Moderator: "moderator",
  Reviewer: "reviewer",
  Admin: "admin",
  SuperAdmin: "super_admin",
  SystemMaintainer: "system_maintainer",
} as const;

export type UserRoleType = (typeof UserRole)[keyof typeof UserRole];

// 角色显示名称映射
export const roleDisplayMap: Record<UserRoleType, string> = {
  [UserRole.Guest]: "访客",
  [UserRole.User]: "用户",
  [UserRole.Member]: "会员",
  [UserRole.Moderator]: "版主",
  [UserRole.Reviewer]: "审核员",
  [UserRole.Admin]: "管理员",
  [UserRole.SuperAdmin]: "超级管理员",
  [UserRole.SystemMaintainer]: "系统管理员",
};

// 角色等级（用于权限判断）
export const roleLevel: Record<UserRoleType, number> = {
  [UserRole.Guest]: 0,
  [UserRole.User]: 1,
  [UserRole.Member]: 2,
  [UserRole.SystemMaintainer]: 3,
  [UserRole.Reviewer]: 10,
  [UserRole.Moderator]: 20,
  [UserRole.Admin]: 50,
  [UserRole.SuperAdmin]: 100,
};

// 可分配的角色（管理员可以设置的角色）
export const assignableRoles: UserRoleType[] = [
  UserRole.User,
  UserRole.Member,
  UserRole.Moderator,
  UserRole.Reviewer,
  UserRole.Admin,
];

// 超级管理员可分配的角色（包含超级管理员）
export const superAdminAssignableRoles: UserRoleType[] = [
  ...assignableRoles,
  UserRole.SuperAdmin,
];
