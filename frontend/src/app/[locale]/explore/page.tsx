// app/[locale]/explore/page.tsx
'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { toast } from 'react-hot-toast';
import {
  MagnifyingGlassIcon,
  FireIcon,
  ClockIcon,
  ChatBubbleLeftRightIcon,
  EyeIcon,
  HeartIcon,
  UserGroupIcon,
  HashtagIcon,
  NewspaperIcon,
  ArrowTrendingUpIcon,
  SparklesIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import { postApi, tagApi, topicApi, userApi } from '@/lib/api';
import type { Post, Tag, Topic, User } from '@/lib/api/types';

// 热门帖子卡片
function HotPostCard({ post, rank }: { post: Post; rank: number }) {
  const [liked, setLiked] = useState(false);
  const [likesCount, setLikesCount] = useState(post.like_count || 0);

  const handleLike = async () => {
    if (!post.id) return;
    try {
      if (liked) {
        await postApi.unlike(post.id);
        setLiked(false);
        setLikesCount(prev => prev - 1);
      } else {
        await postApi.like(post.id);
        setLiked(true);
        setLikesCount(prev => prev + 1);
      }
    } catch (error) {
      console.error('Like failed:', error);
      toast.error('操作失败');
    }
  };

  const rankColors: Record<number, string> = {
    1: 'text-yellow-500',
    2: 'text-gray-400',
    3: 'text-amber-600',
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow">
      <div className="flex items-start gap-3">
        {rank > 0 && (
          <div className={`text-2xl font-bold w-8 ${rankColors[rank] || 'text-gray-300'}`}>
            {rank}
          </div>
        )}
        <div className="flex-1">
          <Link href={`/posts/${post.id}`}>
            <h3 className="font-semibold text-gray-900 hover:text-indigo-600 mb-2 line-clamp-1">
              {post.title}
            </h3>
          </Link>
          <div className="flex items-center gap-3 text-xs text-gray-400">
            <div className="flex items-center gap-1">
              <EyeIcon className="w-3 h-3" />
              {post.view_count}
            </div>
            <div className="flex items-center gap-1">
              <ChatBubbleLeftRightIcon className="w-3 h-3" />
              {post.question?.answer_count || 0}
            </div>
            <button onClick={handleLike} className="flex items-center gap-1 hover:text-red-500">
              {liked ? (
                <HeartSolidIcon className="w-3 h-3 text-red-500" />
              ) : (
                <HeartIcon className="w-3 h-3" />
              )}
              {likesCount}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

// 热门标签卡片
function HotTagCard({ tag }: { tag: Tag }) {
  return (
    <Link href={`/questions?tag_id=${tag.id}`}>
      <div className="bg-white rounded-lg shadow-sm p-3 hover:shadow-md transition-shadow cursor-pointer">
        <div
          className="inline-block px-2 py-1 rounded-md text-xs font-medium mb-2"
          style={{ backgroundColor: `${tag.color}20`, color: tag.color }}
        >
          {tag.name}
        </div>
        <p className="text-xs text-gray-500 line-clamp-2">{tag.description || '暂无描述'}</p>
        <div className="mt-2 text-xs text-gray-400">{tag.post_count} 个帖子</div>
      </div>
    </Link>
  );
}

// 热门话题卡片
function HotTopicCard({ topic }: { topic: Topic }) {
  return (
    <Link href={`/topics/${topic.id}`}>
      <div className="bg-white rounded-lg shadow-sm p-3 hover:shadow-md transition-shadow cursor-pointer">
        <h4 className="font-medium text-gray-900 mb-1 line-clamp-1">{topic.title}</h4>
        <p className="text-xs text-gray-500 line-clamp-2">{topic.description}</p>
        <div className="mt-2 flex items-center gap-2 text-xs text-gray-400">
          <span>{topic.post_count} 帖子</span>
          <span>{topic.follower_count} 关注</span>
        </div>
      </div>
    </Link>
  );
}

// 活跃用户卡片
function ActiveUserCard({ user }: { user: User }) {
  return (
    <Link href={`/users/${user.id}`}>
      <div className="bg-white rounded-lg shadow-sm p-3 hover:shadow-md transition-shadow cursor-pointer flex items-center gap-3">
        {user.avatar ? (
          <img src={user.avatar} alt={user.username} className="w-10 h-10 rounded-full object-cover" />
        ) : (
          <div className="w-10 h-10 bg-indigo-100 rounded-full flex items-center justify-center">
            <span className="text-indigo-600 font-medium">
              {user.username.charAt(0).toUpperCase()}
            </span>
          </div>
        )}
        <div className="flex-1 min-w-0">
          <div className="font-medium text-gray-900 truncate">{user.username}</div>
          <div className="text-xs text-gray-400">积分: {user.score}</div>
        </div>
      </div>
    </Link>
  );
}

// 分类 Tab
const exploreTabs = [
  { id: 'hot', label: '热门', icon: FireIcon, sortBy: 'hot' },
  { id: 'latest', label: '最新', icon: ClockIcon, sortBy: 'latest' },
  { id: 'trending', label: '趋势', icon: ArrowTrendingUpIcon, sortBy: 'views' },
  { id: 'recommended', label: '推荐', icon: SparklesIcon, sortBy: 'recommended' },
];

export default function Explore() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [activeTab, setActiveTab] = useState('hot');
  const [posts, setPosts] = useState<Post[]>([]);
  const [hotTags, setHotTags] = useState<Tag[]>([]);
  const [hotTopics, setHotTopics] = useState<Topic[]>([]);
  const [activeUsers, setActiveUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [searchResults, setSearchResults] = useState<Post[]>([]);
  const [searching, setSearching] = useState(false);

  // 加载探索数据
  const loadExploreData = async () => {
    setLoading(true);
    try {
      // 获取当前 tab 对应的帖子
      const currentTab = exploreTabs.find(tab => tab.id === activeTab);
      const postsResponse = await postApi.list({
        page: 1,
        page_size: 10,
        sort_by: currentTab?.sortBy,
      });

      // 获取热门标签
      const tagsResponse = await tagApi.list();

      // 获取热门话题
      const topicsResponse = await topicApi.list({ page: 1, page_size: 8 });

      // 获取活跃用户排行榜
      const usersResponse = await userApi.leaderboard(8);

      if (postsResponse.data.code === 200) {
        setPosts(postsResponse.data.data.items || []);
      }
      if (tagsResponse.data.code === 200) {
        // 按帖子数量排序，取前12个
        const sortedTags = [...(tagsResponse.data.data || [])].sort(
          (a, b) => (b.post_count || 0) - (a.post_count || 0)
        );
        setHotTags(sortedTags.slice(0, 12));
      }
      if (topicsResponse.data.code === 200) {
        setHotTopics(topicsResponse.data.data.items || []);
      }
      if (usersResponse.data.code === 200) {
        setActiveUsers(usersResponse.data.data || []);
      }
    } catch (error) {
      console.error('Failed to load explore data:', error);
      toast.error('加载失败');
    } finally {
      setLoading(false);
    }
  };

  // 搜索
  const handleSearch = async () => {
    if (!searchKeyword.trim()) {
      toast.error('请输入搜索关键词');
      return;
    }

    setSearching(true);
    try {
      const response = await postApi.list({
        keyword: searchKeyword,
        page: 1,
        page_size: 20,
      });
      if (response.data.code === 200) {
        setSearchResults(response.data.data.items || []);
        if (response.data.data.items?.length === 0) {
          toast('未找到相关结果');
        }
      }
    } catch (error) {
      console.error('Search failed:', error);
      toast.error('搜索失败');
    } finally {
      setSearching(false);
    }
  };

  // 清除搜索
  const clearSearch = () => {
    setSearchKeyword('');
    setSearchResults([]);
  };

  useEffect(() => {
    loadExploreData();
  }, [activeTab]);

  // 获取当前显示的帖子
  const displayPosts = searchKeyword ? searchResults : posts;

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        {/* 头部 */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">探索</h1>
          <p className="text-gray-500">发现热门内容、话题和活跃用户</p>
        </div>

        {/* 搜索框 */}
        <div className="mb-8">
          <div className="flex gap-2">
            <div className="flex-1 relative">
              <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={searchKeyword}
                onChange={(e) => setSearchKeyword(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="搜索帖子..."
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
              />
            </div>
            <button
              onClick={handleSearch}
              disabled={searching}
              className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50"
            >
              {searching ? '搜索中...' : '搜索'}
            </button>
            {searchKeyword && (
              <button
                onClick={clearSearch}
                className="px-4 py-2 text-gray-600 hover:text-gray-800 rounded-lg hover:bg-gray-100 transition-colors"
              >
                清除
              </button>
            )}
          </div>
        </div>

        {/* 内容区域 */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 左侧：帖子列表 */}
          <div className="lg:col-span-2">
            {/* Tab 切换 */}
            {!searchKeyword && (
              <div className="bg-white rounded-lg shadow-sm mb-4">
                <div className="flex border-b">
                  {exploreTabs.map((tab) => {
                    const Icon = tab.icon;
                    return (
                      <button
                        key={tab.id}
                        onClick={() => setActiveTab(tab.id)}
                        className={`flex-1 py-3 text-center font-medium transition-colors flex items-center justify-center gap-2 ${
                          activeTab === tab.id
                            ? 'text-indigo-600 border-b-2 border-indigo-600'
                            : 'text-gray-500 hover:text-gray-700'
                        }`}
                      >
                        <Icon className="w-4 h-4" />
                        {tab.label}
                      </button>
                    );
                  })}
                </div>
              </div>
            )}

            {/* 帖子列表 */}
            {loading && !searchKeyword ? (
              <div className="space-y-3">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div key={i} className="bg-white rounded-lg shadow-sm p-4 animate-pulse">
                    <div className="h-5 bg-gray-200 rounded w-3/4 mb-2" />
                    <div className="h-4 bg-gray-200 rounded w-1/2" />
                  </div>
                ))}
              </div>
            ) : displayPosts.length === 0 ? (
              <div className="bg-white rounded-lg shadow-sm p-12 text-center">
                <div className="text-6xl mb-4">🔍</div>
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  {searchKeyword ? '未找到相关结果' : '暂无内容'}
                </h3>
                <p className="text-gray-500">
                  {searchKeyword ? '尝试使用其他关键词' : '还没有内容，去发布第一个帖子吧！'}
                </p>
                {!searchKeyword && (
                  <Link
                    href="/questions/ask"
                    className="inline-block mt-4 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
                  >
                    发布问题
                  </Link>
                )}
              </div>
            ) : (
              <div className="space-y-3">
                {displayPosts.map((post, index) => (
                  <HotPostCard
                    key={post.id}
                    post={post}
                    rank={activeTab === 'hot' && !searchKeyword ? index + 1 : 0}
                  />
                ))}
              </div>
            )}
          </div>

          {/* 右侧：侧边栏 */}
          <div className="space-y-6">
            {/* 热门标签 */}
            <div className="bg-white rounded-lg shadow-sm p-4">
              <div className="flex items-center gap-2 mb-3">
                <HashtagIcon className="w-5 h-5 text-indigo-500" />
                <h2 className="font-semibold text-gray-900">热门标签</h2>
              </div>
              {loading ? (
                <div className="space-y-2">
                  {[1, 2, 3, 4].map((i) => (
                    <div key={i} className="h-16 bg-gray-100 rounded animate-pulse" />
                  ))}
                </div>
              ) : hotTags.length === 0 ? (
                <p className="text-sm text-gray-400 text-center py-4">暂无标签</p>
              ) : (
                <div className="grid grid-cols-2 gap-2">
                  {hotTags.slice(0, 8).map((tag) => (
                    <HotTagCard key={tag.id} tag={tag} />
                  ))}
                </div>
              )}
            </div>

            {/* 热门话题 */}
            <div className="bg-white rounded-lg shadow-sm p-4">
              <div className="flex items-center gap-2 mb-3">
                <NewspaperIcon className="w-5 h-5 text-green-500" />
                <h2 className="font-semibold text-gray-900">热门话题</h2>
              </div>
              {loading ? (
                <div className="space-y-2">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="h-20 bg-gray-100 rounded animate-pulse" />
                  ))}
                </div>
              ) : hotTopics.length === 0 ? (
                <p className="text-sm text-gray-400 text-center py-4">暂无话题</p>
              ) : (
                <div className="space-y-2">
                  {hotTopics.slice(0, 5).map((topic) => (
                    <HotTopicCard key={topic.id} topic={topic} />
                  ))}
                </div>
              )}
              <Link
                href="/topics"
                className="block text-center text-sm text-indigo-600 hover:text-indigo-700 mt-3"
              >
                查看更多话题 →
              </Link>
            </div>

            {/* 活跃用户 */}
            <div className="bg-white rounded-lg shadow-sm p-4">
              <div className="flex items-center gap-2 mb-3">
                <UserGroupIcon className="w-5 h-5 text-purple-500" />
                <h2 className="font-semibold text-gray-900">积分榜</h2>
              </div>
              {loading ? (
                <div className="space-y-2">
                  {[1, 2, 3, 4].map((i) => (
                    <div key={i} className="h-14 bg-gray-100 rounded animate-pulse" />
                  ))}
                </div>
              ) : activeUsers.length === 0 ? (
                <p className="text-sm text-gray-400 text-center py-4">暂无用户</p>
              ) : (
                <div className="space-y-2">
                  {activeUsers.slice(0, 6).map((user, index) => (
                    <div key={user.id} className="flex items-center gap-2">
                      <div className="w-6 text-sm font-medium text-gray-400">#{index + 1}</div>
                      <ActiveUserCard user={user} />
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* 关于 */}
            <div className="bg-gradient-to-r from-indigo-50 to-purple-50 rounded-lg p-4">
              <h3 className="font-semibold text-gray-900 mb-2">发现更多</h3>
              <p className="text-sm text-gray-600 mb-3">
                关注热门话题、标签和活跃用户，发现更多精彩内容
              </p>
              {!isAuthenticated && (
                <Link
                  href="/login"
                  className="block text-center text-sm text-indigo-600 hover:text-indigo-700"
                >
                  登录后获得个性化推荐 →
                </Link>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}