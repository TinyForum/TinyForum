// app/[locale]/topics/page.tsx (话题列表页)
'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { topicApi } from '@/lib/api/modules/topics';
import { toast } from 'react-hot-toast';
import {
  PlusIcon,
  HeartIcon,
  UserGroupIcon,
  DocumentTextIcon,
  CalendarIcon,
  GlobeAltIcon,
  LockClosedIcon,
} from '@heroicons/react/24/outline';
import type { Topic } from '@/lib/api/types';

// 话题卡片组件
function TopicCard({ topic, onFollowChange }: { topic: Topic; onFollowChange?: () => void }) {
  const { isAuthenticated } = useAuthStore();
  const [following, setFollowing] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);

  const handleFollow = async () => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      return;
    }

    setFollowLoading(true);
    try {
      if (following) {
        await topicApi.unfollow(topic.id);
        setFollowing(false);
        toast.success('已取消收藏');
      } else {
        await topicApi.follow(topic.id);
        setFollowing(true);
        toast.success('已收藏话题');
      }
      onFollowChange?.();
    } catch (error: any) {
      toast.error(error.response?.data?.message || '操作失败');
    } finally {
      setFollowLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6 hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          {/* 话题封面 */}
          {topic.cover && (
            <div className="mb-3">
              <img
                src={topic.cover}
                alt={topic.title}
                className="w-full h-32 object-cover rounded-lg"
              />
            </div>
          )}
          
          {/* 话题标题 */}
          <Link href={`/topics/${topic.id}`}>
            <h3 className="text-xl font-bold text-gray-900 hover:text-indigo-600 transition-colors mb-2">
              {topic.title}
            </h3>
          </Link>
          
          {/* 话题描述 */}
          {topic.description && (
            <p className="text-gray-600 mb-4 line-clamp-2">
              {topic.description}
            </p>
          )}
          
          {/* 统计信息 */}
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
        
        {/* 收藏按钮 */}
        <button
          onClick={handleFollow}
          disabled={followLoading}
          className={`ml-4 px-4 py-2 rounded-lg font-medium transition-colors ${
            following
              ? 'bg-indigo-50 text-indigo-600 hover:bg-indigo-100'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          } disabled:opacity-50`}
        >
          <HeartIcon className={`w-5 h-5 inline mr-1 ${following ? 'fill-current' : ''}`} />
          {following ? '已收藏' : '收藏'}
        </button>
      </div>
    </div>
  );
}

export default function TopicsPage() {
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createForm, setCreateForm] = useState({
    title: '',
    description: '',
    cover: '',
    is_public: true,
  });
  const [creating, setCreating] = useState(false);
  const pageSize = 20;

  // 加载话题列表
  const loadTopics = async () => {
    setLoading(true);
    try {
      const response = await topicApi.list({ page, page_size: pageSize });
      if (response.data.code === 200) {
        const { items, total: totalCount } = response.data.data;
        setTopics(items || []);
        setTotal(totalCount || 0);
      } else {
        toast.error(response.data.message || '加载失败');
      }
    } catch (error: any) {
      console.error('Failed to load topics:', error);
      toast.error(error.response?.data?.message || '加载失败');
    } finally {
      setLoading(false);
    }
  };

  // 创建话题
  const handleCreateTopic = async () => {
    if (!createForm.title.trim()) {
      toast.error('请输入话题标题');
      return;
    }

    setCreating(true);
    try {
      const response = await topicApi.create({
        title: createForm.title.trim(),
        description: createForm.description.trim() || undefined,
        cover: createForm.cover || undefined,
        is_public: createForm.is_public,
      });
      
      if (response.data.code === 200 || response.data.code === 201) {
        toast.success('话题创建成功');
        setShowCreateModal(false);
        setCreateForm({ title: '', description: '', cover: '', is_public: true });
        await loadTopics();
      } else {
        toast.error(response.data.message || '创建失败');
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '创建失败');
    } finally {
      setCreating(false);
    }
  };

  useEffect(() => {
    loadTopics();
  }, [page]);

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* 头部 */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">话题</h1>
              <p className="text-gray-500">发现和参与感兴趣的话题讨论</p>
            </div>
            {isAuthenticated && (
              <button
                onClick={() => setShowCreateModal(true)}
                className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
              >
                <PlusIcon className="w-5 h-5" />
                创建话题
              </button>
            )}
          </div>
        </div>

        {/* 话题列表 */}
        {loading ? (
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-white rounded-lg shadow-sm p-6 animate-pulse">
                <div className="h-6 bg-gray-200 rounded w-1/3 mb-2" />
                <div className="h-4 bg-gray-200 rounded w-2/3 mb-4" />
                <div className="flex gap-4">
                  <div className="h-4 bg-gray-200 rounded w-20" />
                  <div className="h-4 bg-gray-200 rounded w-20" />
                </div>
              </div>
            ))}
          </div>
        ) : topics.length === 0 ? (
          <div className="bg-white rounded-lg shadow-sm p-12 text-center">
            <div className="text-6xl mb-4">📚</div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">暂无话题</h3>
            <p className="text-gray-500 mb-4">成为第一个创建话题的人吧！</p>
            {isAuthenticated && (
              <button
                onClick={() => setShowCreateModal(true)}
                className="inline-block px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
              >
                创建话题
              </button>
            )}
          </div>
        ) : (
          <>
            <div className="space-y-4">
              {topics.map((topic) => (
                <TopicCard key={topic.id} topic={topic} onFollowChange={loadTopics} />
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

      {/* 创建话题模态框 */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
            <div className="p-6 border-b">
              <h2 className="text-xl font-bold text-gray-900">创建话题</h2>
            </div>
            
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  标题 <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={createForm.title}
                  onChange={(e) => setCreateForm(prev => ({ ...prev, title: e.target.value }))}
                  placeholder="请输入话题标题"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
                  maxLength={100}
                />
                <p className="mt-1 text-sm text-gray-400">
                  {createForm.title.length}/100
                </p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  描述
                </label>
                <textarea
                  value={createForm.description}
                  onChange={(e) => setCreateForm(prev => ({ ...prev, description: e.target.value }))}
                  placeholder="描述这个话题的内容和讨论方向"
                  rows={3}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
                  maxLength={500}
                />
                <p className="mt-1 text-sm text-gray-400">
                  {createForm.description.length}/500
                </p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  封面图 URL
                </label>
                <input
                  type="url"
                  value={createForm.cover}
                  onChange={(e) => setCreateForm(prev => ({ ...prev, cover: e.target.value }))}
                  placeholder="https://example.com/cover.jpg"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
                />
              </div>
              
              <div>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={createForm.is_public}
                    onChange={(e) => setCreateForm(prev => ({ ...prev, is_public: e.target.checked }))}
                    className="w-4 h-4 text-indigo-600 rounded focus:ring-indigo-500"
                  />
                  <span className="text-sm text-gray-700">公开话题</span>
                </label>
                <p className="mt-1 text-sm text-gray-400 ml-6">
                  公开话题所有人可见，私密话题仅关注者可查看
                </p>
              </div>
            </div>
            
            <div className="p-6 border-t flex justify-end gap-3">
              <button
                onClick={() => setShowCreateModal(false)}
                className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
              >
                取消
              </button>
              <button
                onClick={handleCreateTopic}
                disabled={creating || !createForm.title.trim()}
                className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {creating ? '创建中...' : '创建'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}