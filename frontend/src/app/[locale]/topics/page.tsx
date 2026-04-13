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
  XMarkIcon,
} from '@heroicons/react/24/outline';
import type { Topic } from '@/lib/api/types';
import { TopicCard } from '@/components/topic/TopicCard';

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
        const { list, total: totalCount } = response.data.data;
        setTopics(list || []);
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
    <div className="min-h-screen bg-base-200 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* 头部 - 使用主题卡片 */}
        <div className="card bg-base-100 shadow-lg mb-6 hover:shadow-xl transition-shadow duration-300">
          <div className="card-body p-6">
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
              <div>
                <h1 className="text-2xl md:text-3xl font-bold text-base-content mb-2">
                  话题
                </h1>
                <p className="text-base-content/60">
                  发现和参与感兴趣的话题讨论
                </p>
              </div>
              {isAuthenticated && (
                <button
                  onClick={() => setShowCreateModal(true)}
                  className="btn btn-primary gap-2"
                >
                  <PlusIcon className="w-5 h-5" />
                  创建话题
                </button>
              )}
            </div>
          </div>
        </div>

        {/* 话题列表 */}
        {loading ? (
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="card bg-base-100 shadow-sm p-6 animate-pulse">
                <div className="skeleton h-6 w-1/3 mb-2" />
                <div className="skeleton h-4 w-2/3 mb-4" />
                <div className="flex gap-4">
                  <div className="skeleton h-4 w-20" />
                  <div className="skeleton h-4 w-20" />
                </div>
              </div>
            ))}
          </div>
        ) : topics.length === 0 ? (
          <div className="card bg-base-100 shadow-sm p-12 text-center">
            <div className="text-6xl mb-4">📚</div>
            <h3 className="text-lg font-bold text-base-content mb-2">暂无话题</h3>
            <p className="text-base-content/60 mb-6">成为第一个创建话题的人吧！</p>
            {isAuthenticated && (
              <button
                onClick={() => setShowCreateModal(true)}
                className="btn btn-primary"
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

            {/* 分页 - 使用 daisyUI 组件 */}
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
                  {totalPages > 5 && page < totalPages - 2 && (
                    <>
                      <button className="join-item btn btn-outline btn-sm disabled">...</button>
                      <button
                        onClick={() => setPage(totalPages)}
                        className="join-item btn btn-outline btn-sm"
                      >
                        {totalPages}
                      </button>
                    </>
                  )}
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
        )}
      </div>

      {/* 创建话题模态框 - 使用 daisyUI 主题 */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="card bg-base-100 shadow-xl w-full max-w-md">
            <div className="card-body p-6">
              {/* 模态框头部 */}
              <div className="flex items-center justify-between mb-2">
                <h2 className="card-title text-xl font-bold text-base-content">
                  创建话题
                </h2>
                <button
                  onClick={() => setShowCreateModal(false)}
                  className="btn btn-ghost btn-sm btn-circle"
                >
                  <XMarkIcon className="w-5 h-5" />
                </button>
              </div>
              
              {/* 表单内容 */}
              <div className="space-y-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      标题 <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="text"
                    value={createForm.title}
                    onChange={(e) => setCreateForm(prev => ({ ...prev, title: e.target.value }))}
                    placeholder="请输入话题标题"
                    className="input input-bordered w-full focus:input-primary"
                    maxLength={100}
                  />
                  <label className="label">
                    <span className="label-text-alt text-base-content/40">
                      {createForm.title.length}/100
                    </span>
                  </label>
                </div>
                
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">描述</span>
                  </label>
                  <textarea
                    value={createForm.description}
                    onChange={(e) => setCreateForm(prev => ({ ...prev, description: e.target.value }))}
                    placeholder="描述这个话题的内容和讨论方向"
                    rows={3}
                    className="textarea textarea-bordered w-full focus:textarea-primary"
                    maxLength={500}
                  />
                  <label className="label">
                    <span className="label-text-alt text-base-content/40">
                      {createForm.description.length}/500
                    </span>
                  </label>
                </div>
                
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">封面图 URL</span>
                  </label>
                  <input
                    type="url"
                    value={createForm.cover}
                    onChange={(e) => setCreateForm(prev => ({ ...prev, cover: e.target.value }))}
                    placeholder="https://example.com/cover.jpg"
                    className="input input-bordered w-full focus:input-primary"
                  />
                  {createForm.cover && (
                    <div className="mt-2">
                      <img 
                        src={createForm.cover} 
                        alt="封面预览"
                        className="w-full h-32 object-cover rounded-lg"
                        onError={(e) => {
                          (e.target as HTMLImageElement).style.display = 'none';
                        }}
                      />
                    </div>
                  )}
                </div>
                
                <div className="form-control">
                  <label className="cursor-pointer label justify-start gap-3">
                    <input
                      type="checkbox"
                      checked={createForm.is_public}
                      onChange={(e) => setCreateForm(prev => ({ ...prev, is_public: e.target.checked }))}
                      className="checkbox checkbox-primary"
                    />
                    <span className="label-text">公开话题</span>
                  </label>
                  <p className="text-sm text-base-content/40 mt-1 ml-7">
                    公开话题所有人可见，私密话题仅关注者可查看
                  </p>
                </div>
              </div>
              
              {/* 按钮区域 */}
              <div className="modal-action mt-6">
                <button
                  onClick={() => setShowCreateModal(false)}
                  className="btn btn-ghost"
                >
                  取消
                </button>
                <button
                  onClick={handleCreateTopic}
                  disabled={creating || !createForm.title.trim()}
                  className="btn btn-primary"
                >
                  {creating && <span className="loading loading-spinner loading-sm"></span>}
                  {creating ? '创建中...' : '创建'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}