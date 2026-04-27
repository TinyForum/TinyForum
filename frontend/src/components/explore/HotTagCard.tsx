import { Tag } from "@/lib/api";
import Link from "next/link";
import {
  HashtagIcon,
  DocumentTextIcon,
  ArrowRightIcon,
} from "@heroicons/react/24/outline";

interface HotTagCardProps {
  tag: Tag;
  size?: "sm" | "md" | "lg";
  showStats?: boolean;
}

export function HotTagCard({
  tag,
  size = "md",
  showStats = true,
}: HotTagCardProps) {
  // 尺寸配置
  const sizeClasses = {
    sm: {
      container: "p-2",
      badge: "px-1.5 py-0.5 text-xs",
      icon: "w-3 h-3",
      title: "text-xs",
      stats: "text-[10px]",
    },
    md: {
      container: "p-3",
      badge: "px-2 py-1 text-xs",
      icon: "w-3.5 h-3.5",
      title: "text-sm",
      stats: "text-xs",
    },
    lg: {
      container: "p-4",
      badge: "px-3 py-1.5 text-sm",
      icon: "w-4 h-4",
      title: "text-base",
      stats: "text-sm",
    },
  };

  const classes = sizeClasses[size];

  // 生成标签颜色（确保可读性）
  const getTagColor = (color: string) => {
    // 如果是红色主题，优先使用主题色
    if (!color || color === "#000000") {
      return {
        bg: "bg-red-100 dark:bg-red-900/30",
        text: "text-red-600 dark:text-red-400",
        border: "border-red-200 dark:border-red-800/50",
      };
    }
    return {
      bg: `${color}10`,
      text: color,
      border: `${color}20`,
    };
  };

  const tagColor = getTagColor(tag.color);

  return (
    <Link href={`/questions?tag_id=${tag.id}`} className="block">
      <div
        className={`group card bg-base-100 border border-base-200 hover:border-primary/20 hover:shadow-md transition-all duration-200 cursor-pointer ${classes.container}`}
      >
        <div className="flex items-start justify-between gap-2">
          <div className="flex-1 min-w-0">
            {/* 标签徽章 */}
            <div className="flex items-center gap-2 mb-2">
              <div
                className={`inline-flex items-center gap-1 rounded-lg font-medium ${tagColor.bg} ${tagColor.text} ${classes.badge}`}
              >
                <HashtagIcon className={classes.icon} />
                <span className="font-mono">{tag.name}</span>
              </div>

              {/* 热门标识 */}
              {tag.post_count !== undefined && tag.post_count > 100 && (
                <span className="badge badge-primary badge-xs">热门</span>
              )}
            </div>

            {/* 描述 */}
            {tag.description && (
              <p
                className={`text-base-content/60 line-clamp-2 mb-2 ${classes.title}`}
              >
                {tag.description}
              </p>
            )}

            {/* 统计信息 */}
            {showStats && tag.post_count !== undefined && (
              <div
                className={`flex items-center gap-1 text-base-content/40 ${classes.stats}`}
              >
                <DocumentTextIcon className={`${classes.icon} opacity-60`} />
                <span>{tag.post_count} 个帖子</span>
              </div>
            )}
          </div>

          {/* 箭头指示 */}
          <div className="opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
            <ArrowRightIcon className="w-3.5 h-3.5 text-base-content/40 group-hover:text-primary group-hover:translate-x-1 transition-all" />
          </div>
        </div>
      </div>
    </Link>
  );
}
