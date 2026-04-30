import { User } from "@/shared/api";
import { Mail, Shield } from "lucide-react";
import { History } from "lucide-react";
interface SecurityOverviewProps {
  user: User;
  lastLoginTime: string;
  securityLevel: number;
}

export function SecurityOverview({
  user,
  lastLoginTime,
  securityLevel,
}: SecurityOverviewProps) {
  const emailStatus = user?.email ? "已验证" : "未验证";
  return (
    <div className="grid gap-4 md:grid-cols-3">
      {/* 安全等级卡片 */}
      <div className="card bg-base-100 border border-base-200 shadow-sm">
        <div className="card-body p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-green-100 dark:bg-green-900/20 rounded-lg">
              <Shield className="w-5 h-5 text-green-600" />
            </div>
            <div>
              <p className="text-xs text-base-content/60">安全等级</p>
              <p className="font-semibold text-sm">{securityLevel + " 分"}</p>
            </div>
          </div>
        </div>
      </div>

      {/* 最后登录卡片 */}
      <div className="card bg-base-100 border border-base-200 shadow-sm">
        <div className="card-body p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-blue-100 dark:bg-blue-900/20 rounded-lg">
              <History className="w-5 h-5 text-blue-600" />
            </div>
            <div>
              <p className="text-xs text-base-content/60">最后登录</p>
              <p className="font-semibold text-sm">{lastLoginTime}</p>
            </div>
          </div>
        </div>
      </div>

      {/* 邮箱状态卡片 */}
      <div className="card bg-base-100 border border-base-200 shadow-sm">
        <div className="card-body p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-purple-100 dark:bg-purple-900/20 rounded-lg">
              <Mail className="w-5 h-5 text-purple-600" />
            </div>
            <div>
              <p className="text-xs text-base-content/60">邮箱状态</p>
              <p className="font-semibold text-sm">{emailStatus}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
