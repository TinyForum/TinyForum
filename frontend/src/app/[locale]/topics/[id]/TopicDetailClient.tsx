// topics/[id]/page.tsx 客户端组件
"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { toast } from "react-hot-toast";
import { useTranslations } from "next-intl";
import { useTopicDetail } from "@/features/topic/hooks/useTopicDetail";
import {
  ArrowLeftIcon,
  HeartIcon,
  UserGroupIcon,
  DocumentTextIcon,
  CalendarIcon,
  GlobeAltIcon,
  LockClosedIcon,
  PlusIcon,
  XMarkIcon,
  ShareIcon,
} from "@heroicons/react/24/outline";
import { HeartIcon as HeartSolidIcon } from "@heroicons/react/24/solid";
import type { Topic } from "@/shared/api/types/topic.model";
import type { Post } from "@/shared/api/types/post.model";
import Image from "next/image";
import { TopicPostCard } from "@/layout/topic/TopicPostCard";

// 类型定义
interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

interface Follower {
  id: number;
  user_id: number;
  topic_id: number;
  user?: {
    id: number;
    username: string;
    avatar_url?: string;
  };
  // created_at:string// 可能不在返回的数据中
}
interface TopicPostItem {
  id: number;
  topic_id: number;
  post_id: number;
  post?: Post;
  added_by: number;
  sort_order: number;
  // created_at 可能不在 TopicPost 中
}

export default function TopicDetailClient({ id }: { id: string }) {
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();
  const [topic, setTopic] = useState<Topic | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [followers, setFollowers] = useState<Follower[]>([]);
  const [following, setFollowing] = useState(false);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"posts" | "followers">("posts");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showAddPostModal, setShowAddPostModal] = useState(false);
  const [postId, setPostId] = useState("");
  const [addingPost, setAddingPost] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const pageSize = 20;

  const t = useTranslations("Topic");

  const topicId = Number(id);

  const {
    fetchTopic,
    fetchTopicPosts,
    fetchTopicFollowers,
    followTopic,
    unfollowTopic,
    addPostToTopic,
  } = useTopicDetail();

  // 加载话题详情
  const loadTopic = useCallback(async () => {
    try {
      const response = await fetchTopic(topicId);
      if (response.data.code === 0) {
        // 修复：确保 response.data.data 存在，否则设为 null
        const data = response.data.data;
        if (data) {
          setTopic(data);
        } else {
          setTopic(null);
          toast.error(t("topic_not_found"));
        }
      } else {
        toast.error(response.data.message || t("load_failed"));
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(error.response?.data?.message || t("load_failed"));
    }
  }, [topicId, fetchTopic]);

  const loadPosts = useCallback(async () => {
    try {
      const response = await fetchTopicPosts(topicId, {
        page,
        page_size: pageSize,
      });
      if (response.data.code === 0) {
        const data = response.data.data;
        const list = data?.list || [];
        // 修复：正确处理 list 中的每个元素
        const mappedPosts: Post[] = list
          .map((item: TopicPostItem) => item.post) // 移除可选链操作符
          .filter((post: Post | undefined): post is Post => post !== undefined);

        setPosts(mappedPosts);
        setTotal(data?.total || 0);
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(error.response?.data?.message || t("load_failed"));
    }
  }, [topicId, page, fetchTopicPosts]);

  // 修复 loadFollowers 函数
  const loadFollowers = useCallback(async () => {
    try {
      const response = await fetchTopicFollowers(topicId, {
        page,
        page_size: pageSize,
      });
      if (response.data.code === 0) {
        const data = response.data.data;
        // 修复：直接使用 API 返回的数据，不要求 created_at
        const list = data?.list || [];
        setFollowers(list);
        setTotal(data?.total || 0);
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(error.response?.data?.message || t("load_failed"));
    }
  }, [topicId, page, fetchTopicFollowers]);

  // 修复 checkFollowStatus 函数
  const checkFollowStatus = useCallback(async () => {
    if (!user?.id) return;

    try {
      const response = await fetchTopicFollowers(topicId, {
        page: 1,
        page_size: 100,
      });
      if (response.data.code === 0) {
        const data = response.data.data;
        const isFollowing = (data?.list || []).some(
          (f: { user_id: number }) => f.user_id === user?.id,
        );
        setFollowing(isFollowing);
      }
    } catch (err) {
      console.error("Failed to check follow status:", err);
    }
  }, [topicId, user?.id, fetchTopicFollowers]);

  // 关注/取消关注
  const handleFollow = async () => {
    if (!isAuthenticated) {
      toast.error(t("please_login"));
      router.push("/login");
      return;
    }

    setFollowLoading(true);
    try {
      if (following) {
        await unfollowTopic(topicId);
        setFollowing(false);
        toast.success(t("unfollow_success"));
        if (topic) {
          setTopic({
            ...topic,
            follower_count: (topic.follower_count || 1) - 1,
          });
        }
      } else {
        await followTopic(topicId);
        setFollowing(true);
        toast.success(t("follow_success"));
        if (topic) {
          setTopic({
            ...topic,
            follower_count: (topic.follower_count || 0) + 1,
          });
        }
      }
      if (activeTab === "followers") {
        await loadFollowers();
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(error.response?.data?.message || t("operation_failed"));
    } finally {
      setFollowLoading(false);
    }
  };

  // 添加帖子到话题
  const handleAddPost = async () => {
    const postIdNum = parseInt(postId);
    if (isNaN(postIdNum)) {
      toast.error(t("enter_valid_post_id"));
      return;
    }

    setAddingPost(true);
    try {
      const response = await addPostToTopic(topicId, postIdNum);
      if (response.data.code === 0) {
        toast.success(t("post_added_success"));
        setShowAddPostModal(false);
        setPostId("");
        await loadPosts();
        if (topic) {
          setTopic({ ...topic, post_count: (topic.post_count || 0) + 1 });
        }
      } else {
        toast.error(response.data.message || t("add_failed"));
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(error.response?.data?.message || t("add_failed"));
    } finally {
      setAddingPost(false);
    }
  };

  // 分享话题
  const handleShare = async () => {
    const url = window.location.href;
    try {
      await navigator.clipboard.writeText(url);
      toast.success(t("link_copied"));
    } catch {
      toast.error(t("copy_failed"));
    }
  };

  useEffect(() => {
    if (topicId) {
      Promise.all([loadTopic(), checkFollowStatus()]).finally(() => {
        setLoading(false);
      });
    }
  }, [topicId, loadTopic, checkFollowStatus]);

  useEffect(() => {
    if (activeTab === "posts") {
      setPage(1);
      loadPosts();
    } else {
      setPage(1);
      loadFollowers();
    }
  }, [activeTab, loadPosts, loadFollowers]);

  useEffect(() => {
    if (activeTab === "posts") {
      loadPosts();
    } else {
      loadFollowers();
    }
  }, [page, activeTab, loadPosts, loadFollowers]);

  if (loading && !topic) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-base-content/60">{t("loading")}</div>
      </div>
    );
  }

  if (!topic) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <p className="text-base-content/60 mb-4">{t("topic_not_found")}</p>
          <Link href="/topics" className="text-primary hover:underline">
            {t("back_to_topic_list")}
          </Link>
        </div>
      </div>
    );
  }

  const totalPages = Math.ceil(total / pageSize);
  const isCreator = user?.id === topic.creator_id;

  // 生成分页按钮
  const renderPaginationButtons = () => {
    const buttons = [];
    const maxVisible = 5;
    let startPage = 1;
    let endPage = totalPages;

    if (totalPages > maxVisible) {
      if (page <= 3) {
        endPage = maxVisible;
      } else if (page >= totalPages - 2) {
        startPage = totalPages - maxVisible + 1;
      } else {
        startPage = page - 2;
        endPage = page + 2;
      }
    }

    for (let i = startPage; i <= endPage; i++) {
      buttons.push(
        <button
          key={i}
          onClick={() => setPage(i)}
          className={`join-item btn btn-outline btn-sm ${
            page === i ? "btn-primary" : ""
          }`}
        >
          {i}
        </button>,
      );
    }
    return buttons;
  };

  return (
    <div className="min-h-screen bg-base-200 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* 返回按钮 */}
        <div className="mb-4">
          <Link
            href="/topics"
            className="inline-flex items-center gap-2 text-base-content/70 hover:text-primary transition-colors"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            {t("back_to_topic_list")}
          </Link>
        </div>

        {/* 话题头部 */}
        <div className="card bg-base-100 shadow-lg mb-6 overflow-hidden hover:shadow-xl transition-shadow duration-300">
          {topic.cover && (
            <div className="w-full h-48 relative">
              <Image
                src={topic.cover}
                alt={topic.title}
                fill
                className="object-cover"
                sizes="100vw"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent" />
            </div>
          )}

          <div className="card-body p-6">
            <div className="flex flex-col lg:flex-row items-start justify-between gap-4">
              <div className="flex-1">
                <h1 className="text-2xl md:text-3xl font-bold text-base-content mb-3">
                  {topic.title}
                </h1>
                {topic.description && (
                  <p className="text-base-content/70 mb-4 leading-relaxed">
                    {topic.description}
                  </p>
                )}

                <div className="flex flex-wrap items-center gap-4 text-sm text-base-content/50">
                  <div className="flex items-center gap-1.5">
                    <DocumentTextIcon className="w-4 h-4" />
                    <span>{t("post_count", { count: topic.post_count || 0 })}</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <UserGroupIcon className="w-4 h-4" />
                    <span>{t("follower_count", { count: topic.follower_count || 0 })}</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <CalendarIcon className="w-4 h-4" />
                    <span>
                      {t("created_at", { date: new Date(topic.created_at).toLocaleDateString() })}
                    </span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    {topic.is_public ? (
                      <>
                        <GlobeAltIcon className="w-4 h-4" />
                        <span className="badge badge-ghost badge-xs">{t("public")}</span>
                      </>
                    ) : (
                      <>
                        <LockClosedIcon className="w-4 h-4" />
                        <span className="badge badge-ghost badge-xs">{t("private")}</span>
                      </>
                    )}
                  </div>
                </div>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={handleShare}
                  className="btn btn-ghost btn-square"
                  title={t("share")}
                >
                  <ShareIcon className="w-5 h-5" />
                </button>

                <button
                  onClick={handleFollow}
                  disabled={followLoading}
                  className={`btn gap-2 ${
                    following ? "btn-primary btn-outline" : "btn-primary"
                  }`}
                >
                  {followLoading ? (
                    <span className="loading loading-spinner loading-xs"></span>
                  ) : following ? (
                    <HeartSolidIcon className="w-5 h-5" />
                  ) : (
                    <HeartIcon className="w-5 h-5" />
                  )}
                  {following ? t("followed") : t("follow")}
                </button>

                {isCreator && (
                  <button
                    onClick={() => setShowAddPostModal(true)}
                    className="btn btn-secondary gap-2"
                  >
                    <PlusIcon className="w-5 h-5" />
                    {t("add_post")}
                  </button>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Tab 切换 */}
        <div className="card bg-base-100 shadow-sm mb-6">
          <div className="tabs tabs-boxed p-1">
            <button
              onClick={() => {
                setActiveTab("posts");
                setPage(1);
              }}
              className={`tab flex-1 ${activeTab === "posts" ? "tab-active" : ""}`}
            >
              {t("post_list")} ({topic.post_count || 0})
            </button>
            <button
              onClick={() => {
                setActiveTab("followers");
                setPage(1);
              }}
              className={`tab flex-1 ${activeTab === "followers" ? "tab-active" : ""}`}
            >
              {t("followers")} ({topic.follower_count || 0})
            </button>
          </div>
        </div>

        {/* 内容区域 - 帖子列表 */}
        {activeTab === "posts" && (
          <>
            {posts.length === 0 ? (
              <div className="card bg-base-100 shadow-sm p-12 text-center">
                <div className="text-6xl mb-4">📝</div>
                <h3 className="text-lg font-bold text-base-content mb-2">
                  {t("no_posts")}
                </h3>
                <p className="text-base-content/60 mb-6">
                  {t("no_posts_desc")}
                </p>
                {isCreator && (
                  <button
                    onClick={() => setShowAddPostModal(true)}
                    className="btn btn-primary"
                  >
                    {t("add_first_post")}
                  </button>
                )}
              </div>
            ) : (
              <>
                <div className="space-y-3">
                  {posts.map((post) => (
                    <TopicPostCard key={post.id} post={post} />
                  ))}
                </div>

                {totalPages > 1 && (
                  <div className="flex justify-center gap-2 mt-8">
                    <button
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                      disabled={page === 1}
                      className="btn btn-outline btn-sm"
                    >
                      {t("previous_page")}
                    </button>
                    <div className="join">{renderPaginationButtons()}</div>
                    <button
                      onClick={() =>
                        setPage((p) => Math.min(totalPages, p + 1))
                      }
                      disabled={page >= totalPages}
                      className="btn btn-outline btn-sm"
                    >
                      {t("next_page")}
                    </button>
                  </div>
                )}
              </>
            )}
          </>
        )}

        {/* 内容区域 - 收藏者列表 */}
        {activeTab === "followers" && (
          <>
            {followers.length === 0 ? (
              <div className="card bg-base-100 shadow-sm p-12 text-center">
                <div className="text-6xl mb-4">👥</div>
                <h3 className="text-lg font-bold text-base-content mb-2">
                  {t("no_followers")}
                </h3>
                <p className="text-base-content/60">
                  {t("no_followers_desc")}
                </p>
              </div>
            ) : (
              <>
                <div className="space-y-3">
                  {followers.map((follow) => (
                    <div
                      key={follow.id}
                      className="card bg-base-100 shadow-sm p-4 hover:shadow-md transition-shadow"
                    >
                      <div className="flex items-center gap-3">
                        <Link href={`/users/${follow.user_id}`}>
                          <div className="avatar placeholder">
                            <div className="w-10 h-10 rounded-full bg-primary/10 text-primary flex items-center justify-center">
                              <span className="font-medium">
                                {follow.user?.username?.[0]?.toUpperCase() ||
                                  `U${follow.user_id}`}
                              </span>
                            </div>
                          </div>
                        </Link>
                        <div className="flex-1">
                          <Link
                            href={`/users/${follow.user_id}`}
                            className="font-medium text-base-content hover:text-primary transition-colors"
                          >
                            {follow.user?.username || t("user_id_format", { id: follow.user_id })}
                          </Link>
                          <div className="text-xs text-base-content/40">
                            {t("followed_at")}{" "}
                            {/* {new Date(follow.created_at).toLocaleDateString()} */}
                          </div>
                        </div>
                        {follow.user_id === user?.id && (
                          <span className="badge badge-primary badge-sm">
                            {t("me")}
                          </span>
                        )}
                      </div>
                    </div>
                  ))}
                </div>

                {totalPages > 1 && (
                  <div className="flex justify-center gap-2 mt-8">
                    <button
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                      disabled={page === 1}
                  className="btn btn-outline btn-sm"
                >
                  {t("previous_page")}
                </button>
                <div className="join">{renderPaginationButtons()}</div>
                <button
                  onClick={() =>
                    setPage((p) => Math.min(totalPages, p + 1))
                  }
                  disabled={page >= totalPages}
                  className="btn btn-outline btn-sm"
                >
                  {t("next_page")}
                </button>
                  </div>
                )}
              </>
            )}
          </>
        )}
      </div>

      {/* 添加帖子模态框 */}
      {showAddPostModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="card bg-base-100 shadow-xl w-full max-w-md">
            <div className="card-body p-6">
              <div className="flex items-center justify-between mb-2">
                <h2 className="card-title text-xl font-bold text-base-content">
                  {t("add_post_to_topic")}
                </h2>
                <button
                  onClick={() => setShowAddPostModal(false)}
                  className="btn btn-ghost btn-sm btn-circle"
                >
                  <XMarkIcon className="w-5 h-5" />
                </button>
              </div>

              <div className="space-y-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      {t("post_id_label")} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="number"
                    value={postId}
                    onChange={(e) => setPostId(e.target.value)}
                    placeholder={t("enter_post_id")}
                    className="input input-bordered w-full focus:input-primary"
                  />
                  <label className="label">
                    <span className="label-text-alt text-base-content/40">
                      {t("enter_post_id_hint")}
                    </span>
                  </label>
                </div>
              </div>

              <div className="modal-action mt-6">
                  <button
                    onClick={() => setShowAddPostModal(false)}
                    className="btn btn-ghost"
                  >
                    {t("cancel")}
                  </button>
                <button
                  onClick={handleAddPost}
                  disabled={addingPost || !postId}
                  className="btn btn-primary"
                >
                  {addingPost && (
                    <span className="loading loading-spinner loading-sm"></span>
                  )}
                  {addingPost ? t("adding") : t("add")}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
