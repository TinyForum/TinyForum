import { Topic } from "@/lib/api";
import Link from "next/link";
import {
  ChatBubbleLeftRightIcon,
  UserGroupIcon,
  ArrowRightIcon,
  HashtagIcon,
  FireIcon,
} from "@heroicons/react/24/outline";

interface HotTopicCardProps {
  topic: Topic;
  rank?: number;
  showTodayCount?: boolean;
}

export function HotTopicCard({
  topic,
  rank,
  showTodayCount = true,
}: HotTopicCardProps) {
  // 获取排名样式
  const getRankStyles = (rank: number) => {
    switch (rank) {
      case 1:
        return {
          bg: "bg-yellow-100 dark:bg-yellow-900/20",
          text: "text-yellow-500",
          icon: "🥇",
        };
      case 2:
        return {
          bg: "bg-gray-100 dark:bg-gray-800",
          text: "text-gray-500",
          icon: "🥈",
        };
      case 3:
        return {
          bg: "bg-orange-100 dark:bg-orange-900/20",
          text: "text-orange-500",
          icon: "🥉",
        };
      default:
        return null;
    }
  };

  const rankStyles = rank && rank <= 3 ? getRankStyles(rank) : null;

  return (
    <Link href={`/topics/${topic.id}`} className="block">
      <div className="group card bg-base-100 border border-base-200 hover:border-primary/20 hover:shadow-md transition-all duration-200 cursor-pointer">
        <div className="card-body p-4">
          <div className="flex items-start gap-3">
            {/* 话题图标/排名 */}
            {rankStyles ? (
              <div
                className={`w-8 h-8 rounded-lg flex items-center justify-center shrink-0 ${rankStyles.bg} ${rankStyles.text}`}
              >
                <span className="text-sm font-bold">{rankStyles.icon}</span>
              </div>
            ) : (
              <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-red-100 to-red-200 dark:from-red-900/30 dark:to-red-800/20 flex items-center justify-center shrink-0 group-hover:scale-110 transition-transform">
                <HashtagIcon className="w-4 h-4 text-red-500" />
              </div>
            )}

            {/* 话题内容 */}
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <h4 className="font-semibold text-base-content group-hover:text-primary transition-colors line-clamp-1">
                  {topic.title}
                </h4>
                {topic.post_count !== undefined && topic.post_count > 50 && (
                  <span className="badge badge-primary badge-xs shrink-0">
                    热门
                  </span>
                )}
              </div>
              {topic.description && (
                <p className="text-sm text-base-content/60 line-clamp-2 mb-2">
                  {topic.description}
                </p>
              )}
              <div className="flex items-center flex-wrap gap-x-3 gap-y-1 text-xs text-base-content/50">
                <span className="flex items-center gap-1">
                  <ChatBubbleLeftRightIcon className="w-3 h-3" />
                  {topic.post_count || 0} 帖子
                </span>
                <span className="flex items-center gap-1">
                  <UserGroupIcon className="w-3 h-3" />
                  {topic.follower_count || 0} 关注
                </span>
                {/* {showTodayCount && topic.today_count !== undefined && topic.today_count > 0 && (
                  <span className="flex items-center gap-1 text-orange-500">
                    <FireIcon className="w-3 h-3" />
                    今日 {topic.today_count}
                  </span>
                )} */}
              </div>
            </div>

            {/* 箭头指示 */}
            <div className="opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
              <ArrowRightIcon className="w-4 h-4 text-base-content/40 group-hover:text-primary group-hover:translate-x-1 transition-all" />
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}
