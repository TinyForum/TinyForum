import { Smartphone } from "lucide-react";

export function TwoFactorAuthCard() {
  return (
    <div className="card bg-primary/5 border border-primary/20">
      <div className="card-body p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/20 rounded-lg">
              <Smartphone className="w-5 h-5 text-primary" />
            </div>
            <div>
              <p className="font-medium">两步验证</p>
              <p className="text-xs text-base-content/60">
                为账户添加额外的安全保护
              </p>
            </div>
          </div>
          <button className="btn btn-sm btn-outline btn-primary">启用</button>
        </div>
      </div>
    </div>
  );
}
