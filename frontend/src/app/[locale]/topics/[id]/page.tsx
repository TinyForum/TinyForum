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
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import type { Topic, TopicPost, Post } from '@/lib/api/types';

// 话题帖子卡片组件
function TopicPostCard({ post }: { post: Post }) {
  return (
    <Link href={`/posts/${post.id}`}>
      <div className="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <h4 className="font-semibold text-gray-900 hover:text-indigo-600 mb-2">
              {post.title}
            </h4>
            {post.summary && (
              <p className="text-gray-600 text-sm line-clamp-2 mb-2">
                {post.summary}
              </p>
            )}
            <div className="flex items-center gap-4 text-xs text-gray-400">
              <span>{new Date(post.created_at).toLocaleDateString()}</span>
              <span>👁️ {post.view_count}</span>
              <span>❤️ {post.like_count}</span>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}

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
        const { items, total: totalCount } = response.data.data;
        setPosts(items.map(item => item.post).filter(Boolean) as Post[]);
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
        const { items, total: totalCount } = response.data.data;
        setFollowers(items || []);
        setTotal(totalCount || 0);
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '加载失败');
    }
  };

  // 检查是否关注
  const checkFollowStatus = async () => {
    try {
      const response = await topicApi.getFollowers(topicId, { page: 1, page_size: 1 });
      if (response.data.code === 200) {
        const isFollowing = response.data.data.items.some(
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
      } else {
        await topicApi.follow(topicId);
        setFollowing(true);
        toast.success('已收藏话题');
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
      } else {
        toast.error(response.data.message || '添加失败');
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '添加失败');
    } finally {
      setAddingPost(false);
    }
  };

  useEffect(() => {
    if (topicId) {
      loadTopic();
      checkFollowStatus();
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
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-gray-500">加载中...</div>
      </div>
    );
  }

  if (!topic) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-500 mb-4">话题不存在</p>
          <Link href="/topics" className="text-indigo-600 hover:underline">
            返回话题列表
          </Link>
        </div>
      </div>
    );
  }

  const totalPages = Math.ceil(total / pageSize);
  const isCreator = user?.id === topic.creator_id;

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* 返回按钮 */}
        <div className="mb-4">
          <Link
            href="/topics"
            className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            返回话题列表
          </Link>
        </div>

        {/* 话题头部 */}
        <div className="bg-white rounded-lg shadow-sm mb-6 overflow-hidden">
          {topic.cover && (
            <div className="w-full h-48">
              <img
                src={topic.cover}
                alt={topic.title}
                className="w-full h-full object-cover"
              />
            </div>
          )}
          
          <div className="p-6">
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <h1 className="text-2xl font-bold text-gray-900 mb-2">
                  {topic.title}
                </h1>
                {topic.description && (
                  <p className="text-gray-600 mb-4">{topic.description}</p>
                )}
                
                <div className="flex items-center gap-4 text-sm text-gray-500">
                  <div className="flex items-center gap-1">
                    <DocumentTextIcon className="w-4 h-4" />
                    {topic.post_count} 个帖子
                  </div>
                  <div className="flex items-center gap-1">
                    <UserGroupIcon className="w-4 h-4" />
                    {topic.follower_count} 人收藏
                  </div>
                  <div className="flex items-center gap-1">
                    <CalendarIcon className="w-4 h-4" />
                    创建于 {new Date(topic.created_at).toLocaleDateString()}
                  </div>
                  <div className="flex items-center gap-1">
                    {topic.is_public ? (
                      <>
                        <GlobeAltIcon className="w-4 h-4" />
                        <span>公开</span>
                      </>
                    ) : (
                      <>
                        <LockClosedIcon className="w-4 h-4" />
                        <span>私密</span>
                      </>
                    )}
                  </div>
                </div>
              </div>
              
              <div className="flex gap-2">
                <button
                  onClick={handleFollow}
                  disabled={followLoading}
                  className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                    following
                      ? 'bg-indigo-50 text-indigo-600 hover:bg-indigo-100'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  } disabled:opacity-50`}
                >
                  <HeartIcon className={`w-5 h-5 inline mr-1 ${following ? 'fill-current' : ''}`} />
                  {following ? '已收藏' : '收藏'}
                </button>
                
                {isCreator && (
                  <button
                    onClick={() => setShowAddPostModal(true)}
                    className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
                  >
                    <PlusIcon className="w-5 h-5 inline mr-1" />
                    添加帖子
                  </button>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Tab 切换 */}
        <div className="bg-white rounded-lg shadow-sm mb-6">
          <div className="flex border-b">
            <button
              onClick={() => {
                setActiveTab('posts');
                setPage(1);
              }}
              className={`flex-1 py-3 text-center font-medium transition-colors ${
                activeTab === 'posts'
                  ? 'text-indigo-600 border-b-2 border-indigo-600'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              帖子列表 ({topic.post_count})
            </button>
            <button
              onClick={() => {
                setActiveTab('followers');
                setPage(1);
              }}
              className={`flex-1 py-3 text-center font-medium transition-colors ${
                activeTab === 'followers'
                  ? 'text-indigo-600 border-b-2 border-indigo-600'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              收藏者 ({topic.follower_count})
            </button>
          </div>
        </div>

        {/* 内容区域 */}
        {activeTab === 'posts' ? (
          posts.length === 0 ? (
            <div className="bg-white rounded-lg shadow-sm p-12 text-center">
              <div className="text-6xl mb-4">📝</div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">暂无帖子</h3>
              <p className="text-gray-500">还没有帖子被添加到这个话题</p>
              {isCreator && (
                <button
                  onClick={() => setShowAddPostModal(true)}
                  className="mt-4 inline-block px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
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
          )
        ) : (
          followers.length === 0 ? (
            <div className="bg-white rounded-lg shadow-sm p-12 text-center">
              <div className="text-6xl mb-4">👥</div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">暂无收藏者</h3>
              <p className="text-gray-500">成为第一个收藏这个话题的人吧！</p>
            </div>
          ) : (
            <>
              <div className="space-y-3">
                {followers.map((follow) => (
                  <div key={follow.id} className="bg-white rounded-lg shadow-sm p-4">
                    <div className="flex items-center gap-3">
                      <Link href={`/users/${follow.user_id}`}>
                        <div className="w-10 h-10 bg-indigo-100 rounded-full flex items-center justify-center">
                          <span className="text-indigo-600 font-medium">
                            U{follow.user_id}
                          </span>
                        </div>
                      </Link>
                      <div>
                        <Link
                          href={`/users/${follow.user_id}`}
                          className="font-medium text-gray-900 hover:text-indigo-600"
                        >
                          用户 {follow.user_id}
                        </Link>
                        <div className="text-xs text-gray-400">
                          收藏于 {new Date(follow.created_at).toLocaleDateString()}
                        </div>
                      </div>
                    </div>
                  </div>
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
          )
        )}
      </div>

      {/* 添加帖子模态框 */}
      {showAddPostModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
            <div className="p-6 border-b flex justify-between items-center">
              <h2 className="text-xl font-bold text-gray-900">添加帖子到话题</h2>
              <button
                onClick={() => setShowAddPostModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <XMarkIcon className="w-6 h-6" />
              </button>
            </div>
            
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  帖子ID <span className="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  value={postId}
                  onChange={(e) => setPostId(e.target.value)}
                  placeholder="请输入要添加的帖子ID"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
                />
                <p className="mt-1 text-sm text-gray-400">
                  输入现有帖子的ID，该帖子将被添加到当前话题
                </p>
              </div>
            </div>
            
            <div className="p-6 border-t flex justify-end gap-3">
              <button
                onClick={() => setShowAddPostModal(false)}
                className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
              >
                取消
              </button>
              <button
                onClick={handleAddPost}
                disabled={addingPost || !postId}
                className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {addingPost ? '添加中...' : '添加'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}