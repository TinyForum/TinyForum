import React from "react";

interface StatCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode; // 直接接收图标元素，不要求 props
  color?: string; // 图标颜色，如 'text-primary/20'
  footer?: React.ReactNode;
  className?: string;
}

export function StatCard({
  title,
  value,
  icon,
  color = "text-primary/20",
  footer,
  className = "",
}: StatCardProps) {
  return (
    <div
      className={`card bg-base-100 border border-base-300 hover:shadow-lg transition-shadow overflow-hidden relative ${className}`}
    >
      {/* 固定位置的图标容器：始终在右下角，固定尺寸 */}
      <div
        className="absolute pointer-events-none"
        style={{ right: "-20px", bottom: "-20px", scale: "3" }}
      >
        <div className={`${color} w-32 h-32 flex items-center justify-center`}>
          {icon}
        </div>
      </div>

      {/* 左侧内容区（保持不变） */}
      <div className="card-body p-5 relative z-10">
        <div>
          <p className="text-base-content/60 text-sm font-medium">{title}</p>
          <p className="text-3xl font-bold mt-1">{value}</p>
        </div>
        {footer && (
          <div className="mt-4 pt-3 border-t border-base-200 text-sm text-base-content/50">
            {footer}
          </div>
        )}
      </div>
    </div>
  );
}
