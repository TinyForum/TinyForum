// app/[locale]/topics/[id]/page.tsx (话题详情页)
'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuthStore } from '@/store/auth';
import { topicApi } from '@/lib/api/modules/topics';
import { toast } from 'react-hot-toast';
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
  FlagIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import type { Topic, TopicPost, Post } from '@/lib/api/types';
import { TopicPostCard } from '@/components/topic/TopicPostCard';

export default function TopicDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();
  const [topic, setTopic] = useState<Topic | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [followers, setFollowers] = useState<any[]>([]);
  const [following, setFollowing] = useState(false);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'posts' | 'followers'>('posts');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showAddPostModal, setShowAddPostModal] = useState(false);
  const [postId, setPostId] = useState('');
  const [addingPost, setAddingPost] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const pageSize = 20;

  const topicId = Number(params.id);

  // 加载话题详情
  const loadTopic = async () => {
    try {
      const response = await topicApi.getById(topicId);
      if (response.data.code === 200) {
        setTopic(response.data.data);
      } else {
        toast.error(response.data.message || '加载失败');
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '加载失败');
    }
  };

  // 加载话题帖子
  const loadPosts = async () => {
    try {
      const response = await topicApi.getPosts(topicId, { page, page_size: pageSize });
      if (response.data.code === 200) {
        const { list, total: totalCount } = response.data.data;
        setPosts(list.map(item => item.post).filter(Boolean) as Post[]);
        setTotal(totalCount || 0);
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '加载失败');
    }
  };

  // 加载关注者
  const loadFollowers = async () => {
    try {
      const response = await topicApi.getFollowers(topicId, { page, page_size: pageSize });
      if (response.data.code === 200) {
        const { list, total: totalCount } = response.data.data;
        setFollowers(list || []);
        setTotal(totalCount || 0);
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '加载失败');
    }
  };

  // 检查是否关注
  const checkFollowStatus = async () => {
    if (!user?.id) return;
    
    try {
      const response = await topicApi.getFollowers(topicId, { page: 1, page_size: 100 });
      if (response.data.code === 200) {
        const isFollowing = response.data.data.list.some(
          (f: any) => f.user_id === user?.id
        );
        setFollowing(isFollowing);
      }
    } catch (error) {
      console.error('Failed to check follow status:', error);
    }
  };

  // 关注/取消关注
  const handleFollow = async () => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      router.push('/login');
      return;
    }

    setFollowLoading(true);
    try {
      if (following) {
        await topicApi.unfollow(topicId);
        setFollowing(false);
        toast.success('已取消收藏');
        // 更新话题的收藏数
        if (topic) {
          setTopic({ ...topic, follower_count: (topic.follower_count || 1) - 1 });
        }
      } else {
        await topicApi.follow(topicId);
        setFollowing(true);
        toast.success('已收藏话题');
        if (topic) {
          setTopic({ ...topic, follower_count: (topic.follower_count || 0) + 1 });
        }
      }
      if (activeTab === 'followers') {
        await loadFollowers();
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '操作失败');
    } finally {
      setFollowLoading(false);
    }
  };

  // 添加帖子到话题
  const handleAddPost = async () => {
    const postIdNum = parseInt(postId);
    if (isNaN(postIdNum)) {
      toast.error('请输入有效的帖子ID');
      return;
    }

    setAddingPost(true);
    try {
      const response = await topicApi.addPost(topicId, { post_id: postIdNum });
      if (response.data.code === 200) {
        toast.success('帖子已添加到话题');
        setShowAddPostModal(false);
        setPostId('');
        await loadPosts();
        // 更新话题的帖子数
        if (topic) {
          setTopic({ ...topic, post_count: (topic.post_count || 0) + 1 });
        }
      } else {
        toast.error(response.data.message || '添加失败');
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '添加失败');
    } finally {
      setAddingPost(false);
    }
  };

  // 分享话题
  const handleShare = async () => {
    const url = window.location.href;
    try {
      await navigator.clipboard.writeText(url);
      toast.success('链接已复制到剪贴板');
    } catch (err) {
      toast.error('复制失败，请手动复制');
    }
  };

  useEffect(() => {
    if (topicId) {
      Promise.all([loadTopic(), checkFollowStatus()]).finally(() => {
        setLoading(false);
      });
    }
  }, [topicId]);

  useEffect(() => {
    if (activeTab === 'posts') {
      loadPosts();
    } else {
      loadFollowers();
    }
  }, [activeTab, page]);

  if (loading && !topic) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-base-content/60">加载中...</div>
      </div>
    );
  }

  if (!topic) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <p className="text-base-content/60 mb-4">话题不存在</p>
          <Link href="/topics" className="text-primary hover:underline">
            返回话题列表
          </Link>
        </div>
      </div>
    );
  }

  const totalPages = Math.ceil(total / pageSize);
  const isCreator = user?.id === topic.creator_id;

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
            返回话题列表
          </Link>
        </div>

        {/* 话题头部 - 使用主题卡片 */}
        <div className="card bg-base-100 shadow-lg mb-6 overflow-hidden hover:shadow-xl transition-shadow duration-300">
          {topic.cover && (
            <div className="w-full h-48 relative">
              <img
                src={topic.cover}
                alt={topic.title}
                className="w-full h-full object-cover"
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
                    <span>{topic.post_count || 0} 个帖子</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <UserGroupIcon className="w-4 h-4" />
                    <span>{topic.follower_count || 0} 人收藏</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <CalendarIcon className="w-4 h-4" />
                    <span>创建于 {new Date(topic.created_at).toLocaleDateString()}</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    {topic.is_public ? (
                      <>
                        <GlobeAltIcon className="w-4 h-4" />
                        <span className="badge badge-ghost badge-xs">公开</span>
                      </>
                    ) : (
                      <>
                        <LockClosedIcon className="w-4 h-4" />
                        <span className="badge badge-ghost badge-xs">私密</span>
                      </>
                    )}
                  </div>
                </div>
              </div>
              
              <div className="flex gap-2">
                <button
                  onClick={handleShare}
                  className="btn btn-ghost btn-square"
                  title="分享"
                >
                  <ShareIcon className="w-5 h-5" />
                </button>
                
                <button
                  onClick={handleFollow}
                  disabled={followLoading}
                  className={`btn gap-2 ${
                    following
                      ? 'btn-primary btn-outline'
                      : 'btn-primary'
                  }`}
                >
                  {followLoading ? (
                    <span className="loading loading-spinner loading-xs"></span>
                  ) : following ? (
                    <HeartSolidIcon className="w-5 h-5" />
                  ) : (
                    <HeartIcon className="w-5 h-5" />
                  )}
                  {following ? '已收藏' : '收藏'}
                </button>
                
                {isCreator && (
                  <button
                    onClick={() => setShowAddPostModal(true)}
                    className="btn btn-secondary gap-2"
                  >
                    <PlusIcon className="w-5 h-5" />
                    添加帖子
                  </button>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Tab 切换 - 使用主题 */}
        <div className="card bg-base-100 shadow-sm mb-6">
          <div className="tabs tabs-boxed p-1">
            <button
              onClick={() => {
                setActiveTab('posts');
                setPage(1);
              }}
              className={`tab flex-1 ${activeTab === 'posts' ? 'tab-active' : ''}`}
            >
              帖子列表 ({topic.post_count || 0})
            </button>
            <button
              onClick={() => {
                setActiveTab('followers');
                setPage(1);
              }}
              className={`tab flex-1 ${activeTab === 'followers' ? 'tab-active' : ''}`}
            >
              收藏者 ({topic.follower_count || 0})
            </button>
          </div>
        </div>

        {/* 内容区域 */}
        {activeTab === 'posts' ? (
          posts.length === 0 ? (
            <div className="card bg-base-100 shadow-sm p-12 text-center">
              <div className="text-6xl mb-4">📝</div>
              <h3 className="text-lg font-bold text-base-content mb-2">暂无帖子</h3>
              <p className="text-base-content/60 mb-6">还没有帖子被添加到这个话题</p>
              {isCreator && (
                <button
                  onClick={() => setShowAddPostModal(true)}
                  className="btn btn-primary"
                >
                  添加第一个帖子
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

              {/* 分页 - 使用 daisyUI */}
              {totalPages > 1 && (
                <div className="flex justify-center gap-2 mt-8">
                  <button
                    onClick={() => setPage(p => Math.max(1, p - 1))}
                    disabled={page === 1}
                    className="btn btn-outline btn-sm"
                  >
                    上一页
                  </button>
                  <div className="join">
                    {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                      let pageNum;
                      if (totalPages <= 5) {
                        pageNum = i + 1;
                      } else if (page <= 3) {
                        pageNum = i + 1;
                      } else if (page >= totalPages - 2) {
                        pageNum = totalPages - 4 + i;
                      } else {
                        pageNum = page - 2 + i;
                      }
                      
                      return (
                        <button
                          key={pageNum}
                          onClick={() => setPage(pageNum)}
                          className={`join-item btn btn-outline btn-sm ${
                            page === pageNum ? 'btn-primary' : ''
                          }`}
                        >
                          {pageNum}
                        </button>
                      );
                    })}
                  </div>
                  <button
                    onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                    disabled={page >= totalPages}
                    className="btn btn-outline btn-sm"
                  >
                    下一页
                  </button>
                </div>
              )}
            </>
          )
        ) : (
          followers.length === 0 ? (
            <div className="card bg-base-100 shadow-sm p-12 text-center">
              <div className="text-6xl mb-4">👥</div>
              <h3 className="text-lg font-bold text-base-content mb-2">暂无收藏者</h3>
              <p className="text-base-content/60">成为第一个收藏这个话题的人吧！</p>
            </div>
          ) : (
            <>
              <div className="space-y-3">
                {followers.map((follow) => (
                  <div key={follow.id} className="card bg-base-100 shadow-sm p-4 hover:shadow-md transition-shadow">
                    <div className="flex items-center gap-3">
                      <Link href={`/users/${follow.user_id}`}>
                        <div className="avatar placeholder">
                          <div className="w-10 h-10 rounded-full bg-primary/10 text-primary">
                            <span className="font-medium">
                              {follow.user?.username?.[0]?.toUpperCase() || `U${follow.user_id}`}
                            </span>
                          </div>
                        </div>
                      </Link>
                      <div className="flex-1">
                        <Link
                          href={`/users/${follow.user_id}`}
                          className="font-medium text-base-content hover:text-primary transition-colors"
                        >
                          {follow.user?.username || `用户 ${follow.user_id}`}
                        </Link>
                        <div className="text-xs text-base-content/40">
                          收藏于 {new Date(follow.created_at).toLocaleDateString()}
                        </div>
                      </div>
                      {follow.user_id === user?.id && (
                        <span className="badge badge-primary badge-sm">我</span>
                      )}
                    </div>
                  </div>
                ))}
              </div>

              {/* 分页 */}
              {totalPages > 1 && (
                <div className="flex justify-center gap-2 mt-8">
                  <button
                    onClick={() => setPage(p => Math.max(1, p - 1))}
                    disabled={page === 1}
                    className="btn btn-outline btn-sm"
                  >
                    上一页
                  </button>
                  <div className="join">
                    {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                      let pageNum;
                      if (totalPages <= 5) {
                        pageNum = i + 1;
                      } else if (page <= 3) {
                        pageNum = i + 1;
                      } else if (page >= totalPages - 2) {
                        pageNum = totalPages - 4 + i;
                      } else {
                        pageNum = page - 2 + i;
                      }
                      
                      return (
                        <button
                          key={pageNum}
                          onClick={() => setPage(pageNum)}
                          className={`join-item btn btn-outline btn-sm ${
                            page === pageNum ? 'btn-primary' : ''
                          }`}
                        >
                          {pageNum}
                        </button>
                      );
                    })}
                  </div>
                  <button
                    onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                    disabled={page >= totalPages}
                    className="btn btn-outline btn-sm"
                  >
                    下一页
                  </button>
                </div>
              )}
            </>
          )
        )}
      </div>

      {/* 添加帖子模态框 - 使用主题 */}
      {showAddPostModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="card bg-base-100 shadow-xl w-full max-w-md">
            <div className="card-body p-6">
              <div className="flex items-center justify-between mb-2">
                <h2 className="card-title text-xl font-bold text-base-content">
                  添加帖子到话题
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
                      帖子ID <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="number"
                    value={postId}
                    onChange={(e) => setPostId(e.target.value)}
                    placeholder="请输入要添加的帖子ID"
                    className="input input-bordered w-full focus:input-primary"
                  />
                  <label className="label">
                    <span className="label-text-alt text-base-content/40">
                      输入现有帖子的ID，该帖子将被添加到当前话题
                    </span>
                  </label>
                </div>
              </div>
              
              <div className="modal-action mt-6">
                <button
                  onClick={() => setShowAddPostModal(false)}
                  className="btn btn-ghost"
                >
                  取消
                </button>
                <button
                  onClick={handleAddPost}
                  disabled={addingPost || !postId}
                  className="btn btn-primary"
                >
                  {addingPost && <span className="loading loading-spinner loading-sm"></span>}
                  {addingPost ? '添加中...' : '添加'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}