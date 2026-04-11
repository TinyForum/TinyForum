// src/components/Avatar.tsx
import { useMemo, useState } from 'react';
import Image from 'next/image';

interface AvatarProps {
  username?: string;
  avatarUrl?: string | null;
  size?: 'sm' | 'md' | 'lg' | number;
  className?: string;
  onError?: () => void;
  unoptimized?: boolean;
  fill?: boolean;
  shape?: 'circle' | 'rounded' | 'square';
  roundedSize?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl' | 'full' | number;
  ring?: boolean;
  ringColor?: string;
  ringOffset?: boolean;
  ringOffsetColor?: string;
}

export default function Avatar({ 
  username, 
  avatarUrl,
  size,
  className = '',
  onError,
  unoptimized,
  fill = false,
  shape = 'circle',
  roundedSize = 'full',
  ring = false,
  ringColor = 'primary',
  ringOffset = false,
  ringOffsetColor = 'base-100'
}: AvatarProps) {
  const [hasError, setHasError] = useState(false);
  const AVATAR_BASE_URL = process.env.NEXT_PUBLIC_AVATAR_BASE_URL || 'https://api.dicebear.com/8.x/initials/svg';
  
  // 计算尺寸
  const sizeInPx = useMemo(() => {
    if (size === undefined) return undefined;
    if (typeof size === 'number') return size;
    const sizes = { sm: 32, md: 40, lg: 56 };
    return sizes[size];
  }, [size]);

  // 判断是否使用 dicebear 头像
  const isDicebearAvatar = useMemo(() => {
    return avatarUrl?.startsWith('https://api.dicebear.com') || false;
  }, [avatarUrl]);

  // 生成头像URL
  const imageUrl = useMemo(() => {
    if (hasError || !avatarUrl) {
      const seed = username || 'default';
      const imageSize = sizeInPx ? `&size=${sizeInPx * 2}` : '';
      return `${AVATAR_BASE_URL}?seed=${encodeURIComponent(seed)}${imageSize}`;
    }
    return avatarUrl;
  }, [avatarUrl, username, sizeInPx, hasError]);

  // 处理图片加载错误
  const handleError = () => {
    setHasError(true);
    onError?.();
  };

  // 判断是否需要 unoptimized
  const shouldBeUnoptimized = useMemo(() => {
    if (unoptimized !== undefined) return unoptimized;
    return isDicebearAvatar;
  }, [unoptimized, isDicebearAvatar]);

  // 获取形状样式
  const getShapeStyles = () => {
    if (shape === 'circle') {
      return 'rounded-full';
    }
    if (shape === 'rounded') {
      if (typeof roundedSize === 'number') {
        return `rounded-[${roundedSize}px]`;
      }
      const roundedSizes = { sm: 'rounded-sm', md: 'rounded-md', lg: 'rounded-lg', xl: 'rounded-xl', '2xl': 'rounded-2xl', '3xl': 'rounded-3xl', full: 'rounded-full' };
      return roundedSizes[roundedSize as keyof typeof roundedSizes] || 'rounded-full';
    }
    return '';
  };

  // 获取环样式
  const getRingStyles = () => {
    if (!ring) return '';
    const ringClass = `ring ring-${ringColor}`;
    const ringOffsetClass = ringOffset ? `ring-offset-${ringOffsetColor} ring-offset-2` : '';
    return `${ringClass} ${ringOffsetClass}`;
  };

  const shapeStyle = getShapeStyles();
  const ringStyle = getRingStyles();

  // 如果使用 fill 模式
  if (fill) {
    return (
      <div className={`avatar ${className}`}>
        <div className={`overflow-hidden bg-base-200 flex items-center justify-center relative w-full h-full ${shapeStyle} ${ringStyle}`}>
          <Image
            src={imageUrl}
            alt={username || '用户头像'}
            fill
            className="object-cover"
            unoptimized={shouldBeUnoptimized}
            onError={handleError}
          />
        </div>
      </div>
    );
  }

  // 如果指定了尺寸，使用固定尺寸的容器
  if (sizeInPx) {
    return (
      <div className={`avatar ${className}`}>
        <div 
          className={`overflow-hidden bg-base-200 flex items-center justify-center ${shapeStyle} ${ringStyle}`}
          style={{ width: sizeInPx, height: sizeInPx }}
        >
          {imageUrl ? (
            <Image
              src={imageUrl}
              alt={username || '用户头像'}
              width={sizeInPx}
              height={sizeInPx}
              className="object-cover"
              unoptimized={shouldBeUnoptimized}
              onError={handleError}
            />
          ) : (
            <span className="text-sm font-medium">
              {username?.charAt(0) || '?'}
            </span>
          )}
        </div>
      </div>
    );
  }

  // 未指定尺寸时，使用父容器控制大小
  return (
    <div className={`avatar ${className}`}>
      <div className={`overflow-hidden bg-base-200 flex items-center justify-center w-full h-full ${shapeStyle} ${ringStyle}`}>
        {imageUrl ? (
          <Image
            src={imageUrl}
            alt={username || '用户头像'}
            width={100}
            height={100}
            className="object-cover w-full h-full"
            unoptimized={shouldBeUnoptimized}
            onError={handleError}
          />
        ) : (
          <span className="text-sm font-medium">
            {username?.charAt(0) || '?'}
          </span>
        )}
      </div>
    </div>
  );
}