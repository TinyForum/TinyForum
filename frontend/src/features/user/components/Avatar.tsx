// src/components/user/Avatar.tsx
import { useMemo, useState } from "react";
import Image from "next/image";

interface AvatarProps {
  username?: string;
  avatarUrl?: string;
  size?: "sm" | "md" | "lg" | number | "full";
  className?: string;
  onError?: () => void;
  shape?: "circle" | "rounded" | "square";
  roundedSize?: "sm" | "md" | "lg" | "xl" | "2xl" | "3xl" | "full" | number;
  ring?: boolean;
  ringColor?: string;
  ringOffset?: boolean;
  ringOffsetColor?: string;
  zIndex?: number | string;
}

const SIZE_CLASS_MAP: Record<string, string> = {
  sm: "w-8 h-8",
  md: "w-10 h-10",
  lg: "w-14 h-14",
  full: "w-full h-full",
};

const ROUNDED_CLASS_MAP: Record<string, string> = {
  sm: "rounded-sm",
  md: "rounded-md",
  lg: "rounded-lg",
  xl: "rounded-xl",
  "2xl": "rounded-2xl",
  "3xl": "rounded-3xl",
  full: "rounded-full",
};

export default function Avatar({
  username,
  avatarUrl,
  size = "full",
  className = "",
  onError,
  shape = "circle",
  roundedSize = "full",
  ring = false,
  ringColor = "primary",
  ringOffset = false,
  ringOffsetColor = "base-100",
  zIndex = "auto",
}: AvatarProps) {
  const [hasError, setHasError] = useState(false);

  const handleError = () => {
    if (!hasError) {
      setHasError(true);
      onError?.();
    }
  };

  const getShapeClass = (): string => {
    if (shape === "circle") return "rounded-full";
    if (shape === "rounded") {
      if (typeof roundedSize === "number") return `rounded-[${roundedSize}px]`;
      return ROUNDED_CLASS_MAP[roundedSize] || "rounded-full";
    }
    return "rounded-none";
  };

  const getRingClass = (): string => {
    if (!ring) return "";
    const ringCls = `ring ring-${ringColor}`;
    const offsetCls = ringOffset
      ? `ring-offset-${ringOffsetColor} ring-offset-2`
      : "";
    return [ringCls, offsetCls].filter(Boolean).join(" ");
  };

  const shapeClass = getShapeClass();
  const ringClass = getRingClass();
  const sizeClass = typeof size !== "number" ? SIZE_CLASS_MAP[size] : "";
  const inlineSizeStyle =
    typeof size === "number" ? { width: size, height: size } : undefined;

  // 重要：外层 div 必须是 relative，且内层 div 使用 absolute 填充
  const wrapperClasses = `relative ${sizeClass} ${className}`.trim();
  const innerClasses = `absolute inset-0 overflow-hidden bg-base-200 flex items-center justify-center ${shapeClass} ${ringClass}`;

  // 占位符或错误状态
  if (!avatarUrl || hasError) {
    return (
      <div className={wrapperClasses} style={{ ...inlineSizeStyle, zIndex }}>
        <div className={innerClasses}>
          <span className="text-sm font-medium text-base-content">
            {username?.charAt(0)?.toUpperCase() || "?"}
          </span>
        </div>
      </div>
    );
  }

  // 正常图片：使用 fill 模式
  return (
    <div className={wrapperClasses} style={{ ...inlineSizeStyle, zIndex }}>
      <div className={innerClasses}>
        <Image
          src={avatarUrl}
          alt={username || "用户头像"}
          fill
          className="object-cover"
          onError={handleError}
          loading="lazy"
          sizes={
            size === "full"
              ? "100%"
              : typeof size === "number"
                ? `${size}px`
                : "40px"
          }
          unoptimized={!avatarUrl.startsWith("/") && !avatarUrl.includes("cdn")}
        />
      </div>
    </div>
  );
}
