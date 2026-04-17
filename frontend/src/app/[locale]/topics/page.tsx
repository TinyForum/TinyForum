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
  FireIcon,
} from '@heroicons/react/24/outline';
import type { Topic } from '@/lib/api/types';
import { TopicCard } from '@/components/topic/TopicCard';

// 加载骨架屏组件
function LoadingSkeleton() {
  return (
    <div className="space-y-4">
      {[1, 2, 3].map((i) => (
        <div key={i} className="card bg-base-100 shadow-sm p-6 animate-pulse">
          <div className="flex gap-4">
            <div className="skeleton h-16 w-16 rounded-xl" />
            <div className="flex-1 space-y-3">
              <div className="skeleton h-6 w-1/3" />
              <div className="skeleton h-4 w-2/3" />
              <div className="flex gap-4">
                <div className="skeleton h-4 w-20" />
                <div className="skeleton h-4 w-20" />
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

// 空状态组件
function EmptyState({ isAuthenticated, onOpenModal }: { isAuthenticated: boolean; onOpenModal: () => void }) {
  return (
    <div className="card bg-base-100 shadow-sm p-12 text-center border border-base-200">
      <div className="text-6xl mb-4 opacity-50">📚</div>
      <h3 className="text-xl font-semibold text-base-content mb-2">暂无话题</h3>
      <p className="text-base-content/60 mb-6">
        成为第一个创建话题的人，开启精彩讨论！
      </p>
      {isAuthenticated && (
        <button onClick={onOpenModal} className="btn btn-primary gap-2">
          <PlusIcon className="w-4 h-4" />
          创建第一个话题
        </button>
      )}
    </div>
  );
}

// 统计卡片组件
function StatsCards({ totalTopics, totalFollowers, hotTopicsCount }: { 
  totalTopics: number; 
  totalFollowers: number;
  hotTopicsCount: number;
}) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
      <div className="bg-gradient-to-br from-primary/10 to-primary/5 rounded-2xl p-4 text-center border border-primary/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <DocumentTextIcon className="w-5 h-5 text-primary" />
          <span className="text-sm font-medium text-primary">话题总数</span>
        </div>
        <div className="text-2xl font-bold text-base-content">{totalTopics}</div>
      </div>
      
      <div className="bg-gradient-to-br from-secondary/10 to-secondary/5 rounded-2xl p-4 text-center border border-secondary/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <UserGroupIcon className="w-5 h-5 text-secondary" />
          <span className="text-sm font-medium text-secondary">总关注数</span>
        </div>
        <div className="text-2xl font-bold text-base-content">{totalFollowers}</div>
      </div>
      
      <div className="bg-gradient-to-br from-warning/10 to-warning/5 rounded-2xl p-4 text-center border border-warning/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <FireIcon className="w-5 h-5 text-warning" />
          <span className="text-sm font-medium text-warning">热门话题</span>
        </div>
        <div className="text-2xl font-bold text-base-content">{hotTopicsCount}</div>
      </div>
    </div>
  );
}

// 分页组件
function Pagination({ currentPage, totalPages, onPageChange }: { 
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
        for (let i = totalPages - maxVisible + 1; i <= totalPages; i++) pages.push(i);
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
              currentPage === pageNum ? 'btn-primary' : 'btn-ghost'
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

// 创建话题模态框组件
function CreateTopicModal({ 
  isOpen, 
  onClose, 
  onCreate 
}: { 
  isOpen: boolean; 
  onClose: () => void; 
  onCreate: (data: { title: string; description: string; cover: string; is_public: boolean }) => Promise<void>;
}) {
  const [form, setForm] = useState({
    title: '',
    description: '',
    cover: '',
    is_public: true,
  });
  const [creating, setCreating] = useState(false);

  const handleSubmit = async () => {
    if (!form.title.trim()) {
      toast.error('请输入话题标题');
      return;
    }

    setCreating(true);
    try {
      await onCreate(form);
      setForm({ title: '', description: '', cover: '', is_public: true });
      onClose();
    } finally {
      setCreating(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 animate-fade-in">
      <div className="card bg-base-100 shadow-xl w-full max-w-md transform transition-all duration-300 scale-100">
        <div className="card-body p-6">
          {/* 模态框头部 */}
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-2">
              <div className="p-2 bg-primary/10 rounded-xl">
                <PlusIcon className="w-5 h-5 text-primary" />
              </div>
              <h2 className="card-title text-xl font-bold text-base-content">
                创建话题
              </h2>
            </div>
            <button
              onClick={onClose}
              className="btn btn-ghost btn-sm btn-circle"
            >
              <XMarkIcon className="w-5 h-5" />
            </button>
          </div>
          
          <div className="divider my-2" />
          
          {/* 表单内容 */}
          <div className="space-y-5">
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  标题 <span className="text-error">*</span>
                </span>
              </label>
              <input
                type="text"
                value={form.title}
                onChange={(e) => setForm(prev => ({ ...prev, title: e.target.value }))}
                placeholder="例如：科技前沿、美食分享、旅行日记"
                className="input input-bordered w-full focus:input-primary"
                maxLength={100}
                autoFocus
              />
              <label className="label">
                <span className="label-text-alt text-base-content/40">
                  {form.title.length}/100
                </span>
              </label>
            </div>
            
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">描述</span>
              </label>
              <textarea
                value={form.description}
                onChange={(e) => setForm(prev => ({ ...prev, description: e.target.value }))}
                placeholder="描述这个话题的讨论方向和内容范围..."
                rows={3}
                className="textarea textarea-bordered w-full focus:textarea-primary resize-none"
                maxLength={500}
              />
              <label className="label">
                <span className="label-text-alt text-base-content/40">
                  {form.description.length}/500
                </span>
              </label>
            </div>
            
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">封面图 URL</span>
              </label>
              <input
                type="url"
                value={form.cover}
                onChange={(e) => setForm(prev => ({ ...prev, cover: e.target.value }))}
                placeholder="https://example.com/cover.jpg"
                className="input input-bordered w-full focus:input-primary"
              />
              {form.cover && (
                <div className="mt-3 rounded-xl overflow-hidden border border-base-200">
                  <img 
                    src={form.cover} 
                    alt="封面预览"
                    className="w-full h-32 object-cover"
                    onError={(e) => {
                      (e.target as HTMLImageElement).style.display = 'none';
                    }}
                  />
                </div>
              )}
            </div>
            
            <div className="form-control bg-base-200/50 rounded-xl p-4">
              <label className="cursor-pointer label justify-start gap-3">
                <input
                  type="checkbox"
                  checked={form.is_public}
                  onChange={(e) => setForm(prev => ({ ...prev, is_public: e.target.checked }))}
                  className="checkbox checkbox-primary"
                />
                <div>
                  <span className="label-text font-medium">公开话题</span>
                  <p className="text-xs text-base-content/40 mt-0.5">
                    公开话题所有人可见，私密话题仅关注者可查看
                  </p>
                </div>
              </label>
            </div>
          </div>
          
          {/* 按钮区域 */}
          <div className="modal-action mt-6 gap-3">
            <button
              onClick={onClose}
              className="btn btn-ghost flex-1"
            >
              取消
            </button>
            <button
              onClick={handleSubmit}
              disabled={creating || !form.title.trim()}
              className="btn btn-primary flex-1 gap-2"
            >
              {creating && <span className="loading loading-spinner loading-sm"></span>}
              {creating ? '创建中...' : '创建话题'}
            </button>
          </div>
        </div>
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
  const handleCreateTopic = async (data: { title: string; description: string; cover: string; is_public: boolean }) => {
    try {
      const response = await topicApi.create({
        title: data.title.trim(),
        description: data.description.trim() || undefined,
        cover: data.cover || undefined,
        is_public: data.is_public,
      });
      
      if (response.data.code === 200 || response.data.code === 201) {
        toast.success('话题创建成功');
        await loadTopics();
        return Promise.resolve();
      } else {
        toast.error(response.data.message || '创建失败');
        return Promise.reject();
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '创建失败');
      return Promise.reject();
    }
  };

  useEffect(() => {
    loadTopics();
  }, [page]);

  const totalPages = Math.ceil(total / pageSize);
  
  // 计算统计数据 - 使用正确的属性名 follower_count
  const totalFollowers = topics.reduce((sum, topic) => sum + (topic.follower_count || 0), 0);
  const hotTopicsCount = topics.filter(t => (t.follower_count || 0) > 100).length;

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-4xl mx-auto px-4 py-8 md:py-12">
        {/* 头部区域 */}
        <div className="text-center mb-8">
          <div className="relative inline-block mb-4">
            <div className="absolute inset-0 bg-gradient-to-r from-primary/20 to-secondary/20 rounded-full blur-2xl" />
            <div className="relative bg-gradient-to-br from-primary to-secondary p-3 rounded-full shadow-lg">
              <FireIcon className="w-8 h-8 text-white" />
            </div>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
            话题广场
          </h1>
          <p className="text-base-content/60 mt-2">
            发现和参与感兴趣的话题讨论
          </p>
        </div>

        {/* 统计卡片 */}
        {!loading && topics.length > 0 && (
          <StatsCards 
            totalTopics={total}
            totalFollowers={totalFollowers}
            hotTopicsCount={hotTopicsCount}
          />
        )}

        {/* 操作栏 */}
        <div className="flex justify-end mb-6">
          {isAuthenticated && (
            <button
              onClick={() => setShowCreateModal(true)}
              className="btn btn-primary gap-2 shadow-md hover:shadow-lg transition-all"
            >
              <PlusIcon className="w-4 h-4" />
              创建话题
            </button>
          )}
        </div>

        {/* 话题列表 */}
        {loading ? (
          <LoadingSkeleton />
        ) : topics.length === 0 ? (
          <EmptyState isAuthenticated={isAuthenticated} onOpenModal={() => setShowCreateModal(true)} />
        ) : (
          <>
            <div className="space-y-4">
              {topics.map((topic, index) => (
                <div key={topic.id} className="animate-fade-in" style={{ animationDelay: `${index * 50}ms` }}>
                  <TopicCard topic={topic} onFollowChange={loadTopics} />
                </div>
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
              共 {total} 个话题，找到你感兴趣的
            </div>
          </>
        )}
      </div>

      {/* 创建话题模态框 */}
      <CreateTopicModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onCreate={handleCreateTopic}
      />
    </div>
  );
}