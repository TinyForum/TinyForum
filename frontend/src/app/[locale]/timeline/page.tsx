// app/[locale]/timeline/page.tsx
'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuthStore } from '@/store/auth';
import { timelineApi } from '@/lib/api/modules/timeline';
import { toast } from 'react-hot-toast';
import {
  UserCircleIcon,
  HeartIcon,
  ChatBubbleLeftRightIcon,
  EyeIcon,
  CalendarIcon,
  TagIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import type { TimelineEvent, Subscription, User, Post, Comment } from '@/lib/api/types';

// 时间线事件类型图标映射
const eventTypeConfig: Record<string, { icon: string; label: string; color: string }> = {
  create_post: { icon: '📝', label: '发布了新帖子', color: 'text-blue-500' },
  create_answer: { icon: '💬', label: '回答了问题', color: 'text-green-500' },
  like_post: { icon: '❤️', label: '点赞了帖子', color: 'text-red-500' },
  like_answer: { icon: '❤️', label: '点赞了回答', color: 'text-red-500' },
  follow_user: { icon: '➕', label: '关注了', color: 'text-purple-500' },
  accept_answer: { icon: '✓', label: '采纳了答案', color: 'text-yellow-500' },
  reward_question: { icon: '💰', label: '获得了悬赏', color: 'text-orange-500' },
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

// 时间线事件卡片组件
function TimelineEventCard({ event, currentUserId }: { event: TimelineEvent; currentUserId?: number }) {
  const [liked, setLiked] = useState(false);
  const payload = parsePayload(event.payload);
  const config = eventTypeConfig[event.action] || { icon: '📄', label: event.action, color: 'text-gray-500' };

  const handleLike = async () => {
    toast.success('功能开发中');
  };

  // 构建目标链接
  const getTargetUrl = () => {
    if (event.target_type === 'post') {
      return `/posts/${event.target_id}`;
    }
    if (event.target_type === 'comment') {
      return `/posts/${event.target_id}`;
    }
    if (event.target_type === 'user') {
      return `/users/${event.target_id}`;
    }
    return '#';
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
      <div className="flex gap-4">
        {/* 用户头像 */}
        <Link href={`/users/${event.actor_id}`} className="flex-shrink-0">
          {event.actor?.avatar ? (
            <img
              src={event.actor.avatar}
              alt={event.actor.username}
              className="w-12 h-12 rounded-full object-cover"
            />
          ) : (
            <UserCircleIcon className="w-12 h-12 text-gray-400" />
          )}
        </Link>

        <div className="flex-1">
          {/* 事件头部 */}
          <div className="flex items-center gap-2 mb-2 flex-wrap">
            <Link
              href={`/users/${event.actor_id}`}
              className="font-semibold text-gray-900 hover:text-indigo-600"
            >
              {event.actor?.username || `用户${event.actor_id}`}
            </Link>
            <span className="text-gray-500">{config.label}</span>
            <span className="text-xl">{config.icon}</span>
          </div>

          {/* 事件内容 */}
          {payload.title && (
            <div className="mb-3">
              <Link href={getTargetUrl()} className="block">
                <div className="bg-gray-50 rounded-lg p-4 hover:bg-gray-100 transition-colors">
                  <h4 className="font-medium text-gray-900 mb-2">
                    {payload.title}
                  </h4>
                  {payload.summary && (
                    <p className="text-gray-600 text-sm line-clamp-2">
                      {payload.summary}
                    </p>
                  )}
                  {payload.content && !payload.summary && (
                    <p className="text-gray-600 text-sm line-clamp-2">
                      {payload.content}
                    </p>
                  )}
                </div>
              </Link>
            </div>
          )}

          {/* 事件元信息 */}
          <div className="flex items-center gap-4 text-sm text-gray-400">
            <div className="flex items-center gap-1">
              <CalendarIcon className="w-4 h-4" />
              {new Date(event.created_at).toLocaleDateString()}
            </div>
            
            {event.target_type === 'post' && (
              <>
                <div className="flex items-center gap-1">
                  <EyeIcon className="w-4 h-4" />
                  {payload.view_count || 0}
                </div>
                <div className="flex items-center gap-1">
                  <ChatBubbleLeftRightIcon className="w-4 h-4" />
                  {payload.comment_count || 0}
                </div>
              </>
            )}

            {event.target_type === 'comment' && (
              <div className="flex items-center gap-1">
                <div className="w-2 h-2 bg-green-500 rounded-full" />
                回答
              </div>
            )}
          </div>

          {/* 操作按钮 */}
          {event.target_type === 'post' && (
            <div className="flex items-center gap-4 mt-3 pt-3 border-t">
              <button
                onClick={handleLike}
                className="flex items-center gap-1 text-gray-500 hover:text-red-500 transition-colors"
              >
                {liked ? (
                  <HeartSolidIcon className="w-5 h-5 text-red-500" />
                ) : (
                  <HeartIcon className="w-5 h-5" />
                )}
                <span className="text-sm">{payload.like_count || 0}</span>
              </button>
              <Link
                href={getTargetUrl()}
                className="flex items-center gap-1 text-gray-500 hover:text-indigo-600 transition-colors"
              >
                <ChatBubbleLeftRightIcon className="w-5 h-5" />
                <span className="text-sm">查看详情</span>
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

// 订阅用户卡片组件
interface SubscribedUser {
  id: number;
  username: string;
  avatar: string;
  bio?: string;
}

function SubscribeCard({ user, onUnsubscribe }: { user: SubscribedUser; onUnsubscribe: (userId: number) => void }) {
  return (
    <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
      <Link href={`/users/${user.id}`} className="flex items-center gap-3 flex-1">
        {user.avatar ? (
          <img src={user.avatar} alt={user.username} className="w-10 h-10 rounded-full object-cover" />
        ) : (
          <UserCircleIcon className="w-10 h-10 text-gray-400" />
        )}
        <div>
          <div className="font-medium text-gray-900">{user.username}</div>
          <div className="text-sm text-gray-500">{user.bio || '这个人很懒，什么都没写'}</div>
        </div>
      </Link>
      <button
        onClick={() => onUnsubscribe(user.id)}
        className="px-3 py-1 text-sm text-red-600 hover:bg-red-50 rounded-md transition-colors"
      >
        取消关注
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
  const [activeTab, setActiveTab] = useState<'home' | 'following'>('home');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showSubscriptions, setShowSubscriptions] = useState(false);
  const pageSize = 20;

  // 加载时间线
  const loadTimeline = async () => {
    setLoading(true);
    try {
      const response = activeTab === 'home'
        ? await timelineApi.getHome({ page, page_size: pageSize })
        : await timelineApi.getFollowing({ page, page_size: pageSize });

      if (response.data.code === 200) {
        const { list, total: totalCount } = response.data.data;
        setEvents(list || []);
        setTotal(totalCount || 0);
      } else {
        toast.error(response.data.message || '加载失败');
      }
    } catch (error: any) {
      console.error('Failed to load timeline:', error);
      toast.error(error.response?.data?.message || '加载失败');
    } finally {
      setLoading(false);
    }
  };

  // 加载订阅列表
  const loadSubscriptions = async () => {
    try {
      const response = await timelineApi.getSubscriptions();
      if (response.data.code === 200) {
        setSubscriptions(response.data.data || []);
      }
    } catch (error) {
      console.error('Failed to load subscriptions:', error);
    }
  };

  // 取消关注
  const handleUnsubscribe = async (userId: number) => {
    try {
      const response = await timelineApi.unsubscribe(userId);
      if (response.data.code === 200) {
        toast.success('已取消关注');
        await loadSubscriptions();
        if (activeTab === 'following') {
          await loadTimeline();
        }
      } else {
        toast.error(response.data.message || '操作失败');
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '操作失败');
    }
  };

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login?redirect=/timeline');
      return;
    }
    loadTimeline();
    loadSubscriptions();
  }, [activeTab, page]);

  if (!isAuthenticated) {
    return null;
  }

  const totalPages = Math.ceil(total / pageSize);

  // 转换订阅数据为用户格式
  const subscribedUsers = subscriptions.map(sub => ({
    id: sub.target_user_id,
    username: `用户${sub.target_user_id}`, // 实际应该从 API 获取用户信息
    avatar: '',
  }));

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-3xl mx-auto px-4">
        {/* 头部 */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">时间线</h1>
          <p className="text-gray-500">关注你感兴趣的人和内容</p>
        </div>

        {/* Tab 切换 */}
        <div className="bg-white rounded-lg shadow-sm mb-6">
          <div className="flex border-b">
            <button
              onClick={() => {
                setActiveTab('home');
                setPage(1);
              }}
              className={`flex-1 py-3 text-center font-medium transition-colors ${
                activeTab === 'home'
                  ? 'text-indigo-600 border-b-2 border-indigo-600'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              推荐
            </button>
            <button
              onClick={() => {
                setActiveTab('following');
                setPage(1);
              }}
              className={`flex-1 py-3 text-center font-medium transition-colors ${
                activeTab === 'following'
                  ? 'text-indigo-600 border-b-2 border-indigo-600'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              关注
            </button>
          </div>
        </div>

        {/* 订阅管理 */}
        <div className="mb-6">
          <button
            onClick={() => setShowSubscriptions(!showSubscriptions)}
            className="flex items-center gap-2 text-gray-600 hover:text-indigo-600 transition-colors"
          >
            <span className="font-medium">我关注的人 ({subscriptions.length})</span>
            <svg
              className={`w-4 h-4 transition-transform ${showSubscriptions ? 'rotate-180' : ''}`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          
          {showSubscriptions && (
            <div className="mt-3 space-y-2">
              {subscriptions.length === 0 ? (
                <div className="bg-white rounded-lg p-6 text-center text-gray-500">
                  <p>还没有关注任何人</p>
                  <Link href="/explore" className="text-indigo-600 hover:underline mt-2 inline-block">
                    发现更多用户 →
                  </Link>
                </div>
              ) : (
                subscriptions.map((sub) => (
                  <SubscribeCard
                    key={sub.id}
                    user={{
                      id: sub.target_user_id,
                      username: `用户${sub.target_user_id}`,
                      avatar: '',
                    }}
                    onUnsubscribe={handleUnsubscribe}
                  />
                ))
              )}
            </div>
          )}
        </div>

        {/* 时间线内容 */}
        {loading ? (
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-white rounded-lg shadow-sm p-6 animate-pulse">
                <div className="flex gap-4">
                  <div className="w-12 h-12 bg-gray-200 rounded-full" />
                  <div className="flex-1">
                    <div className="h-4 bg-gray-200 rounded w-1/4 mb-2" />
                    <div className="h-3 bg-gray-200 rounded w-3/4" />
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : events.length === 0 ? (
          <div className="bg-white rounded-lg shadow-sm p-12 text-center">
            <div className="text-6xl mb-4">📭</div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">暂无动态</h3>
            <p className="text-gray-500 mb-4">
              {activeTab === 'home' 
                ? '还没有任何动态，去探索更多内容吧！'
                : '关注用户后，他们的动态会显示在这里'}
            </p>
            {activeTab === 'following' && (
              <Link
                href="/explore"
                className="inline-block px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
              >
                发现用户
              </Link>
            )}
          </div>
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
              <div className="flex justify-center gap-2 mt-6">
                <button
                  onClick={() => setPage(p => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
                >
                  上一页
                </button>
                <span className="px-3 py-1 text-gray-600">
                  第 {page} / {totalPages} 页
                </span>
                <button
                  onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                  disabled={page >= totalPages}
                  className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
                >
                  下一页
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}