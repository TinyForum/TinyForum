'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import {
  MagnifyingGlassIcon,
  ChatBubbleLeftRightIcon,
  CheckBadgeIcon,
  FireIcon,
  ArrowsRightLeftIcon,
  PlusIcon,
  UserCircleIcon,
  EyeIcon,
  TagIcon,
} from '@heroicons/react/24/outline';
import { questionsApi } from '@/lib/api/questions';
import { useAuthStore } from '@/store/auth';
import type { Post } from '@/types';

type FilterType = 'all' | 'unanswered' | 'answered';
type SortType = 'latest' | 'hot' | 'score';

export default function QuestionsPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { isAuthenticated } = useAuthStore();

  const [questions, setQuestions] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [filter, setFilter] = useState<FilterType>((searchParams.get('filter') as FilterType) || 'all');
  const [sort, setSort] = useState<SortType>('latest');
  const [keyword, setKeyword] = useState(searchParams.get('keyword') || '');

  const pageSize = 15;

  useEffect(() => {
    loadQuestions();
  }, [page, filter, sort, keyword]);

  const loadQuestions = async () => {
    setLoading(true);
    try {
      const response = await questionsApi.list({
        page,
        page_size: pageSize,
        filter: filter === 'all' ? undefined : filter,
        keyword: keyword || undefined,
      });
      if (response.data.code === 200) {
        setQuestions(response.data.data.list);
        setTotal(response.data.data.total);
      }
    } catch (error) {
      console.error('Failed to load questions:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    loadQuestions();
  };

  const formatTime = (time: string) => {
    const date = new Date(time);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    if (days === 0) return '今天';
    if (days === 1) return '昨天';
    if (days < 7) return `${days}天前`;
    return date.toLocaleDateString();
  };

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-5xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex justify-between items-center mb-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">问答社区</h1>
            <p className="text-gray-500 mt-1">提问、回答、分享知识</p>
          </div>
          {isAuthenticated && (
            <Link
              href="/questions/ask"
              className="flex items-center gap-2 px-5 py-2.5 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors font-medium shadow-sm"
            >
              <PlusIcon className="w-5 h-5" />
              提问
            </Link>
          )}
        </div>

        {/* Search Bar */}
        <form onSubmit={handleSearch} className="mb-6">
          <div className="relative">
            <input
              type="text"
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              placeholder="搜索问题..."
              className="w-full px-4 py-3 pl-11 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
            />
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <button
              type="submit"
              className="absolute right-2 top-1/2 -translate-y-1/2 px-4 py-1.5 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors"
            >
              搜索
            </button>
          </div>
        </form>

        {/* Filters */}
        <div className="bg-white rounded-lg shadow-sm mb-6">
          <div className="flex items-center justify-between p-4 border-b">
            <div className="flex gap-2">
              {[
                { value: 'all', label: '全部', icon: ChatBubbleLeftRightIcon },
                { value: 'unanswered', label: '未回答', icon: FireIcon },
                { value: 'answered', label: '已解决', icon: CheckBadgeIcon },
              ].map((item) => {
                const Icon = item.icon;
                const isActive = filter === item.value;
                return (
                  <button
                    key={item.value}
                    onClick={() => {
                      setFilter(item.value as FilterType);
                      setPage(1);
                    }}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg transition-colors ${
                      isActive
                        ? 'bg-indigo-50 text-indigo-600 border border-indigo-200'
                        : 'text-gray-600 hover:bg-gray-50'
                    }`}
                  >
                    <Icon className="w-4 h-4" />
                    {item.label}
                  </button>
                );
              })}
            </div>
            <div className="flex items-center gap-2 text-sm">
              <ArrowsRightLeftIcon className="w-4 h-4 text-gray-400" />
              <select
                value={sort}
                onChange={(e) => setSort(e.target.value as SortType)}
                className="border-none bg-transparent focus:outline-none text-gray-600 cursor-pointer"
              >
                <option value="latest">最新</option>
                <option value="hot">最热</option>
                <option value="score">悬赏最高</option>
              </select>
            </div>
          </div>

          {/* Questions List */}
          <div className="divide-y">
            {loading ? (
              <div className="p-8 text-center text-gray-500">加载中...</div>
            ) : questions.length === 0 ? (
              <div className="p-8 text-center text-gray-500">
                <ChatBubbleLeftRightIcon className="w-12 h-12 mx-auto text-gray-300 mb-3" />
                <p>暂无问题</p>
                {isAuthenticated && (
                  <Link href="/questions/ask" className="text-indigo-600 hover:underline mt-2 inline-block">
                    成为第一个提问者
                  </Link>
                )}
              </div>
            ) : (
              questions.map((question) => (
                <Link
                  key={question.id}
                  href={`/questions/${question.id}`}
                  className="block p-5 hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="text-lg font-medium text-gray-900 hover:text-indigo-600 mb-2">
                        {question.title}
                      </h3>
                      <p className="text-gray-500 text-sm line-clamp-2 mb-3">
                        {question.summary || question.content.replace(/<[^>]*>/g, '').slice(0, 150)}
                      </p>
                      <div className="flex items-center gap-4 text-sm text-gray-400">
                        <span className="flex items-center gap-1">
                          <UserCircleIcon className="w-4 h-4" />
                          {question.author?.username}
                        </span>
                        <span>{formatTime(question.created_at)}</span>
                        <span className="flex items-center gap-1">
                          <ChatBubbleLeftRightIcon className="w-4 h-4" />
                          {(question as any).answer_count || 0} 回答
                        </span>
                        <span className="flex items-center gap-1">
                          <EyeIcon className="w-4 h-4" />
                          {question.view_count} 浏览
                        </span>
                        {(question as any).reward_score > 0 && (
                          <span className="text-orange-500 flex items-center gap-1">
                            💰 {(question as any).reward_score} 积分
                          </span>
                        )}
                      </div>
                      {/* Tags */}
                      {question.tags && question.tags.length > 0 && (
                        <div className="flex flex-wrap gap-2 mt-3">
                          {question.tags.slice(0, 3).map((tag) => (
                            <span
                              key={tag.id}
                              className="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded-md"
                              style={{ borderLeft: `2px solid ${tag.color || '#6366f1'}` }}
                            >
                              <TagIcon className="w-3 h-3 inline mr-0.5" />
                              {tag.name}
                            </span>
                          ))}
                        </div>
                      )}
                    </div>
                    {(question as any).accepted_answer_id && (
                      <div className="ml-4 text-center">
                        <div className="px-3 py-1 bg-green-50 rounded-lg">
                          <CheckBadgeIcon className="w-5 h-5 text-green-500 mx-auto" />
                          <div className="text-xs text-green-600 mt-1">已解决</div>
                        </div>
                      </div>
                    )}
                  </div>
                </Link>
              ))
            )}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-center gap-2 p-4 border-t">
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
        </div>
      </div>
    </div>
  );
}