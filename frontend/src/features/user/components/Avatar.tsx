// src/components/user/Avatar.tsx
"use client";

import { useState, useEffect } from "react";
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
  console.log("avatarUrl: ", avatarUrl);
  const [hasError, setHasError] = useState(false);
  const [normalizedSrc, setNormalizedSrc] = useState<string>("");

  // 将任意格式的图片地址转换为绝对 URL（仅客户端）
  useEffect(() => {
    if (!avatarUrl) {
      setNormalizedSrc("");
      return;
    }

    // 已经是绝对路径
    if (avatarUrl.startsWith("http://") || avatarUrl.startsWith("https://")) {
      setNormalizedSrc(avatarUrl);
      return;
    }

    // 协议相对路径 (//example.com/pic.jpg)
    if (avatarUrl.startsWith("//")) {
      setNormalizedSrc(`${window.location.protocol}${avatarUrl}`);
      return;
    }

    // 相对路径，如 /uploads/avatar.jpg
    if (avatarUrl.startsWith("/")) {
      setNormalizedSrc(`${window.location.origin}${avatarUrl}`);
      return;
    }

    // 默认原样返回（可能无效，但避免报错）
    setNormalizedSrc(avatarUrl);
  }, [avatarUrl]);

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

  // 等待客户端路径转换完成（避免服务端/客户端不一致）
  if (!normalizedSrc) {
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

  // 正常图片
  return (
    <div className={wrapperClasses} style={{ ...inlineSizeStyle, zIndex }}>
      <div className={innerClasses}>
        <Image
          src={normalizedSrc}
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
          unoptimized={normalizedSrc.startsWith("data:")}
        />
      </div>
    </div>
  );
}
