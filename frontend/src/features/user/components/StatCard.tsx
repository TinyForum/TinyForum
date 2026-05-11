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
      {/* 右侧图标容器：占右侧1/3，隐藏超出部分 */}
      <div className="absolute right-0 bottom-0 w-1/3 h-full overflow-hidden pointer-events-none flex items-end justify-end">
        {/* 图标包装：控制尺寸和颜色，使其足够大超出父容器 */}
        <div
          className={`${color} w-32 h-32 flex items-center justify-center scale-[3] pt-2`}
        >
          {icon}
        </div>
      </div>

      {/* 左侧内容区 */}
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
