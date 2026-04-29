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
}

// 尺寸映射
const SIZE_MAP: Record<string, number> = {
  sm: 32,
  md: 40,
  lg: 56,
};

// 圆角映射
const ROUNDED_MAP: Record<string, string> = {
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
}: AvatarProps) {
  const [hasError, setHasError] = useState<boolean>(false);

  // 计算尺寸样式
  const sizeStyle = useMemo(() => {
    if (size === "full") {
      return { width: "100%", height: "100%" };
    }
    if (typeof size === "number") {
      return { width: size, height: size };
    }
    const pxSize = SIZE_MAP[size];
    return { width: pxSize, height: pxSize };
  }, [size]);

  // 获取图片的实际像素尺寸（用于 Image 组件）
  const imageSize = useMemo(() => {
    if (size === "full") {
      return undefined; // 使用 fill 模式
    }
    if (typeof size === "number") {
      return size;
    }
    return SIZE_MAP[size];
  }, [size]);

  // 获取形状样式
  const getShapeStyles = (): string => {
    if (shape === "circle") return "rounded-full";
    if (shape === "rounded") {
      if (typeof roundedSize === "number") return `rounded-[${roundedSize}px]`;
      return ROUNDED_MAP[roundedSize] || "rounded-full";
    }
    if (shape === "square") return "";
    return "";
  };

  // 获取环样式
  const getRingStyles = (): string => {
    if (!ring) return "";
    const ringClass = `ring ring-${ringColor}`;
    const ringOffsetClass = ringOffset
      ? `ring-offset-${ringOffsetColor} ring-offset-2`
      : "";
    return `${ringClass} ${ringOffsetClass}`;
  };

  const handleError = (): void => {
    if (!hasError) {
      setHasError(true);
      onError?.();
    }
  };

  const shapeStyle = getShapeStyles();
  const ringStyle = getRingStyles();

  // 如果没有头像或加载失败，显示占位符
  if (!avatarUrl || hasError) {
    return (
      <div className={`avatar ${className}`}>
        <div
          className={`overflow-hidden bg-base-200 flex items-center justify-center ${shapeStyle} ${ringStyle}`}
          style={sizeStyle}
        >
          <span className="text-sm font-medium">
            {username?.charAt(0)?.toUpperCase() || "?"}
          </span>
        </div>
      </div>
    );
  }

  // 判断是否为完整尺寸（需要 fill 模式）
  const isFullSize = size === "full";

  return (
    <div className={`avatar ${className}`}>
      <div
        className={`relative overflow-hidden bg-base-200 flex items-center justify-center ${shapeStyle} ${ringStyle}`}
        style={sizeStyle}
      >
        {isFullSize ? (
          // 使用 fill 模式（父元素必须有 relative 定位）
          <Image
            src={avatarUrl}
            alt={username || "用户头像"}
            fill
            className="object-cover"
            onError={handleError}
            loading="lazy"
            sizes="100%"
            unoptimized={
              !avatarUrl.startsWith("/") && !avatarUrl.includes("cdn")
            }
          />
        ) : (
          // 使用明确的宽高
          <Image
            src={avatarUrl}
            alt={username || "用户头像"}
            width={imageSize}
            height={imageSize}
            className="object-cover"
            onError={handleError}
            loading="lazy"
            unoptimized={
              !avatarUrl.startsWith("/") && !avatarUrl.includes("cdn")
            }
          />
        )}
      </div>
    </div>
  );
}
