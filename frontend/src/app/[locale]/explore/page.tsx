// app/[locale]/explore/page.tsx
"use client";

import { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { toast } from "react-hot-toast";
import {
  MagnifyingGlassIcon,
  FireIcon,
  ClockIcon,
  ChatBubbleLeftRightIcon,
  UserGroupIcon,
  HashtagIcon,
  NewspaperIcon,
  ArrowTrendingUpIcon,
  SparklesIcon,
  ArrowRightIcon,
} from "@heroicons/react/24/outline";
import { postApi, tagApi, topicApi, userApi } from "@/lib/api";
import { ActiveUserCard } from "@/components/explore/ActiveUserCard";
import { HotPostCard } from "@/components/explore/HotPostCard";
import { HotTagCard } from "@/components/explore/HotTagCard";
import { HotTopicCard } from "@/components/explore/HotTopicCard";
import type { LeaderboardItemResponse } from "@/lib/api/modules/users";
import { Post, Tag, Topic } from "@/shared/api/types";

// 分类 Tab
const exploreTabs = [
  {
    id: "hot",
    label: "热门",
    icon: FireIcon,
    sortBy: "hot",
    color: "text-red-500",
  },
  {
    id: "latest",
    label: "最新",
    icon: ClockIcon,
    sortBy: "latest",
    color: "text-blue-500",
  },
  {
    id: "trending",
    label: "趋势",
    icon: ArrowTrendingUpIcon,
    sortBy: "views",
    color: "text-green-500",
  },
  {
    id: "recommended",
    label: "推荐",
    icon: SparklesIcon,
    sortBy: "recommended",
    color: "text-purple-500",
  },
];

export default function Explore() {
  const { isAuthenticated } = useAuthStore();
  const [activeTab, setActiveTab] = useState("hot");
  const [posts, setPosts] = useState<Post[]>([]);
  const [hotTags, setHotTags] = useState<Tag[]>([]);
  const [hotTopics, setHotTopics] = useState<Topic[]>([]);
  const [activeUsers, setActiveUsers] = useState<LeaderboardItemResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchKeyword, setSearchKeyword] = useState("");
  const [searchResults, setSearchResults] = useState<Post[]>([]);
  const [searching, setSearching] = useState(false);

  // 加载探索数据 - 使用 useCallback 包装
  const loadExploreData = useCallback(async () => {
    setLoading(true);
    try {
      const currentTab = exploreTabs.find((tab) => tab.id === activeTab);
      const [postsResponse, tagsResponse, topicsResponse, usersResponse] =
        await Promise.all([
          postApi.list({ page: 1, page_size: 10, sort_by: currentTab?.sortBy }),
          tagApi.list(),
          topicApi.list({ page: 1, page_size: 8 }),
          userApi.getLeaderboardSimple({ limit: 10 }),
        ]);

      // 添加安全检查
      if (postsResponse.data.code === 200 && postsResponse.data.data) {
        setPosts(postsResponse.data.data.list || []);
      }
      if (tagsResponse.data.code === 200 && tagsResponse.data.data) {
        const sortedTags = [...(tagsResponse.data.data || [])].sort(
          (a, b) => (b.post_count || 0) - (a.post_count || 0),
        );
        setHotTags(sortedTags.slice(0, 12));
      }
      if (topicsResponse.data.code === 200 && topicsResponse.data.data) {
        setHotTopics(topicsResponse.data.data.list || []);
      }
      if (usersResponse.data.code === 200 && usersResponse.data.data) {
        setActiveUsers(usersResponse.data.data || []);
      }
    } catch (error) {
      console.error("Failed to load explore data:", error);
      toast.error("加载失败");
    } finally {
      setLoading(false);
    }
  }, [activeTab]); // 依赖 activeTab

  // 搜索
  const handleSearch = async () => {
    if (!searchKeyword.trim()) {
      toast.error("请输入搜索关键词");
      return;
    }

    setSearching(true);
    try {
      const response = await postApi.list({
        keyword: searchKeyword,
        page: 1,
        page_size: 20,
      });
      if (response.data.code === 200 && response.data.data) {
        setSearchResults(response.data.data.list || []);
        if (response.data.data.list?.length === 0) {
          toast("未找到相关结果");
        }
      }
    } catch (error) {
      console.error("Search failed:", error);
      toast.error("搜索失败");
    } finally {
      setSearching(false);
    }
  };

  // 清除搜索
  const clearSearch = () => {
    setSearchKeyword("");
    setSearchResults([]);
  };

  useEffect(() => {
    loadExploreData();
  }, [loadExploreData]); // 依赖 loadExploreData

  const displayPosts = searchKeyword ? searchResults : posts;

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-200 via-base-100 to-base-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Hero Section */}
        <div className="text-center mb-10">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-red-500 to-red-600 shadow-lg mb-4">
            <SparklesIcon className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-red-600 to-red-500 bg-clip-text text-transparent">
            探索
          </h1>
          <p className="text-base-content/60 mt-2">
            发现热门内容、话题和活跃用户
          </p>
        </div>

        {/* 搜索框 */}
        <div className="mb-8">
          <div className="card bg-base-100 shadow-md border border-base-200 p-2">
            <div className="flex gap-2">
              <div className="flex-1 relative">
                <MagnifyingGlassIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                <input
                  type="text"
                  value={searchKeyword}
                  onChange={(e) => setSearchKeyword(e.target.value)}
                  onKeyPress={(e) => e.key === "Enter" && handleSearch()}
                  placeholder="搜索帖子..."
                  className="w-full pl-11 pr-4 py-2.5 bg-transparent text-base-content placeholder-base-content/40 focus:outline-none"
                />
              </div>
              <button
                onClick={handleSearch}
                disabled={searching}
                className="btn btn-primary min-w-[80px]"
              >
                {searching ? (
                  <span className="loading loading-spinner loading-sm"></span>
                ) : (
                  "搜索"
                )}
              </button>
              {searchKeyword && (
                <button onClick={clearSearch} className="btn btn-ghost">
                  清除
                </button>
              )}
            </div>
          </div>
        </div>

        {/* 内容区域 */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 左侧：帖子列表 */}
          <div className="lg:col-span-2 space-y-4">
            {/* Tab 切换 */}
            {!searchKeyword && (
              <div className="card bg-base-100 shadow-md border border-base-200">
                <div className="card-body p-0">
                  <div className="flex border-b border-base-200">
                    {exploreTabs.map((tab) => {
                      const Icon = tab.icon;
                      const isActive = activeTab === tab.id;
                      return (
                        <button
                          key={tab.id}
                          onClick={() => setActiveTab(tab.id)}
                          className={`flex-1 py-3 text-center font-medium transition-all flex items-center justify-center gap-2 ${
                            isActive
                              ? "text-primary border-b-2 border-primary bg-gradient-to-b from-primary/5 to-transparent"
                              : "text-base-content/60 hover:text-base-content hover:bg-base-200/50"
                          }`}
                        >
                          <Icon
                            className={`w-4 h-4 ${isActive ? tab.color : ""}`}
                          />
                          {tab.label}
                        </button>
                      );
                    })}
                  </div>
                </div>
              </div>
            )}

            {/* 帖子列表 */}
            {loading && !searchKeyword ? (
              <div className="space-y-3">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div
                    key={i}
                    className="card bg-base-100 shadow-md border border-base-200"
                  >
                    <div className="card-body p-5">
                      <div className="h-5 bg-base-200 rounded w-3/4 mb-2 animate-pulse" />
                      <div className="h-4 bg-base-200 rounded w-1/2 animate-pulse" />
                      <div className="flex gap-4 mt-3">
                        <div className="h-3 bg-base-200 rounded w-16 animate-pulse" />
                        <div className="h-3 bg-base-200 rounded w-16 animate-pulse" />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : displayPosts.length === 0 ? (
              <div className="card bg-base-100 shadow-md border border-base-200">
                <div className="card-body py-16 text-center">
                  <div className="text-6xl mb-4">🔍</div>
                  <h3 className="text-lg font-semibold text-base-content mb-2">
                    {searchKeyword ? "未找到相关结果" : "暂无内容"}
                  </h3>
                  <p className="text-base-content/60">
                    {searchKeyword
                      ? "尝试使用其他关键词"
                      : "还没有内容，去发布第一个帖子吧！"}
                  </p>
                  {!searchKeyword && (
                    <Link href="/posts/new" className="btn btn-primary mt-4">
                      <ChatBubbleLeftRightIcon className="w-4 h-4" />
                      发布帖子
                    </Link>
                  )}
                </div>
              </div>
            ) : (
              <div className="space-y-3">
                {displayPosts.map((post, index) => (
                  <HotPostCard
                    key={post.id}
                    post={post}
                    rank={activeTab === "hot" && !searchKeyword ? index + 1 : 0}
                  />
                ))}
              </div>
            )}
          </div>

          {/* 右侧：侧边栏 */}
          <div className="space-y-4">
            {/* 热门标签 */}
            <div className="card bg-base-100 shadow-md border border-base-200">
              <div className="card-body p-5">
                <div className="flex items-center gap-2 mb-3">
                  <div className="w-8 h-8 rounded-lg bg-red-100 dark:bg-red-900/20 flex items-center justify-center">
                    <HashtagIcon className="w-4 h-4 text-red-500" />
                  </div>
                  <h2 className="font-semibold text-base-content">热门标签</h2>
                </div>
                {loading ? (
                  <div className="space-y-2">
                    {[1, 2, 3, 4].map((i) => (
                      <div
                        key={i}
                        className="h-16 bg-base-200 rounded animate-pulse"
                      />
                    ))}
                  </div>
                ) : hotTags.length === 0 ? (
                  <p className="text-sm text-base-content/40 text-center py-4">
                    暂无标签
                  </p>
                ) : (
                  <div className="flex flex-wrap gap-2">
                    {hotTags.slice(0, 12).map((tag) => (
                      <HotTagCard key={tag.id} tag={tag} />
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* 热门话题 */}
            <div className="card bg-base-100 shadow-md border border-base-200">
              <div className="card-body p-5">
                <div className="flex items-center gap-2 mb-3">
                  <div className="w-8 h-8 rounded-lg bg-green-100 dark:bg-green-900/20 flex items-center justify-center">
                    <NewspaperIcon className="w-4 h-4 text-green-500" />
                  </div>
                  <h2 className="font-semibold text-base-content">热门话题</h2>
                </div>
                {loading ? (
                  <div className="space-y-2">
                    {[1, 2, 3].map((i) => (
                      <div
                        key={i}
                        className="h-20 bg-base-200 rounded animate-pulse"
                      />
                    ))}
                  </div>
                ) : hotTopics.length === 0 ? (
                  <p className="text-sm text-base-content/40 text-center py-4">
                    暂无话题
                  </p>
                ) : (
                  <div className="space-y-2">
                    {hotTopics.slice(0, 5).map((topic) => (
                      <HotTopicCard key={topic.id} topic={topic} />
                    ))}
                  </div>
                )}
                <Link
                  href="/topics"
                  className="flex items-center justify-center gap-1 text-sm text-primary hover:text-primary-focus mt-3 transition-colors"
                >
                  查看更多话题
                  <ArrowRightIcon className="w-3 h-3" />
                </Link>
              </div>
            </div>

            {/* 积分榜 */}
            <div className="card bg-base-100 shadow-md border border-base-200">
              <div className="card-body p-5">
                <div className="flex items-center gap-2 mb-3">
                  <div className="w-8 h-8 rounded-lg bg-purple-100 dark:bg-purple-900/20 flex items-center justify-center">
                    <UserGroupIcon className="w-4 h-4 text-purple-500" />
                  </div>
                  <h2 className="font-semibold text-base-content">积分榜</h2>
                </div>
                {loading ? (
                  <div className="space-y-2">
                    {[1, 2, 3, 4].map((i) => (
                      <div
                        key={i}
                        className="h-14 bg-base-200 rounded animate-pulse"
                      />
                    ))}
                  </div>
                ) : activeUsers.length === 0 ? (
                  <p className="text-sm text-base-content/40 text-center py-4">
                    暂无用户
                  </p>
                ) : (
                  <div className="space-y-2">
                    {activeUsers.slice(0, 6).map((user) => (
                      <ActiveUserCard key={user.id} user={user} />
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* 关于 */}
            <div className="card bg-gradient-to-r from-red-50 to-red-100 dark:from-red-900/20 dark:to-red-800/10 border border-red-200 dark:border-red-800/30">
              <div className="card-body p-5">
                <h3 className="font-semibold text-base-content mb-2">
                  发现更多
                </h3>
                <p className="text-sm text-base-content/70 mb-3">
                  关注热门话题、标签和活跃用户，发现更多精彩内容
                </p>
                {!isAuthenticated && (
                  <Link href="/auth/login" className="btn btn-primary btn-sm">
                    登录后获得个性化推荐
                    <ArrowRightIcon className="w-3 h-3" />
                  </Link>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
