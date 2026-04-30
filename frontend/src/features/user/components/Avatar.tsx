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

const SIZE_MAP: Record<string, number> = {
  sm: 32,
  md: 40,
  lg: 56,
};

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
  zIndex = "auto",
}: AvatarProps) {
  const [hasError, setHasError] = useState(false);

  const sizeStyle = useMemo(() => {
    if (size === "full") return { width: "100%", height: "100%" };
    if (typeof size === "number") return { width: size, height: size };
    const pxSize = SIZE_MAP[size];
    return { width: pxSize, height: pxSize };
  }, [size]);

  const imageSize = useMemo(() => {
    if (size === "full") return undefined;
    if (typeof size === "number") return size;
    return SIZE_MAP[size];
  }, [size]);

  const getShapeStyles = (): string => {
    if (shape === "circle") return "rounded-full";
    if (shape === "rounded") {
      if (typeof roundedSize === "number") return `rounded-[${roundedSize}px]`;
      return ROUNDED_MAP[roundedSize] || "rounded-full";
    }
    return "";
  };

  const getRingStyles = (): string => {
    if (!ring) return "";
    return `ring ring-${ringColor} ${ringOffset ? `ring-offset-${ringOffsetColor} ring-offset-2` : ""}`;
  };

  const handleError = () => {
    if (!hasError) {
      setHasError(true);
      onError?.();
    }
  };

  const shapeStyle = getShapeStyles();
  const ringStyle = getRingStyles();
  const isFullSize = size === "full";

  // 占位符 / 错误状态
  if (!avatarUrl || hasError) {
    return (
      <div className={`avatar ${className}`} style={{ zIndex }}>
        <div
          className={`overflow-hidden bg-base-200 flex items-center justify-center ${shapeStyle} ${ringStyle}`}
          style={sizeStyle}
        >
          <span className="text-sm font-medium">
            {username?.charAt(0)?.toUpperCase() || "头像加载失败"}
          </span>
        </div>
      </div>
    );
  }

  return (
    <div className={`avatar ${className}`} style={{ zIndex }}>
      <div
        className={`relative overflow-hidden bg-base-200 ${shapeStyle} ${ringStyle}`}
        style={sizeStyle}
      >
        {isFullSize ? (
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
          <Image
            src={avatarUrl}
            alt={username || "用户头像"}
            fill
            className="object-cover"
            onError={handleError}
            loading="lazy"
            sizes={`${imageSize}px`}
            unoptimized={
              !avatarUrl.startsWith("/") && !avatarUrl.includes("cdn")
            }
          />
        )}
      </div>
    </div>
  );
}
