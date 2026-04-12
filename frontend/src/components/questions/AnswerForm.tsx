'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { questionsApi } from '@/lib/api/questions';
import { tagApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import { toast } from 'react-hot-toast';
import type { Tag } from '@/types';

interface AskForm {
  title: string;
  content: string;
  summary: string;
  reward_score: number;
  tag_ids: number[];
}

export default function AskQuestionPage() {
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [content, setContent] = useState('');
  const [tags, setTags] = useState<Tag[]>([]);
  const [selectedTags, setSelectedTags] = useState<number[]>([]);

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<AskForm>({
    defaultValues: {
      title: '',
      summary: '',
      reward_score: 0,
      tag_ids: [],
    },
  });

  useEffect(() => {
    loadTags();
  }, []);

  const loadTags = async () => {
    try {
      const response = await tagApi.list();
      if (response.data.code === 200) {
        setTags(response.data.data);
      }
    } catch (error) {
      console.error('Failed to load tags:', error);
    }
  };

  // 检查登录状态
  if (!isAuthenticated && typeof window !== 'undefined') {
    router.push('/login?redirect=/questions/ask');
    return null;
  }

  const onSubmit = async (data: AskForm) => {
    if (!content.trim()) {
      toast.error('请输入问题内容');
      return;
    }

    if (data.title.length < 5) {
      toast.error('标题至少需要5个字符');
      return;
    }

    setLoading(true);
    try {
      const response = await questionsApi.create({
        title: data.title,
        content: content,
        summary: data.summary,
        reward_score: data.reward_score,
        tag_ids: selectedTags,
      });

      if (response.data.code === 200) {
        toast.success('问题发布成功！');
        router.push(`/questions/${response.data.data.id}`);
      } else {
        toast.error(response.data.message || '发布失败');
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '发布失败');
    } finally {
      setLoading(false);
    }
  };

  const toggleTag = (tagId: number) => {
    setSelectedTags(prev =>
      prev.includes(tagId)
        ? prev.filter(id => id !== tagId)
        : [...prev, tagId]
    );
    setValue('tag_ids', selectedTags);
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <div className="bg-white rounded-lg shadow-sm">
          <div className="p-6 border-b">
            <h1 className="text-2xl font-bold text-gray-900">提问</h1>
            <p className="text-gray-500 mt-1">详细描述你的问题，获得更精准的回答</p>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="p-6 space-y-6">
            {/* 标题 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                标题 <span className="text-red-500">*</span>
              </label>
              <input
                {...register('title', { 
                  required: '请输入标题', 
                  minLength: { value: 5, message: '标题至少5个字符' },
                  maxLength: { value: 200, message: '标题最多200个字符' }
                })}
                type="text"
                placeholder="用简洁的语言描述你的问题"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              />
              {errors.title && (
                <p className="mt-1 text-sm text-red-500">{errors.title.message}</p>
              )}
            </div>

            {/* 内容 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                问题描述 <span className="text-red-500">*</span>
              </label>
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                rows={10}
                placeholder="详细描述你的问题，包括背景、尝试过的解决方案等..."
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              />
              {!content && (
                <p className="mt-1 text-sm text-red-500">请输入问题内容</p>
              )}
            </div>

            {/* 摘要 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                问题摘要
              </label>
              <textarea
                {...register('summary')}
                rows={2}
                placeholder="简要描述问题（可选，将显示在列表中）"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              />
            </div>

            {/* 标签 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                标签
              </label>
              <div className="flex flex-wrap gap-2">
                {tags.map((tag) => (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => toggleTag(tag.id)}
                    className={`px-3 py-1 rounded-full text-sm transition-colors ${
                      selectedTags.includes(tag.id)
                        ? 'bg-indigo-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    {tag.name}
                  </button>
                ))}
              </div>
            </div>

            {/* 悬赏积分 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                悬赏积分
              </label>
              <div className="flex items-center gap-4">
                <input
                  {...register('reward_score', { min: 0, max: 100 })}
                  type="number"
                  min="0"
                  max="100"
                  className="w-32 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                />
                <span className="text-gray-500 text-sm">
                  当前积分: {user?.score || 0}
                </span>
              </div>
              <p className="mt-1 text-sm text-gray-500">
                悬赏积分可以吸引更多回答，回答被采纳后积分将转给回答者
              </p>
            </div>

            {/* 提交按钮 */}
            <div className="flex justify-end gap-3 pt-4 border-t">
              <button
                type="button"
                onClick={() => router.back()}
                className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
              >
                取消
              </button>
              <button
                type="submit"
                disabled={loading}
                className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? '发布中...' : '发布问题'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}