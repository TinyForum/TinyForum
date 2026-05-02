// SecuritySettings.tsx
"use client";
import { SecurityOverview } from "@/features/settings/components/SecurityOverview";
import { SecuritySettingsHeader } from "@/features/settings/components/SecuritySettingsHeader";
import { TwoFactorAuthCard } from "@/features/settings/components/TwoFactorAuthCard";
import { ChangePasswordForm } from "@/features/settings/components/ChangePasswordForm";
import { formatDate } from "@/shared/lib/utils";
import { UserDO } from "@/shared/api/types/user.model";

// ===================== 主组件 =====================
interface SecuritySettingsProps {
  user: UserDO;
}

export default function SecuritySettings({ user }: SecuritySettingsProps) {
  // 示例数据，实际应从 API 获取
  const lastLoginTime = formatDate(new Date(user.created_at).toString());

  return (
    <div className="space-y-6">
      <SecuritySettingsHeader />
      <SecurityOverview
        user={user}
        lastLoginTime={lastLoginTime}
        securityLevel={0}
      />
      <ChangePasswordForm />
      <TwoFactorAuthCard />
    </div>
  );
}
