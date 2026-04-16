// src/components/Avatar.tsx
import { useMemo, useState } from 'react';

interface AvatarProps {
  username?: string;
  avatarUrl?: string | null;
  size?: 'sm' | 'md' | 'lg' | number;
  className?: string;
  onError?: () => void;
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
  size = 'md',
  className = '',
  onError,
  shape = 'circle',
  roundedSize = 'full',
  ring = false,
  ringColor = 'primary',
  ringOffset = false,
  ringOffsetColor = 'base-100'
}: AvatarProps) {
  const [hasError, setHasError] = useState(false);

  // 计算尺寸
  const sizeInPx = useMemo(() => {
    if (typeof size === 'number') return size;
    const sizes = { sm: 32, md: 40, lg: 56 };
    return sizes[size];
  }, [size]);

  // 获取形状样式
  const getShapeStyles = () => {
    if (shape === 'circle') return 'rounded-full';
    if (shape === 'rounded') {
      if (typeof roundedSize === 'number') return `rounded-[${roundedSize}px]`;
      const roundedSizes = { 
        sm: 'rounded-sm', md: 'rounded-md', lg: 'rounded-lg', 
        xl: 'rounded-xl', '2xl': 'rounded-2xl', '3xl': 'rounded-3xl', 
        full: 'rounded-full' 
      };
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

  const handleError = () => {
    if (!hasError) {
      setHasError(true);
      onError?.();
    }
  };

  const shapeStyle = getShapeStyles();
  const ringStyle = getRingStyles();

  // 获取用户首字母（用于占位符）
  const getInitial = () => {
    if (username && username.length > 0) {
      console.log(username?"用户名不为空":"用户名为空");
      return username.charAt(0).toUpperCase();
    }
    return '?';
  };

  return (
    <div className={`avatar ${className}`}>
      <div 
        className={`overflow-hidden bg-base-200 flex items-center justify-center ${shapeStyle} ${ringStyle}`}
        style={{ width: sizeInPx, height: sizeInPx }}
      >
        {avatarUrl && !hasError ? (
          <img
            src={avatarUrl}
            alt={username || '用户头像'}
            className="w-full h-full object-cover"
            onError={handleError}
            loading="lazy"
          />
        ) : (
          <span className="text-sm font-medium">
            {getInitial()}
          </span>
        )}
      </div>
    </div>
  );
}