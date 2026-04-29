// components/common/RoleBadge.tsx
import { UserRoleType } from "@/shared/type/roles.types";
import { useTranslations } from "next-intl";
// import { UserRoleType, roleLevel } from "@/constants/roles";

interface RoleBadgeProps {
  role: UserRoleType;
  showIcon?: boolean;
  size?: "sm" | "md" | "lg";
}

export function RoleBadge({
  role,
  showIcon = false,
  size = "md",
}: RoleBadgeProps) {
  const t = useTranslations("Common");

  const colorMap: Record<UserRoleType, string> = {
    guest: "bg-gray-100 text-gray-600",
    user: "bg-blue-100 text-blue-600",
    member: "bg-green-100 text-green-600",
    moderator: "bg-purple-100 text-purple-600",
    reviewer: "bg-indigo-100 text-indigo-600",
    admin: "bg-orange-100 text-orange-600",
    super_admin: "bg-red-100 text-red-600",
    bot: "bg-gray-200 text-gray-500",
  };

  const sizeMap = {
    sm: "px-1.5 py-0.5 text-xs",
    md: "px-2 py-1 text-sm",
    lg: "px-3 py-1.5 text-base",
  };

  const iconMap: Partial<Record<UserRoleType, string>> = {
    super_admin: "👑",
    admin: "🛡️",
    moderator: "🔨",
    reviewer: "👁️",
    bot: "🤖",
  };

  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full font-medium ${colorMap[role]} ${sizeMap[size]}`}
    >
      {showIcon && iconMap[role] && <span>{iconMap[role]}</span>}
      {t(`role.${role}`)}
    </span>
  );
}
