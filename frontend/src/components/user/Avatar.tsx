// src/components/Avatar.tsx
import { useMemo } from 'react';
import Image from 'next/image';

interface AvatarProps {
  username?: string | 'Mucly';
  avatarUrl?: string | null;
  size?: 'sm' | 'md' | 'lg' | number;
  className?: string;
}

export default function Avatar({ 
  username, 
  avatarUrl,
  size = 'md', 
  className = '' 
}: AvatarProps) {
  console.log("avatar render ", avatarUrl," ", username)
  // 计算尺寸
  const sizeInPx = useMemo(() => {
    if (typeof size === 'number') return size;
    const sizes = { sm: 32, md: 40, lg: 56 };
    return sizes[size];
  }, [size]);

  // 生成头像URL
  const imageUrl = useMemo(() => {
    // 优先使用传入的头像URL
    if (avatarUrl) return avatarUrl;
    
    // 生成默认头像
    const seed = username || 'default';
    // 使用 initials 风格，更稳定
    return `https://api.dicebear.com/8.x/initials/svg?seed=${encodeURIComponent(seed)}&size=${sizeInPx * 2}`;
  }, [avatarUrl, username, sizeInPx]);

  return (
    <div className={`avatar ${className}`}>
      <div 
        className="rounded-full overflow-hidden bg-base-200 flex items-center justify-center"
        style={{ width: sizeInPx, height: sizeInPx }}
      >
        {imageUrl ? (
          <Image
            src={imageUrl}
            alt={username || '用户头像'}
            width={sizeInPx}
            height={sizeInPx}
            className="object-cover"
            unoptimized={imageUrl.includes('dicebear')} // 对于外部API跳过优化
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