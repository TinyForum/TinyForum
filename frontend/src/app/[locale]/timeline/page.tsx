// app/[locale]/timeline/page.tsx
"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { timelineApi } from "@/lib/api/modules/timeline";
import { toast } from "react-hot-toast";
import {
  UserCircleIcon,
  HeartIcon,
  ChatBubbleLeftRightIcon,
  EyeIcon,
  CalendarIcon,
  TagIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  UsersIcon,
  SparklesIcon,
  UserPlusIcon,
  HomeIcon,
} from "@heroicons/react/24/outline";
import { HeartIcon as HeartSolidIcon } from "@heroicons/react/24/solid";
import type {
  TimelineEvent,
  Subscription,
  User,
  Post,
  Comment,
} from "@/lib/api/types";

// 事件类型配置 - 使用主题色
const eventTypeConfig: Record<
  string,
  { icon: string; label: string; color: string; bgColor: string }
> = {
  create_post: {
    icon: "📝",
    label: "发布了新帖子",
    color: "text-primary",
    bgColor: "bg-primary/10",
  },
  create_answer: {
    icon: "💬",
    label: "回答了问题",
    color: "text-secondary",
    bgColor: "bg-secondary/10",
  },
  like_post: {
    icon: "❤️",
    label: "点赞了帖子",
    color: "text-error",
    bgColor: "bg-error/10",
  },
  like_answer: {
    icon: "❤️",
    label: "点赞了回答",
    color: "text-error",
    bgColor: "bg-error/10",
  },
  follow_user: {
    icon: "➕",
    label: "关注了",
    color: "text-accent",
    bgColor: "bg-accent/10",
  },
  accept_answer: {
    icon: "✓",
    label: "采纳了答案",
    color: "text-warning",
    bgColor: "bg-warning/10",
  },
  reward_question: {
    icon: "💰",
    label: "获得了悬赏",
    color: "text-warning",
    bgColor: "bg-warning/10",
  },
};

// 解析事件负载
interface EventPayload {
  title?: string;
  content?: string;
  summary?: string;
  url?: string;
  view_count?: number;
  like_count?: number;
  comment_count?: number;
}

function parsePayload(payload: string): EventPayload {
  try {
    return JSON.parse(payload);
  } catch {
    return {};
  }
}

// 加载骨架屏
function LoadingSkeleton() {
  return (
    <div className="space-y-4">
      {[1, 2, 3].map((i) => (
        <div key={i} className="card bg-base-100 shadow-sm p-6 animate-pulse">
          <div className="flex gap-4">
            <div className="w-12 h-12 bg-base-200 rounded-full" />
            <div className="flex-1 space-y-3">
              <div className="h-4 bg-base-200 rounded w-1/4" />
              <div className="h-3 bg-base-200 rounded w-3/4" />
              <div className="h-20 bg-base-200 rounded-lg" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

// 时间线事件卡片组件
function TimelineEventCard({
  event,
  currentUserId,
}: {
  event: TimelineEvent;
  currentUserId?: number;
}) {
  const [liked, setLiked] = useState(false);
  const payload = parsePayload(event.payload);
  const config = eventTypeConfig[event.action] || {
    icon: "📄",
    label: event.action,
    color: "text-base-content/60",
    bgColor: "bg-base-200",
  };

  const handleLike = async () => {
    toast.success("功能开发中");
  };

  const getTargetUrl = () => {
    if (event.target_type === "post") return `/posts/${event.target_id}`;
    if (event.target_type === "comment") return `/posts/${event.target_id}`;
    if (event.target_type === "user") return `/users/${event.target_id}`;
    return "#";
  };

  return (
    <div className="card bg-base-100 shadow-sm hover:shadow-md transition-all duration-300 border border-base-200 hover:border-primary/20">
      <div className="card-body p-5">
        <div className="flex gap-4">
          {/* 用户头像 */}
          <Link href={`/users/${event.actor_id}`} className="flex-shrink-0">
            {event.actor?.avatar ? (
              <img
                src={event.actor.avatar}
                alt={event.actor.username}
                className="w-12 h-12 rounded-full object-cover ring-2 ring-primary/20"
              />
            ) : (
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center">
                <UserCircleIcon className="w-8 h-8 text-primary/60" />
              </div>
            )}
          </Link>

          <div className="flex-1 min-w-0">
            {/* 事件头部 */}
            <div className="flex items-center gap-2 mb-3 flex-wrap">
              <Link
                href={`/users/${event.actor_id}`}
                className="font-semibold text-base-content hover:text-primary transition-colors"
              >
                {event.actor?.username || `用户${event.actor_id}`}
              </Link>
              <span className="text-base-content/50">{config.label}</span>
              <span className="text-lg">{config.icon}</span>
            </div>

            {/* 事件内容 */}
            {payload.title && (
              <div className="mb-4">
                <Link href={getTargetUrl()} className="block group">
                  <div className="bg-base-200/50 rounded-xl p-4 hover:bg-base-200 transition-all duration-200 border border-base-200 group-hover:border-primary/20">
                    <h4 className="font-medium text-base-content mb-2 group-hover:text-primary transition-colors">
                      {payload.title}
                    </h4>
                    {(payload.summary || payload.content) && (
                      <p className="text-base-content/60 text-sm line-clamp-2">
                        {payload.summary || payload.content}
                      </p>
                    )}
                  </div>
                </Link>
              </div>
            )}

            {/* 事件元信息 */}
            <div className="flex flex-wrap items-center gap-4 text-sm text-base-content/40">
              <div className="flex items-center gap-1.5">
                <CalendarIcon className="w-3.5 h-3.5" />
                <span>{new Date(event.created_at).toLocaleDateString()}</span>
              </div>

              {event.target_type === "post" && (
                <>
                  <div className="flex items-center gap-1.5">
                    <EyeIcon className="w-3.5 h-3.5" />
                    <span>{payload.view_count || 0}</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <ChatBubbleLeftRightIcon className="w-3.5 h-3.5" />
                    <span>{payload.comment_count || 0}</span>
                  </div>
                </>
              )}

              {event.target_type === "comment" && (
                <div className="flex items-center gap-1.5">
                  <div className="w-1.5 h-1.5 bg-success rounded-full" />
                  <span>回答</span>
                </div>
              )}
            </div>

            {/* 操作按钮 */}
            {event.target_type === "post" && (
              <div className="flex items-center gap-4 mt-4 pt-3 border-t border-base-200">
                <button
                  onClick={handleLike}
                  className="flex items-center gap-1.5 text-base-content/50 hover:text-error transition-all duration-200"
                >
                  {liked ? (
                    <HeartSolidIcon className="w-4 h-4 text-error" />
                  ) : (
                    <HeartIcon className="w-4 h-4" />
                  )}
                  <span className="text-sm">{payload.like_count || 0}</span>
                </button>
                <Link
                  href={getTargetUrl()}
                  className="flex items-center gap-1.5 text-base-content/50 hover:text-primary transition-all duration-200"
                >
                  <ChatBubbleLeftRightIcon className="w-4 h-4" />
                  <span className="text-sm">查看详情</span>
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

// 订阅用户卡片组件
function SubscribeCard({
  user,
  onUnsubscribe,
}: {
  user: { id: number; username: string; avatar?: string; bio?: string };
  onUnsubscribe: (userId: number) => void;
}) {
  return (
    <div className="flex items-center justify-between p-3 bg-base-200/50 rounded-xl hover:bg-base-200 transition-all duration-200 group">
      <Link
        href={`/users/${user.id}`}
        className="flex items-center gap-3 flex-1 min-w-0"
      >
        {user.avatar ? (
          <img
            src={user.avatar}
            alt={user.username}
            className="w-10 h-10 rounded-full object-cover ring-2 ring-primary/20"
          />
        ) : (
          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center">
            <UserCircleIcon className="w-6 h-6 text-primary/60" />
          </div>
        )}
        <div className="flex-1 min-w-0">
          <div className="font-medium text-base-content truncate">
            {user.username}
          </div>
          {user.bio && (
            <div className="text-xs text-base-content/50 truncate">
              {user.bio}
            </div>
          )}
        </div>
      </Link>
      <button
        onClick={() => onUnsubscribe(user.id)}
        className="px-3 py-1.5 text-sm text-error/70 hover:text-error hover:bg-error/10 rounded-lg transition-all duration-200 opacity-0 group-hover:opacity-100"
      >
        取消关注
      </button>
    </div>
  );
}

// 空状态组件
function EmptyState({ activeTab }: { activeTab: "home" | "following" }) {
  return (
    <div className="card bg-base-100 shadow-sm p-12 text-center border border-base-200">
      <div className="text-6xl mb-4 opacity-50">
        {activeTab === "home" ? "📭" : "👥"}
      </div>
      <h3 className="text-lg font-semibold text-base-content mb-2">
        {activeTab === "home" ? "暂无动态" : "暂无关注动态"}
      </h3>
      <p className="text-base-content/60 mb-4">
        {activeTab === "home"
          ? "还没有任何动态，去探索更多内容吧！"
          : "关注用户后，他们的动态会显示在这里"}
      </p>
      {activeTab === "following" && (
        <Link href="/explore" className="btn btn-primary btn-sm gap-2">
          <UserPlusIcon className="w-4 h-4" />
          发现用户
        </Link>
      )}
    </div>
  );
}

// 分页组件
function Pagination({
  currentPage,
  totalPages,
  onPageChange,
}: {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}) {
  const getPageNumbers = () => {
    const pages: number[] = [];
    const maxVisible = 5;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) pages.push(i);
    } else {
      if (currentPage <= 3) {
        for (let i = 1; i <= maxVisible; i++) pages.push(i);
      } else if (currentPage >= totalPages - 2) {
        for (let i = totalPages - maxVisible + 1; i <= totalPages; i++)
          pages.push(i);
      } else {
        for (let i = currentPage - 2; i <= currentPage + 2; i++) pages.push(i);
      }
    }
    return pages;
  };

  return (
    <div className="flex justify-center items-center gap-2 mt-8">
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="btn btn-ghost btn-sm gap-1"
      >
        上一页
      </button>

      <div className="flex gap-1.5 mx-2">
        {getPageNumbers().map((pageNum) => (
          <button
            key={pageNum}
            onClick={() => onPageChange(pageNum)}
            className={`btn btn-sm min-w-[2.5rem] ${
              currentPage === pageNum ? "btn-primary" : "btn-ghost"
            }`}
          >
            {pageNum}
          </button>
        ))}
      </div>

      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage >= totalPages}
        className="btn btn-ghost btn-sm gap-1"
      >
        下一页
      </button>
    </div>
  );
}

export default function Timeline() {
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();
  const [events, setEvents] = useState<TimelineEvent[]>([]);
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"home" | "following">("home");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showSubscriptions, setShowSubscriptions] = useState(false);
  const pageSize = 20;

  const loadTimeline = async () => {
    setLoading(true);
    try {
      const response =
        activeTab === "home"
          ? await timelineApi.getHome({ page, page_size: pageSize })
          : await timelineApi.getFollowing({ page, page_size: pageSize });

      if (response.data.code === 200) {
        const { list, total: totalCount } = response.data.data;
        setEvents(list || []);
        setTotal(totalCount || 0);
      } else {
        toast.error(response.data.message || "加载失败");
      }
    } catch (error: any) {
      console.error("Failed to load timeline:", error);
      toast.error(error.response?.data?.message || "加载失败");
    } finally {
      setLoading(false);
    }
  };

  const loadSubscriptions = async () => {
    try {
      const response = await timelineApi.getSubscriptions();
      if (response.data.code === 200) {
        setSubscriptions(response.data.data || []);
      }
    } catch (error) {
      console.error("Failed to load subscriptions:", error);
    }
  };

  const handleUnsubscribe = async (userId: number) => {
    try {
      const response = await timelineApi.unsubscribe(userId);
      if (response.data.code === 200) {
        toast.success("已取消关注");
        await loadSubscriptions();
        if (activeTab === "following") {
          await loadTimeline();
        }
      } else {
        toast.error(response.data.message || "操作失败");
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || "操作失败");
    }
  };

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login?redirect=/timeline");
      return;
    }
    loadTimeline();
    loadSubscriptions();
  }, [activeTab, page]);

  if (!isAuthenticated) {
    return null;
  }

  const totalPages = Math.ceil(total / pageSize);
  const subscribedUsers = subscriptions.map((sub) => ({
    id: sub.target_user_id,
    username: `用户${sub.target_user_id}`,
    avatar: "",
  }));

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-3xl mx-auto px-4 py-8 md:py-12">
        {/* 头部区域 */}
        <div className="text-center mb-8">
          <div className="relative inline-block mb-4">
            <div className="absolute inset-0 bg-gradient-to-r from-primary/20 to-secondary/20 rounded-full blur-2xl" />
            <div className="relative bg-gradient-to-br from-primary to-secondary p-3 rounded-full shadow-lg">
              <HomeIcon className="w-8 h-8 text-white" />
            </div>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
            时间线
          </h1>
          <p className="text-base-content/60 mt-2">关注你感兴趣的人和内容</p>
        </div>

        {/* Tab 切换 - 优化样式 */}
        <div className="card bg-base-100 shadow-sm mb-6 border border-base-200">
          <div className="flex p-1">
            <button
              onClick={() => {
                setActiveTab("home");
                setPage(1);
              }}
              className={`flex-1 py-2.5 rounded-lg font-medium transition-all duration-200 ${
                activeTab === "home"
                  ? "bg-primary text-white shadow-sm"
                  : "text-base-content/60 hover:text-base-content hover:bg-base-200"
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <SparklesIcon className="w-4 h-4" />
                推荐
              </div>
            </button>
            <button
              onClick={() => {
                setActiveTab("following");
                setPage(1);
              }}
              className={`flex-1 py-2.5 rounded-lg font-medium transition-all duration-200 ${
                activeTab === "following"
                  ? "bg-primary text-white shadow-sm"
                  : "text-base-content/60 hover:text-base-content hover:bg-base-200"
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <UsersIcon className="w-4 h-4" />
                关注
              </div>
            </button>
          </div>
        </div>

        {/* 订阅管理区域 */}
        {subscriptions.length > 0 && (
          <div className="mb-6">
            <button
              onClick={() => setShowSubscriptions(!showSubscriptions)}
              className="flex items-center gap-2 text-base-content/70 hover:text-primary transition-colors group"
            >
              <span className="font-medium">
                我关注的人 ({subscriptions.length})
              </span>
              {showSubscriptions ? (
                <ChevronUpIcon className="w-4 h-4 group-hover:-translate-y-0.5 transition-transform" />
              ) : (
                <ChevronDownIcon className="w-4 h-4 group-hover:translate-y-0.5 transition-transform" />
              )}
            </button>

            {showSubscriptions && (
              <div className="mt-3 space-y-2 animate-fade-in">
                {subscriptions.map((sub) => (
                  <SubscribeCard
                    key={sub.id}
                    user={{
                      id: sub.target_user_id,
                      username: `用户${sub.target_user_id}`,
                      avatar: "",
                    }}
                    onUnsubscribe={handleUnsubscribe}
                  />
                ))}
              </div>
            )}
          </div>
        )}

        {/* 时间线内容 */}
        {loading ? (
          <LoadingSkeleton />
        ) : events.length === 0 ? (
          <EmptyState activeTab={activeTab} />
        ) : (
          <>
            <div className="space-y-4">
              {events.map((event) => (
                <TimelineEventCard
                  key={event.id}
                  event={event}
                  currentUserId={user?.id}
                />
              ))}
            </div>

            {/* 分页 */}
            {totalPages > 1 && (
              <Pagination
                currentPage={page}
                totalPages={totalPages}
                onPageChange={setPage}
              />
            )}

            {/* 底部提示 */}
            <div className="text-center text-xs text-base-content/40 mt-6 pt-4 border-t border-base-200">
              已加载全部内容
            </div>
          </>
        )}
      </div>
    </div>
  );
}
